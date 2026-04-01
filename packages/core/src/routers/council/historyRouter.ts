import { z } from 'zod';
import { t, publicProcedure, adminProcedure } from '../../lib/trpc-core.js';
import { debateHistory } from '../../orchestrator/council/services/debate-history.js';

const debateQueryOptionsSchema = z.object({
  sessionId: z.string().optional(),
  taskType: z.string().optional(),
  approved: z.boolean().optional(),
  supervisorName: z.string().optional(),
  fromTimestamp: z.number().optional(),
  toTimestamp: z.number().optional(),
  minConsensus: z.number().optional(),
  maxConsensus: z.number().optional(),
  limit: z.number().optional(),
  offset: z.number().optional(),
  sortBy: z.enum(['timestamp', 'consensus', 'duration']).optional(),
  sortOrder: z.enum(['asc', 'desc']).optional(),
});

export const historyRouter = t.router({
  status: publicProcedure.query(async () => {
    return {
      enabled: debateHistory.isEnabled(),
      recordCount: await debateHistory.getRecordCount(),
      storageSize: debateHistory.getStorageSize(),
      config: debateHistory.getConfig(),
    };
  }),

  getConfig: publicProcedure.query(async () => {
    return debateHistory.getConfig();
  }),

  updateConfig: adminProcedure.input(z.any()).mutation(async ({ input }) => {
    return debateHistory.updateConfig(input);
  }),

  toggle: adminProcedure.input(z.object({ enabled: z.boolean().optional() })).mutation(async ({ input }) => {
    const enabled = input.enabled ?? !debateHistory.isEnabled();
    debateHistory.updateConfig({ enabled });
    return { enabled: debateHistory.isEnabled() };
  }),

  stats: publicProcedure.query(async () => {
    return debateHistory.getStats();
  }),

  list: publicProcedure.input(debateQueryOptionsSchema).query(async ({ input }) => {
    const records = await debateHistory.queryDebates(input as any);
    return {
      records,
      meta: {
        count: records.length,
        totalRecords: await debateHistory.getRecordCount(),
      },
    };
  }),

  get: publicProcedure.input(z.object({ id: z.string() })).query(async ({ input }) => {
    const record = await debateHistory.getDebate(input.id);
    if (!record) throw new Error('Debate not found');
    return record;
  }),

  delete: adminProcedure.input(z.object({ id: z.string() })).mutation(async ({ input }) => {
    const deleted = await debateHistory.deleteRecord(input.id);
    if (!deleted) throw new Error('Debate not found');
    return { deleted: true, id: input.id };
  }),

  supervisorHistory: publicProcedure.input(z.object({ name: z.string() })).query(async ({ input }) => {
    return debateHistory.getSupervisorVoteHistory(input.name);
  }),

  clear: adminProcedure.mutation(async () => {
    const count = await debateHistory.clearAll();
    return { cleared: count };
  }),

  initialize: adminProcedure.mutation(async () => {
    await debateHistory.initialize();
    return {
      initialized: true,
      recordCount: await debateHistory.getRecordCount(),
    };
  }),
});
