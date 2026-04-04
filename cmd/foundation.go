package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/robertpelloni/hypercode/foundation/assimilation"
	"github.com/robertpelloni/hypercode/foundation/compat"
	foundationpi "github.com/robertpelloni/hypercode/foundation/pi"
	"github.com/robertpelloni/hypercode/foundation/repomap"
	"github.com/spf13/cobra"
)

var foundationCmd = &cobra.Command{
	Use:   "foundation",
	Short: "Inspect and exercise the Go foundation port and assimilation plan",
}

var foundationInventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "List upstream systems targeted for assimilation",
	RunE: func(cmd *cobra.Command, args []string) error {
		items := assimilation.Inventory()
		asJSON, _ := cmd.Flags().GetBool("json")
		if asJSON {
			return json.NewEncoder(os.Stdout).Encode(items)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tCATEGORY\tLANGUAGE\tPATH\tTOP STRENGTH")
		for _, item := range items {
			top := ""
			if len(item.Strengths) > 0 {
				top = item.Strengths[0]
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", item.ID, item.Category, item.Language, item.UpstreamPath, top)
		}
		return w.Flush()
	},
}

var foundationSpecCmd = &cobra.Command{
	Use:   "spec",
	Short: "Print the current Pi-derived Go foundation spec",
	RunE: func(cmd *cobra.Command, args []string) error {
		spec := foundationpi.DefaultFoundationSpec()
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(spec)
	},
}

var foundationToolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Print exact-name tool contracts currently specified in the Go foundation",
	RunE: func(cmd *cobra.Command, args []string) error {
		catalog := compat.DefaultCatalog()
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(map[string]any{
			"count":   catalog.Count(),
			"sources": catalog.Sources(),
			"tools":   catalog.ContractsBySource("pi"),
		})
	},
}

var foundationExecCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a native foundation tool with exact-name contracts",
	RunE: func(cmd *cobra.Command, args []string) error {
		toolName, _ := cmd.Flags().GetString("tool")
		input, _ := cmd.Flags().GetString("input")
		sessionID, _ := cmd.Flags().GetString("session")
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		runtime := foundationpi.NewRuntime(cwd, nil)
		result, execErr := runtime.ExecuteTool(context.Background(), sessionID, toolName, json.RawMessage(input), nil)
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		payload := map[string]any{"tool": toolName, "result": result}
		if execErr != nil {
			payload["error"] = execErr.Error()
		}
		if err := enc.Encode(payload); err != nil {
			return err
		}
		return execErr
	},
}

var foundationRepomapCmd = &cobra.Command{
	Use:   "repomap",
	Short: "Generate a native ranked repo map inspired by Aider-style context condensation",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseDir, _ := cmd.Flags().GetString("dir")
		mentionedFiles, _ := cmd.Flags().GetStringSlice("mention-file")
		mentionedIdents, _ := cmd.Flags().GetStringSlice("mention-ident")
		maxFiles, _ := cmd.Flags().GetInt("max-files")
		includeTests, _ := cmd.Flags().GetBool("include-tests")
		result, err := repomap.Generate(repomap.Options{
			BaseDir:         baseDir,
			MentionedFiles:  mentionedFiles,
			MentionedIdents: mentionedIdents,
			MaxFiles:        maxFiles,
			IncludeTests:    includeTests,
		})
		if err != nil {
			return err
		}
		if asJSON, _ := cmd.Flags().GetBool("json"); asJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}
		_, err = fmt.Fprintln(os.Stdout, result.Map)
		return err
	},
}

var foundationSessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage native foundation sessions",
}

var foundationSessionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new native foundation session",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		runtime := foundationpi.NewRuntime(cwd, nil)
		session, err := runtime.CreateSession(name)
		if err != nil {
			return err
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(session)
	},
}

var foundationSessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List native foundation sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		runtime := foundationpi.NewRuntime(cwd, nil)
		sessions, err := runtime.ListSessions()
		if err != nil {
			return err
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(sessions)
	},
}

var foundationSessionShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a native foundation session",
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID, _ := cmd.Flags().GetString("session")
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		runtime := foundationpi.NewRuntime(cwd, nil)
		session, err := runtime.LoadSession(sessionID)
		if err != nil {
			return err
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(session)
	},
}

var foundationSessionForkCmd = &cobra.Command{
	Use:   "fork",
	Short: "Fork a native foundation session from an entry point",
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID, _ := cmd.Flags().GetString("session")
		entryID, _ := cmd.Flags().GetString("entry")
		name, _ := cmd.Flags().GetString("name")
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		runtime := foundationpi.NewRuntime(cwd, nil)
		session, err := runtime.ForkSession(sessionID, entryID, name)
		if err != nil {
			return err
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(session)
	},
}

func init() {
	foundationInventoryCmd.Flags().Bool("json", false, "emit JSON")
	foundationRepomapCmd.Flags().String("dir", ".", "base directory to scan")
	foundationRepomapCmd.Flags().StringSlice("mention-file", nil, "files to prioritize in ranking")
	foundationRepomapCmd.Flags().StringSlice("mention-ident", nil, "identifiers to prioritize in ranking")
	foundationRepomapCmd.Flags().Int("max-files", 40, "maximum files to include in output")
	foundationRepomapCmd.Flags().Bool("include-tests", false, "include test/spec files")
	foundationRepomapCmd.Flags().Bool("json", false, "emit JSON")
	foundationExecCmd.Flags().String("tool", "", "tool name to execute")
	foundationExecCmd.Flags().String("input", "{}", "tool input as JSON")
	foundationExecCmd.Flags().String("session", "", "optional session id to append tool execution to")
	_ = foundationExecCmd.MarkFlagRequired("tool")

	foundationSessionCreateCmd.Flags().String("name", "", "optional session display name")
	foundationSessionShowCmd.Flags().String("session", "", "session id")
	foundationSessionForkCmd.Flags().String("session", "", "session id")
	foundationSessionForkCmd.Flags().String("entry", "", "entry id to fork from (defaults to latest entry)")
	foundationSessionForkCmd.Flags().String("name", "", "optional fork display name")
	_ = foundationSessionShowCmd.MarkFlagRequired("session")
	_ = foundationSessionForkCmd.MarkFlagRequired("session")

	foundationSessionCmd.AddCommand(foundationSessionCreateCmd)
	foundationSessionCmd.AddCommand(foundationSessionListCmd)
	foundationSessionCmd.AddCommand(foundationSessionShowCmd)
	foundationSessionCmd.AddCommand(foundationSessionForkCmd)

	foundationCmd.AddCommand(foundationInventoryCmd)
	foundationCmd.AddCommand(foundationSpecCmd)
	foundationCmd.AddCommand(foundationToolsCmd)
	foundationCmd.AddCommand(foundationRepomapCmd)
	foundationCmd.AddCommand(foundationExecCmd)
	foundationCmd.AddCommand(foundationSessionCmd)
	rootCmd.AddCommand(foundationCmd)
}
