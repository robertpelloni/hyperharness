
import { EventBus } from '../services/EventBus.js';

export class TerminalSensor {
    private eventBus: EventBus;
    private originalStderrWrite: typeof process.stderr.write | null = null;

    constructor(eventBus: EventBus) {
        this.eventBus = eventBus;
    }

    public start() {
        console.log("[TerminalSensor] Hooking stderr...");

        // Hook stderr to capture errors
        this.originalStderrWrite = process.stderr.write.bind(process.stderr);

        /**
         * Reason: we need to observe terminal error output in real-time for event emission.
         * What: wraps stderr writes, emits `terminal:error` on heuristic matches, then forwards bytes unchanged.
         * Why: preserves normal stderr behavior while enabling lightweight telemetry without unsafe casts.
         */
        const patchedWrite: typeof process.stderr.write = (chunk, encoding?, callback?) => {
            const str = chunk.toString();

            // Heuristic to detect ACTUAL errors vs warnings
            if (str.toLowerCase().includes('error') || str.toLowerCase().includes('exception') || str.includes('❌')) {
                this.eventBus.emitEvent('terminal:error', 'TerminalSensor', { message: str });
            }

            // Pass through to original stderr
            if (!this.originalStderrWrite) {
                return false;
            }
            return this.originalStderrWrite(chunk, encoding as never, callback as never);
        };

        process.stderr.write = patchedWrite;
    }

    public stop() {
        if (this.originalStderrWrite) {
            process.stderr.write = this.originalStderrWrite;
            this.originalStderrWrite = null;
        }
    }
}
