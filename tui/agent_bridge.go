package tui

// ═══════════════════════════════════════════════════════════════════════
// agent_bridge.go — Bridges the real agent.AgentLoop into Bubbletea TUI
// Converts agent.LoopEvent into tea.Msg for real-time streaming display
// ═══════════════════════════════════════════════════════════════════════

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/robertpelloni/hyperharness/agent"
	"github.com/robertpelloni/hyperharness/agents"
	foundationorchestration "github.com/robertpelloni/hyperharness/foundation/orchestration"
	"github.com/robertpelloni/hyperharness/internal/ai"
	"github.com/robertpelloni/hyperharness/tools"
)

// AgentBridge manages a real LLM agent loop and converts its events to Bubbletea messages
type AgentBridge struct {
	mu          sync.Mutex
	loop        *agent.AgentLoop
	program     *tea.Program
	provider    string
	model       string
	workingDir  string
	registry    *tools.Registry
	accumulated string // accumulates assistant text for streaming display
	toolCalls   []ToolCallInfo
}

// ToolCallInfo tracks a single tool call for display
type ToolCallInfo struct {
	Name   string
	Args   string
	Result string
	Error  bool
	Dur    time.Duration
}

// NewAgentBridge creates a new bridge between agent loop and TUI
func NewAgentBridge(provider, model, workingDir string, registry *tools.Registry) *AgentBridge {
	return &AgentBridge{
		provider:   provider,
		model:      model,
		workingDir: workingDir,
		registry:   registry,
	}
}

// SetProgram stores a reference to the Bubbletea program for sending messages
func (b *AgentBridge) SetProgram(p *tea.Program) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.program = p
}

// SetProvider updates the provider and model for user-selected providers
func (b *AgentBridge) SetProvider(provider, model string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.provider = provider
	b.model = model
}

// RunPrompt executes a prompt through the real agent loop and streams events to the TUI
func (b *AgentBridge) RunPrompt(input string) {
	b.mu.Lock()
	b.accumulated = ""
	b.toolCalls = nil
	userProvider := b.provider
	userModel := b.model
	b.mu.Unlock()

	// Ensure we always send a completion message when done
	defer func() {
		b.sendMsg(AgentCompleteMsg{})
	}()

	var provider ai.Provider
	var providerName, modelName string
	var err error

	// If user selected a provider via /provider, use that (no fallback)
	if userProvider != "" && userProvider != "hyperharness" {
		provider, providerName, modelName, err = agent.CreateProvider(userProvider)
		if err != nil || provider == nil {
			b.sendMsg(AgentResponseMsg{
				Content: fmt.Sprintf("No API key for %s. Set %s_API_KEY environment variable.", userProvider, strings.ToUpper(userProvider)),
				Provider: userProvider,
				Model:    userModel,
			})
			return
		}
	} else {
		// No user selection, resolve from env
		provider, providerName, modelName, err = agent.ResolveProvider()
		if err != nil || provider == nil {
			b.sendMsg(AgentResponseMsg{
				Content: fmt.Sprintf("No LLM provider configured. Set one of: ANTHROPIC_API_KEY, OPENAI_API_KEY, GOOGLE_API_KEY, GEMINI_API_KEY, DEEPSEEK_API_KEY, OPENROUTER_API_KEY, XIAOMI_API_KEY\n\nError: %v", err),
				Provider: "none",
				Model:    "none",
			})
			return
		}
	}

	// Override model if user selected one
	if userModel != "" && userModel != "auto" {
		modelName = userModel
	}

	registry := b.registry
	if registry == nil {
		registry = tools.NewRegistry()
	}

	loop := agent.NewAgentLoop(agent.AgentLoopConfig{
		Provider:     provider,
		ProviderName: providerName,
		Model:        modelName,
		WorkingDir:   b.workingDir,
		Registry:     registry,
	})

	// Subscribe to events and convert to TUI messages
	loop.OnEvent(func(event agent.LoopEvent) {
		switch event.Type {
		case agent.LoopEventMessageStart:
			b.sendMsg(ThinkingStartMsg{
				Provider: event.Provider,
				Model:    event.Model,
			})

		case agent.LoopEventMessageDelta:
			b.sendMsg(StreamChunkMsg{Chunk: event.Content})
		case agent.LoopEventMessageEnd:
			b.mu.Lock()
			b.accumulated = event.Content
			b.mu.Unlock()
			b.sendMsg(AgentResponseMsg{
				Content:  event.Content,
				Provider: event.Provider,
				Model:    event.Model,
				InputTok: event.InputTok,
				OutTok:   event.OutputTok,
				Cost:     event.Cost,
			})

		case agent.LoopEventToolCallStart:
			b.sendMsg(ToolExecMsg{
				ToolName: event.ToolName,
				Args:     event.ToolArgs,
				Running:  true,
			})

		case agent.LoopEventToolCallEnd:
			b.mu.Lock()
			b.toolCalls = append(b.toolCalls, ToolCallInfo{
				Name:   event.ToolName,
				Result: event.ToolResult,
				Error:  event.ToolError,
				Dur:    event.ToolDur,
			})
			b.mu.Unlock()
			b.sendMsg(ToolExecMsg{
				ToolName: event.ToolName,
				Args:     event.ToolArgs,
				Output:   event.ToolResult,
				IsError:  event.ToolError,
				Duration: event.ToolDur,
				Running:  false,
			})

		case agent.LoopEventError:
			b.sendMsg(AgentResponseMsg{
				Content:  fmt.Sprintf("Error: %v", event.Err),
				Provider: providerName,
				Model:    modelName,
			})

		case agent.LoopEventComplete:
			b.sendMsg(AgentCompleteMsg{
				Content: event.Content,
			})
		}
	})

	// Run the agent loop with timeout (blocking — but events stream via the listener)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	result, err := loop.Run(ctx, input)
	if err != nil {
		b.sendMsg(AgentResponseMsg{
			Content:  fmt.Sprintf("Agent error: %v", err),
			Provider: providerName,
			Model:    modelName,
		})
		return
	}

	_ = result // result is already sent via events
}

