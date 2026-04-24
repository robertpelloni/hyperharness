import sys

with open("tools/pi_exact_parity.go", "r") as f:
    content = f.read()

# Update getKnowledgeBase to use SQLiteMemoryStore
kb_old = """// Knowledge base singleton for tool integration
var (
	globalKB    *mem.KnowledgeBase
	globalKBMux sync.Once
)

func getKnowledgeBase() (*mem.KnowledgeBase, error) {
	var initErr error
	globalKBMux.Do(func() {
		globalKB, initErr = mem.NewKnowledgeBase("")
	})
	return globalKB, initErr
}"""

kb_new = """// Knowledge base singleton for tool integration
var (
	globalKB    *mem.KnowledgeBase
	globalKBMux sync.Once

	// SQLite FTS5 backend
	globalSQLiteKB    *mem.SQLiteMemoryStore
	globalSQLiteKBMux sync.Once
)

func getKnowledgeBase() (*mem.KnowledgeBase, error) {
	var initErr error
	globalKBMux.Do(func() {
		globalKB, initErr = mem.NewKnowledgeBase("")
	})
	return globalKB, initErr
}

func getSQLiteMemoryStore() (*mem.SQLiteMemoryStore, error) {
	var initErr error
	globalSQLiteKBMux.Do(func() {
		home, _ := os.UserHomeDir()
		dbPath := filepath.Join(home, ".hyperharness", "memory_fts.db")
		globalSQLiteKB, initErr = mem.NewSQLiteMemoryStore(dbPath)
	})
	return globalSQLiteKB, initErr
}"""

content = content.replace(kb_old, kb_new)

# Update memory_store Execute
store_old = """			// Try to use the actual knowledge base
			kb, err := getKnowledgeBase()
			if err == nil && kb != nil {
				entry := &mem.KnowledgeEntry{
					Title:   title,
					Content: content,
					Tags:    tags,
					Scope:   mem.KnowledgeScope(scope),
				}
				if err := kb.Store(entry); err != nil {
					return "", fmt.Errorf("failed to store knowledge: %w", err)
				}
				return fmt.Sprintf("Stored knowledge: %s (id: %s, scope: %s, tags: %v)", title, entry.ID, scope, tags), nil
			}

			return fmt.Sprintf("Stored knowledge: %s (scope: %s)", title, scope), nil"""

store_new = """			// Try to use the SQLite FTS5 knowledge base
			sqlStore, err := getSQLiteMemoryStore()
			if err == nil && sqlStore != nil {
				id := fmt.Sprintf("mem_%d", time.Now().UnixNano())
				entry := mem.MemoryEntry{
					ID:      id,
					Title:   title,
					Content: content,
					Tags:    tags,
					Scope:   scope,
					Type:    "knowledge",
				}
				if err := sqlStore.Store(entry); err != nil {
					return "", fmt.Errorf("failed to store knowledge in sqlite: %w", err)
				}
				return fmt.Sprintf("Stored knowledge in SQLite FTS5: %s (id: %s, scope: %s, tags: %v)", title, entry.ID, scope, tags), nil
			}

			// Fallback to JSON knowledge base
			kb, err := getKnowledgeBase()
			if err == nil && kb != nil {
				entry := &mem.KnowledgeEntry{
					Title:   title,
					Content: content,
					Tags:    tags,
					Scope:   mem.KnowledgeScope(scope),
				}
				if err := kb.Store(entry); err != nil {
					return "", fmt.Errorf("failed to store knowledge: %w", err)
				}
				return fmt.Sprintf("Stored knowledge: %s (id: %s, scope: %s, tags: %v)", title, entry.ID, scope, tags), nil
			}

			return fmt.Sprintf("Stored knowledge: %s (scope: %s)", title, scope), nil"""

content = content.replace(store_old, store_new)

# Update memory_search Execute
search_old = """			// Try to use the actual knowledge base
			kb, err := getKnowledgeBase()
			if err == nil && kb != nil {
				results := kb.Search(query, tags, mem.KnowledgeScope(scope))
				if len(results) > limit {
					results = results[:limit]
				}
				if len(results) == 0 {
					return fmt.Sprintf("No results found for query: %q", query), nil
				}
				var sb strings.Builder
				sb.WriteString(fmt.Sprintf("Found %d results for %q:\\n\\n", len(results), query))
				for i, entry := range results {
					sb.WriteString(fmt.Sprintf("%d. **%s** (scope: %s, tags: %v)\\n   %s\\n\\n",
						i+1, entry.Title, entry.Scope, entry.Tags, truncateString(entry.Content, 200)))
				}
				return sb.String(), nil
			}

			return fmt.Sprintf("Memory search for: %q. Configure memory subsystem for live results.", query), nil"""

search_new = """			limitInt := 10
			if limit > 0 {
				limitInt = limit
			}

			// Try to use the SQLite FTS5 knowledge base
			sqlStore, err := getSQLiteMemoryStore()
			if err == nil && sqlStore != nil {
				var results []mem.MemoryEntry
				var err error
				if len(tags) > 0 {
					results, err = sqlStore.SearchByTags(tags, limitInt)
				} else if scope != "" && scope != "global" {
					results, err = sqlStore.SearchByScope(query, scope, limitInt)
				} else {
					results, err = sqlStore.Search(query, limitInt)
				}

				if err != nil {
					return "", fmt.Errorf("SQLite FTS5 search failed: %w", err)
				}

				if len(results) == 0 {
					return fmt.Sprintf("No results found for query in SQLite FTS5: %q", query), nil
				}
				var sb strings.Builder
				sb.WriteString(fmt.Sprintf("Found %d results in SQLite FTS5 for %q:\\n\\n", len(results), query))
				for i, entry := range results {
					sb.WriteString(fmt.Sprintf("%d. **%s** (scope: %s, tags: %v)\\n   %s\\n\\n",
						i+1, entry.Title, entry.Scope, entry.Tags, truncateString(entry.Content, 200)))
				}
				return sb.String(), nil
			}

			// Fallback to JSON knowledge base
			kb, err := getKnowledgeBase()
			if err == nil && kb != nil {
				results := kb.Search(query, tags, mem.KnowledgeScope(scope))
				if len(results) > limitInt {
					results = results[:limitInt]
				}
				if len(results) == 0 {
					return fmt.Sprintf("No results found for query: %q", query), nil
				}
				var sb strings.Builder
				sb.WriteString(fmt.Sprintf("Found %d results for %q:\\n\\n", len(results), query))
				for i, entry := range results {
					sb.WriteString(fmt.Sprintf("%d. **%s** (scope: %s, tags: %v)\\n   %s\\n\\n",
						i+1, entry.Title, entry.Scope, entry.Tags, truncateString(entry.Content, 200)))
				}
				return sb.String(), nil
			}

			return fmt.Sprintf("Memory search for: %q. Configure memory subsystem for live results.", query), nil"""

content = content.replace(search_old.replace("\\n", "\n"), search_new.replace("\\n", "\n"))

# Add path imports
if '"path/filepath"' not in content:
    content = content.replace(
        '"context"\n\t"encoding/json"',
        '"context"\n\t"encoding/json"\n\t"path/filepath"'
    )

with open("tools/pi_exact_parity.go", "w") as f:
    f.write(content)
print("Updated pi_exact_parity.go")
