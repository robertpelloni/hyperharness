package tools

import (
	"testing"
)

func TestHTMLToMarkdownHeadings(t *testing.T) {
	input := "<h1>Title</h1><h2>Subtitle</h2>"
	result := HTMLToMarkdown(input)
	if !containsStr(result, "# Title") {
		t.Errorf("missing h1: %s", result)
	}
	if !containsStr(result, "## Subtitle") {
		t.Errorf("missing h2: %s", result)
	}
}

func TestHTMLToMarkdownLinks(t *testing.T) {
	input := `<a href="https://example.com">Click here</a>`
	result := HTMLToMarkdown(input)
	if !containsStr(result, "[Click here](https://example.com)") {
		t.Errorf("link: %s", result)
	}
}

func TestHTMLToMarkdownCode(t *testing.T) {
	input := "<code>var x = 1</code>"
	result := HTMLToMarkdown(input)
	if !containsStr(result, "`var x = 1`") {
		t.Errorf("inline code: %s", result)
	}
}

func TestHTMLToMarkdownCodeBlock(t *testing.T) {
	input := "<pre><code>func main() {\n  fmt.Println(\"hi\")\n}</code></pre>"
	result := HTMLToMarkdown(input)
	if !containsStr(result, "```") {
		t.Errorf("code block: %s", result)
	}
}

func TestHTMLToMarkdownLists(t *testing.T) {
	input := "<ul><li>item1</li><li>item2</li></ul>"
	result := HTMLToMarkdown(input)
	if !containsStr(result, "- item1") || !containsStr(result, "- item2") {
		t.Errorf("unordered list: %s", result)
	}
}

func TestHTMLToMarkdownOrderedList(t *testing.T) {
	input := "<ol><li>first</li><li>second</li></ol>"
	result := HTMLToMarkdown(input)
	if !containsStr(result, "1. first") || !containsStr(result, "2. second") {
		t.Errorf("ordered list: %s", result)
	}
}

func TestHTMLToMarkdownEmphasis(t *testing.T) {
	input := "<strong>bold</strong> and <em>italic</em>"
	result := HTMLToMarkdown(input)
	if !containsStr(result, "**bold**") {
		t.Errorf("bold: %s", result)
	}
	if !containsStr(result, "*italic*") {
		t.Errorf("italic: %s", result)
	}
}

func TestHTMLToMarkdownStripsScript(t *testing.T) {
	input := "<p>Hello</p><script>alert('xss')</script><p>World</p>"
	result := HTMLToMarkdown(input)
	if containsStr(result, "alert") {
		t.Errorf("should strip script: %s", result)
	}
}

func TestExtractTextFromHTML(t *testing.T) {
	input := "<html><body><h1>Title</h1><p>Content</p><script>skip</script></body></html>"
	result := ExtractTextFromHTML(input)
	if containsStr(result, "<") {
		t.Errorf("should strip all tags: %s", result)
	}
	if containsStr(result, "skip") {
		t.Errorf("should strip script content: %s", result)
	}
	if !containsStr(result, "Title") || !containsStr(result, "Content") {
		t.Errorf("missing text content: %s", result)
	}
}

func TestDecodeEntities(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a &amp; b", "a & b"},
		{"&lt;tag&gt;", "<tag>"},
		{"a &quot;b&quot;", `a "b"`},
		{"x&#39;y", "x'y"},
		{"a&nbsp;b", "a b"},
	}
	for _, tc := range tests {
		got := decodeEntities(tc.input)
		if got != tc.expected {
			t.Errorf("decodeEntities(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(len(s) > 0 && findSubstring(s, sub)))
}

func findSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
