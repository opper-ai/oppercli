package opperai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the HTTP client for the Opper AI API.
type Client struct {
	APIKey  string
	BaseURL string
	client  *http.Client
	Timeout time.Duration
}

// NewClient creates a new Client with an optional baseURL and timeout.
func NewClient(apiKey string, baseURL ...string) *Client {
	defaultBaseURL := "https://api.opper.ai"
	if len(baseURL) > 0 && baseURL[0] != "" {
		defaultBaseURL = baseURL[0]
	}
	return &Client{
		APIKey:  apiKey,
		BaseURL: defaultBaseURL,
		Timeout: 60 * time.Second,
	}
}

// DoRequest executes an HTTP request.
func (c *Client) DoRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-OPPER-API-KEY", c.APIKey)

	if c.client == nil {
		c.client = &http.Client{
			Timeout: c.Timeout,
		}
	}

	return c.client.Do(req)
}

// Chat initiates a chat session.
func (c *Client) Chat(ctx context.Context, functionPath string, data ChatPayload) (*FunctionResponse, error) {
	serializedData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := c.DoRequest(ctx, http.MethodPost, "/v1/chat/"+functionPath, bytes.NewBuffer(serializedData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimit
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w with status %s", ErrFunctionRunFail, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var functionResponse FunctionResponse
	err = json.Unmarshal(body, &functionResponse)
	if err != nil {
		return nil, err
	}

	return &functionResponse, nil
}

// CreateFunction creates a new function.
func (c *Client) CreateFunction(ctx context.Context, function FunctionDescription) (int, error) {
	serializedData, err := json.Marshal(function)
	if err != nil {
		return 0, err
	}

	resp, err := c.DoRequest(ctx, http.MethodPost, "/api/v1/functions", bytes.NewBuffer(serializedData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to create function %s with status %s", function.Path, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response map[string]int
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	return response["id"], nil
}
