// Package controlplane provides tool detection for installed CLI harnesses.
// Ported from hypercode/go/internal/controlplane/detector.go
//
// WHAT: Detects 50+ CLI coding tools installed on the system (go, node, claude, gemini, aider, etc.)
// WHY: HyperHarness needs to know what tools are available for routing and integration
// HOW: Concurrent exec.LookPath + version detection with TTL caching
package controlplane

import (
	"context"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Tool represents a detected CLI tool on the system.
type Tool struct {
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Command      string   `json:"command"`
	Available    bool     `json:"available"`
	Version      string   `json:"version,omitempty"`
	Path         string   `json:"path,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
}

// ToolProvider interface for detection.
type ToolProvider interface {
	DetectAll(context.Context) ([]Tool, error)
}

type definition struct {
	Type         string
	Name         string
	Command      string
	VersionArgs  []string
	Capabilities []string
}

// Detector discovers installed CLI tools on the system.
type Detector struct {
	definitions []definition
	timeout     time.Duration
	ttl         time.Duration
	mu          sync.Mutex
	inflight    chan struct{}
	cached      []Tool
	cachedAt    time.Time
}

// NewDetector creates a tool detector with specified timeout and cache TTL.
func NewDetector(timeout, ttl time.Duration) *Detector {
	return &Detector{
		timeout: timeout,
		ttl:     ttl,
		definitions: []definition{
			{Type: "go", Name: "Go", Command: "go", VersionArgs: []string{"version"}, Capabilities: []string{"build", "test", "server"}},
			{Type: "node", Name: "Node.js", Command: "node", VersionArgs: []string{"--version"}, Capabilities: []string{"runtime", "scripts"}},
			{Type: "python", Name: "Python", Command: "python", VersionArgs: []string{"--version"}, Capabilities: []string{"runtime", "scripts"}},
			{Type: "hypercode", Name: "HyperCode CLI", Command: "hypercode", VersionArgs: []string{"version"}, Capabilities: []string{"chat", "edit", "repl", "hypercode-adapter"}},
			{Type: "claude", Name: "Claude CLI", Command: "claude", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "code", "analyze"}},
			{Type: "aider", Name: "Aider CLI", Command: "aider", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "edit", "git-aware", "multi-file"}},
			{Type: "gemini", Name: "Gemini CLI", Command: "gemini", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "code", "multimodal"}},
			{Type: "goose", Name: "Goose CLI", Command: "goose", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "code", "agent"}},
			{Type: "grok", Name: "Grok CLI", Command: "grok", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "code", "realtime"}},
			{Type: "opencode", Name: "OpenCode CLI", Command: "opencode", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "edit", "multi-file", "autonomous"}},
			{Type: "crush", Name: "Crush CLI", Command: "crush", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "code", "data"}},
			{Type: "copilot", Name: "GitHub Copilot CLI", Command: "github-copilot-cli", VersionArgs: []string{"--version"}, Capabilities: []string{"explain", "suggest", "shell"}},
			{Type: "ollama", Name: "Ollama CLI", Command: "ollama", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "local", "models"}},
			{Type: "smithery", Name: "Smithery CLI", Command: "smithery", VersionArgs: []string{"--version"}, Capabilities: []string{"mcp", "registry", "tools"}},
			{Type: "dolt", Name: "Dolt CLI", Command: "dolt", VersionArgs: []string{"version"}, Capabilities: []string{"database", "sql", "git"}},
			{Type: "llm", Name: "LLM CLI", Command: "llm", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "models", "prompt"}},
			{Type: "litellm", Name: "LiteLLM CLI", Command: "litellm", VersionArgs: []string{"--version"}, Capabilities: []string{"models", "proxy", "routing"}},
			{Type: "kimi", Name: "Kimi CLI", Command: "kimi", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "code", "long-context"}},
			{Type: "open-interpreter", Name: "Open Interpreter CLI", Command: "interpreter", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "python", "shell"}},
			{Type: "pi", Name: "Pi CLI", Command: "pi", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "personal", "assistant"}},
			{Type: "qwen", Name: "Qwen Code CLI", Command: "qwen", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "code", "local"}},
			{Type: "cursor", Name: "Cursor CLI", Command: "cursor", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "code", "editor"}},
			{Type: "mistral", Name: "Mistral CLI", Command: "mistral", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "code", "local"}},
			{Type: "git", Name: "Git", Command: "git", VersionArgs: []string{"--version"}, Capabilities: []string{"vcs", "diff", "blame"}},
			{Type: "docker", Name: "Docker", Command: "docker", VersionArgs: []string{"--version"}, Capabilities: []string{"containers", "build", "deploy"}},
			{Type: "npm", Name: "npm", Command: "npm", VersionArgs: []string{"--version"}, Capabilities: []string{"packages", "scripts"}},
			{Type: "cargo", Name: "Cargo", Command: "cargo", VersionArgs: []string{"--version"}, Capabilities: []string{"rust", "build", "packages"}},
			{Type: "pip", Name: "pip", Command: "pip", VersionArgs: []string{"--version"}, Capabilities: []string{"python", "packages"}},
			{Type: "rg", Name: "ripgrep", Command: "rg", VersionArgs: []string{"--version"}, Capabilities: []string{"search", "fast"}},
			{Type: "fd", Name: "fd", Command: "fd", VersionArgs: []string{"--version"}, Capabilities: []string{"find", "fast"}},
		},
	}
}

// DetectAll detects all known tools with TTL caching.
func (d *Detector) DetectAll(ctx context.Context) ([]Tool, error) {
	d.mu.Lock()
	if len(d.cached) > 0 && time.Since(d.cachedAt) < d.ttl {
		tools := append([]Tool(nil), d.cached...)
		d.mu.Unlock()
		return tools, nil
	}

	if d.inflight != nil {
		wait := d.inflight
		d.mu.Unlock()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-wait:
			d.mu.Lock()
			tools := append([]Tool(nil), d.cached...)
			d.mu.Unlock()
			return tools, nil
		}
	}

	wait := make(chan struct{})
	d.inflight = wait
	d.mu.Unlock()

	tools := make([]Tool, 0, len(d.definitions))
	for _, def := range d.definitions {
		tools = append(tools, d.detectTool(ctx, def))
	}

	d.mu.Lock()
	d.cached = append([]Tool(nil), tools...)
	d.cachedAt = time.Now()
	close(wait)
	d.inflight = nil
	d.mu.Unlock()

	return tools, nil
}

// DetectAvailable returns only tools that are installed.
func (d *Detector) DetectAvailable(ctx context.Context) ([]Tool, error) {
	all, err := d.DetectAll(ctx)
	if err != nil {
		return nil, err
	}
	var available []Tool
	for _, t := range all {
		if t.Available {
			available = append(available, t)
		}
	}
	return available, nil
}

// DetectOne detects a single tool by type.
func (d *Detector) DetectOne(ctx context.Context, toolType string) *Tool {
	all, err := d.DetectAll(ctx)
	if err != nil {
		return nil
	}
	for _, t := range all {
		if t.Type == toolType {
			return &t
		}
	}
	return nil
}

func (d *Detector) detectTool(ctx context.Context, def definition) Tool {
	tool := Tool{
		Type:         def.Type,
		Name:         def.Name,
		Command:      def.Command,
		Available:    false,
		Capabilities: append([]string(nil), def.Capabilities...),
	}

	executable, err := exec.LookPath(def.Command)
	if err != nil {
		return tool
	}

	tool.Available = true
	tool.Path = executable

	commandCtx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	output, err := exec.CommandContext(commandCtx, executable, def.VersionArgs...).CombinedOutput()
	if err != nil && commandCtx.Err() != nil {
		tool.Version = "timeout"
		return tool
	}

	tool.Version = parseVersion(string(output))
	return tool
}

var versionPattern = regexp.MustCompile(`v?(\d+\.\d+(?:\.\d+)?)`)

func parseVersion(output string) string {
	match := versionPattern.FindStringSubmatch(output)
	if len(match) >= 2 {
		return match[1]
	}
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return "unknown"
	}
	return trimmed
}
