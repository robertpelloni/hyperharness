package tui

// ═══════════════════════════════════════════════════════════════════════
// chat.go — TUI model, update loop, and view composition
// Mirrors pi-mono's InteractiveMode class with all features from:
//   - pi-mono: double-escape, model cycling, follow-up queuing, suspend
//   - goose: permission dialog, ACP streaming, paste mode
//   - opencode: command palette, theme selector, session list
//   - claude-code: /cost, /diff, /rewind, /doctor, /memory, /vim
// ═══════════════════════════════════════════════════════════════════════

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
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
	"github.com/robertpelloni/hyperharness/tools"
)

// ═══════════════════════════════════════════════════════════════════════
// Messages (async events from agent loop, tools, etc.)
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
type AgentResponseMsg struct {
	Content  string
	Provider string
	Model    string
	Plan     *foundationorchestration.PlanResult
}
type GitBranchMsg struct{ Branch string }
type StreamChunkMsg struct{ Chunk string }

// ═══════════════════════════════════════════════════════════════════════
// Model — the central TUI state (mirrors pi-mono's InteractiveMode fields)
// ═══════════════════════════════════════════════════════════════════════

type model struct {
	// Core
	director *agents.Director
	session  *agents.AgentSession
	theme    Theme
	themeName string

	// Chat history (structured entries, not plain strings)
	entries []ChatEntry
	input   string

	// Viewport (scrollable chat area)
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

	// Autocomplete (pi-mono style slash command dropdown)
	showAutocomplete   bool
	autocompleteItems  []SlashCommand
	autocompleteIndex  int
	autocompleteMaxVis int

	// Model selector (pi-mono/opencode style)
	showModelSelector   bool
	modelSelectorIdx    int
	modelSelectorFilter string
	availableModels     []string

	// Command palette (opencode style)
	showCommandPalette  bool
	commandPaletteIdx   int
	commandPaletteFilter string

	// Footer data
	workingDir    string
	gitBranch     string
	provider      string
	modelName     string
	toolCount     int
	totalInputTok int
	totalOutTok   int
	totalCost     float64
	contextPct    float64
	contextWindow int

	// Input history (Up/Down navigation)
	inputHistory []string
	historyIdx   int

	// Follow-up queue (pi-mono: Ctrl+F queues messages during streaming)
	followUpQueue []string

	// Double-escape (pi-mono: Esc Esc triggers tree/fork selector)
	lastEscapeTime int64
	doubleEscAction string // "tree" or "fork" (mirrors pi-mono's doubleEscapeAction setting)

	// Thinking
	thinkingLevel string
	hidingThink   bool
	thinkingLabel string

	// Tool output expand/collapse (pi-mono: Ctrl+O)
	toolOutputExpanded bool

	// Tool registry
	registry *tools.Registry

	// Tool execution tracking
	toolMu   sync.Mutex
	toolRuns []ToolExecMsg

	// Streaming
	streaming     bool
	streamContent string

	// Bash mode (! prefix)
	bashMode bool

	// Permission dialog (goose-style)
	showPermission    bool
	permissionEntry   *ChatEntry

	// Paste mode (goose-style: large paste gets special handling)
	pasteMode    bool
	pasteContent string

	// Debug mode (pi-mono)
	debugMode bool

	// Vim keybindings (claude-code)
	vimMode bool
}

// ═══════════════════════════════════════════════════════════════════════
// Initialization
// ═══════════════════════════════════════════════════════════════════════

func getWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil { return "." }
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
	if err != nil { return "" }
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
	regCount := len(reg.Tools)
	cliToolCount := countInstalledTools()
	gitBranch := getGitBranch(wd)

	// Create AgentSession with tool adapter
	adapter := NewRegistryAdapter(reg)
	session := agents.NewAgentSession(director, adapter, filepath.Join(wd, ".hyperharness", "sessions"))

	m := model{
		director:         director,
		session:          session,
		theme:            DefaultTheme,
		themeName:        "default",
		input:            "",
		entries:          []ChatEntry{},
		loading:          false,
		spinner:          s,
		viewport:         viewport.New(80, 20),
		browserPaneHeight:   8,
		browserPanePosition: "top",
		browserPanePreview:  true,
		workingDir:       wd,
		gitBranch:        gitBranch,
		provider:         "hypercode",
		modelName:        "auto",
		toolCount:        regCount + cliToolCount,
		registry:         reg,
		autocompleteItems:   BuiltinSlashCommands,
		autocompleteMaxVis:  8,
		hidingThink:         true,
		thinkingLabel:       "Thinking...",
		thinkingLevel:       "off",
		contextWindow:       200000,
		doubleEscAction:     "tree", // pi-mono default
		availableModels:     []string{"auto", "gemini-1.5-pro", "gpt-4", "claude-3-5-sonnet", "claude-3-opus", "llama-3", "local"},
	}

	// Subscribe to session events for real-time tool/streaming display
	session.Subscribe(func(event agents.SessionEvent) {
		// Events are handled in the agent loop; this subscription
		// enables future wiring to bubbletea commands
	})

	// Welcome banner
	m.entries = append(m.entries, ChatEntry{
		Type:      EntrySystem,
		Content:   RenderWelcome(wd, gitBranch, "hypercode", "auto", m.toolCount, regCount, DefaultTheme),
		Timestamp: time.Now(),
	})

	return m
}

