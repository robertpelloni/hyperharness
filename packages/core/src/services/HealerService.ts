
import { LLMService } from "@borg/ai";
import { MCPServer } from "../MCPServer.js"; // Circular dep? Ideally pass deps.

export class HealerService {
    private llmService: LLMService;
    private mcpServer: any; // Use any to avoid circular type issues for now

    constructor(llmService: LLMService, mcpServer: any) {
        this.llmService = llmService;
        this.mcpServer = mcpServer;
    }

    public async heal(errorContext: { url: string, error: string, timestamp: string }) {
        console.log(`[HealerService] 🚑 Detecting Error on ${errorContext.url}`);
        console.log(`[HealerService] Error: ${errorContext.error}`);

        // 1. Capture Context (Screenshot)
        let screenshotBase64 = null;
        try {
            if (this.mcpServer.captureScreenshotFromBrowser) {
                const dataUrl = await this.mcpServer.captureScreenshotFromBrowser();
                screenshotBase64 = dataUrl.split(',')[1];
            }
        } catch (e) {
            console.warn("[HealerService] Could not capture screenshot for analysis:", e);
        }

        // 2. Identify and Read Source File from Stack Trace
        let sourceContext = "";
        let targetFilePath = "";
        try {
            // Regex to find file paths in stack trace (Windows/Unix)
            // match strings like: (C:\path\to\file.ts:123:45) or (/path/to/file.ts:123:45)
            const stackMatch = errorContext.error.match(/[\(\s](([a-zA-Z]:\\|\/)[^:)]+\.(ts|js|jsx|tsx)):(\d+)/);

            if (stackMatch && stackMatch[1]) {
                targetFilePath = stackMatch[1];
                const lineNo = parseInt(stackMatch[4]);

                // Read the file
                const fs = await import('fs/promises');
                const content = await fs.readFile(targetFilePath, 'utf-8');
                const lines = content.split('\n');

                // Extract context (+/- 20 lines)
                const start = Math.max(0, lineNo - 20);
                const end = Math.min(lines.length, lineNo + 20);
                const snippet = lines.slice(start, end).map((l, i) => `${start + i + 1}: ${l}`).join('\n');

                sourceContext = `\nSource File: ${targetFilePath}\nContext (around line ${lineNo}):\n\`\`\`typescript\n${snippet}\n\`\`\`\n`;
                console.log(`[HealerService] Identified source file: ${targetFilePath} (Line ${lineNo})`);
            }
        } catch (e) {
            console.warn("[HealerService] Could not extract source context:", e);
        }

        // 3. Diagnose with LLM
        const systemPrompt = `You are an Autonomous Site Reliability Engineer (Healer).
Your job is to diagnose web application errors and propose fixes.
You have access to the error log, a screenshot (optional), and source code context (optional).

Return a JSON object ONLY with the following structure (no markdown backticks):
{
    "analysis": "Explanation of what went wrong",
    "severity": "low|medium|high|critical",
    "fixProposal": {
        "filePath": "${targetFilePath || "path/to/file.ts"}",
        "description": "What this fix does",
        "replacements": [
            {
                "search": "exact string to replace (must match source context exactly)",
                "replace": "new string content"
            }
        ]
    },
    "confidence": 0.0 to 1.0
}`;

        const userPrompt = `Error: ${errorContext.error}
URL: ${errorContext.url}
Timestamp: ${errorContext.timestamp}
${sourceContext}

Please analyze this crash and provide a JSON fix.`;

        try {
            const options: any = { taskComplexity: 'high' };
            if (screenshotBase64) {
                options.images = [{ base64: screenshotBase64, mimeType: 'image/jpeg' }];
            }

            const response = await this.llmService.generateText(
                'google',
                'gemini-1.5-pro',
                systemPrompt,
                userPrompt,
                options
            );

            console.log("\n[HealerService] 🩺 DIAGNOSIS:");
            // Clean markdown blocks if present
            const cleanJson = response.content.replace(/```json\n?|\n?```/g, "").trim();

            let diagnosis;
            try {
                diagnosis = JSON.parse(cleanJson);
                console.log(JSON.stringify(diagnosis, null, 2));

                // 4. Autonomous Healing Logic with Council Approval
                if (diagnosis.fixProposal && diagnosis.fixProposal.replacements.length > 0) {
                    console.log(`[HealerService] 🏛️ Proposing fix to Council (Confidence: ${diagnosis.confidence})...`);

                    const topic = `Healer Fix Proposal for ${diagnosis.fixProposal.filePath}: ${diagnosis.fixProposal.description}\n\nAnalysis: ${diagnosis.analysis}`;

                    // Always consult if confidence is less than absolute certainty, or just always for safety
                    const debate = await this.mcpServer.council.runConsensusSession(topic);

                    if (debate.approved) {
                        console.log("[HealerService] ✅ Council APPROVED fix. Applying...");
                        await this.applyFix(diagnosis.fixProposal);
                    } else {
                        console.log(`[HealerService] ❌ Council REJECTED fix: ${debate.summary}`);
                        await this.logEvent({
                            type: 'FIX_REJECTED',
                            file: diagnosis.fixProposal.filePath,
                            timestamp: new Date().toISOString(),
                            reason: debate.summary
                        });
                    }
                } else if (diagnosis.fixProposal) {
                    console.log("[HealerService] No actionable replacements found in proposal.");
                }

            } catch (jsonErr) {
                console.error("[HealerService] Failed to parse diagnosis JSON:", jsonErr);
                console.log("Raw Response:", response.content);
            }

        } catch (e) {
            console.error("[HealerService] Diagnosis Failed:", e);
        }
    }

