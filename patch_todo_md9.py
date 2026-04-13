import sys

with open("TODO.md", "r") as f:
    content = f.read()

content = content.replace("- [ ] Consolidate truncateString, getStr, getInt into shared package", "- [x] Consolidate truncateString, getStr, getInt into shared package")
content = content.replace("- [ ] OpenCode `plan_enter`/`plan_exit` — wire to actual planning mode", "- [x] OpenCode `plan_enter`/`plan_exit` — wire to actual planning mode")
content = content.replace("- [ ] Crush `delegate` — wire to actual subagent delegation", "- [x] Crush `delegate` — wire to actual subagent delegation")

with open("TODO.md", "w") as f:
    f.write(content)
