package builders

import (
	"fmt"

	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/spf13/cobra"
)

func BuildUsageCommands(executeCommand func(commands.Command) error) *cobra.Command {
	usageCmd := &cobra.Command{
		Use:   "usage",
		Short: "Manage usage information",
		Example: `  # List usage information
  opper usage list

  # List usage with time range and granularity
  opper usage list --from-date=2024-01-01T00:00:00Z --to-date=2024-12-31T23:59:59Z --granularity=day

  # List usage with specific fields and grouping
  opper usage list --fields=completion_tokens,total_tokens --group-by=model,project.name

  # Show count over time as ASCII graph (default)
  opper usage list --graph

  # Show cost over time as ASCII graph
  opper usage list --graph=cost

  # Show count over time by model
  opper usage list --group-by model --graph

  # Export usage as CSV
  opper usage list --out csv`,
	}

	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List usage information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fromDate, _ := cmd.Flags().GetString("from-date")
			toDate, _ := cmd.Flags().GetString("to-date")
			granularity, _ := cmd.Flags().GetString("granularity")
			fields, _ := cmd.Flags().GetStringSlice("fields")
			groupBy, _ := cmd.Flags().GetStringSlice("group-by")
			out, _ := cmd.Flags().GetString("out")
			showGraph, _ := cmd.Flags().GetBool("graph")
			graphType, _ := cmd.Flags().GetString("graph-type")

			// If --graph is used, set the graph type
			var graph string
			if showGraph {
				graph = graphType
				if graph == "" {
					graph = "count"
				}
				if graph != "count" && graph != "cost" {
					return fmt.Errorf("invalid graph type: %s (must be 'count' or 'cost')", graph)
				}
			}

			return executeCommand(&commands.ListUsageCommand{
				FromDate:    fromDate,
				ToDate:      toDate,
				Granularity: granularity,
				Fields:      fields,
				GroupBy:     groupBy,
				Out:         out,
				Graph:       graph,
			})
		},
	}

	// Add flags
	listCmd.Flags().String("from-date", "", "Start date and time (RFC3339 format)")
	listCmd.Flags().String("to-date", "", "End date and time (RFC3339 format)")
	listCmd.Flags().String("granularity", "day", "Time granularity for grouping (minute, hour, day, month, year)")
	listCmd.Flags().StringSlice("fields", nil, "Fields from event_metadata to include and sum")
	listCmd.Flags().StringSlice("group-by", nil, "Fields from tags to group by")
	listCmd.Flags().String("out", "", "Output format (csv)")
	listCmd.Flags().Bool("graph", false, "Show graph")
	listCmd.Flags().String("graph-type", "count", "Graph type (count or cost)")

	usageCmd.AddCommand(listCmd)

	return usageCmd
}
