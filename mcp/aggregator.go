// Package mcp provides MCP tool aggregation from multiple connected servers.
// Ported from hypercode/go/internal/mcp/aggregator.go with enhancements.
//
// WHAT: Aggregates tools from multiple MCP servers into a unified namespace
// WHY: Enables hyperharness to route tool calls to the correct MCP server
// HOW: Maintains a map of server name → StdioClient, broadcasts discovery,
//      routes calls by prefixed tool names (server__tool)
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

// AggregatedTool is a tool discovered from a connected MCP server.
type AggregatedTool struct {
	Name         string      `json:"name"`         // Fully qualified: server__tool
	OriginalName string      `json:"originalName"` // Tool name on the server
	Description  string      `json:"description"`
	Server       string      `json:"server"`       // Server that provides this tool
	InputSchema  interface{} `json:"inputSchema"`
}

// Aggregator manages connections to multiple MCP servers and provides
// unified tool discovery and routing.
type Aggregator struct {
	clients map[string]*StdioClient
	mu      sync.RWMutex
}

// NewAggregator creates a new MCP aggregator.
func NewAggregator() *Aggregator {
	return &Aggregator{
		clients: make(map[string]*StdioClient),
	}
}

// AddServer spawns and connects to an MCP server.
// If the server is already connected, this is a no-op.
func (a *Aggregator) AddServer(name, command string, args []string, env map[string]string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.clients[name]; ok {
		return nil // Already exists
	}

	client := NewStdioClient(name, command, args, env)
	if err := client.Start(); err != nil {
		return fmt.Errorf("failed to start MCP server %s: %w", name, err)
	}

	// Initialize the connection (discover tools, resources, prompts)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := client.Initialize(ctx); err != nil {
		client.Stop()
		return fmt.Errorf("failed to initialize MCP server %s: %w", name, err)
	}

	a.clients[name] = client
	return nil
}

// AddClient registers an already-connected client.
func (a *Aggregator) AddClient(name string, client *StdioClient) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.clients[name] = client
}

// RemoveServer disconnects and removes a server.
func (a *Aggregator) RemoveServer(name string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if client, ok := a.clients[name]; ok {
		client.Stop()
		delete(a.clients, name)
	}
}

// ListTools aggregates tools from all connected servers.
// Tools are namespaced as "server__tool" to avoid collisions.
func (a *Aggregator) ListTools(ctx context.Context) ([]AggregatedTool, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var allTools []AggregatedTool

	for serverName, client := range a.clients {
		for _, tool := range client.Tools {
			allTools = append(allTools, AggregatedTool{
				Name:         serverName + "__" + tool.Name,
				OriginalName: tool.Name,
				Description:  tool.Description,
				Server:       serverName,
				InputSchema:  tool.InputSchema,
			})
		}
	}

	// Sort by qualified name
	sort.Slice(allTools, func(i, j int) bool {
		return allTools[i].Name < allTools[j].Name
	})

	return allTools, nil
}

// ListResources aggregates resources from all connected servers.
func (a *Aggregator) ListResources(ctx context.Context) ([]MCPResource, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var allResources []MCPResource
	for _, client := range a.clients {
		allResources = append(allResources, client.Resources...)
	}
	return allResources, nil
}

// ListPrompts aggregates prompts from all connected servers.
func (a *Aggregator) ListPrompts(ctx context.Context) ([]MCPPrompt, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var allPrompts []MCPPrompt
	for _, client := range a.clients {
		allPrompts = append(allPrompts, client.Prompts...)
	}
	return allPrompts, nil
}

// CallTool routes a tool call to the correct server.
// Supports both "server__tool" and plain "tool" naming.
// For plain names, searches all servers in alphabetical order.
func (a *Aggregator) CallTool(ctx context.Context, toolName string, arguments map[string]interface{}) (interface{}, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Try namespaced name first: "server__tool"
	for serverName, client := range a.clients {
		prefix := serverName + "__"
		if len(toolName) > len(prefix) && toolName[:len(prefix)] == prefix {
			originalName := toolName[len(prefix):]
			return client.CallTool(ctx, originalName, arguments)
		}
	}

	// Fallback: search all servers for the tool
	var servers []string
	for name := range a.clients {
		servers = append(servers, name)
	}
	sort.Strings(servers)

	for _, serverName := range servers {
		client := a.clients[serverName]
		for _, tool := range client.Tools {
			if tool.Name == toolName {
				return client.CallTool(ctx, toolName, arguments)
			}
		}
	}

	return nil, fmt.Errorf("tool %q not found on any connected MCP server", toolName)
}

// HasTool checks if a tool is available on any connected server.
func (a *Aggregator) HasTool(toolName string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, client := range a.clients {
		for _, tool := range client.Tools {
			if tool.Name == toolName || client.Name+"__"+tool.Name == toolName {
				return true
			}
		}
	}
	return false
}

// GetServer returns a specific client by name.
func (a *Aggregator) GetServer(name string) (*StdioClient, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	client, ok := a.clients[name]
	return client, ok
}

// ServerNames returns the names of all connected servers.
func (a *Aggregator) ServerNames() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	names := make([]string, 0, len(a.clients))
	for name := range a.clients {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ServerCount returns the number of connected servers.
func (a *Aggregator) ServerCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.clients)
}

// TotalToolCount returns the total number of tools across all servers.
func (a *Aggregator) TotalToolCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	count := 0
	for _, client := range a.clients {
		count += len(client.Tools)
	}
	return count
}

// Shutdown disconnects all servers.
func (a *Aggregator) Shutdown() {
	a.mu.Lock()
	defer a.mu.Unlock()

	for name, client := range a.clients {
		client.Stop()
		delete(a.clients, name)
	}
}

// Status returns a summary of the aggregator state.
func (a *Aggregator) Status() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	servers := make(map[string]interface{})
	for name, client := range a.clients {
		servers[name] = map[string]interface{}{
			"running":    client.IsRunning(),
			"toolCount":  client.ToolCount(),
			"serverInfo": client.Info,
		}
	}

	return map[string]interface{}{
		"serverCount":  len(a.clients),
		"totalTools":   a.TotalToolCount(),
		"servers":      servers,
	}
}

// MarshalJSON implements custom JSON serialization for status.
func (a *Aggregator) StatusJSON() ([]byte, error) {
	return json.MarshalIndent(a.Status(), "", "  ")
}
