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

	// First argument after program name is the module
	switch args[1] {
	case "functions":
		return p.parseFunctionsCommand(args[2:])
	case "models":
		return p.parseModelsCommand(args[2:])
	case "indexes":
		return p.parseIndexesCommand(args[2:])
	case "help":
		return &HelpCommand{}, nil
	default:
		// Maintain backwards compatibility by treating the first arg as a function name
		return p.parseChatCommand(args[1:])
	}
}

func (p *CommandParser) parseFunctionsCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("functions subcommand required (list, create, delete, get)")
	}

	switch args[0] {
	case "list":
		filter := ""
		if len(args) > 1 {
			filter = args[1]
		}
		return &ListCommand{Filter: filter}, nil
	case "create":
		if len(args) < 2 {
			return nil, fmt.Errorf("usage: functions create <name> [instructions]")
		}
		return &CreateCommand{
			BaseCommand: BaseCommand{
				FunctionPath: args[1],
			},
			Instructions: strings.Join(args[2:], " "),
		}, nil
	case "delete":
		if len(args) < 2 {
			return nil, fmt.Errorf("usage: functions delete <name>")
		}
		return &DeleteCommand{
			BaseCommand: BaseCommand{
				FunctionPath: args[1],
			},
		}, nil
	case "get":
		if len(args) < 2 {
			return nil, fmt.Errorf("usage: functions get <name>")
		}
		return &GetCommand{
			BaseCommand: BaseCommand{
				FunctionPath: args[1],
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown functions subcommand: %s", args[0])
	}
}

func (p *CommandParser) parseModelsCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("models subcommand required (list, create, delete, get)")
	}

	switch args[0] {
	case "list":
		filter := ""
		if len(args) > 1 {
			filter = args[1]
		}
		return &ListModelsCommand{Filter: filter}, nil
	case "create":
		if len(args) < 4 {
			return nil, fmt.Errorf("usage: models create <name> <litellm identifier> <api_key> [extra_json]\n" +
				"Example extra_json: '{\"api_base\": \"https://myoaiservice.azure.com\", \"api_version\": \"2024-06-01\"}'")
		}

		// Join all remaining arguments as they might be part of the JSON
		extra := "{}"
		if len(args) > 4 {
			extra = strings.Join(args[4:], " ")
		}

		return &CreateModelCommand{
			Name:       args[1],
			Identifier: args[2],
			APIKey:     args[3],
			Extra:      extra,
		}, nil
	case "delete":
		if len(args) < 2 {
			return nil, fmt.Errorf("usage: models delete <name>")
		}
		return &DeleteModelCommand{Name: args[1]}, nil
	case "get":
		if len(args) < 2 {
			return nil, fmt.Errorf("usage: models get <name>")
		}
		return &GetModelCommand{Name: args[1]}, nil
	default:
		return nil, fmt.Errorf("unknown models subcommand: %s", args[0])
	}
}

// Keep the existing parseChatCommand for backwards compatibility
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

func (p *CommandParser) parseIndexesCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("indexes subcommand required (list, create, delete, get, query, add, upload)")
	}

	switch args[0] {
	case "list":
		filter := ""
		if len(args) > 1 {
			filter = args[1]
		}
		return &ListIndexesCommand{Filter: filter}, nil
	case "create":
		if len(args) < 2 {
			return nil, fmt.Errorf("usage: indexes create <name>")
		}
		return &CreateIndexCommand{Name: args[1]}, nil
	case "delete":
		if len(args) < 2 {
			return nil, fmt.Errorf("usage: indexes delete <name>")
		}
		return &DeleteIndexCommand{Name: args[1]}, nil
	case "get":
		if len(args) < 2 {
			return nil, fmt.Errorf("usage: indexes get <name>")
		}
		return &GetIndexCommand{Name: args[1]}, nil
	case "query":
		if len(args) < 3 {
			return nil, fmt.Errorf("usage: indexes query <name> <query> [filter_json]")
		}
		filter := "{}"
		if len(args) > 3 {
			filter = args[3]
		}
		return &QueryIndexCommand{
			Name:   args[1],
			Query:  args[2],
			Filter: filter,
		}, nil
	case "add":
		if len(args) < 4 {
			return nil, fmt.Errorf("usage: indexes add <name> <key> <content> [metadata_json]")
		}
		metadata := "{}"
		if len(args) > 4 {
			metadata = args[4]
		}
		return &AddToIndexCommand{
			Name:     args[1],
			Key:      args[2],
			Content:  args[3],
			Metadata: metadata,
		}, nil
	case "upload":
		if len(args) < 3 {
			return nil, fmt.Errorf("usage: indexes upload <name> <file_path>")
		}
		return &UploadToIndexCommand{
			Name:     args[1],
			FilePath: args[2],
		}, nil
	default:
		return nil, fmt.Errorf("unknown indexes subcommand: %s", args[0])
	}
}
