import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { RequestOptions } from "@modelcontextprotocol/sdk/shared/protocol.js";
import {
    CallToolRequestSchema,
    CallToolResult,
    CompatibilityCallToolResultSchema,
    GetPromptRequestSchema,
    GetPromptResultSchema,
    ListPromptsRequestSchema,
    ListPromptsResult,
    ListPromptsResultSchema,
    ListResourcesRequestSchema,
    ListResourcesResult,
    ListResourcesResultSchema,
    ListResourceTemplatesRequestSchema,
    ListResourceTemplatesResult,
    ListResourceTemplatesResultSchema,
    ListToolsRequestSchema,
    ListToolsResultSchema,
    Prompt,
    ReadResourceRequestSchema,
    ReadResourceResultSchema,
    Resource,
    ResourceTemplate,
    Tool,
} from "@modelcontextprotocol/sdk/types.js";
import { z } from "zod";

import { toonSerializer } from "./stubs/toon.serializer.stub.js";
import { jsonConfigProvider } from "./config/JsonConfigProvider.js";
import { toolRegistry } from "./ToolRegistry.js";
import { SavedScriptConfig } from "../interfaces/IConfigProvider.js";

import { codeExecutorService } from "./CodeExecutorService.js";
import { deferredLoadingService } from "./deferred-loading.service.js";
// DB Repositories removed

// Ported Services
import { configImportService } from "./config-import.service.js";
import { configService } from "./config.service.js";
import { toolSetService } from "./tool-set.service.js";
import { ConnectedClient } from "./mcp-client.service.js";
import { getMcpServers } from "./fetch-metamcp.service.js";
import { mcpServerPool } from "./mcp-server-pool.service.js";
import { toolsSyncCache } from "./tools-sync-cache.service.js";
import { parseToolName } from "./tool-name-parser.service.js";
import { sanitizeName } from "./common-utils.js";
import { getAgentMemoryService, getMcpServer } from "../lib/trpc-core.js";

// Middleware
import {
    createFilterCallToolMiddleware,
    createFilterListToolsMiddleware,
} from "./metamcp-middleware/filter-tools.functional.js";
import { createPolicyMiddleware } from "./metamcp-middleware/policy.functional.js";
import {
    CallToolHandler,
    compose,
    ListToolsHandler,
    MetaMCPHandlerContext,
} from "./metamcp-middleware/functional-middleware.js";
import { createLoggingMiddleware } from "./metamcp-middleware/logging.functional.js";
import {
    createToolOverridesCallToolMiddleware,
    createToolOverridesListToolsMiddleware,
    mapOverrideNameToOriginal,
} from "./metamcp-middleware/tool-overrides.functional.js";

const toolsImplementations = {
    create: async ({
        tools,
        mcpServerUuid,
        serverName
    }: {
        tools: Tool[];
        mcpServerUuid: string;
        serverName: string;
    }) => {
        // Register tools in memory
        toolRegistry.registerTools(tools, mcpServerUuid, serverName);
    },
};

const agentService = {
    // ... (unchanged) ...
    runAgent: async (
        task: string,
        toolCallback: (toolName: string, toolArgs: unknown, meta?: Record<string, unknown>) => Promise<unknown>,
        _policyId?: string,
    ) => {
        const mcp = getMcpServer();
        const llm = mcp.llmService;

        const model = await llm.modelSelector.selectModel({ taskComplexity: 'medium' });
        const prompt = `You are an autonomous tool-using assistant.\nTask: ${task}\n\nReturn JSON only:\n{\n  "analysis": "short plan",\n  "tool": { "name": "optional_tool_name", "arguments": {} },\n  "final": "final response"\n}`;
        const response = await llm.generateText(
            model.provider,
            model.modelId,
            'You are a precise JSON-only planner. Use a tool only when required.',
            prompt,
        );

        const raw = response.content?.trim() ?? '{}';
        let parsed: {
            analysis?: string;
            tool?: { name?: string; arguments?: unknown };
            final?: string;
        };

        try {
            parsed = JSON.parse(raw);
        } catch {
            const fenced = raw.match(/```json\s*([\s\S]*?)```/i)?.[1];
            parsed = fenced ? JSON.parse(fenced) : { final: raw };
        }

        let toolResult: unknown = null;
        if (parsed.tool?.name) {
            toolResult = await toolCallback(parsed.tool.name, parsed.tool.arguments ?? {}, {});
        }

        return {
            analysis: parsed.analysis ?? 'No analysis provided.',
            toolResult,
            final: parsed.final ?? 'Task processed.',
            rawModelOutput: raw,
        };
    },
};
const toolSearchService = {
    searchTools: async (query: string, limit: number = 10) => {
        const q = query.trim().toLowerCase();
        const all = toolRegistry.getAllTools(); // Use registry
        const filtered = all
            .filter((rt) => {
                const name = String(rt.tool.name ?? '').toLowerCase();
                const description = String(rt.tool.description ?? '').toLowerCase();
                return name.includes(q) || description.includes(q);
            })
            .slice(0, Math.max(1, limit));

        return filtered.map((rt) => ({
            name: rt.tool.name,
            description: rt.tool.description,
            mcpServerUuid: rt.mcpServerUuid,
        }));
    },
};

const savedScriptService = {
    listScripts: async (): Promise<SavedScriptConfig[]> => {
        return jsonConfigProvider.loadScripts();
    },
    getScript: async (name: string) => {
        const scripts = await savedScriptService.listScripts();
        return scripts.find((script: SavedScriptConfig) => script.name === name) ?? null;
    },
    saveScript: async (name: string, code: string, description?: string) => {
        const script: SavedScriptConfig = {
            name,
            code,
            description: description ?? null,
        };
        await jsonConfigProvider.saveScript(script);
        return script;
    },
};

