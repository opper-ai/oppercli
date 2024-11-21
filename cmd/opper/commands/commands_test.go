package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opper-ai/oppercli/opperai"
)

func setupTestServer(handler http.HandlerFunc) (*httptest.Server, *opperai.Client) {
	server := httptest.NewServer(handler)
	client := opperai.NewClient("test-key", server.URL)
	return server, client
}

func TestExecuteListIndexes(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		response   []opperai.Index
		statusCode int
		wantOutput string
		wantErr    bool
	}{
		{
			name: "successful list",
			args: []string{"indexes", "list"},
			response: []opperai.Index{
				{Name: "index1"},
				{Name: "index2"},
			},
			statusCode: http.StatusOK,
			wantOutput: "index1\nindex2\n",
			wantErr:    false,
		},
		{
			name:       "server error",
			args:       []string{"indexes", "list"},
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.response)
				}
			})
			defer server.Close()

			var buf bytes.Buffer
			err := ExecuteListIndexes(context.Background(), client, &buf, tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteListIndexes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && buf.String() != tt.wantOutput {
				t.Errorf("ExecuteListIndexes() output = %q, want %q", buf.String(), tt.wantOutput)
			}
		})
	}
}
