'use client';

import { trpc } from '@/lib/trpc';
import { useEffect, useState } from 'react';

interface MetricBucket {
    time: number;
    count: number;
    value_avg: number;
}

interface MetricsData {
    windowMs: number;
    totalEvents: number;
    counts: Record<string, number>;
    averages: Record<string, number>;
    series: MetricBucket[];
}

export default function MetricsPage() {
    const [data, setData] = useState<MetricsData | null>(null);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const result = await (trpc as any).metrics.getStats.query({});
                setData(result);
            } catch (e: any) {
                setError(e.message);
            }
        };
        fetchData();
        const interval = setInterval(fetchData, 5000);
        return () => clearInterval(interval);
    }, []);

    const formatBytes = (b: number) => {
        if (b > 1073741824) return `${(b / 1073741824).toFixed(1)} GB`;
        if (b > 1048576) return `${(b / 1048576).toFixed(1)} MB`;
        if (b > 1024) return `${(b / 1024).toFixed(1)} KB`;
        return `${b} B`;
    };

    return (
        <div className="p-6 space-y-6">
            <div className="flex flex-col gap-2">
                <h1 className="text-3xl font-bold tracking-tight">📊 Metrics</h1>
                <p className="text-muted-foreground">
                    System performance telemetry and resource monitoring.
                </p>
            </div>

            {error && (
                <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-4 text-destructive">
                    Error: {error}
                </div>
            )}

            {!data && !error && (
                <div className="text-muted-foreground animate-pulse">Loading metrics...</div>
            )}

            {data && (
                <>
                    {/* Summary Cards */}
                    <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                        <div className="bg-card border rounded-lg p-4">
                            <div className="text-sm text-muted-foreground">Total Events</div>
                            <div className="text-2xl font-bold">{data.totalEvents.toLocaleString()}</div>
                        </div>
                        <div className="bg-card border rounded-lg p-4">
                            <div className="text-sm text-muted-foreground">Heap Usage</div>
                            <div className="text-2xl font-bold">
                                {data.averages.memory_heap ? formatBytes(data.averages.memory_heap) : 'N/A'}
                            </div>
                        </div>
                        <div className="bg-card border rounded-lg p-4">
                            <div className="text-sm text-muted-foreground">RSS</div>
                            <div className="text-2xl font-bold">
                                {data.averages.memory_rss ? formatBytes(data.averages.memory_rss) : 'N/A'}
                            </div>
                        </div>
                        <div className="bg-card border rounded-lg p-4">
                            <div className="text-sm text-muted-foreground">System Load</div>
                            <div className="text-2xl font-bold">
                                {data.averages.system_load?.toFixed(2) ?? 'N/A'}
                            </div>
                        </div>
                    </div>

                    {/* Event Type Breakdown */}
                    <div className="bg-card border rounded-lg p-6">
                        <h2 className="text-lg font-semibold mb-4">Event Type Breakdown</h2>
                        <div className="space-y-2">
                            {Object.entries(data.counts).map(([type, count]) => (
                                <div key={type} className="flex justify-between items-center">
                                    <span className="font-mono text-sm">{type}</span>
                                    <div className="flex items-center gap-4">
                                        <div className="w-32 bg-muted rounded-full h-2">
                                            <div
                                                className="bg-primary rounded-full h-2"
                                                style={{ width: `${Math.min(100, (count / data.totalEvents) * 100)}%` }}
                                            />
                                        </div>
                                        <span className="text-sm text-muted-foreground w-16 text-right">
                                            {count}
                                        </span>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* Sparkline Chart */}
                    {data.series.length > 0 && (
                        <div className="bg-card border rounded-lg p-6">
                            <h2 className="text-lg font-semibold mb-4">Activity Over Time</h2>
                            <div className="flex items-end h-32 gap-px">
                                {data.series.map((bucket, i) => {
                                    const maxCount = Math.max(...data.series.map(s => s.count), 1);
                                    const height = (bucket.count / maxCount) * 100;
                                    return (
                                        <div
                                            key={i}
                                            className="flex-1 bg-primary/60 hover:bg-primary rounded-t transition-colors"
                                            style={{ height: `${Math.max(2, height)}%` }}
                                            title={`${new Date(bucket.time).toLocaleTimeString()}: ${bucket.count} events`}
                                        />
                                    );
                                })}
                            </div>
                            <div className="flex justify-between text-xs text-muted-foreground mt-1">
                                <span>{new Date(data.series[0].time).toLocaleTimeString()}</span>
                                <span>{new Date(data.series[data.series.length - 1].time).toLocaleTimeString()}</span>
                            </div>
                        </div>
                    )}
                </>
            )}
        </div>
    );
}
