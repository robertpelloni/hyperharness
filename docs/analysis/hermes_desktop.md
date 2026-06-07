# Hermes Desktop Analysis

Hermes Desktop (part of the Nous Research Hermes ecosystem) provides a visual and CLI interface for interacting with the Hermes family of models, with a focus on advanced tool-calling and code execution.

## Features

### 1. Advanced Tool Use
- **Structured Outputs**: Precise control over tool parameter generation.
- **Computer Use**: Similar to Codex, often includes browser and system automation.

### 2. Provider Support
- **Ollama/LMStudio**: Strong focus on local provider integration.

## Implementation in Go

- [x] Hermes Tool Parity (`tools/hermes_parity.go`)
- [x] Hermes II Advanced Tools (`tools/hermes_ii_parity.go`)
- [x] Local Provider Routing (`internal/providers/`)
