package httpapi

import (
	"net/http"

	"github.com/hypercodehq/hypercode-go/internal/memorystore"
)

func (s *Server) handleMemorySearch(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "memory.search", nil)
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
