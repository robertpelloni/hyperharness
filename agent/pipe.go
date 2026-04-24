package agent

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/robertpelloni/hyperharness/foundation/adapters"
)

// PipeProcessor mimics Simon Willison's "LLM CLI" and Charmbracelet's "Crush".
// It takes data from standard input, processes it through the LLM with a prompt,
// and streams or returns the result.
func (a *Agent) ProcessPipe(prompt string) (string, error) {
	if a == nil {
		return "", fmt.Errorf("agent is required")
	}
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	return processPipeWithReader(prompt, stat.Mode(), os.Stdin, a.Chat)
}

func processPipeWithReader(prompt string, mode os.FileMode, reader io.Reader, chat func(string) (string, error)) (string, error) {
	if (mode & os.ModeCharDevice) != 0 {
		return "", fmt.Errorf("no data piped to stdin")
	}
	if reader == nil {
		return "", fmt.Errorf("stdin reader is required")
	}
	if chat == nil {
		return "", fmt.Errorf("chat function is required")
	}

	inputData, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read from stdin: %w", err)
	}

	execution := adapters.PrepareProviderExecution(adapters.ProviderExecutionRequest{Prompt: prompt, TaskType: "analysis", CostPreference: "budget"})
	combinedPrompt := buildPipePrompt(prompt, execution.ExecutionHint, string(inputData))
	fmt.Printf("[PipeProcessor] Processing %d bytes of piped data via %s/%s...\n", len(inputData), execution.Route.Provider, execution.Route.Model)

	response, err := chat(combinedPrompt)
	if err != nil {
		return "", err
	}
	return formatPipeResponse(execution.ExecutionHint, response), nil
}

func buildPipePrompt(prompt, executionHint, input string) string {
	return fmt.Sprintf("%s\n\nExecution Hint:\n%s\n\nInput Data:\n%s", prompt, executionHint, input)
}

func formatPipeResponse(executionHint, response string) string {
	if strings.TrimSpace(executionHint) != "" {
		return fmt.Sprintf("[Pipe Execution]\n%s\n\n%s", executionHint, response)
	}
	return response
}