// ═══════════════════════════════════════════════════════════════════════
// Init
// ═══════════════════════════════════════════════════════════════════════

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchGitBranch(m.workingDir))
}

func fetchGitBranch(wd string) tea.Cmd {
	return func() tea.Msg { return GitBranchMsg{Branch: getGitBranch(wd)} }
}

// ═══════════════════════════════════════════════════════════════════════
// Update — the main event loop (mirrors pi-mono's handleEvent)
// ═══════════════════════════════════════════════════════════════════════

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	// ─── Window resize ────────────────────────────────────
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.viewport = viewport.New(msg.Width, max(msg.Height-5, 5))
			m.viewport.SetContent(m.renderAllEntries())
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

		// Permission dialog (goose-style: ↑↓ navigate, Enter confirm, Esc cancel)
		if m.showPermission && m.permissionEntry != nil {
			return m.updatePermissionDialog(msg)
		}

		// Model selector (pi-mono/opencode: ↑↓ navigate, Enter select, Esc cancel)
		if m.showModelSelector {
			return m.updateModelSelector(msg)
		}

		// Command palette (opencode: ↑↓ navigate, Enter select, Esc cancel)
		if m.showCommandPalette {
			return m.updateCommandPalette(msg)
		}

		// Tree browser modal
		if m.browserActive {
			return m.updateTreeBrowser(msg)
		}

		// Pinned pane focus
		if m.browserPinned && m.browserPinnedFocus {
			return m.updatePinnedPane(msg)
		}

		// Autocomplete navigation (pi-mono style)
		if m.showAutocomplete {
			switch msg.Type {
			case tea.KeyUp:
				if m.autocompleteIndex > 0 { m.autocompleteIndex-- }
				return m, nil
			case tea.KeyDown:
				if m.autocompleteIndex < len(m.filteredAutocomplete())-1 { m.autocompleteIndex++ }
				return m, nil
			case tea.KeyTab:
				items := m.filteredAutocomplete()
				if m.autocompleteIndex >= 0 && m.autocompleteIndex < len(items) {
					parts := strings.SplitN(m.input, " ", 2)
					m.input = "/" + items[m.autocompleteIndex].Name
					if len(parts) > 1 { m.input += " " + parts[1] }
					m.showAutocomplete = false
				}
				m.showAutocomplete = false
				return m, nil
			case tea.KeyEsc:
				m.showAutocomplete = false
				return m, nil
			}
		}

		switch msg.Type {

		// ── Ctrl+C: interrupt or quit (pi-mono: app.clear/app.interrupt) ──
		case tea.KeyCtrlC:
			if m.loading {
				m.loading = false
				m.streaming = false
				if m.session != nil { m.session.Abort() }
				m.entries = append(m.entries, ChatEntry{
					Type:      EntrySystem,
					Content:   m.theme.WarningText("⏹ Operation cancelled"),
					Timestamp: time.Now(),
				})
				m.syncViewport()
				return m, nil
			}
			// First Ctrl+C clears input (pi-mono: app.clear)
			if m.input != "" {
				m.input = ""
				return m, nil
			}
			return m, tea.Quit

		// ── Esc: multi-purpose (pi-mono's onEscape chain) ──
		case tea.KeyEsc:
			if m.showAutocomplete { m.showAutocomplete = false; return m, nil }
			if m.dashboardActive {
				m.dashboardActive = false
				m.entries = append(m.entries, ChatEntry{Type: EntrySystem, Content: m.theme.Dim("[Dashboard] closed"), Timestamp: time.Now()})
				m.syncViewport()
				return m, nil
			}
			if m.bashMode {
				m.input = ""
				m.bashMode = false
				return m, nil
			}
			if m.pasteMode {
				m.pasteMode = false
				m.pasteContent = ""
				return m, nil
			}
			// Double-escape with empty editor (pi-mono: triggers tree/fork)
			if m.input == "" {
				now := time.Now().UnixMilli()
				if now-m.lastEscapeTime < 500 {
					// Double-escape triggered!
					m.lastEscapeTime = 0
					if m.doubleEscAction == "tree" {
						return m.handleSlashCmd("/tree-browser")
					} else if m.doubleEscAction == "fork" {
						return m.handleSlashCmd("/fork")
					}
				} else {
					m.lastEscapeTime = now
				}
				return m, nil
			}
			// Single escape clears input
			m.input = ""
			return m, nil

		// ── Ctrl+D: exit/quit (pi-mono: handleCtrlD) ──
		case tea.KeyCtrlD:
			if m.input == "" {
				return m, tea.Quit
			}
			return m, nil

		// ── Ctrl+Z: suspend to shell (pi-mono: handleCtrlZ) ──
		case tea.KeyCtrlZ:
			return m, m.suspendToShell()

		// ── Ctrl+L: tree browser ──
		case tea.KeyCtrlL:
			return m.handleSlashCmd("/tree-pane")

		// ── Ctrl+D when not empty: dashboard toggle ──
		// (handled above for empty input — quit)
		// With input: we use Ctrl+Shift+D for dashboard

		// ── Ctrl+Y: accept shell proposal / commit ──
		case tea.KeyCtrlY:
			return m.handleSlashCmd("/commit")

		// ── Ctrl+O: toggle tool output expansion (pi-mono: app.tools.expand) ──
		case tea.KeyCtrlO:
			m.toolOutputExpanded = !m.toolOutputExpanded
			for i := range m.entries {
				if m.entries[i].Type == EntryTool {
					m.entries[i].Expanded = m.toolOutputExpanded
				}
			}
			m.syncViewport()
			return m, nil

		// ── Ctrl+P: cycle thinking levels (pi-mono: app.thinking.cycle) ──
		case tea.KeyCtrlP:
			if m.session != nil {
				newLevel := m.session.CycleThinkingLevel()
				m.thinkingLevel = string(newLevel)
			} else {
				levels := []string{"off", "minimal", "low", "medium", "high", "xhigh"}
				for i, l := range levels {
					if l == m.thinkingLevel {
						m.thinkingLevel = levels[(i+1)%len(levels)]
						break
					}
				}
			}
			m.entries = append(m.entries, ChatEntry{
				Type:          EntrySystem,
				Content:       m.theme.AccentText(fmt.Sprintf("[Thinking] level set to %s", m.thinkingLevel)),
				Timestamp:     time.Now(),
				ThinkingLevel: m.thinkingLevel,
			})
			m.syncViewport()
			return m, nil

		// ── Ctrl+N: cycle model forward (pi-mono: app.model.cycleForward) ──
		case tea.KeyCtrlN:
			if m.session != nil {
				prov, mdl := m.session.CycleModel("forward")
				m.provider = prov
				m.modelName = mdl
			} else {
				idx := 0
				for i, m2 := range m.availableModels {
					if m2 == m.modelName { idx = i; break }
				}
				idx = (idx + 1) % len(m.availableModels)
				m.modelName = m.availableModels[idx]
			}
			m.entries = append(m.entries, ChatEntry{
				Type:      EntrySystem,
				Content:   m.theme.AccentText(fmt.Sprintf("[Model] switched to %s/%s", m.provider, m.modelName)),
				Timestamp: time.Now(),
			})
			m.syncViewport()
			return m, nil

		// ── Ctrl+R: cycle model backward (pi-mono: app.model.cycleBackward)
		// Note: Ctrl+M = Enter in terminals, so we remap to Ctrl+R
		case tea.KeyCtrlR:
			if m.session != nil {
				prov, mdl := m.session.CycleModel("backward")
				m.provider = prov
				m.modelName = mdl
			} else {
				idx := 0
				for i, m2 := range m.availableModels {
					if m2 == m.modelName { idx = i; break }
				}
				idx = (idx - 1 + len(m.availableModels)) % len(m.availableModels)
				m.modelName = m.availableModels[idx]
			}
			m.entries = append(m.entries, ChatEntry{
				Type:      EntrySystem,
				Content:   m.theme.AccentText(fmt.Sprintf("[Model] switched to %s/%s", m.provider, m.modelName)),
				Timestamp: time.Now(),
			})
			m.syncViewport()
			return m, nil

		// ── Ctrl+E: external editor (pi-mono: app.editor.external) ──
		case tea.KeyCtrlE:
			if m.input != "" {
				editor := os.Getenv("EDITOR")
				if editor == "" { editor = "vim" }
				tmpFile := filepath.Join(os.TempDir(), "hyperharness-input.md")
				os.WriteFile(tmpFile, []byte(m.input), 0644)
				cmd := exec.Command(editor, tmpFile)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
				if content, err := os.ReadFile(tmpFile); err == nil {
					m.input = strings.TrimSpace(string(content))
				}
				os.Remove(tmpFile)
			}
			return m, nil

		// ── Ctrl+F: follow-up / queue message (pi-mono: app.message.followUp) ──
		case tea.KeyCtrlF:
			if m.loading && m.input != "" {
				m.followUpQueue = append(m.followUpQueue, m.input)
				m.entries = append(m.entries, ChatEntry{
					Type:          EntryQueue,
					Content:       m.input,
					QueuePosition: len(m.followUpQueue),
					Timestamp:     time.Now(),
				})
				m.input = ""
				m.syncViewport()
				return m, nil
			}
			return m, nil

		// ── Ctrl+G: dequeue/flush queued messages (pi-mono: app.message.dequeue) ──
		case tea.KeyCtrlG:
			if len(m.followUpQueue) > 0 {
				m.entries = append(m.entries, ChatEntry{
					Type:      EntrySystem,
					Content:   m.theme.Dim(fmt.Sprintf("[Queue] %d messages dequeued", len(m.followUpQueue))),
					Timestamp: time.Now(),
				})
				m.followUpQueue = nil
				m.syncViewport()
			}
			return m, nil

		// ── Ctrl+V: clipboard image paste (pi-mono: app.clipboard.pasteImage) ──
		case tea.KeyCtrlV:
			// Read clipboard content — if it's a file path, insert it
			if m.input == "" {
				clipboard, err := exec.Command("powershell", "-command", "Get-Clipboard").Output()
				if err == nil && len(clipboard) > 0 {
					text := strings.TrimSpace(string(clipboard))
					// Check if it looks like an image path
					if strings.HasSuffix(strings.ToLower(text), ".png") ||
						strings.HasSuffix(strings.ToLower(text), ".jpg") ||
						strings.HasSuffix(strings.ToLower(text), ".gif") ||
						strings.HasSuffix(strings.ToLower(text), ".webp") {
						m.entries = append(m.entries, ChatEntry{
							Type:      EntryImage,
							ImagePath: text,
							ImageMime: "image/" + filepath.Ext(text)[1:],
							Timestamp: time.Now(),
						})
						m.syncViewport()
						return m, nil
					}
					m.input = text
				}
			}
			return m, nil

		// ── Tab: autocomplete / pane focus (pi-mono/opencode) ──
		case tea.KeyTab:
			if m.isSlashContext() {
				m.showAutocomplete = !m.showAutocomplete
				m.autocompleteIndex = 0
				return m, nil
			}
			// Tab toggles tool expansion in goose-style (when viewing a turn with tools)
			if !m.loading && len(m.entries) > 0 {
				lastEntry := m.entries[len(m.entries)-1]
				if lastEntry.Type == EntryTool {
					m.toolOutputExpanded = !m.toolOutputExpanded
					for i := range m.entries {
						if m.entries[i].Type == EntryTool {
							m.entries[i].Expanded = m.toolOutputExpanded
						}
					}
					m.syncViewport()
					return m, nil
				}
			}
			if m.browserPinned {
				m.browserPinnedFocus = !m.browserPinnedFocus
				return m, nil
			}
			return m, nil

		// ── Shift+Tab: open command palette (opencode style) ──
		// Note: bubbletea doesn't distinguish Shift+Tab easily,
		// so Ctrl+K opens the command palette instead

		// ── Ctrl+K: command palette (opencode style) ──
		case tea.KeyCtrlK:
			m.showCommandPalette = true
			m.commandPaletteIdx = 0
			m.commandPaletteFilter = ""
			return m, nil

		// ── Enter ──
		case tea.KeyEnter:
			// Autocomplete: complete and send
			if m.showAutocomplete {
				items := m.filteredAutocomplete()
				if m.autocompleteIndex >= 0 && m.autocompleteIndex < len(items) {
					parts := strings.SplitN(m.input, " ", 2)
					m.input = "/" + items[m.autocompleteIndex].Name
					if len(parts) > 1 { m.input += " " + parts[1] }
				}
				m.showAutocomplete = false
			}
			return m.handleEnter()

		// ── Up arrow: scroll or history ──
		case tea.KeyUp:
			if m.input == "" && len(m.inputHistory) > 0 {
				if m.historyIdx < len(m.inputHistory)-1 {
					m.historyIdx++
					m.input = m.inputHistory[len(m.inputHistory)-1-m.historyIdx]
				}
			} else if m.ready {
				m.viewport.LineUp(1)
			}

		// ── Down arrow: scroll or history ──
		case tea.KeyDown:
			if m.input == "" && m.historyIdx > 0 {
				m.historyIdx--
				m.input = m.inputHistory[len(m.inputHistory)-1-m.historyIdx]
			} else if m.ready {
				m.viewport.LineDown(1)
			}

		// ── PgUp/PgDn/Home/End: viewport scrolling ──
		case tea.KeyPgUp:
			if m.ready { m.viewport.HalfViewUp() }
		case tea.KeyPgDown:
			if m.ready { m.viewport.HalfViewDown() }
		case tea.KeyHome:
			if m.ready { m.viewport.GotoTop() }
		case tea.KeyEnd:
			if m.ready { m.viewport.GotoBottom() }

		// ── Backspace ──
		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.updateAutocomplete()
			}

		// ── Delete ──
		case tea.KeyDelete:
			if len(m.input) > 0 { m.input = m.input[:len(m.input)-1] }

		// ── Regular key input ──
		case tea.KeyRunes, tea.KeySpace:
			runeStr := msg.String()
			m.input += runeStr
			m.updateAutocomplete()
			m.historyIdx = 0

			// Detect paste (goose-style: if delta > threshold, enter paste mode)
			if len(runeStr) > 100 {
				m.pasteMode = true
				m.pasteContent = runeStr
			}

			// Detect bash mode transition (pi-mono: onChange handler)
			wasBash := m.bashMode
			m.bashMode = strings.HasPrefix(strings.TrimSpace(m.input), "!")
			if wasBash != m.bashMode {
				// Border color change is visual-only in View()
			}
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
		// Process follow-up queue (pi-mono: after agent_end, send queued messages)
		if len(m.followUpQueue) > 0 {
			next := m.followUpQueue[0]
			m.followUpQueue = m.followUpQueue[1:]
			m.input = next
			cmds = append(cmds, func() tea.Msg { return tea.KeyMsg{Type: tea.KeyEnter} })
		}
		cmds = append(cmds, fetchGitBranch(m.workingDir))

	case StreamChunkMsg:
		m.streamContent += msg.Chunk
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
			ToolKind:  DetectToolKind(msg.ToolName),
			ToolStatus: "completed",
			Expanded:  m.toolOutputExpanded,
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