// RunDirectorFallback uses the Director (no real LLM) when no API key is available
func (b *AgentBridge) RunDirectorFallback(director *agents.Director, input string) {
	response, err := director.HandleInput(context.Background(), input)
	if err != nil {
		b.sendMsg(AgentResponseMsg{
			Content:  fmt.Sprintf("Error: %v", err),
			Provider: b.provider,
			Model:    b.model,
		})
		return
	}

	if plan, ok := director.State["lastPlan"].(foundationorchestration.PlanResult); ok {
		b.sendMsg(AgentResponseMsg{
			Content:  fmt.Sprintf("[Foundation Route] %s/%s\n%s", plan.Execution.Route.Provider, plan.Execution.Route.Model, response),
			Provider: plan.Execution.Route.Provider,
			Model:    plan.Execution.Route.Model,
		})
	} else {
		b.sendMsg(AgentResponseMsg{
			Content:  response,
			Provider: b.provider,
			Model:    b.model,
		})
	}
}

func (b *AgentBridge) sendMsg(msg tea.Msg) {
	b.mu.Lock()
	p := b.program
	b.mu.Unlock()
	if p != nil {
		p.Send(msg)
	}
}

// HasProvider returns true if a real LLM API key is configured
func (b *AgentBridge) HasProvider() bool {
	_, _, _, err := agent.ResolveProvider()
	return err == nil
}

// ToolCallsSummary returns a summary of tool calls made during the last prompt
func (b *AgentBridge) ToolCallsSummary() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.toolCalls) == 0 {
		return ""
	}
	var parts []string
	for _, tc := range b.toolCalls {
		icon := "✓"
		if tc.Error {
			icon = "✗"
		}
		result := tc.Result
		if len(result) > 80 {
			result = result[:77] + "..."
		}
		result = strings.ReplaceAll(result, "\n", " ")
		parts = append(parts, fmt.Sprintf("%s %s: %s (%v)", icon, tc.Name, result, tc.Dur.Round(time.Millisecond)))
	}
	return strings.Join(parts, "\n")
}
