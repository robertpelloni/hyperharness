package agents

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/robertpelloni/hyperharness/foundation/adapters"
	foundationorchestration "github.com/robertpelloni/hyperharness/foundation/orchestration"
	"github.com/robertpelloni/hyperharness/internal/council"
)

// Director Agent translates the TS core Director orchestrator.
// It acts as the primary task planner, coordinating sub-agents.
type Director struct {
	Name           string
	Provider       ILLMProvider
	State          map[string]interface{}
	History        []Message
	WorkingDir     string
	HyperAdapter   *adapters.HyperCodeAdapter
	Budget         *BudgetTracker
	ToolDispatcher *ToolDispatcher
}

func NewDirector(provider ILLMProvider) *Director {
	cwd, _ := os.Getwd()
	hyperAdapter := adapters.NewHyperCodeAdapter(cwd)
	td := NewToolDispatcher()
	td.RegisterDefaultStubs()
	return &Director{
		Name:           "Director",
		Provider:       provider,
		State:          make(map[string]interface{}),
		WorkingDir:     cwd,
		HyperAdapter:   hyperAdapter,
		Budget:         NewBudgetTracker(DefaultBudgetConfig()),
		ToolDispatcher: td,
		History: []Message{
			{
				Role:    RoleSystem,
				Content: strings.Join([]string{"You are the Borg TechLead Director. Your role is absolute architectural supervision. Plan, delegate, and review.", hyperAdapter.BuildSystemContext()}, "\n\n"),
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

	// Budget enforcement: block if session budget exceeded and enabled
	if d.Budget != nil && d.Budget.GetConfig().Enabled {
		status := d.Budget.BudgetStatus()
		if status.SessionCostPct >= 100 {
			return "", fmt.Errorf("session budget exceeded: $%.2f / $%.2f (%.0f%%)", status.SessionCost, status.SessionCostLimit, status.SessionCostPct)
		}
		if status.DailyCostPct >= 100 {
			return "", fmt.Errorf("daily budget exceeded: $%.2f / $%.2f (%.0f%%)", status.DailyCost, status.DailyCostLimit, status.DailyCostPct)
		}
		if status.SessionTokenPct >= 100 {
			return "", fmt.Errorf("session token budget exceeded: %d / %d (%.0f%%)", status.SessionTokens, status.SessionTokenLimit, status.SessionTokenPct)
		}
	}

	plan, err := foundationorchestration.BuildPlan(foundationorchestration.PlanRequest{
		Prompt:     input,
		WorkingDir: d.WorkingDir,
	})
	if err == nil {
		d.State["lastPlan"] = plan
	}

	providerMessages := append([]Message(nil), d.History...)
	if err == nil {
		providerMessages = append(providerMessages, Message{
			Role:    RoleSystem,
			Content: strings.Join([]string{"Execution planning context:", plan.SystemContextHint, strings.Join(plan.Steps, "\n")}, "\n"),
		})
	}

	responseMsg, providerErr := d.Provider.Chat(ctx, providerMessages, []Tool{})
	if providerErr != nil {
		return "", fmt.Errorf("director execution failed: %w", providerErr)
	}

	d.History = append(d.History, responseMsg)

	// Dispatch any tool calls from the LLM response
	for _, tc := range responseMsg.ToolCalls {
		args, parseErr := parseToolArgs(tc.Args)
		if parseErr != nil {
			d.History = append(d.History, Message{
				Role:    RoleAssistant,
				Content: fmt.Sprintf("[Tool %s] args parse error: %v", tc.Name, parseErr),
			})
			continue
		}
		result, execErr := d.ToolDispatcher.Dispatch(ctx, tc.Name, args)
		if execErr != nil {
			d.History = append(d.History, Message{
				Role:    RoleAssistant,
				Content: fmt.Sprintf("[Tool %s] execution error: %v", tc.Name, execErr),
			})
			// Trigger debate on failure to analyze the error and propose fixes
			director := council.NewDirectorAgent()
			topic := fmt.Sprintf("Tool '%s' failed with args %v. Error: %v. How should we fix this?", tc.Name, args, execErr)
			debateResult, debateErr := director.DebateOnFailure(topic)
			if debateErr == nil {
				d.History = append(d.History, Message{
					Role:    RoleAssistant,
					Content: fmt.Sprintf("[Debate on %s failure]\n%s", tc.Name, debateResult),
				})
			}
		} else {
			d.History = append(d.History, Message{
				Role:    RoleAssistant,
				Content: fmt.Sprintf("[Tool %s result]\n%s", tc.Name, result),
			})
		}
		// Nominal cost per tool invocation
		_ = d.Budget.RecordCost("tool", tc.Name, 0, 0, 0.001)
	}

	if err == nil {
		return fmt.Sprintf("[Director Plan]\n- task type: %s\n- route: %s/%s\n\n%s", plan.TaskType, plan.Execution.Route.Provider, plan.Execution.Route.Model, responseMsg.Content), nil
	}
	return responseMsg.Content, nil
}

func (d *Director) InjectSystemContext(context string) {
	d.History[0].Content += "\n\n" + context
}

func (d *Director) GetState() map[string]interface{} {
	return d.State
}

// Example Stubs for other agents to achieve parity:

type Coder struct{ Director } // Inherits base logic for simplicity in this stub
type MetaArchitect struct{ Director }
type Researcher struct{ Director }
type Council struct{ Director }