// ─── Input handling helpers ───────────────────────────────────────────

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	m.showAutocomplete = false

	// Paste mode: send full paste content (goose-style)
	if m.pasteMode {
		m.pasteMode = false
		content := m.pasteContent
		m.pasteContent = ""
		m.input = ""
		m.entries = append(m.entries, ChatEntry{
			Type:      EntryUser,
			Content:   content,
			Timestamp: time.Now(),
		})
		m.loading = true
		m.syncViewport()
		return m, func() tea.Msg {
			response, err := buildPromptResponse(m.director, content)
			if err != nil {
				return AgentResponseMsg{Content: fmt.Sprintf("Error: %v", err), Provider: m.provider, Model: m.modelName}
			}
			return AgentResponseMsg{Content: response.Display, Provider: m.provider, Model: m.modelName}
		}
	}

	if strings.TrimSpace(m.input) == "" { return m, nil }

	req := strings.TrimSpace(m.input)

	// Store in input history
	m.inputHistory = append(m.inputHistory, req)
	if len(m.inputHistory) > 100 { m.inputHistory = m.inputHistory[len(m.inputHistory)-100:] }
	m.input = ""

	// Bash mode: ! prefix runs shell directly (pi-mono feature)
	if strings.HasPrefix(req, "!!") {
		cmd := strings.TrimSpace(strings.TrimPrefix(req, "!!"))
		m.entries = append(m.entries, ChatEntry{Type: EntryBashMode, Content: cmd, Timestamp: time.Now()})
		m.loading = true
		m.syncViewport()
		return m, func() tea.Msg {
			start := time.Now()
			out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
			dur := time.Since(start)
			return ToolExecMsg{ToolName: "bash", Args: cmd, Output: string(out), IsError: err != nil, Duration: dur}
		}
	}
	if strings.HasPrefix(req, "!") && !strings.HasPrefix(req, "!!") {
		cmd := strings.TrimSpace(strings.TrimPrefix(req, "!"))
		m.entries = append(m.entries, ChatEntry{Type: EntryBashMode, Content: cmd, Timestamp: time.Now()})
		m.loading = true
		m.syncViewport()
		return m, func() tea.Msg {
			start := time.Now()
			out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
			dur := time.Since(start)
			return ToolExecMsg{ToolName: "bash", Args: cmd, Output: string(out), IsError: err != nil, Duration: dur}
		}
	}

	// Slash command
	if strings.HasPrefix(req, "/") { return m.handleSlashCmd(req) }

	// Shell proposal (pi-mono: ?? prefix)
	if strings.HasPrefix(req, "??") {
		query := strings.TrimSpace(strings.TrimPrefix(req, "??"))
		m.entries = append(m.entries, ChatEntry{Type: EntryUser, Content: "?? " + query, Timestamp: time.Now()})
		m.loading = true
		m.syncViewport()
		return m, func() tea.Msg {
			response, err := buildShellProposal(m.director, query)
			if err != nil { return fmt.Sprintf("Error: %v", err) }
			return response
		}
	}

	// Regular user message
	m.entries = append(m.entries, ChatEntry{Type: EntryUser, Content: req, Timestamp: time.Now()})
	m.loading = true
	m.syncViewport()

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
		if plan.Execution.Route.Provider != "" { provider = plan.Execution.Route.Provider }
		if plan.Execution.Route.Model != "" { model = plan.Execution.Route.Model }
		return AgentResponseMsg{Content: response.Display, Provider: provider, Model: model, Plan: &plan}
	}
}

