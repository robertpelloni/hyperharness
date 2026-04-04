package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/hypercodehq/hypercode-go/internal/memorystore"
)

func (s *Server) handleMemorySearch(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	limit := 50
	if limitStr := strings.TrimSpace(r.URL.Query().Get("limit")); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "query parameter is required"})
		return
	}

	payload := map[string]any{"query": query, "limit": limit}
	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "memory.query", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "memory.query",
			},
		})
		return
	}

	results, searchErr := memorystore.Search(s.cfg.WorkspaceRoot, query, limit)
	if searchErr != nil {
		results = []memorystore.SearchResult{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    results,
		"bridge": map[string]any{
			"fallback":  "go-local-memory",
			"procedure": "memory.query",
			"reason":    "upstream unavailable; using local full-text memory fallback",
		},
	})
}

func (s *Server) handleMemoryContexts(w http.ResponseWriter, r *http.Request) {
	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "memory.listContexts", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "memory.listContexts",
			},
		})
		return
	}

	contexts, localErr := s.localMemoryContexts()
	if localErr != nil {
		contexts = []map[string]any{}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    contexts,
		"bridge": map[string]any{
			"fallback":  "go-local-memory",
			"procedure": "memory.listContexts",
			"reason":    "upstream unavailable; using local memory context list",
		},
	})
}

func (s *Server) handleMemorySectionedStatus(w http.ResponseWriter, r *http.Request) {
	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "memory.getSectionedMemoryStatus", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "memory.getSectionedMemoryStatus",
			},
		})
		return
	}

	status, localErr := memorystore.ReadStatus(s.cfg.WorkspaceRoot)
	if localErr != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "error": localErr.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    status,
		"bridge": map[string]any{
			"fallback":  "go-local-memory",
			"procedure": "memory.getSectionedMemoryStatus",
			"reason":    "upstream unavailable; using local sectioned memory status",
		},
	})
}
