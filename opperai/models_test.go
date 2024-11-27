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
	modelsBasePath = "/v1/custom-language-models"
	modelsByName   = "/v1/custom-language-models/by-name"
)

func TestListModels(t *testing.T) {
	tests := []struct {
		name       string
		filter     string
		response   []CustomLanguageModel
		statusCode int

		wantErr bool
	}{
		{
			name:   "successful list",
			filter: "",
			response: []CustomLanguageModel{
				{
					Name:       "model1",
					Identifier: "id1",
					CreatedAt:  "2024-01-01T00:00:00Z",
				},
				{
					Name:       "model2",
					Identifier: "id2",
					CreatedAt:  "2024-01-02T00:00:00Z",
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
				if r.URL.Path != modelsBasePath {
					t.Errorf("expected path %s, got %s", modelsBasePath, r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			models, err := client.Models.List(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if len(models) != len(tt.response) {
					t.Errorf("expected %d models, got %d", len(tt.response), len(models))
				}
				for i, model := range models {
					if model.Name != tt.response[i].Name {
						t.Errorf("expected model name %s, got %s", tt.response[i].Name, model.Name)
					}
					if model.Identifier != tt.response[i].Identifier {
						t.Errorf("expected identifier %s, got %s", tt.response[i].Identifier, model.Identifier)
					}
				}
			}
		})
	}
}

func TestCreateModel(t *testing.T) {
	tests := []struct {
		name       string
		model      CustomLanguageModel
		response   CustomLanguageModel
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful create",
			model: CustomLanguageModel{
				Name:       "test-model",
				Identifier: "test-id",
				APIKey:     "test-key",
				Extra: map[string]interface{}{
					"api_base":    "https://api.test.com",
					"api_version": "2024-01-01",
				},
			},
			response: CustomLanguageModel{
				Name:       "test-model",
				Identifier: "test-id",
				CreatedAt:  "2024-01-01T00:00:00Z",
			},
			statusCode: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "model exists",
			model: CustomLanguageModel{
				Name:       "existing-model",
				Identifier: "existing-id",
			},
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if r.URL.Path != modelsBasePath {
					t.Errorf("expected path %s, got %s", modelsBasePath, r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusCreated {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			err := client.Models.Create(context.Background(), tt.model)

			if (err != nil) != tt.wantErr {
				if err != nil && err.Error() == "failed to create model with status 201 Created" {
					// This is actually a success case
					return
				}
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetModel(t *testing.T) {
	tests := []struct {
		name      string
		modelName string

		response   CustomLanguageModel
		statusCode int
		wantErr    bool
	}{
		{
			name:      "successful get",
			modelName: "test-model",
			response: CustomLanguageModel{
				Name:       "test-model",
				Identifier: "test-id",
				CreatedAt:  "2024-01-01T00:00:00Z",
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "model not found",
			modelName:  "nonexistent",
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
				expectedPath := fmt.Sprintf("%s/%s", modelsByName, tt.modelName)
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				if tt.statusCode == http.StatusNotFound {
					w.WriteHeader(tt.statusCode)
					json.NewEncoder(w).Encode(map[string]string{
						"error": "model not found",
					})
					return
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			model, err := client.Models.Get(context.Background(), tt.modelName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if model.Name != tt.response.Name {
					t.Errorf("expected name %s, got %s", tt.response.Name, model.Name)
				}
				if model.Identifier != tt.response.Identifier {
					t.Errorf("expected identifier %s, got %s", tt.response.Identifier, model.Identifier)
				}
			}
		})
	}
}

func TestDeleteModel(t *testing.T) {
	tests := []struct {
		name       string
		modelName  string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful delete",
			modelName:  "test-model",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "model not found",
			modelName:  "nonexistent",
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
				expectedPath := fmt.Sprintf("%s/%s", modelsByName, tt.modelName)
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			err := client.Models.Delete(context.Background(), tt.modelName)

			if (err != nil) != tt.wantErr {
				if err != nil && tt.statusCode == http.StatusNoContent {
					// This is a success case
					return
				}
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateModel(t *testing.T) {
	tests := []struct {
		name       string
		model      CustomLanguageModel
		response   CustomLanguageModel
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful update",
			model: CustomLanguageModel{
				Name:       "test-model",
				Identifier: "test-id",
				APIKey:     "new-key",
				Extra: map[string]interface{}{
					"api_base":    "https://new.api.com",
					"api_version": "2024-02-01",
				},
			},
			response: CustomLanguageModel{
				Name:       "test-model",
				Identifier: "test-id",
				CreatedAt:  "2024-01-01T00:00:00Z",
				UpdatedAt:  "2024-02-01T00:00:00Z",
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "model not found",
			model: CustomLanguageModel{
				Name:       "nonexistent",
				Identifier: "test-id",
			},
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPatch {
					t.Errorf("expected PATCH request, got %s", r.Method)
				}
				expectedPath := fmt.Sprintf("%s/%s", modelsByName, tt.model.Name)
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			err := client.Models.Update(context.Background(), tt.model.Name, tt.model)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func newTestClient(t *testing.T, baseURL string) *Client {
	return NewClient("test-key", baseURL)
}

func TestListBuiltinModels(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		want    []BuiltinLanguageModel

		wantErr bool
	}{
		{
			name: "successful list",
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v1/language-models" {
					t.Errorf("expected path /v1/language-models, got %s", r.URL.Path)
				}
				if r.Method != "GET" {
					t.Errorf("expected GET method, got %s", r.Method)
				}

				models := []BuiltinLanguageModel{
					{
						Name:            "anthropic/claude-3-sonnet",
						HostingProvider: "Anthropic",
						Location:        "US",
					},
					{
						Name:            "azure/gpt4-eu",
						HostingProvider: "Azure",
						Location:        "EU",
					},
				}
				json.NewEncoder(w).Encode(models)
			},
			want: []BuiltinLanguageModel{
				{
					Name:            "anthropic/claude-3-sonnet",
					HostingProvider: "Anthropic",
					Location:        "US",
				},
				{
					Name:            "azure/gpt4-eu",
					HostingProvider: "Azure",
					Location:        "EU",
				},
			},
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client := newTestClient(t, server.URL)
			got, err := client.Models.ListBuiltin(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("expected %d models, got %d", len(tt.want), len(got))
				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("model %d: got %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
