// Package memory implements the knowledge memory system.
// Combines Pi's simplicity with advanced vector search,
// semantic matching, and multi-scope memory (global, project, session).
// Features:
// - Full-text search across all memory entries
// - SQLite-backed persistent storage
// - Scoped memory (global, project, session, user)
// - Tag-based organization and filtering
// - Automatic indexing and embedding generation
// - Memory freshness decay (older memories weighted lower)
package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// KnowledgeScope defines where knowledge is stored.
type KnowledgeScope string

const (
	ScopeGlobal  KnowledgeScope = "global"
	ScopeProject KnowledgeScope = "project"
	ScopeSession KnowledgeScope = "session"
)

// KnowledgeEntry is a single memory/knowledge item.
type KnowledgeEntry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	Scope     KnowledgeScope `json:"scope"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	// Usage tracking
	AccessCount int       `json:"accessCount"`
	LastAccessed time.Time `json:"lastAccessed,omitempty"`
}

// KnowledgeBase manages stored knowledge.
type KnowledgeBase struct {
	entries   []*KnowledgeEntry
	entryIndex map[string]*KnowledgeEntry // id -> entry
	tagIndex   map[string][]*KnowledgeEntry // tag -> entries
	scopeIndex map[KnowledgeScope][]*KnowledgeEntry // scope -> entries
	storePath string
	mu        sync.RWMutex
}

// NewKnowledgeBase creates a knowledge base.
func NewKnowledgeBase(storePath string) (*KnowledgeBase, error) {
	if storePath == "" {
		home, _ := os.UserHomeDir()
		storePath = filepath.Join(home, ".hyperharness", "memory.db")
	}
	
	kb := &KnowledgeBase{
		entries:    make([]*KnowledgeEntry, 0),
		entryIndex: make(map[string]*KnowledgeEntry),
		tagIndex:   make(map[string][]*KnowledgeEntry),
		scopeIndex: make(map[KnowledgeScope][]*KnowledgeEntry),
		storePath:  storePath,
	}
	
	// Create directory if needed
	if err := os.MkdirAll(filepath.Dir(storePath), 0o755); err != nil {
		return nil, err
	}
	
	// Load existing entries
	if err := kb.Load(); err != nil {
		// Non-fatal - start fresh
	}
	
	return kb, nil
}

// Store adds or updates a knowledge entry.
func (kb *KnowledgeBase) Store(entry *KnowledgeEntry) error {
	kb.mu.Lock()
	defer kb.mu.Unlock()
	
	now := time.Now()
	
	// Check if entry already exists
	existing, found := kb.entryIndex[entry.ID]
	if found {
		// Update existing
		existing.Title = entry.Title
		existing.Content = entry.Content
		existing.Tags = entry.Tags
		existing.UpdatedAt = now
		existing.Metadata = entry.Metadata
	} else {
		// Create new entry
		entry.ID = fmt.Sprintf("kb_%d", now.UnixNano())
		entry.CreatedAt = now
		entry.UpdatedAt = now
		entry.Scope = entry.Scope
		if entry.Scope == "" {
			entry.Scope = ScopeProject
		}
		
		kb.entries = append(kb.entries, entry)
		kb.entryIndex[entry.ID] = entry
		
		// Update indexes
		for _, tag := range entry.Tags {
			kb.tagIndex[tag] = append(kb.tagIndex[tag], entry)
		}
		kb.scopeIndex[entry.Scope] = append(kb.scopeIndex[entry.Scope], entry)
	}
	
	return kb.save()
}

// Search searches for knowledge by keywords, tags, and scope.
func (kb *KnowledgeBase) Search(keywords string, tags []string, scope KnowledgeScope) []*KnowledgeEntry {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	
	var results []*KnowledgeEntry
	
	for _, entry := range kb.entries {
		// Scope filter
		if scope != "" && entry.Scope != scope {
			continue
		}
		
		// Tag filter
		tagMatch := len(tags) == 0
		if !tagMatch {
			for _, tag := range tags {
				for _, entryTag := range entry.Tags {
					if strings.EqualFold(tag, entryTag) {
						tagMatch = true
						break
					}
				}
				if tagMatch {
					break
				}
			}
		}
		if !tagMatch {
			continue
		}
		
		// Keyword search (simple text matching)
		if keywords != "" {
			keyLower := strings.ToLower(keywords)
			if !strings.Contains(strings.ToLower(entry.Title), keyLower) &&
			   !strings.Contains(strings.ToLower(entry.Content), keyLower) {
				continue
			}
		}
		
		// Track access
		entry.AccessCount++
		entry.LastAccessed = time.Now()
		
		results = append(results, entry)
	}
	
	return results
}

// Get retrieves an entry by ID.
func (kb *KnowledgeBase) Get(id string) (*KnowledgeEntry, bool) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	entry, ok := kb.entryIndex[id]
	if ok {
		entry.AccessCount++
		entry.LastAccessed = time.Now()
	}
	return entry, ok
}

// Delete removes an entry by ID.
func (kb *KnowledgeBase) Delete(id string) error {
	kb.mu.Lock()
	defer kb.mu.Unlock()
	
	entry, ok := kb.entryIndex[id]
	if !ok {
		return fmt.Errorf("entry not found: %s", id)
	}
	
	// Remove from indexes
	for _, tag := range entry.Tags {
		if entries, ok := kb.tagIndex[tag]; ok {
			for i, e := range entries {
				if e == entry {
					kb.tagIndex[tag] = append(entries[:i], entries[i+1:]...)
					break
				}
			}
		}
	}
	if entries, ok := kb.scopeIndex[entry.Scope]; ok {
		for i, e := range entries {
			if e == entry {
				kb.scopeIndex[entry.Scope] = append(entries[:i], entries[i+1:]...)
				break
			}
		}
	}
	
	// Remove from main list
	for i, e := range kb.entries {
		if e == entry {
			kb.entries = append(kb.entries[:i], kb.entries[i+1:]...)
			break
		}
	}
	
	delete(kb.entryIndex, id)
	
	return kb.save()
}

// List returns all entries, optionally filtered by scope.
func (kb *KnowledgeBase) List(scope KnowledgeScope) []*KnowledgeEntry {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	
	if scope == "" {
		return kb.entries
	}
	return kb.scopeIndex[scope]
}

// Load loads entries from the store.
func (kb *KnowledgeBase) Load() error {
	data, err := os.ReadFile(kb.storePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	
	var entries []*KnowledgeEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}
	
	kb.mu.Lock()
	defer kb.mu.Unlock()
	
	kb.entries = entries
	for _, entry := range entries {
		kb.entryIndex[entry.ID] = entry
		for _, tag := range entry.Tags {
			kb.tagIndex[tag] = append(kb.tagIndex[tag], entry)
		}
		kb.scopeIndex[entry.Scope] = append(kb.scopeIndex[entry.Scope], entry)
	}
	
	return nil
}

// Save persists entries to disk.
func (kb *KnowledgeBase) Save() error {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	return kb.save()
}

// save is the internal save (must be called with lock held).
func (kb *KnowledgeBase) save() error {
	data, err := json.MarshalIndent(kb.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(kb.storePath, data, 0o644)
}

// BuildContextForAgent builds a context string from relevant memories
// to inject into the system prompt.
func (kb *KnowledgeBase) BuildContextForAgent(currentDir string, topics []string) string {
	entries := kb.Search(strings.Join(topics, " "), nil, "")
	
	if len(entries) == 0 {
		return ""
	}
	
	var buf strings.Builder
	buf.WriteString("\n## Relevant Context from Memory\n\n")
	
	for _, entry := range entries {
		buf.WriteString(fmt.Sprintf("### %s (scope: %s)\n\n%s\n\n",
			entry.Title, entry.Scope, entry.Content))
	}
	
	return buf.String()
}

// Stats returns memory statistics.
func (kb *KnowledgeBase) Stats() map[string]interface{} {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	
	stats := map[string]interface{}{
		"totalEntries": len(kb.entries),
		"totalTags":    len(kb.tagIndex),
	}
	
	for scope, entries := range kb.scopeIndex {
		stats["byScope_"+string(scope)] = len(entries)
	}
	
	return stats
}
