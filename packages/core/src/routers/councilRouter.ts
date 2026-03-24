import { z } from 'zod';
import { t, publicProcedure, adminProcedure } from '../lib/trpc-core.js';
import { getCouncilOrchestrator, getCouncilService } from '../lib/trpc-core.js';
import { readFileSync, writeFileSync, existsSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

/**
 * councilRouter
 *
 * tRPC router for Borg's Council of Agents — the consensus-based advisory system.
 *
 * The Council evaluates Director actions through multi-agent debate.
 * Each member specializes in a domain (architecture, security, QA) and
 * provides opinions + votes on proposed actions.
 *
 * Exposes:
 * - `members` — Real agent definitions from council.json (Phase 63-15)
 * - `runSession` — Triggers a consensus debate on a proposal
 * - `getLatestSession` — Returns the most recent consensus result
 * - `listSessions` — Historical session list
 * - `getSession` — Detailed session by ID
 */

const opinionSchema = z.object({
    agentId: z.string(),
    content: z.string(),
    timestamp: z.number(),
    round: z.number(),
});

const voteSchema = z.object({
    agentId: z.string(),
    choice: z.string(),
    reason: z.string(),
    timestamp: z.number(),
});

const councilSessionSchema = z.object({
    id: z.string(),
    topic: z.string(),
    status: z.enum(['active', 'concluded']),
    round: z.number(),
    opinions: z.array(opinionSchema),
    votes: z.array(voteSchema),
    createdAt: z.number(),
});

/** Schema for a council member entry from council.json */
const memberSchema = z.object({
    name: z.string(),
    provider: z.string(),
    modelId: z.string(),
    systemPrompt: z.string(),
});

export const councilRouter = t.router({
    /**
     * members — Phase 63-15: Council Members Tab
     *
     * Returns the real agent definitions from config/council.json.
     * Each member has: name, provider, modelId, and systemPrompt.
     * Used by the dashboard to render the "Members" tab with live agent data.
     */
    members: publicProcedure
        .output(z.array(memberSchema))
        .query(() => {
            try {
                // Resolve council.json relative to this router file's package
                // Path: packages/core/config/council.json
                const configPath = join(
                    dirname(fileURLToPath(import.meta.url)),
                    '..', '..', 'config', 'council.json'
                );
                const raw = readFileSync(configPath, 'utf-8');
                const parsed = JSON.parse(raw);
                return parsed.members || [];
            } catch (e) {
                console.warn('[councilRouter] Failed to read council.json:', e);
                return [];
            }
        }),

    /**
     * updateMembers — Phase K1: Quota-Aware Routing UX
     * Safely flushes explicit provider drag-and-drop hierarchy straight into the runtime configuration.
     */
    updateMembers: adminProcedure
        .input(z.array(memberSchema))
        .mutation(async ({ input }) => {
            const configPath = join(
                dirname(fileURLToPath(import.meta.url)),
                '..', '..', 'config', 'council.json'
            );
            
            let config: any = { members: [] };
            try {
                if (existsSync(configPath)) {
                    config = JSON.parse(readFileSync(configPath, 'utf-8'));
                }
            } catch (e) {
                console.warn('[councilRouter] Could not parse existing council.json:', e);
            }
            
            config.members = input;
            writeFileSync(configPath, JSON.stringify(config, null, 2));
            return { success: true };
        }),

    /** Triggers a multi-agent consensus debate on a proposal */
    runSession: publicProcedure.input(z.object({ proposal: z.string() })).mutation(async ({ input }) => {
        return getCouncilOrchestrator().runConsensusSession(input.proposal);
    }),

    /** Returns the most recent consensus session result */
    getLatestSession: publicProcedure.query(async () => {
        return getCouncilOrchestrator().lastResult || null;
    }),

    /** Lists all historical council sessions */
    listSessions: publicProcedure.output(z.array(councilSessionSchema)).query(async () => {
        return getCouncilService().listSessions();
    }),

    /** Gets a specific council session by ID */
    getSession: publicProcedure.input(z.object({ id: z.string() })).output(councilSessionSchema.nullable()).query(async ({ input }) => {
        return getCouncilService().getSession(input.id) ?? null;
    })
});

