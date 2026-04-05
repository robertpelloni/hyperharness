package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/robertpelloni/hypercode/agents"
	"github.com/robertpelloni/hypercode/foundation/adapters"
	foundationorchestration "github.com/robertpelloni/hypercode/foundation/orchestration"
	foundationpi "github.com/robertpelloni/hypercode/foundation/pi"
)

type PromptDisplayMsg struct {
	Display string
}

type TreeBrowserItem struct {
	ID                 string
	Kind               string
	Label              string
	Preview            string
	IsLeaf             bool
	SummaryEntries     int
	CommonAncestorID   string
	ReadFilesCount     int
	ModifiedFilesCount int
}

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

func ensureFoundationSession(m *model) (string, error) {
	if m.foundationSessionID != "" {
		return m.foundationSessionID, nil
	}
	runtime := foundationpi.NewRuntime(m.director.WorkingDir, nil)
	session, err := runtime.CreateSession("tui-session")
	if err != nil {
		return "", err
	}
	m.foundationSessionID = session.Metadata.SessionID
	return m.foundationSessionID, nil
}

func appendFoundationUserText(cwd, sessionID, text string) error {
	runtime := foundationpi.NewRuntime(cwd, nil)
	_, err := runtime.AppendUserText(sessionID, text)
	return err
}

func appendFoundationAssistantText(cwd, sessionID, text string) error {
	store := foundationpi.DefaultSessionStore(cwd)
	_, err := store.AppendEntry(sessionID, foundationpi.SessionEntry{Kind: "message", Role: "assistant", Text: text})
	return err
}

func buildFoundationCompactionDisplay(cwd, sessionID string, keepRecentTokens int) (string, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	_, summary, err := runtime.CompactWithGeneratedSummary(context.Background(), sessionID, keepRecentTokens, nil, nil)
	if err != nil {
		return "", err
	}
	return "[Foundation Compaction Summary]\n" + summary, nil
}

func buildFoundationBranchSummaryDisplay(cwd, sessionID, targetID string, maxTokens int) (string, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	_, summary, err := runtime.BranchWithGeneratedSummary(context.Background(), sessionID, targetID, maxTokens, nil, nil)
	if err != nil {
		return "", err
	}
	return "[Foundation Branch Summary]\n" + summary, nil
}

func parseSummaryArgs(input string) (string, int) {
	fields := strings.Fields(strings.TrimSpace(input))
	if len(fields) == 0 {
		return "", 0
	}
	target := fields[0]
	maxTokens := 0
	if len(fields) > 1 {
		fmt.Sscanf(fields[1], "%d", &maxTokens)
	}
	return target, maxTokens
}

func buildFoundationTreeDisplay(cwd, sessionID string) (string, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	session, err := runtime.LoadSession(sessionID)
	if err != nil {
		return "", err
	}
	leafID, err := runtime.GetLeafID(sessionID)
	if err != nil {
		return "", err
	}
	sessionName, _ := runtime.GetSessionName(sessionID)
	var b strings.Builder
	b.WriteString("[Foundation Tree]\n")
	b.WriteString(fmt.Sprintf("session=%s", sessionID))
	if sessionName != "" {
		b.WriteString(fmt.Sprintf(" name=%q", sessionName))
	}
	b.WriteString(fmt.Sprintf(" leaf=%s\n", leafID))
	for _, entry := range session.Entries {
		marker := " "
		if entry.ID == leafID {
			marker = "*"
		}
		preview := entry.Text
		if strings.TrimSpace(preview) == "" {
			preview = entry.ToolName
		}
		if len(preview) > 60 {
			preview = preview[:60] + "..."
		}
		children, _ := runtime.GetChildren(sessionID, entry.ID)
		label, _ := runtime.GetLabel(sessionID, entry.ID)
		labelSuffix := ""
		if strings.TrimSpace(label) != "" {
			labelSuffix = fmt.Sprintf(" label=%q", label)
		}
		childSuffix := ""
		if len(children) > 0 {
			childSuffix = fmt.Sprintf(" children=%d", len(children))
		}
		b.WriteString(fmt.Sprintf("%s %s <- %s [%s]%s%s %s\n", marker, entry.ID, entry.ParentID, entry.Kind, labelSuffix, childSuffix, preview))
	}
	return strings.TrimSpace(b.String()), nil
}

