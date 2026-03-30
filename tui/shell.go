package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/robertpelloni/hypercode/agents"
)

// ProcessShellCommand mimics Copilot CLI's `??` bash translation interceptor.
func ProcessShellCommand(cmd string, m *model) (tea.Model, tea.Cmd) {
	query := strings.TrimSpace(strings.TrimPrefix(cmd, "??"))

	m.history = append(m.history, fmt.Sprintf("You: ?? %s", query))

	// Create the explicit shell assistant
	assistant := agents.NewShellTranslator(m.director.Provider)

	m.loading = true
	var cmds []tea.Cmd

	cmds = append(cmds, func() tea.Msg {
		// Native context abstraction request
		response, err := assistant.Translate(nil, query)
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}

		return ShellProposalMsg{
			Command:     response,
			Explanation: "Automatically derived via shell constraints.",
		}
	})

	return *m, tea.Batch(cmds...)
}

// ShellProposalMsg represents a generated shell string waiting for user approval
type ShellProposalMsg struct {
	Command     string
	Explanation string
}
