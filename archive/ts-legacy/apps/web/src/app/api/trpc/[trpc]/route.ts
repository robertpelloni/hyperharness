import { createHash } from 'node:crypto';
import { existsSync } from 'node:fs';
import fs from 'node:fs/promises';
import os from 'node:os';
import path from 'node:path';

import { normalizeImportedServerType } from '../../../../lib/mcp-import';
import { resolveUpstreamBases } from '../../../../lib/trpc-upstream';

export const runtime = 'nodejs';

function resolveRepoRoot(): string {
  const candidates = [
    process.cwd(),
    path.resolve(process.cwd(), '..'),
    path.resolve(process.cwd(), '..', '..'),
  ];

  for (const candidate of candidates) {
    if (existsSync(path.join(candidate, 'mcp.jsonc')) || existsSync(path.join(candidate, 'mcp.json'))) {
      return candidate;
    }
  }

  for (const candidate of candidates) {
    if (existsSync(path.join(candidate, 'pnpm-workspace.yaml'))) {
      return candidate;
    }
  }

  return path.resolve(process.cwd(), '..', '..');
}

const LEGACY_REPO_ROOT = resolveRepoRoot();
const LEGACY_MCP_JSONC_PATH = path.join(LEGACY_REPO_ROOT, 'mcp.jsonc');
const LEGACY_MCP_JSON_PATH = path.join(LEGACY_REPO_ROOT, 'mcp.json');

function resolveBorgConfigDir(): string {
  const configuredDir = process.env.BORG_CONFIG_DIR?.trim();
  if (configuredDir) {
    return configuredDir;
  }

  return path.join(os.homedir(), '.borg');
}

function resolvePrimaryMcpPaths(): { jsoncPath: string; jsonPath: string } {
  const configDir = resolveBorgConfigDir();
  return {
    jsoncPath: path.join(configDir, 'mcp.jsonc'),
    jsonPath: path.join(configDir, 'mcp.json'),
  };
}

function resolveMcpReadCandidates(): Array<{ filePath: string; allowComments: boolean }> {
  const primaryPaths = resolvePrimaryMcpPaths();
  return [
    { filePath: primaryPaths.jsoncPath, allowComments: true },
    { filePath: primaryPaths.jsonPath, allowComments: false },
    { filePath: LEGACY_MCP_JSONC_PATH, allowComments: true },
    { filePath: LEGACY_MCP_JSON_PATH, allowComments: false },
  ];
}

type LocalMcpServerEntry = {
  command?: string;
  args?: string[];
  env?: Record<string, string>;
  url?: string;
  disabled?: boolean;
  description?: string | null;
  type?: 'STDIO' | 'SSE' | 'STREAMABLE_HTTP';
  bearerToken?: string;
  headers?: Record<string, string>;
  _meta?: {
    uuid?: string;
    status?: string;
    metadataSource?: string;
    toolCount?: number;
    lastSuccessfulBinaryLoadAt?: string | null;
    crashCount?: number;
    maxAttempts?: number;
  } | null;
};

type LocalMcpConfig = {
  mcpServers: Record<string, LocalMcpServerEntry>;
};

type LocalManagedServerRecord = {
  uuid: string;
  name: string;
  description: string | null;
  type: 'STDIO' | 'SSE' | 'STREAMABLE_HTTP';
  command: string | null;
  args: string[];
  env: Record<string, string>;
  url: string | null;
  bearerToken: string | null;
  headers: Record<string, string>;
  _meta: {
    uuid: string;
    status: string;
    metadataSource: string;
    toolCount: number;
    lastSuccessfulBinaryLoadAt: string | null;
    crashCount: number;
    maxAttempts: number;
  };
};

const LOCAL_COMPAT_METADATA_SOURCE = 'local-config-fallback';

const LEGACY_COMPAT_RESPONSE_KEYS = {
  'mcpServers.list': 'mcpServers.list',
  'mcp.listServers': 'mcpServers.list',
  'tools.list': 'tools.list',
  'mcp.listTools': 'tools.list',
  'mcp.getStatus': 'mcp.getStatus',
  'mcp.traffic': 'mcp.traffic',
  'session.list': 'session.list',
  'billing.getProviderQuotas': 'billing.getProviderQuotas',
  'billing.getFallbackChain': 'billing.getFallbackChain',
} as const;

