package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsePatchAdd(t *testing.T) {
	patch := `*** Begin Patch
*** Add File: test.txt
hello world
line 2
*** End Patch`
	hunks, err := ParsePatch(patch)
	if err != nil {
		t.Fatal(err)
	}
	if len(hunks) != 1 {
		t.Fatalf("expected 1 hunk, got %d", len(hunks))
	}
	if hunks[0].Type != "add" {
		t.Errorf("type: %s", hunks[0].Type)
	}
	if hunks[0].Path != "test.txt" {
		t.Errorf("path: %s", hunks[0].Path)
	}
}

func TestParsePatchUpdate(t *testing.T) {
	patch := `*** Begin Patch
*** Update File: main.go
func old() {
-func old() {
+func new() {
*** End Patch`
	hunks, err := ParsePatch(patch)
	if err != nil {
		t.Fatal(err)
	}
	if len(hunks) != 1 {
		t.Fatalf("expected 1 hunk, got %d", len(hunks))
	}
	if hunks[0].Type != "update" {
		t.Errorf("type: %s", hunks[0].Type)
	}
}

func TestParsePatchDelete(t *testing.T) {
	patch := `*** Begin Patch
*** Delete File: old.txt
*** End Patch`
	hunks, err := ParsePatch(patch)
	if err != nil {
		t.Fatal(err)
	}
	if len(hunks) != 1 {
		t.Fatalf("expected 1 hunk, got %d", len(hunks))
	}
	if hunks[0].Type != "delete" {
		t.Errorf("type: %s", hunks[0].Type)
	}
}

func TestParsePatchMultiple(t *testing.T) {
	patch := `*** Begin Patch
*** Add File: new.txt
content
*** Update File: old.txt
-old line
+new line
*** Delete File: gone.txt
*** End Patch`
	hunks, err := ParsePatch(patch)
	if err != nil {
		t.Fatal(err)
	}
	if len(hunks) != 3 {
		t.Fatalf("expected 3 hunks, got %d", len(hunks))
	}
}

func TestApplyPatchAdd(t *testing.T) {
	dir := t.TempDir()
	cwd := dir

	hunks := []PatchHunk{{
		Path:     "newfile.txt",
		Type:     "add",
		Contents: "hello\nworld",
	}}

	result := ApplyPatch(hunks, cwd)
	if len(result.Files) != 1 {
		t.Fatalf("expected 1 file result, got %d", len(result.Files))
	}
	if result.Files[0].Err != nil {
		t.Fatal(result.Files[0].Err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "newfile.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello\nworld\n" {
		t.Errorf("content: %q", string(data))
	}
}

func TestApplyPatchUpdate(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.txt")
	os.WriteFile(filePath, []byte("line1\nline2\nline3\n"), 0o644)

	hunks := []PatchHunk{{
		Path: filePath,
		Type: "update",
		Chunks: []PatchChunk{
			{Type: "context", Content: "line1"},
			{Type: "remove", Content: "line2"},
			{Type: "add", Content: "LINE_TWO"},
			{Type: "context", Content: "line3"},
		},
	}}

	result := ApplyPatch(hunks, dir)
	if result.Files[0].Err != nil {
		t.Fatal(result.Files[0].Err)
	}

	data, _ := os.ReadFile(filePath)
	if string(data) != "line1\nLINE_TWO\nline3\n" {
		t.Errorf("content: %q", string(data))
	}
}

func TestApplyPatchDelete(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.txt")
	os.WriteFile(filePath, []byte("content\n"), 0o644)

	hunks := []PatchHunk{{
		Path: filePath,
		Type: "delete",
	}}

	result := ApplyPatch(hunks, dir)
	if result.Files[0].Err != nil {
		t.Fatal(result.Files[0].Err)
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}

func TestMultiEdit(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.txt")
	os.WriteFile(filePath, []byte("hello world\nfoo bar\nbaz"), 0o644)

	result, err := ApplyMultiEdit(MultiEditParams{
		FilePath: filePath,
		Edits: []MultiEditItem{
			{OldString: "hello", NewString: "HELLO"},
			{OldString: "bar", NewString: "BAR"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Additions != 2 || result.Deletions != 2 {
		t.Errorf("additions=%d deletions=%d", result.Additions, result.Deletions)
	}

	data, _ := os.ReadFile(filePath)
	if string(data) != "HELLO world\nfoo BAR\nbaz" {
		t.Errorf("content: %q", string(data))
	}
}

func TestMultiEditReplaceAll(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.txt")
	os.WriteFile(filePath, []byte("aaa bbb aaa\nccc aaa"), 0o644)

	result, err := ApplyMultiEdit(MultiEditParams{
		FilePath: filePath,
		Edits: []MultiEditItem{
			{OldString: "aaa", NewString: "XXX", ReplaceAll: true},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Additions != 3 {
		t.Errorf("additions=%d", result.Additions)
	}

	data, _ := os.ReadFile(filePath)
	if string(data) != "XXX bbb XXX\nccc XXX" {
		t.Errorf("content: %q", string(data))
	}
}

func TestMultiEditNotFound(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.txt")
	os.WriteFile(filePath, []byte("hello"), 0o644)

	_, err := ApplyMultiEdit(MultiEditParams{
		FilePath: filePath,
		Edits: []MultiEditItem{
			{OldString: "nonexistent", NewString: "xxx"},
		},
	})
	if err == nil {
		t.Error("should error for not found oldString")
	}
}
