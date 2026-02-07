
'use client';

import { trpcc } from '../../../utils/trpc';
import { useState, useEffect } from 'react';

export default function HealerDashboard() {
    const historyQuery = trpcc.healer.getHistory.useQuery(undefined, {
        refetchInterval: 5000 // Poll every 5s
    });

    const [activeInfections, setActiveInfections] = useState<any[]>([]);

    useEffect(() => {
        // In a real app, we might subscribe to an event stream for active infections
        // For now, we'll just mock it or infer from logs if possible. 
        // But since we don't have a stream yet, we will rely on history.
    }, []);

    const history = historyQuery.data || [];

    return (
        <div className="p-8 bg-gray-900 min-h-screen text-gray-100 font-mono">
            <header className="mb-8 border-b border-gray-700 pb-4 flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold text-green-400">THE IMMUNE SYSTEM</h1>
                    <p className="text-gray-400">Self-Healing & Auto-Correction</p>
                </div>
                <div className="text-right">
                    <div className="text-sm text-gray-500">STATUS</div>
                    <div className="text-green-500 font-bold">ACTIVE</div>
                </div>
            </header>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                {/* Active Infections (Mocked/Empty for now unless we add streaming) */}
                <section className="bg-gray-800 rounded-lg border border-gray-700 p-6">
                    <h2 className="text-xl font-bold mb-4 text-red-400 flex items-center">
                        <span className="animate-pulse mr-2">●</span> ACTIVE INFECTIONS
                    </h2>
                    {activeInfections.length === 0 ? (
                        <div className="text-gray-500 italic text-center py-10">
                            No active pathogens detected. System healthy.
                        </div>
                    ) : (
                        <div className="space-y-4">
                            {activeInfections.map((inf, i) => (
                                <div key={i} className="bg-red-900/20 border border-red-500/50 p-4 rounded">
                                    {JSON.stringify(inf)}
                                </div>
                            ))}
                        </div>
                    )}
                </section>

                {/* Immune History */}
                <section className="bg-gray-800 rounded-lg border border-gray-700 p-6">
                    <h2 className="text-xl font-bold mb-4 text-blue-400">IMMUNE HISTORY</h2>
                    <div className="space-y-4 max-h-[600px] overflow-y-auto">
                        {history.length === 0 ? (
                            <div className="text-gray-500 italic text-center py-10">
                                No history available.
                            </div>
                        ) : (
                            history.slice().reverse().map((entry: any, i: number) => (
                                <div key={i} className="bg-gray-700/50 border border-gray-600 p-4 rounded">
                                    <div className="flex justify-between text-xs text-gray-400 mb-2">
                                        <span>{new Date(entry.timestamp).toLocaleString()}</span>
                                        <span className={entry.success ? "text-green-400" : "text-red-400"}>
                                            {entry.success ? "NEUTRALIZED" : "FAILED"}
                                        </span>
                                    </div>
                                    <div className="text-sm font-semibold text-white mb-2">
                                        {entry.error.slice(0, 100)}...
                                    </div>
                                    {entry.fix && (
                                        <div className="bg-black/50 p-2 rounded text-xs overflow-x-auto">
                                            <div className="text-blue-300 mb-1">Fixed File: {entry.fix.diagnosis?.file}</div>
                                            <pre className="text-green-300">{entry.fix.diagnosis?.suggestedFix}</pre>
                                        </div>
                                    )}
                                </div>
                            ))
                        )}
                    </div>
                </section>
            </div>
        </div>
    );
}
