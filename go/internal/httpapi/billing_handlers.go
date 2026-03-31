package httpapi

import (
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) handleBillingStatus(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "billing.getStatus", nil)
}

func (s *Server) handleBillingProviderQuotas(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "billing.getProviderQuotas", nil)
}

func (s *Server) handleBillingCostHistory(w http.ResponseWriter, r *http.Request) {
	days := strings.TrimSpace(r.URL.Query().Get("days"))
	if days == "" {
		s.handleTRPCBridgeCall(w, r, http.MethodGet, "billing.getCostHistory", nil)
		return
	}
	parsed, err := strconv.Atoi(days)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid days query parameter"})
		return
	}
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "billing.getCostHistory", map[string]any{"days": parsed})
}

func (s *Server) handleBillingModelPricing(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "billing.getModelPricing", nil)
}

func (s *Server) handleBillingFallbackChain(w http.ResponseWriter, r *http.Request) {
	taskType := strings.TrimSpace(r.URL.Query().Get("taskType"))
	if taskType == "" {
		s.handleTRPCBridgeCall(w, r, http.MethodGet, "billing.getFallbackChain", nil)
		return
	}
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "billing.getFallbackChain", map[string]any{"taskType": taskType})
}

func (s *Server) handleBillingTaskRoutingRules(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "billing.getTaskRoutingRules", nil)
}

func (s *Server) handleBillingSetRoutingStrategy(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "billing.setRoutingStrategy")
}

func (s *Server) handleBillingSetTaskRoutingRule(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "billing.setTaskRoutingRule")
}

func (s *Server) handleBillingDepletedModels(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "billing.getDepletedModels", nil)
}

func (s *Server) handleBillingFallbackHistory(w http.ResponseWriter, r *http.Request) {
	limit := strings.TrimSpace(r.URL.Query().Get("limit"))
	if limit == "" {
		s.handleTRPCBridgeCall(w, r, http.MethodGet, "billing.getFallbackHistory", nil)
		return
	}
	parsed, err := strconv.Atoi(limit)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid limit query parameter"})
		return
	}
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "billing.getFallbackHistory", map[string]any{"limit": parsed})
}

func (s *Server) handleBillingClearFallbackHistory(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodPost, "billing.clearFallbackHistory", nil)
}
