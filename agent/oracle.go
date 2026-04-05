package agent

import (
	"fmt"
	"strings"
)

// Oracle Mode mimics Amp Code's "Oracle" and Deep Research agents.
// It executes a multi-step reasoning loop before returning a final answer.
func (a *Agent) OracleQuery(prompt string) (string, error) {
	return oracleQueryWithChat(strings.TrimSpace(prompt), a.Chat)
}

func oracleQueryWithChat(prompt string, chat func(string) (string, error)) (string, error) {
	fmt.Println("[Oracle Mode] Initiating deep reasoning loop...")
	if strings.TrimSpace(prompt) == "" {
		return "", fmt.Errorf("prompt is required")
	}
	if chat == nil {
		return "", fmt.Errorf("chat function is required")
	}

	// Step 1: Formulate a research plan
	planPrompt := buildOraclePlanPrompt(prompt)
	plan, err := chat(planPrompt)
	if err != nil {
		return "", err
	}

	fmt.Println("[Oracle Mode] Research Plan Generated:")
	fmt.Println(plan)

	// Step 2: Execute the plan using tools
	executionPrompt := buildOracleExecutionPrompt(plan, prompt)

	// We use the standard Chat function which handles tool execution
	researchData, err := chat(executionPrompt)
	if err != nil {
		return "", err
	}

	// Step 3: Synthesize the final answer
	synthesisPrompt := buildOracleSynthesisPrompt(prompt, researchData)

	finalAnswer, err := chat(synthesisPrompt)
	if err != nil {
		return "", err
	}

	return finalAnswer, nil
}

func buildOraclePlanPrompt(prompt string) string {
	return fmt.Sprintf("You are in Oracle Mode. Create a multi-step research plan to answer this complex query: %s. Output ONLY the numbered steps.", prompt)
}

func buildOracleExecutionPrompt(plan, prompt string) string {
	return fmt.Sprintf("Execute this research plan using your available tools. Gather all necessary context.\nPlan:\n%s\n\nOriginal Query: %s", plan, prompt)
}

func buildOracleSynthesisPrompt(prompt, researchData string) string {
	return fmt.Sprintf("Based on the original query: '%s', and the gathered research data:\n%s\n\nProvide a comprehensive, authoritative, and definitive answer.", prompt, researchData)
}
