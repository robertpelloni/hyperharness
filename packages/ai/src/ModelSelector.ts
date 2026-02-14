import fs from 'fs';
import path from 'path';
import { QuotaService } from './services/QuotaService.js';

export interface ModelSelectionRequest {
    provider?: string;
    taskComplexity?: 'low' | 'medium' | 'high';
    taskType?: 'worker' | 'supervisor'; // explicit role override
}

export interface SelectedModel {
    provider: string;
    modelId: string;
    reason: string;
    systemPrompt?: string; // New: Persona Support
}

interface ModelStatus {
    isDepleted: boolean;
    depletedAt?: number;
    retryAfter?: number; // timestamp
}

// Default Fallback - Robust Chain
// Defines the priority order for model selection based on task type.
// The selector iterates through this list, checking for API keys and depletion status.
const DEFAULT_CHAINS = {
    worker: [
        { provider: 'google', modelId: 'gemini-2.0-flash' },
        { provider: 'anthropic', modelId: 'claude-3-5-sonnet-latest' },
        { provider: 'openai', modelId: 'gpt-4o' },
        { provider: 'deepseek', modelId: 'deepseek-chat' },
        { provider: 'lmstudio', modelId: 'local' },
        { provider: 'ollama', modelId: 'gemma:2b' }
    ],
    supervisor: [
        { provider: 'openai', modelId: 'gpt-4o' },
        { provider: 'anthropic', modelId: 'claude-3-5-sonnet-latest' },
        { provider: 'google', modelId: 'gemini-1.5-pro' },
        { provider: 'lmstudio', modelId: 'local' }
    ]
};

const COOL_DOWN_MS = 60 * 1000;

/**
 * ModelSelector Class
 * Responsible for choosing the best LLM for a given task.
 * Features:
 * - Dynamic Chain Selection: Defaults to hardcoded robust chains but can load overrides from `council.json`.
 * - Quota Management: Integrates with QuotaService to detect budget limits.
 * - Failure Handling: Marks models as depleted upon failure and rotates to the next available provider.
 * - Local Fallback: Gracefully degrades to local models (LM Studio/Ollama) if cloud quotas are hit.
 */
export class ModelSelector {
    private modelStates: Map<string, ModelStatus> = new Map();
    private configPath: string;
    private quotaService: QuotaService;

    constructor() {
        this.quotaService = new QuotaService();
        this.configPath = path.resolve(process.cwd(), 'packages/core/config/council.json');
        console.log("ModelSelector initialized. Config Path:", this.configPath);
    }

    public getQuotaService() { return this.quotaService; }

    public reportFailure(modelId: string) {
        console.warn(`[ModelSelector] Reporting failure for ${modelId}. Marking as DEPLETED.`);
        this.modelStates.set(modelId, {
            isDepleted: true,
            depletedAt: Date.now(),
            retryAfter: Date.now() + COOL_DOWN_MS
        });
    }

    private loadConfig() {
        try {
            if (fs.existsSync(this.configPath)) {
                const raw = fs.readFileSync(this.configPath, 'utf-8');
                return JSON.parse(raw);
            }
        } catch (e) {
            console.error("Failed to load council config:", e);
        }
        return null;
    }

    private hasKey(provider: string): boolean {
        if (provider === 'google') return !!process.env.GOOGLE_API_KEY;
        if (provider === 'openai') return !!process.env.OPENAI_API_KEY;
        if (provider === 'anthropic') return !!process.env.ANTHROPIC_API_KEY;
        if (provider === 'deepseek') return !!process.env.DEEPSEEK_API_KEY;
        return true; // Local/No-Auth
    }

    public async selectModel(req: ModelSelectionRequest): Promise<SelectedModel> {
        const isBudgetExceeded = this.quotaService.isBudgetExceeded();

        if (isBudgetExceeded) {
            console.warn("[ModelSelector] Budget exceeded! Forcing FREE local models.");
            return {
                provider: 'lmstudio',
                modelId: 'local',
                reason: 'BUDGET_EXCEEDED_FORCED_LOCAL'
            };
        }

        let chain = DEFAULT_CHAINS.worker;

        // Dynamic Load
        const config = this.loadConfig();

        if (config && config.members) {
            // Map council members to a chain
            chain = config.members.map((m: any) => ({
                provider: m.provider,
                modelId: m.modelId,
                systemPrompt: m.systemPrompt
            }));
        } else if (req.taskType === 'supervisor' || req.taskComplexity === 'high') {
            chain = DEFAULT_CHAINS.supervisor;
        }

        // Iterate
        for (const candidate of chain) {
            // 1. Check API Key
            if (!this.hasKey(candidate.provider)) {
                // console.warn(`[ModelSelector] Skipping ${candidate.provider} (No Key)`);
                continue;
            }

            // 2. Check Depletion
            const status = this.modelStates.get(candidate.modelId);

            if (status && status.isDepleted) {
                if (Date.now() > (status.retryAfter || 0)) {
                    this.modelStates.delete(candidate.modelId);
                } else {
                    continue;
                }
            }

            return {
                provider: candidate.provider,
                modelId: candidate.modelId,
                reason: status ? 'RECOVERED' : 'PRIMARY_CHOICE',
                // @ts-ignore
                systemPrompt: candidate.systemPrompt
            };
        }

        console.error("[ModelSelector] ALL MODELS DEPLETED! Returning default fallback.");
        return {
            provider: 'lmstudio',
            modelId: 'local',
            reason: 'EMERGENCY_FALLBACK'
        };
    }
}
