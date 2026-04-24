package subagents

import (
	"context"
	"testing"
)

func TestManagerCreateTask(t *testing.T) {
	m := NewManager()
	task := m.CreateTask(TypeCode, "test task", "do something", "parent_1")

	if task.ID == "" {
		t.Error("task ID should not be empty")
	}
	if task.Type != TypeCode {
		t.Errorf("task type: got %s, want %s", task.Type, TypeCode)
	}
	if task.State != StatePending {
		t.Errorf("task state: got %s, want %s", task.State, StatePending)
	}
	if task.ParentID != "parent_1" {
		t.Errorf("parent ID: got %s", task.ParentID)
	}
}

func TestManagerExecuteTask(t *testing.T) {
	m := NewManager()
	task := m.CreateTask(TypeCode, "test task", "write hello world", "")

	result, err := m.ExecuteTask(context.Background(), task)
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if task.State != StateCompleted {
		t.Errorf("task state after execute: got %s, want %s", task.State, StateCompleted)
	}
	if result == "" {
		t.Error("result should not be empty")
	}
}

func TestManagerGetTask(t *testing.T) {
	m := NewManager()
	task := m.CreateTask(TypeResearch, "research task", "find docs", "")

	found, ok := m.GetTask(task.ID)
	if !ok {
		t.Error("task should be found")
	}
	if found.ID != task.ID {
		t.Errorf("task ID mismatch: got %s, want %s", found.ID, task.ID)
	}

	_, ok = m.GetTask("nonexistent")
	if ok {
		t.Error("nonexistent task should not be found")
	}
}

func TestManagerListTasks(t *testing.T) {
	m := NewManager()
	m.CreateTask(TypeCode, "task1", "prompt1", "")
	m.CreateTask(TypeResearch, "task2", "prompt2", "")
	m.CreateTask(TypeCode, "task3", "prompt3", "")

	all := m.ListTasks("")
	if len(all) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(all))
	}
}

func TestManagerCancelTask(t *testing.T) {
	m := NewManager()
	task := m.CreateTask(TypeCode, "cancelable", "prompt", "")

	err := m.CancelTask(task.ID)
	if err != nil {
		t.Fatalf("cancel failed: %v", err)
	}

	found, _ := m.GetTask(task.ID)
	if found.State != StateCancelled {
		t.Errorf("task state after cancel: got %s, want %s", found.State, StateCancelled)
	}
}

func TestDefaultConfigs(t *testing.T) {
	configs := DefaultSubagentConfigs()

	expectedTypes := []SubagentType{
		TypeCode, TypeResearch, TypeReview, TypePlan, TypeBuild,
		TypeTest, TypeDebug, TypeDoc, TypeSecurity, TypeDevOps,
	}

	for _, st := range expectedTypes {
		config, ok := configs[st]
		if !ok {
			t.Errorf("missing config for type: %s", st)
			continue
		}
		if config.Name == "" {
			t.Errorf("config %s missing name", st)
		}
		if config.MaxTokens <= 0 {
			t.Errorf("config %s has invalid MaxTokens: %d", st, config.MaxTokens)
		}
		if config.MaxTurns <= 0 {
			t.Errorf("config %s has invalid MaxTurns: %d", st, config.MaxTurns)
		}
		if len(config.AllowedTools) == 0 {
			t.Errorf("config %s has no allowed tools", st)
		}
	}
}

func TestManagerGetConfig(t *testing.T) {
	m := NewManager()

	config, ok := m.GetConfig(TypeCode)
	if !ok {
		t.Fatal("code config should exist")
	}
	if config.Type != TypeCode {
		t.Errorf("config type mismatch: %s", config.Type)
	}

	_, ok = m.GetConfig(SubagentType("nonexistent"))
	if ok {
		t.Error("nonexistent config should not be found")
	}
}

func TestManagerListConfigs(t *testing.T) {
	m := NewManager()
	configs := m.ListConfigs()
	if len(configs) < 10 {
		t.Errorf("expected at least 10 configs, got %d", len(configs))
	}
}
