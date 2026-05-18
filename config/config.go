package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ── Compatibility aliases (for packages importing the old Config) ──

// Config is an alias for AppConfig for backward compatibility.
type Config = AppConfig

// LoadConfig returns the default config resolved from environment.
// This is the legacy API used by existing packages.
func LoadConfig() *AppConfig {
	m := NewManager()
	m.Load()
	return m.GetConfig()
}

// ConfigLevel defines where a config value comes from.
type ConfigLevel int

const (
	LevelDefault ConfigLevel = iota
	LevelSystem              // /etc/hyperharness or C:\ProgramData
	LevelUser                // ~/.hyperharness
	LevelProject             // .hyperharness/ in project dir
	LevelEnv                 // Environment variable override
	LevelRuntime             // Set at runtime (highest priority)
)

func (l ConfigLevel) String() string {
	switch l {
	case LevelDefault:
		return "default"
	case LevelSystem:
		return "system"
	case LevelUser:
		return "user"
	case LevelProject:
		return "project"
	case LevelEnv:
		return "env"
	case LevelRuntime:
		return "runtime"
	default:
		return "unknown"
	}
}

// AppConfig is the top-level application configuration.
type AppConfig struct {
	// Provider settings
	Provider       string `json:"provider"`
	Model          string `json:"model"`
	APIKey         string `json:"apiKey,omitempty"`
	BaseURL        string `json:"baseUrl,omitempty"`
	ThinkingLevel  string `json:"thinkingLevel,omitempty"`

	// Session settings
	SessionDir     string `json:"sessionDir,omitempty"`
	AutoCompact    bool   `json:"autoCompact"`
	CompactAtPct   int    `json:"compactAtPct"`
	MaxRetries     int    `json:"maxRetries"`

	// TUI settings
	Theme          string `json:"theme"`
	VimMode        bool   `json:"vimMode"`
	DebugMode      bool   `json:"debugMode"`
	AltScreen      bool   `json:"altScreen"`
	MouseSupport   bool   `json:"mouseSupport"`

	// Security settings
	AutonomyLevel  string `json:"autonomyLevel"`
	Sandboxed      bool   `json:"sandboxed"`
	AutoApprove    []string `json:"autoApprove,omitempty"`
	AutoDeny       []string `json:"autoDeny,omitempty"`

	// Tool settings
	ToolTimeout    int      `json:"toolTimeout"`
	MaxToolCalls   int      `json:"maxToolCalls"`
	DisabledTools  []string `json:"disabledTools,omitempty"`

	// MCP settings
	MCPConfigPath  string   `json:"mcpConfigPath,omitempty"`
	MCPServers     []MCPServerConfig `json:"mcpServers,omitempty"`

	// Display settings
	OutputStyle    string `json:"outputStyle"`
	ShowThinking   bool   `json:"showThinking"`
	ShowToolOutput bool   `json:"showToolOutput"`
	PromptFormat   string `json:"promptFormat"`

	// Editor
	Editor         string `json:"editor,omitempty"`

	// History
	MaxHistorySize int `json:"maxHistorySize"`

	// RAG
	RAGEnabled     bool   `json:"ragEnabled"`
	RAGIndexPath   string `json:"ragIndexPath,omitempty"`
}

// MCPServerConfig defines an MCP server connection.
type MCPServerConfig struct {
	Name    string            `json:"name"`
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	URL     string            `json:"url,omitempty"`
}

// DefaultConfig returns the default application configuration.
func DefaultConfig() *AppConfig {
	return &AppConfig{
		Provider:       "anthropic",
		Model:          "claude-sonnet-4-20250514",
		ThinkingLevel:  "medium",
		AutoCompact:    true,
		CompactAtPct:   80,
		MaxRetries:     3,
		Theme:          "default",
		AltScreen:      true,
		MouseSupport:   true,
		AutonomyLevel:  "read",
		Sandboxed:      true,
		ToolTimeout:    60,
		MaxToolCalls:   25,
		OutputStyle:    "default",
		ShowThinking:   true,
		ShowToolOutput: true,
		MaxHistorySize: 1000,
		RAGEnabled:     false,
	}
}

// configEntry tracks the source of a config value.
type configEntry struct {
	Value any
	Level ConfigLevel
}

// Manager is the configuration manager with layered resolution.
type Manager struct {
	mu      sync.RWMutex
	entries map[string]configEntry
	config  *AppConfig
	paths   []string // Config file search paths (in priority order)
}

// NewManager creates a configuration manager with layered resolution.
func NewManager() *Manager {
	m := &Manager{
		entries: make(map[string]configEntry),
		config:  DefaultConfig(),
	}

	// Build search paths
	home, _ := os.UserHomeDir()
	cwd, _ := os.Getwd()

	m.paths = []string{
		filepath.Join(home, ".hyperharness", "config.json"),     // User level
		filepath.Join(cwd, ".hyperharness", "config.json"),      // Project level
	}

	return m
}

