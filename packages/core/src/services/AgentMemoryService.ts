/**
 * Agent Memory Service - Tiered Memory System for Persistent Agent Context
 * 
 * Implements multi-tier memory architecture inspired by Mem0 and Letta:
 * - Session Memory: Ephemeral context within a conversation
 * - Working Memory: Task-relevant facts extracted during execution
 * - Long-term Memory: Persistent learnings across sessions
 * 
 * Features:
 * - Automatic memory extraction from conversations
 * - Relevance-based retrieval with temporal decay
 * - User/Agent/Project namespaces
 * - Memory consolidation (working -> long-term)
 */

import * as fs from 'fs';
import * as path from 'path';
import * as crypto from 'crypto';
import { MemoryManager } from './MemoryManager.js';

// Memory types
export type MemoryType = 'session' | 'working' | 'long_term';
export type MemoryNamespace = 'user' | 'agent' | 'project';

export interface Memory {
    id: string;
    type: MemoryType;
    namespace: MemoryNamespace;
    content: string;
    media?: MemoryMedia[]; // Multi-modal support
    metadata: MemoryMetadata;
    createdAt: Date;
    accessedAt: Date;
    accessCount: number;
    score?: number;  // Relevance score for retrieval
    ttl?: number;    // Time-to-live in ms (for session memory)
}

export interface MemoryMedia {
    type: 'image' | 'audio' | 'video' | 'file';
    url?: string;        // External or Local path
    dataUrl?: string;   // Base64 encoded
    blobHash?: string;  // Reference to binary storage
    mimeType?: string;
    description?: string; // OCR text or AI caption
}

export interface MemoryMetadata {
    source?: string;          // Where memory came from (tool, conversation, etc)
    tags?: string[];          // Categorization tags
    relatedMemories?: string[]; // Links to other memories
    confidence?: number;      // Extraction confidence (0-1)
    userId?: string;          // Associated user
    projectId?: string;       // Associated project
    sessionId?: string;       // Associated session ID
    [key: string]: unknown;
}

export interface MemorySearchOptions {
    namespace?: MemoryNamespace;
    type?: MemoryType;
    limit?: number;
    minScore?: number;
    tags?: string[];
    includeExpired?: boolean;
}

export interface MemoryServiceOptions {
    persistDir: string;
    sessionTTL?: number;        // Default session TTL (30 min)
    consolidationThreshold?: number;  // Access count to promote to long-term
    maxSessionMemories?: number;
    maxWorkingMemories?: number;
}

/**
 * Simple in-memory vector similarity using TF-IDF-like approach
 * Production would use proper embeddings
 */
// SimpleVectorSearch replaced by MemoryManager

/**
 * Agent Memory Service - Main service class
 */
export class AgentMemoryService {
    private memories: Map<string, Memory> = new Map();
    private memoryManager: MemoryManager;
    private options: Required<MemoryServiceOptions>;
    private dirty = false;

    constructor(options: MemoryServiceOptions, memoryManager?: MemoryManager) {
        this.options = {
            persistDir: options.persistDir,
            sessionTTL: options.sessionTTL ?? 30 * 60 * 1000,  // 30 minutes
            consolidationThreshold: options.consolidationThreshold ?? 5,
            maxSessionMemories: options.maxSessionMemories ?? 100,
            maxWorkingMemories: options.maxWorkingMemories ?? 500,
        };

        if (memoryManager) {
            this.memoryManager = memoryManager;
        } else {
            // Initialize MemoryManager (parent of persistDir is workspaceRoot)
            const workspaceRoot = path.dirname(this.options.persistDir);
            this.memoryManager = new MemoryManager(workspaceRoot);
        }

        this.loadFromDisk();
    }

    /**
     * Generate unique memory ID
     */
    private generateId(): string {
        return crypto.randomBytes(8).toString('hex');
    }

    /**
     * Load memories from disk
     */
    private loadFromDisk(): void {
        const filePath = path.join(this.options.persistDir, 'memories.json');

        if (!fs.existsSync(this.options.persistDir)) {
            fs.mkdirSync(this.options.persistDir, { recursive: true });
        }

        if (fs.existsSync(filePath)) {
            try {
                const data = JSON.parse(fs.readFileSync(filePath, 'utf-8'));
                for (const mem of data.memories || []) {
                    mem.createdAt = new Date(mem.createdAt);
                    mem.accessedAt = new Date(mem.accessedAt);

                    if (mem.type === 'session' && mem.ttl) {
                        const age = Date.now() - mem.createdAt.getTime();
                        if (age > mem.ttl) continue;
                    }

                    this.memories.set(mem.id, mem);
                    // No need to manually add to vectorIndex here, 
                    // MemoryManager handles persistence of long-term/working.
                }
                console.log(`[AgentMemoryService] Loaded ${this.memories.size} session memories`);
            } catch (e) {
                console.error('[AgentMemoryService] Failed to load memories:', e);
            }
        }
    }