func (m model) handleSlashCmd(req string) (tea.Model, tea.Cmd) {
	mdl, cmd := ProcessSlashCommand(req, &m)
	m = mdl.(model)
	m.syncViewport()
	return m, cmd
}

// ─── Permission dialog (goose-style) ──────────────────────────────────

func (m model) updatePermissionDialog(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.permissionEntry == nil {
		m.showPermission = false
		return m, nil
	}
	opts := m.permissionEntry.PermissionOptions
	switch msg.Type {
	case tea.KeyUp:
		if m.permissionEntry.PermissionIdx > 0 { m.permissionEntry.PermissionIdx-- }
		return m, nil
	case tea.KeyDown:
		if m.permissionEntry.PermissionIdx < len(opts)-1 { m.permissionEntry.PermissionIdx++ }
		return m, nil
	case tea.KeyEnter:
		if m.permissionEntry.PermissionIdx < len(opts) {
			selected := opts[m.permissionEntry.PermissionIdx]
			m.permissionEntry.PermissionResolved = true
			m.entries = append(m.entries, ChatEntry{
				Type:      EntrySystem,
				Content:   m.theme.AccentText(fmt.Sprintf("[Permission] %s: %s", selected.Kind, selected.Name)),
				Timestamp: time.Now(),
			})
		}
		m.showPermission = false
		m.permissionEntry = nil
		m.syncViewport()
		return m, nil
	case tea.KeyEsc:
		m.showPermission = false
		m.permissionEntry = nil
		m.syncViewport()
		return m, nil
	default:
		// Keyboard shortcuts (goose-style: y=allow_once, a=allow_always, n=reject_once, N=reject_always)
		keyMap := map[string]string{"y": "allow_once", "a": "allow_always", "n": "reject_once", "N": "reject_always"}
		if kind, ok := keyMap[msg.String()]; ok {
			for _, opt := range opts {
				if opt.Kind == kind {
					m.permissionEntry.PermissionResolved = true
					m.entries = append(m.entries, ChatEntry{
						Type:      EntrySystem,
						Content:   m.theme.AccentText(fmt.Sprintf("[Permission] %s: %s", opt.Kind, opt.Name)),
						Timestamp: time.Now(),
					})
					break
				}
			}
			m.showPermission = false
			m.permissionEntry = nil
			m.syncViewport()
			return m, nil
		}
	}
	return m, nil
}

