//go:build !nosqlite
// +build !nosqlite

// Package internal provides SQLite-backed memory storage with FTS5.
// Ported from hypercode/go/internal/memorystore/search.go with enhancements.
//
// WHAT: Full-text search over stored memories using SQLite FTS5
// WHY: Fast, reliable memory retrieval for AI context without external dependencies
// HOW: SQLite database with FTS5 virtual table, automatic migration, batch operations
package memory

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	// Pure-Go SQLite driver (no CGO required)
	_ "github.com/glebarez/sqlite"
)

// SQLiteMemoryStore provides persistent memory storage with full-text search.
type SQLiteMemoryStore struct {
	db   *sql.DB
	path string
	mu   sync.Mutex
}

// MemoryEntry represents a stored memory with metadata.
type MemoryEntry struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"` // "knowledge", "conversation", "decision", "code_pattern"
	Tags      []string  `json:"tags"`
	Source    string    `json:"source"`
	Title     string    `json:"title"`
	Scope     string    `json:"scope"` // "global", "project:<name>", "session:<id>"
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Vector    []float32 `json:"vector,omitempty"`
}

// NewSQLiteMemoryStore creates or opens a SQLite memory store.
func NewSQLiteMemoryStore(dbPath string) (*SQLiteMemoryStore, error) {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrent performance
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	store := &SQLiteMemoryStore{
		db:   db,
		path: dbPath,
	}

	if err := store.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return store, nil
}

// migrate creates the database schema if it doesn't exist.
func (s *SQLiteMemoryStore) migrate() error {
	// Main memories table
	_, err := s.db.Exec(`
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
	`)
	if err != nil {
		return fmt.Errorf("create memories table: %w", err)
	}

	// FTS5 virtual table for full-text search
	_, err = s.db.Exec(`
		CREATE VIRTUAL TABLE IF NOT EXISTS memories_fts USING fts5(
			content,
			title,
			tags,
			type,
			scope,
			content='memories',
			content_rowid='rowid',
			tokenize='porter unicode61'
		)
	`)
	if err != nil {
		// FTS5 might not be available; create fallback index
		_, _ = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_memories_content ON memories(content)`)
		_, _ = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_memories_type ON memories(type)`)
		_, _ = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_memories_scope ON memories(scope)`)
		// Non-fatal: search will use LIKE fallback
		fmt.Fprintf(os.Stderr, "[memory] FTS5 not available, using LIKE fallback: %v\n", err)
		return nil
	}

	// FTS trigger: insert
	_, _ = s.db.Exec(`
		CREATE TRIGGER IF NOT EXISTS memories_ai AFTER INSERT ON memories BEGIN
			INSERT INTO memories_fts(rowid, content, title, tags, type, scope)
			VALUES (new.rowid, new.content, new.title, new.tags, new.type, new.scope);
		END
	`)

	// FTS trigger: delete
	_, _ = s.db.Exec(`
		CREATE TRIGGER IF NOT EXISTS memories_ad AFTER DELETE ON memories BEGIN
			INSERT INTO memories_fts(memories_fts, rowid, content, title, tags, type, scope)
			VALUES ('delete', old.rowid, old.content, old.title, old.tags, old.type, old.scope);
		END
	`)

	// FTS trigger: update
	_, _ = s.db.Exec(`
		CREATE TRIGGER IF NOT EXISTS memories_au AFTER UPDATE ON memories BEGIN
			INSERT INTO memories_fts(memories_fts, rowid, content, title, tags, type, scope)
			VALUES ('delete', old.rowid, old.content, old.title, old.tags, old.type, old.scope);
			INSERT INTO memories_fts(rowid, content, title, tags, type, scope)
			VALUES (new.rowid, new.content, new.title, new.tags, new.type, new.scope);
		END
	`)

	// Additional indexes for common queries
	_, _ = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_memories_created ON memories(created_at)`)
	_, _ = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_memories_source ON memories(source)`)

	return nil
}

// Store saves a memory entry.
func (s *SQLiteMemoryStore) Store(entry MemoryEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry.ID == "" {
		entry.ID = fmt.Sprintf("mem_%d", time.Now().UnixNano())
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	entry.UpdatedAt = time.Now().UTC()

	tagsJSON, _ := json.Marshal(entry.Tags)

	var vectorJSON []byte
	if len(entry.Vector) > 0 {
		vectorJSON, _ = json.Marshal(entry.Vector)
	}

	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO memories (id, content, type, tags, source, title, scope, created_at, updated_at, vector)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.ID, entry.Content, entry.Type, string(tagsJSON),
		entry.Source, entry.Title, entry.Scope,
		entry.CreatedAt, entry.UpdatedAt, vectorJSON,
	)

	return err
}

// Search performs a full-text search over memories.
func (s *SQLiteMemoryStore) Search(query string, limit int) ([]MemoryEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 50
	}

	var rows *sql.Rows
	var err error

	// Try FTS5 first with time decay (older memories ranked lower)
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
	)

	if err != nil {
		// Fallback to LIKE search
		likeQuery := "%" + query + "%"
		rows, err = s.db.Query(`
			SELECT id, content, type, tags, source, title, scope, created_at, updated_at
			FROM memories
			WHERE content LIKE ? OR title LIKE ? OR tags LIKE ?
			ORDER BY updated_at DESC
			LIMIT ?`,
			likeQuery, likeQuery, likeQuery, limit,
		)
		if err != nil {
			return nil, fmt.Errorf("search failed: %w", err)
		}
	}
	defer rows.Close()

	return scanMemoryEntries(rows)
}

// SearchByScope searches memories within a specific scope.
func (s *SQLiteMemoryStore) SearchByScope(query, scope string, limit int) ([]MemoryEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 50
	}

	likeQuery := "%" + query + "%"
	rows, err := s.db.Query(`
		SELECT id, content, type, tags, source, title, scope, created_at, updated_at
		FROM memories
		WHERE scope = ? AND (content LIKE ? OR title LIKE ?)
		ORDER BY updated_at DESC
		LIMIT ?`,
		scope, likeQuery, likeQuery, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMemoryEntries(rows)
}

// Get retrieves a memory by ID.
func (s *SQLiteMemoryStore) Get(id string) (*MemoryEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	row := s.db.QueryRow(`
		SELECT id, content, type, tags, source, title, scope, created_at, updated_at
		FROM memories WHERE id = ?`, id,
	)

	return scanMemoryEntry(row)
}

// Delete removes a memory by ID.
func (s *SQLiteMemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM memories WHERE id = ?", id)
	return err
}

// ListByType returns all memories of a given type.
func (s *SQLiteMemoryStore) ListByType(memType string, limit int) ([]MemoryEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 100
	}

	rows, err := s.db.Query(`
		SELECT id, content, type, tags, source, title, scope, created_at, updated_at
		FROM memories WHERE type = ?
		ORDER BY updated_at DESC LIMIT ?`, memType, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMemoryEntries(rows)
}

// ListRecent returns the most recently updated memories.
func (s *SQLiteMemoryStore) ListRecent(limit int) ([]MemoryEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 20
	}

	rows, err := s.db.Query(`
		SELECT id, content, type, tags, source, title, scope, created_at, updated_at
		FROM memories ORDER BY updated_at DESC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMemoryEntries(rows)
}

