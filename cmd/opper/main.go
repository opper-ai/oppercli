package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/opper-ai/oppercli/cmd/opper/config"
	"github.com/opper-ai/oppercli/opperai"
	"github.com/spf13/cobra"
)

func executeCommand(cmd commands.Command) error {
	ctx := context.Background()

	// Get API key from environment or config
	apiKey, err := config.GetAPIKey("default")
	if err != nil {
		return err
	}

	client := opperai.NewClient(apiKey)
	return cmd.Execute(ctx, client)
}

func main() {
	var rootCmd = &cobra.Command{
		Use:           "opper",
		Short:         "Opper CLI - interact with Opper AI services",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Global flags
	var keyName string
	rootCmd.PersistentFlags().StringVar(&keyName, "key", "default", "API key to use from config")

	// Update executeCommand to use baseUrl
	executeCommand := func(cmd commands.Command) error {
		ctx := context.Background()
		apiKey, baseUrl, err := config.GetAPIKeyAndBaseUrl(keyName)
		if err != nil {
			return err
		}
		client := opperai.NewClient(apiKey, baseUrl)
		return cmd.Execute(ctx, client)
	}

	// Indexes command
	var indexesCmd = &cobra.Command{
		Use:   "indexes",
		Short: "Manage indexes",
	}

	var listIndexesCmd = &cobra.Command{
		Use:   "list [filter]",
		Short: "List indexes",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := ""
			if len(args) > 0 {
				filter = args[0]
			}
			return executeCommand(&commands.ListIndexesCommand{Filter: filter})
		},
	}

	var createIndexCmd = &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.CreateIndexCommand{Name: args[0]})
		},
	}

	var deleteIndexCmd = &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete an index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("yes")
			confirmed, err := commands.ConfirmDeletion("index", args[0], force)
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Println("Operation cancelled")
				return nil
			}
			return executeCommand(&commands.DeleteIndexCommand{Name: args[0]})
		},
	}
	deleteIndexCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	var queryIndexCmd = &cobra.Command{
		Use:   "query <name> <query> [filter_json]",
		Short: "Query an index",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := "{}"
			if len(args) > 2 {
				filter = args[2]
			}
			return executeCommand(&commands.QueryIndexCommand{
				Name:   args[0],
				Query:  args[1],
				Filter: filter,
			})
		},
	}

	var getIndexCmd = &cobra.Command{
		Use:   "get <name>",
		Short: "Get index details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.GetIndexCommand{Name: args[0]})
		},
	}

	var addToIndexCmd = &cobra.Command{
		Use:   "add <name> <key> <content> [metadata_json]",
		Short: "Add content to an index",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			metadata := "{}"
			if len(args) > 3 {
				metadata = args[3]
			}
			return executeCommand(&commands.AddToIndexCommand{
				Name:     args[0],
				Key:      args[1],
				Content:  args[2],
				Metadata: metadata,
			})
		},
	}

	var uploadToIndexCmd = &cobra.Command{
		Use:   "upload <name> <file_path>",
		Short: "Upload and index a file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.UploadToIndexCommand{
				Name:     args[0],
				FilePath: args[1],
			})
		},
	}

	indexesCmd.AddCommand(listIndexesCmd, createIndexCmd, deleteIndexCmd, queryIndexCmd,
		getIndexCmd, addToIndexCmd, uploadToIndexCmd)

	// Models command
	var modelsCmd = &cobra.Command{
		Use:   "models",
		Short: "Manage models",
	}

	var listModelsCmd = &cobra.Command{
		Use:   "list [filter]",
		Short: "List models",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := ""
			if len(args) > 0 {
				filter = args[0]
			}
			return executeCommand(&commands.ListModelsCommand{Filter: filter})
		},
	}

	var createModelCmd = &cobra.Command{
		Use:   "create <name> <litellm-id> <api-key> [extra_json]",
		Short: "Create a new model",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			extra := "{}"
			if len(args) > 3 {
				extra = args[3]
			}
			return executeCommand(&commands.CreateModelCommand{
				Name:       args[0],
				Identifier: args[1],
				APIKey:     args[2],
				Extra:      extra,
			})
		},
	}

	var deleteModelCmd = &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a model",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("yes")
			confirmed, err := commands.ConfirmDeletion("model", args[0], force)
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Println("Operation cancelled")
				return nil
			}
			return executeCommand(&commands.DeleteModelCommand{Name: args[0]})
		},
	}
	deleteModelCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	var getModelCmd = &cobra.Command{
		Use:   "get <name>",
		Short: "Get model details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.GetModelCommand{Name: args[0]})
		},
	}

	var testModelCmd = &cobra.Command{
		Use:   "test <name>",
		Short: "Test a model with an interactive prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.TestModelCommand{Name: args[0]})
		},
	}

	modelsCmd.AddCommand(listModelsCmd, createModelCmd, deleteModelCmd, getModelCmd, testModelCmd)

	// Traces command
	var tracesCmd = &cobra.Command{
		Use:   "traces",
		Short: "Manage traces",
	}

	var listTracesCmd = &cobra.Command{
		Use:   "list",
		Short: "List traces",
		RunE: func(cmd *cobra.Command, args []string) error {
			live, _ := cmd.Flags().GetBool("live")
			return executeCommand(&commands.ListTracesCommand{Live: live})
		},
	}
	listTracesCmd.Flags().Bool("live", false, "Watch for updates")

	var getTraceCmd = &cobra.Command{
		Use:   "get <trace-id>",
		Short: "Get trace details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			live, _ := cmd.Flags().GetBool("live")
			return executeCommand(&commands.GetTraceCommand{
				TraceID: args[0],
				Live:    live,
			})
		},
	}
	getTraceCmd.Flags().Bool("live", false, "Watch for updates")

	tracesCmd.AddCommand(listTracesCmd, getTraceCmd)

	// Functions command
	var functionsCmd = &cobra.Command{
		Use:   "functions",
		Short: "Manage functions",
	}

	var listFunctionsCmd = &cobra.Command{
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

	var createFunctionCmd = &cobra.Command{
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

	var deleteFunctionCmd = &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a function",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("yes")
			confirmed, err := commands.ConfirmDeletion("function", args[0], force)
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Println("Operation cancelled")
				return nil
			}
			return executeCommand(&commands.DeleteCommand{
				BaseCommand: commands.BaseCommand{
					FunctionPath: args[0],
				},
			})
		},
	}
	deleteFunctionCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	var getFunctionCmd = &cobra.Command{
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

	var chatFunctionCmd = &cobra.Command{
		Use:   "chat <name> [message]",
		Short: "Chat with a function",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var message string
			if len(args) > 1 {
				message = strings.Join(args[1:], " ")
			} else {
				// Read from stdin
				stdinData, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("error reading from stdin: %w", err)
				}
				message = string(stdinData)
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

	functionsCmd.AddCommand(listFunctionsCmd, createFunctionCmd, deleteFunctionCmd,
		getFunctionCmd, chatFunctionCmd)

	// Call command
	var callCmd = &cobra.Command{
		Use:   "call [flags] <name> <instructions> <input>",
		Short: "Call a function",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			model, _ := cmd.Flags().GetString("model")

			var input string
			if len(args) > 2 {
				input = args[2]
			} else {
				// Read from stdin
				stdinData, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("error reading from stdin: %w", err)
				}
				input = string(stdinData)
			}

			return executeCommand(&commands.CallCommand{
				Name:         args[0],
				Instructions: args[1],
				Input:        input,
				Model:        model,
			})
		},
	}
	callCmd.Flags().String("model", "", "Custom model to use")

	// Config command
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage API keys and configuration",
	}

	var listConfigCmd = &cobra.Command{
		Use:   "list",
		Short: "List configured API keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.ConfigCommand{Action: "list"})
		},
	}

	var addConfigCmd = &cobra.Command{
		Use:   "add <name> <api-key>",
		Short: "Add a new API key",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseUrl, _ := cmd.Flags().GetString("base-url")
			return executeCommand(&commands.ConfigCommand{
				Action:  "add",
				Name:    args[0],
				Key:     args[1],
				BaseUrl: baseUrl,
			})
		},
	}
	addConfigCmd.Flags().String("base-url", "", "Base URL for the API")

	var removeConfigCmd = &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove an API key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("yes")
			confirmed, err := commands.ConfirmDeletion("API key", args[0], force)
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Println("Operation cancelled")
				return nil
			}
			return executeCommand(&commands.ConfigCommand{
				Action: "remove",
				Name:   args[0],
			})
		},
	}
	removeConfigCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	configCmd.AddCommand(listConfigCmd, addConfigCmd, removeConfigCmd)

	// Add all top-level commands
	rootCmd.AddCommand(indexesCmd, modelsCmd, tracesCmd, functionsCmd, callCmd, configCmd)

	if err := rootCmd.Execute(); err != nil {
		cmd, _, findErr := rootCmd.Find(os.Args[1:])
		if findErr != nil {
			cmd = rootCmd
		}

		fmt.Fprintln(os.Stderr, "Error:", err)

		// Show usage for user errors
		if commands.IsUsageError(err) {
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, cmd.UsageString())
		}

		os.Exit(1)
	}
}
