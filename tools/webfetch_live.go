// Package tools provides the WebFetch live HTTP implementation.
// Supports GET and POST requests with configurable timeouts,
// user-agent spoofing, and content extraction.
package tools

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultHTTPTimeout = 30 * time.Second
	maxHTTPBodySize    = 5 * 1024 * 1024 // 5MB max response
	harnessUserAgent   = "HyperHarness/0.2.0 (Go HTTP Client)"
)

// httpClient is a shared HTTP client with sensible defaults.
var httpClient = &http.Client{
	Timeout: defaultHTTPTimeout,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
}

// fetchURL performs an HTTP request and returns the response body as text.
// This is the actual implementation backing the WebFetch tool surface.
func fetchURL(url, method string, headers map[string]string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("url is required")
	}

	// Default to GET
	if method == "" {
		method = "GET"
	}
	method = strings.ToUpper(method)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", harnessUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,text/plain;q=0.8,*/*;q=0.5")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Read body with size limit
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxHTTPBodySize))
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(body), nil
}