// Count returns the total number of stored memories.
func (s *SQLiteMemoryStore) Count() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM memories").Scan(&count)
	return count, err
}

// CountByType returns the number of memories of each type.
func (s *SQLiteMemoryStore) CountByType() (map[string]int, error) {
	rows, err := s.db.Query("SELECT type, COUNT(*) FROM memories GROUP BY type")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var memType string
		var count int
		if err := rows.Scan(&memType, &count); err == nil {
			result[memType] = count
		}
	}
	return result, nil
}

// Tags returns all unique tags across all memories.
func (s *SQLiteMemoryStore) Tags() ([]string, error) {
	rows, err := s.db.Query("SELECT DISTINCT tags FROM memories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagSet := make(map[string]bool)
	for rows.Next() {
		var tagsJSON string
		if err := rows.Scan(&tagsJSON); err == nil {
			var tags []string
			if json.Unmarshal([]byte(tagsJSON), &tags) == nil {
				for _, t := range tags {
					tagSet[t] = true
				}
			}
		}
	}

	result := make([]string, 0, len(tagSet))
	for t := range tagSet {
		result = append(result, t)
	}
	return result, nil
}

// SearchByTags returns memories matching any of the given tags.
func (s *SQLiteMemoryStore) SearchByTags(tags []string, limit int) ([]MemoryEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 50
	}

	// Build query with LIKE for each tag
	var conditions []string
	var args []interface{}
	for _, tag := range tags {
		conditions = append(conditions, "tags LIKE ?")
		args = append(args, "%\""+tag+"\"%")
	}

	query := fmt.Sprintf(`
		SELECT id, content, type, tags, source, title, scope, created_at, updated_at
		FROM memories WHERE %s
		ORDER BY updated_at DESC LIMIT ?`,
		strings.Join(conditions, " OR "),
	)
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMemoryEntries(rows)
}

// Close closes the database connection.
func (s *SQLiteMemoryStore) Close() error {
	return s.db.Close()
}

// Path returns the database file path.
func (s *SQLiteMemoryStore) Path() string {
	return s.path
}

// scanMemoryEntries scans multiple rows into MemoryEntry slices.
func scanMemoryEntries(rows *sql.Rows) ([]MemoryEntry, error) {
	var entries []MemoryEntry
	for rows.Next() {
		var entry MemoryEntry
		var tagsJSON string
		var createdAt, updatedAt string

		err := rows.Scan(
			&entry.ID, &entry.Content, &entry.Type, &tagsJSON,
			&entry.Source, &entry.Title, &entry.Scope,
			&createdAt, &updatedAt,
		)
		if err != nil {
			continue
		}

		_ = json.Unmarshal([]byte(tagsJSON), &entry.Tags)
		entry.CreatedAt, _ = time.Parse("2006-01-02T15:04:05Z", createdAt)
		entry.UpdatedAt, _ = time.Parse("2006-01-02T15:04:05Z", updatedAt)

		// Try alternate format if first fails
		if entry.CreatedAt.IsZero() {
			entry.CreatedAt, _ = time.Parse("2006-01-02 15:04:05+00:00", createdAt)
		}
		if entry.UpdatedAt.IsZero() {
			entry.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05+00:00", updatedAt)
		}

		entries = append(entries, entry)
	}
	return entries, nil
}

// scanMemoryEntry scans a single row into a MemoryEntry.
func scanMemoryEntry(row *sql.Row) (*MemoryEntry, error) {
	var entry MemoryEntry
	var tagsJSON string
	var createdAt, updatedAt string

	err := row.Scan(
		&entry.ID, &entry.Content, &entry.Type, &tagsJSON,
		&entry.Source, &entry.Title, &entry.Scope,
		&createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal([]byte(tagsJSON), &entry.Tags)
	entry.CreatedAt, _ = time.Parse("2006-01-02T15:04:05Z", createdAt)
	entry.UpdatedAt, _ = time.Parse("2006-01-02T15:04:05Z", updatedAt)

	return &entry, nil
}
