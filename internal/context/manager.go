package context

import (
	"sort"
	"sync"
)

type Message struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	Name       string `json:"name,omitempty"`
	ToolCallID string `json:"tool_call_id,omitempty"`
	Tokens     int    `json:"tokens"`
}

type Status struct {
	TotalMessages    int               `json:"totalMessages"`
	MaxMessages      int               `json:"maxMessages"`
	TotalTokens      int               `json:"totalTokens"`
	MaxTokens        int               `json:"maxTokens"`
	UtilizationPct   float64           `json:"utilizationPct"`
	InjectedContext  int               `json:"injectedContext"`
	ByRole           map[string]int    `json:"byRole"`
}

type Manager struct {
	mu        sync.RWMutex
	messages  []Message
	maxItems  int
	maxTokens int
}

func NewManager(maxItems, maxTokens int) *Manager {
	if maxItems <= 0 {
		maxItems = 100
	}
	if maxTokens <= 0 {
		maxTokens = 200000
	}
	return &Manager{
		messages:  make([]Message, 0),
		maxItems:  maxItems,
		maxTokens: maxTokens,
	}
}

func (m *Manager) Add(msg Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, msg)
	m.trimLocked()
}

func (m *Manager) AddBatch(msgs []Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, msgs...)
	m.trimLocked()
}

func (m *Manager) GetAll() []Message {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]Message, len(m.messages))
	copy(result, m.messages)
	return result
}

func (m *Manager) GetRecent(n int) []Message {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if n <= 0 || n >= len(m.messages) {
		return m.GetAll()
	}
	result := make([]Message, n)
	copy(result, m.messages[len(m.messages)-n:])
	return result
}

func (m *Manager) TotalTokens() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	total := 0
	for _, msg := range m.messages {
		total += msg.Tokens
	}
	return total
}

func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = make([]Message, 0)
}

func (m *Manager) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.messages)
}

func (m *Manager) Compact(maxMsgs int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.messages) <= maxMsgs {
		return false
	}
	m.messages = m.messages[len(m.messages)-maxMsgs:]
	return true
}

func (m *Manager) Inject(role, content string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	msg := Message{Role: role, Content: content, Tokens: len(content) / 4}
	m.messages = append([]Message{msg}, m.messages...)
	return true
}

func (m *Manager) Status() *Status {
	m.mu.RLock()
	defer m.mu.RUnlock()
	total := 0
	byRole := make(map[string]int)
	for _, msg := range m.messages {
		total += msg.Tokens
		byRole[msg.Role]++
	}
	utilPct := 0.0
	if m.maxTokens > 0 {
		utilPct = float64(total) / float64(m.maxTokens) * 100
	}
	injected := 0
	for _, msg := range m.messages {
		if msg.Role == "system" || msg.Role == "injected" {
			injected++
		}
	}
	return &Status{
		TotalMessages:   len(m.messages),
		MaxMessages:     m.maxItems,
		TotalTokens:     total,
		MaxTokens:       m.maxTokens,
		UtilizationPct:  utilPct,
		InjectedContext: injected,
		ByRole:          byRole,
	}
}

func (m *Manager) trimLocked() {
	if len(m.messages) > m.maxItems {
		m.messages = m.messages[len(m.messages)-m.maxItems:]
	}
	total := 0
	for _, msg := range m.messages {
		total += msg.Tokens
	}
	if total > m.maxTokens {
		var system, rest []Message
		for _, msg := range m.messages {
			if msg.Role == "system" {
				system = append(system, msg)
			} else {
				rest = append(rest, msg)
			}
		}
		trimmed := make([]Message, 0, len(m.messages))
		trimmed = append(trimmed, system...)
		tokens := 0
		for i := len(rest) - 1; i >= 0; i-- {
			if tokens+rest[i].Tokens > m.maxTokens {
				break
			}
			tokens += rest[i].Tokens
			trimmed = append([]Message{rest[i]}, trimmed...)
		}
		m.messages = trimmed
	}
}

func (m *Manager) Search(query string) []Message {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []Message
	for _, msg := range m.messages {
		if contains(msg.Content, query) {
			result = append(result, msg)
		}
	}
	return result
}

func (m *Manager) SetMaxTokens(maxTokens int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.maxTokens = maxTokens
	m.trimLocked()
}

func (m *Manager) Summary() map[string]interface{} {
	status := m.Status()
	return map[string]interface{}{
		"messageCount":   status.TotalMessages,
		"totalTokens":    status.TotalTokens,
		"maxTokens":      status.MaxTokens,
		"byRole":         status.ByRole,
		"usagePercent":    status.UtilizationPct,
	}
}

func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

var _ = sort.Ints
