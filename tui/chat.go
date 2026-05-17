package tui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/robertpelloni/hyperharness/agents"
	"github.com/robertpelloni/hyperharness/foundation/adapters"
	foundationorchestration "github.com/robertpelloni/hyperharness/foundation/orchestration"
	"github.com/robertpelloni/hyperharness/internal/controlplane"
	"context"
	"github.com/robertpelloni/hyperharness/tools"
)

// ═══════════════════════════════════════════════════════════════════════
// Theme — mirrors pi-mono's ThemeColor system
// ═══════════════════════════════════════════════════════════════════════

type Theme struct {
	Accent       string
	Border       string
	BorderAccent string
	Success      string
	Error        string
	Warning      string
	Muted        string
	DimColor     string
	TextColor    string
	ThinkingCol  string

	UserMsgBg    string
	UserMsgText  string
	ToolPending  string
	ToolSuccess  string
	ToolError    string
	ToolTitle    string
	ToolOutput   string

	MDHeading    string
	MDLink       string
	MDCode       string
	MDCodeBlock  string
	MDQuote      string
}

var DefaultTheme = Theme{
	Accent:       "#7C3AED",
	Border:       "#374151",
	BorderAccent: "#9333EA",
	Success:      "#10B981",
	Error:        "#EF4444",
	Warning:      "#F59E0B",
	Muted:        "#9CA3AF",
	DimColor:     "#6B7280",
	TextColor:    "#F9FAFB",
	ThinkingCol:  "#8B5CF6",

	UserMsgBg:   "#1E293B",
	UserMsgText: "#93C5FD",
	ToolPending: "#1E1B4B",
	ToolSuccess: "#022C22",
	ToolError:   "#450A0A",
	ToolTitle:   "#F59E0B",
	ToolOutput:  "#9CA3AF",

	MDHeading:   "#C4B5FD",
	MDLink:      "#60A5FA",
	MDCode:      "#FCD34D",
	MDCodeBlock: "#1F2937",
	MDQuote:     "#6B7280",
}

func (t Theme) Fg(color string, text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(text)
}

func (t Theme) Bg(bgColor, text string) string {
	return lipgloss.NewStyle().Background(lipgloss.Color(bgColor)).Render(text)
}

func (t Theme) Bold(text string) string {
	return lipgloss.NewStyle().Bold(true).Render(text)
}

func (t Theme) Italic(text string) string {
	return lipgloss.NewStyle().Italic(true).Render(text)
}

func (t Theme) Dim(text string) string {
	return t.Fg(t.DimColor, text)
}

func (t Theme) AccentText(text string) string {
	return t.Fg(t.Accent, text)
}

func (t Theme) SuccessText(text string) string {
	return t.Fg(t.Success, text)
}

func (t Theme) ErrorText(text string) string {
	return t.Fg(t.Error, text)
}

func (t Theme) WarningText(text string) string {
	return t.Fg(t.Warning, text)
}

// ═══════════════════════════════════════════════════════════════════════
// Messages
// ═══════════════════════════════════════════════════════════════════════

type PromptDisplayMsg struct{ Display string }
type ShellProposalMsg struct {
	Command     string
	Explanation string
}
type ToolExecMsg struct {
	ToolName  string
	Args      string
	Output    string
	Duration  time.Duration
	IsError   bool
	Streaming bool
}
type StreamChunkMsg struct {
	Content string
	Done    bool
}
type AgentResponseMsg struct {
	Content  string
	Provider string
	Model    string
	Plan     *foundationorchestration.PlanResult
}
type GitBranchMsg struct {
	Branch string
}

// ═══════════════════════════════════════════════════════════════════════
// Chat Entry — structured message history (like pi-mono's session entries)
// ═══════════════════════════════════════════════════════════════════════

type EntryType int

const (
	EntryUser EntryType = iota
	EntryAssistant
	EntryTool
	EntrySystem
	EntryThinking
	EntryShellProposal
	EntryDiff
)

type ChatEntry struct {
	Type      EntryType
	Content   string
	ToolName  string
	ToolArgs  string
	ToolOut   string
	ToolDur   time.Duration
	ToolErr   bool
	Expanded  bool
	Provider  string
	Model     string
	Timestamp time.Time
}

// ═══════════════════════════════════════════════════════════════════════
// Model
// ═══════════════════════════════════════════════════════════════════════

type model struct {
	// Core
	director *agents.Director
	theme    Theme

	// Chat history
	entries []ChatEntry
	input   string

	// Viewport
	viewport viewport.Model
	ready    bool
	width    int
	height   int

	// State
	loading  bool
	quitting bool
	spinner  spinner.Model

	// Foundation session
	foundationSessionID     string
	foundationTreeSelection []string

	// Tree browser
	browserActive         bool
	browserItems          []TreeBrowserItem
	browserIndex          int
	browserFilter         string
	browserConfirmPending bool
	browserCollapsed      map[string]bool
	browserGrouped        bool

	// Pinned pane
	browserPinned       bool
	browserPinnedFocus  bool
	browserPaneHeight   int
	browserPanePosition string
	browserPanePreview  bool

	// Dashboard
	dashboardActive bool

	// Autocomplete
	showAutocomplete    bool
	autocompleteItems   []SlashCommand
	autocompleteIndex   int
	autocompleteMaxVis  int

	// Status bar data
	workingDir    string
	gitBranch     string
	provider      string
	modelName     string
	toolCount     int
	sessionCount  int
	totalInputTok int
	totalOutTok   int
	totalCost     float64
	contextPct    float64
	contextWindow int

	// Multi-line editor
	inputLines   []string
	cursorLine   int
	cursorCol    int
	inputHistory []string
	historyIdx   int
	multiline    bool

	// Tool registry
	registry *tools.Registry

	// Thinking
	thinkingLabel string
	hidingThink   bool

	// Tool execution mutex
	toolMu   sync.Mutex
	toolRuns []ToolExecMsg

	// Streaming
	streaming      bool
	streamContent  string
}

