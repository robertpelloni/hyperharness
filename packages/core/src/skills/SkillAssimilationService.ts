
import fs from 'fs/promises';
import path from 'path';
import { SkillRegistry } from './SkillRegistry.js';
import { LLMService } from '@borg/ai';

export interface ResearchItem {
    id: string;
    name: string;
    url: string;
    summary: string;
    relevance: string;
    category?: string;
}

export class SkillAssimilationService {
    constructor(
        private registry: SkillRegistry,
        private llmService: LLMService,
        private skillsRoot: string
    ) { }

    async assimilate(item: ResearchItem, provider?: string): Promise<string> {
        console.log(`[Assimilation] 🧬 Assimilating Skill: ${item.name}...`);

        // 1. Generate Runbook Content via LLM
        const systemPrompt = `You are a Borg Skill Architect. 
        Your job is to convert research metadata about an AI tool or MCP server into a comprehensive Borg Skill (runbook).
        A Borg Skill is a markdown file with frontmatter and instructions on how to use the tool.`;

        const userPrompt = `Tool Name: ${item.name}
        URL: ${item.url}
        Research Summary: ${item.summary}
        Relevance: ${item.relevance}
        
        Generate a comprehensive SKILL.md content. Include:
        - Frontmatter with name and description.
        - Overview section.
        - Setup instructions (if applicable from summary).
        - Usage examples.
        - Troubleshooting tips.
        
        Output ONLY the markdown content including frontmatter.`;

        try {
            const model = await (this.llmService as any).modelSelector.selectModel({
                taskComplexity: 'medium',
                provider: provider
            });
            const response = await this.llmService.generateText(model.provider, model.modelId, systemPrompt, userPrompt);
            const content = response.content;

            // 2. Create the Skill directory and file
            const skillId = item.id || item.name.toLowerCase().replace(/\s+/g, '-');
            const skillDir = path.join(this.skillsRoot, skillId);
            const skillFile = path.join(skillDir, 'SKILL.md');

            await fs.mkdir(skillDir, { recursive: true });
            await fs.writeFile(skillFile, content, 'utf-8');

            // 3. Register it
            // @ts-ignore
            await this.registry.parseSkill(skillFile);

            return `Skill '${item.name}' assimilated successfully at ${skillFile}`;
        } catch (e: any) {
            console.error(`[Assimilation] ❌ Failed to assimilate ${item.name}:`, e.message);
            throw e;
        }
    }
}