    private async applyFix(proposal: { filePath: string, replacements: { search: string, replace: string }[] }) {
        console.log(`[HealerService] 🩹 Applying Fix to ${proposal.filePath}...`);
        try {
            const fs = await import('fs/promises');
            const path = await import('path');
            let content = await fs.readFile(proposal.filePath, 'utf-8');
            let modified = false;

            for (const item of proposal.replacements) {
                if (content.includes(item.search)) {
                    content = content.replace(item.search, item.replace);
                    modified = true;
                } else {
                    console.warn(`[HealerService] Fix failed: Could not find search string in ${proposal.filePath}`);
                    console.warn(`Search: "${item.search}"`);
                }
            }

            if (modified) {
                // Backup
                await fs.writeFile(`${proposal.filePath}.bak`, await fs.readFile(proposal.filePath, 'utf-8'));
                // Write
                await fs.writeFile(proposal.filePath, content);
                console.log(`[HealerService] ✅ Fix Applied! Backup created at ${proposal.filePath}.bak`);

                // Log Event for Dashboard
                await this.logEvent({
                    type: 'FIX_APPLIED',
                    file: proposal.filePath,
                    timestamp: new Date().toISOString(),
                    details: proposal.replacements.length + " replacements"
                });

                // Optional: Notify through WebSocket to refresh
                if (this.mcpServer.broadcastRequestAndAwait) {
                    // this.mcpServer.broadcastRequestAndAwait('TOAST_NOTIFICATION', { message: "Healer applied a fix!", type: "success" });
                }
            }
        } catch (e) {
            console.error(`[HealerService] Failed to apply fix:`, e);
            await this.logEvent({
                type: 'FIX_FAILED',
                file: proposal.filePath,
                timestamp: new Date().toISOString(),
                error: (e as any).message
            });
        }
    }

    public async getHistory(): Promise<any[]> {
        try {
            const fs = await import('fs/promises');
            const path = await import('path');
            const logFile = path.join(process.cwd(), '.borg', 'data', 'healer_events.jsonl');
            try {
                const content = await fs.readFile(logFile, 'utf-8');
                return content.trim().split('\n').map(line => {
                    try { return JSON.parse(line); } catch (e) { return null; }
                }).filter(Boolean).reverse(); // Newest first
            } catch (e) {
                // File likely doesn't exist
                return [];
            }
        } catch (e) {
            console.error("[HealerService] Failed to read history:", e);
            return [];
        }
    }

    private async logEvent(event: any) {
        try {
            const fs = await import('fs/promises');
            const path = await import('path');
            const logDir = path.join(process.cwd(), '.borg', 'data');
            await fs.mkdir(logDir, { recursive: true });
            const logFile = path.join(logDir, 'healer_events.jsonl');
            await fs.appendFile(logFile, JSON.stringify(event) + '\n');
        } catch (e) {
            console.error("[HealerService] Failed to log event:", e);
        }
    }
}
