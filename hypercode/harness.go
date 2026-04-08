package hypercode

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/robertpelloni/hyperharness/agent"
	"github.com/robertpelloni/hyperharness/foundation/pi"
	intborg "github.com/robertpelloni/hyperharness/internal/borg"
	"github.com/robertpelloni/hyperharness/internal/config"
	contextmgr "github.com/robertpelloni/hyperharness/internal/context"
	"github.com/robertpelloni/hyperharness/internal/extensions"
	"github.com/robertpelloni/hyperharness/internal/fs"
	intmcp "github.com/robertpelloni/hyperharness/internal/mcp"
	"github.com/robertpelloni/hyperharness/internal/memory"
	"github.com/robertpelloni/hyperharness/internal/providers"
	"github.com/robertpelloni/hyperharness/internal/sessions"
	"github.com/robertpelloni/hyperharness/internal/skills"
	"github.com/robertpelloni/hyperharness/internal/subagents"
	"github.com/robertpelloni/hyperharness/rpc"
	"github.com/robertpelloni/hyperharness/security"
	"github.com/robertpelloni/hyperharness/tools"
)

// Harness is the unified integration point for the entire HyperHarness system.
// It wires together all subsystems: providers, sessions, memory, tools, MCP,
// skills, extensions, subagents, context management, RPC, and the agent runtime.
type Harness struct {
	// Core subsystems
	Config        *config.Settings
	BorgCore      *intborg.Core
	ProviderReg   *providers.Registry
	ToolRegistry  *tools.Registry
	SessionMgr    *sessions.Manager
	MemoryKB      *memory.KnowledgeBase
	MCPInternal   *intmcp.Registry
	ExtManager    *extensions.Manager
	SkillManager  *skills.Manager
	SubAgentMgr   *subagents.Manager
	ContextMgr    *contextmgr.Manager
	SecurityMgr   *security.PermissionManager

	// RPC (optional)
	RPCServer *rpc.Server

	// Foundation tools
	FoundationRuntime *pi.Runtime

	// State
	homeDir    string
	projectDir string
	mu         sync.RWMutex
}

// Options configures the harness creation.
type Options struct {
	ProjectDir string
	ConfigDir  string
	RPCAddr    string
}

