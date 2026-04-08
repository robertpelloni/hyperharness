package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/robertpelloni/hyperharness/foundation/adapters"
	"github.com/robertpelloni/hyperharness/foundation/compat"
	foundationorchestration "github.com/robertpelloni/hyperharness/foundation/orchestration"
	foundationpi "github.com/robertpelloni/hyperharness/foundation/pi"
	foundationrepomap "github.com/robertpelloni/hyperharness/foundation/repomap"
)

type foundationExecRequest struct {
	Tool    string          `json:"tool"`
	Input   json.RawMessage `json:"input"`
	Session string          `json:"session,omitempty"`
}

type foundationPlanRequest struct {
	Prompt       string `json:"prompt"`
	WorkingDir   string `json:"workingDir,omitempty"`
	IncludeRepo  bool   `json:"includeRepo,omitempty"`
	MaxRepoFiles int    `json:"maxRepoFiles,omitempty"`
	TaskType     string `json:"taskType,omitempty"`
	Cost         string `json:"cost,omitempty"`
	RequireLocal bool   `json:"requireLocal,omitempty"`
}

type foundationSessionCreateRequest struct {
	Name string `json:"name,omitempty"`
}

type foundationSessionForkRequest struct {
	Entry string `json:"entry,omitempty"`
	Name  string `json:"name,omitempty"`
}

type foundationBranchSummaryRequest struct {
	Session   string `json:"session"`
	Target    string `json:"target"`
	MaxTokens int    `json:"maxTokens,omitempty"`
}

type foundationCompactionRequest struct {
	Session          string `json:"session"`
	KeepRecentTokens int    `json:"keepRecentTokens,omitempty"`
}

type foundationMCPCallRequest struct {
	Server    string                 `json:"server"`
	Tool      string                 `json:"tool"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type foundationProviderRouteRequest struct {
	TaskType       string `json:"taskType,omitempty"`
	CostPreference string `json:"costPreference,omitempty"`
	RequireLocal   bool   `json:"requireLocal,omitempty"`
}

type foundationProviderPrepareRequest struct {
	Prompt         string `json:"prompt,omitempty"`
	TaskType       string `json:"taskType,omitempty"`
	CostPreference string `json:"costPreference,omitempty"`
	RequireLocal   bool   `json:"requireLocal,omitempty"`
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

func providerStatusPayload() adapters.ProviderStatus {
	return adapters.BuildProviderStatus()
}

func selectFoundationProviderRoute(body foundationProviderRouteRequest) adapters.ProviderRoute {
	return adapters.SelectProviderRoute(adapters.ProviderRouteRequest{
		TaskType:       body.TaskType,
		CostPreference: body.CostPreference,
		RequireLocal:   body.RequireLocal,
	})
}

func prepareFoundationProviderExecution(body foundationProviderPrepareRequest) adapters.ProviderExecutionResult {
	return adapters.PrepareProviderExecution(adapters.ProviderExecutionRequest{
		Prompt:         body.Prompt,
		TaskType:       body.TaskType,
		CostPreference: body.CostPreference,
		RequireLocal:   body.RequireLocal,
	})
}

func listFoundationMCPTools(cwd string) ([]string, error) {
	adapter := adapters.NewMCPAdapter(cwd)
	return adapter.ListTools()
}

func callFoundationMCPTool(cwd string, body foundationMCPCallRequest) (adapters.MCPCallResult, error) {
	adapter := adapters.NewMCPAdapter(cwd)
	return adapter.CallTool(adapters.MCPCallRequest{
		ServerName: body.Server,
		ToolName:   body.Tool,
		Arguments:  body.Arguments,
	})
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

func generateFoundationPlan(cwd string, body foundationPlanRequest) (foundationorchestration.PlanResult, error) {
	workingDir := body.WorkingDir
	if workingDir == "" {
		workingDir = cwd
	}
	return foundationorchestration.BuildPlan(foundationorchestration.PlanRequest{
		Prompt:       body.Prompt,
		WorkingDir:   workingDir,
		IncludeRepo:  body.IncludeRepo,
		MaxRepoFiles: body.MaxRepoFiles,
		TaskType:     body.TaskType,
		Cost:         body.Cost,
		RequireLocal: body.RequireLocal,
	})
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

func prepareFoundationBranchSummary(cwd string, body foundationBranchSummaryRequest) (*foundationpi.BranchSummaryPreparation, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	if body.MaxTokens > 0 {
		return runtime.PrepareBranchSummaryWithBudget(body.Session, body.Target, body.MaxTokens)
	}
	return runtime.PrepareBranchSummary(body.Session, body.Target)
}

func generateFoundationBranchSummary(cwd string, body foundationBranchSummaryRequest) (map[string]any, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	prep, err := prepareFoundationBranchSummary(cwd, body)
	if err != nil {
		return nil, err
	}
	summary, err := runtime.GenerateBranchSummary(context.Background(), prep, nil)
	if err != nil {
		return nil, err
	}
	session, generated, err := runtime.BranchWithGeneratedSummary(context.Background(), body.Session, body.Target, body.MaxTokens, nil, map[string]any{"readFiles": prep.FileOps.ReadFiles, "modifiedFiles": prep.FileOps.ModifiedFiles})
	if err != nil {
		return nil, err
	}
	return map[string]any{"preparation": prep, "summary": summary, "generated": generated, "session": session}, nil
}

func prepareFoundationCompaction(cwd string, body foundationCompactionRequest) (*foundationpi.CompactionPreparation, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	if body.KeepRecentTokens > 0 {
		return runtime.PrepareCompactionWithBudget(body.Session, body.KeepRecentTokens)
	}
	return runtime.PrepareCompaction(body.Session)
}

func generateFoundationCompaction(cwd string, body foundationCompactionRequest) (map[string]any, error) {
	runtime := foundationpi.NewRuntime(cwd, nil)
	prep, err := prepareFoundationCompaction(cwd, body)
	if err != nil {
		return nil, err
	}
	summary, err := runtime.GenerateCompactionSummary(context.Background(), prep, nil)
	if err != nil {
		return nil, err
	}
	session, generated, err := runtime.CompactWithGeneratedSummary(context.Background(), body.Session, body.KeepRecentTokens, nil, map[string]any{"readFiles": prep.FileOps.ReadFiles, "modifiedFiles": prep.FileOps.ModifiedFiles})
	if err != nil {
		return nil, err
	}
	return map[string]any{"preparation": prep, "summary": summary, "generated": generated, "session": session}, nil
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
