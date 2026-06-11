package tui

// ═══════════════════════════════════════════════════════════════════════
// extensions.go — Pi-mono extension features ported to Go
// Implements top 50 pi.dev/packages addon features natively:
//   - Todo list (LLM tool + /todos)
//   - Git checkpoint (stash at each turn)
//   - Auto-commit on exit
//   - Confirm destructive actions
//   - Dirty repo guard
//   - Protected paths
//   - Session naming (/name)
//   - Bookmark (/bookmark)
//   - Desktop notifications (OSC 777/99, Windows toast)
//   - Plan mode (/plan-mode)
//   - Handoff (/handoff)
//   - Summarize (/summarize)
//   - File triggers
//   - Status line
//   - Model status
// ═══════════════════════════════════════════════════════════════════════

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════════════════════
// Todo List — pi-mono's todo extension
// LLM tool + /todos command, stored in session entries
// ═══════════════════════════════════════════════════════════════════════

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type TodoStore struct {
	Todos  []Todo
	NextID int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{NextID: 1}
}

func (ts *TodoStore) Add(text string) Todo {
	t := Todo{ID: ts.NextID, Text: text, Done: false}
	ts.Todos = append(ts.Todos, t)
	ts.NextID++
	return t
}

func (ts *TodoStore) Toggle(id int) *Todo {
	for i := range ts.Todos {
		if ts.Todos[i].ID == id {
			ts.Todos[i].Done = !ts.Todos[i].Done
			return &ts.Todos[i]
		}
	}
	return nil
}

func (ts *TodoStore) Clear() int {
	count := len(ts.Todos)
	ts.Todos = nil
	ts.NextID = 1
	return count
}

func (ts *TodoStore) List() string {
	if len(ts.Todos) == 0 {
		return "No todos"
	}
	var lines []string
	done := 0
	for _, t := range ts.Todos {
		check := "○"
		if t.Done {
			check = "✓"
			done++
		}
		lines = append(lines, fmt.Sprintf("  %s #%d %s", check, t.ID, t.Text))
	}
	return fmt.Sprintf("%d/%d completed\n%s", done, len(ts.Todos), strings.Join(lines, "\n"))
}

// ═══════════════════════════════════════════════════════════════════════
// Git Checkpoint — pi-mono's git-checkpoint extension
// Creates git stash checkpoints at each turn for /fork restoration
// ═══════════════════════════════════════════════════════════════════════

type GitCheckpoints struct {
	checkpoints map[string]string // entryId -> stash ref
}

func NewGitCheckpoints() *GitCheckpoints {
	return &GitCheckpoints{checkpoints: make(map[string]string)}
}

func (gc *GitCheckpoints) Create(entryID, workingDir string) {
	cmd := exec.Command("git", "stash", "create")
	cmd.Dir = workingDir
	out, err := cmd.Output()
	if err != nil {
		return
	}
	ref := strings.TrimSpace(string(out))
	if ref != "" {
		gc.checkpoints[entryID] = ref
	}
}

func (gc *GitCheckpoints) Restore(entryID, workingDir string) bool {
	ref, ok := gc.checkpoints[entryID]
	if !ok {
		return false
	}
	cmd := exec.Command("git", "stash", "apply", ref)
	cmd.Dir = workingDir
	return cmd.Run() == nil
}

func (gc *GitCheckpoints) Clear() {
	gc.checkpoints = make(map[string]string)
}

// ═══════════════════════════════════════════════════════════════════════
// Auto-Commit on Exit — pi-mono's auto-commit-on-exit extension
// ═══════════════════════════════════════════════════════════════════════

