package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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

func TestTreeBrowserModeNavigation(t *testing.T) {
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
	oldLeaf, err := runtime.GetLeafID(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	mdl, _ := ProcessSlashCommand("/tree-browser", &m)
	updated := mdl.(model)
	if !updated.browserActive || len(updated.browserItems) < 3 {
		t.Fatalf("expected active browser with items, got %#v", updated)
	}
	if view := updated.View(); !strings.Contains(view, "[Foundation Tree Browser :: Modal]") {
		t.Fatalf("expected modal browser title in view, got %s", view)
	} else if !strings.Contains(view, "[Preview]") {
		t.Fatalf("expected preview in browser view, got %s", view)
	} else if !strings.Contains(view, "branchSummaryEntries=") && !strings.Contains(view, "branchSummary=already on active leaf") {
		t.Fatalf("expected branch-summary preview details, got %s", view)
	} else if !strings.Contains(view, "children=") {
		t.Fatalf("expected richer tree layout with child counts, got %s", view)
	} else if !strings.Contains(view, "├─") && !strings.Contains(view, "└─") {
		t.Fatalf("expected graph-style tree connectors, got %s", view)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyLeft})
	updated = mdl.(model)
	if view := updated.View(); !strings.Contains(view, "[+]") {
		t.Fatalf("expected collapsed cue in browser view, got %s", view)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRight})
	updated = mdl.(model)
	if view := updated.View(); !strings.Contains(view, "[-]") {
		t.Fatalf("expected expanded cue in browser view, got %s", view)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyLeft})
	updated = mdl.(model)
	if view := updated.View(); !strings.Contains(view, "[+]") {
		t.Fatalf("expected collapsed cue in browser view, got %s", view)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRight})
	updated = mdl.(model)
	if view := updated.View(); !strings.Contains(view, "[-]") {
		t.Fatalf("expected expanded cue in browser view, got %s", view)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyTab})
	updated = mdl.(model)
	if !updated.browserGrouped {
		t.Fatal("expected browser grouping to toggle on")
	}
	if view := updated.View(); !strings.Contains(view, "[Group]") {
		t.Fatalf("expected grouped browser view, got %s", view)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("B")})
	updated = mdl.(model)
	if updated.browserFilter != "B" {
		t.Fatalf("expected browser filter to be set, got %q", updated.browserFilter)
	}
	if view := updated.View(); !strings.Contains(view, "filter=\"B\"") {
		t.Fatalf("expected filtered browser view, got %s", view)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	updated = mdl.(model)
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated = mdl.(model)
	if updated.browserIndex != 1 {
		t.Fatalf("expected browser index 1, got %d", updated.browserIndex)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated = mdl.(model)
	if !updated.browserActive || !updated.browserConfirmPending {
		t.Fatal("expected browser to enter confirm-pending state after first Enter")
	}
	if view := updated.View(); !strings.Contains(view, "[Confirm]") {
		t.Fatalf("expected confirm prompt in browser view, got %s", view)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated = mdl.(model)
	if updated.browserActive {
		t.Fatal("expected browser to close after confirmed enter")
	}
	newLeaf, err := runtime.GetLeafID(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	if newLeaf == oldLeaf {
		t.Fatalf("expected browser enter to switch leaf, old=%q new=%q", oldLeaf, newLeaf)
	}
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "[Foundation Tree Switch]") {
		t.Fatalf("expected tree switch output after browser enter, got %#v", updated.history)
	}
}

func TestPinnedTreePaneAutoRefreshesAfterSessionMutation(t *testing.T) {
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
	mdl, _ := ProcessSlashCommand("/tree-pane", &m)
	updated := mdl.(model)
	if !updated.browserPinned {
		t.Fatal("expected tree pane pinned")
	}
	before := updated.View()
	updated.foundationSessionID = sessionID
	updated.history = append(updated.history, "You: extra")
	_ = appendFoundationUserText(cwd, sessionID, "extra")
	refreshPinnedFoundationTreeBrowser(&updated)
	after := updated.View()
	if before == after {
		t.Fatalf("expected pinned tree pane to refresh after mutation; before=%s after=%s", before, after)
	}
	if !strings.Contains(after, "extra") {
		t.Fatalf("expected refreshed pane to include new session content, got %s", after)
	}
}

func TestProcessSlashCommandTreePaneToggle(t *testing.T) {
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
	mdl, _ := ProcessSlashCommand("/tree-pane", &m)
	updated := mdl.(model)
	if !updated.browserPinned {
		t.Fatal("expected browser pane to be pinned")
	}
	if view := updated.View(); !strings.Contains(view, "[Foundation Tree Pane :: Passive]") {
		t.Fatalf("expected pinned tree pane in view, got %s", view)
	} else if !strings.Contains(view, "════════") {
		t.Fatalf("expected clearer pane/history divider in view, got %s", view)
	} else if strings.Contains(view, "matches=") && strings.Contains(view, "showing=") == false && strings.Count(view, "[message]") > 8 {
		t.Fatalf("expected viewport metadata for large pinned pane view, got %s", view)
	}
	mdl, _ = ProcessSlashCommand("/tree-pane", &updated)
	updated = mdl.(model)
	if updated.browserPinned {
		t.Fatal("expected browser pane to be hidden")
	}
}

func TestProcessSlashCommandTreePaneFocusNavigation(t *testing.T) {
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
	oldLeaf, err := runtime.GetLeafID(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	mdl, _ := ProcessSlashCommand("/tree-pane", &m)
	updated := mdl.(model)
	mdl, _ = ProcessSlashCommand("/tree-pane-focus", &updated)
	updated = mdl.(model)
	if !updated.browserPinnedFocus {
		t.Fatal("expected tree pane focus enabled")
	}
	if view := updated.View(); !strings.Contains(view, "[Foundation Tree Pane :: Focused]") {
		t.Fatalf("expected focused pane title, got %s", view)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated = mdl.(model)
	if updated.browserIndex != 1 {
		t.Fatalf("expected pane browser index 1, got %d", updated.browserIndex)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated = mdl.(model)
	if !updated.browserConfirmPending {
		t.Fatal("expected pane confirm pending after Enter")
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	updated = mdl.(model)
	if updated.browserConfirmPending {
		t.Fatal("expected pane confirm cleared after Y")
	}
	newLeaf, err := runtime.GetLeafID(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	if newLeaf == oldLeaf {
		t.Fatalf("expected pane-focused navigation to switch leaf, old=%q new=%q", oldLeaf, newLeaf)
	}
}

func TestProcessSlashCommandTreePaneSize(t *testing.T) {
	m := model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPaneHeight: 8}
	mdl, _ := ProcessSlashCommand("/tree-pane-size 12", &m)
	updated := mdl.(model)
	if updated.browserPaneHeight != 12 {
		t.Fatalf("expected pane height 12, got %d", updated.browserPaneHeight)
	}
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "height set to 12") {
		t.Fatalf("expected pane size message, got %#v", updated.history)
	}
}

func TestProcessSlashCommandTreePanePosition(t *testing.T) {
	m := model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPanePosition: "top"}
	mdl, _ := ProcessSlashCommand("/tree-pane-position bottom", &m)
	updated := mdl.(model)
	if updated.browserPanePosition != "bottom" {
		t.Fatalf("expected pane position bottom, got %q", updated.browserPanePosition)
	}
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "position set to bottom") {
		t.Fatalf("expected pane position message, got %#v", updated.history)
	}
}

func TestTreePaneViewportControls(t *testing.T) {
	cwd := t.TempDir()
	m := model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPaneHeight: 4}
	m.director.WorkingDir = cwd
	sessionID, err := ensureFoundationSession(&m)
	if err != nil {
		t.Fatal(err)
	}
	runtime := foundationpi.NewRuntime(cwd, nil)
	for i := 0; i < 12; i++ {
		if _, err := runtime.AppendUserText(sessionID, fmt.Sprintf("entry-%02d", i)); err != nil {
			t.Fatal(err)
		}
	}
	mdl, _ := ProcessSlashCommand("/tree-pane", &m)
	updated := mdl.(model)
	mdl, _ = ProcessSlashCommand("/tree-pane-focus", &updated)
	updated = mdl.(model)
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyEnd})
	updated = mdl.(model)
	if updated.browserIndex == 0 {
		t.Fatal("expected End to move selection to the end")
	}
	if view := updated.View(); !strings.Contains(view, "showing=") {
		t.Fatalf("expected viewport metadata in focused pane view, got %s", view)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyHome})
	updated = mdl.(model)
	if updated.browserIndex != 0 {
		t.Fatalf("expected Home to reset selection to 0, got %d", updated.browserIndex)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	updated = mdl.(model)
	if updated.browserIndex < 4 {
		t.Fatalf("expected PgDown to advance by viewport step, got %d", updated.browserIndex)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	updated = mdl.(model)
	if updated.browserIndex != 0 {
		t.Fatalf("expected PgUp to move back toward 0, got %d", updated.browserIndex)
	}
}

func TestProcessSlashCommandTreePanePreviewToggle(t *testing.T) {
	m := model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPanePreview: true}
	mdl, _ := ProcessSlashCommand("/tree-pane-preview off", &m)
	updated := mdl.(model)
	if updated.browserPanePreview {
		t.Fatal("expected pane preview to be disabled")
	}
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "preview set to off") {
		t.Fatalf("expected pane preview message, got %#v", updated.history)
	}
}

