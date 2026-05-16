// Package cache provides a TTL + LRU in-memory cache.
// Ported from hypercode/go/internal/cache/cache.go
package cache

import (
	"sync"
	"time"
)

// CacheOptions configures the cache.
type CacheOptions struct {
	MaxSize    int
	DefaultTTL time.Duration
	Evict      func(key string, value any)
}

// Cache is a thread-safe in-memory cache with TTL and LRU eviction.
type Cache struct {
	mu       sync.RWMutex
	items    map[string]*cacheItem
	maxSize  int
	defaultTTL time.Duration
	onEvict  func(key string, value any)
	lru      []string
}

// cacheItem holds a cached value with expiration.
type cacheItem struct {
	Value     any
	ExpiresAt time.Time
}

const defaultMaxSize = 1000

// New creates a new cache with the given options.
func New(opts CacheOptions) *Cache {
	maxSize := opts.MaxSize
	if maxSize <= 0 {
		maxSize = defaultMaxSize
	}
	return &Cache{
		items:      make(map[string]*cacheItem),
		maxSize:    maxSize,
		defaultTTL: opts.DefaultTTL,
		onEvict:   opts.Evict,
	}
}

// Get retrieves a value by key.
func (c *Cache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt) {
		c.evictLocked(key)
		return nil, false
	}
	c.touchLocked(key)
	return item.Value, true
}

// Set stores a value with optional TTL.
func (c *Cache) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ttl <= 0 {
		ttl = c.defaultTTL
	}
	c.items[key] = &cacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
	c.touchLocked(key)
	if len(c.items) > c.maxSize {
		c.evictLRULocked()
	}
}

// Delete removes a key.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evictLocked(key)
}

// Clear removes all items.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*cacheItem)
	c.lru = nil
}

// Len returns the number of items.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Has reports whether key exists and is not expired.
func (c *Cache) Has(key string) bool {
	_, ok := c.Get(key)
	return ok
}

func (c *Cache) touchLocked(key string) {
	for i, k := range c.lru {
		if k == key {
			c.lru = append(c.lru[:i], c.lru[i+1:]...)
			break
		}
	}
	c.lru = append(c.lru, key)
}

func (c *Cache) evictLocked(key string) {
	if item, ok := c.items[key]; ok {
		if c.onEvict != nil {
			c.onEvict(key, item.Value)
		}
		delete(c.items, key)
	}
	for i, k := range c.lru {
		if k == key {
			c.lru = append(c.lru[:i], c.lru[i+1:]...)
			break
		}
	}
}

func (c *Cache) evictLRULocked() {
	if len(c.lru) > 0 {
		c.evictLocked(c.lru[0])
	}
}