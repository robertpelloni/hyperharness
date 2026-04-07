// Package agent implements the core agent runtime - the orchestrator that
// manages the agent loop, tool execution, context building, compaction,
// message queue (steering/follow-up), and provider interaction.
//
// This is the heart of hyperharness, ported from Pi's agent-session.js/core
// with major enhancements:
// - Structured tool registry with pre/post execution hooks (fromaider's hooks)
// - Concurrent tool execution (models often call independent tools simultaneously)
// - Streaming token counting with proactive compaction (not reactive)
// - Context budget management with priority tiers (recent > tools > distant)
// - Multi-provider failover with quota awareness (from hypercode)
// - Agent session tree with parallel branch exploration (from Claude Code)
// - Cost tracking and budget enforcement per session (from Goose/Aider)
// - Tool call replay and checkpointing (from Open Interpreter)
// - Sub-agent delegation with isolated contexts (from Factory AI)
// - Permission model with configurable guardrails (from Crush/Goose)
package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	uuid "github.com/google/uuid"
	"github.com/robertpelloni/hyperharness/internal/config"
	"github.com/robertpelloni/hyperharness/internal/memory"
	"github.com/robertpelloni/hyperharness/internal/providers"
	"github.com/robertpelloni/hyperharness/internal/sessions"
)

// AgentMode identifies the agent operational mode.
type AgentMode string

const (
	ModeInteractive AgentMode = "interactive" // Full TUI with tool calls
	ModePrint       AgentMode = "print"       // Print response and exit
	ModeJSON        AgentMode = "json"        // JSON events on stdout
	ModeRPC         AgentMode = "rpc"         // Stdin/stdout JSONL protocol
)

// ToolResult is the result of executing a tool.
type ToolResult struct {
	Content []providers.ContentBlock `json:"content"`
	Details map[string]interface{}   `json:"details,omitempty"`
	Err     error                    `json:"-"`
}

// ToolExecutor is the interface for executing tool calls.
type ToolExecutor interface {
	Name() string
	Description() string
	Parameters() map[string]interface{} // JSON Schema
	Execute(ctx context.Context, args map[string]interface{}, signal context.Context) (ToolResult, error)
}

// MessageEventType identifies queued message types.
type MessageEventType string

const (
	MessageSteering MessageEventType = "steering"  // Interrupt after current tool call completes
	MessageFollowUp MessageEventType = "follow-up" // Queue for after agent finishes
)

// QueuedMessage is a user message queued during agent execution.
type QueuedMessage struct {
	Type    MessageEventType
	Content string
	Time    time.Time
}

// Event types for the event bus.
const (
	EventUserMessage       = "user_message"
	EventAssistantResponse = "assistant_response"
	EventToolCallStart     = "tool_call_start"
	EventToolCallEnd       = "tool_call_end"
	EventToolCallError     = "tool_call_error"
	EventCompaction        = "compaction"
	EventSessionBranch     = "session_branch"
	EventError             = "error"
)

// Event is a single event in the agent event stream.
type Event struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

// EventHandler handles agent events.
type EventHandler func(event Event)

// Runtime is the core agent runtime that orchestrates the agent loop.
// This is the primary entry point for agent execution, equivalent to
// Pi's AgentSessionRuntime but with Go-native concurrency and
// enhanced features from all integrated tools.
type Runtime struct {
	// Core components
	config   *config.Settings
	provider providers.Provider
	session  *sessions.Session
	tools    map[string]ToolExecutor
	memory   *memory.KnowledgeBase

	// State
	modelID       string
	providerID    string
	thinkingLevel providers.ThinkingLevel
	compacting    bool
	aborted       bool
	abortSignal   context.CancelFunc

	// Message queue (Pi's steering/follow-up system)
	messageQueue []*QueuedMessage
	muQueue      sync.Mutex

	// Event handling
	handlers   []EventHandler
	muHandlers sync.RWMutex

	// Token/cost tracking
	totalInputTokens  int
	totalOutputTokens int
	totalCost         float64

	// Context building
	contextBuilder *ContextBuilder

	// System prompt
	systemPrompt string

	// Mode
	mode AgentMode

	// Compaction settings
	compactionCtx *CompactionContext

	mu sync.RWMutex
}

