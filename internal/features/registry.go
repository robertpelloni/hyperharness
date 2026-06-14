package features

import (
	"fmt"
	"sort"
	"sync"
)

// FeatureInitFn is a function that performs any required initialization for a feature.
// It returns an error if initialization fails.
type FeatureInitFn func() error

var (
	registryMu sync.RWMutex
	registry   = make(map[string]FeatureInitFn)
)

// Register adds a feature initialization function to the global registry.
// It must be called from an init() block of the feature implementation or manually
// during startup. Register panics if the name is empty or already registered.
func Register(name string, initFn FeatureInitFn) {
	if name == "" {
		panic("features: Register called with empty name")
	}
	if initFn == nil {
		panic("features: Register called with nil init function for " + name)
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, exists := registry[name]; exists {
		panic("features: duplicate registration for " + name)
	}
	registry[name] = initFn
}

// Get returns the init function for a registered feature, if present.
func Get(name string) (FeatureInitFn, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	fn, ok := registry[name]
	return fn, ok
}

// List returns a sorted slice of registered feature names.
func List() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	// Alphabetical sort for deterministic ordering.
	sort.Strings(names)
	return names
}

// InitAll runs the init function for every registered feature.
// It returns the first error encountered, if any.
func InitAll() error {
	registryMu.RLock()
	defer registryMu.RUnlock()
	for name, fn := range registry {
		if err := fn(); err != nil {
			return fmt.Errorf("feature %s init failed: %w", name, err)
		}
	}
	return nil
}

// Note: callers can use List() to discover available features.
