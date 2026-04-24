package skills

import (
	"strings"
	"testing"
)

func TestManagerRegister(t *testing.T) {
	m := NewManager()
	skill := &Skill{
		Name:        "test-skill",
		Description: "A test skill",
		Version:     "1.0.0",
		Category:    "test",
	}
	m.Register(skill)

	found, ok := m.Get("test-skill")
	if !ok {
		t.Fatal("skill should be found")
	}
	if found.Name != "test-skill" {
		t.Errorf("skill name: got %s", found.Name)
	}
}

func TestManagerList(t *testing.T) {
	m := NewManager()
	m.Register(&Skill{Name: "b-skill", Category: "test"})
	m.Register(&Skill{Name: "a-skill", Category: "test"})
	m.Register(&Skill{Name: "c-skill", Category: "test"})

	skills := m.List()
	if len(skills) != 3 {
		t.Fatalf("expected 3 skills, got %d", len(skills))
	}
	// Should be sorted
	if skills[0].Name != "a-skill" {
		t.Errorf("expected sorted order, first: %s", skills[0].Name)
	}
}

func TestManagerExecute(t *testing.T) {
	m := NewManager()
	m.Register(&Skill{
		Name:    "test-skill",
		Content: "Hello {{name}}, welcome to {{project}}!",
	})

	result, err := m.Execute("test-skill", map[string]interface{}{
		"name":    "World",
		"project": "HyperHarness",
	})
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if !strings.Contains(result.Output, "Hello World") {
		t.Errorf("output missing substitution: %s", result.Output)
	}
	if !strings.Contains(result.Output, "HyperHarness") {
		t.Errorf("output missing project: %s", result.Output)
	}
}

func TestManagerExecuteNotFound(t *testing.T) {
	m := NewManager()
	_, err := m.Execute("nonexistent", nil)
	if err == nil {
		t.Error("expected error for nonexistent skill")
	}
}

func TestManagerMatchByTrigger(t *testing.T) {
	m := NewManager()
	m.Register(&Skill{
		Name:         "memory-skill",
		TriggerWords: []string{"memory", "remember", "recall"},
	})
	m.Register(&Skill{
		Name:         "browser-skill",
		TriggerWords: []string{"browser", "web", "scrape"},
	})

	matches := m.MatchByTrigger("I want to use memory to store something")
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Name != "memory-skill" {
		t.Errorf("wrong match: %s", matches[0].Name)
	}

	matches = m.MatchByTrigger("open the browser and navigate")
	if len(matches) != 1 {
		t.Fatalf("expected 1 match for browser, got %d", len(matches))
	}
}

func TestBuiltins(t *testing.T) {
	builtins := Builtins()
	if len(builtins) < 4 {
		t.Errorf("expected at least 4 builtins, got %d", len(builtins))
	}

	names := make(map[string]bool)
	for _, skill := range builtins {
		if skill.Name == "" {
			t.Error("builtin skill has empty name")
		}
		if skill.Content == "" {
			t.Errorf("builtin %s has empty content", skill.Name)
		}
		names[skill.Name] = true
	}

	for _, expected := range []string{"memory", "context", "agent-browser", "simplify"} {
		if !names[expected] {
			t.Errorf("missing builtin: %s", expected)
		}
	}
}

func TestManagerListByCategory(t *testing.T) {
	m := NewManager()
	m.Register(&Skill{Name: "s1", Category: "core"})
	m.Register(&Skill{Name: "s2", Category: "core"})
	m.Register(&Skill{Name: "s3", Category: "automation"})

	byCat := m.ListByCategory()
	if len(byCat["core"]) != 2 {
		t.Errorf("expected 2 core skills, got %d", len(byCat["core"]))
	}
	if len(byCat["automation"]) != 1 {
		t.Errorf("expected 1 automation skill, got %d", len(byCat["automation"]))
	}
}
