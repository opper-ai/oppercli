package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/opper-ai/oppercli/cmd/opper/commands/output"
	"github.com/opper-ai/oppercli/opperai"
)

func (c *DeleteCommand) Execute(ctx context.Context, client *opperai.Client) error {
	err := client.Functions.Delete(ctx, "", c.FunctionPath)
	if err != nil {
		return fmt.Errorf("error deleting function: %w", err)
	}
	fmt.Println("Function deleted successfully.")
	return nil
}

func (c *ListCommand) Execute(ctx context.Context, client *opperai.Client) error {
	functions, err := client.Functions.List(ctx)
	if err != nil {
		return fmt.Errorf("error listing functions: %w", err)
	}

	// Convert data to rows
	rows := make([][]string, len(functions))
	for i, function := range functions {
		if c.Filter == "" || strings.Contains(function.Path, c.Filter) {
			rows[i] = []string{
				function.Path,
				function.Description,
			}
		}
	}

	// Use the formatter
	output.Table(
		[]string{"PATH", "DESCRIPTION"},
		rows,
	)
	return nil
}

func (c *GetCommand) Execute(ctx context.Context, client *opperai.Client) error {
	function, err := client.Functions.GetByPath(ctx, c.FunctionPath)
	if err != nil {
		return fmt.Errorf("error retrieving function: %w", err)
	}
	if function == nil {
		return fmt.Errorf("function not found")
	}

	// Print basic information
	fmt.Printf("Function: %s\n", function.Path)
	fmt.Printf("%-20s %s\n", "Description:", function.Description)
	if function.UUID != "" {
		fmt.Printf("%-20s %s\n", "UUID:", function.UUID)
	}

	// Print model information
	if function.Model != "" {
		fmt.Printf("%-20s %s", "Model:", function.Model)
		if function.LanguageModelID != 0 {
			fmt.Printf(" (ID: %d)", function.LanguageModelID)
		}
		fmt.Println()
	}

	// Print dataset information
	if function.Dataset.UUID != "" {
		fmt.Printf("\nDataset Information:\n")
		fmt.Printf("%-20s %s\n", "UUID:", function.Dataset.UUID)
		if function.Dataset.EntryCount > 0 {
			fmt.Printf("%-20s %d\n", "Entry Count:", function.Dataset.EntryCount)
		}
	}

	// Print project information
	if function.Project.UUID != "" {
		fmt.Printf("\nProject Information:\n")
		if function.Project.Name != "" {
			fmt.Printf("%-20s %s\n", "Name:", function.Project.Name)
		}
		fmt.Printf("%-20s %s\n", "UUID:", function.Project.UUID)
	}

	// Print few-shot settings
	if function.FewShot || function.FewShotCount > 0 {
		fmt.Printf("\nFew-Shot Settings:\n")
		fmt.Printf("%-20s %v\n", "Enabled:", function.FewShot)
		fmt.Printf("%-20s %d\n", "Count:", function.FewShotCount)
	}

	// Print additional settings
	fmt.Printf("\nAdditional Settings:\n")
	fmt.Printf("%-20s %v\n", "Semantic Search:", function.UseSemanticSearch)
	if function.Revision > 0 {
		fmt.Printf("%-20s %d\n", "Revision:", function.Revision)
	}

	// Print schemas if they exist
	if len(function.InputSchema) > 0 {
		fmt.Printf("\nInput Schema:\n")
		prettyPrintSchema(function.InputSchema)
	}

	if len(function.OutputSchema) > 0 {
		fmt.Printf("\nOutput Schema:\n")
		prettyPrintSchema(function.OutputSchema)
	}

	// Print instructions
	if function.Instructions != "" {
		fmt.Printf("\nInstructions:\n%s\n", function.Instructions)
	}

	return nil
}

func prettyPrintSchema(schema map[string]interface{}) {
	for key, value := range schema {
		fmt.Printf("  %-18s %v\n", key+":", value)
	}
}

func (c *CreateCommand) Execute(ctx context.Context, client *opperai.Client) error {
	if c.Instructions == "" {
		return fmt.Errorf("instructions required")
	}
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

func (c *FunctionChatCommand) Execute(ctx context.Context, client *opperai.Client) error {
	_, err := client.Functions.Chat(ctx, c.FunctionPath, c.Message)
	if err != nil {
		return fmt.Errorf("error chatting with function: %w", err)
	}
	return nil
}

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
		if len(args) == 0 {
			// Check if there's input from stdin
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				// Read from stdin
				scanner := bufio.NewScanner(os.Stdin)
				var input []string
				for scanner.Scan() {
					input = append(input, scanner.Text())
				}
				if err := scanner.Err(); err != nil {
					return nil, fmt.Errorf("error reading from stdin: %w", err)
				}
				message = strings.Join(input, "\n")
			} else {
				return nil, fmt.Errorf("message required (either as arguments or via stdin)")
			}
		} else {
			message = strings.Join(args, " ")
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
