package opperai

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		apiKey   string
		baseURL  []string
		wantBase string
	}{
		{
			name:     "default base URL",
			apiKey:   "test-key",
			baseURL:  []string{},
			wantBase: "https://api.opper.ai",
		},
		{
			name:     "custom base URL",
			apiKey:   "test-key",
			baseURL:  []string{"https://custom.api.com"},
			wantBase: "https://custom.api.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.apiKey, tt.baseURL...)
			if client.BaseURL != tt.wantBase {
				t.Errorf("NewClient() BaseURL = %v, want %v", client.BaseURL, tt.wantBase)
			}
			if client.APIKey != tt.apiKey {
				t.Errorf("NewClient() APIKey = %v, want %v", client.APIKey, tt.apiKey)
			}
			if client.Indexes == nil {
				t.Error("NewClient() Indexes is nil")
			}
			if client.Models == nil {
				t.Error("NewClient() Models is nil")
			}
			if client.Functions == nil {
				t.Error("NewClient() Functions is nil")
			}
		})
	}
}

func TestDoRequest(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful request",
			method:     "GET",
			path:       "/test",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "server error",
			method:     "POST",
			path:       "/error",
			statusCode: http.StatusInternalServerError,
			wantErr:    false, // DoRequest doesn't check status code
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check headers
				if r.Header.Get("Content-Type") != "application/json" {
					t.Error("Content-Type header not set correctly")
				}
				if r.Header.Get("X-OPPER-API-KEY") != "test-key" {
					t.Error("API key header not set correctly")
				}

				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			resp, err := client.DoRequest(context.Background(), tt.method, tt.path, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("DoRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if resp != nil {
				if resp.StatusCode != tt.statusCode {
					t.Errorf("DoRequest() status = %v, want %v", resp.StatusCode, tt.statusCode)
				}
				resp.Body.Close()
			}
		})
	}
}

func TestChat(t *testing.T) {
	tests := []struct {
		name         string
		functionPath string
		payload      ChatPayload
		stream       bool
		responses    []string
		statusCode   int
		wantErr      bool
	}{
		{
			name:         "successful chat",
			functionPath: "test/function",
			payload: ChatPayload{
				Messages: []Message{{Role: "user", Content: "Hello"}},
			},
			stream:     true,
			responses:  []string{`data: {"content":"Hello!"}`},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "rate limit error",
			functionPath: "test/function",
			payload: ChatPayload{
				Messages: []Message{{Role: "user", Content: "Hello"}},
			},
			stream:     true,
			statusCode: http.StatusTooManyRequests,
			wantErr:    true,
		},
		{
			name:         "function error",
			functionPath: "test/function",
			payload: ChatPayload{
				Messages: []Message{{Role: "user", Content: "Hello"}},
			},
			stream:     true,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := "/v1/chat/" + tt.functionPath
				if tt.stream {
					expectedPath += "?stream=true"
				}
				if r.URL.Path != "/v1/chat/"+tt.functionPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					flusher, ok := w.(http.Flusher)
					if !ok {
						t.Fatal("expected http.Flusher")
					}
					w.Header().Set("Content-Type", "text/event-stream")
					for _, resp := range tt.responses {
						fmt.Fprintf(w, "%s\n", resp)
						flusher.Flush()
					}
				}
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			chunks, err := client.Chat(context.Background(), tt.functionPath, tt.payload, tt.stream)

			if (err != nil) != tt.wantErr {
				t.Errorf("Chat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				received := make([]string, 0)
				for chunk := range chunks {
					received = append(received, string(chunk))
				}

				if len(received) != len(tt.responses) {
					t.Errorf("expected %d responses, got %d", len(tt.responses), len(received))
				}
				for i, resp := range tt.responses {
					if received[i] != resp+"\n" {
						t.Errorf("expected response %q, got %q", resp, received[i])
					}
				}
			}
		})
	}
}
