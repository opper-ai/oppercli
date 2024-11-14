package opperai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// TestServer represents a test HTTP server
type TestServer struct {
	*httptest.Server
	// Add any fields needed for tracking requests/responses
}

// NewTestServer creates a new test server with the given handler
func NewTestServer(handler http.HandlerFunc) *TestServer {
	ts := &TestServer{}
	ts.Server = httptest.NewServer(handler)
	return ts
}

// MockJSONResponse creates a handler that returns a JSON response
func MockJSONResponse(statusCode int, response interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if response != nil {
			json.NewEncoder(w).Encode(response)
		}
	}
}
