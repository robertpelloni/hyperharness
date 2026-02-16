"use client";

import { useState } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from "@borg/ui";
import { Loader2, Search, ArrowRight, Zap, Code } from "lucide-react";
import { trpc } from '@/utils/trpc';

export default function SearchDashboard() {
    const [query, setQuery] = useState('');
    const searchQuery = trpc.tools.search.useQuery(
        { query, limit: 30 },
        { enabled: query.trim().length > 0 }
    );

    const results = searchQuery.data || [];
    const isLoading = searchQuery.isLoading;

    return (
        <div className="p-8 space-y-8 h-full flex flex-col">
            <div className="flex justify-between items-center shrink-0">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight text-white">Semantic Search</h1>
                    <p className="text-zinc-500">
                        Find tools and capabilities using natural language
                    </p>
                </div>
            </div>

            <div className="max-w-2xl mx-auto w-full space-y-8 mt-12">
                <div className="relative">
                    <Search className="absolute left-6 top-1/2 -translate-y-1/2 h-6 w-6 text-zinc-500" />
                    <input
                        value={query}
                        onChange={(e) => setQuery(e.target.value)}
                        placeholder="What do you want to achieve? (e.g. 'process csv files')"
                        className="w-full bg-zinc-900 border border-zinc-800 rounded-full p-6 pl-16 text-lg text-white focus:ring-2 focus:ring-blue-500 outline-none shadow-xl"
                        autoFocus
                    />
                </div>

                <div className="space-y-4">
                    {isLoading ? (
                        <div className="flex justify-center p-8">
                            <Loader2 className="h-6 w-6 animate-spin text-zinc-500" />
                        </div>
                    ) : results.length > 0 ? (
                        results.map((tool: any) => (
                            <Card key={tool.uuid} className="bg-zinc-900/50 border-zinc-800 hover:bg-zinc-900 transition-colors group cursor-pointer">
                                <CardContent className="p-6">
                                    <div className="flex justify-between items-start">
                                        <div>
                                            <h3 className="font-mono text-blue-400 font-medium text-lg mb-1 flex items-center gap-2">
                                                {tool.name}
                                                <span className="text-[10px] bg-zinc-800 px-2 py-0.5 rounded text-zinc-500">
                                                    {tool.server}
                                                </span>
                                            </h3>
                                            <p className="text-zinc-400">{tool.description}</p>
                                        </div>
                                        <ArrowRight className="h-5 w-5 text-zinc-600 group-hover:text-blue-500 transition-colors opacity-0 group-hover:opacity-100" />
                                    </div>
                                </CardContent>
                            </Card>
                        ))
                    ) : query && (
                        <div className="text-center text-zinc-500 py-12">
                            No tools found matching "{query}"
                        </div>
                    )}

                    {!query && (
                        <div className="grid grid-cols-2 gap-4 opacity-50">
                            <div className="p-4 rounded border border-dashed border-zinc-800 text-sm text-zinc-500 flex items-center gap-3">
                                <Zap className="h-4 w-4" /> "Generate verify.sh scripts"
                            </div>
                            <div className="p-4 rounded border border-dashed border-zinc-800 text-sm text-zinc-500 flex items-center gap-3">
                                <Code className="h-4 w-4" /> "Refactor React components"
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
