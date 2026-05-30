import Link from 'next/link';

import { MissionControlFunctionToggles } from './mission-control-function-toggles';
import SuggestionsPanel from '../../components/SuggestionsPanel';
import { SessionHandoffWidget } from '../../components/SessionHandoffWidget';
import { ContextHealthWidget } from '../../components/ContextHealthWidget';
import { NeuralPulse } from '../../components/NeuralPulse';
import { HyperCodeOrchestratorWidget } from '../../components/HyperCodeOrchestratorWidget';

export interface DashboardStatusSummary {
    initialized: boolean;
    serverCount: number;
    toolCount: number;
    connectedCount: number;
}

export interface DashboardStartupStatus {
    status: string;
    ready: boolean;
    uptime: number;
    summary?: string;
    blockingReasons?: Array<{
        code: string;
        detail: string;
    }>;
    startupMode?: {
        requestedRuntime?: string;
        activeRuntime?: string;
        requestedPort?: number;
        activePort?: number;
        portDecision?: string;
        portReason?: string;
        launchMode?: string;
        dashboardMode?: string;
        installDecision?: string;
        installReason?: string;
        buildDecision?: string;
        buildReason?: string;
        updatedAt?: string;
    } | null;
            href: '/dashboard/mcp/system',
            hrefLabel: 'Review startup status',
        });
    } else if (startupPendingCount > 0) {
        alerts.push({
            id: 'startup-pending',
            severity: startupStatus.ready ? 'info' : 'warning',
            title: startupStatus.ready ? 'Background startup checks still reporting pending' : 'Startup sequence is still warming up',
            detail: `${startupPendingCount} startup check${startupPendingCount === 1 ? '' : 's'} ${startupPendingCount === 1 ? 'is' : 'are'} not ready yet.`,
            href: '/dashboard',
            hrefLabel: 'Review startup readiness',
        });
    }

    if (options?.serversError) {
        alerts.push({
            id: 'server-inventory-unavailable',
            severity: 'warning',
            title: 'MCP server inventory is unavailable',
            detail: options.serversError,
            href: '/dashboard/mcp',
            hrefLabel: 'Open server health',
        });
    } else if (mcpStatus.initialized && servers.length === 0 && providers.length === 0) {
        alerts.push({
            id: 'first-run-setup',
            severity: 'info',
            title: 'Welcome to HyperCode! Let\'s get started. 🚀',
            detail: 'Your workspace is fresh. Start by configuring an AI Provider and connecting an MCP Server to give your models tools.',
            href: '/dashboard/providers',
            hrefLabel: 'Configure Providers',
        });
    } else if (mcpStatus.initialized && servers.length === 0) {
        alerts.push({
            id: 'no-mcp-servers',
            severity: 'info',
            title: 'No MCP Servers Connected',
            detail: 'Your models have no tools available. Add a server from the registry or sync your VS Code/Cursor configuration.',
            href: '/dashboard/integrations',
            hrefLabel: 'Add MCP Server',
        });
    } else if (mcpStatus.initialized && providers.length === 0) {
        alerts.push({
            id: 'no-providers',
            severity: 'info',
            title: 'No AI Providers Configured',
            detail: 'You need to configure an API key (Anthropic, OpenAI, Gemini) to run autonomous sessions.',
            href: '/dashboard/providers',
            hrefLabel: 'Configure Providers',
        });
    }

    if (options?.providersError) {
        alerts.push({
            id: 'provider-inventory-unavailable',
            severity: 'warning',
            title: 'Provider routing inventory is unavailable',
            detail: options.providersError,
            href: '/dashboard/billing',
            hrefLabel: 'Review providers',
        });
    } else if (degradedProviders > 0) {
        alerts.push({
            id: 'provider-degraded',
            severity: degradedProviders > 1 ? 'critical' : 'warning',
            title: 'Provider routing has degraded capacity',
            detail: `${degradedProviders} configured provider${degradedProviders === 1 ? '' : 's'} ${degradedProviders === 1 ? 'needs' : 'need'} attention before fallback narrows.`,
            href: '/dashboard/billing',
            hrefLabel: 'Review providers',
        });
    }

    if (options?.fallbackChainError) {
        alerts.push({
            id: 'fallback-chain-unavailable',
            severity: 'warning',
            title: 'Provider fallback chain is unavailable',
            detail: options.fallbackChainError,
            href: '/dashboard/billing',
            hrefLabel: 'Review providers',
        });
    }

    if (options?.sessionsError) {
        alerts.push({
            id: 'session-inventory-unavailable',
            severity: 'warning',
            title: 'Supervised session inventory is unavailable',
            detail: options.sessionsError,
            href: '/dashboard/session',
            hrefLabel: 'Open sessions',
        });
    } else if (erroredSessions > 0) {
        alerts.push({
            id: 'session-errors',
            severity: 'critical',
            title: 'Supervised sessions have failed',
            detail: `${erroredSessions} session${erroredSessions === 1 ? '' : 's'} ${erroredSessions === 1 ? 'is' : 'are'} in an error state and may need restart or log review.`,
            href: '/dashboard/session',
            hrefLabel: 'Open sessions',
        });
    }

    return alerts.sort((left, right) => {
        const order = { critical: 0, warning: 1, info: 2 } as const;
        return order[left.severity] - order[right.severity];
    });
}

