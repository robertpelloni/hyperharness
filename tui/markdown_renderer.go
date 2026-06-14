package tui

import (
	"github.com/charmbracelet/glamour"
)

// RenderMarkdown converts a markdown string into a terminal-friendly formatted string.
// It uses the default glamour "auto" style which adapts to the terminal's theme.
func RenderMarkdown(md string) (string, error) {
	// Use the auto style (detects light/dark background) and word‑wrap disabled.
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		return "", err
	}
	return r.Render(md)
}
