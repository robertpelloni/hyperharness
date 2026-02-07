'use client';

import React, { useState, useEffect } from 'react';
import { trpc } from '@/utils/trpc';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { ScrollArea } from "@/components/ui/scroll-area";
import { Badge } from "@/components/ui/badge";
import { Activity, File, Terminal, Zap } from 'lucide-react';

export default function EventsPage() {
    const [events, setEvents] = useState<any[]>([]);

    // In a real implementation, we would use a Subscription
    // trpc.pulse.subscribeEvents.useSubscription(...)
    // For now, we will poll or just show a placeholder interface awaiting the subscription logic

    return (
        <div className="container mx-auto p-6 space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">System Nervous System</h1>
                    <p className="text-muted-foreground">Real-time event stream from sensors and reactors.</p>
                </div>
                <Badge variant="outline" className="text-green-400 border-green-900 bg-green-950/30 animate-pulse">
                    <Activity className="w-3 h-3 mr-2" /> Live
                </Badge>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Active Sensors */}
                <Card className="col-span-1 border-zinc-800 bg-zinc-950/50">
                    <CardHeader>
                        <CardTitle className="text-sm font-medium uppercase text-muted-foreground">Active Sensors</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                                <File className="w-5 h-5 text-blue-400" />
                                <span>File System</span>
                            </div>
                            <Badge className="bg-green-600">Active</Badge>
                        </div>
                        <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                                <Terminal className="w-5 h-5 text-yellow-400" />
                                <span>Terminal Output</span>
                            </div>
                            <Badge className="bg-green-600">Active</Badge>
                        </div>
                    </CardContent>
                </Card>

                {/* Event Stream */}
                <Card className="col-span-1 lg:col-span-2 border-zinc-800 bg-black">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2"><Zap className="w-5 h-5 text-purple-400" /> Event Stream</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <ScrollArea className="h-[500px] w-full pr-4">
                            <div className="space-y-2">
                                {/* Placeholder events */}
                                {[1, 2, 3].map((i) => (
                                    <div key={i} className="flex flex-col p-3 rounded-md border border-zinc-900 bg-zinc-900/30 text-sm">
                                        <div className="flex justify-between items-start mb-1">
                                            <span className="font-mono text-purple-400">agent:heartbeat</span>
                                            <span className="text-xs text-zinc-500">14:02:3{i}</span>
                                        </div>
                                        <div className="text-zinc-400">Source: Director</div>
                                        <div className="font-mono text-xs text-zinc-600 mt-1">{`{ "status": "idle", "memory_usage": "450MB" }`}</div>
                                    </div>
                                ))}
                                <div className="flex flex-col p-3 rounded-md border border-blue-900/30 bg-blue-950/10 text-sm">
                                    <div className="flex justify-between items-start mb-1">
                                        <span className="font-mono text-blue-400">file:change</span>
                                        <span className="text-xs text-zinc-500">14:02:45</span>
                                    </div>
                                    <div className="text-zinc-400">Source: FileSensor</div>
                                    <div className="font-mono text-xs text-zinc-600 mt-1">{`{ "path": "packages/core/src/MCPServer.ts" }`}</div>
                                </div>
                            </div>
                        </ScrollArea>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}
