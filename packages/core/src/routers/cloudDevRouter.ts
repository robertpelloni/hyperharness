/**
 * cloudDevRouter.ts – tRPC router for managing cloud development environments
 *
 * Provides a unified interface for managing cloud dev agents across multiple
 * providers: Jules (Google), Codex (OpenAI), GitHub Copilot Workspace, Devin,
 * and any future cloud coding agents.
 *
 * Each provider is represented as a "Cloud Dev Provider" with standardised
 * session management: create, list, pause, resume, stop, and query sessions.
 *
 * The jules-autopilot submodule (external/jules-autopilot) powers the Jules
 * provider natively.  Other providers use their respective REST APIs.
 */

import { z } from 'zod';
import { t, publicProcedure } from '../lib/trpc-core.js';

// ---- Schema Definitions ----

const CloudDevProviderSchema = z.enum([
    'jules',
    'codex',
    'copilot-workspace',
    'devin',
    'custom',
]);

const CloudDevSessionStatusSchema = z.enum([
    'pending',
    'active',
    'paused',
    'completed',
    'failed',
    'awaiting_approval',
    'cancelled',
]);

const CloudDevSessionSchema = z.object({
    id: z.string(),
    provider: CloudDevProviderSchema,
    projectName: z.string(),
    task: z.string(),
    status: CloudDevSessionStatusSchema,
    createdAt: z.string(),
    updatedAt: z.string(),
    metadata: z.record(z.any()).optional(),
});

type CloudDevSession = z.infer<typeof CloudDevSessionSchema>;

// ---- In-memory session store (would be persisted via DB in production) ----

const sessions: CloudDevSession[] = [];

// ---- Provider Configuration ----

interface ProviderConfig {
    name: string;
    provider: z.infer<typeof CloudDevProviderSchema>;
    apiKeyEnvVar: string;
    baseUrl: string;
    enabled: boolean;
}

const providerConfigs: ProviderConfig[] = [
    {
        name: 'Jules (Google)',
        provider: 'jules',
        apiKeyEnvVar: 'JULES_API_KEY',
        baseUrl: 'https://jules.googleapis.com/v1alpha',
        enabled: true,
    },
    {
        name: 'Codex (OpenAI)',
        provider: 'codex',
        apiKeyEnvVar: 'OPENAI_API_KEY',
        baseUrl: 'https://api.openai.com/v1/codex',
        enabled: false,
    },
    {
        name: 'GitHub Copilot Workspace',
        provider: 'copilot-workspace',
        apiKeyEnvVar: 'GITHUB_TOKEN',
        baseUrl: 'https://api.github.com/copilot/workspace',
        enabled: false,
    },
    {
        name: 'Devin',
        provider: 'devin',
        apiKeyEnvVar: 'DEVIN_API_KEY',
        baseUrl: 'https://api.devin.ai/v1',
        enabled: false,
    },
];

// ---- Router ----

export const cloudDevRouter = t.router({
    /**
     * List all registered cloud dev providers and their status.
     */
    listProviders: publicProcedure.query(() => {
        return providerConfigs.map(p => ({
            ...p,
            hasApiKey: !!process.env[p.apiKeyEnvVar],
        }));
    }),

    /**
     * Create a new cloud dev session.
     */
    createSession: publicProcedure
        .input(
            z.object({
                provider: CloudDevProviderSchema,
                projectName: z.string(),
                task: z.string(),
                metadata: z.record(z.any()).optional(),
            })
        )
        .mutation(({ input }) => {
            const now = new Date().toISOString();
            const session: CloudDevSession = {
                id: `cds-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
                provider: input.provider,
                projectName: input.projectName,
                task: input.task,
                status: 'pending',
                createdAt: now,
                updatedAt: now,
                metadata: input.metadata,
            };
            sessions.push(session);
            return session;
        }),

    /**
     * List all cloud dev sessions, optionally filtered by provider or status.
     */
    listSessions: publicProcedure
        .input(
            z
                .object({
                    provider: CloudDevProviderSchema.optional(),
                    status: CloudDevSessionStatusSchema.optional(),
                })
                .optional()
        )
        .query(({ input }) => {
            let result = [...sessions];
            if (input?.provider) {
                result = result.filter(s => s.provider === input.provider);
            }
            if (input?.status) {
                result = result.filter(s => s.status === input.status);
            }
            // Sort newest first
            result.sort(
                (a, b) =>
                    new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime()
            );
            return result;
        }),

    /**
     * Get a single session by ID.
     */
    getSession: publicProcedure
        .input(z.object({ sessionId: z.string() }))
        .query(({ input }) => {
            return sessions.find(s => s.id === input.sessionId) ?? null;
        }),

    /**
     * Update a session's status (pause, resume, complete, cancel).
     */
    updateSessionStatus: publicProcedure
        .input(
            z.object({
                sessionId: z.string(),
                status: CloudDevSessionStatusSchema,
            })
        )
        .mutation(({ input }) => {
            const session = sessions.find(s => s.id === input.sessionId);
            if (!session) {
                throw new Error(`Session ${input.sessionId} not found`);
            }
            session.status = input.status;
            session.updatedAt = new Date().toISOString();
            return session;
        }),

    /**
     * Delete a session.
     */
    deleteSession: publicProcedure
        .input(z.object({ sessionId: z.string() }))
        .mutation(({ input }) => {
            const idx = sessions.findIndex(s => s.id === input.sessionId);
            if (idx === -1) {
                throw new Error(`Session ${input.sessionId} not found`);
            }
            sessions.splice(idx, 1);
            return { success: true };
        }),

    /**
     * Aggregate stats across all providers.
     */
    stats: publicProcedure.query(() => {
        const byProvider: Record<string, number> = {};
        const byStatus: Record<string, number> = {};

        for (const s of sessions) {
            byProvider[s.provider] = (byProvider[s.provider] || 0) + 1;
            byStatus[s.status] = (byStatus[s.status] || 0) + 1;
        }

        return {
            totalSessions: sessions.length,
            byProvider,
            byStatus,
            providers: providerConfigs.length,
            enabledProviders: providerConfigs.filter(p => p.enabled).length,
        };
    }),
});
