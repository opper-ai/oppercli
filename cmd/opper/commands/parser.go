package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/opper-ai/oppercli/opperai"
)

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

Examples:
opper -c my/function Respond to questions. Be nice, and use emojis.
opper -d my/function
opper -l my/
opper -g my/function
opper my/function Hello, world!
echo "Hello, world!" | opper my/function
echo "Hello, world!" | opper my/function - print the first word`)
	return nil
}

type CommandParser struct{}

func NewCommandParser() *CommandParser {
	return &CommandParser{}
}

func (p *CommandParser) Parse(args []string) (Commander, error) {
	// Show help if no arguments or help flags are provided
	if len(args) < 2 || args[1] == "-h" || args[1] == "--help" || args[1] == "help" {
		return &HelpCommand{}, nil
	}

	switch args[1] {
	case "-c", "--create":
		return p.parseCreateCommand(args[2:])
	case "-d", "--delete":
		return p.parseDeleteCommand(args[2:])
	case "-l", "--list":
		return p.parseListCommand(args[2:])
	case "-g", "--get":
		return p.parseGetCommand(args[2:])
	default:
		return p.parseChatCommand(args[1:])
	}
}

func (p *CommandParser) parseCreateCommand(args []string) (Commander, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("function name required for creation")
	}

	cmd := &CreateCommand{
		BaseCommand: BaseCommand{
			FunctionPath: args[0],
		},
		Instructions: strings.Join(args[1:], " "),
	}

	return cmd, nil
}

func (p *CommandParser) parseDeleteCommand(args []string) (Commander, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("function name required for deletion")
	}

	cmd := &DeleteCommand{
		BaseCommand: BaseCommand{
			FunctionPath: args[0],
		},
	}

	return cmd, nil
}

func (p *CommandParser) parseListCommand(args []string) (Commander, error) {
	filter := ""
	if len(args) > 0 {
		filter = args[0]
	}

	cmd := &ListCommand{
		Filter: filter,
	}

	return cmd, nil
}

func (p *CommandParser) parseGetCommand(args []string) (Commander, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("function path required for retrieval")
	}

	cmd := &GetCommand{
		BaseCommand: BaseCommand{
			FunctionPath: args[0],
		},
	}

	return cmd, nil
}

func (p *CommandParser) parseChatCommand(args []string) (Commander, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("function name required for chat")
	}

	messageContent := ""
	if len(args) > 1 {
		messageContent = strings.Join(args[1:], " ")
	}

	cmd := &ChatCommand{
		BaseCommand: BaseCommand{
			FunctionPath: args[0],
		},
		MessageContent: messageContent,
	}

	return cmd, nil
}
