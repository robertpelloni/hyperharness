package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/robertpelloni/hyperharness/agents"
	"github.com/robertpelloni/hyperharness/tools"
)

func lastEntryContent(m *model) string {
	if len(m.entries) == 0 {
		return ""
	}
	return m.entries[len(m.entries)-1].Content
}

func TestProcessSlashCommandHelp(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/help", m)
	updated := mdl.(model)
	if !strings.Contains(lastEntryContent(&updated), "Slash Commands") {
		t.Fatalf("expected help output, got %s", lastEntryContent(&updated))
	}
}

func TestProcessSlashCommandHotkeys(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/hotkeys", m)
	updated := mdl.(model)
	if !strings.Contains(lastEntryContent(&updated), "Keyboard Shortcuts") {
		t.Fatalf("expected hotkeys output, got %s", lastEntryContent(&updated))
	}
}

func TestProcessSlashCommandClearResetsDirector(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.entries = append(m.entries, ChatEntry{Type: EntryUser, Content: "old"})
	mdl, _ := ProcessSlashCommand("/clear", m)
	updated := mdl.(model)
	if len(updated.entries) != 1 {
		t.Fatalf("expected 1 entry after clear, got %d", len(updated.entries))
	}
	if updated.director == nil {
		t.Fatal("expected director reset")
	}
}

func TestProcessSlashCommandDashboardToggle(t *testing.T) {
	m := &model{dashboardActive: false, director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/dashboard", m)
	updated := mdl.(model)
	if !updated.dashboardActive {
		t.Error("Expected dashboardActive to be true after /dashboard")
	}
	mdl, _ = ProcessSlashCommand("/dashboard", &updated)
	updated = mdl.(model)
	if updated.dashboardActive {
		t.Error("Expected dashboardActive to be false after second /dashboard")
	}
}

func TestProcessSlashCommandSettings(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), provider: "hypercode", modelName: "auto"}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/settings", m)
	updated := mdl.(model)
	if !strings.Contains(lastEntryContent(&updated), "Settings") {
		t.Fatalf("expected settings output, got %s", lastEntryContent(&updated))
	}
}

func TestProcessSlashCommandSession(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/session", m)
	updated := mdl.(model)
	if !strings.Contains(lastEntryContent(&updated), "Session Info") {
		t.Fatalf("expected session output, got %s", lastEntryContent(&updated))
	}
}

func TestProcessSlashCommandModel(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), provider: "test", modelName: "auto"}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/model gpt-4", m)
	updated := mdl.(model)
	if updated.modelName != "gpt-4" {
		t.Fatalf("expected model set to gpt-4, got %s", updated.modelName)
	}
}

func TestProcessSlashCommandCompact(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/compact", m)
	updated := mdl.(model)
	if !strings.Contains(lastEntryContent(&updated), "compacted") {
		t.Fatalf("expected compact output, got %s", lastEntryContent(&updated))
	}
}

func TestProcessSlashCommandTools(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), registry: tools.NewRegistry()}
	m.director.WorkingDir = "."
	m.toolCount = len(m.registry.Tools)
	mdl, _ := ProcessSlashCommand("/tools", m)
	updated := mdl.(model)
	if !strings.Contains(lastEntryContent(&updated), "Registered Tools") {
		t.Fatalf("expected tools output, got %s", lastEntryContent(&updated))
	}
}

func TestProcessSlashCommandTreePaneToggle(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/tree-pane", m)
	updated := mdl.(model)
	if !updated.browserPinned {
		t.Fatal("expected browser pane to be pinned")
	}
	mdl, _ = ProcessSlashCommand("/tree-pane", &updated)
	updated = mdl.(model)
	if updated.browserPinned {
		t.Fatal("expected browser pane to be hidden")
	}
}

func TestProcessSlashCommandTreePaneSize(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPaneHeight: 8}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/tree-pane-size 12", m)
	updated := mdl.(model)
	if updated.browserPaneHeight != 12 {
		t.Fatalf("expected pane height 12, got %d", updated.browserPaneHeight)
	}
}

