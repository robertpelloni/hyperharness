package httpapi

import (
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) handleSupervisorSessionLogs(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(r.URL.Query().Get("id"))
	if sessionID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "missing id query parameter",
		})
		return
	}

	payload := map[string]any{"id": sessionID}
	if limit := strings.TrimSpace(r.URL.Query().Get("limit")); limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   "invalid limit query parameter",
			})
			return
		}
		payload["limit"] = parsedLimit
	}

	s.handleSessionBridgeCall(w, r, http.MethodGet, "session.logs", payload)
}

func (s *Server) handleSupervisorSessionExecuteShell(w http.ResponseWriter, r *http.Request) {
	s.handleSessionBridgeBodyCall(w, r, "session.executeShell")
}

func (s *Server) handleSupervisorSessionAttachInfo(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(r.URL.Query().Get("id"))
	if sessionID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "missing id query parameter",
		})
		return
	}

	s.handleSessionBridgeCall(w, r, http.MethodGet, "session.attachInfo", map[string]any{"id": sessionID})
}

func (s *Server) handleSupervisorSessionHealth(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(r.URL.Query().Get("id"))
	if sessionID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "missing id query parameter",
		})
		return
	}

	s.handleSessionBridgeCall(w, r, http.MethodGet, "session.health", map[string]any{"id": sessionID})
}

func (s *Server) handleSupervisorSessionState(w http.ResponseWriter, r *http.Request) {
	s.handleSessionBridgeCall(w, r, http.MethodGet, "session.getState", nil)
}

func (s *Server) handleSupervisorSessionUpdateState(w http.ResponseWriter, r *http.Request) {
	s.handleSessionBridgeBodyCall(w, r, "session.updateState")
}

func (s *Server) handleSupervisorSessionClear(w http.ResponseWriter, r *http.Request) {
	s.handleSessionBridgeCall(w, r, http.MethodPost, "session.clear", nil)
}

func (s *Server) handleSupervisorSessionHeartbeat(w http.ResponseWriter, r *http.Request) {
	s.handleSessionBridgeCall(w, r, http.MethodPost, "session.heartbeat", nil)
}

func (s *Server) handleSupervisorSessionRestore(w http.ResponseWriter, r *http.Request) {
	s.handleSessionBridgeCall(w, r, http.MethodPost, "session.restore", nil)
}
