package orchestration

/**
 * @file director.go
 * @module go/internal/orchestration
 *
 * WHAT: Go-native implementation of the Director loop.
 * Coordinates autonomous task execution using the SwarmController and CoderAgent.
 */

import (
	"context"
	"fmt"
	"time"
)

type Director struct {
	swarm  *SwarmController
	coder  *CoderAgent
	broker *A2ABroker
}

func NewDirector(swarm *SwarmController, coder *CoderAgent, broker *A2ABroker) *Director {
	return &Director{
		swarm:  swarm,
		coder:  coder,
		broker: broker,
	}
}

func (d *Director) StartAutonomousTask(ctx context.Context, goal string) error {
	fmt.Printf("[Go Director] 🎬 Starting autonomous task: %s\n", goal)

	// 1. Run a Swarm session to plan and review
	fmt.Println("[Go Director] 🧠 Convening Swarm for planning...")
	result, err := d.swarm.StartSession(ctx, goal, SwarmSessionConfig{
		MaxTurns:            3,
		CompletionThreshold: 0.8,
		AutoRotate:          true,
	})
	if err != nil {
		return fmt.Errorf("swarm planning failed: %w", err)
	}

	if !result.Success {
		fmt.Println("[Go Director] ⚠️ Swarm did not reach full consensus, but proceeding with current plan.")
	}

	// 2. Delegate implementation to the Coder Agent via A2A
	fmt.Println("[Go Director] 🤖 Delegating implementation to Go Coder...")
	
	taskID := fmt.Sprintf("task-%d", time.Now().Unix())
	msg := A2AMessage{
		ID:        fmt.Sprintf("a2a-%d", time.Now().UnixMilli()),
		Timestamp: time.Now().UnixMilli(),
		Sender:    "DIRECTOR",
		Recipient: d.coder.ID,
		Type:      TaskRequest,
		Payload: map[string]interface{}{
			"task":   goal,
			"plan":   result.Transcript,
			"taskId": taskID,
		},
	}

	d.broker.RouteMessage(msg)

	// In a real implementation, we would wait for the response or monitor the broker
	return nil
}