type LegacyCompatProcedure = keyof typeof LEGACY_COMPAT_RESPONSE_KEYS;
type LegacyCompatResponseKey = typeof LEGACY_COMPAT_RESPONSE_KEYS[LegacyCompatProcedure];

const LOCAL_COMPAT_RESPONSE_KEYS = {
  ...LEGACY_COMPAT_RESPONSE_KEYS,
  startupStatus: 'startupStatus',
  'mcp.getWorkingSet': 'mcp.getWorkingSet',
  'mcp.getToolSelectionTelemetry': 'mcp.getToolSelectionTelemetry',
  'mcp.searchTools': 'mcp.searchTools',
  'mcp.getToolPreferences': 'mcp.getToolPreferences',
  'mcp.getJsoncEditor': 'mcp.getJsoncEditor',
  'mcpServers.get': 'mcpServers.get',
  'apiKeys.list': 'apiKeys.list',
  'savedScripts.list': 'savedScripts.list',
  'mcpServers.create',
  'mcpServers.update',
  'mcpServers.delete',
  'mcpServers.reloadMetadata',
  'mcpServers.clearMetadataCache',
  'serverHealth.reset',
]);
const LOCAL_COMPAT_MUTATION_PROCEDURES = new Set([
  ...LOCAL_MCP_MUTATION_PROCEDURES,
  ...LOCAL_OPERATOR_MUTATION_PROCEDURES,
]);
    checks: {
      mcpAggregator: {
        ready: serverCount > 0,
        liveReady: serverCount > 0,
        initialization: 'compat-fallback',
        serverCount,
        connectedCount,
        persistedServerCount: serverCount,
        persistedToolCount: 0,
        configuredServerCount: serverCount,
        advertisedServerCount: serverCount,
        advertisedToolCount: 0,
        advertisedAlwaysOnServerCount: 0,
        advertisedAlwaysOnToolCount: 0,
        inventoryReady: false,
        warmupInProgress: true,
      },
      configSync: {
        ready: serverCount > 0,
        status: {
          inProgress: false,
          lastCompletedAt: null,
          lastServerCount: serverCount,
          lastToolCount: 0,
        },
      },
      sessionSupervisor: {
        ready: false,
        sessionCount: 0,
        restore: null,
      },
      browser: {
        ready: false,
        active: false,
        pageCount: 0,
      },
      memory: {
        ready: false,
        initialized: false,
        agentMemory: false,
      },
      extensionBridge: {
        ready: false,
        acceptingConnections: false,
        clientCount: 0,
        hasConnectedClients: false,
        clients: [],
        supportedCapabilities: [],
        supportedHookPhases: [],
      },
      executionEnvironment: {
        ready: false,
        preferredShellId: null,
        preferredShellLabel: null,
        shellCount: 0,
        verifiedShellCount: 0,
        toolCount: 0,
        verifiedToolCount: 0,
        harnessCount: 0,
        verifiedHarnessCount: 0,
        supportsPowerShell: false,
        supportsPosixShell: false,
        notes: [],
      },
    },
  };
}

function resolveNativeStatusBases(): string[] {
  return Array.from(new Set(
    resolveUpstreamBases()
      .map((base) => base.replace(/\/trpc\/?$/i, '').trim())
      .filter(Boolean),
  ));
}

async function fetchNativeControlPlaneData<T>(endpointPath: string, init?: RequestInit): Promise<T | null> {
  for (const base of resolveNativeStatusBases()) {
    try {
      const response = await fetch(`${base}${endpointPath}`, {
        cache: 'no-store',
        ...init,
      });
      if (!response.ok) {
        continue;
      }

      const payload = await response.json() as { success?: unknown; data?: unknown };
      if (payload.success !== true) {
        continue;
      }

      return (payload.data as T | undefined) ?? null;
    } catch {
      // Try the next native control-plane base.
    }
  }

  return null;
}

    },
  });
}

function buildLegacyCompatResponse(req: Request): Response | null {
  return null;
}

