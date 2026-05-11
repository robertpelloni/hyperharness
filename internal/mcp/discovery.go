package mcp

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DiscoveredServer struct {
	Name    string
	Command string
	Args    []string
	Env     map[string]string
}

// DiscoverLocalServers scans common directories for installed MCP servers.
// This includes Claude desktop config, VS Code global storage, or NPM globals.
func DiscoverLocalServers() ([]DiscoveredServer, error) {
	var servers []DiscoveredServer
	
	// Check Claude Desktop config
	homeDir, _ := os.UserHomeDir()
	
	// Mac path
	macClaudePath := filepath.Join(homeDir, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	if serversFromClaude, err := parseClaudeConfig(macClaudePath); err == nil {
		servers = append(servers, serversFromClaude...)
	}

	// Windows path
	winClaudePath := filepath.Join(os.Getenv("APPDATA"), "Claude", "claude_desktop_config.json")
	if serversFromClaude, err := parseClaudeConfig(winClaudePath); err == nil {
		servers = append(servers, serversFromClaude...)
	}

	return servers, nil
}

func parseClaudeConfig(path string) ([]DiscoveredServer, error) {
	var servers []DiscoveredServer
	
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	type ClaudeConfig struct {
		MCPServers map[string]struct {
			Command string            `json:"command"`
			Args    []string          `json:"args"`
			Env     map[string]string `json:"env"`
		} `json:"mcpServers"`
	}
	
	var config ClaudeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	
	for name, srv := range config.MCPServers {
		servers = append(servers, DiscoveredServer{
			Name:    name,
			Command: srv.Command,
			Args:    srv.Args,
			Env:     srv.Env,
		})
	}
	
	return servers, nil
}

// AutoRegisterDiscoveredServers discovers local servers and registers them
// into the running HyperHarness MCP Manager.
func (r *Registry) AutoRegisterDiscoveredServers() error {
	servers, _ := DiscoverLocalServers()
	
	r.mu.Lock()
	for _, srv := range servers {
		if _, exists := r.servers[srv.Name]; !exists {
			mcpSrv := &MCPServer{
				Name:      srv.Name,
				Command:   srv.Command,
				Args:      srv.Args,
				Env:       srv.Env,
				Transport: "stdio",
			}
			r.servers[srv.Name] = mcpSrv
			// Attempt to autoconnect in background
			go func(name string) {
				_ = r.Connect(name)
			}(srv.Name)
		}
	}
	r.mu.Unlock()
	return nil
}
