import sys

with open("TODO.md", "r") as f:
    content = f.read()

content = content.replace("- [ ] Integration test: full agent loop with mock provider", "- [x] Integration test: full agent loop with mock provider")
content = content.replace("- [ ] Integration test: tool execution through Harness.ExecuteTool()", "- [x] Integration test: tool execution through Harness.ExecuteTool()")
content = content.replace("- [ ] Integration test: MCP server connection lifecycle", "- [x] Integration test: MCP server connection lifecycle")

with open("TODO.md", "w") as f:
    f.write(content)
print("Updated TODO.md")
