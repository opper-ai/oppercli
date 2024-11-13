package opperai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Client is the HTTP client for the Opper AI API.
type Client struct {
	APIKey  string
	BaseURL string
	client  *http.Client
	Indexes *IndexesClient
}

// NewClient creates a new Client with an optional baseURL.
func NewClient(apiKey string, baseURL ...string) *Client {
	defaultBaseURL := "https://api.opper.ai"
	if len(baseURL) > 0 && baseURL[0] != "" {
		defaultBaseURL = baseURL[0]
	}
	client := &Client{
		APIKey:  apiKey,
		BaseURL: defaultBaseURL,
	}
	client.Indexes = newIndexesClient(client)
	return client
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
func (c *Client) CreateFunction(ctx context.Context, function *Function) (*FunctionDescription, error) {
	data, err := json.Marshal(function)
	if err != nil {
		return nil, err
	}

	resp, err := c.DoRequest(ctx, http.MethodPost, "/v1/functions", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create function with status %d: %s", resp.StatusCode, string(body))
	}

	var createdFunction FunctionDescription
	if err := json.NewDecoder(resp.Body).Decode(&createdFunction); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &createdFunction, nil
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

// ListFunctions retrieves a list of functions for the organization.
func (c *Client) ListFunctions(ctx context.Context) ([]FunctionDescription, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/api/v1/functions/for_org", nil)
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

// GetFunctionByPath retrieves a function description by its path.
func (c *Client) GetFunctionByPath(ctx context.Context, functionPath string) (*FunctionDescription, error) {
	endpoint := fmt.Sprintf("/api/v1/functions/by_path/%s", functionPath)
	resp, err := c.DoRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Function not found
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

// ListCustomLanguageModels retrieves all custom language models
func (c *Client) ListCustomLanguageModels(ctx context.Context) ([]CustomLanguageModel, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/v1/custom-language-models", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list models with status %s", resp.Status)
	}

	var models []CustomLanguageModel
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, err
	}

	return models, nil
}

// CreateCustomLanguageModel creates a new custom language model
func (c *Client) CreateCustomLanguageModel(ctx context.Context, model CustomLanguageModel) error {
	data, err := json.Marshal(model)
	if err != nil {
		return err
	}

	resp, err := c.DoRequest(ctx, http.MethodPost, "/v1/custom-language-models", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create model with status %s", resp.Status)
	}

	return nil
}

// DeleteCustomLanguageModel deletes a custom language model by name
func (c *Client) DeleteCustomLanguageModel(ctx context.Context, name string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/custom-language-models/by-name/%s", name), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 204 No Content is a success status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete model with status %s", resp.Status)
	}

	return nil
}

// UpdateCustomLanguageModel updates an existing custom language model
func (c *Client) UpdateCustomLanguageModel(ctx context.Context, name string, model CustomLanguageModel) error {
	data, err := json.Marshal(model)
	if err != nil {
		return err
	}

	resp, err := c.DoRequest(ctx, http.MethodPatch, fmt.Sprintf("/v1/custom-language-models/by-name/%s", name), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update model with status %s", resp.Status)
	}

	return nil
}

// GetCustomLanguageModel retrieves a custom language model by name
func (c *Client) GetCustomLanguageModel(ctx context.Context, name string) (*CustomLanguageModel, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/custom-language-models/by-name/%s", name), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Model not found
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get model with status %s", resp.Status)
	}

	var model CustomLanguageModel
	if err := json.NewDecoder(resp.Body).Decode(&model); err != nil {
		return nil, err
	}

	return &model, nil
}

func (c *Client) uploadFile(url string, fields map[string]string, file *os.File) error {
	// Create a pipe to write the multipart form data
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	// Start a goroutine to write the form data
	go func() {
		defer pw.Close()

		// First write all the form fields from the presigned URL
		for key, value := range fields {
			if err := writer.WriteField(key, value); err != nil {
				pw.CloseWithError(fmt.Errorf("failed to write field %s: %v", key, err))
				return
			}
		}

		// Create the form file field
		part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			pw.CloseWithError(fmt.Errorf("failed to create form file: %v", err))
			return
		}

		// Copy the file content to the form field
		if _, err := io.Copy(part, file); err != nil {
			pw.CloseWithError(fmt.Errorf("failed to copy file content: %v", err))
			return
		}

		// Close the multipart writer
		if err := writer.Close(); err != nil {
			pw.CloseWithError(fmt.Errorf("failed to close writer: %v", err))
			return
		}
	}()

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, pr)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set the content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
