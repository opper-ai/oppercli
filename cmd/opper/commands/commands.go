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
	fmt.Println("    create <name> <litellm-id> <key> [extra] Create a new model")
	fmt.Println("      extra: JSON string with additional options")
	fmt.Println("      Example: '{\"api_base\": \"https://myoaiservice.azure.com\", \"api_version\": \"2024-06-01\"}'")
	fmt.Println("    delete <name>              Delete a model")
	fmt.Println("    get <name>                 Get model details")
	fmt.Println("\n  help                         Show this help message")
	fmt.Println("\nCall functions:")
	fmt.Println("  opper <function-name> [message]  Chat with a function")
	fmt.Println("\nExamples:")
	fmt.Println("  opper functions create my/function \"Respond to questions. Be nice.\"")
	fmt.Println("  opper functions list my/")
	fmt.Println("  opper models create my-model my-id api-key '{\"api_base\": \"https://myoaiservice.azure.com\", \"api_version\": \"2024-06-01\"}'")
	fmt.Println("  opper my/function \"Hello, world!\"")
	return nil
}