// ─── Model selector (pi-mono/opencode style) ──────────────────────────

func (m model) updateModelSelector(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.modelSelectorIdx > 0 { m.modelSelectorIdx-- }
		return m, nil
	case tea.KeyDown:
		if m.modelSelectorIdx < len(m.availableModels)-1 { m.modelSelectorIdx++ }
		return m, nil
	case tea.KeyEnter:
		if m.modelSelectorIdx >= 0 && m.modelSelectorIdx < len(m.availableModels) {
			m.modelName = m.availableModels[m.modelSelectorIdx]
			if m.session != nil { m.session.SetModel(m.modelName) }
			m.entries = append(m.entries, ChatEntry{
				Type:      EntrySystem,
				Content:   m.theme.AccentText(fmt.Sprintf("[Model] set to %s/%s", m.provider, m.modelName)),
				Timestamp: time.Now(),
			})
		}
		m.showModelSelector = false
		m.syncViewport()
		return m, nil
	case tea.KeyEsc:
		m.showModelSelector = false
		return m, nil
	default:
		// Filter by typing
		if msg.Type == tea.KeyRunes || msg.Type == tea.KeySpace {
			m.modelSelectorFilter += msg.String()
			m.modelSelectorIdx = 0
		}
		if msg.Type == tea.KeyBackspace && len(m.modelSelectorFilter) > 0 {
			m.modelSelectorFilter = m.modelSelectorFilter[:len(m.modelSelectorFilter)-1]
			m.modelSelectorIdx = 0
		}
		return m, nil
	}
}

