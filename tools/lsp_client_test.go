package tools

import (
	"testing"
)

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"main.go", ".go"},
		{"app.py", ".py"},
		{"index.ts", ".ts"},
		{"component.tsx", ".tsx"},
		{"/path/to/file.rs", ".rs"},
		{"C:\\Users\\test\\file.cpp", ".cpp"},
		{"noext", ""},
		{"Makefile", ""},
		{".gitignore", ".gitignore"},
	}

	for _, tc := range tests {
		got := getFileExtension(tc.path)
		if got != tc.expected {
			t.Errorf("getFileExtension(%q) = %q, want %q", tc.path, got, tc.expected)
		}
	}
}

func TestLSPServerConfigs(t *testing.T) {
	if len(lspServerConfigs) == 0 {
		t.Error("should have LSP server configs")
	}

	// Check key languages are configured
	expectedExts := []string{".go", ".py", ".rs", ".ts", ".tsx", ".js", ".cpp", ".c"}
	for _, ext := range expectedExts {
		config, ok := lspServerConfigs[ext]
		if !ok {
			t.Errorf("missing LSP config for %s", ext)
			continue
		}
		if config.Command == "" {
			t.Errorf("empty command for %s", ext)
		}
	}
}

func TestLSPOperationPlaceholder(t *testing.T) {
	tests := []struct {
		operation  string
		wantSubstr string
	}{
		{"goToDefinition", "goToDefinition"},
		{"findReferences", "findReferences"},
		{"hover", "hover"},
		{"documentSymbol", "documentSymbol"},
		{"workspaceSymbol", "workspaceSymbol"},
		{"goToImplementation", "goToImplementation"},
		{"diagnostics", "diagnostics"},
	}

	for _, tc := range tests {
		result, err := executeLSPOperation(nil, tc.operation, "/tmp/test.go", 5, 10, "MyFunc")
		if err != nil {
			t.Errorf("executeLSPOperation(%s) error: %v", tc.operation, err)
			continue
		}
		if !contains(result, tc.wantSubstr) {
			t.Errorf("executeLSPOperation(%s) = %q, should contain %q", tc.operation, result, tc.wantSubstr)
		}
	}
}

func TestLSPOperationUnknown(t *testing.T) {
	_, err := executeLSPOperation(nil, "unknownOp", "/tmp/test.go", 0, 0, "")
	if err == nil {
		t.Error("expected error for unknown operation")
	}
}

func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
