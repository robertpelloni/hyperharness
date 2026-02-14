import { z } from 'zod';
import { t, publicProcedure, getContextManager } from '../lib/trpc-core.js';

export const contextRouter = t.router({
    list: publicProcedure.query(() => {
        return getContextManager()?.list() ?? [];
    }),

    add: publicProcedure.input(z.object({
        filePath: z.string()
    })).mutation(({ input }) => {
        const contextManager = getContextManager();
        if (contextManager) return contextManager.add(input.filePath);
        return 'ContextManager not initialized';
    }),

    remove: publicProcedure.input(z.object({
        filePath: z.string()
    })).mutation(({ input }) => {
        const contextManager = getContextManager();
        if (contextManager) return contextManager.remove(input.filePath);
        return 'ContextManager not initialized';
    }),

    clear: publicProcedure.mutation(() => {
        const contextManager = getContextManager();
        if (contextManager) return contextManager.clear();
        return 'ContextManager not initialized';
    }),

    getPrompt: publicProcedure.query(() => {
        return getContextManager()?.getContextPrompt() ?? '';
    }),
});

