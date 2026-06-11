package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/robertpelloni/hyperharness/foundation/adapters"
	"github.com/robertpelloni/hyperharness/internal/ai"
	"github.com/robertpelloni/hyperharness/internal/cache"
	"github.com/robertpelloni/hyperharness/internal/git"
	"github.com/robertpelloni/hyperharness/internal/memory"
	"github.com/robertpelloni/hyperharness/tools"
)

// LoopEventType identifies agent loop events
type LoopEventType string

const (
	LoopEventMessageStart  LoopEventType = "message_start"
	LoopEventMessageDelta  LoopEventType = "message_delta"
	LoopEventMessageEnd    LoopEventType = "message_end"
	LoopEventToolCallStart LoopEventType = "tool_call_start"
	LoopEventToolCallEnd   LoopEventType = "tool_call_end"
	LoopEventThinkingStart LoopEventType = "thinking_start"
	LoopEventThinkingEnd   LoopEventType = "thinking_end"
	LoopEventError         LoopEventType = "error"
	LoopEventComplete      LoopEventType = "complete"
)

// LoopEvent is emitted by the agent loop during execution
type LoopEvent struct {
	Type       LoopEventType `json:"type"`
	Content    string        `json:"content,omitempty"`
	ToolName   string        `json:"toolName,omitempty"`
	ToolArgs   string        `json:"toolArgs,omitempty"`
	ToolResult string        `json:"toolResult,omitempty"`
	ToolError  bool          `json:"toolError,omitempty"`
	ToolDur    time.Duration `json:"toolDur,omitempty"`
	Provider   string        `json:"provider,omitempty"`
	Model      string        `json:"model,omitempty"`
	InputTok   int           `json:"inputTok,omitempty"`
	OutputTok  int           `json:"outputTok,omitempty"`
	Cost       float64       `json:"cost,omitempty"`
	Err        error         `json:"err,omitempty"`
}

// LoopEventListener receives agent loop events
type LoopEventListener func(event LoopEvent)

// ParsedToolCall is a tool call extracted from LLM output
type ParsedToolCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

// ToolResult is the result of executing a tool
type ToolResult struct {
	Name   string        `json:"name"`
	Result string        `json:"result"`
	Error  bool          `json:"error"`
	Dur    time.Duration `json:"dur"`
}

// AgentLoopConfig configures a new agent loop
type AgentLoopConfig struct {
	Provider     ai.Provider
	ProviderName string
	Model        string
	WorkingDir   string
	Registry     *tools.Registry
	MaxTurns     int
}

// AgentLoop is a provider-agnostic agent loop with tool calling
type AgentLoop struct {
	mu           sync.RWMutex
	provider     ai.Provider
	providerName string
	modelName    string
	messages     []ai.Message
	registry     *tools.Registry
	toolCache    *cache.Cache
	hyperAdapter *adapters.HyperCodeAdapter
	listeners    []LoopEventListener
	maxTurns     int
	totalInTok   int
	totalOutTok  int
	totalCost    float64
	aborted      bool
	workingDir   string
	kb           *memory.KnowledgeBase
}

// NewAgentLoop creates a new universal agent loop
func NewAgentLoop(cfg AgentLoopConfig) *AgentLoop {
	if cfg.MaxTurns == 0 {
		cfg.MaxTurns = 50
	}
	if cfg.WorkingDir == "" {
		cfg.WorkingDir, _ = os.Getwd()
	}
	registry := cfg.Registry
	if registry == nil {
		registry = tools.NewRegistry()
	}
	hyperAdapter := adapters.NewHyperCodeAdapter(cfg.WorkingDir)
	kb, _ := memory.NewKnowledgeBase("")
	systemPrompt := buildUniversalSystemPrompt(cfg.WorkingDir, hyperAdapter)
	return &AgentLoop{
		provider:     cfg.Provider,
		providerName: cfg.ProviderName,
		modelName:    cfg.Model,
		registry:     registry,
		toolCache:    cache.New(cache.CacheOptions{MaxSize: 1000, DefaultTTL: 300000}),
		hyperAdapter: hyperAdapter,
		maxTurns:     cfg.MaxTurns,
		workingDir:   cfg.WorkingDir,
		kb:           kb,
		messages:     []ai.Message{{Role: "system", Content: systemPrompt}},
	}
}

// OnEvent registers an event listener
func (l *AgentLoop) OnEvent(listener LoopEventListener) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.listeners = append(l.listeners, listener)
}

