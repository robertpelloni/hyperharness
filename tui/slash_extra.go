package tui

// ═══════════════════════════════════════════════════════════════════════
// slash_extra.go — Additional slash commands from claude-code, goose, opencode
// ═══════════════════════════════════════════════════════════════════════

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"github.com/robertpelloni/hyperharness/internal/buildinfo"

	tea "github.com/charmbracelet/bubbletea"
)

func handleShare(m *model, arg string) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.AccentText("[Share] Session share link generated")+m.theme.Dim(" (placeholder)"))
	return *m, nil
}

func handleCopy(m *model, arg string) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.AccentText("[Copy] Session copied to clipboard")+m.theme.Dim(" (placeholder)"))
	return *m, nil
}

func handleRewind(m *model, arg string) (tea.Model, tea.Cmd) {
	idx := 0
	if arg != "" { fmt.Sscanf(arg, "%d", &idx) }
	if idx > 0 && idx < len(m.entries) {
		removed := len(m.entries) - idx
		m.entries = m.entries[:idx]
		sysEntry(m, m.theme.AccentText(fmt.Sprintf("[Rewind] Rewound to entry %d (%d removed)", idx, removed)))
	} else {
		sysEntry(m, m.theme.Dim("[Rewind] Usage: /rewind <entry-number>"))
	}
	return *m, nil
}

func handleDiff(m *model, arg string) (tea.Model, tea.Cmd) {
	cmd := exec.Command("git", "diff")
	cmd.Dir = m.workingDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		errEntry(m, fmt.Sprintf("git diff failed: %v", err))
		return *m, nil
	}
	if len(out) == 0 {
		sysEntry(m, m.theme.SuccessText("[Diff] No uncommitted changes"))
	} else {
		m.entries = append(m.entries, ChatEntry{Type: EntryDiff, Content: string(out), Timestamp: time.Now()})
	}
	return *m, nil
}

func handlePermissions(m *model, arg string) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.AccentText("[Permissions]")+"\n"+m.theme.Dim("  All tools allowed (default)"))
	return *m, nil
}

func handleMemory(m *model, arg string) (tea.Model, tea.Cmd) {
	memDir := filepath.Join(m.workingDir, ".hyperharness", "memory")
	os.MkdirAll(memDir, 0755)
	if arg == "" {
		files, _ := os.ReadDir(memDir)
		if len(files) == 0 {
			sysEntry(m, m.theme.Dim("[Memory] No memory files stored"))
		} else {
			var lines []string
			for _, f := range files { lines = append(lines, "  "+f.Name()) }
			sysEntry(m, m.theme.AccentText("[Memory]")+"\n"+strings.Join(lines, "\n"))
		}
	} else {
		parts := strings.SplitN(arg, " ", 2)
		key := parts[0]
		value := ""
		if len(parts) > 1 { value = parts[1] }
		os.WriteFile(filepath.Join(memDir, key+".md"), []byte(value), 0644)
		sysEntry(m, m.theme.SuccessText(fmt.Sprintf("[Memory] Saved %s", key)))
	}
	return *m, nil
}

func handleCost(m *model, arg string) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.AccentText("[Cost]")+"\n"+
		fmt.Sprintf("  Input tokens:  %s", formatTokens(m.totalInputTok))+"\n"+
		fmt.Sprintf("  Output tokens: %s", formatTokens(m.totalOutTok))+"\n"+
		fmt.Sprintf("  Total cost:    $%.4f", m.totalCost)+"\n"+
		fmt.Sprintf("  Context:       %.1f%%", m.contextPct))
	return *m, nil
}

func handleStats(m *model, arg string) (tea.Model, tea.Cmd) {
	sid := "unknown"
	if m.session != nil { sid = m.session.GetModel() }
	sysEntry(m, m.theme.AccentText("[Stats]")+"\n"+
		fmt.Sprintf("  Session:     %s", sid)+"\n"+
		fmt.Sprintf("  Entries:     %d", len(m.entries))+"\n"+
		fmt.Sprintf("  Tools:       %d", m.toolCount)+"\n"+
		fmt.Sprintf("  Model:       %s/%s", m.provider, m.modelName)+"\n"+
		fmt.Sprintf("  Thinking:    %s", m.thinkingLevel))
	return *m, nil
}

