"use client";

import { useState } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from "@borg/ui";
import { Button } from "@borg/ui";
import { Loader2, Plus, Server, Globe, Terminal } from "lucide-react";
import { trpc } from '@/utils/trpc';

export default function MCPDashboard() {
    const { data: servers, isLoading } = trpc.mcp.listServers.useQuery();
    const { data: status } = trpc.mcp.getStatus.useQuery();
    const [isAddOpen, setIsAddOpen] = useState(false);

    return (
        <div className="p-8 space-y-8">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight text-white">MCP Aggregator</h1>
                    <p className="text-zinc-500">
                        {status ? `${status.serverCount} servers · ${status.toolCount} tools · ${status.connectedCount ?? 0} connected` : 'Manage downstream Model Context Protocol servers'}
                    </p>
                </div>
                <Button onClick={() => setIsAddOpen(!isAddOpen)} className="bg-blue-600 hover:bg-blue-500">
                    <Plus className="mr-2 h-4 w-4" /> Add Server
                </Button>
            </div>

            {isAddOpen && (
                <AddServerForm onSuccess={() => { setIsAddOpen(false); }} />
            )}

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {isLoading ? (
                    <div className="col-span-3 flex justify-center p-12">
                        <Loader2 className="h-8 w-8 animate-spin text-zinc-500" />
                    </div>
                ) : (servers?.length ?? 0) === 0 ? (
                    <div className="col-span-3 text-center p-12 text-zinc-500">
                        <Server className="h-12 w-12 mx-auto mb-4 opacity-30" />
                        <p className="text-lg font-medium">No MCP Servers Registered</p>
                        <p className="text-sm mt-1">Add downstream MCP servers to aggregate their tools into Borg.</p>
                    </div>
                ) : servers?.map((server: any) => (
                    <ServerCard key={server.name} server={server} />
                ))}
            </div>
        </div>
    );
}

function ServerCard({ server }: { server: any }) {
    return (
        <Card className="bg-zinc-900 border-zinc-800">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-lg font-medium text-zinc-200 flex items-center gap-2">
                    <Server className="h-5 w-5 text-zinc-500" />
                    {server.name}
                </CardTitle>
                <StatusBadge status={server.status} />
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    <div className="text-xs font-mono text-zinc-500 bg-black/50 p-2 rounded truncate">
                        {server.config.command} {server.config.args?.join(' ')}
                    </div>
                    <div className="flex justify-between text-sm text-zinc-400">
                        <span>Tools:</span>
                        <span className="font-bold text-white">{server.toolCount}</span>
                    </div>
                    {server.config.env && (
                        <div className="flex flex-wrap gap-1">
                            {Object.keys(server.config.env).map(k => (
                                <span key={k} className="px-1.5 py-0.5 rounded bg-zinc-800 text-[10px] text-zinc-500">{k}</span>
                            ))}
                        </div>
                    )}
                </div>
            </CardContent>
        </Card>
    );
}

function StatusBadge({ status }: { status: string }) {
    const colors = {
        connected: "bg-green-500/10 text-green-500 border-green-500/20",
        stopped: "bg-zinc-500/10 text-zinc-500 border-zinc-500/20",
        error: "bg-red-500/10 text-red-500 border-red-500/20"
    };
    return (
        <span className={`px-2 py-0.5 rounded text-xs border ${colors[status as keyof typeof colors] || colors.stopped}`}>
            {status.toUpperCase()}
        </span>
    );
}

function AddServerForm({ onSuccess }: { onSuccess: () => void }) {
    const [mode, setMode] = useState<'git' | 'manual'>('git');
    const [formData, setFormData] = useState({
        name: '',
        repoUrl: '',
        command: 'npx',
        args: '',
        env: ''
    });
    const [isPending] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("MCP router is not active. Enable it in trpc.ts to add servers.");
    };

    return (
        <Card className="bg-zinc-900 border-zinc-700 mb-8 border-l-4 border-l-blue-600">
            <CardContent className="pt-6">
                <div className="flex gap-4 mb-6">
                    <button
                        onClick={() => setMode('git')}
                        className={`flex items-center gap-2 px-4 py-2 rounded text-sm ${mode === 'git' ? 'bg-blue-600 text-white' : 'bg-zinc-800 text-zinc-400'}`}
                    >
                        <Globe className="h-4 w-4" /> From Git Repo
                    </button>
                    <button
                        onClick={() => setMode('manual')}
                        className={`flex items-center gap-2 px-4 py-2 rounded text-sm ${mode === 'manual' ? 'bg-blue-600 text-white' : 'bg-zinc-800 text-zinc-400'}`}
                    >
                        <Terminal className="h-4 w-4" /> Manual Config
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="text-xs text-zinc-500 uppercase font-bold">Server Name</label>
                            <input
                                required
                                value={formData.name}
                                onChange={e => setFormData({ ...formData, name: e.target.value })}
                                className="w-full bg-zinc-950 border border-zinc-800 rounded p-2 text-sm text-white"
                                placeholder="e.g. weather-server"
                            />
                        </div>
                        {mode === 'git' && (
                            <div>
                                <label className="text-xs text-zinc-500 uppercase font-bold">Repo URL</label>
                                <input
                                    required
                                    value={formData.repoUrl}
                                    onChange={e => setFormData({ ...formData, repoUrl: e.target.value })}
                                    className="w-full bg-zinc-950 border border-zinc-800 rounded p-2 text-sm text-white"
                                    placeholder="https://github.com/..."
                                />
                            </div>
                        )}
                    </div>

                    <div className="grid grid-cols-3 gap-4">
                        <div className="col-span-1">
                            <label className="text-xs text-zinc-500 uppercase font-bold">Command</label>
                            <input
                                required
                                value={formData.command}
                                onChange={e => setFormData({ ...formData, command: e.target.value })}
                                className="w-full bg-zinc-950 border border-zinc-800 rounded p-2 text-sm text-white font-mono"
                            />
                        </div>
                        <div className="col-span-2">
                            <label className="text-xs text-zinc-500 uppercase font-bold">Args (space separated)</label>
                            <input
                                value={formData.args}
                                onChange={e => setFormData({ ...formData, args: e.target.value })}
                                className="w-full bg-zinc-950 border border-zinc-800 rounded p-2 text-sm text-white font-mono"
                                placeholder="-y @modelcontextprotocol/server-memory"
                            />
                        </div>
                    </div>

                    <div>
                        <label className="text-xs text-zinc-500 uppercase font-bold">Environment (JSON)</label>
                        <textarea
                            value={formData.env}
                            onChange={e => setFormData({ ...formData, env: e.target.value })}
                            className="w-full bg-zinc-950 border border-zinc-800 rounded p-2 text-sm text-white font-mono h-20"
                            placeholder='{"API_KEY": "xyz"}'
                        />
                    </div>

                    <div className="flex justify-end pt-2">
                        <Button type="submit" disabled={isPending} className="bg-green-600 hover:bg-green-500">
                            {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {mode === 'git' ? 'Clone & Add Server' : 'Add Server Configuration'}
                        </Button>
                    </div>

                    {error && (
                        <div className="text-red-400 text-sm bg-red-900/20 p-2 rounded">
                            {error}
                        </div>
                    )}
                </form>
            </CardContent>
        </Card>
    );
}
