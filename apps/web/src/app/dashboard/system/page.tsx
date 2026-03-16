"use client";

import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@borg/ui";
import { Badge } from "@borg/ui";
import { Button } from "@borg/ui";
import { Activity, FileText, Shield, Server, Clock, CheckCircle2, AlertCircle, Loader2 } from "lucide-react";
import { trpc } from '@/utils/trpc';
import { PageStatusBanner } from '@/components/PageStatusBanner';

function formatUptime(seconds: number): string {
    if (seconds < 60) return `${Math.floor(seconds)}s`;
    if (seconds < 3600) return `${Math.floor(seconds / 60)}m ${Math.floor(seconds % 60)}s`;
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    return `${h}h ${m}m`;
}

export default function SystemOverview() {
    const { data: startupStatus, isLoading } = trpc.startupStatus.useQuery(undefined, { refetchInterval: 10000 });

    const checks = startupStatus?.checks;
    const subsystems: { label: string; ready: boolean | undefined }[] = [
        { label: 'MCP Aggregator', ready: checks?.mcpAggregator?.liveReady },
        { label: 'Config Sync', ready: checks?.configSync?.ready },
        { label: 'Memory', ready: checks?.memory?.ready },
        { label: 'Browser', ready: checks?.browser?.ready },
        { label: 'Session Supervisor', ready: checks?.sessionSupervisor?.ready },
        { label: 'Extension Bridge', ready: checks?.extensionBridge?.ready },
        { label: 'Execution Env', ready: checks?.executionEnvironment?.ready },
    ];

    const readyCount = subsystems.filter((s) => s.ready === true).length;
    const overallReady = startupStatus?.ready === true;

    return (
        <div className="p-8 space-y-8 h-full overflow-y-auto w-full max-w-[1200px] mx-auto">
            <PageStatusBanner status="beta" message="System overview aggregates subsystem readiness from the startup status contract." />

            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight text-white flex items-center gap-3">
                        <Server className="h-8 w-8 text-violet-500" />
                        System Overview
                    </h1>
                    <p className="text-zinc-500 mt-1">
                        Borg operator console — subsystem health, uptime, and quick navigation.
                    </p>
                </div>
                {startupStatus && (
                    <Badge className={overallReady ? 'bg-green-900 text-green-300' : 'bg-amber-900 text-amber-300'}>
                        {overallReady ? 'Ready' : 'Pending'}
                    </Badge>
                )}
            </div>

            {/* Stats row */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <Card className="bg-zinc-900 border-zinc-800">
                    <CardContent className="pt-6 flex items-center gap-4">
                        <Clock className="h-8 w-8 text-zinc-400 shrink-0" />
                        <div>
                            <p className="text-xs text-zinc-500 uppercase tracking-wider">Uptime</p>
                            <p className="text-2xl font-bold text-white">
                                {isLoading ? '—' : formatUptime(startupStatus?.uptime ?? 0)}
                            </p>
                        </div>
                    </CardContent>
                </Card>
                <Card className="bg-zinc-900 border-zinc-800">
                    <CardContent className="pt-6 flex items-center gap-4">
                        <Activity className="h-8 w-8 text-zinc-400 shrink-0" />
                        <div>
                            <p className="text-xs text-zinc-500 uppercase tracking-wider">Subsystems</p>
                            <p className="text-2xl font-bold text-white">
                                {isLoading ? '—' : `${readyCount}/${subsystems.length}`}
                            </p>
                        </div>
                    </CardContent>
                </Card>
                <Card className="bg-zinc-900 border-zinc-800">
                    <CardContent className="pt-6 flex items-center gap-4">
                        <Server className="h-8 w-8 text-zinc-400 shrink-0" />
                        <div>
                            <p className="text-xs text-zinc-500 uppercase tracking-wider">Version</p>
                            <p className="text-2xl font-bold text-white">
                                {isLoading ? '—' : (startupStatus?.runtime?.version ?? '—')}
                            </p>
                        </div>
                    </CardContent>
                </Card>
            </div>

            {/* Subsystem checks */}
            <Card className="bg-zinc-900 border-zinc-800">
                <CardHeader>
                    <CardTitle className="text-white text-lg">Subsystem Status</CardTitle>
                    <CardDescription>Real-time readiness from the startup contract.</CardDescription>
                </CardHeader>
                <CardContent>
                    {isLoading ? (
                        <div className="flex items-center gap-2 text-zinc-500">
                            <Loader2 className="h-4 w-4 animate-spin" />
                            <span>Loading...</span>
                        </div>
                    ) : (
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                            {subsystems.map((sub) => (
                                <div key={sub.label} className="flex items-center gap-3 p-3 rounded-lg bg-zinc-800/50">
                                    {sub.ready === true ? (
                                        <CheckCircle2 className="h-4 w-4 text-green-500 shrink-0" />
                                    ) : (
                                        <AlertCircle className="h-4 w-4 text-amber-400 shrink-0" />
                                    )}
                                    <span className="text-sm text-zinc-300">{sub.label}</span>
                                    <Badge className={`ml-auto text-xs ${sub.ready ? 'bg-green-900/50 text-green-400' : 'bg-amber-900/50 text-amber-400'}`}>
                                        {sub.ready ? 'Ready' : 'Pending'}
                                    </Badge>
                                </div>
                            ))}
                        </div>
                    )}
                </CardContent>
            </Card>

            {/* Blocking reasons */}
            {!overallReady && startupStatus?.blockingReasons && startupStatus.blockingReasons.length > 0 && (
                <Card className="bg-amber-950/20 border-amber-800/40">
                    <CardHeader>
                        <CardTitle className="text-amber-300 text-base flex items-center gap-2">
                            <AlertCircle className="h-4 w-4" /> Boot Pending
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <ul className="space-y-1">
                            {startupStatus.blockingReasons.map((reason) => (
                                <li key={reason.code} className="text-sm text-amber-200/80">
                                    <span className="font-mono text-amber-400 mr-2">[{reason.code}]</span>
                                    {reason.detail}
                                </li>
                            ))}
                        </ul>
                    </CardContent>
                </Card>
            )}

            {/* Navigation cards */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <Card className="bg-zinc-900 border-zinc-800 hover:border-zinc-600 transition-colors">
                    <CardHeader className="pb-2">
                        <CardTitle className="text-white text-base flex items-center gap-2">
                            <Activity className="h-4 w-4 text-green-500" /> Health
                        </CardTitle>
                        <CardDescription>Server health, MCP servers, crash rates.</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <Button asChild variant="outline" size="sm" className="border-zinc-700 hover:bg-zinc-800">
                            <Link href="/dashboard/health">Open Health</Link>
                        </Button>
                    </CardContent>
                </Card>
                <Card className="bg-zinc-900 border-zinc-800 hover:border-zinc-600 transition-colors">
                    <CardHeader className="pb-2">
                        <CardTitle className="text-white text-base flex items-center gap-2">
                            <FileText className="h-4 w-4 text-cyan-500" /> Logs
                        </CardTitle>
                        <CardDescription>Execution logs, tool calls, error rates.</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <Button asChild variant="outline" size="sm" className="border-zinc-700 hover:bg-zinc-800">
                            <Link href="/dashboard/logs">Open Logs</Link>
                        </Button>
                    </CardContent>
                </Card>
                <Card className="bg-zinc-900 border-zinc-800 hover:border-zinc-600 transition-colors">
                    <CardHeader className="pb-2">
                        <CardTitle className="text-white text-base flex items-center gap-2">
                            <Shield className="h-4 w-4 text-indigo-500" /> Audit
                        </CardTitle>
                        <CardDescription>System audit trail for operator review.</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <Button asChild variant="outline" size="sm" className="border-zinc-700 hover:bg-zinc-800">
                            <Link href="/dashboard/audit">Open Audit</Link>
                        </Button>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}
