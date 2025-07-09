package opperai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsageClient_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/usage/events" {
			t.Errorf("expected path /api/v1/usage/events, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"time_bucket":"2024-01-01T00:00:00Z","cost":"0.001","count":1,"event_type":"generation"}]`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)
	params := &UsageParams{
		EventType: "generation",
		FromDate:  "2024-01-01",
		ToDate:    "2024-01-31",
	}

	usage, err := client.Usage.List(context.Background(), params)
	if err != nil {
		t.Errorf("Usage.List() error = %v", err)
		return
	}

	if len(*usage) != 1 {
		t.Errorf("expected 1 usage event, got %d", len(*usage))
	}

	event := (*usage)[0]
	if event.TimeBucket != "2024-01-01T00:00:00Z" {
		t.Errorf("expected time bucket 2024-01-01T00:00:00Z, got %s", event.TimeBucket)
	}
	if event.Cost != "0.001" {
		t.Errorf("expected cost 0.001, got %s", event.Cost)
	}
	if event.Count != 1 {
		t.Errorf("expected count 1, got %d", event.Count)
	}
	if event.EventType == nil || *event.EventType != "generation" {
		t.Errorf("expected event type generation, got %v", event.EventType)
	}
}

func TestUsageClient_Summary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/usage/summary" {
			t.Errorf("expected path /api/v1/usage/summary, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"total_cost": "0.005",
			"generation_cost": "0.003",
			"platform_cost": "0.002",
			"span_cost": "0.000",
			"embedding_cost": "0.000",
			"metric_cost": "0.000",
			"dataset_storage_cost": "0.000",
			"image_cost": "0.000",
			"total_events": 10,
			"date_range": ["2024-01-01T00:00:00Z", "2024-01-31T23:59:59Z"],
			"event_count_breakdown": {
				"generation": 8,
				"platform": 2
			}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)
	params := &UsageParams{
		FromDate: "2024-01-01",
		ToDate:   "2024-01-31",
	}

	summary, err := client.Usage.Summary(context.Background(), params)
	if err != nil {
		t.Errorf("Usage.Summary() error = %v", err)
		return
	}

	if summary.TotalCost != "0.005" {
		t.Errorf("expected total cost 0.005, got %s", summary.TotalCost)
	}
	if summary.GenerationCost != "0.003" {
		t.Errorf("expected generation cost 0.003, got %s", summary.GenerationCost)
	}
	if summary.PlatformCost != "0.002" {
		t.Errorf("expected platform cost 0.002, got %s", summary.PlatformCost)
	}
	if summary.TotalEvents != 10 {
		t.Errorf("expected total events 10, got %d", summary.TotalEvents)
	}
	if summary.EventCountBreakdown["generation"] != 8 {
		t.Errorf("expected generation count 8, got %d", summary.EventCountBreakdown["generation"])
	}
}

func TestUsageClient_SummaryWithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)
	
	_, err := client.Usage.Summary(context.Background(), nil)
	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
} 