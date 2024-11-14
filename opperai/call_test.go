package opperai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCallClient_Call(t *testing.T) {
	tests := []struct {
		name        string
		response    interface{}
		statusCode  int
		stream      bool
		expectError bool
	}{
		{
			name: "successful call",
			response: map[string]string{
				"message": "Hello, world!",
			},
			statusCode:  http.StatusOK,
			stream:      false,
			expectError: false,
		},
		{
			name:        "server error",
			response:    nil,
			statusCode:  http.StatusInternalServerError,
			stream:      false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if r.URL.Path != "/v1/call" {
					t.Errorf("expected path /v1/call, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				if tt.response != nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			ctx := context.Background()

			result, err := client.Call.Call(ctx, "test-name", "test-instructions", "test-input", "", false)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			expected := "Hello, world!"
			if result.Message != expected {
				t.Errorf("expected '%s', got '%s'", expected, result.Message)
			}
		})
	}
}
