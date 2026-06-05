package council

import (
	"fmt"
	"strings"
	"time"
)

// SubagentTask represents a task dispatched by a Director to a Worker.
type SubagentTask struct {
	ID           string
	Role         string // "researcher", "coder", "reviewer", "planner"
	Objective    string
	ContextFiles []string
	Result       string
	Status       string // "pending", "running", "completed", "failed"
}

// DirectorAgent orchestrates subagents based on the Council architecture.
type DirectorAgent struct {
	workers   map[string]*WorkerAgent
	taskQueue []*SubagentTask
}

func NewDirectorAgent() *DirectorAgent {
	return &DirectorAgent{
		workers:   make(map[string]*WorkerAgent),
		taskQueue: make([]*SubagentTask, 0),
	}
}

// PlanDelegation breaks a high-level goal into subtasks.
func (d *DirectorAgent) PlanDelegation(goal string) []*SubagentTask {
	fmt.Printf("[Director] Planning delegation for goal: %s\n", goal)
	
	plan := []*SubagentTask{
		{ID: "task_1", Role: "researcher", Objective: "Find relevant files for the goal.", Status: "pending"},
		{ID: "task_2", Role: "coder", Objective: "Implement the required changes.", Status: "pending"},
		{ID: "task_3", Role: "reviewer", Objective: "Review the code modifications for errors.", Status: "pending"},
	}
	
	d.taskQueue = append(d.taskQueue, plan...)
	return plan
}

// ExecutePlan runs the queued tasks sequentially with a self-correction loop.
func (d *DirectorAgent) ExecutePlan() (string, error) {
	var results []string
	
	for _, task := range d.taskQueue {
		fmt.Printf("[Director] Dispatching task %s to %s...\n", task.ID, task.Role)
		
		maxAttempts := 3
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			task.Status = "running"
			worker := NewWorkerAgent(task.Role)
			res, err := worker.Execute(task)

			if err != nil {
				fmt.Printf("[Director] Task %s failed (attempt %d/%d): %v\n", task.ID, attempt, maxAttempts, err)
				if attempt == maxAttempts {
					task.Status = "failed"
					return "", fmt.Errorf("task %s failed after %d attempts: %w", task.ID, maxAttempts, err)
				}
				fmt.Println("[Director] Attempting self-correction...")
				task.Objective += " (Retry with correction)"
				continue
			}

			task.Result = res
			task.Status = "completed"
			results = append(results, fmt.Sprintf("[%s output] %s", task.Role, res))
			break
		}
	}
	
	fmt.Println("[Director] Synthesizing final response...")
	return strings.Join(results, "\n\n"), nil
}

// InitiateDebate prompts two agents to debate a topic and summarizes the result.
func (d *DirectorAgent) InitiateDebate(topic string, agentA, agentB *WorkerAgent) (string, error) {
	fmt.Printf("[Director] Initiating debate on topic: %s\n", topic)
	
	responseA, err := agentA.Argue(topic, "Proponent", "")
	if err != nil {
		return "", err
	}
	
	responseB, err := agentB.Argue(topic, "Skeptic", responseA)
	if err != nil {
		return "", err
	}
	
	return fmt.Sprintf("Debate Resolution:\nAgent A: %s\nAgent B: %s\nConsensus reached by Director.", responseA, responseB), nil
}

// WorkerAgent handles individual subtasks in the council.
type WorkerAgent struct {
	Role string
}

func NewWorkerAgent(role string) *WorkerAgent {
	return &WorkerAgent{Role: role}
}

// Execute runs the assigned subtask.
func (w *WorkerAgent) Execute(task *SubagentTask) (string, error) {
	// Mock execution delay
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintf("Completed objective: %s (Simulated execution)", task.Objective), nil
}

// Argue produces an argument from a specific stance.
func (w *WorkerAgent) Argue(topic, stance, previousArgument string) (string, error) {
	// Mock debate reasoning
	time.Sleep(50 * time.Millisecond)
	return fmt.Sprintf("As a %s taking the %s stance, I argue that this approach is optimal.", w.Role, stance), nil
}
