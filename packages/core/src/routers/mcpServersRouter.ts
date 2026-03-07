import { z } from 'zod';
import { t, publicProcedure, adminProcedure } from '../lib/trpc-core.js';
import { mcpServersRepository } from '../db/repositories/index.js';
import fs from 'node:fs/promises';
import path from 'node:path';
import {
    McpServerCreateInputSchema,
    McpServerUpdateInputSchema
} from '../types/metamcp/index.js';
import { metaMCPBridge } from '../services/MetaMCPBridgeService.js';
import { clientConfigSyncService, SUPPORTED_MCP_CLIENTS } from '../mcp/clientConfigSync.js';

const MASTER_INDEX_PATH = path.join(process.cwd(), 'BORG_MASTER_INDEX.jsonc');

const stripJsonComments = (content: string) =>
    content.replace(/\\"|"(?:\\"|[^"])*"|(\/\/.*|\/\*[\s\S]*?\*\/)/g, (m, g) => g ? '' : m);

function getContextUserId(ctx: unknown): string | undefined {
    if (!ctx || typeof ctx !== 'object') {
        return undefined;
    }

    const session = (ctx as { session?: unknown }).session;
    if (!session || typeof session !== 'object') {
        return undefined;
    }

    const user = (session as { user?: unknown }).user;
    if (!user || typeof user !== 'object') {
        return undefined;
    }

    const userId = (user as { id?: unknown }).id;
    return typeof userId === 'string' && userId.trim().length > 0
        ? userId
        : undefined;
}

export const mcpServersRouter = t.router({
    list: publicProcedure.query(async ({ ctx }) => {
        const userId = getContextUserId(ctx);
        return await mcpServersRepository.findAll(userId);
    }),

    get: publicProcedure
        .input(z.object({ uuid: z.string() }))
        .query(async ({ input }) => {
            return await mcpServersRepository.findByUuid(input.uuid);
        }),

    create: adminProcedure
        .input(McpServerCreateInputSchema)
        .mutation(async ({ input }) => {
            return await mcpServersRepository.create(input);
        }),

    update: adminProcedure
        .input(McpServerUpdateInputSchema)
        .mutation(async ({ input }) => {
            if (!input.uuid) throw new Error("UUID required for update");
            return await mcpServersRepository.update(input);
        }),

    delete: adminProcedure
        .input(z.object({ uuid: z.string() }))
        .mutation(async ({ input }) => {
            return await mcpServersRepository.deleteByUuid(input.uuid);
        }),

    bulkImport: adminProcedure
        .input(z.array(McpServerCreateInputSchema))
        .mutation(async ({ input }) => {
            return await mcpServersRepository.bulkCreate(input);
        }),

    syncTargets: publicProcedure.query(async () => {
        return await clientConfigSyncService.listTargets();
    }),

    exportClientConfig: publicProcedure
        .input(z.object({
            client: z.enum(SUPPORTED_MCP_CLIENTS),
            path: z.string().optional(),
        }))
        .query(async ({ input }) => {
            return await clientConfigSyncService.previewSync(input.client, input.path);
        }),

    syncClientConfig: adminProcedure
        .input(z.object({
            client: z.enum(SUPPORTED_MCP_CLIENTS),
            path: z.string().optional(),
        }))
        .mutation(async ({ input }) => {
            return await clientConfigSyncService.syncClientConfig(input.client, input.path);
        }),

    registrySnapshot: publicProcedure.query(async () => {
        try {
            const content = await fs.readFile(MASTER_INDEX_PATH, 'utf-8');
            const parsed = JSON.parse(stripJsonComments(content)) as {
                categories?: Record<string, Array<Record<string, unknown>>>;
            };

            const categories = parsed.categories ?? {};
            const collected: Array<{
                id: string;
                name: string;
                url: string;
                category: string;
                description: string;
                tags: string[];
            }> = [];

            for (const [category, items] of Object.entries(categories)) {
                if (!Array.isArray(items)) continue;

                for (const raw of items) {
                    const url = String(raw?.url ?? '');
                    if (!url) continue;

                    const kind = String(raw?.kind ?? '').toLowerCase();
                    const tags = Array.isArray(raw?.tags)
                        ? raw.tags.map((t) => String(t).toLowerCase())
                        : [];

                    const isMcpLike =
                        category.toLowerCase().includes('mcp') ||
                        kind.includes('mcp') ||
                        tags.some((t) => t.includes('mcp')) ||
                        url.toLowerCase().includes('modelcontextprotocol') ||
                        url.toLowerCase().includes('mcp');

                    if (!isMcpLike) continue;

                    const id = String(raw?.id ?? url);
                    const name = String(raw?.name ?? raw?.id ?? url);
                    const description = String(raw?.description ?? raw?.summary ?? 'No description available.');

                    collected.push({
                        id,
                        name,
                        url,
                        category,
                        description,
                        tags,
                    });
                }
            }

            const deduped = new Map<string, (typeof collected)[number]>();
            for (const item of collected) {
                if (!deduped.has(item.url)) {
                    deduped.set(item.url, item);
                }
            }

            return Array.from(deduped.values()).slice(0, 300);
        } catch {
            return [];
        }
    }),

    /**
     * Lists MCP servers from the running MetaMCP backend (port 12009).
     * Returns an empty array gracefully if MetaMCP is not running.
     */
    listFromMetaMCP: publicProcedure.query(async () => {
        return await metaMCPBridge.listServers();
    }),

    /**
     * Returns the availability status of the MetaMCP backend.
     */
    metamcpStatus: publicProcedure.query(async () => {
        const available = await metaMCPBridge.isAvailable();
        return { available, url: 'http://localhost:12009' };
    }),

    /**
     * Registers a new MCP server in MetaMCP's database via the bridge.
     * MetaMCP handles connection management automatically.
     */
    createInMetaMCP: adminProcedure
        .input(z.object({
            name: z.string().min(1),
            description: z.string().optional(),
            type: z.enum(['STDIO', 'SSE', 'STREAMABLE_HTTP']),
            command: z.string().optional(),
            args: z.array(z.string()).optional(),
            url: z.string().optional(),
            env: z.record(z.string()).optional(),
        }))
        .mutation(async ({ input }) => {
            const result = await metaMCPBridge.createServer(input);
            if (!result) throw new Error('MetaMCP backend unavailable or request failed');
            return result;
        }),

    /**
     * Removes an MCP server from MetaMCP's database via the bridge.
     */
    deleteFromMetaMCP: adminProcedure
        .input(z.object({ uuid: z.string() }))
        .mutation(async ({ input }) => {
            const ok = await metaMCPBridge.deleteServer(input.uuid);
            return { success: ok };
        }),
});
