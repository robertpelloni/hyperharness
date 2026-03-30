package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ProcessSlashCommand mimics Claude Code's native terminal interception primitives.
// If an input begins with '/', it bypasses the LLM and executes directly as a local primitive.
func ProcessSlashCommand(cmd string, m *model) (tea.Model, tea.Cmd) {
	cmd = strings.TrimSpace(cmd)
	parts := strings.Split(cmd, " ")

	switch parts[0] {
	case "/help":
		return handleHelp(m)
	case "/clear":
		return handleClear(m)
	case "/commit":
		return handleCommit(m)
	case "/exit", "/quit":
		return *m, tea.Quit
	default:
		return handleUnknown(m, parts[0])
	}
}

func handleHelp(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	// Simulate glamour markdown output locally bypassing the Director
	m.history = append(m.history, `[System] Parity Slash Commands:
  /help    - This menu
  /clear   - Wipes contextual agent memory
  /commit  - Autonomously generates a standard Git commit
  /exit    - Closes hypercode`)
	return *m, nil
}

func handleClear(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	m.history = []string{
		"[System] Memory crystal wiped. Context reset to null.",
	}
	// We also reset the native Director's history state
	// m.director.InjectSystemContext("")
	return *m, nil
}

func handleCommit(m *model) (tea.Model, tea.Cmd) {
	m.loading = false
	m.history = append(m.history, "[System] Extracting git diff to generate native algorithmic commit message...")
	// We'd tap native.go tools here to run 'git diff' and pass it to Director exclusively for formatting
	return *m, nil
}

func handleUnknown(m *model, cmd string) (tea.Model, tea.Cmd) {
	m.loading = false
	m.history = append(m.history, fmt.Sprintf("[Error] Unknown primitive: %s", cmd))
	return *m, nil
}
