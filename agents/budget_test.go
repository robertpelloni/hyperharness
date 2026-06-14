package agents

import (
	"testing"
	"time"
)

func TestDefaultBudgetConfig(t *testing.T) {
	cfg := DefaultBudgetConfig()
	if cfg.MaxCostPerSession <= 0 {
		t.Errorf("MaxCostPerSession = %f, want > 0", cfg.MaxCostPerSession)
	}
	if cfg.MaxCostPerDay <= 0 {
		t.Errorf("MaxCostPerDay = %f, want > 0", cfg.MaxCostPerDay)
	}
	if cfg.Enabled {
		t.Error("Enabled should be false by default")
	}
	if cfg.WarningThresholdPct != 80.0 {
		t.Errorf("WarningThresholdPct = %f, want 80.0", cfg.WarningThresholdPct)
	}
}

func TestNewBudgetTracker(t *testing.T) {
	cfg := DefaultBudgetConfig()
	bt := NewBudgetTracker(cfg)
	if bt == nil {
		t.Fatal("NewBudgetTracker() returned nil")
	}
	if len(bt.GetEntries()) != 0 {
		t.Error("New tracker should have no entries")
	}
}

func TestBudgetTracker_RecordCost_Disabled(t *testing.T) {
	cfg := DefaultBudgetConfig()
	cfg.Enabled = false
	bt := NewBudgetTracker(cfg)

	err := bt.RecordCost("openai", "gpt-4o", 1000, 500, 0.015)
	if err != nil {
		t.Fatalf("RecordCost() error = %v when disabled", err)
	}
	if len(bt.GetEntries()) != 1 {
		t.Errorf("entries = %d, want 1", len(bt.GetEntries()))
	}
}

func TestBudgetTracker_RecordCost_Enabled(t *testing.T) {
	cfg := DefaultBudgetConfig()
	cfg.Enabled = true
	cfg.MaxCostPerSession = 1.00
	bt := NewBudgetTracker(cfg)

	// First call should succeed
	err := bt.RecordCost("openai", "gpt-4o", 1000, 500, 0.015)
	if err != nil {
		t.Fatalf("first RecordCost() error = %v", err)
	}

	// Second call that exceeds budget should fail
	err = bt.RecordCost("openai", "gpt-4o", 50000, 25000, 0.99)
	if err == nil {
		t.Error("second RecordCost() should error when budget exceeded")
	}
}

func TestBudgetTracker_RecordCost_TokenLimit(t *testing.T) {
	cfg := DefaultBudgetConfig()
	cfg.Enabled = true
	cfg.MaxCostPerSession = 100.00 // High cost limit
	cfg.MaxTokensPerSession = 1000
	bt := NewBudgetTracker(cfg)

	// First call should succeed
	err := bt.RecordCost("openai", "gpt-4o", 500, 400, 0.01)
	if err != nil {
		t.Fatalf("first RecordCost() error = %v", err)
	}

	// Second call that exceeds token limit should fail
	err = bt.RecordCost("openai", "gpt-4o", 200, 200, 0.005)
	if err == nil {
		t.Error("second RecordCost() should error when token budget exceeded")
	}
}

func TestBudgetTracker_BudgetStatus(t *testing.T) {
	cfg := DefaultBudgetConfig()
	cfg.Enabled = true
	bt := NewBudgetTracker(cfg)

	// Record some costs
	bt.RecordCost("openai", "gpt-4o", 1000, 500, 0.015)
	bt.RecordCost("anthropic", "claude-sonnet", 800, 400, 0.012)

	status := bt.BudgetStatus()
	if status.SessionCost != 0.027 {
		t.Errorf("SessionCost = %f, want 0.027", status.SessionCost)
	}
	if status.SessionTokens != 2700 {
		t.Errorf("SessionTokens = %d, want 2700", status.SessionTokens)
	}
	if status.TotalEntries != 2 {
		t.Errorf("TotalEntries = %d, want 2", status.TotalEntries)
	}
}

func TestBudgetTracker_IsNearWarning(t *testing.T) {
	cfg := DefaultBudgetConfig()
	cfg.Enabled = true
	cfg.MaxCostPerSession = 1.00
	cfg.WarningThresholdPct = 80.0
	bt := NewBudgetTracker(cfg)

	// Not near warning yet
	if bt.IsNearWarning() {
		t.Error("IsNearWarning() = true, want false")
	}

	// Record cost near threshold
	bt.RecordCost("openai", "gpt-4o", 0, 0, 0.82)
	if !bt.IsNearWarning() {
		t.Error("IsNearWarning() = false, want true (82% > 80%)")
	}
}

