import { afterEach, describe, expect, it, vi } from 'vitest';
import { Command } from 'commander';

const queryTrpcMock = vi.fn();
const resolveControlPlaneLocationMock = vi.fn(() => ({
  source: 'default',
  baseUrl: 'http://127.0.0.1:4000/trpc',
  host: '127.0.0.1',
  port: 4000,
}));

vi.mock('../control-plane.js', () => ({
  queryTrpc: (...args: unknown[]) => queryTrpcMock(...args),
  resolveControlPlaneLocation: () => resolveControlPlaneLocationMock(),
}));

import { registerToolsCommand } from './tools.js';

const logSpy = vi.spyOn(console, 'log').mockImplementation(() => {});
const errorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

afterEach(() => {
  queryTrpcMock.mockReset();
  resolveControlPlaneLocationMock.mockClear();
  logSpy.mockClear();
  errorSpy.mockClear();
  process.exitCode = 0;
});

function createProgram(): Command {
  const program = new Command();
  registerToolsCommand(program);
  return program;
}

describe('registerToolsCommand', () => {
  it('lists live MCP tools as JSON', async () => {
    queryTrpcMock.mockResolvedValue([
      {
        name: 'search_tools',
        description: 'Search available tools',
        server: 'meta',
        serverDisplayName: 'Meta MCP',
        semanticGroup: 'tool-discovery',
        semanticGroupLabel: 'tool discovery',
        keywords: ['search', 'tools'],
        alwaysOn: true,
      },
    ]);

    const program = createProgram();
    await program.parseAsync(['tools', 'list', '--json'], { from: 'user' });

    expect(queryTrpcMock).toHaveBeenCalledWith('mcp.listTools');
    expect(logSpy).toHaveBeenCalledWith(JSON.stringify({
      tools: [
        {
          name: 'search_tools',
          description: 'Search available tools',
          server: 'meta',
          serverDisplayName: 'Meta MCP',
          semanticGroup: 'tool-discovery',
          semanticGroupLabel: 'tool discovery',
          keywords: ['search', 'tools'],
          alwaysOn: true,
        },
      ],
    }, null, 2));
  });

  it('searches live MCP tools as JSON with top-k limiting', async () => {
    queryTrpcMock.mockResolvedValue([
      {
        name: 'search_tools',
        description: 'Search available tools',
        server: 'meta',
        serverDisplayName: 'Meta MCP',
        loaded: true,
        hydrated: true,
        matchReason: 'matched tool discovery intent',
        rank: 1,
      },
      {
        name: 'list_all_tools',
        description: 'List available tools',
        server: 'meta',
        serverDisplayName: 'Meta MCP',
        loaded: true,
        hydrated: false,
        matchReason: 'matched inventory intent',
        rank: 2,
      },
    ]);

    const program = createProgram();
    await program.parseAsync(['tools', 'search', 'tool discovery', '--top-k', '1', '--json'], { from: 'user' });

    expect(queryTrpcMock).toHaveBeenCalledWith('mcp.searchTools', {
      query: 'tool discovery',
    });
    expect(logSpy).toHaveBeenCalledWith(JSON.stringify({
      query: 'tool discovery',
      profile: null,
      results: [
        {
          name: 'search_tools',
          description: 'Search available tools',
          server: 'meta',
          serverDisplayName: 'Meta MCP',
          loaded: true,
          hydrated: true,
          matchReason: 'matched tool discovery intent',
          rank: 1,
        },
      ],
    }, null, 2));
  });

  it('reports control-plane failures without throwing out of the command', async () => {
    queryTrpcMock.mockRejectedValue(new Error('control plane unavailable'));

    const program = createProgram();
    await program.parseAsync(['tools', 'list'], { from: 'user' });

    expect(errorSpy).toHaveBeenCalled();
    expect(process.exitCode).toBe(1);
  });

  it('shows tool detail as JSON from the live control plane', async () => {
    queryTrpcMock.mockResolvedValue({
      uuid: 'search_tools',
      name: 'search_tools',
      description: 'Search available tools',
      server: 'meta',
      inputSchema: {
        type: 'object',
        properties: {
          query: {
            type: 'string',
            description: 'Search query',
          },
        },
        required: ['query'],
      },
      schemaParamCount: 1,
      always_on: true,
      isDeferred: false,
      mcpServerUuid: 'server-1',
    });

    const program = createProgram();
    await program.parseAsync(['tools', 'info', 'search_tools', '--json'], { from: 'user' });

    expect(queryTrpcMock).toHaveBeenCalledWith('tools.get', { uuid: 'search_tools' });
    expect(logSpy).toHaveBeenCalledWith(JSON.stringify({
      tool: {
        uuid: 'search_tools',
        name: 'search_tools',
        description: 'Search available tools',
        server: 'meta',
        inputSchema: {
          type: 'object',
          properties: {
            query: {
              type: 'string',
              description: 'Search query',
            },
          },
          required: ['query'],
        },
        schemaParamCount: 1,
        always_on: true,
        isDeferred: false,
        mcpServerUuid: 'server-1',
      },
    }, null, 2));
  });
});
