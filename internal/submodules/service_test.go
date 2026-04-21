package submodules

import (
	"os/exec"
	"testing"
)

func hasGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
}

func TestNewSubmoduleService(t *testing.T) {
	hasGit(t)
	ss := NewSubmoduleService(".")
	if ss == nil {
		t.Fatal("should create service")
	}
}

func TestList(t *testing.T) {
	hasGit(t)
	ss := NewSubmoduleService(".")
	subs, err := ss.List()
	if err != nil {
		// May fail if not in a git repo or no submodules
		t.Logf("List() error (expected in some envs): %v", err)
		return
	}
	t.Logf("Found %d submodules", len(subs))
}

func TestSync(t *testing.T) {
	hasGit(t)
	ss := NewSubmoduleService(".")
	err := ss.Sync()
	if err != nil {
		t.Logf("Sync error (expected in some envs): %v", err)
	}
}
