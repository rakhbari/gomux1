package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/AbsaOSS/env-binder/env"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/rakhbari/gomux1/config"
	"github.com/rakhbari/gomux1/utils"
)

type HealthPayload struct {
	Healthy bool `json:"healthy"`
}

type PingPayload struct {
	Response string `json:"response"`
}

type StandardHttpResponse struct {
	RequestId string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	ExecHost  string `json:"execHost"`
	Payload   any    `json:"payload"`
}

func HttpResponseWriter(w http.ResponseWriter, status int, payload any) {
	shr := &StandardHttpResponse{
		RequestId: uuid.New().String(),
		Timestamp: time.Now().String(),
		ExecHost:  readExecHost(),
		Payload:   payload}
	resp, err := json.Marshal(shr)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	io.WriteString(w, string(resp))
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	// Just respond with a "pong!"
	HttpResponseWriter(w, http.StatusOK, &PingPayload{Response: "pong!"})
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Do whatever needed to run health checks
	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	HttpResponseWriter(w, http.StatusOK, &HealthPayload{Healthy: true})
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	responseStatus := http.StatusOK
	// If the Version struct hasn't been loaded for some reason set responseStatus to NotFound
	if version == (utils.Version{}) {
		responseStatus = http.StatusNotFound
	}
	// Responds with the value of the utils.Version struct loaded at app startup
	HttpResponseWriter(w, responseStatus, &version)
}

func readExecHost() string {
	execHost := os.Getenv("POD_NAME")
	if execHost == "" {
		hostname, err := os.Hostname()
		if err == nil {
			execHost = hostname
		} else {
			execHost = "N/A"
		}
	}
	return execHost
}

func configureNewServer(addr string, router *mux.Router, cfg *config.Config) *http.Server {
	srv := &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
		Handler:      router, // Pass in our instance of gorilla/mux.Router
	}
	return srv
}

func cleanup(tlsCertFile string) {
	if !strings.HasSuffix(tlsCertFile, "tlsCertBundle") {
		return
	}
	fmt.Printf("---> Removing file: %s ...\n", tlsCertFile)
	if err := os.Remove(tlsCertFile); err != nil {
		utils.ProcessError(err)
	}
}

var version utils.Version

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Read in the config.Config struct and bind it with env variables (if any passed-in)
	cfg := &config.Config{}
	if err := env.Bind(cfg); err != nil {
		utils.ProcessError(err)
	}
	log.Printf("===> App config: %+v\n", cfg)

	// Load the utils.Version struct from the version.json file (if found)
	utils.LoadVersion(&version)
	log.Printf("===> App version: %+v\n", version)

	httpAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HttpPort)
	httpsAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HttpsPort)

	router := mux.NewRouter()
	// Add routes
	router.HandleFunc("/v1/ping", PingHandler).Methods("GET")
	router.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	router.HandleFunc("/version", VersionHandler).Methods("GET")

	httpSrv := configureNewServer(httpAddr, router, cfg)
	// Run our HTTP server in a goroutine so that it doesn't block.
	go func() {
		log.Println("===> Starting HTTP server ...")
		if err := httpSrv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// If TlsCertPath is passed in, start a TLS server also
	var httpsSrv *http.Server
	if len(cfg.Server.TlsCertPath) > 0 {
		httpsSrv = configureNewServer(httpsAddr, router, cfg)
		// Run our TLS server in a goroutine so that it doesn't block.
		go func() {
			tlsCertFile := utils.GetTlsCertFile(cfg)
			if tlsCertFile == nil {
				log.Println("!!!> ERROR: Problem encountered while loading TLS cert files. Not starting TLS server!")
				return
			}
			log.Printf("===> Starting HTTPS server ... (tlsCertFile: %s)\n", *tlsCertFile)
			if err := httpsSrv.ListenAndServeTLS(*tlsCertFile, cfg.Server.TlsKeyPath); err != nil {
				if !strings.Contains(strings.ToLower(err.Error()), "server closed") {
					utils.ProcessError(err)
				}
			}
			cleanup(*tlsCertFile)
		}()
	}

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	log.Println("===> Shutting down")
	httpSrv.Shutdown(ctx)
	if len(cfg.Server.TlsCertPath) > 0 {
		httpsSrv.Shutdown(ctx)
	}

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	os.Exit(0)
}