/**
 * Filter out tools that are overrides of existing tools to prevent duplicates in database
 * Uses the existing tool overrides cache for optimal performance
 */
async function filterOutOverrideTools(
    tools: Tool[],
    namespaceUuid: string,
    serverName: string,
): Promise<Tool[]> {
    if (!tools || tools.length === 0) {
        return tools;
    }

    const filteredTools: Tool[] = [];

    await Promise.allSettled(
        tools.map(async (tool) => {
            try {
                // Check if this tool name is actually an override name for an existing tool
                // by using the existing mapOverrideNameToOriginal function
                const fullToolName = `${sanitizeName(serverName)}__${tool.name}`;
                const originalName = await mapOverrideNameToOriginal(
                    fullToolName,
                    namespaceUuid,
                    true, // use cache
                );

                // If the original name is different from the current name,
                // this tool is an override and should be filtered out
                if (originalName !== fullToolName) {
                    // This is an override, skip it (don't save to database)
                    return;
                }

                // This is not an override, include it
                filteredTools.push(tool);
            } catch (error) {
                console.error(
                    `Error checking if tool ${tool.name} is an override:`,
                    error,
                );
                // On error, include the tool (fail-safe behavior)
                filteredTools.push(tool);
            }
        }),
    );

    return filteredTools;
}

/**
 * Allows executing a tool programmatically through the full MetaMCP proxy stack.
 * Set during attachTo initialization.
 */
export let executeProxiedTool: ((name: string, args: any) => Promise<CallToolResult>) | null = null;

