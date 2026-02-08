import { z } from 'zod';
import { t, publicProcedure, getMcpServer } from '../lib/trpc-core.js';

export const sessionRouter = t.router({
    getState: publicProcedure.query(() => {
        const mcp = getMcpServer();
        if (mcp && mcp.sessionManager) {
            return mcp.sessionManager.getState();
        }
        return {
            isAutoDriveActive: false,
            activeGoal: null,
            lastObjective: null,
            lastHeartbeat: 0
        };
    }),

    updateState: publicProcedure.input(z.object({
        isAutoDriveActive: z.boolean().optional(),
        activeGoal: z.string().nullable().optional(),
        lastObjective: z.string().nullable().optional(),
    })).mutation(({ input }) => {
        const mcp = getMcpServer();
        if (mcp && mcp.sessionManager) {
            mcp.sessionManager.updateState(input);
            mcp.sessionManager.save();
            return { success: true };
        }
        throw new Error("SessionManager not ready");
    }),

    clear: publicProcedure.mutation(() => {
        const mcp = getMcpServer();
        if (mcp && mcp.sessionManager) {
            mcp.sessionManager.clearSession();
            return { success: true };
        }
        throw new Error("SessionManager not ready");
    }),

    heartbeat: publicProcedure.mutation(() => {
        const mcp = getMcpServer();
        if (mcp && mcp.sessionManager) {
            mcp.sessionManager.touch();
            return { alive: true, timestamp: Date.now() };
        }
        throw new Error("SessionManager not ready");
    })
});
