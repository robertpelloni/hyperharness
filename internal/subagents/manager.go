package subagents

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// SubagentType defines the specialization of a subagent.
type SubagentType string

const (
	TypeCode      SubagentType = "code"
	TypeResearch  SubagentType = "research"
	TypeReview    SubagentType = "review"
	TypePlan      SubagentType = "plan"
	TypeBuild     SubagentType = "build"
	TypeTest      SubagentType = "test"
	TypeDebug     SubagentType = "debug"
	TypeDoc       SubagentType = "doc"
	TypeSecurity  SubagentType = "security"
	TypeDevOps    SubagentType = "devops"
)

// SubagentState represents the state of a running subagent.
type SubagentState string

const (
	StatePending   SubagentState = "pending"
	StateRunning   SubagentState = "running"
	StateCompleted SubagentState = "completed"
	StateFailed    SubagentState = "failed"
	StateCancelled SubagentState = "cancelled"
)

// SubagentTask defines a task for a subagent to execute.
type SubagentTask struct {
	ID          string                 `json:"id"`
	Type        SubagentType           `json:"type"`
	Description string                 `json:"description"`
	Prompt      string                 `json:"prompt"`
	ParentID    string                 `json:"parent_id,omitempty"`
	SessionID   string                 `json:"session_id,omitempty"`
	State       SubagentState          `json:"state"`
	Result      string                 `json:"result,omitempty"`
	Error       string                 `json:"error,omitempty"`
	StartTime   time.Time              `json:"start_time,omitempty"`
	EndTime     time.Time              `json:"end_time,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SubagentConfig configures a subagent type.
type SubagentConfig struct {
	Type         SubagentType `json:"type"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	MaxTokens    int          `json:"max_tokens"`
	MaxTurns     int          `json:"max_turns"`
	AllowedTools []string     `json:"allowed_tools"`
	DeniedTools  []string     `json:"denied_tools"`
	SystemPrompt string       `json:"system_prompt"`
}

// DefaultSubagentConfigs returns the default subagent configurations.
func DefaultSubagentConfigs() map[SubagentType]SubagentConfig {
	return map[SubagentType]SubagentConfig{
		TypeCode: {
			Type:        TypeCode,
			Name:        "Code Agent",
			Description: "Writes, edits, and refactors code. Has full file system and shell access.",
			MaxTokens:   128000,
			MaxTurns:    50,
			AllowedTools: []string{"read", "write", "edit", "bash", "grep", "find", "ls"},
		},
		TypeResearch: {
			Type:        TypeResearch,
			Name:        "Research Agent",
			Description: "Searches codebase and web for information. Read-only access.",
			MaxTokens:   64000,
			MaxTurns:    20,
			AllowedTools: []string{"read", "grep", "find", "ls", "websearch", "webfetch"},
		},
		TypeReview: {
			Type:        TypeReview,
			Name:        "Review Agent",
			Description: "Reviews code for bugs, style, and best practices. Read-only access.",
			MaxTokens:   64000,
			MaxTurns:    15,
			AllowedTools: []string{"read", "grep", "find", "ls"},
		},
		TypePlan: {
			Type:        TypePlan,
			Name:        "Plan Agent",
			Description: "Creates detailed implementation plans. Read-only access to codebase.",
			MaxTokens:   64000,
			MaxTurns:    25,
			AllowedTools: []string{"read", "grep", "find", "ls", "tree"},
		},
		TypeBuild: {
			Type:        TypeBuild,
			Name:        "Build Agent",
			Description: "Executes build plans created by the Plan agent. Full access.",
			MaxTokens:   128000,
			MaxTurns:    75,
			AllowedTools: []string{"read", "write", "edit", "bash", "grep", "find", "ls"},
		},
		TypeTest: {
			Type:        TypeTest,
			Name:        "Test Agent",
			Description: "Writes and runs tests. Has file and shell access.",
			MaxTokens:   96000,
			MaxTurns:    40,
			AllowedTools: []string{"read", "write", "edit", "bash", "grep", "find", "ls"},
		},
		TypeDebug: {
			Type:        TypeDebug,
			Name:        "Debug Agent",
			Description: "Diagnoses and fixes bugs. Has full access with diagnostic tools.",
			MaxTokens:   96000,
			MaxTurns:    50,
			AllowedTools: []string{"read", "write", "edit", "bash", "grep", "find", "ls", "diagnostics"},
		},
		TypeDoc: {
			Type:        TypeDoc,
			Name:        "Documentation Agent",
			Description: "Writes and updates documentation. Has file system access.",
			MaxTokens:   64000,
			MaxTurns:    30,
			AllowedTools: []string{"read", "write", "edit", "grep", "find", "ls"},
		},
		TypeSecurity: {
			Type:        TypeSecurity,
			Name:        "Security Agent",
			Description: "Audits code for security vulnerabilities. Read-only with diagnostic access.",
			MaxTokens:   64000,
			MaxTurns:    25,
			AllowedTools: []string{"read", "grep", "find", "ls", "diagnostics", "safe"},
		},
		TypeDevOps: {
			Type:        TypeDevOps,
			Name:        "DevOps Agent",
			Description: "Handles deployment, CI/CD, and infrastructure tasks.",
			MaxTokens:   64000,
			MaxTurns:    30,
			AllowedTools: []string{"read", "write", "edit", "bash", "grep", "find", "ls"},
		},
	}
}

