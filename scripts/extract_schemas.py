import os
import json
import re

TARGET_DIRS = [
    "claude-code",
    "hermes-agent/tools",
    "crush",
    "opencode",
    "gemini-cli"
]

extracted_tools = []

def extract_ts_tools(content):
    tools = []
    
    # We strictly want something that looks like an object field: name: "ToolName"
    # To avoid regex vars like `const match = line.match(...)` we ensure the value is a pure string literal 
    name_matches = re.finditer(r'name\s*:\s*(["\'])([a-zA-Z0-9_\-]+)\1', content)
    for match in name_matches:
        tool_name = match.group(2)
        if tool_name in ["message", "agent", "system", "user", "run", "config", "test", "index"]: continue
        
        start_idx = max(0, match.start() - 100)
        end_idx = min(len(content), match.end() + 500)
        block = content[start_idx:end_idx]
        
        # Strict description match on a single line
        desc_match = re.search(r'description\s*:\s*(["\'])(.*?)\1', block)
        desc = desc_match.group(2) if desc_match else f"Auto-generated stub for {tool_name}"
        
        params = "{}"
        
        tools.append({
            "name": tool_name,
            "description": desc,
            "parameters": params
        })
        
    return tools

for root_dir in TARGET_DIRS:
    if not os.path.exists(root_dir):
        continue
        
    for root, _, files in os.walk(root_dir):
        for file in files:
            if file.endswith(".ts") or file.endswith(".py"):
                filepath = os.path.join(root, file)
                try:
                    with open(filepath, "r", encoding="utf-8") as f:
                        content = f.read()
                        
                    if file.endswith(".ts"):
                        ts_tools = extract_ts_tools(content)
                        for t in ts_tools:
                            t["source"] = filepath
                            extracted_tools.append(t)
                except Exception:
                    pass

unique_tools = {}
for t in extracted_tools:
    if t["name"] not in unique_tools:
        unique_tools[t["name"]] = t

print(f"Extracted {len(unique_tools)} potential tools safely.")
with open("scripts/extracted_tools.json", "w") as f:
    json.dump(list(unique_tools.values()), f, indent=2)

