package tui

// ═══════════════════════════════════════════════════════════════════════
// renderer.go — Full entry rendering
// Unifies pi-mono's component rendering, goose's tool call panels,
// opencode's markdown rendering, and claude-code's diff/thinkback display
// ═══════════════════════════════════════════════════════════════════════

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// RenderEntry dispatches to the correct renderer based on entry type
func RenderEntry(entry ChatEntry, t Theme) string {
	switch entry.Type {
	case EntryUser:
		return renderUserEntry(entry, t)
	case EntryAssistant:
		return renderAssistantEntry(entry, t)
	case EntryTool:
		return renderToolEntry(entry, t)
	case EntrySystem:
		return renderSystemEntry(entry, t)
	case EntryThinking:
		return renderThinkingEntry(entry, t)
	case EntryShellProposal:
		return renderShellProposalEntry(entry, t)
	case EntryDiff:
		return renderDiffEntry(entry, t)
	case EntryCompactionSummary:
		return renderCompactionSummaryEntry(entry, t)
	case EntryBashMode:
		return renderBashModeEntry(entry, t)
	case EntryPermission:
		return renderPermissionEntry(entry, t)
	case EntryImage:
		return renderImageEntry(entry, t)
	case EntryCustom:
		return renderCustomEntry(entry, t)
	case EntryError:
		return renderErrorEntry(entry, t)
	case EntryQueue:
		return renderQueueEntry(entry, t)
	default:
		return t.Dim(fmt.Sprintf("[unknown entry type %d]", entry.Type))
	}
}

// ═══════════════════════════════════════════════════════════════════════
// User entry — pi-mono's UserMessageComponent
// ═══════════════════════════════════════════════════════════════════════

func renderUserEntry(e ChatEntry, t Theme) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.UserMessageText)).
		Background(lipgloss.Color(t.UserMessageBg)).
		Padding(0, 1)
	return t.Bold(t.Fg(t.Accent, "> ")) + style.Render(e.Content)
}

// ═══════════════════════════════════════════════════════════════════════
// Assistant entry — pi-mono's AssistantMessageComponent with markdown
// ═══════════════════════════════════════════════════════════════════════

func renderAssistantEntry(e ChatEntry, t Theme) string {
	var parts []string
	if e.Provider != "" || e.Model != "" {
		label := ""
		if e.Provider != "" { label += e.Provider }
		if e.Model != "" { label += "/" + e.Model }
		parts = append(parts, t.Dim("┃ "+label))
	}
	md := RenderSimpleMarkdown(e.Content, t)
	parts = append(parts, md)
	return strings.Join(parts, "\n")
}

// ═══════════════════════════════════════════════════════════════════════
// Tool entry — unified pi-mono bordered panel + goose status/kind icons
// ═══════════════════════════════════════════════════════════════════════

func renderToolEntry(e ChatEntry, t Theme) string {
	// Status icon (goose-style: ○ pending, ◑ in_progress, ● completed, ✗ failed)
	statusIcon := t.StatusIcon(e.ToolStatus)
	if e.ToolStatus == "" {
		if e.ToolErr {
			statusIcon = t.ErrorText("✗")
		} else {
			statusIcon = t.SuccessText("✓")
		}
	}

	// Kind icon (goose-style: 📖 read, ✏️ edit, 🔍 search, ▶ execute)
	kind := e.ToolKind
	if kind == "" { kind = DetectToolKind(e.ToolName) }
	kindIcon := ToolKindIcon(kind)

	// Duration
	durStr := ""
	if e.ToolDur > 0 {
		durStr = t.Dim(fmt.Sprintf(" (%v)", e.ToolDur.Round(time.Millisecond)))
	}

	// Title line
	titleLine := statusIcon + " " + kindIcon + " " + t.ToolTitleStyled(e.ToolName) + durStr

	// Args preview
	argsLine := ""
	if e.ToolArgs != "" {
		preview := e.ToolArgs
		if len(preview) > 80 { preview = preview[:77] + "…" }
		argsLine = t.Dim("  ▸ in: " + preview)
	}

	// Locations (goose-style: 📁 path:line)
	locLines := ""
	for _, loc := range e.ToolLocations {
		locStr := loc.Path
		if loc.Line > 0 { locStr += fmt.Sprintf(":%d", loc.Line) }
		locLines += "\n" + t.Dim("  📁 "+locStr)
	}

	// Output
	outputLine := ""
	if e.Expanded && e.Content != "" {
		output := e.Content
		if len(output) > 500 { output = output[:497] + "…" }
		outputLine = "\n" + t.Fg(t.ToolOutput, "  ▸ out: "+output)
	} else if e.Content != "" && !e.Expanded {
		preview := e.Content
		if len(preview) > 60 { preview = preview[:57] + "…" }
		outputLine = "\n" + t.Dim("  ▸ "+preview+" "+t.KeyHint("Ctrl+O", "expand"))
	}

	// Border (goose-style: ╭─╮ │ │ ╰─╯, colored by status)
	borderColor := t.Border
	if e.ToolErr || e.ToolStatus == "failed" { borderColor = t.Cranberry }

	// Combine in bordered panel
	content := titleLine
	if argsLine != "" { content += "\n" + argsLine }
	content += locLines
	if outputLine != "" { content += outputLine }

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1).
		Render(content)

	return panel
}

