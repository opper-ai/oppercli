package opperai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ModelsClient struct {
	client *Client
}

func newModelsClient(client *Client) *ModelsClient {
	return &ModelsClient{client: client}
}

func (c *ModelsClient) List(ctx context.Context) ([]CustomLanguageModel, error) {
	resp, err := c.client.DoRequest(ctx, "GET", "/v1/custom-language-models", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimit
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list models with status %s", resp.Status)
	}

	var models []CustomLanguageModel
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, err
	}

	return models, nil
}

func (c *ModelsClient) Create(ctx context.Context, model CustomLanguageModel) error {
	data, err := json.Marshal(model)
	if err != nil {
		return err
	}

	resp, err := c.client.DoRequest(ctx, "POST", "/v1/custom-language-models", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create model with status %s", resp.Status)
	}

	return nil
}

func (c *ModelsClient) Delete(ctx context.Context, name string) error {
	resp, err := c.client.DoRequest(ctx, "DELETE", fmt.Sprintf("/v1/custom-language-models/by-name/%s", name), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete model with status %s", resp.Status)
	}

	return nil
}

func (c *ModelsClient) Update(ctx context.Context, name string, model CustomLanguageModel) error {
	data, err := json.Marshal(model)
	if err != nil {
		return err
	}

	resp, err := c.client.DoRequest(ctx, "PATCH", fmt.Sprintf("/v1/custom-language-models/by-name/%s", name), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update model with status %s", resp.Status)
	}

	return nil
}

func (c *ModelsClient) Get(ctx context.Context, name string) (*CustomLanguageModel, error) {
	resp, err := c.client.DoRequest(ctx, "GET", fmt.Sprintf("/v1/custom-language-models/by-name/%s", name), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("model not found: %s", name)
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
