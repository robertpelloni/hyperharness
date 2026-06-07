# Integration Testing Guide

HyperHarness includes a robust suite of integration tests to verify the autonomy loop and tool parity across different harnesses.

## Running Tests

### 1. All Tests
```bash
go test ./...
```

### 2. End-to-End (E2E) Harness Tests
Verifies workflows for Tabby, Warp, Wave, and Data Ingestion:
```bash
go test -v ./agent/harness_e2e_test.go
```

### 3. Autonomous Orchestration Tests
Verifies the Council architecture (Director-Worker) and task decomposition:
```bash
go test -v ./internal/council/council_test.go
go test -v ./agent/autopilot_test.go
```

### 4. Tool Registry Tests
Verifies tool registration and execution mapping:
```bash
go test -v ./tools/registry_test.go
```

## Mocking
Integration tests use the `MockExecutor` or mock LLM providers to ensure stability without consuming real tokens. See `agent/autopilot_mock_test.go` for examples.
