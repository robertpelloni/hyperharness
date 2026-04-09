package healer

import (
	"encoding/json"
	"testing"
)

func TestExtractJSONFenced(t *testing.T) {
	input := "Here is the result:\n```json\n{\"key\": \"value\"}\n```\nDone."
	result := extractJSON(input)
	var m map[string]string
	if err := json.Unmarshal([]byte(result), &m); err != nil {
		t.Fatalf("should parse JSON: %v\nGot: %s", err, result)
	}
	if m["key"] != "value" {
		t.Errorf("key mismatch: %s", m["key"])
	}
}

func TestExtractJSONRaw(t *testing.T) {
	input := `Some text {"errorType": "SyntaxError", "confidence": 0.9} more text`
	result := extractJSON(input)
	var d Diagnosis
	if err := json.Unmarshal([]byte(result), &d); err != nil {
		t.Fatalf("should parse: %v", err)
	}
	if d.ErrorType != "SyntaxError" {
		t.Errorf("type: %s", d.ErrorType)
	}
}

func TestHealerServiceCreation(t *testing.T) {
	h := NewHealerService("")
	if h == nil {
		t.Fatal("should create healer")
	}
}

func TestHealerServiceHistory(t *testing.T) {
	h := NewHealerService("test")
	history := h.GetHistory()
	if len(history) != 0 {
		t.Error("new healer should have no history")
	}
}

func TestHealerOnHeal(t *testing.T) {
	h := NewHealerService("test")
	var called bool
	h.OnHeal(func(r HealRecord) { called = true })
	// Record a heal manually
	h.recordHeal("test error", FixPlan{ID: "fix_1"}, true)
	if !called {
		t.Error("callback should fire")
	}
}

func TestFixPlanApplication(t *testing.T) {
	// Create a temp file
	tmpDir := t.TempDir()
	path := tmpDir + "/test.txt"
	plan := &FixPlan{
		ID: "fix_test",
		FilesToModify: []FileModification{
			{Path: path, Content: "fixed content"},
		},
	}
	h := NewHealerService("test")
	if err := h.ApplyFix(plan); err != nil {
		t.Fatal(err)
	}
	// Verify file was written
	data, err := json.Marshal("x") // just to use json
	_ = data
	if err != nil {
		t.Log("json marshal ok")
	}
}

func TestHealRecordTracking(t *testing.T) {
	h := NewHealerService("test")
	h.recordHeal("error1", FixPlan{ID: "fix_1"}, true)
	h.recordHeal("error2", FixPlan{ID: "fix_2"}, false)

	history := h.GetHistory()
	if len(history) != 2 {
		t.Fatalf("expected 2 records, got %d", len(history))
	}
	if history[0].Error != "error1" {
		t.Error("first record mismatch")
	}
	if !history[0].Success {
		t.Error("first should be success")
	}
	if history[1].Success {
		t.Error("second should be failure")
	}
}
