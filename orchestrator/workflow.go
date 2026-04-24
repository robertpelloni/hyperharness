// Package orchestrator provides a DAG-based workflow engine.
// Ported from hypercode/go/internal/workflow/engine.go with enhancements.
//
// WHAT: Directed Acyclic Graph (DAG) workflow engine for multi-step task execution
// WHY: Enables complex multi-tool workflows with dependency management
// HOW: Topological sort for execution order, concurrent step execution,
//      output passing between dependent steps, error propagation
package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// StepStatus represents the execution state of a workflow step.
type StepStatus string

const (
	StepPending   StepStatus = "pending"
	StepRunning   StepStatus = "running"
	StepCompleted StepStatus = "completed"
	StepFailed    StepStatus = "failed"
	StepSkipped   StepStatus = "skipped"
)

// Step represents a single unit of work in a workflow.
type Step struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description,omitempty"`
	DependsOn   []string                `json:"dependsOn,omitempty"`
	Status      StepStatus              `json:"status"`
	StartedAt   *time.Time              `json:"startedAt,omitempty"`
	FinishedAt  *time.Time              `json:"finishedAt,omitempty"`
	Error       string                  `json:"error,omitempty"`
	Output      map[string]interface{}  `json:"output,omitempty"`
	Metadata    map[string]string       `json:"metadata,omitempty"`

	// Execute is the function that performs the step's work.
	// Receives inputs from dependency outputs (keyed as "stepID.field").
	Execute func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) `json:"-"`
}

// Workflow is a directed acyclic graph (DAG) of steps.
type Workflow struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Steps       []*Step    `json:"steps"`
	Status      StepStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	FinishedAt  *time.Time `json:"finishedAt,omitempty"`

	mu sync.RWMutex
}

// NewWorkflow creates a new workflow with the given steps.
// All steps are initialized to pending status.
func NewWorkflow(id, name, description string, steps []*Step) *Workflow {
	for _, s := range steps {
		s.Status = StepPending
	}
	return &Workflow{
		ID:          id,
		Name:        name,
		Description: description,
		Steps:       steps,
		Status:      StepPending,
		CreatedAt:   time.Now().UTC(),
	}
}

// topologicalSort returns steps in execution order, respecting dependencies.
// Returns an error if a circular dependency is detected.
func (w *Workflow) topologicalSort() ([]*Step, error) {
	stepMap := make(map[string]*Step, len(w.Steps))
	for _, s := range w.Steps {
		stepMap[s.ID] = s
	}

	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	var sorted []*Step

	var visit func(id string) error
	visit = func(id string) error {
		if visited[id] {
			return nil
		}
		if visiting[id] {
			return fmt.Errorf("circular dependency detected at step %q", id)
		}

		step, ok := stepMap[id]
		if !ok {
			return fmt.Errorf("step %q not found", id)
		}

		visiting[id] = true
		for _, dep := range step.DependsOn {
			if err := visit(dep); err != nil {
				return err
			}
		}
		visiting[id] = false
		visited[id] = true
		sorted = append(sorted, step)
		return nil
	}

	for _, s := range w.Steps {
		if err := visit(s.ID); err != nil {
			return nil, err
		}
	}

	return sorted, nil
}

