package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opper-ai/oppercli/opperai"
)

// ListModelsCommand handles listing custom language models
type ListModelsCommand struct {
	Filter string
}

func (c *ListModelsCommand) Execute(ctx context.Context, client *opperai.Client) error {
	models, err := client.ListCustomLanguageModels(ctx)
	if err != nil {
		return fmt.Errorf("error listing models: %w", err)
	}

	for _, model := range models {
		if c.Filter == "" || strings.Contains(model.Name, c.Filter) {
			fmt.Printf("Name: %s, Identifier: %s\n", model.Name, model.Identifier)
		}
	}
	return nil
}

// CreateModelCommand handles creating custom language models
type CreateModelCommand struct {
	Name       string
	Identifier string
	APIKey     string
	Extra      string
}

func (c *CreateModelCommand) Execute(ctx context.Context, client *opperai.Client) error {
	var extra map[string]interface{}
	if err := json.Unmarshal([]byte(c.Extra), &extra); err != nil {
		return fmt.Errorf("invalid extra JSON: %w", err)
	}

	model := opperai.CustomLanguageModel{
		Name:       c.Name,
		Identifier: c.Identifier,
		APIKey:     c.APIKey,
		Extra:      extra,
	}

	if err := client.CreateCustomLanguageModel(ctx, model); err != nil {
		return fmt.Errorf("error creating model: %w", err)
	}

	fmt.Printf("Successfully created model: %s\n", c.Name)
	return nil
}

// DeleteModelCommand handles deleting custom language models
type DeleteModelCommand struct {
	Name string
}

func (c *DeleteModelCommand) Execute(ctx context.Context, client *opperai.Client) error {
	if err := client.DeleteCustomLanguageModel(ctx, c.Name); err != nil {
		return fmt.Errorf("error deleting model: %w", err)
	}

	fmt.Printf("Successfully deleted model: %s\n", c.Name)
	return nil
}
