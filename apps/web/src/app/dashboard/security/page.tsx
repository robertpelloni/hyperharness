
'use client';

import React, { useState } from 'react';
import { trpc } from '@/utils/trpc';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@borg/ui';
import { Badge } from '@borg/ui';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@borg/ui';
import { Button } from '@borg/ui';
import { ScrollArea } from '@borg/ui';
import { Shield, Lock } from 'lucide-react';

export default function SecurityPage() {
    const [auditLimit] = useState(50);
    const utils = trpc.useUtils();
    const { data: auditLogs, isLoading: loadingLogs, refetch: refetchLogs } = trpc.audit.query.useQuery({ limit: auditLimit });
    const { data: autonomyLevel } = trpc.autonomy.getLevel.useQuery();
    const lockdownMutation = trpc.autonomy.setLevel.useMutation({
        onSuccess: async () => {
            await utils.autonomy.getLevel.invalidate();
        }
    });

    const policies = [
        {
            toolName: 'Autonomy Guardrail',
            description: autonomyLevel === 'high'
                ? 'System currently in HIGH autonomy. Lockdown will reduce to LOW.'
                : 'System is in LOW autonomy lockdown mode.',
            state: autonomyLevel === 'high' ? 'Elevated' : 'Locked',
        },
        {
            toolName: 'Audit Logging',
            description: 'All tool executions and governance actions are persisted in audit logs.',
            state: 'Active',
        },
    ];


    return (
        <div className="p-6 space-y-6">
            <header className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Security & Governance</h1>
                    <p className="text-muted-foreground">Manage permissions, policies, and view audit logs.</p>
                </div>
                <div className="flex items-center gap-4">
                    <Badge variant={autonomyLevel === 'high' ? 'destructive' : 'secondary'} className="text-md px-3 py-1">
                        Autonomy: {autonomyLevel?.toUpperCase()}
                    </Badge>
                    {/* Lockdown Button */}
                    <Button
                        variant="destructive"
                        disabled={lockdownMutation.isPending || autonomyLevel === 'low'}
                        onClick={() => lockdownMutation.mutate({ level: 'low' })}
                        title={autonomyLevel === 'low' ? 'System already in lockdown mode' : 'Set autonomy to LOW immediately'}
                    >
                        <Lock className="w-4 h-4 mr-2" />
                        {lockdownMutation.isPending ? 'Applying...' : autonomyLevel === 'low' ? 'LOCKDOWN ACTIVE' : 'SYSTEM LOCKDOWN'}
                    </Button>
                </div>
            </header>

            <Tabs defaultValue="audit" className="w-full">
                <TabsList>
                    <TabsTrigger value="audit">Audit Logs</TabsTrigger>
                    <TabsTrigger value="policies">Policies</TabsTrigger>
                </TabsList>

                {/* AUDIT LOGS TAB */}
                <TabsContent value="audit" className="space-y-4">
                    <Card>
                        <CardHeader className="flex flex-row items-center justify-between">
                            <div>
                                <CardTitle>System Audit Log</CardTitle>
                                <CardDescription>Real-time record of all agent actions and tool executions.</CardDescription>
                            </div>
                            <Button variant="outline" size="sm" onClick={() => refetchLogs()}>Refresh</Button>
                        </CardHeader>
                        <CardContent>
                            <ScrollArea className="h-[600px] w-full rounded-md border p-4">
                                {loadingLogs ? <div>Loading logs...</div> : (
                                    <div className="space-y-4">
                                        {auditLogs?.map((log: any, i: number) => (
                                            <div key={i} className="flex flex-col gap-1 border-b pb-2 last:border-0">
                                                <div className="flex justify-between items-start">
                                                    <div className="flex items-center gap-2">
                                                        <span className="text-xs text-muted-foreground font-mono">
                                                            {new Date(log.timestamp).toLocaleTimeString()}
                                                        </span>
                                                        <Badge variant={
                                                            log.level === 'WARN' ? 'secondary' :
                                                                log.level === 'ERROR' ? 'destructive' : 'outline'
                                                        }>
                                                            {log.level}
                                                        </Badge>
                                                        <span className="font-semibold text-sm">{log.action}</span>
                                                    </div>
                                                </div>
                                                <pre className="text-xs bg-muted/50 p-2 rounded overflow-x-auto">
                                                    {typeof log.params === 'string' ? log.params : JSON.stringify(log.params, null, 2)}
                                                </pre>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </ScrollArea>
                        </CardContent>
                    </Card>
                </TabsContent>

                {/* POLICIES TAB */}
                <TabsContent value="policies" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Active Policies</CardTitle>
                            <CardDescription>Rules governing tool execution and resource access.</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="grid gap-4">
                                {policies.map((policy, i: number) => (
                                    <div key={i} className="flex items-center justify-between p-4 border rounded-lg bg-card">
                                        <div className="flex items-center gap-4">
                                            <Shield className="w-5 h-5 text-primary" />
                                            <div>
                                                <h3 className="font-medium">{policy.toolName}</h3>
                                                <p className="text-sm text-muted-foreground">{policy.description}</p>
                                            </div>
                                        </div>
                                        <Badge variant="outline">{policy.state}</Badge>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    );
}
