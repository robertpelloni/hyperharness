package tui

// ═══════════════════════════════════════════════════════════════════════
// adapter.go — Bridge between tools.Registry and agents.ToolExecutor
// This file exists to break the import cycle: tools <-> agents
// ═══════════════════════════════════════════════════════════════════════

import (
	"github.com/robertpelloni/hyperharness/agents"
	"github.com/robertpelloni/hyperharness/tools"
)

// RegistryAdapter wraps a tools.Registry to satisfy agents.ToolExecutor
type RegistryAdapter struct {
	Reg *tools.Registry
}

func NewRegistryAdapter(reg *tools.Registry) *RegistryAdapter {
	return &RegistryAdapter{Reg: reg}
}

func (a *RegistryAdapter) Find(name string) (agents.ToolInfo, bool) {
	t, ok := a.Reg.Find(name)
	if !ok {
		return agents.ToolInfo{}, false
	}
	return agents.ToolInfo{
		Name:        t.Name,
		Description: t.Description,
		Execute:     t.Execute,
	}, true
}

func (a *RegistryAdapter) ListTools() []agents.ToolInfo {
	result := make([]agents.ToolInfo, 0, len(a.Reg.Tools))
	for _, t := range a.Reg.Tools {
		result = append(result, agents.ToolInfo{
			Name:        t.Name,
			Description: t.Description,
			Execute:     t.Execute,
		})
	}
	return result
}
