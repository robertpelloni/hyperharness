package subagents_test

import (
	"sync"
	"testing"

	"github.com/robertpelloni/hyperharness/internal/subagents"
	"github.com/robertpelloni/hyperharness/tools"
)

// TestSubagentLoad tests spawning 50 concurrent agents via the Agent parity tool
// and ensures the GlobalManager tracks and releases them cleanly without panics.
func TestSubagentLoad(t *testing.T) {
	registry := tools.NewRegistry()

	var targetTool *tools.Tool
	for _, tool := range registry.Tools {
		if tool.Name == "Agent" {
			targetTool = &tool
			break
		}
	}

	if targetTool == nil {
		t.Fatalf("Agent tool not found in registry")
	}

	const concurrentAgents = 50
	var wg sync.WaitGroup
	wg.Add(concurrentAgents)

	for i := 0; i < concurrentAgents; i++ {
		go func(agentID int) {
			defer wg.Done()
			_, err := targetTool.Execute(map[string]interface{}{
				"prompt": "Test high load task execution",
				"type":   "Explore",
			})
			if err != nil {
				t.Errorf("Agent %d failed: %v", agentID, err)
			}
		}(i)
	}

	wg.Wait()

	// Evaluate GlobalManager State
	tasks := subagents.GlobalManager.ListTasks()

	// Expect at least 50 tasks (might be more if other tests ran in same binary)
	if len(tasks) < concurrentAgents {
		t.Errorf("Expected at least %d tracked tasks, got %d", concurrentAgents, len(tasks))
	}

	for _, task := range tasks {
		if task.Status == "running" {
			t.Errorf("Found orphaned running task: %s", task.ID)
		}
	}
}
