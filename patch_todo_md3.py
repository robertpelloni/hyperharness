import sys

with open("TODO.md", "r") as f:
    content = f.read()

content = content.replace("- [ ] Add SSE transport support", "- [x] Add SSE transport support")
content = content.replace("- [ ] Tool discovery from connected MCP servers", "- [x] Tool discovery from connected MCP servers")

with open("TODO.md", "w") as f:
    f.write(content)
print("Updated TODO.md")
