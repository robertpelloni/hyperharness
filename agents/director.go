package agents

import (
	"context"
	"fmt"
)

// Director Agent translates the TS core Director orchestrator.
// It acts as the primary task planner, coordinating sub-agents.
type Director struct {
	Name     string
	Provider ILLMProvider
	State    map[string]interface{}
	History  []Message
}

func NewDirector(provider ILLMProvider) *Director {
	return &Director{
		Name:     "Director",
		Provider: provider,
		State:    make(map[string]interface{}),
		History: []Message{
			{
				Role:    RoleSystem,
				Content: "You are the Borg TechLead Director. Your role is absolute architectural supervision. Plan, delegate, and review.",
			},
		},
	}
}

func (d *Director) GetName() string {
	return d.Name
}

func (d *Director) GetRole() string {
	return "supervisor"
}

func (d *Director) HandleInput(ctx context.Context, input string) (string, error) {
	d.History = append(d.History, Message{Role: RoleUser, Content: input})
	
	// Delegate to the provider
	responseMsg, err := d.Provider.Chat(ctx, d.History, []Tool{})
	if err != nil {
		return "", fmt.Errorf("director execution failed: %w", err)
	}

	d.History = append(d.History, responseMsg)
	return responseMsg.Content, nil
}

func (d *Director) InjectSystemContext(context string) {
	d.History[0].Content += "\n\n" + context
}

func (d *Director) GetState() map[string]interface{} {
	return d.State
}

// Example Stubs for other agents to achieve parity:

type Coder struct { Director } // Inherits base logic for simplicity in this stub
type MetaArchitect struct { Director }
type Researcher struct { Director }
type Council struct { Director }
