package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/robertpelloni/hyperharness/agents"
	"github.com/robertpelloni/hyperharness/internal/controlplane"
)

// ─── Styles ───────────────────────────────────────────────────────────

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			MarginBottom(1)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).Bold(true)

	userStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA"))

	assistantStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA"))

	systemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B"))

	toolStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4B5563"))

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")).
			Background(lipgloss.Color("#1F2937")).
			Padding(0, 1)

	statusKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9333EA")).Bold(true)

	statusValStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D5DB"))

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#374151")).
			Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F9FAFB"))
)

// ─── Messages ─────────────────────────────────────────────────────────

type PromptDisplayMsg struct{ Display string }
type ShellProposalMsg struct {
	Command     string
	Explanation string
}
type ToolExecMsg struct {
	ToolName string
	Input    string
	Output   string
	Duration time.Duration
	Err      string
}
type StatusUpdateMsg struct {
	Provider string
	Model    string
	Tools    int
	Sessions int
}

// ─── Model ────────────────────────────────────────────────────────────

type model struct {
	director    *agents.Director
	input       string
	history     []string
	loading     bool
	spinner     spinner.Model
	quitting    bool

	// viewport for scrollable chat
	viewport viewport.Model
	ready    bool

	// foundation session
	foundationSessionID    string
	foundationTreeSelection []string

	// tree browser
	browserActive         bool
	browserItems          []TreeBrowserItem
	browserIndex          int
	browserFilter         string
	browserConfirmPending bool
	browserCollapsed      map[string]bool
	browserGrouped        bool

	// pinned pane
	browserPinned        bool
	browserPinnedFocus   bool
	browserPaneHeight    int
	browserPanePosition  string
	browserPanePreview   bool

	// dashboard
	dashboardActive bool

	// status
	workingDir  string
	provider    string
	modelName   string
	toolCount   int
	sessionCount int

	// multi-line input
	inputMultiline bool

	// width/height tracking
	width  int
	height int

	// tool execution tracking
	toolMu   sync.Mutex
	toolRuns []ToolExecMsg
}

func getWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func countInstalledTools() int {
	detector := controlplane.NewToolDetector(30*time.Second, 10*time.Minute)
	tools := detector.Detect()
	return len(tools)
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	wd := getWorkingDir()
	provider := agents.NewHyperCodeProvider()
	director := agents.NewDirector(provider)
	director.WorkingDir = wd

	toolCount := countInstalledTools()

	m := model{
		director:           director,
		input:              "",
		history:            []string{},
		loading:            false,
		spinner:            s,
		viewport:           viewport.New(80, 20),
		browserPaneHeight:  8,
		browserPanePosition: "top",
		browserPanePreview: true,
		workingDir:         wd,
		provider:           "hypercode",
		modelName:          "auto",
		toolCount:          toolCount,
	}

	// Welcome banner
	m.history = append(m.history, welcomeBanner(wd, toolCount))
	return m
}

func welcomeBanner(wd string, toolCount int) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("╔══════════════════════════════════════════════════════╗"))
	b.WriteString("\n")
	b.WriteString(titleStyle.Render("║          🧠  HyperHarness  —  AI Control Plane      ║"))
	b.WriteString("\n")
	b.WriteString(titleStyle.Render("╚══════════════════════════════════════════════════════╝"))
	b.WriteString("\n\n")

	b.WriteString(systemStyle.Render(fmt.Sprintf("  Working Directory:  %s", wd)))
	b.WriteString("\n")
	b.WriteString(systemStyle.Render(fmt.Sprintf("  Detected Tools:     %d CLI harnesses", toolCount)))
	b.WriteString("\n")
	b.WriteString(systemStyle.Render(fmt.Sprintf("  Provider:           HyperCode (local-first)")))
	b.WriteString("\n")
	b.WriteString(systemStyle.Render(fmt.Sprintf("  Model Routing:      auto (fallback chain)")))
	b.WriteString("\n\n")

	b.WriteString(dimStyle.Render("  ─── Quick Start ─────────────────────────────────────"))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  Type a prompt to chat with the AI director"))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  Type /help for all slash commands"))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  Prefix with ?? for shell command proposals"))
	b.WriteString("\n\n")

	b.WriteString(dimStyle.Render("  ─── Key Commands ────────────────────────────────────"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s  Dashboard (split-pane)", promptStyle.Render("/dashboard")))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s  File tree browser", promptStyle.Render("/tree-browser")))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s  Persistent tree pane", promptStyle.Render("/tree-pane")))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s  Repo map overview", promptStyle.Render("/repomap")))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s  Provider status", promptStyle.Render("/providers")))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s  MCP tool listing", promptStyle.Render("/mcp")))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s  Foundation plan", promptStyle.Render("/plan <prompt>")))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s  Clear history", promptStyle.Render("/clear")))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s  Exit", promptStyle.Render("/exit")))
	b.WriteString("\n\n")

	b.WriteString(dimStyle.Render("  ─────────────────────────────────────────────────────"))
	b.WriteString("\n")

	return b.String()
}

