/**
 * AutoDevService - "Fix until Pass" loops for tests and linters
 * Automatically retries fixing code until tests/lints pass or max attempts reached
 */

import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

export interface LoopConfig {
    maxAttempts: number;
    type: 'test' | 'lint' | 'build';
    target?: string; // Specific file or pattern
    command?: string; // Custom command override
}

export interface LoopResult {
    success: boolean;
    attempts: number;
    output: string;
    errors: string[];
    duration: number;
}

export interface ActiveLoop {
    id: string;
    config: LoopConfig;
    status: 'running' | 'success' | 'failed' | 'cancelled';
    currentAttempt: number;
    startTime: number;
    lastOutput: string;
}

export class AutoDevService {
    private activeLoops: Map<string, ActiveLoop> = new Map();
    private loopCounter = 0;

    constructor(private rootDir: string) { }

    /**
     * Start a "Fix until Pass" loop
     */
    async startLoop(config: LoopConfig): Promise<string> {
        const id = `loop-${++this.loopCounter}`;
        const loop: ActiveLoop = {
            id,
            config,
            status: 'running',
            currentAttempt: 0,
            startTime: Date.now(),
            lastOutput: ''
        };

        this.activeLoops.set(id, loop);
        console.log(`[AutoDev] 🔄 Starting ${config.type} loop (max ${config.maxAttempts} attempts)`);

        // Run the loop asynchronously
        this.runLoop(id).catch(e => {
            console.error(`[AutoDev] Loop ${id} error:`, e);
            const l = this.activeLoops.get(id);
            if (l) l.status = 'failed';
        });

        return id;
    }

    /**
     * Cancel an active loop
     */
    cancelLoop(id: string): boolean {
        const loop = this.activeLoops.get(id);
        if (loop && loop.status === 'running') {
            loop.status = 'cancelled';
            console.log(`[AutoDev] 🛑 Cancelled loop ${id}`);
            return true;
        }
        return false;
    }

    /**
     * Get status of all loops
     */
    getLoops(): ActiveLoop[] {
        return Array.from(this.activeLoops.values());
    }

    /**
     * Get a specific loop
     */
    getLoop(id: string): ActiveLoop | undefined {
        return this.activeLoops.get(id);
    }

    private async runLoop(id: string): Promise<void> {
        const loop = this.activeLoops.get(id);
        if (!loop) return;

        const { config } = loop;
        const command = this.getCommand(config);

        while (loop.currentAttempt < config.maxAttempts && loop.status === 'running') {
            loop.currentAttempt++;
            console.log(`[AutoDev] Attempt ${loop.currentAttempt}/${config.maxAttempts}`);

            try {
                const { stdout, stderr } = await execAsync(command, {
                    cwd: this.rootDir,
                    timeout: 120000 // 2 minute timeout per attempt
                });

                loop.lastOutput = stdout || stderr;

                // Success!
                loop.status = 'success';
                console.log(`[AutoDev] ✅ ${config.type} passed on attempt ${loop.currentAttempt}`);
                return;

            } catch (error: any) {
                loop.lastOutput = error.stdout || error.stderr || error.message;

                if (loop.status !== 'running') {
                    // Cancelled
                    return;
                }

                console.log(`[AutoDev] ❌ Attempt ${loop.currentAttempt} failed`);

                // Don't retry on last attempt
                if (loop.currentAttempt >= config.maxAttempts) {
                    loop.status = 'failed';
                    console.log(`[AutoDev] 💀 Max attempts reached. Loop failed.`);
                    return;
                }

                // Wait before retry (exponential backoff)
                const delay = Math.min(1000 * Math.pow(2, loop.currentAttempt - 1), 30000);
                await new Promise(r => setTimeout(r, delay));
            }
        }
    }

    private getCommand(config: LoopConfig): string {
        if (config.command) return config.command;

        switch (config.type) {
            case 'test':
                return config.target
                    ? `npx vitest run ${config.target}`
                    : 'npm test';
            case 'lint':
                return config.target
                    ? `npx eslint --fix ${config.target}`
                    : 'npm run lint -- --fix';
            case 'build':
                return 'npm run build';
            default:
                return 'npm test';
        }
    }

    /**
     * Clear completed loops
     */
    clearCompleted(): number {
        let count = 0;
        for (const [id, loop] of this.activeLoops) {
            if (loop.status !== 'running') {
                this.activeLoops.delete(id);
                count++;
            }
        }
        return count;
    }
}
