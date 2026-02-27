/**
 * MemoryExportImportService – Export and import memories in multiple formats
 *
 * Supports:
 * - JSON: Full fidelity export/import, preserves all metadata.
 * - CSV: Spreadsheet-friendly; columns = uuid, content, userId, agentId, createdAt, metadata (JSON string).
 * - JSONL: Streaming-friendly; one JSON object per line.
 *
 * Accepts either IMemoryProvider (listMemories/saveMemory) or
 * MemoryManagerRuntime (listContexts/saveContext) as the backend.
 */

import { Memory } from '../../interfaces/IMemoryProvider.js';

export type ExportFormat = 'json' | 'csv' | 'jsonl';

/**
 * Flexible backend interface that works with both IMemoryProvider and
 * MemoryManagerRuntime so the service can be used in either context.
 */
interface MemoryBackend {
    listMemories?: (userId: string, limit?: number, offset?: number) => Promise<Memory[]>;
    listContexts?: () => Promise<unknown[]>;
    saveMemory?: (content: string, metadata: Record<string, unknown>, userId: string, agentId?: string) => Promise<Memory>;
    saveContext?: (content: string, metadata?: Record<string, unknown>) => Promise<unknown>;
}

export class MemoryExportImportService {
    private backend: MemoryBackend;

    constructor(backend: MemoryBackend) {
        this.backend = backend;
    }

    // ---- Export ----

    async exportAll(userId: string, format: ExportFormat): Promise<string> {
        // Fetch all memories (large limit to capture everything)
        let memories: any[] = [];
        if (this.backend.listMemories) {
            memories = await this.backend.listMemories(userId, 100000, 0);
        } else if (this.backend.listContexts) {
            memories = await this.backend.listContexts();
        }

        switch (format) {
            case 'json':
                return JSON.stringify(memories, null, 2);
            case 'csv':
                return this.toCSV(memories);
            case 'jsonl':
                return memories.map(m => JSON.stringify(m)).join('\n');
            default:
                throw new Error(`Unsupported export format: ${format}`);
        }
    }

    // ---- Import ----

    async importBulk(data: string, format: ExportFormat, userId: string): Promise<{ imported: number; errors: number }> {
        let records: any[];

        switch (format) {
            case 'json':
                records = JSON.parse(data);
                break;
            case 'csv':
                records = this.fromCSV(data);
                break;
            case 'jsonl':
                records = data
                    .split('\n')
                    .filter(line => line.trim())
                    .map(line => JSON.parse(line));
                break;
            default:
                throw new Error(`Unsupported import format: ${format}`);
        }

        let imported = 0;
        let errors = 0;

        for (const record of records) {
            try {
                if (this.backend.saveMemory) {
                    await this.backend.saveMemory(
                        record.content || '',
                        record.metadata || {},
                        userId,
                        record.agentId
                    );
                } else if (this.backend.saveContext) {
                    await this.backend.saveContext(
                        record.content || '',
                        record.metadata
                    );
                }
                imported++;
            } catch {
                errors++;
            }
        }

        return { imported, errors };
    }

    // ---- Helpers ----

    private toCSV(memories: any[]): string {
        const header = 'uuid,content,userId,agentId,createdAt,metadata';
        const rows = memories.map(m => {
            const escapedContent = `"${(m.content || '').replace(/"/g, '""')}"`;
            const escapedMeta = `"${JSON.stringify(m.metadata || {}).replace(/"/g, '""')}"`;
            return [
                m.uuid || m.id || '',
                escapedContent,
                m.userId || '',
                m.agentId || '',
                m.createdAt instanceof Date ? m.createdAt.toISOString() : String(m.createdAt || new Date()),
                escapedMeta,
            ].join(',');
        });
        return [header, ...rows].join('\n');
    }

    private fromCSV(csv: string): any[] {
        const lines = csv.split('\n').filter(l => l.trim());
        if (lines.length < 2) return []; // header only or empty

        // Skip header
        const dataLines = lines.slice(1);
        return dataLines.map(line => {
            // Simple CSV parsing (handles quoted fields with embedded commas)
            const fields = this.parseCSVLine(line);
            return {
                uuid: fields[0],
                content: fields[1],
                userId: fields[2],
                agentId: fields[3] || undefined,
                createdAt: fields[4] ? new Date(fields[4]) : new Date(),
                metadata: fields[5] ? (() => { try { return JSON.parse(fields[5]); } catch { return {}; } })() : {},
            };
        });
    }

    private parseCSVLine(line: string): string[] {
        const result: string[] = [];
        let current = '';
        let inQuotes = false;

        for (let i = 0; i < line.length; i++) {
            const char = line[i];
            if (char === '"') {
                if (inQuotes && line[i + 1] === '"') {
                    current += '"';
                    i++; // skip escaped quote
                } else {
                    inQuotes = !inQuotes;
                }
            } else if (char === ',' && !inQuotes) {
                result.push(current);
                current = '';
            } else {
                current += char;
            }
        }
        result.push(current);
        return result;
    }
}