// ═══════════════════════════════════════════════════════════════════════
// Slash commands (mirrors pi-mono's BUILTIN_SLASH_COMMANDS)
// ═══════════════════════════════════════════════════════════════════════

type SlashCommand struct {
	Name        string
	Description string
}

var BuiltinSlashCommands = []SlashCommand{
	{Name: "help", Description: "Show all slash commands"},
	{Name: "clear", Description: "Clear chat history"},
	{Name: "compact", Description: "Compact session context"},
	{Name: "dashboard", Description: "Toggle split-pane dashboard"},
	{Name: "tree", Description: "Navigate session tree"},
	{Name: "tree-browser", Description: "Open file tree browser (modal)"},
	{Name: "tree-pane", Description: "Toggle persistent file tree pane"},
	{Name: "tree-pane-help", Description: "Show tree pane controls"},
	{Name: "repomap", Description: "Generate repository map"},
	{Name: "providers", Description: "Show LLM provider status"},
	{Name: "adapters", Description: "Show adapter status (HyperCode + MCP)"},
	{Name: "mcp", Description: "Show MCP tool listing"},
	{Name: "model", Description: "Show or select model"},
	{Name: "settings", Description: "Open settings menu"},
	{Name: "plan", Description: "Build orchestration plan"},
	{Name: "commit", Description: "Generate git commit message"},
	{Name: "session", Description: "Show session info and stats"},
	{Name: "name", Description: "Set session display name"},
	{Name: "fork", Description: "Create a fork from a previous message"},
	{Name: "export", Description: "Export session to file"},
	{Name: "import", Description: "Import a session from file"},
	{Name: "login", Description: "Login with OAuth provider"},
	{Name: "logout", Description: "Logout from OAuth provider"},
	{Name: "hotkeys", Description: "Show keyboard shortcuts"},
	{Name: "fsession", Description: "Show foundation session info"},
	{Name: "tools", Description: "List all available tools"},
	{Name: "new", Description: "Start a new session"},
	{Name: "resume", Description: "Resume a different session"},
	{Name: "reload", Description: "Reload configuration"},
	{Name: "quit", Description: "Quit HyperHarness"},
	{Name: "exit", Description: "Quit HyperHarness"},
}

// ═══════════════════════════════════════════════════════════════════════
// Initialization
// ═══════════════════════════════════════════════════════════════════════

func getWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func countInstalledTools() int {
	detector := controlplane.NewDetector(30*time.Second, 10*time.Minute)
	detected, _ := detector.DetectAll(context.Background())
	return len(detected)
}

func getGitBranch(wd string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = wd
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func shortenPath(path string) string {
	home, _ := os.UserHomeDir()
	if home != "" && strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(DefaultTheme.Accent))

	wd := getWorkingDir()
	provider := agents.NewHyperCodeProvider()
	director := agents.NewDirector(provider)
	director.WorkingDir = wd

	reg := tools.NewRegistry()
	toolCount := len(reg.Tools)
	cliToolCount := countInstalledTools()
	gitBranch := getGitBranch(wd)

	m := model{
		director:            director,
		theme:               DefaultTheme,
		input:               "",
		entries:             []ChatEntry{},
		loading:             false,
		spinner:             s,
		viewport:            viewport.New(80, 20),
		browserPaneHeight:   8,
		browserPanePosition: "top",
		browserPanePreview:  true,
		workingDir:          wd,
		gitBranch:           gitBranch,
		provider:            "hypercode",
		modelName:           "auto",
		toolCount:           toolCount + cliToolCount,
		registry:            reg,
		autocompleteItems:   BuiltinSlashCommands,
		autocompleteMaxVis:  8,
		inputLines:          []string{""},
		cursorLine:          0,
		cursorCol:           0,
		hidingThink:         true,
		thinkingLabel:       "Thinking...",
		contextWindow:       200000,
	}

	// Welcome banner
	m.entries = append(m.entries, ChatEntry{
		Type:      EntrySystem,
		Content:   m.renderWelcome(),
		Timestamp: time.Now(),
	})

	return m
}

