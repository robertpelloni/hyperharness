package httpapi

import (
	"net/http"
	"strings"
)

func (s *Server) handleConfigList(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "config.list", nil)
}

func (s *Server) handleConfigGet(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSpace(r.URL.Query().Get("key"))
	if key == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "missing key query parameter"})
		return
	}
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "config.get", map[string]any{"key": key})
}

func (s *Server) handleConfigUpsert(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "config.upsert")
}

func (s *Server) handleConfigDelete(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "config.delete")
}

func (s *Server) handleConfigUpdate(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "config.update")
}

func (s *Server) handleConfigGetMCPTimeout(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "config.getMcpTimeout", nil)
}

func (s *Server) handleConfigSetMCPTimeout(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "config.setMcpTimeout")
}

func (s *Server) handleConfigGetMCPMaxAttempts(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "config.getMcpMaxAttempts", nil)
}

func (s *Server) handleConfigSetMCPMaxAttempts(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "config.setMcpMaxAttempts")
}

func (s *Server) handleConfigGetMCPMaxTotalTimeout(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "config.getMcpMaxTotalTimeout", nil)
}

func (s *Server) handleConfigSetMCPMaxTotalTimeout(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "config.setMcpMaxTotalTimeout")
}

func (s *Server) handleConfigGetMCPResetTimeoutOnProgress(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "config.getMcpResetTimeoutOnProgress", nil)
}

func (s *Server) handleConfigSetMCPResetTimeoutOnProgress(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "config.setMcpResetTimeoutOnProgress")
}

func (s *Server) handleConfigGetSessionLifetime(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "config.getSessionLifetime", nil)
}

func (s *Server) handleConfigSetSessionLifetime(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "config.setSessionLifetime")
}

func (s *Server) handleConfigGetSignupDisabled(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "config.getSignupDisabled", nil)
}

func (s *Server) handleConfigSetSignupDisabled(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "config.setSignupDisabled")
}

func (s *Server) handleConfigGetSSOSignupDisabled(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "config.getSsoSignupDisabled", nil)
}

func (s *Server) handleConfigSetSSOSignupDisabled(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "config.setSsoSignupDisabled")
}

func (s *Server) handleConfigGetBasicAuthDisabled(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "config.getBasicAuthDisabled", nil)
}

func (s *Server) handleConfigSetBasicAuthDisabled(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "config.setBasicAuthDisabled")
}

func (s *Server) handleConfigGetAuthProviders(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "config.getAuthProviders", nil)
}

func (s *Server) handleConfigGetAlwaysVisibleTools(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "config.getAlwaysVisibleTools", nil)
}

func (s *Server) handleConfigSetAlwaysVisibleTools(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "config.setAlwaysVisibleTools")
}
