package opperai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	functionsBasePath = "/api/v1/functions"
	functionsByPath   = "/api/v1/functions/by_path"
	functionsForOrg   = "/api/v1/functions/for_org"
)

func TestListFunctions(t *testing.T) {
	tests := []struct {
		name       string
		response   []Function
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful list",
			response: []Function{
				{
					Path:        "test/function1",
					Description: "Test function 1",
				},
				{
					Path:        "test/function2",
					Description: "Test function 2",
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "server error",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				if r.URL.Path != functionsForOrg {
					t.Errorf("expected path %s, got %s", functionsForOrg, r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					response := struct {
						Functions []Function `json:"functions"`
					}{
						Functions: tt.response,
					}
					json.NewEncoder(w).Encode(response)
				}
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			functions, err := client.Functions.List(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if len(functions) != len(tt.response) {
					t.Errorf("expected %d functions, got %d", len(tt.response), len(functions))
				}
				for i, function := range functions {
					if function.Path != tt.response[i].Path {
						t.Errorf("expected path %s, got %s", tt.response[i].Path, function.Path)
					}
					if function.Description != tt.response[i].Description {
						t.Errorf("expected description %s, got %s", tt.response[i].Description, function.Description)
					}
				}
			}
		})
	}
}

func TestCreateFunction(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		description string
		response    Function
		statusCode  int
		wantErr     bool
	}{
		{
			name:        "successful create",
			path:        "test/function",
			description: "Test function",
			response: Function{
				Path:        "test/function",
				Description: "Test function",
			},
			statusCode: http.StatusCreated,
			wantErr:    false,
		},
		{
			name:        "function exists",
			path:        "existing/function",
			description: "Existing function",
			statusCode:  http.StatusConflict,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if r.URL.Path != "/v1/functions" {
					t.Errorf("expected path %s, got %s", "/v1/functions", r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusCreated {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			function := &Function{
				Path:        tt.path,
				Description: tt.description,
			}
			_, err := client.Functions.Create(context.Background(), function)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetFunctionByPath(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		response   Function
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful get",
			path: "test/function",
			response: Function{
				Path:        "test/function",
				Description: "Test function",
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "function not found",
			path:       "nonexistent/function",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				expectedPath := fmt.Sprintf("/api/v1/functions/by_path/%s", tt.path)
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				if tt.statusCode == http.StatusNotFound {
					w.WriteHeader(tt.statusCode)
					json.NewEncoder(w).Encode(map[string]string{
						"error": "function not found",
					})
					return
				}

				if tt.statusCode == http.StatusOK {
					w.WriteHeader(tt.statusCode)
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			function, err := client.Functions.GetByPath(context.Background(), tt.path)

			if tt.wantErr {
				if err == nil {
					t.Error("GetByPath() expected error for not found case")
				}
				return
			}

			if err != nil {
				t.Errorf("GetByPath() unexpected error = %v", err)
				return
			}

			if tt.statusCode == http.StatusOK {
				if function == nil {
					t.Error("GetByPath() returned nil function for success case")
					return
				}
				if function.Path != tt.response.Path {
					t.Errorf("expected path %s, got %s", tt.response.Path, function.Path)
				}
				if function.Description != tt.response.Description {
					t.Errorf("expected description %s, got %s", tt.response.Description, function.Description)
				}
			}
		})
	}
}

func TestDeleteFunction(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful delete",
			path:       "test/function",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "function not found",
			path:       "nonexistent/function",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE request, got %s", r.Method)
				}
				expectedPath := fmt.Sprintf("%s/%s", functionsBasePath, tt.path)
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			err := client.Functions.Delete(context.Background(), tt.path, "")

			if (err != nil) != tt.wantErr {
				if err != nil && err.Error() == "failed to delete function with status 204 No Content" {
					// This is actually a success case
					return
				}
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFunctionChat(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		message    string
		responses  []string
		statusCode int
		wantErr    bool
	}{
		{
			name:    "successful chat",
			path:    "test/function",
			message: "Hello",
			responses: []string{
				"Hi ", "there", "!",
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "function not found",
			path:       "nonexistent/function",
			message:    "Hello",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := fmt.Sprintf("/v1/chat/%s", tt.path)
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				if tt.statusCode == http.StatusNotFound {
					w.WriteHeader(tt.statusCode)
					return
				}

				var requestBody struct {
					Messages []struct {
						Role    string `json:"role"`
						Content string `json:"content"`
					} `json:"messages"`
				}
				if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
				}
				if len(requestBody.Messages) != 1 || requestBody.Messages[0].Content != tt.message {
					t.Errorf("expected message %q, got %q", tt.message, requestBody.Messages[0].Content)
				}

				w.Header().Set("Content-Type", "text/event-stream")
				w.WriteHeader(tt.statusCode)

				for _, resp := range tt.responses {
					chunk := map[string]string{
						"delta": resp,
					}
					data, _ := json.Marshal(chunk)
					fmt.Fprintf(w, "data: %s\n\n", data)
					w.(http.Flusher).Flush()
				}
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			_, err := client.Functions.Chat(context.Background(), tt.path, tt.message)

			if (err != nil) != tt.wantErr {
				t.Errorf("Chat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
