// Package config provides unified configuration for the hyperharness.
// Merges global (~/.hyperharness/config.json) and project (.hyperharness/config.json)
// configs with environment variable overrides, following Pi's settings pattern
// but with Go performance and type safety.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	ConfigDirName  = ".hyperharness"
	GlobalConfig   = "config.json"
	ProjectConfig  = "config.json"
	GlobalSettings = "settings.json"
)

// ProviderConfig holds provider-specific settings.
type ProviderConfig struct {
	APIKey      string `json:"apiKey,omitempty"`
	BaseURL     string `json:"baseURL,omitempty"`
	Timeout     int    `json:"timeout,omitempty"`     // seconds
	MaxRetries  int    `json:"maxRetries,omitempty"`  // default 3
	RetryDelay  int    `json:"retryDelay,omitempty"`  // ms, default 2000
}

// ModelConfig holds model-specific settings.
type ModelConfig struct {
	Provider       string  `json:"provider"`
	ModelID        string  `json:"modelId"`
	MaxTokens      int     `json:"maxTokens,omitempty"`
	Temperature    float64 `json:"temperature,omitempty"`
	TopP           float64 `json:"topP,omitempty"`
	FrequencyPenalty float64 `json:"frequencyPenalty,omitempty"`
}

// CompactionSettings controls context compaction behavior.
type CompactionSettings struct {
	Enabled         bool `json:"enabled"`
	ReserveTokens   int  `json:"reserveTokens"`  // default 16384
	KeepRecentTokens int `json:"keepRecentTokens"` // default 20000
}

// RetrySettings controls retry behavior.
type RetrySettings struct {
	Enabled    bool `json:"enabled"`
	MaxRetries int  `json:"maxRetries"`
	BaseDelay  int  `json:"baseDelayMs"`
	MaxDelay   int  `json:"maxDelayMs"`
}

// SessionSettings controls session behavior.
type SessionSettings struct {
	SessionDir string `json:"sessionDir,omitempty"`
	NoSession  bool   `json:"noSession,omitempty"`
}

// ToolSettings controls which tools are enabled.
type ToolSettings struct {
	Enabled    []string `json:"enabled,omitempty"`
	Disabled   []string `json:"disabled,omitempty"`
	MaxLines         int `json:"maxLines,omitempty"`
	MaxBytes         int64 `json:"maxBytes,omitempty"`
	Timeout  int    `json:"timeout,omitempty"`
	CWD      string `json:"cwd,omitempty"`
}

// UISettings controls terminal UI behavior.
type UISettings struct {
	Theme              string `json:"theme"`
	ShowImages         bool   `json:"showImages"`
	ClearOnShrink      bool   `json:"clearOnShrink"`
	AutoResizeImages   bool   `json:"autoResizeImages"`
	MaxImageSize       int    `json:"maxImageSize"`
	EditorPaddingX     int    `json:"editorPaddingX"`
	AutocompleteMaxVis int    `json:"autocompleteMaxVisible"`
	ShowHardwareCursor bool   `json:"showHardwareCursor"`
	CompactDisplay     string `json:"compactDisplay"` // "full", "summary", "minimal"
}

// MemorySettings controls memory/KB behavior.
type MemorySettings struct {
	Enabled       bool     `json:"enabled"`
	Backend       string   `json:"backend"` // "sqlite", "memory", "file"
	Path          string   `json:"path,omitempty"`
	AutoIndex     bool     `json:"autoIndex"`
	Scope         string   `json:"scope"` // "global", "project", "session"
}

// MCPSettings controls MCP client/server behavior.
type MCPSettings struct {
	Enabled         bool              `json:"enabled"`
	ServerMode      bool              `json:"serverMode"`
	ClientMode      bool              `json:"clientMode"`
	Servers         map[string]MCPConn `json:"servers,omitempty"`
	StdioTimeout    int               `json:"stdioTimeout,omitempty"`
	SSETimeout      int               `json:"sseTimeout,omitempty"`
	RegistryPath    string            `json:"registryPath,omitempty"`
}

// MCPConn defines an MCP server connection.
type MCPConn struct {
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	URL     string            `json:"url,omitempty"`
	Transport string          `json:"transport,omitempty"` // "stdio", "sse"
	Env     map[string]string `json:"env,omitempty"`
}

// SkillSettings controls skill loading.
type SkillSettings struct {
	Enabled    bool     `json:"enabled"`
	Paths      []string `json:"paths,omitempty"`
	AutoLoad   bool     `json:"autoLoad"`
}

// ExtensionSettings controls extension loading.
type ExtensionSettings struct {
	Enabled    bool     `json:"enabled"`
	Paths      []string `json:"paths,omitempty"`
	AutoDiscover bool   `json:"autoDiscover"`
}

// PromptSettings controls prompt templates.
type PromptSettings struct {
	Enabled    bool     `json:"enabled"`
	Paths      []string `json:"paths,omitempty"`
	AutoDiscover bool   `json:"autoDiscover"`
}

// PackageSettings controls package management.
type PackageSettings struct {
	Packages []PackageSource `json:"packages,omitempty"`
}

