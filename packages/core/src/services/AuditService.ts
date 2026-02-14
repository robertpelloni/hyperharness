
import fs from 'fs';
import path from 'path';

export interface AuditLogEntry {
    timestamp: number;
    agentId?: string;
    action: string;
    params: unknown;
    result?: unknown;
    status: 'SUCCESS' | 'FAILURE' | 'DENIED';
    durationMs?: number;
}

export class AuditService {
    private logPath: string;
    private buffer: AuditLogEntry[] = [];
    private flushInterval: NodeJS.Timeout;

    constructor(logDir: string = '.borg/audit') {
        if (!fs.existsSync(logDir)) {
            fs.mkdirSync(logDir, { recursive: true });
        }
        // Rotate logs by day
        const date = new Date().toISOString().split('T')[0];
        this.logPath = path.join(logDir, `audit-${date}.jsonl`);

        // Auto-flush every 5 seconds
        this.flushInterval = setInterval(() => this.flush(), 5000);
    }

    public log(entry: Omit<AuditLogEntry, 'timestamp'>) {
        const fullEntry: AuditLogEntry = {
            timestamp: Date.now(),
            ...entry
        };
        this.buffer.push(fullEntry);
    }

    public async flush() {
        if (this.buffer.length === 0) return;

        const chunk = this.buffer.splice(0, this.buffer.length);
        const lines = chunk.map(e => JSON.stringify(e)).join('\n') + '\n';

        try {
            await fs.promises.appendFile(this.logPath, lines, 'utf8');
        } catch (e) {
            console.error("[AuditService] Failed to flush logs", e);
            // Put back in buffer? Or drop to avoid memory leak? Drop for now.
        }
    }

    public async queryLogs(filter: { action?: string, limit?: number }): Promise<AuditLogEntry[]> {
        // Simple file scan (inefficient for huge logs, acceptable for V1)
        if (!fs.existsSync(this.logPath)) return [];

        const content = await fs.promises.readFile(this.logPath, 'utf8');
        const lines = content.trim().split('\n');
        const logs = lines.map(line => {
            try { return JSON.parse(line); } catch { return null; }
        }).filter(l => l !== null) as AuditLogEntry[];

        return logs
            .filter(l => !filter.action || l.action === filter.action)
            .slice(-(filter.limit || 100));
    }
}
