package opperai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type TracesClient struct {
	client *Client
}

func newTracesClient(client *Client) *TracesClient {
	return &TracesClient{client: client}
}

// List returns all traces
func (c *TracesClient) List(ctx context.Context, limit int) ([]Trace, error) {
	path := "/v1/traces"
	if limit > 0 {
		path = fmt.Sprintf("%s?limit=%d", path, limit)
	}

	resp, err := c.client.DoRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response TraceListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Traces, nil
}

// Get returns a specific trace
func (c *TracesClient) Get(ctx context.Context, traceID string) (*Trace, error) {
	resp, err := c.client.DoRequest(ctx, "GET", fmt.Sprintf("/v1/traces/%s", traceID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var trace Trace
	if err := json.NewDecoder(resp.Body).Decode(&trace); err != nil {
		return nil, err
	}

	return &trace, nil
}

// Add these new types
type TraceUpdate struct {
	Trace *Trace
	Error error
}

type TraceWatcher interface {
	Watch(ctx context.Context) (<-chan TraceUpdate, error)
}

// Add these new methods to TracesClient
func (c *TracesClient) WatchList(ctx context.Context, seenTraces map[string]bool) (<-chan TraceUpdate, error) {
	updates := make(chan TraceUpdate)

	go func() {
		defer close(updates)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				traces, err := c.List(ctx, 10) // Only fetch the 10 most recent traces
				if err != nil {
					updates <- TraceUpdate{Error: err}
					continue
				}

				// Send new traces in reverse order
				for i := len(traces) - 1; i >= 0; i-- {
					trace := traces[i]
					if !seenTraces[trace.UUID] {
						seenTraces[trace.UUID] = true
						updates <- TraceUpdate{Trace: &trace}
					}
				}
			}
		}
	}()

	return updates, nil
}

func (c *TracesClient) WatchTrace(ctx context.Context, traceID string) (<-chan TraceUpdate, error) {
	updates := make(chan TraceUpdate)

	go func() {
		defer close(updates)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		var lastSpanCount int

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				trace, err := c.Get(ctx, traceID)
				if err != nil {
					updates <- TraceUpdate{Error: err}
					continue
				}

				// Send update if number of spans changed or status changed
				if len(trace.Spans) != lastSpanCount {
					lastSpanCount = len(trace.Spans)
					updates <- TraceUpdate{Trace: trace}
				}
			}
		}
	}()

	return updates, nil
}
