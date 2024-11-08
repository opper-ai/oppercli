package commands

import (
	"fmt"
	"strings"
)

type CommandParser struct{}

func NewCommandParser() *CommandParser {
	return &CommandParser{}
}

func (p *CommandParser) Parse(args []string) (Commander, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("insufficient arguments")
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
