export type ManagedServerDiscoveryRecord = {
    uuid?: string;
    name: string;
    _meta?: {
        status?: string | null;
        metadataSource?: string | null;
        toolCount?: number | null;
        lastSuccessfulBinaryLoadAt?: string | null;
    } | null;
    always_on?: boolean;
};

export type BulkMetadataTargetMode = 'all' | 'unresolved';

export type ManagedServerDiscoverySummary = {
    totalCount: number;
    readyCount: number;
    unresolvedCount: number;
    staleReadyCount: number;
    repairableCount: number;
    neverLoadedCount: number;
    localCompatCount: number;
};

export type ServerToolActionLinks = {
    inspectToolsHref: string;
    editToolsHref: string;
    logsHref: string;
};

export type RuntimeServerRecord = {
    name: string;
    status: string;
    toolCount: number;
    runtimeState?: string;
    warmupState?: string;
    runtimeConnected?: boolean;
    advertisedToolCount?: number;
    advertisedSource?: string;
    lastConnectedAt?: string | null;
    lastError?: string | null;
    config?: {
        command?: string;
        args?: string[];
        env?: string[];
    };
};

export type ManagedServerRuntimeRecord = ManagedServerDiscoveryRecord & {
    description?: string | null;
    command?: string | null;
    args?: string[];
    env?: Record<string, string>;
    source_published_server_uuid?: string | null;
    _meta?: (ManagedServerDiscoveryRecord['_meta'] & {
        metadataSource?: string | null;
        toolCount?: number | null;
    }) | null;
};

export type DashboardServerRecord = {
    uuid?: string;
    name: string;
    status: string;
    toolCount: number;
    runtimeState?: string;
    warmupState?: string;
    runtimeConnected?: boolean;
    advertisedToolCount?: number;
    advertisedSource?: string;
    lastConnectedAt?: string | null;
    lastError?: string | null;
    metadataStatus?: string;
    metadataSource?: string;
    metadataToolCount?: number;
    lastSuccessfulBinaryLoadAt?: string;
    always_on?: boolean;
    source_published_server_uuid?: string | null;
            config: {
                command: runtime?.config?.command ?? server.command ?? undefined,
                args: runtime?.config?.args ?? server.args ?? [],
                env: runtime?.config?.env ?? Object.keys(server.env ?? {}),
            },
        };
    });

    const runtimeOnly = runtimeServers
        .filter((server) => !matchedRuntimeNames.has(normalizeServerName(server.name)))
        .map<DashboardServerRecord>((server) => ({
            name: server.name,
            status: server.status,
            toolCount: server.toolCount,
            runtimeState: server.runtimeState ?? server.status,
            warmupState: server.warmupState ?? 'idle',
            runtimeConnected: server.runtimeConnected ?? (server.status === 'connected'),
            advertisedToolCount: server.advertisedToolCount,
            advertisedSource: server.advertisedSource,
            lastConnectedAt: server.lastConnectedAt ?? null,
            lastError: server.lastError ?? null,
            config: server.config,
        }));

    return [...managedFirst, ...runtimeOnly];
}

export function buildServerToolActionLinks(serverName: string): ServerToolActionLinks {
    const encodedServerName = encodeURIComponent(serverName);

    return {
        inspectToolsHref: `/dashboard/mcp/inspector?server=${encodedServerName}`,
        editToolsHref: `/dashboard/mcp/inspector?server=${encodedServerName}&mode=edit-tools`,
        logsHref: `/dashboard/mcp/logs?server=${encodedServerName}`,
    };
}