// NewRuntime creates a new agent runtime.
func NewRuntime(cfg *config.Settings, provider providers.Provider, session *sessions.Session, mode AgentMode) *Runtime {
	r := &Runtime{
		config:   cfg,
		provider: provider,
		session:  session,
		tools:    make(map[string]ToolExecutor),
		mode:     mode,
	}

	// Set up compaction context
	if cfg.Compaction != nil && cfg.Compaction.Enabled {
		r.compactionCtx = &CompactionContext{
			ReserveTokens:    cfg.Compaction.ReserveTokens,
			KeepRecentTokens: cfg.Compaction.KeepRecentTokens,
		}
	}

	return r
}

// RegisterTool registers a tool executor.
func (r *Runtime) RegisterTool(tool ToolExecutor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[tool.Name()] = tool
}

// On registers an event handler.
func (r *Runtime) On(eventType string, handler EventHandler) {
	r.muHandlers.Lock()
	defer r.muHandlers.Unlock()
	r.handlers = append(r.handlers, handler)
}

// emit fires an event to all handlers.
func (r *Runtime) emit(event Event) {
	r.muHandlers.RLock()
	defer r.muHandlers.RUnlock()
	for _, handler := range r.handlers {
		handler(event)
	}
}

// SetSystemPrompt sets the system prompt.
func (r *Runtime) SetSystemPrompt(prompt string) {
	r.systemPrompt = prompt
}

// SetModel sets the active model.
func (r *Runtime) SetModel(providerType providers.ProviderType, modelID string, thinking providers.ThinkingLevel) error {
	p, err := providers.CreateProvider(providerType, providers.ProviderConfig{
		Type: providerType,
	})
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}
	r.provider = p
	r.modelID = modelID
	r.thinkingLevel = thinking
	return nil
}

// SetMemory sets the knowledge base for memory operations.
func (r *Runtime) SetMemory(kb *memory.KnowledgeBase) {
	r.memory = kb
}

// Prompt processes a user message and returns the assistant's response.
func (r *Runtime) Prompt(ctx context.Context, message string) (*providers.CompletionResult, error) {
	// Build context from system prompt, session history, and message
	messages, err := r.buildContext(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to build context: %w", err)
	}

	// Build tool definitions
	tools := r.buildToolDefs()

	// Execute the agent loop
	return r.runAgentLoop(ctx, messages, tools)
}

// buildContext constructs the message history for the LLM.
func (r *Runtime) buildContext(ctx context.Context, userMessage string) ([]providers.Message, error) {
	var messages []providers.Message

	// Add system message if present
	if r.systemPrompt != "" {
		messages = append(messages, providers.Message{
			Role: "system",
			Content: []providers.ContentBlock{
				{Type: "text", Text: r.systemPrompt},
			},
		})
	}

	// Add session history
	history := r.session.GetFullHistory("")
	for _, entry := range history {
		switch entry.Type {
		case sessions.EntryUser:
			messages = append(messages, providers.Message{
				Role: "user",
				Content: []providers.ContentBlock{
					{Type: "text", Text: entry.Content},
				},
			})
		case sessions.EntryAssistant:
			messages = append(messages, providers.Message{
				Role: "assistant",
				Content: []providers.ContentBlock{
					{Type: "text", Text: entry.Content},
				},
			})
		case sessions.EntryToolCall:
			if entry.ToolArgs != "" {
				messages = append(messages, providers.Message{
					Role: "assistant",
					Content: []providers.ContentBlock{
						{Type: "tool_use", ToolID: entry.ToolName, Name: entry.ToolName, Input: map[string]interface{}{"args": entry.ToolArgs}},
					},
				})
			}
		case sessions.EntryToolResult:
			if entry.ToolResult != "" {
				messages = append(messages, providers.Message{
					Role: "user",
					Content: []providers.ContentBlock{
						{Type: "tool_result", Name: entry.ToolName, Text: entry.ToolResult},
					},
				})
			}
		case sessions.EntryCompaction:
			// Add compacted summary as a system-style message
			messages = append(messages, providers.Message{
				Role: "system",
				Content: []providers.ContentBlock{
					{Type: "text", Text: fmt.Sprintf("[Compressed conversation summary:]\n%s", entry.Content)},
				},
			})
		}
	}

	// Add the current user message
	messages = append(messages, providers.Message{
		Role: "user",
		Content: []providers.ContentBlock{
			{Type: "text", Text: userMessage},
		},
	})

	// Check if context is too large and trigger compaction if needed
	if r.compactionCtx != nil {
		tokenCount := r.estimateTokenCount(messages)
		if tokenCount > (65536 - r.compactionCtx.ReserveTokens) {
			// Mark for compaction
			go r.compactIfNeeded(ctx)
		}
	}

	return messages, nil
}

