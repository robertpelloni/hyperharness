package borg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ═══════════════════════════════════════════════════════════════════════
// core.go — Borg Core Engine for HyperHarness
// The central coordination layer that unifies Memory, Context Management,
// MCP routing, and Agent orchestration across all harnesses.
//
// WHAT: The Borg Core is the central nervous system of HyperHarness —
//       it coordinates all subsystems and provides a unified API surface.
// WHY: Without a central coordinator, each subsystem (memory, context,
//       MCP, agents) operates independently, leading to inconsistency.
// HOW: Core provides lifecycle management, event bus, and unified config.
// ═══════════════════════════════════════════════════════════════════════

// State represents the current state of the Borg Core.
type State string

const (
	StateUninitialized State = "uninitialized"
	StateBooting       State = "booting"
	StateReady         State = "ready"
	StateBusy          State = "busy"
	StateShuttingDown  State = "shutting_down"
	StateError         State = "error"
)

// Event represents a system event on the event bus.
type Event struct {
	Type      string         `json:"type"`
	Source    string         `json:"source"`
	Data      map[string]any `json:"data,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

// EventListener receives events from the bus.
type EventListener func(event Event)

// Subsystem represents a pluggable component of the Borg Core.
type Subsystem interface {
	Name() string
	Initialize(core *Core) error
	Shutdown() error
	Status() map[string]any
}

// Core is the central coordination engine.
type Core struct {
	mu         sync.RWMutex
	state      State
	workingDir string
	homeDir    string

	// Subsystems
	subsystems map[string]Subsystem

	// Event bus
	listeners map[string][]EventListener

	// Metadata
	startTime time.Time
	stats     CoreStats
}

// CoreStats tracks core operation metrics.
type CoreStats struct {
	EventsPublished int64     `json:"eventsPublished"`
	Subsystems      int       `json:"subsystems"`
	Uptime          string    `json:"uptime"`
	LastError       string    `json:"lastError,omitempty"`
	LastEventAt     time.Time `json:"lastEventAt,omitempty"`
}

// NewCore creates a new Borg Core instance.
func NewCore(workingDir string) (*Core, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot determine home directory: %w", err)
	}

	// Ensure .hyperharness directory exists
	borgDir := filepath.Join(home, ".hyperharness")
	os.MkdirAll(borgDir, 0755)

	core := &Core{
		state:      StateUninitialized,
		workingDir: workingDir,
		homeDir:    home,
		subsystems: make(map[string]Subsystem),
		listeners:  make(map[string][]EventListener),
	}

	return core, nil
}

// Boot initializes the core and all registered subsystems.
func (c *Core) Boot(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.state = StateBooting
	c.startTime = time.Now()

	c.publishUnlocked(Event{
		Type:      "core:boot_start",
		Source:    "core",
		Timestamp: time.Now(),
		Data:      map[string]any{"workingDir": c.workingDir},
	})

	// Initialize each subsystem in order
	for name, subsystem := range c.subsystems {
		if err := subsystem.Initialize(c); err != nil {
			c.state = StateError
			c.stats.LastError = fmt.Sprintf("subsystem %s failed: %v", name, err)
			c.publishUnlocked(Event{
				Type:      "core:subsystem_error",
				Source:    "core",
				Timestamp: time.Now(),
				Data:      map[string]any{"subsystem": name, "error": err.Error()},
			})
			return fmt.Errorf("subsystem %s initialization failed: %w", name, err)
		}

		c.publishUnlocked(Event{
			Type:      "core:subsystem_ready",
			Source:    "core",
			Timestamp: time.Now(),
			Data:      map[string]any{"subsystem": name},
		})
	}

	c.state = StateReady
	c.stats.Subsystems = len(c.subsystems)

	c.publishUnlocked(Event{
		Type:      "core:boot_complete",
		Source:    "core",
		Timestamp: time.Now(),
		Data:      map[string]any{"subsystems": c.stats.Subsystems},
	})

	return nil
}

// Shutdown gracefully stops all subsystems.
func (c *Core) Shutdown() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.state = StateShuttingDown

	c.publishUnlocked(Event{
		Type:      "core:shutdown_start",
		Source:    "core",
		Timestamp: time.Now(),
	})

	for name, subsystem := range c.subsystems {
		if err := subsystem.Shutdown(); err != nil {
			c.publishUnlocked(Event{
				Type:      "core:subsystem_shutdown_error",
				Source:    "core",
				Timestamp: time.Now(),
				Data:      map[string]any{"subsystem": name, "error": err.Error()},
			})
		}
	}

	c.state = StateUninitialized

	c.publishUnlocked(Event{
		Type:      "core:shutdown_complete",
		Source:    "core",
		Timestamp: time.Now(),
	})

	return nil
}

// RegisterSubsystem adds a subsystem to the core.
func (c *Core) RegisterSubsystem(subsystem Subsystem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.subsystems[subsystem.Name()] = subsystem
}

// GetSubsystem retrieves a subsystem by name.
func (c *Core) GetSubsystem(name string) (Subsystem, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.subsystems[name]
	return s, ok
}

// State returns the current core state.
func (c *Core) State() State {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state
}

// Stats returns current core statistics.
func (c *Core) Stats() CoreStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats := c.stats
	stats.Uptime = time.Since(c.startTime).Round(time.Second).String()
	return stats
}

// WorkingDir returns the configured working directory.
func (c *Core) WorkingDir() string {
	return c.workingDir
}

// HomeDir returns the user's home directory.
func (c *Core) HomeDir() string {
	return c.homeDir
}

// BorgDir returns the .hyperharness directory path.
func (c *Core) BorgDir() string {
	return filepath.Join(c.homeDir, ".hyperharness")
}

// ── Event Bus ──

// Subscribe registers an event listener for a specific event type.
// Use "*" to listen for all events.
func (c *Core) Subscribe(eventType string, listener EventListener) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.listeners[eventType] = append(c.listeners[eventType], listener)
}

// Publish sends an event to all matching listeners.
func (c *Core) Publish(event Event) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.publishUnlocked(event)
}

func (c *Core) publishUnlocked(event Event) {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	c.stats.EventsPublished++
	c.stats.LastEventAt = event.Timestamp

	// Notify specific listeners
	if listeners, ok := c.listeners[event.Type]; ok {
		for _, l := range listeners {
			l(event)
		}
	}

	// Notify wildcard listeners
	if listeners, ok := c.listeners["*"]; ok {
		for _, l := range listeners {
			l(event)
		}
	}
}

// ── Health Check ──

// HealthCheck runs diagnostics on all subsystems.
func (c *Core) HealthCheck() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := map[string]any{
		"state":       string(c.state),
		"uptime":      time.Since(c.startTime).Round(time.Second).String(),
		"workingDir":  c.workingDir,
		"subsystems":  len(c.subsystems),
		"eventsTotal": c.stats.EventsPublished,
	}

	for name, subsystem := range c.subsystems {
		result["subsystem:"+name] = subsystem.Status()
	}

	return result
}

// ── Memory Management ──

// MemoryStore provides persistent key-value storage backed by the borg dir.
type MemoryStore struct {
	mu   sync.RWMutex
	dir  string
	data map[string]string
}

// NewMemoryStore creates a memory store in the borg directory.
func NewMemoryStore(borgDir string) *MemoryStore {
	memDir := filepath.Join(borgDir, "memory")
	os.MkdirAll(memDir, 0755)

	return &MemoryStore{
		dir:  memDir,
		data: make(map[string]string),
	}
}

// Set stores a key-value pair.
func (m *MemoryStore) Set(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
	return os.WriteFile(filepath.Join(m.dir, key+".md"), []byte(value), 0644)
}

// Get retrieves a value by key.
func (m *MemoryStore) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if val, ok := m.data[key]; ok {
		return val, true
	}

	// Try reading from disk
	data, err := os.ReadFile(filepath.Join(m.dir, key+".md"))
	if err != nil {
		return "", false
	}
	return string(data), true
}

// Delete removes a key.
func (m *MemoryStore) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	return os.Remove(filepath.Join(m.dir, key+".md"))
}

// List returns all keys.
func (m *MemoryStore) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}

	// Also check disk
	entries, _ := os.ReadDir(m.dir)
	for _, e := range entries {
		name := e.Name()
		if strings.HasSuffix(name, ".md") {
			key := strings.TrimSuffix(name, ".md")
			found := false
			for _, existing := range keys {
				if existing == key {
					found = true
					break
				}
			}
			if !found {
				keys = append(keys, key)
			}
		}
	}

	return keys
}

// Load reads all memory from disk into the in-memory cache.
func (m *MemoryStore) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		key := strings.TrimSuffix(e.Name(), ".md")
		data, err := os.ReadFile(filepath.Join(m.dir, e.Name()))
		if err != nil {
			continue
		}
		m.data[key] = string(data)
	}

	return nil
}

// strings import needed for List - already imported above
