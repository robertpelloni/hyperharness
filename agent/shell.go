package agent

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// SuggestShellCommand mimics Copilot CLI and Shell Pilot.
// It translates natural language into a runnable shell command.
func (a *Agent) SuggestShellCommand(query string) (string, error) {
	if a == nil || a.client == nil {
		return "", fmt.Errorf("openai client is required")
	}

	prompt := buildShellPrompt(query)
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	}

	resp, err := a.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no shell command choices returned")
	}

	command := resp.Choices[0].Message.Content
	return command, nil
}

func buildShellPrompt(query string) string {
	return fmt.Sprintf("Translate this natural language request into a single, valid shell command for a Windows environment. Output ONLY the command, nothing else. Request: %s", query)
}
