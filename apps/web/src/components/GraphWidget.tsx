"use client";

import React, { useMemo } from 'react';
import { trpc } from '@/utils/trpc';
import { KnowledgeGraph } from '@borg/ui';

export function GraphWidget() {
    // Poll every 10s for updates
    // Using repoGraph directly from tRPC router
    const { data, isLoading } = trpc.repoGraph.get.useQuery(undefined, {
        refetchInterval: 10000,
        refetchOnWindowFocus: false
    });

    const mappedData = useMemo(() => {
        if (!data) return { nodes: [], links: [] };

        // Map RepoGraph D3 format to KnowledgeGraph format
        // RepoGraph usually returns: { id, group }
        // KnowledgeGraph expects: { id, label, type, val }
        return {
            nodes: data.nodes.map((n: any) => ({
                id: n.id,
                label: n.name || n.id.split('/').pop() || n.id,
                type: n.group === 1 ? 'topic' : (n.group === 2 ? 'document' : 'concept'), // Heuristic mapping
                val: n.val || 5
            })),
            links: data.links.map((l: any) => ({
                source: l.source,
                target: l.target,
                value: l.value || 1
            }))
        };
    }, [data]);

    // @ts-ignore
    const openFile = trpc.vscode.open.useMutation();

    return (
        <div className="h-full w-full min-h-[300px] flex flex-col overflow-hidden">
            <div className="flex-1 relative">
                <KnowledgeGraph
                    nodes={mappedData.nodes}
                    links={mappedData.links}
                    loading={isLoading}
                    onNodeClick={(node) => {
                        // Open file in VS Code on click
                        // Only open if it looks like a file or is classified as document
                        if (node.type === 'document' || (node.id && node.id.includes('.'))) {
                            openFile.mutate({ path: node.id });
                        }
                    }}
                />
            </div>
        </div>
    );
}