func TestProcessSlashCommandTreePanePreset(t *testing.T) {
	m := model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPaneHeight: 8, browserPanePreview: true, browserPanePosition: "top", browserGrouped: false}
	mdl, _ := ProcessSlashCommand("/tree-pane-preset compact", &m)
	updated := mdl.(model)
	if updated.browserPaneHeight != 6 || updated.browserPanePreview != false || updated.browserPanePosition != "bottom" || updated.browserGrouped != false {
		t.Fatalf("expected compact preset to apply, got %#v", updated)
	}
	mdl, _ = ProcessSlashCommand("/tree-pane-preset detailed", &updated)
	updated = mdl.(model)
	if updated.browserPaneHeight != 12 || updated.browserPanePreview != true || updated.browserPanePosition != "top" || updated.browserGrouped != false {
		t.Fatalf("expected detailed preset to apply, got %#v", updated)
	}
	mdl, _ = ProcessSlashCommand("/tree-pane-preset navigation", &updated)
	updated = mdl.(model)
	if updated.browserPaneHeight != 10 || updated.browserPanePreview != false || updated.browserPanePosition != "bottom" || updated.browserGrouped != true {
		t.Fatalf("expected navigation preset to apply, got %#v", updated)
	}
	mdl, _ = ProcessSlashCommand("/tree-pane-preset review", &updated)
	updated = mdl.(model)
	if updated.browserPaneHeight != 14 || updated.browserPanePreview != true || updated.browserPanePosition != "top" || updated.browserGrouped != true {
		t.Fatalf("expected review preset to apply, got %#v", updated)
	}
}

