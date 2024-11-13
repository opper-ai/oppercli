package opperai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type IndexesClient struct {
	client *Client
}

func newIndexesClient(client *Client) *IndexesClient {
	return &IndexesClient{client: client}
}

func (c *IndexesClient) List() ([]Index, error) {
	resp, err := c.client.DoRequest(context.Background(), "GET", "/v1/indexes", nil)
	if err != nil {
		return nil, err
	}

	var indexes []Index
	if err := json.NewDecoder(resp.Body).Decode(&indexes); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return indexes, nil
}

func (c *IndexesClient) Create(name string) (*Index, error) {
	body := map[string]interface{}{
		"name": name,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.DoRequest(context.Background(), "POST", "/v1/indexes/by-name", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var index Index
	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &index, nil
}

func (c *IndexesClient) Delete(name string) error {
	_, err := c.client.DoRequest(context.Background(), "DELETE", fmt.Sprintf("/v1/indexes/by-name/%s", name), nil)
	return err
}

func (c *IndexesClient) Get(name string) (*Index, error) {
	resp, err := c.client.DoRequest(context.Background(), "GET", fmt.Sprintf("/v1/indexes/by-name/%s", name), nil)
	if err != nil {
		return nil, err
	}

	var index Index
	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &index, nil
}

func (c *IndexesClient) Query(name string, query string, filters []Filter) ([]RetrievalResponse, error) {
	body := map[string]interface{}{
		"q":       query,
		"filters": filters,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.DoRequest(context.Background(), "POST", fmt.Sprintf("/v1/indexes/query/by-name/%s", name), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var results []RetrievalResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return results, nil
}

func (c *IndexesClient) Add(name string, doc Document) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	_, err = c.client.DoRequest(context.Background(), "POST", fmt.Sprintf("/v1/indexes/index/by-name/%s", name), bytes.NewReader(data))
	return err
}

func (c *IndexesClient) UploadFile(name string, filePath string) error {
	// First get upload URL with filename as query parameter
	filename := filepath.Base(filePath)
	resp, err := c.client.DoRequest(
		context.Background(),
		"GET",
		fmt.Sprintf("/v1/indexes/upload_url/by-name/%s?filename=%s", name, filename),
		nil,
	)
	if err != nil {
		return err
	}

	var uploadData struct {
		URL    string            `json:"url"`
		Fields map[string]string `json:"fields"`
		UUID   string            `json:"uuid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&uploadData); err != nil {
		return fmt.Errorf("failed to decode upload URL response: %v", err)
	}

	// Read file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Upload file to URL
	if err := c.client.uploadFile(uploadData.URL, uploadData.Fields, file); err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	// Register file
	registerData, err := json.Marshal(map[string]string{
		"uuid": uploadData.UUID,
	})
	if err != nil {
		return err
	}

	_, err = c.client.DoRequest(context.Background(), "POST", fmt.Sprintf("/v1/indexes/register_file/by-name/%s", name), bytes.NewReader(registerData))
	return err
}
