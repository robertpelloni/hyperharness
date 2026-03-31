package httpapi

import "net/http"

func (s *Server) handleSquadList(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "squad.list", nil)
}

func (s *Server) handleSquadSpawn(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "squad.spawn")
}

func (s *Server) handleSquadKill(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "squad.kill")
}

func (s *Server) handleSquadChat(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "squad.chat")
}

func (s *Server) handleSquadToggleIndexer(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "squad.toggleIndexer")
}

func (s *Server) handleSquadIndexerStatus(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "squad.getIndexerStatus", nil)
}
