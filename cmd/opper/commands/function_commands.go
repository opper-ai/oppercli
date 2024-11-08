package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/opper-ai/oppercli/opperai"
)

// DeleteCommand handles function deletion
type DeleteCommand struct {
	BaseCommand
}

func (c *DeleteCommand) Execute(ctx context.Context, client *opperai.Client) error {
	err := client.DeleteFunction(ctx, "", c.FunctionPath)
	if err != nil {
		return fmt.Errorf("error deleting function: %w", err)
	}
	fmt.Println("Function deleted successfully.")
	return nil
}

// ListCommand handles function listing
type ListCommand struct {
	BaseCommand
	Filter string
}

func (c *ListCommand) Execute(ctx context.Context, client *opperai.Client) error {
	functions, err := client.ListFunctions(ctx)
	if err != nil {
		return fmt.Errorf("error listing functions: %w", err)
	}

	for _, function := range functions {
		if c.Filter == "" || strings.Contains(function.Path, c.Filter) {
			fmt.Printf("Path: %s, Description: %s\n",
				function.Path, function.Description)
		}
	}
	return nil
}

// GetCommand handles retrieving function details
type GetCommand struct {
	BaseCommand
}

func (c *GetCommand) Execute(ctx context.Context, client *opperai.Client) error {
	function, err := client.GetFunctionByPath(ctx, c.FunctionPath)
	if err != nil {
		return fmt.Errorf("error retrieving function: %w", err)
	}
	if function == nil {
		return fmt.Errorf("function not found")
	}

	fmt.Printf("Path: %s\nDescription: %s\nInstructions: %s\n\n",
		function.Path, function.Description, function.Instructions)
	return nil
}
