package opperai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Client is the HTTP client for the Opper AI API.
type Client struct {
	APIKey  string
	BaseURL string
	client  *http.Client
}

// NewClient creates a new Client with an optional baseURL.
func NewClient(apiKey string, baseURL ...string) *Client {
	defaultBaseURL := "https://api.opper.ai"
	if len(baseURL) > 0 && baseURL[0] != "" {
		defaultBaseURL = baseURL[0]
	}
	return &Client{
		APIKey:  apiKey,
		BaseURL: defaultBaseURL,
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
		c.client = &http.Client{}
	}

	return c.client.Do(req)
}

// Chat initiates a chat session with streaming support and returns a channel for SSE chunks.
func (c *Client) Chat(ctx context.Context, functionPath string, data ChatPayload, stream bool) (<-chan []byte, error) {
	serializedData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	path := "/v1/chat/" + functionPath
	if stream {
		path += "?stream=true"
	}
	resp, err := c.DoRequest(ctx, http.MethodPost, path, bytes.NewBuffer(serializedData))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		resp.Body.Close()
		return nil, ErrRateLimit
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("%w with status %s", ErrFunctionRunFail, resp.Status)
	}

	chunks := make(chan []byte)
	go func() {
		defer close(chunks)
		defer resp.Body.Close()

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break // End of stream
			}
			if err != nil {
				log.Printf("Error reading chunk: %v", err)
				break
			}

			// Filter out non-data lines if necessary
			if bytes.HasPrefix(line, []byte("data:")) {
				chunks <- line
			}
		}
	}()

	return chunks, nil
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

// DeleteFunction deletes a function by its ID or path.
func (c *Client) DeleteFunction(ctx context.Context, id string, path string) error {
	var endpoint string
	if path != "" {
		endpoint = fmt.Sprintf("/api/v1/functions/by_path/%s", path)
	} else if id != "" {
		endpoint = fmt.Sprintf("/api/v1/functions/%s", id)
	} else {
		return fmt.Errorf("either id or path must be provided")
	}

	resp, err := c.DoRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete function with status %s", resp.Status)
	}

	return nil
}
