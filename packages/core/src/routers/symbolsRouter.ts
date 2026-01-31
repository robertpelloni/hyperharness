import { z } from 'zod';
import { t, publicProcedure } from '../lib/trpc-core.js';

export const symbolsRouter = t.router({
    list: publicProcedure.query(() => {
        // @ts-ignore
        if (global.mcpServerInstance?.symbolPinService) {
            // @ts-ignore
            return global.mcpServerInstance.symbolPinService.list();
        }
        return [];
    }),

    find: publicProcedure.input(z.object({
        query: z.string(),
        limit: z.number().default(10)
    })).query(async ({ input }) => {
        // @ts-ignore
        const mcp = global.mcpServerInstance;
        if (mcp && mcp.memoryManager) {
            return await mcp.memoryManager.searchSymbols(input.query, input.limit);
        }
        return [];
    }),

    pin: publicProcedure.input(z.object({
        name: z.string(),
        file: z.string(),
        type: z.enum(['function', 'class', 'method', 'variable', 'interface']),
        lineStart: z.number().optional(),
        lineEnd: z.number().optional(),
        notes: z.string().optional()
    })).mutation(({ input }) => {
        // @ts-ignore
        if (global.mcpServerInstance?.symbolPinService) {
            // @ts-ignore
            return global.mcpServerInstance.symbolPinService.pin(input);
        }
        throw new Error('SymbolPinService not initialized');
    }),

    unpin: publicProcedure.input(z.object({
        id: z.string()
    })).mutation(({ input }) => {
        // @ts-ignore
        if (global.mcpServerInstance?.symbolPinService) {
            // @ts-ignore
            return global.mcpServerInstance.symbolPinService.unpin(input.id);
        }
        return false;
    }),

    updatePriority: publicProcedure.input(z.object({
        id: z.string(),
        priority: z.number()
    })).mutation(({ input }) => {
        // @ts-ignore
        if (global.mcpServerInstance?.symbolPinService) {
            // @ts-ignore
            return global.mcpServerInstance.symbolPinService.updatePriority(input.id, input.priority);
        }
        return false;
    }),

    addNotes: publicProcedure.input(z.object({
        id: z.string(),
        notes: z.string()
    })).mutation(({ input }) => {
        // @ts-ignore
        if (global.mcpServerInstance?.symbolPinService) {
            // @ts-ignore
            return global.mcpServerInstance.symbolPinService.addNotes(input.id, input.notes);
        }
        return false;
    }),

    clear: publicProcedure.mutation(() => {
        // @ts-ignore
        if (global.mcpServerInstance?.symbolPinService) {
            // @ts-ignore
            return global.mcpServerInstance.symbolPinService.clear();
        }
        return 0;
    }),

    forFile: publicProcedure.input(z.object({
        filePath: z.string()
    })).query(({ input }) => {
        // @ts-ignore
        if (global.mcpServerInstance?.symbolPinService) {
            // @ts-ignore
            return global.mcpServerInstance.symbolPinService.forFile(input.filePath);
        }
        return [];
    }),
});