async function buildPreferredSavedScriptsList(): Promise<unknown[]> {
  const scripts = await fetchNativeControlPlaneData<unknown[]>('/api/scripts');
  return Array.isArray(scripts) ? scripts : [];
}


  const dataByResponseKey: Record<LocalCompatResponseKey, unknown> = {
    'mcpServers.list': localServers,
    'tools.list': LEGACY_COMPAT_RESPONSES['tools.list'],
    'mcp.getStatus': localStatus,
    'mcp.traffic': LEGACY_COMPAT_RESPONSES['mcp.traffic'],
    'session.list': LEGACY_COMPAT_RESPONSES['session.list'],
    'billing.getProviderQuotas': LEGACY_COMPAT_RESPONSES['billing.getProviderQuotas'],
    'billing.getFallbackChain': LEGACY_COMPAT_RESPONSES['billing.getFallbackChain'],
    startupStatus: localStartupStatus,
    'mcp.getWorkingSet': {
      tools: [],
      limits: {
        maxLoadedTools: 24,
        maxHydratedSchemas: 8,
      },
    },
    'mcp.getToolSelectionTelemetry': [],
    'mcp.searchTools': [],
    'mcp.getToolPreferences': {
      importantTools: [],
      alwaysLoadedTools: ['search_tools', 'read_file', 'write_file', 'grep_search', 'execute_command', 'browser__open'],
    },
    'mcp.getJsoncEditor': {
      path: localConfigSource.path,
      content: localConfigSource.content,
    },
    'mcpServers.get': undefined,
    'apiKeys.list': [],
    'savedScripts.list': savedScripts,
    },
  });
}

function normalizeTrpcPayload(payload: unknown): unknown[] {
  if (Array.isArray(payload)) {
    return payload;
  }

  if (payload && typeof payload === 'object') {
    return [payload];
  }

  return [];
}

function extractLegacyResponseData(entry: unknown): unknown {
  if (!entry || typeof entry !== 'object') {
    return undefined;
  }

  const result = (entry as { result?: { data?: unknown } }).result;
  if (!result || typeof result !== 'object') {
    return undefined;
  }

  return result.data;
}

function isLegacyResponseEffectivelyEmpty(req: Request, payload: unknown): boolean {
  const procedures = getProcedureNames(req);
  const entries = normalizeTrpcPayload(payload);

  if (procedures.length === 0 || entries.length === 0) {
    return false;
  }

  return procedures.every((procedureName, index) => {
    const responseKey = getLegacyCompatResponseKey(procedureName);
    if (!responseKey) {
      return false;
    }

    const data = extractLegacyResponseData(entries[index]);

    if (responseKey === 'mcpServers.list' || responseKey === 'tools.list' || responseKey === 'mcp.traffic' || responseKey === 'session.list' || responseKey === 'billing.getProviderQuotas') {
      return Array.isArray(data) && data.length === 0;
    }

    if (responseKey === 'billing.getFallbackChain') {
      return Boolean(data) && typeof data === 'object' && Array.isArray((data as { chain?: unknown[] }).chain) && ((data as { chain?: unknown[] }).chain?.length ?? 0) === 0;
    }

    if (responseKey === 'mcp.getStatus') {
      return Boolean(data)
        && typeof data === 'object'
        && Number((data as { serverCount?: unknown }).serverCount ?? 0) === 0
        && Number((data as { toolCount?: unknown }).toolCount ?? 0) === 0
        && Number((data as { connectedCount?: unknown }).connectedCount ?? 0) === 0;
    }

    return false;
  });
}

async function shouldPreferLocalCompat(req: Request, upstreamResponse: Response): Promise<boolean> {
  if (!isLegacyMcpRequest(req) || !upstreamResponse.ok) {
    return false;
  }

  const localConfig = await loadLocalMcpConfig();
  if (Object.keys(localConfig.mcpServers).length === 0) {
    return false;
  }

  try {
    const payload = await upstreamResponse.clone().json();
    return isLegacyResponseEffectivelyEmpty(req, payload);
  } catch {
    return false;
  }
}

