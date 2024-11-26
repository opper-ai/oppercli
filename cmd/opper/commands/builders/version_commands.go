package builders

import (
	"fmt"

	"github.com/spf13/cobra"
)

func BuildVersionCommand(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("opper version %s\n", version)
		},
	}
}
