package workflow

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type StepStatus string

const (
	StatusPending   StepStatus = "pending"
	StatusRunning   StepStatus = "running"
	StatusCompleted StepStatus = "completed"
	StatusFailed    StepStatus = "failed"
	StatusSkipped   StepStatus = "skipped"
)

type Step struct {
	ID          string
	Name        string
	Description string
	DependsOn   []string
	Status      StepStatus
	StartedAt   *time.Time
	FinishedAt  *time.Time
	Error       string
	Output      map[string]any
	Metadata    map[string]string
	Execute     func(ctx context.Context, inputs map[string]any) (map[string]any, error)
}

type Workflow struct {
	ID          string
	Name        string
	Description string
	Steps       []*Step
	Status      StepStatus
	CreatedAt   time.Time
	StartedAt   *time.Time
	FinishedAt  *time.Time
	mu          sync.RWMutex
}

func NewWorkflow(id, name, description string, steps []*Step) *Workflow {
	for _, s := range steps {
		s.Status = StatusPending
	}
	return &Workflow{ID: id, Name: name, Description: description, Steps: steps, Status: StatusPending, CreatedAt: time.Now().UTC()}
}

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

func (w *Workflow) Run(ctx context.Context) error {
	w.mu.Lock()
	now := time.Now().UTC()
	w.StartedAt = &now
	w.Status = StatusRunning
	w.mu.Unlock()

	sorted, err := w.topologicalSort()
	if err != nil {
		w.mu.Lock()
		w.Status = StatusFailed
		fin := time.Now().UTC()
		w.FinishedAt = &fin
		w.mu.Unlock()
		return err
	}

	outputs := &sync.Map{}
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

			for _, dep := range s.DependsOn {
				select {
				case <-ctx.Done():
					w.mu.Lock()
					s.Status = StatusSkipped
					w.mu.Unlock()
					return
				case <-completed[dep]:
					w.mu.RLock()
					for _, ds := range w.Steps {
						if ds.ID == dep && ds.Status == StatusFailed {
							w.mu.RUnlock()
							w.mu.Lock()
							s.Status = StatusSkipped
							s.Error = fmt.Sprintf("dependency %q failed", dep)
							w.mu.Unlock()
							return
						}
					}
					w.mu.RUnlock()
				}
			}

			inputs := make(map[string]any)
			for _, dep := range s.DependsOn {
				if val, ok := outputs.Load(dep); ok {
					if m, ok := val.(map[string]any); ok {
						for k, v := range m {
							inputs[dep+"."+k] = v
						}
					}
				}
			}

			w.mu.Lock()
			startTime := time.Now().UTC()
			s.StartedAt = &startTime
			s.Status = StatusRunning
			w.mu.Unlock()

			var output map[string]any
			var execErr error
			if s.Execute != nil {
				output, execErr = s.Execute(ctx, inputs)
			}

			w.mu.Lock()
			finTime := time.Now().UTC()
			s.FinishedAt = &finTime
			if execErr != nil {
				s.Status = StatusFailed
				s.Error = execErr.Error()
				errChan <- fmt.Errorf("step %q failed: %w", s.ID, execErr)
			} else {
				s.Status = StatusCompleted
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
	var firstErr error
	for err := range errChan {
		if firstErr == nil {
			firstErr = err
		}
	}
	if firstErr != nil {
		w.Status = StatusFailed
	} else {
		w.Status = StatusCompleted
	}
	w.mu.Unlock()
	return firstErr
}

type Engine struct {
	workflows map[string]*Workflow
	mu        sync.RWMutex
}

func NewEngine() *Engine {
	return &Engine{workflows: make(map[string]*Workflow)}
}

func (e *Engine) Register(wf *Workflow) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.workflows[wf.ID] = wf
}

func (e *Engine) Get(id string) (*Workflow, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	wf, ok := e.workflows[id]
	return wf, ok
}

func (e *Engine) List() []*Workflow {
	e.mu.RLock()
	defer e.mu.RUnlock()
	list := make([]*Workflow, 0, len(e.workflows))
	for _, wf := range e.workflows {
		list = append(list, wf)
	}
	return list
}

func (e *Engine) RunWorkflow(ctx context.Context, id string) error {
	wf, ok := e.Get(id)
	if !ok {
		return fmt.Errorf("workflow %q not found", id)
	}
	return wf.Run(ctx)
}
