package tui_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/robertpelloni/hyperharness/internal/subagents"
	"github.com/robertpelloni/hyperharness/tui"
)

// TestDashboardObservability explicitly validates the observability hooks in the TUI
// by verifying that spawned subagents are correctly propagated to the active dashboard.
func TestDashboardObservability(t *testing.T) {
	// Initialize the active dashboard placeholders to reset/verify state
	_, toolContentInit, _ := tui.GenerateDashboardPlaceholders()

	// Spawn a real context-bound task in the subagent manager
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Launch an async task that the dashboard should pick up
	go subagents.GlobalManager.Spawn(ctx, subagents.TypeCode, "Observability validation task", "", "", nil)

	// Wait a moment for it to hit the GlobalManager map
	time.Sleep(10 * time.Millisecond)

	// Fetch new dashboard state
	_, toolContentRunning, _ := tui.GenerateDashboardPlaceholders()

	if strings.Contains(toolContentInit, "Observability validation task") {
		t.Error("Dashboard should not contain the task before it spawns")
	}

	// Wait, the dashboard displays type and ID, not the prompt.
	if !strings.Contains(toolContentRunning, "code - task-") {
		t.Errorf("Expected Dashboard to show running subagent. Got:\n%s", toolContentRunning)
	}
}
