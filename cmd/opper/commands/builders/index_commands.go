package builders

import (
	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/spf13/cobra"
)

func BuildIndexCommands(executeCommand func(commands.Command) error) *cobra.Command {
	indexesCmd := &cobra.Command{
		Use:   "indexes",
		Short: "Manage indexes",
		Example: `  # List all indexes
  opper indexes list

  # Create a new index
  opper indexes create myindex

  # Query an index
  opper indexes query myindex "search term"`,
	}

	// List command
	listCmd := &cobra.Command{
		Use:   "list [filter]",
		Short: "List indexes",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := ""
			if len(args) > 0 {
				filter = args[0]
			}
			format, _ := cmd.Flags().GetString("format")
			return executeCommand(&commands.ListIndexesCommand{
				Filter: filter,
				Format: format,
			})
		},
	}
	listCmd.Flags().String("format", "table", "Output format (table, plain)")

	// Create command
	createCmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.CreateIndexCommand{Name: args[0]})
		},
	}

	// Delete command
	deleteCmd := BuildDeletionCommand(
		"delete <name>",
		"Delete an index",
		"index",
		func(name string) error {
			return executeCommand(&commands.DeleteIndexCommand{Name: name})
		},
	)

	// Query command
	queryCmd := &cobra.Command{
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

	// Get command
	getCmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Get index details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCommand(&commands.GetIndexCommand{Name: args[0]})
		},
	}

	// Add command
	addCmd := &cobra.Command{
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

	// Upload command
	uploadCmd := &cobra.Command{
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

	indexesCmd.AddCommand(
		listCmd,
		createCmd,
		deleteCmd,
		queryCmd,
		getCmd,
		addCmd,
		uploadCmd,
	)

	return indexesCmd
}
