import { z } from 'zod';
import { t, publicProcedure, adminProcedure } from '../lib/trpc-core.js';
import { getMcpServer } from '../lib/mcpHelper.js';

export const mcpRouter = t.router({
    /** List all registered downstream MCP servers */
    listServers: publicProcedure.query(async () => {
        const mcp = getMcpServer();
        try {
            const servers = await (mcp as any).mcpAggregator?.listServers?.() ?? [];
            return servers.map((s: any) => ({
                name: s.name ?? 'unknown',
                status: s.status ?? 'unknown',
                toolCount: s.toolCount ?? s.tools?.length ?? 0,
                config: {
                    command: s.config?.command ?? s.command ?? '',
                    args: s.config?.args ?? s.args ?? [],
                    env: s.config?.env ? Object.keys(s.config.env) : [],
                },
            }));
        } catch {
            return [];
        }
    }),

    /** List all aggregated tools across servers */
    listTools: publicProcedure.query(async () => {
        const mcp = getMcpServer();
        try {
            const tools = await (mcp as any).mcpAggregator?.listAggregatedTools?.() ?? [];
            return tools.map((t: any) => ({
                name: t.name,
                description: t.description ?? '',
                server: t.server ?? 'unknown',
                inputSchema: t.inputSchema ?? null,
            }));
        } catch {
            return [];
        }
    }),

    /** Get aggregator status and stats */
    getStatus: publicProcedure.query(async () => {
        const mcp = getMcpServer();
        const aggregator = (mcp as any).mcpAggregator;
        if (!aggregator) return { initialized: false, serverCount: 0, toolCount: 0 };

        try {
            const servers = await aggregator.listServers?.() ?? [];
            const tools = await aggregator.listAggregatedTools?.() ?? [];
            return {
                initialized: true,
                serverCount: servers.length,
                toolCount: tools.length,
                connectedCount: servers.filter((s: any) => s.status === 'connected').length,
            };
        } catch {
            return { initialized: false, serverCount: 0, toolCount: 0 };
        }
    }),
});
