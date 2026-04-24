package tools

import (
	"testing"
)

func TestGetStr(t *testing.T) {
	m := map[string]interface{}{
		"name":   "Alice",
		"count":  float64(42),
		"active": true,
	}

	if got := GetStr(m, "name"); got != "Alice" {
		t.Errorf("GetStr(name) = %q, want %q", got, "Alice")
	}
	if got := GetStr(m, "missing"); got != "" {
		t.Errorf("GetStr(missing) = %q, want empty", got)
	}
	if got := GetStr(m, "count"); got != "42" {
		t.Errorf("GetStr(count) = %q, want %q", got, "42")
	}
}

func TestGetInt(t *testing.T) {
	m := map[string]interface{}{
		"count":  float64(42),
		"neg":    float64(-5),
		"zero":   float64(0),
		"string": "not_a_number",
	}

	if got := GetInt(m, "count"); got != 42 {
		t.Errorf("GetInt(count) = %d, want 42", got)
	}
	if got := GetInt(m, "neg"); got != -5 {
		t.Errorf("GetInt(neg) = %d, want -5", got)
	}
	if got := GetInt(m, "zero"); got != 0 {
		t.Errorf("GetInt(zero) = %d, want 0", got)
	}
	if got := GetInt(m, "missing"); got != 0 {
		t.Errorf("GetInt(missing) = %d, want 0", got)
	}
}

func TestGetBool(t *testing.T) {
	m := map[string]interface{}{
		"active":  true,
		"enabled": false,
		"count":   float64(1),
	}

	if got := GetBool(m, "active"); !got {
		t.Error("GetBool(active) should be true")
	}
	if got := GetBool(m, "enabled"); got {
		t.Error("GetBool(enabled) should be false")
	}
	if got := GetBool(m, "missing"); got {
		t.Error("GetBool(missing) should be false")
	}
}

func TestGetStringSlice(t *testing.T) {
	m := map[string]interface{}{
		"tags": []interface{}{"go", "testing", "tools"},
		"name": "single",
	}

	result := GetStringSlice(m, "tags")
	if len(result) != 3 {
		t.Fatalf("GetStringSlice(tags) = %d items, want 3", len(result))
	}
	if result[0] != "go" || result[2] != "tools" {
		t.Errorf("GetStringSlice(tags) = %v", result)
	}

	result = GetStringSlice(m, "missing")
	if result != nil {
		t.Errorf("GetStringSlice(missing) = %v, want nil", result)
	}
}

func TestGetStringMap(t *testing.T) {
	m := map[string]interface{}{
		"env": map[string]interface{}{
			"HOME": "/home/user",
			"PATH": "/usr/bin",
		},
	}

	result := GetStringMap(m, "env")
	if result["HOME"] != "/home/user" {
		t.Errorf("GetStringMap(env)[HOME] = %q", result["HOME"])
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello world", 8, "hello..."},
		{"short", 10, "short"},
		{"", 5, ""},
		{"abc", 3, "abc"},
		{"abcd", 3, "abc"},
	}

	for _, tc := range tests {
		got := TruncateString(tc.input, tc.maxLen)
		if got != tc.want {
			t.Errorf("TruncateString(%q, %d) = %q, want %q", tc.input, tc.maxLen, got, tc.want)
		}
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{1500, "1.5 KB"},
		{1500000, "1.4 MB"},
		{1500000000, "1.4 GB"},
	}

	for _, tc := range tests {
		got := FormatFileSize(tc.bytes)
		if got != tc.want {
			t.Errorf("FormatFileSize(%d) = %q, want %q", tc.bytes, got, tc.want)
		}
	}
}

func TestIndentLines(t *testing.T) {
	input := "line1\nline2\nline3"
	want := "  line1\n  line2\n  line3"
	got := IndentLines(input, "  ")
	if got != want {
		t.Errorf("IndentLines = %q, want %q", got, want)
	}
}