func (l *AgentLoop) emit(event LoopEvent) {
	l.mu.RLock()
	listeners := l.listeners
	l.mu.RUnlock()
	for _, fn := range listeners {
		fn(event)
	}
}

// CreateProvider creates a provider by name using the appropriate API key from env
func CreateProvider(providerName string) (ai.Provider, string, string, error) {
	type c struct {
		pName, dModel string
		envVars       []string
		factory       func(string) ai.Provider
	}
	providers := []c{
		{"anthropic", "claude-sonnet-4-20250514", []string{"ANTHROPIC_API_KEY", "CLAUDE_API_KEY"}, func(k string) ai.Provider { return &ai.AnthropicProvider{APIKey: k} }},
		{"openai", "gpt-4o", []string{"OPENAI_API_KEY"}, func(k string) ai.Provider { return &ai.OpenAIProvider{APIKey: k} }},
		{"google", "gemini-2.5-flash", []string{"GOOGLE_API_KEY", "GEMINI_API_KEY"}, func(k string) ai.Provider { return &ai.GeminiProvider{APIKey: k} }},
		{"deepseek", "deepseek-chat", []string{"DEEPSEEK_API_KEY"}, func(k string) ai.Provider { return &ai.DeepSeekProvider{APIKey: k} }},
		{"openrouter", "anthropic/claude-3.5-sonnet", []string{"OPENROUTER_API_KEY"}, func(k string) ai.Provider { return &ai.OpenRouterProvider{APIKey: k} }},
		{"xiaomi", "mimo-v2.5-pro", []string{"XIAOMI_API_KEY", "MIMO_API_KEY"}, func(k string) ai.Provider { return &ai.XiaomiProvider{APIKey: k} }},
		{"ollama", "llama3:8b", []string{}, func(_ string) ai.Provider { return &ai.OllamaProvider{} }},
		{"lmstudio", "local-model", []string{}, func(_ string) ai.Provider { return &ai.LMStudioProvider{} }},
	}
	for _, p := range providers {
		if p.pName == providerName {
			for _, env := range p.envVars {
				if key := os.Getenv(env); key != "" {
					return p.factory(key), p.pName, p.dModel, nil
				}
			}
			// For local providers (ollama, lmstudio), no key needed
			if len(p.envVars) == 0 {
				return p.factory(""), p.pName, p.dModel, nil
			}
			return nil, "", "", fmt.Errorf("no API key for %s (set %s)", providerName, strings.Join(p.envVars, " or "))
		}
	}
	return nil, "", "", fmt.Errorf("unknown provider: %s", providerName)
}

// ResolveProvider creates a provider based on environment variables
func ResolveProvider() (ai.Provider, string, string, error) {
	type c struct {
		envVar, pName, dModel string
		factory               func(string) ai.Provider
	}
	for _, c := range []c{
		{"ANTHROPIC_API_KEY", "anthropic", "claude-sonnet-4-20250514", func(k string) ai.Provider { return &ai.AnthropicProvider{APIKey: k} }},
		{"CLAUDE_API_KEY", "anthropic", "claude-sonnet-4-20250514", func(k string) ai.Provider { return &ai.AnthropicProvider{APIKey: k} }},
		{"GOOGLE_API_KEY", "google", "gemini-2.5-flash", func(k string) ai.Provider { return &ai.GeminiProvider{APIKey: k} }},
		{"GEMINI_API_KEY", "google", "gemini-2.5-flash", func(k string) ai.Provider { return &ai.GeminiProvider{APIKey: k} }},
		{"OPENAI_API_KEY", "openai", "gpt-4o", func(k string) ai.Provider { return &ai.OpenAIProvider{APIKey: k} }},
		{"DEEPSEEK_API_KEY", "deepseek", "deepseek-chat", func(k string) ai.Provider { return &ai.DeepSeekProvider{APIKey: k} }},
		{"OPENROUTER_API_KEY", "openrouter", "anthropic/claude-3.5-sonnet", func(k string) ai.Provider { return &ai.OpenRouterProvider{APIKey: k} }},
		{"XIAOMI_API_KEY", "xiaomi", "mimo-v2.5-pro", func(k string) ai.Provider { return &ai.XiaomiProvider{APIKey: k} }},
		{"MIMO_API_KEY", "xiaomi", "mimo-v2.5-pro", func(k string) ai.Provider { return &ai.XiaomiProvider{APIKey: k} }},
	} {
		if key := os.Getenv(c.envVar); key != "" {
			return c.factory(key), c.pName, c.dModel, nil
		}
	}
	return nil, "", "", fmt.Errorf("no LLM provider configured")
}

