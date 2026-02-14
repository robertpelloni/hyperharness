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

        // Disconnect client if connected
        const clients = aggregator.clients as Map<string, { close?: () => Promise<void> | void }> | undefined;
        if (!clients) {
            return { success: false, error: 'Aggregator client map unavailable' };
        }
        const client = clients.get(input.name);
        if (client) {
            try {
                if (client.close) await client.close();
            } catch {
                /* ignore */
            }
            clients.delete(input.name);
        }

        // Remove from config file
        const fs = await import('fs');
        const configPath = aggregator.configPath;
        if (!configPath) {
            return { success: false, error: 'Aggregator config path unavailable' };
        }
        if (fs.existsSync(configPath)) {
            const config = JSON.parse(fs.readFileSync(configPath, 'utf8'));
            delete config[input.name];
            fs.writeFileSync(configPath, JSON.stringify(config, null, 2));
        }

        return { success: true };
    }),
});
