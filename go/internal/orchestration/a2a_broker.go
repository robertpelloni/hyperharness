package orchestration

/**
 * @file a2a_broker.go
 * @module go/internal/orchestration
 *
 * WHAT: Go-native implementation of the Agent-to-Agent (A2A) message broker.
 *
 * WHY: Total Autonomy — The Go sidecar must be capable of routing messages 
 * between autonomous agents without relying on the Node control plane.
 */

import (
	"fmt"
	"sync"
)

type A2AMessageType string

const (
	TaskRequest     A2AMessageType = "TASK_REQUEST"
	TaskResponse    A2AMessageType = "TASK_RESPONSE"
	ConsensusVote   A2AMessageType = "CONSENSUS_VOTE"
	StateUpdate     A2AMessageType = "STATE_UPDATE"
	Handoff         A2AMessageType = "HANDOFF"
	DebateProposal  A2AMessageType = "DEBATE_PROPOSAL"
	Critique        A2AMessageType = "CRITIQUE"
)

type A2AMessage struct {
	ID        string         `json:"id"`
	Timestamp int64          `json:"timestamp"`
	Sender    string         `json:"sender"`
	Recipient string         `json:"recipient,omitempty"`
	Type      A2AMessageType `json:"type"`
	Payload   interface{}    `json:"payload"`
	ReplyTo   string         `json:"replyTo,omitempty"`
}

type A2ABroker struct {
	mu      sync.RWMutex
	agents  map[string]chan A2AMessage
	history []A2AMessage
}

func NewA2ABroker() *A2ABroker {
	return &A2ABroker{
		agents:  make(map[string]chan A2AMessage),
		history: make([]A2AMessage, 0),
	}
}

func (b *A2ABroker) RegisterAgent(id string) chan A2AMessage {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan A2AMessage, 100)
	b.agents[id] = ch
	fmt.Printf("[Go A2A] Registered agent: %s\n", id)
	return ch
}

func (b *A2ABroker) UnregisterAgent(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if ch, ok := b.agents[id]; ok {
		close(ch)
		delete(b.agents, id)
		fmt.Printf("[Go A2A] Unregistered agent: %s\n", id)
	}
}

func (b *A2ABroker) RouteMessage(msg A2AMessage) {
	b.mu.Lock()
	b.history = append(b.history, msg)
	if len(b.history) > 1000 {
		b.history = b.history[1:]
	}
	b.mu.Unlock()

	b.mu.RLock()
	defer b.mu.RUnlock()

	if msg.Recipient != "" {
		if ch, ok := b.agents[msg.Recipient]; ok {
			select {
			case ch <- msg:
			default:
				fmt.Printf("[Go A2A] Dropped message to %s (buffer full)\n", msg.Recipient)
			}
		}
	} else {
		// Broadcast
		for id, ch := range b.agents {
			if id == msg.Sender {
				continue
			}
			select {
			case ch <- msg:
			default:
				fmt.Printf("[Go A2A] Dropped broadcast to %s (buffer full)\n", id)
			}
		}
	}
}

func (b *A2ABroker) GetHistory() []A2AMessage {
	b.mu.RLock()
	defer b.mu.RUnlock()

	h := make([]A2AMessage, len(b.history))
	copy(h, b.history)
	return h
}

func (b *A2ABroker) ListAgents() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	agents := make([]string, 0, len(b.agents))
	for id := range b.agents {
		agents = append(agents, id)
	}
	return agents
}
