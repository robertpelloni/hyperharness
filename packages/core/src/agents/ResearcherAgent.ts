import { SpecializedAgent } from '../mesh/SpecializedAgent.js';
import { DeepResearchService } from '../services/DeepResearchService.js';

/**
 * ResearcherAgent
 * A specialized agent for gathering and synthesizing information from the web.
 *
 * Architecture:
 * - Extends `SpecializedAgent` to participate in the Mesh (P2P agent network).
 * - Delegates actual research to `DeepResearchService`, which handles:
 *   1. Query generation via LLM
 *   2. Web search via `WebSearchTool` (DuckDuckGo)
 *   3. Recursive sub-topic exploration (depth/breadth controlled)
 *   4. LLM-powered synthesis of findings into a structured report
 *   5. Memory storage via `MemoryManager`
 *
 * Why it exists:
 * - Separates the "agent identity" (capabilities, mesh registration) from the
 *   "service logic" (research pipeline). The agent is the interface; the service
 *   is the engine.
 */
export class ResearcherAgent extends SpecializedAgent {
    /** The research engine that performs the actual web crawling, LLM synthesis, and memory storage */
    private deepResearchService: DeepResearchService;

    constructor(deepResearchService: DeepResearchService) {
        // Register with capabilities ['research', 'search', 'summarization']
        // so the Director can route matching tasks to this agent
        super('Researcher', ['research', 'search', 'summarization']);
        this.deepResearchService = deepResearchService;
    }

    /**
     * Executes a research task from the Mesh/Director.
     *
     * @param offer - Task offer containing `task` (the research query string)
     * @returns Structured result with summary, sources, and the full recursive report
     */
    public async handleTask(offer: any): Promise<any> {
        console.log(`[ResearcherAgent] 🔍 Investigating query: ${offer.task}`);

        try {
            // Perform real recursive research
            // depth=2: research the topic, then research each sub-topic one more level
            // breadth=3: explore up to 3 related sub-topics per level
            const result = await this.deepResearchService.recursiveResearch(offer.task, 2, 3);

            console.log(`[ResearcherAgent] 📄 Report Ready: ${result.topic}`);

            return {
                status: 'completed',
                // Flatten sources into a simple { source, content } array for mesh consumers
                findings: result.sources.map(s => ({ source: s.title, content: s.url })),
                summary: result.summary,
                fullReport: result,
                // Phase 96: Execution Telemetry
                modelMetadata: result.modelMetadata
            };
        } catch (error: any) {
            console.error(`[ResearcherAgent] 💥 Research failed:`, error);
            throw error;
        }
    }
}

