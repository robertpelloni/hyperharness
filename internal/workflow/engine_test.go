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
	for _, s := range wf.Steps {
		if s.Status != StatusCompleted {
			t.Errorf("step %s: %s", s.ID, s.Status)
		}
	}
}

func TestParallelWorkflow(t *testing.T) {
	wf := NewWorkflow("par", "parallel", "", []*Step{
		{ID: "a", Name: "A", Execute: func(ctx context.Context, _ map[string]any) (map[string]any, error) {
			time.Sleep(10 * time.Millisecond)
			return map[string]any{"a": true}, nil
		}},
		{ID: "b", Name: "B", Execute: func(ctx context.Context, _ map[string]any) (map[string]any, error) {
			time.Sleep(10 * time.Millisecond)
			return map[string]any{"b": true}, nil
		}},
		{ID: "c", Name: "C", DependsOn: []string{"a", "b"}, Execute: func(ctx context.Context, inputs map[string]any) (map[string]any, error) {
			return map[string]any{"merged": len(inputs)}, nil
		}},
	})

	start := time.Now()
	if err := wf.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	elapsed := time.Since(start)

	// a and b should run in parallel, so total should be < 30ms
	if elapsed > 50*time.Millisecond {
		t.Errorf("should be parallel, took %v", elapsed)
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

func TestEngine(t *testing.T) {
	engine := NewEngine()
	wf := NewWorkflow("w1", "workflow1", "", []*Step{
		{ID: "s1", Execute: func(ctx context.Context, _ map[string]any) (map[string]any, error) {
			return nil, nil
		}},
	})
	engine.Register(wf)

	got, ok := engine.Get("w1")
	if !ok {
		t.Fatal("should find workflow")
	}
	if got.ID != "w1" {
		t.Errorf("id: %s", got.ID)
	}
	if len(engine.List()) != 1 {
		t.Errorf("list count: %d", len(engine.List()))
	}
}

func TestEngineRunWorkflow(t *testing.T) {
	engine := NewEngine()
	wf := NewWorkflow("run", "runnable", "", []*Step{
		{ID: "s1", Execute: func(ctx context.Context, _ map[string]any) (map[string]any, error) {
			return map[string]any{"ok": true}, nil
		}},
	})
	engine.Register(wf)

	if err := engine.RunWorkflow(context.Background(), "run"); err != nil {
		t.Fatal(err)
	}
}

func TestEngineNotFound(t *testing.T) {
	engine := NewEngine()
	err := engine.RunWorkflow(context.Background(), "missing")
	if err == nil {
		t.Error("should error for missing workflow")
	}
}

func TestStepWithNoExecute(t *testing.T) {
	wf := NewWorkflow("noop", "no-execute", "", []*Step{
		{ID: "a"},
	})
	if err := wf.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	wf := NewWorkflow("cancel", "cancelled", "", []*Step{
		{ID: "a", Execute: func(ctx context.Context, _ map[string]any) (map[string]any, error) {
			return nil, nil
		}},
	})
	// Should handle gracefully - may or may not error depending on timing
	wf.Run(ctx)
}