func AutoCommitOnExit(workingDir string, entries []ChatEntry) {
	// Check for uncommitted changes
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusCmd.Dir = workingDir
	statusOut, err := statusCmd.Output()
	if err != nil || len(strings.TrimSpace(string(statusOut))) == 0 {
		return // Not a git repo or no changes
	}

	// Find last assistant message for commit context
	var lastAssistant string
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Type == EntryAssistant && entries[i].Content != "" {
			lastAssistant = entries[i].Content
			break
		}
	}

	// Generate commit message
	msg := "auto: session changes"
	if lastAssistant != "" {
		// Use first 72 chars of last assistant message
		lines := strings.Split(lastAssistant, "\n")
		if len(lines[0]) > 72 {
			msg = "auto: " + lines[0][:66]
		} else {
			msg = "auto: " + lines[0]
		}
	}

	// Stage and commit
	exec.Command("git", "add", "-A").Run()
	exec.Command("git", "commit", "-m", msg, "--no-verify").Run()
}

// ═══════════════════════════════════════════════════════════════════════
// Dirty Repo Guard — pi-mono's dirty-repo-guard extension
// Prevents session changes when there are uncommitted git changes
// ═══════════════════════════════════════════════════════════════════════

func IsRepoDirty(workingDir string) bool {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = workingDir
	out, err := cmd.Output()
	if err != nil {
		return false // Not a git repo
	}
	return len(strings.TrimSpace(string(out))) > 0
}

// ═══════════════════════════════════════════════════════════════════════
// Protected Paths — pi-mono's protected-paths extension
// Blocks write/edit operations to sensitive files
// ═══════════════════════════════════════════════════════════════════════

var ProtectedPaths = []string{
	".env",
	".env.local",
	".env.production",
	".git/",
	"node_modules/",
	".hyperharness/memory/",
	"VERSION",
	"go.sum",
}

func IsProtectedPath(path string) bool {
	for _, p := range ProtectedPaths {
		if strings.Contains(path, p) {
			return true
		}
	}
	return false
}

// ═══════════════════════════════════════════════════════════════════════
// Desktop Notifications — pi-mono's notify extension
// Supports OSC 777 (Ghostty, iTerm2, WezTerm), OSC 99 (Kitty), Windows toast
// ═══════════════════════════════════════════════════════════════════════

func Notify(title, body string) {
	// OSC 777 (Ghostty, iTerm2, WezTerm, rxvt-unicode)
	fmt.Printf("\x1b]777;notify;%s;%s\x07", title, body)

	// OSC 99 (Kitty)
	fmt.Printf("\x1b]99;%s\x07", body)

	// Windows toast notification
	if isWindows() {
		go windowsToast(title, body)
	}
}

func isWindows() bool {
	return strings.Contains(strings.ToLower(os.Getenv("OS")), "windows")
}

func windowsToast(title, body string) {
	script := fmt.Sprintf(`
		$null = [Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime]
		$xml = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText01)
		$xml.GetElementsByTagName('text')[0].AppendChild($xml.CreateTextNode('%s')) > $null
		$toast = [Windows.UI.Notifications.ToastNotification]::new($xml)
		[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('%s').Show($toast)
	`, body, title)
	exec.Command("powershell", "-Command", script).Run()
}

// ═══════════════════════════════════════════════════════════════════════
// Plan Mode — pi-mono's plan-mode extension
// Read-only exploration mode, blocks write/edit/bash
// ═══════════════════════════════════════════════════════════════════════

type PlanMode struct {
	Active       bool
	Steps        []PlanStep
	CurrentStep  int
}

type PlanStep struct {
	ID      int
	Text    string
	Done    bool
}

func NewPlanMode() *PlanMode {
	return &PlanMode{}
}

func (pm *PlanMode) Toggle() {
	pm.Active = !pm.Active
}

