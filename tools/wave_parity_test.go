package tools

import (
	"testing"
)

func TestWaveParityTools(t *testing.T) {
	r := NewRegistry()

	t.Run("TermGetScrollback", func(t *testing.T) {
		tool, ok := r.Find("TermGetScrollback")
		if !ok {
			t.Fatal("TermGetScrollback tool not found")
		}

		args := map[string]interface{}{
			"blockId": "test-block",
		}
		result, err := tool.Execute(args)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		if result == "" {
			t.Error("Expected non-empty result")
		}
	})

	t.Run("CaptureScreenshot", func(t *testing.T) {
		tool, ok := r.Find("CaptureScreenshot")
		if !ok {
			t.Fatal("CaptureScreenshot tool not found")
		}

		args := map[string]interface{}{}
		result, err := tool.Execute(args)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		if result == "" {
			t.Error("Expected non-empty result")
		}
	})
}
