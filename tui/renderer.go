package tui

// ═══════════════════════════════════════════════════════════════════════
// renderer.go — pi-mono exact-match rendering
// Every visual element matches pi-coding-agent's InteractiveMode output
// ═══════════════════════════════════════════════════════════════════════

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/robertpelloni/hyperharness/internal/buildinfo"
)

// ═══════════════════════════════════════════════════════════════════════
// Welcome banner — pi-mono's showStartupNoticesIfNeeded + init header
// Shows: logo, version, keybinding hints (one per line), onboarding
// ═══════════════════════════════════════════════════════════════════════

func RenderWelcome(cwd, gitBranch, provider, model string, toolCount, regCount int, t Theme) string {
	var b strings.Builder

	// Logo line (pi-mono: theme.bold(theme.fg("accent", APP_NAME)) + theme.fg("dim", ` v${this.version}`))
	b.WriteString(t.AccentText("HyperHarness") + t.Dim(fmt.Sprintf(" v%s", buildinfo.Version)))
	b.WriteString("\n")

	// Keybinding hints (pi-mono: one hint per line, dim parentheses + accent key)
	hints := []struct{ key, desc string }{
		{"Ctrl+C", "to interrupt"},
		{"Ctrl+C", "to clear"},        // first clears, second exits
		{"Ctrl+C twice", "to exit"},   // double-tap exit
		{"Ctrl+D", "to exit (empty)"}, // exit when empty
		{"Ctrl+Z", "to suspend"},
		{"Ctrl+K", "to delete to end"},
		{"Ctrl+P", "to cycle thinking level"},
		{"Ctrl+N/ctrl+R", "to cycle models"},
		{"Ctrl+Shift+P", "to select model"},
		{"Ctrl+O", "to expand tools"},
		{"Ctrl+T", "to expand thinking"},
		{"Ctrl+E", "for external editor"},
		{"/", "for commands"},
		{"!", "to run bash"},
		{"!!", "to run bash (no context)"},
		{"Ctrl+F", "to queue follow-up"},
		{"Ctrl+G", "to edit all queued messages"},
		{"Ctrl+V", "to paste image"},
		{"drop files", "to attach"},
	}
	for _, h := range hints {
		b.WriteString("  " + t.Dim("(") + t.Fg(t.Accent, h.key) + t.Dim(") ") + t.Dim(h.desc) + "\n")
	}

	// Onboarding (pi-mono: theme.fg("dim", `Pi can explain its own features...`))
	b.WriteString("\n")
	b.WriteString(t.Dim("HyperHarness can explain its own features and look up its docs. Ask it how to use or extend HyperHarness."))
	b.WriteString("\n")

	return b.String()
}

// ═══════════════════════════════════════════════════════════════════════
// Entry dispatcher
// ═══════════════════════════════════════════════════════════════════════

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
// Spacer(1) + Markdown with userMessageBg/userMessageText
// ═══════════════════════════════════════════════════════════════════════

func renderUserEntry(e ChatEntry, t Theme) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.UserMessageText)).
		Background(lipgloss.Color(t.UserMessageBg)).
		Padding(0, 1)
	return "\n" + style.Render(e.Content)
}

// ═══════════════════════════════════════════════════════════════════════
// Assistant entry — pi-mono's AssistantMessageComponent
// Spacer(1) + optional thinking label + Markdown(content) + error handling
// ═══════════════════════════════════════════════════════════════════════

func renderAssistantEntry(e ChatEntry, t Theme) string {
	var parts []string

	// Thinking label (pi-mono: if hideThinkingBlock, show italic thinking label)
	if e.ThinkingLevel != "" && e.ThinkingLevel != "off" && e.Hidden {
		parts = append(parts, t.Italic(t.Fg(t.ThinkingCol, "Thinking...")))
	}

	// Content as markdown
	if e.Content != "" {
		parts = append(parts, "\n"+RenderSimpleMarkdown(e.Content, t))
	}

	// Error handling (pi-mono: stopReason === "error")
	if e.ErrorMessage != "" {
		parts = append(parts, t.ErrorText("Error: "+e.ErrorMessage))
	}

	return strings.Join(parts, "\n")
}

