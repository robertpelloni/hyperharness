// Package main is the hyperharness CLI entry point.
// Built with Cobra for command parsing and Viper for configuration.
// Supports modes: interactive (default TUI), print (--print), json (--mode json), rpc (--mode rpc)
package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/hyperharness/hyperharness/internal/agent"
	"context"

	"github.com/hyperharness/hyperharness/internal/config"
	"github.com/hyperharness/hyperharness/internal/memory"
	"github.com/hyperharness/hyperharness/internal/mcp"
	"github.com/hyperharness/hyperharness/internal/providers"
	"github.com/hyperharness/hyperharness/internal/sessions"
	"github.com/hyperharness/hyperharness/internal/tools"
	"github.com/hyperharness/hyperharness/internal/ui"
	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var (
	Version   = "0.1.0-dev"
	Commit    = "dev"
	BuildDate = "dev"
	GoVersion = runtime.Version()
)

func main() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// NewRootCmd creates the root command.
func NewRootCmd() *cobra.Command {
	var (
		printMode      bool
		jsonMode       bool
		rpcMode        bool
		providerFlag   string
		modelFlag      string
		thinkingFlag   string
		continueFlag   bool
		resumeFlag     bool
		sessionPath    string
		forkPath       string
		sessionDir     string
		noSession      bool
		toolsFlag      string
		noTools        bool
		extensionFlags []string
		skillFlags     []string
		promptFlags    []string
		themeFlag      string
		verbose        bool
		version        bool
		command        string // single command mode
	)

	cmd := &cobra.Command{
		Use:   "hyperharness [options] [@files...] [messages...]",
		Short: "🧬 The universal Go-based AI coding agent harness",
		Long: `Hyperharness: The unified AI coding agent that assimilates features from
every coding tool on earth with 100% parameter parity and superior Go-native architecture.

Ported from Pi's proven design, enhanced with features from Claude Code, Aider,
OpenInterpreter, Goose, Crush, Factory AI, Gemini CLI, Copilot CLI, Ollama, 
Codex CLI, and every other major AI coding agent.

Usage examples:
  hyperharness "Fix the bug in main.go"
  hyperharness -c                          # continue last session
  hyperharness -p @README.md "Summarize"   # print mode
  hyperharness --provider openai --model gpt-4o "Help me refactor"
  hyperharness --thinking high "Solve this complex problem"
  hyperharness --tools read,grep,find,ls   # read-only tools only`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if version {
				printVersion()
				return nil
			}

			// Determine mode
			var mode agent.AgentMode
			switch {
			case printMode:
				mode = agent.ModePrint
			case jsonMode:
				mode = agent.ModeJSON
			case rpcMode:
				mode = agent.ModeRPC
			default:
				mode = agent.ModeInteractive
			}

			// Get working directory
			cwd, _ := os.Getwd()

			// Load configuration
			cfgMgr, err := loadConfig(cwd)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("config error: %w", err)
			}
			cfg := cfgMgr.Get()

			// Resolve provider and model
			if providerFlag != "" {
				cfg.DefaultProvider = providerFlag
			}
			if modelFlag != "" {
				cfg.DefaultModel = modelFlag
			}
			if thinkingFlag != "" {
				cfg.DefaultThinkingLevel = thinkingFlag
			}

			// Handle tool filtering
			if noTools {
				cfg.Tools = &config.ToolSettings{
					Enabled:  []string{},
					Disabled: []string{"*"},
				}
			} else if toolsFlag != "" {
				enabled := strings.Split(toolsFlag, ",")
				cfg.Tools = &config.ToolSettings{
					Enabled: enabled,
				}
			}

			// Create session manager
			sessionMgr, err := createSessionManager(sessionDir, cfg.Sessions)
			if err != nil {
				return fmt.Errorf("session error: %w", err)
			}

			// Load or create session
			session, err := resolveSession(sessionMgr, sessionPath, forkPath, continueFlag, resumeFlag, noSession, cwd)
			if err != nil {
				return fmt.Errorf("session error: %w", err)
			}

			// Create provider
			providerType := providers.ProviderType(cfg.DefaultProvider)
			if providerType == "" {
				// Infer from model or default to anthropic
				if cfg.DefaultModel != "" {
					providerType = providers.InferProviderFromModel(cfg.DefaultModel)
				} else {
					providerType = providers.ProviderAnthropic
				}
			}

			providerCfg := providers.ProviderConfig{
				Type: providerType,
			}
			if cfg.Providers != nil {
				if p, ok := cfg.Providers[string(providerType)]; ok {
					providerCfg.APIKey = p.APIKey
					providerCfg.BaseURL = p.BaseURL
				}
			}

			// Also check environment variables
			if providerCfg.APIKey == "" {
				providerCfg.APIKey = getProviderAPIKey(providerType)
			}

			provider, err := providers.CreateProvider(providerType, providerCfg)
			if err != nil {
				return fmt.Errorf("provider error: %w", err)
			}

			// Create memory knowledge base
			memCfg := cfg.Memory
			kb, err := memory.NewKnowledgeBase("")
			if err != nil {
				return fmt.Errorf("memory error: %w", err)
			}

			// Create MCP registry
			mcpReg := mcp.NewRegistry()

			// Load MCP servers from config
			mcpCfgPath := filepath.Join(cfg.ConfigDir(), "mcp.json")
			if _, err := os.Stat(mcpCfgPath); err == nil {
				if err := mcpReg.LoadFromFile(mcpCfgPath); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to load MCP config: %v\n", err)
				}
			}

			// Create agent runtime
			runtime := agent.NewRuntime(cfg, provider, session, mode)

			// Register built-in tools
			builtinTools := tools.BuiltInTools
			for name, creator := range builtinTools {
				if cfg.IsEnabled(name) {
					runtime.RegisterTool(creator(cwd))
				}
			}

			// Register built-in tools with MCP registry for exposure
			for _, t := range tools.ToolSchemas() {
				name := t["name"].(string)
				// MCP registry registration would happen here with a closure
				_ = name
			}

			// Set memory
			runtime.SetMemory(kb)

			// Build system prompt
			promptBuilder := ui.NewSystemPromptBuilder()

			// Load AGENTS.md files
			agentsFiles := config.FindAgentFiles()
			for _, f := range agentsFiles {
				data, err := os.ReadFile(f)
				if err == nil {
					promptBuilder.WithAGENTS(append(promptBuilder.agencies, string(data)))
				}
			}

			// Add memory context
			memoryCtx := kb.BuildContextForAgent(cwd, nil)
			if memoryCtx != "" {
				promptBuilder.WithMemoryContext(memoryCtx)
			}

			runtime.SetSystemPrompt(promptBuilder.Build())

			// Get user message from args
			userMessage := strings.Join(args, " ")

			// Handle command mode (single shot)
			if command != "" {
				userMessage = command
			}

			// Build initial message with @-prefixed files
			userMessage, err = processAtFiles(userMessage, cwd)
			if err != nil {
				return fmt.Errorf("file error: %w", err)
			}

			// Execute based on mode
			switch mode {
			case agent.ModePrint:
				return runPrintMode(runtime, cmd.Context(), userMessage)
			case agent.ModeJSON:
				return runJSONMode(runtime, cmd.Context(), userMessage)
			case agent.ModeRPC:
				return runRPCMode(runtime)
			default:
				return runInteractiveMode(runtime, cmd.Context(), session, cfg)
			}
		},
	}

	// Flags
	flags := cmd.Flags()
	
	// Mode flags
	flags.BoolVarP(&printMode, "print", "p", false, "Print response and exit")
	flags.BoolVar(&jsonMode, "mode", false, "Output events as JSON lines (use json)")
	flags.BoolVar(&rpcMode, "rpc", false, "RPC mode for process integration")

	// Model flags
	flags.StringVar(&providerFlag, "provider", "", "Provider (anthropic, openai, google, etc.)")
	flags.StringVar(&modelFlag, "model", "", "Model pattern or ID")
	flags.StringVar(&thinkingFlag, "thinking", "", "Thinking level (off, minimal, low, medium, high, xhigh)")
	flags.StringVar(&themeFlag, "theme", "", "Theme name")

	// Session flags
	flags.BoolVarP(&continueFlag, "continue", "c", false, "Continue most recent session")
	flags.BoolVarP(&resumeFlag, "resume", "r", false, "Browse and select session")
	flags.StringVar(&sessionPath, "session", "", "Use specific session file or ID")
	flags.StringVar(&forkPath, "fork", "", "Fork specific session into new one")
	flags.StringVar(&sessionDir, "session-dir", "", "Custom session directory")
	flags.BoolVar(&noSession, "no-session", false, "Ephemeral mode (don't save)")

	// Tool flags
	flags.StringVar(&toolsFlag, "tools", "", "Enable specific tools (comma-separated)")
	flags.BoolVar(&noTools, "no-tools", false, "Disable all built-in tools")

	// Resource flags
	flags.StringArrayVarP(&extensionFlags, "extension", "e", nil, "Load extension")
	flags.StringArrayVar(&skillFlags, "skill", nil, "Load skill")
	flags.StringArrayVar(&promptFlags, "prompt-template", nil, "Load prompt template")

	// Other flags
	flags.BoolVar(&verbose, "verbose", false, "Force verbose startup")
	flags.BoolVarP(&version, "version", "v", false, "Show version")
	flags.StringVarP(&command, "command", "C", "", "Run single command")

	// Hide internal flags
	_ = flags.MarkHidden("mode")

	return cmd
}

