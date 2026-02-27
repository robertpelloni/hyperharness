import { z } from 'zod';
import { t, publicProcedure, getMemoryManager, getAgentMemoryService } from '../lib/trpc-core.js';
import path from 'path';
import { MemoryExportImportService } from '../services/memory/MemoryExportImportService.js';

export const memoryRouter = t.router({
    saveContext: publicProcedure.input(z.object({
        source: z.string(),
        url: z.string(),
        title: z.string().optional(),
        content: z.string(),
        metadata: z.record(z.unknown()).optional()
    })).mutation(async ({ input }) => {
        const memoryManager = getMemoryManager();

        const id = await memoryManager.saveContext(input.content, {
            title: input.title,
            source: input.source,
            url: input.url,
            ...input.metadata
        });

        return { success: true, id };
    }),

    query: publicProcedure.input(z.object({
        query: z.string(),
        limit: z.number().optional().default(5)
    })).query(async ({ input }) => {
        return await getMemoryManager().search?.(input.query, input.limit) ?? [];
    }),

    listContexts: publicProcedure.query(async () => {
        return await getMemoryManager().listContexts();
    }),

    getContext: publicProcedure.input(z.object({
        id: z.string()
    })).query(async ({ input }) => {
        return await getMemoryManager().getContext?.(input.id) ?? null;
    }),

    deleteContext: publicProcedure.input(z.object({
        id: z.string()
    })).mutation(async ({ input }) => {
        const memoryManager = getMemoryManager();
        if (!memoryManager.deleteContext) return { success: false };
        await memoryManager.deleteContext(input.id);
        return { success: true };
    }),

    // --- Agent Memory Service (Tiered) ---

    getAgentStats: publicProcedure.query(async () => {
        return getAgentMemoryService()?.getStats() ?? null;
    }),

    searchAgentMemory: publicProcedure.input(z.object({
        query: z.string(),
        type: z.enum(['session', 'working', 'long_term']).optional(),
        limit: z.number().optional().default(10)
    })).query(async ({ input }) => {
        const service = getAgentMemoryService();
        if (!service) return [];
        return await service.search(input.query, {
            type: input.type,
            limit: input.limit
        });
    }),

    addFact: publicProcedure.input(z.object({
        content: z.string(),
        type: z.enum(['working', 'long_term']).default('working')
    })).mutation(async ({ input }) => {
        const service = getAgentMemoryService();
        if (!service) return { success: false };
        await service.add(input.content, input.type, 'user', { source: 'dashboard' });
        return { success: true };
    }),

    // --- Export / Import (Phase 70: Memory Multi-Backend) ---

    exportMemories: publicProcedure.input(z.object({
        userId: z.string().default('default'),
        format: z.enum(['json', 'csv', 'jsonl']).default('json'),
    })).query(async ({ input }) => {
        const manager = getMemoryManager();
        // The manager exposes the IMemoryProvider interface directly
        const exportService = new MemoryExportImportService(manager);
        const data = await exportService.exportAll(input.userId, input.format);
        return { data, format: input.format, exportedAt: new Date().toISOString() };
    }),

    importMemories: publicProcedure.input(z.object({
        userId: z.string().default('default'),
        format: z.enum(['json', 'csv', 'jsonl']).default('json'),
        data: z.string(),
    })).mutation(async ({ input }) => {
        const manager = getMemoryManager();
        const exportService = new MemoryExportImportService(manager);
        const result = await exportService.importBulk(input.data, input.format, input.userId);
        return { ...result, importedAt: new Date().toISOString() };
    }),
});

