// Package tools provides Exa-based web search and code search.
// Ported from superai/opencode/packages/opencode/src/tool/websearch.ts and
// superai/opencode/packages/opencode/src/tool/codesearch.ts
//
// Exa provides high-quality web search and code context retrieval via an MCP
// endpoint. The search returns results optimized for LLM consumption with
// configurable depth, result count, and live crawling.
package tools

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	exaBaseURL        = "https://mcp.exa.ai"
	exaEndpoint       = "/mcp"
	exaDefaultTimeout = 25 * time.Second
	exaMaxResponseSize = 5 * 1024 * 1024 // 5MB
)

// ExaSearchParams defines parameters for web search.
type ExaSearchParams struct {
	Query               string `json:"query"`
	NumResults          int    `json:"numResults,omitempty"`
	Livecrawl           string `json:"livecrawl,omitempty"`  // "fallback" | "preferred"
	Type                string `json:"type,omitempty"`       // "auto" | "fast" | "deep"
	ContextMaxCharacters int   `json:"contextMaxCharacters,omitempty"`
}

// ExaCodeSearchParams defines parameters for code search.
type ExaCodeSearchParams struct {
	Query     string `json:"query"`
	TokensNum int    `json:"tokensNum,omitempty"`
}

// exaRequest is a JSON-RPC request to the Exa MCP endpoint.
type exaRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  exaParams   `json:"params"`
}

type exaParams struct {
	Name      string      `json:"name"`
	Arguments interface{} `json:"arguments"`
}

// exaResponse is a JSON-RPC response from Exa.
type exaResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  *struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"result"`
}

// ExaWebSearch performs a web search via Exa's MCP endpoint.
// Returns the search results as text suitable for LLM consumption.
func ExaWebSearch(ctx context.Context, params ExaSearchParams) (string, error) {
	if params.NumResults == 0 {
		params.NumResults = 8
	}
	if params.Type == "" {
		params.Type = "auto"
	}
	if params.Livecrawl == "" {
		params.Livecrawl = "fallback"
	}

	req := exaRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: exaParams{
			Name:      "web_search_exa",
			Arguments: params,
		},
	}

	return callExa(ctx, req)
}

// ExaCodeSearch performs a code search via Exa's MCP endpoint.
// Returns code snippets and documentation relevant to the query.
func ExaCodeSearch(ctx context.Context, params ExaCodeSearchParams) (string, error) {
	if params.TokensNum == 0 {
		params.TokensNum = 5000
	}

	req := exaRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: exaParams{
			Name:      "get_code_context_exa",
			Arguments: params,
		},
	}

	return callExa(ctx, req)
}

// callExa sends a request to Exa's MCP endpoint and parses the SSE response.
func callExa(ctx context.Context, req exaRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, exaDefaultTimeout)
	defer cancel()

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", exaBaseURL+exaEndpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json, text/event-stream")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("exa request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("exa returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse SSE response
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, exaMaxResponseSize))
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return parseSSEResponse(string(respBody))
}

// parseSSEResponse extracts the text content from an SSE response.
func parseSSEResponse(body string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			var resp exaResponse
			if err := json.Unmarshal([]byte(data), &resp); err != nil {
				continue
			}
			if resp.Result != nil && len(resp.Result.Content) > 0 {
				return resp.Result.Content[0].Text, nil
			}
		}
	}

	return "No results found.", nil
}
