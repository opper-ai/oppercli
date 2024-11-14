package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/opper-ai/oppercli/opperai"
)

// DeleteCommand handles function deletion
type DeleteCommand struct {
	BaseCommand
}

func (c *DeleteCommand) Execute(ctx context.Context, client *opperai.Client) error {
	err := client.Functions.Delete(ctx, "", c.FunctionPath)
	if err != nil {
		return fmt.Errorf("error deleting function: %w", err)
	}
	fmt.Println("Function deleted successfully.")
	return nil
}

// ListCommand handles function listing
type ListCommand struct {
	BaseCommand
	Filter string
}

func (c *ListCommand) Execute(ctx context.Context, client *opperai.Client) error {
	functions, err := client.Functions.List(ctx)
	if err != nil {
		return fmt.Errorf("error listing functions: %w", err)
	}

	// Find the longest path for padding
	maxPathLen := 0
	for _, function := range functions {
		if len(function.Path) > maxPathLen {
			maxPathLen = len(function.Path)
		}
	}

	// Print header
	fmt.Printf("\n%-*s  %s\n", maxPathLen, "PATH", "DESCRIPTION")
	fmt.Printf("%s  %s\n", strings.Repeat("─", maxPathLen), strings.Repeat("─", 50))

	for _, function := range functions {
		if c.Filter == "" || strings.Contains(function.Path, c.Filter) {
			fmt.Printf("%-*s  %s\n",
				maxPathLen,
				function.Path,
				function.Description)
		}
	}
	fmt.Println()
	return nil
}

// GetCommand handles retrieving function details
type GetCommand struct {
	BaseCommand
}

func (c *GetCommand) Execute(ctx context.Context, client *opperai.Client) error {
	function, err := client.Functions.GetByPath(ctx, c.FunctionPath)
	if err != nil {
		return fmt.Errorf("error retrieving function: %w", err)
	}
	if function == nil {
		return fmt.Errorf("function not found")
	}

	fmt.Printf("Path: %s\nDescription: %s\nInstructions: %s\n\n",
		function.Path, function.Description, function.Instructions)
	return nil
}

// CreateCommand handles function creation
type CreateCommand struct {
	BaseCommand
	Instructions string
}

func (c *CreateCommand) Execute(ctx context.Context, client *opperai.Client) error {
	createdFunction, err := client.Functions.Create(ctx, &opperai.Function{
		Path:         c.FunctionPath,
		Instructions: c.Instructions,
	})
	if err != nil {
		return fmt.Errorf("error creating function: %w", err)
	}
	fmt.Printf("Function created successfully: %s\n", createdFunction.Path)
	return nil
}

// FunctionChatCommand handles function chat
type FunctionChatCommand struct {
	BaseCommand
	Message string
}

func (c *FunctionChatCommand) Execute(ctx context.Context, client *opperai.Client) error {
	response, err := client.Functions.Chat(ctx, c.FunctionPath, c.Message)
	if err != nil {
		return fmt.Errorf("error chatting with function: %w", err)
	}

	fmt.Println(response)
	return nil
}

// ParseFunctionCommand handles function command parsing
func ParseFunctionCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("function subcommand required (list, create, delete, get, chat)")
	}

	subcommand := args[0]
	args = args[1:]

	switch subcommand {
	case "chat":
		if len(args) < 1 {
			return nil, fmt.Errorf("function path required")
		}

		functionPath := args[0]
		args = args[1:]

		var message string

		// Check if we should read from stdin
		if len(args) > 0 && args[0] == "-" {
			// Read from stdin
			scanner := bufio.NewScanner(os.Stdin)
			var input []string
			for scanner.Scan() {
				input = append(input, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return nil, fmt.Errorf("error reading from stdin: %w", err)
			}

			// If there are additional args after "-", append them
			if len(args) > 1 {
				input = append(input, args[1:]...)
			}

			message = strings.Join(input, " ")
		} else if len(args) > 0 {
			// Use args directly as message
			message = strings.Join(args, " ")
		} else {
			return nil, fmt.Errorf("message required (either as arguments or via stdin)")
		}

		return &FunctionChatCommand{
			BaseCommand: BaseCommand{
				FunctionPath: functionPath,
			},
			Message: message,
		}, nil

	case "list":
		var filter string
		if len(args) > 0 {
			filter = args[0]
		}
		return &ListCommand{Filter: filter}, nil

	case "create":
		if len(args) < 2 {
			return nil, fmt.Errorf("function name and instructions required")
		}
		return &CreateCommand{
			BaseCommand: BaseCommand{
				FunctionPath: args[0],
			},
			Instructions: args[1],
		}, nil

	case "delete":
		if len(args) < 1 {
			return nil, fmt.Errorf("function name required")
		}
		return &DeleteCommand{
			BaseCommand: BaseCommand{
				FunctionPath: args[0],
			},
		}, nil

	case "get":
		if len(args) < 1 {
			return nil, fmt.Errorf("function name required")
		}
		return &GetCommand{
			BaseCommand: BaseCommand{
				FunctionPath: args[0],
			},
		}, nil

	default:
		return nil, fmt.Errorf("unknown function subcommand: %s", subcommand)
	}
}
