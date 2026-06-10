package agents

import (
	"context"
	"strings"
)

// OfflineProvider is a built-in fallback that provides useful responses
// when no real LLM API is available. It performs basic pattern matching
// and delegates to local tools when possible.
type OfflineProvider struct {
	WorkingDir string
}

func NewOfflineProvider(workingDir string) *OfflineProvider {
	return &OfflineProvider{WorkingDir: workingDir}
}

func (p *OfflineProvider) Chat(ctx context.Context, messages []Message, tools []Tool) (Message, error) {
	// Extract the last user message
	var lastUserMsg string
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == RoleUser {
			lastUserMsg = messages[i].Content
			break
		}
	}

	if strings.TrimSpace(lastUserMsg) == "" {
		return Message{Role: RoleAssistant, Content: "I'm ready to help. What would you like to do?"}, nil
	}

	lower := strings.ToLower(lastUserMsg)

	// Pattern-based responses for common queries
	response := p.generateResponse(lower, lastUserMsg)

	return Message{Role: RoleAssistant, Content: response}, nil
}

func (p *OfflineProvider) generateResponse(lower, original string) string {
	// File operations
	if strings.Contains(lower, "read") && (strings.Contains(lower, "file") || containsPath(original)) {
		return "I'd read that file for you, but I'm running in offline mode. To enable real LLM responses, set one of these environment variables:\n\n" +
			"  • ANTHROPIC_API_KEY or CLAUDE_API_KEY — for Claude models\n" +
			"  • OPENAI_API_KEY — for GPT models\n" +
			"  • GEMINI_API_KEY or GOOGLE_API_KEY — for Gemini models\n" +
			"  • DEEPSEEK_API_KEY — for DeepSeek models\n" +
			"  • OPENROUTER_API_KEY — for OpenRouter (multi-model)\n\n" +
			"You can also use shell commands directly with the ! prefix:\n  !cat " + extractPath(original)
	}

	// Shell commands
	if strings.Contains(lower, "run") || strings.Contains(lower, "execute") || strings.Contains(lower, "command") {
		return "I'm in offline mode — no LLM API key detected. You can run shell commands directly with:\n\n" +
			"  !<command>     — run a shell command\n" +
			"  !!<command>    — run and show output\n\n" +
			"To enable AI-powered responses, set an API key (see /doctor for details)."
	}

	// Help requests
	if strings.Contains(lower, "help") {
		return "HyperHarness Offline Mode\n\n" +
			"I can't provide AI-powered responses without an LLM API key, but you can still:\n\n" +
			"  • Use !<command> to run shell commands directly\n" +
			"  • Use /tools to see all available tools\n" +
			"  • Use /doctor to check your configuration\n" +
			"  • Use /model to switch models\n" +
			"  • Use /diff to see git changes\n" +
			"  • Use /memory to manage memory files\n\n" +
			"To enable AI, set an API key:\n" +
			"  export ANTHROPIC_API_KEY=sk-ant-...\n" +
			"  export OPENAI_API_KEY=sk-...\n" +
			"  export GEMINI_API_KEY=AIza..."
	}

	// Code-related
	if strings.Contains(lower, "code") || strings.Contains(lower, "function") || strings.Contains(lower, "implement") || strings.Contains(lower, "write") {
		return "I'm running in offline mode without an LLM API key. For AI-powered code generation:\n\n" +
			"1. Set an API key:\n" +
			"   export ANTHROPIC_API_KEY=your-key\n\n" +
			"2. Or use the built-in tools:\n" +
			"   !vim file.go       — edit with your $EDITOR\n" +
			"   /tree              — browse files\n" +
			"   /tools             — see all 136+ tools\n\n" +
			"3. Or start the HyperCode server:\n" +
			"   ./hyperharness serve"
	}

	// Questions
	if strings.Contains(lower, "?") || strings.Contains(lower, "what") || strings.Contains(lower, "how") {
		return "I'm in offline mode — no LLM provider configured. To get AI-powered answers:\n\n" +
			"  Set an API key:  export ANTHROPIC_API_KEY=sk-ant-...\n" +
			"  Then restart:    ./hyperharness\n\n" +
			"Available providers (auto-detected from env):\n" +
			"  ANTHROPIC_API_KEY  → Claude Sonnet 4\n" +
			"  OPENAI_API_KEY     → GPT-4o\n" +
			"  GEMINI_API_KEY     → Gemini 2.5 Flash\n" +
			"  DEEPSEEK_API_KEY   → DeepSeek Chat\n" +
			"  OPENROUTER_API_KEY → Multi-model via OpenRouter\n\n" +
			"Run /doctor for a full diagnostic."
	}

	// Default response
	return "HyperHarness is in offline mode. No LLM API key detected.\n\n" +
		"To enable AI-powered responses, set one of these environment variables before starting:\n\n" +
		"  export ANTHROPIC_API_KEY=sk-ant-...    # Claude\n" +
		"  export OPENAI_API_KEY=sk-...           # GPT\n" +
		"  export GEMINI_API_KEY=AIza...          # Gemini\n" +
		"  export DEEPSEEK_API_KEY=...            # DeepSeek\n" +
		"  export OPENROUTER_API_KEY=sk-or-...    # OpenRouter\n\n" +
		"Quick local tools available now:\n" +
		"  !<command>  → run shell commands\n" +
		"  /doctor     → check configuration\n" +
		"  /tools      → list all tools\n" +
		"  /tree       → browse files"
}

func (p *OfflineProvider) Stream(ctx context.Context, messages []Message, tools []Tool, chunkChan chan<- string) error {
	msg, err := p.Chat(ctx, messages, tools)
	if err != nil {
		return err
	}
	chunkChan <- msg.Content
	close(chunkChan)
	return nil
}

func (p *OfflineProvider) GetModelName() string {
	return "offline"
}

func containsPath(s string) bool {
	return strings.Contains(s, "/") || strings.Contains(s, ".") || strings.Contains(s, "\\")
}

func extractPath(s string) string {
	// Simple heuristic: extract what looks like a file path
	words := strings.Fields(s)
	for _, w := range words {
		w = strings.Trim(w, "\"'`.,;:!?()")
		if strings.Contains(w, "/") || strings.HasSuffix(w, ".go") || strings.HasSuffix(w, ".ts") || strings.HasSuffix(w, ".py") {
			return w
		}
	}
	return "path/to/file"
}
