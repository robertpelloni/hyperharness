import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ModelSelector } from './ModelSelector.js';

describe('ModelSelector', () => {
    beforeEach(() => {
        process.env.GOOGLE_API_KEY = 'test-google-key';
        process.env.OPENAI_API_KEY = 'test-openai-key';
        process.env.ANTHROPIC_API_KEY = 'test-anthropic-key';
        process.env.DEEPSEEK_API_KEY = 'test-deepseek-key';
    });

    it('prefers the requested provider when it is available in the active chain', async () => {
        const selector = new ModelSelector();
        vi.spyOn(selector as any, 'getConfiguredChain').mockReturnValue(null);

        const result = await selector.selectModel({ provider: 'anthropic', taskType: 'worker' });

        expect(result.provider).toBe('anthropic');
        expect(result.modelId).toBe('claude-sonnet-4-20250514');
        expect(result.reason).toBe('PROVIDER_PREFERENCE');
    });

    it('skips provider-specific failures and falls forward to the next candidate', async () => {
        const selector = new ModelSelector();
        vi.spyOn(selector as any, 'getConfiguredChain').mockReturnValue(null);

        selector.reportFailure('google', 'gemini-2.0-flash');

        const result = await selector.selectModel({ taskType: 'worker' });

        expect(result.provider).toBe('anthropic');
        expect(result.modelId).toBe('claude-sonnet-4-20250514');
    });

    it('forces a local fallback when the budget is exceeded', async () => {
        const selector = new ModelSelector();
        vi.spyOn(selector as any, 'getConfiguredChain').mockReturnValue(null);
        vi.spyOn(selector.getQuotaService(), 'isBudgetExceeded').mockReturnValue(true);

        const result = await selector.selectModel({ taskType: 'worker' });

        expect(result).toEqual({
            provider: 'lmstudio',
            modelId: 'local',
            reason: 'BUDGET_EXCEEDED_FORCED_LOCAL'
        });
    });
});