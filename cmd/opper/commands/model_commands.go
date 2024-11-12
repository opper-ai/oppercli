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

	// Find the longest name and identifier for padding
	maxNameLen := 4 // minimum length for "NAME"
	maxIdLen := 10  // minimum length for "IDENTIFIER"
	for _, model := range models {
		if len(model.Name) > maxNameLen {
			maxNameLen = len(model.Name)
		}
		if len(model.Identifier) > maxIdLen {
			maxIdLen = len(model.Identifier)
		}
	}

	// Print header
	fmt.Printf("\n%-*s  %-*s  %s\n", maxNameLen, "NAME", maxIdLen, "IDENTIFIER", "CREATED")
	fmt.Printf("%s  %s  %s\n",
		strings.Repeat("─", maxNameLen),
		strings.Repeat("─", maxIdLen),
		strings.Repeat("─", 19))

	for _, model := range models {
		if c.Filter == "" || strings.Contains(model.Name, c.Filter) {
			fmt.Printf("%-*s  %-*s  %s\n",
				maxNameLen,
				model.Name,
				maxIdLen,
				model.Identifier,
				model.CreatedAt)
		}
	}
	fmt.Println()
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

// GetModelCommand handles retrieving custom language model details
type GetModelCommand struct {
	Name string
}

func (c *GetModelCommand) Execute(ctx context.Context, client *opperai.Client) error {
	model, err := client.GetCustomLanguageModel(ctx, c.Name)
	if err != nil {
		return fmt.Errorf("error getting model: %w", err)
	}

	if model == nil {
		fmt.Printf("Model not found: %s\n", c.Name)
		return nil
	}

	// Pretty print the model details
	fmt.Printf("Name: %s\n", model.Name)
	fmt.Printf("Identifier: %s\n", model.Identifier)
	fmt.Printf("Created: %s\n", model.CreatedAt)
	fmt.Printf("Updated: %s\n", model.UpdatedAt)
	if model.Extra != nil {
		extraJSON, err := json.MarshalIndent(model.Extra, "", "  ")
		if err == nil {
			fmt.Printf("Extra:\n%s\n", string(extraJSON))
		}
	}

	return nil
}