// ═══════════════════════════════════════════════════════════════════════
// Tool entry — pi-mono's ToolExecutionComponent
// Bordered panel with status icon, kind icon, title, args, output
// ═══════════════════════════════════════════════════════════════════════

func renderToolEntry(e ChatEntry, t Theme) string {
	// Status icon (pi-mono/goose: ○ pending, ◑ in_progress, ● completed, ✗ failed)
	statusIcon := t.StatusIcon(e.ToolStatus)
	if e.ToolStatus == "" {
		if e.ToolErr {
			statusIcon = t.ErrorText("✗")
		} else {
			statusIcon = t.SuccessText("✓")
		}
	}

	// Kind icon (pi-mono/goose: 📖 read, ✏️ edit, 🔍 search, ▶ execute)
	kind := e.ToolKind
	if kind == "" {
		kind = DetectToolKind(e.ToolName)
	}
	kindIcon := ToolKindIcon(kind)

	// Duration
	durStr := ""
	if e.ToolDur > 0 {
		durStr = t.Dim(fmt.Sprintf(" (%v)", e.ToolDur.Round(time.Millisecond)))
	}

	// Title line (pi-mono: theme.fg("toolTitle", theme.bold(this.toolName)))
	titleLine := statusIcon + " " + kindIcon + " " + t.ToolTitleStyled(e.ToolName) + durStr

	// Args preview (pi-mono: args display in tool panel)
	argsLine := ""
	if e.ToolArgs != "" {
		preview := e.ToolArgs
		if len(preview) > 80 {
			preview = preview[:77] + "…"
		}
		argsLine = t.Dim("  ▸ in: " + preview)
	}

	// Output (pi-mono: expanded shows full, collapsed shows preview + hint)
	outputLine := ""
	if e.Content != "" {
		if e.Expanded {
			output := e.Content
			if len(output) > 500 {
				output = output[:497] + "…"
			}
			outputLine = "\n" + t.Fg(t.ToolOutput, "  ▸ out: "+output)
		} else {
			preview := e.Content
			if len(preview) > 60 {
				preview = preview[:57] + "…"
			}
			outputLine = "\n" + t.Dim("  ▸ "+preview+" "+t.KeyHint("Ctrl+O", "expand"))
		}
	}

	// Border color (pi-mono: toolPendingBg/toolSuccessBg/toolErrorBg)
	borderColor := t.Border
	if e.ToolErr || e.ToolStatus == "failed" {
		borderColor = t.Cranberry
	} else if e.ToolStatus == "completed" || (!e.ToolErr && e.ToolStatus == "") {
		borderColor = t.Teal
	}

	// Combine in bordered panel (pi-mono: Box with rounded border)
	content := titleLine
	if argsLine != "" {
		content += "\n" + argsLine
	}
	if outputLine != "" {
		content += outputLine
	}

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1).
		Render(content)

	return "\n" + panel
}

// ═══════════════════════════════════════════════════════════════════════
// System entry — pi-mono's dim status messages
// ═══════════════════════════════════════════════════════════════════════

func renderSystemEntry(e ChatEntry, t Theme) string {
	return t.Dim("  " + e.Content)
}

// ═══════════════════════════════════════════════════════════════════════
// Thinking entry — pi-mono's collapsible thinking with level-colored border
// ═══════════════════════════════════════════════════════════════════════

