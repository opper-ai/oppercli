package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/opper-ai/oppercli/opperai"
)

func main() {
	// Read the API key from the environment variable
	apiKey := os.Getenv("OPPER_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPPER_API_KEY environment variable not set.")
		os.Exit(1)
	}

	// Initialize the client
	client := opperai.NewClient(apiKey)

	// Create command parser
	parser := commands.NewCommandParser()

	// Parse command
	cmd, err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println("Error parsing command:", err)
		os.Exit(1)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Execute command
	if err := cmd.Execute(ctx, client); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
