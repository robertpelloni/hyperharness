import sys

with open("tools/claude_code_parity.go", "r") as f:
    content = f.read()

pattern = """	// Agent - Claude Code's subagent spawning tool (renamed from "Task")
	r.Tools = append(r.Tools, Tool{
		Name:        "Agent",
		Description: "Launch a new agent that has access to the following tools: Read, Write, Edit, Bash, Glob, Grep, WebSearch, WebFetch, TodoWrite, Skill, ToolSearch, LSP. The agent will work on the given task and return a result. Use this for multi-step research, exploration, or implementation tasks.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"prompt":{"type":"string","description":"The task for the agent to perform"},"type":{"type":"string","description":"Agent type: Explore, Plan, or custom","enum":["Explore","Plan","verification","custom"]}},"required":["prompt"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			prompt := getStr(args, "prompt")
			agentType := getStr(args, "type")
			if agentType == "" {
				agentType = "custom"
			}

			// Delegate to the subagent system
			result := fmt.Sprintf("[Agent|%s] Task spawned: %s\\n\\nAgent completed. Use the results below.\\n(Output from Go-native subagent runtime)", agentType, truncateString(prompt, 200))
			return result, nil
		},
	})"""

replacement = """	// Agent - Claude Code's subagent spawning tool (renamed from "Task")
	r.Tools = append(r.Tools, Tool{
		Name:        "Agent",
		Description: "Launch a new agent that has access to the following tools: Read, Write, Edit, Bash, Glob, Grep, WebSearch, WebFetch, TodoWrite, Skill, ToolSearch, LSP. The agent will work on the given task and return a result. Use this for multi-step research, exploration, or implementation tasks.",
		Parameters:  json.RawMessage(`{"type":"object","properties":{"prompt":{"type":"string","description":"The task for the agent to perform"},"type":{"type":"string","description":"Agent type: Explore, Plan, or custom","enum":["Explore","Plan","verification","custom"]}},"required":["prompt"]}`),
		Execute: func(args map[string]interface{}) (string, error) {
			prompt := getStr(args, "prompt")
			agentType := getStr(args, "type")
			if agentType == "" {
				agentType = "custom"
			}

			// Delegate to the subagent system
			// Map agent type to our internal types
			sType := subagents.TypeResearch
			switch strings.ToLower(agentType) {
			case "plan":
				sType = subagents.TypePlan
			case "explore", "research":
				sType = subagents.TypeResearch
			case "verification", "test":
				sType = subagents.TypeTest
			case "code", "custom":
				sType = subagents.TypeCode
			}

			task := subagents.NewManager().CreateTask(sType, prompt, prompt, "")

			// For this implementation, we would ideally execute it asynchronously or use a background context.
			// But for parity with the old stub, we format a message indicating spawning/completion.
			// In the future, this would call ExecuteTask and stream/wait for actual subagent response.

			// Simulate waiting/executing
			result, err := subagents.NewManager().ExecuteTask(context.Background(), task)
			if err != nil {
				return "", fmt.Errorf("failed to execute subagent: %w", err)
			}

			return fmt.Sprintf("[Agent|%s] Task spawned: %s\n\n%s", agentType, truncateString(prompt, 200), result), nil
		},
	})"""

if pattern in content:
    new_content = content.replace(pattern, replacement)

    # We need to ensure context and subagents are imported
    if '"github.com/robertpelloni/hyperharness/internal/subagents"' not in new_content:
        new_content = new_content.replace('"github.com/robertpelloni/hyperharness/foundation/pi"',
                                          '"github.com/robertpelloni/hyperharness/foundation/pi"\n\t"github.com/robertpelloni/hyperharness/internal/subagents"\n\t"context"')

    with open("tools/claude_code_parity.go", "w") as f:
        f.write(new_content)
    print("Replaced Agent successfully")
else:
    print("Pattern not found")
