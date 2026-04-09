// Package orchestrator provides a persistent task queue for autonomous work.
// Ported from hypercode/go/internal/orchestration/task_queue.go
//
// WHAT: Thread-safe task queue with priorities, status tracking, and persistence
// WHY: Enables autonomous task execution with ordering and retry support
// HOW: In-memory queue with priority ordering and configurable workers
package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// TaskPriority defines task priority levels.
type TaskPriority int

const (
	PriorityLow      TaskPriority = 0
	PriorityNormal   TaskPriority = 1
	PriorityHigh     TaskPriority = 2
	PriorityCritical TaskPriority = 3
)

// TaskStatus defines the lifecycle states of a task.
type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskRunning   TaskStatus = "running"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
	TaskCancelled TaskStatus = "cancelled"
	TaskRetrying  TaskStatus = "retrying"
)

// Task represents a unit of autonomous work.
type Task struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Priority    TaskPriority           `json:"priority"`
	Status      TaskStatus             `json:"status"`
	CreatedAt   time.Time              `json:"createdAt"`
	StartedAt   *time.Time             `json:"startedAt,omitempty"`
	FinishedAt  *time.Time             `json:"finishedAt,omitempty"`
	Result      string                 `json:"result,omitempty"`
	Error       string                 `json:"error,omitempty"`
	MaxRetries  int                    `json:"maxRetries"`
	Retries     int                    `json:"retries"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`

	// Execute performs the task's work
	Execute func(ctx context.Context) (string, error) `json:"-"`
}

// PriorityQueue manages prioritized task execution in memory.
type PriorityQueue struct {
	tasks    map[string]*Task
	pending  []*Task
	mu       sync.RWMutex
	nextID   int64
	workers  int
	running  int64
	stopped  int64
}

// NewPriorityQueue creates a new priority queue.
func NewPriorityQueue(workers int) *PriorityQueue {
	if workers <= 0 {
		workers = 4
	}
	return &PriorityQueue{
		tasks:   make(map[string]*Task),
		pending: make([]*Task, 0),
		workers: workers,
	}
}

// Enqueue adds a task to the priority queue.
func (q *PriorityQueue) Enqueue(task *Task) string {
	q.mu.Lock()
	defer q.mu.Unlock()

	if task.ID == "" {
		q.nextID++
		task.ID = fmt.Sprintf("task-%d", q.nextID)
	}
	task.CreatedAt = time.Now().UTC()
	task.Status = TaskPending

	q.tasks[task.ID] = task
	q.pending = append(q.pending, task)

	// Sort pending by priority (highest first)
	for i := len(q.pending) - 1; i > 0; i-- {
		if q.pending[i].Priority > q.pending[i-1].Priority {
			q.pending[i], q.pending[i-1] = q.pending[i-1], q.pending[i]
		}
	}

	return task.ID
}

// Dequeue gets the highest-priority pending task.
func (q *PriorityQueue) Dequeue() *Task {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.pending) == 0 {
		return nil
	}

	task := q.pending[0]
	q.pending = q.pending[1:]
	return task
}

// Process executes tasks from the queue using worker goroutines.
func (q *PriorityQueue) Process(ctx context.Context) {
	for i := 0; i < q.workers; i++ {
		go q.worker(ctx)
	}
}

// worker runs tasks from the queue.
func (q *PriorityQueue) worker(ctx context.Context) {
	for {
		if atomic.LoadInt64(&q.stopped) == 1 {
			return
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		task := q.Dequeue()
		if task == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		atomic.AddInt64(&q.running, 1)
		q.executeTask(ctx, task)
		atomic.AddInt64(&q.running, -1)
	}
}

// executeTask runs a single task with retry support.
func (q *PriorityQueue) executeTask(ctx context.Context, task *Task) {
	q.mu.Lock()
	task.Status = TaskRunning
	now := time.Now().UTC()
	task.StartedAt = &now
	q.mu.Unlock()

	var result string
	var err error

	if task.Execute != nil {
		result, err = task.Execute(ctx)
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	fin := time.Now().UTC()
	task.FinishedAt = &fin

	if err != nil {
		if task.Retries < task.MaxRetries {
			task.Retries++
			task.Status = TaskRetrying
			// Re-enqueue for retry
			q.pending = append(q.pending, task)
		} else {
			task.Status = TaskFailed
			task.Error = err.Error()
		}
	} else {
		task.Status = TaskCompleted
		task.Result = result
	}
}

// Stop signals all workers to stop.
func (q *PriorityQueue) Stop() {
	atomic.StoreInt64(&q.stopped, 1)
}

// Get retrieves a task by ID.
func (q *PriorityQueue) Get(id string) (*Task, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	task, ok := q.tasks[id]
	return task, ok
}

// List returns all tasks.
func (q *PriorityQueue) List() []*Task {
	q.mu.RLock()
	defer q.mu.RUnlock()
	tasks := make([]*Task, 0, len(q.tasks))
	for _, t := range q.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

// ListByStatus returns tasks filtered by status.
func (q *PriorityQueue) ListByStatus(status TaskStatus) []*Task {
	q.mu.RLock()
	defer q.mu.RUnlock()
	var tasks []*Task
	for _, t := range q.tasks {
		if t.Status == status {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

// PendingCount returns the number of pending tasks.
func (q *PriorityQueue) PendingCount() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.pending)
}

// RunningCount returns the number of running tasks.
func (q *PriorityQueue) RunningCount() int64 {
	return atomic.LoadInt64(&q.running)
}

// TotalCount returns the total number of tasks ever enqueued.
func (q *PriorityQueue) TotalCount() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.tasks)
}

// Cancel cancels a pending task.
func (q *PriorityQueue) Cancel(id string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	task, ok := q.tasks[id]
	if !ok {
		return fmt.Errorf("task %q not found", id)
	}
	if task.Status != TaskPending && task.Status != TaskRetrying {
		return fmt.Errorf("task %q is %s, cannot cancel", id, task.Status)
	}

	task.Status = TaskCancelled

	// Remove from pending
	for i, t := range q.pending {
		if t.ID == id {
			q.pending = append(q.pending[:i], q.pending[i+1:]...)
			break
		}
	}

	return nil
}

// Stats returns queue statistics.
func (q *PriorityQueue) Stats() map[string]interface{} {
	q.mu.RLock()
	defer q.mu.RUnlock()

	stats := map[string]int{
		"pending":   0,
		"running":   0,
		"completed": 0,
		"failed":    0,
		"cancelled": 0,
		"retrying":  0,
	}

	for _, t := range q.tasks {
		stats[string(t.Status)]++
	}

	return map[string]interface{}{
		"total":     len(q.tasks),
		"workers":   q.workers,
		"active":    atomic.LoadInt64(&q.running),
		"byStatus":  stats,
	}
}