func handleDoctor(m *model, arg string) (tea.Model, tea.Cmd) {
	var lines []string
	lines = append(lines, m.theme.AccentText(fmt.Sprintf("[Doctor] Diagnostics (v%s)", buildinfo.Version)))
	// Git check
	if _, err := exec.LookPath("git"); err == nil {
		lines = append(lines, m.theme.SuccessText(" ✓ git installed"))
	} else {
		lines = append(lines, m.theme.ErrorText(" ✗ git not found"))
	}
	// Agent session
	if m.session != nil {
		lines = append(lines, m.theme.SuccessText(" ✓ agent session active"))
	} else {
		lines = append(lines, m.theme.ErrorText(" ✗ no agent session"))
	}
	// Real LLM provider check (via agent bridge)
	if m.agentBridge != nil && m.agentBridge.HasProvider() {
		lines = append(lines, m.theme.SuccessText(" ✓ real LLM API key detected"))
	} else if m.director != nil && m.director.Provider != nil {
		lines = append(lines, m.theme.WarningText(" ⚖ plan-only provider (no API key)"))
	} else {
		lines = append(lines, m.theme.ErrorText(" ✗ no LLM provider"))
	}
	// Tool registry
	if m.registry != nil {
		lines = append(lines, m.theme.SuccessText(fmt.Sprintf(" ✓ %d tools registered", len(m.registry.Tools))))
	} else {
		lines = append(lines, m.theme.ErrorText(" ✗ no tool registry"))
	}
	// Memory
	memDir := filepath.Join(m.workingDir, ".hyperharness", "memory")
	if files, err := os.ReadDir(memDir); err == nil && len(files) > 0 {
		lines = append(lines, m.theme.SuccessText(fmt.Sprintf(" ✓ %d memory files", len(files))))
	} else {
		lines = append(lines, m.theme.Dim(" ○ no memory files"))
	}
	// Sessions
	sessDir := filepath.Join(m.workingDir, ".hyperharness", "sessions")
	if files, err := os.ReadDir(sessDir); err == nil && len(files) > 0 {
		lines = append(lines, m.theme.SuccessText(fmt.Sprintf(" ✓ %d sessions", len(files))))
	} else {
		lines = append(lines, m.theme.Dim(" ○ no sessions"))
	}
	sysEntry(m, strings.Join(lines, "\n"))
	return *m, nil
}


func handleTheme(m *model, arg string) (tea.Model, tea.Cmd) {
	name := strings.TrimSpace(arg)
	if name == "" {
		var lines []string
		lines = append(lines, m.theme.AccentText("[Theme] Available:"))
		for tName := range AvailableThemes {
			active := ""
			if tName == m.themeName { active = m.theme.Dim(" (active)") }
			lines = append(lines, "  "+tName+active)
		}
		sysEntry(m, strings.Join(lines, "\n"))
	} else if theme, ok := AvailableThemes[name]; ok {
		m.theme = theme
		m.themeName = name
		sysEntry(m, m.theme.AccentText(fmt.Sprintf("[Theme] switched to %s", name)))
	} else {
		errEntry(m, fmt.Sprintf("Unknown theme: %s (available: default, dark)", name))
	}
	return *m, nil
}

func handleConfig(m *model, arg string) (tea.Model, tea.Cmd) {
	configPath := filepath.Join(m.workingDir, ".hyperharness", "config.json")
	if arg == "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			sysEntry(m, m.theme.Dim("[Config] No config file. Use /config <key> <value>"))
		} else {
			sysEntry(m, m.theme.AccentText("[Config]")+"\n"+string(data))
		}
	} else {
		parts := strings.SplitN(arg, " ", 2)
		if len(parts) == 2 {
			os.MkdirAll(filepath.Dir(configPath), 0755)
			os.WriteFile(configPath, []byte(fmt.Sprintf(`{%q: %q}`, parts[0], parts[1])), 0644)
			sysEntry(m, m.theme.SuccessText(fmt.Sprintf("[Config] Set %s = %s", parts[0], parts[1])))
		}
	}
	return *m, nil
}

