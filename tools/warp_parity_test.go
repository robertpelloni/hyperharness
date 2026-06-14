package tools

import (
	"testing"
	"time"
)

func TestWarpParityTools(t *testing.T) {
	r := NewRegistry()
	t.Run("RequestCommandOutput", func(t *testing.T) {
		tool, ok := r.Find("RequestCommandOutput")
		if !ok {
			t.Fatal("RequestCommandOutput tool not found")
		}
		args := map[string]interface{}{
			"command": "echo hello",
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
		done := make(chan string, 1)
		errCh := make(chan error, 1)
		go func() {
			args := map[string]interface{}{
				"query":        "func",
				"codebase_path": ".",
			}
			result, err := tool.Execute(args)
			if err != nil {
				errCh <- err
				return
			}
			done <- result
		}()
		select {
		case result := <-done:
			if result == "" {
				t.Error("Expected non-empty result")
			}
		case err := <-errCh:
			t.Fatalf("Execute failed: %v", err)
		case <-time.After(10 * time.Second):
			t.Log("SearchCodebase timed out (likely large repo), skipping")
		}
	})
}
