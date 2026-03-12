"use client";

import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@borg/ui";
import { Button } from "@borg/ui";
import { Badge } from "@borg/ui";
import { Input } from "@borg/ui";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@borg/ui";
import { Activity, Trash2, Search, RefreshCcw, BarChart3, AlertTriangle, CheckCircle2, Clock } from "lucide-react";
import { trpc } from '@/utils/trpc';
import { toast } from 'sonner';

function formatDuration(ms: number) {
    if (ms < 1000) return `${ms}ms`;
    return `${(ms / 1000).toFixed(2)}s`;
}

function formatTime(timestamp: number) {
    return new Date(timestamp).toLocaleTimeString([], { hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit' });
}

export default function LogsDashboard() {
    const [searchQuery, setSearchQuery] = useState("");
    const [limit, setLimit] = useState(100);

    const utils = trpc.useUtils();
    
    // Fetch summary metrics
    const { data: summary, refetch: refetchSummary, isFetching: isFetchingSummary } = trpc.logs.summary.useQuery({ limit: 5000 });
    
    // Fetch log entries
    const { data: logs, refetch: refetchLogs, isFetching: isFetchingLogs } = trpc.logs.list.useQuery({ limit });

    const clearLogs = trpc.logs.clear.useMutation({
        onSuccess: () => {
            toast.success("Logs cleared successfully");
            utils.logs.summary.invalidate();
            utils.logs.list.invalidate();
        },
        onError: (err) => toast.error(`Failed to clear logs: ${err.message}`),
    });

    const handleRefresh = async () => {
        await Promise.all([refetchSummary(), refetchLogs()]);
        toast.success("Logs refreshed");
    };

    const isRefreshing = isFetchingSummary || isFetchingLogs;

    // Filter logs client-side for immediate feedback
    const filteredLogs = logs?.filter((log: any) => {
        if (!searchQuery) return true;
        const query = searchQuery.toLowerCase();
        return (
            (log.toolName && log.toolName.toLowerCase().includes(query)) ||
            (log.serverName && log.serverName.toLowerCase().includes(query)) ||
            (log.error && log.error.toLowerCase().includes(query))
        );
    }) || [];

    return (
        <div className="p-8 space-y-8 h-full overflow-y-auto w-full max-w-[1600px] mx-auto">
            {/* Header */}
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight text-white flex items-center gap-3">
                        <Activity className="h-8 w-8 text-cyan-500" />
                        Execution Logs
                    </h1>
                    <p className="text-zinc-500 mt-1">
                        Review real-time tool execution history, error rates, and performance metrics.
                    </p>
                </div>
                <div className="flex items-center gap-3">
                    <Button 
                        onClick={handleRefresh} 
                        disabled={isRefreshing}
                        variant="outline" 
                        className="border-zinc-700 hover:bg-zinc-800"
                    >
                        <RefreshCcw className={`mr-2 h-4 w-4 ${isRefreshing ? 'animate-spin' : ''}`} /> 
                        Refresh
                    </Button>
                    <Button 
                        onClick={() => clearLogs.mutate()} 
                        disabled={clearLogs.isPending || (summary?.totals?.totalCalls === 0)}
                        variant="destructive" 
                        className="bg-red-950/50 text-red-500 hover:bg-red-900/50 border border-red-900"
                    >
                        <Trash2 className="mr-2 h-4 w-4" /> 
                        Clear Logs
                    </Button>
                </div>
            </div>

            {/* Metrics Row */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <Card className="bg-zinc-900 border-zinc-800">
                    <CardContent className="p-6">
                        <div className="flex items-center justify-between mb-2">
                            <span className="text-zinc-400 font-medium text-sm">Total Calls</span>
                            <Activity className="h-4 w-4 text-blue-500" />
                        </div>
                        <div className="text-3xl font-bold text-white mb-1">
                            {summary?.totals?.totalCalls ?? 0}
                        </div>
                    </CardContent>
                </Card>
                <Card className="bg-zinc-900 border-zinc-800">
                    <CardContent className="p-6">
                        <div className="flex items-center justify-between mb-2">
                            <span className="text-zinc-400 font-medium text-sm">Success Rate</span>
                            <CheckCircle2 className="h-4 w-4 text-green-500" />
                        </div>
                        <div className="text-3xl font-bold text-white mb-1">
                            {summary?.totals?.successRate ?? 100}%
                        </div>
                    </CardContent>
                </Card>
                <Card className="bg-zinc-900 border-zinc-800">
                    <CardContent className="p-6">
                        <div className="flex items-center justify-between mb-2">
                            <span className="text-zinc-400 font-medium text-sm">Error Count</span>
                            <AlertTriangle className="h-4 w-4 text-red-500" />
                        </div>
                        <div className="text-3xl font-bold text-white mb-1">
                            {summary?.totals?.errorCount ?? 0}
                        </div>
                    </CardContent>
                </Card>
                <Card className="bg-zinc-900 border-zinc-800">
                    <CardContent className="p-6">
                        <div className="flex items-center justify-between mb-2">
                            <span className="text-zinc-400 font-medium text-sm">Avg Duration</span>
                            <Clock className="h-4 w-4 text-purple-500" />
                        </div>
                        <div className="text-3xl font-bold text-white mb-1">
                            {summary?.totals?.avgDurationMs ? formatDuration(summary.totals.avgDurationMs) : '0ms'}
                        </div>
                    </CardContent>
                </Card>
            </div>

            {/* Main Content Area */}
            <div className="grid grid-cols-1 xl:grid-cols-4 gap-8">
                
                {/* Left Column: Log Feed */}
                <div className="xl:col-span-3 space-y-4">
                    <div className="flex items-center gap-4">
                        <div className="relative flex-1">
                            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-zinc-500" />
                            <Input 
                                placeholder="Filter by tool, server, or error..." 
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                className="w-full pl-9 bg-zinc-900 border-zinc-800 text-white"
                            />
                        </div>
                        <div className="flex items-center gap-2 text-sm text-zinc-500">
                            <span>Showing {filteredLogs.length} of {limit}</span>
                            <select 
                                value={limit} 
                                onChange={(e) => setLimit(Number(e.target.value))}
                                className="bg-zinc-900 border border-zinc-800 rounded px-2 py-1 text-zinc-300 outline-none"
                            >
                                <option value={50}>50</option>
                                <option value={100}>100</option>
                                <option value={500}>500</option>
                                <option value={1000}>1000</option>
                            </select>
                        </div>
                    </div>

                    <Card className="bg-zinc-900 border-zinc-800 overflow-hidden">
                        <div className="overflow-x-auto">
                            <Table>
                                <TableHeader className="bg-zinc-950">
                                    <TableRow className="border-zinc-800 hover:bg-transparent">
                                        <TableHead className="w-[100px] text-zinc-400">Time</TableHead>
                                        <TableHead className="w-[80px] text-zinc-400">Status</TableHead>
                                        <TableHead className="w-[180px] text-zinc-400">Server</TableHead>
                                        <TableHead className="w-[200px] text-zinc-400">Tool</TableHead>
                                        <TableHead className="w-[100px] text-zinc-400 text-right">Duration</TableHead>
                                        <TableHead className="text-zinc-400">Result / Error</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {filteredLogs.length === 0 ? (
                                        <TableRow>
                                            <TableCell colSpan={6} className="h-32 text-center text-zinc-500">
                                                No logs found.
                                            </TableCell>
                                        </TableRow>
                                    ) : (
                                        filteredLogs.map((log: any) => {
                                            const isError = log.level === 'error' || Boolean(log.error);
                                            return (
                                                <TableRow key={log.id} className="border-zinc-800/50 hover:bg-zinc-800/30">
                                                    <TableCell className="font-mono text-xs text-zinc-500 whitespace-nowrap">
                                                        {formatTime(log.timestamp)}
                                                    </TableCell>
                                                    <TableCell>
                                                        {isError ? (
                                                            <Badge variant="destructive" className="h-5 text-[10px]">ERROR</Badge>
                                                        ) : (
                                                            <Badge variant="outline" className="h-5 text-[10px] border-green-500/30 text-green-400 bg-green-500/10">OK</Badge>
                                                        )}
                                                    </TableCell>
                                                    <TableCell className="font-mono text-xs text-zinc-300 truncate max-w-[180px]" title={log.serverName}>
                                                        {log.serverName || 'System'}
                                                    </TableCell>
                                                    <TableCell className="font-mono text-xs text-blue-400 truncate max-w-[200px]" title={log.toolName}>
                                                        {log.toolName}
                                                    </TableCell>
                                                    <TableCell className="font-mono text-xs text-zinc-400 text-right">
                                                        {formatDuration(log.durationMs || 0)}
                                                    </TableCell>
                                                    <TableCell className="text-xs text-zinc-400 truncate max-w-[400px]" title={isError ? log.error : 'Success'}>
                                                        {isError ? (
                                                            <span className="text-red-400">{log.error || 'Unknown error'}</span>
                                                        ) : (
                                                            <span className="text-zinc-500">Execution completed</span>
                                                        )}
                                                    </TableCell>
                                                </TableRow>
                                            );
                                        })
                                    )}
                                </TableBody>
                            </Table>
                        </div>
                    </Card>
                </div>

                {/* Right Column: Top Tools */}
                <div className="space-y-6">
                    <Card className="bg-zinc-900 border-zinc-800">
                        <CardHeader className="pb-4">
                            <CardTitle className="text-lg font-medium text-white flex items-center gap-2">
                                <BarChart3 className="h-5 w-5 text-zinc-400" />
                                Top Activity
                            </CardTitle>
                            <CardDescription>Most frequently executed tools</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-4">
                                {!summary?.topTools || summary.topTools.length === 0 ? (
                                    <div className="text-center text-sm text-zinc-600 py-4">No tool activity recorded yet.</div>
                                ) : (
                                    summary.topTools.map((tool: any) => (
                                        <div key={tool.name} className="flex flex-col gap-1">
                                            <div className="flex justify-between items-center">
                                                <span className="text-xs font-mono text-zinc-300 truncate max-w-[180px]" title={tool.name}>
                                                    {tool.name}
                                                </span>
                                                <span className="text-xs font-medium text-zinc-400 bg-zinc-800 px-2 py-0.5 rounded">
                                                    {tool.count}
                                                </span>
                                            </div>
                                            {/* Small activity bar relative to max item */}
                                            <div className="h-1.5 w-full bg-zinc-800 rounded-full overflow-hidden flex">
                                                <div 
                                                    className="bg-blue-500 h-full" 
                                                    style={{ width: `${Math.max(2, (tool.count / summary.topTools[0].count) * 100 * (1 - (tool.errors / tool.count)))}%` }} 
                                                />
                                                {tool.errors > 0 && (
                                                    <div 
                                                        className="bg-red-500 h-full opacity-60" 
                                                        style={{ width: `${Math.max(2, (tool.errors / tool.count) * 100)}%` }} 
                                                        title={`${tool.errors} errors`}
                                                    />
                                                )}
                                            </div>
                                        </div>
                                    ))
                                )}
                            </div>
                        </CardContent>
                    </Card>
                </div>

            </div>
        </div>
    );
}
