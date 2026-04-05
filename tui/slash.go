package tui

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/robertpelloni/hypercode/agents"
	"github.com/robertpelloni/hypercode/foundation/adapters"
	foundationorchestration "github.com/robertpelloni/hypercode/foundation/orchestration"
	foundationrepomap "github.com/robertpelloni/hypercode/foundation/repomap"
)

// ProcessSlashCommand mimics Claude Code's native terminal interception primitives.
// If an input begins with '/', it bypasses the LLM and executes directly as a local primitive.
func ProcessSlashCommand(cmd string, m *model) (tea.Model, tea.Cmd) {
	cmd = strings.TrimSpace(cmd)
	parts := strings.Split(cmd, " ")

	switch parts[0] {
	case "/help":
		return handleHelp(m)
	case "/clear":
		return handleClear(m)
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
	case "/summary-compact":
		return handleSummaryCompact(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/summary-compact")))
	case "/summary-branch":
		return handleSummaryBranch(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/summary-branch")))
	case "/tree":
		return handleTree(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/tree")))
	case "/tree-select":
		return handleTreeSelect(m)
	case "/tree-browser":
		return handleTreeBrowser(m)
	case "/tree-pane":
		return handleTreePane(m)
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
	case "/tree-pane-focus-toggle":
		return handleTreePaneFocus(m)
	case "/tree-pane-focus":
		return handleTreePaneFocus(m)
	case "/tree-go":
		return handleTreeGo(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/tree-go")))
	case "/tree-children":
		return handleTreeChildren(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/tree-children")))
	case "/label":
		return handleLabel(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/label")))
	case "/fsession":
		return handleFoundationSession(m)
	case "/exit", "/quit":
		return *m, tea.Quit
	default:
		return handleUnknown(m, parts[0])
	}
}

func handleHelp(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	// Simulate glamour markdown output locally bypassing the Director
	m.history = append(m.history, `[System] Parity Slash Commands:
  /help      - This menu
  /clear     - Wipes contextual agent memory
  /commit    - Autonomously generates a standard Git commit
  /plan      - Build a foundation-backed orchestration plan
  /repomap   - Generate a foundation-backed repo map
  /providers - Show provider visibility and defaults
  /adapters  - Show HyperCode/Borg + MCP adapter status
  /mcp       - Show adapter-backed MCP tool hints
  /fsession  - Show or create the active foundation session
  /tree      - Show the active foundation session tree
  /tree <targetEntryId> [maxTokens] - Switch to a target entry and preserve abandoned branch context
  /tree-select - Show a numbered entry selector for the active foundation session
  /tree-browser - Open a cursor-driven tree browser for the active foundation session
  /tree-pane - Toggle a persistent tree pane while continuing normal prompt interaction
  /tree-pane-show - Explicitly show/pin the persistent tree pane
  /tree-pane-hide - Explicitly hide the persistent tree pane
  /tree-pane-size <n> - Set the persistent tree pane viewport height
  /tree-pane-size-cycle - Quickly cycle common persistent tree pane heights
  /tree-pane-preview <on|off> - Toggle preview details inside the persistent tree pane
  /tree-pane-preview-on - Explicitly enable preview details inside the persistent tree pane
  /tree-pane-preview-off - Explicitly disable preview details inside the persistent tree pane
  /tree-pane-preview-toggle - Quickly toggle preview details for the persistent tree pane
  /tree-pane-grouped <on|off|toggle> - Control grouped rendering for the persistent tree pane
  /tree-pane-grouped-on - Explicitly enable grouped rendering for the persistent tree pane
  /tree-pane-grouped-off - Explicitly disable grouped rendering for the persistent tree pane
  /tree-pane-grouped-toggle - Quickly toggle grouped rendering for the persistent tree pane
  /tree-pane-cycle - Cycle through common pane presets
  /tree-pane-refresh - Manually refresh the persistent tree pane from canonical runtime state
  /tree-browser-clear - Clear transient browser state like filter/collapse/confirm
  /tree-pane-reset - Reset pane configuration to defaults
  /tree-pane-status - Show the current persistent pane configuration
  /tree-pane-preset <compact|detailed|navigation|review> - Apply a named pane layout preset
  /tree-pane-compact - Apply the compact pane preset
  /tree-pane-detailed - Apply the detailed pane preset
  /tree-pane-navigation - Apply the navigation pane preset
  /tree-pane-review - Apply the review pane preset
  /tree-pane-position <top|bottom> - Set the persistent tree pane position
  /tree-pane-top - Explicitly place the persistent tree pane above the main flow
  /tree-pane-bottom - Explicitly place the persistent tree pane below the main flow
  /tree-pane-position-toggle - Quickly toggle pane position between top and bottom
  /tree-pane-focus - Toggle keyboard focus for the pinned tree pane
  /tree-pane-focus-toggle - Quick alias for toggling keyboard focus on the pinned tree pane
  /tree-pane-focus-on - Explicitly enable keyboard focus for the pinned tree pane
  /tree-pane-focus-off - Explicitly disable keyboard focus for the pinned tree pane
  /tree-go <index> [maxTokens] - Switch to an indexed entry from /tree-select
  /tree-children <entryId> - Show direct child branches for an entry
  /label <entryId> <label> - Set a label on an entry (or clear with empty label unsupported in slash)
  /summary-compact [keepRecentTokens] - Generate a native compaction summary for the active foundation session
  /summary-branch <targetEntryId> [maxTokens] - Generate a native branch summary toward a target entry
  /exit      - Closes hypercode`)
	return *m, nil
}

func handleClear(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	workingDir := m.director.WorkingDir
	m.history = []string{
		"[System] Memory crystal wiped. Context reset to null.",
	}
	m.director = agents.NewDirector(&agents.DefaultProvider{})
	m.director.WorkingDir = workingDir
	return *m, nil
}

func handleCommit(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	m.history = append(m.history, "[System] Extracting git diff to generate native algorithmic commit message...")
	// We'd tap native.go tools here to run 'git diff' and pass it to Director exclusively for formatting
	return *m, nil
}

func handlePlan(m *model, prompt string) (tea.Model, tea.Cmd) {
	m.loading = false
	if strings.TrimSpace(prompt) == "" {
		m.history = append(m.history, "[Error] /plan requires a prompt")
		return *m, nil
	}
	plan, err := foundationorchestration.BuildPlan(foundationorchestration.PlanRequest{Prompt: prompt, WorkingDir: m.director.WorkingDir, IncludeRepo: true, MaxRepoFiles: 8})
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] plan generation failed: %v", err))
		return *m, nil
	}
	payload, _ := json.MarshalIndent(plan, "", "  ")
	m.history = append(m.history, "[Foundation Plan]\n"+string(payload))
	return *m, nil
}

func handleRepoMap(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	result, err := foundationrepomap.Generate(foundationrepomap.Options{BaseDir: m.director.WorkingDir, MaxFiles: 12})
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] repomap generation failed: %v", err))
		return *m, nil
	}
	m.history = append(m.history, "[Foundation RepoMap]\n"+result.Map)
	return *m, nil
}

func handleProviders(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	payload, _ := json.MarshalIndent(adapters.BuildProviderStatus(), "", "  ")
	m.history = append(m.history, "[Foundation Providers]\n"+string(payload))
	return *m, nil
}

func handleAdapters(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	hyper := adapters.NewHyperCodeAdapter(m.director.WorkingDir)
	mcpAdapter := adapters.NewMCPAdapter(m.director.WorkingDir)
	payload, _ := json.MarshalIndent(map[string]any{
		"hypercode": hyper.Status(),
		"mcp":       mcpAdapter.Status(),
	}, "", "  ")
	m.history = append(m.history, "[Foundation Adapters]\n"+string(payload))
	return *m, nil
}

func handleMCPTools(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	mcpAdapter := adapters.NewMCPAdapter(m.director.WorkingDir)
	tools, err := mcpAdapter.ListTools()
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] mcp tools unavailable: %v", err))
		return *m, nil
	}
	payload, _ := json.MarshalIndent(map[string]any{"tools": tools}, "", "  ")
	m.history = append(m.history, "[Foundation MCP Tools]\n"+string(payload))
	return *m, nil
}