func (m model) renderWelcome() string {
	t := m.theme
	var b strings.Builder

	b.WriteString(t.Bold(t.AccentText("╭─────────────────────────────────────────────────────╮")))
	b.WriteString("\n")
	b.WriteString(t.Bold(t.AccentText("│         🧠  HyperHarness — AI Coding Agent          │")))
	b.WriteString("\n")
	b.WriteString(t.Bold(t.AccentText("╰─────────────────────────────────────────────────────╯")))
	b.WriteString("\n\n")

	regCount := 0
	if m.registry != nil {
		regCount = len(m.registry.Tools)
	}

	b.WriteString(t.Dim("  cwd   ") + t.Fg(t.TextColor, shortenPath(m.workingDir)))
	if m.gitBranch != "" {
		b.WriteString(t.Dim(" (") + t.Fg(t.Muted, m.gitBranch) + t.Dim(")"))
	}
	b.WriteString("\n")
	b.WriteString(t.Dim("  tools ") + t.Fg(t.TextColor, fmt.Sprintf("%d registered + %d CLI detected", regCount, m.toolCount-regCount)))
	b.WriteString("\n")
	b.WriteString(t.Dim("  model ") + t.Fg(t.TextColor, m.provider+"/"+m.modelName))
	b.WriteString("\n\n")

	b.WriteString(t.Dim("  ─── Key Bindings ────────────────────────────────────────"))
	b.WriteString("\n")
	b.WriteString(t.Dim("  Enter     ") + t.Fg(t.Muted, "Send message"))
	b.WriteString("\n")
	b.WriteString(t.Dim("  Shift+Tab ") + t.Fg(t.Muted, "New line (multi-line input)"))
	b.WriteString("\n")
	b.WriteString(t.Dim("  ↑/↓       ") + t.Fg(t.Muted, "Navigate input history"))
	b.WriteString("\n")
	b.WriteString(t.Dim("  Tab       ") + t.Fg(t.Muted, "Autocomplete slash commands / @file refs"))
	b.WriteString("\n")
	b.WriteString(t.Dim("  Ctrl+C    ") + t.Fg(t.Muted, "Cancel / quit"))
	b.WriteString("\n")
	b.WriteString(t.Dim("  Ctrl+L    ") + t.Fg(t.Muted, "Toggle file tree pane"))
	b.WriteString("\n")
	b.WriteString(t.Dim("  Ctrl+D    ") + t.Fg(t.Muted, "Toggle dashboard"))
	b.WriteString("\n")
	b.WriteString(t.Dim("  Ctrl+Y    ") + t.Fg(t.Muted, "Yank shell proposal"))
	b.WriteString("\n")
	b.WriteString(t.Dim("  ??query   ") + t.Fg(t.Muted, "Shell command proposal (Copilot CLI parity)"))
	b.WriteString("\n")
	b.WriteString(t.Dim("  /command  ") + t.Fg(t.Muted, "Slash commands (Tab to autocomplete)"))
	b.WriteString("\n\n")

	b.WriteString(t.Dim("  ─── Quick Commands ──────────────────────────────────────"))
	b.WriteString("\n")
	quick := []struct{ cmd, desc string }{
		{"/help", "All commands"},
		{"/tree-browser", "File explorer (modal)"},
		{"/tree-pane", "Persistent file pane"},
		{"/repomap", "Repository map"},
		{"/providers", "LLM provider status"},
		{"/mcp", "MCP tool listing"},
		{"/tools", "All registered tools"},
		{"/plan <prompt>", "Orchestration plan"},
		{"/model", "Select model"},
		{"/settings", "Settings menu"},
		{"/compact", "Compact context"},
		{"/hotkeys", "All keybindings"},
	}
	for _, q := range quick {
		b.WriteString("  " + t.AccentText(q.cmd) + t.Dim("  "+q.desc))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(t.Dim("  ─────────────────────────────────────────────────────────"))
	b.WriteString("\n")

	return b.String()
}

// ═══════════════════════════════════════════════════════════════════════
// Init
// ═══════════════════════════════════════════════════════════════════════

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchGitBranch(m.workingDir))
}

func fetchGitBranch(wd string) tea.Cmd {
	return func() tea.Msg {
		return GitBranchMsg{Branch: getGitBranch(wd)}
	}
}

