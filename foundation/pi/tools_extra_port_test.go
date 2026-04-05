package pi

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func mustExec(t *testing.T, runtime *Runtime, tool string, input any) {
	t.Helper()
	_, err := mustExecResultErr(runtime, tool, input)
	if err != nil {
		t.Fatalf("%s failed: %v", tool, err)
	}
}

func mustExecResult(t *testing.T, runtime *Runtime, tool string, input any) *ToolResult {
	t.Helper()
	result, err := mustExecResultErr(runtime, tool, input)
	if err != nil {
		t.Fatalf("%s failed: %v", tool, err)
	}
	return result
}

func mustExecResultErr(runtime *Runtime, tool string, input any) (*ToolResult, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	return runtime.ExecuteTool(context.Background(), "", tool, raw, nil)
}

func TestGrepToolNoMatches(t *testing.T) {
	dir := t.TempDir()
	runtime := NewRuntime(dir, DefaultSessionStore(dir))

	mustExec(t, runtime, "write", WriteToolInput{Path: "c.txt", Content: "nothing here\n"})
	result := mustExecResult(t, runtime, "grep", GrepToolInput{Pattern: "zxcvzxcv"})
	text := textFromResult(t, result)
	if text != "No matches found" {
		t.Fatalf("expected no matches message, got %q", text)
	}
}

func TestFindToolNoFiles(t *testing.T) {
	dir := t.TempDir()
	runtime := NewRuntime(dir, DefaultSessionStore(dir))

	if err := os.WriteFile(filepath.Join(dir, "hello.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	result := mustExecResult(t, runtime, "find", FindToolInput{Pattern: "*.go"})
	text := textFromResult(t, result)
	if text != "No files found matching pattern" {
		t.Fatalf("expected no files message, got %q", text)
	}
}

func TestLsToolEmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	runtime := NewRuntime(dir, DefaultSessionStore(dir))

	subdir := filepath.Join(dir, "empty")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}

	result := mustExecResult(t, runtime, "ls", LSToolInput{Path: "empty"})
	text := textFromResult(t, result)
	if text != "(empty directory)" {
		t.Fatalf("expected empty directory message, got %q", text)
	}
}

func TestGrepToolLimitDetails(t *testing.T) {
	dir := t.TempDir()
	runtime := NewRuntime(dir, DefaultSessionStore(dir))

	content := strings.Repeat("line match\n", 10)
	mustExec(t, runtime, "write", WriteToolInput{Path: "many.txt", Content: content})

	result := mustExecResult(t, runtime, "grep", GrepToolInput{Pattern: "match", Limit: 3})
	details, ok := result.Details.(*GrepToolDetails)
	if !ok || details == nil {
		t.Fatalf("expected grep details with limit info, got %#v", result.Details)
	}
	if details.MatchLimitReached != 3 {
		t.Fatalf("expected match limit 3, got %v", details.MatchLimitReached)
	}
}

func TestFindToolLimitDetails(t *testing.T) {
	dir := t.TempDir()
	runtime := NewRuntime(dir, DefaultSessionStore(dir))

	for i := 0; i < 10; i++ {
		fname := "f" + string(rune('0'+i)) + ".txt"
		if err := os.WriteFile(filepath.Join(dir, fname), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	result := mustExecResult(t, runtime, "find", FindToolInput{Pattern: "*.txt", Limit: 3})
	details, ok := result.Details.(*FindToolDetails)
	if !ok || details == nil {
		t.Fatalf("expected find details with limit info, got %#v", result.Details)
	}
	if details.ResultLimitReached != 3 {
		t.Fatalf("expected result limit 3, got %v", details.ResultLimitReached)
	}
}

func TestLsToolLimitDetails(t *testing.T) {
	dir := t.TempDir()
	runtime := NewRuntime(dir, DefaultSessionStore(dir))

	for i := 0; i < 10; i++ {
		fname := "file" + string(rune('0'+i)) + ".txt"
		if err := os.WriteFile(filepath.Join(dir, fname), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	result := mustExecResult(t, runtime, "ls", LSToolInput{Limit: 3})
	details, ok := result.Details.(*LsToolDetails)
	if !ok || details == nil {
		t.Fatalf("expected ls details with limit info, got %#v", result.Details)
	}
	if details.EntryLimitReached != 3 {
		t.Fatalf("expected entry limit 3, got %v", details.EntryLimitReached)
	}
}

func TestFindToolMatchesRelativeSubdirPattern(t *testing.T) {
	dir := t.TempDir()
	runtime := NewRuntime(dir, DefaultSessionStore(dir))

	if err := os.MkdirAll(filepath.Join(dir, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "src", "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatal(err)
	}

	result := mustExecResult(t, runtime, "find", FindToolInput{Pattern: "src/*.go"})
	text := textFromResult(t, result)
	if !strings.Contains(text, "src/main.go") {
		t.Fatalf("expected relative subdir match, got %q", text)
	}
}
