import type { RegisteredBridgeClient } from '../bridge/bridge-manifest.js';
import type { ExecutionEnvironmentSummary } from '../services/execution-environment.js';

type StartupStatusInput = {
    mcpServer: unknown;
    aggregator?: {
        getInitializationStatus?: () => {
            inProgress: boolean;
            initialized: boolean;
            lastStartedAt?: number;
            lastCompletedAt?: number;
            lastSuccessAt?: number;
            lastError?: string;
            connectedClientCount: number;
            configuredServerCount: number;
        };
    } | null;
    agentMemory?: unknown;
    browserService?: unknown;
    browserStatus: {
        active?: boolean;
        pageCount?: number;
        pageIds?: string[];
    };
    sessionSupervisor?: {
        getRestoreStatus?: () => {
            lastRestoreAt?: number;
            restoredSessionCount: number;
            autoResumeCount: number;
        };
    } | null;
    sessionCount: number;
    mcpConfigService?: {
        getStatus?: () => {
            inProgress: boolean;
            lastStartedAt?: number;
            lastCompletedAt?: number;
            lastSuccessAt?: number;
            lastError?: string;
            lastServerCount: number;
            lastToolCount: number;
        };
    } | null;
    liveServerCount: number;
    persistedServerCount: number;
    persistedToolCount: number;
    persistedAlwaysOnServerCount: number;
    persistedAlwaysOnToolCount: number;
    executionEnvironment?: ExecutionEnvironmentSummary | null;
};

export async function buildStartupStatusSnapshot(input: StartupStatusInput) {
    const {
        mcpServer,
        aggregator,
        agentMemory,
        browserService,
        browserStatus,
        sessionSupervisor,
        sessionCount,
        mcpConfigService,
        liveServerCount,
        persistedServerCount,
        persistedToolCount,
        persistedAlwaysOnServerCount,
        persistedAlwaysOnToolCount,
        executionEnvironment,
    } = input;
    const mcpServerRuntime = mcpServer as {
        memoryManager?: unknown;
        isMemoryInitialized?: boolean;
        getBridgeStatus?: () => {
            ready?: boolean;
            clientCount?: number;
            clients?: RegisteredBridgeClient[];
            supportedCapabilities?: string[];
            supportedHookPhases?: string[];
        };
    };

    const aggregatorStatus = aggregator?.getInitializationStatus?.();
    const configSyncStatus = mcpConfigService?.getStatus?.();
    const restoreStatus = sessionSupervisor?.getRestoreStatus?.();
    const bridgeStatus = mcpServerRuntime.getBridgeStatus?.() ?? {
        ready: false,
        clientCount: 0,
        clients: [],
        supportedCapabilities: [],
        supportedHookPhases: [],
    };

    const bridgeClientCount = Number(bridgeStatus.clientCount ?? 0);
    const bridgeReady = Boolean(bridgeStatus.ready);
    const hasConnectedBridgeClients = bridgeClientCount > 0;
    const memoryInitialized = Boolean(mcpServerRuntime.isMemoryInitialized);
    const memoryReady = Boolean(mcpServerRuntime.memoryManager && agentMemory && memoryInitialized);
    const sessionReady = Boolean(sessionSupervisor);
    const browserReady = Boolean(browserService);
    const mcpReady = Boolean(aggregator);
    const configSyncReady = Boolean(configSyncStatus?.lastCompletedAt) && !configSyncStatus?.inProgress && !configSyncStatus?.lastError;
    const sessionRestoreReady = Boolean(restoreStatus?.lastRestoreAt);
    const configuredServerCount = Number(
        configSyncStatus?.lastServerCount
        ?? aggregatorStatus?.configuredServerCount
        ?? persistedServerCount,
    );
    const knownServerCount = Math.max(liveServerCount, persistedServerCount, configuredServerCount);
    const advertisedServerCount = Math.max(persistedServerCount, configuredServerCount);
    const advertisedToolCount = Math.max(persistedToolCount, Number(configSyncStatus?.lastToolCount ?? 0));
    const inventoryKnown = Boolean(configSyncStatus?.lastCompletedAt || aggregatorStatus?.initialized);
    const inventoryReady = inventoryKnown && (
        configuredServerCount === 0
        || persistedServerCount > 0
        || persistedToolCount > 0
        || liveServerCount > 0
    );
    const aggregatorReady = Boolean(aggregatorStatus?.initialized);
    const liveReady = mcpReady && aggregatorReady;
    const mcpWarmupInProgress = Boolean(
        aggregatorStatus?.inProgress
        || configSyncStatus?.inProgress
        || (inventoryReady && configuredServerCount > 0 && liveServerCount < configuredServerCount),
    );
    const executionReady = Boolean(executionEnvironment?.ready);

    return {
        status: 'running',
        ready: liveReady
            && memoryReady
            && browserReady
            && sessionReady
            && bridgeReady
            && configSyncReady
            && sessionRestoreReady
            && inventoryReady
            && aggregatorReady,
        uptime: process.uptime(),
        checks: {
            mcpAggregator: {
                ready: liveReady,
                liveReady,
                serverCount: knownServerCount,
                connectedCount: liveServerCount,
                initialization: aggregatorStatus ?? null,
                persistedServerCount,
                persistedToolCount,
                advertisedServerCount,
                advertisedToolCount,
                advertisedAlwaysOnServerCount: persistedAlwaysOnServerCount,
                advertisedAlwaysOnToolCount: persistedAlwaysOnToolCount,
                configuredServerCount,
                inventoryReady,
                warmupInProgress: mcpWarmupInProgress,
            },
            configSync: {
                ready: configSyncReady,
                status: configSyncStatus ?? null,
            },
            memory: {
                ready: memoryReady,
                initialized: memoryInitialized,
                agentMemory: Boolean(agentMemory),
            },
            browser: {
                ready: browserReady,
                active: Boolean(browserStatus.active),
                pageCount: Number(browserStatus.pageCount ?? 0),
            },
            sessionSupervisor: {
                ready: sessionReady && sessionRestoreReady,
                sessionCount,
                restore: restoreStatus ?? null,
            },
            extensionBridge: {
                ready: bridgeReady,
                acceptingConnections: bridgeReady,
                clientCount: bridgeClientCount,
                hasConnectedClients: hasConnectedBridgeClients,
                clients: bridgeStatus.clients ?? [],
                supportedCapabilities: bridgeStatus.supportedCapabilities ?? [],
                supportedHookPhases: bridgeStatus.supportedHookPhases ?? [],
            },
            executionEnvironment: {
                ready: executionReady,
                preferredShellId: executionEnvironment?.preferredShellId ?? null,
                preferredShellLabel: executionEnvironment?.preferredShellLabel ?? null,
                shellCount: Number(executionEnvironment?.shellCount ?? 0),
                verifiedShellCount: Number(executionEnvironment?.verifiedShellCount ?? 0),
                toolCount: Number(executionEnvironment?.toolCount ?? 0),
                verifiedToolCount: Number(executionEnvironment?.verifiedToolCount ?? 0),
                harnessCount: Number(executionEnvironment?.harnessCount ?? 0),
                verifiedHarnessCount: Number(executionEnvironment?.verifiedHarnessCount ?? 0),
                supportsPowerShell: Boolean(executionEnvironment?.supportsPowerShell),
                supportsPosixShell: Boolean(executionEnvironment?.supportsPosixShell),
                notes: executionEnvironment?.notes ?? [],
            },
        },
    };
}