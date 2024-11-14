package opperai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MockClient is a test implementation of the Client interface
type MockClient struct {
	*Client         // Embed the Client struct
	ListModelsFunc  func(context.Context) ([]CustomLanguageModel, error)
	CreateModelFunc func(context.Context, CustomLanguageModel) error
	ChatFunc        func(context.Context, string, string) (string, error)
}

func (m *MockClient) ListModels(ctx context.Context) ([]CustomLanguageModel, error) {
	if m.ListModelsFunc != nil {
		return m.ListModelsFunc(ctx)
	}
	return nil, nil
}

func (m *MockClient) Chat(ctx context.Context, functionName string, message string) (string, error) {
	if m.ChatFunc != nil {
		return m.ChatFunc(ctx, functionName, message)
	}
	return "", nil
}

// NewMockClient creates a new mock client for testing
func NewMockClient() *Client {
	mock := &MockClient{
		Client: &Client{
			APIKey:  "test-key",
			BaseURL: "https://api.opper.ai",
			client:  &http.Client{},
		},
	}

	// Create a new CallClient with our mock implementation
	callClient := newCallClient(mock.Client)
	callClient.client = mock.Client
	mock.Client.Call = callClient

	return mock.Client
}

// MockResponse represents a mock response for testing
type MockResponse struct {
	Message string
	Stream  chan string
}

// DoRequest implements the request method for testing
func (m *MockClient) DoRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	if path == "/v1/call" {
		// Parse the request body to get the expected response
		var reqBody struct {
			Name         string `json:"name"`
			Instructions string `json:"instructions"`
			Input        string `json:"input"`
			Stream       bool   `json:"stream,omitempty"`
		}
		if err := json.NewDecoder(body).Decode(&reqBody); err != nil {
			return nil, err
		}

		// Create a mock response
		mockResp := map[string]string{
			"message": fmt.Sprintf("Mock response for %s", reqBody.Name),
		}
		respBody, _ := json.Marshal(mockResp)

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(respBody)),
		}, nil
	}
	return nil, fmt.Errorf("mock: unexpected request to %s", path)
}