function cloneHeaders(req: Request): Headers {
  const headers = new Headers(req.headers);
  headers.delete('host');
  headers.delete('content-length');
  return headers;
}

function isBatchRequest(req: Request): boolean {
  return new URL(req.url).searchParams.get('batch') === '1';
}

function isBulkImportRequest(req: Request): boolean {
  const procedures = getProcedureNames(req);
  return req.method === 'POST' && procedures.length === 1 && procedures[0] === 'mcpServers.bulkImport';
}

function buildTrpcResponse(req: Request, data: unknown, init?: ResponseInit): Response {
  const payload = isBatchRequest(req)
    ? [{ result: { data } }]
    : { result: { data } };

  return new Response(JSON.stringify(payload), {
    status: init?.status ?? 200,
    statusText: init?.statusText,
    headers: {
      'content-type': 'application/json',
      ...(init?.headers ?? {}),
    },
  });
}

function parseTrpcInputEntries(parsed: unknown, req: Request): unknown[] {
  if (isBatchRequest(req) && parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
    return Object.entries(parsed as Record<string, unknown>)
      .sort(([left], [right]) => Number(left) - Number(right))
      .map(([, value]) => {
        if (value && typeof value === 'object' && 'json' in (value as Record<string, unknown>)) {
          return (value as { json?: unknown }).json;
        }

        return value;
      });
  }

  if (parsed && typeof parsed === 'object' && 'json' in (parsed as Record<string, unknown>)) {
    return [(parsed as { json?: unknown }).json];
  }

  return [parsed];
}

function extractTrpcProcedureInputs(body: string | undefined, req: Request): unknown[] {
  const sources = body
    ? [body, new URL(req.url).searchParams.get('input')]
    : [new URL(req.url).searchParams.get('input')];

  for (const source of sources) {
    if (!source) {
      continue;
    }

    try {
      const parsed = JSON.parse(source) as unknown;
      return parseTrpcInputEntries(parsed, req);
    } catch {
      // Try the next tRPC input source.
    }
  }

  return [];
}

function extractTrpcRequestInput(body: string | undefined, req: Request): unknown {
  return extractTrpcProcedureInputs(body, req)[0];
}

function normalizeBulkImportProxyBody(req: Request, body: string | undefined): string | undefined {
  if (!isBulkImportRequest(req)) {
    return body;
  }

  const input = extractTrpcRequestInput(body, req);
  if (!Array.isArray(input)) {
    return body;
  }

  return JSON.stringify(input);
}

async function tryBridgeBulkImport(
  req: Request,
  upstreamBases: string[],
  headers: Headers,
  body: string | undefined,
): Promise<Response | null> {
  const procedures = getProcedureNames(req);
  if (req.method !== 'POST' || procedures.length !== 1 || procedures[0] !== 'mcpServers.bulkImport') {
    return null;
  }

  const input = extractTrpcRequestInput(body, req);
  if (!Array.isArray(input)) {
    return null;
  }

  const mcpServers: Record<string, {
    type: unknown;
    description: unknown;
    command: unknown;
    args: unknown[];
    env: unknown;
    url: unknown;
    bearerToken: unknown;
    headers: unknown;
  }> = {};

  for (const item of input) {
    if (!item || typeof item !== 'object') {
      continue;
    }

    const name = String(item.name ?? '').trim();
    if (!name) {
      continue;
    }

    mcpServers[name] = {
      type: item.type ?? 'STDIO',
      description: item.description ?? undefined,
      command: item.command ?? undefined,
      args: Array.isArray(item.args) ? item.args : [],
      env: item.env ?? {},
      url: item.url ?? undefined,
      bearerToken: item.bearerToken ?? undefined,
      headers: item.headers ?? {},
    };
  }

  if (Object.keys(mcpServers).length === 0) {
    return buildTrpcResponse(req, { imported: 0, errors: ['No valid MCP servers were provided.'] }, { status: 200 });
  }

  for (const upstreamBase of upstreamBases) {
    const normalizedBase = upstreamBase.replace(/\/$/, '');

    for (const procedureName of LEGACY_MCP_BULK_IMPORT_PROCEDURES) {
      try {
        const response = await fetch(`${normalizedBase}/${procedureName}`, {
          method: 'POST',
          headers,
          body: JSON.stringify({ mcpServers }),
        });

        if (response.status === 404) {
          continue;
        }

        if (!response.ok) {
          continue;
        }

        const json = await response.json();
        return buildTrpcResponse(req, extractTrpcData(json), {
          status: response.status,
          statusText: response.statusText,
          headers: { 'x-hypercode-trpc-compat': 'legacy-mcp-bulk-import-bridge' },
        });
      } catch {
        // Try the next candidate upstream/procedure.
      }
    }
  }

  return null;
}

