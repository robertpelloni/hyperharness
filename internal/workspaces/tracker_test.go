package workspaces

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewWorkspaceTracker(t *testing.T) {
	tmpDir := t.TempDir()
	wt := NewWorkspaceTracker(filepath.Join(tmpDir, "workspaces.json"))
	if wt == nil {
		t.Fatal("should create tracker")
	}
}

func TestRegisterAndList(t *testing.T) {
	tmpDir := t.TempDir()
	wt := NewWorkspaceTracker(filepath.Join(tmpDir, "ws.json"))

	// Create real directories to register
	dir1 := filepath.Join(tmpDir, "project1")
	dir2 := filepath.Join(tmpDir, "project2")
	os.MkdirAll(dir1, 0755)
	os.MkdirAll(dir2, 0755)

	wt.RegisterWorkspace(dir1)
	wt.RegisterWorkspace(dir2)

	workspaces, err := wt.ListWorkspaces()
	if err != nil {
		t.Fatal(err)
	}
	if len(workspaces) != 2 {
		t.Fatalf("expected 2, got %d", len(workspaces))
	}
}

func TestRegisterUpdatesAccessTime(t *testing.T) {
	tmpDir := t.TempDir()
	wt := NewWorkspaceTracker(filepath.Join(tmpDir, "ws.json"))

	dir := filepath.Join(tmpDir, "proj")
	os.MkdirAll(dir, 0755)

	wt.RegisterWorkspace(dir)
	wt.RegisterWorkspace(dir) // Register again

	workspaces, _ := wt.ListWorkspaces()
	if len(workspaces) != 1 {
		t.Errorf("should have 1 workspace (updated, not duplicated), got %d", len(workspaces))
	}
}

func TestGetRecent(t *testing.T) {
	tmpDir := t.TempDir()
	wt := NewWorkspaceTracker(filepath.Join(tmpDir, "ws.json"))

	for i := 0; i < 5; i++ {
		dir := filepath.Join(tmpDir, "proj"+string(rune('0'+i)))
		os.MkdirAll(dir, 0755)
		wt.RegisterWorkspace(dir)
	}

	recent, err := wt.GetRecent(3)
	if err != nil {
		t.Fatal(err)
	}
	if len(recent) != 3 {
		t.Errorf("expected 3 recent, got %d", len(recent))
	}
}

func TestRemove(t *testing.T) {
	tmpDir := t.TempDir()
	wt := NewWorkspaceTracker(filepath.Join(tmpDir, "ws.json"))

	dir := filepath.Join(tmpDir, "proj")
	os.MkdirAll(dir, 0755)
	wt.RegisterWorkspace(dir)
	wt.Remove(dir)

	workspaces, _ := wt.ListWorkspaces()
	if len(workspaces) != 0 {
		t.Error("should be empty after remove")
	}
}

func TestNonexistentDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	wt := NewWorkspaceTracker(filepath.Join(tmpDir, "ws.json"))

	wt.RegisterWorkspace("/nonexistent/path")

	workspaces, _ := wt.ListWorkspaces()
	if len(workspaces) != 0 {
		t.Error("nonexistent directories should be filtered out")
	}
}

func TestRegistryPathDefault(t *testing.T) {
	wt := NewWorkspaceTracker("")
	if wt.registryPath == "" {
		t.Error("should have default path")
	}
}