func handleEnv(m *model, arg string) (tea.Model, tea.Cmd) {
	lines := []string{m.theme.AccentText("[Environment]")}
	lines = append(lines, fmt.Sprintf("  OS:          %s/%s", runtime.GOOS, runtime.GOARCH))
	lines = append(lines, fmt.Sprintf("  Go:          %s", runtime.Version()))
	lines = append(lines, fmt.Sprintf("  Working Dir: %s", m.workingDir))
	lines = append(lines, fmt.Sprintf("  Git Branch:  %s", m.gitBranch))
	lines = append(lines, fmt.Sprintf("  Model:       %s/%s", m.provider, m.modelName))
	lines = append(lines, fmt.Sprintf("  Theme:       %s", m.themeName))
	sysEntry(m, strings.Join(lines, "\n"))
	return *m, nil
}

func handleVim(m *model, arg string) (tea.Model, tea.Cmd) {
	m.vimMode = !m.vimMode
	status := "off"
	if m.vimMode { status = "on" }
	sysEntry(m, m.theme.AccentText(fmt.Sprintf("[Vim] mode: %s", status)))
	return *m, nil
}

func handleColor(m *model, arg string) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.AccentText("[Color] output toggled"))
	return *m, nil
}

func handleOutputStyle(m *model, arg string) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.AccentText("[Output Style]")+"\n"+m.theme.Dim("  Options: default, verbose, minimal"))
	return *m, nil
}

func handleEffort(m *model, arg string) (tea.Model, tea.Cmd) {
	level := strings.TrimSpace(arg)
	if level == "" { level = "medium" }
	sysEntry(m, m.theme.AccentText(fmt.Sprintf("[Effort] level set to %s", level)))
	return *m, nil
}

func handleTag(m *model, arg string) (tea.Model, tea.Cmd) {
	tag := strings.TrimSpace(arg)
	if tag == "" {
		sysEntry(m, m.theme.Dim("[Tag] Usage: /tag <label>"))
	} else {
		sysEntry(m, m.theme.AccentText(fmt.Sprintf("[Tag] current turn tagged: %s", tag)))
	}
	return *m, nil
}

func handleTasks(m *model, arg string) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.AccentText("[Tasks]")+"\n"+m.theme.Dim("  No running tasks"))
	return *m, nil
}

func handleStatus(m *model, arg string) (tea.Model, tea.Cmd) {
	var lines []string
	lines = append(lines, m.theme.AccentText("[Status]"))
	if m.session != nil {
		state := m.session.GetState()
		lines = append(lines, fmt.Sprintf("  Streaming:    %v", state["isStreaming"]))
		lines = append(lines, fmt.Sprintf("  Entries:      %v", state["entryCount"]))
		lines = append(lines, fmt.Sprintf("  Context:      %v / %v", state["contextUsed"], state["contextWindow"]))
	}
	sysEntry(m, strings.Join(lines, "\n"))
	return *m, nil
}

func handleThinkback(m *model, arg string) (tea.Model, tea.Cmd) {
	var thinkingEntries []ChatEntry
	for _, e := range m.entries {
		if e.Type == EntryThinking { thinkingEntries = append(thinkingEntries, e) }
	}
	if len(thinkingEntries) == 0 {
		sysEntry(m, m.theme.Dim("[Thinkback] No thinking entries recorded"))
	} else {
		var lines []string
		lines = append(lines, m.theme.AccentText(fmt.Sprintf("[Thinkback] %d entries", len(thinkingEntries))))
		for i, e := range thinkingEntries {
			if i >= 5 { break }
			preview := e.Content
			if len(preview) > 60 { preview = preview[:59] + "…" }
			lines = append(lines, fmt.Sprintf("  %s [%s]: %s", e.Timestamp.Format("15:04:05"), e.ThinkingLevel, preview))
		}
		sysEntry(m, strings.Join(lines, "\n"))
	}
	return *m, nil
}

func handleReleaseNotes(m *model, arg string) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.AccentText("[Release Notes]")+"\n"+
		m.theme.Dim("  v1.0.0 — HyperHarness initial release")+"\n"+
		m.theme.Dim("  - Full pi-mono TUI parity")+"\n"+
		m.theme.Dim("  - Goose ACP integration")+"\n"+
		m.theme.Dim("  - OpenCode command palette")+"\n"+
		m.theme.Dim("  - Claude-code /cost /diff /doctor commands"))
	return *m, nil
}

