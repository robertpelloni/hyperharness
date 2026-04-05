package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertpelloni/hypercode/agents"
	foundationpi "github.com/robertpelloni/hypercode/foundation/pi"
)

func TestProcessSlashCommandPlanAndRepomap(t *testing.T) {
	cwd := t.TempDir()
	if err := os.WriteFile(filepath.Join(cwd, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	m := model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = cwd
	mdl, _ := ProcessSlashCommand("/plan analyze this repo", &m)
	updated := mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "[Foundation Plan]") {
		t.Fatalf("expected plan output, got %#v", updated.history)
	}
	mdl, _ = ProcessSlashCommand("/repomap", &updated)
	updated = mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "<repo_map>") {
		t.Fatalf("expected repomap output, got %#v", updated.history)
	}
}

func TestProcessSlashCommandProvidersAndAdapters(t *testing.T) {
	cwd := t.TempDir()
	m := model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = cwd
	mdl, _ := ProcessSlashCommand("/providers", &m)
	updated := mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "[Foundation Providers]") {
		t.Fatalf("expected provider output, got %#v", updated.history)
	}
	mdl, _ = ProcessSlashCommand("/adapters", &updated)
	updated = mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "[Foundation Adapters]") {
		t.Fatalf("expected adapter output, got %#v", updated.history)
	}
}