// ═══════════════════════════════════════════════════════════════════════
// System entry
// ═══════════════════════════════════════════════════════════════════════

func renderSystemEntry(e ChatEntry, t Theme) string {
	return t.Dim("  " + e.Content)
}

// ═══════════════════════════════════════════════════════════════════════
// Thinking entry — pi-mono's collapsible thinking with level-colored border
// ═══════════════════════════════════════════════════════════════════════

func renderThinkingEntry(e ChatEntry, t Theme) string {
	level := e.ThinkingLevel
	if level == "" { level = "off" }
	borderColor := t.ThinkingBorder(level)

	if e.Hidden || level == "off" {
		label := t.Fg(borderColor, "💭 Thinking")+t.Dim(" ["+level+"]")
		return label
	}

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1).
		Render(t.Fg(t.ThinkingCol, e.Content))

	return t.Fg(borderColor, "💭 Thinking ["+level+"]") + "\n" + panel
}

// ═══════════════════════════════════════════════════════════════════════
// Shell proposal — pi-mono's ShellProposalComponent
// ═══════════════════════════════════════════════════════════════════════

func renderShellProposalEntry(e ChatEntry, t Theme) string {
	cmdLine := t.Fg(t.BashModeCol, "$ ") + t.Bold(t.Fg(t.TextColor, e.Content))
	if e.ToolArgs != "" {
		cmdLine += "\n" + t.Dim("  "+e.ToolArgs)
	}
	hint := t.KeyHint("Ctrl+Y", "accept")
	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.BashModeCol)).
		Padding(0, 1).
		Render("⌘ Shell Proposal\n" + cmdLine + "\n" + hint)
	return panel
}

// ═══════════════════════════════════════════════════════════════════════
// Diff entry — claude-code style with colored +/- lines
// ═══════════════════════════════════════════════════════════════════════

func renderDiffEntry(e ChatEntry, t Theme) string {
	lines := strings.Split(e.Content, "\n")
	var rendered []string
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			rendered = append(rendered, t.Fg(t.ToolDiffAdded, line))
		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			rendered = append(rendered, t.Fg(t.ToolDiffRemoved, line))
		case strings.HasPrefix(line, "@@"):
			rendered = append(rendered, t.Fg(t.ToolDiffContext, line))
		default:
			rendered = append(rendered, t.Fg(t.TextColor, line))
		}
	}
	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.MdCodeBlockBorder)).
		Padding(0, 1).
		Render(strings.Join(rendered, "\n"))
	return panel
}

// ═══════════════════════════════════════════════════════════════════════
// Compaction summary — pi-mono's CompactionSummaryComponent
// ═══════════════════════════════════════════════════════════════════════

func renderCompactionSummaryEntry(e ChatEntry, t Theme) string {
	line := t.Dim("📋 Compacted")
	if e.TokensBefore > 0 || e.TokensAfter > 0 {
		line += t.Dim(fmt.Sprintf(" %d → %d tokens", e.TokensBefore, e.TokensAfter))
	}
	return line
}

// ═══════════════════════════════════════════════════════════════════════
// Bash mode — goose-style with ⚡ indicator
// ═══════════════════════════════════════════════════════════════════════

func renderBashModeEntry(e ChatEntry, t Theme) string {
	return t.Fg(t.BashModeCol, "⚡ ") + t.Fg(t.TextColor, e.Content)
}

// ═══════════════════════════════════════════════════════════════════════
// Permission entry — goose-style bordered permission dialog
// ═══════════════════════════════════════════════════════════════════════

