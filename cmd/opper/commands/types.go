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

// Index Commands
type ListIndexesCommand struct {
	Filter string
	Format string
}

type CreateIndexCommand struct {
	Name string
}

type DeleteIndexCommand struct {
	Name string
}

type GetIndexCommand struct {
	Name string
}

type QueryIndexCommand struct {
	Name   string
	Query  string
	Filter string
}

type AddToIndexCommand struct {
	Name     string
	Key      string
	Content  string
	Metadata string
}

type UploadToIndexCommand struct {
	Name     string
	FilePath string
}

// Model Commands
type ListModelsCommand struct {
	Filter string
}

type CreateModelCommand struct {
	Name       string
	Identifier string
	APIKey     string
	Extra      string
}

type DeleteModelCommand struct {
	Name string
}

type GetModelCommand struct {
	Name string
}

type TestModelCommand struct {
	Name string
}

type ListBuiltinModelsCommand struct {
	Filter string
}

// Function Commands
type ListCommand struct {
	Filter string
}

type CreateCommand struct {
	BaseCommand
	Instructions string
}

type DeleteCommand struct {
	BaseCommand
}

type GetCommand struct {
	BaseCommand
}

type FunctionChatCommand struct {
	BaseCommand
	Message string
}

type ListEvaluationsCommand struct {
	BaseCommand
}

type RunEvaluationCommand struct {
	BaseCommand
}

// Call Commands
type CallCommand struct {
	Name         string
	Instructions string
	Input        string
	Model        string
	Stream       bool
}

// Config Commands
type ConfigCommand struct {
	Action  string
	Name    string
	Key     string
	BaseUrl string
}
