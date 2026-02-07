
"use client";

import { trpcc } from "@/utils/trpc";
import { useState } from "react";
import { motion } from "framer-motion";

export default function KnowledgeDashboard() {
    const submodulesQuery = trpcc.submodule.list.useQuery();
    const resourcesQuery = trpcc.knowledge.getResources.useQuery();
    const updateAllMutation = trpcc.submodule.updateAll.useMutation();

    const [updating, setUpdating] = useState(false);

    const handleUpdateAll = async () => {
        setUpdating(true);
        try {
            await updateAllMutation.mutateAsync();
        } finally {
            setUpdating(false);
            submodulesQuery.refetch();
        }
    };

    const installMutation = trpcc.submodule.install.useMutation();
    const buildMutation = trpcc.submodule.build.useMutation();
    const enableMutation = trpcc.submodule.enable.useMutation();
    const [actioning, setActioning] = useState<string | null>(null);

    const handleAction = async (action: 'install' | 'build' | 'enable', path: string) => {
        setActioning(`${action}-${path}`);
        try {
            if (action === 'install') await installMutation.mutateAsync({ path });
            if (action === 'build') await buildMutation.mutateAsync({ path });
            if (action === 'enable') await enableMutation.mutateAsync({ path });
            submodulesQuery.refetch();
        } catch (e) {
            alert(`Failed to ${action}: ${e}`);
        } finally {
            setActioning(null);
        }
    };

    const submodules = submodulesQuery.data || [];
    const resources = resourcesQuery.data || { categories: [] };

    // Group resources for display
    const mcpServers = resources.categories.find((c: any) => c.name.includes("MCP Directories") || c.name.includes("MCPs"))?.urls || [];
    const skills = resources.categories.find((c: any) => c.name.includes("Skills"))?.urls || [];

    return (
        <div className="p-8 bg-gray-900 min-h-screen text-gray-100 font-mono">
            <header className="mb-8 border-b border-gray-700 pb-4 flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold text-blue-400">KNOWLEDGE BASE</h1>
                    <p className="text-gray-400">External Intelligence & Submodules</p>
                </div>
                <div className="text-right">
                    <div className="text-sm text-gray-500">LAST SYNC</div>
                    <div className="text-blue-500 font-bold">{resources.lastUpdated ? new Date(resources.lastUpdated).toLocaleString() : 'Never'}</div>
                </div>
            </header>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                {/* Submodules Section */}
                <section className="bg-gray-800 rounded-lg border border-gray-700 p-6">
                    <div className="flex justify-between items-center mb-4">
                        <h2 className="text-xl font-bold text-yellow-400 flex items-center">
                            <span className="mr-2">📦</span> Submodules ({submodules.length})
                        </h2>
                        <button
                            onClick={handleUpdateAll}
                            disabled={updating}
                            className={`px-3 py-1 bg-yellow-600 hover:bg-yellow-500 rounded text-sm font-bold transition-colors ${updating ? 'opacity-50 cursor-not-allowed' : ''}`}
                        >
                            {updating ? 'SYNCING...' : 'SYNC ALL'}
                        </button>
                    </div>

                    <div className="overflow-y-auto max-h-[600px] space-y-2">
                        {submodules.length === 0 ? (
                            <div className="text-gray-500 italic">No submodules detected or git status failed.</div>
                        ) : (
                            submodules.map((sub: any) => (
                                <div key={sub.path} className="bg-gray-900 p-3 rounded border-l-2 border-gray-700 hover:border-yellow-500 transition-colors">
                                    <div className="flex justify-between items-start">
                                        <div>
                                            <span className="font-bold text-blue-300 block">{sub.name}</span>
                                            <span className="text-gray-500 text-xs block">{sub.path}</span>
                                        </div>
                                        <span className={`text-xs px-2 py-0.5 rounded ${sub.status === 'clean' ? 'bg-green-900 text-green-300' : 'bg-red-900 text-red-300'}`}>
                                            {sub.status.toUpperCase()}
                                        </span>
                                    </div>

                                    <div className="mt-2 flex flex-wrap gap-2">
                                        {sub.capabilities?.map((cap: string) => (
                                            <span key={cap} className="text-[10px] bg-purple-900 text-purple-200 px-1 rounded border border-purple-700">
                                                {cap}
                                            </span>
                                        ))}
                                    </div>

                                    <div className="mt-3 flex gap-2 border-t border-gray-800 pt-2">
                                        <button
                                            onClick={() => handleAction('install', sub.path)}
                                            disabled={!!actioning}
                                            className={`text-xs px-2 py-1 rounded border ${sub.isInstalled ? 'border-green-800 text-green-500' : 'border-gray-600 text-gray-300 hover:bg-gray-800'}`}
                                        >
                                            {actioning === `install-${sub.path}` ? '...' : (sub.isInstalled ? 'Installed' : 'Install')}
                                        </button>
                                        <button
                                            onClick={() => handleAction('build', sub.path)}
                                            disabled={!!actioning}
                                            className={`text-xs px-2 py-1 rounded border ${sub.isBuilt ? 'border-green-800 text-green-500' : 'border-gray-600 text-gray-300 hover:bg-gray-800'}`}
                                        >
                                            {actioning === `build-${sub.path}` ? '...' : 'Build'}
                                        </button>
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                </section>

                {/* Resources Section - Categorized */}
                <section className="bg-gray-800 rounded-lg border border-gray-700 p-6 flex flex-col gap-6">

                    {resources.categories.map((category: any) => (
                        <div key={category.name}>
                            <h2 className="text-xl font-bold mb-2 text-purple-400 border-b border-gray-700 pb-1">
                                {category.name} <span className="text-gray-500 text-sm">({category.urls.length})</span>
                            </h2>
                            <div className="overflow-y-auto max-h-[200px] space-y-1">
                                {category.urls.map((url: string) => (
                                    <a href={url} target="_blank" rel="noreferrer" key={url} className="block text-sm text-gray-300 hover:text-white truncate">
                                        🔗 {url}
                                    </a>
                                ))}
                            </div>
                        </div>
                    ))}

                </section>
            </div>
        </div>
    );
}
