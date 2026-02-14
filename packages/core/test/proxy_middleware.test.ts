import { describe, test, expect } from 'bun:test';
import { EventEmitter } from 'events';
import { McpProxyManager } from '../src/managers/McpProxyManager.js';

type ProxyCtorArgs = ConstructorParameters<typeof McpProxyManager>;

interface PolicyServiceLike {
    evaluate(ctx: unknown): { allowed: boolean; reason?: string };
}

interface SavedScriptServiceLike {
    getAllScripts(): unknown[];
}

interface ToolListResult {
    tools: Array<{ name?: string }>;
}

interface CallRequestShape {
    params?: { name?: string };
}

class MockMcpManager extends EventEmitter {
    getClient(name: string) {
        return null;
    }
    getAllServers() {
        return [];
    }
}

class MockLogManager {
    log(_entry: unknown) {}
    calculateCost() {
        return 0;
    }
}

describe('McpProxyManager middleware', () => {
    test('list middleware can filter tools', async () => {
        process.env.MCP_DISABLE_METAMCP = 'true';

        const policyService: PolicyServiceLike = { evaluate: () => ({ allowed: true }) };
        const savedScriptService: SavedScriptServiceLike = { getAllScripts: () => [] };

        const proxy = new McpProxyManager(
            new MockMcpManager() as unknown as ProxyCtorArgs[0],
            new MockLogManager() as unknown as ProxyCtorArgs[1],
            { policyService, savedScriptService } as unknown as ProxyCtorArgs[2]
        );

        proxy.useListToolsMiddleware((next: (req: unknown, ctx: unknown) => Promise<ToolListResult>) => async (req: unknown, ctx: unknown) => {
            const res = await next(req, ctx);
            return { tools: res.tools.filter((t) => t.name !== 'search_tools') };
        });

        const tools = await proxy.getAllTools('s');
        expect(tools.find((t: { name?: string }) => t.name === 'search_tools')).toBeUndefined();
    });

    test('call middleware can block tool calls', async () => {
        process.env.MCP_DISABLE_METAMCP = 'true';

        const policyService: PolicyServiceLike = { evaluate: () => ({ allowed: true }) };
        const savedScriptService: SavedScriptServiceLike = { getAllScripts: () => [] };

        const proxy = new McpProxyManager(
            new MockMcpManager() as unknown as ProxyCtorArgs[0],
            new MockLogManager() as unknown as ProxyCtorArgs[1],
            { policyService, savedScriptService } as unknown as ProxyCtorArgs[2]
        );

        proxy.useCallToolMiddleware((_next: (req: unknown, ctx: unknown) => Promise<unknown>) => async (req: CallRequestShape, _ctx: unknown) => {
            if (req.params.name === 'run_code') {
                return { isError: true, content: [{ type: 'text', text: 'blocked by mw' }] };
            }
            return { content: [{ type: 'text', text: 'ok' }] };
        });

        const res = await proxy.callTool('run_code', { code: 'return 1;' }, 's');
        expect(res.isError).toBeTruthy();
    });
});
