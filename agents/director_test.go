package agents

import (
	"context"
	"strings"
	"testing"
)

func TestDirectorInitialization(t *testing.T) {
	provider := &DefaultProvider{}
	director := NewDirector(provider)

	if director.GetName() != "Director" {
		t.Errorf("Expected name 'Director', got '%s'", director.GetName())
	}

	if director.GetRole() != "supervisor" {
		t.Errorf("Expected role 'supervisor', got '%s'", director.GetRole())
	}

	if len(director.History) != 1 {
		t.Errorf("Expected initial prompt history length of 1, got %d", len(director.History))
	}

	if !strings.Contains(director.History[0].Content, "Borg TechLead Director") {
		t.Errorf("System prompt missing core identity")
	}
}

func TestDirectorHandleInput(t *testing.T) {
	provider := &DefaultProvider{}
	director := NewDirector(provider)

	resp, err := director.HandleInput(context.Background(), "Run full diagnostic.")
	if err != nil {
		t.Fatalf("HandleInput failed: %v", err)
	}

	if !strings.Contains(resp, "Native Go Borg Director") {
		t.Errorf("Unexpected default response: %s", resp)
	}

	// 1 (sys) + 1 (user) + 1 (assistant) = 3 messages in history
	if len(director.History) != 3 {
		t.Errorf("History not appending correctly, expected 3, got %d", len(director.History))
	}
}
