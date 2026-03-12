export interface SystemStartupStatusInput {
    checks?: {
        mcpAggregator?: {
            ready?: boolean;
            liveReady?: boolean;
            connectedCount?: number;
            persistedServerCount?: number;
            persistedToolCount?: number;
            configuredServerCount?: number;
            advertisedServerCount?: number;
            advertisedToolCount?: number;
            advertisedAlwaysOnToolCount?: number;
            inventoryReady?: boolean;
            warmupInProgress?: boolean;
        };
        configSync?: {
            ready?: boolean;
            status?: {
                inProgress?: boolean;
                lastServerCount?: number;
                lastToolCount?: number;
            } | null;
        };
        sessionSupervisor?: {
            ready?: boolean;
            sessionCount?: number;
            restore?: {
                restoredSessionCount?: number;
                autoResumeCount?: number;
            } | null;
        };
        memory?: {
            ready?: boolean;
            initialized?: boolean;
        };
        extensionBridge?: {
            ready?: boolean;
            acceptingConnections?: boolean;
            clientCount?: number;
            hasConnectedClients?: boolean;
        };
        executionEnvironment?: {
            ready?: boolean;
            preferredShellLabel?: string | null;
            shellCount?: number;
            verifiedShellCount?: number;
            toolCount?: number;
            verifiedToolCount?: number;
        };
    };
}

export interface SystemStartupCheckRow {
    name: string;
    status: 'Operational' | 'Pending';
    latency: string;
    detail: string;
}

function getRouterInventoryDetail(startupStatus: SystemStartupStatusInput): string {
    const aggregator = startupStatus.checks?.mcpAggregator;
    const persistedServerCount = aggregator?.advertisedServerCount ?? aggregator?.persistedServerCount ?? 0;
    const persistedToolCount = aggregator?.advertisedToolCount ?? aggregator?.persistedToolCount ?? 0;
    const alwaysOnToolCount = aggregator?.advertisedAlwaysOnToolCount ?? 0;

    if (aggregator?.ready && aggregator.inventoryReady && persistedServerCount === 0 && persistedToolCount === 0) {
        return 'No configured servers yet · empty cached inventory is ready';
    }

    if (aggregator?.inventoryReady) {
        const alwaysOnSuffix = alwaysOnToolCount > 0
            ? ` · ${alwaysOnToolCount} always-on advertised immediately`
            : '';
        return `${persistedServerCount} cached servers · ${persistedToolCount} advertised tools${alwaysOnSuffix}`;
    }

    return 'Waiting for the first cached MCP inventory snapshot';
}

function getLiveMcpDetail(startupStatus: SystemStartupStatusInput): string {
    const aggregator = startupStatus.checks?.mcpAggregator;
    const configuredServerCount = Math.max(
        aggregator?.configuredServerCount ?? 0,
        aggregator?.advertisedServerCount ?? aggregator?.persistedServerCount ?? 0,
    );
    const connectedCount = aggregator?.connectedCount ?? 0;
    const liveReady = aggregator?.liveReady ?? aggregator?.ready;

    if (liveReady && configuredServerCount === 0) {
        return 'No downstream servers configured · live MCP runtime is ready';
    }

    if (liveReady) {
        if (configuredServerCount > 0 && connectedCount < configuredServerCount) {
            return `${connectedCount}/${configuredServerCount} live server connections warmed · cached tools stay usable while the rest connect`;
        }

        return `${connectedCount}/${configuredServerCount || connectedCount} live server connections ready`;
    }

    if (aggregator?.inventoryReady) {
        return 'Cached inventory is already advertised · live MCP runtime is still warming';
    }

    return 'Waiting for live MCP runtime initialization';
}

function getExtensionBridgeDetail(startupStatus: SystemStartupStatusInput): string {
    const extensionBridge = startupStatus.checks?.extensionBridge;
    const clientCount = extensionBridge?.clientCount ?? 0;
    const acceptingConnections = extensionBridge?.acceptingConnections ?? extensionBridge?.ready;
    const hasConnectedClients = extensionBridge?.hasConnectedClients ?? clientCount > 0;

    if (acceptingConnections) {
        if (hasConnectedClients) {
            return 'Browser/editor client bridge is accepting connections';
        }

        return 'Browser/editor client bridge is ready, but no IDE or browser adapters have connected yet.';
    }

    return 'Browser/editor client bridge is still coming online';
}

function getExecutionEnvironmentDetail(startupStatus: SystemStartupStatusInput): string {
    const execution = startupStatus.checks?.executionEnvironment;

    if (execution?.preferredShellLabel) {
        return `${execution.preferredShellLabel} preferred · ${execution.verifiedToolCount ?? 0}/${execution.toolCount ?? 0} verified tools`;
    }

    return `${execution?.verifiedShellCount ?? 0}/${execution?.shellCount ?? 0} verified shells`;
}

export function buildSystemStartupChecks(startupStatus: SystemStartupStatusInput): SystemStartupCheckRow[] {
    const aggregator = startupStatus.checks?.mcpAggregator;
    const restore = startupStatus.checks?.sessionSupervisor?.restore;
    const memory = startupStatus.checks?.memory;
    const extensionBridge = startupStatus.checks?.extensionBridge;
    const extensionBridgeOperational = extensionBridge?.acceptingConnections ?? extensionBridge?.ready;
    const persistedToolCount = aggregator?.advertisedToolCount ?? aggregator?.persistedToolCount ?? 0;
    const sessionCount = startupStatus.checks?.sessionSupervisor?.sessionCount ?? 0;
    const bridgeClientCount = extensionBridge?.clientCount ?? 0;
    const execution = startupStatus.checks?.executionEnvironment;
    const connectedCount = aggregator?.connectedCount ?? 0;
    const configuredCount = Math.max(aggregator?.configuredServerCount ?? 0, aggregator?.advertisedServerCount ?? aggregator?.persistedServerCount ?? 0);

    return [
        {
            name: 'Cached Inventory',
            status: aggregator?.inventoryReady ? 'Operational' : 'Pending',
            latency: `${persistedToolCount} tools`,
            detail: getRouterInventoryDetail(startupStatus),
        },
        {
            name: 'Live MCP Runtime',
            status: (aggregator?.liveReady ?? aggregator?.ready) ? 'Operational' : 'Pending',
            latency: `${connectedCount}/${configuredCount || connectedCount} servers`,
            detail: getLiveMcpDetail(startupStatus),
        },
        {
            name: 'Memory / Context',
            status: memory?.ready ? 'Operational' : 'Pending',
            latency: memory?.initialized ? 'initialized' : '-',
            detail: memory?.ready
                ? 'Memory manager initialized and context services are available'
                : memory?.initialized
                    ? 'Memory manager is present, but agent context wiring is still finishing'
                    : 'Memory initialization is still in progress',
        },
        {
            name: 'Session Restore',
            status: startupStatus.checks?.sessionSupervisor?.ready ? 'Operational' : 'Pending',
            latency: `${sessionCount} sessions`,
            detail: restore
                ? `${restore.restoredSessionCount ?? 0} restored · ${restore.autoResumeCount ?? 0} auto-resumed`
                : 'Restore not finished',
        },
        {
            name: 'Client Bridge',
            status: extensionBridgeOperational ? 'Operational' : 'Pending',
            latency: `${bridgeClientCount} clients`,
            detail: getExtensionBridgeDetail(startupStatus),
        },
        {
            name: 'Execution Environment',
            status: execution?.ready ? 'Operational' : 'Pending',
            latency: `${execution?.verifiedToolCount ?? 0} tools`,
            detail: getExecutionEnvironmentDetail(startupStatus),
        },
    ];
}