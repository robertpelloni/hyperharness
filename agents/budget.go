package agents

import (
	"fmt"
	"sync"
	"time"
)

// BudgetConfig holds budget enforcement settings
type BudgetConfig struct {
	MaxCostPerSession  float64 `json:"maxCostPerSession"`  // Max cost per session in USD
	MaxCostPerDay      float64 `json:"maxCostPerDay"`      // Max cost per day in USD
	MaxTokensPerSession int    `json:"maxTokensPerSession"` // Max tokens per session
	MaxTokensPerDay     int    `json:"maxTokensPerDay"`     // Max tokens per day
	WarningThresholdPct float64 `json:"warningThresholdPct"` // Warn at this % of budget (0-100)
	Enabled            bool    `json:"enabled"`             // Whether budget enforcement is active
}

// DefaultBudgetConfig returns sensible defaults
func DefaultBudgetConfig() BudgetConfig {
	return BudgetConfig{
		MaxCostPerSession:   5.00,
		MaxCostPerDay:       20.00,
		MaxTokensPerSession: 500000,
		MaxTokensPerDay:     2000000,
		WarningThresholdPct: 80.0,
		Enabled:             false, // Off by default for backward compatibility
	}
}

// CostEntry records a single LLM API call cost
type CostEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Provider  string    `json:"provider"`
	Model     string    `json:"model"`
	InputTok  int       `json:"inputTok"`
	OutputTok int       `json:"outputTok"`
	Cost      float64   `json:"cost"`
}

// BudgetTracker tracks and enforces spending limits
type BudgetTracker struct {
	mu          sync.RWMutex
	config      BudgetConfig
	entries     []CostEntry
	sessionStart time.Time
}

// NewBudgetTracker creates a new budget tracker with the given config
func NewBudgetTracker(config BudgetConfig) *BudgetTracker {
	return &BudgetTracker{
		config:       config,
		entries:      make([]CostEntry, 0),
		sessionStart: time.Now(),
	}
}

// RecordCost records an LLM API call and returns an error if budget is exceeded
func (bt *BudgetTracker) RecordCost(provider, model string, inputTok, outputTok int, cost float64) error {
	bt.mu.Lock()

	entry := CostEntry{
		Timestamp: time.Now(),
		Provider:  provider,
		Model:     model,
		InputTok:  inputTok,
		OutputTok: outputTok,
		Cost:      cost,
	}
	bt.entries = append(bt.entries, entry)
	bt.mu.Unlock()

	if !bt.config.Enabled {
		return nil
	}

	bt.mu.Lock()
	defer bt.mu.Unlock()

	// Check session cost limit
	sessionCost := bt.sessionCostLocked()
	if bt.config.MaxCostPerSession > 0 && sessionCost > bt.config.MaxCostPerSession {
		return fmt.Errorf("session budget exceeded: $%.2f / $%.2f", sessionCost, bt.config.MaxCostPerSession)
	}

	// Check session token limit
	sessionTokens := bt.sessionTokensLocked()
	if bt.config.MaxTokensPerSession > 0 && sessionTokens > bt.config.MaxTokensPerSession {
		return fmt.Errorf("session token budget exceeded: %d / %d", sessionTokens, bt.config.MaxTokensPerSession)
	}

	// Check daily cost limit
	dayCost := bt.dailyCostLocked()
	if bt.config.MaxCostPerDay > 0 && dayCost > bt.config.MaxCostPerDay {
		return fmt.Errorf("daily budget exceeded: $%.2f / $%.2f", dayCost, bt.config.MaxCostPerDay)
	}

	// Check daily token limit
	dayTokens := bt.dailyTokensLocked()
	if bt.config.MaxTokensPerDay > 0 && dayTokens > bt.config.MaxTokensPerDay {
		return fmt.Errorf("daily token budget exceeded: %d / %d", dayTokens, bt.config.MaxTokensPerDay)
	}

	return nil
}