// PackageSource defines a package to install.
type PackageSource struct {
	Name       string   `json:"name"`
	Version    string   `json:"version,omitempty"`
	Source     string   `json:"source"`
	Skills     []string `json:"skills,omitempty"`
	Extensions []string `json:"extensions,omitempty"`
	Prompts    []string `json:"prompts,omitempty"`
	Themes     []string `json:"themes,omitempty"`
}

// Settings is the complete hyperharness configuration.
type Settings struct {
	// Model & Provider
	DefaultProvider    string                    `json:"defaultProvider"`
	DefaultModel       string                    `json:"defaultModel"`
	DefaultThinkingLevel string                  `json:"defaultThinkingLevel"`
	HideThinkingBlock  bool                      `json:"hideThinkingBlock"`
	Providers          map[string]*ProviderConfig `json:"providers,omitempty"`
	Models             map[string]*ModelConfig    `json:"models,omitempty"`
	ThinkingBudgets    map[string]int            `json:"thinkingBudgets,omitempty"`
	EnabledModels      []string                 `json:"enabledModels,omitempty"`

	// UI
	UI *UISettings `json:"ui,omitempty"`

	// Compaction
	Compaction *CompactionSettings `json:"compaction,omitempty"`

	// Retry
	Retry *RetrySettings `json:"retry,omitempty"`

	// Sessions
	Sessions *SessionSettings `json:"sessions,omitempty"`

	// Tools
	Tools *ToolSettings `json:"tools,omitempty"`

	// Memory
	Memory *MemorySettings `json:"memory,omitempty"`

	// MCP
	MCP *MCPSettings `json:"mcp,omitempty"`

	// Skills
	Skills *SkillSettings `json:"skills,omitempty"`

	// Extensions
	Extensions *ExtensionSettings `json:"extensions,omitempty"`

	// Prompts
	Prompts *PromptSettings `json:"prompts,omitempty"`

	// Packages
	Packages *PackageSettings `json:"packages,omitempty"`

	// Shell
	ShellPath          string   `json:"shellPath,omitempty"`
	ShellCommandPrefix string   `json:"shellCommandPrefix,omitempty"`
	NPMCommand    []string `json:"npmCommand,omitempty"`

	// Message delivery
	SteeringMode  string `json:"steeringMode"`  // "one-at-a-time", "all"
	FollowUpMode  string `json:"followUpMode"`  // "one-at-a-time", "all"
	Transport     string `json:"transport"`     // "sse", "websocket", "auto"

	// Markdown
	MarkdownCodeBlockIndent string `json:"markdownCodeBlockIndent,omitempty"`

	// System prompt
	SystemPrompt        string `json:"systemPrompt,omitempty"`
	AppendSystemPrompt  string `json:"appendSystemPrompt,omitempty"`

	mu sync.RWMutex
}

// Defaults returns the default settings.
func Defaults() *Settings {
	return &Settings{
		DefaultThinkingLevel: "medium",
		UI: &UISettings{
			Theme:              "dark",
			ShowImages:         true,
			AutoResizeImages:   true,
			MaxImageSize:       2000,
			AutocompleteMaxVis: 5,
			CompactDisplay:     "full",
		},
		Compaction: &CompactionSettings{
			Enabled:          true,
			ReserveTokens:    16384,
			KeepRecentTokens: 20000,
		},
		Retry: &RetrySettings{
			Enabled:    true,
			MaxRetries: 3,
			BaseDelay:  2000,
			MaxDelay:   60000,
		},
		Memory: &MemorySettings{
			Enabled:   true,
			Backend:   "sqlite",
			AutoIndex: true,
		},
		MCP: &MCPSettings{
			Enabled:    true,
			ServerMode: true,
			ClientMode: true,
			Servers:    make(map[string]MCPConn),
		},
		Skills: &SkillSettings{
			Enabled:   true,
			AutoLoad: true,
		},
		Extensions: &ExtensionSettings{
			Enabled:      true,
			AutoDiscover: true,
		},
		Prompts: &PromptSettings{
			Enabled:      true,
			AutoDiscover: false,
		},
	}
}

// Manager handles configuration loading and merging.
type Manager struct {
	globalDir  string
	projectDir string
	settings   *Settings
	mu         sync.RWMutex
}

// NewManager creates a new configuration manager.
func NewManager(globalDir, projectDir string) (*Manager, error) {
	m := &Manager{
		globalDir: globalDir,
		projectDir: projectDir,
		settings:   Defaults(),
	}
	if err := m.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return m, nil
}

// Load reads global and project configs, merging them with env overrides.
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Start with defaults
	m.settings = Defaults()

	// Load global config
	if err := m.loadFile(m.globalDir, GlobalSettings, m.settings); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("global config: %w", err)
		}
	}

	// Load project config (overrides global)
	if err := m.loadFile(m.projectDir, ProjectConfig, m.settings); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("project config: %w", err)
		}
	}

	// Apply environment variable overrides
	m.applyEnvOverrides()

	return nil
}

