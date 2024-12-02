package builders

import (
	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/spf13/cobra"
)

func BuildUsageCommands(executeCommand func(commands.Command) error) *cobra.Command {
	usageCmd := &cobra.Command{
		Use:   "usage",
		Short: "Manage usage information",
		Example: `  # List usage information
  opper usage list

  # List usage with filters
  opper usage list --start-date=2024-01-01 --end-date=2024-12-31 --project-name=myproject

  # Export usage as CSV
  opper usage list --out csv`,
	}

	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List usage information",
		RunE: func(cmd *cobra.Command, args []string) error {
			startDate, _ := cmd.Flags().GetString("start-date")
			endDate, _ := cmd.Flags().GetString("end-date")
			projectName, _ := cmd.Flags().GetString("project-name")
			functionName, _ := cmd.Flags().GetString("function-name")
			model, _ := cmd.Flags().GetString("model")
			skip, _ := cmd.Flags().GetInt("skip")
			limit, _ := cmd.Flags().GetInt("limit")
			out, _ := cmd.Flags().GetString("out")

			return executeCommand(&commands.ListUsageCommand{
				StartDate:    startDate,
				EndDate:      endDate,
				ProjectName:  projectName,
				FunctionName: functionName,
				Model:        model,
				Skip:         skip,
				Limit:        limit,
				Out:          out,
			})
		},
	}

	// Add flags
	listCmd.Flags().String("start-date", "", "Filter by start date (YYYY-MM-DD)")
	listCmd.Flags().String("end-date", "", "Filter by end date (YYYY-MM-DD)")
	listCmd.Flags().String("project-name", "", "Filter by project name")
	listCmd.Flags().String("function-name", "", "Filter by function name")
	listCmd.Flags().String("model", "", "Filter by model")
	listCmd.Flags().Int("skip", 0, "Number of items to skip")
	listCmd.Flags().Int("limit", 0, "Maximum number of items to return")
	listCmd.Flags().String("out", "", "Output format (csv)")

	usageCmd.AddCommand(listCmd)

	return usageCmd
}
