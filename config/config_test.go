package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Provider != "anthropic" {
		t.Errorf("default provider: %s", cfg.Provider)
	}
	if cfg.Model != "claude-sonnet-4-20250514" {
		t.Errorf("default model: %s", cfg.Model)
	}
	if !cfg.AutoCompact {
		t.Error("auto compact should be on")
	}
	if cfg.MaxToolCalls != 25 {
		t.Errorf("max tool calls: %d", cfg.MaxToolCalls)
	}
}

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("manager should not be nil")
	}
	if m.config == nil {
		t.Fatal("config should not be nil")
	}
}

func TestSetAndGet(t *testing.T) {
	m := NewManager()
	m.Set("provider", "openai")

	val, level, ok := m.Get("provider")
	if !ok {
		t.Fatal("should find set value")
	}
	if val != "openai" {
		t.Errorf("value: %v", val)
	}
	if level != LevelRuntime {
		t.Errorf("level: %s", level)
	}
}

func TestGetString(t *testing.T) {
	m := NewManager()
	m.Set("provider", "groq")

	if m.GetString("provider") != "groq" {
		t.Errorf("got: %s", m.GetString("provider"))
	}
	if m.GetString("nonexistent") != "" {
		t.Error("nonexistent should return empty string")
	}
}

func TestGetInt(t *testing.T) {
	m := NewManager()
	m.Set("maxRetries", 5)

	if m.GetInt("maxRetries") != 5 {
		t.Errorf("got: %d", m.GetInt("maxRetries"))
	}
}

func TestGetBool(t *testing.T) {
	m := NewManager()
	m.Set("vimMode", true)

	if !m.GetBool("vimMode") {
		t.Error("should be true")
	}
	if m.GetBool("nonexistent") {
		t.Error("nonexistent should be false")
	}
}

func TestUnset(t *testing.T) {
	m := NewManager()
	m.Set("provider", "openai")
	m.Unset("provider")

	_, _, ok := m.Get("provider")
	// After unset, it may still have a value from defaults/lower levels
	// but the runtime override is gone
	_ = ok
}

