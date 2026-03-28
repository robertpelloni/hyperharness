package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/robertpelloni/hypercode/tui"
)

var rootCmd = &cobra.Command{
	Use:   "hypercode",
	Short: "The ultimate AI CLI assistant, assimilated by Borg",
	Long:  `Hypercode is an AI pair programmer and terminal assistant with 100% feature parity with top tools.`,
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
