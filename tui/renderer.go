package tui

// ═══════════════════════════════════════════════════════════════════════
// Renderer — renders ChatEntry objects into styled terminal output
// Mirrors pi-mono's component rendering system
// ═══════════════════════════════════════════════════════════════════════

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderEntry renders a single ChatEntry with pi-mono-style theming
func RenderEntry(e ChatEntry, t Theme) string {
	switch e.Type {
	case EntryUser:
		return renderUserEntry(e, t)
	case EntryAssistant:
		return renderAssistantEntry(e, t)
	case EntryTool:
		return renderToolEntry(e, t)
	case EntrySystem:
		return e.Content
	case EntryThinking:
		return renderThinkingEntry(e, t)
	case EntryShellProposal:
		return renderShellProposal(e, t)
	case EntryDiff:
		return renderDiffEntry(e, t)
	case EntryCompactionSummary:
		return renderCompactionSummary(e, t)
	case EntryBashMode:
		return renderBashModeEntry(e, t)
	default:
		return e.Content
	}
}

// ─── User message ─────────────────────────────────────────────────────

func renderUserEntry(e ChatEntry, t Theme) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.UserMsgText)).
		Background(lipgloss.Color(t.UserMsgBg)).
		Padding(0, 1)
	return style.Render(t.Bold(t.AccentText("> ")) + t.Fg(t.UserMsgText, e.Content))
}

// ─── Assistant message ────────────────────────────────────────────────

func renderAssistantEntry(e ChatEntry, t Theme) string {
	var header string
	if e.Provider != "" || e.Model != "" {
		modelStr := e.Provider
		if e.Model != "" {
			modelStr += "/" + e.Model
		}
		header = t.Dim("  "+modelStr) + "\n"
	}
	content := t.Fg(t.TextColor, e.Content)
	content = RenderSimpleMarkdown(content, t)
	return header + content
}

// ─── Tool execution ───────────────────────────────────────────────────

func renderToolEntry(e ChatEntry, t Theme) string {
	var statusIcon string
	if e.ToolErr {
		statusIcon = t.ErrorText("✗")
	} else if e.Streaming {
		statusIcon = t.WarningText("⟳")
	} else {
		statusIcon = t.SuccessText("✓")
	}

	borderColor := t.Warning
	if e.ToolErr {
		borderColor = t.Error
	} else if e.Streaming {
		borderColor = t.Accent
	}

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1)

	var b strings.Builder
	b.WriteString(fmt.Sprintf(" %s %s ", statusIcon, t.Bold(t.Fg(t.ToolTitle, e.ToolName))))

	if e.ToolDur > 0 {
		b.WriteString(t.Dim(fmt.Sprintf(" (%v)", e.ToolDur)))
	}

	if e.ToolArgs != "" {
		argDisplay := e.ToolArgs
		if len(argDisplay) > 80 {
			argDisplay = argDisplay[:80] + "…"
		}
		b.WriteString("\n  " + t.Dim(argDisplay))
	}

	body := e.Content
	if body != "" {
		if e.Expanded || e.ToolErr {
			if len(body) > 800 && !e.Expanded {
				body = body[:800] + "\n" + t.Dim("  … (truncated, Ctrl+O to expand)")
			}
			b.WriteString("\n" + t.Fg(t.ToolOutput, body))
		} else {
			lines := strings.SplitN(body, "\n", 2)
			preview := lines[0]
			if len(preview) > 120 {
				preview = preview[:120] + "…"
			}
			b.WriteString("\n" + t.Dim(preview))
		}
	}

	return borderStyle.Render(b.String())
}

// ─── Thinking block ───────────────────────────────────────────────────

func renderThinkingEntry(e ChatEntry, t Theme) string {
	levelColor := ThinkingLevelColors[e.ThinkingLevel]
	if levelColor == "" {
		levelColor = t.ThinkingCol
	}

	if e.Hidden {
		return t.Italic(t.Fg(t.ThinkingCol, e.Content))
	}

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(levelColor)).
		Padding(0, 1)

	return borderStyle.Render(t.Italic(t.Fg(t.ThinkingCol, e.Content)))
}

// ─── Shell proposal ───────────────────────────────────────────────────

