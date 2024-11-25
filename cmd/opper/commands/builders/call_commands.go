package builders

import (
	"fmt"
	"io"
	"os"

	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/spf13/cobra"
)

func BuildCallCommand(executeCommand func(commands.Command) error) *cobra.Command {
	callCmd := &cobra.Command{
		Use:   "call [flags] <name> <instructions> <input>",
		Short: "Call a function",
		Example: `  # Call with direct input
  opper call myfunction "respond about X" "what is X?"

  # Call with input from stdin
  echo "what is X?" | opper call myfunction "respond about X"`,
		Args: cobra.MinimumNArgs(2),
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

	return callCmd
}
