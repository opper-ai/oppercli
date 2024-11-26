package builders

import (
	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/spf13/cobra"
)

func BuildTraceCommands(executeCommand func(commands.Command) error) *cobra.Command {
	tracesCmd := &cobra.Command{
		Use:   "traces",
		Short: "Manage traces",
		Example: `  # List all traces
  opper traces list

  # Get trace details
  opper traces get <trace-id>

  # Watch traces in real-time
  opper traces list --live`,
	}

	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List traces",
		RunE: func(cmd *cobra.Command, args []string) error {
			live, _ := cmd.Flags().GetBool("live")
			return executeCommand(&commands.ListTracesCommand{Live: live})
		},
	}
	AddLiveFlags(listCmd)

	// Get command
	getCmd := &cobra.Command{
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
	AddLiveFlags(getCmd)

	tracesCmd.AddCommand(
		listCmd,
		getCmd,
	)

	return tracesCmd
}
