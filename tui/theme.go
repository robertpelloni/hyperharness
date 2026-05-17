package tui

// ═══════════════════════════════════════════════════════════════════════
// Theme — mirrors pi-mono's ThemeColor system
// ═══════════════════════════════════════════════════════════════════════

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Accent       string
	Border       string
	BorderAccent string
	Success      string
	Error        string
	Warning      string
	Muted        string
	DimColor     string
	TextColor    string
	ThinkingCol  string

	UserMsgBg    string
	UserMsgText  string
	ToolPending  string
	ToolSuccess  string
	ToolError    string
	ToolTitle    string
	ToolOutput   string

	MDHeading    string
	MDLink       string
	MDCode       string
	MDCodeBlock  string
	MDQuote      string
	MDListBullet string

	DiffAdded   string
	DiffRemoved string
	DiffContext string

	BashModeCol string
}

var DefaultTheme = Theme{
	Accent:       "#7C3AED",
	Border:       "#374151",
	BorderAccent: "#9333EA",
	Success:      "#10B981",
	Error:        "#EF4444",
	Warning:      "#F59E0B",
	Muted:        "#9CA3AF",
	DimColor:     "#6B7280",
	TextColor:    "#F9FAFB",
	ThinkingCol:  "#8B5CF6",

	UserMsgBg:    "#1E293B",
	UserMsgText:  "#93C5FD",
	ToolPending:  "#1E1B4B",
	ToolSuccess:  "#022C22",
	ToolError:    "#450A0A",
	ToolTitle:    "#F59E0B",
	ToolOutput:   "#9CA3AF",

	MDHeading:    "#C4B5FD",
	MDLink:       "#60A5FA",
	MDCode:       "#FCD34D",
	MDCodeBlock:  "#1F2937",
	MDQuote:      "#6B7280",
	MDListBullet: "#9333EA",

	DiffAdded:   "#22C55E",
	DiffRemoved: "#EF4444",
	DiffContext: "#6B7280",

	BashModeCol: "#F97316",
}

func (t Theme) Fg(color string, text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(text)
}

func (t Theme) Bg(bgColor, text string) string {
	return lipgloss.NewStyle().Background(lipgloss.Color(bgColor)).Render(text)
}

func (t Theme) Bold(text string) string {
	return lipgloss.NewStyle().Bold(true).Render(text)
}

func (t Theme) Italic(text string) string {
	return lipgloss.NewStyle().Italic(true).Render(text)
}

func (t Theme) Dim(text string) string {
	return t.Fg(t.DimColor, text)
}

func (t Theme) AccentText(text string) string {
	return t.Fg(t.Accent, text)
}

func (t Theme) SuccessText(text string) string {
	return t.Fg(t.Success, text)
}

func (t Theme) ErrorText(text string) string {
	return t.Fg(t.Error, text)
}

func (t Theme) WarningText(text string) string {
	return t.Fg(t.Warning, text)
}

func (t Theme) MutedText(text string) string {
	return t.Fg(t.Muted, text)
}

// KeyHint renders a keyboard shortcut hint (pi-mono style)
func (t Theme) KeyHint(key, desc string) string {
	return t.Dim(key) + t.MutedText(" " + desc)
}

// RawKeyHint renders a key hint from a raw string
func (t Theme) RawKeyHint(key, desc string) string {
	return t.Dim(key) + t.MutedText(" " + desc)
}
