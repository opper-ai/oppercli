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
	Options      *CallOptions
	Stream       bool
}

type CallOptions struct {
	Model string
}

func (c *CallCommand) Execute(ctx context.Context, client *opperai.Client) error {
	// Validate required fields
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.Instructions == "" {
		return fmt.Errorf("instructions are required")
	}
	if c.Input == "" {
		return fmt.Errorf("input is required")
	}

	var model string
	if c.Options != nil {
		model = c.Options.Model
	}

	// Always use streaming for better user experience
	response, err := client.Call.Call(ctx, c.Name, c.Instructions, c.Input, model, true)
	if err != nil {
		return err
	}

	if response == nil {
		return fmt.Errorf("received empty response from API")
	}

	// Stream the response to console
	for delta := range response.Stream {
		fmt.Print(delta)
	}
	fmt.Println() // Add newline at the end

	return nil
}