async function tryLocalBulkImport(req: Request, body: string | undefined): Promise<Response | null> {
  const procedures = getProcedureNames(req);
  if (req.method !== 'POST' || procedures.length !== 1 || procedures[0] !== 'mcpServers.bulkImport') {
    return null;
  }

  const input = extractTrpcRequestInput(body, req);
  if (!Array.isArray(input)) {
    return null;
  }

  const localConfig = await loadLocalMcpConfig();
  let imported = 0;

  for (const item of input) {
    if (!item || typeof item !== 'object') {
      continue;
    }

    const name = String(item.name ?? '').trim();
    if (!name) {
      continue;
    }

    localConfig.mcpServers[name] = {
      command: typeof item.command === 'string' ? item.command : undefined,
      args: Array.isArray(item.args) ? item.args.filter((arg): arg is string => typeof arg === 'string') : [],
      env: item.env && typeof item.env === 'object'
        ? Object.fromEntries(Object.entries(item.env as Record<string, unknown>).map(([key, value]) => [key, String(value)]))
        : undefined,
      url: typeof item.url === 'string' ? item.url : undefined,
      description: typeof item.description === 'string' ? item.description : null,
      type: normalizeImportedServerType({
        type: typeof item.type === 'string' ? item.type : undefined,
        url: typeof item.url === 'string' ? item.url : undefined,
      }),
      bearerToken: typeof item.bearerToken === 'string' ? item.bearerToken : undefined,
      headers: item.headers && typeof item.headers === 'object'
        ? Object.fromEntries(Object.entries(item.headers as Record<string, unknown>).map(([key, value]) => [key, String(value)]))
        : undefined,
    };
    imported += 1;
  }

  await writeLocalMcpConfig(localConfig);
  return buildTrpcResponse(req, { imported, errors: [] }, {
    status: 200,
    headers: { 'x-hypercode-trpc-compat': 'local-mcp-config-bulk-import' },
  });
}

async function tryLocalOperatorMutation(req: Request, body: string | undefined): Promise<Response | null> {
  const procedures = getProcedureNames(req);
  const procedureName = procedures[0] ?? '';
  if (req.method !== 'POST' || procedures.length !== 1 || !LOCAL_OPERATOR_MUTATION_PROCEDURES.has(procedureName)) {
    return null;
  }

  const input = extractTrpcRequestInput(body, req);
  if (!input || typeof input !== 'object') {
    return buildTrpcResponse(req, undefined, {
      status: 400,
      statusText: 'Invalid local operator compat input',
      headers: { 'x-hypercode-trpc-compat': 'local-operator-action' },
    });
  }

  let endpointPath = '';
  if (procedureName === 'savedScripts.create') {
    endpointPath = '/api/scripts/create';
  } else if (procedureName === 'savedScripts.update') {
    endpointPath = '/api/scripts/update';
  } else if (procedureName === 'savedScripts.delete') {
    endpointPath = '/api/scripts/delete';
  } else if (procedureName === 'savedScripts.execute') {
    endpointPath = '/api/scripts/execute';
  }

  if (!endpointPath) {
    return null;
  }

  const data = await fetchNativeControlPlaneData<unknown>(endpointPath, {
    method: 'POST',
    headers: {
      'content-type': 'application/json',
    },
    body: JSON.stringify(input),
  });
  if (data === null) {
    return null;
  }

  return buildTrpcResponse(req, data, {
    status: 200,
    headers: { 'x-hypercode-trpc-compat': 'local-operator-action' },
  });
}

