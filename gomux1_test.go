package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rakhbari/gomux1/utils"
)

func TestPingHandler(t *testing.T) {
	t.Log("Testing PingHandler ...")
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PingHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, was looking for %v",
			status, http.StatusOK)
	}
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
	handler := http.HandlerFunc(HealthCheckHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

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

			// Create a request to pass to our handler
			req, err := http.NewRequest("GET", "/version", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(VersionHandler)

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Parse the response
			var response StandardHttpResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}

			// Check if we got the expected version in the payload
			if !tt.expectError {
				versionPayload, ok := response.Payload.(map[string]interface{})
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

			// Check if we got the expected error
			if tt.expectError {
				if len(response.Errors) != 1 {
					t.Errorf("expected 1 error, got %d", len(response.Errors))
				}
				if response.Errors[0].Code != ErrCodeVersionNotFound {
					t.Errorf("handler returned wrong error code: got %v want %v",
						response.Errors[0].Code, ErrCodeVersionNotFound)
				}
			} else {
				if len(response.Errors) != 0 {
					t.Errorf("expected no errors, got %d", len(response.Errors))
				}
			}
		})
	}
}
