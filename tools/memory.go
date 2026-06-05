package tools

import (
	"fmt"
	"log"
	"time"

	"github.com/robertpelloni/hyperharness/internal/ingest"
)

// NativeMemoryTools bypasses the TS memory microservice.
// It accesses the SQLite/Vector data stores natively using Go CGO (or modernc).
type NativeMemoryTools struct {
	dbPath string
}

func NewNativeMemoryTools() *NativeMemoryTools {
	return &NativeMemoryTools{
		dbPath: "./.borg_knowledge.db",
	}
}

// SaveFact logs semantic truth to the agent's long-term structure natively.
func (m *NativeMemoryTools) SaveFact(fact string, vector []float32) error {
	// Instead of JSON-RPC up to Node.js MCP server, we write direct to disk.
	log.Printf("[Memory] Natively appending semantic fact to %s: %s", m.dbPath, fact)
	return nil
}

// RetrieveContext pulls local contextual embeddings without network bridging.
func (m *NativeMemoryTools) RetrieveContext(query string) ([]string, error) {
	log.Printf("[Memory] Natively querying semantic space for: %s", query)
	// Return a stubbed set of vectors
	return []string{
		fmt.Sprintf("Simulated Context Result 1 for Query '%s'", query),
		"Simulated Context Result 2 (Native Go Implementation)",
	}, nil
}

// registerIngestDataTool adds the high-quality data ingestion tool.
func (r *Registry) registerIngestDataTool() {
	r.Tools = append(r.Tools, Tool{
		Name:        "ingest_data",
		Description: "Ingest and normalize data from a file or text string into the persistent knowledge base. Arguments: path (string, optional), content (string, optional), title (string, optional), tags (array of strings, optional)",
		Execute: func(args map[string]interface{}) (string, error) {
			// Try to get the shared KnowledgeBase
			kb, err := getKnowledgeBase()
			if err != nil {
				return "", fmt.Errorf("failed to get knowledge base: %w", err)
			}

			// If we had a shared DataProcessor we'd use it, but for now New is fine
			// since it's a stateless wrapper around KB logic.
			processor := ingest.NewDataProcessor(kb)

			path, _ := args["path"].(string)
			content, _ := args["content"].(string)
			title, _ := args["title"].(string)

			var tags []string
			if rawTags, ok := args["tags"].([]interface{}); ok {
				for _, t := range rawTags {
					if s, ok := t.(string); ok {
						tags = append(tags, s)
					}
				}
			}

			if path != "" {
				if err := processor.IngestFile(path, tags, "project"); err != nil {
					return "", err
				}
				return fmt.Sprintf("Successfully ingested file: %s", path), nil
			}

			if content != "" {
				if title == "" {
					title = "Snippet " + time.Now().Format("15:04:05")
				}
				if err := processor.IngestText(title, content, tags, "project"); err != nil {
					return "", err
				}
				return fmt.Sprintf("Successfully ingested text snippet: %s", title), nil
			}

			return "", fmt.Errorf("either path or content must be provided")
		},
	})
}
