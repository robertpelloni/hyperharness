package extensions

import (
	"encoding/json"
	"testing"
)

func TestManagerRegister(t *testing.T) {
	m := NewManager(t.TempDir())
	ext := &Extension{
		Name:    "test-ext",
		Type:    TypeMCPStdio,
		Command: "echo",
		Enabled: true,
	}

	if err := m.Register(ext); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	found, ok := m.Get("test-ext")
	if !ok {
		t.Fatal("extension should be found")
	}
	if found.Name != "test-ext" {
		t.Errorf("name: %s", found.Name)
	}
}

func TestManagerRegisterEmptyName(t *testing.T) {
	m := NewManager(t.TempDir())
	ext := &Extension{Name: ""}
	if err := m.Register(ext); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestManagerUnregister(t *testing.T) {
	m := NewManager(t.TempDir())
	m.Register(&Extension{Name: "test-ext", Type: TypeBuiltIn, Enabled: true})

	if err := m.Unregister("test-ext"); err != nil {
		t.Fatalf("unregister failed: %v", err)
	}

	_, ok := m.Get("test-ext")
	if ok {
		t.Error("extension should be gone after unregister")
	}
}

func TestManagerEnableDisable(t *testing.T) {
	m := NewManager(t.TempDir())
	m.Register(&Extension{Name: "test-ext", Type: TypeBuiltIn})

	if err := m.Enable("test-ext"); err != nil {
		t.Fatalf("enable failed: %v", err)
	}

	ext, _ := m.Get("test-ext")
	if !ext.Enabled {
		t.Error("should be enabled")
	}

	if err := m.Disable("test-ext"); err != nil {
		t.Fatalf("disable failed: %v", err)
	}

	ext, _ = m.Get("test-ext")
	if ext.Enabled {
		t.Error("should be disabled")
	}
}

func TestManagerList(t *testing.T) {
	m := NewManager(t.TempDir())
	m.Register(&Extension{Name: "b-ext", Enabled: true})
	m.Register(&Extension{Name: "a-ext", Enabled: false})
	m.Register(&Extension{Name: "c-ext", Enabled: true})

	all := m.List()
	if len(all) != 3 {
		t.Fatalf("expected 3, got %d", len(all))
	}

	// Should be sorted
	if all[0].Name != "a-ext" {
		t.Errorf("expected sorted, first: %s", all[0].Name)
	}

	enabled := m.ListEnabled()
	if len(enabled) != 2 {
		t.Errorf("expected 2 enabled, got %d", len(enabled))
	}
}

func TestManagerSaveLoad(t *testing.T) {
	dir := t.TempDir()
	m1 := NewManager(dir)
	m1.Register(&Extension{
		Name:    "test-ext",
		Type:    TypeMCPStdio,
		Command: "test-cmd",
		Args:    []string{"arg1"},
		Enabled: true,
	})

	if err := m1.SaveConfig(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	m2 := NewManager(dir)
	if err := m2.LoadConfig(); err != nil {
		t.Fatalf("load failed: %v", err)
	}

	ext, ok := m2.Get("test-ext")
	if !ok {
		t.Fatal("extension should be found after load")
	}
	if ext.Command != "test-cmd" {
		t.Errorf("command: %s", ext.Command)
	}
}

func TestManagerTools(t *testing.T) {
	m := NewManager(t.TempDir())
	ext := &Extension{
		Name:    "tool-ext",
		Type:    TypeBuiltIn,
		Enabled: true,
		Status:  StatusRunning,
		Tools: []ExtensionTool{
			{Name: "tool1", Description: "First tool"},
			{Name: "tool2", Description: "Second tool"},
		},
	}
	m.Register(ext)

	tools := m.GetAllTools()
	if len(tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(tools))
	}

	// FindTool
	_, tool, ok := m.FindTool("tool1")
	if !ok {
		t.Fatal("tool1 should be found")
	}
	if tool.Name != "tool1" {
		t.Errorf("tool name: %s", tool.Name)
	}
}

func TestManagerString(t *testing.T) {
	m := NewManager(t.TempDir())
	m.Register(&Extension{
		Name:    "ext1",
		Type:    TypeBuiltIn,
		Enabled: true,
		Status:  StatusRunning,
	})

	s := m.String()
	if s == "" {
		t.Error("string should not be empty")
	}
}

func TestInstallFromSmithery(t *testing.T) {
	m := NewManager(t.TempDir())
	err := m.InstallFromSmithery("test-server", map[string]interface{}{
		"port": 8080,
	})
	if err != nil {
		t.Fatalf("install failed: %v", err)
	}

	ext, ok := m.Get("test-server")
	if !ok {
		t.Fatal("server should be registered")
	}
	if !ext.Enabled {
		t.Error("should be enabled by default")
	}
	if !ext.AutoStart {
		t.Error("should auto-start by default")
	}
}

func TestExtensionTypes(t *testing.T) {
	tests := []struct {
		extType ExtensionType
		str     string
	}{
		{TypeMCPStdio, "mcp_stdio"},
		{TypeMCPStream, "mcp_stream"},
		{TypeBuiltIn, "builtin"},
		{TypePlugin, "plugin"},
	}
	for _, tc := range tests {
		if string(tc.extType) != tc.str {
			t.Errorf("extension type: %s != %s", tc.extType, tc.str)
		}
	}
}

func TestExtensionJSONRoundtrip(t *testing.T) {
	ext := &Extension{
		Name:        "json-test",
		Version:     "1.0.0",
		Description: "test extension",
		Type:        TypeMCPStdio,
		Command:     "test",
		Enabled:     true,
		Status:      StatusRunning,
		Tools: []ExtensionTool{
			{Name: "tool1", Description: "test"},
		},
	}

	data, err := json.Marshal(ext)
	if err != nil {
		t.Fatal(err)
	}

	var decoded Extension
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.Name != "json-test" {
		t.Errorf("name: %s", decoded.Name)
	}
	if decoded.Type != TypeMCPStdio {
		t.Errorf("type: %s", decoded.Type)
	}
	if len(decoded.Tools) != 1 {
		t.Errorf("tools count: %d", len(decoded.Tools))
	}
}
