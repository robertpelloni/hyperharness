import sys

with open("internal/memory/sqlite_store.go", "r") as f:
    content = f.read()

struct_old = """type MemoryEntry struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`     // "knowledge", "conversation", "decision", "code_pattern"
	Tags      []string  `json:"tags"`
	Source    string    `json:"source"`
	Title     string    `json:"title"`
	Scope     string    `json:"scope"`    // "global", "project:<name>", "session:<id>"
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}"""

struct_new = """type MemoryEntry struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`     // "knowledge", "conversation", "decision", "code_pattern"
	Tags      []string  `json:"tags"`
	Source    string    `json:"source"`
	Title     string    `json:"title"`
	Scope     string    `json:"scope"`    // "global", "project:<name>", "session:<id>"
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Vector    []float32 `json:"vector,omitempty"`
}"""

content = content.replace(struct_old, struct_new)

table_old = """	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS memories (
			id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			type TEXT NOT NULL DEFAULT 'knowledge',
			tags TEXT NOT NULL DEFAULT '[]',
			source TEXT NOT NULL DEFAULT '',
			title TEXT NOT NULL DEFAULT '',
			scope TEXT NOT NULL DEFAULT 'global',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)"""

table_new = """	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS memories (
			id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			type TEXT NOT NULL DEFAULT 'knowledge',
			tags TEXT NOT NULL DEFAULT '[]',
			source TEXT NOT NULL DEFAULT '',
			title TEXT NOT NULL DEFAULT '',
			scope TEXT NOT NULL DEFAULT 'global',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			vector BLOB
		)
	`)"""

content = content.replace(table_old, table_new)

store_old = """	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO memories (id, content, type, tags, source, title, scope, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.ID, entry.Content, entry.Type, string(tagsJSON),
		entry.Source, entry.Title, entry.Scope,
		entry.CreatedAt, entry.UpdatedAt,
	)"""

store_new = """	var vectorJSON []byte
	if len(entry.Vector) > 0 {
		vectorJSON, _ = json.Marshal(entry.Vector)
	}

	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO memories (id, content, type, tags, source, title, scope, created_at, updated_at, vector)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.ID, entry.Content, entry.Type, string(tagsJSON),
		entry.Source, entry.Title, entry.Scope,
		entry.CreatedAt, entry.UpdatedAt, vectorJSON,
	)"""

content = content.replace(store_old, store_new)

search_old = """	// Try FTS5 first
	rows, err = s.db.Query(`
		SELECT m.id, m.content, m.type, m.tags, m.source, m.title, m.scope, m.created_at, m.updated_at
		FROM memories m
		JOIN memories_fts fts ON m.rowid = fts.rowid
		WHERE memories_fts MATCH ?
		ORDER BY rank
		LIMIT ?`,
		query, limit,
	)"""

search_new = """	// Try FTS5 first with time decay (older memories ranked lower)
	// SQLite FTS5 bm25 ranking returns lower scores for better matches, so we add a penalty
	// based on the age of the document (in days).
	rows, err = s.db.Query(`
		SELECT m.id, m.content, m.type, m.tags, m.source, m.title, m.scope, m.created_at, m.updated_at
		FROM memories m
		JOIN memories_fts fts ON m.rowid = fts.rowid
		WHERE memories_fts MATCH ?
		ORDER BY rank + (julianday('now') - julianday(m.created_at)) * 0.1
		LIMIT ?`,
		query, limit,
	)"""

content = content.replace(search_old, search_new)

with open("internal/memory/sqlite_store.go", "w") as f:
    f.write(content)
print("Updated sqlite_store.go")
