package mcp

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

type ServerEntry struct {
	UUID               string            `json:"uuid"`
	Name               string            `json:"name"`
	DisplayName        string            `json:"displayName"`
	Type               string            `json:"type"`
	Command            string            `json:"command"`
	Args               []string          `json:"args"`
	Env                map[string]string `json:"env"`
	URL                string            `json:"url"`
	Description        string            `json:"description"`
	Enabled            bool              `json:"enabled"`
	AlwaysOn           bool              `json:"alwaysOn"`
	Tags               []string          `json:"tags"`
	AlwaysOnAdvertised bool              `json:"alwaysOnAdvertised"`
}

type ToolEntry struct {
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	Server             string      `json:"server"`
	ServerDisplayName  string      `json:"serverDisplayName"`
	ServerTags         []string    `json:"serverTags"`
	ToolTags           []string    `json:"toolTags"`
	SemanticGroup      string      `json:"semanticGroup"`
	SemanticGroupLabel string      `json:"semanticGroupLabel"`
	AdvertisedName     string      `json:"advertisedName"`
	Keywords           []string    `json:"keywords"`
	AlwaysOn           bool        `json:"alwaysOn"`
	OriginalName       string      `json:"originalName"`
	InputSchema        interface{} `json:"inputSchema"`
}

type Inventory struct {
	Servers []ServerEntry `json:"servers"`
	Tools   []ToolEntry   `json:"tools"`
	Source  string        `json:"source"`
}

func LoadInventory(workspaceRoot, mainConfigDir string) (*Inventory, error) {
	inventory := &Inventory{
		Servers: []ServerEntry{},
		Tools:   []ToolEntry{},
		Source:  "empty",
	}

	// 1. Load from mcp.jsonc/json
	configServers, err := loadConfigServers(mainConfigDir)
	if err == nil && len(configServers) > 0 {
		for name, server := range configServers {
			sEntry := ServerEntry{
				UUID:        "config:" + name,
				Name:        name,
				DisplayName: name,
				Type:        "STDIO",
				Enabled:     true,
			}
			if server.URL != "" {
				sEntry.URL = server.URL
				sEntry.Type = "SSE"
			}
			if server.Command != "" {
				sEntry.Command = server.Command
			}
			sEntry.Args = server.Args
			sEntry.Env = server.Env
			
			inventory.Servers = append(inventory.Servers, sEntry)
			
			for _, tool := range server.Meta.Tools {
				inventory.Tools = append(inventory.Tools, ToolEntry{
					Name:              name + "__" + tool.Name,
					Description:       tool.Description,
					Server:            name,
					ServerDisplayName: name,
					AdvertisedName:    name + "__" + tool.Name,
					OriginalName:      tool.Name,
					InputSchema:       tool.InputSchema,
				})
			}
		}
		inventory.Source = "config"
	}

	// 2. Load from Database
	dbPath := filepath.Join(workspaceRoot, "packages", "core", "metamcp.db")
	db, err := sql.Open("sqlite", dbPath)
	if err == nil {
		defer db.Close()
		
		rows, err := db.Query("SELECT uuid, name, type, command, args, env, url, description, enabled, always_on FROM mcp_servers")
		if err == nil {
			for rows.Next() {
				var s ServerEntry
				var argsRaw, envRaw []byte
				var urlOpt, descOpt sql.NullString
				err := rows.Scan(&s.UUID, &s.Name, &s.Type, &s.Command, &argsRaw, &envRaw, &urlOpt, &descOpt, &s.Enabled, &s.AlwaysOn)
				if err == nil {
					json.Unmarshal(argsRaw, &s.Args)
					json.Unmarshal(envRaw, &s.Env)
					s.URL = urlOpt.String
					s.Description = descOpt.String
					s.DisplayName = s.Name
					s.AlwaysOnAdvertised = s.AlwaysOn
					inventory.Servers = append(inventory.Servers, s)
				}
			}
			rows.Close()
		}

		tRows, err := db.Query("SELECT name, description, mcp_server_uuid, always_on, tool_schema FROM tools")
		if err == nil {
			serverMap := make(map[string]string)
			for _, s := range inventory.Servers {
				serverMap[s.UUID] = s.Name
			}

			for tRows.Next() {
				var t ToolEntry
				var serverUUID string
				var schemaRaw []byte
				err := tRows.Scan(&t.OriginalName, &t.Description, &serverUUID, &t.AlwaysOn, &schemaRaw)
				if err == nil {
					serverName := serverMap[serverUUID]
					if serverName == "" {
						serverName = "unknown"
					}
					t.Server = serverName
					t.ServerDisplayName = serverName
					t.Name = serverName + "__" + t.OriginalName
					t.AdvertisedName = t.Name
					json.Unmarshal(schemaRaw, &t.InputSchema)
					inventory.Tools = append(inventory.Tools, t)
				}
			}
			tRows.Close()
		}
		
		if len(inventory.Tools) > 0 {
			inventory.Source = "database"
		}
	}

	return inventory, nil
}

type configServer struct {
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	Env         map[string]string `json:"env"`
	URL         string            `json:"url"`
	Description string            `json:"description"`
	Disabled    bool              `json:"disabled"`
	Meta        struct {
		Tools []struct {
			Name        string      `json:"name"`
			Description string      `json:"description"`
			InputSchema interface{} `json:"inputSchema"`
		} `json:"tools"`
	} `json:"_meta"`
}

type mcpConfig struct {
	McpServers map[string]configServer `json:"mcpServers"`
}

func loadConfigServers(configDir string) (map[string]configServer, error) {
	path := filepath.Join(configDir, "mcp.jsonc")
	data, err := os.ReadFile(path)
	if err != nil {
		path = filepath.Join(configDir, "mcp.json")
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, err
		}
	}

	// Simple comment stripping
	lines := strings.Split(string(data), "\n")
	var clean []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") {
			continue
		}
		clean = append(clean, line)
	}
	
	var config mcpConfig
	err = json.Unmarshal([]byte(strings.Join(clean, "\n")), &config)
	if err != nil {
		return nil, err
	}

	return config.McpServers, nil
}
