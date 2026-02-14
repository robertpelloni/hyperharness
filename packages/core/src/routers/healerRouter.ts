
import { z } from 'zod';
import { t, publicProcedure, getHealerService } from '../lib/trpc-core.js';
import { observable } from '@trpc/server/observable';

export const healerRouter: any = t.router({
    diagnose: t.procedure.input(z.object({ error: z.string(), context: z.string().optional() })).mutation(async ({ input }) => {
        return getHealerService().analyzeError(input.error, input.context || "");
    }),
    heal: t.procedure.input(z.object({ error: z.string(), context: z.string().optional() })).mutation(async ({ input }) => {
        const success = await getHealerService().heal(input.error, input.context || "");
        return { success };
    }),
    getHistory: t.procedure.query(async () => {
        return getHealerService().getHistory();
    }),
    subscribe: publicProcedure.subscription(() => {
        return observable<unknown>((emit) => {
            const onHeal = (data: unknown) => {
                emit.next(data);
            };
            const service = getHealerService();
            service.on('heal', onHeal);
            return () => {
                service.off('heal', onHeal);
            };
        });
    })
});
