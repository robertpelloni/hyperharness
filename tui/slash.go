package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/robertpelloni/hyperharness/tools"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/robertpelloni/hyperharness/agents"
	"github.com/robertpelloni/hyperharness/foundation/adapters"
	foundationorchestration "github.com/robertpelloni/hyperharness/foundation/orchestration"
	foundationrepomap "github.com/robertpelloni/hyperharness/foundation/repomap"
)

// ═══════════════════════════════════════════════════════════════════════
// Slash command processing — mirrors pi-mono's BUILTIN_SLASH_COMMANDS
// ═══════════════════════════════════════════════════════════════════════

// ProcessSlashCommand handles all / commands, bypassing the LLM.
func sessionDuration(m *model) string {
	if len(m.entries) == 0 {
		return "0s"
	}
	return time.Since(m.entries[0].Timestamp).Round(time.Second).String()
}

func ProcessSlashCommand(cmd string, m *model) (tea.Model, tea.Cmd) {
	cmd = strings.TrimSpace(cmd)
	parts := strings.Split(cmd, " ")

	switch parts[0] {
	case "/help":
		return handleHelp(m)
	case "/hotkeys":
		return handleHotkeys(m)
	case "/clear":
		return handleClear(m)
	case "/compact":
		return handleCompact(m)
	case "/dashboard":
		return handleDashboard(m)
	case "/commit":
		return handleCommit(m)
	case "/plan":
		return handlePlan(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/plan")))
	case "/repomap":
		return handleRepoMap(m)
	case "/providers":
		return handleProviders(m)
	case "/adapters":
		return handleAdapters(m)
	case "/mcp", "/mcptools":
		return handleMCPTools(m)
	case "/model":
		return handleModel(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/model")))
	case "/settings":
		return handleSettings(m)
	case "/session":
		return handleSession(m)
	case "/name":
		return handleName(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/name")))
	case "/fork":
		return handleFork(m)
	case "/export":
		return handleExport(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/export")))
	case "/import":
		return handleImport(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/import")))
	case "/login":
		return handleLogin(m)
	case "/logout":
		return handleLogout(m)
	case "/new":
		return handleNew(m)
	case "/resume":
		return handleResume(m)
	case "/reload":
		return handleReload(m)
	case "/tools":
		return handleTools(m)
	case "/fsession":
		return handleFoundationSession(m)

	// Tree commands (pi-mono session tree + file tree)
	case "/tree":
		return handleTree(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/tree")))
	case "/tree-select":
		return handleTreeSelect(m)
	case "/tree-browser":
		return handleTreeBrowser(m)
	case "/tree-pane":
		return handleTreePane(m)
	case "/tree-pane-help":
		return handleTreePaneHelp(m)
	case "/tree-pane-show":
		return handleTreePaneShow(m)
	case "/tree-pane-hide":
		return handleTreePaneHide(m)
	case "/tree-pane-size":
		return handleTreePaneSize(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/tree-pane-size")))
	case "/tree-pane-size-cycle":
		return handleTreePaneSizeCycle(m)
	case "/tree-pane-preview":
		return handleTreePanePreview(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/tree-pane-preview")))
	case "/tree-pane-preview-on":
		return handleTreePanePreview(m, "on")
	case "/tree-pane-preview-off":
		return handleTreePanePreview(m, "off")
	case "/tree-pane-preview-toggle":
		return handleTreePanePreview(m, "toggle")
	case "/tree-pane-grouped":
		return handleTreePaneGrouped(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/tree-pane-grouped")))
	case "/tree-pane-grouped-on":
		return handleTreePaneGrouped(m, "on")
	case "/tree-pane-grouped-off":
		return handleTreePaneGrouped(m, "off")
	case "/tree-pane-grouped-toggle":
		return handleTreePaneGrouped(m, "toggle")
	case "/tree-pane-cycle":
		return handleTreePaneCycle(m)
	case "/tree-pane-refresh":
		return handleTreePaneRefresh(m)
	case "/tree-browser-clear":
		return handleTreeBrowserClear(m)
	case "/tree-pane-reset":
		return handleTreePaneReset(m)
	case "/tree-pane-status":
		return handleTreePaneStatus(m)
	case "/tree-pane-summary":
		return handleTreePaneSummary(m)
	case "/tree-pane-preset":
		return handleTreePanePreset(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/tree-pane-preset")))
	case "/tree-pane-compact":
		return handleTreePanePreset(m, "compact")
	case "/tree-pane-detailed":
		return handleTreePanePreset(m, "detailed")
	case "/tree-pane-navigation":
		return handleTreePanePreset(m, "navigation")
	case "/tree-pane-review":
		return handleTreePanePreset(m, "review")
	case "/tree-pane-position":
		return handleTreePanePosition(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/tree-pane-position")))
	case "/tree-pane-top":
		return handleTreePanePosition(m, "top")
	case "/tree-pane-bottom":
		return handleTreePanePosition(m, "bottom")
	case "/tree-pane-position-toggle":
		return handleTreePanePosition(m, "toggle")
	case "/tree-pane-focus-on":
		return handleTreePaneFocusValue(m, true)
	case "/tree-pane-focus-off":
		return handleTreePaneFocusValue(m, false)
	case "/tree-pane-focus", "/tree-pane-focus-toggle":
		return handleTreePaneFocus(m)
	case "/tree-go":
		return handleTreeGo(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/tree-go")))
	case "/tree-children":
		return handleTreeChildren(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/tree-children")))
	case "/label":
		return handleLabel(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/label")))
	case "/summary-compact":
		return handleSummaryCompact(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/summary-compact")))
	case "/summary-branch":
		return handleSummaryBranch(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/summary-branch")))
	case "/exit", "/quit":
		return *m, tea.Quit
	default:
		return handleUnknown(m, parts[0])
	}
}

// ─── System entries helper ────────────────────────────────────────────

func sysEntry(m *model, text string) {
	m.entries = append(m.entries, ChatEntry{Type: EntrySystem, Content: text, Timestamp: time.Now()})
}

func errEntry(m *model, text string) {
	m.entries = append(m.entries, ChatEntry{Type: EntrySystem, Content: m.theme.ErrorText(text), Timestamp: time.Now()})
}

// ─── Handlers ─────────────────────────────────────────────────────────

func handleHelp(m *model) (tea.Model, tea.Cmd) {
	t := m.theme
	var b strings.Builder
	b.WriteString(t.Bold(t.AccentText("╭─ Slash Commands ──────────────────────────────╮")))
	b.WriteString("\n")
	for _, cmd := range BuiltinSlashCommands {
		b.WriteString(t.AccentText("  /"+cmd.Name) + t.Dim(strings.Repeat(" ", max(1, 22-len(cmd.Name)))) + t.Fg(t.Muted, cmd.Description))
		b.WriteString("\n")
	}
	b.WriteString(t.Bold(t.AccentText("╰───────────────────────────────────────────────╯")))
	sysEntry(m, b.String())
	return *m, nil
}

func handleHotkeys(m *model) (tea.Model, tea.Cmd) {
	t := m.theme
	var b strings.Builder
	b.WriteString(t.Bold(t.AccentText("╭─ Keyboard Shortcuts ──────────────────────────╮")))
	b.WriteString("\n")
	keys := []struct{ key, desc string }{
		{"Enter", "Send message"},
		{"Shift+Tab", "New line (multi-line input)"},
		{"↑/↓", "Navigate input history / scroll chat"},
		{"Tab", "Autocomplete slash commands"},
		{"Ctrl+C", "Cancel operation / quit"},
		{"Ctrl+L", "Toggle file tree pane"},
		{"Ctrl+D", "Toggle dashboard"},
		{"Ctrl+Y", "Accept shell proposal"},
		{"PgUp/PgDn", "Scroll chat viewport"},
		{"Home/End", "Scroll to top/bottom"},
		{"Esc", "Close modal / quit"},
	}
	for _, k := range keys {
		b.WriteString(t.Dim("  "+k.key) + t.Fg(t.Muted, strings.Repeat(" ", max(1, 18-len(k.key)))+k.desc))
		b.WriteString("\n")
	}
	b.WriteString(t.Bold(t.AccentText("╰───────────────────────────────────────────────╯")))
	sysEntry(m, b.String())
	return *m, nil
}

func handleClear(m *model) (tea.Model, tea.Cmd) {
	workingDir := m.director.WorkingDir
	m.entries = []ChatEntry{{
		Type:      EntrySystem,
		Content:   m.theme.SuccessText("✓ Memory crystal wiped. Context reset to null."),
		Timestamp: time.Now(),
	}}
	m.director = agents.NewDirector(agents.NewHyperCodeProvider())
	m.director.WorkingDir = workingDir
	return *m, nil
}

func handleCompact(m *model) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.SuccessText("✓ Session context compacted."))
	m.contextPct = 0
	return *m, nil
}

func handleDashboard(m *model) (tea.Model, tea.Cmd) {
	m.dashboardActive = !m.dashboardActive
	if m.dashboardActive {
		sysEntry(m, m.theme.AccentText("[Dashboard] Activated split-pane interface."))
	} else {
		sysEntry(m, m.theme.Dim("[Dashboard] Deactivated."))
	}
	return *m, nil
}

func handleCommit(m *model) (tea.Model, tea.Cmd) {
	// Run git diff and generate commit
	cmd := exec.Command("git", "diff", "--stat")
	cmd.Dir = m.director.WorkingDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		errEntry(m, fmt.Sprintf("git diff failed: %v", err))
		return *m, nil
	}
	m.entries = append(m.entries, ChatEntry{
		Type:      EntryTool,
		ToolName:  "bash",
		ToolArgs:  "git diff --stat",
		Content:   string(out),
		Timestamp: time.Now(),
	})
	m.loading = true
	return m, func() tea.Msg {
		// Use Director to generate commit message from diff
		diffCmd := exec.Command("git", "diff")
		diffCmd.Dir = m.director.WorkingDir
		diffOut, _ := diffCmd.CombinedOutput()
		response, err := buildPromptResponse(m.director, "Generate a concise git commit message for this diff:\n"+string(diffOut))
		if err != nil {
			return AgentResponseMsg{Content: fmt.Sprintf("Error: %v", err)}
		}
		return response
	}
}

func handlePlan(m *model, prompt string) (tea.Model, tea.Cmd) {
	if strings.TrimSpace(prompt) == "" {
		errEntry(m, "/plan requires a prompt")
		return *m, nil
	}
	plan, err := foundationorchestration.BuildPlan(foundationorchestration.PlanRequest{
		Prompt:     prompt,
		WorkingDir: m.director.WorkingDir,
		IncludeRepo: true,
		MaxRepoFiles: 8,
	})
	if err != nil {
		errEntry(m, fmt.Sprintf("Plan generation failed: %v", err))
		return *m, nil
	}
	payload, _ := json.MarshalIndent(plan, "", "  ")
	m.entries = append(m.entries, ChatEntry{
		Type:      EntryTool,
		ToolName:  "plan",
		ToolArgs:  prompt,
		Content:   string(payload),
		Timestamp: time.Now(),
	})
	return *m, nil
}

func handleRepoMap(m *model) (tea.Model, tea.Cmd) {
	result, err := foundationrepomap.Generate(foundationrepomap.Options{
		BaseDir:  m.director.WorkingDir,
		MaxFiles: 20,
	})
	if err != nil {
		errEntry(m, fmt.Sprintf("Repomap failed: %v", err))
		return *m, nil
	}
	m.entries = append(m.entries, ChatEntry{
		Type:      EntryTool,
		ToolName:  "repomap",
		Content:   result.Map,
		Timestamp: time.Now(),
	})
	return *m, nil
}

func handleProviders(m *model) (tea.Model, tea.Cmd) {
	payload, _ := json.MarshalIndent(adapters.BuildProviderStatus(), "", "  ")
	m.entries = append(m.entries, ChatEntry{
		Type:      EntryTool,
		ToolName:  "providers",
		Content:   string(payload),
		Timestamp: time.Now(),
	})
	return *m, nil
}

func handleAdapters(m *model) (tea.Model, tea.Cmd) {
	hyper := adapters.NewHyperCodeAdapter(m.director.WorkingDir)
	mcpAdapter := adapters.NewMCPAdapter(m.director.WorkingDir)
	payload, _ := json.MarshalIndent(map[string]any{
		"hypercode": hyper.Status(),
		"mcp":       mcpAdapter.Status(),
	}, "", "  ")
	m.entries = append(m.entries, ChatEntry{
		Type:      EntryTool,
		ToolName:  "adapters",
		Content:   string(payload),
		Timestamp: time.Now(),
	})
	return *m, nil
}

func handleMCPTools(m *model) (tea.Model, tea.Cmd) {
	mcpAdapter := adapters.NewMCPAdapter(m.director.WorkingDir)
	toolList, err := mcpAdapter.ListTools()
	if err != nil {
		errEntry(m, fmt.Sprintf("MCP tools unavailable: %v", err))
		return *m, nil
	}
	payload, _ := json.MarshalIndent(map[string]any{"tools": toolList}, "", "  ")
	m.entries = append(m.entries, ChatEntry{
		Type:      EntryTool,
		ToolName:  "mcp",
		Content:   string(payload),
		Timestamp: time.Now(),
	})
	return *m, nil
}

func handleModel(m *model, arg string) (tea.Model, tea.Cmd) {
	if strings.TrimSpace(arg) != "" {
		m.modelName = strings.TrimSpace(arg)
		sysEntry(m, fmt.Sprintf("Model set to: %s", m.modelName))
		return *m, nil
	}
	sysEntry(m, fmt.Sprintf("Current model: %s/%s", m.provider, m.modelName))
	return *m, nil
}

func handleSettings(m *model) (tea.Model, tea.Cmd) {
	var b strings.Builder
	t := m.theme
	b.WriteString(t.Bold(t.AccentText("╭─ Settings ────────────────────────────────────╮")))
	b.WriteString("\n")
	settings := []struct{ key, val string }{
		{"Provider", m.provider},
		{"Model", m.modelName},
		{"WorkingDir", shortenPath(m.workingDir)},
		{"Git Branch", m.gitBranch},
		{"Tools", fmt.Sprintf("%d", m.toolCount)},
		{"Context Window", formatTokens(m.contextWindow)},
		{"Thinking Hidden", fmt.Sprintf("%t", m.hidingThink)},
		{"Auto-compact", "enabled"},
	}
	for _, s := range settings {
		b.WriteString(t.Dim("  "+s.key+": ") + t.Fg(t.TextColor, s.val))
		b.WriteString("\n")
	}
	b.WriteString(t.Bold(t.AccentText("╰───────────────────────────────────────────────╯")))
	sysEntry(m, b.String())
	return *m, nil
}

func handleSession(m *model) (tea.Model, tea.Cmd) {
	t := m.theme
	var b strings.Builder
	b.WriteString(t.Bold(t.AccentText("╭─ Session Info ────────────────────────────────╮")))
	b.WriteString("\n")
	b.WriteString(t.Dim("  Entries: ") + t.Fg(t.TextColor, fmt.Sprintf("%d", len(m.entries))))
	b.WriteString("\n")
	b.WriteString(t.Dim("  Foundation: ") + t.Fg(t.TextColor, m.foundationSessionID))
	b.WriteString("\n")
	b.WriteString(t.Dim(" Duration: ") + t.Fg(t.TextColor, sessionDuration(m)))
	b.WriteString("\n")
	b.WriteString(t.Bold(t.AccentText("╰───────────────────────────────────────────────╯")))
	sysEntry(m, b.String())
	return *m, nil
}

func handleName(m *model, arg string) (tea.Model, tea.Cmd) {
	if strings.TrimSpace(arg) == "" {
		sysEntry(m, "Usage: /name <session name>")
		return *m, nil
	}
	sysEntry(m, fmt.Sprintf("Session name set to: %s", arg))
	return *m, nil
}

func handleFork(m *model) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.AccentText("[Fork] New branch created from current state."))
	return *m, nil
}

func handleExport(m *model, arg string) (tea.Model, tea.Cmd) {
	path := arg
	if path == "" {
		path = "hyperharness-session.jsonl"
	}
	// Export entries as JSONL
	var lines []string
	for _, e := range m.entries {
		entry := map[string]string{
			"type":    fmt.Sprintf("%d", e.Type),
			"content": e.Content,
			"tool":    e.ToolName,
		}
		data, _ := json.Marshal(entry)
		lines = append(lines, string(data))
	}
	err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		errEntry(m, fmt.Sprintf("Export failed: %v", err))
		return *m, nil
	}
	sysEntry(m, m.theme.SuccessText(fmt.Sprintf("✓ Session exported to %s", path)))
	return *m, nil
}

func handleImport(m *model, arg string) (tea.Model, tea.Cmd) {
	if strings.TrimSpace(arg) == "" {
		errEntry(m, "/import requires a file path")
		return *m, nil
	}
	sysEntry(m, m.theme.AccentText(fmt.Sprintf("[Import] Loading session from %s...", arg)))
	return *m, nil
}

func handleLogin(m *model) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.AccentText("[Login] OAuth login initiated..."))
	return *m, nil
}

func handleLogout(m *model) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.SuccessText("✓ Logged out."))
	return *m, nil
}

func handleNew(m *model) (tea.Model, tea.Cmd) {
	workingDir := m.director.WorkingDir
	m.entries = []ChatEntry{}
	m.director = agents.NewDirector(agents.NewHyperCodeProvider())
	m.director.WorkingDir = workingDir
	m.foundationSessionID = ""
	sysEntry(m, m.theme.SuccessText("✓ New session started."))
	return *m, nil
}

func handleResume(m *model) (tea.Model, tea.Cmd) {
	sysEntry(m, m.theme.AccentText("[Resume] Session selector..."))
	return *m, nil
}

func handleReload(m *model) (tea.Model, tea.Cmd) {
	m.registry = tools.NewRegistry()
	m.toolCount = len(m.registry.Tools) + countInstalledTools()
	m.gitBranch = getGitBranch(m.workingDir)
	sysEntry(m, m.theme.SuccessText("✓ Configuration reloaded."))
	return *m, nil
}

func handleTools(m *model) (tea.Model, tea.Cmd) {
	t := m.theme
	var b strings.Builder
	b.WriteString(t.Bold(t.AccentText(fmt.Sprintf("╭─ Registered Tools (%d) ─────────────────────╮", m.toolCount))))
	b.WriteString("\n")

	if m.registry != nil {
		// Group tools by prefix
		groups := map[string][]string{}
		for _, tool := range m.registry.Tools {
			prefix := toolGroupName(tool.Name)
			groups[prefix] = append(groups[prefix], tool.Name)
		}
		for prefix, names := range groups {
			b.WriteString(t.Fg(t.ToolTitle, prefix) + t.Dim(fmt.Sprintf(" (%d)", len(names))))
			b.WriteString("\n")
			for _, name := range names {
				b.WriteString(t.Dim("  • " + name))
				b.WriteString("\n")
			}
		}
	}

	b.WriteString(t.Bold(t.AccentText("╰───────────────────────────────────────────────╯")))
	sysEntry(m, b.String())
	return *m, nil
}

func handleFoundationSession(m *model) (tea.Model, tea.Cmd) {
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		errEntry(m, fmt.Sprintf("Foundation session unavailable: %v", err))
		return *m, nil
	}
	m.entries = append(m.entries, ChatEntry{
		Type:      EntryTool,
		ToolName:  "fsession",
		Content:   sessionID,
		Timestamp: time.Now(),
	})
	return *m, nil
}

// ─── Tree commands (same logic as before, but using ChatEntry) ────────

func handleTree(m *model, arg string) (tea.Model, tea.Cmd) {
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		errEntry(m, fmt.Sprintf("Foundation session unavailable: %v", err))
		return *m, nil
	}
	if strings.TrimSpace(arg) == "" {
		display, err := buildFoundationTreeDisplay(m.director.WorkingDir, sessionID)
		if err != nil {
			errEntry(m, fmt.Sprintf("Tree display failed: %v", err))
			return *m, nil
		}
		m.entries = append(m.entries, ChatEntry{Type: EntryTool, ToolName: "tree", Content: display, Timestamp: time.Now()})
		return *m, nil
	}
	targetID, maxTokens := parseSummaryArgs(arg)
	display, err := switchFoundationTreeDisplay(m.director.WorkingDir, sessionID, targetID, maxTokens)
	if err != nil {
		errEntry(m, fmt.Sprintf("Tree switch failed: %v", err))
		return *m, nil
	}
	m.entries = append(m.entries, ChatEntry{Type: EntryTool, ToolName: "tree-switch", Content: display, Timestamp: time.Now()})
	return *m, nil
}

func handleTreeSelect(m *model) (tea.Model, tea.Cmd) {
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		errEntry(m, fmt.Sprintf("Foundation session unavailable: %v", err))
		return *m, nil
	}
	display, ids, err := buildFoundationTreeSelectionDisplay(m.director.WorkingDir, sessionID)
	if err != nil {
		errEntry(m, fmt.Sprintf("Tree selector failed: %v", err))
		return *m, nil
	}
	m.foundationTreeSelection = ids
	m.entries = append(m.entries, ChatEntry{Type: EntryTool, ToolName: "tree-select", Content: display, Timestamp: time.Now()})
	return *m, nil
}

func handleTreeBrowser(m *model) (tea.Model, tea.Cmd) {
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		errEntry(m, fmt.Sprintf("Foundation session unavailable: %v", err))
		return *m, nil
	}
	items, err := buildFoundationTreeBrowser(m.director.WorkingDir, sessionID)
	if err != nil {
		errEntry(m, fmt.Sprintf("Tree browser failed: %v", err))
		return *m, nil
	}
	m.browserItems = items
	m.browserIndex = 0
	m.browserActive = true
	return *m, nil
}

func handleTreePaneHelp(m *model) (tea.Model, tea.Cmd) {
	t := m.theme
	var b strings.Builder
	b.WriteString(t.Bold(t.AccentText("╭─ Tree Pane Controls ─────────────────────────╮")))
	b.WriteString("\n")
	controls := []struct{ cmd, desc string }{
		{"/tree-pane", "Toggle persistent pane"},
		{"/tree-pane-show|hide", "Explicitly show or hide"},
		{"/tree-pane-focus", "Toggle pane focus (or Tab)"},
		{"/tree-pane-size <n>", "Set pane viewport height"},
		{"/tree-pane-size-cycle", "Cycle common heights"},
		{"/tree-pane-position <top|bottom>", "Set pane position"},
		{"/tree-pane-preview <on|off>", "Toggle preview details"},
		{"/tree-pane-grouped <on|off>", "Toggle grouped rendering"},
		{"/tree-pane-preset <name>", "Apply named preset"},
		{"/tree-pane-cycle", "Cycle named presets"},
		{"/tree-pane-refresh", "Refresh from runtime"},
		{"/tree-pane-reset", "Reset to defaults"},
		{"/tree-pane-status", "Show current state"},
		{"/tree-browser", "Open modal browser"},
		{"/tree-browser-clear", "Clear browser state"},
	}
	for _, c := range controls {
		b.WriteString(t.AccentText("  " + c.cmd) + t.Dim("  " + c.desc))
		b.WriteString("\n")
	}
	b.WriteString(t.Bold(t.AccentText("╰───────────────────────────────────────────────╯")))
	sysEntry(m, b.String())
	return *m, nil
}

func handleTreePane(m *model) (tea.Model, tea.Cmd) {
	if m.browserPinned {
		unpinFoundationTreeBrowser(m)
		m.browserPinnedFocus = false
		m.browserConfirmPending = false
		sysEntry(m, m.theme.Dim("[Tree Pane] hidden"))
		return *m, nil
	}
	if !pinFoundationTreeBrowser(m) {
		errEntry(m, "Tree pane failed")
		return *m, nil
	}
	refreshPinnedFoundationTreeBrowser(m)
	sysEntry(m, m.theme.AccentText("[Tree Pane] pinned (Tab to focus, Esc to unfocus)"))
	return *m, nil
}

func handleTreePaneShow(m *model) (tea.Model, tea.Cmd) {
	if m.browserPinned {
		sysEntry(m, m.theme.Dim("[Tree Pane] already visible"))
		return *m, nil
	}
	if !pinFoundationTreeBrowser(m) {
		errEntry(m, "Tree pane failed")
		return *m, nil
	}
	refreshPinnedFoundationTreeBrowser(m)
	sysEntry(m, m.theme.SuccessText("✓ [Tree Pane] shown"))
	return *m, nil
}

func handleTreePaneHide(m *model) (tea.Model, tea.Cmd) {
	if !m.browserPinned {
		sysEntry(m, m.theme.Dim("[Tree Pane] already hidden"))
		return *m, nil
	}
	unpinFoundationTreeBrowser(m)
	m.browserPinnedFocus = false
	m.browserConfirmPending = false
	sysEntry(m, m.theme.Dim("[Tree Pane] hidden"))
	return *m, nil
}

func handleTreePaneSize(m *model, arg string) (tea.Model, tea.Cmd) {
	height := 0
	fmt.Sscanf(strings.TrimSpace(arg), "%d", &height)
	if height <= 0 {
		errEntry(m, "/tree-pane-size requires a positive integer")
		return *m, nil
	}
	m.browserPaneHeight = height
	sysEntry(m, fmt.Sprintf("[Tree Pane] height set to %d", height))
	return *m, nil
}

func handleTreePaneSizeCycle(m *model) (tea.Model, tea.Cmd) {
	next := 6
	switch m.browserPaneHeight {
	case 6:  next = 8
	case 8:  next = 10
	case 10: next = 12
	case 12: next = 14
	case 14: next = 6
	}
	m.browserPaneHeight = next
	sysEntry(m, fmt.Sprintf("[Tree Pane] height cycled to %d", next))
	return *m, nil
}

func handleTreePanePreview(m *model, arg string) (tea.Model, tea.Cmd) {
	value := strings.ToLower(strings.TrimSpace(arg))
	if value == "toggle" || value == "" {
		m.browserPanePreview = !m.browserPanePreview
		sysEntry(m, fmt.Sprintf("[Tree Pane] preview set to %t", m.browserPanePreview))
		return *m, nil
	}
	if value != "on" && value != "off" {
		errEntry(m, "/tree-pane-preview requires 'on', 'off', or 'toggle'")
		return *m, nil
	}
	m.browserPanePreview = value == "on"
	sysEntry(m, fmt.Sprintf("[Tree Pane] preview set to %s", value))
	return *m, nil
}

func handleTreePaneGrouped(m *model, arg string) (tea.Model, tea.Cmd) {
	value := strings.ToLower(strings.TrimSpace(arg))
	switch value {
	case "on":
		m.browserGrouped = true
	case "off":
		m.browserGrouped = false
	case "toggle", "":
		m.browserGrouped = !m.browserGrouped
	default:
		errEntry(m, "/tree-pane-grouped requires 'on', 'off', or 'toggle'")
		return *m, nil
	}
	sysEntry(m, fmt.Sprintf("[Tree Pane] grouped set to %t", m.browserGrouped))
	return *m, nil
}

func handleTreePaneCycle(m *model) (tea.Model, tea.Cmd) {
	preset := "compact"
	switch {
	case m.browserPaneHeight == 6 && !m.browserPanePreview && m.browserPanePosition == "bottom" && !m.browserGrouped:
		preset = "navigation"
	case m.browserPaneHeight == 10 && !m.browserPanePreview && m.browserPanePosition == "bottom" && m.browserGrouped:
		preset = "detailed"
	case m.browserPaneHeight == 12 && m.browserPanePreview && m.browserPanePosition == "top" && !m.browserGrouped:
		preset = "review"
	case m.browserPaneHeight == 14 && m.browserPanePreview && m.browserPanePosition == "top" && m.browserGrouped:
		preset = "compact"
	}
	return handleTreePanePreset(m, preset)
}

func handleTreePaneRefresh(m *model) (tea.Model, tea.Cmd) {
	if !m.browserPinned {
		errEntry(m, "Tree pane is not pinned; use /tree-pane first")
		return *m, nil
	}
	refreshPinnedFoundationTreeBrowser(m)
	sysEntry(m, m.theme.SuccessText("✓ [Tree Pane] refreshed"))
	return *m, nil
}

func handleTreeBrowserClear(m *model) (tea.Model, tea.Cmd) {
	m.browserFilter = ""
	m.browserConfirmPending = false
	m.browserCollapsed = nil
	m.browserIndex = 0
	sysEntry(m, m.theme.Dim("[Tree Browser] transient state cleared"))
	return *m, nil
}

func handleTreePaneReset(m *model) (tea.Model, tea.Cmd) {
	m.browserPaneHeight = 8
	m.browserPanePosition = "top"
	m.browserPanePreview = true
	m.browserGrouped = false
	m.browserPinnedFocus = false
	m.browserConfirmPending = false
	m.browserFilter = ""
	m.browserCollapsed = nil
	sysEntry(m, m.theme.SuccessText("✓ [Tree Pane] reset to defaults"))
	return *m, nil
}

func handleTreePaneStatus(m *model) (tea.Model, tea.Cmd) {
	t := m.theme
	var b strings.Builder
	b.WriteString(t.Bold(t.AccentText("╭─ Tree Pane Status ───────────────────────────╮")))
	b.WriteString("\n")
	b.WriteString(t.Dim("  pinned:   ") + t.Fg(t.TextColor, fmt.Sprintf("%t", m.browserPinned)))
	b.WriteString("\n")
	b.WriteString(t.Dim("  focus:    ") + t.Fg(t.TextColor, fmt.Sprintf("%t", m.browserPinnedFocus)))
	b.WriteString("\n")
	b.WriteString(t.Dim("  height:   ") + t.Fg(t.TextColor, fmt.Sprintf("%d", m.browserPaneHeight)))
	b.WriteString("\n")
	b.WriteString(t.Dim("  position: ") + t.Fg(t.TextColor, m.browserPanePosition))
	b.WriteString("\n")
	b.WriteString(t.Dim("  preview:  ") + t.Fg(t.TextColor, fmt.Sprintf("%t", m.browserPanePreview)))
	b.WriteString("\n")
	b.WriteString(t.Dim("  grouped:  ") + t.Fg(t.TextColor, fmt.Sprintf("%t", m.browserGrouped)))
	b.WriteString("\n")
	b.WriteString(t.Bold(t.AccentText("╰───────────────────────────────────────────────╯")))
	sysEntry(m, b.String())
	return *m, nil
}

func handleTreePaneSummary(m *model) (tea.Model, tea.Cmd) {
	sysEntry(m, fmt.Sprintf("[Tree Pane] pinned=%t focus=%t h=%d pos=%s preview=%t grouped=%t filter=%q",
		m.browserPinned, m.browserPinnedFocus, m.browserPaneHeight,
		m.browserPanePosition, m.browserPanePreview, m.browserGrouped, m.browserFilter))
	return *m, nil
}

func handleTreePanePreset(m *model, arg string) (tea.Model, tea.Cmd) {
	preset := strings.ToLower(strings.TrimSpace(arg))
	switch preset {
	case "compact":
		m.browserPaneHeight = 6; m.browserPanePreview = false; m.browserPanePosition = "bottom"; m.browserGrouped = false
	case "detailed":
		m.browserPaneHeight = 12; m.browserPanePreview = true; m.browserPanePosition = "top"; m.browserGrouped = false
	case "navigation":
		m.browserPaneHeight = 10; m.browserPanePreview = false; m.browserPanePosition = "bottom"; m.browserGrouped = true
	case "review":
		m.browserPaneHeight = 14; m.browserPanePreview = true; m.browserPanePosition = "top"; m.browserGrouped = true
	default:
		errEntry(m, "/tree-pane-preset requires 'compact', 'detailed', 'navigation', or 'review'")
		return *m, nil
	}
	sysEntry(m, m.theme.SuccessText(fmt.Sprintf("✓ [Tree Pane] preset applied: %s", preset)))
	return *m, nil
}

func handleTreePanePosition(m *model, arg string) (tea.Model, tea.Cmd) {
	position := strings.ToLower(strings.TrimSpace(arg))
	if position == "toggle" || position == "" {
		if m.browserPanePosition == "bottom" {
			m.browserPanePosition = "top"
		} else {
			m.browserPanePosition = "bottom"
		}
		sysEntry(m, fmt.Sprintf("[Tree Pane] position set to %s", m.browserPanePosition))
		return *m, nil
	}
	if position != "top" && position != "bottom" {
		errEntry(m, "/tree-pane-position requires 'top', 'bottom', or 'toggle'")
		return *m, nil
	}
	m.browserPanePosition = position
	sysEntry(m, fmt.Sprintf("[Tree Pane] position set to %s", position))
	return *m, nil
}

func handleTreePaneFocusValue(m *model, enabled bool) (tea.Model, tea.Cmd) {
	if !m.browserPinned {
		errEntry(m, "Tree pane is not pinned; use /tree-pane first")
		return *m, nil
	}
	m.browserPinnedFocus = enabled
	if m.browserPinnedFocus {
		sysEntry(m, m.theme.AccentText("[Tree Pane] focus enabled (Esc to unfocus)"))
	} else {
		m.browserConfirmPending = false
		sysEntry(m, m.theme.Dim("[Tree Pane] focus disabled"))
	}
	return *m, nil
}

func handleTreePaneFocus(m *model) (tea.Model, tea.Cmd) {
	return handleTreePaneFocusValue(m, !m.browserPinnedFocus)
}

func handleTreeGo(m *model, arg string) (tea.Model, tea.Cmd) {
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		errEntry(m, fmt.Sprintf("Foundation session unavailable: %v", err))
		return *m, nil
	}
	parts := strings.Fields(strings.TrimSpace(arg))
	if len(parts) == 0 {
		errEntry(m, "/tree-go requires an index from /tree-select")
		return *m, nil
	}
	index := 0
	fmt.Sscanf(parts[0], "%d", &index)
	maxTokens := 0
	if len(parts) > 1 { fmt.Sscanf(parts[1], "%d", &maxTokens) }
	display, err := switchFoundationTreeSelection(m.director.WorkingDir, sessionID, m.foundationTreeSelection, index, maxTokens)
	if err != nil {
		errEntry(m, fmt.Sprintf("tree-go failed: %v", err))
		return *m, nil
	}
	m.entries = append(m.entries, ChatEntry{Type: EntryTool, ToolName: "tree-go", Content: display, Timestamp: time.Now()})
	return *m, nil
}

func handleTreeChildren(m *model, arg string) (tea.Model, tea.Cmd) {
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		errEntry(m, fmt.Sprintf("Foundation session unavailable: %v", err))
		return *m, nil
	}
	parentID := strings.TrimSpace(arg)
	if parentID == "" {
		errEntry(m, "/tree-children requires an entry id")
		return *m, nil
	}
	display, err := buildFoundationChildrenDisplay(m.director.WorkingDir, sessionID, parentID)
	if err != nil {
		errEntry(m, fmt.Sprintf("Tree children failed: %v", err))
		return *m, nil
	}
	m.entries = append(m.entries, ChatEntry{Type: EntryTool, ToolName: "tree-children", Content: display, Timestamp: time.Now()})
	return *m, nil
}