// Run executes the agent loop with the given user input
func (l *AgentLoop) Run(ctx context.Context, input string) (string, error) {
	l.mu.Lock()
	l.aborted = false
	l.mu.Unlock()

	l.messages = append(l.messages, ai.Message{Role: "user", Content: input})
	var finalContent string

	for turn := 0; turn < l.maxTurns; turn++ {
		if l.aborted {
			l.emit(LoopEvent{Type: LoopEventComplete, Content: finalContent})
			return finalContent, nil
		}
		select {
		case <-ctx.Done():
			return finalContent, ctx.Err()
		default:
		}

		l.emit(LoopEvent{Type: LoopEventMessageStart, Provider: l.providerName, Model: l.modelName})
		resp, err := l.provider.GenerateText(ctx, l.modelName, l.messages)
		if err != nil {
			l.emit(LoopEvent{Type: LoopEventError, Err: err})
			return finalContent, fmt.Errorf("LLM: %w", err)
		}

		l.totalInTok += resp.Usage.InputTokens
		l.totalOutTok += resp.Usage.OutputTokens
		l.totalCost += estimateCost(l.providerName, l.modelName, resp.Usage.InputTokens, resp.Usage.OutputTokens)

		l.emit(LoopEvent{
			Type:      LoopEventMessageEnd,
			Content:   resp.Content,
			Provider:  resp.Provider,
			Model:     resp.Model,
			InputTok:  resp.Usage.InputTokens,
			OutputTok: resp.Usage.OutputTokens,
			Cost:      l.totalCost,
		})

		l.messages = append(l.messages, ai.Message{Role: "assistant", Content: resp.Content})

		toolCalls := parseToolCalls(resp.Content)
		if len(toolCalls) == 0 {
			finalContent = resp.Content
			break
		}

		toolResults := l.executeToolCalls(ctx, toolCalls)
		var rp []string
		for _, tr := range toolResults {
			rp = append(rp, fmt.Sprintf("Tool: %s\nResult: %s", tr.Name, tr.Result))
		}
		l.messages = append(l.messages, ai.Message{Role: "user", Content: strings.Join(rp, "\n\n")})
		finalContent = resp.Content
	}

	l.emit(LoopEvent{Type: LoopEventComplete, Content: finalContent})
	return finalContent, nil
}

// Abort signals the loop to stop
func (l *AgentLoop) Abort() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.aborted = true
}

// Stats returns usage statistics
func (l *AgentLoop) Stats() (int, int, float64) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.totalInTok, l.totalOutTok, l.totalCost
}

// ListAvailableTools returns the names of all registered tools
func (l *AgentLoop) ListAvailableTools() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	names := make([]string, 0, len(l.registry.Tools))
	for _, t := range l.registry.Tools {
		names = append(names, t.Name)
	}
	return names
}

// Messages returns a copy of the conversation history
func (l *AgentLoop) Messages() []ai.Message {
	l.mu.RLock()
	defer l.mu.RUnlock()
	cp := make([]ai.Message, len(l.messages))
	copy(cp, l.messages)
	return cp
}

// parseToolCalls extracts tool calls from LLM text output
// Supports XML-style tags around JSON tool call specifications
func parseToolCalls(content string) []ParsedToolCall {
	var calls []ParsedToolCall
	// Tool calls are wrapped in special tags (built at runtime to avoid parsing issues)
	opentag := string([]byte{60}) + "tool_call" + string([]byte{62})
	closetag := string([]byte{60}) + "/tool_call" + string([]byte{62})
	idx := 0
	for {
		start := strings.Index(content[idx:], opentag)
		if start == -1 {
			break
		}
		start += idx + len(opentag)
		end := strings.Index(content[start:], closetag)
		if end == -1 {
			break
		}
		callJSON := strings.TrimSpace(content[start : start+end])
		idx = start + end + len(closetag)
		var call struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if err := json.Unmarshal([]byte(callJSON), &call); err == nil && call.Name != "" {
			calls = append(calls, ParsedToolCall{Name: call.Name, Args: call.Arguments})
		}
	}
	return calls
}