func renderPermissionEntry(e ChatEntry, t Theme) string {
	pc := t.PermissionColor()
	hRule := strings.Repeat("─", 56)

	var lines []string
	lines = append(lines, t.Fg(pc, "╭"+hRule+"╮"))
	lines = append(lines, t.Fg(pc, "│ ") + t.Bold(t.Fg(pc, "🔒 Permission required")) + t.Fg(pc, " │"))
	lines = append(lines, t.Fg(pc, "│ ") + " " + t.Fg(pc, " │"))
	lines = append(lines, t.Fg(pc, "│ ") + t.Fg(t.TextColor, e.Content) + t.Fg(pc, " │"))
	lines = append(lines, t.Fg(pc, "│ ") + " " + t.Fg(pc, " │"))

	for i, opt := range e.PermissionOptions {
		icon := " "
		if i == e.PermissionIdx { icon = "▸" }
		color := t.DimColor
		if i == e.PermissionIdx { color = t.TextColor }
		keyHint := ""
		switch opt.Kind {
		case "allow_once": keyHint = "y"
		case "allow_always": keyHint = "a"
		case "reject_once": keyHint = "n"
		case "reject_always": keyHint = "N"
		}
		lines = append(lines, t.Fg(pc, "│ ") + t.Fg(color, icon+" ["+keyHint+"] "+opt.Name) + t.Fg(pc, " │"))
	}

	lines = append(lines, t.Fg(pc, "│ ") + " " + t.Fg(pc, " │"))
	lines = append(lines, t.Fg(pc, "│ ") + t.Dim("↑↓ select · enter confirm · esc cancel") + t.Fg(pc, " │"))
	lines = append(lines, t.Fg(pc, "╰"+hRule+"╯"))

	return strings.Join(lines, "\n")
}

// ═══════════════════════════════════════════════════════════════════════
// Image entry — pi-mono's clipboard image paste
// ═══════════════════════════════════════════════════════════════════════

func renderImageEntry(e ChatEntry, t Theme) string {
	return t.Fg(t.Accent, "🖼 ") + t.Fg(t.TextColor, e.ImagePath) + t.Dim(" ("+e.ImageMime+")")
}

// ═══════════════════════════════════════════════════════════════════════
// Custom/extension entry — pi-mono's custom messages
// ═══════════════════════════════════════════════════════════════════════

func renderCustomEntry(e ChatEntry, t Theme) string {
	label := e.CustomLabel
	if label == "" { label = "custom" }
	bg := e.CustomBg
	if bg == "" { bg = t.CustomMessageBg }
	fg := t.CustomMessageText
	labelColor := t.CustomMessageLabel

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(fg)).
		Background(lipgloss.Color(bg)).
		Padding(0, 1)

	return t.Fg(labelColor, "["+label+"] ") + style.Render(e.Content)
}

// ═══════════════════════════════════════════════════════════════════════
// Error entry — goose-style error display
// ═══════════════════════════════════════════════════════════════════════

func renderErrorEntry(e ChatEntry, t Theme) string {
	msg := e.ErrorMessage
	if msg == "" { msg = e.Content }
	return t.ErrorText("⚠ Error: ") + t.ErrorText(msg)
}

// ═══════════════════════════════════════════════════════════════════════
// Queue entry — goose-style queued message indicator
// ═══════════════════════════════════════════════════════════════════════

func renderQueueEntry(e ChatEntry, t Theme) string {
	return t.Fg(t.DimColor, "> ") + t.Dim(e.Content) + t.Fg(t.Gold, " (queued)")
}

// ═══════════════════════════════════════════════════════════════════════
// Markdown rendering — pi-mono's MarkdownComponent
// ═══════════════════════════════════════════════════════════════════════