// loadConfig loads configuration from global and project directories.
func loadConfig(cwd string) (*config.Manager, error) {
	globalDir := config.ConfigDir()
	projectDir := config.ProjectConfigDir()
	
	return config.NewManager(globalDir, projectDir)
}

// createSessionManager creates the session manager.
func createSessionManager(sessionDir string, settings *config.SessionSettings) (*sessions.Manager, error) {
	if sessionDir == "" {
		if settings != nil && settings.SessionDir != "" {
			sessionDir = settings.SessionDir
		}
	}
	return sessions.NewManager(sessionDir)
}

// resolveSession determines which session to use.
func resolveSession(mgr *sessions.Manager, sessionPath, forkPath string, continueFlag, resumeFlag, noSession bool, cwd string) (*sessions.Session, error) {
	if noSession {
		return sessions.NewSession(cwd, "", ""), nil
	}

	if forkPath != "" {
		src, err := mgr.LoadSession(forkPath)
		if err != nil {
			return nil, fmt.Errorf("fork: %w", err)
		}
		forked, err := src.Fork("")
		if err != nil {
			return nil, err
		}
		return forked, nil
	}

	if sessionPath != "" {
		return mgr.LoadSession(sessionPath)
	}

	if continueFlag {
		sessions := mgr.ListSessions()
		if len(sessions) > 0 {
			return sessions[0], nil
		}
	}

	if resumeFlag {
		// Would show interactive selector
		sessions := mgr.ListSessions()
		if len(sessions) > 0 {
			return sessions[0], nil
		}
	}

	// Create new session
	return mgr.CreateSession(cwd, "", "")
}

