package tui

// ═══════════════════════════════════════════════════════════════════════
// Theme — full pi-mono theme color system
// 50+ colors matching pi-mono's ThemeJsonSchema exactly
// Plus goose's CRANBERRY/TEAL/GOLD palette for permission dialogs
// Plus opencode's theme system for custom color schemes
// ═══════════════════════════════════════════════════════════════════════

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Theme mirrors pi-mono's ThemeJsonSchema with 50+ color fields
// Go field names avoid collision with methods: DimColor, TextColor, ThinkingCol
type Theme struct {
	// ── Core UI (11 colors) ─────────────────────────────────────
	Accent        string
	Border        string
	BorderAccent  string
	BorderMuted   string
	Success       string
	Error         string
	Warning       string
	Muted         string
	DimColor      string // renamed from "Dim" to avoid method collision
	TextColor     string // renamed from "Text" to avoid method collision
	ThinkingCol   string // renamed from "Thinking" to avoid collision

	// ── Backgrounds & Content Text (11 colors) ──────────────────
	SelectedBg        string
	UserMessageBg     string
	UserMessageText   string
	CustomMessageBg   string
	CustomMessageText string
	CustomMessageLabel string
	ToolPendingBg     string
	ToolSuccessBg     string
	ToolErrorBg       string
	ToolTitle         string
	ToolOutput        string

	// ── Markdown (10 colors) ────────────────────────────────────
	MdHeading         string
	MdLink            string
	MdLinkUrl         string
	MdCode            string
	MdCodeBlock       string
	MdCodeBlockBorder string
	MdQuote           string
	MdQuoteBorder     string
	MdHr              string
	MdListBullet      string

	// ── Tool Diffs (3 colors) ───────────────────────────────────
	ToolDiffAdded   string
	ToolDiffRemoved string
	ToolDiffContext string

	// ── Syntax Highlighting (9 colors) ──────────────────────────
	SyntaxComment    string
	SyntaxKeyword    string
	SyntaxFunction   string
	SyntaxVariable   string
	SyntaxString     string
	SyntaxNumber     string
	SyntaxType       string
	SyntaxOperator   string
	SyntaxPunctuation string

	// ── Thinking Level Borders (6 colors) ───────────────────────
	ThinkingOff     string
	ThinkingMinimal string
	ThinkingLow     string
	ThinkingMedium  string
	ThinkingHigh    string
	ThinkingXhigh   string

	// ── Bash Mode (1 color) ────────────────────────────────────
	BashModeCol string

	// ── Goose accent colors ────────────────────────────────────
	Cranberry string // permission/urgent
	Teal      string // success/active
	Gold      string // pending/selected
}

// DefaultTheme — pi-mono's "default" theme with exact color values
var DefaultTheme = Theme{
	// Core UI
	Accent:       "#3B82F6",
	Border:       "#2E3D54",
	BorderAccent: "#3B82F6",
	BorderMuted:  "#1E293B",
	Success:      "#10B981",
	Error:        "#EF4444",
	Warning:      "#F59E0B",
	Muted:        "#8FA4BD",
	DimColor:     "#5A6D84",
	TextColor:    "#E8E4DF",
	ThinkingCol:  "#8B5CF6",

	// Backgrounds & Content Text
	SelectedBg:        "#1E293B",
	UserMessageBg:     "#1E293B",
	UserMessageText:   "#E8E4DF",
	CustomMessageBg:   "#1C1917",
	CustomMessageText: "#E8E4DF",
	CustomMessageLabel: "#C4883A",
	ToolPendingBg:     "#1E293B",
	ToolSuccessBg:     "#0D3320",
	ToolErrorBg:       "#3B1219",
	ToolTitle:         "#8FA4BD",
	ToolOutput:        "#5A6D84",

	// Markdown
	MdHeading:         "#3B82F6",
	MdLink:            "#3B82F6",
	MdLinkUrl:         "#5A6D84",
	MdCode:            "#F59E0B",
	MdCodeBlock:       "#1E293B",
	MdCodeBlockBorder: "#2E3D54",
	MdQuote:           "#8FA4BD",
	MdQuoteBorder:     "#3B82F6",
	MdHr:              "#2E3D54",
	MdListBullet:      "#3B82F6",

	// Tool Diffs
	ToolDiffAdded:   "#10B981",
	ToolDiffRemoved: "#EF4444",
	ToolDiffContext: "#5A6D84",

	// Syntax Highlighting
	SyntaxComment:     "#5A6D84",
	SyntaxKeyword:     "#C084FC",
	SyntaxFunction:    "#3B82F6",
	SyntaxVariable:    "#E8E4DF",
	SyntaxString:      "#10B981",
	SyntaxNumber:      "#F59E0B",
	SyntaxType:        "#C084FC",
	SyntaxOperator:    "#8FA4BD",
	SyntaxPunctuation: "#5A6D84",

	// Thinking Level Borders
	ThinkingOff:     "#6B7280",
	ThinkingMinimal: "#3B82F6",
	ThinkingLow:     "#8B5CF6",
	ThinkingMedium:  "#F59E0B",
	ThinkingHigh:    "#EF4444",
	ThinkingXhigh:   "#DC2626",

	// Bash Mode
	BashModeCol: "#C4883A",

	// Goose accent colors
	Cranberry: "#C0354A",
	Teal:      "#3A7D7B",
	Gold:      "#C4883A",
}

