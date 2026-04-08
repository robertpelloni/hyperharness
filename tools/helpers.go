// Package tools provides shared utility functions for tool implementations.
// Consolidates common helpers that were previously duplicated across parity files.
package tools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetStr extracts a string value from a map[string]interface{}.
// Returns empty string if key is missing or value is not a string.
func GetStr(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case json.Number:
		return string(val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetInt extracts an int value from a map[string]interface{}.
// Supports float64 (JSON default), int, and json.Number.
func GetInt(m map[string]interface{}, key string) int {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case float32:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	default:
		return 0
	}
}

// GetBool extracts a bool value from a map[string]interface{}.
func GetBool(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	b, ok := v.(bool)
	if ok {
		return b
	}
	return false
}

// GetStringSlice extracts a []string from a map[string]interface{}.
// Handles both []string and []interface{} containing strings.
func GetStringSlice(m map[string]interface{}, key string) []string {
	v, ok := m[key]
	if !ok {
		return nil
	}

	switch slice := v.(type) {
	case []string:
		return slice
	case []interface{}:
		result := make([]string, 0, len(slice))
		for _, item := range slice {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	default:
		return nil
	}
}

// GetStringMap extracts a map[string]string from a map[string]interface{}.
func GetStringMap(m map[string]interface{}, key string) map[string]string {
	v, ok := m[key]
	if !ok {
		return nil
	}

	switch m := v.(type) {
	case map[string]string:
		return m
	case map[string]interface{}:
		result := make(map[string]string, len(m))
		for k, v := range m {
			result[k] = fmt.Sprintf("%v", v)
		}
		return result
	default:
		return nil
	}
}

// TruncateString truncates a string to maxLen characters with ellipsis.
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 4 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// FormatFileSize formats bytes into human-readable size.
func FormatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// IndentLines indents all lines in a string by the given prefix.
func IndentLines(s string, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}
