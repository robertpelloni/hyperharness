import { z } from 'zod';
import { t, publicProcedure, getMcpServer, getConfigManager } from '../lib/trpc-core.js';

function getErrorMessage(error: unknown): string {
    if (error instanceof Error) {
        return error.message;
    }
    return typeof error === 'string' ? error : 'Unknown error';
}

export const settingsRouter = t.router({
    /** Get the full configuration object */
    get: publicProcedure.query(() => {
        return getConfigManager()?.loadConfig() || {};
    }),

    /** Update configuration with a partial config object */
    update: publicProcedure.input(z.object({
        config: z.record(z.unknown())
    })).mutation(({ input }) => {
        const configManager = getConfigManager();
        if (configManager) {
            configManager.saveConfig(input.config);
            return { success: true };
        }
        throw new Error("ConfigManager not ready");
    }),

    /** Get detected LLM providers and their API key status (masked) */
    getProviders: publicProcedure.query(() => {
        const providers = [
            { id: 'openai', name: 'OpenAI', envVar: 'OPENAI_API_KEY' },
            { id: 'anthropic', name: 'Anthropic', envVar: 'ANTHROPIC_API_KEY' },
            { id: 'gemini', name: 'Google Gemini', envVar: 'GEMINI_API_KEY' },
            { id: 'xai', name: 'xAI (Grok)', envVar: 'XAI_API_KEY' },
            { id: 'deepseek', name: 'DeepSeek', envVar: 'DEEPSEEK_API_KEY' },
            { id: 'mistral', name: 'Mistral', envVar: 'MISTRAL_API_KEY' },
            { id: 'openrouter', name: 'OpenRouter', envVar: 'OPENROUTER_API_KEY' },
            { id: 'copilot', name: 'GitHub Copilot', envVar: 'COPILOT_API_KEY' },
            { id: 'cohere', name: 'Cohere', envVar: 'COHERE_API_KEY' },
            { id: 'groq', name: 'Groq', envVar: 'GROQ_API_KEY' },
            { id: 'together', name: 'Together AI', envVar: 'TOGETHER_API_KEY' },
            { id: 'fireworks', name: 'Fireworks AI', envVar: 'FIREWORKS_API_KEY' },
        ];

        return providers.map(p => ({
            ...p,
            configured: !!process.env[p.envVar],
            keyPreview: process.env[p.envVar]
                ? `${process.env[p.envVar]!.slice(0, 4)}...${process.env[p.envVar]!.slice(-4)}`
                : null,
        }));
    }),

    /** Test a provider connection — attempts a minimal API call */
    testConnection: publicProcedure.input(z.object({
        provider: z.string(),
    })).mutation(async ({ input }) => {
        const mcp = getMcpServer();
        try {
            const llm = mcp.llmService;
            // Attempt a minimal completion to verify connectivity
            const result = await llm.generate(
                `You are a test. Reply with "ok".`,
                { model: input.provider, maxTokens: 5 }
            );
            return {
                success: true,
                provider: input.provider,
                response: result?.text?.slice(0, 50) ?? 'ok',
                latencyMs: result?.latencyMs ?? 0,
            };
        } catch (err: unknown) {
            return {
                success: false,
                provider: input.provider,
                error: getErrorMessage(err),
            };
        }
    }),

    /** Get environment info for debugging */
    getEnvironment: publicProcedure.query(() => {
        return {
            nodeVersion: process.version,
            platform: process.platform,
            arch: process.arch,
            cwd: process.cwd(),
            env: {
                NODE_ENV: process.env.NODE_ENV ?? 'development',
                PORT: process.env.PORT ?? '3000',
            },
        };
    }),

    /** Get list of configured MCP servers */
    getMcpServers: publicProcedure.query(() => {
        const configManager = getConfigManager();
        if (configManager) {
            const config = configManager.loadConfig();
            return config?.mcpServers ?? [];
        }
        return [];
    }),
});
