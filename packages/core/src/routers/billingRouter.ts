import { z } from 'zod';
import { t, publicProcedure, getLLMService } from '../lib/trpc-core.js';

interface BillingModelRuntime {
    id: string;
    provider: string;
    name?: string;
    inputPrice?: number | null;
    outputPrice?: number | null;
    contextWindow?: number | null;
    tier?: string;
    recommended?: boolean;
}

interface FallbackEntryRuntime {
    provider: string;
    model?: string;
    reason?: string;
}

export const billingRouter = t.router({
    /** Get current billing status — keys, usage, and provider breakdown */
    getStatus: publicProcedure.query(async () => {
        const keys = {
            openai: !!process.env.OPENAI_API_KEY,
            anthropic: !!process.env.ANTHROPIC_API_KEY,
            gemini: !!process.env.GEMINI_API_KEY,
            mistral: !!process.env.MISTRAL_API_KEY,
            deepseek: !!process.env.DEEPSEEK_API_KEY,
            xai: !!process.env.XAI_API_KEY,
            openrouter: !!process.env.OPENROUTER_API_KEY,
            groq: !!process.env.GROQ_API_KEY,
        };

        const llm = getLLMService();
        const stats = llm.getCostStats();
        const quota = llm.modelSelector.getQuotaService();

        let breakdown = quota.getUsageByModel();
        if (breakdown.length === 0) {
            breakdown = [{ provider: 'No Usage Yet', cost: 0, requests: 0 }];
        }

        const usage = {
            currentMonth: stats.estimatedCostUSD,
            limit: 100.00,
            breakdown,
        };

        return { keys, usage };
    }),

    /** Get quota information per provider — limits, remaining, reset dates */
    getProviderQuotas: publicProcedure.query(async () => {
        const llm = getLLMService();
        const quota = llm.modelSelector.getQuotaService();

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
            const quotaInfo = quota.getQuota?.(p.id);
            return {
                provider: p.id,
                name: p.name,
                configured: !!process.env[`${p.id.toUpperCase()}_API_KEY`],
                tier: quotaInfo?.tier ?? 'unknown',
                limit: quotaInfo?.limit ?? null,
                used: quotaInfo?.used ?? 0,
                remaining: quotaInfo?.remaining ?? null,
                resetDate: quotaInfo?.resetDate ?? null,
                rateLimitRpm: quotaInfo?.rateLimitRpm ?? null,
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
        const selector = getLLMService().modelSelector;

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
    getFallbackChain: publicProcedure.query(async () => {
        const selector = getLLMService().modelSelector;
        const chain = (selector.getFallbackChain?.() ?? []) as FallbackEntryRuntime[];
        return {
            chain: chain.map((entry, index: number) => ({
                priority: index + 1,
                provider: entry.provider,
                model: entry.model,
                reason: entry.reason ?? 'configured',
            })),
        };
    }),
});
