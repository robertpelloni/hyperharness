package agents

import (
	"context"
)

// DefaultProvider simulates the LLM locally.
type DefaultProvider struct{}

func (p *DefaultProvider) Chat(ctx context.Context, messages []Message, tools []Tool) (Message, error) {
	return Message{
		Role:    RoleAssistant,
		Content: "I am the new Native Go Borg Director. How can I autonomously assist you today?",
	}, nil
}

func (p *DefaultProvider) Stream(ctx context.Context, messages []Message, tools []Tool, chunkChan chan<- string) error {
	chunkChan <- "I am the "
	chunkChan <- "Native Go Borg Director."
	close(chunkChan)
	return nil
}

func (p *DefaultProvider) GetModelName() string {
	return "borg-native-stub-1.0"
}