func renderShellProposal(e ChatEntry, t Theme) string {
	var b strings.Builder
	b.WriteString(t.Fg(t.ToolTitle, "⌘ Shell Proposal"))
	b.WriteString("\n")
	b.WriteString(t.Fg(t.TextColor, "  $ "+e.Content))
	if e.ToolArgs != "" {
		b.WriteString("\n" + t.Dim("  " + e.ToolArgs))
	}
	b.WriteString("\n" + t.WarningText("  Execute? [Y/n] Ctrl+Y to accept"))
	return b.String()
}

// ─── Diff rendering ───────────────────────────────────────────────────

func renderDiffEntry(e ChatEntry, t Theme) string {
	lines := strings.Split(e.Content, "\n")
	var rendered []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "+") || strings.HasPrefix(trimmed, "→") {
			rendered = append(rendered, t.Fg(t.DiffAdded, line))
		} else if strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "✗") {
			rendered = append(rendered, t.Fg(t.DiffRemoved, line))
		} else if strings.HasPrefix(trimmed, "@@") {
			rendered = append(rendered, t.Fg(t.DiffContext, line))
		} else {
			rendered = append(rendered, t.Dim(line))
		}
	}

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.Accent)).
		Padding(0, 1)

	return borderStyle.Render(strings.Join(rendered, "\n"))
}

// ─── Compaction summary ───────────────────────────────────────────────

func renderCompactionSummary(e ChatEntry, t Theme) string {
	var b strings.Builder
	b.WriteString(t.Bold(t.AccentText("╭─ Compaction Summary ────────────────────────╮")))
	b.WriteString("\n")
	b.WriteString(t.Dim("  Tokens before: ") + t.Fg(t.TextColor, fmt.Sprintf("%d", e.TokensBefore)))
	b.WriteString("\n")
	b.WriteString(t.Dim("  Tokens after:  ") + t.Fg(t.TextColor, fmt.Sprintf("%d", e.TokensAfter)))
	b.WriteString("\n")
	if e.Content != "" {
		b.WriteString(t.Dim("  Summary:       ") + t.Fg(t.Muted, e.Content))
		b.WriteString("\n")
	}
	b.WriteString(t.Bold(t.AccentText("╰───────────────────────────────────────────────╯")))
	return b.String()
}

// ─── Bash mode (! prefix) ─────────────────────────────────────────────

func renderBashModeEntry(e ChatEntry, t Theme) string {
	var b strings.Builder
	b.WriteString(t.Fg(t.BashModeCol, "⚡ Bash Mode"))
	b.WriteString("\n")
	b.WriteString(t.Fg(t.TextColor, "  $ "+e.Content))
	return b.String()
}

// ─── Simple markdown renderer ─────────────────────────────────────────

func RenderSimpleMarkdown(text string, t Theme) string {
	lines := strings.Split(text, "\n")
	var result []string
	inCodeBlock := false

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			if inCodeBlock {
				lang := strings.TrimPrefix(line, "```")
				if lang != "" {
					result = append(result, t.Dim("```"+lang))
				} else {
					result = append(result, t.Dim("```"))
				}
			} else {
				result = append(result, t.Dim("```"))
			}
			continue
		}
		if inCodeBlock {
			result = append(result, t.Fg(t.MDCodeBlock, line))
			continue
		}
		// Headings
		if strings.HasPrefix(line, "### ") {
			result = append(result, t.Bold(t.Fg(t.MDHeading, line)))
			continue
		}
		if strings.HasPrefix(line, "## ") {
			result = append(result, t.Bold(t.Fg(t.MDHeading, line)))
			continue
		}
		if strings.HasPrefix(line, "# ") {
			result = append(result, t.Bold(t.Fg(t.MDHeading, line)))
			continue
		}
		// List items
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			bullet := t.Fg(t.MDListBullet, line[:2])
			result = append(result, bullet+t.Fg(t.TextColor, line[2:]))
			continue
		}
		// Numbered lists
		if len(line) > 2 && line[0] >= '1' && line[0] <= '9' && line[1] == '.' && line[2] == ' ' {
			bullet := t.Fg(t.MDListBullet, line[:3])
			result = append(result, bullet+t.Fg(t.TextColor, line[3:]))
			continue
		}
		// Quote
		if strings.HasPrefix(line, "> ") {
			result = append(result, t.Fg(t.MDQuote, line))
			continue
		}
		// Horizontal rule
		if strings.TrimSpace(line) == "---" || strings.TrimSpace(line) == "***" {
			result = append(result, t.Dim("────────────────────────────────────────"))
			continue
		}

		// Inline formatting: `code`, **bold**, *italic*, [link](url)
		line = renderInlineFormatting(line, t)
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