function getServerTone(status: string): string {
    switch (status) {
        case 'connected':
            return 'border-emerald-500/30 bg-emerald-500/10 text-emerald-200';
        case 'connecting':
        case 'restarting':
            return 'border-amber-500/30 bg-amber-500/10 text-amber-200';
        case 'error':
            return 'border-rose-500/30 bg-rose-500/10 text-rose-200';
        default:
            return 'border-slate-500/30 bg-slate-500/10 text-slate-200';
    }
}

function getSessionTone(status: DashboardSessionSummary['status']): string {
    switch (status) {
        case 'running':
            return 'border-emerald-500/30 bg-emerald-500/10 text-emerald-200';
        case 'starting':
        case 'restarting':
            return 'border-amber-500/30 bg-amber-500/10 text-amber-200';
        case 'error':
            return 'border-rose-500/30 bg-rose-500/10 text-rose-200';
        default:
            return 'border-slate-500/30 bg-slate-500/10 text-slate-200';
    }
}

function getSessionStatusLabel(status: DashboardSessionSummary['status']): string {
    switch (status) {
        case 'running':
            return 'Running';
        case 'starting':
            return 'Starting';
        case 'restarting':
            return 'Restarting';
        case 'stopping':
            return 'Stopping';
        case 'stopped':
            return 'Stopped';
        case 'error':
            return 'Error';
        default:
            return 'Created';
    }
}

function getLatestSessionLogMessage(session: DashboardSessionSummary): string | null {
    if (session.logs.length === 0) {
        return null;
    }

    return [...session.logs]
        .sort((left, right) => right.timestamp - left.timestamp)[0]
        ?.message
        ?.trim() || null;
}

function getSessionRestartPolicyLabel(session: DashboardSessionSummary): string {
    return session.autoRestart === false ? 'Manual restart only' : 'Auto-restart enabled';
}

function getProviderTone(provider: DashboardProviderSummary): string {
    if (!provider.configured) {
        return 'border-slate-500/30 bg-slate-500/10 text-slate-200';
    }

    if (isProviderDegraded(provider)) {
        return 'border-rose-500/30 bg-rose-500/10 text-rose-200';
    }

    return 'border-emerald-500/30 bg-emerald-500/10 text-emerald-200';
}

function formatQuotaValue(value: number | null): string {
    if (value === null) {
        return '—';
    }

    return value.toLocaleString();
}

function formatFallbackLabel(entry: DashboardFallbackSummary): string {
    return entry.model ? `${entry.provider} · ${entry.model}` : entry.provider;
}