// ─── Init ─────────────────────────────────────────────────────────────

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, windowSizeCmd)
}

func windowSizeCmd() tea.Msg {
	return tea.WindowSizeMsg{}
}

// ─── Update ───────────────────────────────────────────────────────────

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.viewport = viewport.New(msg.Width, max(msg.Height-6, 5))
			m.viewport.SetContent(strings.Join(m.history, "\n"))
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = max(msg.Height-6, 5)
		}
		m.viewport.GotoBottom()
		return m, nil

	case tea.KeyMsg:
		// ─── Tree browser modal mode ────────────────────
		if m.browserActive {
			visible := visibleTreeBrowserItems(m.browserItems, m.browserFilter, m.browserCollapsed)
			switch msg.Type {
			case tea.KeyEsc:
				if m.browserConfirmPending {
					m.browserConfirmPending = false
					return m, nil
				}
				m.browserActive = false
				m.browserFilter = ""
				m.browserConfirmPending = false
				m.history = append(m.history, systemStyle.Render("[Tree Browser] closed"))
				m.syncViewport()
				return m, nil
			case tea.KeyLeft:
				if m.browserConfirmPending { return m, nil }
				if m.browserIndex >= 0 && m.browserIndex < len(visible) {
					item := visible[m.browserIndex]
					if item.ChildCount > 0 {
						if m.browserCollapsed == nil { m.browserCollapsed = map[string]bool{} }
						m.browserCollapsed[item.ID] = true
						visible = visibleTreeBrowserItems(m.browserItems, m.browserFilter, m.browserCollapsed)
						if m.browserIndex >= len(visible) { m.browserIndex = max(0, len(visible)-1) }
					}
				}
				return m, nil
			case tea.KeyRight:
				if m.browserConfirmPending { return m, nil }
				if m.browserIndex >= 0 && m.browserIndex < len(visible) {
					item := visible[m.browserIndex]
					if m.browserCollapsed != nil { delete(m.browserCollapsed, item.ID) }
				}
				return m, nil
			case tea.KeyHome:
				if !m.browserConfirmPending { m.browserIndex = 0 }
				return m, nil
			case tea.KeyEnd:
				if !m.browserConfirmPending { m.browserIndex = max(0, len(visible)-1) }
				return m, nil
			case tea.KeyPgUp:
				if !m.browserConfirmPending { m.browserIndex = max(0, m.browserIndex-10) }
				return m, nil
			case tea.KeyPgDown:
				if !m.browserConfirmPending { m.browserIndex = min(len(visible)-1, m.browserIndex+10) }
				return m, nil
			case tea.KeyUp:
				if !m.browserConfirmPending { m.browserIndex = max(0, m.browserIndex-1) }
				return m, nil
			case tea.KeyDown:
				if !m.browserConfirmPending { m.browserIndex = min(len(visible)-1, m.browserIndex+1) }
				return m, nil
			case tea.KeyEnter:
				if m.browserIndex >= 0 && m.browserIndex < len(visible) {
					if !m.browserConfirmPending {
						m.browserConfirmPending = true
						return m, nil
					}
					display, err := openSelectedTreeBrowser(m.director.WorkingDir, m.foundationSessionID, visible, m.browserIndex, 128)
					if err != nil {
						m.history = append(m.history, errorStyle.Render(fmt.Sprintf("[Error] %v", err)))
					} else {
						m.history = append(m.history, successStyle.Render(display))
					}
					m.browserActive = false
					m.browserFilter = ""
					m.browserConfirmPending = false
					m.syncViewport()
				}
				return m, nil
			default:
				if msg.Type == tea.KeyRunes || msg.Type == tea.KeySpace {
					m.browserFilter += msg.String()
					m.browserIndex = 0
				}
				if msg.Type == tea.KeyBackspace && len(m.browserFilter) > 0 {
					m.browserFilter = m.browserFilter[:len(m.browserFilter)-1]
					m.browserIndex = 0
				}
				return m, nil
			}
		}

		// ─── Pinned tree pane focus mode ────────────────
		if m.browserPinned && m.browserPinnedFocus {
			visible := visibleTreeBrowserItems(m.browserItems, m.browserFilter, m.browserCollapsed)
			switch msg.Type {
			case tea.KeyEsc:
				if m.browserConfirmPending { m.browserConfirmPending = false; return m, nil }
				m.browserPinnedFocus = false
				return m, nil
			case tea.KeyHome:
				m.browserIndex = 0; return m, nil
			case tea.KeyEnd:
				m.browserIndex = max(0, len(visible)-1); return m, nil
			case tea.KeyUp:
				m.browserIndex = max(0, m.browserIndex-1); return m, nil
			case tea.KeyDown:
				m.browserIndex = min(len(visible)-1, m.browserIndex+1); return m, nil
			case tea.KeyEnter:
				if m.browserIndex >= 0 && m.browserIndex < len(visible) {
					if !m.browserConfirmPending { m.browserConfirmPending = true; return m, nil }
					display, err := openSelectedTreeBrowser(m.director.WorkingDir, m.foundationSessionID, visible, m.browserIndex, 128)
					if err != nil {
						m.history = append(m.history, errorStyle.Render(fmt.Sprintf("[Error] %v", err)))
					} else {
						m.history = append(m.history, successStyle.Render(display))
					}
					m.browserConfirmPending = false
					refreshPinnedFoundationTreeBrowser(&m)
					m.syncViewport()
				}
				return m, nil
			case tea.KeyTab:
				m.browserPinnedFocus = false
				return m, nil
			}
			return m, nil
		}

		// ─── Normal chat input mode ──────────────────────
		switch msg.Type {
		case tea.KeyCtrlC:
			if m.inputMultiline {
				m.inputMultiline = false
				m.history = append(m.history, systemStyle.Render("[Multi-line input cancelled]"))
				m.syncViewport()
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyEsc:
			if m.inputMultiline {
				m.inputMultiline = false
				m.history = append(m.history, systemStyle.Render("[Multi-line input cancelled]"))
				m.syncViewport()
				return m, nil
			}
			if m.dashboardActive {
				m.dashboardActive = false
				m.history = append(m.history, systemStyle.Render("[Dashboard] closed"))
				m.syncViewport()
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyTab:
			if m.browserPinned && !m.browserPinnedFocus {
				m.browserPinnedFocus = true
				return m, nil
			}
			if m.browserPinned && m.browserPinnedFocus {
				m.browserPinnedFocus = false
				return m, nil
			}

		case tea.KeyEnter:
			if m.inputMultiline {
				m.input += "\n"
				return m, nil
			}
			if strings.TrimSpace(m.input) != "" {
				req := strings.TrimSpace(m.input)

				if strings.HasPrefix(req, "/") {
					m.input = ""
					mdl, cmd := ProcessSlashCommand(req, &m)
					m = mdl.(model)
					m.syncViewport()
					return m, cmd
				}

				if strings.HasPrefix(req, "??") {
					m.input = ""
					mdl, cmd := ProcessShellCommand(req, &m)
					m = mdl.(model)
					m.syncViewport()
					return m, cmd
				}

				m.history = append(m.history, userStyle.Render("You: ")+req)
				m.input = ""
				m.loading = true
				m.syncViewport()

				if sessionID, err := ensureFoundationSession(&m); err == nil {
					m.foundationSessionID = sessionID
					_ = appendFoundationUserText(m.director.WorkingDir, m.foundationSessionID, req)
					refreshPinnedFoundationTreeBrowser(&m)
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

		case tea.KeyPgUp:
			m.viewport.HalfViewUp()

		case tea.KeyPgDown:
			m.viewport.HalfViewDown()
		}

	// ─── Async responses ─────────────────────────────────
	case string:
		m.loading = false
		m.history = append(m.history, assistantStyle.Render("HyperHarness: ")+msg)
		if m.foundationSessionID != "" {
			_ = appendFoundationAssistantText(m.director.WorkingDir, m.foundationSessionID, msg)
			refreshPinnedFoundationTreeBrowser(&m)
		}
		m.syncViewport()

	case PromptDisplayMsg:
		m.loading = false
		m.history = append(m.history, assistantStyle.Render("HyperHarness: ")+msg.Display)
		if m.foundationSessionID != "" {
			_ = appendFoundationAssistantText(m.director.WorkingDir, m.foundationSessionID, msg.Display)
			refreshPinnedFoundationTreeBrowser(&m)
		}
		m.syncViewport()

	case ShellProposalMsg:
		m.loading = false
		display := fmt.Sprintf("%s %s\n%s", toolStyle.Render("[Shell Proposal]"), msg.Command, dimStyle.Render(msg.Explanation))
		display += "\n" + warningStyle.Render("  Execute? (Y/n)")
		m.history = append(m.history, display)
		m.syncViewport()

	case ToolExecMsg:
		m.toolMu.Lock()
		m.toolRuns = append(m.toolRuns, msg)
		if len(m.toolRuns) > 50 { m.toolRuns = m.toolRuns[len(m.toolRuns)-50:] }
		m.toolMu.Unlock()
		var display string
		if msg.Err != "" {
			display = errorStyle.Render(fmt.Sprintf("[Tool Error] %s: %s (%v)", msg.ToolName, msg.Err, msg.Duration))
		} else {
			output := msg.Output
			if len(output) > 500 { output = output[:500] + "..." }
			display = toolStyle.Render(fmt.Sprintf("[Tool] %s (%v)", msg.ToolName, msg.Duration)) + "\n" + dimStyle.Render(output)
		}
		m.history = append(m.history, display)
		m.syncViewport()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// syncViewport updates the viewport content and scrolls to bottom.
func (m *model) syncViewport() {
	m.viewport.SetContent(strings.Join(m.history, "\n"))
	m.viewport.GotoBottom()
}

// ─── View ─────────────────────────────────────────────────────────────

func (m model) View() string {
	if m.quitting {
		return dimStyle.Render("Goodbye! The collective grows.\n")
	}

	// Dashboard mode
	if m.dashboardActive {
		chatContent, toolContent, metrics := GenerateDashboardPlaceholders()
		chatContent = strings.Join(m.history, "\n")
		return RenderDashboard(chatContent, toolContent, metrics)
	}

	var body string

	// Tree browser modal
	if m.browserActive {
		body = m.renderChatArea() + "\n\n" +
			renderTreeBrowser(m.browserItems, m.browserIndex, m.browserFilter,
				m.browserConfirmPending, m.browserCollapsed, m.browserGrouped,
				0, promptStyle.Render("[File Tree Browser :: Modal]"), true)
		return body + "\n" + m.renderStatusBar() + "\n" + m.renderInputLine()
	}

	// Pinned tree pane
	if m.browserPinned {
		paneHeight := m.browserPaneHeight
		if paneHeight <= 0 { paneHeight = 8 }
		title := dimStyle.Render("[File Tree Pane :: Passive]")
		if m.browserPinnedFocus { title = promptStyle.Render("[File Tree Pane :: Focused]") }
		pane := renderTreeBrowser(m.browserItems, m.browserIndex, m.browserFilter,
			m.browserConfirmPending && m.browserPinnedFocus, m.browserCollapsed, m.browserGrouped,
			paneHeight, title, m.browserPanePreview)
		divider := dimStyle.Render("════════════════════════════════════════════════════════════")
		if strings.ToLower(strings.TrimSpace(m.browserPanePosition)) == "bottom" {
			body = m.renderChatArea() + "\n" + divider + "\n" + pane
		} else {
			body = pane + "\n" + divider + "\n" + m.renderChatArea()
		}
		return body + "\n" + m.renderStatusBar() + "\n" + m.renderInputLine()
	}

	// Normal mode
	body = m.renderChatArea()
	return body + "\n" + m.renderStatusBar() + "\n" + m.renderInputLine()
}

func (m model) renderChatArea() string {
	if m.ready {
		return m.viewport.View()
	}
	return strings.Join(m.history, "\n")
}

func (m model) renderInputLine() string {
	if m.loading {
		return fmt.Sprintf("%s %s Thinking...", m.spinner.View(), dimStyle.Render("⏳"))
	}

	prefix := promptStyle.Render("> ")
	inputText := inputStyle.Render(m.input)

	if m.inputMultiline {
		prefix = promptStyle.Render("... ")
	}

	cursor := "▎"
	return prefix + inputText + cursor
}

func (m model) renderStatusBar() string {
	wd := m.workingDir
	if len(wd) > 40 {
		wd = "..." + filepath.Base(filepath.Dir(wd)) + "/" + filepath.Base(wd)
	}

	provider := m.provider
	if provider == "" { provider = "none" }

	model := m.modelName
	if model == "" { model = "auto" }

	left := fmt.Sprintf("%s %s  %s %s  %s %d tools",
		statusKeyStyle.Render("Dir:"), statusValStyle.Render(wd),
		statusKeyStyle.Render("Provider:"), statusValStyle.Render(provider),
		statusKeyStyle.Render("Tools:"), m.toolCount)

	right := fmt.Sprintf("%s %s  %s %s",
		statusKeyStyle.Render("Model:"), statusValStyle.Render(model),
		statusKeyStyle.Render("Sessions:"), statusValStyle.Render(fmt.Sprintf("%d", m.sessionCount)))

	// Pad to fill width
	totalLen := len(stripANSI(left)) + len(stripANSI(right)) + 4
	padding := max(0, m.width-totalLen)
	gap := strings.Repeat(" ", padding)

	return statusBarStyle.Render(left + gap + right)
}

func stripANSI(s string) string {
	// Quick approximation: count only printable runes
	result := 0
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' { inEscape = false }
			continue
		}
		result++
	}
	// Return a string of that length for measurement
	return strings.Repeat("x", result)
}

// ─── Entry Point ──────────────────────────────────────────────────────

func StartREPL() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting TUI: %v\n", err)
		os.Exit(1)
	}
}

func max(a, b int) int {
	if a > b { return a }
	return b
}

func min(a, b int) int {
	if a < b { return a }
	return b
}