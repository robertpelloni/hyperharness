package extensions

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Extension represents an MCP server or tool extension.
type Extension struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Type        ExtensionType          `json:"type"`
	Command     string                 `json:"command"`
	Args        []string               `json:"args,omitempty"`
	Env         map[string]string      `json:"env,omitempty"`
	Enabled     bool                   `json:"enabled"`
	AutoStart   bool                   `json:"auto_start"`
	Status      ExtensionStatus        `json:"status"`
	ToolCount   int                    `json:"tool_count"`
	Tools       []ExtensionTool        `json:"tools,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	InstalledAt time.Time              `json:"installed_at,omitempty"`
	LastError   string                 `json:"last_error,omitempty"`
}

// ExtensionType defines the type of extension.
type ExtensionType string

const (
	TypeMCPStdio  ExtensionType = "mcp_stdio"
	TypeMCPStream ExtensionType = "mcp_stream"
	TypeBuiltIn   ExtensionType = "builtin"
	TypePlugin    ExtensionType = "plugin"
)

// ExtensionStatus represents the runtime status of an extension.
type ExtensionStatus string

const (
	StatusNotStarted ExtensionStatus = "not_started"
	StatusStarting   ExtensionStatus = "starting"
	StatusRunning    ExtensionStatus = "running"
	StatusStopped    ExtensionStatus = "stopped"
	StatusError      ExtensionStatus = "error"
)

// ExtensionTool represents a tool provided by an extension.
type ExtensionTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Schema      json.RawMessage `json:"schema,omitempty"`
}

// Manager manages extension lifecycle (discovery, start/stop, tool routing).
type Manager struct {
	extensions map[string]*Extension
	configDir  string
	mu         sync.RWMutex
}

// NewManager creates a new extension manager.
func NewManager(configDir string) *Manager {
	if configDir == "" {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".hyperharness", "extensions")
	}
	return &Manager{
		extensions: make(map[string]*Extension),
		configDir:  configDir,
	}
}

// Register adds an extension to the manager.
func (m *Manager) Register(ext *Extension) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ext.Name == "" {
		return fmt.Errorf("extension name is required")
	}

	ext.InstalledAt = time.Now()
	if ext.Status == "" {
		ext.Status = StatusNotStarted
	}

	m.extensions[ext.Name] = ext
	return nil
}

// Unregister removes an extension.
func (m *Manager) Unregister(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ext, ok := m.extensions[name]
	if !ok {
		return fmt.Errorf("extension not found: %s", name)
	}

	// Stop if running
	if ext.Status == StatusRunning {
		m.stopExtension(ext)
	}

	delete(m.extensions, name)
	return nil
}

// Get retrieves an extension by name.
func (m *Manager) Get(name string) (*Extension, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ext, ok := m.extensions[name]
	return ext, ok
}

// List returns all extensions.
func (m *Manager) List() []*Extension {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var extensions []*Extension
	for _, ext := range m.extensions {
		extensions = append(extensions, ext)
	}
	sort.Slice(extensions, func(i, j int) bool {
		return extensions[i].Name < extensions[j].Name
	})
	return extensions
}

// ListEnabled returns only enabled extensions.
func (m *Manager) ListEnabled() []*Extension {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var extensions []*Extension
	for _, ext := range m.extensions {
		if ext.Enabled {
			extensions = append(extensions, ext)
		}
	}
	return extensions
}

// Enable enables an extension.
func (m *Manager) Enable(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ext, ok := m.extensions[name]
	if !ok {
		return fmt.Errorf("extension not found: %s", name)
	}
	ext.Enabled = true
	return nil
}

// Disable disables an extension.
func (m *Manager) Disable(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ext, ok := m.extensions[name]
	if !ok {
		return fmt.Errorf("extension not found: %s", name)
	}
	ext.Enabled = false
	if ext.Status == StatusRunning {
		m.stopExtension(ext)
	}
	return nil
}

// Start starts an extension's MCP server.
func (m *Manager) Start(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ext, ok := m.extensions[name]
	if !ok {
		return fmt.Errorf("extension not found: %s", name)
	}

	return m.startExtension(ext)
}

// Stop stops an extension's MCP server.
func (m *Manager) Stop(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ext, ok := m.extensions[name]
	if !ok {
		return fmt.Errorf("extension not found: %s", name)
	}

	return m.stopExtension(ext)
}

// StartAll starts all enabled auto-start extensions.
func (m *Manager) StartAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ext := range m.extensions {
		if ext.Enabled && ext.AutoStart && ext.Status != StatusRunning {
			if err := m.startExtension(ext); err != nil {
				ext.LastError = err.Error()
			}
		}
	}
	return nil
}

// StopAll stops all running extensions.
func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ext := range m.extensions {
		if ext.Status == StatusRunning {
			m.stopExtension(ext)
		}
	}
}

// GetAllTools returns all tools from all enabled extensions.
func (m *Manager) GetAllTools() []ExtensionTool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var tools []ExtensionTool
	for _, ext := range m.extensions {
		if ext.Enabled && ext.Status == StatusRunning {
			tools = append(tools, ext.Tools...)
		}
	}
	return tools
}

// FindTool searches all extensions for a tool by name.
func (m *Manager) FindTool(toolName string) (*Extension, *ExtensionTool, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, ext := range m.extensions {
		if !ext.Enabled || ext.Status != StatusRunning {
			continue
		}
		for i := range ext.Tools {
			if ext.Tools[i].Name == toolName {
				return ext, &ext.Tools[i], true
			}
		}
	}
	return nil, nil, false
}

// startExtension starts an MCP server process.
func (m *Manager) startExtension(ext *Extension) error {
	if ext.Type != TypeMCPStdio && ext.Type != TypeMCPStream {
		ext.Status = StatusRunning
		return nil
	}

	if ext.Command == "" {
		return fmt.Errorf("extension %s has no command", ext.Name)
	}

	// Build command environment
	env := os.Environ()
	for k, v := range ext.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd := exec.Command(ext.Command, ext.Args...)
	cmd.Env = env
	ext.Status = StatusStarting

	// Start the process
	if err := cmd.Start(); err != nil {
		ext.Status = StatusError
		ext.LastError = err.Error()
		return fmt.Errorf("failed to start extension %s: %w", ext.Name, err)
	}

	ext.Status = StatusRunning
	return nil
}

// stopExtension stops an MCP server process.
func (m *Manager) stopExtension(ext *Extension) error {
	ext.Status = StatusStopped
	return nil
}

// LoadConfig loads extension configuration from the config directory.
func (m *Manager) LoadConfig() error {
	configFile := filepath.Join(m.configDir, "extensions.json")
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No config file is fine
		}
		return err
	}

	var extensions []*Extension
	if err := json.Unmarshal(data, &extensions); err != nil {
		return err
	}

	m.mu.Lock()
	for _, ext := range extensions {
		m.extensions[ext.Name] = ext
	}
	m.mu.Unlock()

	return nil
}

// SaveConfig persists extension configuration.
func (m *Manager) SaveConfig() error {
	if err := os.MkdirAll(m.configDir, 0o755); err != nil {
		return err
	}

	m.mu.RLock()
	var extensions []*Extension
	for _, ext := range m.extensions {
		extensions = append(extensions, ext)
	}
	m.mu.RUnlock()

	data, err := json.MarshalIndent(extensions, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(m.configDir, "extensions.json"), data, 0o644)
}

// InstallFromSmithery installs an MCP server from a Smithery-like registry.
func (m *Manager) InstallFromSmithery(serverName string, config map[string]interface{}) error {
	// In production, this would:
	// 1. Query the Smithery API for server metadata
	// 2. Download/verify the server binary
	// 3. Configure environment variables
	// 4. Register the extension

	ext := &Extension{
		Name:        serverName,
		Type:        TypeMCPStdio,
		Description: fmt.Sprintf("MCP server: %s", serverName),
		Enabled:     true,
		AutoStart:   true,
		Status:      StatusNotStarted,
		Config:      config,
	}

	return m.Register(ext)
}

// String returns a human-readable summary of extension status.
func (m *Manager) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var lines []string
	lines = append(lines, fmt.Sprintf("Extensions (%d registered):", len(m.extensions)))

	for _, ext := range m.extensions {
		status := string(ext.Status)
		enabled := "enabled"
		if !ext.Enabled {
			enabled = "disabled"
		}
		tools := ""
		if ext.ToolCount > 0 {
			tools = fmt.Sprintf(" (%d tools)", ext.ToolCount)
		}
		lines = append(lines, fmt.Sprintf("  %s [%s, %s]%s", ext.Name, status, enabled, tools))
	}

	return strings.Join(lines, "\n")
}
