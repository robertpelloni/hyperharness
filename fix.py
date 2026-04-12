import sys

with open("tools/claude_code_parity.go", "r") as f:
    content = f.read()

bad_string = """			return fmt.Sprintf("[Agent|%s] Task spawned: %s

	%s", agentType, truncateString(prompt, 200), result), nil"""

good_string = """			return fmt.Sprintf("[Agent|%s] Task spawned: %s\\n\\n%s", agentType, truncateString(prompt, 200), result), nil"""

content = content.replace(bad_string, good_string)

with open("tools/claude_code_parity.go", "w") as f:
    f.write(content)
