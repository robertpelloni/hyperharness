package tools

import (
	"testing"
)

func TestWarpParityTools(t *testing.T) {
	r := NewRegistry()

	t.Run("RequestCommandOutput", func(t *testing.T) {
		tool, ok := r.Find("RequestCommandOutput")
		if !ok {
			t.Fatal("RequestCommandOutput tool not found")
		}

		args := map[string]interface{}{
			"command": "ls",
		}
		result, err := tool.Execute(args)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		if result == "" {
			t.Error("Expected non-empty result")
		}
	})

	t.Run("SearchCodebase", func(t *testing.T) {
		tool, ok := r.Find("SearchCodebase")
		if !ok {
			t.Fatal("SearchCodebase tool not found")
		}

		args := map[string]interface{}{
			"query": "func",
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
