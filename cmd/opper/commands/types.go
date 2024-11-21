package commands

import (
	"context"

	"github.com/opper-ai/oppercli/opperai"
)

// Command interface defines what all commands must implement
type Command interface {
	Execute(ctx context.Context, client *opperai.Client) error
}

// BaseCommand struct for shared fields
type BaseCommand struct {
	FunctionPath string
}
