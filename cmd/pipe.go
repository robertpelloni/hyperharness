package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/robertpelloni/hyperharness/agent"
	"github.com/robertpelloni/hyperharness/tools"
	"github.com/spf13/cobra"
)

var pipeCmd = &cobra.Command{
	Use:   "pipe [prompt]",
	Short: "Process a prompt through the LLM (non-interactive)",
	Long:  "Send a single prompt to the LLM and print the response. Reads stdin if piped, otherwise uses the argument.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var prompt string

		// Read from stdin if piped
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
				os.Exit(1)
			}
			prompt = string(data)
		}

		// Append any CLI argument as additional prompt
		if len(args) > 0 {
			if prompt != "" {
				prompt = prompt + "\n\n" + args[0]
			} else {
				prompt = args[0]
			}
		}

		if prompt == "" {
			fmt.Fprintln(os.Stderr, "Error: no prompt provided. Pass a prompt as an argument or pipe to stdin.")
			os.Exit(1)
		}

		// Try real LLM provider first
		provider, providerName, modelName, err := agent.ResolveProvider()
		if err != nil || provider == nil {
			fmt.Fprintln(os.Stderr, "No LLM API key detected. Set one of:")
			fmt.Fprintln(os.Stderr, "  ANTHROPIC_API_KEY, OPENAI_API_KEY, GEMINI_API_KEY, DEEPSEEK_API_KEY, OPENROUTER_API_KEY")
			os.Exit(1)
		}

		registry := tools.NewRegistry()
		loop := agent.NewAgentLoop(agent.AgentLoopConfig{
			Provider:     provider,
			ProviderName: providerName,
			Model:        modelName,
			WorkingDir:   ".",
			Registry:     registry,
		})

		result, err := loop.Run(context.Background(), prompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "LLM error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(pipeCmd)
}
