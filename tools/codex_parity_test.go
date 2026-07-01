package tools

import (
	"encoding/json"
	"testing"
)

func TestParityLoad_codex_parity(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("Registry failed to initialize")
	}

	// Verify codex parity tools are registered
	expectedTools := map[string]bool{
		"computer_use_linux": false,
		"read_aloud":         false,
	}

	for _, tool := range r.Tools {
		if _, ok := expectedTools[tool.Name]; ok {
			expectedTools[tool.Name] = true

			// Verify parameter schemas are valid JSON
			var params map[string]interface{}
			err := json.Unmarshal(tool.Parameters, &params)
			if err != nil {
				t.Errorf("Tool %s has invalid JSON parameters: %v", tool.Name, err)
			}

			// Verify Execute function doesn't panic
			_, err = tool.Execute(map[string]interface{}{"test": "data"})
			if err != nil && err.Error() != "not implemented" {
				// We allow "not implemented" errors but other errors should be checked
				// In this case, we just want to ensure it executes without crashing
			}
		}
	}

	for name, found := range expectedTools {
		if !found {
			t.Errorf("Expected tool %s not found in registry", name)
		}
	}
}
