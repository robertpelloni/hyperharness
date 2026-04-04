package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) handleSwarmStart(w http.ResponseWriter, r *http.Request) {
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
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "swarm.startSwarm", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": "swarm.startSwarm"}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.swarmState.startMission(payload),
		"bridge": map[string]any{
			"fallback":  "go-local-swarm",
			"procedure": "swarm.startSwarm",
			"reason":    "upstream unavailable; using native Go swarm mission scaffold",
		},
	})
}

func (s *Server) handleSwarmResumeMission(w http.ResponseWriter, r *http.Request) {
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
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "swarm.resumeMission", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": "swarm.resumeMission"}})
		return
	}
	missionID := strings.TrimSpace(stringValue(payload["missionId"]))
	if missionID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "missing missionId"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    map[string]any{"success": s.swarmState.resumeMission(missionID)},
		"bridge": map[string]any{
			"fallback":  "go-local-swarm",
			"procedure": "swarm.resumeMission",
			"reason":    "upstream unavailable; using native Go swarm mission state",
		},
	})
}

func (s *Server) handleSwarmApproveTask(w http.ResponseWriter, r *http.Request) {
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
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "swarm.approveTask", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": "swarm.approveTask"}})
		return
	}
	missionID := strings.TrimSpace(stringValue(payload["missionId"]))
	taskID := strings.TrimSpace(stringValue(payload["taskId"]))
	approved, _ := payload["approved"].(bool)
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    map[string]any{"success": s.swarmState.approveTask(missionID, taskID, approved)},
		"bridge": map[string]any{
			"fallback":  "go-local-swarm",
			"procedure": "swarm.approveTask",
			"reason":    "upstream unavailable; using native Go swarm mission state",
		},
	})
}

func (s *Server) handleSwarmDecomposeTask(w http.ResponseWriter, r *http.Request) {
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
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "swarm.decomposeTask", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": "swarm.decomposeTask"}})
		return
	}
	subMissionID, ok := s.swarmState.decomposeTask(strings.TrimSpace(stringValue(payload["missionId"])), strings.TrimSpace(stringValue(payload["taskId"])))
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    map[string]any{"success": ok, "subMissionId": nullableString(subMissionID)},
		"bridge": map[string]any{
			"fallback":  "go-local-swarm",
			"procedure": "swarm.decomposeTask",
			"reason":    "upstream unavailable; using native Go swarm mission decomposition scaffold",
		},
	})
}

func (s *Server) handleSwarmUpdateTaskPriority(w http.ResponseWriter, r *http.Request) {
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
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "swarm.updateTaskPriority", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": "swarm.updateTaskPriority"}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    map[string]any{"success": s.swarmState.updateTaskPriority(strings.TrimSpace(stringValue(payload["missionId"])), strings.TrimSpace(stringValue(payload["taskId"])), intNumber(payload["priority"]))},
		"bridge": map[string]any{
			"fallback":  "go-local-swarm",
			"procedure": "swarm.updateTaskPriority",
			"reason":    "upstream unavailable; using native Go swarm mission state",
		},
	})
}

func (s *Server) handleSwarmExecuteDebate(w http.ResponseWriter, r *http.Request) {
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
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "swarm.executeDebate", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": "swarm.executeDebate"}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.swarmState.executeDebate(strings.TrimSpace(stringValue(payload["topic"])), strings.TrimSpace(stringValue(payload["proponentModel"])), strings.TrimSpace(stringValue(payload["opponentModel"])), intNumber(payload["rounds"])),
		"bridge": map[string]any{
			"fallback":  "go-local-swarm",
			"procedure": "swarm.executeDebate",
			"reason":    "upstream unavailable; using native Go swarm debate heuristic",
		},
	})
}

func (s *Server) handleSwarmSeekConsensus(w http.ResponseWriter, r *http.Request) {
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
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "swarm.seekConsensus", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": "swarm.seekConsensus"}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.swarmState.seekConsensus(strings.TrimSpace(stringValue(payload["prompt"])), stringArray(payload["models"])),
		"bridge": map[string]any{
			"fallback":  "go-local-swarm",
			"procedure": "swarm.seekConsensus",
			"reason":    "upstream unavailable; using native Go swarm consensus heuristic",
		},
	})
}

func (s *Server) handleSwarmMissionHistory(w http.ResponseWriter, r *http.Request) {
	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "swarm.getMissionHistory", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": "swarm.getMissionHistory"}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.swarmState.missionHistory(),
		"bridge": map[string]any{
			"fallback":  "go-local-swarm",
			"procedure": "swarm.getMissionHistory",
			"reason":    "upstream unavailable; using native Go swarm mission history",
		},
	})
}

