package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	
	paneStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(40)
		
	mainPaneStyle = paneStyle.Copy().Width(80)
	
	footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Background(lipgloss.Color("235")).
		Padding(0, 1).
		MarginTop(1)
)

// RenderDashboard renders a fully split dashboard view based on current state.
func RenderDashboard(chatContent, toolContent, metricsContent string) string {
	left := mainPaneStyle.Render(chatContent)
	right := paneStyle.Render(toolContent)
	
	splitView := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	footer := footerStyle.Render(metricsContent)
	
	return docStyle.Render(lipgloss.JoinVertical(lipgloss.Left, splitView, footer))
}

// GenerateDashboardPlaceholders provides mock content while the backend wires up.
func GenerateDashboardPlaceholders() (string, string, string) {
	var chat strings.Builder
	chat.WriteString("System: HyperHarness Dashboard Initialized.\n")
	chat.WriteString("> Ready for commands...\n")
	
	var tools strings.Builder
	tools.WriteString("Active Tools:\n")
	tools.WriteString("- Memory Search\n")
	tools.WriteString("- Agent Delegation\n")
	
	metrics := "Tokens: 140 | Cost: $0.001 | Scope: Project"
	
	return chat.String(), tools.String(), metrics
}
