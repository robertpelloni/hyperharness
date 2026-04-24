package borg

import (
	"encoding/json"
	"testing"
)

// mockAdapter is a test adapter
type mockAdapter struct {
	name    string
	started bool
	stopped bool
}

func (m *mockAdapter) Name() string                 { return m.name }
func (m *mockAdapter) Start() error                 { m.started = true; return nil }
func (m *mockAdapter) Stop() error                  { m.stopped = true; return nil }
func (m *mockAdapter) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":    m.name,
		"started": m.started,
		"stopped": m.stopped,
	}
}

func TestCoreStartStop(t *testing.T) {
	core := NewCore(nil)

	if core.IsRunning() {
		t.Error("core should not be running initially")
	}

	if err := core.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}

	if !core.IsRunning() {
		t.Error("core should be running after start")
	}

	if err := core.Stop(); err != nil {
		t.Fatalf("stop failed: %v", err)
	}

	if core.IsRunning() {
		t.Error("core should not be running after stop")
	}
}

func TestCoreAdapters(t *testing.T) {
	core := NewCore(nil)
	adapter := &mockAdapter{name: "test-adapter"}
	core.RegisterAdapter(adapter)

	found, ok := core.GetAdapter("test-adapter")
	if !ok {
		t.Fatal("adapter should be found")
	}
	if found.Name() != "test-adapter" {
		t.Errorf("adapter name: %s", found.Name())
	}

	core.Start()
	if !adapter.started {
		t.Error("adapter should be started")
	}

	core.Stop()
	if !adapter.stopped {
		t.Error("adapter should be stopped")
	}
}

func TestCoreStatus(t *testing.T) {
	core := NewCore(nil)
	core.Start()

	status := core.Status()
	if !status.Active {
		t.Error("status should be active")
	}
	if status.Version == "" {
		t.Error("version should not be empty")
	}

	statusJSON := core.StatusJSON()
	if statusJSON == "" {
		t.Error("status JSON should not be empty")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(statusJSON), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestCoreProperties(t *testing.T) {
	core := NewCore(nil)

	core.SetProperty("test_key", "test_value")
	val := core.GetProperty("test_key")
	if val != "test_value" {
		t.Errorf("property value: %s", val)
	}

	missing := core.GetProperty("nonexistent")
	if missing != "" {
		t.Errorf("missing property should be empty: %s", missing)
	}
}

func TestCoreEvents(t *testing.T) {
	core := NewCore(nil)
	events := make([]string, 0)

	core.On("test:event", func(event string, data interface{}) error {
		events = append(events, event)
		return nil
	})

	core.Emit("test:event", map[string]string{"key": "value"})

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0] != "test:event" {
		t.Errorf("event name: %s", events[0])
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.DataDir == "" {
		t.Error("data dir should not be empty")
	}
	if config.PluginDir == "" {
		t.Error("plugin dir should not be empty")
	}
	if config.MaxMemory == 0 {
		t.Error("max memory should not be zero")
	}
	if config.Properties["engine"] != "hyperharness-go" {
		t.Errorf("engine property: %s", config.Properties["engine"])
	}
}

func TestCoreMultipleAdapters(t *testing.T) {
	core := NewCore(nil)
	a1 := &mockAdapter{name: "adapter-1"}
	a2 := &mockAdapter{name: "adapter-2"}
	a3 := &mockAdapter{name: "adapter-3"}

	core.RegisterAdapter(a1)
	core.RegisterAdapter(a2)
	core.RegisterAdapter(a3)

	core.Start()

	status := core.Status()
	if status.Servers != 3 {
		t.Errorf("expected 3 servers, got %d", status.Servers)
	}

	if !a1.started || !a2.started || !a3.started {
		t.Error("all adapters should be started")
	}

	core.Stop()

	if !a1.stopped || !a2.stopped || !a3.stopped {
		t.Error("all adapters should be stopped")
	}
}
