package context

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Manager handles conversation context lifecycle: compaction, injection, and status.
// This is the Go-native equivalent of the context management systems found in
// Claude Code, Cursor, and other harness tools.
type Manager struct {
	messages    []ContextMessage
	maxMessages int
	maxTokens   int
	systemPrompt string
	injectedContext []string
	mu          sync.RWMutex
}

// ContextMessage represents a single message in the conversation context.
type ContextMessage struct {
	Role       string    `json:"role"`
	Content    string    `json:"content"`
	Name       string    `json:"name,omitempty"`
	ToolCallID string    `json:"tool_call_id,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
	TokenCount int       `json:"token_count,omitempty"`
}

// ContextStatus represents the current state of the context manager.
type ContextStatus struct {
	TotalMessages   int      `json:"total_messages"`
	TotalTokens     int      `json:"total_tokens"`
	MaxMessages     int      `json:"max_messages"`
	MaxTokens       int      `json:"max_tokens"`
	InjectedContext []string `json:"injected_context"`
	UtilizationPct  float64  `json:"utilization_pct"`
}

// NewManager creates a new context manager with sensible defaults.
func NewManager(maxMessages, maxTokens int) *Manager {
	if maxMessages <= 0 {
		maxMessages = 100
	}
	if maxTokens <= 0 {
		maxTokens = 200000
	}
	return &Manager{
		messages:    make([]ContextMessage, 0),
		maxMessages: maxMessages,
		maxTokens:   maxTokens,
		injectedContext: make([]string, 0),
	}
}

// AddMessage appends a message to the context.
func (m *Manager) AddMessage(role, content string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	msg := ContextMessage{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
		TokenCount: estimateTokens(content),
	}
	m.messages = append(m.messages, msg)
}

// AddToolMessage adds a tool result message.
func (m *Manager) AddToolMessage(toolCallID, name, content string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = append(m.messages, ContextMessage{
		Role:       "tool",
		Content:    content,
		Name:       name,
		ToolCallID: toolCallID,
		Timestamp:  time.Now(),
		TokenCount: estimateTokens(content),
	})
}

// GetMessages returns all messages in the context.
func (m *Manager) GetMessages() []ContextMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]ContextMessage, len(m.messages))
	copy(result, m.messages)
	return result
}

// GetMessagesForProvider returns messages formatted for LLM provider consumption.
// Includes system prompt and injected context at the beginning.
func (m *Manager) GetMessagesForProvider() []ContextMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []ContextMessage

	// System prompt
	if m.systemPrompt != "" {
		result = append(result, ContextMessage{
			Role:    "system",
			Content: m.systemPrompt,
		})
	}

	// Injected context
	if len(m.injectedContext) > 0 {
		var sb strings.Builder
		sb.WriteString("## Additional Context\n\n")
		for _, ctx := range m.injectedContext {
			sb.WriteString(ctx)
			sb.WriteString("\n\n")
		}
		result = append(result, ContextMessage{
			Role:    "system",
			Content: sb.String(),
		})
	}

	// Conversation messages
	result = append(result, m.messages...)
	return result
}

// Compact trims the conversation history to fit within limits.
// Preserves the system prompt and the most recent messages.
// Returns the number of messages removed.
func (m *Manager) Compact(maxMessages int) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	if maxMessages <= 0 {
		maxMessages = m.maxMessages
	}

	if len(m.messages) <= maxMessages {
		return 0
	}

	// Keep the most recent messages
	removed := len(m.messages) - maxMessages
	m.messages = m.messages[removed:]
	return removed
}

// CompactToFit compacts messages to fit within a token budget.
func (m *Manager) CompactToFit(maxTokens int) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	if maxTokens <= 0 {
		maxTokens = m.maxTokens
	}

	totalTokens := 0
	for _, msg := range m.messages {
		totalTokens += msg.TokenCount
	}

	if totalTokens <= maxTokens {
		return 0
	}

	// Remove oldest messages until we fit
	removed := 0
	for len(m.messages) > 1 && totalTokens > maxTokens {
		removed++
		totalTokens -= m.messages[0].TokenCount
		m.messages = m.messages[1:]
	}

	return removed
}

// Inject adds contextual information that gets prepended to the conversation.
func (m *Manager) Inject(context string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.injectedContext = append(m.injectedContext, context)
}

// SetSystemPrompt sets the system prompt.
func (m *Manager) SetSystemPrompt(prompt string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.systemPrompt = prompt
}

// Status returns the current context status.
func (m *Manager) Status() ContextStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	totalTokens := 0
	for _, msg := range m.messages {
		totalTokens += msg.TokenCount
	}

	utilization := 0.0
	if m.maxTokens > 0 {
		utilization = float64(totalTokens) / float64(m.maxTokens) * 100
	}

	return ContextStatus{
		TotalMessages:   len(m.messages),
		TotalTokens:     totalTokens,
		MaxMessages:     m.maxMessages,
		MaxTokens:       m.maxTokens,
		InjectedContext: m.injectedContext,
		UtilizationPct:  utilization,
	}
}

// Clear removes all messages and injected context.
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = make([]ContextMessage, 0)
	m.injectedContext = make([]string, 0)
}

// Summarize returns a compact summary of the conversation for context compaction.
// This implements the "context grooming" pattern from Maestro/hypercode.
func (m *Manager) Summarize() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.messages) == 0 {
		return "(empty context)"
	}

	var sb strings.Builder
	userCount := 0
	assistantCount := 0
	toolCount := 0

	for _, msg := range m.messages {
		switch msg.Role {
		case "user":
			userCount++
		case "assistant":
			assistantCount++
		case "tool":
			toolCount++
		}
	}

	sb.WriteString("Context summary:\n")
	sb.WriteString("- User messages: ")
	sb.WriteString(strings.Repeat("|", min(userCount, 50)))
	sb.WriteString(fmt.Sprintf(" (%d)\n", userCount))
	sb.WriteString("- Assistant messages: ")
	sb.WriteString(strings.Repeat("|", min(assistantCount, 50)))
	sb.WriteString(fmt.Sprintf(" (%d)\n", assistantCount))
	sb.WriteString("- Tool calls: ")
	sb.WriteString(strings.Repeat("|", min(toolCount, 50)))
	sb.WriteString(fmt.Sprintf(" (%d)\n", toolCount))

	// Include last assistant message as context
	for i := len(m.messages) - 1; i >= 0; i-- {
		if m.messages[i].Role == "assistant" {
			content := m.messages[i].Content
			if len(content) > 500 {
				content = content[:500] + "..."
			}
			sb.WriteString("\nLast assistant message:\n")
			sb.WriteString(content)
			break
		}
	}

	return sb.String()
}

// estimateTokens provides a rough token count estimate.
// ~4 characters per token is a reasonable approximation for English text.
func estimateTokens(text string) int {
	return len(text) / 4
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