func renderThinkingEntry(e ChatEntry, t Theme) string {
	level := e.ThinkingLevel
	if level == "" {
		level = "off"
	}
	borderColor := t.ThinkingBorder(level)

	if e.Hidden || level == "off" {
		label := t.Fg(borderColor, "💭 Thinking") + t.Dim(" ["+level+"]")
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
		Render("⌘ Shell Proposal\n"+cmdLine+"\n"+hint)
	return "\n" + panel
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
	return "\n" + panel
}

// ═══════════════════════════════════════════════════════════════════════
// Compaction summary — pi-mono's CompactionSummaryMessageComponent
// ═══════════════════════════════════════════════════════════════════════

func renderCompactionSummaryEntry(e ChatEntry, t Theme) string {
	line := t.Dim("📋 Compacted")
	if e.TokensBefore > 0 || e.TokensAfter > 0 {
		line += t.Dim(fmt.Sprintf(" %d → %d tokens", e.TokensBefore, e.TokensAfter))
	}
	return line
}

// ═══════════════════════════════════════════════════════════════════════
// Bash mode — pi-mono's BashExecutionComponent
// Bordered panel with $ command, loader/running indicator, output
// ═══════════════════════════════════════════════════════════════════════

func renderBashModeEntry(e ChatEntry, t Theme) string {
	// Pi-mono: bordered panel with $ command header
 borderColor := t.BashModeCol
 if e.ToolErr {
  borderColor = t.Cranberry
 }

 content := t.Fg(borderColor, t.Bold("$ "+e.Content))

 // Output (if completed)
 if e.Content != "" && e.ToolStatus == "completed" {
  output := e.ToolOut
  if output == "" {
   output = e.ToolArgs // fallback
  }
  if output != "" {
   lines := strings.Split(output, "\n")
   if len(lines) > 20 {
    lines = lines[:20]
    lines = append(lines, t.Dim(fmt.Sprintf("... (%d more lines)", len(strings.Split(output, "\n"))-20)))
   }
   content += "\n" + t.Fg(t.ToolOutput, strings.Join(lines, "\n"))
  }
 }

 // Status line
 statusLine := ""
 if e.ToolStatus == "completed" || e.ToolStatus == "failed" {
  statusLine = "\n" + t.Dim(fmt.Sprintf("exit code: 0 (%v)", e.ToolDur.Round(time.Millisecond)))
 } else if e.ToolStatus == "in_progress" || e.ToolStatus == "" {
  statusLine = "\n" + t.Fg(t.BashModeCol, "Running... (Esc to cancel)")
 }

 panel := lipgloss.NewStyle().
  Border(lipgloss.RoundedBorder()).
  BorderForeground(lipgloss.Color(borderColor)).
  Padding(0, 1).
  Render(content + statusLine)

 return "\n" + panel
}

// ═══════════════════════════════════════════════════════════════════════
// Permission entry — goose-style bordered permission dialog
// ═══════════════════════════════════════════════════════════════════════

func renderPermissionEntry(e ChatEntry, t Theme) string {
	pc := t.PermissionColor()
	hRule := strings.Repeat("─", 56)

	var lines []string
	lines = append(lines, t.Fg(pc, "╭"+hRule+"╮"))
	lines = append(lines, t.Fg(pc, "│ ")+t.Bold(t.Fg(pc, "🔒 Permission required"))+t.Fg(pc, " │"))
	lines = append(lines, t.Fg(pc, "│ ")+" "+t.Fg(pc, " │"))
	lines = append(lines, t.Fg(pc, "│ ")+t.Fg(t.TextColor, e.Content)+t.Fg(pc, " │"))
	lines = append(lines, t.Fg(pc, "│ ")+" "+t.Fg(pc, " │"))

	for i, opt := range e.PermissionOptions {
		icon := " "
		if i == e.PermissionIdx {
			icon = "▸"
		}
		color := t.DimColor
		if i == e.PermissionIdx {
			color = t.TextColor
		}
		keyHint := ""
		switch opt.Kind {
		case "allow_once":
			keyHint = "y"
		case "allow_always":
			keyHint = "a"
		case "reject_once":
			keyHint = "n"
		case "reject_always":
			keyHint = "N"
		}
		lines = append(lines, t.Fg(pc, "│ ")+t.Fg(color, icon+" ["+keyHint+"] "+opt.Name)+t.Fg(pc, " │"))
	}

	lines = append(lines, t.Fg(pc, "│ ")+" "+t.Fg(pc, " │"))
	lines = append(lines, t.Fg(pc, "│ ")+t.Dim("↑↓ select · enter confirm · esc cancel")+t.Fg(pc, " │"))
	lines = append(lines, t.Fg(pc, "╰"+hRule+"╯"))

	return "\n" + strings.Join(lines, "\n")
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
	if label == "" {
		label = "custom"
	}
	bg := e.CustomBg
	if bg == "" {
		bg = t.CustomMessageBg
	}
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
	if msg == "" {
		msg = e.Content
	}
	return t.ErrorText("⚠ Error: ") + t.ErrorText(msg)
}

// ═══════════════════════════════════════════════════════════════════════
// Queue entry — goose-style queued message indicator
// ═══════════════════════════════════════════════════════════════════════

func renderQueueEntry(e ChatEntry, t Theme) string {
	return t.Fg(t.DimColor, "> ") + t.Dim(e.Content) + t.Fg(t.Gold, " (queued)")
}

// ═══════════════════════════════════════════════════════════════════════
// Markdown rendering — pi-mono's Markdown component
// Handles: headings, code blocks, blockquotes, lists, hr, inline code/bold/links
// ═══════════════════════════════════════════════════════════════════════

func RenderSimpleMarkdown(text string, t Theme) string {
	if text == "" {
		return ""
	}
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
		if idx < 0 {
			break
		}
		start := strings.LastIndex(text[:idx], "[")
		if start < 0 {
			break
		}
		end := strings.Index(text[idx+2:], ")")
		if end < 0 {
			break
		}
		label := text[start+1 : idx]
		url := text[idx+2 : idx+2+end]
		text = text[:start] + t.Fg(t.MdLink, label) + t.Dim(" ("+url+")") + text[idx+2+end+1:]
	}
	return text
}