func TestProcessSlashCommandTreePanePosition(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPanePosition: "top"}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/tree-pane-position bottom", m)
	updated := mdl.(model)
	if updated.browserPanePosition != "bottom" {
		t.Fatalf("expected position bottom, got %q", updated.browserPanePosition)
	}
}

func TestProcessSlashCommandTreePanePreset(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPaneHeight: 8, browserPanePreview: true, browserPanePosition: "top", browserGrouped: false}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/tree-pane-preset compact", m)
	updated := mdl.(model)
	if updated.browserPaneHeight != 6 || updated.browserPanePreview || updated.browserPanePosition != "bottom" || updated.browserGrouped {
		t.Fatalf("expected compact preset, got h=%d preview=%t pos=%s grouped=%t",
			updated.browserPaneHeight, updated.browserPanePreview, updated.browserPanePosition, updated.browserGrouped)
	}
}

func TestProcessSlashCommandTreePaneReset(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPaneHeight: 20, browserPanePosition: "bottom", browserPanePreview: false, browserGrouped: true, browserPinnedFocus: true, browserConfirmPending: true, browserFilter: "abc", browserCollapsed: map[string]bool{"x": true}}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/tree-pane-reset", m)
	updated := mdl.(model)
	if updated.browserPaneHeight != 8 || updated.browserPanePosition != "top" || !updated.browserPanePreview || updated.browserGrouped || updated.browserPinnedFocus || updated.browserConfirmPending || updated.browserFilter != "" || updated.browserCollapsed != nil {
		t.Fatalf("expected defaults restored, got h=%d pos=%s preview=%t grouped=%t focus=%t confirm=%t filter=%q collapsed=%v",
			updated.browserPaneHeight, updated.browserPanePosition, updated.browserPanePreview, updated.browserGrouped,
			updated.browserPinnedFocus, updated.browserConfirmPending, updated.browserFilter, updated.browserCollapsed)
	}
}

func TestProcessSlashCommandTreePaneGrouped(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), browserGrouped: false}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/tree-pane-grouped on", m)
	updated := mdl.(model)
	if !updated.browserGrouped {
		t.Fatal("expected grouped enabled")
	}
	mdl, _ = ProcessSlashCommand("/tree-pane-grouped toggle", &updated)
	updated = mdl.(model)
	if updated.browserGrouped {
		t.Fatal("expected grouped toggled off")
	}
}

func TestProcessSlashCommandTreePaneSizeCycle(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPaneHeight: 8}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/tree-pane-size-cycle", m)
	updated := mdl.(model)
	if updated.browserPaneHeight != 10 {
		t.Fatalf("expected first cycle to 10, got %d", updated.browserPaneHeight)
	}
}

func TestProcessSlashCommandTreeBrowserClear(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), browserFilter: "abc", browserConfirmPending: true, browserCollapsed: map[string]bool{"x": true}, browserIndex: 3}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/tree-browser-clear", m)
	updated := mdl.(model)
	if updated.browserFilter != "" || updated.browserConfirmPending || updated.browserCollapsed != nil || updated.browserIndex != 0 {
		t.Fatalf("expected browser state cleared, got filter=%q confirm=%t collapsed=%v index=%d",
			updated.browserFilter, updated.browserConfirmPending, updated.browserCollapsed, updated.browserIndex)
	}
}

func TestProcessSlashCommandTreePaneStatus(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPinned: true, browserPinnedFocus: true, browserPaneHeight: 10, browserPanePosition: "bottom", browserPanePreview: false, browserGrouped: true, browserFilter: "abc"}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/tree-pane-status", m)
	updated := mdl.(model)
	status := lastEntryContent(&updated)
	for _, needle := range []string{"pinned:   true", "focus:    true", "height:   10", "position: bottom"} {
		if !strings.Contains(status, needle) {
			t.Fatalf("expected status to contain %q, got %s", needle, status)
		}
	}
}

func TestProcessSlashCommandTreePaneCycle(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPaneHeight: 8, browserPanePreview: true, browserPanePosition: "top", browserGrouped: false}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/tree-pane-cycle", m)
	updated := mdl.(model)
	if updated.browserPaneHeight != 6 {
		t.Fatalf("expected first cycle to compact (h=6), got h=%d", updated.browserPaneHeight)
	}
}

