import { z } from 'zod';
import { t, publicProcedure, adminProcedure } from '../lib/trpc-core.js';
import { getMcpServer } from '../lib/mcpHelper.js';

function getWorkflowEngine() {
    return (getMcpServer() as any)?.workflowEngine ?? null;
}

export const workflowRouter = t.router({
    // --- Workflow Definitions ---

    list: publicProcedure.query(() => {
        const engine = getWorkflowEngine();
        if (!engine) return [];
        // Access registered workflows from the engine's internal Map
        const workflowsMap = (engine as any).workflows as Map<string, any> | undefined;
        if (!workflowsMap) return [];
        return Array.from(workflowsMap.values()).map((w: any) => ({
            id: w.id,
            name: w.name ?? w.id,
            description: w.description ?? '',
            entryPoint: w.entryPoint,
            nodeCount: w.nodes?.size ?? 0,
            edgeCount: w.edges?.length ?? 0,
        }));
    }),

    getGraph: publicProcedure
        .input(z.object({ workflowId: z.string() }))
        .query(({ input }) => {
            const engine = getWorkflowEngine();
            if (!engine) throw new Error("Workflow Engine not initialized");
            return engine.getGraph(input.workflowId) || { nodes: [], edges: [] };
        }),

    // --- Executions ---

    start: adminProcedure
        .input(z.object({
            workflowId: z.string(),
            initialState: z.record(z.any()).optional()
        }))
        .mutation(async ({ input }) => {
            const engine = getWorkflowEngine();
            if (!engine) throw new Error("Workflow Engine not initialized");
            const execution = await engine.start(input.workflowId, input.initialState || {});
            return execution;
        }),

    listExecutions: publicProcedure.query(() => {
        const engine = getWorkflowEngine();
        if (!engine) return [];
        return engine.listExecutions();
    }),

    getExecution: publicProcedure
        .input(z.object({ executionId: z.string() }))
        .query(({ input }) => {
            const engine = getWorkflowEngine();
            if (!engine) return null;
            return engine.getExecution(input.executionId);
        }),

    getHistory: publicProcedure
        .input(z.object({ executionId: z.string() }))
        .query(({ input }) => {
            const engine = getWorkflowEngine();
            if (!engine) return [];
            return engine.getHistory(input.executionId);
        }),

    // --- Control Flow ---

    resume: adminProcedure
        .input(z.object({ executionId: z.string() }))
        .mutation(async ({ input }) => {
            const engine = getWorkflowEngine();
            if (!engine) throw new Error("Workflow Engine not initialized");
            await engine.resume(input.executionId);
            return { success: true };
        }),

    pause: adminProcedure
        .input(z.object({ executionId: z.string() }))
        .mutation(({ input }) => {
            const engine = getWorkflowEngine();
            if (!engine) throw new Error("Workflow Engine not initialized");
            engine.pause(input.executionId);
            return { success: true };
        }),

    approve: adminProcedure
        .input(z.object({ executionId: z.string() }))
        .mutation(async ({ input }) => {
            const engine = getWorkflowEngine();
            if (!engine) throw new Error("Workflow Engine not initialized");
            await engine.approve(input.executionId);
            return { success: true };
        }),

    reject: adminProcedure
        .input(z.object({
            executionId: z.string(),
            reason: z.string().optional()
        }))
        .mutation(({ input }) => {
            const engine = getWorkflowEngine();
            if (!engine) throw new Error("Workflow Engine not initialized");
            engine.reject(input.executionId, input.reason);
            return { success: true };
        })
});

