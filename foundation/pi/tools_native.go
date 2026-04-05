package pi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type ToolHandler func(ctx context.Context, cwd string, input json.RawMessage) (*ToolResult, error)

const (
	defaultGrepLimit  = 100
	defaultFindLimit  = 1000
	defaultLSLimit    = 500
	grepMaxLineLength = 500
)

func DefaultToolHandlers() map[string]ToolHandler {
	return map[string]ToolHandler{
		"read":  executeReadTool,
		"write": executeWriteTool,
		"edit":  executeEditTool,
		"bash":  executeBashTool,
		"grep":  executeGrepTool,
		"find":  executeFindTool,
		"ls":    executeLSTool,
	}
}

func executeReadTool(ctx context.Context, cwd string, raw json.RawMessage) (*ToolResult, error) {
	var input ReadToolInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, fmt.Errorf("invalid read input: %w", err)
	}
	if strings.TrimSpace(input.Path) == "" {
		return nil, fmt.Errorf("path is required")
	}
	absolutePath, err := resolvePath(cwd, input.Path)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(absolutePath)
	if err != nil {
		return nil, err
	}
	mimeType := mime.TypeByExtension(strings.ToLower(filepath.Ext(absolutePath)))
	if strings.HasPrefix(mimeType, "image/") {
		return &ToolResult{
			ToolName: "read",
			Content: []any{
				TextContent{Type: "text", Text: fmt.Sprintf("Read image file [%s]", mimeType)},
				ImageContent{Type: "image", Data: base64.StdEncoding.EncodeToString(data), MimeType: mimeType},
			},
		}, nil
	}

	lines := strings.Split(string(data), "\n")
	start := 0
	if input.Offset > 0 {
		start = input.Offset - 1
	}
	if start >= len(lines) {
		return nil, fmt.Errorf("offset %d is beyond end of file (%d lines total)", input.Offset, len(lines))
	}
	selected := lines[start:]
	if input.Limit > 0 && input.Limit < len(selected) {
		selected = selected[:input.Limit]
	}
	selectedText := strings.Join(selected, "\n")
	truncation, text := truncateHead(selectedText)
	output := text
	if truncation.FirstLineExceeds {
		lineNo := start + 1
		output = fmt.Sprintf("[Line %d exceeds %d byte limit. Use bash to inspect it safely.]", lineNo, DefaultMaxBytes)
	} else if truncation.Truncated {
		endLine := start + truncation.OutputLines
		output += fmt.Sprintf("\n\n[Showing lines %d-%d of %d. Use offset=%d to continue.]", start+1, endLine, len(lines), endLine+1)
	} else if input.Limit > 0 && start+len(selected) < len(lines) {
		output += fmt.Sprintf("\n\n[%d more lines in file. Use offset=%d to continue.]", len(lines)-(start+len(selected)), start+len(selected)+1)
	}
	var details *ReadToolDetails
	if truncation.Truncated || truncation.FirstLineExceeds {
		nextOffset := 0
		if truncation.OutputLines > 0 {
			nextOffset = start + truncation.OutputLines + 1
		}
		truncation.ContinuationOffset = nextOffset
		details = &ReadToolDetails{Truncation: &truncation}
	}
	result := &ToolResult{ToolName: "read", Content: []any{TextContent{Type: "text", Text: output}}}
	if details != nil {
		result.Details = details
	}
	return result, nil
}

func executeWriteTool(ctx context.Context, cwd string, raw json.RawMessage) (*ToolResult, error) {
	var input WriteToolInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, fmt.Errorf("invalid write input: %w", err)
	}
	if strings.TrimSpace(input.Path) == "" {
		return nil, fmt.Errorf("path is required")
	}
	absolutePath, err := resolvePath(cwd, input.Path)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(absolutePath, []byte(input.Content), 0o644); err != nil {
		return nil, err
	}
	return &ToolResult{
		ToolName: "write",
		Content:  []any{TextContent{Type: "text", Text: fmt.Sprintf("Successfully wrote %d bytes to %s", len(input.Content), input.Path)}},
	}, nil
}