// ─── Command palette (opencode style) ─────────────────────────────────

func (m model) updateCommandPalette(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	filtered := m.filteredCommands()
	switch msg.Type {
	case tea.KeyUp:
		if m.commandPaletteIdx > 0 { m.commandPaletteIdx-- }
		return m, nil
	case tea.KeyDown:
		if m.commandPaletteIdx < len(filtered)-1 { m.commandPaletteIdx++ }
		return m, nil
	case tea.KeyEnter:
		if m.commandPaletteIdx >= 0 && m.commandPaletteIdx < len(filtered) {
			cmd := "/" + filtered[m.commandPaletteIdx].Name
			m.showCommandPalette = false
			m.commandPaletteFilter = ""
			return m.handleSlashCmd(cmd)
		}
		m.showCommandPalette = false
		return m, nil
	case tea.KeyEsc:
		m.showCommandPalette = false
		m.commandPaletteFilter = ""
		return m, nil
	default:
		if msg.Type == tea.KeyRunes || msg.Type == tea.KeySpace {
			m.commandPaletteFilter += msg.String()
			m.commandPaletteIdx = 0
		}
		if msg.Type == tea.KeyBackspace && len(m.commandPaletteFilter) > 0 {
			m.commandPaletteFilter = m.commandPaletteFilter[:len(m.commandPaletteFilter)-1]
			m.commandPaletteIdx = 0
		}
		return m, nil
	}
}

func (m model) filteredCommands() []SlashCommand {
	if m.commandPaletteFilter == "" { return BuiltinSlashCommands }
	prefix := strings.ToLower(m.commandPaletteFilter)
	var filtered []SlashCommand
	for _, cmd := range BuiltinSlashCommands {
		if strings.Contains(strings.ToLower(cmd.Name), prefix) || strings.Contains(strings.ToLower(cmd.Description), prefix) {
			filtered = append(filtered, cmd)
		}
	}
	return filtered
}

// ─── Suspend to shell (pi-mono: handleCtrlZ) ─────────────────────────

