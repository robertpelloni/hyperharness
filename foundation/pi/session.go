package pi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/google/uuid"
)

const CurrentSessionVersion = 3

type SessionMetadata struct {
	SessionID     string `json:"sessionId"`
	Name          string `json:"name,omitempty"`
	WorkingDir    string `json:"workingDir"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updatedAt"`
	Version       int    `json:"version,omitempty"`
	ParentSession string `json:"parentSession,omitempty"`
}

type SessionEntry struct {
	ID        string          `json:"id"`
	ParentID  string          `json:"parentId,omitempty"`
	Kind      string          `json:"kind"`
	Role      string          `json:"role,omitempty"`
	Text      string          `json:"text,omitempty"`
	ToolName  string          `json:"toolName,omitempty"`
	ToolInput json.RawMessage `json:"toolInput,omitempty"`
	Result    *ToolResult     `json:"result,omitempty"`
	CreatedAt int64           `json:"createdAt"`

	ThinkingLevel string          `json:"thinkingLevel,omitempty"`
	Provider      string          `json:"provider,omitempty"`
	ModelID       string          `json:"modelId,omitempty"`
	Summary       string          `json:"summary,omitempty"`
	FromID        string          `json:"fromId,omitempty"`
	FirstKeptID   string          `json:"firstKeptEntryId,omitempty"`
	TokensBefore  int             `json:"tokensBefore,omitempty"`
	CustomType    string          `json:"customType,omitempty"`
	Data          json.RawMessage `json:"data,omitempty"`
	Display       *bool           `json:"display,omitempty"`
	TargetID      string          `json:"targetId,omitempty"`
	Label         string          `json:"label,omitempty"`
	Details       json.RawMessage `json:"details,omitempty"`
}

type SessionFile struct {
	Metadata SessionMetadata `json:"metadata"`
	Entries  []SessionEntry  `json:"entries"`
}

type SessionContext struct {
	Entries        []SessionEntry `json:"entries"`
	ThinkingLevel  string         `json:"thinkingLevel,omitempty"`
	Provider       string         `json:"provider,omitempty"`
	ModelID        string         `json:"modelId,omitempty"`
	SessionName    string         `json:"sessionName,omitempty"`
	LeafID         string         `json:"leafId,omitempty"`
	CompactionUsed bool           `json:"compactionUsed"`
}

type sessionRecord struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type SessionStore struct {
	baseDir string
}

func NewSessionStore(baseDir string) *SessionStore {
	return &SessionStore{baseDir: baseDir}
}

func DefaultSessionStore(cwd string) *SessionStore {
	return NewSessionStore(filepath.Join(cwd, ".hypercode", "foundation", "sessions"))
}

func (s *SessionStore) BaseDir() string {
	return s.baseDir
}

