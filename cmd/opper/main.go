package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/opper-ai/oppercli/cmd/opper/commands"
	"github.com/opper-ai/oppercli/opperai"
)

var Version = "dev"

func main() {
	// Add version flag handling
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("opper version %s\n", Version)
		os.Exit(0)
	}

	// Read the API key from the environment variable
	apiKey := os.Getenv("OPPER_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPPER_API_KEY environment variable not set.")
		os.Exit(1)
	}

	// Initialize the client with base URL from environment or empty string
	baseURL := os.Getenv("OPPER_BASE_URL")
	client := opperai.NewClient(apiKey, baseURL)

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
