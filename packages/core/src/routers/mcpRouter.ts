import { z } from 'zod';
import { t, publicProcedure, adminProcedure, getMcpAggregator } from '../lib/trpc-core.js';

export const mcpRouter = t.router({
    /** List all registered downstream MCP servers */
    listServers: publicProcedure.query(async () => {
        const aggregator = getMcpAggregator();
        try {
            const servers = await aggregator?.listServers?.() ?? [];
            return servers.map((s) => ({
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
        const aggregator = getMcpAggregator();
        try {
            const tools = await aggregator?.listAggregatedTools?.() ?? [];
            return tools.map((tool) => ({
                name: tool.name,
                description: tool.description ?? '',
                server: tool.server ?? 'unknown',
                inputSchema: tool.inputSchema ?? null,
            }));
        } catch {
            return [];
        }
    }),

    traffic: publicProcedure.query(async () => {
        const aggregator = getMcpAggregator();
        try {
            return await aggregator?.getTrafficEvents?.() ?? [];
        } catch {
            return [];
        }
    }),

    searchTools: publicProcedure.input(z.object({
        query: z.string().default(''),
    })).query(async ({ input }) => {
        const aggregator = getMcpAggregator();
        try {
            const tools = await aggregator?.searchTools?.(input.query) ?? [];
            return tools.map((tool) => ({
                name: tool.name,
                description: tool.description ?? '',
                server: tool.server ?? 'unknown',
                inputSchema: tool.inputSchema ?? null,
            }));
        } catch {
            return [];
        }
    }),

    /** Get aggregator status and stats */
    getStatus: publicProcedure.query(async () => {
        const aggregator = getMcpAggregator();
        if (!aggregator) return { initialized: false, serverCount: 0, toolCount: 0, connectedCount: 0 };

        try {
            const servers = await aggregator.listServers?.() ?? [];
            const tools = await aggregator.listAggregatedTools?.() ?? [];
            return {
                initialized: true,
                serverCount: servers.length,
                toolCount: tools.length,
                connectedCount: servers.filter((s) => s.status === 'connected').length,
            };
        } catch {
            return { initialized: false, serverCount: 0, toolCount: 0, connectedCount: 0 };
        }
    }),

    /** Add a new downstream MCP server */
    addServer: adminProcedure.input(z.object({
        name: z.string().min(1),
        command: z.string().min(1),
        args: z.array(z.string()).optional().default([]),
        env: z.record(z.string()).optional().default({}),
    })).mutation(async ({ input }) => {
        const aggregator = getMcpAggregator();
        if (!aggregator) throw new Error('MCP Aggregator not initialized');
        if (!aggregator.addServerConfig) {
            throw new Error('MCP Aggregator addServerConfig unavailable');
        }

        await aggregator.addServerConfig(input.name, {
            command: input.command,
            args: input.args,
            env: input.env,
            enabled: true,
        });
        return { success: true, name: input.name };
    }),

    /** Remove a downstream MCP server */
    removeServer: adminProcedure.input(z.object({
        name: z.string(),
    })).mutation(async ({ input }) => {
        const aggregator = getMcpAggregator();
        if (!aggregator) throw new Error('MCP Aggregator not initialized');

        if (!aggregator.removeServerConfig) {
            throw new Error('MCP Aggregator removeServerConfig unavailable');
        }

        await aggregator.removeServerConfig(input.name);

        return { success: true };
    }),
});
