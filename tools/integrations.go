package tools

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// WebFetchParams defines parameters for web fetching.
type WebFetchParams struct {
	URL     string `json:"url"`
	Format  string `json:"format,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
}

// WebFetchResult is the result of web fetching.
type WebFetchResult struct {
	Content   string `json:"content,omitempty"`
	IsImage   bool   `json:"isImage,omitempty"`
	ImageMime string `json:"imageMime,omitempty"`
	ImageData string `json:"imageData,omitempty"`
}

// FetchWebContent fetches content from a URL.
func FetchWebContent(ctx context.Context, params WebFetchParams) (*WebFetchResult, error) {
	if params.URL == "" {
		return nil, fmt.Errorf("URL is required")
	}
	timeout := params.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", params.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "HyperHarness/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/") {
		data, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024))
		if err != nil {
			return nil, err
		}
		return &WebFetchResult{
			IsImage:   true,
			ImageMime: contentType,
			ImageData: base64.StdEncoding.EncodeToString(data),
		}, nil
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1*1024*1024))
	if err != nil {
		return nil, err
	}
	return &WebFetchResult{Content: string(body)}, nil
}

// DownloadFile downloads a file from a URL to a local path.
func DownloadFile(ctx context.Context, url, filePath string, timeout int) (int64, error) {
	if url == "" || filePath == "" {
		return 0, fmt.Errorf("URL and filePath are required")
	}
	if timeout <= 0 {
		timeout = 60
	}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return 0, err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	written, err := io.Copy(f, resp.Body)
	return written, err
}

// fetchURL makes an HTTP request with custom method, URL, and headers.
func fetchURL(url, method string, headers map[string]string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("URL is required")
	}
	if method == "" {
		method = "GET"
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1*1024*1024))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (r *Registry) registerIntegrationsTools() {
	// Integration tools are registered via the tool functions above
}

func (r *Registry) registerComputerUseTools() {
	// Computer use tools stub
}

func (r *Registry) registerAiderTools() {
	// Aider tools stub
}

// registerClaudeCodeTools is a free function - add wrapper method
func (r *Registry) registerClaudeCodeTools() {
	registerClaudeCodeTools(r)
}
