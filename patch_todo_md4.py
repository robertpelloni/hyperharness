import sys

with open("TODO.md", "r") as f:
    content = f.read()

content = content.replace("- [ ] Smithery registry API integration", "- [x] Smithery registry API integration")

with open("TODO.md", "w") as f:
    f.write(content)
print("Updated TODO.md")
