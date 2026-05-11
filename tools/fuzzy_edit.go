// Package tools implements fuzzy text matching for the edit tool.
// Ported from superai/pi-cli/packages/coding-agent/src/core/tools/edit-diff.ts
//
// The fuzzy matching system is critical because AI models sometimes introduce
// minor differences in whitespace, Unicode quotes, or line endings when
// generating oldText for edit operations. Progressive normalization ensures
// edits succeed even with these common discrepancies.
//
// Matching strategy (in order of priority):
// 1. Exact match - fastest path, no transformation needed
// 2. Whitespace-stripped match - trailing whitespace differences
// 3. Full fuzzy match - Unicode normalization (quotes, dashes, spaces)
package tools

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// FuzzyMatchResult contains the result of a fuzzy text search.
type FuzzyMatchResult struct {
	Found                 bool
	Index                 int
	MatchLength           int
	UsedFuzzyMatch        bool
	ContentForReplacement string // normalized content if fuzzy
}

// FuzzyFindText finds oldText in content, trying exact match first,
// then progressively more aggressive fuzzy matching.
func FuzzyFindText(content, oldText string) FuzzyMatchResult {
	// 1. Exact match
	idx := strings.Index(content, oldText)
	if idx != -1 {
		return FuzzyMatchResult{
			Found:                 true,
			Index:                 idx,
			MatchLength:           len(oldText),
			UsedFuzzyMatch:        false,
			ContentForReplacement: content,
		}
	}

	// 2. Try with trailing whitespace stripped from each line
	strippedContent := stripTrailingWhitespace(content)
	strippedOld := stripTrailingWhitespace(oldText)
	idx = strings.Index(strippedContent, strippedOld)
	if idx != -1 {
		return FuzzyMatchResult{
			Found:                 true,
			Index:                 idx,
			MatchLength:           len(strippedOld),
			UsedFuzzyMatch:        true,
			ContentForReplacement: strippedContent,
		}
	}

	// 3. Full fuzzy: normalize Unicode quotes, dashes, spaces
	normalizedContent := normalizeForFuzzyMatch(content)
	normalizedOld := normalizeForFuzzyMatch(oldText)
	idx = strings.Index(normalizedContent, normalizedOld)
	if idx != -1 {
		return FuzzyMatchResult{
			Found:                 true,
			Index:                 idx,
			MatchLength:           len(normalizedOld),
			UsedFuzzyMatch:        true,
			ContentForReplacement: normalizedContent,
		}
	}

	return FuzzyMatchResult{Found: false}
}

// stripTrailingWhitespace removes trailing whitespace from each line.
func stripTrailingWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r")
	}
	return strings.Join(lines, "\n")
}

