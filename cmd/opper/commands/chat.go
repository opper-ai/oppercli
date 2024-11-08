package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/opper-ai/oppercli/opperai"
)

// ChatCommand handles chat interactions
type ChatCommand struct {
	BaseCommand
	MessageContent string
}

func (c *ChatCommand) Execute(ctx context.Context, client *opperai.Client) error {
	var messageContent string

	// If no message content provided, try reading from stdin
	if c.MessageContent == "" {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		messageContent = string(content)
	} else {
		messageContent = c.MessageContent
	}

	chatPayload := opperai.ChatPayload{
		Messages: []opperai.Message{
			{
				Role:    "user",
				Content: messageContent,
			},
		},
	}

	chunks, err := client.Chat(ctx, c.FunctionPath, chatPayload, true)
	if err != nil {
		return fmt.Errorf("error initiating chat: %w", err)
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

	return nil
}