function buildStartupModeEntries(startupStatus: DashboardStartupStatus): Array<{
    label: string;
    value: string;
    detail?: string;
}> {
    const startupMode = startupStatus.startupMode;
    if (!startupMode) {
        return [];
    }

    return [
        {
            label: 'Requested runtime',
            value: startupMode.requestedRuntime?.trim() || '—',
            detail: startupMode.activeRuntime ? `Active runtime: ${startupMode.activeRuntime}` : undefined,
        },
        {
            label: 'Launch mode',
            value: startupMode.launchMode?.trim() || '—',
            detail: startupMode.dashboardMode ? `Dashboard: ${startupMode.dashboardMode}` : undefined,
        },
        {
            label: 'Control-plane port',
            value: typeof startupMode.activePort === 'number' ? String(startupMode.activePort) : '—',
            detail: [
                typeof startupMode.requestedPort === 'number' ? `Requested: ${startupMode.requestedPort}` : null,
                startupMode.portDecision?.trim() || null,
                startupMode.portReason?.trim() || null,
            ].filter(Boolean).join(' • ') || undefined,
        },
        {
            label: 'Install decision',
            value: startupMode.installDecision?.trim() || '—',
            detail: startupMode.installReason?.trim() || undefined,
        },
        {
            label: 'Build decision',
            value: startupMode.buildDecision?.trim() || '—',
            detail: startupMode.buildReason?.trim() || undefined,
        },
    ];
}

    generatedAtLabel,
    currentTimestamp,
    isBootstrapping = false,
    mcpStatus,
    mcpStatusError,
    startupStatus,
    startupStatusError,
    servers,
    serversError,
    traffic,
    trafficError,
    providers,
    providersError,
    fallbackChain,
    fallbackChainError,
    sessions,
    sessionsError,
    installSurfaceArtifacts,
    onStartSession,
    onStopSession,
    onRestartSession,
    pendingSessionActionId,
}: DashboardHomeViewProps) {
    const overviewMetrics = buildOverviewMetrics(mcpStatus, sessions, providers, isBootstrapping, {
        mcpStatusError,
        sessionsError,
        providersError,
    });
    const startupChecklist = startupStatusError ? [] : buildStartupChecklist(startupStatus, isBootstrapping, installSurfaceArtifacts);
    const startupModeEntries = startupStatusError ? [] : buildStartupModeEntries(startupStatus);
    const startupModeUpdatedAt = startupStatus.startupMode?.updatedAt ? Date.parse(startupStatus.startupMode.updatedAt) : Number.NaN;
                            </div>
                        ) : dashboardAlerts.length === 0 ? (
                            <div className="mt-4 rounded-2xl border border-emerald-500/30 bg-emerald-500/10 p-4 text-sm text-emerald-100">
                                All major systems look healthy. Router, providers, and supervised sessions are not reporting any cross-panel alerts.
                            </div>
                        ) : (
                            <div className="mt-4 grid gap-3 lg:grid-cols-2">
                                {dashboardAlerts.map((alert) => (
                                    <div key={alert.id} className="rounded-2xl border border-slate-800 bg-slate-900/70 p-4">
                                        <div className="flex items-start justify-between gap-3">
                                            <div>
                                                <div className={`inline-flex rounded-full border px-2.5 py-1 text-xs font-medium uppercase tracking-[0.2em] ${getAlertTone(alert.severity)}`}>
                                                    {alert.severity}
                                                </div>
                                                <h3 className="mt-3 text-base font-semibold text-white">{alert.title}</h3>
                                                <p className="mt-2 text-sm text-slate-300">{alert.detail}</p>
                                            </div>
                                        </div>
                                        <Link
                                            href={alert.href}
                                            title={alert.hrefLabel}
                                            aria-label={alert.hrefLabel}
                                            className="mt-4 inline-flex text-sm font-medium text-cyan-200 transition hover:text-cyan-100"
                                        >
                                            {alert.hrefLabel} →
                                        </Link>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </header>

                <MissionControlFunctionToggles />

                <SuggestionsPanel />

                <div className="grid gap-6 md:grid-cols-2">
                    <ContextHealthWidget />
                    <SessionHandoffWidget />
                </div>

                <div className="grid gap-6 xl:grid-cols-2">
                    <div className="flex flex-col gap-6">
                        <section className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6 shadow-lg shadow-slate-950/20">
                            <div className="flex items-start justify-between gap-4">
                                <div>
                                    <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-200">Overview</p>
                                    <h2 className="mt-2 text-xl font-semibold text-white">Router posture</h2>
                                <p className="mt-2 text-sm text-slate-400">Quick health readout for first-time operators.</p>
                            </div>
                        <div className={`rounded-full border px-3 py-1 text-xs font-medium ${routerStatusTone}`}>
                            {routerStatusLabel}
                        </div>
                    </div>

                        {mcpStatusError ? (
                            <div className="mt-4 rounded-2xl border border-rose-500/30 bg-rose-500/10 p-4 text-sm text-rose-200">
                                {mcpStatusError}
                            </div>
                        ) : null}

                        <dl className="mt-6 grid gap-4 sm:grid-cols-2">
                            <div className="rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
                                <dt className="text-sm text-slate-400">Connected servers</dt>
                                <dd className="mt-2 text-2xl font-semibold text-white">{isBootstrapping ? '—' : mcpStatus.connectedCount}</dd>
                            </div>
                            <div className="rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
                                <dt className="text-sm text-slate-400">Indexed tools</dt>
                                <dd className="mt-2 text-2xl font-semibold text-white">{isBootstrapping ? '—' : mcpStatus.toolCount}</dd>
                            </div>
                            <div className="rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
                                <dt className="text-sm text-slate-400">Running sessions</dt>
                                <dd className="mt-2 text-2xl font-semibold text-white">{isBootstrapping || sessionsError ? '—' : sessions.filter((session) => session.status === 'running').length}</dd>
                            </div>
                            <div className="rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
                                <dt className="text-sm text-slate-400">Configured providers</dt>
                                <dd className="mt-2 text-2xl font-semibold text-white">{isBootstrapping || providersError ? '—' : providers.filter((provider) => provider.configured).length}</dd>
                            </div>
                        </dl>

                        <div className="mt-6 rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
                            <div className="flex items-start justify-between gap-4">
                                <div>
                                    <h3 className="text-sm font-semibold uppercase tracking-[0.2em] text-slate-300">Startup readiness</h3>
                                    <p className="mt-1 text-sm text-slate-500">{startupSummary || 'Boot checks reported directly from core startup state.'}</p>
                                </div>
                                <span className={`rounded-full border px-3 py-1 text-xs font-medium ${startupToneClass}`}>
                                    {startupLabel}
                                </span>
                            </div>

                            {startupStatusError ? (
                                <div className="mt-4 rounded-2xl border border-rose-500/30 bg-rose-500/10 p-4 text-sm text-rose-200">
                                    {startupStatusError}
                                </div>
                            ) : (
                                <div className="mt-4 grid gap-3 sm:grid-cols-2">
                                    {startupChecklist.map((item) => (
                                        <div key={item.label} className="rounded-2xl border border-slate-800 bg-slate-900/70 p-3 text-sm">
                                            <div className="flex items-center justify-between gap-3">
                                                <span className="font-medium text-white">{item.label}</span>
                                                <span className={`rounded-full border px-2.5 py-1 text-xs font-medium ${item.ready ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-200' : 'border-amber-500/30 bg-amber-500/10 text-amber-200'}`}>
                                                    {item.ready ? 'Ready' : 'Pending'}
                                                </span>
                                            </div>
                                            <p className="mt-2 text-slate-400">{item.detail}</p>
                                        </div>
                                    ))}
                                </div>
                            )}

                            {startupModeEntries.length > 0 ? (
                                <div className="mt-4 rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
                                    <div className="flex items-start justify-between gap-4">
                                        <div>
                                            <h4 className="text-sm font-semibold uppercase tracking-[0.2em] text-slate-300">Startup mode</h4>
                                            <p className="mt-1 text-sm text-slate-500">Persisted launch provenance captured from the most recent startup handoff.</p>
                                        </div>
                                        {Number.isFinite(startupModeUpdatedAt) ? (
                                            <span className="rounded-full border border-slate-700 bg-slate-900/70 px-3 py-1 text-xs font-medium text-slate-200">
                                                Updated {formatRelativeTimestamp(startupModeUpdatedAt, currentTimestamp)}
                                            </span>
                                        ) : null}
                                    </div>
                                    <div className="mt-4 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
                                        {startupModeEntries.map((entry) => (
                                            <div key={entry.label} className="rounded-2xl border border-slate-800 bg-slate-900/70 p-3 text-sm">
                                                <div className="text-xs font-semibold uppercase tracking-[0.16em] text-slate-400">{entry.label}</div>
                                                <div className="mt-2 text-sm font-medium text-white">{entry.value}</div>
                                                {entry.detail ? (
                                                    <p className="mt-2 text-xs text-slate-400">{entry.detail}</p>
                                                ) : null}
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            ) : null}

                                    <p className="mt-1 text-sm text-slate-500">Fast path for getting browser bridges, editor surfaces, and managed MCP configs into the tools you already use.</p>
                                </div>
                                <Link
                                    href="/dashboard/integrations"
                                    title="Open install surfaces, browser extension artifacts, VS Code packaging, and client sync targets"
                                    aria-label="Open Integration Hub from router posture section"
                                    className="text-sm font-medium text-cyan-200 transition hover:text-cyan-100"
                                >
                                    Open Integration Hub →
                                </Link>
                            </div>

                            <div className="mt-4 grid gap-3 sm:grid-cols-3">
                                <div className="rounded-2xl border border-slate-800 bg-slate-900/70 p-3 text-sm text-slate-300">
                                    <div className="font-medium text-white">Browser extensions</div>
                                    <p className="mt-2 text-slate-400">Load Chromium/Edge and Firefox bundles, then connect them to the live bridge listener.</p>
                                </div>
                                <div className="rounded-2xl border border-slate-800 bg-slate-900/70 p-3 text-sm text-slate-300">
                                    <div className="font-medium text-white">Editor surfaces</div>
                                    <p className="mt-2 text-slate-400">Package and install the VS Code extension, then verify connected bridge clients and hook phases.</p>
                                </div>
                                <div className="rounded-2xl border border-slate-800 bg-slate-900/70 p-3 text-sm text-slate-300">
                                    <div className="font-medium text-white">Client config sync</div>
                                    <p className="mt-2 text-slate-400">Push HyperCode-managed MCP endpoints into Claude Desktop, Cursor, and VS Code without manual JSON surgery.</p>
                                </div>
                            </div>
                        </div>
                    </section>

                    <section className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6 shadow-lg shadow-slate-950/20">
                        <div className="flex items-start justify-between gap-4">
                            <div>
                                <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-200">MCP Router</p>
                                <h2 className="mt-2 text-xl font-semibold text-white">Server health and traffic</h2>
                                <p className="mt-2 text-sm text-slate-400">Live server posture plus the latest router activity.</p>
                            </div>
                            <Link
                                href="/dashboard/mcp"
                                title="Open the full MCP router dashboard with server list, tools, and configuration controls"
                                aria-label="Open detailed MCP dashboard"
                                className="text-sm font-medium text-cyan-200 transition hover:text-cyan-100"
                            >
                                Detailed MCP view →
                            </Link>
                        </div>

                        <div className="mt-6 space-y-3">
                            {serversError ? (
                                <div className="rounded-2xl border border-rose-500/30 bg-rose-500/10 p-4 text-sm text-rose-200">
                                    {serversError}
                                </div>
                            ) : servers.length === 0 ? (
                                <div className="rounded-2xl border border-dashed border-slate-700 bg-slate-950/50 p-4 text-sm text-slate-400">
                                    No MCP servers registered yet.
                                </div>
                            ) : servers.map((server) => (
                                <div key={server.name} className="rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
                                    <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                                        <div>
                                            <div className="flex items-center gap-3">
                                                <h3 className="text-base font-semibold text-white">{server.name}</h3>
                                                <span className={`rounded-full border px-2.5 py-1 text-xs font-medium ${getServerTone(server.status)}`}>
                                                    {sentenceCase(server.status)}
                                                </span>
                                            </div>
                                            <p className="mt-2 break-all font-mono text-xs text-slate-400">
                                                {server.config.command} {server.config.args.join(' ')}
                                            </p>
                                        </div>
                                        <div className="text-right text-sm text-slate-300">
                                            <div>{server.toolCount} tools</div>
                                            <div className="text-xs text-slate-500">{server.config.env.length} env vars</div>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>

                        <div className="mt-6 rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
                            <div className="flex items-center justify-between gap-4">
                                <h3 className="text-sm font-semibold uppercase tracking-[0.2em] text-slate-300">Recent traffic</h3>
                                <Link
                                    href="/dashboard/mcp/inspector"
                                    title="Open the live MCP inspector to trace requests, responses, and tool invocations"
                                    aria-label="Open MCP traffic inspector"
                                    className="text-sm font-medium text-cyan-200 transition hover:text-cyan-100"
                                >
                                    Open inspector →
                                </Link>
                            </div>

                            <div className="mt-4 space-y-3">
                                {trafficError ? (
                                    <p className="text-sm text-rose-300">{trafficError}</p>
                                ) : traffic.length === 0 ? (
                                    <p className="text-sm text-slate-400">No router traffic captured yet.</p>
                                ) : traffic.slice(0, 5).map((event, index) => (
                                    <div key={`${event.server}-${event.method}-${event.timestamp}-${index}`} className="rounded-2xl border border-slate-800 bg-slate-900/70 p-3 text-sm">
                                        <div className="flex items-center justify-between gap-4">
                                            <div className="flex items-center gap-3">
                                                <span className={`inline-flex h-2.5 w-2.5 rounded-full ${event.success ? 'bg-emerald-400' : 'bg-rose-400'}`} />
                                                <span className="font-medium text-white">{event.server}</span>
                                                <span className="rounded-full border border-slate-700 px-2 py-0.5 text-xs text-slate-300">{event.method}</span>
                                            </div>
                                            <span className="text-xs text-slate-500">{formatRelativeTimestamp(event.timestamp, currentTimestamp)}</span>
                                        </div>
                                        <p className="mt-2 text-sm text-slate-300">{summarizeTrafficEvent(event)}</p>
                                        <div className="mt-2 text-xs text-slate-500">Latency {event.latencyMs}ms</div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </section>

                    </div>

                    <div className="flex flex-col gap-6">
                    <section className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6 shadow-lg shadow-slate-950/20">
                        <div className="flex items-start justify-between gap-4">
                            <div>
                                <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-200">Sessions</p>
                                <h2 className="mt-2 text-xl font-semibold text-white">Supervised CLI runtime</h2>
                                <p className="mt-2 text-sm text-slate-400">Live posture for supervised coding sessions, restart policy, and the most recent activity.</p>
                            </div>
                            <Link
                                href="/dashboard/session"
                                title="Open the supervised session dashboard with logs, restart controls, and runtime details"
                                aria-label="Open sessions dashboard"
                                className="text-sm font-medium text-cyan-200 transition hover:text-cyan-100"
                            >
                                Open sessions →
                            </Link>
                        </div>

                        <div className="mt-6 space-y-3">
                            {sessionsError ? (
                                <div className="rounded-2xl border border-rose-500/30 bg-rose-500/10 p-4 text-sm text-rose-200">
                                    {sessionsError}
                                </div>
                            ) : sessions.length === 0 ? (
                                <div className="rounded-2xl border border-dashed border-slate-700 bg-slate-950/50 p-4 text-sm text-slate-400">
                                    No supervised sessions are active yet.
                                </div>
                            ) : sessions.map((session) => {
                                const latestLogMessage = getLatestSessionLogMessage(session);
                                const isPendingAction = pendingSessionActionId === session.id;

                                return (
                                    <div key={session.id} className="rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
                                        <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                                            <div className="min-w-0">
                                                <div className="flex flex-wrap items-center gap-3">
                                                    <h3 className="text-base font-semibold text-white">{session.name}</h3>
                                                    <span className={`rounded-full border px-2.5 py-1 text-xs font-medium ${getSessionTone(session.status)}`}>
                                                        {getSessionStatusLabel(session.status)}
                                                    </span>
                                                    <span className="rounded-full border border-slate-700 px-2.5 py-1 text-xs text-slate-300">
                                                        {session.cliType}
                                                    </span>
                                                </div>
                                                <p className="mt-2 break-all font-mono text-xs text-slate-500">{session.workingDirectory}</p>
                                                {session.lastError ? (
                                                    <p className="mt-2 text-sm text-rose-300">{session.lastError}</p>
                                                ) : null}
                                            </div>
                                            <div className="text-right text-sm text-slate-300">
                                                <div>{formatRelativeTimestamp(session.lastActivityAt, currentTimestamp)}</div>
                                                <div className="text-xs text-slate-500">
                                                    Restarts {session.restartCount}/{session.maxRestartAttempts}
                                                </div>
                                            </div>
                                        </div>

                                        <div className="mt-4 grid gap-3 lg:grid-cols-2">
                                            <div className="rounded-2xl border border-slate-800 bg-slate-900/70 p-3 text-sm text-slate-300">
                                                <div className="font-medium text-white">Restart posture</div>
                                                <p className="mt-2 text-slate-400">{getSessionRestartPolicyLabel(session)}</p>
                                                {session.scheduledRestartAt ? (
                                                    <p className="mt-2 text-cyan-200">Restart queued {formatRestartCountdown(session.scheduledRestartAt, currentTimestamp)}</p>
                                                ) : null}
                                            </div>
                                            <div className="rounded-2xl border border-slate-800 bg-slate-900/70 p-3 text-sm text-slate-300">
                                                <div className="font-medium text-white">Latest activity</div>
                                                <p className="mt-2 text-slate-400">{latestLogMessage ?? 'No session logs captured yet.'}</p>
                                            </div>
                                        </div>

                                        {(onStartSession || onStopSession || onRestartSession) ? (
                                            <div className="mt-4 flex flex-wrap gap-2">
                                                {session.status === 'stopped' || session.status === 'created' ? (
                                                    <button
                                                        type="button"
                                                        onClick={() => onStartSession?.(session.id)}
                                                        disabled={isPendingAction}
                                                        className="rounded-full border border-emerald-500/30 bg-emerald-500/10 px-3 py-1.5 text-sm font-medium text-emerald-100 transition hover:bg-emerald-500/20 disabled:cursor-not-allowed disabled:opacity-60"
                                                    >
                                                        {isPendingAction ? 'Starting…' : 'Start session'}
                                                    </button>
                                                ) : (
                                                    <button
                                                        type="button"
                                                        onClick={() => onStopSession?.(session.id)}
                                                        disabled={isPendingAction}
                                                        className="rounded-full border border-slate-700 bg-slate-900/70 px-3 py-1.5 text-sm font-medium text-slate-100 transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-60"
                                                    >
                                                        {isPendingAction ? 'Stopping…' : 'Stop session'}
                                                    </button>
                                                )}
                                                <button
                                                    type="button"
                                                    onClick={() => onRestartSession?.(session.id)}
                                                    disabled={isPendingAction}
                                                    className="rounded-full border border-cyan-500/30 bg-cyan-500/10 px-3 py-1.5 text-sm font-medium text-cyan-100 transition hover:bg-cyan-500/20 disabled:cursor-not-allowed disabled:opacity-60"
                                                >
                                                    {isPendingAction ? 'Working…' : 'Restart session'}
                                                </button>
                                            </div>
                                        ) : null}
                                    </div>
                                );
                            })}
                        </div>
                    </section>

                    <section className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6 shadow-lg shadow-slate-950/20">
                        <div className="flex items-start justify-between gap-4">
                            <div>
                                <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-200">Providers</p>
                                <h2 className="mt-2 text-xl font-semibold text-white">Quota and fallback posture</h2>
                                <p className="mt-2 text-sm text-slate-400">Which providers are configured, how much headroom remains, and where fallback will go next.</p>
                            </div>
                            <Link
                                href="/dashboard/billing"
                                title="Open provider billing and quota analytics with fallback chain controls"
                                aria-label="Open detailed provider billing dashboard"
                                className="text-sm font-medium text-cyan-200 transition hover:text-cyan-100"
                            >
                                Detailed provider view →
                            </Link>
                        </div>

                        <div className="mt-6 space-y-3">
                            {providersError ? (
                                <div className="rounded-2xl border border-rose-500/30 bg-rose-500/10 p-4 text-sm text-rose-200">
                                    {providersError}
                                </div>
                            ) : providers.length === 0 ? (
                                <div className="rounded-2xl border border-dashed border-slate-700 bg-slate-950/50 p-4 text-sm text-slate-400">
                                    No provider data available yet. Configure an API key or OAuth-backed provider in Billing to unlock fallback routing.
                                </div>
                            ) : providers.map((provider) => {
                                const usagePercent = getQuotaUsagePercent(provider);
                                return (
                                    <div key={provider.provider} className="rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
                                        <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                                            <div>
                                                <div className="flex items-center gap-3">
                                                    <h3 className="text-base font-semibold text-white">{provider.name}</h3>
                                                    <span className={`rounded-full border px-2.5 py-1 text-xs font-medium ${getProviderTone(provider)}`}>
                                                        {provider.configured ? sentenceCase(provider.availability ?? (provider.authenticated ? 'healthy' : 'degraded')) : 'Not configured'}
                                                    </span>
                                                </div>
                                                <p className="mt-2 text-sm text-slate-400">
                                                    {provider.authMethod ? `${provider.authMethod} · ` : ''}{provider.tier}
                                                    {provider.resetDate ? ` · resets ${provider.resetDate}` : ''}
                                                </p>
                                                {provider.lastError ? (
                                                    <p className="mt-2 text-sm text-rose-300">{provider.lastError}</p>
                                                ) : null}
                                            </div>
                                            <div className="text-right text-sm text-slate-300">
                                                <div>Used {formatQuotaValue(provider.used)}</div>
                                                <div>Remaining {formatQuotaValue(provider.remaining)}</div>
                                            </div>
                                        </div>
                                        <div className="mt-4">
                                            <div className="h-2 rounded-full bg-slate-800">
                                                <div
                                                    className="h-2 rounded-full bg-cyan-400 transition-all"
                                                    style={{ width: `${usagePercent ?? 100}%` }}
                                                />
                                            </div>
                                            <div className="mt-2 flex items-center justify-between text-xs text-slate-500">
                                                <span>Limit {formatQuotaValue(provider.limit)}</span>
                                                <span>{usagePercent === null ? 'Usage limit unavailable' : `${usagePercent}% used`}</span>
                                            </div>
                                        </div>
                                    </div>
                                );
                            })}
                        </div>

                        <div className="mt-6 rounded-2xl border border-slate-800 bg-slate-950/70 p-4">
                            <h3 className="text-sm font-semibold uppercase tracking-[0.2em] text-slate-300">Fallback chain</h3>
                            <div className="mt-4 space-y-2">
                                {fallbackChainError ? (
                                    <p className="text-sm text-rose-300">{fallbackChainError}</p>
                                ) : fallbackChain.length === 0 ? (
                                    <p className="text-sm text-slate-400">No fallback chain is exposed yet. Configure providers to populate the routing order.</p>
                                ) : fallbackChain.map((entry) => (
                                    <div key={`${entry.priority}-${entry.provider}-${entry.model ?? 'default'}`} className="flex items-center justify-between gap-4 rounded-2xl border border-slate-800 bg-slate-900/70 px-4 py-3 text-sm">
                                        <div>
                                            <span className="font-medium text-white">{entry.priority}. {formatFallbackLabel(entry)}</span>
                                            <p className="mt-1 text-xs text-slate-500">{entry.reason}</p>
                                        </div>
                                        <span className="rounded-full border border-slate-700 px-2.5 py-1 text-xs text-slate-300">priority {entry.priority}</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </section>

                        <HyperCodeOrchestratorWidget />

                    </div>

                    <div className="flex flex-col gap-6">
                        <HyperCodeOrchestratorWidget />
                    </div>
                </div>
            </div>
        </div>
    );
}

export function getStartupBlockingReasons(startupStatus: DashboardStartupStatus): StartupBlockingReasonView[] {
    if (!Array.isArray(startupStatus.blockingReasons)) {
        return [];
    }

    return startupStatus.blockingReasons
        .filter((reason): reason is StartupBlockingReasonView => Boolean(reason && typeof reason.code === 'string' && typeof reason.detail === 'string'))
        .map((reason) => ({
            code: reason.code,
            detail: reason.detail,
        }));
}

export function getStartupBlockingReasonAction(code: string): StartupBlockingReasonAction {
    switch (code) {
        case 'mcp_aggregator_not_initialized':
        case 'mcp_inventory_not_ready':
        case 'mcp_resident_runtime_not_ready':
        case 'mcp_config_sync_pending':
            return {
                href: '/dashboard/mcp/system',
                label: 'Open MCP system',
            };
        case 'memory_not_ready':
        case 'sectioned_memory_not_ready':
            return {
                href: '/dashboard/memory',
                label: 'Open memory dashboard',
            };
        case 'browser_service_not_ready':
        case 'extension_bridge_not_ready':
        case 'execution_environment_not_ready':
            return {
                href: '/dashboard/integrations',
                label: 'Open Integration Hub',
            };
        case 'session_restore_not_ready':
            return {
                href: '/dashboard/session',
                label: 'Open sessions',
            };
        default:
            return {
                href: '/dashboard',
                label: 'Open startup overview',
            };
    }
}

export function getStartupBlockingReasonImpactedChecks(code: string): StartupBlockingReasonImpactedCheck[] {
    switch (code) {
        case 'mcp_aggregator_not_initialized':
        case 'mcp_inventory_not_ready':
            return [
                { key: 'cached-inventory', label: 'Cached inventory' },
                { key: 'resident-runtime', label: 'Resident MCP runtime' },
            ];
        case 'mcp_resident_runtime_not_ready':
            return [
                { key: 'resident-runtime', label: 'Resident MCP runtime' },
            ];
        case 'mcp_config_sync_pending':
            return [
                { key: 'cached-inventory', label: 'Cached inventory' },
            ];
        case 'memory_not_ready':
        case 'sectioned_memory_not_ready':
            return [
                { key: 'memory-context', label: 'Memory / context' },
            ];
        case 'session_restore_not_ready':
            return [
                { key: 'session-restore', label: 'Session restore' },
            ];
        case 'browser_service_not_ready':
        case 'extension_bridge_not_ready':
            return [
                { key: 'client-bridge', label: 'Client bridge' },
            ];
        case 'execution_environment_not_ready':
            return [
                { key: 'execution-environment', label: 'Execution environment' },
            ];
        default:
            return [];
    }
}

export function getStartupBlockingReasonGroupImpactedChecks(
    reasons: StartupBlockingReasonWithPriority[],
): StartupBlockingReasonImpactedCheck[] {
    const seen = new Set<string>();
    const impactedChecks: StartupBlockingReasonImpactedCheck[] = [];

    for (const reason of reasons) {
        const checks = getStartupBlockingReasonImpactedChecks(reason.code);
        for (const check of checks) {
            if (seen.has(check.key)) {
                continue;
            }

            seen.add(check.key);
            impactedChecks.push(check);
        }
    }

    return impactedChecks;
}

export function getStartupBlockingReasonSubsystem(code: string): { key: string; label: string } {
    switch (code) {
        case 'mcp_aggregator_not_initialized':
        case 'mcp_inventory_not_ready':
        case 'mcp_resident_runtime_not_ready':
        case 'mcp_config_sync_pending':
            return {
                key: 'mcp',
                label: 'MCP router',
            };
        case 'memory_not_ready':
        case 'sectioned_memory_not_ready':
            return {
                key: 'memory',
                label: 'Memory / context',
            };
        case 'session_restore_not_ready':
            return {
                key: 'sessions',
                label: 'Session supervisor',
            };
        case 'browser_service_not_ready':
        case 'extension_bridge_not_ready':
        case 'execution_environment_not_ready':
            return {
                key: 'integrations',
                label: 'Integrations',
            };
        default:
            return {
                key: 'startup',
                label: 'Startup platform',
            };
    }
}

export function getStartupBlockingReasonTitle(code: string): string {
    switch (code) {
        case 'mcp_aggregator_not_initialized':
            return 'MCP router is not initialized';
        case 'mcp_inventory_not_ready':
            return 'Cached MCP inventory is not ready';
        case 'mcp_resident_runtime_not_ready':
            return 'Resident MCP runtime is still warming';
        case 'mcp_config_sync_pending':
            return 'MCP config sync is still pending';
        case 'memory_not_ready':
            return 'Memory manager is still initializing';
        case 'sectioned_memory_not_ready':
            return 'Sectioned memory default sections are not ready';
        case 'browser_service_not_ready':
            return 'Browser service bridge is not ready';
        case 'extension_bridge_not_ready':
            return 'Extension bridge listener is offline';
        case 'execution_environment_not_ready':
            return 'Execution environment verification is incomplete';
        case 'session_restore_not_ready':
            return 'Session restore has not completed yet';
        default:
            return 'Startup blocker requires operator attention';
    }
}

export function getStartupBlockingReasonPriority(code: string): number {
    switch (code) {
        case 'mcp_aggregator_not_initialized':
        case 'mcp_resident_runtime_not_ready':
        case 'execution_environment_not_ready':
            return 100;
        case 'mcp_inventory_not_ready':
        case 'mcp_config_sync_pending':
        case 'extension_bridge_not_ready':
            return 80;
        case 'memory_not_ready':
        case 'sectioned_memory_not_ready':
        case 'session_restore_not_ready':
            return 60;
        case 'browser_service_not_ready':
            return 40;
        default:
            return 20;
    }
}

export function getStartupBlockingReasonPriorityLabel(priority: number): 'High' | 'Medium' | 'Low' {
    if (priority >= 80) {
        return 'High';
    }

    if (priority >= 50) {
        return 'Medium';
    }

    return 'Low';
}

export function getStartupBlockingReasonPriorityTone(priorityLabel: 'High' | 'Medium' | 'Low'): string {
    switch (priorityLabel) {
        case 'High':
            return 'border-rose-500/40 bg-rose-500/10 text-rose-100';
        case 'Medium':
            return 'border-amber-500/40 bg-amber-500/10 text-amber-100';
        default:
            return 'border-emerald-500/40 bg-emerald-500/10 text-emerald-100';
    }
}

export function getStartupBlockingReasonPriorityCounts(
    startupBlockingReasons: StartupBlockingReasonWithPriority[],
): StartupBlockingReasonPriorityCounts {
    return startupBlockingReasons.reduce<StartupBlockingReasonPriorityCounts>((counts, reason) => {
        const label = getStartupBlockingReasonPriorityLabel(reason.priority);
        if (label === 'High') {
            counts.high += 1;
        } else if (label === 'Medium') {
            counts.medium += 1;
        } else {
            counts.low += 1;
        }

        return counts;
    }, {
        high: 0,
        medium: 0,
        low: 0,
    });
}

export function getPrioritizedStartupBlockingReasons(
    startupBlockingReasons: StartupBlockingReasonView[],
): StartupBlockingReasonWithPriority[] {
    return startupBlockingReasons
        .map((reason, index) => ({
            ...reason,
            priority: getStartupBlockingReasonPriority(reason.code),
            index,
        }))
        .sort((left, right) => {
            if (right.priority !== left.priority) {
                return right.priority - left.priority;
            }

            return left.index - right.index;
        })
        .map(({ index: _index, ...reason }) => reason);
}

export function getGroupedStartupBlockingReasons(
    startupBlockingReasons: StartupBlockingReasonWithPriority[],
): StartupBlockingReasonGroup[] {
    const groups = new Map<string, StartupBlockingReasonGroup>();

    for (const reason of startupBlockingReasons) {
        const subsystem = getStartupBlockingReasonSubsystem(reason.code);
        const existingGroup = groups.get(subsystem.key);
        if (existingGroup) {
            existingGroup.reasons.push(reason);
            continue;
        }

        groups.set(subsystem.key, {
            key: subsystem.key,
            label: subsystem.label,
            reasons: [reason],
        });
    }

    return Array.from(groups.values()).sort((left, right) => {
        const leftOrder = STARTUP_BLOCKING_REASON_GROUP_ORDER[left.key] ?? Number.MAX_SAFE_INTEGER;
        const rightOrder = STARTUP_BLOCKING_REASON_GROUP_ORDER[right.key] ?? Number.MAX_SAFE_INTEGER;
        if (leftOrder !== rightOrder) {
            return leftOrder - rightOrder;
        }

        return left.label.localeCompare(right.label);
    });
}

export function getStartupBlockingReasonGroupSeverity(
    reasons: StartupBlockingReasonWithPriority[],
): 'High' | 'Medium' | 'Low' {
    const maxPriority = reasons.reduce((highest, reason) => Math.max(highest, reason.priority), 0);
    return getStartupBlockingReasonPriorityLabel(maxPriority);
}

export function getStartupBlockingReasonGroupTopAction(
    reasons: StartupBlockingReasonWithPriority[],
): StartupBlockingReasonAction | null {
    if (reasons.length === 0) {
        return null;
    }

    const topReason = reasons.reduce((selected, reason) => {
        if (!selected) {
            return reason;
        }

        return reason.priority > selected.priority ? reason : selected;
    }, null as StartupBlockingReasonWithPriority | null);

    return topReason ? getStartupBlockingReasonAction(topReason.code) : null;
}

export function getStartupBlockingReasonGroupPrimaryReason(
    reasons: StartupBlockingReasonWithPriority[],
): StartupBlockingReasonWithPriority | null {
    if (reasons.length === 0) {
        return null;
    }

    return reasons.reduce((selected, reason) => {
        if (!selected) {
            return reason;
        }

        return reason.priority > selected.priority ? reason : selected;
    }, null as StartupBlockingReasonWithPriority | null);
}

export function getStartupBlockingReasonGroupPriorityCounts(
    reasons: StartupBlockingReasonWithPriority[],
): StartupBlockingReasonPriorityCounts {
    return getStartupBlockingReasonPriorityCounts(reasons);
}

export function getStartupBlockingReasonActions(
    startupBlockingReasons: StartupBlockingReasonView[],
): StartupBlockingReasonAction[] {
    const seen = new Set<string>();
    const actions: StartupBlockingReasonAction[] = [];

    for (const reason of startupBlockingReasons) {
        const action = getStartupBlockingReasonAction(reason.code);
        const key = `${action.href}|${action.label}`;
        if (seen.has(key)) {
            continue;
        }

        seen.add(key);
        actions.push(action);
    }

    return actions;
}
