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
	ParentID           string
	Kind               string
	Label              string
	Preview            string
	IsLeaf             bool
	Depth              int
	ChildCount         int
	Prefix             string
	IsLastChild        bool
	GroupKey           string
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
	depthByID := map[string]int{}
	entryByID := map[string]foundationpi.SessionEntry{}
	childrenByParent := map[string][]string{}
	for _, entry := range session.Entries {
		entryByID[entry.ID] = entry
		childrenByParent[entry.ParentID] = append(childrenByParent[entry.ParentID], entry.ID)
	}
	for _, entry := range session.Entries {
		depth := 0
		if entry.ParentID != "" {
			depth = depthByID[entry.ParentID] + 1
		}
		depthByID[entry.ID] = depth
		preview := entry.Text
		if strings.TrimSpace(preview) == "" {
			preview = entry.ToolName
		}
		if len(preview) > 60 {
			preview = preview[:60] + "..."
		}
		label, _ := runtime.GetLabel(sessionID, entry.ID)
		children, _ := runtime.GetChildren(sessionID, entry.ID)
		item := TreeBrowserItem{ID: entry.ID, ParentID: entry.ParentID, Kind: entry.Kind, Label: label, Preview: preview, IsLeaf: entry.ID == leafID, Depth: depth, ChildCount: len(children)}
		item.GroupKey = buildGroupKey(entry, entryByID)
		siblings := childrenByParent[entry.ParentID]
		if len(siblings) > 0 && siblings[len(siblings)-1] == entry.ID {
			item.IsLastChild = true
		}
		item.Prefix = buildTreePrefix(entry, entryByID, childrenByParent)
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

func buildGroupKey(entry foundationpi.SessionEntry, entryByID map[string]foundationpi.SessionEntry) string {
	if entry.ParentID == "" {
		return "root"
	}
	current := entry
	for current.ParentID != "" {
		parent, ok := entryByID[current.ParentID]
		if !ok {
			break
		}
		if parent.ParentID == "" {
			return "branch:" + parent.ID
		}
		current = parent
	}
	return "root"
}

func buildTreePrefix(entry foundationpi.SessionEntry, entryByID map[string]foundationpi.SessionEntry, childrenByParent map[string][]string) string {
	if entry.ParentID == "" {
		return ""
	}
	segments := make([]string, 0)
	currentParentID := entry.ParentID
	for currentParentID != "" {
		parent, ok := entryByID[currentParentID]
		if !ok {
			break
		}
		siblings := childrenByParent[parent.ParentID]
		isLast := len(siblings) > 0 && siblings[len(siblings)-1] == parent.ID
		if parent.ParentID == "" {
			if isLast {
				segments = append([]string{"  "}, segments...)
			} else {
				segments = append([]string{"│ "}, segments...)
			}
		} else {
			if isLast {
				segments = append([]string{"  "}, segments...)
			} else {
				segments = append([]string{"│ "}, segments...)
			}
		}
		currentParentID = parent.ParentID
	}
	return strings.Join(segments, "")
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

func visibleTreeBrowserItems(items []TreeBrowserItem, filter string, collapsed map[string]bool) []TreeBrowserItem {
	filtered := filterTreeBrowserItems(items, filter)
	if len(collapsed) == 0 {
		return filtered
	}
	out := make([]TreeBrowserItem, 0, len(filtered))
	skipDepth := -1
	for _, item := range filtered {
		if skipDepth >= 0 {
			if item.Depth > skipDepth {
				continue
			}
			skipDepth = -1
		}
		out = append(out, item)
		if collapsed[item.ID] {
			skipDepth = item.Depth
		}
	}
	return out
}

func renderTreeBrowser(items []TreeBrowserItem, selected int, filter string, confirmPending bool, collapsed map[string]bool, grouped bool, maxVisible int, title string) string {
	visible := visibleTreeBrowserItems(items, filter, collapsed)
	if selected >= len(visible) {
		selected = max(0, len(visible)-1)
	}
	start := 0
	end := len(visible)
	if maxVisible > 0 && len(visible) > maxVisible {
		start = selected - maxVisible/2
		if start < 0 {
			start = 0
		}
		end = start + maxVisible
		if end > len(visible) {
			end = len(visible)
			start = max(0, end-maxVisible)
		}
	}
	window := visible[start:end]
	var b strings.Builder
	b.WriteString(title + "\n")
	if confirmPending {
		b.WriteString("Confirm switch: Enter/Y = confirm, N/Esc/Backspace = cancel.\n")
	} else {
		b.WriteString("Use ↑/↓ to move, PgUp/PgDn/Home/End for viewport jumps, type to filter, Backspace to clear, Enter to arm switch, Esc to close, Tab to toggle grouping.\n")
	}
	b.WriteString(fmt.Sprintf("filter=%q matches=%d grouped=%t\n", filter, len(visible), grouped))
	if maxVisible > 0 && len(visible) > maxVisible {
		b.WriteString(fmt.Sprintf("showing=%d-%d of %d\n", start+1, end, len(visible)))
	}
	lastGroup := ""
	for offset, item := range window {
		i := start + offset
		if grouped && item.GroupKey != lastGroup {
			b.WriteString(fmt.Sprintf("\n[Group] %s\n", item.GroupKey))
			lastGroup = item.GroupKey
		}
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
		childSuffix := ""
		if item.ChildCount > 0 {
			childSuffix = fmt.Sprintf(" children=%d", item.ChildCount)
		}
		glyph := "•"
		if item.Depth > 0 {
			if item.IsLastChild {
				glyph = "└─"
			} else {
				glyph = "├─"
			}
		}
		fold := "  "
		if item.ChildCount > 0 {
			if collapsed != nil && collapsed[item.ID] {
				fold = "[+]"
			} else {
				fold = "[-]"
			}
		}
		b.WriteString(fmt.Sprintf("%s%s [%d] %s%s %s %s [%s]%s%s %s\n", cursor, leaf, i+1, item.Prefix, glyph, fold, item.ID, item.Kind, labelSuffix, childSuffix, item.Preview))
	}
	if len(visible) > 0 && selected >= 0 && selected < len(visible) {
		item := visible[selected]
		b.WriteString("\n[Preview]\n")
		b.WriteString(fmt.Sprintf("id=%s\nkind=%s\ndepth=%d\nchildren=%d\n", item.ID, item.Kind, item.Depth, item.ChildCount))
		if strings.TrimSpace(item.Label) != "" {
			b.WriteString(fmt.Sprintf("label=%q\n", item.Label))
		}
		b.WriteString(fmt.Sprintf("preview=%s\n", item.Preview))
		if item.IsLeaf {
			b.WriteString("branchSummary=already on active leaf\n")
		} else {
			b.WriteString(fmt.Sprintf("branchSummaryEntries=%d\n", item.SummaryEntries))
			if strings.TrimSpace(item.CommonAncestorID) != "" {
				b.WriteString(fmt.Sprintf("commonAncestor=%s\n", item.CommonAncestorID))
			}
			b.WriteString(fmt.Sprintf("readFiles=%d modifiedFiles=%d\n", item.ReadFilesCount, item.ModifiedFilesCount))
		}
		if confirmPending {
			b.WriteString("\n[Confirm]\nSwitch to this entry? Press Enter or Y to confirm, N/Esc/Backspace to cancel.")
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

func pinFoundationTreeBrowser(m *model) error {
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		return err
	}
	items, err := buildFoundationTreeBrowser(m.director.WorkingDir, sessionID)
	if err != nil {
		return err
	}
	m.browserItems = items
	m.browserPinned = true
	return nil
}

func refreshPinnedFoundationTreeBrowser(m *model) {
	if !m.browserPinned {
		return
	}
	sessionID, err := ensureFoundationSession(m)
	if err != nil {
		return
	}
	items, err := buildFoundationTreeBrowser(m.director.WorkingDir, sessionID)
	if err != nil {
		return
	}
	m.browserItems = items
	visible := visibleTreeBrowserItems(m.browserItems, m.browserFilter, m.browserCollapsed)
	if m.browserIndex >= len(visible) {
		m.browserIndex = max(0, len(visible)-1)
	}
}

func unpinFoundationTreeBrowser(m *model) {
	m.browserPinned = false
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