export const attachTo = async (
    server: Server,
    namespaceUuid: string,
    sessionId: string,
    nativeToolDefinitions: Tool[],
    nativeToolHandler: (name: string, args: unknown) => Promise<CallToolResult>,
    includeInactiveServers: boolean = false,
) => {
    const toolToClient: Record<string, ConnectedClient> = {};
    const toolToServerUuid: Record<string, string> = {};
    const promptToClient: Record<string, ConnectedClient> = {};
    const resourceToClient: Record<string, ConnectedClient> = {};

    // Session-specific map of "loaded" tools that should be exposed to the client
    // Key: toolName, Value: true
    // Limited to 200 items to prevent unbounded growth in long sessions
    const loadedTools = new Set<string>();
    const hydratedToolSchemas = new Set<string>();
    const MAX_LOADED_TOOLS = 200;
    const deferredSchemaMode = process.env.MCP_DEFER_TOOL_SCHEMAS === 'true'
        || process.env.MCP_PROGRESSIVE_MODE === 'true';

    const addToLoadedTools = (name: string) => {
        if (loadedTools.size >= MAX_LOADED_TOOLS && !loadedTools.has(name)) {
            // Remove the first item (oldest) if limit reached - effectively a FIFO eviction
            const first = loadedTools.values().next().value;
            if (first) loadedTools.delete(first);
        }
        loadedTools.add(name);
    };

    const addToHydratedSchemas = (name: string) => {
        if (hydratedToolSchemas.size >= MAX_LOADED_TOOLS && !hydratedToolSchemas.has(name)) {
            const first = hydratedToolSchemas.values().next().value;
            if (first) hydratedToolSchemas.delete(first);
        }
        hydratedToolSchemas.add(name);
    };

    // Helper function to detect if a server is the same instance
    const isSameServerInstance = (
        params: { name?: string; url?: string | null },
        _serverUuid: string,
    ): boolean => {
        // Check if server name is exactly the same as our current server instance
        // This prevents exact recursive calls to the same server
        if (params.name === `metamcp-unified-${namespaceUuid}`) {
            return true;
        }

        return false;
    };

    // Server instance is passed in, so we don't create it here.
    // However, we might want to ensure capabilities are set if not already?
    // The caller (MCPServer) handles Server creation and capability definition.

    // Create the handler context
    const handlerContext: MetaMCPHandlerContext = {
        namespaceUuid,
        sessionId,
    };

    // ----------------------------------------------------------------------
    // Handler Implementations (Unwrapped)
    // ----------------------------------------------------------------------

    const originalListToolsHandler: ListToolsHandler = async (
        request,
        context,
    ) => {
        // 1. Meta Tools
        const metaTools: Tool[] = [
            {
                name: "search_tools",
                description: "Semantically search for available tools across all connected MCP servers. Use this to find tools for a specific task.",
                inputSchema: {
                    type: "object",
                    properties: {
                        query: {
                            type: "string",
                            description: "The search query describing what you want to do (e.g., 'manage github issues', 'query database')",
                        },
                        limit: {
                            type: "number",
                            description: "Max number of results to return (default: 10)",
                        },
                    },
                    required: ["query"],
                },
            },
            {
                name: "load_tool",
                description: "Load a specific tool by name into your context so you can use it. In progressive mode this loads lightweight metadata first; use get_tool_schema to hydrate the full parameter schema when needed.",
                inputSchema: {
                    type: "object",
                    properties: {
                        name: {
                            type: "string",
                            description: "The full name of the tool to load (e.g., 'github__create_issue')",
                        },
                    },
                    required: ["name"],
                },
            },
            {
                name: "get_tool_schema",
                description: "Explicitly fetch and hydrate the full JSON schema for a deferred tool after search/load, reducing default token overhead for sub-agents.",
                inputSchema: {
                    type: "object",
                    properties: {
                        name: {
                            type: "string",
                            description: "The full name of the tool to hydrate (e.g. 'github__create_issue').",
                        },
                    },
                    required: ["name"],
                },
            },
            {
                name: "run_code",
                description: "Execute TypeScript/JavaScript code in a secure sandbox. Use this to chain multiple tool calls, process data, or perform logic. You can call other tools from within this code using `await mcp.call('tool_name', args)`.",
                inputSchema: {
                    type: "object",
                    properties: {
                        code: {
                            type: "string",
                            description: "The TypeScript/JavaScript code to execute. Top-level await is supported.",
                        },
                    },
                    required: ["code"],
                },
            },
            {
                name: "run_python",
                description: "Execute Python 3 code. Suitable for data processing or simple scripts. No direct tool calling integration yet (use run_code for tool chaining).",
                inputSchema: {
                    type: "object",
                    properties: {
                        code: {
                            type: "string",
                            description: "The Python 3 code to execute.",
                        },
                    },
                    required: ["code"],
                },
            },
            {
                name: "run_agent",
                description: "Run an autonomous AI agent to perform a task. The agent will analyze your request, find relevant tools, write its own code, and execute it.",
                inputSchema: {
                    type: "object",
                    properties: {
                        task: {
                            type: "string",
                            description: "The natural language description of the task (e.g., 'Find the latest issue in repo X and summarize it').",
                        },
                        policyId: {
                            type: "string",
                            description: "Optional UUID of a Policy to restrict the agent's tool access.",
                        }
                    },
                    required: ["task"],
                },
            },
            {
                name: "save_script",
                description: "Save a successful code snippet as a reusable tool (Saved Script). The script will be available as a tool in future sessions.",
                inputSchema: {
                    type: "object",
                    properties: {
                        name: {
                            type: "string",
                            description: "The name of the new tool (must be unique, alphanumeric).",
                        },
                        description: {
                            type: "string",
                            description: "Description of what this script does.",
                        },
                        code: {
                            type: "string",
                            description: "The code to save.",
                        },
                    },
                    required: ["name", "code"],
                },
            },
            {
                name: "save_tool_set",
                description: "Save the currently loaded tools as a 'Tool Set' (Profile). This allows you to restore this working environment later.",
                inputSchema: {
                    type: "object",
                    properties: {
                        name: {
                            type: "string",
                            description: "Name of the tool set (e.g., 'web_dev', 'data_analysis').",
                        },
                        description: {
                            type: "string",
                            description: "Description of the tool set.",
                        },
                    },
                    required: ["name"],
                },
            },
            {
                name: "load_tool_set",
                description: "Load a previously saved Tool Set (Profile). This will add all tools in the set to your current context.",
                inputSchema: {
                    type: "object",
                    properties: {
                        name: {
                            type: "string",
                            description: "Name of the tool set to load.",
                        },
                    },
                    required: ["name"],
                },
            },
            {
                name: "import_mcp_config",
                description: "Import MCP servers from a JSON configuration file content (e.g., claude_desktop_config.json).",
                inputSchema: {
                    type: "object",
                    properties: {
                        configJson: {
                            type: "string",
                            description: "The content of the JSON configuration file.",
                        },
                    },
                    required: ["configJson"],
                },
            },
            // --- Phase 7: Memory Tools (claude-mem parity) ---
            {
                name: "save_memory",
                description: "Save a fact, instruction, or knowledge snippet to your long-term memory. Use this to remember user preferences, project details, or learnings for future sessions.",
                inputSchema: {
                    type: "object",
                    properties: {
                        content: {
                            type: "string",
                            description: "The information to remember.",
                        },
                        type: {
                            type: "string",
                            enum: ["working", "long_term"],
                            description: "Type of memory. 'long_term' persists across sessions. 'working' is for the current task context.",
                            default: "long_term"
                        },
                    },
                    required: ["content"],
                },
            },
            {
                name: "search_memory",
                description: "Search your memory for relevant facts or context.",
                inputSchema: {
                    type: "object",
                    properties: {
                        query: {
                            type: "string",
                            description: "The search query.",
                        },
                        limit: {
                            type: "number",
                            description: "Max results.",
                            default: 5
                        },
                    },
                    required: ["query"],
                },
            },
        ];

        // 2. Native Tools (Pre-loaded / Standard Lib)
        // Add native tools to metaTools list so they are always available
        // We prefix them if needed, or assume they are global.
        // For Borg, standard tools are global.
        if (nativeToolDefinitions && nativeToolDefinitions.length > 0) {
            metaTools.push(...nativeToolDefinitions);
        }

        // 3. Saved Scripts
        // Fetch user-defined saved scripts and expose them as tools
        try {
            const savedScripts = await savedScriptService.listScripts();
            const scriptTools: Tool[] = savedScripts.map((script: SavedScriptConfig) => ({
                name: `script__${script.name}`,
                description: `[Saved Script] ${script.description || "No description"}`,
                inputSchema: {
                    type: "object",
                    properties: {}, // Scripts currently take no args
                    additionalProperties: true
                }
            }));
            metaTools.push(...scriptTools);
        } catch (e) {
            console.error("Error fetching saved scripts", e);
        }

        const serverParams = await getMcpServers(
            context.namespaceUuid,
            includeInactiveServers,
        );

        const visitedServers = new Set<string>();
        const allServerEntries = Object.entries(serverParams);
        const allAvailableTools: Tool[] = [];

        console.log(`[DEBUG-TOOLS] 📋 Processing ${allServerEntries.length} servers`);

        await Promise.allSettled(
            allServerEntries.map(async ([mcpServerUuid, params]) => {
                if (visitedServers.has(mcpServerUuid)) return;

                const session = await mcpServerPool.getSession(
                    context.sessionId,
                    mcpServerUuid,
                    params,
                    namespaceUuid,
                );
                if (!session) {
                    console.log(`[DEBUG-TOOLS] ❌ No session for: ${params.name}`);
                    return;
                }

                const serverVersion = session.client.getServerVersion();
                const actualServerName = serverVersion?.name || params.name || "";
                const ourServerName = `metamcp-unified-${namespaceUuid}`;
                if (actualServerName === ourServerName) return;
                if (isSameServerInstance(params, mcpServerUuid)) return;

                visitedServers.add(mcpServerUuid);

                const capabilities = session.client.getServerCapabilities();
                if (!capabilities?.tools) return;

                const serverName = params.name || session.client.getServerVersion()?.name || "";

                try {
                    const allServerTools: Tool[] = [];
                    let cursor: string | undefined = undefined;
                    let hasMore = true;

                    while (hasMore) {
                        const result = await session.client.request(
                            {
                                method: "tools/list",
                                params: { cursor, _meta: request.params?._meta }
                            },
                            ListToolsResultSchema as unknown as import("zod").ZodType<any>
                        ) as import("@modelcontextprotocol/sdk/types.js").ListToolsResult;
                        if (result.tools) allServerTools.push(...result.tools);
                        cursor = result.nextCursor;
                        hasMore = !!result.nextCursor;
                    }

                    if (allServerTools.length > 0) {
                        try {
                            const toolsToSave = await filterOutOverrideTools(
                                allServerTools,
                                namespaceUuid,
                                serverName
                            );
                            if (toolsToSave.length > 0) {
                                // Use in-memory registry with namespaced names to match Proxy/Dashboard expectations
                                const namespacedTools = toolsToSave.map(t => ({
                                    ...t,
                                    name: `${sanitizeName(serverName)}__${t.name}`
                                }));
                                if (deferredSchemaMode) {
                                    namespacedTools.forEach((tool) => {
                                        deferredLoadingService.cacheToolSchema(tool.name, tool);
                                        const deferred = deferredLoadingService.createDeferredTool(tool, mcpServerUuid);
                                        const minimalTool = deferredLoadingService.createMinimalTool(deferred);
                                        toolRegistry.registerTool(minimalTool, mcpServerUuid, serverName, {
                                            isDeferred: true,
                                            fullTool: tool,
                                        });
                                    });
                                } else {
                                    toolRegistry.registerTools(
                                        namespacedTools,
                                        mcpServerUuid,
                                        serverName
                                    );
                                }
                            }
                        } catch (e) {
                            console.error("DB Save Error", e);
                        }
                    }

                    allServerTools.forEach(tool => {
                        const toolName = `${sanitizeName(serverName)}__${tool.name}`;
                        toolToClient[toolName] = session;
                        toolToServerUuid[toolName] = mcpServerUuid;
                        const namespacedTool = {
                            ...tool,
                            name: toolName
                        };

                        if (deferredSchemaMode) {
                            deferredLoadingService.cacheToolSchema(toolName, namespacedTool);
                        }

                        allAvailableTools.push({
                            ...(deferredSchemaMode
                                ? deferredLoadingService.createMinimalTool(
                                    deferredLoadingService.createDeferredTool(namespacedTool, mcpServerUuid)
                                )
                                : namespacedTool)
                        });
                    });

                } catch (error) {
                    console.error(`Error fetching tools from ${serverName}:`, error);
                }
            })
        );

        const resultTools = [...metaTools];

        allAvailableTools.forEach(tool => {
            if (loadedTools.has(tool.name)) {
                if (deferredSchemaMode && hydratedToolSchemas.has(tool.name)) {
                    const cachedSchema = deferredLoadingService.getCachedSchema(tool.name);
                    resultTools.push(cachedSchema ?? tool);
                    return;
                }

                resultTools.push(tool);
            }
        });

        return { tools: resultTools };
    };

    // ----------------------------------------------------------------------
    // Middleware Composition & Recursive Handling
    // ----------------------------------------------------------------------

    // We need a mechanism to allow _internalCallToolImpl to call the *final composed function* (recursiveCallToolHandler).
    // However, recursiveCallToolHandler is composed *using* a handler that calls _internalCallToolImpl.
    // This creates a circular dependency:
    // recursiveCallToolHandler -> middleware -> internalHandler -> recursiveCallToolHandler

    // To solve this cleanly, we use a mutable reference pattern.
    let recursiveCallToolHandlerRef: CallToolHandler | null = null;

    // The "delegate" handler simply calls whatever function is currently in the reference.
    const delegateHandler: CallToolHandler = async (request, context) => {
        if (!recursiveCallToolHandlerRef) {
            throw new Error("Handler not initialized");
        }
        return recursiveCallToolHandlerRef(request, context);
    };

    /**
     * Internal implementation that does the actual work.
     */
    const _internalCallToolImpl = async (
        name: string,
        args: any,
        meta?: any
    ): Promise<CallToolResult> => {

        // Check for TOON request
        const useToon = meta?.toon === true || meta?.toon === "true";

        const formatResult = (result: CallToolResult): CallToolResult => {
            if (!useToon) return result;

            // Attempt to compress JSON content
            const newContent = result.content.map(item => {
                if (item.type === "text") {
                    try {
                        // Try to parse as JSON first
                        const data = JSON.parse(item.text);
                        const serialized = toonSerializer.serialize(data);
                        return { ...item, text: serialized };
                    } catch (e) {
                        // Not JSON, return as is
                        return item;
                    }
                }
                return item;
            });

            return {
                ...result,
                content: newContent
            };
        };

        // 1. Meta Tools
        if (name === "search_tools") {
            const { query, limit } = args as { query: string; limit?: number };
            const results = await toolSearchService.searchTools(query, limit);
            return formatResult({
                content: [
                    {
                        type: "text",
                        text: JSON.stringify(results, null, 2),
                    },
                ],
            });
        }

        if (name === "load_tool") {
            const { name: toolName } = args as { name: string };
            if (toolToClient[toolName]) {
                addToLoadedTools(toolName);
                return {
                    content: [
                        {
                            type: "text",
                            text: `Tool '${toolName}' loaded.`,
                        }
                    ]
                };
            } else {
                return {
                    content: [
                        {
                            type: "text",
                            text: `Tool '${toolName}' not found.`,
                        },
                    ],
                    isError: true,
                };
            }
        }

        if (name === "get_tool_schema") {
            const { name: toolName } = args as { name: string };
            const cachedSchema = deferredLoadingService.getCachedSchema(toolName);

            if (!cachedSchema) {
                return {
                    content: [
                        {
                            type: "text",
                            text: `Schema for '${toolName}' is not available in cache. Load or rediscover the tool first.`,
                        },
                    ],
                    isError: true,
                };
            }

            addToLoadedTools(toolName);
            addToHydratedSchemas(toolName);

            const registeredTool = toolRegistry.getTool(toolName);
            if (registeredTool) {
                toolRegistry.registerTool(cachedSchema, registeredTool.mcpServerUuid, registeredTool.serverName, {
                    isDeferred: false,
                    fullTool: cachedSchema,
                });
            }

            return formatResult({
                content: [
                    {
                        type: "text",
                        text: JSON.stringify(cachedSchema.inputSchema ?? { type: "object", properties: {} }, null, 2),
                    },
                ],
            });
        }

        if (name === "save_script") {
            const { name: scriptName, code, description } = args as { name: string; code: string; description?: string };
            try {
                const saved = await savedScriptService.saveScript(scriptName, code, description);
                return {
                    content: [{ type: "text", text: `Script '${saved.name}' saved successfully.` }]
                };
            } catch (error: any) {
                return {
                    content: [{ type: "text", text: `Failed to save script: ${error.message}` }],
                    isError: true
                };
            }
        }

        if (name === "run_python") {
            const { code } = args as { code: string };
            try {
                const output = await codeExecutorService.executeCode(code);
                return {
                    content: [{ type: "text", text: output }]
                };
            } catch (error: any) {
                return {
                    content: [{ type: "text", text: `Execution failed: ${error.message}` }],
                    isError: true
                };
            }
        }

        if (name === "save_tool_set") {
            const { name: setName, description } = args as { name: string; description?: string };
            try {
                const toolsToSave = Array.from(loadedTools);
                if (toolsToSave.length === 0) {
                    return {
                        content: [{ type: "text", text: `No tools currently loaded to save.` }],
                        isError: true
                    };
                }
                const saved = await toolSetService.createToolSet(setName, toolsToSave, description);
                return {
                    content: [{ type: "text", text: `Tool Set '${saved.name}' saved with ${saved.tools.length} tools.` }]
                };
            } catch (error: any) {
                return {
                    content: [{ type: "text", text: `Failed to save tool set: ${error.message}` }],
                    isError: true
                };
            }
        }

        if (name === "load_tool_set") {
            const { name: setName } = args as { name: string };
            try {
                const set = await toolSetService.getToolSet(setName);
                if (!set) {
                    return {
                        content: [{ type: "text", text: `Tool Set '${setName}' not found.` }],
                        isError: true
                    };
                }

                // Add tools to loadedTools
                let count = 0;
                const missing = [];
                for (const toolName of set.tools) {
                    if (toolToClient[toolName]) {
                        addToLoadedTools(toolName);
                        count++;
                    } else {
                        missing.push(toolName);
                    }
                }

                let msg = `Loaded ${count} tools from set '${setName}'.`;
                if (missing.length > 0) {
                    msg += ` Warning: ${missing.length} tools could not be found (might be offline): ${missing.join(", ")}`;
                }

                return {
                    content: [{ type: "text", text: msg }]
                };
            } catch (error: any) {
                return {
                    content: [{ type: "text", text: `Failed to load tool set: ${error.message}` }],
                    isError: true
                };
            }
        }

        if (name === "import_mcp_config") {
            const { configJson } = args as { configJson: string };
            try {
                const result = await configImportService.importClaudeConfig(configJson);
                return {
                    content: [{
                        type: "text",
                        text: `Imported ${result.imported} servers. Skipped: ${JSON.stringify(result.skipped)}`
                    }]
                };
            } catch (error: any) {
                return {
                    content: [{ type: "text", text: `Import failed: ${error.message}` }],
                    isError: true
                };
            }
        }

        if (name === "run_code") {
            const { code } = args as { code: string };
            try {
                // RECURSION MAGIC: We pass the *delegate* handler to the sandbox.
                // This ensures that when the sandbox calls 'mcp.call', it goes through
                // the full middleware stack (logging, auditing, etc.) just like a request from the client.
                const result = await codeExecutorService.executeCode(
                    code,
                    async (toolName: string, toolArgs: any) => {
                        if (toolName === "run_code" || toolName === "run_agent") {
                            throw new Error("Cannot call run_code/run_agent from within sandbox");
                        }
                        // Call the delegate handler which points to the composed stack
                        const res = await delegateHandler({
                            method: "tools/call",
                            params: {
                                name: toolName,
                                arguments: toolArgs,
                                _meta: meta
                            }
                        }, handlerContext);

                        return res;
                    }
                );
                return formatResult({
                    content: [{ type: "text", text: JSON.stringify(result, null, 2) }]
                });
            } catch (error: any) {
                const errorInfo = {
                    message: error?.message || String(error),
                    name: error?.name || "Error",
                    stack: error?.stack || undefined,
                };
                return {
                    content: [{
                        type: "text",
                        text: `Error: ${errorInfo.message}\nName: ${errorInfo.name}${errorInfo.stack ? `\nStack: ${errorInfo.stack}` : ""}`
                    }],
                    isError: true
                };
            }
        }

        if (name === "run_agent") {
            const { task, policyId } = args as { task: string; policyId?: string };
            try {
                const result = await agentService.runAgent(
                    task,
                    async (toolName: string, toolArgs: any, meta: any) => {
                        if (toolName === "run_code" || toolName === "run_agent") {
                            throw new Error("Recursive agent calls restricted.");
                        }

                        // Inject policyId into meta if provided
                        const callMeta = { ...meta, ...(policyId ? { policyId } : {}) };

                        const res = await delegateHandler({
                            method: "tools/call",
                            params: {
                                name: toolName,
                                arguments: toolArgs,
                                _meta: callMeta
                            }
                        }, handlerContext);
                        return res;
                    },
                    policyId
                );
                return formatResult({
                    content: [{ type: "text", text: JSON.stringify(result, null, 2) }]
                });
            } catch (error: any) {
                const errorInfo = {
                    message: error?.message || String(error),
                    name: error?.name || "Error",
                    stack: error?.stack || undefined,
                };
                return {
                    content: [{
                        type: "text",
                        text: `Agent Error: ${errorInfo.message}\nName: ${errorInfo.name}${errorInfo.stack ? `\nStack: ${errorInfo.stack}` : ""}`
                    }],
                    isError: true
                };
            }
        }

        // 3. Memory Tools
        if (name === "save_memory") {
            const { content, type = "long_term" } = args as { content: string, type?: "working" | "long_term" };
            const service = getAgentMemoryService();
            if (!service) return { content: [{ type: "text", text: "Memory service not available." }], isError: true };

            await service.add(content, type, 'agent', { source: 'tool_call', sessionId });
            return {
                content: [{ type: "text", text: `Memory saved (${type}).` }]
            };
        }

        if (name === "search_memory") {
            const { query, limit = 5 } = args as { query: string, limit?: number };
            const service = getAgentMemoryService();
            if (!service) return { content: [{ type: "text", text: "Memory service not available." }], isError: true };

            const results = await service.search(query, { limit });
            return formatResult({
                content: [{ type: "text", text: JSON.stringify(results, null, 2) }]
            });
        }

        // 2. Saved Scripts execution
        if (name.startsWith("script__")) {
            const scriptName = name.replace("script__", "");
            const script = await savedScriptService.getScript(scriptName);

            if (script) {
                try {
                    // Execute saved script using the SAME logic as run_code
                    const result = await codeExecutorService.executeCode(
                        script.code,
                        async (toolName: string, toolArgs: any) => {
                            if (toolName === "run_code" || toolName.startsWith("script__")) {
                                throw new Error("Recursion restricted in saved scripts");
                            }
                            const res = await delegateHandler({
                                method: "tools/call",
                                params: {
                                    name: toolName,
                                    arguments: toolArgs,
                                    _meta: meta
                                }
                            }, handlerContext);
                            return res;
                        }
                    );
                    return formatResult({
                        content: [{ type: "text", text: JSON.stringify(result, null, 2) }]
                    });
                } catch (error: any) {
                    const errorInfo = {
                        message: error?.message || String(error),
                        name: error?.name || "Error",
                        stack: error?.stack || undefined,
                    };
                    return {
                        content: [{
                            type: "text",
                            text: `Script Error: ${errorInfo.message}\nName: ${errorInfo.name}${errorInfo.stack ? `\nStack: ${errorInfo.stack}` : ""}`
                        }],
                        isError: true
                    };
                }
            }
        }


        // 3. Native Tools
        const nativeTool = nativeToolDefinitions.find(t => t.name === name);
        if (nativeTool) {
            try {
                // Call native handler
                // We wrap it in formatResult just in case, though usually native tools return standard results
                const res = await nativeToolHandler(name, args);
                return formatResult(res);
            } catch (nativeErr: any) {
                return {
                    content: [{
                        type: "text",
                        text: `Native Tool Error: ${nativeErr.message}\nStack: ${nativeErr.stack}`
                    }],
                    isError: true
                };
            }
        }

        // 4. Downstream Tools
        const clientForTool = toolToClient[name];

        if (!clientForTool) {
            throw new Error(`Unknown tool: ${name}`);
        }

        const parsed = parseToolName(name);
        if (!parsed) throw new Error(`Invalid tool name: ${name}`);

        try {
            const abortController = new AbortController();
            const mcpRequestOptions: RequestOptions = {
                signal: abortController.signal,
                timeout: await configService.getMcpTimeout(),
            };

            const result = await clientForTool.client.request(
                {
                    method: "tools/call",
                    params: {
                        name: parsed.originalToolName,
                        arguments: args || {},
                        _meta: meta,
                    }
                },
                CompatibilityCallToolResultSchema,
                mcpRequestOptions
            );
            return formatResult(result as CallToolResult);
        } catch (error) {
            console.error(`Error calling ${name}:`, error);
            throw error;
        }
    };

    const implCallToolHandler: CallToolHandler = async (
        request,
        _context,
    ) => {
        const { name, arguments: args, _meta } = request.params;
        return await _internalCallToolImpl(name, args, _meta);
    };

    // Compose the middleware
    // The composed handler calls implCallToolHandler, which calls _internalCallToolImpl,
    // which might call delegateHandler, which calls recursiveCallToolHandlerRef (this composed stack).
    recursiveCallToolHandlerRef = compose(
        createLoggingMiddleware({ enabled: true }),
        createPolicyMiddleware({ enabled: true }), // Add Policy Middleware
        createFilterCallToolMiddleware({
            cacheEnabled: true,
            customErrorMessage: (toolName, reason) => `Access denied: ${reason}`,
        }),
        createToolOverridesCallToolMiddleware({ cacheEnabled: true }),
    )(implCallToolHandler);


    const listToolsWithMiddleware = compose(
        createToolOverridesListToolsMiddleware({
            cacheEnabled: true,
            persistentCacheOnListTools: true,
        }),
        createFilterListToolsMiddleware({ cacheEnabled: true }),
    )(originalListToolsHandler);


    // Set up the handlers
    server.setRequestHandler(ListToolsRequestSchema, async (request) => {
        return await listToolsWithMiddleware(request, handlerContext);
    });

    executeProxiedTool = async (name: string, args: any) => {
        if (!recursiveCallToolHandlerRef) throw new Error("Proxy Handler not initialized");
        return await recursiveCallToolHandlerRef({
            method: "tools/call",
            params: { name, arguments: args || {}, _meta: {} }
        } as any, handlerContext);
    };

    server.setRequestHandler(CallToolRequestSchema, async (request) => {
        if (!recursiveCallToolHandlerRef) throw new Error("Handler not initialized");
        return await recursiveCallToolHandlerRef(request, handlerContext);
    });

    // Get Prompt Handler
    server.setRequestHandler(GetPromptRequestSchema, async (request) => {
        const { name } = request.params;
        const clientForPrompt = promptToClient[name];

        if (!clientForPrompt) {
            throw new Error(`Unknown prompt: ${name}`);
        }

        try {
            // Parse the prompt name using shared utility
            const parsed = parseToolName(name);
            if (!parsed) {
                throw new Error(`Invalid prompt name format: ${name}`);
            }

            const promptName = parsed.originalToolName;
            const response = await clientForPrompt.client.request(
                {
                    method: "prompts/get",
                    params: {
                        name: promptName,
                        arguments: request.params.arguments || {},
                        _meta: request.params._meta,
                    },
                },
                GetPromptResultSchema,
            );

            return response;
        } catch (error) {
            console.error(
                `Error getting prompt through ${clientForPrompt.client.getServerVersion()?.name
                }:`,
                error,
            );
            throw error;
        }
    });

    // List Prompts Handler
    server.setRequestHandler(ListPromptsRequestSchema, async (request) => {
        const serverParams = await getMcpServers(
            namespaceUuid,
            includeInactiveServers,
        );
        const allPrompts: ListPromptsResult["prompts"] = [];

        // Track visited servers to detect circular references - reset on each call
        const visitedServers = new Set<string>();

        // Filter out self-referencing servers before processing
        const validPromptServers = Object.entries(serverParams).filter(
            ([uuid, params]) => {
                // Skip if we've already visited this server to prevent circular references
                if (visitedServers.has(uuid)) {
                    console.log(
                        `Skipping already visited server in prompts: ${params.name || uuid}`,
                    );
                    return false;
                }

                // Check if this server is the same instance to prevent self-referencing
                if (isSameServerInstance(params, uuid)) {
                    console.log(
                        `Skipping self-referencing server in prompts: ${params.name || uuid}`,
                    );
                    return false;
                }

                // Mark this server as visited
                visitedServers.add(uuid);
                return true;
            },
        );

        await Promise.allSettled(
            validPromptServers.map(async ([uuid, params]) => {
                const session = await mcpServerPool.getSession(
                    sessionId,
                    uuid,
                    params,
                    namespaceUuid,
                );
                if (!session) return;

                // Now check for self-referencing using the actual MCP server name
                const serverVersion = session.client.getServerVersion();
                const actualServerName = serverVersion?.name || params.name || "";
                const ourServerName = `metamcp-unified-${namespaceUuid}`;

                if (actualServerName === ourServerName) {
                    console.log(
                        `Skipping self-referencing MetaMCP server in prompts: "${actualServerName}"`,
                    );
                    return;
                }

                const capabilities = session.client.getServerCapabilities();
                if (!capabilities?.prompts) return;

                // Use name assigned by user, fallback to name from server
                const serverName =
                    params.name || session.client.getServerVersion()?.name || "";
                try {
                    const result = await session.client.request(
                        {
                            method: "prompts/list",
                            params: {
                                cursor: request.params?.cursor,
                                _meta: request.params?._meta,
                            },
                        },
                        ListPromptsResultSchema as unknown as import("zod").ZodType<any>,
                    );

                    if (result.prompts) {
                        const promptsWithSource = result.prompts.map((prompt: Prompt) => {
                            const promptName = `${sanitizeName(serverName)}__${prompt.name}`;
                            promptToClient[promptName] = session;
                            return {
                                ...prompt,
                                name: promptName,
                                description: prompt.description || "",
                            };
                        });
                        allPrompts.push(...promptsWithSource);
                    }
                } catch (error) {
                    console.error(`Error fetching prompts from: ${serverName}`, error);
                }
            }),
        );

        return {
            prompts: allPrompts,
            nextCursor: request.params?.cursor,
        };
    });

    // List Resources Handler
    server.setRequestHandler(ListResourcesRequestSchema, async (request) => {
        const serverParams = await getMcpServers(
            namespaceUuid,
            includeInactiveServers,
        );
        const allResources: ListResourcesResult["resources"] = [];

        // Track visited servers to detect circular references - reset on each call
        const visitedServers = new Set<string>();

        // Filter out self-referencing servers before processing
        const validResourceServers = Object.entries(serverParams).filter(
            ([uuid, params]) => {
                // Skip if we've already visited this server to prevent circular references
                if (visitedServers.has(uuid)) {
                    console.log(
                        `Skipping already visited server in resources: ${params.name || uuid}`,
                    );
                    return false;
                }

                // Check if this server is the same instance to prevent self-referencing
                if (isSameServerInstance(params, uuid)) {
                    console.log(
                        `Skipping self-referencing server in resources: ${params.name || uuid}`,
                    );
                    return false;
                }

                // Mark this server as visited
                visitedServers.add(uuid);
                return true;
            },
        );

        await Promise.allSettled(
            validResourceServers.map(async ([uuid, params]) => {
                const session = await mcpServerPool.getSession(
                    sessionId,
                    uuid,
                    params,
                    namespaceUuid,
                );
                if (!session) return;

                // Now check for self-referencing using the actual MCP server name
                const serverVersion = session.client.getServerVersion();
                const actualServerName = serverVersion?.name || params.name || "";
                const ourServerName = `metamcp-unified-${namespaceUuid}`;

                if (actualServerName === ourServerName) {
                    console.log(
                        `Skipping self-referencing MetaMCP server in resources: "${actualServerName}"`,
                    );
                    return;
                }

                const capabilities = session.client.getServerCapabilities();
                if (!capabilities?.resources) return;

                // Use name assigned by user, fallback to name from server
                const serverName =
                    params.name || session.client.getServerVersion()?.name || "";
                try {
                    const result = await session.client.request(
                        {
                            method: "resources/list",
                            params: {
                                cursor: request.params?.cursor,
                                _meta: request.params?._meta,
                            },
                        },
                        ListResourcesResultSchema as unknown as import("zod").ZodType<any>,
                    );

                    if (result.resources) {
                        const resourcesWithSource = result.resources.map((resource: Resource) => {
                            resourceToClient[resource.uri] = session;
                            return {
                                ...resource,
                                name: resource.name || "",
                            };
                        });
                        allResources.push(...resourcesWithSource);
                    }
                } catch (error) {
                    console.error(`Error fetching resources from: ${serverName}`, error);
                }
            }),
        );

        return {
            resources: allResources,
            nextCursor: request.params?.cursor,
        };
    });

    // Read Resource Handler
    server.setRequestHandler(ReadResourceRequestSchema, async (request) => {
        const { uri } = request.params;
        const clientForResource = resourceToClient[uri];

        if (!clientForResource) {
            throw new Error(`Unknown resource: ${uri}`);
        }

        try {
            return await clientForResource.client.request(
                {
                    method: "resources/read",
                    params: {
                        uri,
                        _meta: request.params._meta,
                    },
                },
                ReadResourceResultSchema,
            );
        } catch (error) {
            console.error(
                `Error reading resource through ${clientForResource.client.getServerVersion()?.name
                }:`,
                error,
            );
            throw error;
        }
    });

    // List Resource Templates Handler
    server.setRequestHandler(
        ListResourceTemplatesRequestSchema,
        async (request) => {
            const serverParams = await getMcpServers(
                namespaceUuid,
                includeInactiveServers,
            );
            const allTemplates: ResourceTemplate[] = [];

            // Track visited servers to detect circular references - reset on each call
            const visitedServers = new Set<string>();

            // Filter out self-referencing servers before processing
            const validTemplateServers = Object.entries(serverParams).filter(
                ([uuid, params]) => {
                    // Skip if we've already visited this server to prevent circular references
                    if (visitedServers.has(uuid)) {
                        console.log(
                            `Skipping already visited server in resource templates: ${params.name || uuid}`,
                        );
                        return false;
                    }

                    // Check if this server is the same instance to prevent self-referencing
                    if (isSameServerInstance(params, uuid)) {
                        console.log(
                            `Skipping self-referencing server in resource templates: ${params.name || uuid}`,
                        );
                        return false;
                    }

                    // Mark this server as visited
                    visitedServers.add(uuid);
                    return true;
                },
            );

            await Promise.allSettled(
                validTemplateServers.map(async ([uuid, params]) => {
                    const session = await mcpServerPool.getSession(
                        sessionId,
                        uuid,
                        params,
                        namespaceUuid,
                    );
                    if (!session) return;

                    // Now check for self-referencing using the actual MCP server name
                    const serverVersion = session.client.getServerVersion();
                    const actualServerName = serverVersion?.name || params.name || "";
                    const ourServerName = `metamcp-unified-${namespaceUuid}`;

                    if (actualServerName === ourServerName) {
                        console.log(
                            `Skipping self-referencing MetaMCP server in resource templates: "${actualServerName}"`,
                        );
                        return;
                    }

                    const capabilities = session.client.getServerCapabilities();
                    if (!capabilities?.resources) return;

                    const serverName =
                        params.name || session.client.getServerVersion()?.name || "";

                    try {
                        const result = await session.client.request(
                            {
                                method: "resources/templates/list",
                                params: {
                                    cursor: request.params?.cursor,
                                    _meta: request.params?._meta,
                                },
                            },
                            ListResourceTemplatesResultSchema as unknown as import("zod").ZodType<any>,
                        );

                        if (result.resourceTemplates) {
                            const templatesWithSource = result.resourceTemplates.map(
                                (template: ResourceTemplate) => ({
                                    ...template,
                                    name: template.name || "",
                                }),
                            );
                            allTemplates.push(...templatesWithSource);
                        }
                    } catch (error) {
                        console.error(
                            `Error fetching resource templates from: ${serverName}`,
                            error,
                        );
                        return;
                    }
                }),
            );

            return {
                resourceTemplates: allTemplates,
                nextCursor: request.params?.cursor,
            };
        },
    );

    const cleanup = async () => {
        // Cleanup is now handled by the pool
        await mcpServerPool.cleanupSession(sessionId);
    };

    return { cleanup };
};
