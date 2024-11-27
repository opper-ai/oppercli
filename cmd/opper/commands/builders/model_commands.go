package builders

import (
	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/spf13/cobra"
)

func BuildModelCommands(executeCommand func(commands.Command) error) *cobra.Command {
	modelsCmd := &cobra.Command{
		Use:   "models",
		Short: "Manage models",
		Example: `  # List all models
  opper models list

  # Create a new model
  opper models create mymodel litellm-id api-key

  # Test a model
  opper models test mymodel`,
	}

	// List command
	listCmd := &cobra.Command{
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

	// Create command
	createCmd := &cobra.Command{
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

	// Delete command
	deleteCmd := BuildDeletionCommand(
		"delete <name>",
		"Delete a model",
		"model",
		func(name string) error {
			return executeCommand(&commands.DeleteModelCommand{Name: name})
		},
	)

	// Get command
	getCmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Get model details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.GetModelCommand{Name: args[0]})
		},
	}

	// Test command
	testCmd := &cobra.Command{
		Use:   "test <name>",
		Short: "Test a model with an interactive prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.TestModelCommand{Name: args[0]})
		},
	}

	// Builtin command
	builtinCmd := &cobra.Command{
		Use:   "builtin",
		Short: "List built-in models",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.ListBuiltinModelsCommand{})
		},
	}

	modelsCmd.AddCommand(
		listCmd,
		createCmd,
		deleteCmd,
		getCmd,
		testCmd,
		builtinCmd,
	)

	return modelsCmd
}
