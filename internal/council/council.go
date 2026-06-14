package council

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/robertpelloni/hyperharness/internal/debate"
)

// Executor defines the interface for executing LLM requests.
type Executor interface {
	Chat(input string) (string, error)
}

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
	exec      Executor
	workers   map[string]*WorkerAgent
	taskQueue []*SubagentTask
}

func NewDirectorAgent() *DirectorAgent {
	return &DirectorAgent{
		workers:   make(map[string]*WorkerAgent),
		taskQueue: make([]*SubagentTask, 0),
	}
}

// PlanDelegation breaks a high-level goal into subtasks using an LLM.
func (d *DirectorAgent) PlanDelegation(ctx context.Context, exec Executor, goal string) ([]*SubagentTask, error) {
	fmt.Printf("[Director] Planning delegation for goal: %s\n", goal)
	d.exec = exec

	prompt := fmt.Sprintf(`Break down the following goal into a sequence of sub-tasks for specialized agents.
Available roles: researcher, coder, reviewer, planner.
Output MUST be a JSON array of objects with "id", "role", and "objective" fields.

Goal: %s`, goal)

	response, err := exec.Chat(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get delegation plan: %w", err)
	}

	// Try to extract JSON if there's markdown wrapping
	jsonStr := response
	if idx := strings.Index(jsonStr, "["); idx != -1 {
		if lastIdx := strings.LastIndex(jsonStr, "]"); lastIdx != -1 {
			jsonStr = jsonStr[idx : lastIdx+1]
		}
	}

	var plan []*SubagentTask
	if err := json.Unmarshal([]byte(jsonStr), &plan); err != nil {
		// Fallback to basic plan if parsing fails
		fmt.Printf("[Director] Failed to parse LLM plan, using fallback: %v\n", err)
		plan = []*SubagentTask{
			{ID: "task_1", Role: "researcher", Objective: "Analyze the goal: " + goal, Status: "pending"},
			{ID: "task_2", Role: "coder", Objective: "Implement solution for: " + goal, Status: "pending"},
		}
	}

	for _, t := range plan {
		t.Status = "pending"
	}

	d.taskQueue = append(d.taskQueue, plan...)
	return plan, nil
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
			worker.exec = d.exec
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

// DebateOnFailure triggers a structured debate between two subagents on a
// given topic (typically a failed code edit or validation error) and returns
// a formatted result. This is used when code edits fail validation or tests.
func (d *DirectorAgent) DebateOnFailure(topic string) (string, error) {
	fmt.Printf("[Director] Debate on failure triggered for: %s\n", topic)
	pro := NewWorkerAgent("coder")
	con := NewWorkerAgent("reviewer")
	res, err := debate.Run(topic, pro, con)
	if err != nil {
		return "", fmt.Errorf("debate on failure: %w", err)
	}
	return debate.FormatResult(res), nil
}

// WorkerAgent handles individual subtasks in the council.
type WorkerAgent struct {
	Role string
	exec Executor
}

func NewWorkerAgent(role string) *WorkerAgent {
	return &WorkerAgent{Role: role}
}

// Execute runs the assigned subtask.
func (w *WorkerAgent) Execute(task *SubagentTask) (string, error) {
	if w.exec == nil {
		// Mock execution delay for backward compatibility/testing if no executor
		time.Sleep(100 * time.Millisecond)
		return fmt.Sprintf("Completed objective: %s (Simulated execution)", task.Objective), nil
	}

	prompt := fmt.Sprintf("You are a %s agent. Objective: %s\nPerform the task and provide a summary of your work.", w.Role, task.Objective)
	return w.exec.Chat(prompt)
}

// Argue produces an argument from a specific stance.
func (w *WorkerAgent) Argue(topic, stance, previousArgument string) (string, error) {
	// Mock debate reasoning
	time.Sleep(50 * time.Millisecond)
	return fmt.Sprintf("As a %s taking the %s stance, I argue that this approach is optimal.", w.Role, stance), nil
}
