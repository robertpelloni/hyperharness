package agents

// ═══════════════════════════════════════════════════════════════════════
// session.go — Agent session management
// Mirrors pi-mono's AgentSession: prompt, tool execution, compaction, abort
// ═══════════════════════════════════════════════════════════════════════

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ThinkingLevel mirrors pi-mono's thinking levels
type ThinkingLevel string

const (
	ThinkingOff     ThinkingLevel = "off"
	ThinkingMinimal ThinkingLevel = "minimal"
	ThinkingLow     ThinkingLevel = "low"
	ThinkingMedium  ThinkingLevel = "medium"
	ThinkingHigh    ThinkingLevel = "high"
	ThinkingXhigh   ThinkingLevel = "xhigh"
)

// AllThinkingLevels is the ordered list of thinking levels for cycling
var AllThinkingLevels = []ThinkingLevel{ThinkingOff, ThinkingMinimal, ThinkingLow, ThinkingMedium, ThinkingHigh, ThinkingXhigh}

// SessionEvent types (mirrors pi-mono's AgentSessionEvent)
type SessionEventType string

const (
	EventAgentStart       SessionEventType = "agent_start"
	EventMessageStart     SessionEventType = "message_start"
	EventMessageUpdate    SessionEventType = "message_update"
	EventMessageEnd       SessionEventType = "message_end"
	EventToolExecStart    SessionEventType = "tool_execution_start"
	EventToolExecUpdate   SessionEventType = "tool_execution_update"
	EventToolExecEnd      SessionEventType = "tool_execution_end"
	EventAgentEnd         SessionEventType = "agent_end"
	EventCompactionStart  SessionEventType = "compaction_start"
	EventCompactionEnd    SessionEventType = "compaction_end"
	EventAutoRetryStart   SessionEventType = "auto_retry_start"
	EventAutoRetryEnd     SessionEventType = "auto_retry_end"
	EventQueueUpdate      SessionEventType = "queue_update"
)

// SessionEvent mirrors pi-mono's AgentSessionEvent
type SessionEvent struct {
	Type       SessionEventType
	Message    Message
	ToolCallID string
	ToolName   string
	ToolArgs   string
	ToolOutput string
	ToolIsErr  bool
	ToolDur    time.Duration
	Err        error
	Reason     string
	Aborted    bool
	Attempt    int
	MaxAttempt int
	DelayMs    int
}

// SessionEventListener receives session events
type SessionEventListener func(event SessionEvent)

// AgentSession mirrors pi-mono's AgentSession
type AgentSession struct {
	director      *Director
	registry      ToolExecutor
	sessionDir    string
	workingDir     string
	model          string
	provider       string
	thinkingLevel  ThinkingLevel
	maxIterations  int
	isStreaming     bool
	isCompacting   bool
	isAborted      bool

	mu             sync.Mutex
	listeners      []SessionListener
	entries        []SessionEntry
	pendingTools   map[string]*ToolExecution

	// Auto-compaction
	autoCompact    bool
	contextWindow  int
	contextUsed    int

	// Auto-retry
	autoRetry      bool
	maxRetries     int
	retryAttempt   int

	// Session metadata
	sessionID      string
	sessionName    string
	createdAt      time.Time
	totalInputTok  int
	totalOutputTok int
	totalCost      float64
}

// ToolExecutor is an interface for tool execution (breaks import cycle with tools package)
type ToolExecutor interface {
	Find(name string) (ToolInfo, bool)
	ListTools() []ToolInfo
}

// ToolInfo describes a tool for the session
type ToolInfo struct {
	Name        string
	Description string
	Execute     func(args map[string]interface{}) (string, error)
}
	ID   int
	Fn   SessionEventListener
}

