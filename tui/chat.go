package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/robertpelloni/hypercode/agents"
)

type model struct {
	director                *agents.Director
	input                   string
	history                 []string
	loading                 bool
	spinner                 spinner.Model
	foundationSessionID     string
	foundationTreeSelection []string
	browserActive           bool
	browserItems            []TreeBrowserItem
	browserIndex            int
	browserFilter           string
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		director: agents.NewDirector(agents.NewHyperCodeProvider()),
		input:    "",
		history:  []string{},
		loading:  false,
		spinner:  s,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.browserActive {
			visible := filterTreeBrowserItems(m.browserItems, m.browserFilter)
			switch msg.Type {
			case tea.KeyEsc:
				m.browserActive = false
				m.browserFilter = ""
				m.history = append(m.history, "[Foundation Tree Browser] closed")
				return m, nil
			case tea.KeyUp:
				if m.browserIndex > 0 {
					m.browserIndex--
				}
				return m, nil
			case tea.KeyDown:
				if m.browserIndex < len(visible)-1 {
					m.browserIndex++
				}
				return m, nil
			case tea.KeyBackspace, tea.KeyDelete:
				if len(m.browserFilter) > 0 {
					m.browserFilter = m.browserFilter[:len(m.browserFilter)-1]
					visible = filterTreeBrowserItems(m.browserItems, m.browserFilter)
					if m.browserIndex >= len(visible) {
						m.browserIndex = max(0, len(visible)-1)
					}
				}
				return m, nil
			case tea.KeyRunes, tea.KeySpace:
				m.browserFilter += msg.String()
				visible = filterTreeBrowserItems(m.browserItems, m.browserFilter)
				if m.browserIndex >= len(visible) {
					m.browserIndex = max(0, len(visible)-1)
				}
				return m, nil
			case tea.KeyEnter:
				if m.browserIndex >= 0 && m.browserIndex < len(visible) {
					display, err := openSelectedTreeBrowser(m.director.WorkingDir, m.foundationSessionID, visible, m.browserIndex, 128)
					if err != nil {
						m.history = append(m.history, fmt.Sprintf("[Error] tree browser switch failed: %v", err))
					} else {
						m.history = append(m.history, display)
					}
					m.browserActive = false
					m.browserFilter = ""
				}
				return m, nil
			}
		}
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if strings.TrimSpace(m.input) != "" {
				req := strings.TrimSpace(m.input)

				if strings.HasPrefix(m.input, "/") {
					m.input = ""
					mdl, cmd := ProcessSlashCommand(req, &m)
					m = mdl.(model)
					return m, cmd
				}

				if strings.HasPrefix(m.input, "??") {
					m.input = ""
					mdl, cmd := ProcessShellCommand(req, &m)
					m = mdl.(model)
					return m, cmd
				}

				m.history = append(m.history, "You: "+req)
				m.input = ""
				m.loading = true
				if sessionID, err := ensureFoundationSession(&m); err == nil {
					m.foundationSessionID = sessionID
					_ = appendFoundationUserText(m.director.WorkingDir, m.foundationSessionID, req)
				}
				cmds = append(cmds, func() tea.Msg {
					response, err := buildPromptResponse(m.director, req)
					if err != nil {
						return fmt.Sprintf("Error: %v", err)
					}
					return response
				})
				return m, tea.Batch(cmds...)
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
		m.history = append(m.history, "Borg-Go-Director: "+msg)
		if m.foundationSessionID != "" {
			_ = appendFoundationAssistantText(m.director.WorkingDir, m.foundationSessionID, msg)
		}

	case PromptDisplayMsg:
		m.loading = false
		m.history = append(m.history, "Borg-Go-Director: "+msg.Display)
		if m.foundationSessionID != "" {
			_ = appendFoundationAssistantText(m.director.WorkingDir, m.foundationSessionID, msg.Display)
		}

	case ShellProposalMsg:
		m.loading = false
		display := fmt.Sprintf("[Shell Proposal] %s", msg.Command)
		if strings.TrimSpace(msg.Explanation) != "" {
			display += fmt.Sprintf("\n[Route] %s", msg.Explanation)
		}
		display += "\n> Execute? (Y/n)"
		m.history = append(m.history, display)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := strings.Join(m.history, "\n")
	s += "\n\n"
	if m.browserActive {
		s += renderTreeBrowser(m.browserItems, m.browserIndex, m.browserFilter)
		return s
	}
	if m.loading {
		s += fmt.Sprintf("%s Processing neural inputs...\n", m.spinner.View())
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
