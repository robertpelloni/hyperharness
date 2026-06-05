package agent

import (
	"github.com/robertpelloni/hyperharness/internal/providers"
	"github.com/robertpelloni/hyperharness/internal/cache"
	"sync"
	"github.com/robertpelloni/hyperharness/internal/subagents"
	"github.com/robertpelloni/hyperharness/internal/memory"
	"github.com/robertpelloni/hyperharness/internal/ingest"
	"github.com/robertpelloni/hyperharness/internal/git"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/robertpelloni/hyperharness/foundation/adapters"
	"github.com/robertpelloni/hyperharness/tools"
	"github.com/sashabaranov/go-openai"
)

type Agent struct {
	client       *openai.Client
	messages     []openai.ChatCompletionMessage
	tools        *tools.Registry
	HyperAdapter *adapters.HyperCodeAdapter
	toolCache    *cache.Cache
	processor    *ingest.DataProcessor
}

func NewAgent() *Agent {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = "dummy"
	}

	registry := tools.NewRegistry()
	cwd, _ := os.Getwd()
	hyperAdapter := adapters.NewHyperCodeAdapter(cwd)
	
	gitCtx, _ := git.GetAwarenessContext(cwd)
	baseCtx := hyperAdapter.BuildSystemContext()
	if gitCtx != "" {
		baseCtx = baseCtx + "\n\n" + gitCtx
	}
	
	systemPrompt := buildAgentSystemPrompt(baseCtx)

	// Initialize data processor with a temporary knowledge base if one isn't provided
	// In a real implementation, this would be passed from the session/context
	kb, _ := memory.NewKnowledgeBase("")
	processor := ingest.NewDataProcessor(kb)

	config := openai.DefaultConfig(apiKey)
	config.HTTPClient = providers.GetPooledHTTPClient()
	
	return &Agent{
		client: openai.NewClientWithConfig(config),
		messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
		},
		tools:        registry,
		HyperAdapter: hyperAdapter,
		toolCache:    cache.New(cache.CacheOptions{MaxSize: 1000, DefaultTTL: 300000}), // 5 min TTL
		processor:    processor,
	}
}

func (a *Agent) buildOpenAITools(st subagents.SubagentType) []openai.Tool {
	openAITools := make([]openai.Tool, 0)
	for _, t := range a.tools.Tools {
		if st != "" && !subagents.IsToolAllowed(st, t.Name) {
			continue
		}
		openAITools = append(openAITools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  append(json.RawMessage(nil), t.Parameters...),
			},
		})
	}
	return openAITools
}

func (a *Agent) Chat(input string) (string, error) {
	if a == nil {
		return "", fmt.Errorf("agent is required")
	}
	if a.client == nil {
		return "", fmt.Errorf("openai client is required")
	}
	if a.tools == nil {
		return "", fmt.Errorf("tool registry is required")
	}

	// Normalize input using the data processor
	normalizedInput := input
	if a.processor != nil {
		normalizedInput = a.processor.Normalize(input)
	}

	a.messages = append(a.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: normalizedInput,
	})

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: a.messages,
		Tools:    a.buildOpenAITools(""),
	}

	resp, err := a.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}

	msg, err := firstChoiceMessage(resp)
	if err != nil {
		return "", err
	}
	a.messages = append(a.messages, msg)

	if len(msg.ToolCalls) > 0 {
		return a.handleToolCalls(msg.ToolCalls)
	}

	return msg.Content, nil
}

func (a *Agent) handleToolCalls(toolCalls []openai.ToolCall) (string, error) {
	if a == nil {
		return "", fmt.Errorf("agent is required")
	}
	if a.client == nil {
		return "", fmt.Errorf("openai client is required")
	}
	if a.tools == nil {
		return "", fmt.Errorf("tool registry is required")
	}

	var summary strings.Builder
	results := make([]openai.ChatCompletionMessage, len(toolCalls))
	var wg sync.WaitGroup
	var mu sync.Mutex // For summary builder
	
	for i, tc := range toolCalls {
		wg.Add(1)
		go func(index int, call openai.ToolCall) {
			defer wg.Done()
			toolResult := executeToolCall(a.tools, call, a.toolCache)
			
			results[index] = openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    toolResult,
				Name:       call.Function.Name,
				ToolCallID: call.ID,
			}
			
			mu.Lock()
			summary.WriteString(fmt.Sprintf("Executed tool %s\n", call.Function.Name))
			mu.Unlock()
		}(i, tc)
	}
	
	wg.Wait()
	a.messages = append(a.messages, results...)


	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: a.messages,
		Tools:    a.buildOpenAITools(""),
	}

	resp, err := a.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return summary.String(), err
	}

	msg, err := firstChoiceMessage(resp)
	if err != nil {
		return summary.String(), err
	}
	a.messages = append(a.messages, msg)

	if summary.Len() == 0 {
		return msg.Content, nil
	}
	return summary.String() + "\n" + msg.Content, nil
}

func buildAgentSystemPrompt(systemContext string) string {
	return strings.Join([]string{
		"You are Hypercode, a Go-native coding and terminal assistant integrated with Borg and HyperCode.",
		"Prefer the exact-name Pi-compatible tools read, write, edit, and bash when solving coding tasks.",
		"Use repomap for repository-wide context when a condensed map would help.",
		"Additional legacy tools may exist for compatibility, but exact-contract tools are preferred.",
		systemContext,
	}, "\n\n")
}

func firstChoiceMessage(resp openai.ChatCompletionResponse) (openai.ChatCompletionMessage, error) {
	if len(resp.Choices) == 0 {
		return openai.ChatCompletionMessage{}, fmt.Errorf("no completion choices returned")
	}
	return resp.Choices[0].Message, nil
}

func executeToolCall(registry *tools.Registry, tc openai.ToolCall, toolCache *cache.Cache) string {
	if registry == nil {
		return fmt.Sprintf("Unknown tool: %s", tc.Function.Name)
	}
	var args map[string]interface{}
	_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
	tool, ok := registry.Find(tc.Function.Name)
	if !ok {
		return fmt.Sprintf("Unknown tool: %s", tc.Function.Name)
	}
	
	// Cache read-only tools
	isCacheable := false
	cacheKey := ""
	
	readOnlyTools := map[string]bool{
		"read": true, "read_file": true, "file_read": true, "ls": true, "list_directory": true, "tree": true,
		"grep": true, "find": true, "search": true, "search_code": true, "search_files": true, "codesearch": true,
		"websearch": true, "webfetch": true, "web_search": true, "web_fetch": true, "memory_search": true,
		"get_system_stats": true, "diagnostics": true,
	}
	
	if readOnlyTools[tool.Name] && toolCache != nil {
		isCacheable = true
		// Make sure it's fully deterministic, key on toolName + exact args JSON string
		cacheKey = tool.Name + ":" + tc.Function.Arguments
		
		if cachedOut, ok := toolCache.Get(cacheKey); ok {
			if strOut, ok := cachedOut.(string); ok {
				return strOut
			}
		}
	}

	out, err := tool.Execute(args)
	if err != nil {
		if out != "" {
			return out
		}
		return fmt.Sprintf("Error executing %s: %v", tool.Name, err)
	}
	
	if isCacheable && toolCache != nil {
		toolCache.Set(cacheKey, out, 5*time.Minute)
	}
	
	return out
}