func (s *SessionStore) Create(name, workingDir string) (*SessionFile, error) {
	if err := os.MkdirAll(s.baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("create session directory: %w", err)
	}
	now := time.Now().UnixMilli()
	session := &SessionFile{
		Metadata: SessionMetadata{
			SessionID:  uuid.NewString(),
			Name:       name,
			WorkingDir: workingDir,
			CreatedAt:  now,
			UpdatedAt:  now,
			Version:    CurrentSessionVersion,
		},
	}
	if err := s.Save(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionStore) Save(session *SessionFile) error {
	if session == nil {
		return fmt.Errorf("session is nil")
	}
	if err := os.MkdirAll(s.baseDir, 0o755); err != nil {
		return fmt.Errorf("create session directory: %w", err)
	}
	if session.Metadata.SessionID == "" {
		session.Metadata.SessionID = uuid.NewString()
	}
	if session.Metadata.Version == 0 {
		session.Metadata.Version = CurrentSessionVersion
	}
	session.Metadata.UpdatedAt = time.Now().UnixMilli()
	path := s.Path(session.Metadata.SessionID)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create session file: %w", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	writeRecord := func(recordType string, value any) error {
		payload, err := json.Marshal(value)
		if err != nil {
			return err
		}
		line, err := json.Marshal(sessionRecord{Type: recordType, Data: payload})
		if err != nil {
			return err
		}
		if _, err := writer.Write(append(line, '\n')); err != nil {
			return err
		}
		return nil
	}
	if err := writeRecord("session", session.Metadata); err != nil {
		return fmt.Errorf("write session metadata: %w", err)
	}
	for _, entry := range session.Entries {
		if err := writeRecord("entry", entry); err != nil {
			return fmt.Errorf("write session entry: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush session file: %w", err)
	}
	return nil
}

func (s *SessionStore) Load(sessionID string) (*SessionFile, error) {
	file, err := os.Open(s.Path(sessionID))
	if err != nil {
		return nil, fmt.Errorf("open session file: %w", err)
	}
	defer file.Close()

	session := &SessionFile{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record sessionRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, fmt.Errorf("decode session record: %w", err)
		}
		switch record.Type {
		case "session":
			if err := json.Unmarshal(record.Data, &session.Metadata); err != nil {
				return nil, fmt.Errorf("decode session metadata: %w", err)
			}
			if session.Metadata.Version == 0 {
				session.Metadata.Version = 1
			}
		case "entry":
			var entry SessionEntry
			if err := json.Unmarshal(record.Data, &entry); err != nil {
				return nil, fmt.Errorf("decode session entry: %w", err)
			}
			session.Entries = append(session.Entries, entry)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan session file: %w", err)
	}
	return session, nil
}

func (s *SessionStore) AppendEntry(sessionID string, entry SessionEntry) (*SessionFile, error) {
	session, err := s.Load(sessionID)
	if err != nil {
		return nil, err
	}
	if entry.ID == "" {
		entry.ID = uuid.NewString()
	}
	if entry.CreatedAt == 0 {
		entry.CreatedAt = time.Now().UnixMilli()
	}
	if entry.ParentID == "" && len(session.Entries) > 0 {
		entry.ParentID = session.Entries[len(session.Entries)-1].ID
	}
	session.Entries = append(session.Entries, entry)
	if err := s.Save(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionStore) AppendThinkingLevelChange(sessionID, level string) (*SessionFile, error) {
	return s.AppendEntry(sessionID, SessionEntry{Kind: "thinking_level_change", ThinkingLevel: level})
}

func (s *SessionStore) AppendModelChange(sessionID, provider, modelID string) (*SessionFile, error) {
	return s.AppendEntry(sessionID, SessionEntry{Kind: "model_change", Provider: provider, ModelID: modelID})
}

func (s *SessionStore) AppendCompaction(sessionID, summary, firstKeptEntryID string, tokensBefore int, details any) (*SessionFile, error) {
	var raw json.RawMessage
	if details != nil {
		payload, err := json.Marshal(details)
		if err != nil {
			return nil, err
		}
		raw = payload
	}
	return s.AppendEntry(sessionID, SessionEntry{Kind: "compaction", Summary: summary, FirstKeptID: firstKeptEntryID, TokensBefore: tokensBefore, Details: raw})
}

func (s *SessionStore) AppendBranchSummary(sessionID, fromID, summary string, details any) (*SessionFile, error) {
	var raw json.RawMessage
	if details != nil {
		payload, err := json.Marshal(details)
		if err != nil {
			return nil, err
		}
		raw = payload
	}
	return s.AppendEntry(sessionID, SessionEntry{Kind: "branch_summary", FromID: fromID, Summary: summary, Details: raw})
}

func (s *SessionStore) AppendCustomEntry(sessionID, customType string, data any) (*SessionFile, error) {
	var raw json.RawMessage
	if data != nil {
		payload, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		raw = payload
	}
	return s.AppendEntry(sessionID, SessionEntry{Kind: "custom", CustomType: customType, Data: raw})
}

func (s *SessionStore) AppendCustomMessage(sessionID, customType, text string, display bool, details any) (*SessionFile, error) {
	var raw json.RawMessage
	if details != nil {
		payload, err := json.Marshal(details)
		if err != nil {
			return nil, err
		}
		raw = payload
	}
	return s.AppendEntry(sessionID, SessionEntry{Kind: "custom_message", CustomType: customType, Text: text, Display: &display, Details: raw})
}

func (s *SessionStore) AppendSessionInfo(sessionID, name string) (*SessionFile, error) {
	return s.AppendEntry(sessionID, SessionEntry{Kind: "session_info", Text: name})
}

func (s *SessionStore) AppendLabelChange(sessionID, targetID, label string) (*SessionFile, error) {
	return s.AppendEntry(sessionID, SessionEntry{Kind: "label", TargetID: targetID, Label: label})
}

func (s *SessionStore) Fork(sessionID, fromEntryID, name string) (*SessionFile, error) {
	session, err := s.Load(sessionID)
	if err != nil {
		return nil, err
	}
	forked, err := s.Create(name, session.Metadata.WorkingDir)
	if err != nil {
		return nil, err
	}
	forked.Metadata.ParentSession = session.Metadata.SessionID
	if fromEntryID == "" && len(session.Entries) > 0 {
		fromEntryID = session.Entries[len(session.Entries)-1].ID
	}
	for _, entry := range session.Entries {
		forked.Entries = append(forked.Entries, entry)
		if entry.ID == fromEntryID {
			break
		}
	}
	if err := s.Save(forked); err != nil {
		return nil, err
	}
	return forked, nil
}

func (s *SessionStore) Path(sessionID string) string {
	return filepath.Join(s.baseDir, sessionID+".jsonl")
}

func (s *SessionStore) List() ([]SessionMetadata, error) {
	if err := os.MkdirAll(s.baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("create session directory: %w", err)
	}
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, fmt.Errorf("read session directory: %w", err)
	}
	result := make([]SessionMetadata, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".jsonl" {
			continue
		}
		sessionID := entry.Name()[:len(entry.Name())-len(filepath.Ext(entry.Name()))]
		session, err := s.Load(sessionID)
		if err != nil {
			continue
		}
		result = append(result, session.Metadata)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].UpdatedAt > result[j].UpdatedAt
	})
	return result, nil
}

