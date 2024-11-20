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
}

type CallOptions struct {
	Model string
}

func (c *CallCommand) Execute(ctx context.Context, client *opperai.Client) error {
	var model string
	if c.Options != nil {
		model = c.Options.Model
	}

	response, err := client.Call.Call(ctx, c.Name, c.Instructions, c.Input, model, false)
	if err != nil {
		return err
	}

	if response == nil {
		return fmt.Errorf("received empty response from API")
	}

	fmt.Println(response.Message)
	return nil
}
