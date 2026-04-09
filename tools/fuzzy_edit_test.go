package tools

import (
	"testing"
)

func TestFuzzyFindExactMatch(t *testing.T) {
	content := "hello world\nfoo bar\nbaz"
	result := FuzzyFindText(content, "foo bar")
	if !result.Found {
		t.Fatal("should find exact match")
	}
	if result.UsedFuzzyMatch {
		t.Error("exact match should not use fuzzy")
	}
}

func TestFuzzyFindTrailingWhitespace(t *testing.T) {
	// oldText has NO trailing spaces, but content lines DO.
	// Exact match still succeeds because "foo bar" is a substring of "foo bar  "
	content := "hello world  \nfoo bar  \nbaz"
	result := FuzzyFindText(content, "foo bar")
	if !result.Found {
		t.Fatal("should find 'foo bar' in content")
	}
	// Exact match works because it's a substring
	if result.UsedFuzzyMatch {
		t.Error("should be exact match since it's a substring")
	}
}

func TestFuzzyFindSmartQuotes(t *testing.T) {
	content := "it's a test\nwith \u201Csmart quotes\u201D"
	result := FuzzyFindText(content, "it's a test")
	if !result.Found {
		t.Fatal("should find with smart quote normalization")
	}
}

func TestFuzzyFindDashes(t *testing.T) {
	content := "hello\u2013world"
	result := FuzzyFindText(content, "hello-world")
	if !result.Found {
		t.Fatal("should find with dash normalization")
	}
}

func TestFuzzyFindNotFound(t *testing.T) {
	content := "hello world"
	result := FuzzyFindText(content, "goodbye world")
	if result.Found {
		t.Error("should not find non-matching text")
	}
}

func TestNormalizeForFuzzyMatch(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\u201Chello\u201D", "\"hello\""},
		{"it\u2019s", "it's"},
		{"a\u2013b", "a-b"},
		{"a\u2014b", "a-b"},
		{"hello\u00A0world", "hello world"},
		{"line1  \nline2  ", "line1\nline2"},
	}
	for _, tc := range tests {
		got := normalizeForFuzzyMatch(tc.input)
		if got != tc.expected {
			t.Errorf("normalizeForFuzzyMatch(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestDetectLineEnding(t *testing.T) {
	if ending := DetectLineEnding("hello\nworld"); ending != "\n" {
		t.Errorf("LF: got %q", ending)
	}
	if ending := DetectLineEnding("hello\r\nworld"); ending != "\r\n" {
		t.Errorf("CRLF: got %q", ending)
	}
	if ending := DetectLineEnding("no newlines"); ending != "\n" {
		t.Errorf("no newlines: got %q", ending)
	}
}

func TestNormalizeToLF(t *testing.T) {
	input := "hello\r\nworld\r\n"
	got := NormalizeToLF(input)
	if got != "hello\nworld\n" {
		t.Errorf("got %q", got)
	}
}

func TestRestoreLineEndings(t *testing.T) {
	input := "hello\nworld\n"
	got := RestoreLineEndings(input, "\r\n")
	if got != "hello\r\nworld\r\n" {
		t.Errorf("got %q", got)
	}
	got = RestoreLineEndings(input, "\n")
	if got != input {
		t.Errorf("got %q", got)
	}
}

func TestStripBOM(t *testing.T) {
	bom, text := StripBOM("\xEF\xBB\xBFhello")
	if bom != "\xEF\xBB\xBF" {
		t.Errorf("BOM: %q", bom)
	}
	if text != "hello" {
		t.Errorf("text: %q", text)
	}
	bom2, text2 := StripBOM("hello")
	if bom2 != "" {
		t.Errorf("no BOM should return empty")
	}
	if text2 != "hello" {
		t.Errorf("text: %q", text2)
	}
}

func TestApplyEdits(t *testing.T) {
	content := "line1\nline2\nline3\nline4"
	edits := []Edit{
		{OldText: "line2", NewText: "LINE_TWO"},
		{OldText: "line4", NewText: "LINE_FOUR"},
	}
	base, newContent, err := ApplyEdits(content, edits)
	if err != nil {
		t.Fatalf("ApplyEdits error: %v", err)
	}
	if base != content {
		t.Error("base should be original content")
	}
	if newContent != "line1\nLINE_TWO\nline3\nLINE_FOUR" {
		t.Errorf("result: %q", newContent)
	}
}

func TestApplyEditsNotFound(t *testing.T) {
	content := "hello world"
	edits := []Edit{
		{OldText: "not found", NewText: "replacement"},
	}
	_, _, err := ApplyEdits(content, edits)
	if err == nil {
		t.Error("should error for not found oldText")
	}
}

func TestApplyEditsOverlapping(t *testing.T) {
	content := "hello world foo bar"
	edits := []Edit{
		{OldText: "hello world", NewText: "HELLO"},
		{OldText: "world foo", NewText: "WORLD"},
	}
	_, _, err := ApplyEdits(content, edits)
	if err == nil {
		t.Error("should error for overlapping edits")
	}
}
