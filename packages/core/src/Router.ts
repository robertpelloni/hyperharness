import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";
import { ListToolsResult, CallToolResult, Tool } from "@modelcontextprotocol/sdk/types.js";

export interface RouterConfig {
    defaultProvider?: string;
}

export class Router {
    private clients: Map<string, Client> = new Map();
    private toolCache: Map<string, string> = new Map(); // toolName -> clientName

    constructor(private config: RouterConfig = {}) {
        console.log("Router initialized");
    }

    /**
     * Connects to a local MCP server via stdio.
     */
    async connectToServer(name: string, command: string, args: string[]) {
        if (this.clients.has(name)) return this.clients.get(name);

        const transport = new StdioClientTransport({
            command,
            args,
        });

        const client = new Client(
            {
                name: "borg-router",
                version: "0.1.0",
            },
            {
                capabilities: {},
            }
        );

        await client.connect(transport);
        this.clients.set(name, client);
        console.log(`Connected to MCP Server: ${name}`);
        return client;
    }

    async listTools(): Promise<Tool[]> {
        const allTools: Tool[] = [];
        this.toolCache.clear(); // Invalidate cache on refresh

        for (const [name, client] of this.clients.entries()) {
            try {
                const result = await client.listTools();
                // Simple namespacing: if tool is "read_file", keep it as is for now
                // In future: prefix with `${name}__` if needed
                for (const tool of result.tools) {
                    this.toolCache.set(tool.name, name);
                    allTools.push(tool);
                }
            } catch (e) {
                console.error(`Failed to list tools from ${name}`, e);
            }
        }
        return allTools;
    }

    async callTool(name: string, args: any): Promise<CallToolResult> {
        // 1. Try Cache
        let clientName = this.toolCache.get(name);

        // 2. If miss, refresh cache (lazy load)
        if (!clientName) {
            await this.listTools();
            clientName = this.toolCache.get(name);
        }

        // 3. Execute
        if (clientName) {
            const client = this.clients.get(clientName);
            if (client) {
                return (await client.callTool({ name, arguments: args })) as unknown as CallToolResult;
            }
        }

        throw new Error(`Tool ${name} not found in any connected MCP server.`);
    }

    getClient(name: string): Client | undefined {
        return this.clients.get(name);
    }
}