func switchFoundationTreeDisplay(cwd, sessionID, targetID string, maxTokens int) (string, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	currentLeaf, err := runtime.GetLeafID(sessionID)
	if err != nil {
		return "", err
	}
	if targetID == "" {
		return "", fmt.Errorf("target entry id is required")
	}
	if currentLeaf == targetID {
		return fmt.Sprintf("[Foundation Tree]\nAlready on branch leaf %s", targetID), nil
	}
	_, summary, err := runtime.BranchWithGeneratedSummary(context.Background(), sessionID, targetID, maxTokens, nil, nil)
	if err != nil {
		return "", err
	}
	return "[Foundation Tree Switch]\n" + summary, nil
}

func buildFoundationChildrenDisplay(cwd, sessionID, parentID string) (string, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	children, err := runtime.GetChildren(sessionID, parentID)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString("[Foundation Tree Children]\n")
	b.WriteString(fmt.Sprintf("parent=%s\n", parentID))
	for _, entry := range children {
		preview := entry.Text
		if strings.TrimSpace(preview) == "" {
			preview = entry.ToolName
		}
		if len(preview) > 60 {
			preview = preview[:60] + "..."
		}
		label, _ := runtime.GetLabel(sessionID, entry.ID)
		labelSuffix := ""
		if strings.TrimSpace(label) != "" {
			labelSuffix = fmt.Sprintf(" label=%q", label)
		}
		b.WriteString(fmt.Sprintf("- %s [%s]%s %s\n", entry.ID, entry.Kind, labelSuffix, preview))
	}
	return strings.TrimSpace(b.String()), nil
}

func setFoundationLabel(cwd, sessionID, targetID, label string) (string, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	_, err := runtime.AppendLabelChange(sessionID, targetID, label)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(label) == "" {
		return fmt.Sprintf("[Foundation Label]\nCleared label on %s", targetID), nil
	}
	return fmt.Sprintf("[Foundation Label]\nSet label %q on %s", label, targetID), nil
}

func buildFoundationTreeSelectionDisplay(cwd, sessionID string) (string, []string, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	session, err := runtime.LoadSession(sessionID)
	if err != nil {
		return "", nil, err
	}
	leafID, err := runtime.GetLeafID(sessionID)
	if err != nil {
		return "", nil, err
	}
	ids := make([]string, 0, len(session.Entries))
	var b strings.Builder
	b.WriteString("[Foundation Tree Select]\n")
	b.WriteString("Use /tree-go <index> [maxTokens] to switch to a numbered entry.\n")
	for i, entry := range session.Entries {
		ids = append(ids, entry.ID)
		marker := " "
		if entry.ID == leafID {
			marker = "*"
		}
		preview := entry.Text
		if strings.TrimSpace(preview) == "" {
			preview = entry.ToolName
		}
		if len(preview) > 60 {
			preview = preview[:60] + "..."
		}
		label, _ := runtime.GetLabel(sessionID, entry.ID)
		labelSuffix := ""
		if strings.TrimSpace(label) != "" {
			labelSuffix = fmt.Sprintf(" label=%q", label)
		}
		b.WriteString(fmt.Sprintf("%s [%d] %s [%s]%s %s\n", marker, i+1, entry.ID, entry.Kind, labelSuffix, preview))
	}
	return strings.TrimSpace(b.String()), ids, nil
}

func switchFoundationTreeSelection(cwd, sessionID string, ids []string, index int, maxTokens int) (string, error) {
	if index <= 0 || index > len(ids) {
		return "", fmt.Errorf("selection index %d out of range", index)
	}
	return switchFoundationTreeDisplay(cwd, sessionID, ids[index-1], maxTokens)
}

