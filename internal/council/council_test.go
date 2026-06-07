package council

import (
	"context"
	"strings"
	"testing"
)

type MockExecutor struct {
	Responses []string
	Index     int
}

func (m *MockExecutor) Chat(input string) (string, error) {
	if m.Index >= len(m.Responses) {
		return "Mock completed", nil
	}
	res := m.Responses[m.Index]
	m.Index++
	return res, nil
}

func TestDirectorExecutePlan(t *testing.T) {
	d := NewDirectorAgent()
	exec := &MockExecutor{
		Responses: []string{
			`[{"id": "t1", "role": "coder", "objective": "Write code"}]`, // PlanDelegation response
			"Code written successfully",                                   // Worker execution response
		},
	}

	plan, err := d.PlanDelegation(context.Background(), exec, "Build a feature")
	if err != nil {
		t.Fatalf("failed to plan: %v", err)
	}

	if len(plan) != 1 {
		t.Fatalf("expected 1 task in plan, got %d", len(plan))
	}

	result, err := d.ExecutePlan()
	if err != nil {
		t.Fatalf("expected execution to succeed, got error: %v", err)
	}

	if !strings.Contains(result, "Code written successfully") {
		t.Errorf("missing expected output in execution results: %s", result)
	}
}

func TestInitiateDebate(t *testing.T) {
	d := NewDirectorAgent()
	exec := &MockExecutor{}
	a := NewWorkerAgent("researcher")
	a.exec = exec
	b := NewWorkerAgent("reviewer")
	b.exec = exec

	result, err := d.InitiateDebate("use a map vs a slice", a, b)
	if err != nil {
		t.Fatalf("expected debate to succeed, got error: %v", err)
	}

	if !strings.Contains(result, "Consensus reached by Director") {
		t.Errorf("missing consensus marker")
	}
}
