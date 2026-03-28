package agent

import (
	"fmt"
	"io"
	"os"
)

// PipeProcessor mimics Simon Willison's "LLM CLI" and Charmbracelet's "Crush".
// It takes data from standard input, processes it through the LLM with a prompt,
// and streams or returns the result.
func (a *Agent) ProcessPipe(prompt string) (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	// Check if data is being piped to Stdin
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", fmt.Errorf("no data piped to stdin")
	}

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read from stdin: %w", err)
	}

	combinedPrompt := fmt.Sprintf("%s\n\nInput Data:\n%s", prompt, string(inputData))
	
	fmt.Printf("[PipeProcessor] Processing %d bytes of piped data...\n", len(inputData))
	
	return a.Chat(combinedPrompt)
}
