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
    "path"
    "strings"
    "time"

    "github.com/AbsaOSS/env-binder/env"
    "github.com/google/uuid"
    "github.com/gorilla/mux"

    "github.com/rakhbari/gomux1/config"
    utils "github.com/rakhbari/gomux1/utils"
)

type PingPayload struct {
    Response string `json:"response"`
}

type HealthPayload struct {
    Healthy bool `json:"healthy"`
}

type StandardApiResponse struct {
    RequestId string  `json:"requestId"`
    Timestamp string  `json:"timestamp"`
    ExecHost  string  `json:"execHost"`
    Payload   any     `json:"payload"`
    Errors    []Error `json:"errors"`
}

type Error struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail"`
    HelpUrl string `json:"helpUrl"`
}

func HttpResponseWriter(w http.ResponseWriter, status int, apiResp *StandardApiResponse) {
    apiResp.RequestId = uuid.New().String()
    apiResp.Timestamp = time.Now().String()
    apiResp.ExecHost = readExecHost()
    resp, err := json.Marshal(apiResp)
    if err != nil {
        log.Printf("!!!> ERROR: json.Marshall failed: %v", err)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    io.WriteString(w, string(resp))
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
    // Just respond with a "pong!"
    apiResponse := &StandardApiResponse{Payload: PingPayload{Response: "pong!"}}
    HttpResponseWriter(w, http.StatusOK, apiResponse)
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    // Do whatever needed to run health checks
    // In the future we could report back on the status of our DB, or our cache
    // (e.g. Redis) by performing a simple PING, and include them in the response.
    apiResponse := &StandardApiResponse{Payload: HealthPayload{Healthy: true}}
    HttpResponseWriter(w, http.StatusOK, apiResponse)
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
    responseStatus := http.StatusOK
    // If the Version struct hasn't been loaded for some reason set responseStatus to NotFound
    if version == (utils.Version{}) {
        responseStatus = http.StatusNotFound
    }
    // Responds with the value of the utils.Version struct loaded at app startup
    HttpResponseWriter(w, responseStatus, &StandardApiResponse{Payload: &version})
}

func BearerTokenFormHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("scheme: %s", r.URL.Scheme)
    log.Printf("path: %s", r.URL.Path)
    log.Printf("url_long: %s", r.Form["url_long"])

    // NOTE: If you do not call ParseForm method, the following data can not be obtained
    r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
    namespace := r.FormValue("namespace")
    svcAcct := r.FormValue("service_acct")
    argoBaseUrl := r.FormValue("argo_base_url")

    home, _ := os.UserHomeDir()

    bearerToken, err := utils.GetSvcAcctToken(path.Join(home, ".kube/config"), namespace, svcAcct)
    if err != nil {
        error := &Error{Code: "E0001", Message: err.Error()}
        HttpResponseWriter(w, http.StatusInternalServerError, &StandardApiResponse{Errors: []Error{*error}})
        return
    }

    // https://argo.akhbari.us:9443/workflows/app1?limit=50
    argoUrl := argoBaseUrl + "/workflows/" + namespace + "?limit=50"
    r.Header.Add("Authorization", "Bearer "+*bearerToken)
    http.Redirect(w, r, argoUrl, http.StatusSeeOther)
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

func ConfigureAppRouter() *mux.Router {
    router := mux.NewRouter()
    // Add routes
    router.HandleFunc("/v1/ping", PingHandler).Methods("GET")
    router.HandleFunc("/health", HealthCheckHandler).Methods("GET")
    router.HandleFunc("/version", VersionHandler).Methods("GET")
    router.HandleFunc("/v1/bearer-token", BearerTokenFormHandler).Methods("POST")
    return router
}

func configureAppServer(addr string, router *mux.Router, cfg *config.Config) *http.Server {
    return &http.Server{
        Addr: addr,
        // Good practice to set timeouts to avoid Slowloris attacks.
        WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
        ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
        IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
        Handler:      router, // Pass in our instance of gorilla/mux.Router
    }
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

    router := ConfigureAppRouter()

    ServeStatic(router, cfg.WebApp.ContentDir)

    httpSrv := configureAppServer(httpAddr, router, cfg)
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
        httpsSrv = configureAppServer(httpsAddr, router, cfg)
        // Run our TLS server in a goroutine so that it doesn't block.
        go func() {
            tlsCertFile := utils.GetTlsCertFile(cfg)
            if tlsCertFile == nil {
                log.Println("!!!> ERROR: Problem encountered while loading TLS cert files. Not starting TLS server!")
                return
            }
            log.Printf("===> Starting TLS server ... (tlsCertFile: %s)\n", *tlsCertFile)
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

    // Block until we receive our os.Signal. (Ctrl-C)
    <-c

    // Create a deadline to wait for.
    ctx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()

    // Doesn't block if no connections, but will otherwise wait
    // until the timeout deadline.
    log.Println("===> Shutting down HTTP server ...")
    httpSrv.Shutdown(ctx)
    if httpsSrv != nil {
        log.Println("===> Shutting down TLS server ...")
        httpsSrv.Shutdown(ctx)
    }

    // Optionally, you could run srv.Shutdown in a goroutine and block on
    // <-ctx.Done() if your application should wait for other services
    // to finalize based on context cancellation.
    os.Exit(0)
}

func ServeStatic(router *mux.Router, staticDirectory string) {
    staticPaths := map[string]string{
        "/app/":     staticDirectory + "/",
        "/styles/":  staticDirectory + "/styles/",
        "/images/":  staticDirectory + "/images/",
        "/scripts/": staticDirectory + "/scripts/",
    }
    for pathName, pathValue := range staticPaths {
        router.PathPrefix(pathName).Handler(http.StripPrefix(pathName, http.FileServer(http.Dir(pathValue))))
    }
}
