
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
            // We need to access captureScreenshotFromBrowser from MCPServer
            // Since it's private, we might need a public accessor or just use the tool execution flow?
            // Better to make a public method on MCPServer or pass a capture callback.
            // For now, assume we can call the tool handler or similar.
            // Or better: MCPServer passes a screenshot function.

            // Let's rely on MCPServer instance having public capture method or we refactor MCPServer to expose it.
            // I'll assume I'll make captureScreenshotFromBrowser public or internal access.
            if (this.mcpServer.captureScreenshotFromBrowser) {
                const dataUrl = await this.mcpServer.captureScreenshotFromBrowser();
                screenshotBase64 = dataUrl.split(',')[1];
            }
        } catch (e) {
            console.warn("[HealerService] Could not capture screenshot for analysis:", e);
        }

        // 2. Diagnose with LLM
        const systemPrompt = `You are an Autonomous Site Reliability Engineer (Healer).
Your job is to diagnose web application errors and propose fixes.
You have access to the error log and a screenshot of the state.

Return a JSON object with:
{
    "analysis": "Explanation of what went wrong",
    "severity": "low|medium|high|critical",
    "suggestedFix": "Description of code change needed",
    "fileLikelyInvolved": "filename.ts (guess based on stack trace or context)"
}`;

        const userPrompt = `Error: ${errorContext.error}
URL: ${errorContext.url}
Timestamp: ${errorContext.timestamp}

Please analyze this crash.`;

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
            console.log(response.content);

            // FUTURE: Parse JSON and trigger Auto-Dev or notify user via Notification Center

        } catch (e) {
            console.error("[HealerService] Diagnosis Failed:", e);
        }
    }
}
