package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/robertpelloni/hypercode/foundation/assimilation"
	"github.com/robertpelloni/hypercode/foundation/compat"
	foundationpi "github.com/robertpelloni/hypercode/foundation/pi"
	"github.com/spf13/cobra"
)

var foundationCmd = &cobra.Command{
	Use:   "foundation",
	Short: "Inspect the Go foundation port and assimilation plan",
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

func init() {
	foundationInventoryCmd.Flags().Bool("json", false, "emit JSON")
	foundationCmd.AddCommand(foundationInventoryCmd)
	foundationCmd.AddCommand(foundationSpecCmd)
	foundationCmd.AddCommand(foundationToolsCmd)
	rootCmd.AddCommand(foundationCmd)
}
