package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ═══════════════════════════════════════════════════════════════════════
// Dashboard — split-pane view mirroring pi-mono's layout
// ═══════════════════════════════════════════════════════════════════════

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)

	chatPaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(DefaultTheme.BorderAccent)).
			Padding(0, 1)

	toolPaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(DefaultTheme.Border)).
			Padding(0, 1)

	footerDashStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(DefaultTheme.DimColor)).
				Background(lipgloss.Color("#1F2937")).
				Padding(0, 1).
				MarginTop(1)
)

// RenderDashboard renders a split-pane dashboard view.
// Left: chat viewport, Right: tool sidebar, Bottom: metrics
func RenderDashboard(chatContent, toolContent, metricsContent string) string {
	// Calculate widths: 60/40 split
	chatWidth := 80
	toolWidth := 40

	left := chatPaneStyle.Width(chatWidth).Render(chatContent)
	right := toolPaneStyle.Width(toolWidth).Render(toolContent)
	splitView := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	footer := footerDashStyle.Render(metricsContent)
	return docStyle.Render(lipgloss.JoinVertical(lipgloss.Left, splitView, footer))
}

// GenerateDashboardPlaceholders provides initial content.
func GenerateDashboardPlaceholders() (string, string, string) {
	t := DefaultTheme
	var chat strings.Builder
	chat.WriteString(t.AccentText("System: Dashboard Initialized") + "\n")
	chat.WriteString(t.Dim("> Ready for commands...") + "\n")

	var tools strings.Builder
	tools.WriteString(t.Bold(t.AccentText("Active Tools")) + "\n")
	tools.WriteString(t.Fg(t.ToolTitle, "- Memory Search") + "\n")
	tools.WriteString(t.Fg(t.ToolTitle, "- Agent Delegation") + "\n")
	tools.WriteString(t.Fg(t.ToolTitle, "- Code Execution") + "\n")

	metrics := "Tokens: 0 | Cost: $0.000 | Scope: Project"
	return chat.String(), tools.String(), metrics
}
