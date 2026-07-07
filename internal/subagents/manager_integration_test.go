package subagents

import (
	"context"
	"testing"
	"time"
)

func TestIntegration_SubagentManagerStreamLoop(t *testing.T) {
	mgr := NewManager()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var chunks []string
	callback := func(chunk string) {
		chunks = append(chunks, chunk)
	}

	prompt := "What is 2+2? Only answer with the number."
	output, err := mgr.Spawn(ctx, TypeResearch, prompt, prompt, "Test Context", callback)

	if err != nil {
		// It's acceptable for it to fallback if the API is offline
		t.Logf("Offline or API key missing fallback handled: %v", err)
	}

	if output == "" {
		t.Fatal("Expected non-empty output")
	}

	// Wait, we need to check if the stream actually ran
	if len(chunks) == 0 {
		t.Fatal("Expected stream chunks to be collected")
	}

	t.Logf("Integration stream successfully completed with %d chunks.", len(chunks))
}
