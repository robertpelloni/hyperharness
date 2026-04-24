package eventbus

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestEventBusExactMatch(t *testing.T) {
	bus := New(100)
	var received SystemEvent
	bus.Subscribe("tool:call", func(e SystemEvent) { received = e })

	bus.EmitEvent("tool:call", "test", map[string]string{"tool": "bash"})
	if received.Source != "test" {
		t.Error("should receive event")
	}
}

func TestEventBusWildcard(t *testing.T) {
	bus := New(100)
	var count int32
	bus.Subscribe("tool:*", func(e SystemEvent) { atomic.AddInt32(&count, 1) })

	bus.EmitEvent("tool:call", "test", nil)
	bus.EmitEvent("tool:result", "test", nil)
	bus.EmitEvent("agent:start", "test", nil) // should NOT match

	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&count) != 2 {
		t.Errorf("expected 2 wildcard matches, got %d", count)
	}
}

func TestEventBusGlobal(t *testing.T) {
	bus := New(100)
	var count int32
	bus.OnGlobal(func(e SystemEvent) { atomic.AddInt32(&count, 1) })

	bus.EmitEvent("tool:call", "test", nil)
	bus.EmitEvent("agent:start", "test", nil)
	bus.EmitEvent("file:change", "test", nil)

	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&count) != 3 {
		t.Errorf("expected 3 global events, got %d", count)
	}
}

func TestEventBusHistory(t *testing.T) {
	bus := New(5)
	for i := 0; i < 10; i++ {
		bus.EmitEvent("tool:call", "test", i)
	}
	history := bus.GetHistory(5)
	if len(history) != 5 {
		t.Errorf("expected 5 history items, got %d", len(history))
	}
}

func TestEventBusNoMatch(t *testing.T) {
	bus := New(100)
	bus.Subscribe("tool:call", func(e SystemEvent) {
		t.Error("should not be called")
	})
	bus.EmitEvent("agent:start", "test", nil)
}

func TestEventBusListenerCount(t *testing.T) {
	bus := New(100)
	bus.Subscribe("a", func(_ SystemEvent) {})
	bus.Subscribe("b", func(_ SystemEvent) {})
	bus.Subscribe("c:*", func(_ SystemEvent) {})
	bus.OnGlobal(func(_ SystemEvent) {})

	if bus.ListenerCount() != 4 {
		t.Errorf("expected 4 listeners, got %d", bus.ListenerCount())
	}
}

func TestEventBusSafeCall(t *testing.T) {
	bus := New(100)
	bus.Subscribe("crash", func(e SystemEvent) { panic("oops") })
	// Should not panic
	bus.EmitEvent("crash", "test", nil)
}
