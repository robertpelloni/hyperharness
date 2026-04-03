package mcp

import (
	"context"
	"fmt"
	"sync"
)

type Aggregator struct {
	clients map[string]*StdioClient
	mu      sync.RWMutex
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		clients: make(map[string]*StdioClient),
	}
}

func (a *Aggregator) AddServer(name string, command string, args []string, env map[string]string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.clients[name]; ok {
		return nil // Already exists
	}

	client := NewStdioClient(name, command, args, env)
	if err := client.Start(); err != nil {
		return err
	}

	a.clients[name] = client
	return nil
}

func (a *Aggregator) ListTools(ctx context.Context) ([]ToolEntry, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var allTools []ToolEntry
	for name, client := range a.clients {
		resp, err := client.Call(ctx, "tools/list", nil)
		if err != nil {
			fmt.Printf("Error listing tools for %s: %v\n", name, err)
			continue
		}

		if resp.Result != nil {
			// Parse result and add to allTools
			// ...
		}
	}

	return allTools, nil
}
