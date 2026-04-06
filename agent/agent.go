package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/robertpelloni/hyperharness/borg"
	"github.com/robertpelloni/hyperharness/foundation/adapters"
	"github.com/robertpelloni/hyperharness/tools"
	"github.com/sashabaranov/go-openai"
)

type Agent struct {
	client       *openai.Client
	messages     []openai.ChatCompletionMessage
	tools        *tools.Registry
	BorgAdapter  *borg.Adapter
	HyperAdapter *adapters.HyperCodeAdapter
}

func NewAgent() *Agent {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = "dummy"
	}

	registry := tools.NewRegistry()
	borgAdapter := borg.NewAdapter()
	cwd, _ := os.Getwd()
	hyperAdapter := adapters.NewHyperCodeAdapter(cwd)
	systemPrompt := buildAgentSystemPrompt(hyperAdapter.BuildSystemContext())

	return &Agent{
		client: openai.NewClient(apiKey),
		messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
		},
		tools:        registry,
		BorgAdapter:  borgAdapter,
		HyperAdapter: hyperAdapter,
	}
}

func (a *Agent) buildOpenAITools() []openai.Tool {
	openAITools := make([]openai.Tool, 0, len(a.tools.Tools))
	for _, t := range a.tools.Tools {
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

	a.messages = append(a.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: a.messages,
		Tools:    a.buildOpenAITools(),
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
	for _, tc := range toolCalls {
		toolResult := executeToolCall(a.tools, tc)
		a.messages = append(a.messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    toolResult,
			Name:       tc.Function.Name,
			ToolCallID: tc.ID,
		})
		summary.WriteString(fmt.Sprintf("Executed tool %s\n", tc.Function.Name))
	}

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: a.messages,
		Tools:    a.buildOpenAITools(),
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

func executeToolCall(registry *tools.Registry, tc openai.ToolCall) string {
	if registry == nil {
		return fmt.Sprintf("Unknown tool: %s", tc.Function.Name)
	}
	var args map[string]interface{}
	_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
	for _, t := range registry.Tools {
		if t.Name != tc.Function.Name {
			continue
		}
		out, err := t.Execute(args)
		if err != nil {
			if out != "" {
				return out
			}
			return fmt.Sprintf("Error executing %s: %v", t.Name, err)
		}
		return out
	}
	return fmt.Sprintf("Unknown tool: %s", tc.Function.Name)
}
