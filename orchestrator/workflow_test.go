package orchestrator

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkflowCreation(t *testing.T) {
	steps := []*Step{
		{ID: "s1", Name: "Step 1"},
		{ID: "s2", Name: "Step 2", DependsOn: []string{"s1"}},
	}
	wf := NewWorkflow("test-wf", "Test", "Test workflow", steps)
	if wf.Status != StepPending {
		t.Error("workflow should start as pending")
	}
	if len(wf.Steps) != 2 {
		t.Error("should have 2 steps")
	}
	for _, s := range wf.Steps {
		if s.Status != StepPending {
			t.Errorf("step %s should be pending", s.ID)
		}
	}
}

func TestWorkflowLinearExecution(t *testing.T) {
	var order []string
	steps := []*Step{
		{
			ID: "s1",
			Execute: func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
				order = append(order, "s1")
				return map[string]interface{}{"result": "hello"}, nil
			},
		},
		{
			ID:        "s2",
			DependsOn: []string{"s1"},
			Execute: func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
				order = append(order, "s2")
				return map[string]interface{}{"result": inputs["s1.result"]}, nil
			},
		},
	}

	wf := NewWorkflow("test", "Test", "", steps)
	err := wf.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if wf.Status != StepCompleted {
		t.Errorf("expected completed, got %s", wf.Status)
	}
	if len(order) != 2 || order[0] != "s1" || order[1] != "s2" {
		t.Errorf("execution order: %v", order)
	}
	if wf.Duration() < 0 {
		t.Error("duration should be >= 0")
	}
}

func TestWorkflowParallelExecution(t *testing.T) {
	var started int64
	steps := []*Step{
		{
			ID: "a",
			Execute: func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
				atomic.AddInt64(&started, 1)
				time.Sleep(50 * time.Millisecond)
				return map[string]interface{}{"a": 1}, nil
			},
		},
		{
			ID: "b",
			Execute: func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
				atomic.AddInt64(&started, 1)
				time.Sleep(50 * time.Millisecond)
				return map[string]interface{}{"b": 2}, nil
			},
		},
		{
			ID:        "c",
			DependsOn: []string{"a", "b"},
			Execute: func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{"combined": inputs}, nil
			},
		},
	}

	wf := NewWorkflow("parallel", "Parallel", "", steps)
	err := wf.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// a and b should have run concurrently, so total time < 150ms
	if wf.Duration() > 140*time.Millisecond {
		t.Errorf("steps should have run in parallel, took %v", wf.Duration())
	}
}

func TestWorkflowFailureSkipsDependents(t *testing.T) {
	steps := []*Step{
		{
			ID: "fail",
			Execute: func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
				return nil, fmt.Errorf("intentional failure")
			},
		},
		{
			ID:        "after-fail",
			DependsOn: []string{"fail"},
			Execute: func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
				return nil, nil
			},
		},
	}

	wf := NewWorkflow("fail-test", "Fail Test", "", steps)
	err := wf.Run(context.Background())
	if err == nil {
		t.Error("should have error from failed step")
	}

	if wf.Status != StepFailed {
		t.Errorf("workflow should be failed, got %s", wf.Status)
	}

	// Check the dependent step was skipped
	step, ok := wf.GetStep("after-fail")
	if !ok {
		t.Fatal("step not found")
	}
	if step.Status != StepSkipped {
		t.Errorf("dependent should be skipped, got %s", step.Status)
	}

	if wf.FailedCount() != 1 {
		t.Errorf("expected 1 failed, got %d", wf.FailedCount())
	}
}

func TestWorkflowCircularDependency(t *testing.T) {
	steps := []*Step{
		{ID: "a", DependsOn: []string{"b"}},
		{ID: "b", DependsOn: []string{"a"}},
	}

	wf := NewWorkflow("circular", "Circular", "", steps)
	err := wf.Run(context.Background())
	if err == nil {
		t.Error("should fail with circular dependency")
	}
}

func TestWorkflowCancellation(t *testing.T) {
	steps := []*Step{
		{
			ID: "slow",
			Execute: func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
				<-ctx.Done()
				return nil, ctx.Err()
			},
		},
		{
			ID:        "after-slow",
			DependsOn: []string{"slow"},
			Execute: func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
				return nil, nil
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	wf := NewWorkflow("cancel", "Cancel Test", "", steps)
	wf.Run(ctx)

	step, _ := wf.GetStep("after-slow")
	if step.Status == StepCompleted {
		t.Error("dependent step should not complete after cancellation")
	}
}

func TestWorkflowGetStep(t *testing.T) {
	steps := []*Step{
		{ID: "s1"},
		{ID: "s2"},
	}
	wf := NewWorkflow("test", "Test", "", steps)

	step, ok := wf.GetStep("s1")
	if !ok || step.ID != "s1" {
		t.Error("should find s1")
	}

	_, ok = wf.GetStep("nonexistent")
	if ok {
		t.Error("should not find nonexistent step")
	}
}

func TestWorkflowCounts(t *testing.T) {
	steps := []*Step{
		{ID: "s1", Execute: func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) { return nil, nil }},
		{ID: "s2", Execute: func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) { return nil, nil }},
	}
	wf := NewWorkflow("test", "Test", "", steps)
	wf.Run(context.Background())

	if wf.CompletedCount() != 2 {
		t.Errorf("expected 2 completed, got %d", wf.CompletedCount())
	}
	if wf.FailedCount() != 0 {
		t.Errorf("expected 0 failed, got %d", wf.FailedCount())
	}
}

func TestWorkflowEngine(t *testing.T) {
	engine := NewWorkflowEngine()

	wf := NewWorkflow("test", "Test", "", []*Step{
		{ID: "s1", Execute: func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{"value": 42}, nil
		}},
	})

	engine.Register(wf)

	// Get
	retrieved, ok := engine.Get("test")
	if !ok || retrieved.ID != "test" {
		t.Error("should retrieve registered workflow")
	}

	// List
	list := engine.List()
	if len(list) != 1 {
		t.Errorf("expected 1 workflow, got %d", len(list))
	}

	// Run
	err := engine.RunWorkflow(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}

	// Delete
	engine.Delete("test")
	_, ok = engine.Get("test")
	if ok {
		t.Error("should not find deleted workflow")
	}

	// Run nonexistent
	err = engine.RunWorkflow(context.Background(), "nonexistent")
	if err == nil {
		t.Error("should error for nonexistent workflow")
	}
}

func TestWorkflowNoExecutor(t *testing.T) {
	steps := []*Step{
		{ID: "noop"}, // No Execute function
	}
	wf := NewWorkflow("noop", "No-op", "", steps)
	err := wf.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if wf.Status != StepCompleted {
		t.Errorf("expected completed, got %s", wf.Status)
	}
}

// Ensure fmt is used
var _ = fmt.Sprintf
