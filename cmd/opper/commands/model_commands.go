package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opper-ai/oppercli/opperai"
)

func (c *ListModelsCommand) Execute(ctx context.Context, client *opperai.Client) error {
	models, err := client.Models.List(ctx)
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

	if err := client.Models.Create(ctx, model); err != nil {
		return fmt.Errorf("error creating model: %w", err)
	}

	fmt.Printf("Successfully created model: %s\n", c.Name)
	fmt.Printf("To test your model, run: opper models test %s\n", c.Name)

	return nil
}

func (c *DeleteModelCommand) Execute(ctx context.Context, client *opperai.Client) error {
	if err := client.Models.Delete(ctx, c.Name); err != nil {
		return fmt.Errorf("error deleting model: %w", err)
	}

	fmt.Printf("Successfully deleted model: %s\n", c.Name)
	return nil
}

func (c *GetModelCommand) Execute(ctx context.Context, client *opperai.Client) error {
	model, err := client.Models.Get(ctx, c.Name)
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

func (c *TestModelCommand) Execute(ctx context.Context, client *opperai.Client) error {
	// First verify the model exists
	model, err := client.Models.Get(ctx, c.Name)
	if err != nil {
		return fmt.Errorf("error getting model: %w", err)
	}

	fmt.Printf("Testing model %s (%s)...\n\n", c.Name, model.Identifier)

	// Create a call command to test the model
	callCmd := &CallCommand{
		Name:         "opper/cli/model-test",
		Instructions: "The user will input a model name. just confirm that it is working, return the model name, confirm it's working, keep it short and do not ask questions.",
		Input:        c.Name,
		Model:        c.Name,
	}

	if err := callCmd.Execute(ctx, client); err != nil {
		return err
	}

	return nil
}

func (c *ListBuiltinModelsCommand) Execute(ctx context.Context, client *opperai.Client) error {
	models, err := client.Models.ListBuiltin(ctx)
	if err != nil {
		return fmt.Errorf("error listing built-in models: %w", err)
	}

	// Find the longest name and provider for padding
	maxNameLen := 4     // minimum length for "NAME"
	maxProviderLen := 8 // minimum length for "PROVIDER"
	for _, model := range models {
		// Only consider models that match the filter
		if c.Filter != "" && !strings.Contains(strings.ToLower(model.Name), strings.ToLower(c.Filter)) {
			continue
		}
		if len(model.Name) > maxNameLen {
			maxNameLen = len(model.Name)
		}
		if len(model.HostingProvider) > maxProviderLen {
			maxProviderLen = len(model.HostingProvider)
		}
	}

	// Print header
	fmt.Printf("\n%-*s  %-*s  %s\n", maxNameLen, "NAME", maxProviderLen, "PROVIDER", "LOCATION")
	fmt.Printf("%s  %s  %s\n",
		strings.Repeat("─", maxNameLen),
		strings.Repeat("─", maxProviderLen),
		strings.Repeat("─", 8))

	for _, model := range models {
		// Only print models that match the filter
		if c.Filter != "" && !strings.Contains(strings.ToLower(model.Name), strings.ToLower(c.Filter)) {
			continue
		}
		fmt.Printf("%-*s  %-*s  %s\n",
			maxNameLen,
			model.Name,
			maxProviderLen,
			model.HostingProvider,
			model.Location)
	}
	fmt.Println()
	return nil
}
