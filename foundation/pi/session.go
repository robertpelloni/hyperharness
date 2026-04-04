package pi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	LeafID        string `json:"leafId,omitempty"`
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

type BranchSummaryFileOps struct {
	ReadFiles     []string `json:"readFiles,omitempty"`
	ModifiedFiles []string `json:"modifiedFiles,omitempty"`
}

type BranchSummaryPreparation struct {
	TargetID               string               `json:"targetId"`
	OldLeafID              string               `json:"oldLeafId"`
	CommonAncestorID       string               `json:"commonAncestorId,omitempty"`
	EntriesToSummarize     []SessionEntry       `json:"entriesToSummarize"`
	SerializedConversation string               `json:"serializedConversation,omitempty"`
	FileOps                BranchSummaryFileOps `json:"fileOps,omitempty"`
	EstimatedTokens        int                  `json:"estimatedTokens,omitempty"`
	MaxTokens              int                  `json:"maxTokens,omitempty"`
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
	if entry.ParentID == "" {
		if session.Metadata.LeafID != "" {
			entry.ParentID = session.Metadata.LeafID
		} else if len(session.Entries) > 0 {
			entry.ParentID = session.Entries[len(session.Entries)-1].ID
		}
	}
	session.Entries = append(session.Entries, entry)
	session.Metadata.LeafID = entry.ID
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
	if fromEntryID == "" {
		if session.Metadata.LeafID != "" {
			fromEntryID = session.Metadata.LeafID
		} else if len(session.Entries) > 0 {
			fromEntryID = session.Entries[len(session.Entries)-1].ID
		}
	}
	for _, entry := range session.Entries {
		forked.Entries = append(forked.Entries, entry)
		if entry.ID == fromEntryID {
			forked.Metadata.LeafID = entry.ID
			break
		}
	}
	if err := s.Save(forked); err != nil {
		return nil, err
	}
	return forked, nil
}

func (s *SessionStore) GetLeafID(sessionID string) (string, error) {
	session, err := s.Load(sessionID)
	if err != nil {
		return "", err
	}
	if session.Metadata.LeafID != "" {
		return session.Metadata.LeafID, nil
	}
	if len(session.Entries) == 0 {
		return "", nil
	}
	return session.Entries[len(session.Entries)-1].ID, nil
}

