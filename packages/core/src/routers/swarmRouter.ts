/**
 * swarmRouter.ts
 * 
 * Exposes the Multi-Agent Orchestration features (Swarm, Debate, Consensus) to the Next.js UI.
 */

import { z } from 'zod';
import { t, publicProcedure } from '../trpc.js';
import { SwarmOrchestrator } from '../agents/swarm/SwarmOrchestrator.js';
import { DebateProtocol } from '../agents/swarm/DebateProtocol.js';
import { ConsensusEngine } from '../agents/swarm/ConsensusEngine.js';

export const swarmRouter = t.router({

    startSwarm: publicProcedure
        .input(
            z.object({
                masterPrompt: z.string(),
                maxConcurrency: z.number().optional()
            })
        )
        .mutation(async ({ input }: { input: { masterPrompt: string; maxConcurrency?: number } }) => {
            const orchestrator = new SwarmOrchestrator({ maxConcurrency: input.maxConcurrency || 5 });
            const tasks = await orchestrator.decomposeGoal(input.masterPrompt);

            // We don't await the execution so TRPC can return immediately,
            // actual implementation would publish events.
            // E.g: orchestrator.executeSwarm();

            return { status: 'initialized', tasksCount: tasks.length };
        }),

    executeDebate: publicProcedure
        .input(
            z.object({
                topic: z.string(),
                proponentModel: z.string(),
                opponentModel: z.string(),
                judgeModel: z.string(),
                rounds: z.number().optional()
            })
        )
        .mutation(async ({ input }: { input: { topic: string; proponentModel: string; opponentModel: string; judgeModel: string; rounds?: number } }) => {
            const debate = new DebateProtocol({
                topic: input.topic,
                proponentModel: input.proponentModel,
                opponentModel: input.opponentModel,
                judgeModel: input.judgeModel,
                maxRounds: input.rounds || 3
            });

            // Synchronous wait for simulation, standard TRPC streaming is better for prod
            const result = await debate.conductDebate();
            return result;
        }),

    seekConsensus: publicProcedure
        .input(
            z.object({
                prompt: z.string(),
                models: z.array(z.string()).min(2),
                requiredAgreement: z.number().optional()
            })
        )
        .mutation(async ({ input }: { input: { prompt: string; models: string[]; requiredAgreement?: number } }) => {
            const engine = new ConsensusEngine();

            const result = await engine.seekConsensus({
                prompt: input.prompt,
                models: input.models,
                requiredAgreementPercentage: input.requiredAgreement
            });

            return result;
        }),
});
