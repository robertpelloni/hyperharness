package cache

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestCacheSetGet(t *testing.T) {
	c := New(CacheOptions{MaxSize: 10, DefaultTTL: 60000})
	c.Set("key1", "value1")
	val, ok := c.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("expected value1, got %v, ok=%v", val, ok)
	}
}

func TestCacheMiss(t *testing.T) {
	c := New(CacheOptions{})
	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("should be a miss")
	}
}

func TestCacheTTLExpiry(t *testing.T) {
	c := New(CacheOptions{DefaultTTL: 50}) // 50ms
	c.Set("short", "data")
	time.Sleep(80 * time.Millisecond)
	_, ok := c.Get("short")
	if ok {
		t.Error("should have expired")
	}
}

func TestCacheCustomTTL(t *testing.T) {
	c := New(CacheOptions{DefaultTTL: 60000})
	c.SetTTL("short", "data", 30)
	time.Sleep(50 * time.Millisecond)
	_, ok := c.Get("short")
	if ok {
		t.Error("should have expired")
	}
}

func TestCacheLRUEviction(t *testing.T) {
	c := New(CacheOptions{MaxSize: 3, DefaultTTL: 60000})
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Set("d", 4) // should evict one entry

	size, _ := c.Stats()
	if size != 3 {
		t.Errorf("should have 3 entries after eviction, got %d", size)
	}

	_, okD := c.Get("d")
	if !okD {
		t.Error("d should exist (just added)")
	}
}

func TestCacheDelete(t *testing.T) {
	c := New(CacheOptions{})
	c.Set("key", "val")
	if !c.Delete("key") {
		t.Error("should have existed")
	}
	if c.Delete("key") {
		t.Error("should no longer exist")
	}
}

func TestCacheHas(t *testing.T) {
	c := New(CacheOptions{})
	c.Set("exists", true)
	if !c.Has("exists") {
		t.Error("should have key")
	}
	if c.Has("missing") {
		t.Error("should not have key")
	}
}

func TestCacheClear(t *testing.T) {
	c := New(CacheOptions{})
	c.Set("a", 1)
	c.Set("b", 2)
	c.Clear()
	if size, _ := c.Stats(); size != 0 {
		t.Errorf("should be empty, got %d", size)
	}
}

func TestCacheEvents(t *testing.T) {
	var events []CacheEvent
	c := New(CacheOptions{})
	c.OnEvent(func(e CacheEvent) { events = append(events, e) })

	c.Set("k", "v")
	c.Get("k")
	c.Get("missing")
	c.Delete("k")

	if len(events) != 4 {
		t.Fatalf("expected 4 events, got %d", len(events))
	}
	if events[0].Type != "set" { t.Error() }
	if events[1].Type != "hit" { t.Error() }
	if events[2].Type != "miss" { t.Error() }
	if events[3].Type != "delete" { t.Error() }
}

func TestCacheSingleton(t *testing.T) {
	c1 := GetInstance("test-ns")
	c2 := GetInstance("test-ns")
	if c1 != c2 {
		t.Error("should be same instance")
	}
}

func TestCached(t *testing.T) {
	c := New(CacheOptions{DefaultTTL: 60000})
	var calls int32
	fn := func() (interface{}, error) {
		atomic.AddInt32(&calls, 1)
		return "computed", nil
	}

	v1, _ := Cached(c, "key", fn, 0)
	v2, _ := Cached(c, "key", fn, 0)
	if v1 != "computed" || v2 != "computed" {
		t.Error("values mismatch")
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Error("should only call fn once")
	}
}
