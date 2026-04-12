import sys

with open("TODO.md", "r") as f:
    content = f.read()

content = content.replace("- [ ] Implement actual stdio transport for MCP server connections", "- [x] Implement actual stdio transport for MCP server connections")
content = content.replace("- [ ] Bidirectional routing (expose internal tools via MCP server)", "- [x] Bidirectional routing (expose internal tools via MCP server)")

with open("TODO.md", "w") as f:
    f.write(content)
print("Updated TODO.md")
