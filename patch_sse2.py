import sys
import re

with open("internal/mcp/mcp.go", "r") as f:
    content = f.read()

# Make sure imports are there
if '"time"' not in content:
    content = content.replace(
        '"net/http"\n\t"bytes"\n\t"context"',
        '"net/http"\n\t"bytes"\n\t"context"\n\t"time"'
    )

server_pattern = """type MCPServer struct {
	Name      string            `json:"name"`
	Command   string            `json:"command,omitempty"`
	Args      []string          `json:"args,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	URL       string            `json:"url,omitempty"`
	Transport string            `json:"transport"` // "stdio" or "sse"
	Connected bool              `json:"connected"`
	Tools     []Tool            `json:"tools,omitempty"`

	// Internal runtime fields
	stdin  io.WriteCloser `json:"-"`
	stdout *bufio.Reader  `json:"-"`
	cmd    *exec.Cmd      `json:"-"`
	reqID  int            `json:"-"`
	mu     sync.Mutex     `json:"-"`
}"""

server_replacement = """type MCPServer struct {
	Name      string            `json:"name"`
	Command   string            `json:"command,omitempty"`
	Args      []string          `json:"args,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	URL       string            `json:"url,omitempty"`
	Transport string            `json:"transport"` // "stdio" or "sse"
	Connected bool              `json:"connected"`
	Tools     []Tool            `json:"tools,omitempty"`

	// Internal runtime fields
	stdin  io.WriteCloser `json:"-"`
	stdout *bufio.Reader  `json:"-"`
	cmd    *exec.Cmd      `json:"-"`
	reqID  int            `json:"-"`
	mu     sync.Mutex     `json:"-"`

	// SSE specific
	endpoint  string        `json:"-"`
	client    *http.Client  `json:"-"`
}"""

content = content.replace(server_pattern, server_replacement)

with open("internal/mcp/mcp.go", "w") as f:
    f.write(content)
print("Updated mcp.go with SSE fixes")