// ═══════════════════════════════════════════════════════════════════════
// Update
// ═══════════════════════════════════════════════════════════════════════

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	// ─── Window resize ────────────────────────────────────
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			// Reserve: 2 lines footer + 1 line input + 1 blank = 4
			m.viewport = viewport.New(msg.Width, max(msg.Height-5, 5))
			m.viewport.SetContent(m.renderEntries())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = max(msg.Height-5, 5)
		}
		m.viewport.GotoBottom()
		return m, nil

	case GitBranchMsg:
		m.gitBranch = msg.Branch
		return m, nil

	// ─── Keyboard ─────────────────────────────────────────
	case tea.KeyMsg:
		// Tree browser modal
		if m.browserActive {
			return m.updateTreeBrowser(msg)
		}
		// Pinned pane focus
		if m.browserPinned && m.browserPinnedFocus {
			return m.updatePinnedPane(msg)
		}
		// Autocomplete navigation
		if m.showAutocomplete {
			switch msg.Type {
			case tea.KeyUp:
				if m.autocompleteIndex > 0 {
					m.autocompleteIndex--
				}
				return m, nil
			case tea.KeyDown:
				if m.autocompleteIndex < len(m.filteredAutocomplete())-1 {
					m.autocompleteIndex++
				}
				return m, nil
			case tea.KeyTab, tea.KeyEnter:
				items := m.filteredAutocomplete()
				if m.autocompleteIndex >= 0 && m.autocompleteIndex < len(items) {
					// Complete the slash command
					parts := strings.SplitN(m.input, " ", 2)
					m.input = "/" + items[m.autocompleteIndex].Name
					if len(parts) > 1 {
						m.input += " " + parts[1]
					}
					m.showAutocomplete = false
					if msg.Type == tea.KeyEnter {
						return m.handleEnter()
					}
				}
				m.showAutocomplete = false
				return m, nil
			case tea.KeyEsc:
				m.showAutocomplete = false
				return m, nil
			}
			// Fall through for continued typing
		}

		switch msg.Type {
		case tea.KeyCtrlC:
			if m.loading {
				m.loading = false
				m.streaming = false
				m.entries = append(m.entries, ChatEntry{
					Type:      EntrySystem,
					Content:   m.theme.WarningText("⏹ Operation cancelled"),
					Timestamp: time.Now(),
				})
				m.syncViewport()
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyEsc:
			if m.showAutocomplete {
				m.showAutocomplete = false
				return m, nil
			}
			if m.dashboardActive {
				m.dashboardActive = false
				m.entries = append(m.entries, ChatEntry{Type: EntrySystem, Content: m.theme.Dim("[Dashboard] closed"), Timestamp: time.Now()})
				m.syncViewport()
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyCtrlL:
			return m.handleSlashCmd("/tree-pane")

		case tea.KeyCtrlD:
			return m.handleSlashCmd("/dashboard")

		case tea.KeyCtrlY:
			// Yank last shell proposal
			for i := len(m.entries) - 1; i >= 0; i-- {
				if m.entries[i].Type == EntryShellProposal {
					return m.handleSlashCmd("/commit")
				}
			}
			return m, nil

		case tea.KeyTab:
			if m.isSlashContext() {
				m.showAutocomplete = !m.showAutocomplete
				m.autocompleteIndex = 0
				return m, nil
			}
			// Default: toggle pinned pane focus
			if m.browserPinned {
				m.browserPinnedFocus = !m.browserPinnedFocus
				return m, nil
			}
			return m, nil

		case tea.KeyEnter:
			return m.handleEnter()

		case tea.KeyUp:
			if m.input == "" && len(m.inputHistory) > 0 {
				if m.historyIdx < len(m.inputHistory)-1 {
					m.historyIdx++
					m.input = m.inputHistory[len(m.inputHistory)-1-m.historyIdx]
				}
			} else if m.ready {
				m.viewport.LineUp(1)
			}

		case tea.KeyDown:
			if m.input == "" && m.historyIdx > 0 {
				m.historyIdx--
				m.input = m.inputHistory[len(m.inputHistory)-1-m.historyIdx]
			} else if m.ready {
				m.viewport.LineDown(1)
			}

		case tea.KeyPgUp:
			if m.ready { m.viewport.HalfViewUp() }

		case tea.KeyPgDown:
			if m.ready { m.viewport.HalfViewDown() }

		case tea.KeyHome:
			if m.ready { m.viewport.GotoTop() }

		case tea.KeyEnd:
			if m.ready { m.viewport.GotoBottom() }

		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.updateAutocomplete()
			}

		case tea.KeyDelete:
			if len(m.input) > 0 && m.cursorCol < len(m.input) {
				m.input = m.input[:m.cursorCol] + m.input[m.cursorCol+1:]
			}

		case tea.KeyRunes, tea.KeySpace:
			m.input += msg.String()
			m.updateAutocomplete()
			m.historyIdx = 0
		}

	// ─── Async responses ─────────────────────────────────
	case AgentResponseMsg:
		m.loading = false
		m.streaming = false

		entry := ChatEntry{
			Type:      EntryAssistant,
			Content:   msg.Content,
			Provider:  msg.Provider,
			Model:     msg.Model,
			Timestamp: time.Now(),
		}
		if msg.Plan != nil {
			entry.Content = fmt.Sprintf("[Director Plan]\n  task: %s\n  route: %s/%s\n\n%s",
				msg.Plan.TaskType, msg.Plan.Execution.Route.Provider, msg.Plan.Execution.Route.Model, msg.Content)
		}
		m.entries = append(m.entries, entry)
		m.syncViewport()
		// Refresh git branch
		cmds = append(cmds, fetchGitBranch(m.workingDir))

	case StreamChunkMsg:
		if !m.streaming {
			m.streaming = true
			m.streamContent = ""
		}
		m.streamContent += msg.Content
		if msg.Done {
			m.streaming = false
			m.loading = false
			m.entries = append(m.entries, ChatEntry{
				Type:      EntryAssistant,
				Content:   m.streamContent,
				Provider:  m.provider,
				Model:     m.modelName,
				Timestamp: time.Now(),
			})
			m.streamContent = ""
		}
		m.syncViewport()

	case string:
		m.loading = false
		m.entries = append(m.entries, ChatEntry{
			Type:      EntryAssistant,
			Content:   msg,
			Timestamp: time.Now(),
		})
		m.syncViewport()

	case PromptDisplayMsg:
		m.loading = false
		m.entries = append(m.entries, ChatEntry{
			Type:      EntryAssistant,
			Content:   msg.Display,
			Provider:  m.provider,
			Model:     m.modelName,
			Timestamp: time.Now(),
		})
		m.syncViewport()

	case ShellProposalMsg:
		m.loading = false
		m.entries = append(m.entries, ChatEntry{
			Type:      EntryShellProposal,
			Content:   msg.Command,
			ToolName:  "shell",
			ToolArgs:  msg.Explanation,
			Timestamp: time.Now(),
		})
		m.syncViewport()

	case ToolExecMsg:
		m.toolMu.Lock()
		m.toolRuns = append(m.toolRuns, msg)
		if len(m.toolRuns) > 50 { m.toolRuns = m.toolRuns[len(m.toolRuns)-50:] }
		m.toolMu.Unlock()
		m.entries = append(m.entries, ChatEntry{
			Type:      EntryTool,
			Content:   msg.Output,
			ToolName:  msg.ToolName,
			ToolArgs:  msg.Args,
			ToolDur:   msg.Duration,
			ToolErr:   msg.IsError,
			Timestamp: time.Now(),
		})
		m.syncViewport()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// ═══════════════════════════════════════════════════════════════════════
// Input handling helpers
// ═══════════════════════════════════════════════════════════════════════

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	m.showAutocomplete = false
	if strings.TrimSpace(m.input) == "" {
		return m, nil
	}

	req := strings.TrimSpace(m.input)

	// Store in history
	m.inputHistory = append(m.inputHistory, req)
	if len(m.inputHistory) > 100 {
		m.inputHistory = m.inputHistory[len(m.inputHistory)-100:]
	}
	m.input = ""

	// Slash command
	if strings.HasPrefix(req, "/") {
		return m.handleSlashCmd(req)
	}

	// Shell proposal
	if strings.HasPrefix(req, "??") {
		query := strings.TrimSpace(strings.TrimPrefix(req, "??"))
		m.entries = append(m.entries, ChatEntry{
			Type:      EntryUser,
			Content:   "?? " + query,
			Timestamp: time.Now(),
		})
		m.loading = true
		m.syncViewport()
		return m, func() tea.Msg {
			response, err := buildShellProposal(m.director, query)
			if err != nil {
				return fmt.Sprintf("Error: %v", err)
			}
			return response
		}
	}

	// Regular user message
	m.entries = append(m.entries, ChatEntry{
		Type:      EntryUser,
		Content:   req,
		Timestamp: time.Now(),
	})
	m.loading = true
	m.syncViewport()

	// Foundation session
	if sessionID, err := ensureFoundationSession(&m); err == nil {
		m.foundationSessionID = sessionID
		_ = appendFoundationUserText(m.director.WorkingDir, m.foundationSessionID, req)
		refreshPinnedFoundationTreeBrowser(&m)
	}

	return m, func() tea.Msg {
		response, err := buildPromptResponse(m.director, req)
		if err != nil {
			return AgentResponseMsg{Content: fmt.Sprintf("Error: %v", err), Provider: m.provider, Model: m.modelName}
		}
		plan, _ := m.director.State["lastPlan"].(foundationorchestration.PlanResult)
		provider := m.provider
		model := m.modelName
		if plan.Execution.Route.Provider != "" {
			provider = plan.Execution.Route.Provider
		}
		if plan.Execution.Route.Model != "" {
			model = plan.Execution.Route.Model
		}
		return AgentResponseMsg{
			Content:  response.Display,
			Provider: provider,
			Model:    model,
			Plan:     &plan,
		}
	}
}

