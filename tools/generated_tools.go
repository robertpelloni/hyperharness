package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/robertpelloni/hyperharness/internal/subagents"
)

// registerGeneratedTools dynamically generates stub Tool definitions from extracted upstream JSON schemas.
// These act as empty parity facades, and can be manually wired up to the internal Foundation and MCP services later.
func registerGeneratedTools(r *Registry) {

	// insights extracted from claude-code/src/commands.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "insights",
		Description: "Generate a report analyzing your Claude Code sessions",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "insights"), nil
		},
	})

	// brief extracted from claude-code/src/commands/brief.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "brief",
		Description: "Toggle brief-only mode",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "brief"), nil
		},
	})

	// review extracted from claude-code/src/commands/review.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "review",
		Description: "Review a pull request",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "review"), nil
		},
	})

	// ultrareview extracted from claude-code/src/commands/review.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ultrareview",
		Description: "Auto-generated stub for ultrareview",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ultrareview"), nil
		},
	})

	// security-review extracted from claude-code/src/commands/security-review.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "security-review",
		Description: "Complete a security review of the pending changes on the current branch",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "security-review"), nil
		},
	})

	// version extracted from claude-code/src/commands/version.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "version",
		Description: "Print the version this session is running (not what autoupdate downloaded)",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "version"), nil
		},
	})

	// init extracted from claude-code/src/commands/init.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "init",
		Description: "Auto-generated stub for init",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "init"), nil
		},
	})

	// commit-push-pr extracted from claude-code/src/commands/commit-push-pr.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "commit-push-pr",
		Description: "Commit, push, and open a PR",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "commit-push-pr"), nil
		},
	})

	// init-verifiers extracted from claude-code/src/commands/init-verifiers.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "init-verifiers",
		Description: "Create verifier skill(s) for automated verification of code changes",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "init-verifiers"), nil
		},
	})

	// bridge-kick extracted from claude-code/src/commands/bridge-kick.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "bridge-kick",
		Description: "Inject bridge failure states for manual recovery testing",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "bridge-kick"), nil
		},
	})

	// advisor extracted from claude-code/src/commands/advisor.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "advisor",
		Description: "Configure the advisor model",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "advisor"), nil
		},
	})

	// commit extracted from claude-code/src/commands/commit.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "commit",
		Description: "Create a git commit",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "commit"), nil
		},
	})

	// project_areas extracted from claude-code/src/commands/insights.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "project_areas",
		Description: "Auto-generated stub for project_areas",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "project_areas"), nil
		},
	})

	// interaction_style extracted from claude-code/src/commands/insights.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "interaction_style",
		Description: "Auto-generated stub for interaction_style",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "interaction_style"), nil
		},
	})

	// what_works extracted from claude-code/src/commands/insights.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "what_works",
		Description: "Auto-generated stub for what_works",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "what_works"), nil
		},
	})

	// friction_analysis extracted from claude-code/src/commands/insights.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "friction_analysis",
		Description: "Auto-generated stub for friction_analysis",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "friction_analysis"), nil
		},
	})

	// suggestions extracted from claude-code/src/commands/insights.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "suggestions",
		Description: "Auto-generated stub for suggestions",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "suggestions"), nil
		},
	})

	// on_the_horizon extracted from claude-code/src/commands/insights.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "on_the_horizon",
		Description: "Auto-generated stub for on_the_horizon",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "on_the_horizon"), nil
		},
	})

	// cc_team_improvements extracted from claude-code/src/commands/insights.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "cc_team_improvements",
		Description: "Auto-generated stub for cc_team_improvements",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "cc_team_improvements"), nil
		},
	})

	// model_behavior_improvements extracted from claude-code/src/commands/insights.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "model_behavior_improvements",
		Description: "Auto-generated stub for model_behavior_improvements",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "model_behavior_improvements"), nil
		},
	})

	// fun_ending extracted from claude-code/src/commands/insights.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "fun_ending",
		Description: "Auto-generated stub for fun_ending",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "fun_ending"), nil
		},
	})

	// at_a_glance extracted from claude-code/src/commands/insights.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "at_a_glance",
		Description: "Auto-generated stub for at_a_glance",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "at_a_glance"), nil
		},
	})

	// tag extracted from claude-code/src/commands/tag/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tag",
		Description: "Toggle a searchable tag on the current session",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tag"), nil
		},
	})

	// passes extracted from claude-code/src/commands/passes/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "passes",
		Description: "Auto-generated stub for passes",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "passes"), nil
		},
	})

	// remote-env extracted from claude-code/src/commands/remote-env/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "remote-env",
		Description: "Configure the default remote environment for teleport sessions",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "remote-env"), nil
		},
	})

	// thinkback-play extracted from claude-code/src/commands/thinkback-play/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "thinkback-play",
		Description: "Play the thinkback animation",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "thinkback-play"), nil
		},
	})

	// usage extracted from claude-code/src/commands/usage/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "usage",
		Description: "Show plan usage limits",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "usage"), nil
		},
	})

	// rate-limit-options extracted from claude-code/src/commands/rate-limit-options/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "rate-limit-options",
		Description: "Show options when rate limit is reached",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "rate-limit-options"), nil
		},
	})

	// compact extracted from claude-code/src/commands/compact/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "compact",
		Description: "Clear conversation history but keep a summary in context. Optional: /compact [instructions for summarization]",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "compact"), nil
		},
	})

	// terminal-setup extracted from claude-code/src/commands/terminalSetup/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "terminal-setup",
		Description: "Auto-generated stub for terminal-setup",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "terminal-setup"), nil
		},
	})

	// export extracted from claude-code/src/commands/export/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "export",
		Description: "Export the current conversation to a file or clipboard",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "export"), nil
		},
	})

	// agents extracted from claude-code/src/commands/agents/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "agents",
		Description: "Manage agent configurations",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Manage agent configurations", "Configure active agents", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("agents config failed: %w", err)
			}
			return fmt.Sprintf("Agent configurations updated:\n%s", result), nil
		},
	})

	// permissions extracted from claude-code/src/commands/permissions/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "permissions",
		Description: "Manage allow & deny tool permission rules",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "permissions"), nil
		},
	})

	// branch extracted from claude-code/src/commands/branch/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "branch",
		Description: "Create a branch of the current conversation at this point",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "branch"), nil
		},
	})

	// reload-plugins extracted from claude-code/src/commands/reload-plugins/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "reload-plugins",
		Description: "Activate pending plugin changes in the current session",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "reload-plugins"), nil
		},
	})

	// mobile extracted from claude-code/src/commands/mobile/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mobile",
		Description: "Show QR code to download the Claude mobile app",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mobile"), nil
		},
	})

	// install-slack-app extracted from claude-code/src/commands/install-slack-app/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "install-slack-app",
		Description: "Install the Claude Slack app",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "install-slack-app"), nil
		},
	})

	// remote-control extracted from claude-code/src/commands/bridge/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "remote-control",
		Description: "Connect this terminal for remote-control sessions",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "remote-control"), nil
		},
	})

	// rewind extracted from claude-code/src/commands/rewind/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "rewind",
		Description: "Auto-generated stub for rewind",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "rewind"), nil
		},
	})

	// upgrade extracted from claude-code/src/commands/upgrade/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "upgrade",
		Description: "Upgrade to Max for higher rate limits and more Opus",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "upgrade"), nil
		},
	})

	// vim extracted from claude-code/src/commands/vim/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "vim",
		Description: "Toggle between Vim and Normal editing modes",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "vim"), nil
		},
	})

	// hooks extracted from claude-code/src/commands/hooks/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "hooks",
		Description: "View hook configurations for tool events",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "hooks"), nil
		},
	})

	// color extracted from claude-code/src/commands/color/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "color",
		Description: "Set the prompt bar color for this session",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "color"), nil
		},
	})

	// tasks extracted from claude-code/src/commands/tasks/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tasks",
		Description: "List and manage background tasks",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tasks"), nil
		},
	})

	// heapdump extracted from claude-code/src/commands/heapdump/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "heapdump",
		Description: "Dump the JS heap to ~/Desktop",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "heapdump"), nil
		},
	})

	// btw extracted from claude-code/src/commands/btw/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "btw",
		Description: "Ask a quick side question without interrupting the main conversation",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "btw"), nil
		},
	})

	// effort extracted from claude-code/src/commands/effort/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "effort",
		Description: "Set effort level for model usage",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "effort"), nil
		},
	})

	// release-notes extracted from claude-code/src/commands/release-notes/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "release-notes",
		Description: "View release notes",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "release-notes"), nil
		},
	})

	// diff extracted from claude-code/src/commands/diff/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "diff",
		Description: "View uncommitted changes and per-turn diffs",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "diff"), nil
		},
	})

	// mcp extracted from claude-code/src/commands/mcp/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp",
		Description: "Manage MCP servers",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp"), nil
		},
	})

	// help extracted from claude-code/src/commands/help/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "help",
		Description: "Show help and available commands",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "help"), nil
		},
	})

	// privacy-settings extracted from claude-code/src/commands/privacy-settings/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "privacy-settings",
		Description: "View and update your privacy settings",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "privacy-settings"), nil
		},
	})

	// cost extracted from claude-code/src/commands/cost/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "cost",
		Description: "Show the total cost and duration of the current session",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "cost"), nil
		},
	})

	// copy extracted from claude-code/src/commands/copy/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "copy",
		Description: "Copy Claude's last response to clipboard (or /copy N for the Nth-latest)",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "copy"), nil
		},
	})

	// stickers extracted from claude-code/src/commands/stickers/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "stickers",
		Description: "Order Claude Code stickers",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "stickers"), nil
		},
	})

	// extra-usage extracted from claude-code/src/commands/extra-usage/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "extra-usage",
		Description: "Configure extra usage to keep working when limits are hit",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "extra-usage"), nil
		},
	})

	// clear extracted from claude-code/src/commands/clear/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "clear",
		Description: "Clear conversation history and free up context",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "clear"), nil
		},
	})

	// stats extracted from claude-code/src/commands/stats/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "stats",
		Description: "Show your Claude Code usage statistics and activity",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "stats"), nil
		},
	})

	// resume extracted from claude-code/src/commands/resume/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "resume",
		Description: "Resume a previous conversation",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "resume"), nil
		},
	})

	// doctor extracted from claude-code/src/commands/doctor/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "doctor",
		Description: "Diagnose and verify your Claude Code installation and settings",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "doctor"), nil
		},
	})

	// think-back extracted from claude-code/src/commands/thinkback/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "think-back",
		Description: "Your 2025 Claude Code Year in Review",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "think-back"), nil
		},
	})

	// voice extracted from claude-code/src/commands/voice/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "voice",
		Description: "Toggle voice mode",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "voice"), nil
		},
	})

	// feedback extracted from claude-code/src/commands/feedback/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "feedback",
		Description: "Auto-generated stub for feedback",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "feedback"), nil
		},
	})

	// desktop extracted from claude-code/src/commands/desktop/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "desktop",
		Description: "Continue the current session in Claude Desktop",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "desktop"), nil
		},
	})

	// sandbox extracted from claude-code/src/commands/sandbox-toggle/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "sandbox",
		Description: "Auto-generated stub for sandbox",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "sandbox"), nil
		},
	})

	// chrome extracted from claude-code/src/commands/chrome/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "chrome",
		Description: "Claude in Chrome (Beta) settings",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "chrome"), nil
		},
	})

	// web-setup extracted from claude-code/src/commands/remote-setup/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "web-setup",
		Description: "Setup Claude Code on the web (requires connecting your GitHub account)",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "web-setup"), nil
		},
	})

	// Default extracted from claude-code/src/commands/remote-setup/api.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Default",
		Description: "Default - trusted network access",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Default"), nil
		},
	})

	// python extracted from claude-code/src/commands/remote-setup/api.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "python",
		Description: "Auto-generated stub for python",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "python"), nil
		},
	})

	// node extracted from claude-code/src/commands/remote-setup/api.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "node",
		Description: "Auto-generated stub for node",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "node"), nil
		},
	})

	// session extracted from claude-code/src/commands/session/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "session",
		Description: "Show remote session URL and QR code",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "session"), nil
		},
	})

	// logout extracted from claude-code/src/commands/logout/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "logout",
		Description: "Sign out from your Anthropic account",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "logout"), nil
		},
	})

	// exit extracted from claude-code/src/commands/exit/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "exit",
		Description: "Exit the REPL",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "exit"), nil
		},
	})

	// model extracted from claude-code/src/commands/model/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "model",
		Description: "Auto-generated stub for model",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "model"), nil
		},
	})

	// output-style extracted from claude-code/src/commands/output-style/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "output-style",
		Description: "Deprecated: use /config to change output style",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "output-style"), nil
		},
	})

	// memory extracted from claude-code/src/commands/memory/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "memory",
		Description: "Edit Claude memory files",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "memory"), nil
		},
	})

	// pr-comments extracted from claude-code/src/commands/pr_comments/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "pr-comments",
		Description: "Get comments from a GitHub pull request",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "pr-comments"), nil
		},
	})

	// install-github-app extracted from claude-code/src/commands/install-github-app/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "install-github-app",
		Description: "Set up Claude GitHub Actions for a repository",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "install-github-app"), nil
		},
	})

	// keybindings extracted from claude-code/src/commands/keybindings/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "keybindings",
		Description: "Open or create your keybindings configuration file",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "keybindings"), nil
		},
	})

	// skills extracted from claude-code/src/commands/skills/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "skills",
		Description: "List available skills",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "skills"), nil
		},
	})

	// fast extracted from claude-code/src/commands/fast/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "fast",
		Description: "Auto-generated stub for fast",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "fast"), nil
		},
	})

	// status extracted from claude-code/src/commands/status/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "status",
		Description: "Show Claude Code status including version, model, account, API connectivity, and tool statuses",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "status"), nil
		},
	})

	// login extracted from claude-code/src/commands/login/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "login",
		Description: "Auto-generated stub for login",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "login"), nil
		},
	})

	// context extracted from claude-code/src/commands/context/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "context",
		Description: "Visualize current context usage as a colored grid",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "context"), nil
		},
	})

	// theme extracted from claude-code/src/commands/theme/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "theme",
		Description: "Change the theme",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "theme"), nil
		},
	})

	// ide extracted from claude-code/src/commands/ide/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ide",
		Description: "Manage IDE integrations and show status",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ide"), nil
		},
	})

	// rename extracted from claude-code/src/commands/rename/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "rename",
		Description: "Rename the current conversation",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "rename"), nil
		},
	})

	// add-dir extracted from claude-code/src/commands/add-dir/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "add-dir",
		Description: "Add a new working directory",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "add-dir"), nil
		},
	})

	// plan extracted from claude-code/src/commands/plan/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "plan",
		Description: "Enable plan mode or view the current session plan",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "plan"), nil
		},
	})

	// files extracted from claude-code/src/commands/files/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "files",
		Description: "List all files currently in context",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "files"), nil
		},
	})

	// Unknown extracted from claude-code/src/services/PromptSuggestion/speculation.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Unknown",
		Description: "Auto-generated stub for Unknown",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Unknown"), nil
		},
	})

	// unknown extracted from claude-code/src/services/api/promptCacheBreakDetection.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "unknown",
		Description: "Auto-generated stub for unknown",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "unknown"), nil
		},
	})

	// claude-code extracted from claude-code/src/services/mcp/client.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "claude-code",
		Description: "Anthropic's agentic coding tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "claude-code"), nil
		},
	})

	// mcp__ extracted from claude-code/src/services/mcp/client.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp__",
		Description: "Auto-generated stub for mcp__",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp__"), nil
		},
	})

	// none extracted from claude-code/src/components/permissions/utils.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "none",
		Description: "Auto-generated stub for none",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "none"), nil
		},
	})

	// red extracted from claude-code/src/ink/termio.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "red",
		Description: "Auto-generated stub for red",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "red"), nil
		},
	})

	// displayText extracted from claude-code/src/hooks/unifiedSuggestions.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "displayText",
		Description: "Auto-generated stub for displayText",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "displayText"), nil
		},
	})

	// name extracted from claude-code/src/hooks/unifiedSuggestions.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "name",
		Description: "Auto-generated stub for name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "name"), nil
		},
	})

	// server extracted from claude-code/src/hooks/unifiedSuggestions.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "server",
		Description: "Auto-generated stub for server",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "server"), nil
		},
	})

	// description extracted from claude-code/src/hooks/unifiedSuggestions.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "description",
		Description: "Auto-generated stub for description",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "description"), nil
		},
	})

	// agentType extracted from claude-code/src/hooks/unifiedSuggestions.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "agentType",
		Description: "Auto-generated stub for agentType",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "agentType"), nil
		},
	})

	// Skills extracted from claude-code/src/utils/analyzeContext.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Skills",
		Description: "Auto-generated stub for Skills",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Skills"), nil
		},
	})

	// Messages extracted from claude-code/src/utils/analyzeContext.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Messages",
		Description: "Auto-generated stub for Messages",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Messages"), nil
		},
	})

	// claude-local extracted from claude-code/src/utils/localInstaller.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "claude-local",
		Description: "Auto-generated stub for claude-local",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "claude-local"), nil
		},
	})

	// Microcompact extracted from claude-code/src/utils/queryProfiler.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Microcompact",
		Description: "Auto-generated stub for Microcompact",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Microcompact"), nil
		},
	})

	// Autocompact extracted from claude-code/src/utils/queryProfiler.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Autocompact",
		Description: "Auto-generated stub for Autocompact",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Autocompact"), nil
		},
	})

	// zsh extracted from claude-code/src/utils/completionCache.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "zsh",
		Description: "Auto-generated stub for zsh",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "zsh"), nil
		},
	})

	// bash extracted from claude-code/src/utils/completionCache.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "bash",
		Description: "Auto-generated stub for bash",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "bash"), nil
		},
	})

	// fish extracted from claude-code/src/utils/completionCache.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "fish",
		Description: "Auto-generated stub for fish",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "fish"), nil
		},
	})

	// x extracted from claude-code/src/utils/sideQuery.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "x",
		Description: "Auto-generated stub for x",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "x"), nil
		},
	})

	// PreToolUse extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "PreToolUse",
		Description: "Auto-generated stub for PreToolUse",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "PreToolUse"), nil
		},
	})

	// PostToolUse extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "PostToolUse",
		Description: "Auto-generated stub for PostToolUse",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "PostToolUse"), nil
		},
	})

	// PostToolUseFailure extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "PostToolUseFailure",
		Description: "Auto-generated stub for PostToolUseFailure",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "PostToolUseFailure"), nil
		},
	})

	// PermissionDenied extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "PermissionDenied",
		Description: "Auto-generated stub for PermissionDenied",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "PermissionDenied"), nil
		},
	})

	// Notification extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Notification",
		Description: "Auto-generated stub for Notification",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Notification"), nil
		},
	})

	// StopFailure extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "StopFailure",
		Description: "Auto-generated stub for StopFailure",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "StopFailure"), nil
		},
	})

	// SubagentStop extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "SubagentStop",
		Description: "Auto-generated stub for SubagentStop",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "SubagentStop", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "SubagentStop", result), nil
		},
	})

	// Stop extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Stop",
		Description: "Auto-generated stub for Stop",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Stop"), nil
		},
	})

	// TeammateIdle extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "TeammateIdle",
		Description: "Auto-generated stub for TeammateIdle",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "TeammateIdle"), nil
		},
	})

	// TaskCreated extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "TaskCreated",
		Description: "Auto-generated stub for TaskCreated",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "TaskCreated"), nil
		},
	})

	// TaskCompleted extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "TaskCompleted",
		Description: "Auto-generated stub for TaskCompleted",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "TaskCompleted"), nil
		},
	})

	// UserPromptSubmit extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "UserPromptSubmit",
		Description: "Auto-generated stub for UserPromptSubmit",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "UserPromptSubmit"), nil
		},
	})

	// SessionStart extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "SessionStart",
		Description: "Auto-generated stub for SessionStart",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "SessionStart"), nil
		},
	})

	// Setup extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Setup",
		Description: "Auto-generated stub for Setup",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Setup"), nil
		},
	})

	// SubagentStart extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "SubagentStart",
		Description: "Auto-generated stub for SubagentStart",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "SubagentStart", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "SubagentStart", result), nil
		},
	})

	// PreCompact extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "PreCompact",
		Description: "Auto-generated stub for PreCompact",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "PreCompact"), nil
		},
	})

	// PostCompact extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "PostCompact",
		Description: "Auto-generated stub for PostCompact",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "PostCompact"), nil
		},
	})

	// SessionEnd extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "SessionEnd",
		Description: "Auto-generated stub for SessionEnd",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "SessionEnd"), nil
		},
	})

	// PermissionRequest extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "PermissionRequest",
		Description: "Auto-generated stub for PermissionRequest",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "PermissionRequest"), nil
		},
	})

	// ConfigChange extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ConfigChange",
		Description: "Auto-generated stub for ConfigChange",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ConfigChange"), nil
		},
	})

	// CwdChanged extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "CwdChanged",
		Description: "Auto-generated stub for CwdChanged",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "CwdChanged"), nil
		},
	})

	// FileChanged extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "FileChanged",
		Description: "Auto-generated stub for FileChanged",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "FileChanged"), nil
		},
	})

	// InstructionsLoaded extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "InstructionsLoaded",
		Description: "Auto-generated stub for InstructionsLoaded",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "InstructionsLoaded"), nil
		},
	})

	// Elicitation extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Elicitation",
		Description: "Auto-generated stub for Elicitation",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Elicitation"), nil
		},
	})

	// ElicitationResult extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ElicitationResult",
		Description: "Auto-generated stub for ElicitationResult",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ElicitationResult"), nil
		},
	})

	// WorktreeCreate extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "WorktreeCreate",
		Description: "Auto-generated stub for WorktreeCreate",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "WorktreeCreate"), nil
		},
	})

	// WorktreeRemove extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "WorktreeRemove",
		Description: "Auto-generated stub for WorktreeRemove",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "WorktreeRemove"), nil
		},
	})

	// function extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "function",
		Description: "Auto-generated stub for function",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "function"), nil
		},
	})

	// callback extracted from claude-code/src/utils/hooks.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "callback",
		Description: "Auto-generated stub for callback",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "callback"), nil
		},
	})

	// explain_command extracted from claude-code/src/utils/permissions/permissionExplainer.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "explain_command",
		Description: "Provide an explanation of a shell command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "explain_command"), nil
		},
	})

	// skill_improvement extracted from claude-code/src/utils/hooks/skillImprovement.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "skill_improvement",
		Description: "Auto-generated stub for skill_improvement",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "skill_improvement"), nil
		},
	})

	// settings extracted from claude-code/src/utils/plugins/marketplaceManager.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "settings",
		Description: "Auto-generated stub for settings",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "settings"), nil
		},
	})

	// HEAD extracted from claude-code/src/utils/git/gitFilesystem.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "HEAD",
		Description: "Auto-generated stub for HEAD",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "HEAD"), nil
		},
	})

	// claude-native-installer extracted from claude-code/src/utils/nativeInstaller/download.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "claude-native-installer",
		Description: "Auto-generated stub for claude-native-installer",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "claude-native-installer"), nil
		},
	})

	// plaintext extracted from claude-code/src/utils/secureStorage/plainTextStorage.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "plaintext",
		Description: "Auto-generated stub for plaintext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "plaintext"), nil
		},
	})

	// keychain extracted from claude-code/src/utils/secureStorage/macOsKeychainStorage.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "keychain",
		Description: "Auto-generated stub for keychain",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "keychain"), nil
		},
	})

	// commandName extracted from claude-code/src/utils/suggestions/commandSuggestions.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "commandName",
		Description: "Auto-generated stub for commandName",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "commandName"), nil
		},
	})

	// partKey extracted from claude-code/src/utils/suggestions/commandSuggestions.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "partKey",
		Description: "Auto-generated stub for partKey",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "partKey"), nil
		},
	})

	// aliasKey extracted from claude-code/src/utils/suggestions/commandSuggestions.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "aliasKey",
		Description: "Auto-generated stub for aliasKey",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "aliasKey"), nil
		},
	})

	// descriptionKey extracted from claude-code/src/utils/suggestions/commandSuggestions.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "descriptionKey",
		Description: "Auto-generated stub for descriptionKey",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "descriptionKey"), nil
		},
	})

	// workflowName extracted from claude-code/src/utils/task/framework.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "workflowName",
		Description: "Auto-generated stub for workflowName",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "workflowName"), nil
		},
	})

	// trace_truncated extracted from claude-code/src/utils/telemetry/perfettoTracing.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "trace_truncated",
		Description: "Auto-generated stub for trace_truncated",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "trace_truncated"), nil
		},
	})

	// process_name extracted from claude-code/src/utils/telemetry/perfettoTracing.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "process_name",
		Description: "Auto-generated stub for process_name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "process_name"), nil
		},
	})

	// thread_name extracted from claude-code/src/utils/telemetry/perfettoTracing.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "thread_name",
		Description: "Auto-generated stub for thread_name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "thread_name"), nil
		},
	})

	// parent_agent extracted from claude-code/src/utils/telemetry/perfettoTracing.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "parent_agent",
		Description: "Auto-generated stub for parent_agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "parent_agent"), nil
		},
	})

	// Sampling extracted from claude-code/src/utils/telemetry/perfettoTracing.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Sampling",
		Description: "Auto-generated stub for Sampling",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Sampling"), nil
		},
	})

	// Interaction extracted from claude-code/src/utils/telemetry/perfettoTracing.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Interaction",
		Description: "Auto-generated stub for Interaction",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Interaction"), nil
		},
	})

	// third-party extracted from claude-code/src/utils/telemetry/pluginTelemetry.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "third-party",
		Description: "Auto-generated stub for third-party",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "third-party"), nil
		},
	})

	// Brave extracted from claude-code/src/utils/claudeInChrome/common.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Brave",
		Description: "Auto-generated stub for Brave",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Brave"), nil
		},
	})

	// Arc extracted from claude-code/src/utils/claudeInChrome/common.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Arc",
		Description: "Auto-generated stub for Arc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Arc"), nil
		},
	})

	// Chromium extracted from claude-code/src/utils/claudeInChrome/common.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Chromium",
		Description: "Auto-generated stub for Chromium",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Chromium"), nil
		},
	})

	// Vivaldi extracted from claude-code/src/utils/claudeInChrome/common.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Vivaldi",
		Description: "Auto-generated stub for Vivaldi",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Vivaldi"), nil
		},
	})

	// Opera extracted from claude-code/src/utils/claudeInChrome/common.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Opera",
		Description: "Auto-generated stub for Opera",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Opera"), nil
		},
	})

	// iTerm2 extracted from claude-code/src/utils/deepLink/terminalLauncher.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "iTerm2",
		Description: "Auto-generated stub for iTerm2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "iTerm2"), nil
		},
	})

	// Ghostty extracted from claude-code/src/utils/deepLink/terminalLauncher.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Ghostty",
		Description: "Auto-generated stub for Ghostty",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Ghostty"), nil
		},
	})

	// Kitty extracted from claude-code/src/utils/deepLink/terminalLauncher.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Kitty",
		Description: "Auto-generated stub for Kitty",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Kitty"), nil
		},
	})

	// Alacritty extracted from claude-code/src/utils/deepLink/terminalLauncher.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Alacritty",
		Description: "Auto-generated stub for Alacritty",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Alacritty"), nil
		},
	})

	// WezTerm extracted from claude-code/src/utils/deepLink/terminalLauncher.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "WezTerm",
		Description: "Auto-generated stub for WezTerm",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "WezTerm"), nil
		},
	})

	// x-terminal-emulator extracted from claude-code/src/utils/deepLink/terminalLauncher.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "x-terminal-emulator",
		Description: "Auto-generated stub for x-terminal-emulator",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "x-terminal-emulator"), nil
		},
	})

	// PowerShell extracted from claude-code/src/utils/deepLink/terminalLauncher.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "PowerShell",
		Description: "Auto-generated stub for PowerShell",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "PowerShell"), nil
		},
	})

	// pyright extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "pyright",
		Description: "Type checker for Python",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "pyright"), nil
		},
	})

	// --version extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--version",
		Description: "Show help message",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--version"), nil
		},
	})

	// - extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "-",
		Description: "Read file or directory list from stdin",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "-"), nil
		},
	})

	// --createstub extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--createstub",
		Description: "Read file or directory list from stdin",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--createstub"), nil
		},
	})

	// IMPORT extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "IMPORT",
		Description: "Create type stub file(s) for import",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "IMPORT"), nil
		},
	})

	// DIRECTORY extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "DIRECTORY",
		Description: "Use typeshed type stubs at this location",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "DIRECTORY"), nil
		},
	})

	// --verifytypes extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--verifytypes",
		Description: "Verify completeness of types in py.typed package",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--verifytypes"), nil
		},
	})

	// --ignoreexternal extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--ignoreexternal",
		Description: "Ignore external imports for --verifytypes",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--ignoreexternal"), nil
		},
	})

	// --pythonpath extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--pythonpath",
		Description: "Ignore external imports for --verifytypes",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--pythonpath"), nil
		},
	})

	// FILE extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "FILE",
		Description: "Path to the Python interpreter",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "FILE"), nil
		},
	})

	// --pythonplatform extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--pythonplatform",
		Description: "Path to the Python interpreter",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--pythonplatform"), nil
		},
	})

	// PLATFORM extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "PLATFORM",
		Description: "Analyze for platform",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "PLATFORM"), nil
		},
	})

	// --pythonversion extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--pythonversion",
		Description: "Analyze for platform",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--pythonversion"), nil
		},
	})

	// VERSION extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "VERSION",
		Description: "Analyze for Python version",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "VERSION"), nil
		},
	})

	// --outputjson extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--outputjson",
		Description: "Output results in JSON format",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--outputjson"), nil
		},
	})

	// --verbose extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--verbose",
		Description: "Output results in JSON format",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--verbose"), nil
		},
	})

	// --stats extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--stats",
		Description: "Emit verbose diagnostics",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--stats"), nil
		},
	})

	// --dependencies extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--dependencies",
		Description: "Print detailed performance stats",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--dependencies"), nil
		},
	})

	// --level extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--level",
		Description: "Emit import dependency information",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--level"), nil
		},
	})

	// LEVEL extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "LEVEL",
		Description: "Minimum diagnostic level",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "LEVEL"), nil
		},
	})

	// --skipunannotated extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--skipunannotated",
		Description: "Minimum diagnostic level",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--skipunannotated"), nil
		},
	})

	// --warnings extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--warnings",
		Description: "Skip type analysis of unannotated functions",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--warnings"), nil
		},
	})

	// --threads extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "--threads",
		Description: "Use exit code of 1 if warnings are reported",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "--threads"), nil
		},
	})

	// N extracted from claude-code/src/utils/bash/specs/pyright.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "N",
		Description: "Use up to N threads to parallelize type checking",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "N"), nil
		},
	})

	// sleep extracted from claude-code/src/utils/bash/specs/sleep.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "sleep",
		Description: "Delay for a specified amount of time",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "sleep"), nil
		},
	})

	// duration extracted from claude-code/src/utils/bash/specs/sleep.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "duration",
		Description: "Delay for a specified amount of time",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "duration"), nil
		},
	})

	// srun extracted from claude-code/src/utils/bash/specs/srun.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "srun",
		Description: "Run a command on SLURM cluster nodes",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "srun"), nil
		},
	})

	// count extracted from claude-code/src/utils/bash/specs/srun.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "count",
		Description: "Number of tasks",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "count"), nil
		},
	})

	// command extracted from claude-code/src/utils/bash/specs/srun.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "command",
		Description: "Number of nodes to allocate",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "command"), nil
		},
	})

	// timeout extracted from claude-code/src/utils/bash/specs/timeout.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "timeout",
		Description: "Run a command with a time limit",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "timeout"), nil
		},
	})

	// nohup extracted from claude-code/src/utils/bash/specs/nohup.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "nohup",
		Description: "Run a command immune to hangups",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "nohup"), nil
		},
	})

	// alias extracted from claude-code/src/utils/bash/specs/alias.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "alias",
		Description: "Create or list command aliases",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "alias"), nil
		},
	})

	// definition extracted from claude-code/src/utils/bash/specs/alias.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "definition",
		Description: "Create or list command aliases",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "definition"), nil
		},
	})

	// time extracted from claude-code/src/utils/bash/specs/time.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "time",
		Description: "Time a command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "time"), nil
		},
	})

	// DreamTask extracted from claude-code/src/tasks/DreamTask/DreamTask.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "DreamTask",
		Description: "Auto-generated stub for DreamTask",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "DreamTask"), nil
		},
	})

	// verify extracted from claude-code/src/skills/bundled/verify.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "verify",
		Description: "Auto-generated stub for verify",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "verify"), nil
		},
	})

	// claude-in-chrome extracted from claude-code/src/skills/bundled/claudeInChrome.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "claude-in-chrome",
		Description: "Automates your Chrome browser to interact with web pages - clicking elements, filling forms, capturing screenshots, reading console logs, and navigating sites. Opens pages in new tabs within your existing Chrome session. Requires site-level permissions before executing (configured in the extension).",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "claude-in-chrome"), nil
		},
	})

	// keybindings-help extracted from claude-code/src/skills/bundled/keybindings.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "keybindings-help",
		Description: "Use when the user wants to customize keyboard shortcuts, rebind keys, add chord bindings, or modify ~/.claude/keybindings.json. Examples: \"rebind ctrl+s\", \"add a chord shortcut\", \"change the submit key\", \"customize keybindings\".",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "keybindings-help"), nil
		},
	})

	// loop extracted from claude-code/src/skills/bundled/loop.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "loop",
		Description: "Run a prompt or slash command on a recurring interval (e.g. /loop 5m /foo, defaults to 10m)",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "loop"), nil
		},
	})

	// skillify extracted from claude-code/src/skills/bundled/skillify.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "skillify",
		Description: "Capture this session's repeatable process into a skill. Call at end of the process you want to capture with an optional description.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "skillify"), nil
		},
	})

	// simplify extracted from claude-code/src/skills/bundled/simplify.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "simplify",
		Description: "Review changed code for reuse, quality, and efficiency, then fix any issues found.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "simplify"), nil
		},
	})

	// lorem-ipsum extracted from claude-code/src/skills/bundled/loremIpsum.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "lorem-ipsum",
		Description: "Generate filler text for long context testing. Specify token count as argument (e.g., /lorem-ipsum 50000). Outputs approximately the requested number of tokens. Ant-only.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "lorem-ipsum"), nil
		},
	})

	// batch extracted from claude-code/src/skills/bundled/batch.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "batch",
		Description: "Research and plan a large-scale change, then execute it in parallel across 5–30 isolated worktree agents that each open a PR.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "batch"), nil
		},
	})

	// stuck extracted from claude-code/src/skills/bundled/stuck.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "stuck",
		Description: "[ANT-ONLY] Investigate frozen/stuck/slow Claude Code sessions on this machine and post a diagnostic report to #claude-code-feedback.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "stuck"), nil
		},
	})

	// update-config extracted from claude-code/src/skills/bundled/updateConfig.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "update-config",
		Description: "Auto-generated stub for update-config",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "update-config"), nil
		},
	})

	// debug extracted from claude-code/src/skills/bundled/debug.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "debug",
		Description: "Auto-generated stub for debug",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "debug"), nil
		},
	})

	// claude-api extracted from claude-code/src/skills/bundled/claudeApi.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "claude-api",
		Description: "Build apps with the Claude API or Anthropic SDK.\n",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "claude-api"), nil
		},
	})

	// schedule extracted from claude-code/src/skills/bundled/scheduleRemoteAgents.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "schedule",
		Description: "Create, update, list, or run scheduled remote agents (triggers) that execute on a cron schedule.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "schedule"), nil
		},
	})

	// remember extracted from claude-code/src/skills/bundled/remember.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "remember",
		Description: "Review auto-memory entries and propose promotions to CLAUDE.md, CLAUDE.local.md, or shared memory. Also detects outdated, conflicting, and duplicate entries across memory layers.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "remember"), nil
		},
	})

	// ship-audit extracted from claude-code/src/tools/AgentTool/prompt.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ship-audit",
		Description: "Branch ship-readiness audit",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ship-audit"), nil
		},
	})

	// migration-review extracted from claude-code/src/tools/AgentTool/prompt.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "migration-review",
		Description: "Independent migration review",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "migration-review"), nil
		},
	})

	// web_search extracted from claude-code/src/tools/WebSearchTool/WebSearchTool.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "web_search",
		Description: "Auto-generated stub for web_search",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "web_search"), nil
		},
	})

	// in-process extracted from claude-code/src/tools/shared/spawnMultiAgent.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "in-process",
		Description: "Auto-generated stub for in-process",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "in-process"), nil
		},
	})

	// ReadMcpResourceTool extracted from claude-code/src/tools/ReadMcpResourceTool/ReadMcpResourceTool.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ReadMcpResourceTool",
		Description: "Auto-generated stub for ReadMcpResourceTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ReadMcpResourceTool"), nil
		},
	})

	// remote_skill extracted from claude-code/src/tools/SkillTool/SkillTool.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "remote_skill",
		Description: "Auto-generated stub for remote_skill",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "remote_skill"), nil
		},
	})

	// Explanatory extracted from claude-code/src/constants/outputStyles.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Explanatory",
		Description: "Claude explains its implementation choices and codebase patterns",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Explanatory"), nil
		},
	})

	// Learning extracted from claude-code/src/constants/outputStyles.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Learning",
		Description: "Claude pauses and asks you to write small pieces of code for hands-on practice",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Learning"), nil
		},
	})

	// opencode extracted from opencode/sst.config.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "opencode",
		Description: "Auto-generated stub for opencode",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "opencode"), nil
		},
	})

	// SYNC_SERVER extracted from opencode/infra/app.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "SYNC_SERVER",
		Description: "Auto-generated stub for SYNC_SERVER",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "SYNC_SERVER"), nil
		},
	})

	// production extracted from opencode/infra/console.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "production",
		Description: "Auto-generated stub for production",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "production"), nil
		},
	})

	// f extracted from opencode/packages/opencode/test/keybind.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "f",
		Description: "Auto-generated stub for f",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "f"), nil
		},
	})

	// g extracted from opencode/packages/opencode/test/keybind.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "g",
		Description: "Auto-generated stub for g",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "g"), nil
		},
	})

	// h extracted from opencode/packages/opencode/test/keybind.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "h",
		Description: "Auto-generated stub for h",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "h"), nil
		},
	})

	// return extracted from opencode/packages/opencode/test/keybind.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "return",
		Description: "Auto-generated stub for return",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "return"), nil
		},
	})

	// f2 extracted from opencode/packages/opencode/test/keybind.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "f2",
		Description: "Auto-generated stub for f2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "f2"), nil
		},
	})

	// pgup extracted from opencode/packages/opencode/test/keybind.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "pgup",
		Description: "Auto-generated stub for pgup",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "pgup"), nil
		},
	})

	// z extracted from opencode/packages/opencode/test/keybind.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "z",
		Description: "Auto-generated stub for z",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "z"), nil
		},
	})

	// a extracted from opencode/packages/opencode/test/keybind.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "a",
		Description: "Auto-generated stub for a",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "a"), nil
		},
	})

	// c extracted from opencode/packages/opencode/test/keybind.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "c",
		Description: "Auto-generated stub for c",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "c"), nil
		},
	})

	// y extracted from opencode/packages/opencode/test/keybind.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "y",
		Description: "Auto-generated stub for y",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "y"), nil
		},
	})

	// u extracted from opencode/packages/opencode/test/keybind.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "u",
		Description: "Auto-generated stub for u",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "u"), nil
		},
	})

	// q extracted from opencode/packages/opencode/test/keybind.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "q",
		Description: "Auto-generated stub for q",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "q"), nil
		},
	})

	// j extracted from opencode/packages/opencode/test/keybind.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "j",
		Description: "Auto-generated stub for j",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "j"), nil
		},
	})

	// testuser extracted from opencode/packages/opencode/test/config/config.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "testuser",
		Description: "Auto-generated stub for testuser",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "testuser"), nil
		},
	})

	// base extracted from opencode/packages/opencode/test/config/config.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "base",
		Description: "Auto-generated stub for base",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "base"), nil
		},
	})

	// helper extracted from opencode/packages/opencode/test/config/config.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "helper",
		Description: "Auto-generated stub for helper",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "helper"), nil
		},
	})

	// config-fixture extracted from opencode/packages/opencode/test/config/config.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "config-fixture",
		Description: "Auto-generated stub for config-fixture",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "config-fixture"), nil
		},
	})

	// demo-plugin extracted from opencode/packages/opencode/test/config/config.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "demo-plugin",
		Description: "Auto-generated stub for demo-plugin",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "demo-plugin"), nil
		},
	})

	// project-user extracted from opencode/packages/opencode/test/config/config.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "project-user",
		Description: "Auto-generated stub for project-user",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "project-user"), nil
		},
	})

	// test_tool extracted from opencode/packages/opencode/test/mcp/lifecycle.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test_tool",
		Description: "A test tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test_tool"), nil
		},
	})

	// do_thing extracted from opencode/packages/opencode/test/mcp/lifecycle.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "do_thing",
		Description: "does a thing",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "do_thing"), nil
		},
	})

	// next_tool extracted from opencode/packages/opencode/test/mcp/lifecycle.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "next_tool",
		Description: "next",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "next_tool"), nil
		},
	})

	// my_tool extracted from opencode/packages/opencode/test/mcp/lifecycle.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my_tool",
		Description: "a tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my_tool"), nil
		},
	})

	// good_tool extracted from opencode/packages/opencode/test/mcp/lifecycle.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "good_tool",
		Description: "works",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "good_tool"), nil
		},
	})

	// my-prompt extracted from opencode/packages/opencode/test/mcp/lifecycle.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-prompt",
		Description: "A test prompt",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-prompt"), nil
		},
	})

	// my-resource extracted from opencode/packages/opencode/test/mcp/lifecycle.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-resource",
		Description: "A test resource",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-resource"), nil
		},
	})

	// hidden-prompt extracted from opencode/packages/opencode/test/mcp/lifecycle.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "hidden-prompt",
		Description: "Should not appear",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "hidden-prompt"), nil
		},
	})

	// tool-a extracted from opencode/packages/opencode/test/mcp/lifecycle.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tool-a",
		Description: "Tool A",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tool-a"), nil
		},
	})

	// One extracted from opencode/packages/opencode/test/cli/account.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "One",
		Description: "Auto-generated stub for One",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "One"), nil
		},
	})

	// portkey extracted from opencode/packages/opencode/test/cli/plugin-auth-picker.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "portkey",
		Description: "Auto-generated stub for portkey",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "portkey"), nil
		},
	})

	// acme-plugin extracted from opencode/packages/opencode/test/cli/tui/plugin-loader-entrypoint.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "acme-plugin",
		Description: "Auto-generated stub for acme-plugin",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "acme-plugin"), nil
		},
	})

	// dir-plugin extracted from opencode/packages/opencode/test/cli/tui/plugin-loader-entrypoint.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "dir-plugin",
		Description: "Auto-generated stub for dir-plugin",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "dir-plugin"), nil
		},
	})

	// demo-install-plugin extracted from opencode/packages/opencode/test/cli/tui/plugin-install.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "demo-install-plugin",
		Description: "Auto-generated stub for demo-install-plugin",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "demo-install-plugin"), nil
		},
	})

	// github extracted from opencode/packages/opencode/test/cli/tui/plugin-loader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "github",
		Description: "Auto-generated stub for github",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "github"), nil
		},
	})

	// Anthropic extracted from opencode/packages/opencode/test/cli/tui/transcript.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Anthropic",
		Description: "Auto-generated stub for Anthropic",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Anthropic"), nil
		},
	})

	// acme extracted from opencode/packages/opencode/test/plugin/install.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "acme",
		Description: "Auto-generated stub for acme",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "acme"), nil
		},
	})

	// scope-plugin extracted from opencode/packages/opencode/test/plugin/loader-shared.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "scope-plugin",
		Description: "Auto-generated stub for scope-plugin",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "scope-plugin"), nil
		},
	})

	// GPT-4o extracted from opencode/packages/opencode/test/plugin/github-copilot-models.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "GPT-4o",
		Description: "Auto-generated stub for GPT-4o",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "GPT-4o"), nil
		},
	})

	// plug extracted from opencode/packages/opencode/test/plugin/workspace-adaptor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "plug",
		Description: "plugin workspace adaptor",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "plug"), nil
		},
	})

	// typescript extracted from opencode/packages/opencode/test/util/module.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "typescript",
		Description: "Auto-generated stub for typescript",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "typescript"), nil
		},
	})

	// eslint extracted from opencode/packages/opencode/test/util/module.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "eslint",
		Description: "Auto-generated stub for eslint",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "eslint"), nil
		},
	})

	// biome extracted from opencode/packages/opencode/test/util/module.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "biome",
		Description: "Auto-generated stub for biome",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "biome"), nil
		},
	})

	// Test extracted from opencode/packages/opencode/test/session/processor-effect.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Test",
		Description: "Auto-generated stub for Test",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Test"), nil
		},
	})

	// build extracted from opencode/packages/opencode/test/session/processor-effect.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "build",
		Description: "Auto-generated stub for build",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "build"), nil
		},
	})

	// other extracted from opencode/packages/opencode/test/session/structured-output.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "other",
		Description: "Auto-generated stub for other",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "other"), nil
		},
	})

	// John extracted from opencode/packages/opencode/test/session/structured-output.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "John",
		Description: "Auto-generated stub for John",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "John"), nil
		},
	})

	// dir extracted from opencode/packages/opencode/test/session/message-v2.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "dir",
		Description: "desc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "dir"), nil
		},
	})

	// ContextOverflowError extracted from opencode/packages/opencode/test/session/message-v2.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ContextOverflowError",
		Description: "Auto-generated stub for ContextOverflowError",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ContextOverflowError"), nil
		},
	})

	// APIError extracted from opencode/packages/opencode/test/session/message-v2.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "APIError",
		Description: "Auto-generated stub for APIError",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "APIError"), nil
		},
	})

	// UnknownError extracted from opencode/packages/opencode/test/session/message-v2.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "UnknownError",
		Description: "Auto-generated stub for UnknownError",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "UnknownError"), nil
		},
	})

	// OpenAI extracted from opencode/packages/opencode/test/session/llm.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "OpenAI",
		Description: "Auto-generated stub for OpenAI",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "OpenAI"), nil
		},
	})

	// big-pickle extracted from opencode/packages/opencode/test/acp/event-subscription.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "big-pickle",
		Description: "build",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "big-pickle"), nil
		},
	})

	// home extracted from opencode/packages/opencode/test/fixture/tui-plugin.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "home",
		Description: "Auto-generated stub for home",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "home"), nil
		},
	})

	// test-workspace extracted from opencode/packages/opencode/test/project/worktree.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-workspace",
		Description: "Auto-generated stub for test-workspace",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-workspace"), nil
		},
	})

	// first extracted from opencode/packages/opencode/test/sync/index.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "first",
		Description: "Auto-generated stub for first",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "first"), nil
		},
	})

	// second extracted from opencode/packages/opencode/test/sync/index.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "second",
		Description: "Auto-generated stub for second",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "second"), nil
		},
	})

	// replayed extracted from opencode/packages/opencode/test/sync/index.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "replayed",
		Description: "Auto-generated stub for replayed",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "replayed"), nil
		},
	})

	// bad extracted from opencode/packages/opencode/test/sync/index.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "bad",
		Description: "Auto-generated stub for bad",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "bad"), nil
		},
	})

	// GPT-4 extracted from opencode/packages/opencode/test/provider/provider.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "GPT-4",
		Description: "Auto-generated stub for GPT-4",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "GPT-4"), nil
		},
	})

	// Random extracted from opencode/packages/opencode/test/provider/provider.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Random",
		Description: "Auto-generated stub for Random",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Random"), nil
		},
	})

	// Other extracted from opencode/packages/opencode/test/provider/provider.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Other",
		Description: "Auto-generated stub for Other",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Other"), nil
		},
	})

	// Model extracted from opencode/packages/opencode/test/provider/provider.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Model",
		Description: "Auto-generated stub for Model",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Model"), nil
		},
	})

	// GPT-5 extracted from opencode/packages/opencode/test/provider/transform.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "GPT-5",
		Description: "Auto-generated stub for GPT-5",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "GPT-5"), nil
		},
	})

	// get_weather extracted from opencode/packages/opencode/test/provider/copilot/copilot-chat-model.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "get_weather",
		Description: "Get the weather for a location",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "get_weather"), nil
		},
	})

	// calculator extracted from opencode/packages/opencode/test/provider/copilot/convert-to-copilot-messages.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "calculator",
		Description: "Auto-generated stub for calculator",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "calculator"), nil
		},
	})

	// thwomp extracted from opencode/packages/opencode/test/provider/copilot/convert-to-copilot-messages.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "thwomp",
		Description: "Auto-generated stub for thwomp",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "thwomp"), nil
		},
	})

	// searchTool extracted from opencode/packages/opencode/test/provider/copilot/convert-to-copilot-messages.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "searchTool",
		Description: "Auto-generated stub for searchTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "searchTool"), nil
		},
	})

	// mapsTool extracted from opencode/packages/opencode/test/provider/copilot/convert-to-copilot-messages.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mapsTool",
		Description: "Auto-generated stub for mapsTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mapsTool"), nil
		},
	})

	// workspace-test extracted from opencode/packages/opencode/test/control-plane/adaptors.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "workspace-test",
		Description: "Auto-generated stub for workspace-test",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "workspace-test"), nil
		},
	})

	// Builder extracted from opencode/packages/opencode/test/agent/agent.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Builder",
		Description: "Auto-generated stub for Builder",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Builder"), nil
		},
	})

	// tool-skill extracted from opencode/packages/opencode/test/tool/skill.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tool-skill",
		Description: "Auto-generated stub for tool-skill",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tool-skill"), nil
		},
	})

	// custom-tools extracted from opencode/packages/opencode/test/tool/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "custom-tools",
		Description: "Auto-generated stub for custom-tools",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "custom-tools"), nil
		},
	})

	// cowsay extracted from opencode/packages/opencode/test/tool/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "cowsay",
		Description: "Auto-generated stub for cowsay",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "cowsay"), nil
		},
	})

	// npm extracted from opencode/packages/opencode/src/installation/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "npm",
		Description: "Auto-generated stub for npm",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "npm"), nil
		},
	})

	// yarn extracted from opencode/packages/opencode/src/installation/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "yarn",
		Description: "Auto-generated stub for yarn",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "yarn"), nil
		},
	})

	// pnpm extracted from opencode/packages/opencode/src/installation/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "pnpm",
		Description: "Auto-generated stub for pnpm",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "pnpm"), nil
		},
	})

	// bun extracted from opencode/packages/opencode/src/installation/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "bun",
		Description: "Auto-generated stub for bun",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "bun"), nil
		},
	})

	// brew extracted from opencode/packages/opencode/src/installation/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "brew",
		Description: "Auto-generated stub for brew",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "brew"), nil
		},
	})

	// scoop extracted from opencode/packages/opencode/src/installation/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "scoop",
		Description: "Auto-generated stub for scoop",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "scoop"), nil
		},
	})

	// choco extracted from opencode/packages/opencode/src/installation/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "choco",
		Description: "Auto-generated stub for choco",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "choco"), nil
		},
	})

	// gofmt extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "gofmt",
		Description: "Auto-generated stub for gofmt",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "gofmt"), nil
		},
	})

	// mix extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mix",
		Description: "Auto-generated stub for mix",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mix"), nil
		},
	})

	// prettier extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "prettier",
		Description: "Auto-generated stub for prettier",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "prettier"), nil
		},
	})

	// oxfmt extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "oxfmt",
		Description: "Auto-generated stub for oxfmt",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "oxfmt"), nil
		},
	})

	// zig extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "zig",
		Description: "Auto-generated stub for zig",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "zig"), nil
		},
	})

	// clang-format extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "clang-format",
		Description: "Auto-generated stub for clang-format",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "clang-format"), nil
		},
	})

	// ktlint extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ktlint",
		Description: "Auto-generated stub for ktlint",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ktlint"), nil
		},
	})

	// ruff extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ruff",
		Description: "Auto-generated stub for ruff",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ruff"), nil
		},
	})

	// air extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "air",
		Description: "Auto-generated stub for air",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "air"), nil
		},
	})

	// uv extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "uv",
		Description: "Auto-generated stub for uv",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "uv"), nil
		},
	})

	// rubocop extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "rubocop",
		Description: "Auto-generated stub for rubocop",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "rubocop"), nil
		},
	})

	// standardrb extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "standardrb",
		Description: "Auto-generated stub for standardrb",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "standardrb"), nil
		},
	})

	// htmlbeautifier extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "htmlbeautifier",
		Description: "Auto-generated stub for htmlbeautifier",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "htmlbeautifier"), nil
		},
	})

	// dart extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "dart",
		Description: "Auto-generated stub for dart",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "dart"), nil
		},
	})

	// ocamlformat extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ocamlformat",
		Description: "Auto-generated stub for ocamlformat",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ocamlformat"), nil
		},
	})

	// terraform extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "terraform",
		Description: "Auto-generated stub for terraform",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "terraform"), nil
		},
	})

	// latexindent extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "latexindent",
		Description: "Auto-generated stub for latexindent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "latexindent"), nil
		},
	})

	// gleam extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "gleam",
		Description: "Auto-generated stub for gleam",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "gleam"), nil
		},
	})

	// shfmt extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "shfmt",
		Description: "Auto-generated stub for shfmt",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "shfmt"), nil
		},
	})

	// nixfmt extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "nixfmt",
		Description: "Auto-generated stub for nixfmt",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "nixfmt"), nil
		},
	})

	// rustfmt extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "rustfmt",
		Description: "Auto-generated stub for rustfmt",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "rustfmt"), nil
		},
	})

	// pint extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "pint",
		Description: "Auto-generated stub for pint",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "pint"), nil
		},
	})

	// ormolu extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ormolu",
		Description: "Auto-generated stub for ormolu",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ormolu"), nil
		},
	})

	// cljfmt extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "cljfmt",
		Description: "Auto-generated stub for cljfmt",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "cljfmt"), nil
		},
	})

	// dfmt extracted from opencode/packages/opencode/src/format/formatter.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "dfmt",
		Description: "Auto-generated stub for dfmt",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "dfmt"), nil
		},
	})

	// main extracted from opencode/packages/opencode/src/git/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "main",
		Description: "Auto-generated stub for main",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "main"), nil
		},
	})

	// master extracted from opencode/packages/opencode/src/git/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "master",
		Description: "Auto-generated stub for master",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "master"), nil
		},
	})

	// OpenCode extracted from opencode/packages/opencode/src/mcp/oauth-provider.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "OpenCode",
		Description: "Auto-generated stub for OpenCode",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "OpenCode"), nil
		},
	})

	// StreamableHTTP extracted from opencode/packages/opencode/src/mcp/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "StreamableHTTP",
		Description: "Auto-generated stub for StreamableHTTP",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "StreamableHTTP"), nil
		},
	})

	// SSE extracted from opencode/packages/opencode/src/mcp/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "SSE",
		Description: "Auto-generated stub for SSE",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "SSE"), nil
		},
	})

	// opencode-debug extracted from opencode/packages/opencode/src/cli/cmd/mcp.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "opencode-debug",
		Description: "Auto-generated stub for opencode-debug",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "opencode-debug"), nil
		},
	})

	// xterm-256color extracted from opencode/packages/opencode/src/pty/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "xterm-256color",
		Description: "Auto-generated stub for xterm-256color",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "xterm-256color"), nil
		},
	})

	// workspace extracted from opencode/packages/opencode/src/lsp/client.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "workspace",
		Description: "Auto-generated stub for workspace",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "workspace"), nil
		},
	})

	// Reject extracted from opencode/packages/opencode/src/acp/agent.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Reject",
		Description: "Auto-generated stub for Reject",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Reject"), nil
		},
	})

	// image extracted from opencode/packages/opencode/src/acp/agent.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "image",
		Description: "Auto-generated stub for image",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "image"), nil
		},
	})

	// file extracted from opencode/packages/opencode/src/acp/agent.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "file",
		Description: "Auto-generated stub for file",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "file"), nil
		},
	})

	// Windsurf extracted from opencode/packages/opencode/src/ide/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Windsurf",
		Description: "Auto-generated stub for Windsurf",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Windsurf"), nil
		},
	})

	// Cursor extracted from opencode/packages/opencode/src/ide/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Cursor",
		Description: "Auto-generated stub for Cursor",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Cursor"), nil
		},
	})

	// VSCodium extracted from opencode/packages/opencode/src/ide/index.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "VSCodium",
		Description: "Auto-generated stub for VSCodium",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "VSCodium"), nil
		},
	})

	// Worktree extracted from opencode/packages/opencode/src/control-plane/adaptors/worktree.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Worktree",
		Description: "Create a git worktree",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Worktree"), nil
		},
	})

	// general extracted from opencode/packages/opencode/src/agent/agent.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "general",
		Description: "Auto-generated stub for general",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "general"), nil
		},
	})

	// explore extracted from opencode/packages/opencode/src/agent/agent.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "explore",
		Description: "Auto-generated stub for explore",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "explore"), nil
		},
	})

	// compaction extracted from opencode/packages/opencode/src/agent/agent.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "compaction",
		Description: "Auto-generated stub for compaction",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "compaction"), nil
		},
	})

	// title extracted from opencode/packages/opencode/src/agent/agent.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "title",
		Description: "Auto-generated stub for title",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "title"), nil
		},
	})

	// summary extracted from opencode/packages/opencode/src/agent/agent.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "summary",
		Description: "Auto-generated stub for summary",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "summary"), nil
		},
	})

	// ProviderAuthError extracted from opencode/packages/sdk/js/src/gen/types.gen.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ProviderAuthError",
		Description: "Auto-generated stub for ProviderAuthError",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ProviderAuthError"), nil
		},
	})

	// MessageOutputLengthError extracted from opencode/packages/sdk/js/src/gen/types.gen.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "MessageOutputLengthError",
		Description: "Auto-generated stub for MessageOutputLengthError",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "MessageOutputLengthError"), nil
		},
	})

	// MessageAbortedError extracted from opencode/packages/sdk/js/src/gen/types.gen.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "MessageAbortedError",
		Description: "Auto-generated stub for MessageAbortedError",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "MessageAbortedError"), nil
		},
	})

	// NotFoundError extracted from opencode/packages/sdk/js/src/gen/types.gen.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "NotFoundError",
		Description: "Auto-generated stub for NotFoundError",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "NotFoundError"), nil
		},
	})

	// StructuredOutputError extracted from opencode/packages/sdk/js/src/v2/gen/types.gen.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "StructuredOutputError",
		Description: "Auto-generated stub for StructuredOutputError",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "StructuredOutputError"), nil
		},
	})

	// chromium extracted from opencode/packages/app/playwright.config.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "chromium",
		Description: "Auto-generated stub for chromium",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "chromium"), nil
		},
	})

	// Provider extracted from opencode/packages/app/src/components/dialog-custom-provider.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Provider",
		Description: "Auto-generated stub for Provider",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Provider"), nil
		},
	})

	// planner extracted from opencode/packages/app/src/components/prompt-input/build-request-parts.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "planner",
		Description: "Auto-generated stub for planner",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "planner"), nil
		},
	})

	// ConfigInvalidError extracted from opencode/packages/app/src/utils/server-errors.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ConfigInvalidError",
		Description: "Auto-generated stub for ConfigInvalidError",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ConfigInvalidError"), nil
		},
	})

	// ServerTimeoutError extracted from opencode/packages/app/src/utils/server-errors.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ServerTimeoutError",
		Description: "Auto-generated stub for ServerTimeoutError",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ServerTimeoutError"), nil
		},
	})

	// ProviderModelNotFoundError extracted from opencode/packages/app/src/utils/server-errors.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ProviderModelNotFoundError",
		Description: "Auto-generated stub for ProviderModelNotFoundError",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ProviderModelNotFoundError"), nil
		},
	})

	// AbortError extracted from opencode/packages/app/src/context/global-sync/utils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "AbortError",
		Description: "Auto-generated stub for AbortError",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "AbortError"), nil
		},
	})

	// src extracted from opencode/packages/app/src/context/file/watcher.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "src",
		Description: "Auto-generated stub for src",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "src"), nil
		},
	})

	// Files extracted from opencode/packages/app/src/constants/file-picker.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Files",
		Description: "Auto-generated stub for Files",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Files"), nil
		},
	})

	// Send extracted from opencode/packages/app/e2e/fixtures.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Send",
		Description: "Auto-generated stub for Send",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Send"), nil
		},
	})

	// Settings extracted from opencode/packages/app/e2e/actions.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Settings",
		Description: "Auto-generated stub for Settings",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Settings"), nil
		},
	})

	// Shortcuts extracted from opencode/packages/app/e2e/settings/settings.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Shortcuts",
		Description: "Auto-generated stub for Shortcuts",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Shortcuts"), nil
		},
	})

	// Providers extracted from opencode/packages/app/e2e/settings/settings-providers.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Providers",
		Description: "Auto-generated stub for Providers",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Providers"), nil
		},
	})

	// Connect extracted from opencode/packages/app/e2e/settings/settings-providers.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Connect",
		Description: "Auto-generated stub for Connect",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Connect"), nil
		},
	})

	// Models extracted from opencode/packages/app/e2e/settings/settings-models.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Models",
		Description: "Auto-generated stub for Models",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Models"), nil
		},
	})

	// Back extracted from opencode/packages/app/e2e/app/titlebar-history.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Back",
		Description: "Auto-generated stub for Back",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Back"), nil
		},
	})

	// Forward extracted from opencode/packages/app/e2e/app/titlebar-history.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Forward",
		Description: "Auto-generated stub for Forward",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Forward"), nil
		},
	})

	// Status extracted from opencode/packages/app/e2e/app/server-default.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Status",
		Description: "Auto-generated stub for Status",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Status"), nil
		},
	})

	// Save extracted from opencode/packages/app/e2e/projects/project-edit.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Save",
		Description: "Auto-generated stub for Save",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Save"), nil
		},
	})

	// Cancel extracted from opencode/packages/app/e2e/projects/project-edit.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Cancel",
		Description: "Auto-generated stub for Cancel",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Cancel"), nil
		},
	})

	// Publish extracted from opencode/packages/app/e2e/session/session.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Publish",
		Description: "Auto-generated stub for Publish",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Publish"), nil
		},
	})

	// Unpublish extracted from opencode/packages/app/e2e/session/session.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Unpublish",
		Description: "Auto-generated stub for Unpublish",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Unpublish"), nil
		},
	})

	// Context extracted from opencode/packages/app/e2e/prompt/context.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Context",
		Description: "Auto-generated stub for Context",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Context"), nil
		},
	})

	// storybook-solidjs-vite extracted from opencode/packages/storybook/.storybook/main.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "storybook-solidjs-vite",
		Description: "Auto-generated stub for storybook-solidjs-vite",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "storybook-solidjs-vite"), nil
		},
	})

	// playground-css extracted from opencode/packages/storybook/.storybook/playground-css-plugin.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "playground-css",
		Description: "Auto-generated stub for playground-css",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "playground-css"), nil
		},
	})

	// fix extracted from opencode/packages/storybook/.storybook/mocks/app/context/sync.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "fix",
		Description: "Run fix command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "fix"), nil
		},
	})

	// auth extracted from opencode/packages/console/app/src/context/auth.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "auth",
		Description: "Auto-generated stub for auth",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "auth"), nil
		},
	})

	// auto extracted from opencode/packages/console/core/src/billing.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "auto",
		Description: "Auto-generated stub for auto",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "auto"), nil
		},
	})

	// Folder extracted from opencode/packages/plugin/src/example-workspace.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Folder",
		Description: "Create a blank folder",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Folder"), nil
		},
	})

	// provider-icons-plugin extracted from opencode/packages/ui/vite.config.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "provider-icons-plugin",
		Description: "Auto-generated stub for provider-icons-plugin",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "provider-icons-plugin"), nil
		},
	})

	// user1 extracted from opencode/packages/enterprise/test/core/storage.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "user1",
		Description: "Auto-generated stub for user1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "user1"), nil
		},
	})

	// user2 extracted from opencode/packages/enterprise/test/core/storage.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "user2",
		Description: "Auto-generated stub for user2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "user2"), nil
		},
	})

	// user3 extracted from opencode/packages/enterprise/test/core/storage.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "user3",
		Description: "Auto-generated stub for user3",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "user3"), nil
		},
	})

	// user4 extracted from opencode/packages/enterprise/test/core/storage.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "user4",
		Description: "Auto-generated stub for user4",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "user4"), nil
		},
	})

	// user5 extracted from opencode/packages/enterprise/test/core/storage.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "user5",
		Description: "Auto-generated stub for user5",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "user5"), nil
		},
	})

	// gemini-3-flash-preview extracted from gemini-cli/integration-tests/test-mcp-support.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "gemini-3-flash-preview",
		Description: "Auto-generated stub for gemini-3-flash-preview",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "gemini-3-flash-preview"), nil
		},
	})

	// addition-server extracted from gemini-cli/integration-tests/simple-mcp-server.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "addition-server",
		Description: "Add two numbers",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "addition-server"), nil
		},
	})

	// add extracted from gemini-cli/integration-tests/simple-mcp-server.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "add",
		Description: "Add two numbers",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "add"), nil
		},
	})

	// write_file extracted from gemini-cli/integration-tests/hooks-system.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "write_file",
		Description: "Auto-generated stub for write_file",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "write_file"), nil
		},
	})

	// hook-a extracted from gemini-cli/integration-tests/hooks-system.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "hook-a",
		Description: "Auto-generated stub for hook-a",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "hook-a"), nil
		},
	})

	// hook-b extracted from gemini-cli/integration-tests/hooks-system.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "hook-b",
		Description: "Auto-generated stub for hook-b",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "hook-b"), nil
		},
	})

	// multi-hook-active extracted from gemini-cli/integration-tests/hooks-system.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "multi-hook-active",
		Description: "Auto-generated stub for multi-hook-active",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "multi-hook-active"), nil
		},
	})

	// multi-hook-disabled extracted from gemini-cli/integration-tests/hooks-system.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "multi-hook-disabled",
		Description: "Auto-generated stub for multi-hook-disabled",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "multi-hook-disabled"), nil
		},
	})

	// test-mcp-server extracted from gemini-cli/integration-tests/test-mcp-server.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-mcp-server",
		Description: "Auto-generated stub for test-mcp-server",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-mcp-server"), nil
		},
	})

	// xterm-color extracted from gemini-cli/integration-tests/symlink-install.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "xterm-color",
		Description: "Auto-generated stub for xterm-color",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "xterm-color"), nil
		},
	})

	// cyclic-schema-server extracted from gemini-cli/integration-tests/mcp_server_cyclic_schema.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "cyclic-schema-server",
		Description: "Auto-generated stub for cyclic-schema-server",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "cyclic-schema-server"), nil
		},
	})

	// tool_with_cyclic_schema extracted from gemini-cli/integration-tests/mcp_server_cyclic_schema.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tool_with_cyclic_schema",
		Description: "Auto-generated stub for tool_with_cyclic_schema",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tool_with_cyclic_schema"), nil
		},
	})

	// test-extension extracted from gemini-cli/integration-tests/extensions-reload.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-extension",
		Description: "Auto-generated stub for test-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-extension"), nil
		},
	})

	// get_context extracted from gemini-cli/packages/sdk/examples/session-context.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "get_context",
		Description: "Get information about the current session context.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "get_context"), nil
		},
	})

	// failVisible extracted from gemini-cli/packages/sdk/src/tool.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "failVisible",
		Description: "Fails with a visible error if input is \"fail\"",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "failVisible"), nil
		},
	})

	// checkSystemStatus extracted from gemini-cli/packages/sdk/src/tool.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "checkSystemStatus",
		Description: "Checks the current system status",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "checkSystemStatus"), nil
		},
	})

	// testTool extracted from gemini-cli/packages/sdk/src/tool.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "testTool",
		Description: "A test tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "testTool"), nil
		},
	})

	// successTool extracted from gemini-cli/packages/sdk/src/tool.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "successTool",
		Description: "Always succeeds",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "successTool"), nil
		},
	})

	// failTool extracted from gemini-cli/packages/sdk/src/tool.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "failTool",
		Description: "Always fails",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "failTool"), nil
		},
	})

	// visibleErrorTool extracted from gemini-cli/packages/sdk/src/tool.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "visibleErrorTool",
		Description: "Fails with visible error",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "visibleErrorTool"), nil
		},
	})

	// catchAllTool extracted from gemini-cli/packages/sdk/src/tool.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "catchAllTool",
		Description: "Catches all errors",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "catchAllTool"), nil
		},
	})

	// ext1 extracted from gemini-cli/packages/core/src/commands/extensions.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ext1",
		Description: "Auto-generated stub for ext1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ext1"), nil
		},
	})

	// ext2 extracted from gemini-cli/packages/core/src/commands/extensions.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ext2",
		Description: "Auto-generated stub for ext2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ext2"), nil
		},
	})

	// save_memory extracted from gemini-cli/packages/core/src/agents/memory-manager-agent.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "save_memory",
		Description: "Auto-generated stub for save_memory",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "save_memory"), nil
		},
	})

	// codebase_investigator extracted from gemini-cli/packages/core/src/agents/codebase-investigator.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "codebase_investigator",
		Description: "Auto-generated stub for codebase_investigator",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "codebase_investigator"), nil
		},
	})

	// write_file_interactive extracted from gemini-cli/packages/core/src/agents/local-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "write_file_interactive",
		Description: "Auto-generated stub for write_file_interactive",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "write_file_interactive"), nil
		},
	})

	// TestAgent extracted from gemini-cli/packages/core/src/agents/local-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "TestAgent",
		Description: "The final result.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "TestAgent"), nil
		},
	})

	// other-tool extracted from gemini-cli/packages/core/src/agents/local-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "other-tool",
		Description: "Auto-generated stub for other-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "other-tool"), nil
		},
	})

	// instantiated_tool extracted from gemini-cli/packages/core/src/agents/local-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "instantiated_tool",
		Description: "Auto-generated stub for instantiated_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "instantiated_tool"), nil
		},
	})

	// BrowserLikeAgent extracted from gemini-cli/packages/core/src/agents/local-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "BrowserLikeAgent",
		Description: "An agent using instance tools.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "BrowserLikeAgent"), nil
		},
	})

	// click extracted from gemini-cli/packages/core/src/agents/local-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "click",
		Description: "Auto-generated stub for click",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "click"), nil
		},
	})

	// fill extracted from gemini-cli/packages/core/src/agents/local-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "fill",
		Description: "Auto-generated stub for fill",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "fill"), nil
		},
	})

	// take_snapshot extracted from gemini-cli/packages/core/src/agents/local-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "take_snapshot",
		Description: "Auto-generated stub for take_snapshot",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "take_snapshot"), nil
		},
	})

	// navigate_page extracted from gemini-cli/packages/core/src/agents/local-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "navigate_page",
		Description: "Auto-generated stub for navigate_page",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "navigate_page"), nil
		},
	})

	// MixedAgent extracted from gemini-cli/packages/core/src/agents/local-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "MixedAgent",
		Description: "Uses both patterns.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "MixedAgent"), nil
		},
	})

	// cli_help extracted from gemini-cli/packages/core/src/agents/cli-help-agent.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "cli_help",
		Description: "Specialized agent for answering questions about the Gemini CLI application. Invoke this agent for questions regarding CLI features, configuration schemas (e.g., policies), or instructions on how to create custom subagents. It queries internal documentation to provide accurate usage guidance.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "cli_help"), nil
		},
	})

	// generalist extracted from gemini-cli/packages/core/src/agents/generalist-agent.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "generalist",
		Description: "A general-purpose AI agent with access to all tools. Highly recommended for tasks that are turn-intensive or involve processing large amounts of data. Use this to keep the main session history lean and efficient. Excellent for: batch refactoring/error fixing across multiple files, running commands with high-volume output, and speculative investigations.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "generalist"), nil
		},
	})

	// MockAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "MockAgent",
		Description: "Mock Description V1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "MockAgent"), nil
		},
	})

	// common-agent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "common-agent",
		Description: "User version",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "common-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "common-agent", result), nil
		},
	})

	// project-only extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "project-only",
		Description: "Project only",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "project-only"), nil
		},
	})

	// extension-agent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "extension-agent",
		Description: "Auto-generated stub for extension-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "extension-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "extension-agent", result), nil
		},
	})

	// RemoteAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "RemoteAgent",
		Description: "A remote agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "RemoteAgent"), nil
		},
	})

	// AutoAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "AutoAgent",
		Description: "Auto-generated stub for AutoAgent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "AutoAgent"), nil
		},
	})

	// RemoteAgentWithAuth extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "RemoteAgentWithAuth",
		Description: "A remote agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "RemoteAgentWithAuth"), nil
		},
	})

	// RemoteAgentBadAuth extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "RemoteAgentBadAuth",
		Description: "A remote agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "RemoteAgentBadAuth"), nil
		},
	})

	// FailAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "FailAgent",
		Description: "An agent that fails to load",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "FailAgent"), nil
		},
	})

	// SecuredAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "SecuredAgent",
		Description: "A secured remote agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "SecuredAgent"), nil
		},
	})

	// x-api-key extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "x-api-key",
		Description: "Auto-generated stub for x-api-key",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "x-api-key"), nil
		},
	})

	// FailingRemoteAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "FailingRemoteAgent",
		Description: "A remote agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "FailingRemoteAgent"), nil
		},
	})

	// RemoteAgentWithDescription extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "RemoteAgentWithDescription",
		Description: "User-provided description",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "RemoteAgentWithDescription"), nil
		},
	})

	// Skill1 extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Skill1",
		Description: "Card-provided description",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Skill1"), nil
		},
	})

	// Skill2 extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Skill2",
		Description: "Desc1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Skill2"), nil
		},
	})

	// RemoteAgentWithSkillsOnly extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "RemoteAgentWithSkillsOnly",
		Description: "User-provided description",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "RemoteAgentWithSkillsOnly"), nil
		},
	})

	// RemoteAgentWithEmptyAgentDescription extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "RemoteAgentWithEmptyAgentDescription",
		Description: "User-provided description",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "RemoteAgentWithEmptyAgentDescription"), nil
		},
	})

	// RemoteAgentAccumulationTest extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "RemoteAgentAccumulationTest",
		Description: "User-provided description",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "RemoteAgentAccumulationTest"), nil
		},
	})

	// EmptyDescAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "EmptyDescAgent",
		Description: "",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "EmptyDescAgent"), nil
		},
	})

	// SkillFallbackAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "SkillFallbackAgent",
		Description: "User description",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "SkillFallbackAgent"), nil
		},
	})

	// SkillNoDesc extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "SkillNoDesc",
		Description: "Card description",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "SkillNoDesc"), nil
		},
	})

	// RemotePolicyAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "RemotePolicyAgent",
		Description: "A remote agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "RemotePolicyAgent"), nil
		},
	})

	// invoke_agent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "invoke_agent",
		Description: "Auto-generated stub for invoke_agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "invoke_agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "invoke_agent", result), nil
		},
	})

	// LocalPolicyAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "LocalPolicyAgent",
		Description: "Auto-generated stub for LocalPolicyAgent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "LocalPolicyAgent"), nil
		},
	})

	// OverwrittenAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "OverwrittenAgent",
		Description: "A remote agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "OverwrittenAgent"), nil
		},
	})

	// InitialAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "InitialAgent",
		Description: "Auto-generated stub for InitialAgent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "InitialAgent"), nil
		},
	})

	// NewAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "NewAgent",
		Description: "Auto-generated stub for NewAgent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "NewAgent"), nil
		},
	})

	// InheritingAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "InheritingAgent",
		Description: "Auto-generated stub for InheritingAgent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "InheritingAgent"), nil
		},
	})

	// AnotherAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "AnotherAgent",
		Description: "Auto-generated stub for AnotherAgent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "AnotherAgent"), nil
		},
	})

	// EnabledAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "EnabledAgent",
		Description: "Auto-generated stub for EnabledAgent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "EnabledAgent"), nil
		},
	})

	// DisabledAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "DisabledAgent",
		Description: "Auto-generated stub for DisabledAgent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "DisabledAgent"), nil
		},
	})

	// GetterAgent extracted from gemini-cli/packages/core/src/agents/registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "GetterAgent",
		Description: "Auto-generated stub for GetterAgent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "GetterAgent"), nil
		},
	})

	// VeryLongAgentNameThatTakesUpSpace extracted from gemini-cli/packages/core/src/agents/local-invocation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "VeryLongAgentNameThatTakesUpSpace",
		Description: "Auto-generated stub for VeryLongAgentNameThatTakesUpSpace",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "VeryLongAgentNameThatTakesUpSpace"), nil
		},
	})

	// ls extracted from gemini-cli/packages/core/src/agents/local-invocation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ls",
		Description: "Auto-generated stub for ls",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ls"), nil
		},
	})

	// ProjectAgent extracted from gemini-cli/packages/core/src/agents/registry_acknowledgement.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ProjectAgent",
		Description: "Project Agent Desc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ProjectAgent"), nil
		},
	})

	// World extracted from gemini-cli/packages/core/src/agents/utils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "World",
		Description: "Auto-generated stub for World",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "World"), nil
		},
	})

	// Jo extracted from gemini-cli/packages/core/src/agents/utils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Jo",
		Description: "Auto-generated stub for Jo",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Jo"), nil
		},
	})

	// Alice extracted from gemini-cli/packages/core/src/agents/utils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Alice",
		Description: "Auto-generated stub for Alice",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Alice"), nil
		},
	})

	// test-agent extracted from gemini-cli/packages/core/src/agents/a2a-client-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-agent",
		Description: "A test agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "test-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "test-agent", result), nil
		},
	})

	// test-agent-md extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-agent-md",
		Description: "A markdown agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "test-agent-md", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "test-agent-md", result), nil
		},
	})

	// complex-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "complex-agent",
		Description: "A complex markdown agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "complex-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "complex-agent", result), nil
		},
	})

	// mcp-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp-agent",
		Description: "An agent with MCP servers",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "mcp-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "mcp-agent", result), nil
		},
	})

	// remote-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "remote-agent",
		Description: "A remote agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "remote-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "remote-agent", result), nil
		},
	})

	// inferred-remote extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "inferred-remote",
		Description: "Inferred",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "inferred-remote"), nil
		},
	})

	// no-body-remote extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "no-body-remote",
		Description: "Auto-generated stub for no-body-remote",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "no-body-remote"), nil
		},
	})

	// remote-1 extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "remote-1",
		Description: "Auto-generated stub for remote-1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "remote-1"), nil
		},
	})

	// remote-2 extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "remote-2",
		Description: "Auto-generated stub for remote-2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "remote-2"), nil
		},
	})

	// no-trailing-newline extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "no-trailing-newline",
		Description: "Auto-generated stub for no-trailing-newline",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "no-trailing-newline"), nil
		},
	})

	// json-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "json-agent",
		Description: "A JSON-based remote agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "json-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "json-agent", result), nil
		},
	})

	// json-remote extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "json-remote",
		Description: "A JSON-based remote agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "json-remote"), nil
		},
	})

	// inferred-json-remote extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "inferred-json-remote",
		Description: "Auto-generated stub for inferred-json-remote",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "inferred-json-remote"), nil
		},
	})

	// spanner-test-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "spanner-test-agent",
		Description: "An agent to test Spanner MCP with auth",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "spanner-test-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "spanner-test-agent", result), nil
		},
	})

	// no-card-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "no-card-agent",
		Description: "Missing card info",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "no-card-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "no-card-agent", result), nil
		},
	})

	// api-key-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "api-key-agent",
		Description: "Auto-generated stub for api-key-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "api-key-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "api-key-agent", result), nil
		},
	})

	// X-Custom-Key extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "X-Custom-Key",
		Description: "Auto-generated stub for X-Custom-Key",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "X-Custom-Key"), nil
		},
	})

	// bearer-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "bearer-agent",
		Description: "Auto-generated stub for bearer-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "bearer-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "bearer-agent", result), nil
		},
	})

	// basic-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "basic-agent",
		Description: "Auto-generated stub for basic-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "basic-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "basic-agent", result), nil
		},
	})

	// digest-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "digest-agent",
		Description: "Auto-generated stub for digest-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "digest-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "digest-agent", result), nil
		},
	})

	// raw-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "raw-agent",
		Description: "Auto-generated stub for raw-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "raw-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "raw-agent", result), nil
		},
	})

	// auth-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "auth-agent",
		Description: "Auto-generated stub for auth-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "auth-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "auth-agent", result), nil
		},
	})

	// oauth2-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "oauth2-agent",
		Description: "Auto-generated stub for oauth2-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "oauth2-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "oauth2-agent", result), nil
		},
	})

	// oauth2-full-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "oauth2-full-agent",
		Description: "Auto-generated stub for oauth2-full-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "oauth2-full-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "oauth2-full-agent", result), nil
		},
	})

	// oauth2-minimal-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "oauth2-minimal-agent",
		Description: "Auto-generated stub for oauth2-minimal-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "oauth2-minimal-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "oauth2-minimal-agent", result), nil
		},
	})

	// oauth2-convert-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "oauth2-convert-agent",
		Description: "Auto-generated stub for oauth2-convert-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "oauth2-convert-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "oauth2-convert-agent", result), nil
		},
	})

	// unknown-auth-agent extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "unknown-auth-agent",
		Description: "Auto-generated stub for unknown-auth-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "unknown-auth-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "unknown-auth-agent", result), nil
		},
	})

	// orphan extracted from gemini-cli/packages/core/src/agents/agentLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "orphan",
		Description: "Auto-generated stub for orphan",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "orphan"), nil
		},
	})

	// my-agent extracted from gemini-cli/packages/core/src/agents/a2aUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-agent",
		Description: "Auto-generated stub for my-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "my-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "my-agent", result), nil
		},
	})

	// priority-test extracted from gemini-cli/packages/core/src/agents/a2aUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "priority-test",
		Description: "Auto-generated stub for priority-test",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "priority-test"), nil
		},
	})

	// Code extracted from gemini-cli/packages/core/src/agents/a2aUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Code",
		Description: "Auto-generated stub for Code",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Code"), nil
		},
	})

	// Data extracted from gemini-cli/packages/core/src/agents/a2aUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Data",
		Description: "Auto-generated stub for Data",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Data"), nil
		},
	})

	// confucius extracted from gemini-cli/packages/core/src/agents/skill-extraction-agent.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "confucius",
		Description: "Extracts reusable skills from past conversation sessions and writes them as SKILL.md files.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "confucius"), nil
		},
	})

	// test-tool extracted from gemini-cli/packages/core/src/agents/agent-scheduler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-tool",
		Description: "Auto-generated stub for test-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-tool"), nil
		},
	})

	// new_page extracted from gemini-cli/packages/core/src/agents/agent-scheduler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "new_page",
		Description: "Auto-generated stub for new_page",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "new_page"), nil
		},
	})

	// TestLocalAgent extracted from gemini-cli/packages/core/src/agents/agent-tool.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "TestLocalAgent",
		Description: "A local test agent.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "TestLocalAgent"), nil
		},
	})

	// TestRemoteAgent extracted from gemini-cli/packages/core/src/agents/agent-tool.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "TestRemoteAgent",
		Description: "A remote test agent.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "TestRemoteAgent"), nil
		},
	})

	// UnknownAgent extracted from gemini-cli/packages/core/src/agents/agent-tool.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "UnknownAgent",
		Description: "Auto-generated stub for UnknownAgent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "UnknownAgent"), nil
		},
	})

	// admin extracted from gemini-cli/packages/core/src/agents/remote-invocation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "admin",
		Description: "Auto-generated stub for admin",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "admin"), nil
		},
	})

	// Result extracted from gemini-cli/packages/core/src/agents/remote-invocation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Result",
		Description: "Auto-generated stub for Result",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Result"), nil
		},
	})

	// browser_agent extracted from gemini-cli/packages/core/src/agents/browser/browserAgentInvocation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "browser_agent",
		Description: "mock definition",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "browser_agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "browser_agent", result), nil
		},
	})

	// navigate_browser extracted from gemini-cli/packages/core/src/agents/browser/browserAgentInvocation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "navigate_browser",
		Description: "Auto-generated stub for navigate_browser",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "navigate_browser"), nil
		},
	})

	// fill_form extracted from gemini-cli/packages/core/src/agents/browser/browserAgentInvocation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "fill_form",
		Description: "Auto-generated stub for fill_form",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "fill_form"), nil
		},
	})

	// click_element extracted from gemini-cli/packages/core/src/agents/browser/browserAgentInvocation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "click_element",
		Description: "Auto-generated stub for click_element",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "click_element"), nil
		},
	})

	// tool_a extracted from gemini-cli/packages/core/src/agents/browser/browserAgentInvocation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tool_a",
		Description: "Auto-generated stub for tool_a",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tool_a"), nil
		},
	})

	// tool_b extracted from gemini-cli/packages/core/src/agents/browser/browserAgentInvocation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tool_b",
		Description: "Auto-generated stub for tool_b",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tool_b"), nil
		},
	})

	// upload_file extracted from gemini-cli/packages/core/src/agents/browser/mcpToolWrapper.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "upload_file",
		Description: "Upload a file",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "upload_file"), nil
		},
	})

	// type_text extracted from gemini-cli/packages/core/src/agents/browser/browserAgentFactory.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "type_text",
		Description: "Navigate to URL",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "type_text"), nil
		},
	})

	// click_at extracted from gemini-cli/packages/core/src/agents/browser/browserAgentFactory.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "click_at",
		Description: "Type text into an element",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "click_at"), nil
		},
	})

	// close_page extracted from gemini-cli/packages/core/src/agents/browser/browserAgentFactory.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "close_page",
		Description: "Open a new page/tab",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "close_page"), nil
		},
	})

	// select_page extracted from gemini-cli/packages/core/src/agents/browser/browserAgentFactory.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "select_page",
		Description: "Close page",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "select_page"), nil
		},
	})

	// press_key extracted from gemini-cli/packages/core/src/agents/browser/browserAgentFactory.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "press_key",
		Description: "Close page",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "press_key"), nil
		},
	})

	// hover extracted from gemini-cli/packages/core/src/agents/browser/browserAgentFactory.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "hover",
		Description: "Type text into an element",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "hover"), nil
		},
	})

	// take_screenshot extracted from gemini-cli/packages/core/src/agents/browser/browserAgentFactory.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "take_screenshot",
		Description: "Take snapshot",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "take_screenshot"), nil
		},
	})

	// list_pages extracted from gemini-cli/packages/core/src/agents/browser/browserAgentFactory.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "list_pages",
		Description: "Take screenshot",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "list_pages"), nil
		},
	})

	// gemini-cli-browser-agent extracted from gemini-cli/packages/core/src/agents/browser/browserManager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "gemini-cli-browser-agent",
		Description: "Auto-generated stub for gemini-cli-browser-agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "gemini-cli-browser-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "gemini-cli-browser-agent", result), nil
		},
	})

	// X-API-Key extracted from gemini-cli/packages/core/src/agents/auth-provider/factory.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "X-API-Key",
		Description: "Auto-generated stub for X-API-Key",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "X-API-Key"), nil
		},
	})

	// X-Custom-Auth extracted from gemini-cli/packages/core/src/agents/auth-provider/api-key-provider.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "X-Custom-Auth",
		Description: "Auto-generated stub for X-Custom-Auth",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "X-Custom-Auth"), nil
		},
	})

	// someTool extracted from gemini-cli/packages/core/src/services/loopDetectionService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "someTool",
		Description: "Auto-generated stub for someTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "someTool"), nil
		},
	})

	// packages extracted from gemini-cli/packages/core/src/services/sandboxManager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "packages",
		Description: "Auto-generated stub for packages",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "packages"), nil
		},
	})

	// mock-pty extracted from gemini-cli/packages/core/src/services/shellExecutionService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mock-pty",
		Description: "Auto-generated stub for mock-pty",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mock-pty"), nil
		},
	})

	// skill-extraction extracted from gemini-cli/packages/core/src/services/memoryService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "skill-extraction",
		Description: "Auto-generated stub for skill-extraction",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "skill-extraction"), nil
		},
	})

	// list_files extracted from gemini-cli/packages/core/src/services/chatRecordingService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "list_files",
		Description: "Auto-generated stub for list_files",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "list_files"), nil
		},
	})

	// read_file extracted from gemini-cli/packages/core/src/services/chatRecordingService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "read_file",
		Description: "Auto-generated stub for read_file",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "read_file"), nil
		},
	})

	// feature-x extracted from gemini-cli/packages/core/src/services/worktreeService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "feature-x",
		Description: "Auto-generated stub for feature-x",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "feature-x"), nil
		},
	})

	// plainTool extracted from gemini-cli/packages/core/src/scheduler/tool-modifier.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "plainTool",
		Description: "Auto-generated stub for plainTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "plainTool"), nil
		},
	})

	// mockModifiableTool extracted from gemini-cli/packages/core/src/scheduler/tool-modifier.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mockModifiableTool",
		Description: "Auto-generated stub for mockModifiableTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mockModifiableTool"), nil
		},
	})

	// webSearchTool extracted from gemini-cli/packages/core/src/scheduler/tool-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "webSearchTool",
		Description: "Mock web search",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "webSearchTool"), nil
		},
	})

	// slowTool extracted from gemini-cli/packages/core/src/scheduler/tool-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "slowTool",
		Description: "Auto-generated stub for slowTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "slowTool"), nil
		},
	})

	// actualToolName extracted from gemini-cli/packages/core/src/scheduler/tool-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "actualToolName",
		Description: "Auto-generated stub for actualToolName",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "actualToolName"), nil
		},
	})

	// get_big_text extracted from gemini-cli/packages/core/src/scheduler/tool-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "get_big_text",
		Description: "Auto-generated stub for get_big_text",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "get_big_text"), nil
		},
	})

	// remote_agent_call extracted from gemini-cli/packages/core/src/scheduler/tool-executor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "remote_agent_call",
		Description: "Remote agent call",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "remote_agent_call", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "remote_agent_call", result), nil
		},
	})

	// mockTool extracted from gemini-cli/packages/core/src/scheduler/scheduler_hooks.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mockTool",
		Description: "Auto-generated stub for mockTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mockTool"), nil
		},
	})

	// test-stop-hook extracted from gemini-cli/packages/core/src/scheduler/scheduler_hooks.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-stop-hook",
		Description: "Auto-generated stub for test-stop-hook",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-stop-hook"), nil
		},
	})

	// test-block-hook extracted from gemini-cli/packages/core/src/scheduler/scheduler_hooks.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-block-hook",
		Description: "Auto-generated stub for test-block-hook",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-block-hook"), nil
		},
	})

	// test-modify-input-hook extracted from gemini-cli/packages/core/src/scheduler/scheduler_hooks.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-modify-input-hook",
		Description: "Auto-generated stub for test-modify-input-hook",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-modify-input-hook"), nil
		},
	})

	// read-tool-1 extracted from gemini-cli/packages/core/src/scheduler/scheduler_parallel.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "read-tool-1",
		Description: "Auto-generated stub for read-tool-1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "read-tool-1"), nil
		},
	})

	// read-tool-2 extracted from gemini-cli/packages/core/src/scheduler/scheduler_parallel.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "read-tool-2",
		Description: "Auto-generated stub for read-tool-2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "read-tool-2"), nil
		},
	})

	// write-tool extracted from gemini-cli/packages/core/src/scheduler/scheduler_parallel.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "write-tool",
		Description: "Auto-generated stub for write-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "write-tool"), nil
		},
	})

	// agent-tool-1 extracted from gemini-cli/packages/core/src/scheduler/scheduler_parallel.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "agent-tool-1",
		Description: "Auto-generated stub for agent-tool-1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "agent-tool-1", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "agent-tool-1", result), nil
		},
	})

	// agent-tool-2 extracted from gemini-cli/packages/core/src/scheduler/scheduler_parallel.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "agent-tool-2",
		Description: "Auto-generated stub for agent-tool-2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "agent-tool-2", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "agent-tool-2", result), nil
		},
	})

	// tool extracted from gemini-cli/packages/core/src/scheduler/confirmation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tool",
		Description: "Auto-generated stub for tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tool"), nil
		},
	})

	// original-tool-name extracted from gemini-cli/packages/core/src/scheduler/scheduler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "original-tool-name",
		Description: "Auto-generated stub for original-tool-name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "original-tool-name"), nil
		},
	})

	// tool-b extracted from gemini-cli/packages/core/src/scheduler/scheduler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tool-b",
		Description: "Auto-generated stub for tool-b",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tool-b"), nil
		},
	})

	// missing-tool extracted from gemini-cli/packages/core/src/scheduler/scheduler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "missing-tool",
		Description: "Auto-generated stub for missing-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "missing-tool"), nil
		},
	})

	// mcp-tool extracted from gemini-cli/packages/core/src/scheduler/scheduler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp-tool",
		Description: "Auto-generated stub for mcp-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp-tool"), nil
		},
	})

	// replace extracted from gemini-cli/packages/core/src/scheduler/policy.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "replace",
		Description: "Auto-generated stub for replace",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "replace"), nil
		},
	})

	// run_shell_command extracted from gemini-cli/packages/core/src/scheduler/policy.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "run_shell_command",
		Description: "Auto-generated stub for run_shell_command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "run_shell_command"), nil
		},
	})

	// skill1 extracted from gemini-cli/packages/core/src/config/config.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "skill1",
		Description: "Auto-generated stub for skill1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "skill1"), nil
		},
	})

	// Global extracted from gemini-cli/packages/core/src/config/memory.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Global",
		Description: "Auto-generated stub for Global",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Global"), nil
		},
	})

	// Extension extracted from gemini-cli/packages/core/src/config/memory.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Extension",
		Description: "Auto-generated stub for Extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Extension"), nil
		},
	})

	// Project extracted from gemini-cli/packages/core/src/config/memory.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Project",
		Description: "Auto-generated stub for Project",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Project"), nil
		},
	})

	// hook1 extracted from gemini-cli/packages/core/src/hooks/trustedHooks.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "hook1",
		Description: "Auto-generated stub for hook1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "hook1"), nil
		},
	})

	// trusted-hook extracted from gemini-cli/packages/core/src/hooks/trustedHooks.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "trusted-hook",
		Description: "Auto-generated stub for trusted-hook",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "trusted-hook"), nil
		},
	})

	// new-hook extracted from gemini-cli/packages/core/src/hooks/trustedHooks.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "new-hook",
		Description: "Auto-generated stub for new-hook",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "new-hook"), nil
		},
	})

	// my-hook extracted from gemini-cli/packages/core/src/hooks/trustedHooks.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-hook",
		Description: "Auto-generated stub for my-hook",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-hook"), nil
		},
	})

	// hook2 extracted from gemini-cli/packages/core/src/hooks/hookPlanner.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "hook2",
		Description: "Auto-generated stub for hook2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "hook2"), nil
		},
	})

	// test-hook extracted from gemini-cli/packages/core/src/hooks/runtimeHooks.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-hook",
		Description: "Auto-generated stub for test-hook",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-hook"), nil
		},
	})

	// TestTool extracted from gemini-cli/packages/core/src/hooks/runtimeHooks.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "TestTool",
		Description: "Auto-generated stub for TestTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "TestTool"), nil
		},
	})

	// BeforeTool extracted from gemini-cli/packages/core/src/hooks/runtimeHooks.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "BeforeTool",
		Description: "Auto-generated stub for BeforeTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "BeforeTool"), nil
		},
	})

	// fail-hook extracted from gemini-cli/packages/core/src/hooks/runtimeHooks.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "fail-hook",
		Description: "Auto-generated stub for fail-hook",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "fail-hook"), nil
		},
	})

	// persist-hook extracted from gemini-cli/packages/core/src/hooks/runtimeHooks.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "persist-hook",
		Description: "Auto-generated stub for persist-hook",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "persist-hook"), nil
		},
	})

	// disabled-hook extracted from gemini-cli/packages/core/src/hooks/hookRegistry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "disabled-hook",
		Description: "Auto-generated stub for disabled-hook",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "disabled-hook"), nil
		},
	})

	// friendly-name extracted from gemini-cli/packages/core/src/hooks/hookRegistry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "friendly-name",
		Description: "Auto-generated stub for friendly-name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "friendly-name"), nil
		},
	})

	// InvalidEvent extracted from gemini-cli/packages/core/src/hooks/hookRegistry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "InvalidEvent",
		Description: "Auto-generated stub for InvalidEvent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "InvalidEvent"), nil
		},
	})

	// my-friendly-hook extracted from gemini-cli/packages/core/src/hooks/hookRunner.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-friendly-hook",
		Description: "Auto-generated stub for my-friendly-hook",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-friendly-hook"), nil
		},
	})

	// EditTool extracted from gemini-cli/packages/core/src/hooks/hookEventHandler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "EditTool",
		Description: "Auto-generated stub for EditTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "EditTool"), nil
		},
	})

	// my-mcp-server extracted from gemini-cli/packages/core/src/hooks/hookEventHandler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-mcp-server",
		Description: "Auto-generated stub for my-mcp-server",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-mcp-server"), nil
		},
	})

	// my-mcp-server__read_file extracted from gemini-cli/packages/core/src/hooks/hookEventHandler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-mcp-server__read_file",
		Description: "Auto-generated stub for my-mcp-server__read_file",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-mcp-server__read_file"), nil
		},
	})

	// failing-hook extracted from gemini-cli/packages/core/src/hooks/hookEventHandler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "failing-hook",
		Description: "Auto-generated stub for failing-hook",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "failing-hook"), nil
		},
	})

	// new-extension extracted from gemini-cli/packages/core/src/utils/memoryDiscovery.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "new-extension",
		Description: "Auto-generated stub for new-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "new-extension"), nil
		},
	})

	// test-checker extracted from gemini-cli/packages/core/src/utils/extensionLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-checker",
		Description: "Auto-generated stub for test-checker",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-checker"), nil
		},
	})

	// another-extension extracted from gemini-cli/packages/core/src/utils/extensionLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "another-extension",
		Description: "Auto-generated stub for another-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "another-extension"), nil
		},
	})

	// myFunction extracted from gemini-cli/packages/core/src/utils/partUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "myFunction",
		Description: "Auto-generated stub for myFunction",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "myFunction"), nil
		},
	})

	// func1 extracted from gemini-cli/packages/core/src/utils/partUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "func1",
		Description: "Auto-generated stub for func1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "func1"), nil
		},
	})

	// do_stuff extracted from gemini-cli/packages/core/src/utils/partUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "do_stuff",
		Description: "Auto-generated stub for do_stuff",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "do_stuff"), nil
		},
	})

	// keep extracted from gemini-cli/packages/core/src/utils/partUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "keep",
		Description: "Auto-generated stub for keep",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "keep"), nil
		},
	})

	// lydell-node-pty extracted from gemini-cli/packages/core/src/utils/getPty.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "lydell-node-pty",
		Description: "Auto-generated stub for lydell-node-pty",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "lydell-node-pty"), nil
		},
	})

	// node-pty extracted from gemini-cli/packages/core/src/utils/getPty.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "node-pty",
		Description: "Auto-generated stub for node-pty",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "node-pty"), nil
		},
	})

	// Get-ChildItem extracted from gemini-cli/packages/core/src/utils/shell-utils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Get-ChildItem",
		Description: "Auto-generated stub for Get-ChildItem",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Get-ChildItem"), nil
		},
	})

	// Select-Object extracted from gemini-cli/packages/core/src/utils/shell-utils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Select-Object",
		Description: "Auto-generated stub for Select-Object",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Select-Object"), nil
		},
	})

	// echo extracted from gemini-cli/packages/core/src/utils/shell-utils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "echo",
		Description: "Auto-generated stub for echo",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "echo"), nil
		},
	})

	// Bob extracted from gemini-cli/packages/core/src/utils/textUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Bob",
		Description: "Auto-generated stub for Bob",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Bob"), nil
		},
	})

	// Charlie extracted from gemini-cli/packages/core/src/utils/textUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Charlie",
		Description: "Auto-generated stub for Charlie",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Charlie"), nil
		},
	})

	// myTool extracted from gemini-cli/packages/core/src/utils/apiConversionUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "myTool",
		Description: "Auto-generated stub for myTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "myTool"), nil
		},
	})

	// foo extracted from gemini-cli/packages/core/src/utils/tokenCalculation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "foo",
		Description: "Auto-generated stub for foo",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "foo"), nil
		},
	})

	// multimodal_tool extracted from gemini-cli/packages/core/src/utils/tokenCalculation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "multimodal_tool",
		Description: "Auto-generated stub for multimodal_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "multimodal_tool"), nil
		},
	})

	// d0 extracted from gemini-cli/packages/core/src/utils/tokenCalculation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "d0",
		Description: "Auto-generated stub for d0",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "d0"), nil
		},
	})

	// d1 extracted from gemini-cli/packages/core/src/utils/tokenCalculation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "d1",
		Description: "Auto-generated stub for d1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "d1"), nil
		},
	})

	// d2 extracted from gemini-cli/packages/core/src/utils/tokenCalculation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "d2",
		Description: "Auto-generated stub for d2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "d2"), nil
		},
	})

	// d3 extracted from gemini-cli/packages/core/src/utils/tokenCalculation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "d3",
		Description: "Auto-generated stub for d3",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "d3"), nil
		},
	})

	// d4 extracted from gemini-cli/packages/core/src/utils/tokenCalculation.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "d4",
		Description: "Auto-generated stub for d4",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "d4"), nil
		},
	})

	// HeadersTimeoutError extracted from gemini-cli/packages/core/src/utils/errors_timeout.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "HeadersTimeoutError",
		Description: "Auto-generated stub for HeadersTimeoutError",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "HeadersTimeoutError"), nil
		},
	})

	// testFunc extracted from gemini-cli/packages/core/src/utils/generateContentResponseUtilities.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "testFunc",
		Description: "Auto-generated stub for testFunc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "testFunc"), nil
		},
	})

	// testFunc1 extracted from gemini-cli/packages/core/src/utils/generateContentResponseUtilities.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "testFunc1",
		Description: "Auto-generated stub for testFunc1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "testFunc1"), nil
		},
	})

	// testFunc2 extracted from gemini-cli/packages/core/src/utils/generateContentResponseUtilities.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "testFunc2",
		Description: "Auto-generated stub for testFunc2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "testFunc2"), nil
		},
	})

	// processData extracted from gemini-cli/packages/core/src/utils/generateContentResponseUtilities.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "processData",
		Description: "Auto-generated stub for processData",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "processData"), nil
		},
	})

	// test-pkg extracted from gemini-cli/packages/core/src/utils/package.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-pkg",
		Description: "no package.json is found",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-pkg"), nil
		},
	})

	// A extracted from gemini-cli/packages/core/src/utils/markdownUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "A",
		Description: "Auto-generated stub for A",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "A"), nil
		},
	})

	// B extracted from gemini-cli/packages/core/src/utils/markdownUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "B",
		Description: "Line 1\nLine 2\nLine 3",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "B"), nil
		},
	})

	// some_tool extracted from gemini-cli/packages/core/src/utils/tool-visibility.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "some_tool",
		Description: "Auto-generated stub for some_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "some_tool"), nil
		},
	})

	// tool1 extracted from gemini-cli/packages/core/src/core/turn.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tool1",
		Description: "Auto-generated stub for tool1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tool1"), nil
		},
	})

	// tool2 extracted from gemini-cli/packages/core/src/core/turn.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tool2",
		Description: "Auto-generated stub for tool2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tool2"), nil
		},
	})

	// undefined_tool_name extracted from gemini-cli/packages/core/src/core/turn.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "undefined_tool_name",
		Description: "Auto-generated stub for undefined_tool_name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "undefined_tool_name"), nil
		},
	})

	// ReadFile extracted from gemini-cli/packages/core/src/core/turn.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ReadFile",
		Description: "Auto-generated stub for ReadFile",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ReadFile"), nil
		},
	})

	// debugTool extracted from gemini-cli/packages/core/src/core/turn.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "debugTool",
		Description: "Auto-generated stub for debugTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "debugTool"), nil
		},
	})

	// find_restaurant extracted from gemini-cli/packages/core/src/core/geminiChat.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "find_restaurant",
		Description: "Auto-generated stub for find_restaurant",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "find_restaurant"), nil
		},
	})

	// Vesuvio extracted from gemini-cli/packages/core/src/core/geminiChat.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Vesuvio",
		Description: "Auto-generated stub for Vesuvio",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Vesuvio"), nil
		},
	})

	// test_function extracted from gemini-cli/packages/core/src/core/geminiChat.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test_function",
		Description: "Auto-generated stub for test_function",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test_function"), nil
		},
	})

	// old_tool extracted from gemini-cli/packages/core/src/core/geminiChat.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "old_tool",
		Description: "Auto-generated stub for old_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "old_tool"), nil
		},
	})

	// find_restaurant_2 extracted from gemini-cli/packages/core/src/core/geminiChat.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "find_restaurant_2",
		Description: "Auto-generated stub for find_restaurant_2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "find_restaurant_2"), nil
		},
	})

	// tool_with_sig extracted from gemini-cli/packages/core/src/core/geminiChat.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tool_with_sig",
		Description: "Auto-generated stub for tool_with_sig",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tool_with_sig"), nil
		},
	})

	// another_tool extracted from gemini-cli/packages/core/src/core/geminiChat.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "another_tool",
		Description: "Auto-generated stub for another_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "another_tool"), nil
		},
	})

	// tail-tool extracted from gemini-cli/packages/core/src/core/geminiChat.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tail-tool",
		Description: "Auto-generated stub for tail-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tail-tool"), nil
		},
	})

	// mock-tool extracted from gemini-cli/packages/core/src/core/geminiChat.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mock-tool",
		Description: "Auto-generated stub for mock-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mock-tool"), nil
		},
	})

	// original-tool extracted from gemini-cli/packages/core/src/core/geminiChat.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "original-tool",
		Description: "Auto-generated stub for original-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "original-tool"), nil
		},
	})

	// tool-name extracted from gemini-cli/packages/core/src/core/geminiChat.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tool-name",
		Description: "Auto-generated stub for tool-name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tool-name"), nil
		},
	})

	// streamTool extracted from gemini-cli/packages/core/src/core/loggingContentGenerator.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "streamTool",
		Description: "Auto-generated stub for streamTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "streamTool"), nil
		},
	})

	// mcp_myserver_search extracted from gemini-cli/packages/core/src/core/loggingContentGenerator.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_myserver_search",
		Description: "Search via MCP",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_myserver_search"), nil
		},
	})

	// __leading extracted from gemini-cli/packages/core/src/core/loggingContentGenerator.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "__leading",
		Description: "test",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "__leading"), nil
		},
	})

	// trailing__ extracted from gemini-cli/packages/core/src/core/loggingContentGenerator.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "trailing__",
		Description: "test",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "trailing__"), nil
		},
	})

	// a__b__c extracted from gemini-cli/packages/core/src/core/loggingContentGenerator.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "a__b__c",
		Description: "test",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "a__b__c"), nil
		},
	})

	// mock-agent extracted from gemini-cli/packages/core/src/core/prompts.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mock-agent",
		Description: "Mock Agent Description",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "mock-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "mock-agent", result), nil
		},
	})

	// test-skill extracted from gemini-cli/packages/core/src/core/prompts.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-skill",
		Description: "A test skill description",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-skill"), nil
		},
	})

	// glob extracted from gemini-cli/packages/core/src/core/prompts.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "glob",
		Description: "Auto-generated stub for glob",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "glob"), nil
		},
	})

	// grep_search extracted from gemini-cli/packages/core/src/core/prompts.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "grep_search",
		Description: "Auto-generated stub for grep_search",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "grep_search"), nil
		},
	})

	// ask_user extracted from gemini-cli/packages/core/src/core/prompts.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ask_user",
		Description: "Auto-generated stub for ask_user",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ask_user"), nil
		},
	})

	// exit_plan_mode extracted from gemini-cli/packages/core/src/core/prompts.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "exit_plan_mode",
		Description: "Auto-generated stub for exit_plan_mode",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "exit_plan_mode"), nil
		},
	})

	// huge_tool extracted from gemini-cli/packages/core/src/core/client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "huge_tool",
		Description: "Auto-generated stub for huge_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "huge_tool"), nil
		},
	})

	// Read extracted from gemini-cli/packages/core/src/output/stream-json-formatter.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Read",
		Description: "Auto-generated stub for Read",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Read"), nil
		},
	})

	// noise extracted from gemini-cli/packages/core/src/routing/strategies/classifierStrategy.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "noise",
		Description: "Auto-generated stub for noise",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "noise"), nil
		},
	})

	// strategy1 extracted from gemini-cli/packages/core/src/routing/strategies/compositeStrategy.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "strategy1",
		Description: "Auto-generated stub for strategy1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "strategy1"), nil
		},
	})

	// strategy2 extracted from gemini-cli/packages/core/src/routing/strategies/compositeStrategy.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "strategy2",
		Description: "Auto-generated stub for strategy2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "strategy2"), nil
		},
	})

	// terminal extracted from gemini-cli/packages/core/src/routing/strategies/compositeStrategy.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "terminal",
		Description: "Auto-generated stub for terminal",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "terminal"), nil
		},
	})

	// allow_git_status_rule extracted from gemini-cli/packages/core/src/policy/shell-safety.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "allow_git_status_rule",
		Description: "Auto-generated stub for allow_git_status_rule",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "allow_git_status_rule"), nil
		},
	})

	// ask_another_unknown_command_rule extracted from gemini-cli/packages/core/src/policy/shell-safety.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ask_another_unknown_command_rule",
		Description: "Auto-generated stub for ask_another_unknown_command_rule",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ask_another_unknown_command_rule"), nil
		},
	})

	// ask_rule_1 extracted from gemini-cli/packages/core/src/policy/shell-safety.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ask_rule_1",
		Description: "Auto-generated stub for ask_rule_1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ask_rule_1"), nil
		},
	})

	// ask_rule_2 extracted from gemini-cli/packages/core/src/policy/shell-safety.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ask_rule_2",
		Description: "Auto-generated stub for ask_rule_2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ask_rule_2"), nil
		},
	})

	// shell extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "shell",
		Description: "Auto-generated stub for shell",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "shell"), nil
		},
	})

	// edit extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "edit",
		Description: "Auto-generated stub for edit",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "edit"), nil
		},
	})

	// mcp_my-server_tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_my-server_tool",
		Description: "Auto-generated stub for mcp_my-server_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_my-server_tool"), nil
		},
	})

	// safe-tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "safe-tool",
		Description: "Auto-generated stub for safe-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "safe-tool"), nil
		},
	})

	// any-other-tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "any-other-tool",
		Description: "Auto-generated stub for any-other-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "any-other-tool"), nil
		},
	})

	// interactive-tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "interactive-tool",
		Description: "Auto-generated stub for interactive-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "interactive-tool"), nil
		},
	})

	// allowed-tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "allowed-tool",
		Description: "Auto-generated stub for allowed-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "allowed-tool"), nil
		},
	})

	// unknown-tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "unknown-tool",
		Description: "Auto-generated stub for unknown-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "unknown-tool"), nil
		},
	})

	// any-tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "any-tool",
		Description: "Auto-generated stub for any-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "any-tool"), nil
		},
	})

	// dangerous-tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "dangerous-tool",
		Description: "Auto-generated stub for dangerous-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "dangerous-tool"), nil
		},
	})

	// new-tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "new-tool",
		Description: "Auto-generated stub for new-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "new-tool"), nil
		},
	})

	// mcp_mcp_tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_mcp_tool",
		Description: "Auto-generated stub for mcp_mcp_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_mcp_tool"), nil
		},
	})

	// mcp_other_tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_other_tool",
		Description: "Auto-generated stub for mcp_other_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_other_tool"), nil
		},
	})

	// mcp_my-server_tool1 extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_my-server_tool1",
		Description: "Auto-generated stub for mcp_my-server_tool1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_my-server_tool1"), nil
		},
	})

	// mcp_my-server_another_tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_my-server_another_tool",
		Description: "Auto-generated stub for mcp_my-server_another_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_my-server_another_tool"), nil
		},
	})

	// mcp_blocked-server_tool1 extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_blocked-server_tool1",
		Description: "Auto-generated stub for mcp_blocked-server_tool1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_blocked-server_tool1"), nil
		},
	})

	// mcp_blocked-server_dangerous extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_blocked-server_dangerous",
		Description: "Auto-generated stub for mcp_blocked-server_dangerous",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_blocked-server_dangerous"), nil
		},
	})

	// mcp_other-server_tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_other-server_tool",
		Description: "Auto-generated stub for mcp_other-server_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_other-server_tool"), nil
		},
	})

	// my-server-tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-server-tool",
		Description: "Auto-generated stub for my-server-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-server-tool"), nil
		},
	})

	// my-server extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-server",
		Description: "Auto-generated stub for my-server",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-server"), nil
		},
	})

	// mcp_my-server_dangerous-tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_my-server_dangerous-tool",
		Description: "Auto-generated stub for mcp_my-server_dangerous-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_my-server_dangerous-tool"), nil
		},
	})

	// mcp_my-server_safe-tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_my-server_safe-tool",
		Description: "Auto-generated stub for mcp_my-server_safe-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_my-server_safe-tool"), nil
		},
	})

	// mcp_mcp_safe_server_malicious_tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_mcp_safe_server_malicious_tool",
		Description: "Auto-generated stub for mcp_mcp_safe_server_malicious_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_mcp_safe_server_malicious_tool"), nil
		},
	})

	// mcp_other_server_tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_other_server_tool",
		Description: "Auto-generated stub for mcp_other_server_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_other_server_tool"), nil
		},
	})

	// mcp_safe_server_tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_safe_server_tool",
		Description: "Auto-generated stub for mcp_safe_server_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_safe_server_tool"), nil
		},
	})

	// read extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "read",
		Description: "Auto-generated stub for read",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "read"), nil
		},
	})

	// api extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "api",
		Description: "Auto-generated stub for api",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "api"), nil
		},
	})

	// deep extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "deep",
		Description: "Auto-generated stub for deep",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "deep"), nil
		},
	})

	// unknown_subagent extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "unknown_subagent",
		Description: "Auto-generated stub for unknown_subagent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "unknown_subagent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "unknown_subagent", result), nil
		},
	})

	// checker1 extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "checker1",
		Description: "Auto-generated stub for checker1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "checker1"), nil
		},
	})

	// checker2 extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "checker2",
		Description: "Auto-generated stub for checker2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "checker2"), nil
		},
	})

	// matching extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "matching",
		Description: "Auto-generated stub for matching",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "matching"), nil
		},
	})

	// non-matching extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "non-matching",
		Description: "Auto-generated stub for non-matching",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "non-matching"), nil
		},
	})

	// global extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "global",
		Description: "Auto-generated stub for global",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "global"), nil
		},
	})

	// any_tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "any_tool",
		Description: "Auto-generated stub for any_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "any_tool"), nil
		},
	})

	// mcp_server_tool extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_server_tool",
		Description: "Auto-generated stub for mcp_server_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_server_tool"), nil
		},
	})

	// wildcard extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "wildcard",
		Description: "Auto-generated stub for wildcard",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "wildcard"), nil
		},
	})

	// activate_skill extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "activate_skill",
		Description: "Auto-generated stub for activate_skill",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "activate_skill"), nil
		},
	})

	// enter_plan_mode extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "enter_plan_mode",
		Description: "Auto-generated stub for enter_plan_mode",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "enter_plan_mode"), nil
		},
	})

	// c1 extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "c1",
		Description: "Auto-generated stub for c1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "c1"), nil
		},
	})

	// c2 extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "c2",
		Description: "Auto-generated stub for c2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "c2"), nil
		},
	})

	// c3 extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "c3",
		Description: "Auto-generated stub for c3",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "c3"), nil
		},
	})

	// write extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "write",
		Description: "Auto-generated stub for write",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "write"), nil
		},
	})

	// mcp_mcp_test extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_mcp_test",
		Description: "Auto-generated stub for mcp_mcp_test",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_mcp_test"), nil
		},
	})

	// mcp_mcp_stable extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_mcp_stable",
		Description: "Auto-generated stub for mcp_mcp_stable",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_mcp_stable"), nil
		},
	})

	// h1 extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "h1",
		Description: "Auto-generated stub for h1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "h1"), nil
		},
	})

	// h2 extracted from gemini-cli/packages/core/src/policy/policy-engine.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "h2",
		Description: "Auto-generated stub for h2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "h2"), nil
		},
	})

	// github__list_issues extracted from gemini-cli/packages/core/src/policy/toml-loader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "github__list_issues",
		Description: "Auto-generated stub for github__list_issues",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "github__list_issues"), nil
		},
	})

	// mcp_github_create_issue extracted from gemini-cli/packages/core/src/policy/toml-loader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_github_create_issue",
		Description: "Auto-generated stub for mcp_github_create_issue",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_github_create_issue"), nil
		},
	})

	// mcp_github_delete_issue extracted from gemini-cli/packages/core/src/policy/toml-loader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_github_delete_issue",
		Description: "Auto-generated stub for mcp_github_delete_issue",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_github_delete_issue"), nil
		},
	})

	// mcp_github_list_repos extracted from gemini-cli/packages/core/src/policy/toml-loader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_github_list_repos",
		Description: "Auto-generated stub for mcp_github_list_repos",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_github_list_repos"), nil
		},
	})

	// some_random_tool extracted from gemini-cli/packages/core/src/policy/toml-loader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "some_random_tool",
		Description: "Auto-generated stub for some_random_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "some_random_tool"), nil
		},
	})

	// list_directory extracted from gemini-cli/packages/core/src/policy/memory-manager-policy.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "list_directory",
		Description: "Auto-generated stub for list_directory",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "list_directory"), nil
		},
	})

	// python-checker extracted from gemini-cli/packages/core/src/safety/checker-runner.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "python-checker",
		Description: "Auto-generated stub for python-checker",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "python-checker"), nil
		},
	})

	// unknownTool extracted from gemini-cli/packages/core/src/safety/conseca/policy-enforcer.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "unknownTool",
		Description: "Auto-generated stub for unknownTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "unknownTool"), nil
		},
	})

	// prompt1 extracted from gemini-cli/packages/core/src/prompts/prompt-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "prompt1",
		Description: "Auto-generated stub for prompt1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "prompt1"), nil
		},
	})

	// prompt2 extracted from gemini-cli/packages/core/src/prompts/prompt-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "prompt2",
		Description: "Auto-generated stub for prompt2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "prompt2"), nil
		},
	})

	// server2_prompt1 extracted from gemini-cli/packages/core/src/prompts/prompt-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "server2_prompt1",
		Description: "Auto-generated stub for server2_prompt1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "server2_prompt1"), nil
		},
	})

	// test-function extracted from gemini-cli/packages/core/src/telemetry/semantic.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-function",
		Description: "Auto-generated stub for test-function",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-function"), nil
		},
	})

	// test_phase extracted from gemini-cli/packages/core/src/telemetry/startupProfiler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test_phase",
		Description: "Auto-generated stub for test_phase",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test_phase"), nil
		},
	})

	// test-span extracted from gemini-cli/packages/core/src/telemetry/file-exporters.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-span",
		Description: "Auto-generated stub for test-span",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-span"), nil
		},
	})

	// circular-span extracted from gemini-cli/packages/core/src/telemetry/file-exporters.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "circular-span",
		Description: "Auto-generated stub for circular-span",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "circular-span"), nil
		},
	})

	// ext-one extracted from gemini-cli/packages/core/src/telemetry/loggers.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ext-one",
		Description: "Auto-generated stub for ext-one",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ext-one"), nil
		},
	})

	// ext-two extracted from gemini-cli/packages/core/src/telemetry/loggers.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ext-two",
		Description: "Auto-generated stub for ext-two",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ext-two"), nil
		},
	})

	// test-worktree extracted from gemini-cli/packages/core/src/telemetry/loggers.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-worktree",
		Description: "Auto-generated stub for test-worktree",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-worktree"), nil
		},
	})

	// mock_mcp_tool extracted from gemini-cli/packages/core/src/telemetry/loggers.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mock_mcp_tool",
		Description: "Auto-generated stub for mock_mcp_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mock_mcp_tool"), nil
		},
	})

	// mock_mcp_server extracted from gemini-cli/packages/core/src/telemetry/loggers.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mock_mcp_server",
		Description: "Auto-generated stub for mock_mcp_server",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mock_mcp_server"), nil
		},
	})

	// testing extracted from gemini-cli/packages/core/src/telemetry/loggers.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "testing",
		Description: "Auto-generated stub for testing",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "testing"), nil
		},
	})

	// before-tool extracted from gemini-cli/packages/core/src/telemetry/loggers.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "before-tool",
		Description: "Auto-generated stub for before-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "before-tool"), nil
		},
	})

	// new_space extracted from gemini-cli/packages/core/src/telemetry/memory-monitor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "new_space",
		Description: "Auto-generated stub for new_space",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "new_space"), nil
		},
	})

	// old_space extracted from gemini-cli/packages/core/src/telemetry/memory-monitor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "old_space",
		Description: "Auto-generated stub for old_space",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "old_space"), nil
		},
	})

	// fn extracted from gemini-cli/packages/core/src/telemetry/metrics.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "fn",
		Description: "Auto-generated stub for fn",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "fn"), nil
		},
	})

	// my-fn extracted from gemini-cli/packages/core/src/telemetry/metrics.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-fn",
		Description: "Auto-generated stub for my-fn",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-fn"), nil
		},
	})

	// Bash extracted from gemini-cli/packages/core/src/telemetry/metrics.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Bash",
		Description: "Auto-generated stub for Bash",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Bash"), nil
		},
	})

	// AfterTool extracted from gemini-cli/packages/core/src/telemetry/metrics.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "AfterTool",
		Description: "Auto-generated stub for AfterTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "AfterTool"), nil
		},
	})

	// type extracted from gemini-cli/packages/core/src/telemetry/metrics.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "type",
		Description: "Auto-generated stub for type",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "type"), nil
		},
	})

	// Devin extracted from gemini-cli/packages/core/src/telemetry/clearcut-logger/clearcut-logger.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Devin",
		Description: "Auto-generated stub for Devin",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Devin"), nil
		},
	})

	// unidentified extracted from gemini-cli/packages/core/src/telemetry/clearcut-logger/clearcut-logger.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "unidentified",
		Description: "Auto-generated stub for unidentified",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "unidentified"), nil
		},
	})

	// some_other_tool extracted from gemini-cli/packages/core/src/telemetry/clearcut-logger/clearcut-logger.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "some_other_tool",
		Description: "Auto-generated stub for some_other_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "some_other_tool"), nil
		},
	})

	// CONCORD extracted from gemini-cli/packages/core/src/telemetry/clearcut-logger/clearcut-logger.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "CONCORD",
		Description: "Auto-generated stub for CONCORD",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "CONCORD"), nil
		},
	})

	// calculate extracted from gemini-cli/packages/core/src/code_assist/converter.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "calculate",
		Description: "Auto-generated stub for calculate",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "calculate"), nil
		},
	})

	// Free extracted from gemini-cli/packages/core/src/code_assist/server.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Free",
		Description: "free tier",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Free"), nil
		},
	})

	// tier extracted from gemini-cli/packages/core/src/code_assist/server.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tier",
		Description: "Auto-generated stub for tier",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tier"), nil
		},
	})

	// t1 extracted from gemini-cli/packages/core/src/code_assist/telemetry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "t1",
		Description: "Auto-generated stub for t1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "t1"), nil
		},
	})

	// t2 extracted from gemini-cli/packages/core/src/code_assist/telemetry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "t2",
		Description: "Auto-generated stub for t2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "t2"), nil
		},
	})

	// paid extracted from gemini-cli/packages/core/src/code_assist/setup.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "paid",
		Description: "Paid tier",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "paid"), nil
		},
	})

	// free extracted from gemini-cli/packages/core/src/code_assist/setup.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "free",
		Description: "Free tier",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "free"), nil
		},
	})

	// user-gemini extracted from gemini-cli/packages/core/src/skills/skillManagerAlias.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "user-gemini",
		Description: "desc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "user-gemini"), nil
		},
	})

	// user-agent extracted from gemini-cli/packages/core/src/skills/skillManagerAlias.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "user-agent",
		Description: "desc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "user-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "user-agent", result), nil
		},
	})

	// project-gemini extracted from gemini-cli/packages/core/src/skills/skillManagerAlias.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "project-gemini",
		Description: "desc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "project-gemini"), nil
		},
	})

	// project-agent extracted from gemini-cli/packages/core/src/skills/skillManagerAlias.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "project-agent",
		Description: "desc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "project-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "project-agent", result), nil
		},
	})

	// same-skill extracted from gemini-cli/packages/core/src/skills/skillManagerAlias.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "same-skill",
		Description: "gemini-desc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "same-skill"), nil
		},
	})

	// test-ext extracted from gemini-cli/packages/core/src/skills/skillManager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-ext",
		Description: "ext-desc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-ext"), nil
		},
	})

	// skill-extension extracted from gemini-cli/packages/core/src/skills/skillManager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "skill-extension",
		Description: "ext-desc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "skill-extension"), nil
		},
	})

	// same-name extracted from gemini-cli/packages/core/src/skills/skillManager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "same-name",
		Description: "ext-desc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "same-name"), nil
		},
	})

	// builtin-skill extracted from gemini-cli/packages/core/src/skills/skillManager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "builtin-skill",
		Description: "builtin-desc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "builtin-skill"), nil
		},
	})

	// regular-skill extracted from gemini-cli/packages/core/src/skills/skillManager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "regular-skill",
		Description: "regular",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "regular-skill"), nil
		},
	})

	// disabled-builtin extracted from gemini-cli/packages/core/src/skills/skillManager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "disabled-builtin",
		Description: "disabled builtin",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "disabled-builtin"), nil
		},
	})

	// web_fetch extracted from gemini-cli/packages/core/src/tools/web-fetch.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "web_fetch",
		Description: "Auto-generated stub for web_fetch",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "web_fetch"), nil
		},
	})

	// test-server extracted from gemini-cli/packages/core/src/tools/mcp-client-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-server",
		Description: "Auto-generated stub for test-server",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-server"), nil
		},
	})

	// extension-1 extracted from gemini-cli/packages/core/src/tools/mcp-client-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "extension-1",
		Description: "Auto-generated stub for extension-1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "extension-1"), nil
		},
	})

	// extension-2 extracted from gemini-cli/packages/core/src/tools/mcp-client-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "extension-2",
		Description: "Auto-generated stub for extension-2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "extension-2"), nil
		},
	})

	// blocked-server extracted from gemini-cli/packages/core/src/tools/mcp-client-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "blocked-server",
		Description: "Auto-generated stub for blocked-server",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "blocked-server"), nil
		},
	})

	// non-existent extracted from gemini-cli/packages/core/src/tools/activate-skill.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "non-existent",
		Description: "Auto-generated stub for non-existent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "non-existent"), nil
		},
	})

	// gemini-cli-mcp-client extracted from gemini-cli/packages/core/src/tools/mcp-client.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "gemini-cli-mcp-client",
		Description: "Auto-generated stub for gemini-cli-mcp-client",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "gemini-cli-mcp-client"), nil
		},
	})

	// FILE_NOT_FOUND extracted from gemini-cli/packages/core/src/tools/edit.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "FILE_NOT_FOUND",
		Description: "Auto-generated stub for FILE_NOT_FOUND",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "FILE_NOT_FOUND"), nil
		},
	})

	// ATTEMPT_TO_CREATE_EXISTING_FILE extracted from gemini-cli/packages/core/src/tools/edit.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ATTEMPT_TO_CREATE_EXISTING_FILE",
		Description: "Auto-generated stub for ATTEMPT_TO_CREATE_EXISTING_FILE",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ATTEMPT_TO_CREATE_EXISTING_FILE"), nil
		},
	})

	// NO_OCCURRENCE_FOUND extracted from gemini-cli/packages/core/src/tools/edit.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "NO_OCCURRENCE_FOUND",
		Description: "Auto-generated stub for NO_OCCURRENCE_FOUND",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "NO_OCCURRENCE_FOUND"), nil
		},
	})

	// EXPECTED_OCCURRENCE_MISMATCH extracted from gemini-cli/packages/core/src/tools/edit.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "EXPECTED_OCCURRENCE_MISMATCH",
		Description: "Auto-generated stub for EXPECTED_OCCURRENCE_MISMATCH",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "EXPECTED_OCCURRENCE_MISMATCH"), nil
		},
	})

	// WriteFileTool extracted from gemini-cli/packages/core/src/tools/confirmation-policy.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "WriteFileTool",
		Description: "Auto-generated stub for WriteFileTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "WriteFileTool"), nil
		},
	})

	// WebFetchTool extracted from gemini-cli/packages/core/src/tools/confirmation-policy.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "WebFetchTool",
		Description: "Auto-generated stub for WebFetchTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "WebFetchTool"), nil
		},
	})

	// testFunction extracted from gemini-cli/packages/core/src/tools/mcp-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "testFunction",
		Description: "Auto-generated stub for testFunction",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "testFunction"), nil
		},
	})

	// validTool extracted from gemini-cli/packages/core/src/tools/mcp-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "validTool",
		Description: "a param with no type",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "validTool"), nil
		},
	})

	// invalidTool extracted from gemini-cli/packages/core/src/tools/mcp-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "invalidTool",
		Description: "a param with no type",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "invalidTool"), nil
		},
	})

	// readOnlyTool extracted from gemini-cli/packages/core/src/tools/mcp-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "readOnlyTool",
		Description: "A read-only tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "readOnlyTool"), nil
		},
	})

	// writeTool extracted from gemini-cli/packages/core/src/tools/mcp-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "writeTool",
		Description: "A write tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "writeTool"), nil
		},
	})

	// multiAnnotationTool extracted from gemini-cli/packages/core/src/tools/mcp-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "multiAnnotationTool",
		Description: "A tool with multiple annotations",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "multiAnnotationTool"), nil
		},
	})

	// toolWithDefs extracted from gemini-cli/packages/core/src/tools/mcp-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "toolWithDefs",
		Description: "A tool using $defs",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "toolWithDefs"), nil
		},
	})

	// resource extracted from gemini-cli/packages/core/src/tools/mcp-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "resource",
		Description: "Test Resource",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "resource"), nil
		},
	})

	// one extracted from gemini-cli/packages/core/src/tools/mcp-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "one",
		Description: "first",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "one"), nil
		},
	})

	// two extracted from gemini-cli/packages/core/src/tools/mcp-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "two",
		Description: "second",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "two"), nil
		},
	})

	// newTool extracted from gemini-cli/packages/core/src/tools/mcp-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "newTool",
		Description: "Auto-generated stub for newTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "newTool"), nil
		},
	})

	// ext-setting extracted from gemini-cli/packages/core/src/tools/mcp-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ext-setting",
		Description: "Auto-generated stub for ext-setting",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ext-setting"), nil
		},
	})

	// excluded-tool-class extracted from gemini-cli/packages/core/src/tools/tool-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "excluded-tool-class",
		Description: "Auto-generated stub for excluded-tool-class",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "excluded-tool-class"), nil
		},
	})

	// current_test_tool extracted from gemini-cli/packages/core/src/tools/tool-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "current_test_tool",
		Description: "Auto-generated stub for current_test_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "current_test_tool"), nil
		},
	})

	// c-tool extracted from gemini-cli/packages/core/src/tools/tool-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "c-tool",
		Description: "Auto-generated stub for c-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "c-tool"), nil
		},
	})

	// a-tool extracted from gemini-cli/packages/core/src/tools/tool-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "a-tool",
		Description: "Auto-generated stub for a-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "a-tool"), nil
		},
	})

	// b-tool extracted from gemini-cli/packages/core/src/tools/tool-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "b-tool",
		Description: "Auto-generated stub for b-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "b-tool"), nil
		},
	})

	// regular-tool extracted from gemini-cli/packages/core/src/tools/tool-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "regular-tool",
		Description: "Auto-generated stub for regular-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "regular-tool"), nil
		},
	})

	// builtin-1 extracted from gemini-cli/packages/core/src/tools/tool-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "builtin-1",
		Description: "Auto-generated stub for builtin-1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "builtin-1"), nil
		},
	})

	// builtin-2 extracted from gemini-cli/packages/core/src/tools/tool-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "builtin-2",
		Description: "Auto-generated stub for builtin-2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "builtin-2"), nil
		},
	})

	// tool-with-bad-format extracted from gemini-cli/packages/core/src/tools/tool-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tool-with-bad-format",
		Description: "A tool with an invalid format property",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tool-with-bad-format"), nil
		},
	})

	// failing-tool extracted from gemini-cli/packages/core/src/tools/tool-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "failing-tool",
		Description: "A tool that fails",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "failing-tool"), nil
		},
	})

	// policy-test-tool extracted from gemini-cli/packages/core/src/tools/tool-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "policy-test-tool",
		Description: "tests policy",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "policy-test-tool"), nil
		},
	})

	// resource-name extracted from gemini-cli/packages/core/src/tools/mcp-tool.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "resource-name",
		Description: "Auto-generated stub for resource-name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "resource-name"), nil
		},
	})

	// error extracted from gemini-cli/packages/core/src/tools/mcp-tool.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "error",
		Description: "Auto-generated stub for error",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "error"), nil
		},
	})

	// grep_search_ripgrep extracted from gemini-cli/packages/core/src/tools/definitions/coreToolsModelSnapshots.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "grep_search_ripgrep",
		Description: "Auto-generated stub for grep_search_ripgrep",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "grep_search_ripgrep"), nil
		},
	})

	// google_web_search extracted from gemini-cli/packages/core/src/tools/definitions/coreToolsModelSnapshots.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "google_web_search",
		Description: "Auto-generated stub for google_web_search",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "google_web_search"), nil
		},
	})

	// read_many_files extracted from gemini-cli/packages/core/src/tools/definitions/coreToolsModelSnapshots.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "read_many_files",
		Description: "Auto-generated stub for read_many_files",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "read_many_files"), nil
		},
	})

	// write_todos extracted from gemini-cli/packages/core/src/tools/definitions/coreToolsModelSnapshots.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "write_todos",
		Description: "Auto-generated stub for write_todos",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "write_todos"), nil
		},
	})

	// get_internal_docs extracted from gemini-cli/packages/core/src/tools/definitions/coreToolsModelSnapshots.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "get_internal_docs",
		Description: "Auto-generated stub for get_internal_docs",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "get_internal_docs"), nil
		},
	})

	// activate_skill_empty extracted from gemini-cli/packages/core/src/tools/definitions/coreToolsModelSnapshots.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "activate_skill_empty",
		Description: "Auto-generated stub for activate_skill_empty",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "activate_skill_empty"), nil
		},
	})

	// activate_skill_single extracted from gemini-cli/packages/core/src/tools/definitions/coreToolsModelSnapshots.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "activate_skill_single",
		Description: "Auto-generated stub for activate_skill_single",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "activate_skill_single"), nil
		},
	})

	// grep extracted from gemini-cli/packages/core/src/context/chatCompressionService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "grep",
		Description: "Auto-generated stub for grep",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "grep"), nil
		},
	})

	// raw_tool extracted from gemini-cli/packages/core/src/context/chatCompressionService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "raw_tool",
		Description: "Auto-generated stub for raw_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "raw_tool"), nil
		},
	})

	// massive_preserved extracted from gemini-cli/packages/core/src/context/chatCompressionService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "massive_preserved",
		Description: "Auto-generated stub for massive_preserved",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "massive_preserved"), nil
		},
	})

	// t3 extracted from gemini-cli/packages/core/src/context/toolOutputMaskingService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "t3",
		Description: "Auto-generated stub for t3",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "t3"), nil
		},
	})

	// p extracted from gemini-cli/packages/core/src/context/toolOutputMaskingService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "p",
		Description: "Auto-generated stub for p",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "p"), nil
		},
	})

	// l extracted from gemini-cli/packages/core/src/context/toolOutputMaskingService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "l",
		Description: "Auto-generated stub for l",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "l"), nil
		},
	})

	// padding extracted from gemini-cli/packages/core/src/context/toolOutputMaskingService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "padding",
		Description: "Auto-generated stub for padding",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "padding"), nil
		},
	})

	// tiny_tool extracted from gemini-cli/packages/core/src/context/toolOutputMaskingService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tiny_tool",
		Description: "Auto-generated stub for tiny_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tiny_tool"), nil
		},
	})

	// bulky_tool extracted from gemini-cli/packages/core/src/context/toolOutputMaskingService.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "bulky_tool",
		Description: "Auto-generated stub for bulky_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "bulky_tool"), nil
		},
	})

	// ModifyingProcessor extracted from gemini-cli/packages/core/src/context/pipeline/orchestrator.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ModifyingProcessor",
		Description: "Auto-generated stub for ModifyingProcessor",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ModifyingProcessor"), nil
		},
	})

	// Throwing extracted from gemini-cli/packages/core/src/context/pipeline/orchestrator.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Throwing",
		Description: "Auto-generated stub for Throwing",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Throwing"), nil
		},
	})

	// MockAsyncProcessor extracted from gemini-cli/packages/core/src/context/pipeline/orchestrator.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "MockAsyncProcessor",
		Description: "Auto-generated stub for MockAsyncProcessor",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "MockAsyncProcessor"), nil
		},
	})

	// TestPipeline extracted from gemini-cli/packages/core/src/context/pipeline/orchestrator.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "TestPipeline",
		Description: "Auto-generated stub for TestPipeline",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "TestPipeline"), nil
		},
	})

	// FailingPipeline extracted from gemini-cli/packages/core/src/context/pipeline/orchestrator.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "FailingPipeline",
		Description: "Auto-generated stub for FailingPipeline",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "FailingPipeline"), nil
		},
	})

	// TestAsync extracted from gemini-cli/packages/core/src/context/pipeline/orchestrator.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "TestAsync",
		Description: "Auto-generated stub for TestAsync",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "TestAsync"), nil
		},
	})

	// Normalization extracted from gemini-cli/packages/core/src/context/config/profiles.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Normalization",
		Description: "Auto-generated stub for Normalization",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Normalization"), nil
		},
	})

	// Async extracted from gemini-cli/packages/core/src/context/system-tests/lifecycle.golden.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Async",
		Description: "Auto-generated stub for Async",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Async"), nil
		},
	})

	// NodeTruncationProcessor extracted from gemini-cli/packages/core/src/context/processors/nodeTruncationProcessor.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "NodeTruncationProcessor",
		Description: "Auto-generated stub for NodeTruncationProcessor",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "NodeTruncationProcessor"), nil
		},
	})

	// StateSnapshotAsyncProcessor extracted from gemini-cli/packages/core/src/context/processors/stateSnapshotAsyncProcessor.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "StateSnapshotAsyncProcessor",
		Description: "Auto-generated stub for StateSnapshotAsyncProcessor",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "StateSnapshotAsyncProcessor"), nil
		},
	})

	// HistoryTruncationProcessor extracted from gemini-cli/packages/core/src/context/processors/historyTruncationProcessor.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "HistoryTruncationProcessor",
		Description: "Auto-generated stub for HistoryTruncationProcessor",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "HistoryTruncationProcessor"), nil
		},
	})

	// RollingSummaryProcessor extracted from gemini-cli/packages/core/src/context/processors/rollingSummaryProcessor.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "RollingSummaryProcessor",
		Description: "Auto-generated stub for RollingSummaryProcessor",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "RollingSummaryProcessor"), nil
		},
	})

	// ToolMaskingProcessor extracted from gemini-cli/packages/core/src/context/processors/toolMaskingProcessor.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ToolMaskingProcessor",
		Description: "Auto-generated stub for ToolMaskingProcessor",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ToolMaskingProcessor"), nil
		},
	})

	// StateSnapshotProcessor extracted from gemini-cli/packages/core/src/context/processors/stateSnapshotProcessor.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "StateSnapshotProcessor",
		Description: "Auto-generated stub for StateSnapshotProcessor",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "StateSnapshotProcessor"), nil
		},
	})

	// NodeDistillationProcessor extracted from gemini-cli/packages/core/src/context/processors/nodeDistillationProcessor.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "NodeDistillationProcessor",
		Description: "Auto-generated stub for NodeDistillationProcessor",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "NodeDistillationProcessor"), nil
		},
	})

	// BlobDegradationProcessor extracted from gemini-cli/packages/core/src/context/processors/blobDegradationProcessor.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "BlobDegradationProcessor",
		Description: "Auto-generated stub for BlobDegradationProcessor",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "BlobDegradationProcessor"), nil
		},
	})

	// streamable-http-client extracted from gemini-cli/packages/core/src/ide/ide-client.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "streamable-http-client",
		Description: "Auto-generated stub for streamable-http-client",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "streamable-http-client"), nil
		},
	})

	// stdio-client extracted from gemini-cli/packages/core/src/ide/ide-client.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "stdio-client",
		Description: "Auto-generated stub for stdio-client",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "stdio-client"), nil
		},
	})

	// custom-ide extracted from gemini-cli/packages/core/src/ide/detect-ide.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "custom-ide",
		Description: "Auto-generated stub for custom-ide",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "custom-ide"), nil
		},
	})

	// devin extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "devin",
		Description: "Auto-generated stub for devin",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "devin"), nil
		},
	})

	// replit extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "replit",
		Description: "Auto-generated stub for replit",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "replit"), nil
		},
	})

	// cursor extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "cursor",
		Description: "Auto-generated stub for cursor",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "cursor"), nil
		},
	})

	// cloudshell extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "cloudshell",
		Description: "Auto-generated stub for cloudshell",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "cloudshell"), nil
		},
	})

	// codespaces extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "codespaces",
		Description: "Auto-generated stub for codespaces",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "codespaces"), nil
		},
	})

	// firebasestudio extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "firebasestudio",
		Description: "Auto-generated stub for firebasestudio",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "firebasestudio"), nil
		},
	})

	// trae extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "trae",
		Description: "Auto-generated stub for trae",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "trae"), nil
		},
	})

	// vscode extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "vscode",
		Description: "Auto-generated stub for vscode",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "vscode"), nil
		},
	})

	// vscodefork extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "vscodefork",
		Description: "Auto-generated stub for vscodefork",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "vscodefork"), nil
		},
	})

	// positron extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "positron",
		Description: "Auto-generated stub for positron",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "positron"), nil
		},
	})

	// antigravity extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "antigravity",
		Description: "Auto-generated stub for antigravity",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "antigravity"), nil
		},
	})

	// sublimetext extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "sublimetext",
		Description: "Auto-generated stub for sublimetext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "sublimetext"), nil
		},
	})

	// jetbrains extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "jetbrains",
		Description: "Auto-generated stub for jetbrains",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "jetbrains"), nil
		},
	})

	// intellijidea extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "intellijidea",
		Description: "Auto-generated stub for intellijidea",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "intellijidea"), nil
		},
	})

	// webstorm extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "webstorm",
		Description: "Auto-generated stub for webstorm",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "webstorm"), nil
		},
	})

	// pycharm extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "pycharm",
		Description: "Auto-generated stub for pycharm",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "pycharm"), nil
		},
	})

	// goland extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "goland",
		Description: "Auto-generated stub for goland",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "goland"), nil
		},
	})

	// androidstudio extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "androidstudio",
		Description: "Auto-generated stub for androidstudio",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "androidstudio"), nil
		},
	})

	// clion extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "clion",
		Description: "Auto-generated stub for clion",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "clion"), nil
		},
	})

	// rustrover extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "rustrover",
		Description: "Auto-generated stub for rustrover",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "rustrover"), nil
		},
	})

	// datagrip extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "datagrip",
		Description: "Auto-generated stub for datagrip",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "datagrip"), nil
		},
	})

	// phpstorm extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "phpstorm",
		Description: "Auto-generated stub for phpstorm",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "phpstorm"), nil
		},
	})

	// zed extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "zed",
		Description: "Auto-generated stub for zed",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "zed"), nil
		},
	})

	// xcode extracted from gemini-cli/packages/core/src/ide/detect-ide.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "xcode",
		Description: "Auto-generated stub for xcode",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "xcode"), nil
		},
	})

	// someOtherTool extracted from gemini-cli/packages/core/src/ide/ide-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "someOtherTool",
		Description: "Auto-generated stub for someOtherTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "someOtherTool"), nil
		},
	})

	// openDiff extracted from gemini-cli/packages/core/src/ide/ide-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "openDiff",
		Description: "Auto-generated stub for openDiff",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "openDiff"), nil
		},
	})

	// closeDiff extracted from gemini-cli/packages/core/src/ide/ide-client.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "closeDiff",
		Description: "Auto-generated stub for closeDiff",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "closeDiff"), nil
		},
	})

	// raw-name extracted from gemini-cli/packages/core/src/agent/tool-display-utils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "raw-name",
		Description: "Auto-generated stub for raw-name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "raw-name"), nil
		},
	})

	// myFunc extracted from gemini-cli/packages/core/src/agent/content-utils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "myFunc",
		Description: "Auto-generated stub for myFunc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "myFunc"), nil
		},
	})

	// extensions extracted from gemini-cli/packages/a2a-server/src/commands/command-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "extensions",
		Description: "Manage extensions.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "extensions"), nil
		},
	})

	// restore extracted from gemini-cli/packages/a2a-server/src/commands/command-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "restore",
		Description: "Restores the server.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "restore"), nil
		},
	})

	// test-command extracted from gemini-cli/packages/a2a-server/src/commands/command-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-command",
		Description: "",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-command"), nil
		},
	})

	// test-command-sub-sub extracted from gemini-cli/packages/a2a-server/src/commands/command-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-command-sub-sub",
		Description: "",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-command-sub-sub"), nil
		},
	})

	// test-command-sub extracted from gemini-cli/packages/a2a-server/src/commands/command-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-command-sub",
		Description: "",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-command-sub"), nil
		},
	})

	// cyclic-command extracted from gemini-cli/packages/a2a-server/src/commands/command-registry.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "cyclic-command",
		Description: "",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "cyclic-command"), nil
		},
	})

	// test-tool-1 extracted from gemini-cli/packages/a2a-server/src/http/app.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-tool-1",
		Description: "Auto-generated stub for test-tool-1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-tool-1"), nil
		},
	})

	// test-tool-2 extracted from gemini-cli/packages/a2a-server/src/http/app.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-tool-2",
		Description: "Auto-generated stub for test-tool-2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-tool-2"), nil
		},
	})

	// test-tool-no-approval extracted from gemini-cli/packages/a2a-server/src/http/app.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-tool-no-approval",
		Description: "Auto-generated stub for test-tool-no-approval",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-tool-no-approval"), nil
		},
	})

	// test-tool-yolo extracted from gemini-cli/packages/a2a-server/src/http/app.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-tool-yolo",
		Description: "Auto-generated stub for test-tool-yolo",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-tool-yolo"), nil
		},
	})

	// arg1 extracted from gemini-cli/packages/a2a-server/src/http/app.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "arg1",
		Description: "A test command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "arg1"), nil
		},
	})

	// sub-command extracted from gemini-cli/packages/a2a-server/src/http/app.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "sub-command",
		Description: "Argument 1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "sub-command"), nil
		},
	})

	// another-command extracted from gemini-cli/packages/a2a-server/src/http/app.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "another-command",
		Description: "Another test command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "another-command"), nil
		},
	})

	// not-top-level extracted from gemini-cli/packages/a2a-server/src/http/app.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "not-top-level",
		Description: "Not a top level command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "not-top-level"), nil
		},
	})

	// workspace-command extracted from gemini-cli/packages/a2a-server/src/http/app.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "workspace-command",
		Description: "A command that requires a workspace",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "workspace-command"), nil
		},
	})

	// context-check-command extracted from gemini-cli/packages/a2a-server/src/http/app.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "context-check-command",
		Description: "checks context",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "context-check-command"), nil
		},
	})

	// stream-test extracted from gemini-cli/packages/a2a-server/src/http/app.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "stream-test",
		Description: "A test streaming command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "stream-test"), nil
		},
	})

	// non-stream-test extracted from gemini-cli/packages/a2a-server/src/http/app.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "non-stream-test",
		Description: "A test non-streaming command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "non-stream-test"), nil
		},
	})

	// call_mcp_tool extracted from gemini-cli/packages/a2a-server/src/agent/task-event-driven.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "call_mcp_tool",
		Description: "Auto-generated stub for call_mcp_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "call_mcp_tool"), nil
		},
	})

	// writeFile extracted from gemini-cli/packages/a2a-server/src/agent/task-event-driven.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "writeFile",
		Description: "Auto-generated stub for writeFile",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "writeFile"), nil
		},
	})

	// pwd extracted from gemini-cli/packages/a2a-server/src/agent/task-event-driven.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "pwd",
		Description: "Auto-generated stub for pwd",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "pwd"), nil
		},
	})

	// not_restorable extracted from gemini-cli/packages/a2a-server/src/agent/task.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "not_restorable",
		Description: "Auto-generated stub for not_restorable",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "not_restorable"), nil
		},
	})

	// errorTool extracted from gemini-cli/packages/cli/src/nonInteractiveCli.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "errorTool",
		Description: "Auto-generated stub for errorTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "errorTool"), nil
		},
	})

	// nonexistentTool extracted from gemini-cli/packages/cli/src/nonInteractiveCli.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "nonexistentTool",
		Description: "Auto-generated stub for nonexistentTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "nonexistentTool"), nil
		},
	})

	// testcommand extracted from gemini-cli/packages/cli/src/nonInteractiveCli.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "testcommand",
		Description: "a test command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "testcommand"), nil
		},
	})

	// confirm extracted from gemini-cli/packages/cli/src/nonInteractiveCli.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "confirm",
		Description: "a command that needs confirmation",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "confirm"), nil
		},
	})

	// noaction extracted from gemini-cli/packages/cli/src/nonInteractiveCli.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "noaction",
		Description: "unhandled type",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "noaction"), nil
		},
	})

	// testargs extracted from gemini-cli/packages/cli/src/nonInteractiveCli.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "testargs",
		Description: "a test command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "testargs"), nil
		},
	})

	// ShellTool extracted from gemini-cli/packages/cli/src/nonInteractiveCli.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ShellTool",
		Description: "A shell tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ShellTool"), nil
		},
	})

	// stopTool extracted from gemini-cli/packages/cli/src/nonInteractiveCli.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "stopTool",
		Description: "Auto-generated stub for stopTool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "stopTool"), nil
		},
	})

	// my-extension extracted from gemini-cli/packages/cli/src/commands/extensions/disable.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-extension",
		Description: "Auto-generated stub for my-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-extension"), nil
		},
	})

	// S1 extracted from gemini-cli/packages/cli/src/commands/extensions/configure.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "S1",
		Description: "Auto-generated stub for S1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "S1"), nil
		},
	})

	// S2 extracted from gemini-cli/packages/cli/src/commands/extensions/configure.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "S2",
		Description: "Auto-generated stub for S2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "S2"), nil
		},
	})

	// custom-commands extracted from gemini-cli/packages/cli/src/commands/extensions/new.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "custom-commands",
		Description: "Auto-generated stub for custom-commands",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "custom-commands"), nil
		},
	})

	// mcp-server extracted from gemini-cli/packages/cli/src/commands/extensions/new.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp-server",
		Description: "Auto-generated stub for mcp-server",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp-server"), nil
		},
	})

	// mock-extension extracted from gemini-cli/packages/cli/src/commands/extensions/install.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mock-extension",
		Description: "Auto-generated stub for mock-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mock-extension"), nil
		},
	})

	// http-extension extracted from gemini-cli/packages/cli/src/commands/extensions/install.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "http-extension",
		Description: "Auto-generated stub for http-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "http-extension"), nil
		},
	})

	// https-extension extracted from gemini-cli/packages/cli/src/commands/extensions/install.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "https-extension",
		Description: "Auto-generated stub for https-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "https-extension"), nil
		},
	})

	// git-extension extracted from gemini-cli/packages/cli/src/commands/extensions/install.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "git-extension",
		Description: "Auto-generated stub for git-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "git-extension"), nil
		},
	})

	// sso-extension extracted from gemini-cli/packages/cli/src/commands/extensions/install.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "sso-extension",
		Description: "Auto-generated stub for sso-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "sso-extension"), nil
		},
	})

	// local-extension extracted from gemini-cli/packages/cli/src/commands/extensions/install.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "local-extension",
		Description: "Auto-generated stub for local-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "local-extension"), nil
		},
	})

	// local-ext-name extracted from gemini-cli/packages/cli/src/commands/extensions/validate.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "local-ext-name",
		Description: "Auto-generated stub for local-ext-name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "local-ext-name"), nil
		},
	})

	// INVALID_NAME extracted from gemini-cli/packages/cli/src/commands/extensions/validate.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "INVALID_NAME",
		Description: "Auto-generated stub for INVALID_NAME",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "INVALID_NAME"), nil
		},
	})

	// valid-name extracted from gemini-cli/packages/cli/src/commands/extensions/validate.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "valid-name",
		Description: "Auto-generated stub for valid-name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "valid-name"), nil
		},
	})

	// missing-extension extracted from gemini-cli/packages/cli/src/commands/extensions/update.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "missing-extension",
		Description: "Auto-generated stub for missing-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "missing-extension"), nil
		},
	})

	// my-linked-extension extracted from gemini-cli/packages/cli/src/commands/extensions/link.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-linked-extension",
		Description: "Auto-generated stub for my-linked-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-linked-extension"), nil
		},
	})

	// mcp-test-client extracted from gemini-cli/packages/cli/src/commands/mcp/list.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp-test-client",
		Description: "Auto-generated stub for mcp-test-client",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp-test-client"), nil
		},
	})

	// skill2 extracted from gemini-cli/packages/cli/src/commands/skills/list.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "skill2",
		Description: "desc2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "skill2"), nil
		},
	})

	// regular extracted from gemini-cli/packages/cli/src/commands/skills/list.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "regular",
		Description: "desc1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "regular"), nil
		},
	})

	// builtin extracted from gemini-cli/packages/cli/src/commands/skills/list.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "builtin",
		Description: "desc2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "builtin"), nil
		},
	})

	// deploy extracted from gemini-cli/packages/cli/src/services/FileCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "deploy",
		Description: "Auto-generated stub for deploy",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "deploy"), nil
		},
	})

	// active-ext extracted from gemini-cli/packages/cli/src/services/FileCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "active-ext",
		Description: "Auto-generated stub for active-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "active-ext"), nil
		},
	})

	// inactive-ext extracted from gemini-cli/packages/cli/src/services/FileCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "inactive-ext",
		Description: "Auto-generated stub for inactive-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "inactive-ext"), nil
		},
	})

	// no-commands extracted from gemini-cli/packages/cli/src/services/FileCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "no-commands",
		Description: "Auto-generated stub for no-commands",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "no-commands"), nil
		},
	})

	// my-test-ext extracted from gemini-cli/packages/cli/src/services/FileCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-test-ext",
		Description: "Auto-generated stub for my-test-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-test-ext"), nil
		},
	})

	// shorthand extracted from gemini-cli/packages/cli/src/services/FileCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "shorthand",
		Description: "Auto-generated stub for shorthand",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "shorthand"), nil
		},
	})

	// model_led extracted from gemini-cli/packages/cli/src/services/FileCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "model_led",
		Description: "Auto-generated stub for model_led",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "model_led"), nil
		},
	})

	// pipeline extracted from gemini-cli/packages/cli/src/services/FileCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "pipeline",
		Description: "Auto-generated stub for pipeline",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "pipeline"), nil
		},
	})

	// at-file extracted from gemini-cli/packages/cli/src/services/FileCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "at-file",
		Description: "Auto-generated stub for at-file",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "at-file"), nil
		},
	})

	// profile extracted from gemini-cli/packages/cli/src/services/BuiltinCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "profile",
		Description: "Profile command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "profile"), nil
		},
	})

	// about extracted from gemini-cli/packages/cli/src/services/BuiltinCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "about",
		Description: "About the CLI",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "about"), nil
		},
	})

	// chat extracted from gemini-cli/packages/cli/src/services/BuiltinCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "chat",
		Description: "Auto-generated stub for chat",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "chat"), nil
		},
	})

	// list extracted from gemini-cli/packages/cli/src/services/BuiltinCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "list",
		Description: "Auto-generated stub for list",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "list"), nil
		},
	})

	// save extracted from gemini-cli/packages/cli/src/services/BuiltinCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "save",
		Description: "Auto-generated stub for save",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "save"), nil
		},
	})

	// delete extracted from gemini-cli/packages/cli/src/services/BuiltinCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "delete",
		Description: "Auto-generated stub for delete",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "delete"), nil
		},
	})

	// share extracted from gemini-cli/packages/cli/src/services/BuiltinCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "share",
		Description: "Auto-generated stub for share",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "share"), nil
		},
	})

	// checkpoints extracted from gemini-cli/packages/cli/src/services/BuiltinCommandLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "checkpoints",
		Description: "Auto-generated stub for checkpoints",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "checkpoints"), nil
		},
	})

	// test-prompt extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-prompt",
		Description: "A test prompt.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-prompt"), nil
		},
	})

	// age extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "age",
		Description: "The animal's name.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "age"), nil
		},
	})

	// species extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "species",
		Description: "The animal's age.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "species"), nil
		},
	})

	// enclosure extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "enclosure",
		Description: "The animal's species.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "enclosure"), nil
		},
	})

	// trail extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "trail",
		Description: "The animal's enclosure.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "trail"), nil
		},
	})

	// arg2 extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "arg2",
		Description: "Auto-generated stub for arg2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "arg2"), nil
		},
	})

	// named extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "named",
		Description: "Auto-generated stub for named",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "named"), nil
		},
	})

	// pos extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "pos",
		Description: "Auto-generated stub for pos",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "pos"), nil
		},
	})

	// pos1 extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "pos1",
		Description: "Auto-generated stub for pos1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "pos1"), nil
		},
	})

	// pos2 extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "pos2",
		Description: "Auto-generated stub for pos2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "pos2"), nil
		},
	})

	// named1 extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "named1",
		Description: "Auto-generated stub for named1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "named1"), nil
		},
	})

	// named2 extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "named2",
		Description: "Auto-generated stub for named2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "named2"), nil
		},
	})

	// pos3 extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "pos3",
		Description: "Auto-generated stub for pos3",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "pos3"), nil
		},
	})

	// test-name extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-name",
		Description: "Auto-generated stub for test-name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-name"), nil
		},
	})

	// optional extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "optional",
		Description: "Auto-generated stub for optional",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "optional"), nil
		},
	})

	// required extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "required",
		Description: "Auto-generated stub for required",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "required"), nil
		},
	})

	// find extracted from gemini-cli/packages/cli/src/services/McpPromptLoader.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "find",
		Description: "Auto-generated stub for find",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "find"), nil
		},
	})

	// pickle extracted from gemini-cli/packages/cli/src/services/SlashCommandConflictHandler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "pickle",
		Description: "Auto-generated stub for pickle",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "pickle"), nil
		},
	})

	// launch extracted from gemini-cli/packages/cli/src/services/SlashCommandConflictHandler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "launch",
		Description: "Auto-generated stub for launch",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "launch"), nil
		},
	})

	// b extracted from gemini-cli/packages/cli/src/services/SlashCommandConflictHandler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "b",
		Description: "Auto-generated stub for b",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "b"), nil
		},
	})

	// cmd extracted from gemini-cli/packages/cli/src/services/prompt-processors/shellProcessor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "cmd",
		Description: "Auto-generated stub for cmd",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "cmd"), nil
		},
	})

	// spaces extracted from gemini-cli/packages/cli/src/services/prompt-processors/shellProcessor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "spaces",
		Description: "Auto-generated stub for spaces",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "spaces"), nil
		},
	})

	// mycommand extracted from gemini-cli/packages/cli/src/services/prompt-processors/argumentProcessor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mycommand",
		Description: "Auto-generated stub for mycommand",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mycommand"), nil
		},
	})

	// good-agents-ext extracted from gemini-cli/packages/cli/src/config/extension-manager-agents.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "good-agents-ext",
		Description: "Auto-generated stub for good-agents-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "good-agents-ext", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "good-agents-ext", result), nil
		},
	})

	// bad-agents-ext extracted from gemini-cli/packages/cli/src/config/extension-manager-agents.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "bad-agents-ext",
		Description: "Auto-generated stub for bad-agents-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "bad-agents-ext", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "bad-agents-ext", result), nil
		},
	})

	// terafox extracted from gemini-cli/packages/cli/src/config/settings_validation_warning.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "terafox",
		Description: "Auto-generated stub for terafox",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "terafox"), nil
		},
	})

	// unknown_tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "unknown_tool",
		Description: "Auto-generated stub for unknown_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "unknown_tool"), nil
		},
	})

	// mcp_allowed-server_tool1 extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_allowed-server_tool1",
		Description: "Auto-generated stub for mcp_allowed-server_tool1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_allowed-server_tool1"), nil
		},
	})

	// mcp_allowed-server_another_tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_allowed-server_another_tool",
		Description: "Auto-generated stub for mcp_allowed-server_another_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_allowed-server_another_tool"), nil
		},
	})

	// mcp_trusted-server_tool1 extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_trusted-server_tool1",
		Description: "Auto-generated stub for mcp_trusted-server_tool1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_trusted-server_tool1"), nil
		},
	})

	// mcp_trusted-server_special_tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_trusted-server_special_tool",
		Description: "Auto-generated stub for mcp_trusted-server_special_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_trusted-server_special_tool"), nil
		},
	})

	// mcp_blocked-server_any_tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_blocked-server_any_tool",
		Description: "Auto-generated stub for mcp_blocked-server_any_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_blocked-server_any_tool"), nil
		},
	})

	// mcp_unknown-server_tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_unknown-server_tool",
		Description: "Auto-generated stub for mcp_unknown-server_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_unknown-server_tool"), nil
		},
	})

	// mcp_mcp-server_tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_mcp-server_tool",
		Description: "Auto-generated stub for mcp_mcp-server_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_mcp-server_tool"), nil
		},
	})

	// mcp_another-server_tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_another-server_tool",
		Description: "Auto-generated stub for mcp_another-server_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_another-server_tool"), nil
		},
	})

	// custom-tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "custom-tool",
		Description: "Auto-generated stub for custom-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "custom-tool"), nil
		},
	})

	// mcp_my-server_special-tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_my-server_special-tool",
		Description: "Auto-generated stub for mcp_my-server_special-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_my-server_special-tool"), nil
		},
	})

	// mcp_allowed-server_tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_allowed-server_tool",
		Description: "Auto-generated stub for mcp_allowed-server_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_allowed-server_tool"), nil
		},
	})

	// mcp_trusted-server_tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_trusted-server_tool",
		Description: "Auto-generated stub for mcp_trusted-server_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_trusted-server_tool"), nil
		},
	})

	// mcp_blocked-server_tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_blocked-server_tool",
		Description: "Auto-generated stub for mcp_blocked-server_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_blocked-server_tool"), nil
		},
	})

	// unknown_agent extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "unknown_agent",
		Description: "Auto-generated stub for unknown_agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "unknown_agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "unknown_agent", result), nil
		},
	})

	// blocked-tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "blocked-tool",
		Description: "Auto-generated stub for blocked-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "blocked-tool"), nil
		},
	})

	// mcp_blocked-server_any extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_blocked-server_any",
		Description: "Auto-generated stub for mcp_blocked-server_any",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_blocked-server_any"), nil
		},
	})

	// specific-tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "specific-tool",
		Description: "Auto-generated stub for specific-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "specific-tool"), nil
		},
	})

	// mcp_trusted-server_any extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_trusted-server_any",
		Description: "Auto-generated stub for mcp_trusted-server_any",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_trusted-server_any"), nil
		},
	})

	// mcp_mcp-server_any extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_mcp-server_any",
		Description: "Auto-generated stub for mcp_mcp-server_any",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_mcp-server_any"), nil
		},
	})

	// mcp_conflicted-server_tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_conflicted-server_tool",
		Description: "Auto-generated stub for mcp_conflicted-server_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_conflicted-server_tool"), nil
		},
	})

	// mcp_my-server_other-tool extracted from gemini-cli/packages/cli/src/config/policy-engine.integration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mcp_my-server_other-tool",
		Description: "Auto-generated stub for mcp_my-server_other-tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "mcp_my-server_other-tool"), nil
		},
	})

	// https_proxy extracted from gemini-cli/packages/cli/src/config/config.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "https_proxy",
		Description: "Auto-generated stub for https_proxy",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "https_proxy"), nil
		},
	})

	// http_proxy extracted from gemini-cli/packages/cli/src/config/config.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "http_proxy",
		Description: "Auto-generated stub for http_proxy",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "http_proxy"), nil
		},
	})

	// HTTPS_PROXY extracted from gemini-cli/packages/cli/src/config/config.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "HTTPS_PROXY",
		Description: "Auto-generated stub for HTTPS_PROXY",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "HTTPS_PROXY"), nil
		},
	})

	// HTTP_PROXY extracted from gemini-cli/packages/cli/src/config/config.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "HTTP_PROXY",
		Description: "Auto-generated stub for HTTP_PROXY",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "HTTP_PROXY"), nil
		},
	})

	// ext3 extracted from gemini-cli/packages/cli/src/config/config.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ext3",
		Description: "Auto-generated stub for ext3",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ext3"), nil
		},
	})

	// ext-plan extracted from gemini-cli/packages/cli/src/config/config.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ext-plan",
		Description: "Auto-generated stub for ext-plan",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ext-plan"), nil
		},
	})

	// Value extracted from gemini-cli/packages/cli/src/config/extension-manager-hydration.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Value",
		Description: "Val",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Value"), nil
		},
	})

	// skills-ext extracted from gemini-cli/packages/cli/src/config/extension-manager-skills.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "skills-ext",
		Description: "Auto-generated stub for skills-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "skills-ext"), nil
		},
	})

	// skills-ext-load extracted from gemini-cli/packages/cli/src/config/extension-manager-skills.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "skills-ext-load",
		Description: "Auto-generated stub for skills-ext-load",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "skills-ext-load"), nil
		},
	})

	// good-skills-ext extracted from gemini-cli/packages/cli/src/config/extension-manager-skills.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "good-skills-ext",
		Description: "Auto-generated stub for good-skills-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "good-skills-ext"), nil
		},
	})

	// traversal-extension extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "traversal-extension",
		Description: "Auto-generated stub for traversal-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "traversal-extension"), nil
		},
	})

	// disabled-extension extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "disabled-extension",
		Description: "Auto-generated stub for disabled-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "disabled-extension"), nil
		},
	})

	// enabled-extension extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "enabled-extension",
		Description: "Auto-generated stub for enabled-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "enabled-extension"), nil
		},
	})

	// policy-extension extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "policy-extension",
		Description: "Auto-generated stub for policy-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "policy-extension"), nil
		},
	})

	// security-test-extension extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "security-test-extension",
		Description: "Auto-generated stub for security-test-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "security-test-extension"), nil
		},
	})

	// my-linked-extension-with-path extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-linked-extension-with-path",
		Description: "Auto-generated stub for my-linked-extension-with-path",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-linked-extension-with-path"), nil
		},
	})

	// good-ext extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "good-ext",
		Description: "Auto-generated stub for good-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "good-ext"), nil
		},
	})

	// bad_name extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "bad_name",
		Description: "Auto-generated stub for bad_name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "bad_name"), nil
		},
	})

	// my-ext extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-ext",
		Description: "Auto-generated stub for my-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-ext"), nil
		},
	})

	// link-ext-name extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "link-ext-name",
		Description: "Auto-generated stub for link-ext-name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "link-ext-name"), nil
		},
	})

	// no-meta-name extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "no-meta-name",
		Description: "Auto-generated stub for no-meta-name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "no-meta-name"), nil
		},
	})

	// hook-extension extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "hook-extension",
		Description: "Auto-generated stub for hook-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "hook-extension"), nil
		},
	})

	// hook-extension-disabled extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "hook-extension-disabled",
		Description: "Auto-generated stub for hook-extension-disabled",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "hook-extension-disabled"), nil
		},
	})

	// hook-extension-install extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "hook-extension-install",
		Description: "Auto-generated stub for hook-extension-install",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "hook-extension-install"), nil
		},
	})

	// my-local-extension extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-local-extension",
		Description: "Auto-generated stub for my-local-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-local-extension"), nil
		},
	})

	// missing-name-ext extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "missing-name-ext",
		Description: "Auto-generated stub for missing-name-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "missing-name-ext"), nil
		},
	})

	// origin extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "origin",
		Description: "Auto-generated stub for origin",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "origin"), nil
		},
	})

	// my-auto-update-ext extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-auto-update-ext",
		Description: "An old setting",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-auto-update-ext"), nil
		},
	})

	// OLD_SETTING extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "OLD_SETTING",
		Description: "An old setting",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "OLD_SETTING"), nil
		},
	})

	// NEW_SETTING extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "NEW_SETTING",
		Description: "A new setting",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "NEW_SETTING"), nil
		},
	})

	// other-extension extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "other-extension",
		Description: "Auto-generated stub for other-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "other-extension"), nil
		},
	})

	// My-Local-Extension extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "My-Local-Extension",
		Description: "Auto-generated stub for My-Local-Extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "My-Local-Extension"), nil
		},
	})

	// gemini-sql-extension extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "gemini-sql-extension",
		Description: "Auto-generated stub for gemini-sql-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "gemini-sql-extension"), nil
		},
	})

	// no-metadata-extension extracted from gemini-cli/packages/cli/src/config/extension.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "no-metadata-extension",
		Description: "Auto-generated stub for no-metadata-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "no-metadata-extension"), nil
		},
	})

	// my-theme-extension extracted from gemini-cli/packages/cli/src/config/extension-manager-themes.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-theme-extension",
		Description: "Auto-generated stub for my-theme-extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-theme-extension"), nil
		},
	})

	// My-Awesome-Theme extracted from gemini-cli/packages/cli/src/config/extension-manager-themes.spec.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "My-Awesome-Theme",
		Description: "Auto-generated stub for My-Awesome-Theme",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "My-Awesome-Theme"), nil
		},
	})

	// MyTheme extracted from gemini-cli/packages/cli/src/config/extension-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "MyTheme",
		Description: "Auto-generated stub for MyTheme",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "MyTheme"), nil
		},
	})

	// duplicate-ext extracted from gemini-cli/packages/cli/src/config/extension-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "duplicate-ext",
		Description: "Auto-generated stub for duplicate-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "duplicate-ext"), nil
		},
	})

	// integrity-ext extracted from gemini-cli/packages/cli/src/config/extension-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "integrity-ext",
		Description: "Auto-generated stub for integrity-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "integrity-ext"), nil
		},
	})

	// themed-ext extracted from gemini-cli/packages/cli/src/config/extension-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "themed-ext",
		Description: "Auto-generated stub for themed-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "themed-ext"), nil
		},
	})

	// disabled-ext extracted from gemini-cli/packages/cli/src/config/extension-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "disabled-ext",
		Description: "Auto-generated stub for disabled-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "disabled-ext"), nil
		},
	})

	// local-ext extracted from gemini-cli/packages/cli/src/config/extensions/github.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "local-ext",
		Description: "Auto-generated stub for local-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "local-ext"), nil
		},
	})

	// s1 extracted from gemini-cli/packages/cli/src/config/extensions/extensionSettings.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "s1",
		Description: "d1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "s1"), nil
		},
	})

	// s2 extracted from gemini-cli/packages/cli/src/config/extensions/extensionSettings.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "s2",
		Description: "d1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "s2"), nil
		},
	})

	// Username extracted from gemini-cli/packages/cli/src/config/extensions/extensionSettings.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Username",
		Description: "Your public username",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Username"), nil
		},
	})

	// value extracted from gemini-cli/packages/cli/src/config/extensions/extensionSettings.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "value",
		Description: "Test desc",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "value"), nil
		},
	})

	// s3 extracted from gemini-cli/packages/cli/src/config/extensions/extensionSettings.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "s3",
		Description: "d2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "s3"), nil
		},
	})

	// old-ext extracted from gemini-cli/packages/cli/src/config/extensions/consent.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "old-ext",
		Description: "Auto-generated stub for old-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "old-ext"), nil
		},
	})

	// locked-skill extracted from gemini-cli/packages/cli/src/config/extensions/consent.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "locked-skill",
		Description: "A skill in a locked dir",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "locked-skill"), nil
		},
	})

	// ext extracted from gemini-cli/packages/cli/src/config/extensions/update.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ext",
		Description: "Auto-generated stub for ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ext"), nil
		},
	})

	// some-dir extracted from gemini-cli/packages/cli/src/utils/logCleanup.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "some-dir",
		Description: "Auto-generated stub for some-dir",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "some-dir"), nil
		},
	})

	// my-feature extracted from gemini-cli/packages/cli/src/utils/worktreeSetup.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-feature",
		Description: "Auto-generated stub for my-feature",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-feature"), nil
		},
	})

	// generated-name extracted from gemini-cli/packages/cli/src/utils/worktreeSetup.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "generated-name",
		Description: "Auto-generated stub for generated-name",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "generated-name"), nil
		},
	})

	// update_topic extracted from gemini-cli/packages/cli/src/utils/sessionUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "update_topic",
		Description: "Updating the topic",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "update_topic"), nil
		},
	})

	// test_tool_2 extracted from gemini-cli/packages/cli/src/utils/sessionUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test_tool_2",
		Description: "Auto-generated stub for test_tool_2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test_tool_2"), nil
		},
	})

	// static extracted from gemini-cli/packages/cli/src/utils/envVarResolver.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "static",
		Description: "Auto-generated stub for static",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "static"), nil
		},
	})

	// gemini-sandbox extracted from gemini-cli/packages/cli/src/utils/sandbox.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "gemini-sandbox",
		Description: "Auto-generated stub for gemini-sandbox",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "gemini-sandbox"), nil
		},
	})

	// parent extracted from gemini-cli/packages/cli/src/utils/commands.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "parent",
		Description: "Parent command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "parent"), nil
		},
	})

	// notakes extracted from gemini-cli/packages/cli/src/utils/commands.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "notakes",
		Description: "Subcommand that does not take arguments",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "notakes"), nil
		},
	})

	// takes extracted from gemini-cli/packages/cli/src/utils/commands.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "takes",
		Description: "Subcommand that takes arguments",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "takes"), nil
		},
	})

	// CLI_TITLE extracted from gemini-cli/packages/cli/src/utils/windowTitle.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "CLI_TITLE",
		Description: "Auto-generated stub for CLI_TITLE",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "CLI_TITLE"), nil
		},
	})

	// file-reader extracted from gemini-cli/packages/cli/src/ui/commands/toolsCommand.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "file-reader",
		Description: "Reads files from the local system.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "file-reader"), nil
		},
	})

	// code-editor extracted from gemini-cli/packages/cli/src/ui/commands/toolsCommand.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "code-editor",
		Description: "Reads files from the local system.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "code-editor"), nil
		},
	})

	// policies extracted from gemini-cli/packages/cli/src/ui/commands/policiesCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "policies",
		Description: "Manage policies",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "policies"), nil
		},
	})

	// show extracted from gemini-cli/packages/cli/src/ui/commands/memoryCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "show",
		Description: "Show the current memory contents",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "show"), nil
		},
	})

	// reload extracted from gemini-cli/packages/cli/src/ui/commands/memoryCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "reload",
		Description: "Reload the memory from the source",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "reload"), nil
		},
	})

	// inbox extracted from gemini-cli/packages/cli/src/ui/commands/memoryCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "inbox",
		Description: "Review skills extracted from past sessions and move them to global or project skills",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "inbox"), nil
		},
	})

	// skill3 extracted from gemini-cli/packages/cli/src/ui/commands/skillsCommand.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "skill3",
		Description: "Auto-generated stub for skill3",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "skill3"), nil
		},
	})

	// commands extracted from gemini-cli/packages/cli/src/ui/commands/commandsCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "commands",
		Description: "Manage custom slash commands. Usage: /commands [reload]",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "commands"), nil
		},
	})

	// bug extracted from gemini-cli/packages/cli/src/ui/commands/bugCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "bug",
		Description: "Submit a bug report",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "bug"), nil
		},
	})

	// enable extracted from gemini-cli/packages/cli/src/ui/commands/agentsCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "enable",
		Description: "Enable a disabled agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "enable"), nil
		},
	})

	// disable extracted from gemini-cli/packages/cli/src/ui/commands/agentsCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "disable",
		Description: "Disable an enabled agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "disable"), nil
		},
	})

	// agent1 extracted from gemini-cli/packages/cli/src/ui/commands/agentsCommand.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "agent1",
		Description: "desc1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "agent1", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "agent1", result), nil
		},
	})

	// agent2 extracted from gemini-cli/packages/cli/src/ui/commands/agentsCommand.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "agent2",
		Description: "desc1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "agent2", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "agent2", result), nil
		},
	})

	// set extracted from gemini-cli/packages/cli/src/ui/commands/modelCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "set",
		Description: "Set the model to use. Usage: /model set <model-name> [--persist]",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "set"), nil
		},
	})

	// manage extracted from gemini-cli/packages/cli/src/ui/commands/modelCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "manage",
		Description: "Opens a dialog to configure the model",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "manage"), nil
		},
	})

	// docs extracted from gemini-cli/packages/cli/src/ui/commands/docsCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "docs",
		Description: "Open full Gemini CLI documentation in your browser",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "docs"), nil
		},
	})

	// quit extracted from gemini-cli/packages/cli/src/ui/commands/quitCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "quit",
		Description: "Exit the cli",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "quit"), nil
		},
	})

	// update extracted from gemini-cli/packages/cli/src/ui/commands/extensionsCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "update",
		Description: "Update extensions. Usage: update <extension-names>|--all",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "update"), nil
		},
	})

	// install extracted from gemini-cli/packages/cli/src/ui/commands/extensionsCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "install",
		Description: "Install an extension from a git repo or local path",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "install"), nil
		},
	})

	// link extracted from gemini-cli/packages/cli/src/ui/commands/extensionsCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "link",
		Description: "Link an extension from a local path",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "link"), nil
		},
	})

	// uninstall extracted from gemini-cli/packages/cli/src/ui/commands/extensionsCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "uninstall",
		Description: "Uninstall an extension",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "uninstall"), nil
		},
	})

	// desc extracted from gemini-cli/packages/cli/src/ui/commands/toolsCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "desc",
		Description: "List available Gemini CLI tools with descriptions.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "desc"), nil
		},
	})

	// tools extracted from gemini-cli/packages/cli/src/ui/commands/toolsCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tools",
		Description: "List available Gemini CLI tools. Use /tools desc to include descriptions.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tools"), nil
		},
	})

	// shortcuts extracted from gemini-cli/packages/cli/src/ui/commands/shortcutsCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "shortcuts",
		Description: "Toggle the shortcuts panel above the input",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "shortcuts"), nil
		},
	})

	// schema extracted from gemini-cli/packages/cli/src/ui/commands/mcpCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "schema",
		Description: "List configured MCP servers and tools with descriptions and schemas",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "schema"), nil
		},
	})

	// all-ext extracted from gemini-cli/packages/cli/src/ui/commands/extensionsCommand.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "all-ext",
		Description: "Auto-generated stub for all-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "all-ext"), nil
		},
	})

	// setting1 extracted from gemini-cli/packages/cli/src/ui/commands/extensionsCommand.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "setting1",
		Description: "Auto-generated stub for setting1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "setting1"), nil
		},
	})

	// another-inactive-ext extracted from gemini-cli/packages/cli/src/ui/commands/extensionsCommand.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "another-inactive-ext",
		Description: "Auto-generated stub for another-inactive-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "another-inactive-ext"), nil
		},
	})

	// another-active-ext extracted from gemini-cli/packages/cli/src/ui/commands/extensionsCommand.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "another-active-ext",
		Description: "Auto-generated stub for another-active-ext",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "another-active-ext"), nil
		},
	})

	// trust extracted from gemini-cli/packages/cli/src/ui/commands/permissionsCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "trust",
		Description: "Manage folder trust settings. Usage: /permissions trust [<directory-path>]",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "trust"), nil
		},
	})

	// editor extracted from gemini-cli/packages/cli/src/ui/commands/editorCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "editor",
		Description: "Set external editor preference",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "editor"), nil
		},
	})

	// compress extracted from gemini-cli/packages/cli/src/ui/commands/compressCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "compress",
		Description: "Compresses the context by replacing it with a summary",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "compress"), nil
		},
	})

	// setup-github extracted from gemini-cli/packages/cli/src/ui/commands/setupGithubCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "setup-github",
		Description: "Set up GitHub Actions",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "setup-github"), nil
		},
	})

	// test1 extracted from gemini-cli/packages/cli/src/ui/commands/chatCommand.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test1",
		Description: "Auto-generated stub for test1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test1"), nil
		},
	})

	// test2 extracted from gemini-cli/packages/cli/src/ui/commands/chatCommand.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test2",
		Description: "Auto-generated stub for test2",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test2"), nil
		},
	})

	// my-function extracted from gemini-cli/packages/cli/src/ui/commands/chatCommand.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-function",
		Description: "Auto-generated stub for my-function",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-function"), nil
		},
	})

	// signin extracted from gemini-cli/packages/cli/src/ui/commands/authCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "signin",
		Description: "Sign in or change the authentication method",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "signin"), nil
		},
	})

	// signout extracted from gemini-cli/packages/cli/src/ui/commands/authCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "signout",
		Description: "Sign out and clear all cached credentials",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "signout"), nil
		},
	})

	// privacy extracted from gemini-cli/packages/cli/src/ui/commands/privacyCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "privacy",
		Description: "Display the privacy notice",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "privacy"), nil
		},
	})

	// corgi extracted from gemini-cli/packages/cli/src/ui/commands/corgiCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "corgi",
		Description: "Toggles corgi mode",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "corgi"), nil
		},
	})

	// panel extracted from gemini-cli/packages/cli/src/ui/commands/hooksCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "panel",
		Description: "Display all registered hooks with their status",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "panel"), nil
		},
	})

	// enable-all extracted from gemini-cli/packages/cli/src/ui/commands/hooksCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "enable-all",
		Description: "Enable all disabled hooks",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "enable-all"), nil
		},
	})

	// disable-all extracted from gemini-cli/packages/cli/src/ui/commands/hooksCommand.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "disable-all",
		Description: "Disable all enabled hooks",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "disable-all"), nil
		},
	})

	// up extracted from gemini-cli/packages/cli/src/ui/components/shared/text-buffer.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "up",
		Description: "Auto-generated stub for up",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "up"), nil
		},
	})

	// down extracted from gemini-cli/packages/cli/src/ui/components/shared/text-buffer.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "down",
		Description: "Auto-generated stub for down",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "down"), nil
		},
	})

	// i extracted from gemini-cli/packages/cli/src/ui/components/shared/text-buffer.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "i",
		Description: "Auto-generated stub for i",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "i"), nil
		},
	})

	// enter extracted from gemini-cli/packages/cli/src/ui/components/shared/text-buffer.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "enter",
		Description: "Auto-generated stub for enter",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "enter"), nil
		},
	})

	// tab extracted from gemini-cli/packages/cli/src/ui/components/shared/text-buffer.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tab",
		Description: "Auto-generated stub for tab",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tab"), nil
		},
	})

	// backspace extracted from gemini-cli/packages/cli/src/ui/components/shared/text-buffer.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "backspace",
		Description: "Auto-generated stub for backspace",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "backspace"), nil
		},
	})

	// left extracted from gemini-cli/packages/cli/src/ui/components/shared/text-buffer.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "left",
		Description: "Auto-generated stub for left",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "left"), nil
		},
	})

	// right extracted from gemini-cli/packages/cli/src/ui/components/shared/text-buffer.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "right",
		Description: "Auto-generated stub for right",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "right"), nil
		},
	})

	// f1 extracted from gemini-cli/packages/cli/src/ui/components/shared/text-buffer.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "f1",
		Description: "Auto-generated stub for f1",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "f1"), nil
		},
	})

	// left-press extracted from gemini-cli/packages/cli/src/ui/hooks/useMouseClick.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "left-press",
		Description: "Auto-generated stub for left-press",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "left-press"), nil
		},
	})

	// CodebaseInvestigator extracted from gemini-cli/packages/cli/src/ui/hooks/atCommandProcessor_agents.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "CodebaseInvestigator",
		Description: "Investigates codebase",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "CodebaseInvestigator"), nil
		},
	})

	// Expansion extracted from gemini-cli/packages/cli/src/ui/hooks/slashCommandProcessor.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Expansion",
		Description: "Command expansion needs shell access",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Expansion"), nil
		},
	})

	// get_time extracted from gemini-cli/packages/cli/src/ui/hooks/useSessionBrowser.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "get_time",
		Description: "Auto-generated stub for get_time",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "get_time"), nil
		},
	})

	// OtherAgent extracted from gemini-cli/packages/cli/src/ui/hooks/useAtCompletion_agents.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "OtherAgent",
		Description: "Another agent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "OtherAgent"), nil
		},
	})

	// comma extracted from gemini-cli/packages/cli/src/ui/hooks/atCommandProcessor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "comma",
		Description: "Auto-generated stub for comma",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "comma"), nil
		},
	})

	// period extracted from gemini-cli/packages/cli/src/ui/hooks/atCommandProcessor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "period",
		Description: "Auto-generated stub for period",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "period"), nil
		},
	})

	// semicolon extracted from gemini-cli/packages/cli/src/ui/hooks/atCommandProcessor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "semicolon",
		Description: "Auto-generated stub for semicolon",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "semicolon"), nil
		},
	})

	// logs extracted from gemini-cli/packages/cli/src/ui/hooks/atCommandProcessor.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "logs",
		Description: "Auto-generated stub for logs",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "logs"), nil
		},
	})

	// tail_tool extracted from gemini-cli/packages/cli/src/ui/hooks/useToolScheduler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "tail_tool",
		Description: "Auto-generated stub for tail_tool",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "tail_tool"), nil
		},
	})

	// research extracted from gemini-cli/packages/cli/src/ui/hooks/useToolScheduler.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "research",
		Description: "Auto-generated stub for research",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "research"), nil
		},
	})

	// review-frontend extracted from gemini-cli/packages/cli/src/ui/hooks/useSlashCompletion.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "review-frontend",
		Description: "Review frontend code",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "review-frontend"), nil
		},
	})

	// question-mark extracted from gemini-cli/packages/cli/src/ui/hooks/useSlashCompletion.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "question-mark",
		Description: "Alternative name for help",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "question-mark"), nil
		},
	})

	// visible extracted from gemini-cli/packages/cli/src/ui/hooks/useSlashCompletion.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "visible",
		Description: "A visible command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "visible"), nil
		},
	})

	// hidden extracted from gemini-cli/packages/cli/src/ui/hooks/useSlashCompletion.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "hidden",
		Description: "A visible command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "hidden"), nil
		},
	})

	// footer extracted from gemini-cli/packages/cli/src/ui/hooks/useSlashCompletion.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "footer",
		Description: "Configure footer",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "footer"), nil
		},
	})

	// memory-leak extracted from gemini-cli/packages/cli/src/ui/hooks/useSlashCompletion.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "memory-leak",
		Description: "Show memory",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "memory-leak"), nil
		},
	})

	// summarize extracted from gemini-cli/packages/cli/src/ui/hooks/useSlashCompletion.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "summarize",
		Description: "Summarize content",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "summarize"), nil
		},
	})

	// custom-script extracted from gemini-cli/packages/cli/src/ui/hooks/useSlashCompletion.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "custom-script",
		Description: "Run custom script",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "custom-script"), nil
		},
	})

	// left-release extracted from gemini-cli/packages/cli/src/ui/utils/mouse.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "left-release",
		Description: "Auto-generated stub for left-release",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "left-release"), nil
		},
	})

	// sub extracted from gemini-cli/packages/cli/src/ui/utils/directoryUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "sub",
		Description: "Auto-generated stub for sub",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "sub"), nil
		},
	})

	// Downloads extracted from gemini-cli/packages/cli/src/ui/utils/directoryUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Downloads",
		Description: "Auto-generated stub for Downloads",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Downloads"), nil
		},
	})

	// other-project extracted from gemini-cli/packages/cli/src/ui/utils/directoryUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "other-project",
		Description: "Auto-generated stub for other-project",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "other-project"), nil
		},
	})

	// Documents extracted from gemini-cli/packages/cli/src/ui/utils/directoryUtils.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Documents",
		Description: "Auto-generated stub for Documents",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Documents"), nil
		},
	})

	// test-package extracted from gemini-cli/packages/cli/src/ui/utils/updateCheck.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-package",
		Description: "Auto-generated stub for test-package",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-package"), nil
		},
	})

	// MyCustomTheme extracted from gemini-cli/packages/cli/src/ui/themes/theme-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "MyCustomTheme",
		Description: "Auto-generated stub for MyCustomTheme",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "MyCustomTheme"), nil
		},
	})

	// NonExistent extracted from gemini-cli/packages/cli/src/ui/themes/theme-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "NonExistent",
		Description: "Auto-generated stub for NonExistent",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "NonExistent"), nil
		},
	})

	// ExtensionTheme extracted from gemini-cli/packages/cli/src/ui/themes/theme-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "ExtensionTheme",
		Description: "Auto-generated stub for ExtensionTheme",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "ExtensionTheme"), nil
		},
	})

	// Theme extracted from gemini-cli/packages/cli/src/ui/themes/theme-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Theme",
		Description: "Auto-generated stub for Theme",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Theme"), nil
		},
	})

	// SettingsTheme extracted from gemini-cli/packages/cli/src/ui/themes/theme-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "SettingsTheme",
		Description: "Auto-generated stub for SettingsTheme",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "SettingsTheme"), nil
		},
	})

	// MyDark extracted from gemini-cli/packages/cli/src/ui/themes/theme-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "MyDark",
		Description: "Auto-generated stub for MyDark",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "MyDark"), nil
		},
	})

	// MyLight extracted from gemini-cli/packages/cli/src/ui/themes/theme-manager.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "MyLight",
		Description: "Auto-generated stub for MyLight",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "MyLight"), nil
		},
	})

	// None extracted from gemini-cli/packages/cli/src/ui/editors/editorSettingsManager.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "None",
		Description: "Auto-generated stub for None",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "None"), nil
		},
	})

	// YOLO extracted from gemini-cli/packages/cli/src/acp/acpResume.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "YOLO",
		Description: "Auto-approves all tools",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "YOLO"), nil
		},
	})

	// Plan extracted from gemini-cli/packages/cli/src/acp/acpResume.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Plan",
		Description: "Read-only mode",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Plan"), nil
		},
	})

	// gemini-cli extracted from gemini-cli/packages/cli/src/acp/acpClient.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "gemini-cli",
		Description: "Auto-generated stub for gemini-cli",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "gemini-cli"), nil
		},
	})

	// Deny extracted from gemini-cli/packages/cli/src/acp/acpClient.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Deny",
		Description: "Auto-generated stub for Deny",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Deny"), nil
		},
	})

	// Allow extracted from gemini-cli/packages/cli/src/acp/acpClient.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "Allow",
		Description: "Auto-generated stub for Allow",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "Allow"), nil
		},
	})

	// KEY extracted from gemini-cli/packages/cli/src/acp/acpClient.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "KEY",
		Description: "Auto-generated stub for KEY",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "KEY"), nil
		},
	})

	// bravo extracted from gemini-cli/packages/cli/src/acp/commands/help.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "bravo",
		Description: "Bravo command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "bravo"), nil
		},
	})

	// alpha extracted from gemini-cli/packages/cli/src/acp/commands/help.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "alpha",
		Description: "Alpha command",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "alpha"), nil
		},
	})

	// docs-agent extracted from gemini-cli/packages/test-utils/src/fixtures/agents.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "docs-agent",
		Description: "An agent with expertise in updating documentation.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "docs-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "docs-agent", result), nil
		},
	})

	// testing-agent extracted from gemini-cli/packages/test-utils/src/fixtures/agents.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "testing-agent",
		Description: "An agent with expertise in writing and updating tests.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "testing-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "testing-agent", result), nil
		},
	})

	// database-agent extracted from gemini-cli/packages/test-utils/src/fixtures/agents.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "database-agent",
		Description: "An expert in database schemas, SQL, and creating database migrations.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "database-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "database-agent", result), nil
		},
	})

	// css-agent extracted from gemini-cli/packages/test-utils/src/fixtures/agents.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "css-agent",
		Description: "An expert in CSS, styling, and UI design.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "css-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "css-agent", result), nil
		},
	})

	// i18n-agent extracted from gemini-cli/packages/test-utils/src/fixtures/agents.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "i18n-agent",
		Description: "An expert in internationalization and translations.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "i18n-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "i18n-agent", result), nil
		},
	})

	// security-agent extracted from gemini-cli/packages/test-utils/src/fixtures/agents.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "security-agent",
		Description: "An expert in security audits and vulnerability patches.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "security-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "security-agent", result), nil
		},
	})

	// devops-agent extracted from gemini-cli/packages/test-utils/src/fixtures/agents.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "devops-agent",
		Description: "An expert in CI/CD, Docker, and deployment scripts.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "devops-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "devops-agent", result), nil
		},
	})

	// analytics-agent extracted from gemini-cli/packages/test-utils/src/fixtures/agents.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "analytics-agent",
		Description: "An expert in tracking, analytics, and metrics.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "analytics-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "analytics-agent", result), nil
		},
	})

	// accessibility-agent extracted from gemini-cli/packages/test-utils/src/fixtures/agents.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "accessibility-agent",
		Description: "An expert in web accessibility and ARIA roles.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "accessibility-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "accessibility-agent", result), nil
		},
	})

	// mobile-agent extracted from gemini-cli/packages/test-utils/src/fixtures/agents.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "mobile-agent",
		Description: "An expert in React Native and mobile app development.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "mobile-agent", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "mobile-agent", result), nil
		},
	})

	// gemini-cli-companion-mcp-server extracted from gemini-cli/packages/vscode-ide-companion/src/ide-server.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "gemini-cli-companion-mcp-server",
		Description: "(IDE Tool) Open a diff view to create or modify a file. Returns a notification once the diff has been accepted or rejected.",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "gemini-cli-companion-mcp-server"), nil
		},
	})

	// test-project extracted from gemini-cli/evals/validation_fidelity.eval.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-project",
		Description: "Auto-generated stub for test-project",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-project"), nil
		},
	})

	// test-location-repro extracted from gemini-cli/evals/edit-locations-eval.eval.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-location-repro",
		Description: "Auto-generated stub for test-location-repro",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-location-repro"), nil
		},
	})

	// test-api-failure extracted from gemini-cli/evals/test-helper.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-api-failure",
		Description: "Auto-generated stub for test-api-failure",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-api-failure"), nil
		},
	})

	// test-logic-failure extracted from gemini-cli/evals/test-helper.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-logic-failure",
		Description: "Auto-generated stub for test-logic-failure",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-logic-failure"), nil
		},
	})

	// test-recovery extracted from gemini-cli/evals/test-helper.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-recovery",
		Description: "Auto-generated stub for test-recovery",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-recovery"), nil
		},
	})

	// test-api-503 extracted from gemini-cli/evals/test-helper.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-api-503",
		Description: "Auto-generated stub for test-api-503",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-api-503"), nil
		},
	})

	// test-absolute-path extracted from gemini-cli/evals/test-helper.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-absolute-path",
		Description: "Auto-generated stub for test-absolute-path",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-absolute-path"), nil
		},
	})

	// test-traversal extracted from gemini-cli/evals/test-helper.test.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "test-traversal",
		Description: "Auto-generated stub for test-traversal",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "test-traversal"), nil
		},
	})

	// subagent-eval-project extracted from gemini-cli/evals/subagents.eval.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "subagent-eval-project",
		Description: "Auto-generated stub for subagent-eval-project",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			streamCallback := func(chunk string) {
				fmt.Print(chunk)
			}

			result, err := subagents.GlobalManager.Spawn(context.Background(), subagents.TypePlan, "Autogenerated subagent task", "Execute delegated task", "", streamCallback)
			if err != nil {
				return "", fmt.Errorf("subagent %s failed: %w", "subagent-eval-project", err)
			}
			return fmt.Sprintf("Subagent %s completed:\n%s", "subagent-eval-project", result), nil
		},
	})

	// example extracted from gemini-cli/evals/interactive-hang.eval.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "example",
		Description: "Auto-generated stub for example",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "example"), nil
		},
	})

	// my-app extracted from gemini-cli/evals/ask_user.eval.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "my-app",
		Description: "Auto-generated stub for my-app",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "my-app"), nil
		},
	})

	// users-api extracted from gemini-cli/evals/update_topic.eval.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "users-api",
		Description: "Auto-generated stub for users-api",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "users-api"), nil
		},
	})

	// typescript-project extracted from gemini-cli/evals/automated-tool-use.eval.ts
	r.Tools = append(r.Tools, Tool{
		Name:        "typescript-project",
		Description: "Auto-generated stub for typescript-project",
		Parameters:  json.RawMessage(`{}`),
		Execute: func(args map[string]interface{}) (string, error) {
			return fmt.Sprintf("Tool %s not fully wired yet", "typescript-project"), nil
		},
	})
}
