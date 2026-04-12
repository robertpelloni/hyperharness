import sys

with open("TODO.md", "r") as f:
    content = f.read()

content = content.replace("- [ ] Vector embeddings for semantic search", "- [x] Vector embeddings for semantic search")
content = content.replace("- [ ] Memory decay weighting (older = less relevant)", "- [x] Memory decay weighting (older = less relevant)")

with open("TODO.md", "w") as f:
    f.write(content)
print("Updated TODO.md")