// Run executes all steps in topological order.
// Independent steps run concurrently. Failed dependencies skip downstream steps.
func (w *Workflow) Run(ctx context.Context) error {
	w.mu.Lock()
	now := time.Now().UTC()
	w.StartedAt = &now
	w.Status = StepRunning
	w.mu.Unlock()

	sorted, err := w.topologicalSort()
	if err != nil {
		w.mu.Lock()
		w.Status = StepFailed
		fin := time.Now().UTC()
		w.FinishedAt = &fin
		w.mu.Unlock()
		return err
	}

	// Build outputs map for passing between steps
	outputs := &sync.Map{}

	// Build a completion tracker per step
	completed := make(map[string]chan struct{})
	for _, s := range sorted {
		completed[s.ID] = make(chan struct{})
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(sorted))

	for _, step := range sorted {
		wg.Add(1)
		go func(s *Step) {
			defer wg.Done()
			defer close(completed[s.ID])

			// Wait for all dependencies to complete
			for _, dep := range s.DependsOn {
				select {
				case <-ctx.Done():
					w.mu.Lock()
					s.Status = StepSkipped
					w.mu.Unlock()
					return
				case <-completed[dep]:
					// Check if dependency failed
					w.mu.RLock()
					for _, ds := range w.Steps {
						if ds.ID == dep && ds.Status == StepFailed {
							w.mu.RUnlock()
							w.mu.Lock()
							s.Status = StepSkipped
							s.Error = fmt.Sprintf("dependency %q failed", dep)
							w.mu.Unlock()
							return
						}
					}
					w.mu.RUnlock()
				}
			}

			// Gather inputs from dependency outputs
			inputs := make(map[string]interface{})
			for _, dep := range s.DependsOn {
				if val, ok := outputs.Load(dep); ok {
					if m, ok := val.(map[string]interface{}); ok {
						for k, v := range m {
							inputs[dep+"."+k] = v
						}
					}
				}
			}

			// Execute the step
			w.mu.Lock()
			startTime := time.Now().UTC()
			s.StartedAt = &startTime
			s.Status = StepRunning
			w.mu.Unlock()

			var output map[string]interface{}
			var execErr error

			if s.Execute != nil {
				output, execErr = s.Execute(ctx, inputs)
			}

			w.mu.Lock()
			finTime := time.Now().UTC()
			s.FinishedAt = &finTime
			if execErr != nil {
				s.Status = StepFailed
				s.Error = execErr.Error()
				errChan <- fmt.Errorf("step %q failed: %w", s.ID, execErr)
			} else {
				s.Status = StepCompleted
				s.Output = output
				outputs.Store(s.ID, output)
			}
			w.mu.Unlock()
		}(step)
	}

	wg.Wait()
	close(errChan)

	w.mu.Lock()
	finTime := time.Now().UTC()
	w.FinishedAt = &finTime

	// Check for any failures
	var firstErr error
	for err := range errChan {
		if firstErr == nil {
			firstErr = err
		}
	}

	if firstErr != nil {
		w.Status = StepFailed
	} else {
		w.Status = StepCompleted
	}
	w.mu.Unlock()

	return firstErr
}

// GetStep retrieves a step by ID.
func (w *Workflow) GetStep(id string) (*Step, bool) {
	for _, s := range w.Steps {
		if s.ID == id {
			return s, true
		}
	}
	return nil, false
}

// CompletedCount returns the number of completed steps.
func (w *Workflow) CompletedCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	count := 0
	for _, s := range w.Steps {
		if s.Status == StepCompleted {
			count++
		}
	}
	return count
}

// FailedCount returns the number of failed steps.
func (w *Workflow) FailedCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	count := 0
	for _, s := range w.Steps {
		if s.Status == StepFailed {
			count++
		}
	}
	return count
}

// Duration returns the total workflow execution time.
func (w *Workflow) Duration() time.Duration {
	if w.StartedAt == nil {
		return 0
	}
	end := time.Now()
	if w.FinishedAt != nil {
		end = *w.FinishedAt
	}
	return end.Sub(*w.StartedAt)
}

// WorkflowEngine manages multiple workflows.
type WorkflowEngine struct {
	workflows map[string]*Workflow
	mu        sync.RWMutex
}

// NewWorkflowEngine creates a new workflow engine.
func NewWorkflowEngine() *WorkflowEngine {
	return &WorkflowEngine{
		workflows: make(map[string]*Workflow),
	}
}

// Register adds a workflow to the engine.
func (e *WorkflowEngine) Register(wf *Workflow) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.workflows[wf.ID] = wf
}

// Get retrieves a workflow by ID.
func (e *WorkflowEngine) Get(id string) (*Workflow, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	wf, ok := e.workflows[id]
	return wf, ok
}

// List returns all registered workflows.
func (e *WorkflowEngine) List() []*Workflow {
	e.mu.RLock()
	defer e.mu.RUnlock()
	list := make([]*Workflow, 0, len(e.workflows))
	for _, wf := range e.workflows {
		list = append(list, wf)
	}
	return list
}

// RunWorkflow executes a workflow by ID.
func (e *WorkflowEngine) RunWorkflow(ctx context.Context, id string) error {
	wf, ok := e.Get(id)
	if !ok {
		return fmt.Errorf("workflow %q not found", id)
	}
	return wf.Run(ctx)
}

// Delete removes a workflow from the engine.
func (e *WorkflowEngine) Delete(id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.workflows, id)
}
