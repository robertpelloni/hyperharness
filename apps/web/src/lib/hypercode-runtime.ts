import { existsSync, readFileSync } from 'node:fs';
import os from 'node:os';
import path from 'node:path';

export type HypercodeStartupProvenance = {
  requestedRuntime?: string;
  activeRuntime?: string;
  launchMode?: string;
  dashboardMode?: string;
  installDecision?: string;
  installReason?: string;
  buildDecision?: string;
  buildReason?: string;
  updatedAt?: string;
};

type HypercodeLockRecord = {
  port?: number;
  host?: string;
  startup?: HypercodeStartupProvenance | null;
};

function resolveBrowserHost(host: string): string {
  return host === '0.0.0.0' || host === '::' || host === '[::]'
    ? '127.0.0.1'
    : host;
}

export function resolveHypercodeConfigDir(): string {
  const configuredDir = process.env.HYPERCODE_CONFIG_DIR?.trim();
  if (configuredDir) {
    return configuredDir;
  }

  return path.join(os.homedir(), '.hypercode');
}

export function resolveHypercodeLockPath(): string {
  return path.join(resolveHypercodeConfigDir(), 'lock');
}

function readHypercodeLockRecord(): HypercodeLockRecord | null {
  const lockPath = resolveHypercodeLockPath();
  if (!existsSync(lockPath)) {
    return null;
  }

  try {
    return JSON.parse(readFileSync(lockPath, 'utf8')) as HypercodeLockRecord;
  } catch {
    return null;
  }
}

export function readLocalStartupProvenance(): HypercodeStartupProvenance | null {
  return readHypercodeLockRecord()?.startup ?? null;
}

export function resolveLockedHypercodeBase(): string | null {
  const parsed = readHypercodeLockRecord();
  if (!parsed || typeof parsed.port !== 'number' || parsed.port <= 0) {
    return null;
  }

  const host = typeof parsed.host === 'string' && parsed.host.trim().length > 0
    ? resolveBrowserHost(parsed.host.trim())
    : '127.0.0.1';

  return `http://${host}:${parsed.port}`;
}