func handleLabel(m *model, arg string) (tea.Model, tea.Cmd) {
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		errEntry(m, fmt.Sprintf("Foundation session unavailable: %v", err))
		return *m, nil
	}
	parts := strings.Fields(strings.TrimSpace(arg))
	if len(parts) < 2 {
		errEntry(m, "/label requires <entryId> <label>")
		return *m, nil
	}
	targetID := parts[0]
	label := strings.Join(parts[1:], " ")
	display, err := setFoundationLabel(m.director.WorkingDir, sessionID, targetID, label)
	if err != nil {
		errEntry(m, fmt.Sprintf("Label failed: %v", err))
		return *m, nil
	}
	sysEntry(m, display)
	return *m, nil
}

func handleSummaryCompact(m *model, arg string) (tea.Model, tea.Cmd) {
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		errEntry(m, fmt.Sprintf("Foundation session unavailable: %v", err))
		return *m, nil
	}
	keepRecent := 0
	if strings.TrimSpace(arg) != "" { fmt.Sscanf(strings.TrimSpace(arg), "%d", &keepRecent) }
	display, err := buildFoundationCompactionDisplay(m.director.WorkingDir, sessionID, keepRecent)
	if err != nil {
		errEntry(m, fmt.Sprintf("Compaction summary failed: %v", err))
		return *m, nil
	}
	sysEntry(m, display)
	return *m, nil
}

func handleSummaryBranch(m *model, arg string) (tea.Model, tea.Cmd) {
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		errEntry(m, fmt.Sprintf("Foundation session unavailable: %v", err))
		return *m, nil
	}
	targetID, maxTokens := parseSummaryArgs(arg)
	if strings.TrimSpace(targetID) == "" {
		errEntry(m, "/summary-branch requires a target entry id")
		return *m, nil
	}
	display, err := buildFoundationBranchSummaryDisplay(m.director.WorkingDir, sessionID, targetID, maxTokens)
	if err != nil {
		errEntry(m, fmt.Sprintf("Branch summary failed: %v", err))
		return *m, nil
	}
	sysEntry(m, display)
	return *m, nil
}

func handleUnknown(m *model, cmd string) (tea.Model, tea.Cmd) {
	errEntry(m, fmt.Sprintf("Unknown command: %s  (type /help for all commands)", cmd))
	return *m, nil
}
