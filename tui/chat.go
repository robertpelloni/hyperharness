package tui

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/robertpelloni/hypercode/agent"
)

type model struct {
	agent   *agent.Agent
	input   string
	history []string
	loading bool
}

func initialModel() model {
	return model{
		agent:   agent.NewAgent(),
		input:   "",
		history: []string{},
		loading: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if strings.TrimSpace(m.input) != "" {
				m.history = append(m.history, "You: "+m.input)
				req := m.input
				m.input = ""
				m.loading = true
				return m, func() tea.Msg {
					response, err := m.agent.Chat(req)
					if err != nil {
						return fmt.Sprintf("Error: %v", err)
					}
					return response
				}
			}
		case tea.KeyBackspace, tea.KeyDelete:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case tea.KeyRunes, tea.KeySpace:
			m.input += msg.String()
		}
	case string:
		m.loading = false
		m.history = append(m.history, "SuperCLI: "+msg)
	}
	return m, nil
}

func (m model) View() string {
	s := strings.Join(m.history, "\n")
	s += "\n\n"
	if m.loading {
		s += "Thinking...\n"
	} else {
		s += "> " + m.input
	}
	return s
}

func StartREPL() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
