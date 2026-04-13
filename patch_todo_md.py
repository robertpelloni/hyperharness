import sys

with open("TODO.md", "r") as f:
    content = f.read()

content = content.replace("- [ ] Wire `TodoWrite` tool to session-level state persistence", "- [x] Wire `TodoWrite` tool to session-level state persistence")
content = content.replace("- [ ] Wire `Agent` tool to actual `internal/subagents` Manager", "- [x] Wire `Agent` tool to actual `internal/subagents` Manager")
content = content.replace("- [ ] Wire `LSP` tool to gopls/other language servers via stdio", "- [x] Wire `LSP` tool to gopls/other language servers via stdio")
content = content.replace("- [ ] Wire `WebSearch` to Exa or Brave Search API", "- [x] Wire `WebSearch` to Exa or Brave Search API")
content = content.replace("- [ ] Wire `WebFetch` to actual HTTP client", "- [x] Wire `WebFetch` to actual HTTP client")
content = content.replace("- [ ] Wire `Config` tool to actual `internal/config` Manager", "- [x] Wire `Config` tool to actual `internal/config` Manager")
content = content.replace("- [ ] Wire `Skill` tool to actual `internal/skills` Manager", "- [x] Wire `Skill` tool to actual `internal/skills` Manager")

with open("TODO.md", "w") as f:
    f.write(content)
print("Updated TODO.md")
