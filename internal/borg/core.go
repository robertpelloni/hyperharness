package borg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// CoreStatus represents the status of the Borg core engine.
type CoreStatus struct {
	Active       bool                   `json:"active"`
	Version      string                 `json:"version"`
	Uptime       time.Duration          `json:"uptime"`
	ToolsLoaded  int                    `json:"toolsLoaded"`
	Servers      int                    `json:"servers"`
	MemoryUsage  int64                  `json:"memoryUsage"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// Core is the central Borg engine that coordinates all subsystems.
type Core struct {
	status    CoreStatus
	adapters  map[string]Adapter
	hooks     map[string][]HookFunc
	startTime time.Time
	config    *Config
	mu        sync.RWMutex
}

// Adapter is the interface for Borg subsystem adapters.
type Adapter interface {
	Name() string
	Start() error
	Stop() error
	Status() map[string]interface{}
}

// HookFunc is a callback for event hooks.
type HookFunc func(event string, data interface{}) error

// Config holds Borg core configuration.
type Config struct {
	DataDir       string            `json:"dataDir"`
	PluginDir     string            `json:"pluginDir"`
	AutoStart     bool              `json:"autoStart"`
	MaxMemory     int64             `json:"maxMemory"`
	Properties    map[string]string `json:"properties,omitempty"`
}

// DefaultConfig returns a default configuration.
func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		DataDir:   filepath.Join(home, ".hyperharness", "borg"),
		PluginDir: filepath.Join(home, ".hyperharness", "plugins"),
		AutoStart: true,
		MaxMemory: 512 * 1024 * 1024, // 512MB
		Properties: map[string]string{
			"engine":     "hyperharness-go",
			"version":    "0.1.0",
			"protocol":   "mcp-2024-11-05",
		},
	}
}

// NewCore creates a new Borg core engine.
func NewCore(config *Config) *Core {
	if config == nil {
		config = DefaultConfig()
	}

	return &Core{
		status: CoreStatus{
			Active:   false,
			Version:  "0.1.0",
			Metadata: make(map[string]interface{}),
		},
		adapters:  make(map[string]Adapter),
		hooks:     make(map[string][]HookFunc),
		startTime: time.Now(),
		config:    config,
	}
}

// Start initializes the Borg core and all registered adapters.
func (c *Core) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ensure data directory exists
	if err := os.MkdirAll(c.config.DataDir, 0o755); err != nil {
		return fmt.Errorf("failed to create data dir: %w", err)
	}

	// Start all adapters
	for name, adapter := range c.adapters {
		if err := adapter.Start(); err != nil {
			c.status.Metadata["error_"+name] = err.Error()
		}
	}

	c.status.Active = true
	c.startTime = time.Now()
	c.emit("core:start", nil)
	return nil
}

// Stop shuts down the Borg core and all adapters.
func (c *Core) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, adapter := range c.adapters {
		if err := adapter.Stop(); err != nil {
			c.status.Metadata["error_"+name] = err.Error()
		}
	}

	c.status.Active = false
	c.emit("core:stop", nil)
	return nil
}

// RegisterAdapter adds a subsystem adapter.
func (c *Core) RegisterAdapter(adapter Adapter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.adapters[adapter.Name()] = adapter
}

// GetAdapter retrieves an adapter by name.
func (c *Core) GetAdapter(name string) (Adapter, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	a, ok := c.adapters[name]
	return a, ok
}

// On registers an event hook.
func (c *Core) On(event string, hook HookFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hooks[event] = append(c.hooks[event], hook)
}

// Emit emits an event to all registered hooks.
func (c *Core) Emit(event string, data interface{}) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.emit(event, data)
}

// emit is the internal emit (must be called with read lock held).
func (c *Core) emit(event string, data interface{}) {
	hooks, ok := c.hooks[event]
	if !ok {
		return
	}
	for _, hook := range hooks {
		hook(event, data)
	}
}

// Status returns the current core status.
func (c *Core) Status() CoreStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()

	status := c.status
	status.Uptime = time.Since(c.startTime)
	status.Servers = len(c.adapters)

	// Collect adapter statuses
	for name, adapter := range c.adapters {
		adapterStatus := adapter.Status()
		status.Metadata["adapter_"+name] = adapterStatus
	}

	return status
}

// StatusJSON returns the status as JSON.
func (c *Core) StatusJSON() string {
	status := c.Status()
	data, _ := json.MarshalIndent(status, "", "  ")
	return string(data)
}

// SetProperty sets a configuration property.
func (c *Core) SetProperty(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.Properties[key] = value
}

// GetProperty retrieves a configuration property.
func (c *Core) GetProperty(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config.Properties[key]
}

// IsRunning returns whether the core is active.
func (c *Core) IsRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status.Active
}
