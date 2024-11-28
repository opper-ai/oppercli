package builders

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/spf13/cobra"
)

func BuildFunctionCommands(executeCommand func(commands.Command) error) *cobra.Command {
	functionsCmd := &cobra.Command{
		Use:   "functions",
		Short: "Manage functions",
		Example: `  # List all functions
  opper functions list

  # Create a new function
  opper functions create myfunction "respond to questions about X"

  # Chat with a function
  opper functions chat myfunction "Hello"
  echo "Hello" | opper functions chat myfunction`,
	}

	// List command
	listCmd := &cobra.Command{
		Use:   "list [filter]",
		Short: "List functions",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := ""
			if len(args) > 0 {
				filter = args[0]
			}
			return executeCommand(&commands.ListCommand{Filter: filter})
		},
	}

	// Create command
	createCmd := &cobra.Command{
		Use:   "create <name> [instructions]",
		Short: "Create a new function",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			instructions := ""
			if len(args) > 1 {
				instructions = strings.Join(args[1:], " ")
			}
			return executeCommand(&commands.CreateCommand{
				BaseCommand: commands.BaseCommand{
					FunctionPath: args[0],
				},
				Instructions: instructions,
			})
		},
	}

	// Delete command
	deleteCmd := BuildDeletionCommand(
		"delete <name>",
		"Delete a function",
		"function",
		func(name string) error {
			return executeCommand(&commands.DeleteCommand{
				BaseCommand: commands.BaseCommand{
					FunctionPath: name,
				},
			})
		},
	)

	// Get command
	getCmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Get function details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.GetCommand{
				BaseCommand: commands.BaseCommand{
					FunctionPath: args[0],
				},
			})
		},
	}

	// Chat command
	chatCmd := &cobra.Command{
		Use:   "chat <name> [message]",
		Short: "Chat with a function",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var message string
			if len(args) > 1 {
				message = strings.Join(args[1:], " ")
			} else {
				// Check if there's input from stdin
				stat, _ := os.Stdin.Stat()
				if (stat.Mode() & os.ModeCharDevice) == 0 {
					// Read from stdin
					stdinData, err := io.ReadAll(os.Stdin)
					if err != nil {
						return fmt.Errorf("error reading from stdin: %w", err)
					}
					message = string(stdinData)
				} else {
					// Interactive mode
					scanner := bufio.NewScanner(os.Stdin)
					fmt.Print("> ")
					if scanner.Scan() {
						message = scanner.Text()
					}
					if err := scanner.Err(); err != nil {
						return fmt.Errorf("error reading input: %w", err)
					}
				}
			}
			if message == "" {
				return fmt.Errorf("message required (either as arguments or via stdin)")
			}

			return executeCommand(&commands.FunctionChatCommand{
				BaseCommand: commands.BaseCommand{
					FunctionPath: args[0],
				},
				Message: message,
			})
		},
	}

	// Evaluations command
	evaluationsCmd := &cobra.Command{
		Use:   "evaluations",
		Short: "Manage function evaluations",
	}

	// List evaluations command
	listEvaluationsCmd := &cobra.Command{
		Use:   "list <name>",
		Short: "List evaluations for a function",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.ListEvaluationsCommand{
				BaseCommand: commands.BaseCommand{
					FunctionPath: args[0],
				},
			})
		},
	}

	evaluationsCmd.AddCommand(listEvaluationsCmd)

	functionsCmd.AddCommand(
		listCmd,
		createCmd,
		deleteCmd,
		getCmd,
		chatCmd,
		evaluationsCmd,
	)

	return functionsCmd
}
