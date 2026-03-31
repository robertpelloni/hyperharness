package httpapi

import (
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) handleBrowserStatus(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "browser.status", nil)
}

func (s *Server) handleBrowserClosePage(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "browser.closePage")
}

func (s *Server) handleBrowserCloseAll(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodPost, "browser.closeAll", nil)
}

func (s *Server) handleBrowserSearchHistory(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	if query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "missing query parameter"})
		return
	}

	payload := map[string]any{"query": query}
	if maxResults := strings.TrimSpace(r.URL.Query().Get("maxResults")); maxResults != "" {
		parsed, err := strconv.Atoi(maxResults)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid maxResults query parameter"})
			return
		}
		payload["maxResults"] = parsed
	}
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "browser.searchHistory", payload)
}

func (s *Server) handleBrowserScrapePage(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodGet, "browser.scrapePage", nil)
}

func (s *Server) handleBrowserScreenshot(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeCall(w, r, http.MethodPost, "browser.screenshot", nil)
}

func (s *Server) handleBrowserDebug(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "browser.debug")
}

func (s *Server) handleBrowserProxyFetch(w http.ResponseWriter, r *http.Request) {
	s.handleTRPCBridgeBodyCall(w, r, "browser.proxyFetch")
}
