import { mkdtempSync, rmSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { tmpdir } from 'node:os';

import { afterEach, describe, expect, it, vi } from 'vitest';

import { queryTrpc, readLocalStartupProvenance, resolveControlPlaneLocation } from './control-plane.js';

const tempDirs: string[] = [];

afterEach(() => {
  while (tempDirs.length > 0) {
    const dir = tempDirs.pop();
    if (dir) {
      rmSync(dir, { recursive: true, force: true });
    }
  }
});

function createTempDir(): string {
  const dir = mkdtempSync(join(tmpdir(), 'hypercode-control-plane-'));
  tempDirs.push(dir);
  return dir;
}

describe('resolveControlPlaneLocation', () => {
  it('prefers the explicit upstream and normalizes missing /trpc suffixes', () => {
    expect(resolveControlPlaneLocation({
      upstream: 'http://localhost:3100',
    })).toEqual({
      source: 'env',
      baseUrl: 'http://localhost:3100/trpc',
      host: 'localhost',
      port: 3100,
    });
  });

  it('uses the lock file when available and rewrites wildcard hosts for local access', () => {
    const dataDir = createTempDir();
    writeFileSync(join(dataDir, 'lock'), JSON.stringify({
      instanceId: 'hypercode-test',
      pid: 1234,
      port: 4555,
      host: '0.0.0.0',
      createdAt: '2026-03-30T00:00:00.000Z',
    }), 'utf8');

    expect(resolveControlPlaneLocation({ dataDir })).toEqual({
      source: 'lock',
      baseUrl: 'http://127.0.0.1:4555/trpc',
      host: '127.0.0.1',
      port: 4555,
    });
  });

  it('falls back to the default control-plane address when no upstream is configured', () => {
    expect(resolveControlPlaneLocation({
      dataDir: createTempDir(),
    })).toEqual({
      source: 'default',
      baseUrl: 'http://127.0.0.1:4000/trpc',
      host: '127.0.0.1',
      port: 4000,
    });
  });
});

describe('readLocalStartupProvenance', () => {
  it('derives control-plane port provenance from legacy lock files', () => {
    const dataDir = createTempDir();
    writeFileSync(join(dataDir, 'lock'), JSON.stringify({
      instanceId: 'hypercode-legacy',
      pid: 1234,
      port: 4555,
      host: '127.0.0.1',
      createdAt: '2026-04-06T09:00:00.000Z',
    }), 'utf8');

    expect(readLocalStartupProvenance(dataDir)).toEqual({
      requestedPort: 4555,
      activePort: 4555,
      portDecision: 'derived from lock record',
      portReason: 'Detailed startup port provenance was unavailable; using the current control-plane lock port.',
      updatedAt: '2026-04-06T09:00:00.000Z',
    });
  });
});

  });

  it('surfaces tRPC error payloads clearly', async () => {
    await expect(queryTrpc('mesh.getStatus', undefined, {
      upstream: 'http://localhost:4000',
      fetchImpl: vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
        json: async () => ({
          error: {
            message: 'mesh unavailable',
          },
        }),
      }) as unknown as typeof fetch,
    })).rejects.toThrow('HyperCode control plane query failed for mesh.getStatus: mesh unavailable');
  });
});