    /**
     * Save memories to disk
     */
    private saveToDisk(): void {
        if (!this.dirty) return;

        const filePath = path.join(this.options.persistDir, 'memories.json');
        const data = {
            version: 1,
            savedAt: new Date().toISOString(),
            memories: Array.from(this.memories.values()),
        };

        fs.writeFileSync(filePath, JSON.stringify(data, null, 2));
        this.dirty = false;
    }

    /**
     * Add a memory
     */
    async add(
        content: string,
        type: MemoryType,
        namespace: MemoryNamespace,
        metadata: MemoryMetadata = {}
    ): Promise<Memory> {
        const memory: Memory = {
            id: this.generateId(),
            type,
            namespace,
            content,
            metadata,
            createdAt: new Date(),
            accessedAt: new Date(),
            accessCount: 0,
            ttl: type === 'session' ? this.options.sessionTTL : undefined,
        };

        this.memories.set(memory.id, memory);
        this.dirty = true;

        // If working or long-term, also persist to vector DB
        if (type !== 'session') {
            await this.memoryManager.saveContext(content, {
                ...metadata,
                id: memory.id,
                type,
                namespace,
                createdAt: memory.createdAt.getTime()
            });
        }

        this.enforceMemoryLimits();
        this.scheduleSave();
        return memory;
    }

    /**
     * Get a memory by ID
     */
    get(id: string): Memory | null {
        const memory = this.memories.get(id);
        if (!memory) return null;

        // Check expiration
        if (memory.type === 'session' && memory.ttl) {
            const age = Date.now() - memory.createdAt.getTime();
            if (age > memory.ttl) {
                this.delete(id);
                return null;
            }
        }

        // Update access metadata
        memory.accessedAt = new Date();
        memory.accessCount++;
        this.dirty = true;

        // Check for consolidation
        this.checkConsolidation(memory);

        return memory;
    }

    /**
     * Delete a memory
     */
    async delete(id: string): Promise<boolean> {
        const existed = this.memories.delete(id);
        if (existed) {
            this.dirty = true;
        }

        // Also delete from vector DB
        try {
            await this.memoryManager.deleteContext(id);
        } catch (e) {
            // Might not exist in vector DB if it was only a session memory
        }

        return existed;
    }

    /**
     * Search memories using hybrid strategy (Local sessions + Vector DB)
     */
    async search(query: string, options: MemorySearchOptions = {}): Promise<Memory[]> {
        const limit = options.limit ?? 10;

        // 1. Get results from Vector DB
        const vectorResults = await this.memoryManager.search(query, limit * 2);
        const mappedResults: Memory[] = vectorResults.map(r => ({
            id: r.id,
            content: r.content,
            type: (r.metadata as any)?.type || 'long_term',
            namespace: (r.metadata as any)?.namespace || 'project',
            metadata: r.metadata as any,
            createdAt: new Date((r.metadata as any)?.createdAt || Date.now()),
            accessedAt: new Date(),
            accessCount: 0, // We could track this in metadata if needed
            score: r.score
        }));

        // 2. Add current session memories (which might not be indexed yet)
        const sessionMemories = Array.from(this.memories.values())
            .filter(m => m.type === 'session');

        // For session memories, we do a simple title/content match if needed, 
        // or just include the most recent ones. 
        // Better: Index session memories too? For now, we rely on the fact that 
        // they are small.

        let results = [...mappedResults, ...sessionMemories];

        // Filter by type/namespace
        if (options.type) {
            results = results.filter(m => m.type === options.type);
        }
        if (options.namespace) {
            results = results.filter(m => m.namespace === options.namespace);
        }
        if (options.tags?.length) {
            results = results.filter(m =>
                options.tags!.some(tag => m.metadata.tags?.includes(tag))
            );
        }

        // Simple relevance boost for session memories if they mention query terms
        const queryTerms = query.toLowerCase().split(/\s+/).filter(t => t.length > 2);
        for (const res of results) {
            if (res.type === 'session' && !res.score) {
                let sessionScore = 0;
                for (const term of queryTerms) {
                    if (res.content.toLowerCase().includes(term)) sessionScore += 1;
                }
                res.score = sessionScore;
            }
        }

        // Apply temporal decay
        const now = Date.now();
        for (const memory of results) {
            const ageHours = (now - memory.accessedAt.getTime()) / (1000 * 60 * 60);
            const decayFactor = Math.exp(-ageHours / 24);
            memory.score = (memory.score ?? 0) * (0.5 + 0.5 * decayFactor);
        }

        return results
            .sort((a, b) => (b.score ?? 0) - (a.score ?? 0))
            .slice(0, limit);
    }

    /**
     * Get recent memories
     */
    getRecent(limit: number = 10, options: MemorySearchOptions = {}): Memory[] {
        let memories = Array.from(this.memories.values());

        if (options.type) {
            memories = memories.filter(m => m.type === options.type);
        }
        if (options.namespace) {
            memories = memories.filter(m => m.namespace === options.namespace);
        }

        return memories
            .sort((a, b) => b.createdAt.getTime() - a.createdAt.getTime())
            .slice(0, limit);
    }

