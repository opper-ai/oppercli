package main

import (
	"context"
	"os"

	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/opper-ai/oppercli/cmd/opper/commands/builders"
	"github.com/opper-ai/oppercli/cmd/opper/config"
	"github.com/opper-ai/oppercli/opperai"
	"github.com/spf13/cobra"
)

var version = "dev"

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
	rootCmd := &cobra.Command{
		Use:   "opper",
		Short: "Opper CLI - interact with Opper AI services",
	}

	// Global flags
	var keyName string
	rootCmd.PersistentFlags().StringVar(&keyName, "key", "default", "API key to use from config")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug output")

	// Create executor function
	executeCommand := func(cmd commands.Command) error {
		ctx := context.Background()
		apiKey, baseUrl, err := config.GetAPIKeyAndBaseUrl(keyName)
		if err != nil {
			return builders.FormatError(err)
		}
		client := opperai.NewClient(apiKey, baseUrl)
		if err := cmd.Execute(ctx, client); err != nil {
			return builders.FormatError(err)
		}
		return nil
	}

	// Add command groups using builders
	rootCmd.AddCommand(
		builders.BuildIndexCommands(executeCommand),
		builders.BuildModelCommands(executeCommand),
		builders.BuildTraceCommands(executeCommand),
		builders.BuildFunctionCommands(executeCommand),
		builders.BuildConfigCommands(executeCommand),
		builders.BuildVersionCommand(version),
		builders.BuildCallCommand(executeCommand),
		builders.BuildUsageCommands(executeCommand),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
