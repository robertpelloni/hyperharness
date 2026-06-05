package subagents

import (
	"context"
	"testing"
	"time"
)

func TestManagerTaskLifecycle(t *testing.T) {
	mgr := NewManager()

	t.Run("CreateTask", func(t *testing.T) {
		task := mgr.CreateTask(TypeCode, "Test Prompt", "Test Input", "Test Context")
		if task.ID == "" {
			t.Fatal("Task ID should not be empty")
		}
		if task.Type != TypeCode {
			t.Errorf("Expected task type %s, got %s", TypeCode, task.Type)
		}
		if task.Status != "created" {
			t.Errorf("Expected status 'created', got %s", task.Status)
		}
	})

	t.Run("ExecuteTask", func(t *testing.T) {
		task := mgr.CreateTask(TypePlan, "Generate Plan", "", "")

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		output, err := mgr.ExecuteTask(ctx, task)
		if err != nil {
			t.Fatalf("ExecuteTask failed: %v", err)
		}

		if task.Status != "completed" {
			t.Errorf("Expected status 'completed', got %s", task.Status)
		}
		if output == "" {
			t.Error("Expected non-empty output")
		}
	})

	t.Run("ExecuteTaskTimeout", func(t *testing.T) {
		// This test expects a timeout because we mocked 500ms delay and
		// we will provide a very short context timeout.
		task := mgr.CreateTask(TypeCode, "Slow Task", "", "")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err := mgr.ExecuteTask(ctx, task)
		if err == nil {
			t.Error("Expected timeout error")
		}
	})
}
