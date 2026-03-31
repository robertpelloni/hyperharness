package httpapi

import "net/http"

func (s *Server) handleDarwinEvolve(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "darwin.evolve")
}

func (s *Server) handleDarwinExperiment(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "darwin.experiment")
}

func (s *Server) handleDarwinStatus(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "darwin.getStatus", nil)
}
