package agent

import (
	"errors"
	"strings"
	"testing"

	foundationorchestration "github.com/robertpelloni/hyperharness/foundation/orchestration"
)

func TestSortedAgentNamesReturnsDeterministicOrder(t *testing.T) {
	names := sortedAgentNames(map[string]*Agent{
		"zeta":  {},
		"alpha": {},
		"beta":  {},
	})
	joined := strings.Join(names, ",")
	if joined != "alpha,beta,zeta" {
		t.Fatalf("unexpected order: %q", joined)
	}
}

func TestBuildDelegatedExecutionTaskIncludesAllSteps(t *testing.T) {
	plan := foundationorchestration.PlanResult{Steps: []string{"step one", "step two"}}
	task := buildDelegatedExecutionTask(plan)
	if !strings.Contains(task, "Execute the following plan:") {
		t.Fatalf("expected execution prefix, got %q", task)
	}
	if !strings.Contains(task, "step one\nstep two") {
		t.Fatalf("expected joined steps, got %q", task)
	}
}

func TestRenderOrchestrationPlanIncludesRepoMapAndExecutionHeader(t *testing.T) {
	plan := foundationorchestration.PlanResult{
		TaskType:        "analysis",
		RepoMapIncluded: true,
		RepoMap:         "<repo_map>demo</repo_map>",
		Steps:           []string{"inspect", "summarize"},
	}
	plan.Execution.Route.Provider = "openai"
	plan.Execution.Route.Model = "gpt-4o"

	text := renderOrchestrationPlan(plan)
	for _, needle := range []string{"### Orchestration Plan ###", "Task Type: analysis", "Provider Route: openai/gpt-4o", "1. inspect", "2. summarize", "### Repo Map ###", "<repo_map>demo</repo_map>", "### Execution ###"} {
		if !strings.Contains(text, needle) {
			t.Fatalf("expected rendered plan to contain %q, got %q", needle, text)
		}
	}
}

func TestExecuteDelegatedPlanUsesSortedAgentsAndFormatsFailures(t *testing.T) {
	plan := foundationorchestration.PlanResult{Steps: []string{"inspect", "summarize"}}
	var calls []string
	delegate := func(name, task string) (string, error) {
		calls = append(calls, name+"::"+task)
		if name == "beta" {
			return "", errors.New("boom")
		}
		return "ok-" + name, nil
	}

	text := executeDelegatedPlan(plan, sortedAgentNames(map[string]*Agent{"beta": {}, "alpha": {}}), delegate)
	if len(calls) != 2 || !strings.HasPrefix(calls[0], "alpha::") || !strings.HasPrefix(calls[1], "beta::") {
		t.Fatalf("expected deterministic delegate order, got %#v", calls)
	}
	for _, needle := range []string{"Agent alpha result:\nok-alpha", "Agent beta failed: boom"} {
		if !strings.Contains(text, needle) {
			t.Fatalf("expected output to contain %q, got %q", needle, text)
		}
	}
}

func TestOrchestratorDelegateRejectsNilReceiver(t *testing.T) {
	var orch *Orchestrator
	_, err := orch.Delegate("executor", "task")
	if err == nil || !strings.Contains(err.Error(), "orchestrator is required") {
		t.Fatalf("expected nil orchestrator error, got %v", err)
	}
}

func TestOrchestratorPlanAndExecuteRejectsNilReceiver(t *testing.T) {
	var orch *Orchestrator
	_, err := orch.PlanAndExecute("task")
	if err == nil || !strings.Contains(err.Error(), "orchestrator is required") {
		t.Fatalf("expected nil orchestrator error, got %v", err)
	}
}
