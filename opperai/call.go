package opperai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

			reader := bufio.NewReader(resp.Body)
			for {
				line, err := reader.ReadBytes('\n')
				if err != nil {
					if err != io.EOF {
						fmt.Fprintf(os.Stderr, "Error reading stream: %v\n", err)
					}
					return
				}

				// Skip empty lines
				if len(line) == 0 {
					continue
				}

				// Remove "data: " prefix if present
				data := bytes.TrimPrefix(line, []byte("data: "))

				// Try to parse the JSON
				var chunk struct {
					Delta string `json:"delta"`
				}
				if err := json.Unmarshal(data, &chunk); err != nil {
					continue // Skip malformed JSON
				}

				if chunk.Delta != "" {
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
