// Package tools provides session-level todo/task tracking.
// Implements the state backend for the TodoWrite tool from Claude Code,
// supporting pending/in_progress/completed states with persistence.
package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SessionTodo represents a tracked task in the session.
type SessionTodo struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`    // Imperative form: "Fix auth bug"
	ActiveForm string    `json:"activeForm"` // Present continuous: "Fixing auth bug"
	Status     string    `json:"status"`     // "pending" | "in_progress" | "completed"
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// SessionTodoStore manages session-level todos with persistence.
type SessionTodoStore struct {
	todos  []SessionTodo
	path   string
	nextID int
	mu     sync.Mutex
}

// NewSessionTodoStore creates a new todo store, loading from the given path if it exists.
func NewSessionTodoStore(storePath string) *SessionTodoStore {
	ts := &SessionTodoStore{
		todos: make([]SessionTodo, 0),
		path:  storePath,
	}
	ts.load()
	return ts
}

// Update replaces all todos with the given list.
// Returns old and new todo lists for the TodoWrite tool response.
func (ts *SessionTodoStore) Update(items []SessionTodo) (old, new_ []SessionTodo) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	old = make([]SessionTodo, len(ts.todos))
	copy(old, ts.todos)

	now := time.Now()
	new_ = make([]SessionTodo, len(items))
	for i := range items {
		if items[i].ID == "" {
			ts.nextID++
			items[i].ID = fmt.Sprintf("todo_%d", ts.nextID)
			items[i].CreatedAt = now
		}
		items[i].UpdatedAt = now
		new_[i] = items[i]
	}

	ts.todos = new_
	ts.save()

	return old, new_
}

// List returns all current todos.
func (ts *SessionTodoStore) List() []SessionTodo {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	result := make([]SessionTodo, len(ts.todos))
	copy(result, ts.todos)
	return result
}

// GetInProgress returns the currently in-progress todo (should be at most 1).
func (ts *SessionTodoStore) GetInProgress() *SessionTodo {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	for i := range ts.todos {
		if ts.todos[i].Status == "in_progress" {
			return &ts.todos[i]
		}
	}
	return nil
}

// CountByStatus returns the count of todos in each status.
func (ts *SessionTodoStore) CountByStatus() map[string]int {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	counts := map[string]int{"pending": 0, "in_progress": 0, "completed": 0}
	for _, t := range ts.todos {
		counts[t.Status]++
	}
	return counts
}

// Summary returns a human-readable summary.
func (ts *SessionTodoStore) Summary() string {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if len(ts.todos) == 0 {
		return "No tasks tracked."
	}

	counts := make(map[string]int)
	for _, t := range ts.todos {
		counts[t.Status]++
	}

	var inProgress string
	for _, t := range ts.todos {
		if t.Status == "in_progress" {
			inProgress = fmt.Sprintf("\n  ▶ IN PROGRESS: %s", t.ActiveForm)
		}
	}

	return fmt.Sprintf("Tasks: %d pending, %d in_progress, %d completed%s",
		counts["pending"], counts["in_progress"], counts["completed"], inProgress)
}

// load reads todos from disk.
func (ts *SessionTodoStore) load() {
	if ts.path == "" {
		return
	}
	data, err := os.ReadFile(ts.path)
	if err != nil {
		return
	}
	var todos []SessionTodo
	if err := json.Unmarshal(data, &todos); err != nil {
		return
	}
	ts.todos = todos
	// Compute next ID
	for _, t := range todos {
		if t.ID != "" {
			ts.nextID++
		}
	}
}

// save writes todos to disk.
func (ts *SessionTodoStore) save() {
	if ts.path == "" {
		return
	}
	data, _ := json.MarshalIndent(ts.todos, "", "  ")
	os.MkdirAll(filepath.Dir(ts.path), 0o755)
	os.WriteFile(ts.path, data, 0o644)
}
