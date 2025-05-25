package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	utils "github.com/rakhbari/gomux1/utils"
)

var router *mux.Router

type ExpectedHttpResponse struct {
	RequestId string  `json:"requestId"`
	Timestamp string  `json:"timestamp"`
	ExecHost  string  `json:"execHost"`
	Payload   any     `json:"payload"`
	Errors    []Error `json:"errors"`
}

func init() {
	log.Println("init ...")
	router = ConfigureAppRouter()
}

func TestPingHandler(t *testing.T) {
	t.Log("Testing PingHandler ...")
	req, err := http.NewRequest("GET", "/v1/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, was looking for %v",
			status, http.StatusOK)
	}

	resp := ExpectedHttpResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("Got error while trying to unmarshal the response. Error: %+v", err)
	}

	// Check the response structure
	if resp.RequestId == "" {
		t.Error("RequestId is empty")
	}
	if resp.Timestamp == "" {
		t.Error("Timestamp is empty")
	}
	if resp.ExecHost == "" {
		t.Error("ExecHost is empty")
	}

	// Check the payload
	payload, ok := resp.Payload.(map[string]interface{})
	if !ok {
		t.Error("Payload is not a map")
	}
	if response, ok := payload["response"].(string); !ok || response != "pong!" {
		t.Errorf("Expected 'pong!' response, got %v", payload["response"])
	}
}

func TestHealthCheckHandler(t *testing.T) {
	t.Log("Testing HealthCheckHandler ...")
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, was looking for %v",
			status, http.StatusOK)
	}

	resp := ExpectedHttpResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("Got error while trying to unmarshal the response. Error: %+v", err)
	}

	// Check the response structure
	if resp.RequestId == "" {
		t.Error("RequestId is empty")
	}
	if resp.Timestamp == "" {
		t.Error("Timestamp is empty")
	}
	if resp.ExecHost == "" {
		t.Error("ExecHost is empty")
	}

	// Check the payload
	payload, ok := resp.Payload.(map[string]interface{})
	if !ok {
		t.Error("Payload is not a map")
	}
	if healthy, ok := payload["healthy"].(bool); !ok || !healthy {
		t.Errorf("Expected healthy=true, got %v", payload["healthy"])
	}
}

func TestVersionHandler(t *testing.T) {
	tests := []struct {
		name           string
		version        utils.Version
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Success - Version Available",
			version: utils.Version{
				Timestamp: "2024-01-01 12:00:00 UTC",
				GitSha:    "abc123",
				GitBranch: "main",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Error - Version Not Available",
			version:        utils.Version{},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the global version variable for the test
			version = tt.version

			req, err := http.NewRequest("GET", "/version", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			resp := ExpectedHttpResponse{}
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			if err != nil {
				t.Fatal(err)
			}

			// Check the response structure
			if resp.RequestId == "" {
				t.Error("RequestId is empty")
			}
			if resp.Timestamp == "" {
				t.Error("Timestamp is empty")
			}
			if resp.ExecHost == "" {
				t.Error("ExecHost is empty")
			}

			if !tt.expectError {
				// Check if we got the expected version in the payload
				versionPayload, ok := resp.Payload.(map[string]interface{})
				if !ok {
					t.Fatal("payload is not a map")
				}
				if versionPayload["timestamp"] != tt.version.Timestamp {
					t.Errorf("handler returned wrong timestamp: got %v want %v",
						versionPayload["timestamp"], tt.version.Timestamp)
				}
				if versionPayload["gitSha"] != tt.version.GitSha {
					t.Errorf("handler returned wrong gitSha: got %v want %v",
						versionPayload["gitSha"], tt.version.GitSha)
				}
				if versionPayload["gitBranch"] != tt.version.GitBranch {
					t.Errorf("handler returned wrong gitBranch: got %v want %v",
						versionPayload["gitBranch"], tt.version.GitBranch)
				}
			}
		})
	}
}
