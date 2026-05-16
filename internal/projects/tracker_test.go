package projects

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRegisterAndGet(t *testing.T) {
	dir := t.TempDir()
	pt := NewProjectTracker("")
	proj, err := pt.Register(dir)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := pt.Get(dir)
	if !ok {
		t.Fatal("should find project")
	}
	if got.ID != proj.ID {
		t.Errorf("id: %s", got.ID)
	}
}

func TestDetectGoProject(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
	pt := NewProjectTracker("")
	proj, _ := pt.Register(dir)
	if proj.Language != "go" {
		t.Errorf("language: %s", proj.Language)
	}
}

func TestPersistence(t *testing.T) {
	storeDir := t.TempDir()
	storePath := filepath.Join(storeDir, "projects.json")
	dir := t.TempDir()
	pt1 := NewProjectTracker(storePath)
	pt1.Register(dir)
	pt2 := NewProjectTracker(storePath)
	if pt2.Count() != 1 {
		t.Errorf("should load from disk, count: %d", pt2.Count())
	}
}