func TestProcessSlashCommandTreePaneStatus(t *testing.T) {
	m := model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPinned: true, browserPinnedFocus: true, browserPaneHeight: 10, browserPanePosition: "bottom", browserPanePreview: false, browserGrouped: true, browserFilter: "abc"}
	mdl, _ := ProcessSlashCommand("/tree-pane-status", &m)
	updated := mdl.(model)
	if len(updated.history) == 0 || !strings.Contains(updated.history[len(updated.history)-1], "[Foundation Tree Pane Status]") {
		t.Fatalf("expected pane status output, got %#v", updated.history)
	}
	status := updated.history[len(updated.history)-1]
	for _, needle := range []string{"pinned=true", "focus=true", "height=10", "position=bottom", "preview=false", "grouped=true", "filter=\"abc\""} {
		if !strings.Contains(status, needle) {
			t.Fatalf("expected status to contain %q, got %s", needle, status)
		}
	}
}

func TestProcessSlashCommandTreeBrowserClear(t *testing.T) {
	m := model{director: agents.NewDirector(&agents.DefaultProvider{}), browserFilter: "abc", browserConfirmPending: true, browserCollapsed: map[string]bool{"x": true}, browserIndex: 3}
	mdl, _ := ProcessSlashCommand("/tree-browser-clear", &m)
	updated := mdl.(model)
	if updated.browserFilter != "" || updated.browserConfirmPending || updated.browserCollapsed != nil || updated.browserIndex != 0 {
		t.Fatalf("expected browser transient state cleared, got %#v", updated)
	}
}

func TestProcessSlashCommandTreePaneReset(t *testing.T) {
	m := model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPaneHeight: 20, browserPanePosition: "bottom", browserPanePreview: false, browserGrouped: true, browserPinnedFocus: true, browserConfirmPending: true, browserFilter: "abc", browserCollapsed: map[string]bool{"x": true}}
	mdl, _ := ProcessSlashCommand("/tree-pane-reset", &m)
	updated := mdl.(model)
	if updated.browserPaneHeight != 8 || updated.browserPanePosition != "top" || updated.browserPanePreview != true || updated.browserGrouped != false || updated.browserPinnedFocus != false || updated.browserConfirmPending != false || updated.browserFilter != "" || updated.browserCollapsed != nil {
		t.Fatalf("expected pane defaults restored, got %#v", updated)
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
