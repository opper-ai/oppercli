package commands

import (
	"testing"

	"github.com/opper-ai/oppercli/opperai"
)

type MockClient struct {
	*opperai.Client
	// Add mock implementations
}

func NewTestClient(t *testing.T) *MockClient {
	return &MockClient{
		Client: &opperai.Client{}, // Return a basic mock client
	}
}
