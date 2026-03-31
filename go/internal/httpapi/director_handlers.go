package httpapi

import "net/http"

func (s *Server) handleDirectorMemorize(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "director.memorize")
}

func (s *Server) handleDirectorChat(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "director.chat")
}

func (s *Server) handleDirectorStatus(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "director.status", nil)
}

func (s *Server) handleDirectorUpdateConfig(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "director.updateConfig")
}

func (s *Server) handleDirectorConfigGet(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "directorConfig.get", nil)
}

func (s *Server) handleDirectorConfigTest(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "directorConfig.test", nil)
}

func (s *Server) handleDirectorConfigUpdate(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "directorConfig.update")
}

func (s *Server) handleDirectorStopAutoDrive(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodPost, "director.stopAutoDrive", nil)
}

func (s *Server) handleDirectorStartAutoDrive(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodPost, "director.startAutoDrive", nil)
}