// ═══════════════════════════════════════════════════════════════════════
// Footer — pi-mono's FooterComponent (2-line footer)
// Line 1: pwd (git branch) │ ↑in/↓out │ $cost │ ctx %/window(auto)
// Line 2: (provider) model • thinking level │ tools N │ Ctrl+C to interrupt
// ═══════════════════════════════════════════════════════════════════════

func RenderFooter(cwd, gitBranch, provider, model string, toolCount, inputTok, outTok int, cost float64, contextPct float64, contextWindow int, thinkingLevel string, t Theme, width int) string {
	if width <= 0 {
		width = 80
	}

	// Line 1: pwd + git + tokens + cost + context
	pwd := shortenPath(cwd)
	if gitBranch != "" {
		pwd = pwd + " (" + gitBranch + ")"
	}

	// Token stats (pi-mono: ↑in ↓out $cost)
	var statsParts []string
	if inputTok > 0 {
		statsParts = append(statsParts, "↑"+formatTokens(inputTok))
	}
	if outTok > 0 {
		statsParts = append(statsParts, "↓"+formatTokens(outTok))
	}
	if cost > 0 {
		statsParts = append(statsParts, fmt.Sprintf("$%.3f", cost))
	}

	// Context percentage (pi-mono: color-coded, with (auto) indicator)
	ctxColor := t.DimColor
	if contextPct > 90 {
		ctxColor = t.Error
	} else if contextPct > 70 {
		ctxColor = t.Warning
	}
	ctxStr := fmt.Sprintf("%.1f%%/%s(auto)", contextPct, formatTokens(contextWindow))
	if contextPct == 0 {
		ctxStr = fmt.Sprintf("?/%s(auto)", formatTokens(contextWindow))
	}
	statsParts = append(statsParts, t.Fg(ctxColor, ctxStr))

	statsLeft := strings.Join(statsParts, " ")

	// Line 1: pwd on left, stats on right (pi-mono: right-aligned)
	pwdLine := t.Dim(pwd)
	if utf8.RuneCountInString(pwdLine)+2+utf8.RuneCountInString(statsLeft) > width {
		// Truncate pwd
		maxPwd := width - utf8.RuneCountInString(statsLeft) - 5
		if maxPwd > 0 {
			pwdRunes := []rune(pwd)
			if len(pwdRunes) > maxPwd {
				pwd = string(pwdRunes[:maxPwd-3]) + "..."
			}
		}
		pwdLine = t.Dim(pwd)
	}
	padding := width - utf8.RuneCountInString(pwdLine) - utf8.RuneCountInString(statsLeft)
	if padding < 2 {
		padding = 2
	}
	line1 := pwdLine + strings.Repeat(" ", padding) + t.Dim(statsLeft)

	// Line 2: model on right, thinking on left (pi-mono: right-aligned model)
	rightSide := model
	if thinkingLevel != "" && thinkingLevel != "off" {
		rightSide = model + " • " + thinkingLevel
	}
	if provider != "" {
		rightSide = "(" + provider + ") " + rightSide
	}

	leftSide := t.Dim(fmt.Sprintf("tools %d", toolCount))
	padding2 := width - utf8.RuneCountInString(leftSide) - utf8.RuneCountInString(rightSide)
	if padding2 < 2 {
		padding2 = 2
	}
	line2 := leftSide + strings.Repeat(" ", padding2) + t.Dim(rightSide)

	// Horizontal rule
	hr := t.Dim("─" + strings.Repeat("─", width-1))

	return hr + "\n" + line1 + "\n" + line2
}

