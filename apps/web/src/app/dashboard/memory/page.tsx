'use client';

import React, { useState } from 'react';
import { trpc } from '@/utils/trpc';
import { MemoryGraph } from '@/components/memory/MemoryGraph';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Search, Brain, Database, Share2 } from 'lucide-react';

export default function MemoryPage() {
    const [query, setQuery] = useState('');
    const [activeTab, setActiveTab] = useState('graph');

    const { data: stats } = trpc.knowledge.getStats.useQuery();
    const { data: graphData, refetch } = trpc.knowledge.getGraph.useQuery({ query });

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        refetch();
    };

    return (
        <div className="container mx-auto p-6 space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Memory & Knowledge</h1>
                    <p className="text-muted-foreground">Vizualize the agent's semantic memory.</p>
                </div>
                <div className="flex gap-4">
                    <Card className="w-[150px]">
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium">Total Memories</CardTitle>
                            <Database className="h-4 w-4 text-muted-foreground" />
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">{stats?.count || 0}</div>
                        </CardContent>
                    </Card>
                </div>
            </div>

            <div className="flex items-center space-x-2">
                <form onSubmit={handleSearch} className="flex w-full max-w-sm items-center space-x-2">
                    <Input
                        placeholder="Search memory..."
                        value={query}
                        onChange={(e) => setQuery(e.target.value)}
                    />
                    <Button type="submit" size="icon">
                        <Search className="h-4 w-4" />
                    </Button>
                </form>
            </div>

            <Tabs defaultValue="graph" className="space-y-4" onValueChange={setActiveTab}>
                <TabsList>
                    <TabsTrigger value="graph" className="flex gap-2"><Share2 className="w-4 h-4" /> Graph View</TabsTrigger>
                    <TabsTrigger value="list" className="flex gap-2"><Brain className="w-4 h-4" /> List View</TabsTrigger>
                </TabsList>
                <TabsContent value="graph" className="space-y-4">
                    <Card className="h-[600px] overflow-hidden relative border-zinc-800 bg-black/50">
                        <MemoryGraph data={graphData || { nodes: [], edges: [] }} />
                    </Card>
                </TabsContent>
                <TabsContent value="list">
                    <Card>
                        <CardContent className="pt-6">
                            <p className="text-muted-foreground">List view coming soon. Use the Graph view to explore relationships.</p>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    );
}