func executeEditTool(ctx context.Context, cwd string, raw json.RawMessage) (*ToolResult, error) {
	var input EditToolInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, fmt.Errorf("invalid edit input: %w", err)
	}
	if strings.TrimSpace(input.Path) == "" {
		return nil, fmt.Errorf("path is required")
	}
	if len(input.Edits) == 0 {
		return nil, fmt.Errorf("edit tool input is invalid. edits must contain at least one replacement")
	}
	absolutePath, err := resolvePath(cwd, input.Path)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(absolutePath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %s", input.Path)
	}
	original := string(data)
	type replacement struct {
		start   int
		end     int
		newText string
	}
	replacements := make([]replacement, 0, len(input.Edits))
	for _, edit := range input.Edits {
		if edit.OldText == "" {
			return nil, fmt.Errorf("oldText must not be empty")
		}
		matches := findAllOccurrences(original, edit.OldText)
		if len(matches) == 0 {
			return nil, fmt.Errorf("oldText not found in %s", input.Path)
		}
		if len(matches) > 1 {
			return nil, fmt.Errorf("oldText must match a unique region in %s", input.Path)
		}
		replacements = append(replacements, replacement{start: matches[0], end: matches[0] + len(edit.OldText), newText: edit.NewText})
	}
	sort.Slice(replacements, func(i, j int) bool { return replacements[i].start < replacements[j].start })
	for i := 1; i < len(replacements); i++ {
		if replacements[i].start < replacements[i-1].end {
			return nil, fmt.Errorf("edits overlap or are nested in %s", input.Path)
		}
	}
	var out strings.Builder
	cursor := 0
	firstChangedLine := 1
	if len(replacements) > 0 {
		firstChangedLine = 1 + strings.Count(original[:replacements[0].start], "\n")
	}
	for _, replacement := range replacements {
		out.WriteString(original[cursor:replacement.start])
		out.WriteString(replacement.newText)
		cursor = replacement.end
	}
	out.WriteString(original[cursor:])
	updated := out.String()
	if err := os.WriteFile(absolutePath, []byte(updated), 0o644); err != nil {
		return nil, err
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(original, updated, false)
	return &ToolResult{
		ToolName: "edit",
		Content:  []any{TextContent{Type: "text", Text: fmt.Sprintf("Successfully replaced %d block(s) in %s.", len(input.Edits), input.Path)}},
		Details: &EditToolDetails{
			Diff:             dmp.DiffPrettyText(diffs),
			FirstChangedLine: firstChangedLine,
		},
	}, nil
}

func executeBashTool(ctx context.Context, cwd string, raw json.RawMessage) (*ToolResult, error) {
	var input BashToolInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, fmt.Errorf("invalid bash input: %w", err)
	}
	if strings.TrimSpace(input.Command) == "" {
		return nil, fmt.Errorf("command is required")
	}
	commandCtx := ctx
	var cancel context.CancelFunc
	if input.Timeout > 0 {
		commandCtx, cancel = context.WithTimeout(ctx, time.Duration(input.Timeout*float64(time.Second)))
		defer cancel()
	}
	cmd := shellCommand(commandCtx, input.Command)
	cmd.Dir = cwd
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	combined := stdout.String() + stderr.String()
	truncation, rendered := truncateTail(combined)
	var fullOutputPath string
	var details *BashToolDetails
	if truncation.Truncated {
		fullOutputPath = filepath.Join(os.TempDir(), fmt.Sprintf("pi-bash-%s.log", uuid.NewString()))
		if writeErr := os.WriteFile(fullOutputPath, []byte(combined), 0o600); writeErr == nil {
			truncation.FullOutputPath = fullOutputPath
			rendered += fmt.Sprintf("\n\n[Showing lines %d-%d of %d. Full output: %s]", truncation.TotalLines-truncation.OutputLines+1, truncation.TotalLines, truncation.TotalLines, fullOutputPath)
		}
		details = &BashToolDetails{Truncation: &truncation, FullOutputPath: fullOutputPath}
	}
	if rendered == "" {
		rendered = "(no output)"
	}
	if err != nil {
		if commandCtx.Err() == context.DeadlineExceeded {
			rendered += fmt.Sprintf("\n\nCommand timed out after %g seconds", input.Timeout)
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			rendered += fmt.Sprintf("\n\nCommand exited with code %d", exitErr.ExitCode())
		}
		result := &ToolResult{ToolName: "bash", Content: []any{TextContent{Type: "text", Text: rendered}}, IsError: true}
		if details != nil {
			result.Details = details
		}
		return result, fmt.Errorf("%s", rendered)
	}
	result := &ToolResult{ToolName: "bash", Content: []any{TextContent{Type: "text", Text: rendered}}}
	if details != nil {
		result.Details = details
	}
	return result, nil
}

