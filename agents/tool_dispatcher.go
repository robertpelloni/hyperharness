package agents

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "strings"
    "sync"
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

    // edit – placeholder that simply returns the newText (no in‑place edit).
    d.Register("edit", func(ctx context.Context, args map[string]interface{}) (string, error) {
        nt, ok := args["newText"].(string)
        if !ok {
            return "", fmt.Errorf("edit: missing newText")
        }
        return nt, nil
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

    // tree – stub that just returns a placeholder.
    d.Register("tree", func(ctx context.Context, args map[string]interface{}) (string, error) {
        return "[tree view not implemented]", nil
    })

    // websearch – placeholder that returns a static message.
    d.Register("websearch", func(ctx context.Context, args map[string]interface{}) (string, error) {
        q, ok := args["query"].(string)
        if !ok {
            return "", fmt.Errorf("websearch: missing query")
        }
        return fmt.Sprintf("search results for %s (stub)", q), nil
    })

    // webfetch – stub that returns a static message.
    d.Register("webfetch", func(ctx context.Context, args map[string]interface{}) (string, error) {
        url, ok := args["url"].(string)
        if !ok {
            return "", fmt.Errorf("webfetch: missing url")
        }
        return fmt.Sprintf("fetched content from %s (stub)", url), nil
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
