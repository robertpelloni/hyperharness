// Package orchestrator provides council-based multi-model orchestration.
// Ported from hypercode/go/internal/orchestration/council.go
//
// WHAT: Council of LLMs that debate, vote, and reach consensus on decisions
// WHY: Multi-model deliberation produces better results than single-model inference
// HOW: Each council member proposes a response, a chairman synthesizes the best answer
package orchestrator

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// CouncilRole defines a member's role in the council.
type CouncilRole string

const (
	CouncilChairman   CouncilRole = "chairman"
	CouncilDelegate   CouncilRole = "member"
	CouncilCritic     CouncilRole = "critic"
	CouncilSpecialist CouncilRole = "specialist"
)

// CouncilMember represents a member of the council.
type CouncilMember struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Role     CouncilRole `json:"role"`
	Provider string      `json:"provider"`
	ModelID  string      `json:"modelId"`
}

// CouncilVote represents a vote from a council member.
type CouncilVote struct {
	MemberID  string  `json:"memberId"`
	Proposal  string  `json:"proposal"`
	Confidence float64 `json:"confidence"`
	Rationale string  `json:"rationale"`
}

// CouncilSession represents a deliberation session.
type CouncilSession struct {
	ID         string          `json:"id"`
	Topic      string          `json:"topic"`
	Members    []CouncilMember `json:"members"`
	Votes      []CouncilVote   `json:"votes"`
	Consensus  string          `json:"consensus,omitempty"`
	Status     string          `json:"status"` // "deliberating", "consensus", "failed"
	Turns      int             `json:"turns"`
	StartedAt  time.Time       `json:"startedAt"`
	FinishedAt *time.Time      `json:"finishedAt,omitempty"`
}

// Council manages multi-model deliberation.
type Council struct {
	members  map[string]CouncilMember
	sessions map[string]*CouncilSession
	mu       sync.RWMutex
}

// NewCouncil creates a new council.
func NewCouncil() *Council {
	return &Council{
		members:  make(map[string]CouncilMember),
		sessions: make(map[string]*CouncilSession),
	}
}

// AddMember adds a member to the council.
func (c *Council) AddMember(member CouncilMember) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.members[member.ID] = member
}

// RemoveMember removes a member from the council.
func (c *Council) RemoveMember(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.members, id)
}

// MemberCount returns the number of council members.
func (c *Council) MemberCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.members)
}

// ListMembers returns all council members.
func (c *Council) ListMembers() []CouncilMember {
	c.mu.RLock()
	defer c.mu.RUnlock()
	members := make([]CouncilMember, 0, len(c.members))
	for _, m := range c.members {
		members = append(members, m)
	}
	return members
}

// Deliberate starts a council session on a topic.
// Each member proposes a response, then the chairman synthesizes consensus.
// This is a simplified version that works without actual LLM calls.
func (c *Council) Deliberate(ctx context.Context, topic string, maxTurns int) (*CouncilSession, error) {
	c.mu.Lock()
	session := &CouncilSession{
		ID:        fmt.Sprintf("council-%d-%d", time.Now().UnixNano(), len(c.sessions)),
		Topic:     topic,
		Members:   c.listMembersLocked(),
		Status:    "deliberating",
		StartedAt: time.Now().UTC(),
	}
	c.sessions[session.ID] = session
	c.mu.Unlock()

	if maxTurns <= 0 {
		maxTurns = 3
	}

	for turn := 0; turn < maxTurns; turn++ {
		session.Turns++

		// Collect proposals from all members
		votes := c.collectProposals(ctx, session, topic)
		session.Votes = append(session.Votes, votes...)

		// Check for consensus
		if c.hasConsensus(votes) {
			session.Consensus = c.synthesizeConsensus(votes)
			session.Status = "consensus"
			break
		}
	}

	if session.Status != "consensus" {
		// Force consensus from best available
		if len(session.Votes) > 0 {
			session.Consensus = c.synthesizeConsensus(session.Votes)
			session.Status = "consensus"
		} else {
			session.Status = "failed"
		}
	}

	fin := time.Now().UTC()
	session.FinishedAt = &fin
	return session, nil
}

// collectProposals simulates collecting proposals from council members.
// In production, each member would call their assigned LLM.
func (c *Council) collectProposals(ctx context.Context, session *CouncilSession, topic string) []CouncilVote {
	c.mu.RLock()
	members := make([]CouncilMember, 0, len(c.members))
	for _, m := range c.members {
		members = append(members, m)
	}
	c.mu.RUnlock()

	var votes []CouncilVote
	for _, m := range members {
		select {
		case <-ctx.Done():
			return votes
		default:
		}

		// Simulate proposal with confidence
		confidence := 0.5 + rand.Float64()*0.5
		votes = append(votes, CouncilVote{
			MemberID:   m.ID,
			Proposal:   fmt.Sprintf("[Proposal from %s (%s)] Analyzing: %s", m.Name, m.Role, topic),
			Confidence: confidence,
			Rationale:  fmt.Sprintf("Based on %s analysis with %.0f%% confidence", m.Role, confidence*100),
		})
	}

	return votes
}

// hasConsensus checks if votes indicate consensus (all above threshold).
func (c *Council) hasConsensus(votes []CouncilVote) bool {
	if len(votes) < 2 {
		return false
	}
	for _, v := range votes {
		if v.Confidence < 0.7 {
			return false
		}
	}
	return true
}

// synthesizeConsensus picks the best proposal based on confidence.
func (c *Council) synthesizeConsensus(votes []CouncilVote) string {
	if len(votes) == 0 {
		return "No proposals received"
	}

	best := votes[0]
	for _, v := range votes[1:] {
		if v.Confidence > best.Confidence {
			best = v
		}
	}

	return fmt.Sprintf("Council Consensus: %s (confidence: %.0f%%)\n%s",
		best.Proposal, best.Confidence*100, best.Rationale)
}

// GetSession retrieves a session by ID.
func (c *Council) GetSession(id string) (*CouncilSession, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	session, ok := c.sessions[id]
	return session, ok
}

// ListSessions returns all council sessions.
func (c *Council) ListSessions() []*CouncilSession {
	c.mu.RLock()
	defer c.mu.RUnlock()
	sessions := make([]*CouncilSession, 0, len(c.sessions))
	for _, s := range c.sessions {
		sessions = append(sessions, s)
	}
	return sessions
}

// listMembersLocked returns members (must hold lock).
func (c *Council) listMembersLocked() []CouncilMember {
	members := make([]CouncilMember, 0, len(c.members))
	for _, m := range c.members {
		members = append(members, m)
	}
	return members
}