// ═══════════════════════════════════════════════════════════════════════
// Input bar — pi-mono's editor with mode indicators
// ═══════════════════════════════════════════════════════════════════════

func RenderInputBar(m model) string {
	t := m.theme
	var prompt, inputDisplay string

	if m.loading {
		prompt = t.Bold(t.Fg(t.ThinkingCol, "… "))
		inputDisplay = t.Fg(t.ThinkingCol, t.Italic(m.thinkingLabel))
		if len(m.followUpQueue) > 0 {
			inputDisplay += t.Fg(t.Gold, t.Italic(fmt.Sprintf(" [%d queued]", len(m.followUpQueue))))
		}
		inputDisplay += " " + m.spinner.View()
	} else if m.pasteMode {
		prompt = t.Bold(t.Fg(t.BashModeCol, "📋 "))
		preview := m.pasteContent
		if len(preview) > 60 {
			preview = preview[:57] + "…"
		}
		charCount := len(m.pasteContent)
		inputDisplay = t.Fg(t.TextColor, preview) + t.Dim(fmt.Sprintf(" (%d chars)", charCount))
		inputDisplay += "\n" + t.Dim(t.Italic("enter to send · esc to clear"))
	} else if m.bashMode || strings.HasPrefix(m.input, "!") {
		prompt = t.Bold(t.Fg(t.BashModeCol, "⚡ "))
		inputDisplay = t.Fg(t.TextColor, m.input) + "▎"
	} else {
		prompt = t.Bold(t.AccentText("> "))
		inputDisplay = t.Fg(t.TextColor, m.input) + "▎"
	}

	// Queue hint (pi-mono: message queued — will send when finished)
	queueHint := ""
	if m.loading && len(m.followUpQueue) > 0 {
		queueHint = "\n" + t.Fg(t.Gold, t.Italic("message queued — will send when agent finishes"))
	}

	// Autocomplete dropdown
	autocomplete := ""
	if m.showAutocomplete && m.isSlashContext() {
		autocomplete = "\n" + RenderAutocomplete(m.filteredAutocomplete(), m.autocompleteIndex, m.autocompleteMaxVis, t)
	}

	return prompt + inputDisplay + queueHint + autocomplete
}

// ═══════════════════════════════════════════════════════════════════════
// Autocomplete dropdown — pi-mono's CombinedAutocompleteProvider
// ═══════════════════════════════════════════════════════════════════════