// NewHarness creates and initializes the complete HyperHarness system.
func NewHarness(opts Options) (*Harness, error) {
	h := &Harness{}

	// 1. Resolve directories
	var err error
	h.homeDir, err = os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("resolve home: %w", err)
	}
	if opts.ProjectDir == "" {
		h.projectDir, _ = os.Getwd()
	} else {
		h.projectDir = opts.ProjectDir
	}

	configDir := opts.ConfigDir
	if configDir == "" {
		configDir = filepath.Join(h.homeDir, ".hyperharness")
	}

	// 2. Load configuration
	h.Config = config.Defaults()

	// 3. Initialize Borg Core
	h.BorgCore = intborg.NewCore(intborg.DefaultConfig())

	// 4. Initialize provider registry
	h.ProviderReg = providers.NewRegistry()

	// 5. Initialize tool registry (with all 122+ parity tools)
	h.ToolRegistry = tools.NewRegistry()

	// 6. Initialize session manager
	sessionDir := filepath.Join(configDir, "sessions")
	h.SessionMgr, err = sessions.NewManager(sessionDir)
	if err != nil {
		return nil, fmt.Errorf("session manager: %w", err)
	}

	// 7. Initialize memory system
	memoryPath := filepath.Join(configDir, "memory", "knowledge.json")
	h.MemoryKB, err = memory.NewKnowledgeBase(memoryPath)
	if err != nil {
		return nil, fmt.Errorf("memory system: %w", err)
	}

	// 8. Initialize MCP registries
	h.MCPInternal = intmcp.NewRegistry()

	// 9. Initialize extension manager
	extDir := filepath.Join(configDir, "extensions")
	h.ExtManager = extensions.NewManager(extDir)

	// 10. Initialize skill manager
	h.SkillManager = skills.NewManager()
	h.SkillManager.Discover(h.projectDir)
	for _, builtin := range skills.Builtins() {
		h.SkillManager.Register(builtin)
	}

	// 11. Initialize subagent manager
	h.SubAgentMgr = subagents.NewManager()

	// 12. Initialize context manager
	h.ContextMgr = contextmgr.NewManager(200, 200000)

	// 13. Initialize security
	h.SecurityMgr = security.NewPermissionManager(security.AutonomyLevelExecute)

	// 14. Initialize foundation runtime
	h.FoundationRuntime = pi.NewRuntime(h.projectDir, nil)

	// 15. Register Borg adapters for each subsystem
	h.BorgCore.RegisterAdapter(borgAdapter("memory", h.MemoryKB))
	h.BorgCore.RegisterAdapter(borgAdapter("tools", h.ToolRegistry))
	h.BorgCore.RegisterAdapter(borgAdapter("sessions", h.SessionMgr))
	h.BorgCore.RegisterAdapter(borgAdapter("extensions", h.ExtManager))
	h.BorgCore.RegisterAdapter(borgAdapter("mcp", h.MCPInternal))

	// 16. Load MCP servers from config files
	mcpConfigPath := filepath.Join(h.projectDir, ".hypercode", "mcp.json")
	if fs.FileExists(mcpConfigPath) {
		if err := h.MCPInternal.LoadFromFile(mcpConfigPath); err != nil {
			log.Printf("[Harness] MCP config load failed: %v", err)
		}
	}

	// 17. Load extensions config
	if err := h.ExtManager.LoadConfig(); err != nil {
		log.Printf("[Harness] Extensions config load failed: %v", err)
	}

	// 18. Start RPC server if configured
	if opts.RPCAddr != "" {
		h.RPCServer = rpc.NewServer()
		h.registerRPCHandlers()
		if err := h.RPCServer.Listen(opts.RPCAddr); err != nil {
			log.Printf("[Harness] RPC server failed: %v", err)
			h.RPCServer = nil
		}
	}

	// 19. Start Borg Core
	if err := h.BorgCore.Start(); err != nil {
		log.Printf("[Harness] Borg core start warning: %v", err)
	}

	return h, nil
}

// Close gracefully shuts down all subsystems.
func (h *Harness) Close() error {
	var errs []error

	if h.RPCServer != nil {
		if err := h.RPCServer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if err := h.MemoryKB.Save(); err != nil {
		errs = append(errs, err)
	}

	if err := h.SessionMgr.Save(); err != nil {
		errs = append(errs, err)
	}

	if err := h.ExtManager.SaveConfig(); err != nil {
		errs = append(errs, err)
	}

	h.BorgCore.Stop()

	if len(errs) > 0 {
		return fmt.Errorf("close errors: %v", errs)
	}
	return nil
}

// CreateSession creates a new agent session with full subsystem integration.
func (h *Harness) CreateSession(cwd, providerType, modelID string) (*sessions.Session, error) {
	session, err := h.SessionMgr.CreateSession(cwd, providerType, modelID)
	if err != nil {
		return nil, err
	}

	// Inject memory context
	memoryCtx := h.MemoryKB.BuildContextForAgent(cwd, nil)
	if memoryCtx != "" {
		h.ContextMgr.Inject(memoryCtx)
	}

	return session, nil
}

// ExecuteTool executes a tool by name through the unified tool registry.
func (h *Harness) ExecuteTool(name string, args map[string]interface{}) (string, error) {
	// Try internal tool registry first
	if tool, ok := h.ToolRegistry.Find(name); ok {
		return tool.Execute(args)
	}

	// Try internal MCP
	internalResult, err := h.MCPInternal.CallTool(name, args)
	if err == nil {
		var texts []string
		for _, c := range internalResult.Content {
			texts = append(texts, c.Text)
		}
		result := ""
		for i, t := range texts {
			if i > 0 {
				result += "\n"
			}
			result += t
		}
		return result, nil
	}

	return "", fmt.Errorf("tool not found: %s", name)
}

// StoreMemory stores a knowledge entry.
func (h *Harness) StoreMemory(title, content string, tags []string, scope memory.KnowledgeScope) error {
	return h.MemoryKB.Store(&memory.KnowledgeEntry{
		Title:   title,
		Content: content,
		Tags:    tags,
		Scope:   scope,
	})
}

// SearchMemory searches the knowledge base.
func (h *Harness) SearchMemory(keywords string, tags []string, scope memory.KnowledgeScope) []*memory.KnowledgeEntry {
	return h.MemoryKB.Search(keywords, tags, scope)
}

// Status returns a comprehensive status of all subsystems.
func (h *Harness) Status() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	status := map[string]interface{}{
		"projectDir": h.projectDir,
		"homeDir":    h.homeDir,
		"borg":       h.BorgCore.Status(),
		"tools":      len(h.ToolRegistry.Tools),
		"memory":     h.MemoryKB.Stats(),
		"sessions":   len(h.SessionMgr.ListSessions()),
		"extensions": len(h.ExtManager.List()),
		"skills":     len(h.SkillManager.List()),
		"subagents":  len(h.SubAgentMgr.ListTasks("")),
		"context":    h.ContextMgr.Status(),
	}

	if h.RPCServer != nil {
		status["rpc"] = map[string]interface{}{
			"clients": len(h.RPCServer.Clients()),
		}
	}

	return status
}

