"use client";

import { useEffect, useState, useRef } from 'react';
import { createReconnectPolicy, getReconnectDelayMs, resolveCoreWsUrl, shouldRetryReconnect } from '@borg/ui';
import { trpc } from '@/utils/trpc';

interface Packet {
    id: string;
    type: 'TOOL_CALL_START' | 'TOOL_CALL_END' | 'BROWSER_LOG' | 'KNOWLEDGE_CAPTURED' | 'BROWSER_DEBUG_EVENT' | 'MCP_TRAFFIC';
    tool?: string;
    args?: any;
    result?: string;
    duration?: number;
    success?: boolean;
    timestamp: number;
    level?: string;
    content?: string;
    url?: string;
    source?: string;
    title?: string;
    preview?: string;
    method?: string;
    params?: unknown;
    tabId?: number | null;
    server?: string;
}

type PacketTypeFilter = Packet['type'] | 'ALL';
type PacketSourceFilter = 'ALL' | 'browser_extension' | 'vscode_extension' | 'unknown';

const PACKET_TYPE_FILTERS: Array<{ value: PacketTypeFilter; label: string }> = [
    { value: 'ALL', label: 'All events' },
    { value: 'TOOL_CALL_START', label: 'Tool starts' },
    { value: 'TOOL_CALL_END', label: 'Tool results' },
    { value: 'MCP_TRAFFIC', label: 'MCP traffic' },
    { value: 'BROWSER_LOG', label: 'Logs' },
    { value: 'BROWSER_DEBUG_EVENT', label: 'CDP events' },
    { value: 'KNOWLEDGE_CAPTURED', label: 'Knowledge' },
];

const PACKET_SOURCE_FILTERS: Array<{ value: PacketSourceFilter; label: string }> = [
    { value: 'ALL', label: 'All sources' },
    { value: 'browser_extension', label: 'Browser ext' },
    { value: 'vscode_extension', label: 'VS Code ext' },
    { value: 'unknown', label: 'Unknown' },
];

function getPacketSearchText(packet: Packet): string {
    return [
        packet.type,
        packet.tool,
        packet.result,
        packet.content,
        packet.url,
        packet.source,
        packet.title,
        packet.preview,
        packet.method,
        typeof packet.params === 'string' ? packet.params : JSON.stringify(packet.params ?? {}),
    ].filter(Boolean).join(' ').toLowerCase();
}

