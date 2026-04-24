// Package orchestrator provides Agent-to-Agent (A2A) communication.
// Ported from hypercode/go/internal/orchestration/a2a_broker.go
//
// WHAT: A2A message broker for inter-agent communication
// WHY: Enables agents to coordinate, negotiate, and share state
// HOW: In-process message routing with typed messages and subscriptions
package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// A2AMessage represents a message between agents.
type A2AMessage struct {
	ID        string                 `json:"id"`
	Timestamp int64                  `json:"timestamp"`
	Sender    string                 `json:"sender"`
	Recipient string                 `json:"recipient,omitempty"`
	Type      A2AMessageType         `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
}

// A2AMessageType defines the type of inter-agent message.
type A2AMessageType string

const (
	TaskRequest      A2AMessageType = "task_request"
	TaskResponse     A2AMessageType = "task_response"
	TaskNegotiation  A2AMessageType = "task_negotiation"
	StateUpdate      A2AMessageType = "state_update"
	ErrorReport      A2AMessageType = "error_report"
	HandshakeRequest A2AMessageType = "handshake_request"
	HandshakeAccept  A2AMessageType = "handshake_accept"
)

// A2AHandler processes incoming A2A messages.
type A2AHandler func(msg A2AMessage) (*A2AMessage, error)

// A2ABroker routes messages between agents in the same process.
type A2ABroker struct {
	handlers map[string]A2AHandler
	inbox    map[string]chan A2AMessage
	history  []A2AMessage
	mu       sync.RWMutex
}

// NewA2ABroker creates a new A2A message broker.
func NewA2ABroker() *A2ABroker {
	return &A2ABroker{
		handlers: make(map[string]A2AHandler),
		inbox:    make(map[string]chan A2AMessage),
	}
}

// Register registers an agent with the broker.
func (b *A2ABroker) Register(agentID string, handler A2AHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[agentID] = handler
	b.inbox[agentID] = make(chan A2AMessage, 100)
}

// Unregister removes an agent from the broker.
func (b *A2ABroker) Unregister(agentID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.handlers, agentID)
	if ch, ok := b.inbox[agentID]; ok {
		close(ch)
		delete(b.inbox, agentID)
	}
}

// RouteMessage delivers a message to the recipient's inbox.
func (b *A2ABroker) RouteMessage(msg A2AMessage) {
	b.mu.Lock()
	b.history = append(b.history, msg)
	if len(b.history) > 1000 {
		b.history = b.history[len(b.history)-1000:]
	}
	b.mu.Unlock()

	if msg.Recipient != "" {
		b.mu.RLock()
		inbox, ok := b.inbox[msg.Recipient]
		b.mu.RUnlock()
		if ok {
			select {
			case inbox <- msg:
			default:
				// Inbox full, drop message
			}
		}
	}

	// Also deliver to handler if registered
	if msg.Recipient != "" {
		b.mu.RLock()
		handler, ok := b.handlers[msg.Recipient]
		b.mu.RUnlock()
		if ok && handler != nil {
			go handler(msg)
		}
	}
}

// Broadcast sends a message to all registered agents.
func (b *A2ABroker) Broadcast(msg A2AMessage) {
	b.mu.RLock()
	agents := make([]string, 0, len(b.inbox))
	for id := range b.inbox {
		agents = append(agents, id)
	}
	b.mu.RUnlock()

	for _, agentID := range agents {
		msgCopy := msg
		msgCopy.Recipient = agentID
		b.RouteMessage(msgCopy)
	}
}

// Query sends a request and waits for a response.
func (b *A2ABroker) Query(ctx context.Context, msg A2AMessage) (*A2AMessage, error) {
	// Try direct handler first
	b.mu.RLock()
	handler, ok := b.handlers[msg.Recipient]
	b.mu.RUnlock()

	if ok && handler != nil {
		resp, err := handler(msg)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	// Fallback: broadcast and wait
	b.RouteMessage(msg)
	return nil, fmt.Errorf("no handler for agent %q", msg.Recipient)
}

// History returns recent messages.
func (b *A2ABroker) History(limit int) []A2AMessage {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if limit <= 0 || limit > len(b.history) {
		limit = len(b.history)
	}
	result := make([]A2AMessage, limit)
	copy(result, b.history[len(b.history)-limit:])
	return result
}

// AgentCount returns the number of registered agents.
func (b *A2ABroker) AgentCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers)
}

// NewA2AMessage creates a new message with auto-generated ID and timestamp.
func NewA2AMessage(sender string, msgType A2AMessageType, payload map[string]interface{}) A2AMessage {
	return A2AMessage{
		ID:        fmt.Sprintf("msg-%d", time.Now().UnixNano()),
		Timestamp: time.Now().UTC().UnixMilli(),
		Sender:    sender,
		Type:      msgType,
		Payload:   payload,
	}
}
