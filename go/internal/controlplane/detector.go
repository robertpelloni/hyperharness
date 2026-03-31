package controlplane

import (
	"context"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Tool struct {
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Command      string   `json:"command"`
	Available    bool     `json:"available"`
	Version      string   `json:"version,omitempty"`
	Path         string   `json:"path,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
}

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

type Detector struct {
	definitions []definition
	timeout     time.Duration
	ttl         time.Duration

	mu       sync.Mutex
	inflight chan struct{}
	cached   []Tool
	cachedAt time.Time
}

func NewDetector(timeout, ttl time.Duration) *Detector {
	return &Detector{
		timeout: timeout,
		ttl:     ttl,
		definitions: []definition{
			{Type: "go", Name: "Go", Command: "go", VersionArgs: []string{"version"}, Capabilities: []string{"build", "test", "server"}},
			{Type: "node", Name: "Node.js", Command: "node", VersionArgs: []string{"--version"}, Capabilities: []string{"runtime", "scripts"}},
			{Type: "python", Name: "Python", Command: "python", VersionArgs: []string{"--version"}, Capabilities: []string{"runtime", "scripts"}},
			{Type: "codex", Name: "Codex CLI", Command: "codex", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "code"}},
			{Type: "opencode", Name: "OpenCode CLI", Command: "opencode", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "edit"}},
			{Type: "gemini", Name: "Gemini CLI", Command: "gemini", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "multimodal"}},
			{Type: "goose", Name: "Goose CLI", Command: "goose", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "automation"}},
			{Type: "claude-code", Name: "Claude Code CLI", Command: "claude-code", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "code"}},
			{Type: "aider", Name: "Aider CLI", Command: "aider", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "git-aware"}},
			{Type: "qwen", Name: "Qwen CLI", Command: "qwen", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "code"}},
			{Type: "superai-cli", Name: "SuperAI CLI", Command: "superai", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "code"}},
			{Type: "codebuff", Name: "Codebuff CLI", Command: "codebuff", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "refactor"}},
			{Type: "codemachine", Name: "Codemachine CLI", Command: "codemachine", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "generate"}},
			{Type: "factory-droid", Name: "Factory Droid CLI", Command: "droid", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "automation"}},
			{Type: "cursor", Name: "Cursor CLI", Command: "cursor", VersionArgs: []string{"--version"}, Capabilities: []string{"editor", "chat"}},
			{Type: "copilot", Name: "GitHub Copilot CLI", Command: "github-copilot-cli", VersionArgs: []string{"--version"}, Capabilities: []string{"chat", "terminal"}},
		},
	}
}

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