// registerRPCHandlers registers all RPC method handlers.
func (h *Harness) registerRPCHandlers() {
	h.RPCServer.Handle("status", func(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
		return json.Marshal(h.Status())
	})

	h.RPCServer.Handle("execute", func(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
		var req struct {
			Tool string                 `json:"tool"`
			Args map[string]interface{} `json:"args"`
		}
		if err := json.Unmarshal(params, &req); err != nil {
			return nil, err
		}
		result, err := h.ExecuteTool(req.Tool, req.Args)
		if err != nil {
			return nil, err
		}
		return json.Marshal(map[string]string{"result": result})
	})

	h.RPCServer.Handle("memory/store", func(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
		var req struct {
			Title   string   `json:"title"`
			Content string   `json:"content"`
			Tags    []string `json:"tags"`
			Scope   string   `json:"scope"`
		}
		if err := json.Unmarshal(params, &req); err != nil {
			return nil, err
		}
		err := h.StoreMemory(req.Title, req.Content, req.Tags, memory.KnowledgeScope(req.Scope))
		if err != nil {
			return nil, err
		}
		return json.Marshal(map[string]string{"status": "stored"})
	})

	h.RPCServer.Handle("memory/search", func(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
		var req struct {
			Keywords string   `json:"keywords"`
			Tags     []string `json:"tags"`
			Scope    string   `json:"scope"`
		}
		if err := json.Unmarshal(params, &req); err != nil {
			return nil, err
		}
		results := h.SearchMemory(req.Keywords, req.Tags, memory.KnowledgeScope(req.Scope))
		return json.Marshal(results)
	})
}

// borgAdapter wraps a subsystem as a Borg adapter.
type borgAdapterWrapper struct {
	name      string
	subsystem interface{}
}

func borgAdapter(name string, subsystem interface{}) *borgAdapterWrapper {
	return &borgAdapterWrapper{name: name, subsystem: subsystem}
}

func (a *borgAdapterWrapper) Name() string                 { return a.name }
func (a *borgAdapterWrapper) Start() error                 { return nil }
func (a *borgAdapterWrapper) Stop() error                  { return nil }
func (a *borgAdapterWrapper) Status() map[string]interface{} {
	return map[string]interface{}{"name": a.name, "active": true}
}

// Ensure Agent type is referenced (used by integration layer)
var _ *agent.Agent = nil