func buildFoundationTreeBrowser(cwd, sessionID string) ([]TreeBrowserItem, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	session, err := runtime.LoadSession(sessionID)
	if err != nil {
		return nil, err
	}
	leafID, err := runtime.GetLeafID(sessionID)
	if err != nil {
		return nil, err
	}
	items := make([]TreeBrowserItem, 0, len(session.Entries))
	for _, entry := range session.Entries {
		preview := entry.Text
		if strings.TrimSpace(preview) == "" {
			preview = entry.ToolName
		}
		if len(preview) > 60 {
			preview = preview[:60] + "..."
		}
		label, _ := runtime.GetLabel(sessionID, entry.ID)
		item := TreeBrowserItem{ID: entry.ID, Kind: entry.Kind, Label: label, Preview: preview, IsLeaf: entry.ID == leafID}
		if entry.ID != leafID {
			prep, prepErr := runtime.PrepareBranchSummaryWithBudget(sessionID, entry.ID, 128)
			if prepErr == nil && prep != nil {
				item.SummaryEntries = len(prep.EntriesToSummarize)
				item.CommonAncestorID = prep.CommonAncestorID
				item.ReadFilesCount = len(prep.FileOps.ReadFiles)
				item.ModifiedFilesCount = len(prep.FileOps.ModifiedFiles)
			}
		}
		items = append(items, item)
	}
	return items, nil
}

func filterTreeBrowserItems(items []TreeBrowserItem, filter string) []TreeBrowserItem {
	filter = strings.TrimSpace(strings.ToLower(filter))
	if filter == "" {
		return append([]TreeBrowserItem(nil), items...)
	}
	out := make([]TreeBrowserItem, 0, len(items))
	for _, item := range items {
		haystack := strings.ToLower(strings.Join([]string{item.Kind, item.Label, item.Preview}, " "))
		if strings.Contains(haystack, filter) {
			out = append(out, item)
		}
	}
	return out
}

func renderTreeBrowser(items []TreeBrowserItem, selected int, filter string) string {
	visible := filterTreeBrowserItems(items, filter)
	if selected >= len(visible) {
		selected = max(0, len(visible)-1)
	}
	var b strings.Builder
	b.WriteString("[Foundation Tree Browser]\n")
	b.WriteString("Use ↑/↓ to move, type to filter, Backspace to clear, Enter to switch, Esc to close.\n")
	b.WriteString(fmt.Sprintf("filter=%q matches=%d\n", filter, len(visible)))
	for i, item := range visible {
		cursor := " "
		if i == selected {
			cursor = ">"
		}
		leaf := " "
		if item.IsLeaf {
			leaf = "*"
		}
		labelSuffix := ""
		if strings.TrimSpace(item.Label) != "" {
			labelSuffix = fmt.Sprintf(" label=%q", item.Label)
		}
		b.WriteString(fmt.Sprintf("%s%s [%d] %s [%s]%s %s\n", cursor, leaf, i+1, item.ID, item.Kind, labelSuffix, item.Preview))
	}
	if len(visible) > 0 && selected >= 0 && selected < len(visible) {
		item := visible[selected]
		b.WriteString("\n[Preview]\n")
		b.WriteString(fmt.Sprintf("id=%s\nkind=%s\n", item.ID, item.Kind))
		if strings.TrimSpace(item.Label) != "" {
			b.WriteString(fmt.Sprintf("label=%q\n", item.Label))
		}
		b.WriteString(fmt.Sprintf("preview=%s\n", item.Preview))
		if item.IsLeaf {
			b.WriteString("branchSummary=already on active leaf")
		} else {
			b.WriteString(fmt.Sprintf("branchSummaryEntries=%d\n", item.SummaryEntries))
			if strings.TrimSpace(item.CommonAncestorID) != "" {
				b.WriteString(fmt.Sprintf("commonAncestor=%s\n", item.CommonAncestorID))
			}
			b.WriteString(fmt.Sprintf("readFiles=%d modifiedFiles=%d", item.ReadFilesCount, item.ModifiedFilesCount))
		}
	}
	return strings.TrimSpace(b.String())
}

func openSelectedTreeBrowser(cwd, sessionID string, items []TreeBrowserItem, selected, maxTokens int) (string, error) {
	if selected < 0 || selected >= len(items) {
		return "", fmt.Errorf("selected tree item out of range")
	}
	return switchFoundationTreeDisplay(cwd, sessionID, items[selected].ID, maxTokens)
}

func buildShellProposal(director *agents.Director, query string) (ShellProposalMsg, error) {
	execution := adapters.PrepareProviderExecution(adapters.ProviderExecutionRequest{Prompt: query, TaskType: "analysis", CostPreference: "budget"})
	assistant := agents.NewShellTranslator(director.Provider)
	response, err := assistant.Translate(context.Background(), query)
	if err != nil {
		return ShellProposalMsg{}, err
	}
	return ShellProposalMsg{
		Command:     response,
		Explanation: execution.ExecutionHint,
	}, nil
}