func handleFoundationSession(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] foundation session unavailable: %v", err))
		return *m, nil
	}
	m.history = append(m.history, fmt.Sprintf("[Foundation Session]\n%s", sessionID))
	return *m, nil
}

func handleSummaryCompact(m *model, arg string) (tea.Model, tea.Cmd) {
	m.loading = false
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] foundation session unavailable: %v", err))
		return *m, nil
	}
	keepRecent := 0
	if strings.TrimSpace(arg) != "" {
		fmt.Sscanf(strings.TrimSpace(arg), "%d", &keepRecent)
	}
	display, err := buildFoundationCompactionDisplay(m.director.WorkingDir, sessionID, keepRecent)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] compaction summary failed: %v", err))
		return *m, nil
	}
	m.history = append(m.history, display)
	return *m, nil
}

func handleSummaryBranch(m *model, arg string) (tea.Model, tea.Cmd) {
	m.loading = false
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] foundation session unavailable: %v", err))
		return *m, nil
	}
	targetID, maxTokens := parseSummaryArgs(arg)
	if strings.TrimSpace(targetID) == "" {
		m.history = append(m.history, "[Error] /summary-branch requires a target entry id")
		return *m, nil
	}
	display, err := buildFoundationBranchSummaryDisplay(m.director.WorkingDir, sessionID, targetID, maxTokens)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] branch summary failed: %v", err))
		return *m, nil
	}
	m.history = append(m.history, display)
	return *m, nil
}