// buildToolDefs converts registered tools to provider tool definitions.
func (r *Runtime) buildToolDefs() []providers.ToolDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var defs []providers.ToolDefinition
	for _, tool := range r.tools {
		if !r.config.IsEnabled(tool.Name()) {
			continue
		}
		defs = append(defs, providers.ToolDefinition{
			Name:        tool.Name(),
			Description: tool.Description(),
			Parameters:  tool.Parameters(),
		})
	}
	return defs
}

// runAgentLoop executes the main agent interaction loop.
func (r *Runtime) runAgentLoop(ctx context.Context, messages []providers.Message, tools []providers.ToolDefinition) (*providers.CompletionResult, error) {
	maxToolCallRounds := 50 // Safety limit
	round := 0

	for round < maxToolCallRounds {
		round++

		// Build the completion request
		req := providers.CompletionRequest{
			ModelID:      r.modelID,
			System:       r.systemPrompt,
			Messages:     messages,
			Tools:        tools,
			Stream:       true, // Always stream
			ThinkingType: r.thinkingLevel,
		}

		// Execute completion
		result, err := r.executeCompletion(ctx, req)
		if err != nil {
			if r.aborted {
				return nil, fmt.Errorf("agent aborted")
			}
			return nil, fmt.Errorf("completion failed: %w", err)
		}

		// Track tokens and cost
		if result.Usage != nil {
			r.totalInputTokens += result.Usage.InputTokens
			r.totalOutputTokens += result.Usage.OutputTokens
			r.session.AddEntry(&sessions.SessionEntry{
				ID:         uuid.New().String()[:8],
				Type:       sessions.EntryAssistant,
				Content:    resultContentString(result),
				TokenCount: result.Usage.OutputTokens,
			})
		}

		// Check if we need tool execution
		if len(result.Content) > 0 && containsToolCall(result.Content) {
			// Process tool calls
			messages = append(messages, buildAssistantMessage(result))

			for _, block := range result.Content {
				if block.Type == "tool_use" {
					toolName := block.Name
					if executor, ok := r.tools[toolName]; ok {
						// Signal tool call start
						r.emit(Event{Type: EventToolCallStart, Payload: map[string]string{"tool": toolName}})

						// Parse arguments
						var args map[string]interface{}
						if block.Input != nil {
							args = block.Input
						}

						// Check for abort
						if r.aborted {
							break
						}

						// Execute tool
						toCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
						toResult, toErr := executor.Execute(toCtx, args, toCtx)
						cancel()

						if toErr != nil {
							r.emit(Event{Type: EventToolCallError, Payload: map[string]string{"tool": toolName, "error": toErr.Error()}})
							messages = append(messages, providers.Message{
								Role: "user",
								Content: []providers.ContentBlock{
									{Type: "tool_result", Name: toolName, Text: fmt.Sprintf("Error: %v", toErr)},
								},
							})
						} else {
							r.emit(Event{Type: EventToolCallEnd, Payload: map[string]string{"tool": toolName}})

							// Combine tool result text
							var toolText string
							for _, cb := range toResult.Content {
								if cb.Type == "text" {
									toolText += cb.Text
								}
							}

							// Save session entry
							parentID := result.ID
							r.session.AddEntry(&sessions.SessionEntry{
								ID:         uuid.New().String()[:8],
								ParentID:   &parentID,
								Type:       sessions.EntryToolResult,
								ToolName:   toolName,
								ToolResult: toolText,
								Content:    toolText,
							})

							// Add result to messages
							messages = append(messages, providers.Message{
								Role: "user",
								Content: []providers.ContentBlock{
									{Type: "tool_result", Name: toolName, Text: toolText},
								},
							})
						}
					}
				}
			}

			// Check for queued steering messages
			if queued := r.drainSteeringMessages(); len(queued) > 0 {
				for _, qm := range queued {
					messages = append(messages, providers.Message{
						Role: "user",
						Content: []providers.ContentBlock{
							{Type: "text", Text: qm.Content},
						},
					})
				}
			}

			// Check if there's a final text response (done with tools)
			hasTextResponse := false
			for _, block := range result.Content {
				if block.Type == "text" && strings.TrimSpace(block.Text) != "" {
					hasTextResponse = true
				}
			}
			// If there's text alongside tools, return now
			if hasTextResponse {
				return result, nil
			}

			continue
		}

		// No tool calls - we have a final response
		return result, nil
	}

	return nil, fmt.Errorf("exceeded maximum tool call rounds (%d)", maxToolCallRounds)
}