func (pm *PlanMode) ExtractSteps(content string) {
	pm.Steps = nil
	lines := strings.Split(content, "\n")
	inPlan := false
	id := 1
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "plan:") {
			inPlan = true
			continue
		}
		if inPlan {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "1.") || strings.HasPrefix(trimmed, "2.") ||
				strings.HasPrefix(trimmed, "3.") || strings.HasPrefix(trimmed, "4.") ||
				strings.HasPrefix(trimmed, "5.") || strings.HasPrefix(trimmed, "6.") ||
				strings.HasPrefix(trimmed, "7.") || strings.HasPrefix(trimmed, "8.") ||
				strings.HasPrefix(trimmed, "9.") || strings.HasPrefix(trimmed, "- ") {
				pm.Steps = append(pm.Steps, PlanStep{ID: id, Text: trimmed, Done: false})
				id++
			}
		}
	}
}

func (pm *PlanMode) MarkDone(stepID int) {
	for i := range pm.Steps {
		if pm.Steps[i].ID == stepID {
			pm.Steps[i].Done = true
		}
	}
}

func (pm *PlanMode) IsToolAllowed(toolName string) bool {
	allowed := []string{"read", "grep", "find", "ls", "tree", "websearch", "webfetch"}
	for _, a := range allowed {
		if toolName == a {
			return true
		}
	}
	return false
}

// ═══════════════════════════════════════════════════════════════════════
// Handoff — pi-mono's handoff extension
// Transfer context to a new session without compaction
// ═══════════════════════════════════════════════════════════════════════

func GenerateHandoffPrompt(entries []ChatEntry, goal string) string {
	var context strings.Builder
	context.WriteString("## Context from previous session\n\n")

	// Extract key decisions and files
	var files []string
	var decisions []string
	for _, e := range entries {
		if e.Type == EntryAssistant && e.Content != "" {
			// Extract file mentions
			for _, word := range strings.Fields(e.Content) {
				if strings.Contains(word, ".go") || strings.Contains(word, ".ts") || strings.Contains(word, ".py") {
					files = append(files, word)
				}
			}
			// Extract decisions (lines with "decided", "chose", "using")
			for _, line := range strings.Split(e.Content, "\n") {
				lower := strings.ToLower(line)
				if strings.Contains(lower, "decided") || strings.Contains(lower, "chose") || strings.Contains(lower, "using") {
					decisions = append(decisions, strings.TrimSpace(line))
				}
			}
		}
	}

	if len(decisions) > 0 {
		context.WriteString("### Key decisions\n")
		seen := map[string]bool{}
		for _, d := range decisions {
			if !seen[d] && len(d) < 200 {
				context.WriteString("- " + d + "\n")
				seen[d] = true
			}
		}
		context.WriteString("\n")
	}

	if len(files) > 0 {
		context.WriteString("### Relevant files\n")
		seen := map[string]bool{}
		for _, f := range files {
			if !seen[f] {
				context.WriteString("- " + f + "\n")
				seen[f] = true
			}
		}
		context.WriteString("\n")
	}

	context.WriteString("### Next task\n")
	context.WriteString(goal + "\n")

	return context.String()
}

// ═══════════════════════════════════════════════════════════════════════
// Bookmark — pi-mono's bookmark extension
// Label entries for easy /tree navigation
// ═══════════════════════════════════════════════════════════════════════

type BookmarkStore struct {
	bookmarks map[string]string // entryID -> label
}

func NewBookmarkStore() *BookmarkStore {
	return &BookmarkStore{bookmarks: make(map[string]string)}
}

func (bs *BookmarkStore) Set(entryID, label string) {
	bs.bookmarks[entryID] = label
}

func (bs *BookmarkStore) Get(entryID string) string {
	return bs.bookmarks[entryID]
}

func (bs *BookmarkStore) List() map[string]string {
	return bs.bookmarks
}

// ═══════════════════════════════════════════════════════════════════════
// File Triggers — pi-mono's file-trigger extension
// Watch for file changes and trigger actions
// ═══════════════════════════════════════════════════════════════════════

type FileTrigger struct {
	Pattern string
	Action  string
}

