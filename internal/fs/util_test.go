package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "a", "b", "c")

	if DirExists(newDir) {
		t.Error("dir should not exist yet")
	}

	if err := EnsureDir(newDir); err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}

	if !DirExists(newDir) {
		t.Error("dir should exist after EnsureDir")
	}
}

func TestEnsureParentDir(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "sub", "dir", "file.txt")

	if err := EnsureParentDir(filePath); err != nil {
		t.Fatalf("EnsureParentDir failed: %v", err)
	}

	if !DirExists(filepath.Join(tmpDir, "sub", "dir")) {
		t.Error("parent dir should exist")
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "exists.txt")
	os.WriteFile(existingFile, []byte("hello"), 0o644)

	if !FileExists(existingFile) {
		t.Error("existing file should be found")
	}
	if FileExists(filepath.Join(tmpDir, "nope.txt")) {
		t.Error("nonexistent file should not be found")
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"main.go", "go"},
		{"app.py", "python"},
		{"index.js", "javascript"},
		{"app.ts", "typescript"},
		{"main.rs", "rust"},
		{"Main.java", "java"},
		{"style.css", "css"},
		{"config.json", "json"},
		{"values.yaml", "yaml"},
		{"settings.toml", "toml"},
		{"Dockerfile", "dockerfile"},
		{"Makefile", "makefile"},
		{"go.mod", "go"},
		{"README.md", "markdown"},
		{"unknown.xyz", ""},
	}

	for _, tc := range tests {
		result := DetectLanguage(tc.filename)
		if result != tc.expected {
			t.Errorf("DetectLanguage(%q) = %q, want %q", tc.filename, result, tc.expected)
		}
	}
}

func TestIsGitIgnored(t *testing.T) {
	ignored := []string{
		".git/config",
		"node_modules/react/index.js",
		"__pycache__/module.pyc",
		"dist/bundle.js",
		"build/output.o",
		"target/debug/main",
		"vendor/lib.go",
	}
	for _, path := range ignored {
		if !IsGitIgnored(path) {
			t.Errorf("expected %q to be ignored", path)
		}
	}

	notIgnored := []string{
		"main.go",
		"src/index.ts",
		"README.md",
		".github/workflows/ci.yml",
		".agents/AGENTS.md",
	}
	for _, path := range notIgnored {
		if IsGitIgnored(path) {
			t.Errorf("expected %q to NOT be ignored", path)
		}
	}
}

func TestGetFileInfo(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	os.WriteFile(testFile, []byte("package main\n\nfunc main() {}\n"), 0o644)

	info, err := GetFileInfo(testFile)
	if err != nil {
		t.Fatalf("GetFileInfo failed: %v", err)
	}

	if info.IsDir {
		t.Error("should not be a directory")
	}
	if info.Language != "go" {
		t.Errorf("language: got %q, want 'go'", info.Language)
	}
	if info.LineCount != 4 {
		t.Errorf("line count: got %d, want 4", info.LineCount)
	}
	if info.Size == 0 {
		t.Error("size should not be 0")
	}
}

func TestWalkDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.go"), []byte("package test\n"), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("hello\n"), 0o644)
	os.Mkdir(filepath.Join(tmpDir, "sub"), 0o755)
	os.WriteFile(filepath.Join(tmpDir, "sub", "c.py"), []byte("print('hi')\n"), 0o644)

	tree, err := WalkDirectory(tmpDir, 0)
	if err != nil {
		t.Fatalf("WalkDirectory failed: %v", err)
	}

	if tree.FileCount != 3 {
		t.Errorf("file count: got %d, want 3", tree.FileCount)
	}
	if tree.DirCount != 1 {
		t.Errorf("dir count: got %d, want 1 (sub/), dirs: %+v", tree.DirCount, tree.Dirs)
	}
	if tree.TotalSize == 0 {
		t.Error("total size should not be 0")
	}
}

func TestFindProjectRoot(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "src", "pkg")
	os.MkdirAll(subDir, 0o755)
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n"), 0o644)

	root := FindProjectRoot(subDir)
	if root != tmpDir {
		t.Errorf("project root: got %q, want %q", root, tmpDir)
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		contains string
	}{
		{500, "B"},
		{1500, "KB"},
		{1500000, "MB"},
		{1500000000, "GB"},
	}
	for _, tc := range tests {
		result := FormatSize(tc.bytes)
		if !contains(result, tc.contains) {
			t.Errorf("FormatSize(%d) = %q, should contain %q", tc.bytes, result, tc.contains)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && (s[:len(sub)] == sub || contains(s[1:], sub)))
}

func TestPlatform(t *testing.T) {
	p := Platform()
	if p == "" {
		t.Error("Platform() should not be empty")
	}
}

func TestIsTextFile(t *testing.T) {
	if !IsTextFile("main.go") {
		t.Error("main.go should be text file")
	}
	if IsTextFile("image.png") {
		t.Error("image.png should not be text file")
	}
}
