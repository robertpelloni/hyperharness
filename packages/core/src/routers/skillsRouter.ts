
import { z } from 'zod';
import { t, publicProcedure } from '../lib/trpc-core.js';

export const skillsRouter = t.router({
    list: publicProcedure.query(async () => {
        // @ts-ignore
        const mcpServer = global.mcpServerInstance;
        if (!mcpServer || !mcpServer.skillRegistry) return [];

        return mcpServer.skillRegistry.getSkills();
    }),

    read: publicProcedure.input(z.object({
        name: z.string()
    })).query(async ({ input }) => {
        // @ts-ignore
        const mcpServer = global.mcpServerInstance;
        if (!mcpServer || !mcpServer.skillRegistry) {
            return { content: [{ type: "text", text: "Skill registry not available" }] };
        }
        return mcpServer.skillRegistry.readSkill(input.name);
    }),

    create: publicProcedure.input(z.object({
        id: z.string(),
        name: z.string(),
        description: z.string()
    })).mutation(async ({ input }) => {
        // @ts-ignore
        const mcpServer = global.mcpServerInstance;
        if (!mcpServer || !mcpServer.skillRegistry) {
            return { content: [{ type: "text", text: "Skill registry not available" }] };
        }
        return mcpServer.skillRegistry.createSkill(input.id, input.name, input.description);
    }),

    save: publicProcedure.input(z.object({
        id: z.string(),
        content: z.string()
    })).mutation(async ({ input }) => {
        // @ts-ignore
        const mcpServer = global.mcpServerInstance;
        if (!mcpServer || !mcpServer.skillRegistry) {
            return { content: [{ type: "text", text: "Skill registry not available" }] };
        }
        return mcpServer.skillRegistry.saveSkill(input.id, input.content);
    }),

    assimilate: publicProcedure.input(z.object({
        topic: z.string(),
        docsUrl: z.string().optional()
    })).mutation(async ({ input }) => {
        // @ts-ignore
        const mcpServer = global.mcpServerInstance;
        if (!mcpServer || !mcpServer.skillAssimilationService) {
            return { success: false, logs: ["Service not ready"] };
        }

        return await mcpServer.skillAssimilationService.assimilate({
            topic: input.topic,
            docsUrl: input.docsUrl,
            autoInstall: true
        });
    }),
});
