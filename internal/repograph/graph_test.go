package repograph

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildGraph(t *testing.T) {
	dir := t.TempDir()
	// Create test files
	os.MkdirAll(filepath.Join(dir, "pkg"), 0755)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte(`package main
import "fmt"
func main() { fmt.Println("hello") }
`), 0644)
	os.WriteFile(filepath.Join(dir, "pkg", "util.go"), []byte(`package pkg
type Widget struct { Name string }
func NewWidget() *Widget { return &Widget{} }
`), 0644)

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
	if graph.Stats.TotalTypes < 1 {
		t.Errorf("types: %d", graph.Stats.TotalTypes)
	}
}

func TestBuildIgnoresVendor(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "vendor", "lib"), 0755)
	os.WriteFile(filepath.Join(dir, "vendor", "lib", "a.go"), []byte(`package lib
func Vendored() {}`), 0644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte(`package main
func Main() {}`), 0644)

	rgs := NewRepoGraphService(dir)
	graph, _ := rgs.Build(context.Background())
	if graph.Stats.TotalFiles != 1 {
		t.Errorf("should skip vendor, got %d files", graph.Stats.TotalFiles)
	}
}

func TestTypeScriptIndexing(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "app.ts"), []byte(`import { foo } from './foo';
export function hello(): string { return "hi"; }
export class App { run() {} }
export interface IApp { start(): void; }
`), 0644)

	rgs := NewRepoGraphService(dir)
	graph, _ := rgs.Build(context.Background())
	if graph.Stats.TotalFunctions < 1 {
		t.Errorf("functions: %d", graph.Stats.TotalFunctions)
	}
}

func TestPythonIndexing(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "main.py"), []byte(`import os
from sys import argv

def hello():
    pass

class MyClass:
    def method(self):
        pass
`), 0644)

	rgs := NewRepoGraphService(dir)
	graph, _ := rgs.Build(context.Background())
	if graph.Stats.TotalFunctions < 1 {
		t.Errorf("functions: %d", graph.Stats.TotalFunctions)
	}
}

func TestSearchSymbols(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "util.go"), []byte(`package main
func FindUser() {}
func DeleteUser() {}
`), 0644)

	rgs := NewRepoGraphService(dir)
	rgs.Build(context.Background())
	results := rgs.SearchSymbols("User", 10)
	if len(results) < 2 {
		t.Errorf("results: %d", len(results))
	}
}

func TestGetGraphBeforeBuild(t *testing.T) {
	rgs := NewRepoGraphService(t.TempDir())
	if rgs.GetGraph() != nil {
		t.Error("should be nil before build")
	}
}

func TestSearchBeforeBuild(t *testing.T) {
	rgs := NewRepoGraphService(t.TempDir())
	if rgs.SearchSymbols("test", 10) != nil {
		t.Error("should be nil before build")
	}
}

func TestFindDependentsBeforeBuild(t *testing.T) {
	rgs := NewRepoGraphService(t.TempDir())
	if rgs.FindDependents("main.go") != nil {
		t.Error("should be nil before build")
	}
}

func TestGoExportedDetection(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "api.go"), []byte(`package api
func PublicFunc() {}
func privateFunc() {}
type PublicType struct {}
type privateType struct {}
`), 0644)

	rgs := NewRepoGraphService(dir)
	graph, _ := rgs.Build(context.Background())

	pub, ok := graph.Nodes["api.go#PublicFunc"]
	if !ok || !pub.IsExported {
		t.Error("PublicFunc should be exported")
	}
	priv, ok := graph.Nodes["api.go#privateFunc"]
	if !ok || priv.IsExported {
		t.Error("privateFunc should not be exported")
	}
}