func renderInlineFormatting(line string, t Theme) string {
	// Inline code: `text`
	for {
		start := strings.Index(line, "`")
		if start == -1 {
			break
		}
		end := strings.Index(line[start+1:], "`")
		if end == -1 {
			break
		}
		codeContent := line[start+1 : start+1+end]
		replacement := t.Fg(t.MDCode, codeContent)
		line = line[:start] + replacement + line[start+1+end+1:]
	}

	// Bold: **text**
	for {
		start := strings.Index(line, "**")
		if start == -1 {
			break
		}
		end := strings.Index(line[start+2:], "**")
		if end == -1 {
			break
		}
		boldContent := line[start+2 : start+2+end]
		replacement := t.Bold(t.Fg(t.TextColor, boldContent))
		line = line[:start] + replacement + line[start+2+end+2:]
	}

	// Links: [text](url)
	for {
		start := strings.Index(line, "[")
		if start == -1 {
			break
		}
		end := strings.Index(line[start:], "](")
		if end == -1 {
			break
		}
		urlStart := start + end + 2
		urlEnd := strings.Index(line[urlStart:], ")")
		if urlEnd == -1 {
			break
		}
		linkText := line[start+1 : start+end]
		url := line[urlStart : urlStart+urlEnd]
		replacement := t.Fg(t.MDLink, linkText) + t.Dim("("+url+")")
		line = line[:start] + replacement + line[urlStart+urlEnd+1:]
	}

	return line
}

// ─── Welcome banner ───────────────────────────────────────────────────