func TestProcessSlashCommandFoundationSummarySurfaces(t *testing.T) {
	cwd := t.TempDir()
	m := model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = cwd
	sessionID, err := ensureFoundationSession(&m)
	if err != nil {
		t.Fatal(err)
	}
	runtime := foundationpi.NewRuntime(cwd, nil)
	if _, err := runtime.AppendUserText(sessionID, strings.Repeat("older ", 40)); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.AppendUserText(sessionID, strings.Repeat("recent ", 40)); err != nil {
		t.Fatal(err)
	}
	mdl, _ := ProcessSlashCommand("/summary-compact 60", &m)
	updated := mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "[Foundation Compaction Summary]") {
		t.Fatalf("expected compaction summary output, got %#v", updated.history)
	}

	// Build a small branch topology for /summary-branch.
	m = updated
	session, err := runtime.LoadSession(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.AppendUserText(sessionID, "A"); err != nil {
		t.Fatal(err)
	}
	aID, err := runtime.GetLeafID(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.AppendUserText(sessionID, "B"); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.AppendUserText(sessionID, "C"); err != nil {
		t.Fatal(err)
	}
	oldLeaf, err := runtime.GetLeafID(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.BranchSession(sessionID, aID); err != nil {
		t.Fatal(err)
	}
	session, err = runtime.AppendUserText(sessionID, "E")
	if err != nil {
		t.Fatal(err)
	}
	targetID := session.Metadata.LeafID
	if _, err := runtime.BranchSession(sessionID, oldLeaf); err != nil {
		t.Fatal(err)
	}
	mdl, _ = ProcessSlashCommand("/summary-branch "+targetID+" 128", &m)
	updated = mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "[Foundation Branch Summary]") {
		t.Fatalf("expected branch summary output, got %#v", updated.history)
	}
}

func TestProcessSlashCommandTreeSurfaces(t *testing.T) {
	cwd := t.TempDir()
	m := model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = cwd
	sessionID, err := ensureFoundationSession(&m)
	if err != nil {
		t.Fatal(err)
	}
	runtime := foundationpi.NewRuntime(cwd, nil)
	if _, err := runtime.AppendUserText(sessionID, "A"); err != nil {
		t.Fatal(err)
	}
	aID, err := runtime.GetLeafID(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.AppendUserText(sessionID, "B"); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.AppendUserText(sessionID, "C"); err != nil {
		t.Fatal(err)
	}
	oldLeaf, err := runtime.GetLeafID(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.BranchSession(sessionID, aID); err != nil {
		t.Fatal(err)
	}
	session, err := runtime.AppendUserText(sessionID, "E")
	if err != nil {
		t.Fatal(err)
	}
	targetID := session.Metadata.LeafID
	if _, err := runtime.BranchSession(sessionID, oldLeaf); err != nil {
		t.Fatal(err)
	}
	mdl, _ := ProcessSlashCommand("/tree", &m)
	updated := mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "[Foundation Tree]") {
		t.Fatalf("expected tree output, got %#v", updated.history)
	}
	mdl, _ = ProcessSlashCommand("/tree "+targetID+" 128", &updated)
	updated = mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "[Foundation Tree Switch]") {
		t.Fatalf("expected tree switch output, got %#v", updated.history)
	}
	if !strings.Contains(updated.history[len(updated.history)-1], "## Goal") {
		t.Fatalf("expected structured branch summary in tree switch output, got %#v", updated.history[len(updated.history)-1])
	}
	newLeaf, err := runtime.GetLeafID(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	if newLeaf == oldLeaf {
		t.Fatalf("expected leaf to move away from old leaf %q", oldLeaf)
	}
}

func TestProcessSlashCommandTreeExplorerLabelAndChildren(t *testing.T) {
	cwd := t.TempDir()
	m := model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = cwd
	sessionID, err := ensureFoundationSession(&m)
	if err != nil {
		t.Fatal(err)
	}
	runtime := foundationpi.NewRuntime(cwd, nil)
	if _, err := runtime.AppendUserText(sessionID, "A"); err != nil {
		t.Fatal(err)
	}
	aID, err := runtime.GetLeafID(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.AppendUserText(sessionID, "B"); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.AppendUserText(sessionID, "C"); err != nil {
		t.Fatal(err)
	}
	oldLeaf, err := runtime.GetLeafID(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.BranchSession(sessionID, aID); err != nil {
		t.Fatal(err)
	}
	_, err = runtime.AppendUserText(sessionID, "E")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.BranchSession(sessionID, oldLeaf); err != nil {
		t.Fatal(err)
	}
	mdl, _ := ProcessSlashCommand("/label "+aID+" checkpoint-1", &m)
	updated := mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "[Foundation Label]") {
		t.Fatalf("expected label output, got %#v", updated.history)
	}
	mdl, _ = ProcessSlashCommand("/tree", &updated)
	updated = mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "label=\"checkpoint-1\"") {
		t.Fatalf("expected labeled tree output, got %#v", updated.history[len(updated.history)-1])
	}
	mdl, _ = ProcessSlashCommand("/tree-children "+aID, &updated)
	updated = mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "[Foundation Tree Children]") {
		t.Fatalf("expected tree children output, got %#v", updated.history)
	}
	if !strings.Contains(updated.history[len(updated.history)-1], "B") && !strings.Contains(updated.history[len(updated.history)-1], "E") {
		t.Fatalf("expected child branch previews in output, got %#v", updated.history[len(updated.history)-1])
	}
}

func TestProcessSlashCommandTreeSelectorSurfaces(t *testing.T) {
	cwd := t.TempDir()
	m := model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = cwd
	sessionID, err := ensureFoundationSession(&m)
	if err != nil {
		t.Fatal(err)
	}
	runtime := foundationpi.NewRuntime(cwd, nil)
	if _, err := runtime.AppendUserText(sessionID, "A"); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.AppendUserText(sessionID, "B"); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.AppendUserText(sessionID, "C"); err != nil {
		t.Fatal(err)
	}
	mdl, _ := ProcessSlashCommand("/tree-select", &m)
	updated := mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "[Foundation Tree Select]") {
		t.Fatalf("expected tree select output, got %#v", updated.history)
	}
	if len(updated.foundationTreeSelection) < 3 {
		t.Fatalf("expected selection ids, got %#v", updated.foundationTreeSelection)
	}
	oldLeaf, err := runtime.GetLeafID(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	mdl, _ = ProcessSlashCommand("/tree-go 1 128", &updated)
	updated = mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "[Foundation Tree Switch]") {
		t.Fatalf("expected tree-go switch output, got %#v", updated.history)
	}
	newLeaf, err := runtime.GetLeafID(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	if newLeaf == oldLeaf {
		t.Fatalf("expected leaf to change after tree-go, old=%q new=%q", oldLeaf, newLeaf)
	}
	if !strings.Contains(updated.history[len(updated.history)-1], "## Goal") {
		t.Fatalf("expected structured summary in tree-go output, got %#v", updated.history[len(updated.history)-1])
	}
}

func TestProcessSlashCommandClearResetsDirector(t *testing.T) {
	m := model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.history = []string{"old"}
	mdl, _ := ProcessSlashCommand("/clear", &m)
	updated := mdl.(model)
	if len(updated.history) != 1 {
		t.Fatalf("expected reset history, got %#v", updated.history)
	}
	if updated.director == nil {
		t.Fatal("expected director reset")
	}
}
