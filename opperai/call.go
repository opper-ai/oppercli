package opperai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CallClient struct {
	client *Client
}

func newCallClient(client *Client) *CallClient {
	return &CallClient{client: client}
}

type CallResponse struct {
	Message string `json:"message"`
	Stream  chan string
}

func (c *CallClient) Call(ctx context.Context, name string, instructions string, input string, model string, stream bool) (*CallResponse, error) {
	payload := map[string]interface{}{
		"name":         name,
		"instructions": instructions,
		"input":        input,
	}

	if model != "" {
		payload["model"] = model
	}

	if stream {
		payload["stream"] = true
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.DoRequest(ctx, "POST", "/v1/call", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Handle non-200 status codes first
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("%s", string(body))
	}

	if stream {
		streamChan := make(chan string)
		go func() {
			defer close(streamChan)
			defer resp.Body.Close()
			decoder := json.NewDecoder(resp.Body)
			for decoder.More() {
				var chunk struct {
					Delta string `json:"delta"`
				}
				line, err := decoder.Token()
				if err != nil {
					// Handle error
					return
				}
				if str, ok := line.(string); ok && str == "data" {
					if err := decoder.Decode(&chunk); err != nil {
						// Handle error
						return
					}
					streamChan <- chunk.Delta
				}
			}
		}()
		return &CallResponse{Stream: streamChan}, nil
	}

	// For non-streaming responses
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var result struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &CallResponse{Message: result.Message}, nil
}
