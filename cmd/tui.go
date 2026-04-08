package cmd

import (
	"log"

	"github.com/robertpelloni/hyperharness/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the highly-advanced Native Go BubbleTea Orchestrator Interface",
	Long:  "Starts the interactive visual TUI leveraging the Native Node.js-Parity HyperCode Orchestrator engine complete with AutoDrive autonomy.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("[Boot] Spinning up Native HyperCode TUI...")
		tui.StartREPL()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