// Load reads configuration from all layers.
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Load system-level config
	if systemPath := systemConfigPath(); systemPath != "" {
		if err := m.loadFile(systemPath, LevelSystem); err == nil {
			// System config is optional
		}
	}

	// Load user-level config
	if len(m.paths) > 0 {
		m.loadFile(m.paths[0], LevelUser)
	}

	// Load project-level config
	if len(m.paths) > 1 {
		m.loadFile(m.paths[1], LevelProject)
	}

	// Apply environment variable overrides
	m.applyEnvOverrides()

	// Build the final config from entries
	m.rebuildConfig()

	return nil
}

// Get retrieves a configuration value with its source level.
func (m *Manager) Get(key string) (any, ConfigLevel, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.entries[key]
	if !ok {
		return nil, LevelDefault, false
	}
	return entry.Value, entry.Level, true
}

// Set sets a runtime configuration value (highest priority).
func (m *Manager) Set(key string, value any) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries[key] = configEntry{Value: value, Level: LevelRuntime}
	m.rebuildConfig()
}

// GetString returns a config value as a string.
func (m *Manager) GetString(key string) string {
	val, _, ok := m.Get(key)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", val)
}

// GetInt returns a config value as an int.
func (m *Manager) GetInt(key string) int {
	val, _, ok := m.Get(key)
	if !ok {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case json.Number:
		i, _ := v.Int64()
		return int(i)
	default:
		return 0
	}
}

// GetBool returns a config value as a bool.
func (m *Manager) GetBool(key string) bool {
	val, _, ok := m.Get(key)
	if !ok {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	default:
		return fmt.Sprintf("%v", v) == "true"
	}
}

// GetConfig returns the resolved application config.
func (m *Manager) GetConfig() *AppConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// Save saves the current config to the project-level file.
func (m *Manager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.paths) < 2 {
		return fmt.Errorf("no project config path")
	}

	dir := filepath.Dir(m.paths[1])
	os.MkdirAll(dir, 0755)

	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.paths[1], data, 0644)
}

// SaveUser saves the current config to the user-level file.
func (m *Manager) SaveUser() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.paths) < 1 {
		return fmt.Errorf("no user config path")
	}

	dir := filepath.Dir(m.paths[0])
	os.MkdirAll(dir, 0755)

	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.paths[0], data, 0644)
}

// List shows all config values with their sources.
func (m *Manager) List() map[string]ConfigLevel {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]ConfigLevel)
	for key, entry := range m.entries {
		result[key] = entry.Level
	}
	return result
}

// Unset removes a runtime configuration value.
func (m *Manager) Unset(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.entries, key)
	m.rebuildConfig()
}

// Reset restores all config values to defaults.
func (m *Manager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries = make(map[string]configEntry)
	m.config = DefaultConfig()
}

// ── Internal helpers ──

func (m *Manager) loadFile(path string, level ConfigLevel) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("invalid config file %s: %w", path, err)
	}

	for key, value := range raw {
		// Only set if not already set at a higher priority level
		if existing, ok := m.entries[key]; ok && existing.Level > level {
			continue
		}
		m.entries[key] = configEntry{Value: value, Level: level}
	}

	return nil
}

func (m *Manager) applyEnvOverrides() {
	envMappings := map[string]string{
		"HYPERHARNESS_PROVIDER":  "provider",
		"HYPERHARNESS_MODEL":     "model",
		"HYPERHARNESS_API_KEY":   "apiKey",
		"HYPERHARNESS_BASE_URL":  "baseUrl",
		"HYPERHARNESS_THINKING":  "thinkingLevel",
		"HYPERHARNESS_THEME":     "theme",
		"HYPERHARNESS_EDITOR":    "editor",
		"HYPERHARNESS_AUTONOMY":  "autonomyLevel",
		"HYPERHARNESS_DEBUG":     "debugMode",
		"HYPERHARNESS_VIM":       "vimMode",
		"HYPERHARNESS_OUTPUT":    "outputStyle",

		// Legacy env var support
		"SUPERCLI_PROVIDER":      "provider",
		"SUPERCLI_MODEL":         "model",
		"OPENAI_API_KEY":         "apiKey",
		"OPENAI_BASE_URL":        "baseUrl",
		"ANTHROPIC_API_KEY":      "apiKey",
		"EDITOR":                 "editor",
	}

	for envVar, configKey := range envMappings {
		if val := os.Getenv(envVar); val != "" {
			m.entries[configKey] = configEntry{Value: val, Level: LevelEnv}
		}
	}
}