func (m model) suspendToShell() tea.Cmd {
	return func() tea.Msg {
		// Save TUI state, suspend, open subshell
		shell := os.Getenv("SHELL")
		if shell == "" {
			if runtime.GOOS == "windows" {
				shell = "cmd.exe"
			} else {
				shell = "/bin/sh"
			}
		}
		cmd := exec.Command(shell)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Temporarily restore terminal state
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan)
		defer signal.Reset()

		cmd.Run()
		return ChatEntry{Type: EntrySystem, Content: DefaultTheme.Dim("[Suspend] returned from shell"), Timestamp: time.Now()}
	}
}

// ─── Autocomplete helpers ─────────────────────────────────────────────

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
	if !strings.HasPrefix(prefix, "/") { return BuiltinSlashCommands }
	prefix = prefix[1:]
	var filtered []SlashCommand
	for _, cmd := range BuiltinSlashCommands {
		if strings.HasPrefix(cmd.Name, prefix) { filtered = append(filtered, cmd) }
	}
	return filtered
}

// ─── Tree browser / pinned pane keyboard handling ─────────────────────

func (m model) updateTreeBrowser(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visible := visibleTreeBrowserItems(m.browserItems, m.browserFilter, m.browserCollapsed)
	switch msg.Type {
	case tea.KeyEsc:
		if m.browserConfirmPending { m.browserConfirmPending = false; return m, nil }
		m.browserActive = false; m.browserFilter = ""; m.browserConfirmPending = false
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
			if !m.browserConfirmPending { m.browserConfirmPending = true; return m, nil }
			display, err := openSelectedTreeBrowser(m.director.WorkingDir, m.foundationSessionID, visible, m.browserIndex, 128)
			if err != nil {
				m.entries = append(m.entries, ChatEntry{Type: EntrySystem, Content: m.theme.ErrorText(fmt.Sprintf("[Error] %v", err)), Timestamp: time.Now()})
			} else {
				m.entries = append(m.entries, ChatEntry{Type: EntrySystem, Content: m.theme.SuccessText(display), Timestamp: time.Now()})
			}
			m.browserActive = false; m.browserFilter = ""; m.browserConfirmPending = false
			m.syncViewport()
		}
		return m, nil
	default:
		if msg.Type == tea.KeyRunes || msg.Type == tea.KeySpace {
			m.browserFilter += msg.String(); m.browserIndex = 0
		}
		if msg.Type == tea.KeyBackspace && len(m.browserFilter) > 0 {
			m.browserFilter = m.browserFilter[:len(m.browserFilter)-1]; m.browserIndex = 0
		}
		return m, nil
	}
}

