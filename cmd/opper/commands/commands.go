package commands

import (
	"context"
	"fmt"

	"github.com/opper-ai/oppercli/opperai"
)

// Command interface defines what all commands must implement
type Command interface {
	Execute(ctx context.Context, client *opperai.Client) error
}

// Base command struct for shared fields
type BaseCommand struct {
	FunctionPath string
}

// HelpCommand shows usage information
type HelpCommand struct{}

func (c *HelpCommand) Execute(ctx context.Context, client *opperai.Client) error {
	fmt.Println("Usage:")
	fmt.Println("  opper <command> <subcommand> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  functions:")
	fmt.Println("    list [filter]              List functions, optionally filtering by name")
	fmt.Println("    create <name> [instructions] Create a new function")
	fmt.Println("    delete <name>              Delete a function")
	fmt.Println("    get <name>                 Get function details")
	fmt.Println("\n  models:")
	fmt.Println("    list [filter]              List custom language models")
	fmt.Println("    create <name> <id> <key> [extra] Create a new model")
	fmt.Println("      extra: JSON string with additional options")
	fmt.Println("      Example: '{\"temperature\": 0.7, \"model\": \"gpt-4\"}'")
	fmt.Println("    delete <name>              Delete a model")
	fmt.Println("    get <name>                 Get model details")
	fmt.Println("\n  help                         Show this help message")
	fmt.Println("\nLegacy Usage (still supported):")
	fmt.Println("  opper <function-name> [message]  Chat with a function")
	fmt.Println("\nExamples:")
	fmt.Println("  opper functions create my/function \"Respond to questions. Be nice.\"")
	fmt.Println("  opper functions list my/")
	fmt.Println("  opper models create my-model my-id api-key '{\"temperature\": 0.7}'")
	fmt.Println("  opper my/function \"Hello, world!\"")
	return nil
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
