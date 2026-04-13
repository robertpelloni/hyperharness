import sys

with open("tools/goose_opencode_kimi_parity.go", "r") as f:
    content = f.read()

plan_exit_old = """		// plan_exit - OpenCode's plan mode exit tool
		r.Tools = append(r.Tools, Tool{
			Name:        "plan_exit",
			Description: "Exit plan mode and switch to build agent. Confirms the plan is complete and asks user if they want to start implementing.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {},
				"additionalProperties": false
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				return "Plan exit requested. Switching to build agent.\\nUser approved switching to build agent. Wait for further instructions.", nil
			},
		})"""

plan_exit_new = """		// plan_exit - OpenCode's plan mode exit tool
		r.Tools = append(r.Tools, Tool{
			Name:        "plan_exit",
			Description: "Exit plan mode and switch to build agent. Confirms the plan is complete and asks user if they want to start implementing.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {},
				"additionalProperties": false
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				// Signal to the orchestrator to transition from Plan to Build phase
				return "Plan exit requested. Switching to build agent.\\nUser approved switching to build agent. Awaiting instructions from build agent.", nil
			},
		})"""

plan_enter_old = """		// plan_enter - Enter plan mode
		r.Tools = append(r.Tools, Tool{
			Name:        "plan_enter",
			Description: "Enter plan mode for research and planning. Switches to plan agent that can read but not modify files.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {},
				"additionalProperties": false
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				return "Plan mode entered. Switched to plan agent for research and planning.", nil
			},
		})"""

plan_enter_new = """		// plan_enter - Enter plan mode
		r.Tools = append(r.Tools, Tool{
			Name:        "plan_enter",
			Description: "Enter plan mode for research and planning. Switches to plan agent that can read but not modify files.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {},
				"additionalProperties": false
			}`),
			Execute: func(args map[string]interface{}) (string, error) {
				// Signal to the orchestrator to transition to a Plan agent
				return "Plan mode entered. Switched to plan agent for research and planning. You now have read-only tools.", nil
			},
		})"""

content = content.replace(plan_exit_old.replace("\\n", "\n"), plan_exit_new.replace("\\n", "\n"))
content = content.replace(plan_enter_old.replace("\\n", "\n"), plan_enter_new.replace("\\n", "\n"))

with open("tools/goose_opencode_kimi_parity.go", "w") as f:
    f.write(content)
print("Updated plan_enter/exit")
