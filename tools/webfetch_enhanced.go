// Package tools provides an enhanced WebFetch with format conversion.
// Ported from superai/opencode/packages/opencode/src/tool/webfetch.ts and
// superai/crush/internal/agent/tools/fetch.go
//
// Supports three output formats:
// - "text": plain text extraction from HTML
// - "markdown": HTML-to-Markdown conversion (default)
// - "html": raw HTML (body only)
//
// Also handles image responses with base64 encoding.
package tools

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	webFetchMaxSize    = 5 * 1024 * 1024 // 5MB
	webFetchTimeout    = 30 * time.Second
	webFetchMaxTimeout = 120 * time.Second
)

// WebFetchParams defines parameters for the enhanced web fetch.
type WebFetchParams struct {
	URL     string `json:"url"`
	Format  string `json:"format"`            // "text", "markdown", "html"
	Timeout int    `json:"timeout,omitempty"` // seconds, max 120
}

// WebFetchResult contains the fetch response with metadata.
type WebFetchResult struct {
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
	IsImage     bool   `json:"isImage,omitempty"`
	ImageData   string `json:"imageData,omitempty"`
	ImageMime   string `json:"imageMime,omitempty"`
}

// FetchWebContent fetches a URL and converts the response to the requested format.
// This implements the OpenCode/Crush webfetch tool with full format support.
func FetchWebContent(ctx context.Context, params WebFetchParams) (*WebFetchResult, error) {
	if params.URL == "" {
		return nil, fmt.Errorf("url is required")
	}
	if !strings.HasPrefix(params.URL, "http://") && !strings.HasPrefix(params.URL, "https://") {
		return nil, fmt.Errorf("URL must start with http:// or https://")
	}

	// Validate format
	format := strings.ToLower(params.Format)
	if format == "" {
		format = "markdown"
	}
	if format != "text" && format != "markdown" && format != "html" {
		return nil, fmt.Errorf("format must be one of: text, markdown, html")
	}

	// Set up timeout
	timeout := webFetchTimeout
	if params.Timeout > 0 {
		timeout = time.Duration(params.Timeout) * time.Second
		if timeout > webFetchMaxTimeout {
			timeout = webFetchMaxTimeout
		}
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Build Accept header based on requested format
	acceptHeader := "*/*"
	switch format {
	case "markdown":
		acceptHeader = "text/markdown;q=1.0, text/x-markdown;q=0.9, text/plain;q=0.8, text/html;q=0.7, */*;q=0.1"
	case "text":
		acceptHeader = "text/plain;q=1.0, text/markdown;q=0.9, text/html;q=0.8, */*;q=0.1"
	case "html":
		acceptHeader = "text/html;q=1.0, application/xhtml+xml;q=0.9, text/plain;q=0.8, text/markdown;q=0.7, */*;q=0.1"
	}

	req, err := http.NewRequestWithContext(ctx, "GET", params.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "HyperHarness/0.2.0")
	req.Header.Set("Accept", acceptHeader)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(io.LimitReader(resp.Body, webFetchMaxSize))
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	mime := contentType
	if idx := strings.Index(contentType, ";"); idx != -1 {
		mime = strings.TrimSpace(contentType[:idx])
	}
	mime = strings.ToLower(mime)

	// Check if response is an image
	isImage := strings.HasPrefix(mime, "image/") && mime != "image/svg+xml"
	if isImage {
		return &WebFetchResult{
			Content:     "Image fetched successfully",
			ContentType: contentType,
			IsImage:     true,
			ImageData:   base64.StdEncoding.EncodeToString(body),
			ImageMime:   mime,
		}, nil
	}

	content := string(body)
	if !utf8.ValidString(content) {
		return nil, fmt.Errorf("response content is not valid UTF-8")
	}

	// Convert based on format
	switch format {
	case "markdown":
		if strings.Contains(contentType, "text/html") {
			content = HTMLToMarkdown(content)
		}
	case "text":
		if strings.Contains(contentType, "text/html") {
			content = ExtractTextFromHTML(content)
		}
	case "html":
		// Return as-is (or extract body)
		if strings.Contains(contentType, "text/html") {
			content = extractHTMLBody(content)
		}
	}

	// Truncate if too large
	if len(content) > webFetchMaxSize {
		content = content[:webFetchMaxSize]
		content += fmt.Sprintf("\n\n[Content truncated to %d bytes]", webFetchMaxSize)
	}

	return &WebFetchResult{
		Content:     content,
		ContentType: contentType,
	}, nil
}

// extractHTMLBody extracts the <body> content from an HTML document.
func extractHTMLBody(html string) string {
	lower := strings.ToLower(html)
	startIdx := strings.Index(lower, "<body")
	if startIdx == -1 {
		return html
	}
	// Find end of <body...> tag
	gtIdx := strings.Index(html[startIdx:], ">")
	if gtIdx == -1 {
		return html
	}
	contentStart := startIdx + gtIdx + 1

	endIdx := strings.Index(strings.ToLower(html[contentStart:]), "</body>")
	if endIdx == -1 {
		return html[contentStart:]
	}
	return html[contentStart : contentStart+endIdx]
}

// DownloadFile downloads a file from a URL and saves it to disk.
// Ported from superai/crush/internal/agent/tools/download.go
func DownloadFile(ctx context.Context, url, filePath string, timeoutSec int) (int64, error) {
	if url == "" {
		return 0, fmt.Errorf("url is required")
	}
	if filePath == "" {
		return 0, fmt.Errorf("file_path is required")
	}

	timeout := 5 * time.Minute
	if timeoutSec > 0 {
		if timeoutSec > 600 {
			timeoutSec = 600
		}
		timeout = time.Duration(timeoutSec) * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "HyperHarness/0.2.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Create parent directories
	if err := createParentDirs(filePath); err != nil {
		return 0, fmt.Errorf("failed to create parent directories: %w", err)
	}

	// Write file using the writeToolFile abstraction from tools_native.go
	written, err := writeDownloadFile(filePath, resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to write file: %w", err)
	}

	return written, nil
}
