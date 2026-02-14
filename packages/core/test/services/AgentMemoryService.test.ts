import { describe, test, expect, beforeEach, afterEach, mock } from 'bun:test';
import { AgentMemoryService } from '../../src/services/AgentMemoryService.ts';
import { MemoryManager } from '../../src/services/MemoryManager.ts';
import * as fs from 'fs';
import * as path from 'path';
import * as os from 'os';

interface SessionMediaMetadata {
    media: Array<{
        type: 'image' | 'audio' | 'video' | 'file';
        url?: string;
        description?: string;
    }>;
}

// Mock MemoryManager
class MockMemoryManager {
    public storage: Map<string, any> = new Map();

    async saveContext(content: string, metadata: any = {}) {
        const id = metadata.id || `mock-${Date.now()}`;
        this.storage.set(id, { content, metadata });
        return id;
    }

    async deleteContext(id: string) {
        this.storage.delete(id);
    }

    async search(query: string, limit: number = 5) {
        // Simple mock search: filter by content inclusion
        const results = [];
        for (const [id, doc] of this.storage.entries()) {
            if (doc.content.includes(query)) {
                results.push({
                    id,
                    content: doc.content,
                    metadata: doc.metadata,
                    score: 0.9
                });
            }
        }
        return results.slice(0, limit);
    }

    // Stub other methods if needed
    async initialize() { }
    async getContext(id: string) { return this.storage.get(id); }
}

describe('AgentMemoryService', () => {
    let memoryService: AgentMemoryService;
    let mockMemoryManager: MockMemoryManager;
    let tempDir: string;

    beforeEach(() => {
        tempDir = path.join(os.tmpdir(), `borg-test-${Date.now()}`);
        if (!fs.existsSync(tempDir)) fs.mkdirSync(tempDir, { recursive: true });

        mockMemoryManager = new MockMemoryManager();

        memoryService = new AgentMemoryService({
            persistDir: path.join(tempDir, 'memories'),
            consolidationThreshold: 2 // Low for testing
        }, mockMemoryManager as unknown as MemoryManager);
    });

    afterEach(() => {
        if (fs.existsSync(tempDir)) {
            fs.rmSync(tempDir, { recursive: true, force: true });
        }
    });

    test('should add and retrieve session memory', async () => {
        const content = 'This is a session memory';
        const mem = await memoryService.addSession(content);

        expect(mem.type).toBe('session');
        expect(mem.content).toBe(content);

        const retrieved = memoryService.get(mem.id);
        expect(retrieved).not.toBeNull();
        expect(retrieved?.id).toBe(mem.id);
    });

    test('should add working memory and consolidate to long-term', async () => {
        const content = 'Important project fact';
        const mem = await memoryService.addWorking(content);

        expect(mem.type).toBe('working');

        // Check if persisted to mock manager
        expect(mockMemoryManager.storage.has(mem.id)).toBe(true);

        // Access memory to trigger consolidation (threshold is 2)
        // Access 1
        memoryService.get(mem.id);
        // Access 2 -> triggers consolidation
        await memoryService.get(mem.id);

        const consolidated = memoryService.get(mem.id);
        expect(consolidated?.type).toBe('long_term');

        // Verify mock manager updated
        const mockDoc = mockMemoryManager.storage.get(mem.id);
        expect(mockDoc.metadata.type).toBe('long_term');
    });

    test('should support multi-modal metadata', async () => {
        const mem = await memoryService.addSession('Look at this image', {
            media: [{
                type: 'image',
                url: 'http://example.com/img.png',
                description: 'A test image'
            }]
        } as unknown as SessionMediaMetadata);

        const retrieved = memoryService.get(mem.id);
        expect(retrieved?.metadata.media).toBeDefined();
        // @ts-ignore
        expect(retrieved?.metadata.media[0].type).toBe('image');
    });

    test('should search across tiers (hybrid search)', async () => {
        // Add session memory
        await memoryService.addSession('Current context about pizza');

        // Add long-term memory (persists to mock manager)
        await memoryService.addLongTerm('Historic context about pasta');

        // verify Mock Manager acts like vector DB
        const searchResultsMock = await mockMemoryManager.search('pasta');
        expect(searchResultsMock.length).toBe(1);

        // Search via Service (Hybrid)
        const results = await memoryService.search('context');

        // Should find both (pizza from session, pasta from mock manager)
        expect(results.length).toBeGreaterThanOrEqual(1);

        // Verify we found the 'pasta' (from Long Term)
        const foundPasta = results.find(r => r.content.includes('pasta'));
        expect(foundPasta).toBeDefined();

        // Verify we found 'pizza' (from Session)
        const foundPizza = results.find(r => r.content.includes('pizza'));
        expect(foundPizza).toBeDefined();
    });
});
