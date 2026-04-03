package httpapi

import (
	"net/http"

	"github.com/hypercodehq/hypercode-go/internal/mcp"
)

func (s *Server) handleMCPStatus(w http.ResponseWriter, r *http.Request) {
	inventory, err := mcp.LoadInventory(s.cfg.WorkspaceRoot, s.cfg.MainConfigDir)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"initialized":    true,
			"serverCount":    len(inventory.Servers),
			"toolCount":      len(inventory.Tools),
			"connectedCount": len(inventory.Servers), // Mocking connected as all for now
			"source":         inventory.Source,
		},
	})
}

func (s *Server) handleMCPTools(w http.ResponseWriter, r *http.Request) {
	inventory, err := mcp.LoadInventory(s.cfg.WorkspaceRoot, s.cfg.MainConfigDir)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    inventory.Tools,
	})
}

func (s *Server) handleMCPSearchTools(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	inventory, err := mcp.LoadInventory(s.cfg.WorkspaceRoot, s.cfg.MainConfigDir)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "error": err.Error()})
		return
	}

	// Simple search implementation
	var filtered []mcp.ToolEntry
	for _, tool := range inventory.Tools {
		if query == "" || 
		   (len(tool.Name) > 0 && (tool.Name == query)) ||
		   (len(tool.Description) > 0 && (tool.Description == query)) {
			filtered = append(filtered, tool)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    filtered,
	})
}

func (s *Server) handleMCPRuntimeServers(w http.ResponseWriter, r *http.Request) {
	inventory, err := mcp.LoadInventory(s.cfg.WorkspaceRoot, s.cfg.MainConfigDir)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    inventory.Servers,
	})
}
