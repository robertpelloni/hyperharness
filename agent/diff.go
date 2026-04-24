package agent

import (
	"fmt"
	"os"
	"strings"
)

// ApplyInlineDiff mimics Opencode / Cursor style inline file modifications.
// It parses a unified diff format and applies the changes directly to the file.
func (a *Agent) ApplyInlineDiff(filePath, diffContent string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file for diff application: %w", err)
	}

	strContent := string(content)

	// Extremely simplified diff application for parity demonstration.
	// In a full implementation, we'd use a robust diff matching algorithm.
	lines := strings.Split(diffContent, "\n")
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "---"), strings.HasPrefix(line, "+++"), strings.HasPrefix(line, "@@"):
			continue
		case strings.HasPrefix(line, "-"):
			toRemove := strings.TrimPrefix(line, "-")
			if strings.Contains(strContent, toRemove+"\n") {
				strContent = strings.Replace(strContent, toRemove+"\n", "", 1)
			} else {
				strContent = strings.Replace(strContent, toRemove, "", 1)
			}
		case strings.HasPrefix(line, "+"):
			toAdd := strings.TrimPrefix(line, "+")
			// This is still intentionally naive; we only skip unified-diff headers
			// and append added lines rather than interpreting hunk positions.
			strContent += toAdd + "\n"
		}
	}

	err = os.WriteFile(filePath, []byte(strContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write applied diff: %w", err)
	}

	fmt.Printf("[InlineDiff] Successfully applied diff to %s\n", filePath)
	return nil
}
