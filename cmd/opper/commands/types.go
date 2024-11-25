package commands

import (
	"context"

	"github.com/opper-ai/oppercli/opperai"
)

// Command interface that all commands must implement
type Command interface {
	Execute(ctx context.Context, client *opperai.Client) error
}

// BaseCommand contains common fields
type BaseCommand struct {
	FunctionPath string
}
