
import { EventBus } from '../services/EventBus.js';

export class TerminalSensor {
    private eventBus: EventBus;
    private originalStderrWrite: any;

    constructor(eventBus: EventBus) {
        this.eventBus = eventBus;
    }

    public start() {
        console.log("[TerminalSensor] Hooking stderr...");

        // Hook stderr to capture errors
        this.originalStderrWrite = process.stderr.write;

        // @ts-ignore
        process.stderr.write = (chunk: any, encoding?: any, callback?: any) => {
            const str = chunk.toString();

            // Heuristic to detect ACTUAL errors vs warnings
            if (str.toLowerCase().includes('error') || str.toLowerCase().includes('exception') || str.includes('❌')) {
                this.eventBus.emitEvent('terminal:error', 'TerminalSensor', { message: str });
            }

            // Pass through to original stderr
            return this.originalStderrWrite.call(process.stderr, chunk, encoding, callback);
        };
    }

    public stop() {
        if (this.originalStderrWrite) {
            process.stderr.write = this.originalStderrWrite;
        }
    }
}
