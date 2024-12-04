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

	// Sub-clients
	Indexes   *IndexesClient
	Models    *ModelsClient
	Functions *FunctionsClient
	Call      *CallClient
	Traces    *TracesClient
	Usage     *UsageClient
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
		client:  &http.Client{},
	}

	// Initialize sub-clients
	client.Indexes = newIndexesClient(client)
	client.Models = newModelsClient(client)
	client.Functions = newFunctionsClient(client)
	client.Call = newCallClient(client)
	client.Traces = newTracesClient(client)
	client.Usage = newUsageClient(client)

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

	return c.client.Do(req)
}

// Chat initiates a chat session with streaming support
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

// uploadFile is a helper function for file uploads
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