func executeGrepTool(ctx context.Context, cwd string, raw json.RawMessage) (*ToolResult, error) {
	var input GrepToolInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, fmt.Errorf("invalid grep input: %w", err)
	}
	if strings.TrimSpace(input.Pattern) == "" {
		return nil, fmt.Errorf("pattern is required")
	}

	searchRoot := cwd
	if strings.TrimSpace(input.Path) != "" {
		var err error
		searchRoot, err = resolvePath(cwd, input.Path)
		if err != nil {
			return nil, err
		}
	}

	effectiveLimit := input.Limit
	if effectiveLimit <= 0 {
		effectiveLimit = defaultGrepLimit
	}

	matcher, err := compileGrepMatcher(input)
	if err != nil {
		return nil, err
	}

	var outputLines []string
	matchCount := 0
	err = filepath.Walk(searchRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if input.Glob != "" {
			matched, globErr := filepath.Match(input.Glob, filepath.Base(path))
			if globErr != nil || !matched {
				return nil
			}
		}
		lines, readErr := readFileLines(path)
		if readErr != nil {
			return nil
		}
		relPath := filepath.ToSlash(relPathOrBase(path, searchRoot))
		for i, line := range lines {
			if matchCount >= effectiveLimit {
				return nil
			}
			if !matcher(line) {
				continue
			}
			matchCount++
			if input.Context > 0 {
				start := max(1, i+1-input.Context)
				end := min(len(lines), i+1+input.Context)
				for j := start; j <= end; j++ {
					text, _ := truncateLineStr(strings.ReplaceAll(lines[j-1], "\r", ""))
					if j == i+1 {
						outputLines = append(outputLines, fmt.Sprintf("%s:%d:%s", relPath, j, text))
					} else {
						outputLines = append(outputLines, fmt.Sprintf("%s-%d-%s", relPath, j, text))
					}
				}
			} else {
				text, _ := truncateLineStr(strings.ReplaceAll(line, "\r", ""))
				outputLines = append(outputLines, fmt.Sprintf("%s:%d:%s", relPath, i+1, text))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if matchCount == 0 {
		return &ToolResult{ToolName: "grep", Content: []any{TextContent{Type: "text", Text: "No matches found"}}}, nil
	}

	rawOutput := strings.Join(outputLines, "\n")
	truncation, output := truncateHead(rawOutput)
	var details *GrepToolDetails
	var notices []string
	if matchCount >= effectiveLimit {
		notices = append(notices, fmt.Sprintf("%d matches limit reached. Use limit=%d for more, or refine pattern", effectiveLimit, effectiveLimit*2))
		details = &GrepToolDetails{MatchLimitReached: effectiveLimit}
	}
	if truncation.Truncated {
		notices = append(notices, fmt.Sprintf("%d byte limit reached", DefaultMaxBytes))
		if details == nil {
			details = &GrepToolDetails{}
		}
		details.Truncation = &truncation
	}
	if len(notices) > 0 {
		output += fmt.Sprintf("\n\n[%s]", strings.Join(notices, ". "))
	}
	result := &ToolResult{ToolName: "grep", Content: []any{TextContent{Type: "text", Text: output}}}
	if details != nil {
		result.Details = details
	}
	return result, nil
}

func executeFindTool(ctx context.Context, cwd string, raw json.RawMessage) (*ToolResult, error) {
	var input FindToolInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, fmt.Errorf("invalid find input: %w", err)
	}
	if strings.TrimSpace(input.Pattern) == "" {
		return nil, fmt.Errorf("pattern is required")
	}
	searchRoot := cwd
	if strings.TrimSpace(input.Path) != "" {
		var err error
		searchRoot, err = resolvePath(cwd, input.Path)
		if err != nil {
			return nil, err
		}
	}
	effectiveLimit := input.Limit
	if effectiveLimit <= 0 {
		effectiveLimit = defaultFindLimit
	}

	var results []string
	err := filepath.Walk(searchRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if len(results) >= effectiveLimit {
			return nil
		}
		matched, globErr := filepath.Match(input.Pattern, filepath.Base(path))
		if globErr != nil || !matched {
			relPath, relErr := filepath.Rel(searchRoot, path)
			if relErr != nil {
				return nil
			}
			matched, _ = filepath.Match(input.Pattern, filepath.ToSlash(relPath))
		}
		if matched {
			results = append(results, filepath.ToSlash(relPathOrBase(path, searchRoot)))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return &ToolResult{ToolName: "find", Content: []any{TextContent{Type: "text", Text: "No files found matching pattern"}}}, nil
	}

	rawOutput := strings.Join(results, "\n")
	truncation, output := truncateHead(rawOutput)
	var details *FindToolDetails
	var notices []string
	if len(results) >= effectiveLimit {
		notices = append(notices, fmt.Sprintf("%d results limit reached. Use limit=%d for more, or refine pattern", effectiveLimit, effectiveLimit*2))
		details = &FindToolDetails{ResultLimitReached: effectiveLimit}
	}
	if truncation.Truncated {
		notices = append(notices, fmt.Sprintf("%d byte limit reached", DefaultMaxBytes))
		if details == nil {
			details = &FindToolDetails{}
		}
		details.Truncation = &truncation
	}
	if len(notices) > 0 {
		output += fmt.Sprintf("\n\n[%s]", strings.Join(notices, ". "))
	}
	result := &ToolResult{ToolName: "find", Content: []any{TextContent{Type: "text", Text: output}}}
	if details != nil {
		result.Details = details
	}
	return result, nil
}

func executeLSTool(ctx context.Context, cwd string, raw json.RawMessage) (*ToolResult, error) {
	var input LSToolInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, fmt.Errorf("invalid ls input: %w", err)
	}
	listRoot := cwd
	if strings.TrimSpace(input.Path) != "" {
		var err error
		listRoot, err = resolvePath(cwd, input.Path)
		if err != nil {
			return nil, err
		}
	}
	info, err := os.Stat(listRoot)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", input.Path)
	}
	entries, err := os.ReadDir(listRoot)
	if err != nil {
		return nil, err
	}
	effectiveLimit := input.Limit
	if effectiveLimit <= 0 {
		effectiveLimit = defaultLSLimit
	}
	results := make([]string, 0, min(effectiveLimit, len(entries)))
	entryLimitReached := false
	for _, entry := range entries {
		if len(results) >= effectiveLimit {
			entryLimitReached = true
			break
		}
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		results = append(results, name)
	}
	sort.Slice(results, func(i, j int) bool {
		return strings.ToLower(results[i]) < strings.ToLower(results[j])
	})
	if len(results) == 0 {
		return &ToolResult{ToolName: "ls", Content: []any{TextContent{Type: "text", Text: "(empty directory)"}}}, nil
	}
	rawOutput := strings.Join(results, "\n")
	truncation, output := truncateHead(rawOutput)
	var details *LsToolDetails
	var notices []string
	if entryLimitReached {
		notices = append(notices, fmt.Sprintf("%d entries limit reached. Use limit=%d for more", effectiveLimit, effectiveLimit*2))
		details = &LsToolDetails{EntryLimitReached: effectiveLimit}
	}
	if truncation.Truncated {
		notices = append(notices, fmt.Sprintf("%d byte limit reached", DefaultMaxBytes))
		if details == nil {
			details = &LsToolDetails{}
		}
		details.Truncation = &truncation
	}
	if len(notices) > 0 {
		output += fmt.Sprintf("\n\n[%s]", strings.Join(notices, ". "))
	}
	result := &ToolResult{ToolName: "ls", Content: []any{TextContent{Type: "text", Text: output}}}
	if details != nil {
		result.Details = details
	}
	return result, nil
}

