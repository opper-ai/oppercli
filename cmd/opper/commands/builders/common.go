package builders

import (
	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/spf13/cobra"
)

// AddDeletionFlags adds common deletion-related flags
func AddDeletionFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
}

// AddLiveFlags adds common live-update related flags
func AddLiveFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("live", false, "Watch for updates")
}

// BuildDeletionCommand creates a command with standard deletion behavior
func BuildDeletionCommand(use, short, resourceType string, executor func(string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("yes")
			confirmed, err := commands.ConfirmDeletion(resourceType, args[0], force)
			if err != nil || !confirmed {
				return err
			}
			return executor(args[0])
		},
	}
	AddDeletionFlags(cmd)
	return cmd
}
