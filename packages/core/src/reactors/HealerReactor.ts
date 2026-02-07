
import { EventBus, SystemEvent } from '../services/EventBus.js';
import { HealerService } from '../services/HealerService.js';

export class HealerReactor {
    private eventBus: EventBus;
    private healerService: HealerService;
    private isHealing: boolean = false;
    private lastErrorTime: number = 0;
    private readonly COOLDOWN_MS = 10000; // 10s cooldown to prevent loops

    constructor(eventBus: EventBus, healerService: HealerService) {
        this.eventBus = eventBus;
        this.healerService = healerService;
    }

    public start() {
        console.log("[HealerReactor] 🛡️ Immune System Active. Listening for pathogens...");
        // @ts-ignore
        this.eventBus.subscribe('terminal:error', this.handleError.bind(this));
    }

    private async handleError(event: SystemEvent) {
        const errorLog = event.payload.message || event.payload.error;

        // Prevent reaction loops (e.g., if the healer itself causes an error)
        const now = Date.now();
        if (this.isHealing || (now - this.lastErrorTime < this.COOLDOWN_MS)) {
            console.log("[HealerReactor] ⏳ Healing in progress or cooldown. Skipping error.");
            return;
        }

        console.log(`[HealerReactor] 🩺 Detected Pathogen! Initiating immune response...`);
        this.isHealing = true;
        this.lastErrorTime = now;

        try {
            this.eventBus.emitEvent('task:update', 'HealerReactor', { message: 'Diagnosing error...' });

            // Trigger the Healer Service
            const report = await this.healerService.autoHeal(errorLog);

            if (report.success) {
                console.log(`[HealerReactor] ✅ Pathogen neutralized. Fix applied to ${report.file}`);
                this.eventBus.emitEvent('system:healed', 'HealerReactor', { file: report.file, fix: report.fix });
            } else {
                console.log(`[HealerReactor] ⚠️ Integration failed. Could not auto-heal.`);
            }

        } catch (error) {
            console.error("[HealerReactor] ❌ Immune System Failure:", error);
        } finally {
            this.isHealing = false;
        }
    }
}