func TestProcessSlashCommandUnknown(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/nonexistent", m)
	updated := mdl.(model)
	if !strings.Contains(lastEntryContent(&updated), "Unknown command") {
		t.Fatalf("expected unknown command error, got %s", lastEntryContent(&updated))
	}
}

func TestTreeBrowserNavigation(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), browserActive: true, browserItems: []TreeBrowserItem{{ID: "1", Name: "a", IsDir: false}, {ID: "2", Name: "b", IsDir: false}}, browserIndex: 0}
	m.director.WorkingDir = "."
	mdl, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated := mdl.(model)
	if updated.browserIndex != 1 {
		t.Fatalf("expected browser index 1, got %d", updated.browserIndex)
	}
	mdl, _ = updated.Update(tea.KeyMsg{Type: tea.KeyUp})
	updated = mdl.(model)
	if updated.browserIndex != 0 {
		t.Fatalf("expected browser index 0, got %d", updated.browserIndex)
	}
}

func TestProcessSlashCommandTreePaneFocusToggle(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPinned: true, browserPinnedFocus: false}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/tree-pane-focus-toggle", m)
	updated := mdl.(model)
	if !updated.browserPinnedFocus {
		t.Fatal("expected focus toggle to enable focus")
	}
	mdl, _ = ProcessSlashCommand("/tree-pane-focus-toggle", &updated)
	updated = mdl.(model)
	if updated.browserPinnedFocus {
		t.Fatal("expected second focus toggle to disable focus")
	}
}

func TestProcessSlashCommandTreePanePreviewToggle(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{}), browserPanePreview: true}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/tree-pane-preview off", m)
	updated := mdl.(model)
	if updated.browserPanePreview {
		t.Fatal("expected preview disabled")
	}
}

func TestProcessSlashCommandNewSession(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = "."
	m.entries = append(m.entries, ChatEntry{Type: EntryUser, Content: "old stuff"})
	mdl, _ := ProcessSlashCommand("/new", m)
	updated := mdl.(model)
	if len(updated.entries) != 1 {
		t.Fatalf("expected 1 entry after new, got %d", len(updated.entries))
	}
}

func TestProcessSlashCommandReload(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/reload", m)
	updated := mdl.(model)
	if !strings.Contains(lastEntryContent(&updated), "reloaded") {
		t.Fatalf("expected reload message, got %s", lastEntryContent(&updated))
	}
}

func TestProcessSlashCommandFork(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/fork", m)
	updated := mdl.(model)
	if !strings.Contains(lastEntryContent(&updated), "Fork") {
		t.Fatalf("expected fork output, got %s", lastEntryContent(&updated))
	}
}

func TestProcessSlashCommandName(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/name my-session", m)
	updated := mdl.(model)
	if !strings.Contains(lastEntryContent(&updated), "my-session") {
		t.Fatalf("expected name set, got %s", lastEntryContent(&updated))
	}
}

func TestProcessSlashCommandExport(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = "."
	m.entries = append(m.entries, ChatEntry{Type: EntryUser, Content: "test"})
	mdl, _ := ProcessSlashCommand("/export", m)
	updated := mdl.(model)
	if !strings.Contains(lastEntryContent(&updated), "exported") {
		t.Fatalf("expected export output, got %s", lastEntryContent(&updated))
	}
}

func TestProcessSlashCommandTreePaneShowHide(t *testing.T) {
	m := &model{director: agents.NewDirector(&agents.DefaultProvider{})}
	m.director.WorkingDir = "."
	mdl, _ := ProcessSlashCommand("/tree-pane-show", m)
	updated := mdl.(model)
	if !updated.browserPinned {
		t.Fatal("expected tree pane shown")
	}
	mdl, _ = ProcessSlashCommand("/tree-pane-hide", &updated)
	updated = mdl.(model)
	if updated.browserPinned {
		t.Fatal("expected tree pane hidden")
	}
}
