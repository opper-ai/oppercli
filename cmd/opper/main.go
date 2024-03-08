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

type Args struct {
	FunctionName   string
	MessageContent string
	IsCreation     bool
	IsDeletion     bool
	IsListing      bool
	IsGetByPath    bool
	ListFilter     string
}

func parseArgs() Args {
	// Check if at least the function name argument is provided.
	if len(os.Args) < 2 {
		fmt.Println("Error: Function name argument not provided.")
		os.Exit(1)
	}

	args := Args{}

	// The first argument after the program name is the function name or a flag.
	firstArg := os.Args[1]

	if firstArg == "-c" {
		args.IsCreation = true
		if len(os.Args) < 3 {
			fmt.Println("Error: Function name for creation not provided.")
			os.Exit(1)
		}
		args.FunctionName = os.Args[2]
		args.MessageContent = strings.Join(os.Args[3:], " ")
	} else if firstArg == "-d" {
		args.IsDeletion = true
		if len(os.Args) < 3 {
			fmt.Println("Error: Function name for deletion not provided.")
			os.Exit(1)
		}
		args.FunctionName = os.Args[2]
	} else if firstArg == "-l" {
		args.IsListing = true
		if len(os.Args) > 2 {
			args.ListFilter = os.Args[2]
		}
	} else if firstArg == "-g" {
		args.IsGetByPath = true
		if len(os.Args) < 3 {
			fmt.Println("Error: Function path for retrieval not provided.")
			os.Exit(1)
		}
		args.FunctionName = os.Args[2]
	} else {
		args.FunctionName = firstArg
		messageContentParts := []string{}

		if len(os.Args) > 2 {
			if os.Args[2] == "-" {
				messageContentParts = append(messageContentParts, strings.Join(os.Args[3:], " "))
			} else {
				messageContentParts = append(messageContentParts, strings.Join(os.Args[2:], " "))
			}
		}

		if len(os.Args) == 2 || os.Args[2] == "-" {
			content, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
				os.Exit(1)
			}
			if len(messageContentParts) > 0 {
				messageContentParts = append(messageContentParts, "\n")
			}
			messageContentParts = append(messageContentParts, string(content))
		}

		args.MessageContent = strings.Join(messageContentParts, "")
	}

	return args
}

func main() {
	args := parseArgs()

	// Read the API key from the environment variable.
	apiKey := os.Getenv("OPPER_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPPER_API_KEY environment variable not set.")
		os.Exit(1)
	}

	// Initialize the client with your API key.
	client := opperai.NewClient(apiKey)

	if args.IsCreation {
		function := opperai.FunctionDescription{
			Path:         args.FunctionName,
			Instructions: args.MessageContent,
		}

		_, err := client.CreateFunction(context.Background(), function)
		if err != nil {
			fmt.Println("Error creating function:", err)
			return
		}

		fmt.Println("Function created successfully.")
		return
	} else if args.IsDeletion {
		err := client.DeleteFunction(context.Background(), "", args.FunctionName)
		if err != nil {
			fmt.Println("Error deleting function:", err)
			return
		}
		fmt.Println("Function deleted successfully.")
		return
	} else if args.IsListing {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		functions, err := client.ListFunctions(ctx)
		if err != nil {
			fmt.Println("Error listing functions:", err)
			return
		}

		for _, function := range functions {
			if args.ListFilter == "" || strings.Contains(function.Path, args.ListFilter) {
				fmt.Printf("Path: %s, Description: %s\n",
					function.Path, function.Description)
			}
		}
		return
	} else if args.IsGetByPath {
		function, err := client.GetFunctionByPath(context.Background(), args.FunctionName)
		if err != nil {
			fmt.Println("Error retrieving function by path:", err)
			return
		}
		if function == nil {
			fmt.Println("Function not found.")
			return
		}
		fmt.Printf("Path: %s\nDescription: %s\nInstructions: %s\n\n", function.Path, function.Description, function.Instructions)
		return
	} else {
		// Prepare the chat payload with the message content.
		chatPayload := opperai.ChatPayload{
			Messages: []opperai.Message{
				{
					Role:    "user",
					Content: args.MessageContent,
				},
			},
		}

		// Create a context with a timeout to avoid hanging the request indefinitely.
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		chunks, err := client.Chat(ctx, args.FunctionName, chatPayload, true)
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
}
