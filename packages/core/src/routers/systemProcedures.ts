import { z } from 'zod';
import {
    t,
    publicProcedure,
    adminProcedure,
    getMcpServer,
    getProjectTracker,
    getLspService,
} from '../lib/trpc-core.js';

export const systemProcedures = {
    health: publicProcedure.query(() => {
        return { status: 'running', service: '@borg/core' };
    }),
    getTaskStatus: publicProcedure
        .input(z.object({ taskId: z.string().optional() }))
        .query(() => {
            const tracker = getProjectTracker();
            if (!tracker) {
                return { taskId: 'offline', status: 'offline', progress: 0, currentTask: 'Offline' };
            }

            const status = tracker.getStatus();
            return {
                taskId: status.currentTask,
                currentTask: status.currentTask,
                status: status.status,
                progress: status.progress
            };
        }),
    indexingStatus: t.procedure.query(() => {
        const lspService = getLspService();
        if (!lspService) return { status: 'offline', filesIndexed: 0, totalFiles: 0 };
        return lspService.getStatus();
    }),
    executeTool: adminProcedure.input(z.object({
        name: z.string(),
        args: z.record(z.unknown())
    })).mutation(async ({ input }) => {
        const result = await getMcpServer().executeTool(input.name, input.args);
        if (result.isError) throw new Error(result.content[0].text);
        return result.content[0].text;
    }),
};
