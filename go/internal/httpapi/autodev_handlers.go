package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (s *Server) handleAutoDevStartLoop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"success": false, "error": "method not allowed"})
		return
	}

	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid JSON body"})
		return
	}

	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "autoDev.startLoop", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "autoDev.startLoop",
			},
		})
		return
	}

	loopType := strings.TrimSpace(stringValue(payload["type"]))
	if loopType == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "missing type"})
		return
	}
	maxAttempts := intNumber(payload["maxAttempts"])
	if maxAttempts <= 0 {
		maxAttempts = 3
	}
	loop := s.autoDevState.startLoop(localAutoDevLoopConfig{
		Type:        loopType,
		MaxAttempts: maxAttempts,
		Target:      strings.TrimSpace(stringValue(payload["target"])),
		Command:     strings.TrimSpace(stringValue(payload["command"])),
	})
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"success": true,
			"loopId":  loop.ID,
			"loop":    loop,
		},
		"bridge": map[string]any{
			"fallback":  "go-local-autodev",
			"procedure": "autoDev.startLoop",
			"reason":    "upstream unavailable; using native Go AutoDev fallback loop manager",
		},
	})
}

func (s *Server) handleAutoDevCancelLoop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"success": false, "error": "method not allowed"})
		return
	}

	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid JSON body"})
		return
	}

	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "autoDev.cancelLoop", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "autoDev.cancelLoop",
			},
		})
		return
	}

	loopID := strings.TrimSpace(stringValue(payload["loopId"]))
	if loopID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "missing loopId"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.autoDevState.cancelLoop(loopID),
		"bridge": map[string]any{
			"fallback":  "go-local-autodev",
			"procedure": "autoDev.cancelLoop",
			"reason":    "upstream unavailable; using native Go AutoDev fallback loop manager",
		},
	})
}

func (s *Server) handleAutoDevGetLoops(w http.ResponseWriter, r *http.Request) {
	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "autoDev.getLoops", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "autoDev.getLoops",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.autoDevState.getLoops(),
		"bridge": map[string]any{
			"fallback":  "go-local-autodev",
			"procedure": "autoDev.getLoops",
			"reason":    "upstream unavailable; using native Go AutoDev fallback loop state",
		},
	})
}

func (s *Server) handleAutoDevGetLoop(w http.ResponseWriter, r *http.Request) {
	loopID := strings.TrimSpace(r.URL.Query().Get("loopId"))
	if loopID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "missing loopId query parameter"})
		return
	}

	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "autoDev.getLoop", map[string]any{"loopId": loopID}, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "autoDev.getLoop",
			},
		})
		return
	}

	loop, ok := s.autoDevState.getLoop(loopID)
	if !ok {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    nil,
			"bridge": map[string]any{
				"fallback":  "go-local-autodev",
				"procedure": "autoDev.getLoop",
				"reason":    "upstream unavailable; loop not present in native Go AutoDev fallback state",
			},
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    loop,
		"bridge": map[string]any{
			"fallback":  "go-local-autodev",
			"procedure": "autoDev.getLoop",
			"reason":    "upstream unavailable; using native Go AutoDev fallback loop state",
		},
	})
}

func (s *Server) handleAutoDevClearCompleted(w http.ResponseWriter, r *http.Request) {
	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "autoDev.clearCompleted", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "autoDev.clearCompleted",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.autoDevState.clearCompleted(),
		"bridge": map[string]any{
			"fallback":  "go-local-autodev",
			"procedure": "autoDev.clearCompleted",
			"reason":    "upstream unavailable; using native Go AutoDev fallback loop state",
		},
	})
}
