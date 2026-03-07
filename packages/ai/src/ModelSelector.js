import fs from 'fs';
import path from 'path';
import { QuotaService } from './services/QuotaService.js';

const DEFAULT_CHAINS = {
    worker: [
        { provider: 'google', modelId: 'gemini-2.0-flash' },
        { provider: 'anthropic', modelId: 'claude-sonnet-4-20250514' },
        { provider: 'openai', modelId: 'gpt-4o' },
        { provider: 'deepseek', modelId: 'deepseek-chat' },
        { provider: 'lmstudio', modelId: 'local' },
        { provider: 'ollama', modelId: 'gemma:2b' }
    ],
    supervisor: [
        { provider: 'openai', modelId: 'gpt-4o' },
        { provider: 'anthropic', modelId: 'claude-sonnet-4-20250514' },
        { provider: 'google', modelId: 'gemini-1.5-pro' },
        { provider: 'lmstudio', modelId: 'local' }
    ]
};

const COOL_DOWN_MS = 60 * 1000;

export class ModelSelector {
    modelStates = new Map();
    configPath;
    quotaService;

    constructor() {
        this.quotaService = new QuotaService();
        this.configPath = path.resolve(process.cwd(), 'packages/core/config/council.json');
        console.log('ModelSelector initialized. Config Path:', this.configPath);
    }

    getQuotaService() {
        return this.quotaService;
    }

    getCandidateKey(provider, modelId) {
        return `${provider}:${modelId}`;
    }

    reportFailure(provider, modelId) {
        const candidateKey = this.getCandidateKey(provider, modelId);
        console.warn(`[ModelSelector] Reporting failure for ${candidateKey}. Marking as DEPLETED.`);
        this.modelStates.set(candidateKey, {
            isDepleted: true,
            depletedAt: Date.now(),
            retryAfter: Date.now() + COOL_DOWN_MS
        });
    }

    loadConfig() {
        try {
            if (fs.existsSync(this.configPath)) {
                const raw = fs.readFileSync(this.configPath, 'utf-8');
                return JSON.parse(raw);
            }
        }
        catch (error) {
            console.error('Failed to load council config:', error);
        }

        return null;
    }

    getConfiguredChain() {
        const config = this.loadConfig();

        if (config && Array.isArray(config.members)) {
            return config.members.map((member) => ({
                provider: member.provider,
                modelId: member.modelId,
                systemPrompt: member.systemPrompt
            }));
        }

        return null;
    }

    buildChain(req) {
        const configuredChain = this.getConfiguredChain();

        let chain = configuredChain
            ? configuredChain
            : (req.taskType === 'supervisor' || req.taskComplexity === 'high'
                ? DEFAULT_CHAINS.supervisor
                : DEFAULT_CHAINS.worker);

        if (req.provider) {
            const preferred = chain.filter((candidate) => candidate.provider === req.provider);
            const remaining = chain.filter((candidate) => candidate.provider !== req.provider);
            chain = [...preferred, ...remaining];
        }

        return chain;
    }

    isCoolingDown(provider, modelId) {
        const candidateKey = this.getCandidateKey(provider, modelId);
        const status = this.modelStates.get(candidateKey);

        if (!status || !status.isDepleted) {
            return false;
        }

        if (Date.now() > (status.retryAfter || 0)) {
            this.modelStates.delete(candidateKey);
            return false;
        }

        return true;
    }

    hasKey(provider) {
        if (provider === 'google') return !!process.env.GOOGLE_API_KEY;
        if (provider === 'openai') return !!process.env.OPENAI_API_KEY;
        if (provider === 'anthropic') return !!process.env.ANTHROPIC_API_KEY;
        if (provider === 'deepseek') return !!process.env.DEEPSEEK_API_KEY;
        return true;
    }

    async selectModel(req) {
        const isBudgetExceeded = this.quotaService.isBudgetExceeded();

        if (isBudgetExceeded) {
            console.warn('[ModelSelector] Budget exceeded! Forcing FREE local models.');
            return {
                provider: 'lmstudio',
                modelId: 'local',
                reason: 'BUDGET_EXCEEDED_FORCED_LOCAL'
            };
        }

        const chain = this.buildChain(req);
        const excludedCandidates = new Set(req.exclude || []);

        for (const candidate of chain) {
            const candidateKey = this.getCandidateKey(candidate.provider, candidate.modelId);

            if (excludedCandidates.has(candidateKey)) {
                continue;
            }

            if (!this.hasKey(candidate.provider)) {
                continue;
            }

            if (this.isCoolingDown(candidate.provider, candidate.modelId)) {
                continue;
            }

            return {
                provider: candidate.provider,
                modelId: candidate.modelId,
                reason: req.provider === candidate.provider ? 'PROVIDER_PREFERENCE' : 'PRIMARY_CHOICE',
                systemPrompt: candidate.systemPrompt
            };
        }

        console.error('[ModelSelector] ALL MODELS DEPLETED! Returning default fallback.');
        return {
            provider: 'lmstudio',
            modelId: 'local',
            reason: 'EMERGENCY_FALLBACK'
        };
    }
}