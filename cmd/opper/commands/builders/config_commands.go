package builders

import (
	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/spf13/cobra"
)

func BuildConfigCommands(executeCommand func(commands.Command) error) *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage API keys and configuration",
		Example: `  # List configured API keys
  opper config list

  # Add a new API key
  opper config add mykey sk-xxx

  # Add a new API key with custom base URL
  opper config add mykey sk-xxx --base-url https://api.custom.com`,
	}

	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List configured API keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.ConfigCommand{Action: "list"})
		},
	}

	// Add command
	addCmd := &cobra.Command{
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
	addCmd.Flags().String("base-url", "", "Base URL for the API")

	// Remove command
	removeCmd := BuildDeletionCommand(
		"remove <name>",
		"Remove an API key",
		"API key",
		func(name string) error {
			return executeCommand(&commands.ConfigCommand{
				Action: "remove",
				Name:   name,
			})
		},
	)

	configCmd.AddCommand(
		listCmd,
		addCmd,
		removeCmd,
	)

	return configCmd
}