func TestBudgetTracker_ResetSession(t *testing.T) {
	cfg := DefaultBudgetConfig()
	cfg.Enabled = true
	cfg.MaxCostPerSession = 0.05
	bt := NewBudgetTracker(cfg)

	// Record a cost
	bt.RecordCost("openai", "gpt-4o", 1000, 500, 0.015)

	// Reset session
	time.Sleep(10 * time.Millisecond) // Ensure time difference
	bt.ResetSession()

	// Session cost should now be ~0 (old entries are before new session start)
	status := bt.BudgetStatus()
	if status.SessionCost > 0.001 {
		t.Errorf("SessionCost after reset = %f, want ~0", status.SessionCost)
	}
}

func TestBudgetTracker_SetConfig(t *testing.T) {
	cfg := DefaultBudgetConfig()
	bt := NewBudgetTracker(cfg)

	newCfg := BudgetConfig{
		MaxCostPerSession:   10.00,
		MaxCostPerDay:       50.00,
		MaxTokensPerSession: 100000,
		MaxTokensPerDay:     500000,
		WarningThresholdPct: 90.0,
		Enabled:             true,
	}
	bt.SetConfig(newCfg)

	got := bt.GetConfig()
	if got.MaxCostPerSession != newCfg.MaxCostPerSession {
		t.Errorf("MaxCostPerSession = %f, want %f", got.MaxCostPerSession, newCfg.MaxCostPerSession)
	}
	if got.WarningThresholdPct != newCfg.WarningThresholdPct {
		t.Errorf("WarningThresholdPct = %f, want %f", got.WarningThresholdPct, newCfg.WarningThresholdPct)
	}
}

func TestBudgetTracker_DailyLimits(t *testing.T) {
	cfg := DefaultBudgetConfig()
	cfg.Enabled = true
	cfg.MaxCostPerDay = 0.02
	cfg.MaxCostPerSession = 100.00 // High session limit
	bt := NewBudgetTracker(cfg)

	// First call should succeed
	err := bt.RecordCost("openai", "gpt-4o", 1000, 500, 0.015)
	if err != nil {
		t.Fatalf("first RecordCost() error = %v", err)
	}

	// Second call that exceeds daily limit should fail
	err = bt.RecordCost("openai", "gpt-4o", 1000, 500, 0.010)
	if err == nil {
		t.Error("second RecordCost() should error when daily budget exceeded")
	}
}

func TestBudgetTracker_GetEntries(t *testing.T) {
	cfg := DefaultBudgetConfig()
	bt := NewBudgetTracker(cfg)

	bt.RecordCost("openai", "gpt-4o", 1000, 500, 0.015)
	bt.RecordCost("anthropic", "claude-sonnet", 800, 400, 0.012)

	entries := bt.GetEntries()
	if len(entries) != 2 {
		t.Errorf("entries = %d, want 2", len(entries))
	}
	if entries[0].Provider != "openai" {
		t.Errorf("entries[0].Provider = %q, want 'openai'", entries[0].Provider)
	}
	if entries[1].Provider != "anthropic" {
		t.Errorf("entries[1].Provider = %q, want 'anthropic'", entries[1].Provider)
	}
}

func Test_percentage(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		max   float64
		want  float64
	}{
		{"half", 50, 100, 50},
		{"full", 100, 100, 100},
		{"over", 150, 100, 100}, // Should cap at 100
		{"zero max", 50, 0, 0},
		{"negative max", 50, -10, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := percentage(tt.value, tt.max)
			if got != tt.want {
				t.Errorf("percentage(%f, %f) = %f, want %f", tt.value, tt.max, got, tt.want)
			}
		})
	}
}

func TestBudgetStatus_Structure(t *testing.T) {
	cfg := DefaultBudgetConfig()
	bt := NewBudgetTracker(cfg)
	status := bt.BudgetStatus()

	// Just verify the structure is populated
	if status.SessionStart.IsZero() {
		t.Error("SessionStart should not be zero")
	}
	if status.WarningThreshold != cfg.WarningThresholdPct {
		t.Errorf("WarningThreshold = %f, want %f", status.WarningThreshold, cfg.WarningThresholdPct)
	}
}