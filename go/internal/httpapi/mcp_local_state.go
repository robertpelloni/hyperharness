package httpapi

import (
	"sort"
	"sync"
	"time"
)

type localMCPTool struct {
	Name        string
	Description string
	InputSchema any
	AlwaysOn    bool
	ServerName  string
	Source      string
}

type localMCPToolState struct {
	LoadedAt     int64
	AccessedAt   int64
	HydratedAt   *int64
	AlwaysLoaded bool
}

type localMCPTelemetryEvent struct {
	Type                   string   `json:"type"`
	ToolName               string   `json:"toolName,omitempty"`
	Source                 string   `json:"source,omitempty"`
	Status                 string   `json:"status,omitempty"`
	Message                string   `json:"message,omitempty"`
	EvictedTools           []string `json:"evictedTools,omitempty"`
	LoadedToolCount        int      `json:"loadedToolCount,omitempty"`
	HydratedToolCount      int      `json:"hydratedToolCount,omitempty"`
	LoadedUtilizationPct   float64  `json:"loadedUtilizationPct,omitempty"`
	HydratedUtilizationPct float64  `json:"hydratedUtilizationPct,omitempty"`
	Timestamp              int64    `json:"timestamp"`
}

type localMCPEvictionEvent struct {
	ToolName       string `json:"toolName"`
	Timestamp      int64  `json:"timestamp"`
	Tier           string `json:"tier"`
	IdleEvicted    bool   `json:"idleEvicted"`
	IdleDurationMs int64  `json:"idleDurationMs"`
}

type localMCPStateManager struct {
	mu         sync.RWMutex
	loaded     map[string]*localMCPToolState
	telemetry  []localMCPTelemetryEvent
	evictions  []localMCPEvictionEvent
	maxHistory int
}

func newLocalMCPStateManager() *localMCPStateManager {
	return &localMCPStateManager{
		loaded:     map[string]*localMCPToolState{},
		telemetry:  []localMCPTelemetryEvent{},
		evictions:  []localMCPEvictionEvent{},
		maxHistory: 200,
	}
}

func (m *localMCPStateManager) snapshot(limits map[string]any, available map[string]localMCPTool) map[string]any {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.applyAlwaysLoadedLocked(available)
	tools := make([]map[string]any, 0, len(m.loaded))
	for name, state := range m.loaded {
		tool, ok := available[name]
		if !ok {
			continue
		}
		tools = append(tools, map[string]any{
			"name":           name,
			"description":    tool.Description,
			"serverName":     tool.ServerName,
			"loaded":         true,
			"hydrated":       state.HydratedAt != nil,
			"alwaysOn":       state.AlwaysLoaded || tool.AlwaysOn,
			"lastLoadedAt":   state.LoadedAt,
			"lastAccessedAt": state.AccessedAt,
			"lastHydratedAt": nullableInt64Value(state.HydratedAt),
		})
	}
	sort.Slice(tools, func(i, j int) bool {
		left := stringValue(tools[i]["name"])
		right := stringValue(tools[j]["name"])
		return left < right
	})
	return map[string]any{
		"limits": limits,
		"tools":  tools,
	}
}

func (m *localMCPStateManager) loadTool(name string, limits map[string]any, available map[string]localMCPTool) (string, []string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.applyAlwaysLoadedLocked(available)
	tool, ok := available[name]
	if !ok {
		return "", nil, false
	}
	now := time.Now().UTC().UnixMilli()
	state, exists := m.loaded[name]
	if !exists {
		state = &localMCPToolState{LoadedAt: now, AccessedAt: now, AlwaysLoaded: tool.AlwaysOn}
		m.loaded[name] = state
	} else {
		state.AccessedAt = now
	}
	evicted := m.enforceLoadedCapacityLocked(intNumber(limits["maxLoadedTools"]))
	return "loaded", evicted, true
}

func (m *localMCPStateManager) unloadTool(name string, available map[string]localMCPTool) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if tool, ok := available[name]; ok && tool.AlwaysOn {
		if state, exists := m.loaded[name]; exists {
			state.HydratedAt = nil
			state.AccessedAt = time.Now().UTC().UnixMilli()
		}
		return "tool remains loaded because it is configured as always-loaded; cleared hydrated schema state only", true
	}
	if _, ok := m.loaded[name]; !ok {
		return "", false
	}
	delete(m.loaded, name)
	return "unloaded", true
}

func (m *localMCPStateManager) hydrateTool(name string, limits map[string]any, available map[string]localMCPTool) (map[string]any, []string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.applyAlwaysLoadedLocked(available)
	tool, ok := available[name]
	if !ok {
		return nil, nil, false
	}
	now := time.Now().UTC().UnixMilli()
	state, exists := m.loaded[name]
	if !exists {
		state = &localMCPToolState{LoadedAt: now, AccessedAt: now, AlwaysLoaded: tool.AlwaysOn}
		m.loaded[name] = state
	}
	state.AccessedAt = now
	state.HydratedAt = &now
	evicted := m.enforceHydratedCapacityLocked(intNumber(limits["maxHydratedSchemas"]))
	return map[string]any{"inputSchema": tool.InputSchema, "evictedHydratedTools": evicted}, evicted, true
}

