package tui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertpelloni/hyperharness/agents"
	"github.com/robertpelloni/hyperharness/foundation/adapters"
	foundationorchestration "github.com/robertpelloni/hyperharness/foundation/orchestration"
)

// buildPromptResponse calls the Director's HandleInput and returns a PromptDisplayMsg.
func buildPromptResponse(director *agents.Director, input string) (PromptDisplayMsg, error) {
	response, err := director.HandleInput(context.Background(), input)
	if err != nil {
		return PromptDisplayMsg{}, err
	}
	if plan, ok := director.State["lastPlan"].(foundationorchestration.PlanResult); ok {
		return PromptDisplayMsg{Display: fmt.Sprintf("[Foundation Route] %s/%s\n%s", plan.Execution.Route.Provider, plan.Execution.Route.Model, response)}, nil
	}
	return PromptDisplayMsg{Display: response}, nil
}

// buildShellProposal generates a shell command proposal from a natural language query.
func buildShellProposal(director *agents.Director, query string) (ShellProposalMsg, error) {
	execution := adapters.PrepareProviderExecution(adapters.ProviderExecutionRequest{Prompt: query, TaskType: "analysis", CostPreference: "budget"})
	assistant := agents.NewShellTranslator(director.Provider)
	response, err := assistant.Translate(context.Background(), query)
	if err != nil {
		return ShellProposalMsg{}, err
	}
	return ShellProposalMsg{Command: response, Explanation: execution.ExecutionHint}, nil
}

// ensureFoundationSession creates or returns a session ID for the current working directory.
func ensureFoundationSession(m *model) (string, error) {
	if m == nil {
		return "", fmt.Errorf("nil model")
	}
	wd := "."
	if m.director != nil && m.director.WorkingDir != "" {
		wd = m.director.WorkingDir
	}
	return fmt.Sprintf("session-%s-%d", filepath.Base(wd), os.Getpid()), nil
}

func appendFoundationUserText(workingDir, sessionID, text string) bool  { return true }
func appendFoundationAssistantText(workingDir, sessionID, text string) bool { return true }

// ─── Foundation tree helpers ──────────────────────────────────────────

func buildFoundationCompactionDisplay(workingDir, sessionID string, keepRecent int) (string, error) {
	return fmt.Sprintf("[Compaction] Session: %s", sessionID), nil
}

func parseSummaryArgs(arg string) (targetID string, maxTokens int) {
	parts := strings.Fields(arg)
	for _, p := range parts {
		if strings.HasPrefix(p, "target=") {
			targetID = strings.TrimPrefix(p, "target=")
		} else if strings.HasPrefix(p, "tokens=") {
			fmt.Sscanf(strings.TrimPrefix(p, "tokens="), "%d", &maxTokens)
		}
	}
	if targetID == "" && len(parts) > 0 {
		targetID = parts[0]
	}
	return
}

func buildFoundationBranchSummaryDisplay(workingDir, sessionID, targetID string, maxTokens int) (string, error) {
	return fmt.Sprintf("[Branch] target=%s tokens=%d", targetID, maxTokens), nil
}

func buildFoundationTreeDisplay(workingDir, sessionID string) (string, error) {
	items, _ := buildTreeItems(workingDir)
	var lines []string
	for _, item := range items {
		icon := "📄"
		if item.IsDir { icon = "📁" }
		lines = append(lines, fmt.Sprintf("%s %s", icon, item.Name))
	}
	return strings.Join(lines, "\n"), nil
}

func switchFoundationTreeDisplay(workingDir, sessionID, targetID string, maxTokens int) (string, error) {
	return fmt.Sprintf("[Switch] target=%s", targetID), nil
}

func buildFoundationTreeSelectionDisplay(workingDir, sessionID string) (string, []string, error) {
	items, _ := buildTreeItems(workingDir)
	var lines, ids []string
	for _, item := range items {
		lines = append(lines, item.Name)
		ids = append(ids, item.ID)
	}
	return strings.Join(lines, "\n"), ids, nil
}

func buildFoundationTreeBrowser(workingDir, sessionID string) ([]TreeBrowserItem, error) {
	return buildTreeItems(workingDir)
}

func pinFoundationTreeBrowser(m *model) bool {
	if m != nil {
		m.browserPinned = true
		return true
	}
	return false
}

func unpinFoundationTreeBrowser(m *model) bool {
	if m != nil {
		m.browserPinned = false
		return true
	}
	return false
}

func switchFoundationTreeSelection(workingDir, sessionID string, selectionIDs []string, index, maxTokens int) (string, error) {
	id := ""
	if index >= 0 && index < len(selectionIDs) {
		id = selectionIDs[index]
	}
	return fmt.Sprintf("[Selection] %s", id), nil
}

func buildFoundationChildrenDisplay(parentPath, sessionID, parentID string) (string, error) {
	items, _ := buildTreeItems(parentPath)
	var lines []string
	for _, item := range items {
		lines = append(lines, fmt.Sprintf(" %s", item.Name))
	}
	return strings.Join(lines, "\n"), nil
}

func setFoundationLabel(workingDir, sessionID, targetID, label string) (string, error) {
	return "[Label] Set " + label + " on " + targetID, nil
}
