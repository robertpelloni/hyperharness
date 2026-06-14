package debate

import (
	"fmt"
	"strings"
)

// Arguer defines the interface for agents that can participate in a debate.
type Arguer interface {
	Argue(topic, stance, previousArgument string) (string, error)
}

// DebateResult contains the outcome of a structured debate between subagents.
type DebateResult struct {
	Topic       string
	ProArgument string
	ConArgument string
	Consensus   string
}

// Run executes a debate on the given topic between two Arguers.
// pro argues as the proponent first; con responds as the skeptic.
// It returns the arguments from both sides and a simple consensus determination.
func Run(topic string, pro, con Arguer) (*DebateResult, error) {
	// Proponent argues first
	proArg, err := pro.Argue(topic, "Proponent", "")
	if err != nil {
		return nil, fmt.Errorf("debate: proponent argument failed: %w", err)
	}

	// Skeptic responds to the proponent's argument
	conArg, err := con.Argue(topic, "Skeptic", proArg)
	if err != nil {
		return nil, fmt.Errorf("debate: skeptic argument failed: %w", err)
	}

	// Simple consensus heuristic: if the skeptic's argument contains
	// agreement language or references the proponent's points positively,
	// we consider it a consensus.
	consensus := "No clear consensus"
	conLower := strings.ToLower(conArg)
	if strings.Contains(conLower, "agree") ||
		strings.Contains(conLower, "valid") ||
		strings.Contains(conLower, "acceptable") ||
		strings.Contains(conLower, "optimal") {
		consensus = "Agreement reached"
	}

	return &DebateResult{
		Topic:       topic,
		ProArgument: proArg,
		ConArgument: conArg,
		Consensus:   consensus,
	}, nil
}

// FormatResult returns a human-readable summary of the debate.
func FormatResult(r *DebateResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🔥 Debate on: %s\n", r.Topic))
	sb.WriteString(fmt.Sprintf("\n📢 Proponent (Coder):\n%s\n", r.ProArgument))
	sb.WriteString(fmt.Sprintf("\n🛡️ Skeptic (Reviewer):\n%s\n", r.ConArgument))
	sb.WriteString(fmt.Sprintf("\n✅ Consensus: %s\n", r.Consensus))
	return sb.String()
}