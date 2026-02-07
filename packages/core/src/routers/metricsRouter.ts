import { z } from 'zod';
import { t, getMcpServer } from '../lib/trpc-core.js';

export const metricsRouter = t.router({
    getStats: t.procedure.input(z.object({
        windowMs: z.number().optional()
    }).optional()).query(async ({ input }) => {
        const server = getMcpServer();
        return server.metricsService.getStats(input?.windowMs);
    }),

    track: t.procedure.input(z.object({
        type: z.string(),
        value: z.number(),
        tags: z.record(z.string()).optional()
    })).mutation(async ({ input }) => {
        const server = getMcpServer();
        server.metricsService.track(input.type, input.value, input.tags);
        return { success: true };
    })
});
