import fs from 'fs';
import path from 'path';
// Remove yaml dependency to avoid overhead if simple parsing works, or use a regex parser for MVP.
// We'll use a simple JSON fallback if YAML fails, or check if we can add 'js-yaml' to core.
// For now, let's use a simple JSON format policy.json for robustness, or parse simple YAML manually.
// Let's rely on standard JSON for policy definition first to avoid deps issues.

export interface PolicyRule {
    action: string; // "read", "write", "execute", "tool_use"
    resource: string; // Glob pattern or specific tool name
    effect: "ALLOW" | "DENY" | "ASK";
    reason?: string;
}

export interface PolicyConfig {
    version: string;
    rules: PolicyRule[];
}

export class PolicyService {
    private policyPath: string;
    private config: PolicyConfig;

    constructor(cwd: string) {
        this.policyPath = path.join(cwd, '.borg', 'policy.json');
        this.config = { version: "1.0", rules: [] };
        this.loadPolicy();
    }

    private loadPolicy() {
        try {
            if (fs.existsSync(this.policyPath)) {
                const content = fs.readFileSync(this.policyPath, 'utf-8');
                this.config = JSON.parse(content);
                console.log(`[PolicyService] Loaded ${this.config.rules.length} rules.`);
            } else {
                // Create default policy
                this.config = {
                    version: "1.0",
                    rules: [
                        { action: "write", resource: ".borg/policy.json", effect: "DENY", reason: "Prevent self-modification of policy" },
                        { action: "execute", resource: "rm -rf /", effect: "DENY", reason: "Prevent catastrophe" },
                        // Default Allow All for now to not break existing flow until configured
                        { action: "*", resource: "*", effect: "ALLOW" }
                    ]
                };
                this.savePolicy();
            }
        } catch (e) {
            console.error("[PolicyService] Failed to load policy:", e);
        }
    }

    private savePolicy() {
        try {
            const dir = path.dirname(this.policyPath);
            if (!fs.existsSync(dir)) fs.mkdirSync(dir, { recursive: true });
            fs.writeFileSync(this.policyPath, JSON.stringify(this.config, null, 2));
        } catch (e) {
            console.error("[PolicyService] Failed to save policy:", e);
        }
    }

    public check(action: string, resource: string): "ALLOW" | "DENY" | "ASK" {
        // Iterate rules top to bottom. First match wins.
        for (const rule of this.config.rules) {
            if (this.matches(rule.action, action) && this.matches(rule.resource, resource)) {
                return rule.effect;
            }
        }
        return "ALLOW"; // Default allow if no rule matches (or deny if paranoid mode?)
    }

    private matches(pattern: string, value: string): boolean {
        if (pattern === "*") return true;
        if (pattern === value) return true;

        // Simple wildcard suffix matching (e.g. "tool_*")
        if (pattern.endsWith('*')) {
            const prefix = pattern.slice(0, -1);
            return value.startsWith(prefix);
        }

        return false;
    }

    public updateRules(newRules: PolicyRule[]) {
        this.config.rules = newRules;
        this.savePolicy();
    }

    public getRules() {
        return this.config.rules;
    }
}
