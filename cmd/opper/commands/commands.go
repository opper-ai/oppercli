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
	fmt.Println(`Usage:
-c <function name> [instructions]     Create a function with the specified name and optional instructions.
-d <function name>                    Delete the specified function.
-l [list filter]                      List functions, optionally filtering by the provided filter.
-g <function path>                    Retrieve a function by its path.
<function name> [prompt]              Initiate a chat with the specified function name and optional prompt.
                                      If message content is not provided directly, it can be read from stdin.
                                      Pass both stdin and a prompt by passing '-' before the prompt.
-lm [filter]                         List custom language models, optionally filtering by name.
-cm <name> <id> <key> [extra]        Create a custom language model.
-dm <name>                           Delete a custom language model.

Examples:
opper -c my/function Respond to questions. Be nice, and use emojis.
opper -d my/function
opper -l my/
opper -g my/function
opper my/function Hello, world!
echo "Hello, world!" | opper my/function
echo "Hello, world!" | opper my/function - print the first word
opper -lm my/
opper -cm my/model my-id api-key "{}"
opper -dm my/model`)
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
