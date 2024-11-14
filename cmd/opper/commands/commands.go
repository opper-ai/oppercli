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
      Message can be provided as arguments or via stdin:
      Examples:
        opper functions chat my/function "Hello, world!"
        echo "Hello" | opper functions chat my/function
        echo '{"data":"test"}' | opper functions chat my/function - format as JSON

  models:
    list [filter]              List custom language models
    create <name> <litellm-id> <key> [extra] Create a new model
      extra: JSON string with additional options
      Example: '{"api_base": "https://myoaiservice.azure.com", "api_version": "2024-06-01"}'
    delete <name>              Delete a model
    get <name>                 Get model details

  indexes:
    list [filter]              List indexes, optionally filtering by name
    create <name>              Create a new index
    delete <name>              Delete an index
    get <name>                 Get index details
    query <name> <query> [filter_json]  Query an index
    add <name> <key> <content> [metadata_json]  Add content to an index
    upload <name> <file_path>  Upload and index a file (PDF, CSV, TXT)

  help                         Show this help message

Examples:
  opper functions create my/function "Respond to questions. Be nice."
  opper functions list my/
  opper models create my-model my-id api-key '{"api_base": "https://myoaiservice.azure.com", "api_version": "2024-06-01"}'`
}