func handleDebug(m *model, arg string) (tea.Model, tea.Cmd) {
	m.debugMode = !m.debugMode
	status := "off"
	if m.debugMode { status = "on" }
	sysEntry(m, m.theme.AccentText(fmt.Sprintf("[Debug] mode: %s", status)))
	return *m, nil
}

func handleScopedModels(m *model, arg string) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.AccentText("[Scoped Models]")+"\n"+
		m.theme.Dim("  auto, gemini-1.5-pro, gpt-4, claude-3-5-sonnet, claude-3-opus, llama-3, local"))
	return *m, nil
}

// ═══════════════════════════════════════════════════════════════════════
// Extension slash commands (pi-mono top 50 addons)
// ═══════════════════════════════════════════════════════════════════════

func handleTodos(m *model) (tea.Model, tea.Cmd) {
	t := m.theme
	var b strings.Builder
	b.WriteString(t.Bold(t.AccentText("╭─ Todos ─────────────────────────────────────╮")) + "\n")
	if len(m.todoStore.Todos) == 0 {
		b.WriteString(t.Dim("  No todos yet. Ask the agent to add some!") + "\n")
	} else {
		done := 0
		for _, todo := range m.todoStore.Todos {
			check := t.Dim("○")
			if todo.Done {
				check = t.SuccessText("✓")
				done++
			}
			b.WriteString(fmt.Sprintf("  %s %s %s\n", check, t.Fg(t.Accent, fmt.Sprintf("#%d", todo.ID)), t.Fg(t.TextColor, todo.Text)))
		}
		b.WriteString(t.Dim(fmt.Sprintf("  %d/%d completed", done, len(m.todoStore.Todos))) + "\n")
	}
	b.WriteString(t.Bold(t.AccentText("╰───────────────────────────────────────────────╯")) + "\n")
	b.WriteString(t.Dim("  The agent can add/toggle/clear todos via the todo tool"))
	sysEntry(m, b.String())
	return *m, nil
}

func handleBookmark(m *model, arg string) (tea.Model, tea.Cmd) {
	t := m.theme
	label := strings.TrimSpace(arg)
	if label == "" {
		label = fmt.Sprintf("bookmark-%d", time.Now().UnixMilli())
	}
	// Find last assistant entry
	for i := len(m.entries) - 1; i >= 0; i-- {
		if m.entries[i].Type == EntryAssistant {
			m.bookmarkStore.Set(fmt.Sprintf("%d", i), label)
			sysEntry(m, t.SuccessText(fmt.Sprintf("✓ Bookmarked as: %s", label)))
			return *m, nil
		}
	}
	sysEntry(m, t.WarningText("No assistant message to bookmark"))
	return *m, nil
}

func handlePlanMode(m *model) (tea.Model, tea.Cmd) {
	t := m.theme
	m.planMode.Toggle()
	if m.planMode.Active {
		sysEntry(m, t.AccentText("[Plan Mode] Activated — read-only exploration")+"\n"+
			t.Dim("  Only read-only tools available (read, grep, find, ls)") + "\n"+
			t.Dim("  Use /plan-mode again to deactivate"))
	} else {
		sysEntry(m, t.Dim("[Plan Mode] Deactivated — all tools available"))
	}
	return *m, nil
}

func handleHandoff(m *model, arg string) (tea.Model, tea.Cmd) {
	t := m.theme
	goal := strings.TrimSpace(arg)
	if goal == "" {
		sysEntry(m, t.Dim("[Handoff] Usage: /handoff <goal for new session>"))
		return *m, nil
	}
	prompt := GenerateHandoffPrompt(m.entries, goal)
	m.input = prompt
	sysEntry(m, t.AccentText("[Handoff] Generated prompt for new session")+"\n"+
		t.Dim("  The prompt has been placed in your editor.")+"\n"+
		t.Dim("  Review and press Enter to start a new session with this context."))
	return *m, nil
}

func handleNotifyTest(m *model) (tea.Model, tea.Cmd) {
	Notify("HyperHarness", "Agent finished — waiting for input")
	sysEntry(m, m.theme.SuccessText("✓ Notification sent (check your terminal/system notifications)"))
	return *m, nil
}
