import { TRPCError } from '@trpc/server';
import { z } from 'zod';
import { adminProcedure, t, publicProcedure, getLLMService } from '../lib/trpc-core.js';
import type { ProviderRoutingStrategy, ProviderTaskType } from '../providers/types.js';

interface ProviderQuotaRuntime {
    provider: string;
    name: string;
    configured: boolean;
    authenticated: boolean;
    authMethod: string;
    /** @see ProviderAuthTruth */
    authTruth?: string;
    tier: string;
    limit: number | null;
    used: number;
    remaining: number | null;
    resetDate: string | null;
    rateLimitRpm: number | null;
    availability: string;
    lastError?: string;
    windows?: Array<{
        key: string;
        label: string;
        used: number;
        limit: number | null;
        remaining: number | null;
        resetDate: string | null;
        unit: string;
    }>;
    source?: string;
    connectionId?: string | null;
    /** @see QuotaDataConfidence */
    quotaConfidence?: string;
    /** ISO-8601 timestamp of last real-time fetch from provider, or null. */
    quotaRefreshedAt?: string | null;
            }));
        }

        // Build quota info for each configured provider
        const providers = [
            { id: 'openai', name: 'OpenAI' },
            { id: 'anthropic', name: 'Anthropic' },
            { id: 'gemini', name: 'Google Gemini' },
            { id: 'deepseek', name: 'DeepSeek' },
            { id: 'xai', name: 'xAI (Grok)' },
            { id: 'mistral', name: 'Mistral' },
            { id: 'openrouter', name: 'OpenRouter' },
            { id: 'groq', name: 'Groq' },
        ];

        return providers.map(p => {
            const hasKey = !!process.env[`${p.id.toUpperCase()}_API_KEY`];
            const quotaInfo = quota.getQuota?.(p.id);
            return {
                provider: p.id,
                name: p.name,
                configured: hasKey,
                authenticated: hasKey,
                authMethod: hasKey ? 'api_key' : 'none',
                authTruth: quotaInfo?.authTruth ?? (hasKey ? 'authenticated' : 'not_configured'),
                tier: quotaInfo?.tier ?? 'unknown',
                limit: quotaInfo?.limit ?? null,
                used: quotaInfo?.used ?? 0,
                remaining: quotaInfo?.remaining ?? null,
                resetDate: quotaInfo?.resetDate ?? null,
                rateLimitRpm: quotaInfo?.rateLimitRpm ?? null,
                availability: hasKey ? 'available' : 'missing_config',
                lastError: null,
                quotaConfidence: quotaInfo?.quotaConfidence ?? 'estimated',
                quotaRefreshedAt: quotaInfo?.quotaRefreshedAt ?? null,
            };
        });
    }),

    /** Get cost history over time for charts */
    getCostHistory: publicProcedure.input(z.object({
        days: z.number().min(1).max(90).default(30),
    }).optional()).query(async ({ input }) => {
        const llm = getLLMService();
        const stats = llm.getCostStats();

        // Build daily breakdown from metrics if available
        const days = input?.days ?? 30;
        const history: { date: string; cost: number; requests: number }[] = [];
        const statsAny = stats;

        // If there's a cost history, use it; otherwise generate from current stats
        if (statsAny.dailyHistory) {
            return { history: statsAny.dailyHistory.slice(-days) };
        }

        // Fallback: single entry for today
        const today = new Date().toISOString().split('T')[0];
        history.push({
            date: today,
            cost: stats.estimatedCostUSD ?? 0,
            requests: statsAny.totalRequests ?? 0,
        });

        return { history };
    }),

    /** Get model-level pricing and efficiency data */
    getModelPricing: publicProcedure.query(async () => {
        const llm = getLLMService();
        const selector = llm.modelSelector as typeof llm.modelSelector & BillingSelectorRuntime;

        // Get available models with their pricing info
        const models = (selector.getAvailableModels?.() ?? []) as BillingModelRuntime[];
        return {
            models: models.map((m) => ({
                id: m.id,
                provider: m.provider,
                name: m.name ?? m.id,
                inputPricePer1k: m.inputPrice ?? null,
                outputPricePer1k: m.outputPrice ?? null,
                contextWindow: m.contextWindow ?? null,
                tier: m.tier ?? 'standard',
                recommended: m.recommended ?? false,
            })),
        };
    }),

    /** Get current fallback chain configuration */
    getFallbackChain: publicProcedure
        .input(z.object({
            taskType: z.enum(TASK_TYPE_VALUES).optional(),
        }).optional())
        .query(async ({ input }) => {
        const llm = getLLMService();
        const selector = llm.modelSelector as typeof llm.modelSelector & BillingSelectorRuntime;
        return buildFallbackChainResponse(selector, input?.taskType);
    }),

    /** Get task-specific routing rules and top-ranked fallback previews */
    getTaskRoutingRules: publicProcedure.query(async () => {
        const llm = getLLMService();
        const selector = llm.modelSelector as typeof llm.modelSelector & BillingSelectorRuntime;
        const taskRules = selector.getTaskRoutingRules?.() ?? {
            coding: 'cheapest',
            planning: 'best',
            research: 'best',
            general: 'round-robin',
            worker: 'cheapest',
            supervisor: 'best',
        } satisfies Record<ProviderTaskType, ProviderRoutingStrategy>;

        const defaultStrategy = selector.getRoutingStrategy?.() ?? 'best';

        return {
            defaultStrategy,
            rules: TASK_TYPES.map((taskType) => ({
                taskType,
                strategy: taskRules[taskType] ?? defaultStrategy,
                fallbackPreview: (selector.getFallbackChain?.({ routingTaskType: taskType }) ?? [])
                    .slice(0, 3)
                    .map((entry) => ({
                        provider: entry.provider,
                        model: entry.model,
                        reason: entry.reason ?? 'configured',
                    })),
            })),
        };
    }),

    setRoutingStrategy: adminProcedure
        .input(z.object({
            strategy: z.enum(['cheapest', 'best', 'round-robin']),
        }))
        .mutation(async ({ input }) => {
            const llm = getLLMService();
            const selector = llm.modelSelector as typeof llm.modelSelector & BillingSelectorRuntime;

            selector.setRoutingStrategy?.(input.strategy);

            return {
                ok: true,
                strategy: selector.getRoutingStrategy?.() ?? input.strategy,
            };
        }),

    setTaskRoutingRule: adminProcedure
        .input(z.object({
            taskType: z.enum(TASK_TYPE_VALUES),
            strategy: z.enum(['cheapest', 'best', 'round-robin']).nullable(),
        }))
        .mutation(async ({ input }) => {
            const llm = getLLMService();
            const selector = llm.modelSelector as typeof llm.modelSelector & BillingSelectorRuntime;

            if (input.strategy) {
                selector.setTaskRoutingStrategy?.(input.taskType, input.strategy);
            } else {
                selector.resetTaskRoutingStrategy?.(input.taskType);
            }

            const rules = selector.getTaskRoutingRules?.() ?? {
                coding: 'cheapest',
                planning: 'best',
                research: 'best',
                general: 'round-robin',
                worker: 'cheapest',
                supervisor: 'best',
            } satisfies Record<ProviderTaskType, ProviderRoutingStrategy>;

            return {
                ok: true,
                taskType: input.taskType,
                strategy: rules[input.taskType],
            };
        }),

    /** Get models currently depleted or on cooldown in the ModelSelector session state. */
    getDepletedModels: publicProcedure.query(async () => {
        const llm = getLLMService();
        try {
            return llm.modelSelector.getDepletedModels?.() ?? [];
        } catch (error) {
            throw buildBillingUnavailableError('Depleted model status', error);
        }
    }),

    /**
     * Get recent provider fallback decisions from the in-process ring buffer.
     *
     * Each entry represents a routing decision where the preferred provider was not
     * honored (due to quota exhaustion, rate limits, budget cap, or no candidates
     * available at all).  Returns entries in reverse-chronological order.
     *
     * The billing dashboard "Recent Fallback Decisions" card consumes this to let
     * operators understand why HyperCode substituted a different provider/model than
     * configured.
     */
    getFallbackHistory: publicProcedure
        .input(z.object({ limit: z.number().min(1).max(50).optional() }).optional())
        .query(async ({ input }) => {
            const llm = getLLMService();
            try {
                const selector = llm.modelSelector as unknown as {
                    getFallbackHistory?: (limit: number) => Array<{
                        id: number;
                        timestamp: number;
                        requestedProvider?: string;
                        selectedProvider: string;
                        selectedModelId: string;
                        taskType: string;
                        strategy: string;
                        reason: string;
                        causeCode: string;
                    }>;
                };
                return selector.getFallbackHistory?.(input?.limit ?? 20) ?? [];
            } catch (error) {
                throw buildBillingUnavailableError('Fallback history', error);
            }
        }),

    /** Clear in-memory provider fallback history ring buffer. */
    clearFallbackHistory: adminProcedure
        .mutation(async () => {
            const llm = getLLMService();
            try {
                const selector = llm.modelSelector as unknown as {
                    clearFallbackHistory?: () => void;
                };
                selector.clearFallbackHistory?.();
                return { ok: true };
            } catch (error) {
                throw buildBillingUnavailableError('Fallback history clearing', error);
            }
        }),
});

