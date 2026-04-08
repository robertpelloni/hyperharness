# HyperHarness Deployment Guide

## Prerequisites

- **Go 1.21+** (tested with Go 1.26.2)
- **Git** (for version info and worktree management)
- **C compiler** (for sqlite3 CGO dependency)

## Building from Source

```bash
# Clone
git clone https://github.com/robertpelloni/hyperharness.git
cd hyperharness

# Build
go build -buildvcs=false -o hyperharness ./cmd/hyperharness

# Install globally
go install -buildvcs=false ./cmd/hyperharness
```

### Build with Version Info

```bash
VERSION=$(cat VERSION)
COMMIT=$(git rev-parse --short HEAD)
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

go build -buildvcs=false -ldflags \
  "-X main.Version=$VERSION \
   -X main.Commit=$COMMIT \
   -X main.BuildDate=$DATE" \
  -o hyperharness ./cmd/hyperharness
```

## Running Tests

```bash
# Full test suite (all 25 packages)
go test -buildvcs=false ./... -count=1 -timeout 180s

# Specific package
go test -buildvcs=false ./tools/ -v -count=1

# With coverage
go test -buildvcs=false ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Configuration

### API Keys (Environment Variables)

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
export OPENAI_API_KEY="sk-..."
export GOOGLE_API_KEY="AIza..."
export GROQ_API_KEY="gsk_..."
export MISTRAL_API_KEY="..."
export XAI_API_KEY="xai-..."
```

### MCP Server Configuration

Create `.hypercode/mcp.json` in your project:

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/dir"]
    }
  }
}
```

### Global Config

Located at `~/.hyperharness/config.json`:

```json
{
  "defaultProvider": "anthropic",
  "defaultModel": "claude-3-5-sonnet-20241022",
  "defaultThinkingLevel": "medium",
  "sessions": {
    "sessionDir": "~/.hyperharness/sessions"
  }
}
```

## Usage Modes

```bash
# Interactive mode (default)
hyperharness

# Print mode (single response)
hyperharness -p "Fix the bug in main.go"

# JSON mode (machine-readable)
hyperharness --mode json "Explain this code"

# RPC mode (process integration)
hyperharness --rpc

# With specific provider/model
hyperharness --provider openai --model gpt-4o "Help me"

# Continue last session
hyperharness -c

# Read-only tools only
hyperharness --tools read,grep,find,ls
```

## Deployment as MCP Server

HyperHarness can act as an MCP server, exposing all 145+ tools:

```json
{
  "mcpServers": {
    "hyperharness": {
      "command": "hyperharness",
      "args": ["--rpc"]
    }
  }
}
```

## Docker Deployment

```dockerfile
FROM golang:1.26-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY . .
RUN go build -buildvcs=false -o hyperharness ./cmd/hyperharness

FROM alpine:latest
RUN apk add --no-cache ca-certificates git
COPY --from=builder /app/hyperharness /usr/local/bin/
ENTRYPOINT ["hyperharness"]
```

```bash
docker build -t hyperharness .
docker run -e ANTHROPIC_API_KEY=sk-... hyperharness -p "Hello"
```