func (m model) handleSlashCmd(req string) (tea.Model, tea.Cmd) {
	mdl, cmd := ProcessSlashCommand(req, &m)
	m = mdl.(model)
	m.syncViewport()
	return m, cmd
}

func (m *model) updateAutocomplete() {
	if m.isSlashContext() {
		m.showAutocomplete = true
		m.autocompleteIndex = 0
	} else {
		m.showAutocomplete = false
	}
}

func (m model) isSlashContext() bool {
	return strings.HasPrefix(m.input, "/") && !strings.Contains(m.input, " ")
}

func (m model) filteredAutocomplete() []SlashCommand {
	prefix := strings.ToLower(m.input)
	if !strings.HasPrefix(prefix, "/") {
		return BuiltinSlashCommands
	}
	prefix = prefix[1:] // strip /

	var filtered []SlashCommand
	for _, cmd := range BuiltinSlashCommands {
		if strings.HasPrefix(cmd.Name, prefix) {
			filtered = append(filtered, cmd)
		}
	}
	return filtered
}

// ═══════════════════════════════════════════════════════════════════════
// Tree browser / pinned pane keyboard handling
// ═══════════════════════════════════════════════════════════════════════

func (m model) updateTreeBrowser(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		m.entries = append(m.entries, ChatEntry{Type: EntrySystem, Content: m.theme.Dim("[Tree Browser] closed"), Timestamp: time.Now()})
		m.syncViewport()
		return m, nil
	case tea.KeyLeft:
		if !m.browserConfirmPending && m.browserIndex >= 0 && m.browserIndex < len(visible) {
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
		if !m.browserConfirmPending && m.browserIndex >= 0 && m.browserIndex < len(visible) {
			item := visible[m.browserIndex]
			if m.browserCollapsed != nil { delete(m.browserCollapsed, item.ID) }
		}
		return m, nil
	case tea.KeyUp:
		if !m.browserConfirmPending { m.browserIndex = max(0, m.browserIndex-1) }
		return m, nil
	case tea.KeyDown:
		if !m.browserConfirmPending { m.browserIndex = min(len(visible)-1, m.browserIndex+1) }
		return m, nil
	case tea.KeyHome:
		m.browserIndex = 0; return m, nil
	case tea.KeyEnd:
		m.browserIndex = max(0, len(visible)-1); return m, nil
	case tea.KeyPgUp:
		m.browserIndex = max(0, m.browserIndex-10); return m, nil
	case tea.KeyPgDown:
		m.browserIndex = min(len(visible)-1, m.browserIndex+10); return m, nil
	case tea.KeyEnter:
		if m.browserIndex >= 0 && m.browserIndex < len(visible) {
			if !m.browserConfirmPending {
				m.browserConfirmPending = true
				return m, nil
			}
			display, err := openSelectedTreeBrowser(m.director.WorkingDir, m.foundationSessionID, visible, m.browserIndex, 128)
			if err != nil {
				m.entries = append(m.entries, ChatEntry{Type: EntrySystem, Content: m.theme.ErrorText(fmt.Sprintf("[Error] %v", err)), Timestamp: time.Now()})
			} else {
				m.entries = append(m.entries, ChatEntry{Type: EntrySystem, Content: m.theme.SuccessText(display), Timestamp: time.Now()})
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

func (m model) updatePinnedPane(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visible := visibleTreeBrowserItems(m.browserItems, m.browserFilter, m.browserCollapsed)
	switch msg.Type {
	case tea.KeyEsc, tea.KeyTab:
		m.browserPinnedFocus = false
		return m, nil
	case tea.KeyUp:
		m.browserIndex = max(0, m.browserIndex-1); return m, nil
	case tea.KeyDown:
		m.browserIndex = min(len(visible)-1, m.browserIndex+1); return m, nil
	case tea.KeyHome:
		m.browserIndex = 0; return m, nil
	case tea.KeyEnd:
		m.browserIndex = max(0, len(visible)-1); return m, nil
	case tea.KeyEnter:
		if m.browserIndex >= 0 && m.browserIndex < len(visible) {
			if !m.browserConfirmPending { m.browserConfirmPending = true; return m, nil }
			display, err := openSelectedTreeBrowser(m.director.WorkingDir, m.foundationSessionID, visible, m.browserIndex, 128)
			if err != nil {
				m.entries = append(m.entries, ChatEntry{Type: EntrySystem, Content: m.theme.ErrorText(fmt.Sprintf("[Error] %v", err)), Timestamp: time.Now()})
			} else {
				m.entries = append(m.entries, ChatEntry{Type: EntrySystem, Content: m.theme.SuccessText(display), Timestamp: time.Now()})
			}
			m.browserConfirmPending = false
			refreshPinnedFoundationTreeBrowser(&m)
			m.syncViewport()
		}
		return m, nil
	}
	return m, nil
}

// ═══════════════════════════════════════════════════════════════════════
// View — pi-mono layout: [tree-pane?] [chat viewport] [footer] [input]
// ═══════════════════════════════════════════════════════════════════════

func (m model) View() string {
	if m.quitting {
		return m.theme.Dim("Goodbye.\n")
	}

	// Dashboard mode — full split-pane
	if m.dashboardActive {
		chatContent := m.renderEntries()
		toolContent := m.renderToolSidebar()
		metrics := m.renderMetrics()
		return RenderDashboard(chatContent, toolContent, metrics)
	}

	var body string

	// Tree browser modal overlay
	if m.browserActive {
		chatView := m.renderChatArea()
		browserView := renderTreeBrowser(m.browserItems, m.browserIndex, m.browserFilter,
			m.browserConfirmPending, m.browserCollapsed, m.browserGrouped,
			0, m.theme.AccentText("[File Tree Browser]"), true)
		body = chatView + "\n" + browserView
		return body + "\n" + m.renderFooter() + "\n" + m.renderInputBar()
	}

	// Pinned tree pane
	if m.browserPinned {
		paneHeight := m.browserPaneHeight
		if paneHeight <= 0 { paneHeight = 8 }
		title := m.theme.Dim("[Tree Pane]")
		if m.browserPinnedFocus { title = m.theme.AccentText("[Tree Pane :: Focused]") }
		pane := renderTreeBrowser(m.browserItems, m.browserIndex, m.browserFilter,
			m.browserConfirmPending && m.browserPinnedFocus, m.browserCollapsed, m.browserGrouped,
			paneHeight, title, m.browserPanePreview)
		divider := m.theme.Dim("─" + strings.Repeat("─", max(0, m.width-1)))

		if strings.ToLower(strings.TrimSpace(m.browserPanePosition)) == "bottom" {
			body = m.renderChatArea() + "\n" + divider + "\n" + pane
		} else {
			body = pane + "\n" + divider + "\n" + m.renderChatArea()
		}
		return body + "\n" + m.renderFooter() + "\n" + m.renderInputBar()
	}

	// Normal mode
	body = m.renderChatArea()
	return body + "\n" + m.renderFooter() + "\n" + m.renderInputBar()
}

// ─── Chat area ────────────────────────────────────────────────────────

func (m model) renderChatArea() string {
	if m.ready {
		return m.viewport.View()
	}
	return m.renderEntries()
}

// ─── Entry rendering — mirrors pi-mono's component rendering ──────────

func (m model) renderEntries() string {
	var lines []string

	for _, entry := range m.entries {
		switch entry.Type {
		case EntryUser:
			lines = append(lines, m.renderUserEntry(entry))
		case EntryAssistant:
			lines = append(lines, m.renderAssistantEntry(entry))
		case EntryTool:
			lines = append(lines, m.renderToolEntry(entry))
		case EntrySystem:
			lines = append(lines, entry.Content)
		case EntryThinking:
			lines = append(lines, m.theme.Italic(m.theme.Fg(m.theme.ThinkingCol, entry.Content)))
		case EntryShellProposal:
			lines = append(lines, m.renderShellProposal(entry))
		case EntryDiff:
			lines = append(lines, m.renderDiffEntry(entry))
		}
		lines = append(lines, "") // blank line between entries
	}

	// Streaming indicator
	if m.loading && m.streaming && m.streamContent != "" {
		lines = append(lines, m.theme.Fg(m.theme.Accent, "┃ ")+m.streamContent)
	} else if m.loading {
		lines = append(lines, m.theme.Fg(m.theme.ThinkingCol, m.theme.Italic(m.thinkingLabel))+" "+m.spinner.View())
	}

	return strings.Join(lines, "\n")
}

func (m model) renderUserEntry(e ChatEntry) string {
	t := m.theme
	promptText := t.Bold(t.Fg(t.Accent, "> ")) + t.Fg(t.UserMsgText, e.Content)
	// Wrap in a user message box (pi-mono style with subtle bg)
	boxStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.UserMsgText)).
		Background(lipgloss.Color(t.UserMsgBg)).
		Padding(0, 1)
	return boxStyle.Render(promptText)
}

func (m model) renderAssistantEntry(e ChatEntry) string {
	t := m.theme
	var header string
	if e.Provider != "" || e.Model != "" {
		modelStr := e.Provider
		if e.Model != "" {
			modelStr += "/" + e.Model
		}
		header = t.Dim("  "+modelStr) + "\n"
	}
	content := t.Fg(t.TextColor, e.Content)
	// Simple markdown-like rendering
	content = m.renderSimpleMarkdown(content)
	return header + content
}

func (m model) renderToolEntry(e ChatEntry) string {
	t := m.theme
	var statusIcon string
	if e.ToolErr {
		statusIcon = t.ErrorText("✗")
	} else {
		statusIcon = t.SuccessText("✓")
	}

	var borderStyle lipgloss.Style
	if e.ToolErr {
		borderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(t.Error))
	} else {
		borderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(t.Warning))
	}

	header := fmt.Sprintf(" %s %s ", statusIcon, t.Bold(t.Fg(t.ToolTitle, e.ToolName)))
	if e.ToolDur > 0 {
		header += t.Dim(fmt.Sprintf(" (%v)", e.ToolDur))
	}
	if e.ToolArgs != "" {
		argDisplay := e.ToolArgs
		if len(argDisplay) > 80 { argDisplay = argDisplay[:80] + "…" }
		header += "\n  " + t.Dim(argDisplay)
	}

	var body string
	if e.Expanded || e.ToolErr {
		output := e.Content
		if len(output) > 800 && !e.Expanded {
			output = output[:800] + "\n" + t.Dim("  ... (truncated)")
		}
		if output != "" {
			body = "\n" + t.Fg(t.ToolOutput, output)
		}
	} else {
		// Collapsed: show first line only
		lines := strings.SplitN(e.Content, "\n", 2)
		preview := lines[0]
		if len(preview) > 120 { preview = preview[:120] + "…" }
		if preview != "" {
			body = "\n" + t.Dim(preview)
		}
	}

	return borderStyle.Render(header + body)
}