export function TrafficInspector() {
    const [packets, setPackets] = useState<Packet[]>([]);
    const [isConnected, setIsConnected] = useState(false);
    const [customWsUrl, setCustomWsUrl] = useState<string>('');
    const [connectTrigger, setConnectTrigger] = useState(0); // Used to force reconnect
    const [typeFilter, setTypeFilter] = useState<PacketTypeFilter>('ALL');
    const [sourceFilter, setSourceFilter] = useState<PacketSourceFilter>('ALL');
    const [searchQuery, setSearchQuery] = useState('');
    const wsRef = useRef<WebSocket | null>(null);
    const reconnectAttemptsRef = useRef(0);
    const reconnectPolicy = createReconnectPolicy();

    useEffect(() => {
        // Prefer custom URL if set, otherwise fallback to env
        const targetUrl = customWsUrl || process.env.NEXT_PUBLIC_CORE_WS_URL || 'ws://localhost:3000';
        const wsUrl = resolveCoreWsUrl(targetUrl);

        console.log(`[TrafficInspector] Connecting to: ${wsUrl}`);

        const connect = () => {
            if (wsRef.current?.readyState === WebSocket.OPEN) {
                wsRef.current.close();
            }

            // Connect to Borg Core Bridge
            const ws = new WebSocket(wsUrl);

            ws.onopen = () => {
                setIsConnected(true);
                reconnectAttemptsRef.current = 0;
            };

            ws.onclose = () => {
                setIsConnected(false);
                if (shouldRetryReconnect(reconnectAttemptsRef.current, reconnectPolicy)) {
                    reconnectAttemptsRef.current += 1;
                    const delayMs = getReconnectDelayMs(reconnectAttemptsRef.current, reconnectPolicy);
                    setTimeout(connect, delayMs); // Reconnect with capped backoff
                }
            };

            ws.onerror = () => {
                // Let onclose handle capped retries.
                ws.close();
            };

            ws.onmessage = (event) => {
                try {
                    const msg = JSON.parse(event.data);
                    // Filter for interesting events
                    if (msg.type === 'TOOL_CALL_START' || msg.type === 'TOOL_CALL_END') {
                        addPacket({
                            ...msg,
                            timestamp: msg.timestamp ?? Date.now()
                        });
                    } else if (msg.type === 'BROWSER_LOG') {
                        const payload = msg.payload ?? msg;
                        const timestamp = Number(payload.timestamp) || Date.now();
                        addPacket({
                            id: `browser-log-${timestamp}-${payload.url ?? 'unknown'}`,
                            type: 'BROWSER_LOG',
                            timestamp,
                            level: String(payload.level ?? 'log'),
                            content: String(payload.content ?? payload.message ?? ''),
                            url: String(payload.url ?? 'unknown'),
                            source: String(payload.source ?? 'browser_extension')
                        });
                    } else if (msg.type === 'KNOWLEDGE_CAPTURED') {
                        const payload = msg.payload ?? msg;
                        const timestamp = Number(payload.timestamp) || Date.now();
                        addPacket({
                            id: `knowledge-captured-${payload.id ?? timestamp}`,
                            type: 'KNOWLEDGE_CAPTURED',
                            timestamp,
                            title: String(payload.title ?? 'Captured Page'),
                            url: String(payload.url ?? ''),
                            source: String(payload.source ?? 'browser_extension'),
                            preview: String(payload.preview ?? '')
                        });
                    } else if (msg.type === 'BROWSER_DEBUG_EVENT') {
                        const payload = msg.payload ?? msg;
                        const timestamp = Number(payload.timestamp) || Date.now();
                        addPacket({
                            id: `browser-debug-${payload.tabId ?? 'unknown'}-${payload.method ?? 'event'}-${timestamp}`,
                            type: 'BROWSER_DEBUG_EVENT',
                            timestamp,
                            source: String(payload.source ?? 'browser_extension'),
                            method: String(payload.method ?? 'unknown'),
                            params: payload.params,
                            tabId: typeof payload.tabId === 'number' ? payload.tabId : null,
                        });
                    }
                } catch (e) { }
            };

            wsRef.current = ws;
        };

        connect();

        return () => {
            wsRef.current?.close();
        };
    }, [connectTrigger]);

    const addPacket = (packet: Packet) => {
        setPackets(prev => {
            // Avoid duplicates
            if (prev.some(p => p.id === packet.id && p.type === packet.type)) return prev;
            const newPackets = [packet, ...prev].slice(0, 50);
            return newPackets;
        });
    };

    const handleReplay = async () => {
        // logs router is not active — replay not available
        console.warn('[TrafficInspector] Log replay not available — logs router is disabled');
    };

    const filteredPackets = packets.filter((packet) => {
        const normalizedSource = packet.source === 'browser_extension' || packet.source === 'vscode_extension'
            ? packet.source
            : 'unknown';
        const matchesType = typeFilter === 'ALL' || packet.type === typeFilter;
        const matchesSource = sourceFilter === 'ALL' || normalizedSource === sourceFilter;
        const matchesSearch = !searchQuery.trim() || getPacketSearchText(packet).includes(searchQuery.trim().toLowerCase());
        return matchesType && matchesSource && matchesSearch;
    });

    const packetCounts = packets.reduce<Record<string, number>>((acc, packet) => {
        acc[packet.type] = (acc[packet.type] ?? 0) + 1;
        return acc;
    }, {});

    return (
        <div className="bg-black/80 rounded-xl border border-zinc-800 overflow-hidden flex flex-col h-[600px]">
            {/* Header */}
            <div className="p-4 border-b border-zinc-800 bg-zinc-900/50 flex flex-col gap-2">
                <div className="flex justify-between items-center">
                    <div className="flex items-center gap-2">
                        <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'}`} />
                        <h2 className="font-mono font-bold text-zinc-300">NETWORK TRAFFIC (MCP)</h2>
                    </div>
                    <div className="flex gap-2">
                        <button
                            onClick={handleReplay}
                            className="text-xs px-2 py-1 bg-blue-900/30 text-blue-400 hover:bg-blue-900/50 rounded"
                        >
                            REPLAY LOGS
                        </button>
                        <button
                            onClick={() => setPackets([])}
                            className="text-xs text-zinc-500 hover:text-white"
                        >
                            CLEAR
                        </button>
                    </div>
                </div>

                {/* Connection Config Bar */}
                <div className="flex gap-2 items-center">
                    <input
                        type="text"
                        value={customWsUrl}
                        onChange={(e) => setCustomWsUrl(e.target.value)}
                        placeholder="ws://localhost:3000"
                        className="flex-1 bg-black/50 border border-zinc-800 rounded px-2 py-1 text-xs font-mono text-zinc-400 focus:border-blue-500 outline-none"
                    />
                    <button
                        onClick={() => setConnectTrigger(prev => prev + 1)}
                        className="text-xs px-3 py-1 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded font-mono"
                    >
                        CONNECT
                    </button>
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-[1fr_auto_auto] gap-2 items-center">
                    <input
                        type="text"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        placeholder="Filter by text, URL, method, tool..."
                        className="bg-black/50 border border-zinc-800 rounded px-2 py-1 text-xs font-mono text-zinc-400 focus:border-blue-500 outline-none"
                    />
                    <select
                        value={typeFilter}
                        onChange={(e) => setTypeFilter(e.target.value as PacketTypeFilter)}
                        className="bg-black/50 border border-zinc-800 rounded px-2 py-1 text-xs font-mono text-zinc-300 outline-none"
                    >
                        {PACKET_TYPE_FILTERS.map((filter) => (
                            <option key={filter.value} value={filter.value}>{filter.label}</option>
                        ))}
                    </select>
                    <select
                        value={sourceFilter}
                        onChange={(e) => setSourceFilter(e.target.value as PacketSourceFilter)}
                        className="bg-black/50 border border-zinc-800 rounded px-2 py-1 text-xs font-mono text-zinc-300 outline-none"
                    >
                        {PACKET_SOURCE_FILTERS.map((filter) => (
                            <option key={filter.value} value={filter.value}>{filter.label}</option>
                        ))}
                    </select>
                </div>

                <div className="flex flex-wrap gap-2 text-[11px] font-mono">
                    {PACKET_TYPE_FILTERS.filter((filter) => filter.value !== 'ALL').map((filter) => (
                        <span key={filter.value} className="px-2 py-1 rounded bg-zinc-800/80 text-zinc-400 border border-zinc-700">
                            {filter.label}: <span className="text-zinc-200">{packetCounts[filter.value] ?? 0}</span>
                        </span>
                    ))}
                    <span className="px-2 py-1 rounded bg-blue-900/30 text-blue-300 border border-blue-900/40">
                        Visible: {filteredPackets.length}
                    </span>
                </div>
            </div>

            {/* Packet Stream */}
            <div className="flex-1 overflow-y-auto p-4 space-y-2 font-mono text-sm">
                {filteredPackets.length === 0 && (
                    <div className="text-zinc-600 text-center mt-20">Waiting for traffic...</div>
                )}
                {filteredPackets.map((p, i) => (
                    <PacketRow key={`${p.id}-${p.type}-${i}`} packet={p} />
                ))}
            </div>
        </div>
    );
}

