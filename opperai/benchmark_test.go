package opperai

import (
	"context"
	"testing"
)

func BenchmarkIndexQuery(b *testing.B) {
	client := NewClient("test-key")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Indexes.Query("test-index", "test query", []Filter{})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkChat(b *testing.B) {
	client := NewClient("test-key")
	ctx := context.Background()
	payload := ChatPayload{
		Messages: []Message{{Role: "user", Content: "Hello"}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Chat(ctx, "test/function", payload, false)
		if err != nil {
			b.Fatal(err)
		}
	}
}