// SessionEntry is a persistent session record
type SessionEntry struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	ToolName  string    `json:"tool_name,omitempty"`
	ToolArgs  string    `json:"tool_args,omitempty"`
	ToolOut   string    `json:"tool_out,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// ToolExecution tracks an in-flight tool call
type ToolExecution struct {
	CallID   string
	Name     string
	Args     string
	Start    time.Time
	Duration time.Duration
	Output   string
	IsError  bool
	Done     bool
}

// NewAgentSession creates a new agent session (mirrors pi-mono's AgentSession constructor)
func NewAgentSession(director *Director, registry ToolExecutor, sessionDir string) *AgentSession {
	wd := "."
	if director != nil && director.WorkingDir != "" {
		wd = director.WorkingDir
	}
	return &AgentSession{
		director:      director,
		registry:      registry,
		sessionDir:    sessionDir,
		workingDir:    wd,
		model:         "auto",
		provider:      "hypercode",
		thinkingLevel: ThinkingOff,
		maxIterations: 25,
		pendingTools:  make(map[string]*ToolExecution),
		autoCompact:   true,
		contextWindow: 200000,
		autoRetry:     true,
		maxRetries:    3,
		createdAt:     time.Now(),
		sessionID:     fmt.Sprintf("session-%d", time.Now().Unix()),
	}
}

// Subscribe adds an event listener (mirrors pi-mono's subscribe)
func (s *AgentSession) Subscribe(fn SessionEventListener) func() {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := len(s.listeners)
	s.listeners = append(s.listeners, SessionListener{ID: id, Fn: fn})
	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		for i, l := range s.listeners {
			if l.ID == id {
				s.listeners = append(s.listeners[:i], s.listeners[i+1:]...)
				break
			}
		}
	}
}

func (s *AgentSession) emit(event SessionEvent) {
	s.mu.Lock()
	listeners := make([]SessionEventListener, 0, len(s.listeners))
	for _, l := range s.listeners {
		listeners = append(listeners, l.Fn)
	}
	s.mu.Unlock()

	for _, fn := range listeners {
		fn(event)
	}
}

// Prompt sends a user message and runs the agent loop (mirrors pi-mono's prompt)
func (s *AgentSession) Prompt(ctx context.Context, text string) error {
	s.isStreaming = true
	s.isAborted = false

	// Add user message to history
	s.entries = append(s.entries, SessionEntry{
		Role:      "user",
		Content:   text,
		Timestamp: time.Now(),
	})

	s.emit(SessionEvent{Type: EventAgentStart})

	// Add to director history
	if s.director != nil {
		s.director.History = append(s.director.History, Message{
			Role:    RoleUser,
			Content: text,
		})
	}

	// Run agent loop with tool execution
	for i := 0; i < s.maxIterations; i++ {
		if s.isAborted {
			s.emit(SessionEvent{Type: EventAgentEnd, Err: fmt.Errorf("aborted")})
			s.isStreaming = false
			return fmt.Errorf("aborted")
		}

		// Get LLM response
		s.emit(SessionEvent{
			Type:    EventMessageStart,
			Message: Message{Role: RoleAssistant, Content: ""},
		})

		if s.director == nil || s.director.Provider == nil {
			s.emit(SessionEvent{
				Type:    EventMessageEnd,
				Message: Message{Role: RoleAssistant, Content: "No provider configured"},
			})
			break
		}

		responseMsg, err := s.director.Provider.Chat(ctx, s.director.History, s.buildToolSchemas())
		if err != nil {
			s.emit(SessionEvent{
				Type:    EventMessageEnd,
				Message: Message{Role: RoleAssistant, Content: fmt.Sprintf("Error: %v", err)},
				Err:     err,
			})
			// Auto-retry on error
			if s.autoRetry && s.retryAttempt < s.maxRetries {
				s.retryAttempt++
				s.emit(SessionEvent{
					Type:       EventAutoRetryStart,
					Attempt:    s.retryAttempt,
					MaxAttempt: s.maxRetries,
					DelayMs:    2000,
				})
				select {
				case <-time.After(2 * time.Second):
				case <-ctx.Done():
					s.isStreaming = false
					return ctx.Err()
				}
				s.emit(SessionEvent{Type: EventAutoRetryEnd})
				continue
			}
			break
		}

		s.director.History = append(s.director.History, responseMsg)
		s.entries = append(s.entries, SessionEntry{
			Role:      "assistant",
			Content:   responseMsg.Content,
			Timestamp: time.Now(),
		})

		s.emit(SessionEvent{
			Type:    EventMessageUpdate,
			Message: responseMsg,
		})

		// Execute tool calls if any
		if len(responseMsg.ToolCalls) == 0 {
			s.emit(SessionEvent{Type: EventMessageEnd, Message: responseMsg})
			break
		}

		// Process each tool call
		for _, tc := range responseMsg.ToolCalls {
			callID := tc.ID
			if callID == "" {
				callID = fmt.Sprintf("tc-%d", time.Now().UnixNano())
			}

			s.emit(SessionEvent{
				Type:       EventToolExecStart,
				ToolCallID: callID,
				ToolName:   tc.Name,
				ToolArgs:   tc.Args,
			})

			output, dur, toolErr := s.executeTool(ctx, tc.Name, tc.Args)

			s.director.History = append(s.director.History, Message{
				Role:      RoleTool,
				ToolCallID: callID,
				Name:      tc.Name,
				Content:   output,
			})

			s.emit(SessionEvent{
				Type:       EventToolExecEnd,
				ToolCallID: callID,
				ToolName:   tc.Name,
				ToolOutput: output,
				ToolIsErr:  toolErr != nil,
				ToolDur:    dur,
			})
		}

		// Check context overflow for auto-compaction
		s.contextUsed += len(responseMsg.Content)
		if s.autoCompact && s.contextUsed > s.contextWindow {
			s.emit(SessionEvent{Type: EventCompactionStart, Reason: "overflow"})
			s.contextUsed = s.contextUsed / 2
			s.emit(SessionEvent{Type: EventCompactionEnd, Aborted: false})
		}
	}

	s.emit(SessionEvent{Type: EventAgentEnd})
	s.isStreaming = false
	s.retryAttempt = 0
	return nil
}

// executeTool finds and runs a tool from the registry
func (s *AgentSession) executeTool(ctx context.Context, name, argsStr string) (string, time.Duration, error) {
	start := time.Now()

	if s.registry == nil {
		return "no tool registry", 0, fmt.Errorf("no registry")
	}

	tool, found := s.registry.Find(name)
	if !found {
		return fmt.Sprintf("tool %q not found", name), time.Since(start), fmt.Errorf("not found")
	}

	// Parse args
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(argsStr), &args); err != nil {
		// Try as simple command arg
		args = map[string]interface{}{"command": argsStr, "query": argsStr, "path": argsStr, "content": argsStr}
	}

	result, err := tool.Execute(args)
	dur := time.Since(start)
	if err != nil {
		return result, dur, err
	}
	return result, dur, nil
}

// buildToolSchemas converts registry tools to the ILLMProvider Tool format
func (s *AgentSession) buildToolSchemas() []Tool {
	if s.registry == nil {
		return nil
	}
	var result []Tool
	for _, t := range s.registry.ListTools() {
		result = append(result, Tool{
			Name:        t.Name,
			Description: t.Description,
			Schema:      nil, // Would need to parse t.Parameters into map[string]interface{}
		})
	}
	// Limit to avoid overwhelming the LLM
	if len(result) > 100 {
		result = result[:100]
	}
	return result
}

// Abort cancels the current operation (mirrors pi-mono's abort)
func (s *AgentSession) Abort() {
	s.isAborted = true
	s.isStreaming = false
}

// CycleModel cycles through available models (mirrors pi-mono's cycleModel)
func (s *AgentSession) CycleModel(direction string) (string, string) {
	models := []string{"auto", "gemini-1.5-pro", "gpt-4", "claude-3-5-sonnet", "local"}
	idx := 0
	for i, m := range models {
		if m == s.model {
			idx = i
			break
		}
	}
	if direction == "backward" {
		idx = (idx - 1 + len(models)) % len(models)
	} else {
		idx = (idx + 1) % len(models)
	}
	s.model = models[idx]
	return s.provider, s.model
}

// CycleThinkingLevel cycles through thinking levels (mirrors pi-mono's cycleThinkingLevel)
func (s *AgentSession) CycleThinkingLevel() ThinkingLevel {
	for i, l := range AllThinkingLevels {
		if l == s.thinkingLevel {
			s.thinkingLevel = AllThinkingLevels[(i+1)%len(AllThinkingLevels)]
			return s.thinkingLevel
		}
	}
	s.thinkingLevel = ThinkingMinimal
	return s.thinkingLevel
}

// Compact runs context compaction (mirrors pi-mono's compact)
func (s *AgentSession) Compact() {
	s.emit(SessionEvent{Type: EventCompactionStart, Reason: "manual"})
	// Truncate early history, keep last N entries
	if len(s.entries) > 20 {
		s.entries = s.entries[len(s.entries)-20:]
	}
	if s.director != nil && len(s.director.History) > 10 {
		// Keep system prompt + recent history
		sys := s.director.History[0]
		recent := s.director.History[len(s.director.History)-9:]
		s.director.History = append([]Message{sys}, recent...)
	}
	s.contextUsed = s.contextUsed / 2
	s.emit(SessionEvent{Type: EventCompactionEnd, Aborted: false})
}

// GetState returns session state for display
func (s *AgentSession) GetState() map[string]interface{} {
	return map[string]interface{}{
		"model":          s.model,
		"provider":       s.provider,
		"thinkingLevel":  string(s.thinkingLevel),
		"isStreaming":     s.isStreaming,
		"isCompacting":   s.isCompacting,
		"entryCount":     len(s.entries),
		"contextUsed":    s.contextUsed,
		"contextWindow":  s.contextWindow,
		"autoCompact":    s.autoCompact,
		"sessionID":      s.sessionID,
		"sessionName":    s.sessionName,
		"totalInputTok":  s.totalInputTok,
		"totalOutputTok": s.totalOutputTok,
		"totalCost":      s.totalCost,
	}
}

// Save persists session to disk
func (s *AgentSession) Save() error {
	if s.sessionDir == "" {
		return nil
	}
	os.MkdirAll(s.sessionDir, 0755)
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.sessionDir, s.sessionID+".jsonl"), data, 0644)
}

// Load restores session from disk
func (s *AgentSession) Load(sessionID string) error {
	if s.sessionDir == "" {
		return fmt.Errorf("no session dir")
	}
	data, err := os.ReadFile(filepath.Join(s.sessionDir, sessionID+".jsonl"))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.entries)
}

// SetModel sets the model
func (s *AgentSession) SetModel(model string) {
	s.model = model
}

// SetProvider sets the provider
func (s *AgentSession) SetProvider(provider string) {
	s.provider = provider
}

// SetThinkingLevel sets the thinking level
func (s *AgentSession) SetThinkingLevel(level ThinkingLevel) {
	s.thinkingLevel = level
}

// SetSessionName sets the display name
func (s *AgentSession) SetSessionName(name string) {
	s.sessionName = name
}

// IsStreaming returns whether the agent is currently processing
func (s *AgentSession) IsStreaming() bool {
	return s.isStreaming
}

// GetModel returns the current model
func (s *AgentSession) GetModel() string {
	return s.model
}

// GetProvider returns the current provider
func (s *AgentSession) GetProvider() string {
	return s.provider
}

// GetThinkingLevel returns the current thinking level
func (s *AgentSession) GetThinkingLevel() ThinkingLevel {
	return s.thinkingLevel
}

// GetEntries returns all session entries
func (s *AgentSession) GetEntries() []SessionEntry {
	return s.entries
}

// GetContextUsage returns context usage info
func (s *AgentSession) GetContextUsage() (used, window int, pct float64) {
	return s.contextUsed, s.contextWindow, float64(s.contextUsed) / float64(s.contextWindow) * 100
}

// silence unused import
var _ = log.Println
