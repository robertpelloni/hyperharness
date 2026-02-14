import { describe, test, expect } from 'bun:test';
import { EventEmitter } from 'events';
import { McpProxyManager } from '../src/managers/McpProxyManager.js';
import { createLoggingMiddleware } from '../src/middleware/logging-middleware.js';

type ProxyCtorArgs = ConstructorParameters<typeof McpProxyManager>;

interface PolicyServiceLike {
  evaluate(ctx: unknown): { allowed: boolean; reason?: string };
}

interface SavedScriptServiceLike {
  getAllScripts(): unknown[];
}

interface LogEntry {
  type?: string;
  args?: {
    namespaceId?: string;
    endpointPath?: string;
  };
}

class MockMcpManager extends EventEmitter {
  getClient(_name: string) {
    return null;
  }
  getAllServers() {
    return [];
  }
}

class MockLogManager {
  entries: LogEntry[] = [];
  log(entry: LogEntry) {
    this.entries.push(entry);
  }
  calculateCost() {
    return 0;
  }
}

describe('proxy logging middleware', () => {
  test('logs include namespaceId and endpointPath', async () => {
    const logManager = new MockLogManager();

    const policyService: PolicyServiceLike = { evaluate: () => ({ allowed: true }) };
    const savedScriptService: SavedScriptServiceLike = { getAllScripts: () => [] };

    const proxy = new McpProxyManager(
      new MockMcpManager() as unknown as ProxyCtorArgs[0],
      logManager as unknown as ProxyCtorArgs[1],
      { policyService, savedScriptService } as unknown as ProxyCtorArgs[2]
    );

    proxy.useCallToolMiddleware(createLoggingMiddleware({ enabled: true, logManager }));

    await proxy.start();

    proxy.setNamespaceForSession('s', 'ns1');
    proxy.setEndpointForSession('s', '/api/mcp/coding');

    await proxy.callTool('namespace_get', {}, 's');

    const requestLog = logManager.entries.find(e => e.type === 'request');
    expect(requestLog.args.namespaceId).toBe('ns1');
    expect(requestLog.args.endpointPath).toBe('/api/mcp/coding');
  });
});
