package agent

import (
	"fmt"
	"os"
	"sort"
	"strings"

	foundationorchestration "github.com/robertpelloni/hyperharness/foundation/orchestration"
)

const (
	defaultExecutorAgentName    = "executor"
	defaultExecutorSystemPrompt = "You are the Executor Agent. Complete the tasks assigned to you."
)

// Orchestrator manages multiple sub-agents (Maestro/Codemachine parity)
type Orchestrator struct {
	Agents     map[string]*Agent
	Lead       *Agent
	WorkingDir string
}

func NewOrchestrator() *Orchestrator {
	lead := NewAgent()
	lead.messages[0].Content = "You are the TechLead Orchestrator. Break down user requests and delegate to specialized sub-agents."
	cwd, _ := os.Getwd()

	return &Orchestrator{
		Agents:     make(map[string]*Agent),
		Lead:       lead,
		WorkingDir: cwd,
	}
}

// Spawn creates a new specialized sub-agent
func (o *Orchestrator) Spawn(name, systemPrompt string) *Agent {
	if o == nil {
		return nil
	}
	if o.Agents == nil {
		o.Agents = make(map[string]*Agent)
	}
	subAgent := NewAgent()
	subAgent.messages[0].Content = systemPrompt
	o.Agents[name] = subAgent
	return subAgent
}

// Delegate passes a task to a specific sub-agent
func (o *Orchestrator) Delegate(agentName, task string) (string, error) {
	if o == nil {
		return "", fmt.Errorf("orchestrator is required")
	}
	agent, exists := o.Agents[agentName]
	if !exists {
		return "", fmt.Errorf("sub-agent %s not found", agentName)
	}
	return agent.Chat(fmt.Sprintf("[Delegated Task]: %s", task))
}

func (o *Orchestrator) BuildPlan(task string) (foundationorchestration.PlanResult, error) {
	return foundationorchestration.BuildPlan(foundationorchestration.PlanRequest{
		Prompt:       task,
		WorkingDir:   o.WorkingDir,
		IncludeRepo:  true,
		MaxRepoFiles: 8,
	})
}

// PlanAndExecute uses the local orchestration planner and then delegates it.
func (o *Orchestrator) PlanAndExecute(task string) (string, error) {
	if o == nil {
		return "", fmt.Errorf("orchestrator is required")
	}
	plan, err := o.BuildPlan(task)
	if err != nil {
		return "", err
	}
	if len(o.Agents) == 0 {
		o.Spawn(defaultExecutorAgentName, defaultExecutorSystemPrompt)
	}

	var builder strings.Builder
	builder.WriteString(renderOrchestrationPlan(plan))
	builder.WriteString(executeDelegatedPlan(plan, sortedAgentNames(o.Agents), o.Delegate))
	return builder.String(), nil
}

func renderOrchestrationPlan(plan foundationorchestration.PlanResult) string {
	var builder strings.Builder
	builder.WriteString("### Orchestration Plan ###\n")
	builder.WriteString(fmt.Sprintf("Task Type: %s\n", plan.TaskType))
	builder.WriteString(fmt.Sprintf("Provider Route: %s/%s\n", plan.Execution.Route.Provider, plan.Execution.Route.Model))
	for i, step := range plan.Steps {
		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, step))
	}
	if plan.RepoMapIncluded {
		builder.WriteString("\n### Repo Map ###\n")
		builder.WriteString(plan.RepoMap)
		builder.WriteString("\n")
	}
	builder.WriteString("\n### Execution ###\n")
	return builder.String()
}

func executeDelegatedPlan(plan foundationorchestration.PlanResult, agentNames []string, delegate func(string, string) (string, error)) string {
	var builder strings.Builder
	task := buildDelegatedExecutionTask(plan)
	for _, name := range agentNames {
		result, err := delegate(name, task)
		if err != nil {
			builder.WriteString(fmt.Sprintf("Agent %s failed: %v\n", name, err))
		} else {
			builder.WriteString(fmt.Sprintf("Agent %s result:\n%s\n", name, result))
		}
	}
	return builder.String()
}

func buildDelegatedExecutionTask(plan foundationorchestration.PlanResult) string {
	return fmt.Sprintf("Execute the following plan:\n%s", strings.Join(plan.Steps, "\n"))
}

func sortedAgentNames(agents map[string]*Agent) []string {
	names := make([]string, 0, len(agents))
	for name := range agents {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
