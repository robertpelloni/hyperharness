
import { MemoryManager } from '../services/MemoryManager.js';
import { MCPServer } from '../MCPServer.js';
import { search, SafeSearchType } from 'duck-duck-scrape';

/**
 * Deep Research Service
 * Autonomously research a topic by searching, reading, and memorizing.
 */
export class ResearchService {
    private memory: MemoryManager;
    private server: MCPServer;
    private visitedUrls: Set<string> = new Set();

    constructor(server: MCPServer, memory: MemoryManager) {
        this.server = server;
        this.memory = memory;
    }

    /**
     * Conduct deep research on a topic
     * @param topic - The research query
     * @param depth - How many pages to read (default 3)
     */
    public async research(topic: string, depth: number = 3): Promise<string> {
        console.log(`[Research] Starting deep dive on: "${topic}"`);
        const report: string[] = [];
        this.visitedUrls.clear();

        // 1. Initial Search
        const searchResults = await search(topic, {
            safeSearch: SafeSearchType.MODERATE
        });
        const targets = searchResults.results.slice(0, depth);

        report.push(`# Research Report: ${topic}`);
        report.push(`Found ${targets.length} primary sources.`);

        // 2. Process Sources
        for (const target of targets) {
            if (this.visitedUrls.has(target.url)) continue;
            this.visitedUrls.add(target.url);

            console.log(`[Research] Reading: ${target.title}`);
            try {
                // Use existing read_page tool logic via MCPServer
                // We use executeTool to leverage the existing browser bridge if available
                // If browser is not available, we might want a fallback (cheerio/fetch)
                // For now, assume browser or fail gracefully.

                // Note: read_page requires a connected browser extension usually.
                // We'll try to fetch using the server's read_page capability.

                const result = await this.server.executeTool("read_page", { url: target.url });

                // Extract content from result
                // Result format: { content: [{ type: 'text', text: '...' }] }
                // Or error string
                const contentText = result.content?.[0]?.text || "";

                if (contentText.startsWith("Error")) {
                    report.push(`- [FAILED] ${target.title}: ${contentText}`);
                    continue;
                }

                // 3. Memorize
                const ctxId = await this.memory.saveContext(
                    `RESEARCH TOPIC: ${topic}\nSOURCE: ${target.url}\n\n${contentText}`,
                    {
                        title: target.title,
                        source: target.url,
                        type: 'research'
                    }
                );

                report.push(`- [MEMORIZED] ${target.title} (ID: ${ctxId})`);

            } catch (err: any) {
                console.error(`[Research] Failed to process ${target.url}:`, err);
                report.push(`- [ERROR] ${target.title}: ${err.message}`);
            }
        }

        return report.join('\n');
    }
}