func (m model) updatePinnedPane(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visible := visibleTreeBrowserItems(m.browserItems, m.browserFilter, m.browserCollapsed)
	switch msg.Type {
	case tea.KeyEsc, tea.KeyTab:
		m.browserPinnedFocus = false; return m, nil
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
	if m.quitting { return m.theme.Dim("Goodbye.\n") }

	// Permission dialog overlay (goose-style)
	if m.showPermission && m.permissionEntry != nil {
		overlay := RenderEntry(*m.permissionEntry, m.theme)
		chatView := m.renderChatArea()
		return chatView + "\n" + overlay + "\n" + m.renderInputBar()
	}

	// Model selector overlay (pi-mono/opencode style)
	if m.showModelSelector {
		overlay := RenderModelSelector(m.modelName, m.availableModels, m.theme)
		chatView := m.renderChatArea()
		return chatView + "\n" + overlay + "\n" + m.renderInputBar()
	}

	// Command palette overlay (opencode style)
	if m.showCommandPalette {
		filtered := m.filteredCommands()
		overlay := RenderAutocomplete(filtered, m.commandPaletteIdx, 10, m.theme)
		chatView := m.renderChatArea()
		return chatView + "\n" + overlay + "\n" + m.renderInputBar()
	}

	// Dashboard mode
	if m.dashboardActive {
		chatContent := m.renderAllEntries()
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

// ─── Chat area (scrollable viewport) ──────────────────────────────────

func (m model) renderChatArea() string {
	if m.ready { return m.viewport.View() }
	return m.renderAllEntries()
}

// ─── Render all entries ────────────────────────────────────────────────

func (m model) renderAllEntries() string {
	var lines []string
	for _, entry := range m.entries {
		lines = append(lines, RenderEntry(entry, m.theme))
		lines = append(lines, "") // blank line between entries
	}

	// Loading / streaming indicator
	if m.loading && m.streaming && m.streamContent != "" {
		lines = append(lines, m.theme.Fg(m.theme.Accent, "┃ ")+m.streamContent)
	} else if m.loading {
		msg := m.thinkingLabel
		if len(m.followUpQueue) > 0 {
			msg += fmt.Sprintf(" (%d queued)", len(m.followUpQueue))
		}
		msg += fmt.Sprintf(" (%s to interrupt)", m.theme.KeyHint("Ctrl+C", ""))
		lines = append(lines, m.theme.Fg(m.theme.ThinkingCol, m.theme.Italic(msg))+" "+m.spinner.View())
	}

	return strings.Join(lines, "\n")
}

// ─── Footer ───────────────────────────────────────────────────────────

func (m model) renderFooter() string {
	return RenderFooter(m.workingDir, m.gitBranch, m.provider, m.modelName, m.toolCount,
		m.totalInputTok, m.totalOutTok, m.totalCost, m.contextPct, m.contextWindow,
		m.thinkingLevel, m.theme, m.width)
}

// ─── Input bar ────────────────────────────────────────────────────────

func (m model) renderInputBar() string {
	t := m.theme
	var prompt, inputDisplay string

	if m.loading {
		prompt = t.Bold(t.Fg(t.ThinkingCol, "… "))
		inputDisplay = t.Fg(t.ThinkingCol, t.Italic(m.thinkingLabel))
		if len(m.followUpQueue) > 0 {
			inputDisplay += t.Fg(t.Gold, t.Italic(fmt.Sprintf(" [%d queued]", len(m.followUpQueue))))
		}
		inputDisplay += " " + m.spinner.View()
	} else if m.pasteMode {
		prompt = t.Bold(t.Fg(t.BashModeCol, "📋 "))
		preview := m.pasteContent
		if len(preview) > 60 { preview = preview[:57] + "…" }
		charCount := len(m.pasteContent)
		inputDisplay = t.Fg(t.TextColor, preview) + t.Dim(fmt.Sprintf(" (%d chars)", charCount))
		inputDisplay += "\n" + t.Dim(t.Italic("enter to send · esc to clear"))
	} else if m.bashMode || strings.HasPrefix(m.input, "!") {
		prompt = t.Bold(t.Fg(t.BashModeCol, "⚡ "))
		inputDisplay = t.Fg(t.TextColor, m.input) + "▎"
	} else {
		prompt = t.Bold(t.AccentText("> "))
		inputDisplay = t.Fg(t.TextColor, m.input) + "▎"
	}

	// Queue indicator (goose-style: "message queued — will send when finished")
	queueHint := ""
	if m.loading && len(m.followUpQueue) > 0 {
		queueHint = "\n" + t.Fg(t.Gold, t.Italic("message queued — will send when agent finishes"))
	}

	// Autocomplete dropdown
	autocomplete := ""
	if m.showAutocomplete && m.isSlashContext() {
		autocomplete = "\n" + RenderAutocomplete(m.filteredAutocomplete(), m.autocompleteIndex, m.autocompleteMaxVis, t)
	}

	return prompt + inputDisplay + queueHint + autocomplete
}

// ─── Tool sidebar (dashboard mode) ────────────────────────────────────

func (m model) renderToolSidebar() string {
	t := m.theme
	var lines []string
	lines = append(lines, t.Bold(t.AccentText("Active Tools")))
	lines = append(lines, "")

	if m.registry != nil {
		groups := map[string]int{}
		for _, tool := range m.registry.Tools {
			groups[toolGroupName(tool.Name)]++
		}
		for name, count := range groups {
			kindIcon := ToolKindIcon(DetectToolKind(name))
			lines = append(lines, kindIcon+" "+t.Fg(t.ToolTitle, name)+t.Dim(fmt.Sprintf(" (%d)", count)))
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
	if len(recent) > 5 { recent = recent[len(recent)-5:] }
	if len(recent) > 0 {
		lines = append(lines, "")
		lines = append(lines, t.Bold(t.Dim("Recent Executions")))
		for _, run := range recent {
			icon := t.SuccessText("✓")
			if run.IsError { icon = t.ErrorText("✗") }
			kindIcon := ToolKindIcon(DetectToolKind(run.ToolName))
			lines = append(lines, icon+" "+kindIcon+" "+t.Fg(t.ToolTitle, run.ToolName)+t.Dim(fmt.Sprintf(" (%v)", run.Duration.Round(time.Millisecond))))
		}
	}

	return strings.Join(lines, "\n")
}

// ─── Metrics ──────────────────────────────────────────────────────────

func (m model) renderMetrics() string {
	return fmt.Sprintf("Tokens: %s in / %s out │ Cost: $%.3f │ Scope: Project │ Tools: %d",
		formatTokens(m.totalInputTok), formatTokens(m.totalOutTok), m.totalCost, m.toolCount)
}

// ─── Helpers ──────────────────────────────────────────────────────────

func (m *model) syncViewport() {
	m.viewport.SetContent(m.renderAllEntries())
	m.viewport.GotoBottom()
}

func formatTokens(count int) string {
	if count < 1000 { return fmt.Sprintf("%d", count) }
	if count < 10000 { return fmt.Sprintf("%.1fk", float64(count)/1000) }
	if count < 1000000 { return fmt.Sprintf("%dk", count/1000) }
	return fmt.Sprintf("%.1fM", float64(count)/1000000)
}

func toolGroupName(name string) string {
	parts := strings.SplitN(name, "_", 2)
	if len(parts) > 0 { return parts[0] }
	return name
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
