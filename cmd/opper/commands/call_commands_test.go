package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opper-ai/oppercli/opperai"
)

func TestCallCommand_Execute(t *testing.T) {
	tests := []struct {
		name         string
		command      *CallCommand
		expectError  bool
		mockResponse string
	}{
		{
			name: "successful call",
			command: &CallCommand{
				Name:         "test-name",
				Instructions: "test-instructions",
				Input:        "test-input",
				Stream:       false,
			},
			expectError:  false,
			mockResponse: "Hello, world!",
		},
		{
			name: "successful streaming call",
			command: &CallCommand{
				Name:         "test-name",
				Instructions: "test-instructions",
				Input:        "test-input",
				Stream:       true,
			},
			expectError:  false,
			mockResponse: "Hello, world!",
		},
		{
			name: "empty input",
			command: &CallCommand{
				Name:         "test-name",
				Instructions: "test-instructions",
				Input:        "",
				Stream:       false,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// First request (non-streaming) always returns JSON
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{
					"message": tt.mockResponse,
				})
			}))
			defer server.Close()

			client := opperai.NewClient("test-key", server.URL)
			err := tt.command.Execute(context.Background(), client)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
