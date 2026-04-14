import sys

with open("tools/claude_code_parity.go", "r") as f:
    lines = f.readlines()

out = []
for line in lines:
    out.append(line.replace("\\n", "\\n").replace("\\n", "\\n")) # That does nothing. Wait.

with open("tools/claude_code_parity.go", "r") as f:
    content = f.read()

bad_str = r'return fmt.Sprintf("[Agent|%s] Task spawned: %s\n\n%s", agentType, truncateString(prompt, 200), result), nil'
good_str = 'return fmt.Sprintf("[Agent|%s] Task spawned: %s\\n\\n%s", agentType, truncateString(prompt, 200), result), nil'

content = content.replace(bad_str, good_str)

with open("tools/claude_code_parity.go", "w") as f:
    f.write(content)