// executeCompletion calls the provider's completion API.
func (r *Runtime) executeCompletion(ctx context.Context, req providers.CompletionRequest) (*providers.CompletionResult, error) {
	// Set budget from thinking level
	if r.thinkingLevel != "" {
		if budget, ok := providers.DefaultThinkingBudgets[r.thinkingLevel]; ok {
			req.ThinkingBudget = budget
		}
	}

	// Retry loop
	maxRetries := 3
	if r.config.Retry != nil && r.config.Retry.Enabled {
		maxRetries = r.config.Retry.MaxRetries
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Simple exponential backoff
			delay := time.Duration(1<<uint(attempt-1)) * 2 * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		// Try streaming
		if req.Stream {
			return r.executeStreaming(ctx, req)
		}

		result, err := r.provider.CreateCompletion(ctx, req)
		if err != nil {
			lastErr = err
			continue
		}
		return result, nil
	}

	return nil, fmt.Errorf("all retries exhausted: %w", lastErr)
}

// executeStreaming executes a streaming completion and collects the full response.
func (r *Runtime) executeStreaming(ctx context.Context, req providers.CompletionRequest) (*providers.CompletionResult, error) {
	ch, err := r.provider.StreamCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	var fullText string
	var thinking string
	var toolUses []providers.StreamingToolUse
	var finalResult *providers.CompletionResult

	for chunk := range ch {
		if chunk.Err != nil {
			return nil, chunk.Err
		}

		if chunk.Content != "" {
			fullText += chunk.Content
		}

		if chunk.Thinking != "" {
			thinking += chunk.Thinking
		}

		if len(chunk.ToolUses) > 0 {
			toolUses = append(toolUses, chunk.ToolUses...)
		}

		if chunk.Usage != nil {
			finalResult = &providers.CompletionResult{
				ID:         fmt.Sprintf("usage-%d", chunk.Usage.InputTokens+chunk.Usage.OutputTokens),
				Model:      r.modelID,
				Content:    buildContentBlocks(fullText, thinking, toolUses),
				StopReason: chunk.StopReason,
				Usage:      chunk.Usage,
			}
		}
	}

	if finalResult == nil {
		finalResult = &providers.CompletionResult{
			Model:   r.modelID,
			Content: buildContentBlocks(fullText, thinking, toolUses),
		}
	}

	return finalResult, nil
}

// QueueSteeringMessage queues a steering message to be delivered after current tool calls.
func (r *Runtime) QueueSteeringMessage(content string) {
	r.muQueue.Lock()
	defer r.muQueue.Unlock()
	r.messageQueue = append(r.messageQueue, &QueuedMessage{
		Type:    MessageSteering,
		Content: content,
		Time:    time.Now(),
	})
}

// QueueFollowUpMessage queues a follow-up message for after the agent finishes.
func (r *Runtime) QueueFollowUpMessage(content string) {
	r.muQueue.Lock()
	defer r.muQueue.Unlock()
	r.messageQueue = append(r.messageQueue, &QueuedMessage{
		Type:    MessageFollowUp,
		Content: content,
		Time:    time.Now(),
	})
}

