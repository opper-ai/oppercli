package commands

import (
	"context"
	"fmt"

	"github.com/opper-ai/oppercli/opperai"
)

type CallCommand struct {
	Name         string
	Instructions string
	Input        string
	Model        string
	Stream       bool
}

func (c *CallCommand) Execute(ctx context.Context, client *opperai.Client) error {
	if c.Input == "" {
		return NewUsageError(fmt.Errorf("input required (either as arguments or via stdin)"))
	}

	// Validate required fields
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.Instructions == "" {
		return fmt.Errorf("instructions are required")
	}

	// Try with non-streaming first
	response, err := client.Call.Call(ctx, c.Name, c.Instructions, c.Input, c.Model, false)
	if err != nil {
		return err // Return the error directly to preserve the error message
	}

	if response == nil {
		return fmt.Errorf("received empty response from API")
	}

	// If non-streaming succeeded, use streaming for output
	if c.Stream {
		streamResponse, err := client.Call.Call(ctx, c.Name, c.Instructions, c.Input, c.Model, true)
		if err != nil {
			return err
		}

		for delta := range streamResponse.Stream {
			fmt.Print(delta)
		}
		fmt.Println()
	} else {
		fmt.Println(response.Message)
	}

	return nil
}
