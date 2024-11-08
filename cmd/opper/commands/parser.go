package commands

import (
	"fmt"
	"strings"
)

type CommandParser struct{}

func NewCommandParser() *CommandParser {
	return &CommandParser{}
}

func (p *CommandParser) Parse(args []string) (Command, error) {
	if len(args) < 2 {
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
	case "-lm", "--list-models":
		filter := ""
		if len(args) > 2 {
			filter = args[2]
		}
		return &ListModelsCommand{Filter: filter}, nil
	case "-cm", "--create-model":
		if len(args) < 5 {
			return nil, fmt.Errorf("usage: -cm <name> <identifier> <api_key> [extra_json]")
		}
		extra := "{}"
		if len(args) > 5 {
			extra = args[5]
		}
		return &CreateModelCommand{
			Name:       args[2],
			Identifier: args[3],
			APIKey:     args[4],
			Extra:      extra,
		}, nil
	case "-dm", "--delete-model":
		if len(args) < 3 {
			return nil, fmt.Errorf("usage: -dm <model_name>")
		}
		return &DeleteModelCommand{Name: args[2]}, nil
	default:
		return p.parseChatCommand(args[1:])
	}
}

func (p *CommandParser) parseCreateCommand(args []string) (Command, error) {
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

func (p *CommandParser) parseDeleteCommand(args []string) (Command, error) {
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

func (p *CommandParser) parseListCommand(args []string) (Command, error) {
	filter := ""
	if len(args) > 0 {
		filter = args[0]
	}

	cmd := &ListCommand{
		Filter: filter,
	}

	return cmd, nil
}

func (p *CommandParser) parseGetCommand(args []string) (Command, error) {
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

func (p *CommandParser) parseChatCommand(args []string) (Command, error) {
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