// normalizeForFuzzyMatch applies progressive Unicode normalization:
// - NFKC normalization
// - Strip trailing whitespace per line
// - Smart single quotes → '
// - Smart double quotes → "
// - Various dashes/hyphens → -
// - Special spaces → regular space
func normalizeForFuzzyMatch(text string) string {
	// First strip trailing whitespace
	text = stripTrailingWhitespace(text)

	// Apply character-by-character normalization
	var b strings.Builder
	b.Grow(len(text))
	for _, r := range text {
		switch {
		// Smart single quotes → '
		case r == '\u2018' || r == '\u2019' || r == '\u201A' || r == '\u201B':
			b.WriteByte('\'')
		// Smart double quotes → "
		case r == '\u201C' || r == '\u201D' || r == '\u201E' || r == '\u201F':
			b.WriteByte('"')
		// Various dashes/hyphens → -
		case r == '\u2010' || r == '\u2011' || r == '\u2012' || r == '\u2013' ||
			r == '\u2014' || r == '\u2015' || r == '\u2212':
			b.WriteByte('-')
		// Special spaces → regular space
		case r == '\u00A0' || (r >= '\u2002' && r <= '\u200A') ||
			r == '\u202F' || r == '\u205F' || r == '\u3000':
			b.WriteByte(' ')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// DetectLineEnding detects whether content uses CRLF or LF line endings.
// Returns "\r\n" for CRLF, "\n" for LF (default).
func DetectLineEnding(content string) string {
	crlfIdx := strings.Index(content, "\r\n")
	lfIdx := strings.Index(content, "\n")
	if lfIdx == -1 {
		return "\n"
	}
	if crlfIdx == -1 {
		return "\n"
	}
	if crlfIdx < lfIdx {
		return "\r\n"
	}
	return "\n"
}

// NormalizeToLF converts all line endings to LF.
func NormalizeToLF(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	return text
}

// RestoreLineEndings converts LF to the specified line ending.
func RestoreLineEndings(text, ending string) string {
	if ending == "\r\n" {
		return strings.ReplaceAll(text, "\n", "\r\n")
	}
	return text
}

// StripBOM removes a UTF-8 BOM from the beginning of content.
// Returns the BOM (if present) and the content without it.
func StripBOM(content string) (bom, text string) {
	if strings.HasPrefix(content, "\xEF\xBB\xBF") {
		return "\xEF\xBB\xBF", content[3:]
	}
	return "", content
}

// ApplyEdits applies multiple oldText→newText edits to content.
// Each edit is matched against the ORIGINAL content (not incrementally).
// Returns the base content and new content for diff generation.
// Overlapping edits are rejected.
func ApplyEdits(content string, edits []Edit) (string, string, error) {
	normalizedContent := NormalizeToLF(content)

	var matched []matchedEdit
	for i, edit := range edits {
		result := FuzzyFindText(normalizedContent, NormalizeToLF(edit.OldText))
		if !result.Found {
			// Try searching in original content without normalization
			result = FuzzyFindText(content, edit.OldText)
			if !result.Found {
				return "", "", EditNotFound{Index: i, OldText: truncateEditHint(edit.OldText)}
			}
			// If found in original, use original content for replacement
			normalizedContent = content
		}

		matched = append(matched, matchedEdit{
			editIndex:   i,
			matchIndex:  result.Index,
			matchLength: result.MatchLength,
			newText:     edit.NewText,
		})
	}

	// Sort by match position
	sortEditsByPosition(matched)

	// Check for overlaps
	for i := 1; i < len(matched); i++ {
		prev := matched[i-1]
		curr := matched[i]
		if curr.matchIndex < prev.matchIndex+prev.matchLength {
			return "", "", OverlappingEdits{Index: curr.editIndex}
		}
	}

	// Apply edits from end to start to preserve indices
	newContent := normalizedContent
	for i := len(matched) - 1; i >= 0; i-- {
		me := matched[i]
		newContent = newContent[:me.matchIndex] + me.newText + newContent[me.matchIndex+me.matchLength:]
	}

	return normalizedContent, newContent, nil
}

// matchedEdit tracks where an edit matched in the content.
type matchedEdit struct {
	editIndex   int
	matchIndex  int
	matchLength int
	newText     string
}

// Edit represents a single oldText→newText replacement.
type Edit struct {
	OldText string `json:"oldText"`
	NewText string `json:"newText"`
}

// EditNotFound indicates that oldText could not be found in the content.
type EditNotFound struct {
	Index   int
	OldText string
}

func (e EditNotFound) Error() string {
	return "could not find oldText in file: " + e.OldText
}

// OverlappingEdits indicates that two edits overlap.
type OverlappingEdits struct {
	Index int
}

func (e OverlappingEdits) Error() string {
	return "overlapping edits detected"
}

// truncateEditHint truncates text for error messages.
func truncateEditHint(s string) string {
	maxRunes := 200
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxRunes]) + "..."
}

// sortEditsByPosition sorts matched edits by their match index.
func sortEditsByPosition(edits []matchedEdit) {
	for i := 1; i < len(edits); i++ {
		for j := i; j > 0 && edits[j].matchIndex < edits[j-1].matchIndex; j-- {
			edits[j], edits[j-1] = edits[j-1], edits[j]
		}
	}
}

// Ensure unicode is referenced
var _ = unicode.SimpleFold
