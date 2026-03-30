package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/robertpelloni/hypercode/borg"
	"github.com/robertpelloni/hypercode/tools"
	"github.com/sashabaranov/go-openai"
)

type Agent struct {
	client      *openai.Client
	messages    []openai.ChatCompletionMessage
	tools       *tools.Registry
	BorgAdapter *borg.Adapter
}

func NewAgent() *Agent {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = "dummy"
	}

	registry := tools.NewRegistry()
	borgAdapter := borg.NewAdapter()

	return &Agent{
		client: openai.NewClient(apiKey),
		messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are Hypercode (formerly SuperCLI), the world's most advanced AI CLI assistant. You have full feature parity with Aider, Claude Code, and Open Interpreter. You can: 1) Execute shell commands. 2) Read/Write files. 3) Use a stateful Code Interpreter for Python/Node.js. 4) Search the codebase. 5) View an AST-lite Repository Map. Use these tools aggressively to solve the user's problems. You are ASSIMILATED by Borg for memory and context management.",
			},
		},
		tools:       registry,
		BorgAdapter: borgAdapter,
	}
}

func (a *Agent) buildOpenAITools() []openai.Tool {
	var openAITools []openai.Tool
	for _, t := range a.tools.Tools {
		// Using a simplified schema for demonstration
		openAITools = append(openAITools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        t.Name,
				Description: t.Description,
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"command": { "type": "string" },
						"file_path": { "type": "string" },
						"content": { "type": "string" }
					}
				}`),
			},
		})
	}
	return openAITools
}

func (a *Agent) Chat(input string) (string, error) {
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

	msg := resp.Choices[0].Message
	a.messages = append(a.messages, msg)

	if len(msg.ToolCalls) > 0 {
		return a.handleToolCalls(msg.ToolCalls)
	}

	return msg.Content, nil
}

func (a *Agent) handleToolCalls(toolCalls []openai.ToolCall) (string, error) {
	resultSummary := ""

	for _, tc := range toolCalls {
		var args map[string]interface{}
		json.Unmarshal([]byte(tc.Function.Arguments), &args)

		var toolResult string
		for _, t := range a.tools.Tools {
			if t.Name == tc.Function.Name {
				out, err := t.Execute(args)
				if err != nil {
					toolResult = fmt.Sprintf("Error executing %s: %v", t.Name, err)
				} else {
					toolResult = out
				}
				break
			}
		}

		a.messages = append(a.messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    toolResult,
			Name:       tc.Function.Name,
			ToolCallID: tc.ID,
		})

		resultSummary += fmt.Sprintf("Executed tool %s\n", tc.Function.Name)
	}

	// Recurse to get the model's final response after tool execution
	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: a.messages,
		Tools:    a.buildOpenAITools(),
	}

	resp, err := a.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return resultSummary, err
	}

	msg := resp.Choices[0].Message
	a.messages = append(a.messages, msg)

	return resultSummary + "\n" + msg.Content, nil
}
