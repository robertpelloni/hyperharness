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
	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "session.getState", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "session.getState",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.sessionState.snapshot(),
		"bridge": map[string]any{
			"fallback":  "go-local-session-state",
			"procedure": "session.getState",
			"reason":    "upstream unavailable; using native Go session-state snapshot",
		},
	})
}

func (s *Server) handleSupervisorSessionUpdateState(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any
	if err := decodeJSONBody(r, &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid JSON body"})
		return
	}

	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "session.updateState", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "session.updateState",
			},
		})
		return
	}

	nextState := s.sessionState.update(payload)
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"success":            true,
			"toolAdvertisements": []string{},
			"memoryBootstrap":    nil,
			"state":              nextState,
		},
		"bridge": map[string]any{
			"fallback":  "go-local-session-state",
			"procedure": "session.updateState",
			"reason":    "upstream unavailable; using native Go session-state persistence without TS memory/bootstrap enrichment",
		},
	})
}

func (s *Server) handleSupervisorSessionClear(w http.ResponseWriter, r *http.Request) {
	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "session.clear", map[string]any{}, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "session.clear",
			},
		})
		return
	}

	s.sessionState.clear()
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    map[string]any{"success": true},
		"bridge": map[string]any{
			"fallback":  "go-local-session-state",
			"procedure": "session.clear",
			"reason":    "upstream unavailable; using native Go session-state reset",
		},
	})
}

func (s *Server) handleSupervisorSessionHeartbeat(w http.ResponseWriter, r *http.Request) {
	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "session.heartbeat", map[string]any{}, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "session.heartbeat",
			},
		})
		return
	}

	state := s.sessionState.touch()
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"alive":     true,
			"timestamp": state["lastHeartbeat"],
		},
		"bridge": map[string]any{
			"fallback":  "go-local-session-state",
			"procedure": "session.heartbeat",
			"reason":    "upstream unavailable; using native Go session-state heartbeat",
		},
	})
}

func (s *Server) handleSupervisorSessionRestore(w http.ResponseWriter, r *http.Request) {
	s.handleSessionBridgeCall(w, r, http.MethodPost, "session.restore", nil)
}
