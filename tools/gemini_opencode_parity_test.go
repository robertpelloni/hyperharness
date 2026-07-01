package tools

import (
	"testing"
)

func TestParityLoad_gemini_opencode_parity(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("Registry failed to initialize")
	}

	for _, tool := range r.Tools {
		if tool.Name == "" {
			t.Errorf("Tool without a name found")
		}
		if tool.Execute == nil {
			t.Errorf("Tool %s has no Execute function", tool.Name)
		}
	}
}
