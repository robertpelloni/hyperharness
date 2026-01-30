import { VectorProvider, Document, SearchResult } from '../interfaces/VectorProvider.js';
import path from 'path';

export class MemoryManager {
    private provider: VectorProvider | null = null;
    private initialized: boolean = false;
    private dbPath: string;

    constructor(workspaceRoot: string = process.cwd()) {
        this.dbPath = path.join(workspaceRoot, '.borg', 'db');
    }

    private async initialize() {
        if (this.initialized) return;

        console.log("[MemoryManager] Initializing Vector Backend...");

        // Lazy load the default LanceDB provider from @borg/memory
        // In the future, this could allow switching to Chroma/SQLite based on config
        const { VectorStore } = await import('@borg/memory');

        // Adapt existing VectorStore to VectorProvider interface
        // @ts-ignore - The types might not match perfectly yet, acting as an adapter
        const store = new VectorStore(this.dbPath);

        this.provider = {
            initialize: async () => store.initialize(),
            add: async (docs: Document[]) => {
                await store.addDocuments(docs.map(d => ({
                    id: d.id,
                    file_path: (d.metadata as any)?.file_path || d.id,
                    content: d.content,
                    hash: 'dynamic', // TODO: Calculate hash
                    vector: d.vector
                })));
            },
            search: async (query: string, limit: number = 5) => {
                const results = await store.search(query, limit);
                return results.map((r: any) => ({
                    id: r.id,
                    content: r.content,
                    metadata: { file_path: r.file_path, hash: r.hash },
                    score: 0 // LanceDB wrapper didn't return score in previous interface, need update if needed
                }));
            },
            delete: async (ids: string[]) => {
                // TODO: Implement delete in VectorStore
            },
            reset: async () => store.reset()
        };

        if (this.provider) {
            await this.provider.initialize();
        }

        this.initialized = true;
        console.log("[MemoryManager] Ready.");
    }

    public async saveContext(content: string, metadata: any = {}) {
        if (!this.initialized) await this.initialize();
        if (!this.provider) throw new Error("Provider failed to init");

        const docId = `ctx/${Date.now()}/${Math.random().toString(36).substring(7)}`;

        // TODO: Integrate CodeSplitter here for better chunking

        await this.provider.add([{
            id: docId,
            content: content,
            metadata: metadata
        }]);

        return docId;
    }

    public async search(query: string, limit: number = 5): Promise<SearchResult[]> {
        if (!this.initialized) await this.initialize();
        if (!this.provider) return [];

        return await this.provider.search(query, limit);
    }
}
