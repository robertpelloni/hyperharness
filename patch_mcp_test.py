import sys
import re

with open("internal/mcp/mcp_integration_test.go", "r") as f:
    content = f.read()

content = content.replace("mcp.NewRegistry", "NewRegistry")
content = content.replace("mcp.MCPServer", "MCPServer")
content = content.replace('"github.com/robertpelloni/hyperharness/internal/mcp"', "")

with open("internal/mcp/mcp_integration_test.go", "w") as f:
    f.write(content)
