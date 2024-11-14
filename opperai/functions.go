package opperai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	resp, err := c.client.DoRequest(ctx, "GET", "/api/v1/functions/for_org", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list functions with status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Functions []FunctionDescription `json:"functions"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Functions, nil
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

// Add other function methods (Delete, List, Get, etc.)...
