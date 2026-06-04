package tools

import (
	"testing"
)

func TestTabbyParityTools(t *testing.T) {
	r := NewRegistry()

	t.Run("tabby_completion", func(t *testing.T) {
		tool, ok := r.Find("tabby_completion")
		if !ok {
			t.Fatal("tabby_completion tool not found")
		}

		args := map[string]interface{}{
			"segments": map[string]interface{}{
				"prefix": "func main() {",
			},
		}
		result, err := tool.Execute(args)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		if result == "" {
			t.Error("Expected non-empty result")
		}
	})

	t.Run("tabby_chat", func(t *testing.T) {
		tool, ok := r.Find("tabby_chat")
		if !ok {
			t.Fatal("tabby_chat tool not found")
		}

		args := map[string]interface{}{
			"messages": []interface{}{
				map[string]interface{}{"role": "user", "content": "Hello"},
			},
		}
		result, err := tool.Execute(args)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		if result == "" {
			t.Error("Expected non-empty result")
		}
	})
}
