import { z } from 'zod';
import { t, publicProcedure, getSessionManager } from '../lib/trpc-core.js';

export const sessionRouter = t.router({
    getState: publicProcedure.query(() => {
        return getSessionManager().getState();
    }),

    updateState: publicProcedure.input(z.object({
        isAutoDriveActive: z.boolean().optional(),
        activeGoal: z.string().nullable().optional(),
        lastObjective: z.string().nullable().optional(),
    })).mutation(({ input }) => {
        const manager = getSessionManager();
        manager.updateState(input);
        manager.save();
        return { success: true };
    }),

    clear: publicProcedure.mutation(() => {
        getSessionManager().clearSession();
        return { success: true };
    }),

    heartbeat: publicProcedure.mutation(() => {
        getSessionManager().touch();
        return { alive: true, timestamp: Date.now() };
    })
});
