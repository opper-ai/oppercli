package commands

import (
	"context"
	"fmt"

	"github.com/opper-ai/oppercli/opperai"
)

type Commander interface {
	Execute(ctx context.Context, client *opperai.Client) error
}

// Base command struct for shared fields
type BaseCommand struct {
	FunctionPath string
}

// CreateCommand handles function creation
type CreateCommand struct {
	BaseCommand
	Instructions string
	Model        string // For future custom model support
}

func (c *CreateCommand) Execute(ctx context.Context, client *opperai.Client) error {
	function := opperai.FunctionDescription{
		Path:         c.FunctionPath,
		Instructions: c.Instructions,
		// Model field will be added to FunctionDescription in types.go
	}

	_, err := client.CreateFunction(ctx, function)
	if err != nil {
		return fmt.Errorf("error creating function: %w", err)
	}

	fmt.Println("Function created successfully.")
	return nil
}

// Add other command structs (DeleteCommand, ListCommand, etc.)