async function tryLocalManagedServerMutation(req: Request, body: string | undefined): Promise<Response | null> {
  const procedures = getProcedureNames(req);
  if (req.method !== 'POST' || procedures.length !== 1 || !LOCAL_MCP_MUTATION_PROCEDURES.has(procedures[0] ?? '')) {
    return null;
  }

  const input = extractTrpcRequestInput(body, req);
  if (!input || typeof input !== 'object') {
    return buildTrpcResponse(req, undefined, {
      status: 400,
      statusText: 'Invalid local MCP compat input',
      headers: { 'x-hypercode-trpc-compat': 'local-mcp-managed-action' },
    });
  }

  const localConfig = await loadLocalMcpConfig();
  const procedureName = procedures[0] ?? '';

  if (procedureName === 'mcpServers.create') {
    const name = String((input as { name?: unknown }).name ?? '').trim();
    if (!name) {
      return buildTrpcResponse(req, undefined, {
        status: 400,
        statusText: 'Server name is required',
        headers: { 'x-hypercode-trpc-compat': 'local-mcp-managed-action' },
      });
    }

    if (localConfig.mcpServers[name]) {
      return buildTrpcResponse(req, undefined, {
        status: 409,
        statusText: 'Server already exists',
        headers: { 'x-hypercode-trpc-compat': 'local-mcp-managed-action' },
      });
    }

    localConfig.mcpServers[name] = {
      command: typeof (input as { command?: unknown }).command === 'string' ? (input as { command?: string }).command : undefined,
      args: Array.isArray((input as { args?: unknown[] }).args)
        ? (input as { args?: unknown[] }).args?.filter((arg): arg is string => typeof arg === 'string')
        : [],
      env: (input as { env?: unknown }).env && typeof (input as { env?: unknown }).env === 'object'
        ? Object.fromEntries(Object.entries((input as { env?: Record<string, unknown> }).env ?? {}).map(([key, value]) => [key, String(value)]))
        : {},
      url: typeof (input as { url?: unknown }).url === 'string' ? (input as { url?: string }).url : undefined,
      disabled: false,
      description: typeof (input as { description?: unknown }).description === 'string' ? (input as { description?: string }).description : null,
      type: normalizeImportedServerType({
        type: typeof (input as { type?: unknown }).type === 'string' ? (input as { type?: string }).type : undefined,
        url: typeof (input as { url?: unknown }).url === 'string' ? (input as { url?: string }).url : undefined,
      }),
      bearerToken: typeof (input as { bearerToken?: unknown }).bearerToken === 'string' ? (input as { bearerToken?: string }).bearerToken : undefined,
      headers: (input as { headers?: unknown }).headers && typeof (input as { headers?: unknown }).headers === 'object'
        ? Object.fromEntries(Object.entries((input as { headers?: Record<string, unknown> }).headers ?? {}).map(([key, value]) => [key, String(value)]))
        : {},
    };

    upsertLocalServerMeta(localConfig, name);
    await writeLocalMcpConfig(localConfig);

    return buildTrpcResponse(req, buildLocalManagedServerRecord(name, localConfig.mcpServers[name]), {
      status: 200,
      headers: { 'x-hypercode-trpc-compat': 'local-mcp-managed-action' },
    });
  }

  const localUuid = String((input as { uuid?: unknown; serverUuid?: unknown }).uuid ?? (input as { serverUuid?: unknown }).serverUuid ?? '').trim();
  const match = localUuid ? findLocalServerByUuid(localConfig, localUuid) : null;
  if (!match) {
    return buildTrpcResponse(req, undefined, {
      status: 404,
      statusText: 'Local managed MCP server not found',
      headers: { 'x-hypercode-trpc-compat': 'local-mcp-managed-action' },
    });
  }

  if (procedureName === 'mcpServers.update') {
    const nextName = typeof (input as { name?: unknown }).name === 'string' && String((input as { name?: unknown }).name).trim()
      ? String((input as { name?: unknown }).name).trim()
      : match.name;
    const existingMeta = buildLocalServerMeta(match.name, match.server);
    const updatedServer: LocalMcpServerEntry = {
      ...match.server,
      description: typeof (input as { description?: unknown }).description === 'string'
        ? (input as { description?: string }).description
        : (input as { description?: null }).description === null
          ? null
          : match.server.description ?? null,
      type: normalizeImportedServerType({
        type: typeof (input as { type?: unknown }).type === 'string'
          ? String((input as { type?: unknown }).type)
          : match.server.type,
        url: typeof (input as { url?: unknown }).url === 'string'
          ? String((input as { url?: unknown }).url)
          : match.server.url,
      }),
      command: typeof (input as { command?: unknown }).command === 'string'
        ? (input as { command?: string }).command
        : (input as { command?: null }).command === null
          ? undefined
          : match.server.command,
      args: Array.isArray((input as { args?: unknown[] }).args)
        ? (input as { args?: unknown[] }).args?.filter((arg): arg is string => typeof arg === 'string')
        : match.server.args ?? [],
      env: (input as { env?: unknown }).env && typeof (input as { env?: unknown }).env === 'object'
        ? Object.fromEntries(Object.entries((input as { env?: Record<string, unknown> }).env ?? {}).map(([key, value]) => [key, String(value)]))
        : match.server.env ?? {},
      url: typeof (input as { url?: unknown }).url === 'string'
        ? (input as { url?: string }).url
        : (input as { url?: null }).url === null
          ? undefined
          : match.server.url,
      bearerToken: typeof (input as { bearerToken?: unknown }).bearerToken === 'string'
        ? (input as { bearerToken?: string }).bearerToken
        : (input as { bearerToken?: null }).bearerToken === null
          ? undefined
          : match.server.bearerToken,
      headers: (input as { headers?: unknown }).headers && typeof (input as { headers?: unknown }).headers === 'object'
        ? Object.fromEntries(Object.entries((input as { headers?: Record<string, unknown> }).headers ?? {}).map(([key, value]) => [key, String(value)]))
        : match.server.headers ?? {},
      _meta: existingMeta,
    };

    if (nextName !== match.name) {
      delete localConfig.mcpServers[match.name];
    }

    localConfig.mcpServers[nextName] = updatedServer;
    upsertLocalServerMeta(localConfig, nextName, { uuid: existingMeta.uuid });
    await writeLocalMcpConfig(localConfig);

    return buildTrpcResponse(req, buildLocalManagedServerRecord(nextName, localConfig.mcpServers[nextName]), {
      status: 200,
      headers: { 'x-hypercode-trpc-compat': 'local-mcp-managed-action' },
    });
  }

  if (procedureName === 'mcpServers.delete') {
    const deletedServer = buildLocalManagedServerRecord(match.name, match.server);
    delete localConfig.mcpServers[match.name];
    await writeLocalMcpConfig(localConfig);

    return buildTrpcResponse(req, deletedServer, {
      status: 200,
      headers: { 'x-hypercode-trpc-compat': 'local-mcp-managed-action' },
    });
  }

  if (procedureName === 'mcpServers.reloadMetadata') {
    const mode = typeof (input as { mode?: unknown }).mode === 'string' ? String((input as { mode?: unknown }).mode) : 'binary';
    const metadata = upsertLocalServerMeta(localConfig, match.name, {
      status: 'ready',
      metadataSource: `local-${mode}`,
      toolCount: sanitizeNumericMetadata(match.server._meta?.toolCount, 0),
      lastSuccessfulBinaryLoadAt: new Date().toISOString(),
    });
    await writeLocalMcpConfig(localConfig);

    return buildTrpcResponse(req, {
      server: buildLocalManagedServerRecord(match.name, localConfig.mcpServers[match.name]),
      metadata,
    }, {
      status: 200,
      headers: { 'x-hypercode-trpc-compat': 'local-mcp-managed-action' },
    });
  }

  if (procedureName === 'mcpServers.clearMetadataCache') {
    const metadata = upsertLocalServerMeta(localConfig, match.name, {
      status: match.server.disabled ? 'disabled' : 'pending',
      metadataSource: LOCAL_COMPAT_METADATA_SOURCE,
      toolCount: 0,
      lastSuccessfulBinaryLoadAt: null,
    });
    await writeLocalMcpConfig(localConfig);

    return buildTrpcResponse(req, {
      server: buildLocalManagedServerRecord(match.name, localConfig.mcpServers[match.name]),
      metadata,
    }, {
      status: 200,
      headers: { 'x-hypercode-trpc-compat': 'local-mcp-managed-action' },
    });
  }

  if (procedureName === 'serverHealth.reset') {
    upsertLocalServerMeta(localConfig, match.name, {
      crashCount: 0,
      maxAttempts: 0,
    });
    await writeLocalMcpConfig(localConfig);

    return buildTrpcResponse(req, {
      status: 'unavailable',
      crashCount: 0,
      maxAttempts: 0,
    }, {
      status: 200,
      headers: { 'x-hypercode-trpc-compat': 'local-mcp-managed-action' },
    });
  }

  return null;
}

