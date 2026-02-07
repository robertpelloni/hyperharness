import { z } from 'zod';
import { t, publicProcedure } from '../lib/trpc-core.js';
import { KnowledgeService } from '../services/KnowledgeService.js';

export const knowledgeRouter = t.router({
    getGraph: publicProcedure
        .input(z.object({
            query: z.string().optional(),
            depth: z.number().optional().default(1)
        }))
        .query(async ({ input }) => {
            // @ts-ignore
            const mcp = global.mcpServerInstance;
            if (!mcp || !mcp.knowledgeService) {
                return { nodes: [], edges: [] };
            }

            // Get raw content from service
            const result = await mcp.knowledgeService.getGraph(input.query, input.depth);

            try {
                // Parse the JSON string returned by the service
                const content = result.content[0].text;
                const graphData = JSON.parse(content);
                return graphData;
            } catch (e) {
                console.error("Failed to parse graph data", e);
                return { nodes: [], edges: [] };
            }
        }),

    getStats: publicProcedure.query(async () => {
        // @ts-ignore
        const mcp = global.mcpServerInstance;
        if (!mcp || !mcp.memoryManager) return { count: 0 };

        const contexts = await mcp.memoryManager.listContexts();
        return { count: contexts.length };
    }),
    getResources: publicProcedure.query(async () => {
        const fs = await import('fs/promises');
        const path = await import('path');
        const resourcePath = path.join(process.cwd(), 'knowledge', 'resources.json');

        try {
            const content = await fs.readFile(resourcePath, 'utf-8');
            return JSON.parse(content);
        } catch (e) {
            console.warn("Resources file not found or invalid", e);
            return { lastUpdated: "Never", categories: [] };
        }
    })
});
