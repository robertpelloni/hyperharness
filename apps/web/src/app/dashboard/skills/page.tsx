"use client";

import { trpc } from "@/utils/trpc";
import { useState, useMemo } from "react";
import { Search, Library, CheckCircle, Clock, Zap, ExternalLink, RefreshCw, FileText } from "lucide-react";

export default function SkillsPage() {
    const { data: library, isLoading, refetch } = trpc.skills.listLibrary.useQuery();
    const [searchTerm, setSearchTerm] = useState("");
    const [activeTab, setActiveTab] = useState<'assimilated' | 'backlog'>('assimilated');
    const [selectedItem, setSelectedItem] = useState<any | null>(null);

    const assimilateMutation = trpc.skills.assimilate.useMutation({
        onSuccess: () => {
            refetch();
            setSelectedItem(null);
            // In a real app we'd trigger a toast here
        }
    });

    const allItems = useMemo(() => {
        if (!library) return [];
        return [
            ...library.mcp_servers,
            ...library.universal_harness,
            ...library.skills
        ];
    }, [library]);

    const filteredItems = useMemo(() => {
        return allItems.filter(item => {
            const matchesSearch = item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                item.summary?.toLowerCase().includes(searchTerm.toLowerCase());
            const matchesTab = activeTab === 'assimilated' ? item.isAssimilated : !item.isAssimilated;
            return matchesSearch && matchesTab;
        });
    }, [allItems, searchTerm, activeTab]);

    return (
        <div className="space-y-6 max-w-7xl mx-auto">
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight text-white flex items-center gap-2">
                        <Library className="w-8 h-8 text-blue-500" />
                        Skills Library
                    </h1>
                    <p className="text-zinc-400 mt-1">The Great Absorption: Manage and expand Borg's autonomous capabilities.</p>
                </div>
                <div className="flex items-center gap-2">
                    <div className="relative">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-500" />
                        <input
                            type="text"
                            placeholder="Search library..."
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            className="pl-10 pr-4 py-2 bg-zinc-900 border border-zinc-800 rounded-lg text-sm text-white focus:outline-none focus:ring-2 focus:ring-blue-500/50 w-64"
                        />
                    </div>
                    <button
                        onClick={() => refetch()}
                        className="p-2 bg-zinc-900 border border-zinc-800 hover:bg-zinc-800 text-zinc-400 hover:text-white rounded-lg transition-all"
                        title="Refresh Index"
                    >
                        <RefreshCw className={`w-5 h-5 ${isLoading ? 'animate-spin' : ''}`} />
                    </button>
                </div>
            </div>

            <div className="flex border-b border-zinc-800 gap-8">
                <button
                    onClick={() => setActiveTab('assimilated')}
                    className={`pb-4 text-sm font-medium transition-all relative ${activeTab === 'assimilated' ? 'text-blue-500' : 'text-zinc-500 hover:text-zinc-300'}`}
                >
                    Assimilated
                    {activeTab === 'assimilated' && <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-500" />}
                </button>
                <button
                    onClick={() => setActiveTab('backlog')}
                    className={`pb-4 text-sm font-medium transition-all relative ${activeTab === 'backlog' ? 'text-blue-500' : 'text-zinc-500 hover:text-zinc-300'}`}
                >
                    Research Backlog
                    {activeTab === 'backlog' && <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-500" />}
                </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {isLoading ? (
                    Array.from({ length: 6 }).map((_, i) => (
                        <div key={i} className="h-40 bg-zinc-900/50 border border-zinc-800 rounded-xl animate-pulse" />
                    ))
                ) : filteredItems.length === 0 ? (
                    <div className="col-span-full py-20 text-center text-zinc-500 bg-zinc-900/30 rounded-xl border border-dashed border-zinc-800">
                        {searchTerm ? "No tools match your search criteria." : `No items found in the ${activeTab} category.`}
                    </div>
                ) : (
                    filteredItems.map((item: any) => (
                        <div
                            key={item.id}
                            onClick={() => setSelectedItem(item)}
                            className="group p-5 bg-zinc-900/50 border border-zinc-800 hover:border-blue-500/50 hover:bg-zinc-900 rounded-xl cursor-pointer transition-all shadow-lg hover:shadow-blue-500/10 flex flex-col justify-between"
                        >
                            <div>
                                <div className="flex items-start justify-between mb-3">
                                    <div className="p-2 bg-zinc-950 rounded-lg border border-zinc-800 group-hover:border-blue-500/30 transition-colors">
                                        <Zap className={`w-5 h-5 ${item.isAssimilated ? 'text-blue-500' : 'text-zinc-600'}`} />
                                    </div>
                                    <span className={`px-2 py-1 text-[10px] font-bold uppercase rounded-full ${item.isAssimilated ? 'bg-blue-500/10 text-blue-500 border border-blue-500/20' : 'bg-zinc-800 text-zinc-500 border border-zinc-700'
                                        }`}>
                                        {item.isAssimilated ? 'Active' : 'Awaiting Ingest'}
                                    </span>
                                </div>
                                <h3 className="text-white font-semibold group-hover:text-blue-400 transition-colors">{item.name}</h3>
                                <p className="text-xs text-zinc-500 mt-1 line-clamp-2">{item.summary}</p>
                            </div>

                            <div className="mt-4 pt-4 border-t border-zinc-800 flex items-center justify-between text-xs text-zinc-500">
                                <span className="flex items-center gap-1 italic">
                                    {item.relevance?.split(' - ')[0] || 'Prioritized'}
                                </span>
                                <span className="flex items-center gap-1 text-blue-500/70 opacity-0 group-hover:opacity-100 transition-opacity">
                                    Click to view detail <ExternalLink className="w-3 h-3" />
                                </span>
                            </div>
                        </div>
                    ))
                )}
            </div>

            {selectedItem && (
                <div className="fixed inset-0 bg-black/80 backdrop-blur-md flex items-center justify-center p-4 z-50 animate-in fade-in duration-200" onClick={() => setSelectedItem(null)}>
                    <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-0 max-w-3xl w-full max-h-[85vh] overflow-hidden shadow-2xl flex flex-col" onClick={e => e.stopPropagation()}>
                        <div className="p-6 border-b border-zinc-800 flex justify-between items-start bg-zinc-900/50">
                            <div>
                                <div className="flex items-center gap-2 mb-1">
                                    <h2 className="text-2xl font-bold text-white">{selectedItem.name}</h2>
                                    <span className={`px-2 py-0.5 text-[10px] uppercase rounded-full border ${selectedItem.isAssimilated ? 'bg-blue-500/10 text-blue-500 border-blue-500/20' : 'bg-zinc-800 text-zinc-500 border-zinc-700'
                                        }`}>
                                        {selectedItem.isAssimilated ? 'Assimilated' : 'Awaiting Ingest'}
                                    </span>
                                </div>
                                <a href={selectedItem.url} target="_blank" className="text-xs text-blue-500 hover:underline flex items-center gap-1">
                                    {selectedItem.url} <ExternalLink className="w-3 h-3" />
                                </a>
                            </div>
                            <button onClick={() => setSelectedItem(null)} className="p-2 hover:bg-zinc-800 rounded-lg text-zinc-500 hover:text-white transition-colors">✕</button>
                        </div>

                        <div className="flex-1 overflow-y-auto p-6 space-y-6">
                            <section>
                                <h4 className="text-xs font-bold text-zinc-500 uppercase tracking-widest mb-2">Summary</h4>
                                <p className="text-zinc-300 leading-relaxed bg-zinc-950/50 p-4 rounded-xl border border-zinc-800">{selectedItem.summary}</p>
                            </section>

                            <section>
                                <h4 className="text-xs font-bold text-zinc-500 uppercase tracking-widest mb-2">Strategic Relevance</h4>
                                <p className="text-zinc-300 italic">{selectedItem.relevance}</p>
                            </section>

                            {selectedItem.isAssimilated ? (
                                <section>
                                    <h4 className="text-xs font-bold text-zinc-500 uppercase tracking-widest mb-2 flex items-center gap-2">
                                        <FileText className="w-3 h-3" /> Installed Runbook (SKILL.md)
                                    </h4>
                                    <SkillDetails skillName={selectedItem.id} />
                                </section>
                            ) : (
                                <div className="p-6 bg-blue-500/5 rounded-2xl border border-blue-500/20 text-center space-y-4">
                                    <Zap className="w-10 h-10 text-blue-500 mx-auto" />
                                    <div>
                                        <h3 className="text-white font-bold text-lg">Knowledge Assimilation Available</h3>
                                        <p className="text-sm text-zinc-400 max-w-md mx-auto mt-1">This tool has been researched but not yet ingested. Triggering assimilation will generate a functional Borg Skill runbook using LLM synthesis.</p>
                                    </div>
                                    <button
                                        onClick={() => assimilateMutation.mutate(selectedItem)}
                                        disabled={assimilateMutation.isLoading}
                                        className="px-8 py-3 bg-blue-600 hover:bg-blue-500 disabled:bg-zinc-800 text-white rounded-xl font-bold transition-all shadow-lg shadow-blue-500/20 flex items-center gap-2 mx-auto"
                                    >
                                        {assimilateMutation.isLoading ? (
                                            <>
                                                <RefreshCw className="w-5 h-5 animate-spin" />
                                                Assimilating Knowledge...
                                            </>
                                        ) : (
                                            <>
                                                <Zap className="w-5 h-5 fill-current" />
                                                Trigger Assimilation
                                            </>
                                        )}
                                    </button>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

function SkillDetails({ skillName }: { skillName: string }) {
    const { data: details, isLoading } = trpc.skills.read.useQuery({ name: skillName });

    if (isLoading) return <div className="py-10 text-center text-zinc-500">Retrieving runbook...</div>;
    if (!details || details.content[0].text.includes("not found")) return <div className="p-4 bg-red-900/10 text-red-500 rounded-lg text-sm border border-red-500/20">Skill runbook not found on disk.</div>;

    const content = details.content[0].text;

    return (
        <div className="bg-zinc-950 border border-zinc-800 rounded-xl overflow-hidden">
            <div className="p-4 bg-zinc-900/50 border-b border-zinc-800 flex items-center justify-between">
                <span className="text-[10px] font-mono text-zinc-500 lowercase">{skillName}/SKILL.md</span>
                <CheckCircle className="w-4 h-4 text-green-500" />
            </div>
            <pre className="p-6 overflow-x-auto text-sm font-mono text-zinc-300 whitespace-pre-wrap leading-relaxed max-h-[400px]">
                {content}
            </pre>
        </div>
    );
}
