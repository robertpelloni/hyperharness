import { z } from 'zod';
import { t, publicProcedure } from '../lib/trpc-core.js';
import { getCouncilOrchestrator, getCouncilService } from '../lib/trpc-core.js';

export const councilRouter = t.router({
    runSession: publicProcedure.input(z.object({ proposal: z.string() })).mutation(async ({ input }) => {
        return getCouncilOrchestrator().runConsensusSession(input.proposal);
    }),
    getLatestSession: publicProcedure.query(async () => {
        return getCouncilOrchestrator().lastResult || null;
    }),
    listSessions: publicProcedure.query(async () => {
        return getCouncilService().listSessions();
    }),
    getSession: publicProcedure.input(z.object({ id: z.string() })).query(async ({ input }) => {
        return getCouncilService().getSession(input.id);
    })
});
