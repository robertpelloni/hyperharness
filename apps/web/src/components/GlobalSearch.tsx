"use client";
import React, { useState } from 'react';
import { trpc } from '@/utils/trpc';
import { motion, AnimatePresence } from 'framer-motion';

export const GlobalSearch: React.FC = () => {
    const [query, setQuery] = useState('');
    const [results, setResults] = useState<any[]>([]);
    const [isOpen, setIsOpen] = useState(false);

    // @ts-ignore
    const searchMutation = trpc.search.query.useMutation({
        onSuccess: (data: any) => {
            // executeTool('search_codebase') returns { content: [{ type, text }] }
            // We usually parse it. If it returns standard string, we might need to parse lines.
            // If the tool returns JSON, great.
            // Assuming `search_codebase` returns text with "File: ..." lines.

            // Let's assume data.content is the output array from MCP tool
            const text = data.content?.[0]?.text || "";
            // Parse simplistic file paths
            const lines = text.split('\n');
            const parsed = [];
            let currentFile = "";
            let snippet = "";

            // Heuristic Parsing of ripgrep/grep output or search_codebase output
            // Usually: "packages/core/foo.ts:10: const bar = 1;"

            for (const line of lines) {
                if (line.includes(':')) {
                    const parts = line.split(':');
                    // Very rough heuristic
                    if (parts[0].includes('/') || parts[0].includes('\\') || parts[0].includes('.')) {
                        parsed.push({
                            file: parts[0].trim(),
                            snippet: line.substring(parts[0].length + 1).trim().substring(0, 100)
                        });
                    }
                }
            }

            // If vector search returns structured data, adapt here.
            // Fallback for "No matches"
            if (parsed.length === 0 && text.length > 0) {
                parsed.push({ file: 'No exact file match', snippet: text.substring(0, 100) });
            }

            setResults(parsed.slice(0, 10)); // Top 10
        }
    });

    // @ts-ignore
    const openFile = trpc.vscode.open.useMutation({
        onSuccess: () => {
            setIsOpen(false);
            setQuery('');
        }
    });

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        if (!query.trim()) return;
        searchMutation.mutate({ query });
        setIsOpen(true);
    };

    return (
        <div className="relative z-50">
            <form onSubmit={handleSearch} className="relative">
                <input
                    type="text"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    placeholder="Search codebase..."
                    className="w-64 bg-zinc-100 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-full px-4 py-1.5 text-sm text-zinc-800 dark:text-zinc-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:w-96 transition-all"
                />
                <button type="submit" className="absolute right-3 top-1.5 text-zinc-400 hover:text-blue-500">
                    🔍
                </button>
            </form>

            <AnimatePresence>
                {isOpen && (results.length > 0 || (searchMutation as any).isLoading) && (
                    <motion.div
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0 }}
                        className="absolute right-0 top-12 w-96 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg shadow-2xl overflow-hidden"
                    >
                        {(searchMutation as any).isLoading ? (
                            <div className="p-4 text-center text-zinc-500 text-sm">Searching vector index...</div>
                        ) : (
                            <div className="max-h-96 overflow-y-auto">
                                {results.map((res, i) => (
                                    <button
                                        key={i}
                                        onClick={() => openFile.mutate({ path: res.file })}
                                        className="w-full text-left p-3 hover:bg-zinc-100 dark:hover:bg-zinc-800 border-b border-zinc-100 dark:border-zinc-800 last:border-0 transition-colors"
                                    >
                                        <div className="text-xs font-bold text-blue-500 break-all">{res.file}</div>
                                        <div className="text-xs text-zinc-500 mt-1 line-clamp-2 font-mono bg-zinc-50 dark:bg-black p-1 rounded">
                                            {res.snippet}
                                        </div>
                                    </button>
                                ))}
                            </div>
                        )}
                        <div className="bg-zinc-50 dark:bg-black p-2 text-[10px] text-center text-zinc-400 border-t border-zinc-200 dark:border-zinc-800">
                            Press ESC to close
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>

            {/* Backdrop to close */}
            {isOpen && (
                <div
                    className="fixed inset-0 z-[-1]"
                    onClick={() => setIsOpen(false)}
                />
            )}
        </div>
    );
};
