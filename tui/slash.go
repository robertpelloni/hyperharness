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
	case "/tree-pane-size":
		return handleTreePaneSize(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/tree-pane-size")))
	case "/tree-pane-position":
		return handleTreePanePosition(m, strings.TrimSpace(strings.TrimPrefix(cmd, "/tree-pane-position")))
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
  /tree-pane-size <n> - Set the persistent tree pane viewport height
  /tree-pane-position <top|bottom> - Set the persistent tree pane position
  /tree-pane-focus - Toggle keyboard focus for the pinned tree pane
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

func handleTreePanePosition(m *model, arg string) (tea.Model, tea.Cmd) {
	m.loading = false
	position := strings.ToLower(strings.TrimSpace(arg))
	if position != "top" && position != "bottom" {
		m.history = append(m.history, "[Error] /tree-pane-position requires 'top' or 'bottom'")
		return *m, nil
	}
	m.browserPanePosition = position
	m.history = append(m.history, fmt.Sprintf("[Foundation Tree Pane] position set to %s", position))
	return *m, nil
}

func handleTreePaneFocus(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	if !m.browserPinned {
		m.history = append(m.history, "[Error] tree pane is not pinned; use /tree-pane first")
		return *m, nil
	}
	m.browserPinnedFocus = !m.browserPinnedFocus
	if m.browserPinnedFocus {
		m.history = append(m.history, "[Foundation Tree Pane] focus enabled")
	} else {
		m.browserConfirmPending = false
		m.history = append(m.history, "[Foundation Tree Pane] focus disabled")
	}
	return *m, nil
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
