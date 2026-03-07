import fs from 'fs';
import path from 'path';

import type { MCPServerConfig } from './types.js';

export class MCPConfigStore {
    constructor(private readonly configPath: string) {}

    public readAll(): Record<string, MCPServerConfig> {
        if (!fs.existsSync(this.configPath)) {
            return {};
        }

        const raw = fs.readFileSync(this.configPath, 'utf8');
        return JSON.parse(raw) as Record<string, MCPServerConfig>;
    }

    public writeAll(config: Record<string, MCPServerConfig>): void {
        const dir = path.dirname(this.configPath);
        if (!fs.existsSync(dir)) {
            fs.mkdirSync(dir, { recursive: true });
        }

        fs.writeFileSync(this.configPath, JSON.stringify(config, null, 2));
    }

    public upsert(name: string, config: MCPServerConfig): void {
        const current = this.readAll();
        current[name] = config;
        this.writeAll(current);
    }

    public remove(name: string): void {
        const current = this.readAll();
        delete current[name];
        this.writeAll(current);
    }
}
