package workflow

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestLinearWorkflow(t *testing.T) {
	wf := NewWorkflow("test", "linear", "test linear", []*Step{
		{ID: "a", Name: "A", Execute: func(ctx context.Context, _ map[string]any) (map[string]any, error) {
			return map[string]any{"val": 1}, nil
		}},
		{ID: "b", Name: "B", DependsOn: []string{"a"}, Execute: func(ctx context.Context, inputs map[string]any) (map[string]any, error) {
			return map[string]any{"val": inputs["a.val"]}, nil
		}},
	})
	if err := wf.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if wf.Status != StatusCompleted {
		t.Errorf("status: %s", wf.Status)
	}
}

func TestCircularDependency(t *testing.T) {
	wf := NewWorkflow("circ", "circular", "", []*Step{
		{ID: "a", DependsOn: []string{"b"}},
		{ID: "b", DependsOn: []string{"a"}},
	})
	if err := wf.Run(context.Background()); err == nil {
		t.Error("should fail with circular dep")
	}
}

func TestStepFailureSkipsDependents(t *testing.T) {
	wf := NewWorkflow("fail", "failing", "", []*Step{
		{ID: "a", Execute: func(ctx context.Context, _ map[string]any) (map[string]any, error) {
			return nil, errors.New("boom")
		}},
		{ID: "b", DependsOn: []string{"a"}, Execute: func(ctx context.Context, _ map[string]any) (map[string]any, error) {
			return map[string]any{"b": true}, nil
		}},
	})
	err := wf.Run(context.Background())
	if err == nil {
		t.Error("should fail")
	}
	if wf.Steps[0].Status != StatusFailed {
		t.Errorf("a should be failed: %s", wf.Steps[0].Status)
	}
	if wf.Steps[1].Status != StatusSkipped {
		t.Errorf("b should be skipped: %s", wf.Steps[1].Status)
	}
}

func TestParallelWorkflow(t *testing.T) {
	wf := NewWorkflow("par", "parallel", "", []*Step{
		{ID: "a", Execute: func(ctx context.Context, _ map[string]any) (map[string]any, error) {
			time.Sleep(10 * time.Millisecond)
			return map[string]any{"a": true}, nil
		}},
		{ID: "b", Execute: func(ctx context.Context, _ map[string]any) (map[string]any, error) {
			time.Sleep(10 * time.Millisecond)
			return map[string]any{"b": true}, nil
		}},
		{ID: "c", DependsOn: []string{"a", "b"}, Execute: func(ctx context.Context, inputs map[string]any) (map[string]any, error) {
			return map[string]any{"merged": len(inputs)}, nil
		}},
	})
	start := time.Now()
	if err := wf.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	elapsed := time.Since(start)
	if elapsed > 50*time.Millisecond {
		t.Errorf("should be parallel, took %v", elapsed)
	}
}