// executeToolCalls runs tool calls in parallel using goroutines
func (l *AgentLoop) executeToolCalls(ctx context.Context, calls []ParsedToolCall) []ToolResult {
	results := make([]ToolResult, len(calls))
	var wg sync.WaitGroup
	for i, call := range calls {
		wg.Add(1)
		go func(idx int, tc ParsedToolCall) {
			defer wg.Done()
			l.emit(LoopEvent{Type: LoopEventToolCallStart, ToolName: tc.Name, ToolArgs: formatArgs(tc.Args)})
			start := time.Now()
			result, isErr := l.executeTool(ctx, tc)
			dur := time.Since(start)
			results[idx] = ToolResult{Name: tc.Name, Result: result, Error: isErr, Dur: dur}
			l.emit(LoopEvent{Type: LoopEventToolCallEnd, ToolName: tc.Name, ToolResult: result, ToolError: isErr, ToolDur: dur})
		}(i, call)
	}
	wg.Wait()
	return results
}

// executeTool runs a single tool with caching for read-only tools
func (l *AgentLoop) executeTool(ctx context.Context, call ParsedToolCall) (string, bool) {
	readOnly := map[string]bool{
		"read": true, "read_file": true, "file_read": true,
		"ls": true, "list_directory": true, "tree": true,
		"grep": true, "find": true, "search": true,
		"websearch": true, "webfetch": true, "memory_search": true,
	}
	if readOnly[call.Name] && l.toolCache != nil {
		cacheKey := call.Name + ":" + formatArgs(call.Args)
		if cached, ok := l.toolCache.Get(cacheKey); ok {
			if s, ok := cached.(string); ok {
				return s, false
			}
		}
	}
	tool, ok := l.registry.Find(call.Name)
	if !ok {
		return fmt.Sprintf("Unknown tool: %s", call.Name), true
	}
	if p, ok := call.Args["path"].(string); ok {
		if !l.isPathAllowed(p) {
			return "Access denied: path outside working dir", true
		}
	}
	out, err := tool.Execute(call.Args)
	if err != nil {
		if out != "" {
			return out, true
		}
		return fmt.Sprintf("Error: %v", err), true
	}
	if readOnly[call.Name] && l.toolCache != nil {
		l.toolCache.Set(call.Name+":"+formatArgs(call.Args), out, 5*time.Minute)
	}
	return out, false
}

// isPathAllowed checks if a path is within the working directory
func (l *AgentLoop) isPathAllowed(path string) bool {
	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "\\") && !strings.Contains(path, ":") {
		return true
	}
	return strings.HasPrefix(path, l.workingDir)
}

// buildUniversalSystemPrompt creates the system prompt with tool instructions
func buildUniversalSystemPrompt(workingDir string, adapter *adapters.HyperCodeAdapter) string {
	baseCtx := adapter.BuildSystemContext()
	gitCtx, _ := git.GetAwarenessContext(workingDir)
	if gitCtx != "" {
		baseCtx = baseCtx + "\n\n" + gitCtx
	}
	ot := string([]byte{60}) + "tool_call" + string([]byte{62})
	ct := string([]byte{60}) + "/tool_call" + string([]byte{62})
	return strings.Join([]string{
		"You are HyperCode, a Go-native coding and terminal assistant.",
		"You have access to tools. Invoke them using XML tool_call tags.",
		"To call a tool: " + ot + `{"name":"tool_name","arguments":{}}` + ct,
		"You may call multiple tools. After receiving results, continue reasoning.",
		"Available: read, write, edit, bash, grep, ls, find, tree, glob, websearch, webfetch, memory_search, git_status, git_diff, git_log",
		"Rules: Use read before editing. Use edit for precise edits. Verify changes. Report errors honestly.",
		baseCtx,
	}, "\n")
}

// formatArgs serializes tool arguments to JSON
func formatArgs(args map[string]interface{}) string {
	data, _ := json.Marshal(args)
	return string(data)
}

// estimateCost provides rough cost estimates per provider
func estimateCost(provider, model string, inputTok, outputTok int) float64 {
	costs := map[string][2]float64{
		"anthropic": {3.0, 15.0},
		"openai":    {2.5, 10.0},
		"google":    {0.075, 0.3},
		"deepseek":  {0.14, 0.28},
		"openrouter": {3.0, 15.0},
	}
	pc, ok := costs[provider]
	if !ok {
		return 0
	}
	return float64(inputTok)/1e6*pc[0] + float64(outputTok)/1e6*pc[1]
}
