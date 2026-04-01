package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/borghq/hypercode-go/internal/providers"
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
	var payload any
	if taskType != "" {
		payload = map[string]any{"taskType": taskType}
	}

	var result map[string]any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "billing.getFallbackChain", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "billing.getFallbackChain",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    buildLocalFallbackChainResponse(taskType),
		"bridge": map[string]any{
			"fallback":  "go-local-provider-routing",
			"procedure": "billing.getFallbackChain",
			"reason":    "upstream unavailable; using local provider fallback chain preview",
		},
	})
}

func (s *Server) handleBillingTaskRoutingRules(w http.ResponseWriter, r *http.Request) {
	var result map[string]any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "billing.getTaskRoutingRules", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "billing.getTaskRoutingRules",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    buildLocalTaskRoutingRulesResponse(),
		"bridge": map[string]any{
			"fallback":  "go-local-provider-routing",
			"procedure": "billing.getTaskRoutingRules",
			"reason":    "upstream unavailable; using local provider routing rules preview",
		},
	})
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

func buildLocalFallbackChainResponse(taskType string) map[string]any {
	summary := providers.BuildRoutingSummary(providers.Snapshot())
	catalogByProvider := make(map[string]providers.CatalogEntry)
	for _, entry := range providers.Catalog(providers.Snapshot()) {
		catalogByProvider[entry.Provider] = entry
	}

	selectedTaskType := any(nil)
	if taskType != "" {
		selectedTaskType = taskType
	}

	chain := make([]map[string]any, 0)
	for _, task := range summary.Tasks {
		if taskType != "" && task.TaskType != taskType {
			continue
		}
		for index, candidate := range task.Candidates {
			model := ""
			if entry, ok := catalogByProvider[candidate.Provider]; ok {
				model = entry.DefaultModel
			}
			chain = append(chain, map[string]any{
				"priority": index + 1,
				"provider": candidate.Provider,
				"model":    model,
				"reason":   candidate.Reason,
			})
		}
		break
	}

	return map[string]any{
		"selectedTaskType": selectedTaskType,
		"chain":            chain,
	}
}

func buildLocalTaskRoutingRulesResponse() map[string]any {
	summary := providers.BuildRoutingSummary(providers.Snapshot())
	catalogByProvider := make(map[string]providers.CatalogEntry)
	for _, entry := range providers.Catalog(providers.Snapshot()) {
		catalogByProvider[entry.Provider] = entry
	}

	rules := make([]map[string]any, 0, len(summary.Tasks))
	for _, task := range summary.Tasks {
		fallbackPreview := make([]map[string]any, 0, 3)
		for _, candidate := range task.Candidates {
			if len(fallbackPreview) >= 3 {
				break
			}
			model := ""
			if entry, ok := catalogByProvider[candidate.Provider]; ok {
				model = entry.DefaultModel
			}
			fallbackPreview = append(fallbackPreview, map[string]any{
				"provider": candidate.Provider,
				"model":    model,
				"reason":   candidate.Reason,
			})
		}

		rules = append(rules, map[string]any{
			"taskType":        task.TaskType,
			"strategy":        task.Strategy,
			"fallbackPreview": fallbackPreview,
		})
	}

	return map[string]any{
		"defaultStrategy": summary.DefaultStrategy,
		"rules":           rules,
	}
}