func (m model) renderShellProposal(e ChatEntry) string {
	t := m.theme
	var b strings.Builder
	b.WriteString(t.Fg(t.ToolTitle, "⌘ Shell Proposal"))
	b.WriteString("\n")
	b.WriteString(t.Fg(t.TextColor, "  $ "+e.Content))
	if e.ToolArgs != "" {
		b.WriteString("\n" + t.Dim("  "+e.ToolArgs))
	}
	b.WriteString("\n" + t.WarningText("  Execute? [Y/n] Ctrl+Y to accept"))
	return b.String()
}

func (m model) renderDiffEntry(e ChatEntry) string {
	t := m.theme
	lines := strings.Split(e.Content, "\n")
	var rendered []string
	for _, line := range lines {
		if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "→") {
			rendered = append(rendered, t.SuccessText(line))
		} else if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "✗") {
			rendered = append(rendered, t.ErrorText(line))
		} else if strings.HasPrefix(line, "@@") {
			rendered = append(rendered, t.Fg(t.Muted, line))
		} else {
			rendered = append(rendered, t.Fg(t.DimColor, line))
		}
	}
	return strings.Join(rendered, "\n")
}

func (m model) renderSimpleMarkdown(text string) string {
	// Minimal markdown rendering — code blocks, inline code, bold, italic, headings, lists
	lines := strings.Split(text, "\n")
	var result []string
	inCodeBlock := false

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			if inCodeBlock {
				result = append(result, m.theme.Fg(m.theme.DimColor, line))
			} else {
				result = append(result, m.theme.Fg(m.theme.DimColor, "```"))
			}
			continue
		}
		if inCodeBlock {
			result = append(result, m.theme.Fg(m.theme.MDCodeBlock, line))
			continue
		}
		// Headings
		if strings.HasPrefix(line, "### ") {
			result = append(result, m.theme.Bold(m.theme.Fg(m.theme.MDHeading, line)))
			continue
		}
		if strings.HasPrefix(line, "## ") {
			result = append(result, m.theme.Bold(m.theme.Fg(m.theme.MDHeading, line)))
			continue
		}
		if strings.HasPrefix(line, "# ") {
			result = append(result, m.theme.Bold(m.theme.Fg(m.theme.MDHeading, line)))
			continue
		}
		// List items
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			result = append(result, m.theme.Fg(m.theme.DimColor, line[:2])+m.theme.Fg(m.theme.TextColor, line[2:]))
			continue
		}
		// Quote
		if strings.HasPrefix(line, "> ") {
			result = append(result, m.theme.Fg(m.theme.MDQuote, line))
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

// ─── Footer — mirrors pi-mono's FooterComponent ──────────────────────

func (m model) renderFooter() string {
	t := m.theme

	// Line 1: pwd + git branch + session name
	pwd := shortenPath(m.workingDir)
	if m.gitBranch != "" {
		pwd += " (" + m.gitBranch + ")"
	}

	// Line 2: token stats + context % + model
	var stats []string
	if m.totalInputTok > 0 {
		stats = append(stats, fmt.Sprintf("↑%s", formatTokens(m.totalInputTok)))
	}
	if m.totalOutTok > 0 {
		stats = append(stats, fmt.Sprintf("↓%s", formatTokens(m.totalOutTok)))
	}
	if m.totalCost > 0 {
		stats = append(stats, fmt.Sprintf("$%.3f", m.totalCost))
	}
	contextStr := fmt.Sprintf("%.0f%%/%s", m.contextPct, formatTokens(m.contextWindow))
	if m.contextPct > 90 {
		contextStr = t.ErrorText(contextStr)
	} else if m.contextPct > 70 {
		contextStr = t.WarningText(contextStr)
	}
	stats = append(stats, contextStr)

	statsLeft := strings.Join(stats, " ")

	modelStr := m.provider + "/" + m.modelName
	if len(m.provider) > 0 && m.provider != "hypercode" {
		modelStr = "(" + m.provider + ") " + m.modelName
	}

	// Compose footer
	pwdLine := t.Dim(pwd)
	if len(pwdLine) > m.width {
		pwdLine = pwdLine[:m.width]
	}

	rightWidth := len(modelStr)
	leftWidth := len(statsLeft)
	gap := max(0, m.width-leftWidth-rightWidth-2)
	statsLine := t.Dim(statsLeft) + strings.Repeat(" ", gap) + t.Dim(modelStr)

	// Extension/MCP status line
	var extStatus string
	adapter := adapters.NewMCPAdapter(m.workingDir)
	status := adapter.Status()
		if len(status.Servers) > 0 {
			extStatus = t.Dim(fmt.Sprintf("  mcp: %d servers", len(status.Servers)))
		}

	result := pwdLine + "\n" + statsLine
	if extStatus != "" {
		result += "\n" + extStatus
	}
	return result
}

// ─── Input bar — pi-mono style with prompt ────────────────────────────

func (m model) renderInputBar() string {
	t := m.theme
	prompt := t.Bold(t.AccentText("> "))

	var inputDisplay string
	if m.loading {
		inputDisplay = t.Fg(t.ThinkingCol, t.Italic(m.thinkingLabel)) + " " + m.spinner.View()
	} else {
		inputDisplay = t.Fg(t.TextColor, m.input) + "▎"
	}

	// Autocomplete dropdown
	autocomplete := ""
	if m.showAutocomplete && m.isSlashContext() {
		autocomplete = "\n" + m.renderAutocomplete()
	}

	return prompt + inputDisplay + autocomplete
}

func (m model) renderAutocomplete() string {
	t := m.theme
	items := m.filteredAutocomplete()
	if len(items) == 0 {
		return t.Dim("  No matching commands")
	}

	var lines []string
	start := 0
	if m.autocompleteIndex >= m.autocompleteMaxVis {
		start = m.autocompleteIndex - m.autocompleteMaxVis + 1
	}
	end := min(len(items), start+m.autocompleteMaxVis)

	for i := start; i < end; i++ {
		item := items[i]
		if i == m.autocompleteIndex {
			lines = append(lines, t.Bold(t.AccentText("  ▶ /"+item.Name))+t.Dim("  "+item.Description))
		} else {
			lines = append(lines, t.Dim("    /"+item.Name)+"  "+t.Dim(item.Description))
		}
	}

	if len(items) > m.autocompleteMaxVis {
		lines = append(lines, t.Dim(fmt.Sprintf("  ↑↓ scroll  (%d/%d)", m.autocompleteIndex+1, len(items))))
	}

	return strings.Join(lines, "\n")
}

// ─── Tool sidebar (dashboard mode) ────────────────────────────────────

func (m model) renderToolSidebar() string {
	t := m.theme
	var lines []string
	lines = append(lines, t.Bold(t.AccentText("Active Tools")))
	lines = append(lines, "")

	// Show registered tools grouped
	if m.registry != nil {
		groups := map[string]int{}
		for _, tool := range m.registry.Tools {
			groups[toolGroupName(tool.Name)]++
		}
		for name, count := range groups {
			lines = append(lines, t.Fg(t.ToolTitle, name)+t.Dim(fmt.Sprintf(" (%d)", count)))
		}
	}

	// MCP tools
	adapter := adapters.NewMCPAdapter(m.workingDir)
	status := adapter.Status()
		if len(status.Servers) > 0 {
			lines = append(lines, t.Fg(t.ToolTitle, "mcp")+t.Dim(fmt.Sprintf(" (%d)", len(status.Servers))))
		}

	// Recent tool executions
	m.toolMu.Lock()
	recent := m.toolRuns
	m.toolMu.Unlock()
	if len(recent) > 5 {
		recent = recent[len(recent)-5:]
	}
	if len(recent) > 0 {
		lines = append(lines, "")
		lines = append(lines, t.Bold(t.Dim("Recent Executions")))
		for _, run := range recent {
			icon := t.SuccessText("✓")
			if run.IsError { icon = t.ErrorText("✗") }
			lines = append(lines, icon+" "+t.Fg(t.ToolTitle, run.ToolName)+t.Dim(fmt.Sprintf(" (%v)", run.Duration)))
		}
	}

	return strings.Join(lines, "\n")
}

func toolGroupName(name string) string {
	parts := strings.SplitN(name, "_", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return name
}

// ─── Metrics line ─────────────────────────────────────────────────────

func (m model) renderMetrics() string {
	return fmt.Sprintf("Tokens: %s in / %s out │ Cost: $%.3f │ Scope: Project │ Tools: %d",
		formatTokens(m.totalInputTok), formatTokens(m.totalOutTok), m.totalCost, m.toolCount)
}

// ─── Helpers ──────────────────────────────────────────────────────────

func (m *model) syncViewport() {
	m.viewport.SetContent(m.renderEntries())
	m.viewport.GotoBottom()
}

func formatTokens(count int) string {
	if count < 1000 { return fmt.Sprintf("%d", count) }
	if count < 10000 { return fmt.Sprintf("%.1fk", float64(count)/1000) }
	if count < 1000000 { return fmt.Sprintf("%dk", count/1000) }
	return fmt.Sprintf("%.1fM", float64(count)/1000000)
}

func max(a, b int) int {
	if a > b { return a }
	return b
}

func min(a, b int) int {
	if a < b { return a }
	return b
}

// ═══════════════════════════════════════════════════════════════════════
// Entry Point
// ═══════════════════════════════════════════════════════════════════════

func StartREPL() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Unused: kept for build compatibility
var _ = runtime.NumCPU