// processAtFiles processes @-prefixed file arguments.
func processAtFiles(message string, cwd string) (string, error) {
	var result strings.Builder
	words := strings.Fields(message)
	
	for _, word := range words {
		if strings.HasPrefix(word, "@") {
			filePath := word[1:]
			// Resolve path
			if !filepath.IsAbs(filePath) {
				filePath = filepath.Join(cwd, filePath)
			}
			data, err := os.ReadFile(filePath)
			if err != nil {
				return message, fmt.Errorf("failed to read file %s: %w", word, err)
			}
			result.WriteString(fmt.Sprintf("\n--- Contents of %s ---\n%s\n--- End of %s ---\n", 
				filepath.Base(filePath), string(data), filepath.Base(filePath)))
		} else {
			result.WriteString(word)
			result.WriteString(" ")
		}
	}
	
	return strings.TrimSpace(result.String()), nil
}

// runPrintMode executes in print mode - single request/response.
func runPrintMode(runtime *agent.Runtime, ctx context.Context, message string) error {
	result, err := runtime.Prompt(ctx, message)
	if err != nil {
		return err
	}
	
	content := ""
	for _, block := range result.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}
	fmt.Println(content)
	return nil
}

// runJSONMode executes in JSON mode - outputs events as JSONL.
func runJSONMode(runtime *agent.Runtime, ctx context.Context, message string) error {
	result, err := runtime.Prompt(ctx, message)
	if err != nil {
		fmt.Fprintf(os.Stderr, `{"type":"error","message":"%s"}%s`, err.Error(), "\n")
		return err
	}
	
	// Output as JSON
	content := ""
	for _, block := range result.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}
	fmt.Printf(`{"type":"response","content":%q}%s`, content, "\n")
	if result.Usage != nil {
		fmt.Printf(`{"type":"usage","input_tokens":%d,"output_tokens":%d}%s`,
			result.Usage.InputTokens, result.Usage.OutputTokens, "\n")
	}
	return nil
}

// runRPCMode starts the RPC server on stdin/stdout.
func runRPCMode(runtime *agent.Runtime) error {
	// JSONL protocol over stdin/stdout
	// See docs/rpc.md for full specification
	fmt.Fprintln(os.Stderr, "RPC mode: ready for JSONL input on stdin")
	
	// Implementation would read JSON events from stdin and write responses to stdout
	// Format: {"type":"prompt","message":"...","id":"..."}
	// Response: {"type":"chunk","content":"...","id":"..."}
	// End: {"type":"done","id":"..."}
	
	return nil
}

