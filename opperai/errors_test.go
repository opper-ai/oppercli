package opperai

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		operation   func(*Client) error
		statusCode  int
		wantErrType error
	}{
		{
			name: "invalid api key",
			operation: func(c *Client) error {
				c.APIKey = "invalid"
				_, err := c.Models.List(context.Background())
				return err
			},
			statusCode:  http.StatusUnauthorized,
			wantErrType: ErrUnauthorized,
		},
		{
			name: "rate limit",
			operation: func(c *Client) error {
				_, err := c.Models.List(context.Background())
				return err
			},
			statusCode:  http.StatusTooManyRequests,
			wantErrType: ErrRateLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusUnauthorized {
					json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
				}
				if tt.statusCode == http.StatusTooManyRequests {
					json.NewEncoder(w).Encode(map[string]string{"error": "rate limit exceeded"})
				}
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			err := tt.operation(client)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, tt.wantErrType) {
				t.Errorf("got error %v, want error type %v", err, tt.wantErrType)
			}
		})
	}
}
