// Package eventbus provides a typed, thread-safe, pattern-matching event bus.
// Ported from hypercode/go/internal/eventbus/eventbus.go
//
// WHAT: Pub/sub event system with exact match, wildcard patterns, and bounded history
// WHY: Decoupled communication between all subsystems (tools, memory, agents, MCP)
// HOW: Listener registration with glob patterns, safe concurrent emit, bounded history buffer
package eventbus

import (
	"regexp"
	"strings"
	"sync"
	"time"
)

// SystemEventType enumerates well-known event categories.
type SystemEventType string

const (
	EventAgentHeartbeat SystemEventType = "agent:heartbeat"
	EventAgentStart     SystemEventType = "agent:start"
	EventAgentStop      SystemEventType = "agent:stop"
	EventTaskUpdate     SystemEventType = "task:update"
	EventTaskComplete   SystemEventType = "task:complete"
	EventToolCall       SystemEventType = "tool:call"
	EventToolResult     SystemEventType = "tool:result"
	EventMemoryPrune    SystemEventType = "memory:prune"
	EventMemoryStore    SystemEventType = "memory:store"
	EventFileChange     SystemEventType = "file:change"
	EventTerminalError  SystemEventType = "terminal:error"
	EventMCPConnect     SystemEventType = "mcp:connect"
	EventMCPDisconnect  SystemEventType = "mcp:disconnect"
)

// SystemEvent is a single structured event emitted by any component.
type SystemEvent struct {
	Type      SystemEventType `json:"type"`
	Timestamp int64           `json:"timestamp"`
	Source    string          `json:"source"`
	Payload   interface{}     `json:"payload,omitempty"`
}

// wildcardListener binds a compiled glob pattern to a callback.
type wildcardListener struct {
	pattern  *regexp.Regexp
	listener func(SystemEvent)
}

// EventBus is a typed, thread-safe, pattern-matching event bus.
type EventBus struct {
	mu                sync.RWMutex
	exactListeners    map[string][]func(SystemEvent)
	wildcardListeners []wildcardListener
	globalListeners   []func(SystemEvent)
	history           []SystemEvent
	maxHistory        int
}

// New creates a new EventBus with the given maximum history size.
func New(maxHistory int) *EventBus {
	if maxHistory <= 0 {
		maxHistory = 1000
	}
	return &EventBus{
		exactListeners: make(map[string][]func(SystemEvent)),
		history:        make([]SystemEvent, 0, maxHistory),
		maxHistory:     maxHistory,
	}
}

// Subscribe registers a listener for events matching pattern.
// Patterns with '*' are treated as wildcards (e.g. "file:*" matches "file:change").
func (eb *EventBus) Subscribe(pattern string, listener func(SystemEvent)) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if containsGlob(pattern) {
		eb.wildcardListeners = append(eb.wildcardListeners, wildcardListener{
			pattern:  globToRegex(pattern),
			listener: listener,
		})
	} else {
		eb.exactListeners[pattern] = append(eb.exactListeners[pattern], listener)
	}
}

// OnGlobal registers a listener that fires for every event.
func (eb *EventBus) OnGlobal(listener func(SystemEvent)) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.globalListeners = append(eb.globalListeners, listener)
}

// EmitEvent publishes an event to all matching subscribers.
func (eb *EventBus) EmitEvent(eventType SystemEventType, source string, payload interface{}) {
	event := SystemEvent{
		Type:      eventType,
		Timestamp: time.Now().UnixMilli(),
		Source:    source,
		Payload:   payload,
	}

	eb.mu.Lock()
	eb.history = append(eb.history, event)
	if len(eb.history) > eb.maxHistory {
		eb.history = eb.history[len(eb.history)-eb.maxHistory:]
	}
	exact := eb.exactListeners[string(eventType)]
	wcs := append([]wildcardListener(nil), eb.wildcardListeners...)
	var globals []func(SystemEvent); globals = append(globals, eb.globalListeners...)
	eb.mu.Unlock()

	for _, fn := range globals {
		safeCall(fn, event)
	}
	for _, fn := range exact {
		safeCall(fn, event)
	}
	for _, wl := range wcs {
		if wl.pattern.MatchString(string(eventType)) {
			safeCall(wl.listener, event)
		}
	}
}

// GetHistory returns the last `limit` events from history.
func (eb *EventBus) GetHistory(limit int) []SystemEvent {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	if limit <= 0 || limit > len(eb.history) {
		limit = len(eb.history)
	}
	result := make([]SystemEvent, limit)
	copy(result, eb.history[len(eb.history)-limit:])
	return result
}

// ListenerCount returns total number of registered listeners.
func (eb *EventBus) ListenerCount() int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	count := len(eb.globalListeners) + len(eb.wildcardListeners)
	for _, v := range eb.exactListeners {
		count += len(v)
	}
	return count
}

func containsGlob(s string) bool {
	return strings.Contains(s, "*")
}

func globToRegex(glob string) *regexp.Regexp {
	escaped := regexp.QuoteMeta(glob)
	regex := "^" + strings.ReplaceAll(escaped, `\*`, ".*") + "$"
	re, _ := regexp.Compile(regex)
	if re == nil {
		re = regexp.MustCompile(".*")
	}
	return re
}

func safeCall(fn func(SystemEvent), event SystemEvent) {
	defer func() { recover() }()
	fn(event)
}
