package council

import (
	"strings"
	"testing"
)

func TestDirectorExecutePlan(t *testing.T) {
	d := NewDirectorAgent()
	plan := d.PlanDelegation("Build a feature")

	if len(plan) != 3 {
		t.Fatalf("expected 3 tasks in plan, got %d", len(plan))
	}

	result, err := d.ExecutePlan()
	if err != nil {
		t.Fatalf("expected execution to succeed, got error: %v", err)
	}

	if !strings.Contains(result, "Completed objective: Find relevant files for the goal") {
		t.Errorf("missing expected output in execution results")
	}
}

func TestInitiateDebate(t *testing.T) {
	d := NewDirectorAgent()
	a := NewWorkerAgent("researcher")
	b := NewWorkerAgent("reviewer")

	result, err := d.InitiateDebate("use a map vs a slice", a, b)
	if err != nil {
		t.Fatalf("expected debate to succeed, got error: %v", err)
	}

	if !strings.Contains(result, "Consensus reached by Director") {
		t.Errorf("missing consensus marker")
	}
}
