import { describe, expect, it } from 'vitest';

import type { DashboardStartupStatus } from '../../dashboard-home-view';
import { buildSystemStartupChecks } from './system-status-helpers';

type StartupCheckOverrides = {
    [Key in keyof DashboardStartupStatus['checks']]?: Partial<DashboardStartupStatus['checks'][Key]>;
};

function createStartupStatus(overrides?: StartupCheckOverrides): DashboardStartupStatus {
    return {
        status: 'running',
        ready: true,
        uptime: 42,
        checks: {
            mcpAggregator: {
                ready: true,
                liveReady: true,
                serverCount: 0,
                connectedCount: 0,
                initialization: {
                    inProgress: false,
                    initialized: true,
                    connectedClientCount: 0,
                    configuredServerCount: 0,
                },
                persistedServerCount: 0,
                persistedToolCount: 0,
                configuredServerCount: 0,
                advertisedServerCount: 0,
                advertisedToolCount: 0,
                advertisedAlwaysOnToolCount: 0,
                inventoryReady: true,
                ...overrides?.mcpAggregator,
            },
            configSync: {
                ready: true,
                status: {
                    inProgress: false,
                    lastCompletedAt: 1_700_000_000_000,
                    lastServerCount: 0,
                    lastToolCount: 0,
                },
                ...overrides?.configSync,
            },
            memory: {
                ready: true,
                initialized: true,
                agentMemory: true,
                ...overrides?.memory,
            },
            browser: {
                ready: true,
                active: false,
                pageCount: 0,
                ...overrides?.browser,
            },
            sessionSupervisor: {
                ready: true,
                sessionCount: 0,
                restore: {
                    restoredSessionCount: 0,
                    autoResumeCount: 0,
                },
                ...overrides?.sessionSupervisor,
            },
            extensionBridge: {
                ready: true,
                acceptingConnections: true,
                clientCount: 0,
                hasConnectedClients: false,
                ...overrides?.extensionBridge,
            },
            executionEnvironment: {
                ready: true,
                preferredShellLabel: 'PowerShell 7',
                shellCount: 2,
                verifiedShellCount: 2,
                toolCount: 5,
                verifiedToolCount: 5,
                ...overrides?.executionEnvironment,
            },
        },
    };
}

describe('system status startup helpers', () => {
    it('treats an empty but initialized router inventory as operational', () => {
        const checks = buildSystemStartupChecks(createStartupStatus());

        expect(checks[1]).toEqual({
            name: 'Live MCP Runtime',
            status: 'Operational',
            latency: '0/0 servers',
            detail: 'No downstream servers configured · live MCP runtime is ready',
        });
    });

    it('keeps router inventory pending until the router itself is ready', () => {
        const checks = buildSystemStartupChecks(createStartupStatus({
            mcpAggregator: {
                ready: false,
                liveReady: false,
                inventoryReady: true,
                persistedServerCount: 2,
                persistedToolCount: 18,
                advertisedServerCount: 2,
                advertisedToolCount: 18,
            },
        }));

        expect(checks[0]).toEqual({
            name: 'Cached Inventory',
            status: 'Operational',
            latency: '18 tools',
            detail: '2 cached servers · 18 advertised tools',
        });

        expect(checks[1]).toEqual({
            name: 'Live MCP Runtime',
            status: 'Pending',
            latency: '0/2 servers',
            detail: 'Cached inventory is already advertised · live MCP runtime is still warming',
        });
    });

    it('shows the bridge listener as operational while idle before clients attach', () => {
        const checks = buildSystemStartupChecks(createStartupStatus({
            extensionBridge: {
                ready: true,
                acceptingConnections: true,
                clientCount: 0,
                hasConnectedClients: false,
            },
        }));

        expect(checks[3]).toEqual({
            name: 'Session Restore',
            status: 'Operational',
            latency: '0 sessions',
            detail: '0 restored · 0 auto-resumed',
        });

        expect(checks[4]).toEqual({
            name: 'Client Bridge',
            status: 'Operational',
            latency: '0 clients',
            detail: 'Browser/editor client bridge is ready, but no IDE or browser adapters have connected yet.',
        });
    });

    it('shows the bridge listener as pending while the listener is still booting', () => {
        const checks = buildSystemStartupChecks(createStartupStatus({
            extensionBridge: {
                ready: false,
                acceptingConnections: false,
                clientCount: 0,
                hasConnectedClients: false,
            },
        }));

        expect(checks[3]).toEqual({
            name: 'Session Restore',
            status: 'Operational',
            latency: '0 sessions',
            detail: '0 restored · 0 auto-resumed',
        });

        expect(checks[4]).toEqual({
            name: 'Client Bridge',
            status: 'Pending',
            latency: '0 clients',
            detail: 'Browser/editor client bridge is still coming online',
        });
    });

    it('shows the execution environment posture with preferred shell context', () => {
        const checks = buildSystemStartupChecks(createStartupStatus({
            executionEnvironment: {
                ready: true,
                preferredShellLabel: 'PowerShell 7',
                toolCount: 6,
                verifiedToolCount: 5,
            },
        }));

        expect(checks[5]).toEqual({
            name: 'Execution Environment',
            status: 'Operational',
            latency: '5 tools',
            detail: 'PowerShell 7 preferred · 5/6 verified tools',
        });
    });
});