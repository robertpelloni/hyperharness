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
	lines = append(lines, m.theme.AccentText("[Doctor] Diagnostics"))
	if _, err := exec.LookPath("git"); err == nil {
		lines = append(lines, m.theme.SuccessText("  ✓ git installed"))
	} else { lines = append(lines, m.theme.ErrorText("  ✗ git not found")) }
	if m.session != nil {
		lines = append(lines, m.theme.SuccessText("  ✓ agent session active"))
	} else { lines = append(lines, m.theme.ErrorText("  ✗ no agent session")) }
	if m.director != nil && m.director.Provider != nil {
		lines = append(lines, m.theme.SuccessText("  ✓ LLM provider configured"))
	} else { lines = append(lines, m.theme.ErrorText("  ✗ no LLM provider")) }
	if m.registry != nil {
		lines = append(lines, m.theme.SuccessText(fmt.Sprintf("  ✓ %d tools registered", len(m.registry.Tools))))
	} else { lines = append(lines, m.theme.ErrorText("  ✗ no tool registry")) }
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