func handleTree(m *model, arg string) (tea.Model, tea.Cmd) {
	m.loading = false
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] foundation session unavailable: %v", err))
		return *m, nil
	}
	if strings.TrimSpace(arg) == "" {
		display, err := buildFoundationTreeDisplay(m.director.WorkingDir, sessionID)
		if err != nil {
			m.history = append(m.history, fmt.Sprintf("[Error] tree display failed: %v", err))
			return *m, nil
		}
		m.history = append(m.history, display)
		return *m, nil
	}
	targetID, maxTokens := parseSummaryArgs(arg)
	display, err := switchFoundationTreeDisplay(m.director.WorkingDir, sessionID, targetID, maxTokens)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] tree switch failed: %v", err))
		return *m, nil
	}
	m.history = append(m.history, display)
	return *m, nil
}

func handleTreeSelect(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] foundation session unavailable: %v", err))
		return *m, nil
	}
	display, ids, err := buildFoundationTreeSelectionDisplay(m.director.WorkingDir, sessionID)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] tree selector failed: %v", err))
		return *m, nil
	}
	m.foundationTreeSelection = ids
	m.history = append(m.history, display)
	return *m, nil
}

func handleTreeBrowser(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] foundation session unavailable: %v", err))
		return *m, nil
	}
	items, err := buildFoundationTreeBrowser(m.director.WorkingDir, sessionID)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] tree browser failed: %v", err))
		return *m, nil
	}
	m.browserItems = items
	m.browserIndex = 0
	m.browserActive = true
	return *m, nil
}

func handleTreePane(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	if m.browserPinned {
		unpinFoundationTreeBrowser(m)
		m.browserPinnedFocus = false
		m.browserConfirmPending = false
		m.history = append(m.history, "[Foundation Tree Pane] hidden")
		return *m, nil
	}
	if err := pinFoundationTreeBrowser(m); err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] tree pane failed: %v", err))
		return *m, nil
	}
	m.history = append(m.history, "[Foundation Tree Pane] pinned")
	return *m, nil
}

func handleTreePaneShow(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	if m.browserPinned {
		m.history = append(m.history, "[Foundation Tree Pane] already visible")
		return *m, nil
	}
	if err := pinFoundationTreeBrowser(m); err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] tree pane failed: %v", err))
		return *m, nil
	}
	m.history = append(m.history, "[Foundation Tree Pane] shown")
	return *m, nil
}

func handleTreePaneHide(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	if !m.browserPinned {
		m.history = append(m.history, "[Foundation Tree Pane] already hidden")
		return *m, nil
	}
	unpinFoundationTreeBrowser(m)
	m.browserPinnedFocus = false
	m.browserConfirmPending = false
	m.history = append(m.history, "[Foundation Tree Pane] hidden")
	return *m, nil
}

func handleTreePaneSize(m *model, arg string) (tea.Model, tea.Cmd) {
	m.loading = false
	height := 0
	fmt.Sscanf(strings.TrimSpace(arg), "%d", &height)
	if height <= 0 {
		m.history = append(m.history, "[Error] /tree-pane-size requires a positive integer")
		return *m, nil
	}
	m.browserPaneHeight = height
	m.history = append(m.history, fmt.Sprintf("[Foundation Tree Pane] height set to %d", height))
	return *m, nil
}

