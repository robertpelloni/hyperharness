package pi

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestEnhancedGrepTool tests the ripgrep/native fallback grep implementation.
func TestEnhancedGrepTool(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "grep_test.txt")
	os.WriteFile(testFile, []byte("hello world\nfoo bar\nhello foo\n"), 0o644)

	handlers := EnhancedToolHandlers()
	grepHandler, ok := handlers["grep"]
	if !ok {
		t.Fatal("enhanced grep handler not found")
	}

	input := GrepToolInput{
		Pattern: "hello",
		Path:    tmpDir,
		Limit:   10,
	}
	raw, _ := json.Marshal(input)

	result, err := grepHandler(context.Background(), tmpDir, raw)
	if err != nil {
		t.Fatalf("enhanced grep failed: %v", err)
	}

	if result == nil {
		t.Fatal("result is nil")
	}

	found := false
	for _, block := range result.Content {
		if tc, ok := block.(TextContent); ok {
			if tc.Text != "" {
				found = true
				t.Logf("grep result: %s", tc.Text)
			}
		}
	}
	if !found {
		t.Error("enhanced grep returned no content")
	}
}

// TestEnhancedFindTool tests the fd/native fallback find implementation.
func TestEnhancedFindTool(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package test\n"), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("hello\n"), 0o644)

	handlers := EnhancedToolHandlers()
	findHandler, ok := handlers["find"]
	if !ok {
		t.Fatal("enhanced find handler not found")
	}

	input := FindToolInput{
		Pattern: "*.go",
		Path:    tmpDir,
	}
	raw, _ := json.Marshal(input)

	result, err := findHandler(context.Background(), tmpDir, raw)
	if err != nil {
		t.Fatalf("enhanced find failed: %v", err)
	}

	if result == nil {
		t.Fatal("result is nil")
	}

	found := false
	for _, block := range result.Content {
		if tc, ok := block.(TextContent); ok {
			if tc.Text != "" {
				found = true
				if tc.Text == "No files found matching pattern" {
					t.Errorf("find returned no results in %s", tmpDir)
				} else {
					t.Logf("find result: %s", tc.Text)
				}
			}
		}
	}
	if !found {
		t.Error("enhanced find returned no content")
	}
}

// TestEnhancedLsTool tests the enhanced ls implementation.
func TestEnhancedLsTool(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("a\n"), 0o644)
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0o755)

	handlers := EnhancedToolHandlers()
	lsHandler, ok := handlers["ls"]
	if !ok {
		t.Fatal("enhanced ls handler not found")
	}

	input := LSToolInput{
		Path: tmpDir,
	}
	raw, _ := json.Marshal(input)

	result, err := lsHandler(context.Background(), tmpDir, raw)
	if err != nil {
		t.Fatalf("enhanced ls failed: %v", err)
	}

	found := false
	for _, block := range result.Content {
		if tc, ok := block.(TextContent); ok {
			if tc.Text != "" {
				found = true
				if !containsAll(tc.Text, "file1.txt", "subdir/") {
					t.Errorf("ls missing expected entries: %s", tc.Text)
				}
			}
		}
	}
	if !found {
		t.Error("enhanced ls returned no content")
	}
}

// TestResolveOrDefaultPath tests path resolution helper.
func TestResolveOrDefaultPath(t *testing.T) {
	tests := []struct {
		cwd, path, expected string
	}{
		{"/home/user", ".", "/home/user"},
		{"/home/user", "subdir", "/home/user/subdir"},
		{"/home/user", "/abs/path", "/abs/path"},
		{"/home/user", "../parent", "/home/parent"},
	}

	for _, tc := range tests {
		result := resolveOrDefaultPath(tc.cwd, tc.path)
		expected := filepath.Clean(tc.expected)
		if result != expected {
			t.Errorf("resolveOrDefaultPath(%q, %q) = %q, want %q", tc.cwd, tc.path, result, expected)
		}
	}
}

// TestLineScanner tests the line scanner utility.
func TestLineScanner(t *testing.T) {
	input := "line1\nline2\nline3\n"
	scanner := newLineScanner(stringReader(input))

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "line1" {
		t.Errorf("line 0: got %q, want %q", lines[0], "line1")
	}
	if lines[2] != "line3" {
		t.Errorf("line 2: got %q, want %q", lines[2], "line3")
	}
}

// TestMaxMinInt tests the helper functions.
func TestMaxMinInt(t *testing.T) {
	if maxInt(1, 2) != 2 {
		t.Error("maxInt(1,2) should be 2")
	}
	if maxInt(3, 1) != 3 {
		t.Error("maxInt(3,1) should be 3")
	}
	if minInt(1, 2) != 1 {
		t.Error("minInt(1,2) should be 1")
	}
	if minInt(3, 1) != 1 {
		t.Error("minInt(3,1) should be 1")
	}
}

// TestFormatFindResults tests the find results formatter.
func TestFormatFindResults(t *testing.T) {
	// Empty results
	result, err := formatFindResults(nil, false, 100)
	if err != nil {
		t.Fatalf("formatFindResults failed: %v", err)
	}
	found := false
	for _, block := range result.Content {
		if tc, ok := block.(TextContent); ok {
			if tc.Text == "No files found matching pattern" {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected 'No files found' for empty results")
	}

	// Non-empty results
	result, err = formatFindResults([]string{"file1.go", "file2.go"}, false, 100)
	if err != nil {
		t.Fatalf("formatFindResults failed: %v", err)
	}
	found = false
	for _, block := range result.Content {
		if tc, ok := block.(TextContent); ok {
			if tc.Text != "" && tc.Text != "No files found matching pattern" {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected non-empty results")
	}
}

// TestGrepNativeWithLiteral tests native grep with literal mode.
func TestGrepNativeWithLiteral(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "literal_test.txt")
	os.WriteFile(testFile, []byte("func TestFoo(t *testing.T) {\nvar regex = \"(a|b)*\"\n"), 0o644)

	input := GrepToolInput{
		Pattern: "(a|b)*",
		Path:    testFile,
		Literal: true,
		Limit:   10,
	}

	result, err := executeGrepNative(context.Background(), tmpDir, testFile, input, 10)
	if err != nil {
		t.Fatalf("native grep with literal failed: %v", err)
	}

	found := false
	for _, block := range result.Content {
		if tc, ok := block.(TextContent); ok {
			if tc.Text != "" && tc.Text != "No matches found" {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected literal grep to find regex pattern")
	}
}

// TestLsWithNativeEmptyDir tests ls with an empty directory.
func TestLsWithNativeEmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	result, err := executeLsWithNative(tmpDir, 500)
	if err != nil {
		t.Fatalf("ls empty dir failed: %v", err)
	}

	found := false
	for _, block := range result.Content {
		if tc, ok := block.(TextContent); ok {
			if tc.Text == "(empty directory)" {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected empty directory message")
	}
}

// Helper type for string-based io.Reader
type stringReaderImpl string

func (s stringReaderImpl) Read(p []byte) (n int, err error) {
	copy(p, s)
	return len(s), nil
}

func stringReader(s string) stringReaderImpl {
	return stringReaderImpl(s)
}

func containsAll(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if !containsStr(s, sub) {
			return false
		}
	}
	return true
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s[1:], sub) || s[:len(sub)] == sub)
}