// Manager manages subagent lifecycle and execution.
type Manager struct {
	tasks    map[string]*SubagentTask
	configs  map[SubagentType]SubagentConfig
	mu       sync.RWMutex
	taskSeq  int64
}

// NewManager creates a new subagent manager.
func NewManager() *Manager {
	return &Manager{
		tasks:   make(map[string]*SubagentTask),
		configs: DefaultSubagentConfigs(),
	}
}

// CreateTask creates a new subagent task.
func (m *Manager) CreateTask(agentType SubagentType, description, prompt, parentID string) *SubagentTask {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.taskSeq++
	task := &SubagentTask{
		ID:          fmt.Sprintf("sub_%d_%d", time.Now().UnixNano(), m.taskSeq),
		Type:        agentType,
		Description: description,
		Prompt:      prompt,
		ParentID:    parentID,
		State:       StatePending,
		StartTime:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}

	m.tasks[task.ID] = task
	return task
}

// ExecuteTask runs a subagent task synchronously.
// In a full implementation, this would spin up an agent loop with
// the subagent's configuration and tool restrictions.
func (m *Manager) ExecuteTask(ctx context.Context, task *SubagentTask) (string, error) {
	m.mu.Lock()
	task.State = StateRunning
	m.mu.Unlock()

	// Get config for this agent type
	config, ok := m.configs[task.Type]
	if !ok {
		m.mu.Lock()
		task.State = StateFailed
		task.Error = fmt.Sprintf("unknown agent type: %s", task.Type)
		m.mu.Unlock()
		return "", fmt.Errorf("unknown agent type: %s", task.Type)
	}

	// In production, this would:
	// 1. Create an isolated agent session
	// 2. Configure allowed/denied tools from config
	// 3. Set the system prompt from config
	// 4. Run the agent loop with the prompt
	// 5. Collect and return results

	result := fmt.Sprintf(
		"Subagent task completed (type: %s, max_turns: %d, allowed_tools: %v)\n\nTask: %s\nPrompt: %s",
		config.Type, config.MaxTurns, config.AllowedTools,
		task.Description, task.Prompt,
	)

	m.mu.Lock()
	task.State = StateCompleted
	task.Result = result
	task.EndTime = time.Now()
	m.mu.Unlock()

	return result, nil
}

// GetTask retrieves a task by ID.
func (m *Manager) GetTask(id string) (*SubagentTask, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	task, ok := m.tasks[id]
	return task, ok
}

// ListTasks returns all tasks, optionally filtered by state.
func (m *Manager) ListTasks(state SubagentState) []*SubagentTask {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var tasks []*SubagentTask
	for _, task := range m.tasks {
		if state == "" || task.State == state {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// CancelTask cancels a running task.
func (m *Manager) CancelTask(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, ok := m.tasks[id]
	if !ok {
		return fmt.Errorf("task not found: %s", id)
	}

	if task.State != StateRunning && task.State != StatePending {
		return fmt.Errorf("task %s is not cancellable (state: %s)", id, task.State)
	}

	task.State = StateCancelled
	task.EndTime = time.Now()
	return nil
}

// GetConfig returns the configuration for a subagent type.
func (m *Manager) GetConfig(agentType SubagentType) (SubagentConfig, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	config, ok := m.configs[agentType]
	return config, ok
}

// ListConfigs returns all available subagent configurations.
func (m *Manager) ListConfigs() []SubagentConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var configs []SubagentConfig
	for _, config := range m.configs {
		configs = append(configs, config)
	}
	return configs
}
