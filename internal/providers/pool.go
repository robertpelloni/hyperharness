package providers

import (
	"net/http"
	"time"
)

var globalHTTPClient *http.Client

// GetPooledHTTPClient returns a globally shared HTTP client configured with aggressive
// connection pooling to reduce TCP/TLS overhead on high-frequency LLM dispatch loops.
func GetPooledHTTPClient() *http.Client {
	if globalHTTPClient == nil {
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.MaxIdleConns = 100
		t.MaxConnsPerHost = 100
		t.MaxIdleConnsPerHost = 100
		t.IdleConnTimeout = 90 * time.Second

		globalHTTPClient = &http.Client{
			Timeout:   120 * time.Second, // Global timeout to prevent hangs
			Transport: t,
		}
	}
	return globalHTTPClient
}
