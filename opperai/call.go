package opperai

import (
	"bytes"
	"context"
	"encoding/json"
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
	defer resp.Body.Close()

	if stream {
		streamChan := make(chan string)
		go func() {
			defer close(streamChan)
			decoder := json.NewDecoder(resp.Body)
			for decoder.More() {
				var chunk struct {
					Delta string `json:"delta"`
				}
				if err := decoder.Decode(&chunk); err != nil {
					// Handle error
					return
				}
				streamChan <- chunk.Delta
			}
		}()
		return &CallResponse{Stream: streamChan}, nil
	}

	var result struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &CallResponse{Message: result.Message}, nil
}
