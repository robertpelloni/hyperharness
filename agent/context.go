package agent

import (
	"github.com/robertpelloni/hyperharness/internal/ast"
	"strings"
	"github.com/sashabaranov/go-openai"
)

// TrimHistory ensures the message history doesn't exceed a certain length,
// mimicking the context window management of advanced agents like Claude Code.
func (a *Agent) TrimHistory(maxMessages int) {
	if maxMessages <= 0 || len(a.messages) == 0 || len(a.messages) <= maxMessages {
		return
	}

	// Keep the system prompt (index 0) and the most recent messages.
	trimmed := make([]openai.ChatCompletionMessage, 0, maxMessages)
	trimmed = append(trimmed, a.messages[0])
	if maxMessages == 1 {
		a.messages = trimmed
		return
	}

	startIndex := len(a.messages) - (maxMessages - 1)
	if startIndex < 1 {
		startIndex = 1
	}
	trimmed = append(trimmed, a.messages[startIndex:]...)

	a.messages = trimmed
}

// CompactTokenHeavyMessages scans the current messages array for extremely long tool outputs
// (like reading massive source files). If a message is too large, it replaces the raw
// Go file output with an AST structural summary to conserve tokens.
func (a *Agent) CompactTokenHeavyMessages(maxCharThreshold int) {
	if a == nil || len(a.messages) <= 1 {
		return
	}

	for i := 1; i < len(a.messages); i++ {
		msg := a.messages[i]
		if msg.Role == openai.ChatMessageRoleTool && len(msg.Content) > maxCharThreshold {
			// If it looks like a go file or has package main, we can attempt an AST summarization.
			if strings.Contains(msg.Content, "package ") && strings.Contains(msg.Content, "func ") {
				summary, err := ast.SummarizeGoFile("dynamic.go", []byte(msg.Content))
				if err == nil && summary != "" {
					a.messages[i].Content = "[Context Window Auto-Compacted AST Summary]\n" + summary
				} else {
					// Fallback to simple truncation
					a.messages[i].Content = msg.Content[:maxCharThreshold] + "\n... (output truncated due to context limits)"
				}
			} else {
				a.messages[i].Content = msg.Content[:maxCharThreshold] + "\n... (output truncated due to context limits)"
			}
		}
	}
}
