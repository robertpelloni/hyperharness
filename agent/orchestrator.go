package agent

import (
	"fmt"
	"strings"
)

// Orchestrator manages multiple sub-agents (Maestro/Codemachine parity)
type Orchestrator struct {
	Agents map[string]*Agent
	Lead   *Agent
}

func NewOrchestrator() *Orchestrator {
	lead := NewAgent()
	lead.messages[0].Content = "You are the TechLead Orchestrator. Break down user requests and delegate to specialized sub-agents."

	return &Orchestrator{
		Agents: make(map[string]*Agent),
		Lead:   lead,
	}
}

// Spawn creates a new specialized sub-agent
func (o *Orchestrator) Spawn(name, systemPrompt string) *Agent {
	subAgent := NewAgent()
	subAgent.messages[0].Content = systemPrompt
	o.Agents[name] = subAgent
	return subAgent
}

// Delegate passes a task to a specific sub-agent
func (o *Orchestrator) Delegate(agentName, task string) (string, error) {
	agent, exists := o.Agents[agentName]
	if !exists {
		return "", fmt.Errorf("sub-agent %s not found", agentName)
	}
	return agent.Chat(fmt.Sprintf("[Delegated Task]: %s", task))
}

// PlanAndExecute uses the Lead agent to plan a task and then delegates it
func (o *Orchestrator) PlanAndExecute(task string) (string, error) {
	// 1. Lead creates a plan
	planResponse, err := o.Lead.Chat(fmt.Sprintf("Create a step-by-step plan for: %s. Output only the plan.", task))
	if err != nil {
		return "", err
	}

	// 2. Here we would parse the plan and delegate. For now, we simulate execution.
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("### Orchestration Plan ###\n%s\n\n### Execution ###\n", planResponse))

	// Create a generic execution agent if none exist
	if len(o.Agents) == 0 {
		o.Spawn("executor", "You are the Executor Agent. Complete the tasks assigned to you.")
	}

	for name := range o.Agents {
		result, err := o.Delegate(name, fmt.Sprintf("Execute the following plan:\n%s", planResponse))
		if err != nil {
			builder.WriteString(fmt.Sprintf("Agent %s failed: %v\n", name, err))
		} else {
			builder.WriteString(fmt.Sprintf("Agent %s result:\n%s\n", name, result))
		}
	}

	return builder.String(), nil
}