func resolvePath(cwd, toolPath string) (string, error) {
	if strings.TrimSpace(toolPath) == "" {
		return "", fmt.Errorf("path is required")
	}
	if filepath.IsAbs(toolPath) {
		return filepath.Clean(toolPath), nil
	}
	return filepath.Join(cwd, filepath.Clean(toolPath)), nil
}

func findAllOccurrences(haystack, needle string) []int {
	var matches []int
	for start := 0; start <= len(haystack)-len(needle); {
		idx := strings.Index(haystack[start:], needle)
		if idx < 0 {
			break
		}
		absolute := start + idx
		matches = append(matches, absolute)
		start = absolute + len(needle)
	}
	return matches
}

func truncateHead(text string) (TruncationDetails, string) {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return TruncationDetails{}, ""
	}
	if len(lines[0]) > DefaultMaxBytes {
		return TruncationDetails{Truncated: true, FirstLineExceeds: true, MaxBytes: DefaultMaxBytes, TotalLines: len(lines)}, ""
	}
	var out []string
	bytesUsed := 0
	for i, line := range lines {
		lineBytes := len(line)
		if i < len(lines)-1 {
			lineBytes++
		}
		if len(out) >= DefaultMaxLines {
			return TruncationDetails{Truncated: true, TruncatedBy: "lines", TotalLines: len(lines), OutputLines: len(out), OutputBytes: bytesUsed, MaxLines: DefaultMaxLines, MaxBytes: DefaultMaxBytes}, strings.Join(out, "\n")
		}
		if bytesUsed+lineBytes > DefaultMaxBytes {
			return TruncationDetails{Truncated: true, TruncatedBy: "bytes", TotalLines: len(lines), OutputLines: len(out), OutputBytes: bytesUsed, MaxLines: DefaultMaxLines, MaxBytes: DefaultMaxBytes}, strings.Join(out, "\n")
		}
		out = append(out, line)
		bytesUsed += lineBytes
	}
	return TruncationDetails{TotalLines: len(lines), OutputLines: len(out), OutputBytes: bytesUsed, MaxLines: DefaultMaxLines, MaxBytes: DefaultMaxBytes}, strings.Join(out, "\n")
}