func (m *localMCPStateManager) telemetryList() []map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]map[string]any, 0, len(m.telemetry))
	for _, event := range m.telemetry {
		result = append(result, map[string]any{
			"type":                   event.Type,
			"toolName":               nullableString(event.ToolName),
			"source":                 nullableString(event.Source),
			"status":                 nullableString(event.Status),
			"message":                nullableString(event.Message),
			"evictedTools":           event.EvictedTools,
			"loadedToolCount":        event.LoadedToolCount,
			"hydratedToolCount":      event.HydratedToolCount,
			"loadedUtilizationPct":   event.LoadedUtilizationPct,
			"hydratedUtilizationPct": event.HydratedUtilizationPct,
			"timestamp":              event.Timestamp,
		})
	}
	return result
}

func (m *localMCPStateManager) clearTelemetry() map[string]any {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := len(m.telemetry)
	m.telemetry = nil
	return map[string]any{"ok": true, "cleared": count}
}

func (m *localMCPStateManager) evictionList() []map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]map[string]any, 0, len(m.evictions))
	for _, event := range m.evictions {
		result = append(result, map[string]any{
			"toolName":       event.ToolName,
			"timestamp":      event.Timestamp,
			"tier":           event.Tier,
			"idleEvicted":    event.IdleEvicted,
			"idleDurationMs": event.IdleDurationMs,
		})
	}
	return result
}

func (m *localMCPStateManager) clearEvictions() map[string]any {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := len(m.evictions)
	m.evictions = nil
	return map[string]any{"ok": true, "message": "cleared", "cleared": count}
}

func (m *localMCPStateManager) recordTelemetry(event localMCPTelemetryEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.telemetry = append([]localMCPTelemetryEvent{event}, m.telemetry...)
	if len(m.telemetry) > m.maxHistory {
		m.telemetry = m.telemetry[:m.maxHistory]
	}
}

func (m *localMCPStateManager) applyAlwaysLoadedLocked(available map[string]localMCPTool) {
	now := time.Now().UTC().UnixMilli()
	for name, tool := range available {
		if !tool.AlwaysOn {
			continue
		}
		state, ok := m.loaded[name]
		if !ok {
			m.loaded[name] = &localMCPToolState{LoadedAt: now, AccessedAt: now, AlwaysLoaded: true}
			continue
		}
		state.AlwaysLoaded = true
	}
}

func (m *localMCPStateManager) enforceLoadedCapacityLocked(limit int) []string {
	if limit <= 0 {
		return nil
	}
	names := m.loadedNamesLocked(false)
	evicted := []string{}
	for len(names) > limit {
		victim := names[len(names)-1]
		delete(m.loaded, victim)
		evicted = append(evicted, victim)
		m.evictions = append([]localMCPEvictionEvent{{ToolName: victim, Timestamp: time.Now().UTC().UnixMilli(), Tier: "loaded", IdleEvicted: false, IdleDurationMs: 0}}, m.evictions...)
		names = m.loadedNamesLocked(false)
	}
	if len(m.evictions) > m.maxHistory {
		m.evictions = m.evictions[:m.maxHistory]
	}
	return evicted
}

func (m *localMCPStateManager) enforceHydratedCapacityLocked(limit int) []string {
	if limit <= 0 {
		return nil
	}
	hydrated := make([]string, 0)
	for name, state := range m.loaded {
		if state.HydratedAt != nil {
			hydrated = append(hydrated, name)
		}
	}
	sort.Slice(hydrated, func(i, j int) bool {
		left := m.loaded[hydrated[i]]
		right := m.loaded[hydrated[j]]
		leftTS := int64(0)
		rightTS := int64(0)
		if left != nil && left.HydratedAt != nil {
			leftTS = *left.HydratedAt
		}
		if right != nil && right.HydratedAt != nil {
			rightTS = *right.HydratedAt
		}
		return leftTS > rightTS
	})
	evicted := []string{}
	for len(hydrated) > limit {
		victim := hydrated[len(hydrated)-1]
		if state := m.loaded[victim]; state != nil {
			state.HydratedAt = nil
		}
		evicted = append(evicted, victim)
		m.evictions = append([]localMCPEvictionEvent{{ToolName: victim, Timestamp: time.Now().UTC().UnixMilli(), Tier: "hydrated", IdleEvicted: false, IdleDurationMs: 0}}, m.evictions...)
		hydrated = hydrated[:len(hydrated)-1]
	}
	if len(m.evictions) > m.maxHistory {
		m.evictions = m.evictions[:m.maxHistory]
	}
	return evicted
}

func (m *localMCPStateManager) loadedNamesLocked(includeAlwaysOnly bool) []string {
	names := make([]string, 0, len(m.loaded))
	for name, state := range m.loaded {
		if state == nil {
			continue
		}
		if includeAlwaysOnly && !state.AlwaysLoaded {
			continue
		}
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		left := m.loaded[names[i]]
		right := m.loaded[names[j]]
		if left.AlwaysLoaded != right.AlwaysLoaded {
			return left.AlwaysLoaded
		}
		return left.AccessedAt > right.AccessedAt
	})
	return names
}

func nullableInt64Value(value *int64) any {
	if value == nil {
		return nil
	}
	return *value
}
