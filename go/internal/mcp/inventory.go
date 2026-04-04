package mcp

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

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
	Servers  []ServerEntry `json:"servers"`
	Tools    []ToolEntry   `json:"tools"`
	Source   string        `json:"source"`
	CachedAt string        `json:"cachedAt,omitempty"`
}

type persistedInventoryCache struct {
	Version   int       `json:"version"`
	CachedAt  string    `json:"cachedAt"`
	Inventory Inventory `json:"inventory"`
}

func LoadInventory(workspaceRoot, mainConfigDir string) (*Inventory, error) {
	return LoadInventoryWithCache(workspaceRoot, mainConfigDir, "")
}

func LoadInventoryWithCache(workspaceRoot, mainConfigDir, cachePath string) (*Inventory, error) {
	liveInventory, err := loadLiveInventory(workspaceRoot, mainConfigDir)
	if err != nil {
		return nil, err
	}
	if hasInventoryContents(liveInventory) {
		if cachePath != "" {
			persistedAt := time.Now().UTC().Format(time.RFC3339)
			liveInventory.CachedAt = persistedAt
			_ = saveInventoryCache(cachePath, liveInventory, persistedAt)
		}
		return liveInventory, nil
	}
	if cachePath != "" {
		cached, err := loadInventoryCache(cachePath)
		if err == nil && cached != nil {
			return cached, nil
		}
	}
	return liveInventory, nil
}

func SyncInventoryCacheFromLiveSources(workspaceRoot, mainConfigDir, cachePath string) (*Inventory, error) {
	liveInventory, err := loadLiveInventory(workspaceRoot, mainConfigDir)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(cachePath) == "" {
		return liveInventory, nil
	}
	if hasInventoryContents(liveInventory) {
		persistedAt := time.Now().UTC().Format(time.RFC3339)
		liveInventory.CachedAt = persistedAt
		if err := saveInventoryCache(cachePath, liveInventory, persistedAt); err != nil {
			return nil, err
		}
		return liveInventory, nil
	}
	if err := removeInventoryCache(cachePath); err != nil {
		return nil, err
	}
	return liveInventory, nil
}

func loadLiveInventory(workspaceRoot, mainConfigDir string) (*Inventory, error) {
	inventory := &Inventory{
		Servers: []ServerEntry{},
		Tools:   []ToolEntry{},
		Source:  "empty",
	}

	configServers, err := loadConfigServers(mainConfigDir)
	if err == nil && len(configServers) > 0 {
		for name, server := range configServers {
			sEntry := ServerEntry{
				UUID:        "config:" + name,
				Name:        name,
				DisplayName: name,
				Type:        "STDIO",
				Enabled:     !server.Disabled,
			}
			if server.URL != "" {
				sEntry.URL = server.URL
				sEntry.Type = "SSE"
			}
			if server.Command != "" {
				sEntry.Command = server.Command
			}
			sEntry.Args = append([]string(nil), server.Args...)
			sEntry.Env = cloneStringMap(server.Env)
			sEntry.Description = server.Description
			inventory.Servers = append(inventory.Servers, sEntry)

			for _, tool := range server.Meta.Tools {
				inventory.Tools = append(inventory.Tools, ToolEntryFromMetadata(name, MetadataTool{
					Name:        tool.Name,
					Description: tool.Description,
					InputSchema: tool.InputSchema,
					AlwaysOn:    tool.AlwaysOn,
				}))
			}
		}
		inventory.Source = "config"
	}

	dbPath := filepath.Join(workspaceRoot, "packages", "core", "metamcp.db")
	db, err := sql.Open("sqlite", dbPath)
	if err == nil {
		defer db.Close()
		dbToolCountBefore := len(inventory.Tools)

		rows, err := db.Query("SELECT uuid, name, type, command, args, env, url, description, enabled, always_on FROM mcp_servers")
		if err == nil {
			for rows.Next() {
				var s ServerEntry
				var argsRaw, envRaw []byte
				var urlOpt, descOpt sql.NullString
				err := rows.Scan(&s.UUID, &s.Name, &s.Type, &s.Command, &argsRaw, &envRaw, &urlOpt, &descOpt, &s.Enabled, &s.AlwaysOn)
				if err == nil {
					_ = json.Unmarshal(argsRaw, &s.Args)
					_ = json.Unmarshal(envRaw, &s.Env)
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
					_ = json.Unmarshal(schemaRaw, &t.InputSchema)
					inventory.Tools = append(inventory.Tools, t)
				}
			}
			tRows.Close()
		}

		if len(inventory.Tools) > dbToolCountBefore {
			inventory.Source = "database"
		}
	}

	sortInventory(inventory)
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
			AlwaysOn    bool        `json:"alwaysOn"`
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

func loadInventoryCache(cachePath string) (*Inventory, error) {
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}
	var persisted persistedInventoryCache
	if err := json.Unmarshal(data, &persisted); err != nil {
		return nil, err
	}
	inventory := persisted.Inventory
	inventory.Source = "cache"
	if inventory.CachedAt == "" {
		inventory.CachedAt = persisted.CachedAt
	}
	sortInventory(&inventory)
	return &inventory, nil
}

