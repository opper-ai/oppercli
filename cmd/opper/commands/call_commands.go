package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/opper-ai/oppercli/opperai"
)

type CallCommand struct {
	BaseCommand
	Name         string
	Instructions string
	Input        string
	Model        string
	Stream       bool
}

func (c *CallCommand) Execute(ctx context.Context, client *opperai.Client) error {
	if client == nil {
		// Initialize client using environment variable
		apiKey := os.Getenv("OPPER_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("OPPER_API_KEY environment variable not set")
		}
		client = opperai.NewClient(apiKey)
	}

	if c.Stream {
		res, err := client.Call.Call(ctx, c.Name, c.Instructions, c.Input, c.Model, true)
		if err != nil {
			return err
		}

		for chunk := range res.Stream {
			fmt.Print(chunk)
		}
		return nil
	}

	result, err := client.Call.Call(ctx, c.Name, c.Instructions, c.Input, c.Model, false)
	if err != nil {
		return err
	}

	fmt.Println(result.Message)
	return nil
}