func handleTreePaneSizeCycle(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	next := 6
	switch m.browserPaneHeight {
	case 6:
		next = 8
	case 8:
		next = 10
	case 10:
		next = 12
	case 12:
		next = 14
	case 14:
		next = 6
	}
	m.browserPaneHeight = next
	m.history = append(m.history, fmt.Sprintf("[Foundation Tree Pane] height cycled to %d", next))
	return *m, nil
}

func handleTreePanePreview(m *model, arg string) (tea.Model, tea.Cmd) {
	m.loading = false
	value := strings.ToLower(strings.TrimSpace(arg))
	if value == "toggle" || value == "" {
		m.browserPanePreview = !m.browserPanePreview
		m.history = append(m.history, fmt.Sprintf("[Foundation Tree Pane] preview set to %t", m.browserPanePreview))
		return *m, nil
	}
	if value != "on" && value != "off" {
		m.history = append(m.history, "[Error] /tree-pane-preview requires 'on', 'off', or use /tree-pane-preview-toggle")
		return *m, nil
	}
	m.browserPanePreview = value == "on"
	m.history = append(m.history, fmt.Sprintf("[Foundation Tree Pane] preview set to %s", value))
	return *m, nil
}

func handleTreePaneGrouped(m *model, arg string) (tea.Model, tea.Cmd) {
	m.loading = false
	value := strings.ToLower(strings.TrimSpace(arg))
	switch value {
	case "on":
		m.browserGrouped = true
	case "off":
		m.browserGrouped = false
	case "toggle", "":
		m.browserGrouped = !m.browserGrouped
	default:
		m.history = append(m.history, "[Error] /tree-pane-grouped requires 'on', 'off', or 'toggle'")
		return *m, nil
	}
	m.history = append(m.history, fmt.Sprintf("[Foundation Tree Pane] grouped set to %t", m.browserGrouped))
	return *m, nil
}

func handleTreePaneCycle(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
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
	m.loading = false
	if !m.browserPinned {
		m.history = append(m.history, "[Error] tree pane is not pinned; use /tree-pane first")
		return *m, nil
	}
	refreshPinnedFoundationTreeBrowser(m)
	m.history = append(m.history, "[Foundation Tree Pane] refreshed")
	return *m, nil
}

func handleTreeBrowserClear(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	m.browserFilter = ""
	m.browserConfirmPending = false
	m.browserCollapsed = nil
	m.browserIndex = 0
	m.history = append(m.history, "[Foundation Tree Browser] transient state cleared")
	return *m, nil
}

func handleTreePaneReset(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	m.browserPaneHeight = 8
	m.browserPanePosition = "top"
	m.browserPanePreview = true
	m.browserGrouped = false
	m.browserPinnedFocus = false
	m.browserConfirmPending = false
	m.browserFilter = ""
	m.browserCollapsed = nil
	m.history = append(m.history, "[Foundation Tree Pane] reset to defaults")
	return *m, nil
}

func handleTreePaneStatus(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	m.history = append(m.history, fmt.Sprintf("[Foundation Tree Pane Status]\npinned=%t\nfocus=%t\nheight=%d\nposition=%s\npreview=%t\ngrouped=%t\nfilter=%q", m.browserPinned, m.browserPinnedFocus, m.browserPaneHeight, m.browserPanePosition, m.browserPanePreview, m.browserGrouped, m.browserFilter))
	return *m, nil
}

func handleTreePanePreset(m *model, arg string) (tea.Model, tea.Cmd) {
	m.loading = false
	preset := strings.ToLower(strings.TrimSpace(arg))
	switch preset {
	case "compact":
		m.browserPaneHeight = 6
		m.browserPanePreview = false
		m.browserPanePosition = "bottom"
		m.browserGrouped = false
		m.history = append(m.history, "[Foundation Tree Pane] preset applied: compact")
	case "detailed":
		m.browserPaneHeight = 12
		m.browserPanePreview = true
		m.browserPanePosition = "top"
		m.browserGrouped = false
		m.history = append(m.history, "[Foundation Tree Pane] preset applied: detailed")
	case "navigation":
		m.browserPaneHeight = 10
		m.browserPanePreview = false
		m.browserPanePosition = "bottom"
		m.browserGrouped = true
		m.history = append(m.history, "[Foundation Tree Pane] preset applied: navigation")
	case "review":
		m.browserPaneHeight = 14
		m.browserPanePreview = true
		m.browserPanePosition = "top"
		m.browserGrouped = true
		m.history = append(m.history, "[Foundation Tree Pane] preset applied: review")
	default:
		m.history = append(m.history, "[Error] /tree-pane-preset requires 'compact', 'detailed', 'navigation', or 'review'")
	}
	return *m, nil
}

