package repograph

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildGraph(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "pkg"), 0755)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\nimport \"fmt\"\nfunc main() { fmt.Println(\"hello\") }\n"), 0644)
	os.WriteFile(filepath.Join(dir, "pkg", "util.go"), []byte("package pkg\ntype Widget struct { Name string }\nfunc NewWidget() *Widget { return &Widget{} }\n"), 0644)

	rgs := NewRepoGraphService(dir)
	graph, err := rgs.Build(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if graph.Stats.TotalFiles < 2 {
		t.Errorf("files: %d", graph.Stats.TotalFiles)
	}
	if graph.Stats.TotalFunctions < 2 {
		t.Errorf("functions: %d", graph.Stats.TotalFunctions)
	}
}

func TestBuildIgnoresVendor(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "vendor", "lib"), 0755)
	os.WriteFile(filepath.Join(dir, "vendor", "lib", "a.go"), []byte("package lib\nfunc Vendored() {}\n"), 0644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\nfunc Main() {}\n"), 0644)

	rgs := NewRepoGraphService(dir)
	graph, _ := rgs.Build(context.Background())
	if graph.Stats.TotalFiles != 1 {
		t.Errorf("should skip vendor, got %d files", graph.Stats.TotalFiles)
	}
}

func TestSearchSymbols(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "util.go"), []byte("package main\nfunc FindUser() {}\nfunc DeleteUser() {}\n"), 0644)

	rgs := NewRepoGraphService(dir)
	rgs.Build(context.Background())
	results := rgs.SearchSymbols("User", 10)
	if len(results) < 2 {
		t.Errorf("results: %d", len(results))
	}
}