package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/opper-ai/opper-go/opperai"
)

func main() {
	// Check if at least the function name argument is provided.
	if len(os.Args) < 2 {
		fmt.Println("Error: Function name argument not provided.")
		os.Exit(1)
	}

	// The first argument after the program name is the function name.
	functionName := os.Args[1]

	// Read the API key from the environment variable.
	apiKey := os.Getenv("OPPER_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPPER_API_KEY environment variable not set.")
		os.Exit(1)
	}

	// Initialize the client with your API key.
	client := opperai.NewClient(apiKey)

	if functionName == "-c" {
		// Handle function creation
		if len(os.Args) < 3 {
			fmt.Println("Error: Function name for creation not provided.")
			os.Exit(1)
		}

		functionName := os.Args[2]
		instructions := strings.Join(os.Args[3:], " ")

		function := opperai.FunctionDescription{
			Path:         functionName,
			Instructions: instructions,
		}

		_, err := client.CreateFunction(context.Background(), function)
		if err != nil {
			fmt.Println("Error creating function:", err)
			return
		}

		fmt.Println("Function created successfully.")
		return
	}

	var messageContent string
	if len(os.Args) == 2 || os.Args[2] == "-" {
		// If the second argument is "-" or not provided, read from stdin.
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
			os.Exit(1)
		}
		messageContent = string(content)
	} else {
		// Combine all remaining arguments as the message content.
		messageContent = strings.Join(os.Args[2:], " ")
	}

	// Prepare the chat payload with the message content.
	chatPayload := opperai.ChatPayload{
		Messages: []opperai.Message{
			{
				Role:    "user",
				Content: messageContent,
			},
		},
	}

	// Create a context with a timeout to avoid hanging the request indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	chunks, err := client.Chat(ctx, functionName, chatPayload, true)
	if err != nil {
		fmt.Println("Error initiating chat function with streaming:", err)
		return
	}

	for chunk := range chunks {
		trimmedChunk := strings.TrimPrefix(string(chunk), "data: ")

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(trimmedChunk), &result); err != nil {
			fmt.Fprintln(os.Stderr, "Error unmarshalling chunk:", err)
			continue
		}
		if delta, ok := result["delta"].(string); ok {
			fmt.Print(delta)
		}
	}
	fmt.Println()
}
