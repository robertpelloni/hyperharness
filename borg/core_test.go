package borg

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

type testSubsystem struct {
	name     string
	initErr  error
	shutErr  error
	inited   bool
	shutdown bool
}

func (t *testSubsystem) Name() string                          { return t.name }
func (t *testSubsystem) Initialize(_ *Core) error              { t.inited = true; return t.initErr }
func (t *testSubsystem) Shutdown() error                       { t.shutdown = true; return t.shutErr }
func (t *testSubsystem) Status() map[string]any                { return map[string]any{"initialized": t.inited} }

func TestNewCore(t *testing.T) {
	core, err := NewCore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if core.State() != StateUninitialized {
		t.Errorf("state: %s", core.State())
	}
	if core.WorkingDir() == "" {
		t.Error("working dir should be set")
	}
}

func TestBootAndShutdown(t *testing.T) {
	core, _ := NewCore(t.TempDir())

	sub := &testSubsystem{name: "test"}
	core.RegisterSubsystem(sub)

	err := core.Boot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if core.State() != StateReady {
		t.Errorf("state after boot: %s", core.State())
	}
	if !sub.inited {
		t.Error("subsystem should be initialized")
	}

	err = core.Shutdown()
	if err != nil {
		t.Fatal(err)
	}
	if !sub.shutdown {
		t.Error("subsystem should be shut down")
	}
}

func TestBootSubsystemError(t *testing.T) {
	core, _ := NewCore(t.TempDir())

	sub := &testSubsystem{name: "failing", initErr: fmt.Errorf("init failed")}
	core.RegisterSubsystem(sub)

	err := core.Boot(context.Background())
	if err == nil {
		t.Error("should fail when subsystem init fails")
	}
	if core.State() != StateError {
		t.Errorf("state should be error: %s", core.State())
	}
}

func TestGetSubsystem(t *testing.T) {
	core, _ := NewCore(t.TempDir())
	sub := &testSubsystem{name: "memory"}
	core.RegisterSubsystem(sub)

	got, ok := core.GetSubsystem("memory")
	if !ok {
		t.Error("should find registered subsystem")
	}
	if got.Name() != "memory" {
		t.Errorf("name: %s", got.Name())
	}

	_, ok = core.GetSubsystem("nonexistent")
	if ok {
		t.Error("should not find unregistered subsystem")
	}
}

func TestEventBus(t *testing.T) {
	core, _ := NewCore(t.TempDir())

	var received atomic.Int32
	core.Subscribe("test:event", func(event Event) {
		if event.Type != "test:event" {
			t.Errorf("event type: %s", event.Type)
		}
		received.Add(1)
	})

	core.Publish(Event{
		Type:   "test:event",
		Source: "test",
	})

	if received.Load() != 1 {
		t.Errorf("expected 1 event, got %d", received.Load())
	}
}

func TestWildcardListener(t *testing.T) {
	core, _ := NewCore(t.TempDir())

	var count atomic.Int32
	core.Subscribe("*", func(event Event) {
		count.Add(1)
	})

	core.Publish(Event{Type: "a", Source: "test"})
	core.Publish(Event{Type: "b", Source: "test"})
	core.Publish(Event{Type: "c", Source: "test"})

	if count.Load() != 3 {
		t.Errorf("wildcard should receive all events: %d", count.Load())
	}
}

func TestStats(t *testing.T) {
	core, _ := NewCore(t.TempDir())
	core.Boot(context.Background())

	stats := core.Stats()
	if stats.Subsystems != 0 {
		t.Errorf("no subsystems: %d", stats.Subsystems)
	}
	if stats.Uptime == "" {
		t.Error("uptime should be set")
	}
}

func TestHealthCheck(t *testing.T) {
	core, _ := NewCore(t.TempDir())
	sub := &testSubsystem{name: "health_test"}
	core.RegisterSubsystem(sub)
	core.Boot(context.Background())

	health := core.HealthCheck()
	if health["state"] != "ready" {
		t.Errorf("health state: %v", health["state"])
	}
}

func TestMemoryStore(t *testing.T) {
	dir := t.TempDir()
	store := NewMemoryStore(dir)

	err := store.Set("test_key", "test_value")
	if err != nil {
		t.Fatal(err)
	}

	val, ok := store.Get("test_key")
	if !ok {
		t.Fatal("should find stored key")
	}
	if val != "test_value" {
		t.Errorf("value: %s", val)
	}
}

func TestMemoryStorePersistence(t *testing.T) {
	dir := t.TempDir()

	// Create and write
	store1 := NewMemoryStore(dir)
	store1.Set("persist", "value1")

	// Create new store and load
	store2 := NewMemoryStore(dir)
	err := store2.Load()
	if err != nil {
		t.Fatal(err)
	}

	val, ok := store2.Get("persist")
	if !ok {
		t.Fatal("should find persisted key")
	}
	if val != "value1" {
		t.Errorf("persisted value: %s", val)
	}
}

func TestMemoryStoreDelete(t *testing.T) {
	dir := t.TempDir()
	store := NewMemoryStore(dir)

	store.Set("temp", "will be deleted")
	store.Delete("temp")

	_, ok := store.Get("temp")
	if ok {
		t.Error("should not find deleted key")
	}
}

func TestMemoryStoreList(t *testing.T) {
	dir := t.TempDir()
	store := NewMemoryStore(dir)

	store.Set("key1", "val1")
	store.Set("key2", "val2")
	store.Set("key3", "val3")

	keys := store.List()
	if len(keys) < 3 {
		t.Errorf("should have at least 3 keys, got %d", len(keys))
	}
}

func TestCoreStateTransitions(t *testing.T) {
	core, _ := NewCore(t.TempDir())

	if core.State() != StateUninitialized {
		t.Error("should start uninitialized")
	}

	core.Boot(context.Background())
	if core.State() != StateReady {
		t.Errorf("should be ready: %s", core.State())
	}

	core.Shutdown()
	if core.State() != StateUninitialized {
		t.Errorf("should be uninitialized after shutdown: %s", core.State())
	}
}

func TestBootWithMultipleSubsystems(t *testing.T) {
	core, _ := NewCore(t.TempDir())

	core.RegisterSubsystem(&testSubsystem{name: "memory"})
	core.RegisterSubsystem(&testSubsystem{name: "context"})
	core.RegisterSubsystem(&testSubsystem{name: "mcp"})

	err := core.Boot(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	stats := core.Stats()
	if stats.Subsystems != 3 {
		t.Errorf("should have 3 subsystems: %d", stats.Subsystems)
	}
}

func TestEventTimestamp(t *testing.T) {
	core, _ := NewCore(t.TempDir())

	var eventTime time.Time
	core.Subscribe("timed", func(e Event) {
		eventTime = e.Timestamp
	})

	core.Publish(Event{Type: "timed", Source: "test"})

	if eventTime.IsZero() {
		t.Error("event should have timestamp")
	}
}