func truncateTail(text string) (TruncationDetails, string) {
	lines := strings.Split(text, "\n")
	if len(lines) <= DefaultMaxLines && len(text) <= DefaultMaxBytes {
		return TruncationDetails{TotalLines: len(lines), OutputLines: len(lines), OutputBytes: len(text), MaxLines: DefaultMaxLines, MaxBytes: DefaultMaxBytes}, text
	}
	selected := make([]string, 0, min(DefaultMaxLines, len(lines)))
	bytesUsed := 0
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		lineBytes := len(line)
		if len(selected) > 0 {
			lineBytes++
		}
		if len(selected) >= DefaultMaxLines {
			break
		}
		if bytesUsed+lineBytes > DefaultMaxBytes {
			break
		}
		selected = append(selected, line)
		bytesUsed += lineBytes
	}
	for i, j := 0, len(selected)-1; i < j; i, j = i+1, j-1 {
		selected[i], selected[j] = selected[j], selected[i]
	}
	truncatedBy := "lines"
	if len(text) > DefaultMaxBytes {
		truncatedBy = "bytes"
	}
	return TruncationDetails{Truncated: true, TruncatedBy: truncatedBy, TotalLines: len(lines), OutputLines: len(selected), OutputBytes: bytesUsed, MaxLines: DefaultMaxLines, MaxBytes: DefaultMaxBytes}, strings.Join(selected, "\n")
}

func compileGrepMatcher(input GrepToolInput) (func(string) bool, error) {
	if input.Literal {
		pattern := input.Pattern
		if input.IgnoreCase {
			pattern = strings.ToLower(pattern)
			return func(line string) bool {
				return strings.Contains(strings.ToLower(line), pattern)
			}, nil
		}
		return func(line string) bool {
			return strings.Contains(line, pattern)
		}, nil
	}
	pattern := input.Pattern
	if input.IgnoreCase {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex: %w", err)
	}
	return re.MatchString, nil
}

func relPathOrBase(path, base string) string {
	if rel, err := filepath.Rel(base, path); err == nil && !strings.HasPrefix(rel, "..") {
		return rel
	}
	return filepath.Base(path)
}

func readFileLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n"), nil
}

func truncateLineStr(s string) (string, bool) {
	runes := []rune(s)
	if len(runes) > grepMaxLineLength {
		return string(runes[:grepMaxLineLength]), true
	}
	return s, false
}

func truncateDisplay(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func shellCommand(ctx context.Context, command string) *exec.Cmd {
	if isWindows() {
		return exec.CommandContext(ctx, "cmd", "/C", command)
	}
	return exec.CommandContext(ctx, "sh", "-lc", command)
}
