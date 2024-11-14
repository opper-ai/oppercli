package commands

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type CommandParser struct{}

func NewCommandParser() *CommandParser {
	return &CommandParser{}
}

func (p *CommandParser) parseFunctionsCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("functions subcommand required (list, create, delete, get, chat)")
	}

	switch args[0] {
	case "list":
		var filter string
		if len(args) > 1 {
			filter = args[1]
		}
		return &ListCommand{Filter: filter}, nil

	case "create":
		if len(args) < 2 {
			return nil, fmt.Errorf("function name and instructions required")
		}
		return &CreateCommand{
			BaseCommand: BaseCommand{
				FunctionPath: args[1],
			},
			Instructions: strings.Join(args[2:], " "),
		}, nil

	case "delete":
		if len(args) < 2 {
			return nil, fmt.Errorf("function name required")
		}
		return &DeleteCommand{
			BaseCommand: BaseCommand{
				FunctionPath: args[1],
			},
		}, nil

	case "get":
		if len(args) < 2 {
			return nil, fmt.Errorf("function name required")
		}
		return &GetCommand{
			BaseCommand: BaseCommand{
				FunctionPath: args[1],
			},
		}, nil

	case "chat":
		if len(args) < 2 {
			return nil, fmt.Errorf("usage: functions chat <name> [message]")
		}
		functionPath := args[1]
		args = args[2:] // Remove the function path from args

		var message string
		if len(args) > 0 && args[0] == "-" {
			// Read from stdin
			stdinData, err := io.ReadAll(os.Stdin)
			if err != nil {
				return nil, fmt.Errorf("error reading from stdin: %w", err)
			}

			// If there are additional args after "-", append them
			if len(args) > 1 {
				message = fmt.Sprintf("%s %s", string(stdinData), strings.Join(args[1:], " "))
			} else {
				message = string(stdinData)
			}
		} else if len(args) > 0 {
			// Use args directly as message
			message = strings.Join(args, " ")
		} else {
			// No args, try reading from stdin
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				scanner := bufio.NewScanner(os.Stdin)
				var input []string
				for scanner.Scan() {
					input = append(input, scanner.Text())
				}
				if err := scanner.Err(); err != nil {
					return nil, fmt.Errorf("error reading from stdin: %w", err)
				}
				message = strings.Join(input, "\n")
			}
		}

		if message == "" {
			return nil, fmt.Errorf("message required (either as arguments or via stdin)")
		}

		return &FunctionChatCommand{
			BaseCommand: BaseCommand{
				FunctionPath: functionPath,
			},
			Message: message,
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
			return nil, fmt.Errorf("usage: models create <name> <litellm identifier> <api_key> [extra_json]")
		}
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

func (p *CommandParser) parseChatCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("function name required")
	}

	functionPath := args[0]
	message := ""
	if len(args) > 1 {
		message = strings.Join(args[1:], " ")
	} else {
		// Check if there's input from stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			var input []string
			for scanner.Scan() {
				input = append(input, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return nil, fmt.Errorf("error reading from stdin: %w", err)
			}
			message = strings.Join(input, "\n")
		}
	}

	if message == "" {
		return nil, fmt.Errorf("message required (either as arguments or via stdin)")
	}

	return &FunctionChatCommand{
		BaseCommand: BaseCommand{
			FunctionPath: functionPath,
		},
		Message: message,
	}, nil
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
	case "traces":
		if len(args) < 3 {
			return nil, fmt.Errorf("traces command requires a subcommand (list, get)")
		}

		// Parse flags
		flagSet := flag.NewFlagSet("traces", flag.ContinueOnError)
		live := flagSet.Bool("live", false, "Watch for updates")
		if err := flagSet.Parse(args[3:]); err != nil {
			return nil, err
		}

		switch args[2] {
		case "list":
			return &ListTracesCommand{Live: *live}, nil
		case "get":
			remainingArgs := flagSet.Args()
			if len(remainingArgs) < 1 {
				return nil, fmt.Errorf("traces get requires a trace ID")
			}
			return &GetTraceCommand{
				TraceID: remainingArgs[0],
				Live:    *live,
			}, nil
		default:
			return nil, fmt.Errorf("unknown traces subcommand: %s", args[2])
		}
	case "help":
		return &HelpCommand{}, nil
	case "call":
		if len(args) < 4 {
			return nil, fmt.Errorf("usage: call <name> <instructions>")
		}

		name := args[2]
		instructions := args[3]
		var input string

		// Check if we have data on stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			var inputLines []string
			for scanner.Scan() {
				inputLines = append(inputLines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return nil, fmt.Errorf("error reading from stdin: %w", err)
			}
			input = strings.Join(inputLines, "\n")
		} else if len(args) > 4 {
			// If no stdin, use remaining args as input
			input = strings.Join(args[4:], " ")
		}

		return &CallCommand{
			Name:         name,
			Instructions: instructions,
			Input:        input,
		}, nil
	default:
		// Maintain backwards compatibility by treating the first arg as a function name
		return p.parseChatCommand(args[1:])
	}
}
