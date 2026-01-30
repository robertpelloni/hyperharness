import { z } from 'zod';
import { t, publicProcedure } from '../lib/trpc-core.js';

export const autoDevRouter = t.router({
    startLoop: publicProcedure.input(z.object({
        type: z.enum(['test', 'lint', 'build']),
        maxAttempts: z.number().min(1).max(10).default(5),
        target: z.string().optional(),
        command: z.string().optional()
    })).mutation(async ({ input }) => {
        // @ts-ignore
        if (global.mcpServerInstance?.autoDevService) {
            // @ts-ignore
            const id = await global.mcpServerInstance.autoDevService.startLoop(input);
            return { success: true, loopId: id };
        }
        throw new Error('AutoDevService not initialized');
    }),

    cancelLoop: publicProcedure.input(z.object({
        loopId: z.string()
    })).mutation(({ input }) => {
        // @ts-ignore
        if (global.mcpServerInstance?.autoDevService) {
            // @ts-ignore
            return global.mcpServerInstance.autoDevService.cancelLoop(input.loopId);
        }
        return false;
    }),

    getLoops: publicProcedure.query(() => {
        // @ts-ignore
        if (global.mcpServerInstance?.autoDevService) {
            // @ts-ignore
            return global.mcpServerInstance.autoDevService.getLoops();
        }
        return [];
    }),

    getLoop: publicProcedure.input(z.object({
        loopId: z.string()
    })).query(({ input }) => {
        // @ts-ignore
        if (global.mcpServerInstance?.autoDevService) {
            // @ts-ignore
            return global.mcpServerInstance.autoDevService.getLoop(input.loopId);
        }
        return null;
    }),

    clearCompleted: publicProcedure.mutation(() => {
        // @ts-ignore
        if (global.mcpServerInstance?.autoDevService) {
            // @ts-ignore
            return global.mcpServerInstance.autoDevService.clearCompleted();
        }
        return 0;
    }),
});
