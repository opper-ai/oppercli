package opperai_test // Use different package to ensure we're testing the public API

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/opper-ai/oppercli/opperai"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	apiKey := os.Getenv("OPPER_API_KEY")
	if apiKey == "" {
		t.Skip("OPPER_API_KEY not set")
	}

	// Create a test server instead of using real API for integration tests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/indexes/by-name":
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(opperai.Index{Name: "test-index"})
			}
		case "/v1/indexes/by-name/test-index":
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
			}
		case "/v1/indexes/index/by-name/test-index":
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusOK)
			}
		case "/v1/indexes/query/by-name/test-index":
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]opperai.RetrievalResponse{{Content: "test content", Score: 1.0}})
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := opperai.NewClient(apiKey, server.URL)

	t.Run("full index lifecycle", func(t *testing.T) {
		// Create index
		index, err := client.Indexes.Create("test-index")
		if err != nil {
			t.Fatalf("failed to create index: %v", err)
		}

		// Add content
		err = client.Indexes.Add(index.Name, opperai.Document{
			Key:     "test-doc",
			Content: "test content",
		})
		if err != nil {
			t.Fatalf("failed to add content: %v", err)
		}

		// Query
		results, err := client.Indexes.Query(index.Name, "test", nil)
		if err != nil {
			t.Fatalf("failed to query: %v", err)
		}
		if len(results) == 0 {
			t.Error("expected results, got none")
		}

		// Cleanup
		err = client.Indexes.Delete(index.Name)
		if err != nil {
			t.Fatalf("failed to delete index: %v", err)
		}
	})
}
