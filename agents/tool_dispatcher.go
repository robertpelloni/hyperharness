package agents

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "strings"
    "sync"
    "time"
)

// ToolHandler is the signature for functions that implement a tool.
// args is a map of argument name → value as supplied by the LLM.
type ToolHandler func(ctx context.Context, args map[string]interface{}) (string, error)

// ToolDispatcher routes tool calls from the LLM to concrete Go implementations.
type ToolDispatcher struct {
    mu       sync.RWMutex
    handlers map[string]ToolHandler
}

// NewToolDispatcher creates a new dispatcher with an empty registry.
func NewToolDispatcher() *ToolDispatcher {
    return &ToolDispatcher{handlers: make(map[string]ToolHandler)}
}

// Register adds a handler for a specific tool name. Panics if name is empty
// or already registered – this mirrors the strict registration model used
// throughout HyperHarness.
func (d *ToolDispatcher) Register(name string, h ToolHandler) {
    if name == "" {
        panic("tool dispatcher: Register called with empty name")
    }
    if h == nil {
        panic("tool dispatcher: Register called with nil handler for " + name)
    }
    d.mu.Lock()
    defer d.mu.Unlock()
    if _, exists := d.handlers[name]; exists {
        panic("tool dispatcher: duplicate registration for " + name)
    }
    d.handlers[name] = h
}

// Dispatch looks up a handler by name and executes it with the supplied
// argument map. If the tool is not registered, an error is returned.
func (d *ToolDispatcher) Dispatch(ctx context.Context, name string, args map[string]interface{}) (string, error) {
    d.mu.RLock()
    h, ok := d.handlers[name]
    d.mu.RUnlock()
    if !ok {
        return "", fmt.Errorf("tool %s not registered", name)
    }
    return h(ctx, args)
}