func RenderSimpleMarkdown(text string, t Theme) string {
	if text == "" { return "" }
	lines := strings.Split(text, "\n")
	var rendered []string
	inCodeBlock := false

	for _, line := range lines {
		// Code block toggle
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			if inCodeBlock {
				lang := strings.TrimPrefix(line, "```")
				if lang != "" {
					rendered = append(rendered, t.Fg(t.MdCodeBlockBorder, "─ "+lang+" ─"))
				}
			}
			continue
		}
		if inCodeBlock {
			rendered = append(rendered, t.Fg(t.MdCode, line))
			continue
		}

		// Headings
		if strings.HasPrefix(line, "### ") {
			rendered = append(rendered, t.Fg(t.MdHeading, t.Bold("   "+line[4:])))
		} else if strings.HasPrefix(line, "## ") {
			rendered = append(rendered, t.Fg(t.MdHeading, t.Bold("  "+line[3:])))
		} else if strings.HasPrefix(line, "# ") {
			rendered = append(rendered, t.Fg(t.MdHeading, t.Bold(line[2:])))
		} else if strings.HasPrefix(line, "> ") {
			rendered = append(rendered, t.Fg(t.MdQuoteBorder, "│ ")+t.Fg(t.MdQuote, line[2:]))
		} else if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			rendered = append(rendered, t.Fg(t.MdListBullet, "  • ")+t.Fg(t.TextColor, line[2:]))
		} else if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "***") {
			rendered = append(rendered, t.Fg(t.MdHr, strings.Repeat("─", 60)))
		} else {
			// Inline formatting
			l := line
			l = renderInlineCode(l, t)
			l = renderInlineBold(l, t)
			l = renderInlineLinks(l, t)
			rendered = append(rendered, t.Fg(t.TextColor, l))
		}
	}
	return strings.Join(rendered, "\n")
}

func renderInlineCode(text string, t Theme) string {
	var result strings.Builder
	parts := strings.Split(text, "`")
	for i, part := range parts {
		if i%2 == 1 {
			result.WriteString(t.Fg(t.MdCode, part))
		} else {
			result.WriteString(part)
		}
	}
	return result.String()
}

func renderInlineBold(text string, t Theme) string {
	var result strings.Builder
	i := 0
	for i < len(text) {
		if i+1 < len(text) && text[i] == '*' && text[i+1] == '*' {
			end := strings.Index(text[i+2:], "**")
			if end >= 0 {
				result.WriteString(t.Bold(text[i+2 : i+2+end]))
				i = i + 2 + end + 2
				continue
			}
		}
		result.WriteByte(text[i])
		i++
	}
	return result.String()
}

func renderInlineLinks(text string, t Theme) string {
	// Simple [text](url) rendering
	for {
		idx := strings.Index(text, "](")
		if idx < 0 { break }
		start := strings.LastIndex(text[:idx], "[")
		if start < 0 { break }
		end := strings.Index(text[idx+2:], ")")
		if end < 0 { break }
		label := text[start+1 : idx]
		url := text[idx+2 : idx+2+end]
		text = text[:start] + t.Fg(t.MdLink, label) + t.Dim(" ("+url+")") + text[idx+2+end+1:]
	}
	return text
}

// ═══════════════════════════════════════════════════════════════════════
// Welcome banner — pi-mono startup display
// ═══════════════════════════════════════════════════════════════════════

