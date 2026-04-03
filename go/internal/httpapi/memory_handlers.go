package httpapi

import (
	"net/http"
	"strconv"

	"github.com/hypercodehq/hypercode-go/internal/memorystore"
)

func (s *Server) handleMemorySearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "query parameter is required"})
		return
	}

	results, err := memorystore.Search(s.cfg.WorkspaceRoot, query, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "error": "fallback memory search failed: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    results,
		"bridge": map[string]any{
			"fallback":  "go-local-memory",
			"procedure": "memory.search",
			"reason":    "upstream unavailable; using local full-text memory fallback",
		},
	})
}

func (s *Server) handleMemoryContexts(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "memory.contexts", nil)
}

func (s *Server) handleMemorySectionedStatus(w http.ResponseWriter, r *http.Request) {
	status, err := memorystore.ReadStatus(s.cfg.WorkspaceRoot)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    status,
	})
}
