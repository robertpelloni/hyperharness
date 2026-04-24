import sys

with open("TODO.md", "r") as f:
    content = f.read()

content = content.replace("- [ ] SQLite backend with FTS5 for knowledge base", "- [ ] SQLite backend with FTS5 for knowledge base") # Ensure no changes

with open("TODO.md", "w") as f:
    f.write(content)
