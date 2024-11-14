package opperai

import (
	"context"
	"encoding/json"
	"fmt"
)

type TracesClient struct {
	client *Client
}

func newTracesClient(client *Client) *TracesClient {
	return &TracesClient{client: client}
}

// List returns all traces
func (c *TracesClient) List(ctx context.Context) ([]Trace, error) {
	resp, err := c.client.DoRequest(ctx, "GET", "/v1/traces", nil)
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