// RegisterDefaultStubs registers a minimal set of stub tool handlers that
// provide basic functionality required for the TUI to operate without the
// full implementation of every Pi‑tool. These stubs can be replaced later
// with real implementations.
func (d *ToolDispatcher) RegisterDefaultStubs() {
    // read – returns the contents of a file path supplied via "path".
    d.Register("read", func(ctx context.Context, args map[string]interface{}) (string, error) {
        p, ok := args["path"].(string)
        if !ok || p == "" {
            return "", fmt.Errorf("read: missing path argument")
        }
        data, err := os.ReadFile(p)
        if err != nil {
            return "", fmt.Errorf("read: %w", err)
        }
        return string(data), nil
    })

    // write – writes content to a file. Expects "path" and "content".
    d.Register("write", func(ctx context.Context, args map[string]interface{}) (string, error) {
        p, ok1 := args["path"].(string)
        c, ok2 := args["content"].(string)
        if !ok1 || !ok2 {
            return "", fmt.Errorf("write: missing path or content")
        }
        if err := os.WriteFile(p, []byte(c), 0644); err != nil {
            return "", fmt.Errorf("write: %w", err)
        }
        return "ok", nil
    })

    // edit – replaces text in a file. Supports:
    //   - path + oldText + newText: replace first occurrence of oldText with newText
    //   - path + newText: overwrite entire file
    d.Register("edit", func(ctx context.Context, args map[string]interface{}) (string, error) {
        path, okPath := args["path"].(string)
        newText, okNew := args["newText"].(string)
        oldText, okOld := args["oldText"].(string)
        if !okPath {
            return "", fmt.Errorf("edit: missing path")
        }
        if !okNew {
            return "", fmt.Errorf("edit: missing newText")
        }

        // Read current content
        data, err := os.ReadFile(path)
        if err != nil {
            return "", fmt.Errorf("edit: read: %w", err)
        }

        // If oldText is provided, do replacement
        if okOld && oldText != "" {
            content := string(data)
            if !strings.Contains(content, oldText) {
                return "", fmt.Errorf("edit: oldText not found in file")
            }
            newContent := strings.Replace(content, oldText, newText, 1)
            if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
                return "", fmt.Errorf("edit: write: %w", err)
            }
            return fmt.Sprintf("✓ Replaced 1 occurrence in %s", path), nil
        }

        // Otherwise, overwrite entire file
        if err := os.WriteFile(path, []byte(newText), 0644); err != nil {
            return "", fmt.Errorf("edit: write: %w", err)
        }
        return fmt.Sprintf("✓ Wrote %d bytes to %s", len(newText), path), nil
    })

    // bash – executes a shell command and returns stdout.
    // NOTE: This is intentional for a coding agent, but we add basic safety checks.
    d.Register("bash", func(ctx context.Context, args map[string]interface{}) (string, error) {
        cmdStr, ok := args["command"].(string)
        if !ok {
            return "", fmt.Errorf("bash: missing command")
        }
        // Basic safety: reject obviously destructive commands
        dangerous := []string{"rm -rf /", "mkfs", ":(){:|:&};:", "> /dev/sda"}
        for _, d := range dangerous {
            if strings.Contains(cmdStr, d) {
                return "", fmt.Errorf("bash: blocked dangerous command pattern")
            }
        }
        // Allow context cancellation via ctx.
        cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
        out, err := cmd.CombinedOutput()
        if err != nil {
            return string(out), fmt.Errorf("bash: %w", err)
        }
        return string(out), nil
    })

    // grep – simple wrapper around Go's regexp FindAllString.
    d.Register("grep", func(ctx context.Context, args map[string]interface{}) (string, error) {
        pattern, ok1 := args["pattern"].(string)
        text, ok2 := args["text"].(string)
        if !ok1 || !ok2 {
            return "", fmt.Errorf("grep: missing pattern or text")
        }
        re, err := regexp.Compile(pattern)
        if err != nil {
            return "", fmt.Errorf("grep: %w", err)
        }
        matches := re.FindAllString(text, -1)
        return strings.Join(matches, "\n"), nil
    })

    // find – returns file paths matching a glob pattern.
    d.Register("find", func(ctx context.Context, args map[string]interface{}) (string, error) {
        pattern, ok := args["pattern"].(string)
        if !ok {
            return "", fmt.Errorf("find: missing pattern")
        }
        matches, err := filepath.Glob(pattern)
        if err != nil {
            return "", fmt.Errorf("find: %w", err)
        }
        return strings.Join(matches, "\n"), nil
    })

    // ls – list directory entries.
    d.Register("ls", func(ctx context.Context, args map[string]interface{}) (string, error) {
        dir, ok := args["path"].(string)
        if !ok || dir == "" {
            dir = "."
        }
        entries, err := os.ReadDir(dir)
        if err != nil {
            return "", fmt.Errorf("ls: %w", err)
        }
        var out []string
        for _, e := range entries {
            name := e.Name()
            if e.IsDir() {
                name += "/"
            }
            out = append(out, name)
        }
        return strings.Join(out, "\n"), nil
    })

    // tree – returns a simple directory tree view.
    d.Register("tree", func(ctx context.Context, args map[string]interface{}) (string, error) {
        root, ok := args["path"].(string)
        if !ok || root == "" {
            root = "."
        }
        // Walk the directory tree
        var lines []string
        err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            // Compute relative path for indentation
            rel, _ := filepath.Rel(root, p)
            if rel == "." {
                // Skip root itself
                return nil
            }
            depth := strings.Count(rel, string(os.PathSeparator))
            indent := strings.Repeat("│   ", depth)
            name := info.Name()
            if info.IsDir() {
                name += "/"
            }
            lines = append(lines, fmt.Sprintf("%s%s", indent, name))
            return nil
        })
        if err != nil {
            return "", fmt.Errorf("tree: %w", err)
        }
        if len(lines) == 0 {
            return "(empty)", nil
        }
        return strings.Join(lines, "\n"), nil
    })

    // websearch – searches the web using DuckDuckGo's Instant Answer API
    d.Register("websearch", func(ctx context.Context, args map[string]interface{}) (string, error) {
        query, ok := args["query"].(string)
        if !ok {
            return "", fmt.Errorf("websearch: missing query")
        }
        apiURL := fmt.Sprintf("https://api.duckduckgo.com/?q=%s&format=json&no_html=1", url.QueryEscape(query))
        req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
        if err != nil {
            return "", fmt.Errorf("websearch: %w", err)
        }
        client := &http.Client{Timeout: 30 * time.Second}
        resp, err := client.Do(req)
        if err != nil {
            return "", fmt.Errorf("websearch: %w", err)
        }
        defer resp.Body.Close()
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            return "", fmt.Errorf("websearch: read: %w", err)
        }
        // Parse JSON response
        var result struct {
            AbstractText string `json:"AbstractText"`
            Heading      string `json:"Heading"`
            Results      []struct {
                Title string `json:"Title"`
                Text  string `json:"Text"`
                URL   string `json:"URL"`
            } `json:"Results"`
        }
        if err := json.Unmarshal(body, &result); err != nil {
            return string(body), nil
        }
        var lines []string
        if result.Heading != "" {
            lines = append(lines, fmt.Sprintf("**%s**", result.Heading))
        }
        if result.AbstractText != "" {
            lines = append(lines, result.AbstractText)
            lines = append(lines, "")
        }
        for _, r := range result.Results {
            lines = append(lines, fmt.Sprintf("• %s\n  %s\n  %s", r.Title, r.Text, r.URL))
        }
        if len(lines) == 0 {
            return "No results found", nil
        }
        return strings.Join(lines, "\n\n"), nil
    })

    // webfetch – fetches content from a URL
    d.Register("webfetch", func(ctx context.Context, args map[string]interface{}) (string, error) {
        urlStr, ok := args["url"].(string)
        if !ok {
            return "", fmt.Errorf("webfetch: missing url")
        }
        req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
        if err != nil {
            return "", fmt.Errorf("webfetch: %w", err)
        }
        client := &http.Client{Timeout: 30 * time.Second}
        resp, err := client.Do(req)
        if err != nil {
            return "", fmt.Errorf("webfetch: %w", err)
        }
        defer resp.Body.Close()
        if resp.StatusCode != http.StatusOK {
            return "", fmt.Errorf("webfetch: %s", resp.Status)
        }
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            return "", fmt.Errorf("webfetch: read: %w", err)
        }
        return string(body), nil
    })
}

// Helper to unmarshal JSON args safely – used by the Director when
// dispatching tool calls.
func parseToolArgs(raw string) (map[string]interface{}, error) {
    var args map[string]interface{}
    if err := json.Unmarshal([]byte(raw), &args); err != nil {
        return nil, fmt.Errorf("failed to parse tool args: %w", err)
    }
    return args, nil
}