func RenderWelcome(cwd, gitBranch, provider, model string, toolCount, regCount int, t Theme) string {
	var b strings.Builder

	// ASCII logo
	b.WriteString(t.AccentText("  ╦ ╦┌─┐┌─┐┌┐ ┬┌─┐┬ ┬┌┬┐  ╔╗╔┌─┐─┐ ┬┬ ┬┌┬┐\n"))
	b.WriteString(t.AccentText("  ║║║├┤ ├─┤├┴┐│├┤ └┬┘│││  ║║║├┤ ┌─┐│ │ │││\n"))
	b.WriteString(t.AccentText("  ╚╩╝└─┘┴ ┴└─┘┴└   ┴ ┴ ┴  ╝╚╝└─┘└─┘└─┘─┴┘\n"))
	b.WriteString("\n")

	// Status line
	b.WriteString(t.Dim("  cwd:   ") + t.Fg(t.TextColor, cwd) + "\n")
	if gitBranch != "" {
		b.WriteString(t.Dim("  git:   ") + t.Fg(t.Accent, "⎇ "+gitBranch) + "\n")
	}
	b.WriteString(t.Dim("  tools: ") + t.Fg(t.Success, fmt.Sprintf("%d registered", regCount)) + t.Dim(" + ") + t.Fg(t.Success, fmt.Sprintf("%d CLI detected", toolCount-regCount)) + "\n")
	b.WriteString(t.Dim("  model: ") + t.Fg(t.TextColor, provider+"/"+model) + "\n")
	b.WriteString("\n")

	// Key bindings (pi-mono style)
	b.WriteString(t.Dim("  Key bindings:\n"))
	b.WriteString("  " + t.KeyHint("Enter", "send") + "  ")
	b.WriteString(t.KeyHint("Ctrl+C", "interrupt") + "  ")
	b.WriteString(t.KeyHint("Esc Esc", "tree/fork") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+L", "tree pane") + "  ")
	b.WriteString(t.KeyHint("Ctrl+D", "dashboard") + "  ")
	b.WriteString(t.KeyHint("Ctrl+O", "expand tools") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+P", "thinking") + "  ")
	b.WriteString(t.KeyHint("Ctrl+E", "editor") + "  ")
	b.WriteString(t.KeyHint("Ctrl+F", "follow-up") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+N", "model→") + "  ")
	b.WriteString(t.KeyHint("Ctrl+M", "model←") + "  ")
	b.WriteString(t.KeyHint("Tab", "autocomplete") + "\n")
	b.WriteString("  " + t.KeyHint("↑↓", "history") + "  ")
	b.WriteString(t.KeyHint("PgUp/Dn", "scroll") + "  ")
	b.WriteString(t.KeyHint("!", "bash mode") + "\n")
	b.WriteString("\n")

	// Quick commands
	b.WriteString(t.Dim("  Quick commands: "))
	commands := []string{"/help", "/model", "/tools", "/compact", "/commit", "/tree"}
	for i, cmd := range commands {
		if i > 0 { b.WriteString(t.Dim(" · ")) }
		b.WriteString(t.Fg(t.Accent, cmd))
	}
	b.WriteString("\n")

	return b.String()
}

// ═══════════════════════════════════════════════════════════════════════
// Footer — pi-mono's 2-line footer with all stats
// ═══════════════════════════════════════════════════════════════════════

func RenderFooter(cwd, gitBranch, provider, model string, toolCount, inputTok, outTok int, cost float64, contextPct float64, contextWindow int, thinkingLevel string, t Theme, width int) string {
	if width <= 0 { width = 80 }

	// Line 1: pwd + git + tokens
	pwd := shortenPath(cwd)
	line1 := t.Dim(" " + pwd)
	if gitBranch != "" { line1 += t.Dim(" ⎇ ") + t.Fg(t.Accent, gitBranch) }
	line1 += t.Dim(" │ ")
	line1 += t.Dim("↑") + t.Fg(t.TextColor, formatTokens(inputTok)) + t.Dim("/↓") + t.Fg(t.TextColor, formatTokens(outTok))
	line1 += t.Dim(" │ $") + t.Fg(t.TextColor, fmt.Sprintf("%.3f", cost))

	// Context percentage (color-coded)
	ctxColor := t.DimColor
	if contextPct > 90 { ctxColor = t.Error } else if contextPct > 70 { ctxColor = t.Warning }
	line1 += t.Dim(" │ ctx ") + t.Fg(ctxColor, fmt.Sprintf("%.0f%%", contextPct))

	// Truncate if too wide
	if len(line1) > width {
		line1 = line1[:width]
	}

	// Line 2: model + thinking + tools
	line2 := t.Dim(" ") + t.Fg(t.Accent, provider) + t.Dim("/") + t.Fg(t.TextColor, model)
	if thinkingLevel != "" && thinkingLevel != "off" {
		line2 += t.Dim(" │ 💭 ") + t.Fg(t.ThinkingBorder(thinkingLevel), thinkingLevel)
	}
	line2 += t.Dim(" │ tools ") + t.Fg(t.TextColor, fmt.Sprintf("%d", toolCount))
	line2 += t.Dim(" │ " + t.KeyHint("Ctrl+C", "int"))

	if len(line2) > width {
		line2 = line2[:width]
	}

	return t.Dim("─"+strings.Repeat("─", max(width-1, 0))) + "\n" + line1 + "\n" + line2
}

// ═══════════════════════════════════════════════════════════════════════
// Autocomplete dropdown — pi-mono's slash command autocomplete
// ═══════════════════════════════════════════════════════════════════════

func RenderAutocomplete(items []SlashCommand, selected, maxVisible int, t Theme) string {
	if len(items) == 0 { return "" }
	if maxVisible <= 0 { maxVisible = 8 }

	start := 0
	if selected >= maxVisible {
		start = selected - maxVisible + 1
	}
	end := min(len(items), start+maxVisible)

	var lines []string
	lines = append(lines, t.Dim("  ╭─ commands ─╮"))
	for i := start; i < end; i++ {
		cmd := items[i]
		prefix := "  │ "
		if i == selected {
			prefix = t.Fg(t.Accent, "  ▸ ")
			lines = append(lines, prefix+t.Bold(t.Fg(t.Accent, "/"+cmd.Name))+t.Dim(" — "+cmd.Description))
		} else {
			lines = append(lines, prefix+t.Fg(t.TextColor, "/"+cmd.Name)+t.Dim(" — "+cmd.Description))
		}
	}
	lines = append(lines, t.Dim("  ╰────────────╯"))

	return strings.Join(lines, "\n")
}

// ═══════════════════════════════════════════════════════════════════════
// Dashboard — split-pane view (pi-mono + opencode)
// ═══════════════════════════════════════════════════════════════════════

func RenderDashboard(chatContent, toolSidebar, metrics string) string {
	width := 80
	leftW := width * 2 / 3
	rightW := width - leftW - 1

	chatLines := strings.Split(chatContent, "\n")
	toolLines := strings.Split(toolSidebar, "\n")
	maxLines := max(len(chatLines), len(toolLines))

	var b strings.Builder
	for i := 0; i < maxLines; i++ {
		left := ""
		if i < len(chatLines) { left = chatLines[i] }
		right := ""
		if i < len(toolLines) { right = toolLines[i] }
		// Pad left to leftW
		if len(left) < leftW {
			left += strings.Repeat(" ", leftW-len(left))
		} else if len(left) > leftW {
			left = left[:leftW]
		}
		if len(right) > rightW { right = right[:rightW] }
		b.WriteString(left + "│" + right + "\n")
	}
	b.WriteString(strings.Repeat("─", leftW) + "│" + strings.Repeat("─", rightW) + "\n")
	b.WriteString(metrics)

	return b.String()
}

// ═══════════════════════════════════════════════════════════════════════
// Hotkeys display — pi-mono's /hotkeys command
// ═══════════════════════════════════════════════════════════════════════

func RenderHotkeys(t Theme) string {
	var b strings.Builder
	b.WriteString(t.Bold(t.AccentText("⌨  Keyboard Shortcuts\n\n")))
	b.WriteString(t.Dim("── Editing ──────────────────────────────────\n"))
	b.WriteString("  " + t.KeyHint("Enter", "send message") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+E", "external editor") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+C", "interrupt / quit") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+D", "dashboard toggle") + "\n")
	b.WriteString("  " + t.KeyHint("Esc Esc", "tree/fork selector") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+Z", "suspend to shell") + "\n")
	b.WriteString("\n")
	b.WriteString(t.Dim("── Agent ─────────────────────────────────────\n"))
	b.WriteString("  " + t.KeyHint("Ctrl+P", "cycle thinking level") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+N", "cycle model forward") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+M", "cycle model backward") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+O", "toggle tool expansion") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+Y", "accept shell proposal") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+F", "queue follow-up message") + "\n")
	b.WriteString("\n")
	b.WriteString(t.Dim("── Navigation ────────────────────────────────\n"))
	b.WriteString("  " + t.KeyHint("↑↓", "scroll / input history") + "\n")
	b.WriteString("  " + t.KeyHint("PgUp/PgDn", "scroll half page") + "\n")
	b.WriteString("  " + t.KeyHint("Home/End", "scroll to top/bottom") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+L", "tree browser") + "\n")
	b.WriteString("  " + t.KeyHint("Tab", "autocomplete / pane focus") + "\n")
	b.WriteString("\n")
	b.WriteString(t.Dim("── Bash Mode ─────────────────────────────────\n"))
	b.WriteString("  " + t.KeyHint("!", "shell command (in context)") + "\n")
	b.WriteString("  " + t.KeyHint("!!", "shell command (excluded from context)") + "\n")
	b.WriteString("  " + t.KeyHint("??", "natural language → shell command") + "\n")
	return b.String()
}

// ═══════════════════════════════════════════════════════════════════════
// Model selector — pi-mono/opencode style
// ═══════════════════════════════════════════════════════════════════════

func RenderModelSelector(current string, models []string, t Theme) string {
	var b strings.Builder
	b.WriteString(t.Bold(t.AccentText("🤖 Model Selector\n\n")))
	for _, m := range models {
		icon := "  "
		if m == current { icon = t.SuccessText("▸ ") }
		b.WriteString(icon + t.Fg(t.TextColor, m))
		if m == current { b.WriteString(t.Dim(" (active)")) }
		b.WriteString("\n")
	}
	b.WriteString("\n" + t.Dim("Ctrl+N / Ctrl+M to cycle, /model <name> to set"))
	return b.String()
}

// ═══════════════════════════════════════════════════════════════════════
// Helpers
// ═══════════════════════════════════════════════════════════════════════

// formatTokens, max, min are defined in chat.go
