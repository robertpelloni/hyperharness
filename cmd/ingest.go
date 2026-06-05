package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/robertpelloni/hyperharness/internal/ingest"
	"github.com/robertpelloni/hyperharness/internal/memory"
	"github.com/spf13/cobra"
)

var ingestCmd = &cobra.Command{
	Use:   "ingest [path]",
	Short: "Ingest and index files or directories into the knowledge base",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		tags, _ := cmd.Flags().GetStringSlice("tags")
		scopeStr, _ := cmd.Flags().GetString("scope")

		scope := memory.KnowledgeScope(scopeStr)
		if scope == "" {
			scope = memory.ScopeProject
		}

		// Initialize knowledge base
		kb, err := memory.NewKnowledgeBase("")
		if err != nil {
			return fmt.Errorf("failed to initialize knowledge base: %w", err)
		}
		processor := ingest.NewDataProcessor(kb)

		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			fmt.Printf("Ingesting directory: %s\n", path)
			return filepath.Walk(path, func(p string, d os.FileInfo, err error) error {
				if err != nil || d.IsDir() {
					return err
				}
				// Skip dotfiles and common ignores
				if filepath.HasPrefix(filepath.Base(p), ".") {
					return nil
				}

				fmt.Printf("  Processing %s...\n", p)
				return processor.IngestFile(p, tags, scope)
			})
		}

		fmt.Printf("Ingesting file: %s\n", path)
		return processor.IngestFile(path, tags, scope)
	},
}

func init() {
	ingestCmd.Flags().StringSliceP("tags", "t", []string{}, "Tags to apply to ingested entries")
	ingestCmd.Flags().StringP("scope", "s", "project", "Scope for the knowledge (global, project, session)")
	rootCmd.AddCommand(ingestCmd)
}