func RenderAutocomplete(items []SlashCommand, selected, maxVisible int, t Theme) string {
	if len(items) == 0 {
		return ""
	}
	if maxVisible <= 0 {
		maxVisible = 8
	}

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
		if i < len(chatLines) {
			left = chatLines[i]
		}
		right := ""
		if i < len(toolLines) {
			right = toolLines[i]
		}
		// Pad left to leftW
		leftRunes := utf8.RuneCountInString(left)
		if leftRunes < leftW {
			left += strings.Repeat(" ", leftW-leftRunes)
		} else if leftRunes > leftW {
			left = string([]rune(left)[:leftW])
		}
		rightRunes := utf8.RuneCountInString(right)
		if rightRunes > rightW {
			right = string([]rune(right)[:rightW])
		}
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
	b.WriteString("  " + t.KeyHint("Shift+Enter", "new line") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+E", "external editor") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+C", "interrupt / quit") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+D", "exit (empty)") + "\n")
	b.WriteString("  " + t.KeyHint("Esc Esc", "tree/fork selector") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+Z", "suspend to shell") + "\n")
	b.WriteString("\n")
	b.WriteString(t.Dim("── Agent ─────────────────────────────────────\n"))
	b.WriteString("  " + t.KeyHint("Ctrl+P", "cycle thinking level") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+N", "cycle model forward") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+R", "cycle model backward") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+Shift+P", "model selector") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+O", "toggle tool expansion") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+T", "toggle thinking visibility") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+F", "queue follow-up message") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+G", "restore queued messages") + "\n")
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
	b.WriteString("\n")
	b.WriteString(t.Dim("── Overlays ──────────────────────────────────\n"))
	b.WriteString("  " + t.KeyHint("Ctrl+K", "command palette") + "\n")
	b.WriteString("  " + t.KeyHint("Ctrl+V", "paste image from clipboard") + "\n")
	return b.String()
}

// ═══════════════════════════════════════════════════════════════════════
// Provider selector — pi-mono style (up/down to navigate, enter to select)
// ═══════════════════════════════════════════════════════════════════════

func RenderProviderSelector(current string, providers []string, selectedIdx int, t Theme) string {
	var b strings.Builder
	b.WriteString(t.Bold(t.AccentText("Provider Selector")) + "\n\n")
	for i, p := range providers {
		icon := "  "
		if i == selectedIdx {
			icon = t.SuccessText("> ")
		}
		b.WriteString(icon + t.Fg(t.TextColor, p))
		if p == current {
			b.WriteString(t.Dim(" (active)"))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n" + t.Dim("Up/Down to navigate, Enter to select, Esc to cancel"))
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
		if m == current {
			icon = t.SuccessText("▸ ")
		}
		b.WriteString(icon + t.Fg(t.TextColor, m))
		if m == current {
			b.WriteString(t.Dim(" (active)"))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n" + t.Dim("Ctrl+N / Ctrl+R to cycle, /model <name> to set"))
	return b.String()
}

// ═══════════════════════════════════════════════════════════════════════
// Helpers
// ═══════════════════════════════════════════════════════════════════════

// Format file size for display
func formatFileSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1fMB", float64(bytes)/(1024*1024))
}

// Format path for display (replace home with ~)
func shortenPathForDisplay(path string) string {
	home, _ := os.UserHomeDir()
	if home != "" && strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}

// TruncateString truncates a string to a maximum width, adding ellipsis if needed
func TruncateString(s string, maxWidth int) string {
	runes := []rune(s)
	if len(runes) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return string(runes[:maxWidth])
	}
	return string(runes[:maxWidth-3]) + "..."
}

// VisibleWidth returns the visible width of a string (excluding ANSI escape sequences)
func VisibleWidth(s string) int {
	// Simple implementation - count runes, ignoring ANSI sequences
	inEscape := false
	width := 0
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		width++
	}
	return width
}

// PadRight pads a string to the right with spaces to reach the target width
func PadRight(s string, width int) string {
	currentWidth := VisibleWidth(s)
	if currentWidth >= width {
		return s
	}
	return s + strings.Repeat(" ", width-currentWidth)
}

// PadLeft pads a string to the left with spaces to reach the target width
func PadLeft(s string, width int) string {
	currentWidth := VisibleWidth(s)
	if currentWidth >= width {
		return s
	}
	return strings.Repeat(" ", width-currentWidth) + s
}

// Center centers a string within the given width
func Center(s string, width int) string {
	currentWidth := VisibleWidth(s)
	if currentWidth >= width {
		return s
	}
	leftPad := (width - currentWidth) / 2
	rightPad := width - currentWidth - leftPad
	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}