// loadFile reads a JSON config file and merges into settings.
func (m *Manager) loadFile(dir, filename string, target *Settings) error {
	path := filepath.Join(dir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// applyEnvOverrides applies environment variable settings.
func (m *Manager) applyEnvOverrides() {
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		if m.settings.Providers == nil {
			m.settings.Providers = make(map[string]*ProviderConfig)
		}
		m.settings.Providers["anthropic"] = &ProviderConfig{APIKey: key}
	}
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		if m.settings.Providers == nil {
			m.settings.Providers = make(map[string]*ProviderConfig)
		}
		m.settings.Providers["openai"] = &ProviderConfig{APIKey: key}
	}
	if key := os.Getenv("GOOGLE_API_KEY"); key != "" {
		if m.settings.Providers == nil {
			m.settings.Providers = make(map[string]*ProviderConfig)
		}
		m.settings.Providers["google"] = &ProviderConfig{APIKey: key}
	}
	if model := os.Getenv("HYPERHARNESS_MODEL"); model != "" {
		m.settings.DefaultModel = model
	}
	if provider := os.Getenv("HYPERHARNESS_PROVIDER"); provider != "" {
		m.settings.DefaultProvider = provider
	}
}

// Get returns the current settings (thread-safe).
func (m *Manager) Get() *Settings {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings
}

// Update merges new settings into current settings.
func (m *Manager) Update(newSettings *Settings) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Merge at JSON level for simplicity
	merged, err := json.Marshal(m.settings)
	if err != nil {
		return err
	}
	overlay, err := json.Marshal(newSettings)
	if err != nil {
		return err
	}
	var mergedMap map[string]interface{}
	if err := json.Unmarshal(merged, &mergedMap); err != nil {
		return err
	}
	var overlayMap map[string]interface{}
	if err := json.Unmarshal(overlay, &overlayMap); err != nil {
		return err
	}
	deepMerge(mergedMap, overlayMap)
	mergedBytes, err := json.Marshal(mergedMap)
	if err != nil {
		return err
	}
	return json.Unmarshal(mergedBytes, m.settings)
}

// Save persists current settings to project config.
func (m *Manager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	dir := m.projectDir
	if dir == "" {
		dir = m.globalDir
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(m.settings, "", "  ")
	if err != nil {
		return err
	}
	
	path := filepath.Join(dir, ProjectConfig)
	return os.WriteFile(path, data, 0o644)
}

// deepMerge recursively merges overlay into base.
func deepMerge(base, overlay map[string]interface{}) {
	for key, overlayVal := range overlay {
		if baseVal, ok := base[key]; ok {
			if baseMap, ok := baseVal.(map[string]interface{}); ok {
				if overlayMap, ok := overlayVal.(map[string]interface{}); ok {
					deepMerge(baseMap, overlayMap)
					continue
				}
			}
		}
		base[key] = overlayVal
	}
}

// ConfigDir returns the global config directory.
func ConfigDir() string {
	if dir := os.Getenv("HYPERHARNESS_DIR"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".hyperharness")
}

// ProjectConfigDir returns the project config directory.
func ProjectConfigDir() string {
	cwd, _ := os.Getwd()
	for dir := cwd; dir != "" && dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, ".hyperharness")); err == nil {
			return filepath.Join(dir, ".hyperharness")
		}
	}
	return ""
}

// FindAgentFiles walks up from cwd finding AGENTS.md files.
func FindAgentFiles() []string {
	var files []string
	cwd, err := os.Getwd()
	if err != nil {
		return files
	}
	
	checkDirs := []string{cwd}
	// Also check standard agent skill dirs
	home, _ := os.UserHomeDir()
	standardDirs := []string{
		filepath.Join(home, ".hyperharness", "AGENTS.md"),
		filepath.Join(home, ".agents", "skills"),
	}
	for _, d := range standardDirs {
		if _, err := os.Stat(d); err == nil {
			checkDirs = append(checkDirs, d)
		}
	}
	
	for dir := cwd; dir != "" && dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		// Check for AGENTS.md or CLAUDE.md
		for _, name := range []string{"AGENTS.md", "CLAUDE.md"} {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err == nil {
				files = append(files, path)
			}
		}
	}
	return files
}

// IsEnabled checks if a tool/feature is enabled given the settings.
func (s *Settings) IsEnabled(name string) bool {
	if s.Tools == nil {
		return true
	}
	// Explicitly disabled?
	for _, d := range s.Tools.Disabled {
		if matchGlob(d, name) {
			return false
		}
	}
	// Explicitly enabled list?
	if len(s.Tools.Enabled) > 0 {
		for _, e := range s.Tools.Enabled {
			if matchGlob(e, name) {
				return true
			}
		}
		return false
	}
	return true
}

// matchGlob performs simple glob matching (* only).
func matchGlob(pattern, name string) bool {
	if pattern == name {
		return true
	}
	if pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(name, strings.TrimSuffix(pattern, "*"))
	}
	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(name, strings.TrimPrefix(pattern, "*"))
	}
	return false
}