// GetStatus returns the current budget status
func (bt *BudgetTracker) BudgetStatus() BudgetStatus {
	bt.mu.RLock()
	defer bt.mu.RUnlock()

	sessionCost := bt.sessionCostLocked()
	sessionTokens := bt.sessionTokensLocked()
	dayCost := bt.dailyCostLocked()
	dayTokens := bt.dailyTokensLocked()

	return BudgetStatus{
		SessionCost:       sessionCost,
		SessionCostLimit:  bt.config.MaxCostPerSession,
		SessionCostPct:    percentage(sessionCost, bt.config.MaxCostPerSession),
		SessionTokens:     sessionTokens,
		SessionTokenLimit: bt.config.MaxTokensPerSession,
		SessionTokenPct:   percentage(float64(sessionTokens), float64(bt.config.MaxTokensPerSession)),
		DailyCost:         dayCost,
		DailyCostLimit:    bt.config.MaxCostPerDay,
		DailyCostPct:      percentage(dayCost, bt.config.MaxCostPerDay),
		DailyTokens:       dayTokens,
		DailyTokenLimit:   bt.config.MaxTokensPerDay,
		DailyTokenPct:     percentage(float64(dayTokens), float64(bt.config.MaxTokensPerDay)),
		TotalEntries:      len(bt.entries),
		SessionStart:      bt.sessionStart,
		WarningThreshold:  bt.config.WarningThresholdPct,
	}
}

// IsNearWarning returns true if any budget is near the warning threshold
func (bt *BudgetTracker) IsNearWarning() bool {
	status := bt.BudgetStatus()
	return status.SessionCostPct >= bt.config.WarningThresholdPct ||
		status.SessionTokenPct >= bt.config.WarningThresholdPct ||
		status.DailyCostPct >= bt.config.WarningThresholdPct ||
		status.DailyTokenPct >= bt.config.WarningThresholdPct
}

// GetEntries returns a copy of all cost entries
func (bt *BudgetTracker) GetEntries() []CostEntry {
	bt.mu.RLock()
	defer bt.mu.RUnlock()
	result := make([]CostEntry, len(bt.entries))
	copy(result, bt.entries)
	return result
}

// ResetSession resets the session start time (but keeps daily tracking)
func (bt *BudgetTracker) ResetSession() {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	bt.sessionStart = time.Now()
}

// SetConfig updates the budget configuration
func (bt *BudgetTracker) SetConfig(config BudgetConfig) {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	bt.config = config
}

// GetConfig returns the current budget configuration
func (bt *BudgetTracker) GetConfig() BudgetConfig {
	bt.mu.RLock()
	defer bt.mu.RUnlock()
	return bt.config
}

func (bt *BudgetTracker) sessionCostLocked() float64 {
	var total float64
	for _, e := range bt.entries {
		if !e.Timestamp.Before(bt.sessionStart) {
			total += e.Cost
		}
	}
	return total
}

func (bt *BudgetTracker) sessionTokensLocked() int {
	var total int
	for _, e := range bt.entries {
		if !e.Timestamp.Before(bt.sessionStart) {
			total += e.InputTok + e.OutputTok
		}
	}
	return total
}

func (bt *BudgetTracker) dailyCostLocked() float64 {
	var total float64
	cutoff := time.Now().Add(-24 * time.Hour)
	for _, e := range bt.entries {
		if e.Timestamp.After(cutoff) {
			total += e.Cost
		}
	}
	return total
}

func (bt *BudgetTracker) dailyTokensLocked() int {
	var total int
	cutoff := time.Now().Add(-24 * time.Hour)
	for _, e := range bt.entries {
		if e.Timestamp.After(cutoff) {
			total += e.InputTok + e.OutputTok
		}
	}
	return total
}

func percentage(value, max float64) float64 {
	if max <= 0 {
		return 0
	}
	pct := (value / max) * 100
	if pct > 100 {
		return 100
	}
	return pct
}

// BudgetStatus represents the current state of all budgets
type BudgetStatus struct {
	SessionCost       float64   `json:"sessionCost"`
	SessionCostLimit  float64   `json:"sessionCostLimit"`
	SessionCostPct    float64   `json:"sessionCostPct"`
	SessionTokens     int       `json:"sessionTokens"`
	SessionTokenLimit int       `json:"sessionTokenLimit"`
	SessionTokenPct   float64   `json:"sessionTokenPct"`
	DailyCost         float64   `json:"dailyCost"`
	DailyCostLimit    float64   `json:"dailyCostLimit"`
	DailyCostPct      float64   `json:"dailyCostPct"`
	DailyTokens       int       `json:"dailyTokens"`
	DailyTokenLimit   int       `json:"dailyTokenLimit"`
	DailyTokenPct     float64   `json:"dailyTokenPct"`
	TotalEntries      int       `json:"totalEntries"`
	SessionStart      time.Time `json:"sessionStart"`
	WarningThreshold  float64   `json:"warningThreshold"`
}
