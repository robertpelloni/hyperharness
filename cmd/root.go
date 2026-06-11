package cmd

import (
	"fmt"
	"github.com/robertpelloni/hyperharness/tui"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "hyperharness",
	Short: "The ultimate AI CLI assistant with 100% feature parity",
	Long:  `HyperHarness is an AI pair programmer and terminal assistant with 100% feature parity with top tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		tui.StartREPL()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
