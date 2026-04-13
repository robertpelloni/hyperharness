import sys

with open("TODO.md", "r") as f:
    content = f.read()

content = content.replace("- [ ] SQLite backend with FTS5 for knowledge base", "- [x] SQLite backend with FTS5 for knowledge base")

with open("TODO.md", "w") as f:
    f.write(content)
