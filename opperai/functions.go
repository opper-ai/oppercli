package opperai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type FunctionsClient struct {
	client *Client
}

func newFunctionsClient(client *Client) *FunctionsClient {
	return &FunctionsClient{client: client}
}

func (c *FunctionsClient) Create(ctx context.Context, function *Function) (*FunctionDescription, error) {
	data, err := json.Marshal(function)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.DoRequest(ctx, "POST", "/v1/functions", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create function with status %d: %s", resp.StatusCode, string(body))
	}

	var createdFunction FunctionDescription
	if err := json.NewDecoder(resp.Body).Decode(&createdFunction); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &createdFunction, nil
}

func (c *FunctionsClient) Delete(ctx context.Context, id string, path string) error {
	var endpoint string
	if path != "" {
		endpoint = fmt.Sprintf("/api/v1/functions/by_path/%s", path)
	} else if id != "" {
		endpoint = fmt.Sprintf("/api/v1/functions/%s", id)
	} else {
		return fmt.Errorf("either id or path must be provided")
	}

	resp, err := c.client.DoRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete function with status %s", resp.Status)
	}

	return nil
}

func (c *FunctionsClient) List(ctx context.Context) ([]FunctionDescription, error) {
	resp, err := c.client.DoRequest(ctx, "GET", "/v1/functions?limit=1000", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list functions with status %s", resp.Status)
	}

	// Create a struct to match the API response structure
	var response struct {
		Meta struct {
			TotalCount int `json:"total_count"`
		} `json:"meta"`
		Data []FunctionDescription `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response.Data, nil
}

func (c *FunctionsClient) GetByPath(ctx context.Context, functionPath string) (*FunctionDescription, error) {
	endpoint := fmt.Sprintf("/api/v1/functions/by_path/%s", functionPath)
	resp, err := c.client.DoRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("function not found: %s", functionPath)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get function %s with status %s", functionPath, resp.Status)
	}

	var function FunctionDescription
	err = json.NewDecoder(resp.Body).Decode(&function)
	if err != nil {
		return nil, err
	}

	return &function, nil
}

func (c *FunctionsClient) Chat(ctx context.Context, functionPath string, message string) (string, error) {
	functionPath = strings.Trim(functionPath, "/")

	chatPayload := ChatPayload{
		Messages: []Message{
			{
				Role:    "user",
				Content: message,
			},
		},
	}

	chunks, err := c.client.Chat(ctx, functionPath, chatPayload, true)
	if err != nil {
		return "", err
	}

	for chunk := range chunks {
		trimmedChunk := strings.TrimPrefix(string(chunk), "data: ")
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(trimmedChunk), &result); err != nil {
			continue
		}
		if delta, ok := result["delta"].(string); ok {
			fmt.Print(delta)
		}
	}
	fmt.Println()

	return "", nil
}

func (c *FunctionsClient) ListEvaluations(ctx context.Context, functionUUID string) (*EvaluationsResponse, error) {
	endpoint := fmt.Sprintf("/api/v1/functions/%s/evaluations", functionUUID)
	resp, err := c.client.DoRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list evaluations with status %s", resp.Status)
	}

	var evaluations EvaluationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&evaluations); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &evaluations, nil
}
