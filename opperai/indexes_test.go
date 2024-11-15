package opperai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	indexesBasePath = "/v1/indexes"
	indexesByName   = "/v1/indexes/by-name"
	indexesQuery    = "/v1/indexes/query/by-name"
)

func TestListIndexes(t *testing.T) {
	tests := []struct {
		name       string
		response   []Index
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful list",
			response:   []Index{{Name: "index1"}, {Name: "index2"}},
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
				if r.URL.Path != indexesBasePath {
					t.Errorf("expected path %s, got %s", indexesBasePath, r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			indexes, err := client.Indexes.List("")

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if len(indexes) != len(tt.response) {
					t.Errorf("expected %d indexes, got %d", len(tt.response), len(indexes))
				}
				for i, index := range indexes {
					if index.Name != tt.response[i].Name {
						t.Errorf("expected index %s, got %s", tt.response[i].Name, index.Name)
					}
				}
			}
		})
	}
}

func TestCreateIndex(t *testing.T) {
	tests := []struct {
		name       string
		indexName  string
		response   Index
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful create",
			indexName:  "test-index",
			response:   Index{Name: "test-index"},
			statusCode: http.StatusCreated,
			wantErr:    false,
		},
		{
			name:       "index exists",
			indexName:  "existing-index",
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
		{
			name:       "server error",
			indexName:  "test-index",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := indexesBasePath
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusCreated || tt.statusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			index, err := client.Indexes.Create(tt.indexName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if index.Name != tt.response.Name {
					t.Errorf("expected index name %s, got %s", tt.response.Name, index.Name)
				}
			}
		})
	}
}

func TestDeleteIndex(t *testing.T) {
	tests := []struct {
		name       string
		indexName  string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful delete",
			indexName:  "test-index",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "index not found",
			indexName:  "nonexistent",
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
				expectedPath := fmt.Sprintf("%s/%s", indexesByName, tt.indexName)
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				if tt.statusCode == http.StatusNotFound {
					w.WriteHeader(tt.statusCode)
					json.NewEncoder(w).Encode(map[string]string{
						"error": "index not found",
					})
					return
				}

				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			err := client.Indexes.Delete(tt.indexName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueryIndex(t *testing.T) {
	tests := []struct {
		name       string
		indexName  string
		query      string
		filters    []Filter
		response   []RetrievalResponse
		statusCode int
		wantErr    bool
	}{
		{
			name:      "successful query",
			indexName: "test-index",
			query:     "test query",
			filters: []Filter{
				{Field: "key", Value: "value"},
			},
			response: []RetrievalResponse{
				{Content: "result 1", Score: 0.9},
				{Content: "result 2", Score: 0.8},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "index not found",
			indexName:  "nonexistent",
			query:      "test",
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
				expectedPath := fmt.Sprintf("%s/%s", indexesQuery, tt.indexName)
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
			results, err := client.Indexes.Query(tt.indexName, tt.query, tt.filters)

			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if len(results) != len(tt.response) {
					t.Errorf("expected %d results, got %d", len(tt.response), len(results))
				}
				for i, result := range results {
					if result.Content != tt.response[i].Content {
						t.Errorf("expected content %s, got %s", tt.response[i].Content, result.Content)
					}
					if result.Score != tt.response[i].Score {
						t.Errorf("expected score %f, got %f", tt.response[i].Score, result.Score)
					}
				}
			}
		})
	}
}
