package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/robertpelloni/hyperharness/internal/subagents"
)

// The dashboard package handles rendering real-time observability telemetry.
// Based on Phase 6 specifications, this dashboard visualizes subagent tracking,
// memory connections, and MCP states.

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#383838")).
			Padding(1, 2)

	statusRunning = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render("● Running")
	statusFailed  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4D4D")).Render("● Failed")
	statusDone    = lipgloss.NewStyle().Foreground(lipgloss.Color("#5694F4")).Render("● Completed")

	docStyle      = lipgloss.NewStyle().Margin(1, 2)
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

// RenderDashboardLipgloss renders a split-pane dashboard using lipgloss join.
// The simple text-based RenderDashboard is in renderer.go for fallback.
func RenderDashboardLipgloss(chatContent, toolContent, metricsContent string) string {
	chatWidth := 80
	toolWidth := 40
	left := chatPaneStyle.Width(chatWidth).Render(chatContent)
	right := toolPaneStyle.Width(toolWidth).Render(toolContent)
	splitView := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	footer := footerDashStyle.Render(metricsContent)
	return docStyle.Render(lipgloss.JoinVertical(lipgloss.Left, splitView, footer))
}

type DashboardModel struct {
	width  int
	height int
}

func NewDashboard() DashboardModel {
	return DashboardModel{}
}

func (m DashboardModel) Init() tea.Cmd {
	return nil
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m DashboardModel) View() string {
	// Fetch real-time subagents
	tasks := subagents.GlobalManager.ListTasks()

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Borg Control Plane") + "\n\n")

	// Render Active Subagents Panel
	sb.WriteString("Active Subagents:\n")
	if len(tasks) == 0 {
		sb.WriteString("  No active tasks.\n")
	} else {
		for i, t := range tasks {
			if i > 5 {
				sb.WriteString(fmt.Sprintf("  ... and %d more\n", len(tasks)-6))
				break
			}
			statusStr := statusDone
			if t.Status == "running" {
				statusStr = statusRunning
			} else if t.Status == "failed" {
				statusStr = statusFailed
			}

			sb.WriteString(fmt.Sprintf("  [%s] %s - %s\n", statusStr, string(t.Type), t.ID))
		}
	}

	panel := panelStyle.Render(sb.String())
	return panel
}

// GenerateDashboardPlaceholders provides initial content.
// Now it overrides the toolContent with the new real-time DashboardModel view.
func GenerateDashboardPlaceholders() (string, string, string) {
	t := DefaultTheme
	var chat strings.Builder
	chat.WriteString(t.AccentText("System: Dashboard Initialized") + "\n")
	chat.WriteString(t.Dim("> Ready for commands...") + "\n")

	// Real-time Dashboard
	dash := NewDashboard()
	toolsView := dash.View()

	metrics := "Tokens: 0 | Cost: $0.000 | Scope: Project"
	return chat.String(), toolsView, metrics
}
