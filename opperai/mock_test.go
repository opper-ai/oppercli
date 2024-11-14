package opperai

import (
	"context"
	"testing"
)

// MockClient implements the Client interface for testing
type MockClient struct {
	ListModelsFunc  func(context.Context) ([]CustomLanguageModel, error)
	CreateModelFunc func(context.Context, CustomLanguageModel) error
	ChatFunc        func(context.Context, string, string) (string, error)
	// Add other methods
}

func (m *MockClient) ListModels(ctx context.Context) ([]CustomLanguageModel, error) {
	if m.ListModelsFunc != nil {
		return m.ListModelsFunc(ctx)
	}
	return nil, nil
}

// Add the Chat method to the mock
func (m *MockClient) Chat(ctx context.Context, functionName string, message string) (string, error) {
	if m.ChatFunc != nil {
		return m.ChatFunc(ctx, functionName, message)
	}
	return "", nil
}

// Example test using mock
func TestWithMock(t *testing.T) {
	mock := &MockClient{
		ListModelsFunc: func(ctx context.Context) ([]CustomLanguageModel, error) {
			return []CustomLanguageModel{{Name: "test"}}, nil
		},
	}

	models, err := mock.ListModels(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(models) != 1 {
		t.Errorf("expected 1 model, got %d", len(models))
	}
	if models[0].Name != "test" {
		t.Errorf("expected model name 'test', got %s", models[0].Name)
	}
}