func (s *Server) handleSwarmMissionRiskSummary(w http.ResponseWriter, r *http.Request) {
	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "swarm.getMissionRiskSummary", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": "swarm.getMissionRiskSummary"}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.swarmState.missionRiskSummary(),
		"bridge": map[string]any{
			"fallback":  "go-local-swarm",
			"procedure": "swarm.getMissionRiskSummary",
			"reason":    "upstream unavailable; using native Go swarm mission risk summary",
		},
	})
}

func (s *Server) handleSwarmMissionRiskRows(w http.ResponseWriter, r *http.Request) {
	payload := map[string]any{}
	statusFilter := strings.TrimSpace(r.URL.Query().Get("statusFilter"))
	if statusFilter != "" {
		payload["statusFilter"] = statusFilter
	}
	if sortBy := strings.TrimSpace(r.URL.Query().Get("sortBy")); sortBy != "" {
		payload["sortBy"] = sortBy
	}
	minRiskValue := 0.0
	if minRisk := strings.TrimSpace(r.URL.Query().Get("minRisk")); minRisk != "" {
		if parsed, err := strconv.ParseFloat(minRisk, 64); err == nil {
			minRiskValue = parsed
			payload["minRisk"] = parsed
		}
	}
	limitValue := 0
	if limit := strings.TrimSpace(r.URL.Query().Get("limit")); limit != "" {
		if parsed, err := strconv.Atoi(limit); err == nil {
			limitValue = parsed
			payload["limit"] = parsed
		}
	}
	var result any
	proc := "swarm.getMissionRiskRows"
	upstreamBase, err := s.callUpstreamJSON(r.Context(), proc, payloadOrNil(payload), &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": proc}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.swarmState.missionRiskRows(statusFilter, minRiskValue, limitValue),
		"bridge": map[string]any{
			"fallback":  "go-local-swarm",
			"procedure": proc,
			"reason":    "upstream unavailable; using native Go swarm mission risk rows",
		},
	})
}

func (s *Server) handleSwarmMissionRiskFacets(w http.ResponseWriter, r *http.Request) {
	payload := map[string]any{}
	statusFilter := strings.TrimSpace(r.URL.Query().Get("statusFilter"))
	if statusFilter != "" {
		payload["statusFilter"] = statusFilter
	}
	minRiskValue := 0.0
	if minRisk := strings.TrimSpace(r.URL.Query().Get("minRisk")); minRisk != "" {
		if parsed, err := strconv.ParseFloat(minRisk, 64); err == nil {
			minRiskValue = parsed
			payload["minRisk"] = parsed
		}
	}
	var result any
	proc := "swarm.getMissionRiskFacets"
	upstreamBase, err := s.callUpstreamJSON(r.Context(), proc, payloadOrNil(payload), &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": proc}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.swarmState.missionRiskFacets(statusFilter, minRiskValue),
		"bridge": map[string]any{
			"fallback":  "go-local-swarm",
			"procedure": proc,
			"reason":    "upstream unavailable; using native Go swarm mission risk facets",
		},
	})
}

func (s *Server) handleSwarmMeshCapabilities(w http.ResponseWriter, r *http.Request) {
	var result any
	proc := "swarm.getMeshCapabilities"
	upstreamBase, err := s.callUpstreamJSON(r.Context(), proc, nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": proc}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.swarmState.meshCapabilities(),
		"bridge": map[string]any{
			"fallback":  "go-local-swarm",
			"procedure": proc,
			"reason":    "upstream unavailable; using native Go swarm mesh capability snapshot",
		},
	})
}

func (s *Server) handleSwarmSendDirectMessage(w http.ResponseWriter, r *http.Request) {
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
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "swarm.sendDirectMessage", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": "swarm.sendDirectMessage"}})
		return
	}
	target := strings.TrimSpace(stringValue(payload["targetNodeId"]))
	if target == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "missing targetNodeId"})
		return
	}
	messagePayload, _ := payload["payload"].(map[string]any)
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.swarmState.sendDirectMessage(target, messagePayload),
		"bridge": map[string]any{
			"fallback":  "go-local-swarm",
			"procedure": "swarm.sendDirectMessage",
			"reason":    "upstream unavailable; using native Go swarm direct-message acknowledgment",
		},
	})
}

func payloadOrNil(payload map[string]any) map[string]any {
	if len(payload) == 0 {
		return nil
	}
	return payload
}
