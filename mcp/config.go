package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config binds the ~/.borg/mcp.json native parsing
type Config struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

type ServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

// ReadScopedClient strictly constructs OS Subprocess environment lists without pulling globally
func ParseMetadataContext() (*Config, error) {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".borg", "mcp.json")
	
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("missing .borg/mcp.json definition: %w", err)
	}
	
	var conf Config
	if err := json.Unmarshal(b, &conf); err != nil {
		return nil, err
	}
	
	return &conf, nil
}

// FlattenEnv constructs the isolated subprocess environmental bindings
func (s *ServerConfig) FlattenEnv() []string {
	var envList []string
	
	// Start with host required bindings silently merged
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "PATH=") || strings.HasPrefix(e, "NODE_ENV=") {
			envList = append(envList, e)
		}
	}

	for k, v := range s.Env {
		envList = append(envList, fmt.Sprintf("%s=%s", k, v))
	}
	return envList
}