func (m *Manager) rebuildConfig() {
	// Start from defaults and apply entries
	cfg := DefaultConfig()

	if entry, ok := m.entries["provider"]; ok {
		v := entry.Value
		cfg.Provider = fmt.Sprintf("%v", v)
	}
	if entry, ok := m.entries["model"]; ok {
		v := entry.Value
		cfg.Model = fmt.Sprintf("%v", v)
	}
	if entry, ok := m.entries["apiKey"]; ok {
		v := entry.Value
		cfg.APIKey = fmt.Sprintf("%v", v)
	}
	if entry, ok := m.entries["baseUrl"]; ok {
		v := entry.Value
		cfg.BaseURL = fmt.Sprintf("%v", v)
	}
	if entry, ok := m.entries["thinkingLevel"]; ok {
		v := entry.Value
		cfg.ThinkingLevel = fmt.Sprintf("%v", v)
	}
	if entry, ok := m.entries["sessionDir"]; ok {
		v := entry.Value
		cfg.SessionDir = fmt.Sprintf("%v", v)
	}
	if entry, ok := m.entries["autoCompact"]; ok {
		v := entry.Value
		cfg.AutoCompact = toBool(v)
	}
	if entry, ok := m.entries["compactAtPct"]; ok {
		v := entry.Value
		cfg.CompactAtPct = toInt(v)
	}
	if entry, ok := m.entries["maxRetries"]; ok {
		v := entry.Value
		cfg.MaxRetries = toInt(v)
	}
	if entry, ok := m.entries["theme"]; ok {
		v := entry.Value
		cfg.Theme = fmt.Sprintf("%v", v)
	}
	if entry, ok := m.entries["vimMode"]; ok {
		v := entry.Value
		cfg.VimMode = toBool(v)
	}
	if entry, ok := m.entries["debugMode"]; ok {
		v := entry.Value
		cfg.DebugMode = toBool(v)
	}
	if entry, ok := m.entries["altScreen"]; ok {
		v := entry.Value
		cfg.AltScreen = toBool(v)
	}
	if entry, ok := m.entries["mouseSupport"]; ok {
		v := entry.Value
		cfg.MouseSupport = toBool(v)
	}
	if entry, ok := m.entries["autonomyLevel"]; ok {
		v := entry.Value
		cfg.AutonomyLevel = fmt.Sprintf("%v", v)
	}
	if entry, ok := m.entries["sandboxed"]; ok {
		v := entry.Value
		cfg.Sandboxed = toBool(v)
	}
	if entry, ok := m.entries["toolTimeout"]; ok {
		v := entry.Value
		cfg.ToolTimeout = toInt(v)
	}
	if entry, ok := m.entries["maxToolCalls"]; ok {
		v := entry.Value
		cfg.MaxToolCalls = toInt(v)
	}
	if entry, ok := m.entries["outputStyle"]; ok {
		v := entry.Value
		cfg.OutputStyle = fmt.Sprintf("%v", v)
	}
	if entry, ok := m.entries["showThinking"]; ok {
		v := entry.Value
		cfg.ShowThinking = toBool(v)
	}
	if entry, ok := m.entries["showToolOutput"]; ok {
		v := entry.Value
		cfg.ShowToolOutput = toBool(v)
	}
	if entry, ok := m.entries["editor"]; ok {
		v := entry.Value
		cfg.Editor = fmt.Sprintf("%v", v)
	}
	if entry, ok := m.entries["maxHistorySize"]; ok {
		v := entry.Value
		cfg.MaxHistorySize = toInt(v)
	}
	if entry, ok := m.entries["ragEnabled"]; ok {
		v := entry.Value
		cfg.RAGEnabled = toBool(v)
	}
	if entry, ok := m.entries["ragIndexPath"]; ok {
		v := entry.Value
		cfg.RAGIndexPath = fmt.Sprintf("%v", v)
	}
	if entry, ok := m.entries["mcpConfigPath"]; ok {
		v := entry.Value
		cfg.MCPConfigPath = fmt.Sprintf("%v", v)
	}
	if entry, ok := m.entries["promptFormat"]; ok {
		v := entry.Value
		cfg.PromptFormat = fmt.Sprintf("%v", v)
	}

	m.config = cfg
}

func systemConfigPath() string {
	// Check for system-level config
	paths := []string{
		"/etc/hyperharness/config.json",
	}
	if home := os.Getenv("ProgramData"); home != "" {
		paths = append(paths, filepath.Join(home, "hyperharness", "config.json"))
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func toBool(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return strings.ToLower(val) == "true"
	default:
		return fmt.Sprintf("%v", val) == "true"
	}
}

func toInt(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	case json.Number:
		i, _ := val.Int64()
		return int(i)
	default:
		return 0
	}
}
