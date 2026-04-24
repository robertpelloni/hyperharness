package tools

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSessionTodoStoreCreate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")
	store := NewSessionTodoStore(path)

	if store == nil {
		t.Fatal("store should not be nil")
	}
	if len(store.List()) != 0 {
		t.Error("new store should have no todos")
	}
}

func TestSessionTodoStoreUpdate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")
	store := NewSessionTodoStore(path)

	todos := []SessionTodo{
		{Content: "Task 1", ActiveForm: "Doing task 1", Status: "pending"},
		{Content: "Task 2", ActiveForm: "Doing task 2", Status: "in_progress"},
		{Content: "Task 3", ActiveForm: "Doing task 3", Status: "completed"},
	}

	old, new_ := store.Update(todos)

	if len(old) != 0 {
		t.Errorf("old should be empty, got %d", len(old))
	}
	if len(new_) != 3 {
		t.Fatalf("new should have 3 items, got %d", len(new_))
	}
	if new_[0].ID == "" {
		t.Error("todos should get IDs")
	}
	if new_[1].Status != "in_progress" {
		t.Errorf("status: %s", new_[1].Status)
	}
}

func TestSessionTodoStorePersistence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")

	// Create and populate
	store1 := NewSessionTodoStore(path)
	store1.Update([]SessionTodo{
		{Content: "Persistent task", ActiveForm: "Persisting", Status: "pending"},
	})

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("todo file should be persisted")
	}

	// Load in new store
	store2 := NewSessionTodoStore(path)
	todos := store2.List()
	if len(todos) != 1 {
		t.Fatalf("expected 1 persisted todo, got %d", len(todos))
	}
	if todos[0].Content != "Persistent task" {
		t.Errorf("content: %s", todos[0].Content)
	}
}

func TestSessionTodoStoreGetInProgress(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")
	store := NewSessionTodoStore(path)

	store.Update([]SessionTodo{
		{Content: "Pending", ActiveForm: "Pending", Status: "pending"},
		{Content: "Active", ActiveForm: "Activating", Status: "in_progress"},
		{Content: "Done", ActiveForm: "Done", Status: "completed"},
	})

	active := store.GetInProgress()
	if active == nil {
		t.Fatal("should have an in-progress todo")
	}
	if active.Content != "Active" {
		t.Errorf("active content: %s", active.Content)
	}
}

func TestSessionTodoStoreCountByStatus(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")
	store := NewSessionTodoStore(path)

	store.Update([]SessionTodo{
		{Content: "P1", ActiveForm: "P1", Status: "pending"},
		{Content: "P2", ActiveForm: "P2", Status: "pending"},
		{Content: "A1", ActiveForm: "A1", Status: "in_progress"},
		{Content: "D1", ActiveForm: "D1", Status: "completed"},
		{Content: "D2", ActiveForm: "D2", Status: "completed"},
		{Content: "D3", ActiveForm: "D3", Status: "completed"},
	})

	counts := store.CountByStatus()
	if counts["pending"] != 2 {
		t.Errorf("pending: %d", counts["pending"])
	}
	if counts["in_progress"] != 1 {
		t.Errorf("in_progress: %d", counts["in_progress"])
	}
	if counts["completed"] != 3 {
		t.Errorf("completed: %d", counts["completed"])
	}
}

func TestSessionTodoStoreSummary(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")
	store := NewSessionTodoStore(path)

	// Empty store
	s := store.Summary()
	if s != "No tasks tracked." {
		t.Errorf("empty summary: %s", s)
	}

	// With tasks
	store.Update([]SessionTodo{
		{Content: "Active", ActiveForm: "Working on Active", Status: "in_progress"},
	})

	s = store.Summary()
	if s == "" {
		t.Error("summary should not be empty")
	}
}

func TestSessionTodoStoreUpdateReplaces(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")
	store := NewSessionTodoStore(path)

	// First update
	store.Update([]SessionTodo{
		{Content: "Old task", ActiveForm: "Old", Status: "pending"},
	})

	// Second update replaces entirely
	old, new_ := store.Update([]SessionTodo{
		{Content: "New task 1", ActiveForm: "New 1", Status: "pending"},
		{Content: "New task 2", ActiveForm: "New 2", Status: "pending"},
	})

	if len(old) != 1 || old[0].Content != "Old task" {
		t.Errorf("old should contain first update: %v", old)
	}
	if len(new_) != 2 {
		t.Errorf("new should have 2 items: %d", len(new_))
	}

	// Store should now have 2 items
	if len(store.List()) != 2 {
		t.Errorf("store should have 2 items: %d", len(store.List()))
	}
}

func TestSessionTodoStoreEmptyPath(t *testing.T) {
	store := NewSessionTodoStore("")
	store.Update([]SessionTodo{
		{Content: "No path", ActiveForm: "No path", Status: "pending"},
	})
	// Should not panic, just won't persist
	_ = time.Now // suppress unused import
}
