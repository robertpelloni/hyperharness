import { mcpServersRepository, toolsRepository } from '../db/repositories/index.js';

import { deriveSemanticCatalogForServer } from './catalogMetadata.js';

export type CachedMcpServerInventory = Awaited<ReturnType<typeof mcpServersRepository.findAll>>[number] & {
    displayName: string;
    tags: string[];
    alwaysOnAdvertised: boolean;
};

export type CachedMcpToolInventory = {
    name: string;
    description: string;
    server: string;
    serverDisplayName: string;
    serverTags: string[];
    toolTags: string[];
    semanticGroup: string;
    semanticGroupLabel: string;
    advertisedName: string;
    keywords: string[];
    alwaysOn: boolean;
    originalName: string;
    inputSchema: unknown;
};

export type CachedMcpInventorySource = 'database' | 'config' | 'empty';

type CachedMcpInventorySnapshot = {
    servers: CachedMcpServerInventory[];
    toolCounts: Map<string, number>;
    tools: CachedMcpToolInventory[];
    source: CachedMcpInventorySource;
    snapshotUpdatedAt: string | null;
    databaseAvailable: boolean;
    databaseError: string | null;
    fallbackUsed: boolean;
=======
function buildConfigSnapshot(configServers: Record<string, BorgMcpServerEntry>): CachedMcpInventorySnapshot {
>>>>>>> origin/rewrite/main-sanitized:packages/core/src/mcp/cachedToolInventory.ts
    const servers = Object.entries(configServers).map(([name, server], index) => {
        const metadata = server._meta;
        return {
            uuid: `config:${name}:${index}`,
            name,
            type: server.type ?? (server.url ? 'SSE' : 'STDIO'),
            command: server.command ?? '',
            args: server.args ?? [],
            env: server.env ?? {},
            bearerToken: null,
            headers: {},
            url: server.url ?? null,
            description: metadata?.description ?? server.description ?? null,
            created_at: toDateOrNull(metadata?.discoveredAt),
            updated_at: toDateOrNull(metadata?.cacheHydratedAt ?? metadata?.discoveredAt),
            error_status: metadata?.status === 'failed' ? 'ERROR' : 'NONE',
            error_message: metadata?.error ?? null,
            enabled: !server.disabled,
            always_on: Boolean(metadata?.alwaysOn),
            user_id: null,
            source_published_server_uuid: null,
=======
            const typedTool = tool as BorgMcpToolMetadata;
>>>>>>> origin/rewrite/main-sanitized:packages/core/src/mcp/cachedToolInventory.ts
            tools.push({
                name: namespaceToolName(server.name, typedTool.name),
                description: typedTool.description ?? '',
                server: server.name,
                serverDisplayName: typedTool.serverDisplayName ?? server.displayName ?? server.name,
                serverTags: typedTool.serverTags ?? server.tags ?? [],
                toolTags: typedTool.toolTags ?? [],
                semanticGroup: typedTool.semanticGroup ?? 'general-utility',
                semanticGroupLabel: typedTool.semanticGroupLabel ?? 'general utility',
                advertisedName: typedTool.advertisedName ?? namespaceToolName(server.name, typedTool.name),
                keywords: typedTool.keywords ?? [],
                alwaysOn: Boolean(typedTool.alwaysOn ?? server.alwaysOnAdvertised),
                originalName: typedTool.name,
                inputSchema: typedTool.inputSchema ?? null,
            });
        });
    });

    return {
        servers,
        toolCounts,
        tools,
        source: tools.length > 0 || servers.length > 0 ? 'config' : 'empty',
        snapshotUpdatedAt: toSnapshotUpdatedAt(updatedAtCandidates),
        databaseAvailable: true,
        databaseError: null,
        fallbackUsed: false,
    };
}

export async function getCachedToolInventory() {
    const configSnapshot = buildConfigSnapshot(config.mcpServers ?? {});

    try {
        const databaseSnapshot = await buildDatabaseSnapshot();
        
        // MERGE both sources so they don't cross-contaminate but both are available
        const mergedServers = [...configSnapshot.servers, ...databaseSnapshot.servers];
        const mergedTools = [...configSnapshot.tools, ...databaseSnapshot.tools];
        const mergedToolCounts = new Map<string, number>();
        
        for (const [key, val] of configSnapshot.toolCounts.entries()) mergedToolCounts.set(key, val);
        for (const [key, val] of databaseSnapshot.toolCounts.entries()) mergedToolCounts.set(key, val);

        const source: CachedMcpInventorySource = mergedTools.length > 0 ? 'database' : 'empty';

        return {
            servers: mergedServers,
            toolCounts: mergedToolCounts,
            tools: mergedTools,
            source,
            snapshotUpdatedAt: databaseSnapshot.snapshotUpdatedAt ?? configSnapshot.snapshotUpdatedAt,
            databaseAvailable: true,
            databaseError: null,
            fallbackUsed: false,
        };
    } catch (error) {
        return {
            ...configSnapshot,
            databaseAvailable: false,
            databaseError: formatOptionalSqliteFailure('Persisted MCP inventory is unavailable', error),
            fallbackUsed: true,
        };
    }
}
