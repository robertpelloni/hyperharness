import json
import os

with open("scripts/extracted_tools.json", "r") as f:
    tools = json.load(f)

# Generates boilerplate Go structs for parity tools
go_output = """package tools

import (
	"encoding/json"
	"fmt"
)

// registerGeneratedTools dynamically generates stub Tool definitions from extracted upstream JSON schemas.
// These act as empty parity facades, and can be manually wired up to the internal Foundation and MCP services later.
func registerGeneratedTools(r *Registry) {
"""

for t in tools:
    name = t.get("name", "").strip()
    desc = t.get("description", "").replace('"', '\\"').replace("\n", " ")
    params = t.get("parameters", "{}").strip()
    
    # Very basic escaping to ensure syntactically valid Go
    # Since regex extraction was naive, we'll wrap params inside a literal string,
    # but we must be careful with backticks inside params
    safe_params = params.replace('`', '` + "`" + `')
    
    # Just skip if params string is malformed or massive due to bad regex parse
    if len(safe_params) > 5000:
        safe_params = "{}"
        
    go_output += f"""
	// {name} extracted from {t.get("source", "unknown")}
	r.Tools = append(r.Tools, Tool{{
		Name:        "{name}",
		Description: "{desc}",
		Parameters: json.RawMessage(`{safe_params}`),
		Execute: func(args map[string]interface{{}}) (string, error) {{
			return fmt.Sprintf("Tool %s not fully wired yet", "{name}"), nil
		}},
	}})
"""

go_output += "}\n"

with open("tools/generated_tools.go", "w") as f:
    f.write(go_output)

# Now for TypeScript
ts_output = """import type { AgentTool, AgentToolResult, AgentToolUpdateCallback } from "@mariozechner/pi-agent-core";

// Dynamically generated stub Tool definitions from extracted upstream JSON schemas.
// These act as empty parity facades, and can be manually wired up to the internal Pi CLI tools later.

export function createGeneratedParityTools(): AgentTool<any>[] {
    return [
"""

for t in tools:
    name = t.get("name", "").strip()
    desc = t.get("description", "").replace('"', '\\"').replace("\n", " ")
    
    ts_output += f"""        {{
            name: "{name}",
			label: "{name}",
            description: "{desc}",
            parameters: {{ type: "object", properties: {{}} }}, // Fallback empty schema stub
            execute: async (_toolCallId: string, _params: any) => {{
                return {{ isError: false, type: "text", text: "Tool {name} not fully wired yet" }} as any;
            }}
        }},
"""

ts_output += "    ];\n}\n"

os.makedirs("pi-cli/packages/coding-agent/src/core/tools/parity/generated", exist_ok=True)
with open("pi-cli/packages/coding-agent/src/core/tools/parity/generated/index.ts", "w") as f:
    f.write(ts_output)

print(f"Successfully generated {len(tools)} stub tools for Go and TypeScript.")
