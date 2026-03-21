import { Hono } from 'hono';
import { zValidator } from '@hono/zod-validator';
import { z } from 'zod';
import { systemManager } from '../services/system-manager.js';
import { selfEvolution } from '../services/self-evolution.js';

const app = new Hono();

app.get('/submodules', async (c) => {
  const submodules = await systemManager.getSubmodules();
  return c.json({ success: true, data: submodules });
});

app.get('/version', async (c) => {
  const version = await systemManager.getProjectVersion();
  return c.json({ success: true, data: version });
});

app.post('/evolve', zValidator('json', z.object({
  description: z.string().min(10)
})), async (c) => {
  const { description } = c.req.valid('json');
  const sessionId = await selfEvolution.evolveSystem(description);
  return c.json({ success: true, sessionId, message: 'Self-evolution task initiated' });
});

app.post('/optimize-weights', async (c) => {
  selfEvolution.optimizeWeights();
  return c.json({ success: true, message: 'Weights optimized based on historical performance' });
});

export default app;
