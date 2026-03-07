export interface MCPServerConfig {
    command: string;
    args: string[];
    enabled: boolean;
    env?: Record<string, string>;
}

export interface MCPToolDefinition {
    name: string;
    description?: string;
    inputSchema?: unknown;
    [key: string]: unknown;
}

export interface MCPAggregatedTool extends MCPToolDefinition {
    name: string;
    server: string;
    _originalName: string;
}

export interface MCPClientLike {
    connect(): Promise<void>;
    listTools(): Promise<MCPToolDefinition[]>;
    callTool(toolName: string, args: unknown): Promise<unknown>;
    close(): Promise<void>;
}

export type MCPServerStatus = 'stopped' | 'connecting' | 'connected' | 'error' | 'restarting';

export interface MCPServerState {
    name: string;
    config: MCPServerConfig;
    status: MCPServerStatus;
    toolCount: number;
    lastError?: string;
    restartCount: number;
}

export interface MCPTrafficEvent {
    server: string;
    method: 'tools/list' | 'tools/call';
    paramsSummary: string;
    latencyMs: number;
    success: boolean;
    timestamp: number;
    toolName?: string;
    error?: string;
}

export interface MCPAggregatorOptions {
    configPath?: string;
    createClient?: (name: string, config: MCPServerConfig) => MCPClientLike;
    restartDelayMs?: number;
    now?: () => number;
    maxTrafficEvents?: number;
}