func RenderWelcome(workingDir, gitBranch, provider, modelName string, toolCount, regCount int, t Theme) string {
	var b strings.Builder

	// Logo
	b.WriteString(t.Bold(t.AccentText("╭─────────────────────────────────────────────────────╮")))
	b.WriteString("\n")
	b.WriteString(t.Bold(t.AccentText("│         🧠  HyperHarness — AI Coding Agent          │")))
	b.WriteString("\n")
	b.WriteString(t.Bold(t.AccentText("╰─────────────────────────────────────────────────────╯")))
	b.WriteString("\n\n")

	// Status line
	b.WriteString(t.Dim("  cwd   ") + t.Fg(t.TextColor, shortenPath(workingDir)))
	if gitBranch != "" {
		b.WriteString(t.Dim(" (") + t.Fg(t.Muted, gitBranch) + t.Dim(")"))
	}
	b.WriteString("\n")
	b.WriteString(t.Dim("  tools ") + t.Fg(t.TextColor, fmt.Sprintf("%d registered + %d CLI detected", regCount, toolCount-regCount)))
	b.WriteString("\n")
	b.WriteString(t.Dim("  model ") + t.Fg(t.TextColor, provider+"/"+modelName))
	b.WriteString("\n\n")

	// Key bindings (pi-mono style)
	b.WriteString(t.Dim("  ─── Key Bindings ────────────────────────────────────────"))
	b.WriteString("\n")
	keys := []struct{ k, d string }{
		{"Enter", "Send message"},
		{"Esc×2", "Quit"},
		{"↑/↓", "Input history / scroll chat"},
		{"Tab", "Autocomplete commands"},
		{"Ctrl+C", "Interrupt agent / quit"},
		{"Ctrl+O", "Expand/collapse tool output"},
		{"Ctrl+P", "Cycle models"},
		{"Ctrl+L", "Toggle file tree pane"},
		{"Ctrl+D", "Toggle dashboard"},
		{"Ctrl+Y", "Accept shell proposal"},
		{"Ctrl+E", "Open external editor"},
		{"!", "Bash mode (direct shell)"},
		{"!!", "Bash mode (no context)"},
		{"/", "Slash commands"},
		{"??", "Shell command proposal"},
	}
	for _, k := range keys {
		b.WriteString("  " + t.KeyHint(k.k, k.d))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(t.Dim("  ─── Slash Commands ──────────────────────────────────────"))
	b.WriteString("\n")
	quick := []struct{ cmd, desc string }{
		{"/help", "All commands"},
		{"/tree-browser", "File explorer (modal)"},
		{"/tree-pane", "Persistent file pane"},
		{"/repomap", "Repository map"},
		{"/providers", "LLM provider status"},
		{"/mcp", "MCP tool listing"},
		{"/tools", "All registered tools"},
		{"/plan <prompt>", "Orchestration plan"},
		{"/model", "Select model"},
		{"/settings", "Settings menu"},
		{"/compact", "Compact context"},
		{"/hotkeys", "All keybindings"},
		{"/changelog", "Show changelog"},
		{"/quit", "Quit HyperHarness"},
	}
	for _, q := range quick {
		b.WriteString("  " + t.AccentText(q.cmd) + t.Dim("  "+q.desc))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(t.Dim("  ─────────────────────────────────────────────────────────"))
	b.WriteString("\n")
	b.WriteString(t.Dim("  Ask HyperHarness to explain its own features."))
	b.WriteString("\n")

	return b.String()
}

// ─── Footer rendering ─────────────────────────────────────────────────

func RenderFooter(workingDir, gitBranch, provider, modelName string, toolCount, totalInputTok, totalOutTok int, totalCost, contextPct float64, contextWindow int, thinkingLevel string, t Theme, width int) string {
	// Line 1: pwd + git branch
	pwd := shortenPath(workingDir)
	if gitBranch != "" {
		pwd += " (" + gitBranch + ")"
	}

	// Line 2: token stats + context + model
	var stats []string
	if totalInputTok > 0 {
		stats = append(stats, fmt.Sprintf("↑%s", formatTokens(totalInputTok)))
	}
	if totalOutTok > 0 {
		stats = append(stats, fmt.Sprintf("↓%s", formatTokens(totalOutTok)))
	}
	if totalCost > 0 {
		stats = append(stats, fmt.Sprintf("$%.3f", totalCost))
	}

	contextStr := fmt.Sprintf("%.0f%%/%s", contextPct, formatTokens(contextWindow))
	if contextPct > 90 {
		contextStr = t.ErrorText(contextStr)
	} else if contextPct > 70 {
		contextStr = t.WarningText(contextStr)
	}
	stats = append(stats, contextStr)

	// Thinking level indicator
	if thinkingLevel != "" && thinkingLevel != "off" {
		levelColor := ThinkingLevelColors[thinkingLevel]
		if levelColor == "" {
			levelColor = t.ThinkingCol
		}
		stats = append(stats, t.Fg(levelColor, "think:"+thinkingLevel))
	}

	statsLeft := strings.Join(stats, " ")

	modelStr := modelName
	if provider != "" && provider != "hypercode" {
		modelStr = "(" + provider + ") " + modelName
	}

	pwdLine := t.Dim(pwd)
	if len(pwdLine) > width {
		pwdLine = pwdLine[:width]
	}

	rightWidth := len(modelStr)
	leftWidth := len(statsLeft)
	gap := max(0, width-leftWidth-rightWidth-2)
	statsLine := t.Dim(statsLeft) + strings.Repeat(" ", gap) + t.Dim(modelStr)

	return pwdLine + "\n" + statsLine
}

// ─── Autocomplete rendering ───────────────────────────────────────────

func RenderAutocomplete(items []SlashCommand, index, maxVis int, t Theme) string {
	if len(items) == 0 {
		return t.Dim("  No matching commands")
	}

	var lines []string
	start := 0
	if index >= maxVis {
		start = index - maxVis + 1
	}
	end := min(len(items), start+maxVis)

	for i := start; i < end; i++ {
		item := items[i]
		if i == index {
			lines = append(lines, t.Bold(t.AccentText("  ▶ /"+item.Name))+t.Dim("  "+item.Description))
		} else {
			lines = append(lines, t.Dim("    /"+item.Name)+"  "+t.Dim(item.Description))
		}
	}

	if len(items) > maxVis {
		lines = append(lines, t.Dim(fmt.Sprintf("  ↑↓ scroll  (%d/%d)", index+1, len(items))))
	}

	return strings.Join(lines, "\n")
}