func (s *SessionStore) GetEntry(sessionID, entryID string) (*SessionEntry, error) {
	session, err := s.Load(sessionID)
	if err != nil {
		return nil, err
	}
	for i := range session.Entries {
		if session.Entries[i].ID == entryID {
			return &session.Entries[i], nil
		}
	}
	return nil, fmt.Errorf("entry not found: %s", entryID)
}

func (s *SessionStore) GetChildren(sessionID, parentID string) ([]SessionEntry, error) {
	session, err := s.Load(sessionID)
	if err != nil {
		return nil, err
	}
	children := make([]SessionEntry, 0)
	for _, entry := range session.Entries {
		if entry.ParentID == parentID {
			children = append(children, entry)
		}
	}
	sort.Slice(children, func(i, j int) bool { return children[i].CreatedAt < children[j].CreatedAt })
	return children, nil
}

func (s *SessionStore) GetBranch(sessionID, fromID string) ([]SessionEntry, error) {
	session, err := s.Load(sessionID)
	if err != nil {
		return nil, err
	}
	if fromID == "" && len(session.Entries) > 0 {
		fromID = session.Entries[len(session.Entries)-1].ID
	}
	byID := make(map[string]SessionEntry, len(session.Entries))
	for _, entry := range session.Entries {
		byID[entry.ID] = entry
	}
	branch := make([]SessionEntry, 0)
	currentID := fromID
	for currentID != "" {
		entry, ok := byID[currentID]
		if !ok {
			break
		}
		branch = append(branch, entry)
		currentID = entry.ParentID
	}
	for i, j := 0, len(branch)-1; i < j; i, j = i+1, j-1 {
		branch[i], branch[j] = branch[j], branch[i]
	}
	return branch, nil
}

func (s *SessionStore) GetLabel(sessionID, targetID string) (string, error) {
	session, err := s.Load(sessionID)
	if err != nil {
		return "", err
	}
	for i := len(session.Entries) - 1; i >= 0; i-- {
		entry := session.Entries[i]
		if entry.Kind == "label" && entry.TargetID == targetID {
			return entry.Label, nil
		}
	}
	return "", nil
}

func (s *SessionStore) GetSessionName(sessionID string) (string, error) {
	session, err := s.Load(sessionID)
	if err != nil {
		return "", err
	}
	for i := len(session.Entries) - 1; i >= 0; i-- {
		entry := session.Entries[i]
		if entry.Kind == "session_info" && entry.Text != "" {
			return entry.Text, nil
		}
	}
	return session.Metadata.Name, nil
}

func (s *SessionStore) BuildSessionContext(sessionID, leafID string) (*SessionContext, error) {
	session, err := s.Load(sessionID)
	if err != nil {
		return nil, err
	}
	branch, err := s.GetBranch(sessionID, leafID)
	if err != nil {
		return nil, err
	}
	ctx := &SessionContext{SessionName: session.Metadata.Name, LeafID: leafID}
	if ctx.SessionName == "" {
		if name, err := s.GetSessionName(sessionID); err == nil {
			ctx.SessionName = name
		}
	}
	firstKeptID := ""
	for _, entry := range branch {
		switch entry.Kind {
		case "model_change":
			if entry.Provider != "" {
				ctx.Provider = entry.Provider
			}
			if entry.ModelID != "" {
				ctx.ModelID = entry.ModelID
			}
		case "thinking_level_change":
			if entry.ThinkingLevel != "" {
				ctx.ThinkingLevel = entry.ThinkingLevel
			}
		case "compaction":
			ctx.CompactionUsed = true
			ctx.Entries = append(ctx.Entries, entry)
			firstKeptID = entry.FirstKeptID
		case "custom":
			// persistent extension state, not part of LLM context
		case "session_info", "label":
			// metadata only
		default:
			if firstKeptID != "" {
				if entry.ID == firstKeptID {
					firstKeptID = "__passed__"
				} else if firstKeptID != "__passed__" {
					continue
				}
			}
			ctx.Entries = append(ctx.Entries, entry)
		}
	}
	if len(branch) > 0 && ctx.LeafID == "" {
		ctx.LeafID = branch[len(branch)-1].ID
	}
	return ctx, nil
}
