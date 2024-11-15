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
	fmt.Println(GetHelpText())
	return nil
}

func GetHelpText() string {
	return `Usage:
  opper <command> <subcommand> [arguments]

Commands:
  functions:
    list [filter]              List functions, optionally filtering by name
    create <name> [instructions] Create a new function
    delete <name>              Delete a function
    get <name>                 Get function details
    chat <name> [message]      Chat with a function
      Input: echo "message" | opper functions chat <name>
             opper functions chat <name> <message...>

  models:
    list [filter]              List custom language models
    create <name> <litellm-id> <key> [extra] Create a new model
    delete <name>              Delete a model
    get <name>                 Get model details

  indexes:
    list [filter]              List indexes
    create <name>              Create a new index
    delete <name>              Delete an index
    get <name>                 Get index details
    query <name> <query>       Query an index
    add <name> <key> <content> Add content to an index
    upload <name> <file>       Upload and index a file

  traces:
    list [--live]              List all traces with optional live updates
    get [--live] <trace-id>    Get details and spans of a trace with optional live updates

  call <name> <instructions>   Call a function
    Input: echo "input" | opper call <name> <instructions>
           opper call <name> <instructions> <input...>

  help                         Show this help message`
}
