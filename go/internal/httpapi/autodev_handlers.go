package httpapi

import (
	"net/http"
	"strings"
)

func (s *Server) handleAutoDevStartLoop(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "autoDev.startLoop")
}

func (s *Server) handleAutoDevCancelLoop(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "autoDev.cancelLoop")
}

func (s *Server) handleAutoDevGetLoops(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "autoDev.getLoops", nil)
}

func (s *Server) handleAutoDevGetLoop(w http.ResponseWriter, r *http.Request) {
	loopID := strings.TrimSpace(r.URL.Query().Get("loopId"))
	if loopID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "missing loopId query parameter"})
		return
	}
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "autoDev.getLoop", map[string]any{"loopId": loopID})
}

func (s *Server) handleAutoDevClearCompleted(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodPost, "autoDev.clearCompleted", nil)
}