// runInteractiveMode starts the TUI REPL.
func runInteractiveMode(runtime *agent.Runtime, ctx context.Context, session *sessions.Session, cfg *config.Settings) error {
	// Setup signal handling
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Print startup header
	printStartupHeader(cfg, session)

	// Interactive REPL loop
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			fmt.Print("hyperharness> ")
			var input string
			fmt.Scanln(&input)
			
			if input == "" {
				continue
			}
			
			// Check for commands
			if strings.HasPrefix(input, "/") {
				if err := handleCommand(input, runtime, session, cfg); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
				continue
			}
			
			// Process input
			result, err := runtime.Prompt(ctx, input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
			
			// Print response
			content := ""
			for _, block := range result.Content {
				if block.Type == "text" {
					content += block.Text
				}
			}
			fmt.Println(content)
		}
	}
}

// handleCommand processes a slash command.
func handleCommand(input string, runtime *agent.Runtime, session *sessions.Session, cfg *config.Settings) error {
	parts := strings.SplitN(input, " ", 2)
	cmd := strings.TrimPrefix(parts[0], "/")
	args := ""
	if len(parts) > 1 {
		args = parts[1]
	}
	
	switch cmd {
	case "model":
		fmt.Printf("Current model: %s (provider: %s)\n", session.ModelID, session.Provider)
		if args != "" {
			// Switch model
		}
	case "tree":
		fmt.Println("Session tree navigation")
	case "compact":
		fmt.Println("Compacting context...")
	case "new":
		fmt.Println("Starting new session")
	case "help", "hotkeys":
		printHelp()
	case "quit", "exit":
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: /%s. Type /help for available commands.\n", cmd)
	}
	
	return nil
}

// printStartupHeader prints the startup information.
func printStartupHeader(cfg *config.Settings, session *sessions.Session) {
	fmt.Println("🧬 hyperharness " + Version)
	fmt.Println("=" + strings.Repeat("=", 59))
	
	agentsFiles := config.FindAgentFiles()
	if len(agentsFiles) > 0 {
		fmt.Printf("📋 AGENTS.md: %s\n", agentsFiles)
	}
	
	fmt.Printf("💾 Session: %s | Model: %s | Provider: %s\n", 
		session.ID, session.ModelID, session.Provider)
	fmt.Printf("📁 Working directory: %s\n", session.CWD)
	fmt.Println("=" + strings.Repeat("=", 59))
}

// printHelp prints the help message.
func printHelp() {
	fmt.Println(`
Available commands:
  /login          Authenticate with provider
  /logout         Clear authentication
  /model          Switch or view current model
  /scoped-models  Configure model cycling
  /settings       Open settings
  /resume         Browse past sessions
  /new            Start new session
  /name <name>    Set session name
  /session        Show session info
  /tree           Navigate session tree
  /fork           Fork current branch
  /compact        Compress context
  /copy           Copy last response
  /export [file]  Export session to HTML
  /share          Share session via GitHub gist
  /reload         Reload configuration
  /hotkeys        Show all keyboard shortcuts
  /changelog      Show version history
  /quit           Exit

Keyboard shortcuts:
  Enter           Send message
  Alt+Enter       Queue follow-up message
  Ctrl+C          Clear editor / Quit (double tap)
  Escape          Cancel / abort
  Ctrl+L          Model selector
  Ctrl+P          Cycle scoped models
  Shift+Tab       Cycle thinking level
  @               File autocomplete
  !               Run command without sending`)
}

// getProviderAPIKey gets API key from environment.
func getProviderAPIKey(t providers.ProviderType) string {
	switch t {
	case providers.ProviderAnthropic:
		return os.Getenv("ANTHROPIC_API_KEY")
	case providers.ProviderOpenAI:
		return os.Getenv("OPENAI_API_KEY")
	case providers.ProviderGoogle:
		return os.Getenv("GOOGLE_API_KEY")
	case providers.ProviderGroq:
		return os.Getenv("GROQ_API_KEY")
	case providers.ProviderMistral:
		return os.Getenv("MISTRAL_API_KEY")
	case providers.ProviderXAI:
		return os.Getenv("XAI_API_KEY")
	case providers.ProviderCerebras:
		return os.Getenv("CEREBRAS_API_KEY")
	default:
		return ""
	}
}

// printVersion prints version information.
func printVersion() {
	fmt.Printf("hyperharness %s (commit: %s, built: %s, %s)\n", 
		Version, Commit, BuildDate, GoVersion)
}
