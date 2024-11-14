package opperai

import (
	"context"
	"testing"
)

// MockClient implements the Client interface for testing
type MockClient struct {
	ListModelsFunc  func(context.Context) ([]CustomLanguageModel, error)
	CreateModelFunc func(context.Context, CustomLanguageModel) error
	// Add other methods
}

func (m *MockClient) ListModels(ctx context.Context) ([]CustomLanguageModel, error) {
	if m.ListModelsFunc != nil {
		return m.ListModelsFunc(ctx)
	}
	return nil, nil
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