func handleTreePanePosition(m *model, arg string) (tea.Model, tea.Cmd) {
	m.loading = false
	position := strings.ToLower(strings.TrimSpace(arg))
	if position == "toggle" || position == "" {
		if strings.TrimSpace(m.browserPanePosition) == "bottom" {
			m.browserPanePosition = "top"
		} else {
			m.browserPanePosition = "bottom"
		}
		m.history = append(m.history, fmt.Sprintf("[Foundation Tree Pane] position set to %s", m.browserPanePosition))
		return *m, nil
	}
	if position != "top" && position != "bottom" {
		m.history = append(m.history, "[Error] /tree-pane-position requires 'top', 'bottom', or use /tree-pane-position-toggle")
		return *m, nil
	}
	m.browserPanePosition = position
	m.history = append(m.history, fmt.Sprintf("[Foundation Tree Pane] position set to %s", position))
	return *m, nil
}

func handleTreePaneFocusValue(m *model, enabled bool) (tea.Model, tea.Cmd) {
	m.loading = false
	if !m.browserPinned {
		m.history = append(m.history, "[Error] tree pane is not pinned; use /tree-pane first")
		return *m, nil
	}
	m.browserPinnedFocus = enabled
	if m.browserPinnedFocus {
		m.history = append(m.history, "[Foundation Tree Pane] focus enabled")
	} else {
		m.browserConfirmPending = false
		m.history = append(m.history, "[Foundation Tree Pane] focus disabled")
	}
	return *m, nil
}

func handleTreePaneFocus(m *model) (tea.Model, tea.Cmd) {
	return handleTreePaneFocusValue(m, !m.browserPinnedFocus)
}

func handleTreeGo(m *model, arg string) (tea.Model, tea.Cmd) {
	m.loading = false
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] foundation session unavailable: %v", err))
		return *m, nil
	}
	parts := strings.Fields(strings.TrimSpace(arg))
	if len(parts) == 0 {
		m.history = append(m.history, "[Error] /tree-go requires an index from /tree-select")
		return *m, nil
	}
	index := 0
	fmt.Sscanf(parts[0], "%d", &index)
	maxTokens := 0
	if len(parts) > 1 {
		fmt.Sscanf(parts[1], "%d", &maxTokens)
	}
	display, err := switchFoundationTreeSelection(m.director.WorkingDir, sessionID, m.foundationTreeSelection, index, maxTokens)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] tree-go failed: %v", err))
		return *m, nil
	}
	m.history = append(m.history, display)
	return *m, nil
}

func handleTreeChildren(m *model, arg string) (tea.Model, tea.Cmd) {
	m.loading = false
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] foundation session unavailable: %v", err))
		return *m, nil
	}
	parentID := strings.TrimSpace(arg)
	if parentID == "" {
		m.history = append(m.history, "[Error] /tree-children requires an entry id")
		return *m, nil
	}
	display, err := buildFoundationChildrenDisplay(m.director.WorkingDir, sessionID, parentID)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] tree children failed: %v", err))
		return *m, nil
	}
	m.history = append(m.history, display)
	return *m, nil
}

func handleLabel(m *model, arg string) (tea.Model, tea.Cmd) {
	m.loading = false
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] foundation session unavailable: %v", err))
		return *m, nil
	}
	parts := strings.Fields(strings.TrimSpace(arg))
	if len(parts) < 2 {
		m.history = append(m.history, "[Error] /label requires <entryId> <label>")
		return *m, nil
	}
	targetID := parts[0]
	label := strings.Join(parts[1:], " ")
	display, err := setFoundationLabel(m.director.WorkingDir, sessionID, targetID, label)
	if err != nil {
		m.history = append(m.history, fmt.Sprintf("[Error] label failed: %v", err))
		return *m, nil
	}
	m.history = append(m.history, display)
	return *m, nil
}

func handleUnknown(m *model, cmd string) (tea.Model, tea.Cmd) {
	m.loading = false
	m.history = append(m.history, fmt.Sprintf("[Error] Unknown primitive: %s", cmd))
	return *m, nil
}
