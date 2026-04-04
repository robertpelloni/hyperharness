package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/robertpelloni/hypercode/foundation/adapters"
	"github.com/robertpelloni/hypercode/foundation/compat"
	foundationpi "github.com/robertpelloni/hypercode/foundation/pi"
	foundationrepomap "github.com/robertpelloni/hypercode/foundation/repomap"
)

type foundationExecRequest struct {
	Tool    string          `json:"tool"`
	Input   json.RawMessage `json:"input"`
	Session string          `json:"session,omitempty"`
}

type foundationSessionCreateRequest struct {
	Name string `json:"name,omitempty"`
}

type foundationSessionForkRequest struct {
	Entry string `json:"entry,omitempty"`
	Name  string `json:"name,omitempty"`
}

func currentFoundationRuntime() (*foundationpi.Runtime, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, "", err
	}
	return foundationpi.NewRuntime(cwd, nil), cwd, nil
}

func foundationAdaptersPayload(cwd string) map[string]any {
	hyperAdapter := adapters.NewHyperCodeAdapter(cwd)
	mcpAdapter := adapters.NewMCPAdapter(cwd)
	return map[string]any{
		"hypercode": hyperAdapter.Status(),
		"mcp":       mcpAdapter.Status(),
	}
}

func mcpToolContracts() []compat.ToolContract {
	catalog := compat.DefaultCatalog()
	return catalog.ContractsBySource("pi")
}

func executeFoundationTool(cwd string, req foundationExecRequest) (map[string]any, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	result, execErr := runtime.ExecuteTool(context.Background(), req.Session, req.Tool, req.Input, nil)
	payload := map[string]any{
		"tool":   req.Tool,
		"result": result,
	}
	if execErr != nil {
		payload["error"] = execErr.Error()
	}
	return payload, execErr
}

func generateFoundationRepomap(cwd string, body foundationrepomap.Options) (foundationrepomap.Result, error) {
	if body.BaseDir == "" {
		body.BaseDir = cwd
	}
	return foundationrepomap.Generate(body)
}

func createFoundationSession(cwd string, body foundationSessionCreateRequest) (*foundationpi.SessionFile, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	return runtime.CreateSession(body.Name)
}

func listFoundationSessions(cwd string) ([]foundationpi.SessionMetadata, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	return runtime.ListSessions()
}

func getFoundationSession(cwd, sessionID string) (*foundationpi.SessionFile, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	return runtime.LoadSession(sessionID)
}

func forkFoundationSession(cwd, sessionID string, body foundationSessionForkRequest) (*foundationpi.SessionFile, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	return runtime.ForkSession(sessionID, body.Entry, body.Name)
}

func encodeFoundationReadAsString(cwd, requestedPath string) (string, error) {
	input, err := json.Marshal(foundationpi.ReadToolInput{Path: requestedPath})
	if err != nil {
		return "", err
	}
	payload, err := executeFoundationTool(cwd, foundationExecRequest{Tool: "read", Input: input})
	if err != nil {
		return "", err
	}
	result, _ := payload["result"].(*foundationpi.ToolResult)
	if result == nil || len(result.Content) == 0 {
		return "", nil
	}
	if block, ok := result.Content[0].(foundationpi.TextContent); ok {
		return block.Text, nil
	}
	return "", fmt.Errorf("unexpected read result content")
}
