package projects

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewProjectTracker(t *testing.T) {
	pt := NewProjectTracker("")
	if pt.Count() != 0 {
		t.Error("should start empty")
	}
}

func TestRegisterAndGet(t *testing.T) {
	dir := t.TempDir()
	pt := NewProjectTracker("")

	proj, err := pt.Register(dir)
	if err != nil {
		t.Fatal(err)
	}
	if proj.Name != filepath.Base(dir) {
		t.Errorf("name: %s", proj.Name)
	}

	got, ok := pt.Get(dir)
	if !ok {
		t.Fatal("should find project")
	}
	if got.ID != proj.ID {
		t.Errorf("id mismatch: %s vs %s", got.ID, proj.ID)
	}
}

func TestRegisterUpdatesSession(t *testing.T) {
	dir := t.TempDir()
	pt := NewProjectTracker("")

	pt.Register(dir)
	pt.Register(dir)

	got, _ := pt.Get(dir)
	if got.SessionCount != 2 {
		t.Errorf("session count: %d", got.SessionCount)
	}
}

func TestList(t *testing.T) {
	pt := NewProjectTracker("")
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	pt.Register(dir1)
	pt.Register(dir2)

	list := pt.List()
	if len(list) != 2 {
		t.Errorf("count: %d", len(list))
	}
}

func TestRecent(t *testing.T) {
	pt := NewProjectTracker("")
	for i := 0; i < 5; i++ {
		pt.Register(t.TempDir())
	}
	recent := pt.Recent(3)
	if len(recent) != 3 {
		t.Errorf("recent: %d", len(recent))
	}
}

func TestRemove(t *testing.T) {
	dir := t.TempDir()
	pt := NewProjectTracker("")
	pt.Register(dir)
	if !pt.Remove(dir) {
		t.Error("should remove")
	}
	if pt.Remove(dir) {
		t.Error("should not remove again")
	}
}

func TestStats(t *testing.T) {
	dir := t.TempDir()
	pt := NewProjectTracker("")
	pt.Register(dir)
	stats := pt.Stats()
	total, _ := stats["total"].(int)
	if total != 1 {
		t.Errorf("total: %d", total)
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