async function handler(req: Request): Promise<Response> {
  const upstreamBases = resolveUpstreamBases();
  const headers = cloneHeaders(req);
  const hasBody = req.method !== 'GET' && req.method !== 'HEAD';
  const body = hasBody ? await req.text() : undefined;
  const upstreamBody = normalizeBulkImportProxyBody(req, body);

  let upstreamResponse: Response | null = null;
  let lastUpstreamUrl = '';
  let lastError: unknown;
  let saw404 = false;

  for (const upstreamBase of upstreamBases) {
    const upstreamUrl = buildUpstreamUrl(req, upstreamBase);
    lastUpstreamUrl = upstreamUrl.toString();

    try {
      const response = await fetch(upstreamUrl, {
        method: req.method,
        headers,
        body: upstreamBody,
      });

      if (response.status === 404) {
        saw404 = true;
        continue;
      }

      if (response.status >= 500 && canFallbackLocally(req)) {
        continue;
      }

      if (isLegacyMcpRequest(req) && (response.status === 401 || response.status === 403)) {
        continue;
      }

      upstreamResponse = response;
      break;
    } catch (error) {
      lastError = error;
    }
  }

  if (!upstreamResponse) {
    const bulkImportBridgeResponse = await tryBridgeBulkImport(req, upstreamBases, headers, body);
    if (bulkImportBridgeResponse) {
      return bulkImportBridgeResponse;
    }

    const localManagedMutationResponse = await tryLocalManagedServerMutation(req, body);
    if (localManagedMutationResponse) {
      return localManagedMutationResponse;
    }

    const localBulkImportResponse = await tryLocalBulkImport(req, body);
    if (localBulkImportResponse) {
      return localBulkImportResponse;
    }

    const bridgeResponse = await tryResolveLegacyMcpResponse(req, upstreamBases, headers);
    if (bridgeResponse) {
      return bridgeResponse;
    }

    const compatResponse = await buildLocalCompatResponse(req, body);
    if (compatResponse) {
      return compatResponse;
    }

    const message = saw404
      ? 'No configured tRPC upstream exposed the requested procedure path'
      : lastError instanceof Error
        ? lastError.message
        : String(lastError ?? 'No upstream responded');
    return new Response(
      JSON.stringify({
        error: 'TRPC_UPSTREAM_UNAVAILABLE',
        message,
        upstream: lastUpstreamUrl,
      }),
      {
        status: 502,
        headers: { 'content-type': 'application/json' },
      },
    );
  }

  if (await shouldPreferLocalCompat(req, upstreamResponse)) {
    const compatResponse = await buildLocalCompatResponse(req, body);
    if (compatResponse) {
      return compatResponse;
    }
  }

  // Detect and handle SSE (Server-Sent Events) for subscriptions
  const isSse = upstreamResponse.headers.get('content-type')?.includes('text/event-stream');
  
  const responseHeaders = new Headers(upstreamResponse.headers);
  if (isSse) {
    responseHeaders.set('Connection', 'keep-alive');
    responseHeaders.set('Cache-Control', 'no-cache, no-transform');
    responseHeaders.set('X-Accel-Buffering', 'no');
  }

  return new Response(upstreamResponse.body, {
    status: upstreamResponse.status,
    statusText: upstreamResponse.statusText,
    headers: responseHeaders,
  });
}

export { handler as GET, handler as POST };
