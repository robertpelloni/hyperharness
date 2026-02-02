"use client";
import React from 'react';
import { trpc } from '../utils/trpc';
import { motion } from 'framer-motion';

export function TestStatusWidget() {
    const { data: results } = trpc.autoTest.getResults.useQuery(undefined, { refetchInterval: 2000 });
    const { data: status, refetch: refetchStatus } = trpc.autoTest.getStatus.useQuery();

    const utils = trpc.useContext();
    const startMutation = trpc.autoTest.start.useMutation({
        onSuccess: () => { refetchStatus(); }
    });
    const stopMutation = trpc.autoTest.stop.useMutation({
        onSuccess: () => { refetchStatus(); }
    });
    // @ts-ignore
    const clearMutation = (trpc as any).autoTest.clear.useMutation({
        onSuccess: () => {
            // @ts-ignore
            (utils as any).autoTest.getResults.invalidate();
        }
    });

    const passCount = results?.filter((r: any) => r.status === 'pass').length || 0;
    const failCount = results?.filter((r: any) => r.status === 'fail').length || 0;
    const runningCount = results?.filter((r: any) => r.status === 'running').length || 0;

    return (
        <div className="relative overflow-hidden rounded-xl bg-gradient-to-br from-zinc-900 to-black border border-zinc-800 p-4 h-full flex flex-col">
            {/* Background Glow */}
            <div className={`absolute inset-0 opacity-10 blur-3xl ${failCount > 0 ? 'bg-red-500' : 'bg-green-500'
                }`} />

            <div className="relative z-10 flex flex-col h-full">
                {/* Header */}
                <div className="flex justify-between items-center mb-4">
                    <h3 className="text-lg font-bold text-white flex items-center gap-2">
                        <span className="text-2xl">🧪</span>
                        Auto-Test Watcher
                    </h3>
                    <div className="flex gap-2">
                        {!status?.isRunning ? (
                            <button
                                onClick={() => startMutation.mutate()}
                                className="bg-green-600 hover:bg-green-500 text-white px-3 py-1 rounded-lg text-xs font-bold transition-all shadow-lg shadow-green-500/20"
                            >
                                ▶ Start
                            </button>
                        ) : (
                            <button
                                onClick={() => stopMutation.mutate()}
                                className="bg-yellow-600 hover:bg-yellow-500 text-white px-3 py-1 rounded-lg text-xs font-bold transition-all"
                            >
                                ⏹ Stop
                            </button>
                        )}
                        <button
                            onClick={() => clearMutation.mutate()}
                            className="text-zinc-500 hover:text-zinc-300 text-xs px-2 py-1 bg-zinc-800 hover:bg-zinc-700 rounded-lg transition-all"
                        >
                            Clear
                        </button>
                    </div>
                </div>

                {/* Stats Bar */}
                <div className="grid grid-cols-3 gap-2 mb-4">
                    <div className="bg-green-500/10 border border-green-500/30 rounded-lg p-2 text-center">
                        <div className="text-xl font-bold text-green-400">{passCount}</div>
                        <div className="text-[10px] text-green-300/70 uppercase">Passed</div>
                    </div>
                    <div className="bg-red-500/10 border border-red-500/30 rounded-lg p-2 text-center">
                        <div className="text-xl font-bold text-red-400">{failCount}</div>
                        <div className="text-[10px] text-red-300/70 uppercase">Failed</div>
                    </div>
                    <div className="bg-blue-500/10 border border-blue-500/30 rounded-lg p-2 text-center">
                        <div className="text-xl font-bold text-blue-400">{runningCount}</div>
                        <div className="text-[10px] text-blue-300/70 uppercase">Running</div>
                    </div>
                </div>

                {/* Results List */}
                <div className="flex-1 overflow-y-auto space-y-2 font-mono text-xs custom-scrollbar min-h-0">
                    {results && results.length > 0 ? (
                        results.map((r: any, i: number) => (
                            <motion.div
                                key={r.file}
                                initial={{ opacity: 0, y: 10 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ delay: i * 0.05 }}
                                className={`p-3 rounded-lg border ${r.status === 'pass' ? 'bg-green-500/5 border-green-500/20' :
                                        r.status === 'fail' ? 'bg-red-500/5 border-red-500/20' :
                                            'bg-blue-500/5 border-blue-500/20'
                                    }`}
                            >
                                <div className="flex justify-between items-center mb-1">
                                    <span className="text-zinc-300 truncate font-medium" title={r.file}>
                                        {r.file.split(/[\\/]/).pop()}
                                    </span>
                                    <span className={`font-bold text-[10px] px-2 py-0.5 rounded-full ${r.status === 'pass' ? 'bg-green-500/20 text-green-400' :
                                            r.status === 'fail' ? 'bg-red-500/20 text-red-400' :
                                                'bg-blue-500/20 text-blue-400 animate-pulse'
                                        }`}>
                                        {r.status.toUpperCase()}
                                    </span>
                                </div>
                                <div className="text-[10px] text-zinc-600 truncate">{r.file}</div>
                                {r.status === 'fail' && r.output && (
                                    <pre className="mt-2 text-red-300 bg-black/50 p-2 rounded-lg overflow-x-auto text-[10px] max-h-20 overflow-y-auto">
                                        {r.output.substring(0, 300)}...
                                    </pre>
                                )}
                            </motion.div>
                        ))
                    ) : (
                        <div className="flex flex-col items-center justify-center h-full text-zinc-600">
                            {status?.isRunning ? (
                                <>
                                    <motion.div
                                        animate={{ rotate: 360 }}
                                        transition={{ repeat: Infinity, duration: 2, ease: "linear" }}
                                        className="w-10 h-10 border-4 border-pink-500 border-t-transparent rounded-full mb-3"
                                    />
                                    <p className="text-sm">Watching for changes...</p>
                                </>
                            ) : (
                                <>
                                    <span className="text-4xl mb-2">🧪</span>
                                    <p className="text-sm">Watcher stopped</p>
                                    <p className="text-[10px]">Click Start to begin</p>
                                </>
                            )}
                        </div>
                    )}
                </div>

                {/* Footer */}
                <div className="mt-3 pt-3 border-t border-zinc-800 flex justify-between items-center">
                    <span className={`text-[10px] px-2 py-1 rounded-full ${status?.isRunning ? 'bg-green-500/20 text-green-400' : 'bg-zinc-800 text-zinc-500'
                        }`}>
                        {status?.isRunning ? '● WATCHING' : '○ STOPPED'}
                    </span>
                    <span className="text-[10px] text-zinc-600">
                        {results?.length || 0} test results
                    </span>
                </div>
            </div>
        </div>
    );
}
