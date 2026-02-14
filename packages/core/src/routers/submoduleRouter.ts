import { z } from 'zod';
import { t, publicProcedure, adminProcedure, getSubmoduleService } from '../lib/trpc-core.js';

export const submoduleRouter = t.router({
    list: publicProcedure.query(async () => {
        const service = getSubmoduleService();
        if (service) return await service.listSubmodules();
        return [];
    }),
    updateAll: adminProcedure.mutation(async () => {
        const service = getSubmoduleService();
        if (service) return await service.updateAll();
        return { success: false, output: "SubmoduleService not available" };
    }),
    installDependencies: adminProcedure.input(z.object({
        path: z.string()
    })).mutation(async ({ input }) => {
        const service = getSubmoduleService();
        if (service) return await service.installDependencies(input.path);
        return { success: false, output: "SubmoduleService not available" };
    }),
    build: adminProcedure.input(z.object({
        path: z.string()
    })).mutation(async ({ input }) => {
        const service = getSubmoduleService();
        if (service) return await service.buildSubmodule(input.path);
        return { success: false, output: "SubmoduleService not available" };
    }),
    enable: adminProcedure.input(z.object({
        path: z.string()
    })).mutation(async ({ input }) => {
        const service = getSubmoduleService();
        if (service) return await service.enableSubmodule(input.path);
        return { success: false, output: "SubmoduleService not available" };
    }),
    detectCapabilities: publicProcedure.input(z.object({
        path: z.string()
    })).query(({ input }) => {
        const service = getSubmoduleService();
        if (service) return service.detectCapabilities(input.path);
        return { caps: [], startCommand: undefined };
    })
});