func saveInventoryCache(cachePath string, inventory *Inventory, cachedAt string) error {
	if strings.TrimSpace(cachePath) == "" || inventory == nil {
		return nil
	}
	persisted := persistedInventoryCache{
		Version:  1,
		CachedAt: cachedAt,
		Inventory: Inventory{
			Servers:  cloneServerEntries(inventory.Servers),
			Tools:    cloneToolEntries(inventory.Tools),
			Source:   inventory.Source,
			CachedAt: cachedAt,
		},
	}
	data, err := json.MarshalIndent(persisted, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(cachePath, data, 0o644)
}

func removeInventoryCache(cachePath string) error {
	if strings.TrimSpace(cachePath) == "" {
		return nil
	}
	err := os.Remove(cachePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func hasInventoryContents(inventory *Inventory) bool {
	if inventory == nil {
		return false
	}
	return len(inventory.Servers) > 0 || len(inventory.Tools) > 0
}

func sortInventory(inventory *Inventory) {
	if inventory == nil {
		return
	}
	sort.Slice(inventory.Servers, func(i, j int) bool {
		left := inventory.Servers[i]
		right := inventory.Servers[j]
		if left.Name == right.Name {
			return left.UUID < right.UUID
		}
		return left.Name < right.Name
	})
	sort.Slice(inventory.Tools, func(i, j int) bool {
		left := inventory.Tools[i]
		right := inventory.Tools[j]
		leftName := left.AdvertisedName
		if strings.TrimSpace(leftName) == "" {
			leftName = left.Name
		}
		rightName := right.AdvertisedName
		if strings.TrimSpace(rightName) == "" {
			rightName = right.Name
		}
		if leftName == rightName {
			return left.Server < right.Server
		}
		return leftName < rightName
	})
}

func cloneServerEntries(source []ServerEntry) []ServerEntry {
	if len(source) == 0 {
		return []ServerEntry{}
	}
	cloned := make([]ServerEntry, 0, len(source))
	for _, entry := range source {
		copyEntry := entry
		copyEntry.Args = append([]string(nil), entry.Args...)
		copyEntry.Env = cloneStringMap(entry.Env)
		copyEntry.Tags = append([]string(nil), entry.Tags...)
		cloned = append(cloned, copyEntry)
	}
	return cloned
}

func cloneToolEntries(source []ToolEntry) []ToolEntry {
	if len(source) == 0 {
		return []ToolEntry{}
	}
	cloned := make([]ToolEntry, 0, len(source))
	for _, entry := range source {
		copyEntry := entry
		copyEntry.ServerTags = append([]string(nil), entry.ServerTags...)
		copyEntry.ToolTags = append([]string(nil), entry.ToolTags...)
		copyEntry.Keywords = append([]string(nil), entry.Keywords...)
		cloned = append(cloned, copyEntry)
	}
	return cloned
}

func cloneStringMap(source map[string]string) map[string]string {
	if len(source) == 0 {
		return map[string]string{}
	}
	cloned := make(map[string]string, len(source))
	for key, value := range source {
		cloned[key] = value
	}
	return cloned
}