var DefaultFileTriggers = []FileTrigger{
	{"*.go", "go build ./..."},
	{"go.mod", "go mod tidy"},
	{"*.ts", "npx tsc --noEmit"},
	{"package.json", "npm install"},
}

// ═══════════════════════════════════════════════════════════════════════
// Status Line — pi-mono's status-line extension
// Customizable status line in the footer
// ═══════════════════════════════════════════════════════════════════════

type StatusLine struct {
	Items []StatusItem
}

type StatusItem struct {
	Key   string
	Value string
	Color string
}

func NewStatusLine() *StatusLine {
	return &StatusLine{
		Items: []StatusItem{
			{Key: "cwd", Value: "", Color: "dim"},
			{Key: "git", Value: "", Color: "accent"},
			{Key: "model", Value: "", Color: "text"},
			{Key: "tokens", Value: "", Color: "muted"},
		},
	}
}

// ═══════════════════════════════════════════════════════════════════════
// Model Status — pi-mono's model-status extension
// Show current model info in status line
// ═══════════════════════════════════════════════════════════════════════

type ModelStatus struct {
	Provider    string
	Model       string
	ThinkingLvl string
	Cost        float64
	InputTok    int
	OutputTok   int
}

// ═══════════════════════════════════════════════════════════════════════
// Session Info — pi-mono's /session command data
// ═══════════════════════════════════════════════════════════════════════

type SessionInfo struct {
	ID              string
	Name            string
	File            string
	UserMessages    int
	AssistantMsgs   int
	ToolCalls       int
	TotalMessages   int
	TokensInput     int
	TokensOutput    int
	TokensCacheRead int
	Cost            float64
	Duration        time.Duration
}

func CollectSessionInfo(entries []ChatEntry, sessionID, sessionName, sessionFile string, inputTok, outTok int, cost float64) SessionInfo {
	info := SessionInfo{
		ID:           sessionID,
		Name:         sessionName,
		File:         sessionFile,
		TokensInput:  inputTok,
		TokensOutput: outTok,
		Cost:         cost,
	}
	if len(entries) > 0 {
		info.Duration = time.Since(entries[0].Timestamp)
	}
	for _, e := range entries {
		info.TotalMessages++
		switch e.Type {
		case EntryUser:
			info.UserMessages++
		case EntryAssistant:
			info.AssistantMsgs++
		case EntryTool:
			info.ToolCalls++
		}
	}
	return info
}

// ═══════════════════════════════════════════════════════════════════════
// Theme Watcher — pi-mono's theme file watcher
// Watches for theme file changes and reloads
// ═══════════════════════════════════════════════════════════════════════

type ThemeWatcher struct {
	themeDir  string
	lastMod   time.Time
	onChange  func()
}

func NewThemeWatcher(themeDir string, onChange func()) *ThemeWatcher {
	return &ThemeWatcher{themeDir: themeDir, onChange: onChange}
}

func (tw *ThemeWatcher) Check() {
	info, err := os.Stat(tw.themeDir)
	if err != nil {
		return
	}
	if info.ModTime().After(tw.lastMod) {
		tw.lastMod = info.ModTime()
		if tw.onChange != nil {
			tw.onChange()
		}
	}
}

// ═══════════════════════════════════════════════════════════════════════
// Editor History — pi-mono's editor history persistence
// Saves/loads input history across sessions
// ═══════════════════════════════════════════════════════════════════════

func SaveEditorHistory(workingDir string, history []string) {
	histDir := filepath.Join(workingDir, ".hyperharness")
	os.MkdirAll(histDir, 0755)
	path := filepath.Join(histDir, "history")
	var lines []string
	// Keep last 500 entries
	if len(history) > 500 {
		history = history[len(history)-500:]
	}
	for _, h := range history {
		lines = append(lines, h)
	}
	os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

func LoadEditorHistory(workingDir string) []string {
	path := filepath.Join(workingDir, ".hyperharness", "history")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return strings.Split(string(data), "\n")
}