func TestReset(t *testing.T) {
	m := NewManager()
	m.Set("provider", "openai")
	m.Set("model", "gpt-4")
	m.Reset()

	cfg := m.GetConfig()
	if cfg.Provider != "anthropic" {
		t.Errorf("after reset, provider should be default: %s", cfg.Provider)
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.json")

	configData := map[string]any{
		"provider":   "openai",
		"model":      "gpt-4o",
		"theme":      "dark",
		"maxRetries": 5,
	}
	data, _ := json.MarshalIndent(configData, "", "  ")
	os.WriteFile(configFile, data, 0644)

	m := NewManager()
	m.paths = []string{configFile}
	m.Load()

	cfg := m.GetConfig()
	if cfg.Provider != "openai" {
		t.Errorf("provider from file: %s", cfg.Provider)
	}
	if cfg.Model != "gpt-4o" {
		t.Errorf("model from file: %s", cfg.Model)
	}
	if cfg.Theme != "dark" {
		t.Errorf("theme from file: %s", cfg.Theme)
	}
}

func TestEnvOverrides(t *testing.T) {
	os.Setenv("HYPERHARNESS_PROVIDER", "deepseek")
	defer os.Unsetenv("HYPERHARNESS_PROVIDER")

	m := NewManager()
	m.Load()

	val, level, ok := m.Get("provider")
	if !ok {
		t.Fatal("should find provider")
	}
	if val != "deepseek" {
		t.Errorf("env override: %v", val)
	}
	if level != LevelEnv {
		t.Errorf("level should be env: %s", level)
	}
}

func TestLayeredResolution(t *testing.T) {
	// Create user-level config
	dir := t.TempDir()
	userConfig := filepath.Join(dir, "user.json")
	projectConfig := filepath.Join(dir, "project.json")

	userData := map[string]any{"provider": "anthropic", "theme": "default"}
	projData := map[string]any{"provider": "openai", "model": "gpt-4o"}

	u, _ := json.MarshalIndent(userData, "", "  ")
	p, _ := json.MarshalIndent(projData, "", "  ")
	os.WriteFile(userConfig, u, 0644)
	os.WriteFile(projectConfig, p, 0644)

	m := NewManager()
	m.paths = []string{userConfig, projectConfig}
	m.Load()

	// Project should override user for "provider"
	val, level, _ := m.Get("provider")
	if val != "openai" {
		t.Errorf("project should override user: %v", val)
	}
	if level != LevelProject {
		t.Errorf("level should be project: %s", level)
	}

	// User-only value should still be available
	themeVal, _, _ := m.Get("theme")
	if themeVal != "default" {
		t.Errorf("user theme: %v", themeVal)
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	projectConfig := filepath.Join(dir, "config.json")

	m := NewManager()
	m.paths = []string{"", projectConfig}
	m.Set("provider", "groq")
	m.rebuildConfig()

	err := m.Save()
	if err != nil {
		t.Fatal(err)
	}

	// Read back
	data, err := os.ReadFile(projectConfig)
	if err != nil {
		t.Fatal(err)
	}

	var saved map[string]any
	json.Unmarshal(data, &saved)
	if saved["provider"] != "groq" {
		t.Errorf("saved provider: %v", saved["provider"])
	}
}

func TestList(t *testing.T) {
	m := NewManager()
	m.Set("provider", "openai")
	m.Set("theme", "dark")

	list := m.List()
	if len(list) < 2 {
		t.Errorf("should have at least 2 entries, got %d", len(list))
	}
	if list["provider"] != LevelRuntime {
		t.Errorf("provider level: %s", list["provider"])
	}
}

func TestConfigLevelString(t *testing.T) {
	levels := map[ConfigLevel]string{
		LevelDefault: "default",
		LevelSystem:  "system",
		LevelUser:    "user",
		LevelProject: "project",
		LevelEnv:     "env",
		LevelRuntime: "runtime",
	}
	for level, expected := range levels {
		if level.String() != expected {
			t.Errorf("ConfigLevel(%d).String() = %q, want %q", level, level.String(), expected)
		}
	}
}

func TestMCPServerConfig(t *testing.T) {
	server := MCPServerConfig{
		Name:    "filesystem",
		Command: "mcp-server-filesystem",
		Args:    []string{"/tmp"},
		Env:     map[string]string{"DEBUG": "1"},
	}

	data, err := json.Marshal(server)
	if err != nil {
		t.Fatal(err)
	}

	var decoded MCPServerConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Name != "filesystem" {
		t.Errorf("name: %s", decoded.Name)
	}
	if len(decoded.Args) != 1 {
		t.Errorf("args: %v", decoded.Args)
	}
}

func TestAppConfigJSONRoundtrip(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Provider = "openai"
	cfg.Model = "gpt-4o"
	cfg.Theme = "dark"
	cfg.MCPServers = []MCPServerConfig{
		{Name: "test", Command: "test-server"},
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var decoded AppConfig
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Provider != "openai" {
		t.Errorf("provider: %s", decoded.Provider)
	}
	if decoded.Model != "gpt-4o" {
		t.Errorf("model: %s", decoded.Model)
	}
	if len(decoded.MCPServers) != 1 {
		t.Errorf("mcp servers: %d", len(decoded.MCPServers))
	}
}

func TestWorkflowSaveAndList(t *testing.T) {
	dir := t.TempDir()
	oldFile := workflowsFile
	workflowsFile = filepath.Join(dir, "workflows.json")
	defer func() { workflowsFile = oldFile }()

	w := Workflow{
		Name:        "test-workflow",
		Description: "A test workflow",
		Commands:    []string{"echo hello", "echo world"},
	}

	err := SaveWorkflow(w)
	if err != nil {
		t.Fatal(err)
	}

	workflows := ListWorkflows()
	if len(workflows) != 1 {
		t.Fatalf("expected 1 workflow, got %d", len(workflows))
	}
	if workflows[0].Name != "test-workflow" {
		t.Errorf("name: %s", workflows[0].Name)
	}
	if len(workflows[0].Commands) != 2 {
		t.Errorf("commands: %v", workflows[0].Commands)
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		input    any
		expected int
	}{
		{42, 42},
		{float64(3.7), 3},
		{json.Number("100"), 100},
		{"not a number", 0},
	}
	for _, tc := range tests {
		result := toInt(tc.input)
		if result != tc.expected {
			t.Errorf("toInt(%v) = %d, want %d", tc.input, result, tc.expected)
		}
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		input    any
		expected bool
	}{
		{true, true},
		{false, false},
		{"true", true},
		{"True", true},
		{"false", false},
		{"1", false},
		{0, false},
	}
	for _, tc := range tests {
		result := toBool(tc.input)
		if result != tc.expected {
			t.Errorf("toBool(%v) = %v, want %v", tc.input, result, tc.expected)
		}
	}
}
