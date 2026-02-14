import { SpecializedAgent } from '../mesh/SpecializedAgent.js';
import { DeepResearchService } from '../services/DeepResearchService.js';

export class ResearcherAgent extends SpecializedAgent {
    private deepResearchService: DeepResearchService;

    constructor(deepResearchService: DeepResearchService) {
        super('Researcher', ['research', 'search', 'summarization']);
        this.deepResearchService = deepResearchService;
    }

    public async handleTask(offer: any): Promise<any> {
        console.log(`[ResearcherAgent] 🔍 Investigating query: ${offer.task}`);

        try {
            // Perform real recursive research
            // Default depth 2, breadth 3
            const result = await this.deepResearchService.recursiveResearch(offer.task, 2, 3);

            console.log(`[ResearcherAgent] 📄 Report Ready: ${result.topic}`);

            return {
                status: 'completed',
                findings: result.sources.map(s => ({ source: s.title, content: s.url })),
                summary: result.summary,
                fullReport: result
            };
        } catch (error: any) {
            console.error(`[ResearcherAgent] 💥 Research failed:`, error);
            throw error;
        }
    }
}

