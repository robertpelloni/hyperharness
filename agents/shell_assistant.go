package agents

import (
	"context"
	"fmt"
	"strings"
)

// ShellTranslator represents a strictly constrained Agent prompt designed to output ONLY bash/pwsh.
type ShellTranslator struct {
	Provider ILLMProvider
}

func NewShellTranslator(provider ILLMProvider) *ShellTranslator {
	return &ShellTranslator{
		Provider: provider,
	}
}

// Translate converts natural language into shell syntax natively.
func (s *ShellTranslator) Translate(ctx context.Context, intent string) (string, error) {
	if s == nil || s.Provider == nil {
		return "", fmt.Errorf("provider is required")
	}
	intent = strings.TrimSpace(intent)
	if intent == "" {
		return "", fmt.Errorf("intent is required")
	}

	resp, err := s.Provider.Chat(ctx, buildShellTranslationMessages(intent), nil)
	if err != nil {
		return "", fmt.Errorf("shell LLM translation error: %w", err)
	}
	return resp.Content, nil
}

func buildShellTranslationMessages(intent string) []Message {
	return []Message{
		{
			Role:    RoleSystem,
			Content: "You are the native Copilot CLI equivalence engine. Translate the user intent into a single native shell command. Do not wrap in markdown quotes. Output ONLY the raw executable text.",
		},
		{
			Role:    RoleUser,
			Content: intent,
		},
	}
}
