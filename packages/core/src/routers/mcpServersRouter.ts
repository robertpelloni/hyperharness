import { z } from 'zod';
import { t, publicProcedure, adminProcedure } from '../lib/trpc-core.js';
import { mcpServersRepository } from '../db/repositories/index.js';
import fs from 'node:fs/promises';
import path from 'node:path';
import {
    McpServerCreateInputSchema,
    McpServerUpdateInputSchema
} from '../types/metamcp/index.js';

const MASTER_INDEX_PATH = path.join(process.cwd(), 'BORG_MASTER_INDEX.jsonc');

const stripJsonComments = (content: string) =>
    content.replace(/\\"|"(?:\\"|[^"])*"|(\/\/.*|\/\*[\s\S]*?\*\/)/g, (m, g) => g ? '' : m);

export const mcpServersRouter = t.router({
    list: publicProcedure.query(async () => {
        // TODO: Pass userId if auth context available
        return await mcpServersRepository.findAll();
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
});
