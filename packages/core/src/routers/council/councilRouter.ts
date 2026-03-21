import { z } from 'zod';
import { t, publicProcedure, adminProcedure } from '../../lib/trpc-core.js';
import { getSupervisorCouncil, getCouncilWsManager, getCouncilHierarchy } from '../../lib/trpc-core.js';
import { createSupervisor, createSupervisors, createMockSupervisor } from '../../orchestrator/council/supervisors/index.js';

const supervisorConfigSchema = z.object({
  name: z.string(),
  provider: z.enum(['openai', 'anthropic', 'google', 'xai', 'moonshot', 'deepseek', 'qwen', 'custom', 'gemini', 'grok', 'kimi']),
  apiKey: z.string().optional(),
  model: z.string().optional(),
  temperature: z.number().optional(),
  baseURL: z.string().optional(),
  systemPrompt: z.string().optional(),
  weight: z.number().optional(),
});

const configSchema = z.object({
  supervisors: z.array(supervisorConfigSchema).optional(),
  debateRounds: z.number().int().min(1).max(10).optional(),
  consensusThreshold: z.number().min(0).max(1).optional(),
  enabled: z.boolean().optional(),
  smartPilot: z.boolean().optional(),
  weightedVoting: z.boolean().optional(),
  consensusMode: z.enum(['simple-majority', 'supermajority', 'unanimous', 'weighted', 'ceo-override', 'ceo-veto', 'hybrid-ceo-majority', 'ranked-choice']).optional(),
  leadSupervisor: z.string().optional(),
  fallbackSupervisors: z.array(z.string()).optional(),
});

const developmentTaskSchema = z.object({
  id: z.string(),
  description: z.string(),
  context: z.string(),
  files: z.array(z.string()),
});

export const councilRouter = t.router({
  status: publicProcedure.query(async () => {
    const instance = await getSupervisorCouncil();
    const available = await instance.getAvailableSupervisors();
    const hierarchy = (await getCouncilHierarchy()).getAllSpecializedCouncils().map(s => ({
      id: s.id,
      name: s.name,
      specialties: s.specialties,
      supervisorCount: s.council.getSupervisors().length
    }));
    
    const config = instance.getConfig();
    
    return {
      enabled: config.enabled ?? true,
      supervisorCount: instance.getSupervisors().length,
      availableCount: available.length,
      config,
      hierarchy,
    };
  }),

  updateConfig: adminProcedure.input(configSchema).mutation(async ({ input }) => {
    const instance = await getSupervisorCouncil();
    
    if (input.consensusMode) instance.setConsensusMode(input.consensusMode);
    if (input.leadSupervisor) instance.setLeadSupervisor(input.leadSupervisor);
    if (input.fallbackSupervisors) instance.setFallbackChain(input.fallbackSupervisors);
    if (input.debateRounds) instance.setDebateRounds(input.debateRounds);
    if (input.consensusThreshold) instance.setConsensusThreshold(input.consensusThreshold);
    
    const ws = await getCouncilWsManager();
    ws.broadcast({
      type: 'log',
      payload: { level: 'info', message: `Council config updated`, timestamp: Date.now() },
      timestamp: Date.now(),
    });
    
    return instance.getConfig();
  }),

  addSupervisors: adminProcedure.input(z.object({
    supervisors: z.array(supervisorConfigSchema)
  })).mutation(async ({ input }) => {
    const instance = await getSupervisorCouncil();
    const added: string[] = [];
    const failed: string[] = [];

    for (const supervisorConfig of input.supervisors) {
      try {
        const supervisor = createSupervisor(supervisorConfig as any);
        instance.addSupervisor(supervisor);
        added.push(supervisorConfig.name);
      } catch (error) {
        const msg = error instanceof Error ? error.message : 'Unknown error';
        failed.push(`${supervisorConfig.name}: ${msg}`);
      }
    }

    const ws = await getCouncilWsManager();
    ws.broadcast({
      type: 'log',
      payload: { level: 'info', message: `Added ${added.length} supervisors: ${added.join(', ')}`, timestamp: Date.now() },
      timestamp: Date.now(),
    });

    return { added, failed };
  }),

  clearSupervisors: adminProcedure.mutation(async () => {
    const instance = await getSupervisorCouncil();
    instance.clearSupervisors();
    
    const ws = await getCouncilWsManager();
    ws.broadcast({
      type: 'log',
      payload: { level: 'info', message: 'All supervisors removed', timestamp: Date.now() },
      timestamp: Date.now(),
    });
    
    return { success: true };
  }),

  debate: adminProcedure.input(developmentTaskSchema).mutation(async ({ input }) => {
    const instance = await getSupervisorCouncil();
    const ws = await getCouncilWsManager();
    
    ws.broadcast({
      type: 'log',
      payload: { level: 'info', message: `Debate started: ${input.description}`, timestamp: Date.now(), source: 'council' },
      timestamp: Date.now(),
    });
    
    try {
      const decision = await instance.debate(input);
      
      ws.broadcast({
        type: 'council_decision',
        payload: decision,
        timestamp: Date.now(),
      });
      
      return decision;
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unknown error';
      ws.notifyError(`Debate failed: ${message}`);
      throw new Error(`Debate failed: ${message}`);
    }
  }),

  toggle: adminProcedure.mutation(async () => {
    const instance = await getSupervisorCouncil();
    const config = instance.getConfig();
    const newEnabled = !config.enabled;
    // Note: Need to add setEnabled to SupervisorCouncil if not present, 
    // or just update via config. For now, assuming direct config access or adding a method.
    // In migrated code, it was a local variable in the route.
    // I'll assume we add a setEnabled method to the instance for parity.
    (instance as any).config.enabled = newEnabled; 
    
    const ws = await getCouncilWsManager();
    ws.broadcast({
      type: 'log',
      payload: { level: 'info', message: `Council ${newEnabled ? 'enabled' : 'disabled'}`, timestamp: Date.now() },
      timestamp: Date.now(),
    });
    
    return { enabled: newEnabled };
  }),

  addMock: adminProcedure.mutation(async () => {
    const instance = await getSupervisorCouncil();
    const mockName = `MockSupervisor-${Date.now()}`;
    instance.addSupervisor(createMockSupervisor(mockName));
    
    const ws = await getCouncilWsManager();
    ws.broadcast({
      type: 'log',
      payload: { level: 'info', message: `Mock supervisor added: ${mockName}`, timestamp: Date.now() },
      timestamp: Date.now(),
    });
    
    return { added: mockName };
  }),
});
