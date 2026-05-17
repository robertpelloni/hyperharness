package tui

import (
	"fmt"
	"strings"
	"time"
	tea "github.com/charmbracelet/bubbletea"
)

// ProcessShellCommand mimics Copilot CLI's `??` bash translation interceptor.
func ProcessShellCommand(cmd string, m *model) (tea.Model, tea.Cmd) {
	query := strings.TrimSpace(strings.TrimPrefix(cmd, "??"))
	m.entries = append(m.entries, ChatEntry{
		Type:      EntryUser,
		Content:   "?? " + query,
		Timestamp: time.Now(),
	})
	m.loading = true
	return *m, func() tea.Msg {
		response, err := buildShellProposal(m.director, query)
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
		return response
	}
}
