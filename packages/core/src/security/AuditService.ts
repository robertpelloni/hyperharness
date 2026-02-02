
import fs from 'fs';
import path from 'path';

export interface AuditEntry {
    timestamp: string;
    level: 'INFO' | 'WARN' | 'ERROR';
    agentId?: string;
    event: string;
    details?: any;
    metadata?: any;
}

export class AuditService {
    private logPath: string;

    constructor(cwd: string) {
        const auditDir = path.join(cwd, '.borg', 'audit');
        if (!fs.existsSync(auditDir)) {
            // Ensure .borg exists first
            if (!fs.existsSync(path.join(cwd, '.borg'))) {
                fs.mkdirSync(path.join(cwd, '.borg'));
            }
            fs.mkdirSync(auditDir);
        }
        // Rotate inputs? For MVP, single file.
        this.logPath = path.join(auditDir, 'audit.jsonl');
    }

    log(event: string, details?: any, level: 'INFO' | 'WARN' | 'ERROR' = 'INFO') {
        const entry: AuditEntry = {
            timestamp: new Date().toISOString(),
            level,
            event,
            details,
        };

        const line = JSON.stringify(entry) + '\n';
        try {
            fs.appendFileSync(this.logPath, line, 'utf-8');
        } catch (e) {
            console.error("[AuditService] Failed to write log:", e);
        }
    }

    getLogs(limit: number = 50): AuditEntry[] {
        return this.query({ limit });
    }

    query(filters: { level?: string; agentId?: string; limit?: number; event?: string }): AuditEntry[] {
        if (!fs.existsSync(this.logPath)) return [];

        try {
            const content = fs.readFileSync(this.logPath, 'utf-8');
            const lines = content.trim().split('\n');
            const entries: AuditEntry[] = [];

            // Read backwards for latest first
            for (let i = lines.length - 1; i >= 0; i--) {
                const line = lines[i].trim();
                if (!line) continue;

                try {
                    const entry = JSON.parse(line);

                    // Filters
                    if (filters.level && entry.level !== filters.level) continue;
                    if (filters.agentId && entry.agentId !== filters.agentId) continue;
                    if (filters.event && !entry.event.includes(filters.event)) continue;

                    entries.push(entry);

                    if (filters.limit && entries.length >= filters.limit) break;
                } catch (e) {
                    // Ignore corrupted lines
                }
            }
            return entries;
        } catch (e) {
            console.error("[AuditService] Failed to read logs:", e);
            return [];
        }
    }
}
