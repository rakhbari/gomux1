package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var router *mux.Router

type ExpectedHttpResponse struct {
	RequestId string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	ExecHost  string `json:"execHost"`
	Payload   any    `json:"payload"`
}

func init() {
	log.Println("init ...")
	router = ConfigureAppRouter()
}

func TestPingHandler(t *testing.T) {
	t.Log("Testing PingHandler ...")
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/v1/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, was looking for %v",
			status, http.StatusOK)
	}

	resp := ExpectedHttpResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("Got error while trying to unmarshal the response. Error: %+v", err)
	}
	t.Logf("resp: %+v", resp)
}

func TestHealthCheckHandler(t *testing.T) {
	t.Log("Testing HealthCheckHandler ...")
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, was looking for %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	/*
	   expected := `{"alive": true}`
	   if rr.Body.String() != expected {
	       t.Errorf("handler returned unexpected body: got %v want %v",
	           rr.Body.String(), expected)
	   }
	*/
}
