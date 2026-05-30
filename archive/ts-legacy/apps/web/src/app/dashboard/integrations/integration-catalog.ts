export type SyncTargetSummary = {
    client: string;
    path: string;
    candidates: string[];
    exists: boolean;
};

export type CliHarnessDetectionSummary = {
    installed?: boolean;
};

export type StartupStatusSummary = {
    status?: string;
    summary?: string;
export function getExternalClientRows(syncTargets?: SyncTargetSummary[] | null): ExternalClientRow[] {
    const targetMap = new Map((syncTargets ?? []).map((target) => [target.client, target]));

    return EXTERNAL_CLIENTS.map((client) => {
        const target = client.syncClient ? targetMap.get(client.syncClient) : undefined;
        const detected = Boolean(target?.exists);

        return {
            id: client.id,
            label: client.label,
            windowsPath: client.windowsPath,
            notes: client.notes,
            autoSyncSupported: client.autoSyncSupported,
            detected,
            resolvedPath: target?.path ?? client.windowsPath,
            statusLabel: detected
                ? 'Detected'
                : client.autoSyncSupported
                    ? 'Ready to sync'
                    : 'Reference path',
            statusTone: detected
                ? 'success'
                : client.autoSyncSupported
                    ? 'warning'
                    : 'muted',
        };
    });
}

export function getStatusBadgeClasses(tone: 'success' | 'warning' | 'muted'): string {
    if (tone === 'success') {
        return 'border-emerald-500/30 bg-emerald-500/10 text-emerald-300';
    }

    if (tone === 'warning') {
        return 'border-amber-500/30 bg-amber-500/10 text-amber-300';
    }

    return 'border-zinc-700 bg-zinc-900 text-zinc-300';
}

function formatRelativeTime(timestamp: number, now = Date.now()): string {
    const deltaMs = Math.max(0, now - timestamp);
    const deltaSeconds = Math.floor(deltaMs / 1000);

    if (deltaSeconds < 5) {
        return 'just now';
    }

    if (deltaSeconds < 60) {
        return `${deltaSeconds}s ago`;
    }

    const deltaMinutes = Math.floor(deltaSeconds / 60);
    if (deltaMinutes < 60) {
        return `${deltaMinutes}m ago`;
    }

    const deltaHours = Math.floor(deltaMinutes / 60);
    if (deltaHours < 24) {
        return `${deltaHours}h ago`;
    }

    return `${Math.floor(deltaHours / 24)}d ago`;
}

export function getConnectedBridgeClientRows(startupStatus?: StartupStatusSummary | null): ConnectedBridgeClientRow[] {
    const clients = startupStatus?.checks?.extensionBridge?.clients ?? [];

    return [...clients]
        .sort((left, right) => (left.clientName ?? left.clientId ?? '').localeCompare(right.clientName ?? right.clientId ?? ''))
        .map((client) => ({
            clientId: client.clientId ?? 'unknown-client',
            clientName: client.clientName ?? client.clientId ?? 'Unknown bridge client',
            clientType: client.clientType ?? 'unknown',
            version: client.version,
            platform: client.platform,
            capabilities: client.capabilities ?? [],
            hookPhases: client.hookPhases ?? [],
            lastSeenLabel: typeof client.lastSeenAt === 'number' ? formatRelativeTime(client.lastSeenAt) : 'unknown',
        }));
}