// DarkTheme — pi-mono's "dark" theme variant
var DarkTheme = Theme{
	Accent:       "#6366F1",
	Border:       "#1E293B",
	BorderAccent: "#6366F1",
	BorderMuted:  "#0F172A",
	Success:      "#22C55E",
	Error:        "#F87171",
	Warning:      "#FBBF24",
	Muted:        "#94A3B8",
	DimColor:     "#475569",
	TextColor:    "#F1F5F9",
	ThinkingCol:  "#A78BFA",

	SelectedBg:        "#0F172A",
	UserMessageBg:     "#0F172A",
	UserMessageText:   "#F1F5F9",
	CustomMessageBg:   "#0C0A09",
	CustomMessageText: "#F1F5F9",
	CustomMessageLabel: "#D97706",
	ToolPendingBg:     "#0F172A",
	ToolSuccessBg:     "#052E16",
	ToolErrorBg:       "#450A0A",
	ToolTitle:         "#94A3B8",
	ToolOutput:        "#475569",

	MdHeading:         "#6366F1",
	MdLink:            "#6366F1",
	MdLinkUrl:         "#475569",
	MdCode:            "#FBBF24",
	MdCodeBlock:       "#0F172A",
	MdCodeBlockBorder: "#1E293B",
	MdQuote:           "#94A3B8",
	MdQuoteBorder:     "#6366F1",
	MdHr:              "#1E293B",
	MdListBullet:      "#6366F1",

	ToolDiffAdded:   "#22C55E",
	ToolDiffRemoved: "#F87171",
	ToolDiffContext: "#475569",

	SyntaxComment:     "#475569",
	SyntaxKeyword:     "#C084FC",
	SyntaxFunction:    "#6366F1",
	SyntaxVariable:    "#F1F5F9",
	SyntaxString:      "#22C55E",
	SyntaxNumber:      "#FBBF24",
	SyntaxType:        "#C084FC",
	SyntaxOperator:    "#94A3B8",
	SyntaxPunctuation: "#475569",

	ThinkingOff:     "#6B7280",
	ThinkingMinimal: "#6366F1",
	ThinkingLow:     "#A78BFA",
	ThinkingMedium:  "#FBBF24",
	ThinkingHigh:    "#F87171",
	ThinkingXhigh:   "#EF4444",

	BashModeCol: "#D97706",
	Cranberry:   "#DC2626",
	Teal:        "#14B8A6",
	Gold:        "#D97706",
}

// AvailableThemes for /theme command (opencode-style)
var AvailableThemes = map[string]Theme{
	"default": DefaultTheme,
	"dark":    DarkTheme,
}

// ═══════════════════════════════════════════════════════════════════════
// Theme helper methods — style generation (mirrors pi-mono's theme.fg/bg)
// ═══════════════════════════════════════════════════════════════════════

func (t Theme) Fg(color, text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(text)
}

func (t Theme) Bg(fgColor, bgColor, text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(fgColor)).Background(lipgloss.Color(bgColor)).Render(text)
}

func (t Theme) Bold(text string) string {
	return lipgloss.NewStyle().Bold(true).Render(text)
}

func (t Theme) Italic(text string) string {
	return lipgloss.NewStyle().Italic(true).Render(text)
}

func (t Theme) Underline(text string) string {
	return lipgloss.NewStyle().Underline(true).Render(text)
}

func (t Theme) Strikethrough(text string) string {
	return lipgloss.NewStyle().Strikethrough(true).Render(text)
}

func (t Theme) Dim(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(t.DimColor)).Render(text)
}

func (t Theme) AccentText(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Bold(true).Render(text)
}

func (t Theme) SuccessText(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(t.Success)).Render(text)
}

func (t Theme) ErrorText(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(t.Error)).Render(text)
}

func (t Theme) WarningText(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(t.Warning)).Render(text)
}

func (t Theme) MutedText(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted)).Render(text)
}

// KeyHint renders a keyboard shortcut hint (pi-mono style: dim parentheses, accent key)
func (t Theme) KeyHint(key, description string) string {
	return t.Dim("(") + t.Fg(t.Accent, key) + t.Dim(") ") + t.Dim(description)
}

// UserMsg renders user message with bg (pi-mono style)
func (t Theme) UserMsg(text string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.UserMessageText)).
		Background(lipgloss.Color(t.UserMessageBg)).
		Padding(0, 1).
		Render(text)
}

// ToolTitleStyled renders tool name with title color (pi-mono style)
func (t Theme) ToolTitleStyled(name string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(t.ToolTitle)).Bold(true).Render(name)
}

// ThinkingBorder returns the border color for a given thinking level
func (t Theme) ThinkingBorder(level string) string {
	switch level {
	case "off":
		return t.ThinkingOff
	case "minimal":
		return t.ThinkingMinimal
	case "low":
		return t.ThinkingLow
	case "medium":
		return t.ThinkingMedium
	case "high":
		return t.ThinkingHigh
	case "xhigh":
		return t.ThinkingXhigh
	default:
		return t.ThinkingOff
	}
}

// PermissionColor returns the color for permission dialog borders (goose-style: gold)
func (t Theme) PermissionColor() string {
	return t.Gold
}

// StatusIcon returns a styled status indicator (goose-style)
func (t Theme) StatusIcon(status string) string {
	ind, ok := ToolStatusIndicators[status]
	if !ok {
		ind = ToolStatusIndicators["pending"]
	}
	return t.Fg(ind.Color, ind.Icon)
}

// ToolKindIcon returns the emoji for a tool kind
func ToolKindIcon(kind string) string {
	if icon, ok := ToolKindIcons[kind]; ok {
		return icon
	}
	return ToolKindIcons["other"]
}

// WrapWidth truncates or wraps text to a given width (utility)
func WrapWidth(text string, width int) string {
	if width <= 0 {
		width = 80
	}
	lines := strings.Split(text, "\n")
	var result []string
	for _, line := range lines {
		for len(line) > width {
			result = append(result, line[:width-1]+"…")
			line = line[width:]
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}
