package agents

import (
	"context"
	"fmt"

	"github.com/robertpelloni/hyperharness/foundation/adapters"
)

// DefaultProvider simulates the LLM locally.
type DefaultProvider struct{}

func (p *DefaultProvider) Chat(ctx context.Context, messages []Message, tools []Tool) (Message, error) {
	if p == nil {
		return Message{}, fmt.Errorf("provider is required")
	}
	prompt := lastProviderMessageContent(messages)
	execution := adapters.PrepareProviderExecution(adapters.ProviderExecutionRequest{Prompt: prompt, CostPreference: "budget"})
	return Message{
		Role:    RoleAssistant,
		Content: fmt.Sprintf("I am the new Native Go Borg Director. %s", execution.ExecutionHint),
	}, nil
}

func (p *DefaultProvider) Stream(ctx context.Context, messages []Message, tools []Tool, chunkChan chan<- string) error {
	if p == nil {
		return fmt.Errorf("provider is required")
	}
	if chunkChan == nil {
		return fmt.Errorf("chunkChan is required")
	}
	chunkChan <- "I am the "
	chunkChan <- "Native Go Borg Director."
	close(chunkChan)
	return nil
}

func (p *DefaultProvider) GetModelName() string {
	return "borg-native-stub-1.0"
}
