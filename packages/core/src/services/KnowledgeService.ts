
import { MemoryManager } from './MemoryManager.js';


export interface GraphNode {
    id: string;
    label: string;
    type: 'topic' | 'document' | 'concept';
    val: number; // Size/Importance
}

export interface GraphEdge {
    source: string;
    target: string;
    value: number; // Similarity/Weight
}

export interface GraphData {
    nodes: GraphNode[];
    edges: GraphEdge[];
}

export class KnowledgeService {
    private memory: MemoryManager;

    constructor(memory: MemoryManager) {
        this.memory = memory;
    }


    public async getGraph(query?: string, depth: number = 1): Promise<{ content: any[] }> {
        let nodes: GraphNode[] = [];
        let edges: GraphEdge[] = [];

        // 1. Get Root Nodes (Recent or Query-based)
        // Combine Vector Search with Graph Search
        const vectorDocs = await this.memory.search(query || "AI", 10);

        let graphResults: string[] = [];
        if (this.memory.graph) {
            try {
                graphResults = await this.memory.graph.search(query || "AI");
            } catch (e) {
                console.error("Graph search failed (bridge might be offline):", e);
            }
        }

        // 2. Map to Graph Format
        // Process Vector Results
        vectorDocs.forEach(doc => {
            const nodeId = doc.id || doc.metadata?.source || "unknown";
            nodes.push({
                id: nodeId,
                label: (doc.metadata?.title as string) || "Untitled " + nodeId.substring(0, 8),
                type: 'document',
                val: 5
            });

            // Mock edge to central query if exists
            if (query) {
                edges.push({
                    source: query,
                    target: nodeId,
                    value: 1
                });
            }
        });

        // Process Graph Results (Cognee returns formatted strings currently, we'd parse them if they were structured)
        // For now, treat them as insights linked to the query
        graphResults.forEach((res, idx) => {
            const nodeId = `insight-${idx}`;
            nodes.push({
                id: nodeId,
                label: res.substring(0, 50) + "...",
                type: 'concept',
                val: 8
            });
            if (query) {
                edges.push({ source: query, target: nodeId, value: 2 });
            }
        });

        if (query) {
            nodes.push({ id: query, label: query, type: 'topic', val: 10 });
        }

        return {
            content: [{
                type: "text",
                text: JSON.stringify({ nodes, edges })
            }]
        };
    }
}