func (s *SessionStore) Branch(sessionID, entryID string) (*SessionFile, error) {
	session, err := s.Load(sessionID)
	if err != nil {
		return nil, err
	}
	if entryID == "" {
		return nil, fmt.Errorf("entry id is required")
	}
	found := false
	for _, entry := range session.Entries {
		if entry.ID == entryID {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("entry not found: %s", entryID)
	}
	session.Metadata.LeafID = entryID
	if err := s.Save(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionStore) ResetLeaf(sessionID string) (*SessionFile, error) {
	session, err := s.Load(sessionID)
	if err != nil {
		return nil, err
	}
	session.Metadata.LeafID = ""
	if err := s.Save(session); err != nil {
		return nil, err
	}
	return session, nil
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
	if fromID == "" {
		if session.Metadata.LeafID != "" {
			fromID = session.Metadata.LeafID
		} else if len(session.Entries) > 0 {
			fromID = session.Entries[len(session.Entries)-1].ID
		}
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

func (s *SessionStore) GetCommonAncestor(sessionID, firstLeafID, secondLeafID string) (string, error) {
	if firstLeafID == "" || secondLeafID == "" {
		return "", nil
	}
	firstBranch, err := s.GetBranch(sessionID, firstLeafID)
	if err != nil {
		return "", err
	}
	secondBranch, err := s.GetBranch(sessionID, secondLeafID)
	if err != nil {
		return "", err
	}
	ancestor := ""
	for i := 0; i < len(firstBranch) && i < len(secondBranch); i++ {
		if firstBranch[i].ID != secondBranch[i].ID {
			break
		}
		ancestor = firstBranch[i].ID
	}
	return ancestor, nil
}

func (s *SessionStore) PrepareBranchSummary(sessionID, targetID string) (*BranchSummaryPreparation, error) {
	return s.PrepareBranchSummaryWithBudget(sessionID, targetID, 0)
}

func (s *SessionStore) PrepareBranchSummaryWithBudget(sessionID, targetID string, maxTokens int) (*BranchSummaryPreparation, error) {
	session, err := s.Load(sessionID)
	if err != nil {
		return nil, err
	}
	oldLeafID := session.Metadata.LeafID
	if oldLeafID == "" && len(session.Entries) > 0 {
		oldLeafID = session.Entries[len(session.Entries)-1].ID
	}
	if oldLeafID == "" || targetID == "" {
		return &BranchSummaryPreparation{TargetID: targetID, OldLeafID: oldLeafID, MaxTokens: maxTokens}, nil
	}
	commonAncestorID, err := s.GetCommonAncestor(sessionID, oldLeafID, targetID)
	if err != nil {
		return nil, err
	}
	oldBranch, err := s.GetBranch(sessionID, oldLeafID)
	if err != nil {
		return nil, err
	}
	entries := make([]SessionEntry, 0)
	startCollecting := commonAncestorID == ""
	for _, entry := range oldBranch {
		if !startCollecting {
			if entry.ID == commonAncestorID {
				startCollecting = true
			}
			continue
		}
		if entry.ID == commonAncestorID {
			continue
		}
		entries = append(entries, entry)
	}
	trimmed := trimEntriesToBudget(entries, maxTokens)
	fileOps := collectBranchFileOps(entries)
	serialized := serializeConversation(trimmed)
	estimated := estimateSerializedTokens(serialized)
	return &BranchSummaryPreparation{
		TargetID:               targetID,
		OldLeafID:              oldLeafID,
		CommonAncestorID:       commonAncestorID,
		EntriesToSummarize:     trimmed,
		SerializedConversation: serialized,
		FileOps:                fileOps,
		EstimatedTokens:        estimated,
		MaxTokens:              maxTokens,
	}, nil
}

func (s *SessionStore) BranchWithSummary(sessionID, targetID, summary string, details any) (*SessionFile, error) {
	prep, err := s.PrepareBranchSummary(sessionID, targetID)
	if err != nil {
		return nil, err
	}
	session, err := s.Branch(sessionID, targetID)
	if err != nil {
		return nil, err
	}
	if summary == "" {
		return session, nil
	}
	return s.AppendBranchSummary(sessionID, prep.OldLeafID, summary, details)
}

func trimEntriesToBudget(entries []SessionEntry, maxTokens int) []SessionEntry {
	if maxTokens <= 0 || len(entries) == 0 {
		return append([]SessionEntry(nil), entries...)
	}
	selected := make([]SessionEntry, 0, len(entries))
	used := 0
	for i := len(entries) - 1; i >= 0; i-- {
		serialized := serializeSingleEntry(entries[i])
		tokens := estimateSerializedTokens(serialized)
		if len(selected) > 0 && used+tokens > maxTokens {
			break
		}
		selected = append(selected, entries[i])
		used += tokens
	}
	for i, j := 0, len(selected)-1; i < j; i, j = i+1, j-1 {
		selected[i], selected[j] = selected[j], selected[i]
	}
	return selected
}

func estimateSerializedTokens(text string) int {
	if text == "" {
		return 0
	}
	return (len(text) + 3) / 4
}

func serializeConversation(entries []SessionEntry) string {
	parts := make([]string, 0, len(entries))
	for _, entry := range entries {
		serialized := serializeSingleEntry(entry)
		if serialized != "" {
			parts = append(parts, serialized)
		}
	}
	return strings.Join(parts, "\n")
}

func serializeSingleEntry(entry SessionEntry) string {
	switch entry.Kind {
	case "message":
		role := entry.Role
		if role == "" {
			role = "message"
		}
		return fmt.Sprintf("[%s]: %s", strings.Title(role), truncateSummaryText(entry.Text, 2000))
	case "tool_call":
		return fmt.Sprintf("[Assistant tool calls]: %s(%s)", entry.ToolName, truncateSummaryText(string(entry.ToolInput), 1000))
	case "tool_result":
		if entry.Result == nil {
			return "[Tool result]:"
		}
		return fmt.Sprintf("[Tool result]: %s", truncateSummaryText(flattenToolResultText(entry.Result), 2000))
	case "branch_summary":
		return fmt.Sprintf("[Branch summary]: %s", truncateSummaryText(entry.Summary, 2000))
	case "compaction":
		return fmt.Sprintf("[Compaction summary]: %s", truncateSummaryText(entry.Summary, 2000))
	case "custom_message":
		return fmt.Sprintf("[Custom message %s]: %s", entry.CustomType, truncateSummaryText(entry.Text, 2000))
	default:
		return ""
	}
}

func flattenToolResultText(result *ToolResult) string {
	if result == nil {
		return ""
	}
	parts := make([]string, 0, len(result.Content))
	for _, block := range result.Content {
		switch v := block.(type) {
		case TextContent:
			parts = append(parts, v.Text)
		case map[string]any:
			if text, ok := v["text"].(string); ok {
				parts = append(parts, text)
			}
		}
	}
	return strings.Join(parts, "\n")
}

func truncateSummaryText(text string, max int) string {
	if len(text) <= max {
		return text
	}
	if max <= 32 {
		return text[:max]
	}
	return text[:max] + fmt.Sprintf(" ...(truncated %d chars)", len(text)-max)
}

func collectBranchFileOps(entries []SessionEntry) BranchSummaryFileOps {
	readSet := map[string]bool{}
	modifiedSet := map[string]bool{}
	for _, entry := range entries {
		collectFileOpsFromEntry(entry, readSet, modifiedSet)
	}
	readFiles := make([]string, 0, len(readSet))
	for path := range readSet {
		readFiles = append(readFiles, path)
	}
	modifiedFiles := make([]string, 0, len(modifiedSet))
	for path := range modifiedSet {
		modifiedFiles = append(modifiedFiles, path)
	}
	sort.Strings(readFiles)
	sort.Strings(modifiedFiles)
	return BranchSummaryFileOps{ReadFiles: readFiles, ModifiedFiles: modifiedFiles}
}

func collectFileOpsFromEntry(entry SessionEntry, readSet, modifiedSet map[string]bool) {
	if entry.Kind == "tool_call" && len(entry.ToolInput) > 0 {
		var payload map[string]any
		if err := json.Unmarshal(entry.ToolInput, &payload); err == nil {
			if path, ok := payload["path"].(string); ok && path != "" {
				switch entry.ToolName {
				case "read":
					readSet[path] = true
				case "write", "edit":
					modifiedSet[path] = true
				}
			}
		}
	}
	if (entry.Kind == "branch_summary" || entry.Kind == "compaction") && len(entry.Details) > 0 {
		var details struct {
			ReadFiles     []string `json:"readFiles"`
			ModifiedFiles []string `json:"modifiedFiles"`
		}
		if err := json.Unmarshal(entry.Details, &details); err == nil {
			for _, path := range details.ReadFiles {
				if path != "" {
					readSet[path] = true
				}
			}
			for _, path := range details.ModifiedFiles {
				if path != "" {
					modifiedSet[path] = true
				}
			}
		}
	}
}

func DefaultStructuredSummaryTemplate(prep *BranchSummaryPreparation) string {
	if prep == nil {
		return ""
	}
	return fmt.Sprintf("## Goal\n[What the user is trying to accomplish]\n\n## Constraints & Preferences\n- [Requirements mentioned by user]\n\n## Progress\n### Done\n- [x] [Completed tasks]\n\n### In Progress\n- [ ] [Current work]\n\n### Blocked\n- [Issues, if any]\n\n## Key Decisions\n- **[Decision]**: [Rationale]\n\n## Next Steps\n1. [What should happen next]\n\n## Critical Context\n- [Data needed to continue]\n\n<read-files>\n%s\n</read-files>\n\n<modified-files>\n%s\n</modified-files>", strings.Join(prep.FileOps.ReadFiles, "\n"), strings.Join(prep.FileOps.ModifiedFiles, "\n"))
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
	if leafID == "" {
		leafID = session.Metadata.LeafID
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