function PacketRow({ packet }: { packet: Packet }) {
    if (packet.type === 'BROWSER_LOG') {
        const level = packet.level ?? 'log';
        const icon = level === 'error' ? '🔴' : level === 'warn' ? '🟡' : '🔵';
        const borderColor = level === 'error'
            ? 'border-red-900/30 bg-red-900/10'
            : level === 'warn'
                ? 'border-yellow-900/30 bg-yellow-900/10'
                : 'border-sky-900/30 bg-sky-900/10';

        return (
            <div className={`p-3 rounded border ${borderColor} transition-all hover:bg-zinc-800/50`}>
                <div className="flex justify-between items-start gap-3">
                    <div className="flex items-center gap-3 min-w-0">
                        <span className="font-bold">{icon}</span>
                        <span className="text-zinc-200 font-bold">BROWSER LOG</span>
                        <span className="text-xs text-zinc-500 uppercase">{level}</span>
                        <span className="text-xs text-emerald-400/80 uppercase">{packet.source ?? 'browser_extension'}</span>
                        <span className="text-xs text-zinc-600 truncate">{packet.url}</span>
                    </div>
                    <span className="text-xs text-zinc-600 whitespace-nowrap">
                        {new Date(packet.timestamp).toLocaleTimeString().split(' ')[0]}
                    </span>
                </div>
                <div className="mt-2 pl-6 text-xs text-zinc-300 break-all">{packet.content}</div>
            </div>
        );
    }

    if (packet.type === 'KNOWLEDGE_CAPTURED') {
        return (
            <div className="p-3 rounded border border-emerald-900/30 bg-emerald-900/10 transition-all hover:bg-zinc-800/50">
                <div className="flex justify-between items-start gap-3">
                    <div className="flex items-center gap-3 min-w-0">
                        <span className="font-bold text-emerald-400">🧠</span>
                        <span className="text-zinc-200 font-bold truncate">{packet.title}</span>
                        <span className="text-xs text-zinc-600 truncate">{packet.url}</span>
                    </div>
                    <span className="text-xs text-zinc-600 whitespace-nowrap">
                        {new Date(packet.timestamp).toLocaleTimeString().split(' ')[0]}
                    </span>
                </div>
                <div className="mt-2 pl-6 text-xs text-zinc-400 break-all">
                    Source: <span className="text-emerald-300/80">{packet.source}</span>
                </div>
                {packet.preview && (
                    <div className="mt-1 pl-6 text-xs text-zinc-500 break-all">
                        Preview: {packet.preview}
                    </div>
                )}
            </div>
        );
    }

    if (packet.type === 'BROWSER_DEBUG_EVENT') {
        return (
            <div className="p-3 rounded border border-amber-900/30 bg-amber-900/10 transition-all hover:bg-zinc-800/50">
                <div className="flex justify-between items-start gap-3">
                    <div className="flex items-center gap-3 min-w-0">
                        <span className="font-bold text-amber-400">🪲</span>
                        <span className="text-zinc-200 font-bold">CDP EVENT</span>
                        <span className="text-xs text-amber-300 uppercase">{packet.method}</span>
                        <span className="text-xs text-zinc-500 uppercase">tab {packet.tabId ?? 'n/a'}</span>
                        <span className="text-xs text-emerald-400/80 uppercase">{packet.source ?? 'browser_extension'}</span>
                    </div>
                    <span className="text-xs text-zinc-600 whitespace-nowrap">
                        {new Date(packet.timestamp).toLocaleTimeString().split(' ')[0]}
                    </span>
                </div>
                <pre className="mt-2 pl-6 text-xs text-zinc-300 break-all whitespace-pre-wrap">{JSON.stringify(packet.params ?? {}, null, 2)}</pre>
            </div>
        );
    }

    if (packet.type === 'MCP_TRAFFIC') {
        const borderColor = packet.success ? 'border-violet-900/30 bg-violet-900/10' : 'border-red-900/30 bg-red-900/10';
        const iconColor = packet.success ? 'text-violet-400' : 'text-red-400';

        return (
            <div className={`p-3 rounded border ${borderColor} transition-all hover:bg-zinc-800/50`}>
                <div className="flex justify-between items-start gap-3">
                    <div className="flex items-center gap-3 min-w-0">
                        <span className={`font-bold ${iconColor}`}>⬢</span>
                        <span className="text-zinc-200 font-bold">{packet.method}</span>
                        <span className="text-xs text-violet-300 uppercase">{packet.server}</span>
                        {packet.tool ? <span className="text-xs text-zinc-500 uppercase">tool {packet.tool}</span> : null}
                    </div>
                    <span className="text-xs text-zinc-600 whitespace-nowrap">
                        {new Date(packet.timestamp).toLocaleTimeString().split(' ')[0]}
                    </span>
                </div>
                <div className="mt-2 pl-6 text-xs text-zinc-300 break-all">
                    Params: <span className="text-violet-200/80">{packet.content || '—'}</span>
                    <div className="mt-1 text-zinc-500">Latency: {packet.duration}ms</div>
                </div>
            </div>
        );
    }

    const isStart = packet.type === 'TOOL_CALL_START';

    // Color coding
    const borderColor = isStart ? 'border-blue-900/30' : (packet.success ? 'border-green-900/30' : 'border-red-900/30');
    const bgColor = isStart ? 'bg-blue-900/10' : (packet.success ? 'bg-green-900/10' : 'bg-red-900/10');
    const icon = isStart ? '→' : (packet.success ? '✓' : '✗');
    const iconColor = isStart ? 'text-blue-400' : (packet.success ? 'text-green-400' : 'text-red-400');

    return (
        <div className={`p-3 rounded border ${borderColor} ${bgColor} transition-all hover:bg-zinc-800/50`}>
            <div className="flex justify-between items-start">
                <div className="flex items-center gap-3">
                    <span className={`font-bold ${iconColor}`}>{icon}</span>
                    <span className="text-zinc-300 font-bold">{packet.tool}</span>
                    <span className="text-xs text-zinc-600">#{packet.id.substring(0, 4)}</span>
                </div>
                <span className="text-xs text-zinc-600">
                    {new Date(packet.timestamp).toLocaleTimeString().split(' ')[0]}
                </span>
            </div>

            {/* Details */}
            <div className="mt-2 pl-6">
                {isStart ? (
                    <div className="text-zinc-400 break-all text-xs">
                        Args: <span className="text-blue-300/80">{JSON.stringify(packet.args).substring(0, 200)}</span>
                    </div>
                ) : (
                    <div className="text-zinc-400 break-all text-xs">
                        Result: <span className="text-zinc-500">{packet.result}</span>
                        <div className="mt-1 text-zinc-600">
                            Duration: {packet.duration}ms
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
