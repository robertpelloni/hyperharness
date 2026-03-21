import { Hono } from 'hono';
import { debateHistory, type DebateQueryOptions } from '../services/debate-history.js';
import type { TaskType } from './types.js';

export const debateHistoryRoutes = new Hono();

debateHistoryRoutes.get('/status', (c) => {
  return c.json({
    success: true,
    data: {
      enabled: debateHistory.isEnabled(),
      recordCount: debateHistory.getRecordCount(),
      storageSize: debateHistory.getStorageSize(),
      config: debateHistory.getConfig(),
    },
  });
});

debateHistoryRoutes.get('/config', (c) => {
  return c.json({
    success: true,
    data: debateHistory.getConfig(),
  });
});

debateHistoryRoutes.post('/config', async (c) => {
  const body = await c.req.json();
  const config = debateHistory.updateConfig(body);
  return c.json({
    success: true,
    data: config,
  });
});

debateHistoryRoutes.post('/toggle', async (c) => {
  const body = await c.req.json<{ enabled?: boolean }>();
  const enabled = body.enabled ?? !debateHistory.isEnabled();
  debateHistory.updateConfig({ enabled });
  return c.json({
    success: true,
    data: { enabled: debateHistory.isEnabled() },
  });
});

debateHistoryRoutes.get('/stats', (c) => {
  return c.json({
    success: true,
    data: debateHistory.getStats(),
  });
});

debateHistoryRoutes.get('/list', (c) => {
  const query = c.req.query();
  
  const options: DebateQueryOptions = {
    sessionId: query.sessionId,
    taskType: query.taskType as TaskType | undefined,
    approved: query.approved === 'true' ? true : query.approved === 'false' ? false : undefined,
    supervisorName: query.supervisorName,
    fromTimestamp: query.fromTimestamp ? parseInt(query.fromTimestamp, 10) : undefined,
    toTimestamp: query.toTimestamp ? parseInt(query.toTimestamp, 10) : undefined,
    minConsensus: query.minConsensus ? parseFloat(query.minConsensus) : undefined,
    maxConsensus: query.maxConsensus ? parseFloat(query.maxConsensus) : undefined,
    limit: query.limit ? parseInt(query.limit, 10) : undefined,
    offset: query.offset ? parseInt(query.offset, 10) : undefined,
    sortBy: query.sortBy as 'timestamp' | 'consensus' | 'duration' | undefined,
    sortOrder: query.sortOrder as 'asc' | 'desc' | undefined,
  };

  const records = debateHistory.queryDebates(options);
  return c.json({
    success: true,
    data: records,
    meta: {
      count: records.length,
      totalRecords: debateHistory.getRecordCount(),
    },
  });
});

debateHistoryRoutes.get('/debates/:id', (c) => {
  const id = c.req.param('id');
  const record = debateHistory.getDebate(id);
  
  if (!record) {
    return c.json({ success: false, error: 'Debate not found' }, 404);
  }
  
  return c.json({
    success: true,
    data: record,
  });
});

debateHistoryRoutes.delete('/debates/:id', (c) => {
  const id = c.req.param('id');
  const deleted = debateHistory.deleteRecord(id);
  
  if (!deleted) {
    return c.json({ success: false, error: 'Debate not found' }, 404);
  }
  
  return c.json({
    success: true,
    data: { deleted: true, id },
  });
});

debateHistoryRoutes.get('/supervisor/:name', (c) => {
  const name = c.req.param('name');
  const history = debateHistory.getSupervisorVoteHistory(name);
  return c.json({
    success: true,
    data: history,
  });
});

debateHistoryRoutes.get('/export/json', (c) => {
  const query = c.req.query();
  
  const options: DebateQueryOptions = {
    sessionId: query.sessionId,
    approved: query.approved === 'true' ? true : query.approved === 'false' ? false : undefined,
    limit: query.limit ? parseInt(query.limit, 10) : undefined,
  };

  const json = debateHistory.exportToJson(options);
  
  c.header('Content-Type', 'application/json');
  c.header('Content-Disposition', 'attachment; filename="debate-history.json"');
  return c.body(json);
});

debateHistoryRoutes.get('/export/csv', (c) => {
  const query = c.req.query();
  
  const options: DebateQueryOptions = {
    sessionId: query.sessionId,
    approved: query.approved === 'true' ? true : query.approved === 'false' ? false : undefined,
    limit: query.limit ? parseInt(query.limit, 10) : undefined,
  };

  const csv = debateHistory.exportToCsv(options);
  
  c.header('Content-Type', 'text/csv');
  c.header('Content-Disposition', 'attachment; filename="debate-history.csv"');
  return c.body(csv);
});

debateHistoryRoutes.delete('/clear', (c) => {
  const count = debateHistory.clearAll();
  return c.json({
    success: true,
    data: { cleared: count },
  });
});

debateHistoryRoutes.post('/initialize', (c) => {
  debateHistory.initialize();
  return c.json({
    success: true,
    data: {
      initialized: true,
      recordCount: debateHistory.getRecordCount(),
    },
  });
});
