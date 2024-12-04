package opperai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

type UsageClient struct {
	client *Client
}

func newUsageClient(client *Client) *UsageClient {
	return &UsageClient{client: client}
}

func (c *UsageClient) List(ctx context.Context, params *UsageParams) (*UsageResponse, error) {
	url := "/api/v1/usage/events"
	if params != nil {
		values, err := query.Values(params)
		if err != nil {
			return nil, fmt.Errorf("failed to encode query parameters: %w", err)
		}
		if len(values) > 0 {
			url = fmt.Sprintf("%s?%s", url, values.Encode())
		}
	}

	resp, err := c.client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimit
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list usage with status %s", resp.Status)
	}

	var usage UsageResponse
	if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
		return nil, err
	}

	return &usage, nil
}