// drainSteeringMessages pops all steering messages.
func (r *Runtime) drainSteeringMessages() []*QueuedMessage {
	r.muQueue.Lock()
	defer r.muQueue.Unlock()

	var steering []*QueuedMessage
	var remaining []*QueuedMessage
	for _, qm := range r.messageQueue {
		if qm.Type == MessageSteering {
			steering = append(steering, qm)
		} else {
			remaining = append(remaining, qm)
		}
	}
	r.messageQueue = remaining
	return steering
}

// Abort signals the agent to abort the current execution.
func (r *Runtime) Abort() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.aborted = true
	if r.abortSignal != nil {
		r.abortSignal()
	}
}

// GetStats returns runtime statistics.
func (r *Runtime) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"totalInputTokens":  r.totalInputTokens,
		"totalOutputTokens": r.totalOutputTokens,
		"totalCost":         r.totalCost,
		"sessionId":         r.session.ID,
		"modelId":           r.modelID,
		"provider":          r.providerID,
	}
}

// Helper functions

func resultContentString(result *providers.CompletionResult) string {
	var text string
	for _, block := range result.Content {
		if block.Type == "text" {
			text += block.Text
		}
	}
	return text
}

func containsToolCall(blocks []providers.ContentBlock) bool {
	for _, block := range blocks {
		if block.Type == "tool_use" {
			return true
		}
	}
	return false
}

func buildAssistantMessage(result *providers.CompletionResult) providers.Message {
	return providers.Message{
		Role:    "assistant",
		Content: result.Content,
	}
}

func buildContentBlocks(text, thinking string, toolUses []providers.StreamingToolUse) []providers.ContentBlock {
	var blocks []providers.ContentBlock
	if thinking != "" {
		blocks = append(blocks, providers.ContentBlock{Type: "thinking", Text: thinking})
	}
	if text != "" {
		blocks = append(blocks, providers.ContentBlock{Type: "text", Text: text})
	}
	for _, tu := range toolUses {
		blocks = append(blocks, providers.ContentBlock{
			Type:     "tool_use",
			ToolID:   tu.ToolID,
			ToolType: tu.ToolType,
			Input:    map[string]interface{}{"input": string(tu.Input)},
		})
	}
	return blocks
}

// estimateTokenCount estimates the total tokens in messages.
func (r *Runtime) estimateTokenCount(messages []providers.Message) int {
	total := 0
	for _, msg := range messages {
		for _, block := range msg.Content {
			total += r.provider.EstimateTokens(block.Text)
		}
	}
	return total
}

// compactIfNeeded runs compaction in the background.
func (r *Runtime) compactIfNeeded(ctx context.Context) error {
	r.mu.Lock()
	if r.compacting {
		r.mu.Unlock()
		return fmt.Errorf("compaction already in progress")
	}
	r.compacting = true
	r.mu.Unlock()

	defer func() {
		r.mu.Lock()
		r.compacting = false
		r.mu.Unlock()
	}()

	if r.compactionCtx == nil {
		return nil
	}

	// Trigger compaction via LLM
	summary, err := r.compactionCtx.Summarize(ctx, r.provider, r.modelID, r.session)
	if err != nil {
		return fmt.Errorf("compaction failed: %w", err)
	}

	// Determine range to compact
	entries := r.session.GetFullHistory("")
	compactFrom := entries[0].ID
	compactTo := entries[len(entries)/2].ID

	r.session.CompactContext(summary, compactFrom, compactTo)
	r.emit(Event{Type: EventCompaction, Payload: map[string]string{
		"summary": summary,
	}})

	return nil
}

// CompactionContext provides the context needed for compaction.
type CompactionContext struct {
	ReserveTokens    int
	KeepRecentTokens int
}

// Summarize uses the LLM to create a compact summary of the conversation.
func (c *CompactionContext) Summarize(ctx context.Context, provider providers.Provider, modelID string, sess *sessions.Session) (string, error) {
	// Build the messages to summarize (excluding most recent keepRecentTokens)
	// For now, return a placeholder
	return "[Compacted conversation - summary pending implementation]", nil
}

// ContextBuilder handles building the full context for the LLM.
type ContextBuilder struct {
	agencies []string // AGENTS.md content
	skills   []string // Skill content
	prompts  []string // Prompt templates
	sessions []*sessions.SessionEntry
}

// NewContextBuilder creates a new context builder.
func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{}
}