    /**
     * Add session memory (ephemeral)
     */
    async addSession(content: string, metadata: MemoryMetadata = {}): Promise<Memory> {
        return await this.add(content, 'session', 'agent', metadata);
    }

    /**
     * Add working memory (task-relevant)
     */
    async addWorking(content: string, namespace: MemoryNamespace = 'project', metadata: MemoryMetadata = {}): Promise<Memory> {
        return await this.add(content, 'working', namespace, metadata);
    }

    /**
     * Add long-term memory (persistent)
     */
    async addLongTerm(content: string, namespace: MemoryNamespace = 'project', metadata: MemoryMetadata = {}): Promise<Memory> {
        return await this.add(content, 'long_term', namespace, metadata);
    }

    /**
     * Check if memory should be consolidated to long-term
     */
    private async checkConsolidation(memory: Memory): Promise<void> {
        if (memory.type === 'working' &&
            memory.accessCount >= this.options.consolidationThreshold) {
            // Promote to long-term memory
            memory.type = 'long_term';
            memory.ttl = undefined;
            this.dirty = true;

            // Sync update to vector DB
            await this.memoryManager.saveContext(memory.content, {
                ...memory.metadata,
                id: memory.id,
                type: 'long_term',
                consolidatedAt: Date.now()
            });

            console.log(`[AgentMemoryService] Consolidated memory ${memory.id} to long-term`);
        }
    }

    /**
     * Enforce memory limits by removing old memories
     */
    private enforceMemoryLimits(): void {
        const session = this.getByType('session');
        const working = this.getByType('working');

        // Remove excess session memories (oldest first)
        if (session.length > this.options.maxSessionMemories) {
            const toRemove = session
                .sort((a, b) => a.createdAt.getTime() - b.createdAt.getTime())
                .slice(0, session.length - this.options.maxSessionMemories);
            for (const m of toRemove) {
                this.delete(m.id);
            }
        }

        // Remove excess working memories (least accessed first)
        if (working.length > this.options.maxWorkingMemories) {
            const toRemove = working
                .sort((a, b) => a.accessCount - b.accessCount)
                .slice(0, working.length - this.options.maxWorkingMemories);
            for (const m of toRemove) {
                this.delete(m.id);
            }
        }
    }

    /**
     * Get memories by type
     */
    getByType(type: MemoryType): Memory[] {
        return Array.from(this.memories.values()).filter(m => m.type === type);
    }

    /**
     * Get memories by namespace
     */
    getByNamespace(namespace: MemoryNamespace): Memory[] {
        return Array.from(this.memories.values()).filter(m => m.namespace === namespace);
    }

    /**
     * Clear all session memories
     */
    clearSession(): void {
        for (const [id, memory] of this.memories) {
            if (memory.type === 'session') {
                this.memories.delete(id);
            }
        }
        this.dirty = true;
    }

    /**
     * Get memory statistics
     */
    getStats(): Record<string, number> {
        const byType: Record<string, number> = { session: 0, working: 0, long_term: 0 };
        const byNamespace: Record<string, number> = { user: 0, agent: 0, project: 0 };

        for (const memory of this.memories.values()) {
            byType[memory.type]++;
            byNamespace[memory.namespace]++;
        }

        // We could also mix in stats from memoryManager if needed

        return {
            sessionCount: this.memories.size,
            ...byType,
            ...byNamespace,
        };
    }

    /**
     * Export memories to JSON
     */
    export(): string {
        return JSON.stringify({
            exportedAt: new Date().toISOString(),
            memories: Array.from(this.memories.values()),
        }, null, 2);
    }

    /**
     * Import memories from JSON
     */
    async import(jsonData: string): Promise<number> {
        const data = JSON.parse(jsonData);
        let count = 0;

        for (const mem of data.memories || []) {
            mem.createdAt = new Date(mem.createdAt);
            mem.accessedAt = new Date(mem.accessedAt);
            mem.id = this.generateId();  // Assign new ID

            this.memories.set(mem.id, mem);

            if (mem.type !== 'session') {
                await this.memoryManager.saveContext(mem.content, {
                    ...mem.metadata,
                    id: mem.id,
                    type: mem.type
                });
            }
            count++;
        }

        this.dirty = true;
        return count;
    }

    /**
     * Schedule auto-save
     */
    private saveTimeout: NodeJS.Timeout | null = null;
    private scheduleSave(): void {
        if (this.saveTimeout) return;
        this.saveTimeout = setTimeout(() => {
            this.saveToDisk();
            this.saveTimeout = null;
        }, 5000);  // Batch saves every 5 seconds
    }

    /**
     * Force save to disk
     */
    flush(): void {
        if (this.saveTimeout) {
            clearTimeout(this.saveTimeout);
            this.saveTimeout = null;
        }
        this.saveToDisk();
    }

    /**
     * Shutdown and save
     */
    shutdown(): void {
        this.flush();
    }
}

export default AgentMemoryService;
