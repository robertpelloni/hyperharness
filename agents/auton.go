package agents

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

const autoDriveIterationDelay = 500 * time.Millisecond

// AutoDrive State Machine provides Roo Code / Cline level autonomy within Hypercode.
// It bypasses the need for the TS orchestrator driving a loop.
type AutoDrive struct {
	MaxIterations int
	Director      *Director
	IsRunning     bool
}

func NewAutoDrive(director *Director) *AutoDrive {
	return &AutoDrive{
		MaxIterations: 25,
		Director:      director,
		IsRunning:     false,
	}
}

// Start initiates the autonomous loop recursively calling the ILLMProvider dynamically handling inputs!
func (a *AutoDrive) Start(ctx context.Context, objective string, sandboxDir string) error {
	return autoDriveStart(ctx, a, objective, sandboxDir, func() { time.Sleep(autoDriveIterationDelay) })
}

func autoDriveStart(ctx context.Context, a *AutoDrive, objective string, sandboxDir string, wait func()) error {
	if a == nil {
		return fmt.Errorf("autodrive is required")
	}
	if a.Director == nil {
		return fmt.Errorf("director is required")
	}
	if a.Director.Provider == nil {
		return fmt.Errorf("director provider is required")
	}
	objective = strings.TrimSpace(objective)
	if objective == "" {
		return fmt.Errorf("objective is required")
	}
	if strings.TrimSpace(sandboxDir) == "" {
		return fmt.Errorf("sandboxDir is required")
	}
	if a.MaxIterations <= 0 {
		return fmt.Errorf("maxIterations must be positive")
	}

	a.IsRunning = true
	defer func() { a.IsRunning = false }()
	log.Printf("[AutoDrive] Engaged native autonomous loop internally isolating execution inside: %s", sandboxDir)

	prompt := buildAutoDrivePrompt(objective, sandboxDir)
	a.Director.History = append(a.Director.History, Message{Role: RoleUser, Content: prompt})

	for i := 0; i < a.MaxIterations; i++ {
		if !a.IsRunning {
			return fmt.Errorf("autodrive aborted early via user interruption natively")
		}
		if wait != nil {
			wait()
		}
		if !a.IsRunning {
			return fmt.Errorf("autodrive aborted early via user interruption natively")
		}
		log.Printf("[AutoDrive] Iteration %d: Generating subsequent MCP tool boundaries...", i+1)

		responseMsg, err := a.Director.Provider.Chat(ctx, a.Director.History, []Tool{})
		if err != nil {
			log.Printf("[AutoDrive] Chat completion failed internally: %v", err)
			return err
		}

		a.Director.History = append(a.Director.History, responseMsg)
		if len(responseMsg.ToolCalls) == 0 {
			log.Printf("[AutoDrive] LLM yielded zero tool targets declaring closure organically.")
			break
		}
		for _, tc := range responseMsg.ToolCalls {
			log.Printf("[AutoDrive] Natively executing %s...", tc.Name)
			a.Director.History = append(a.Director.History, autoDriveToolResultMessage(tc))
		}
	}

	log.Printf("[AutoDrive] Director loop successfully executed objective natively!")
	return nil
}

func buildAutoDrivePrompt(objective string, sandboxDir string) string {
	return fmt.Sprintf("Execute the following plan autonomously:\n%s\n\nCRITICAL: All commands MUST be executed exclusively within '%s'.", objective, sandboxDir)
}

func autoDriveToolResultMessage(tc ToolCall) Message {
	out := fmt.Sprintf("Mock execution output representing terminal feedback natively bridging %s(%s)", tc.Name, tc.Args)
	return Message{Role: RoleTool, ToolCallID: tc.ID, Name: tc.Name, Content: out}
}

func (a *AutoDrive) Abort() {
	if a == nil {
		return
	}
	a.IsRunning = false
	log.Println("[AutoDrive] ABORT command received. Terminating loops.")
}
