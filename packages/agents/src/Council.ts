import { ModelSelector } from "@borg/ai";
import { IAgent } from "@borg/ai";
import { IMCPServer } from "@borg/adk";

export enum CouncilRole {
    ARCHITECT = 'Architect',
    CRITIC = 'Critic',
    PRODUCT = 'Product',
    META_ARCHITECT = 'MetaArchitect'
}

interface DebateResult {
    approved: boolean;
    transcripts: { speaker: string, text: string }[];
    summary: string;
}

export class Council {
    private agents: Map<string, IAgent> = new Map();
    private server: IMCPServer | undefined;
    public lastResult: DebateResult | null = null;

    constructor(private modelSelector: ModelSelector, agents?: Map<string, IAgent>, server?: IMCPServer) {
        if (agents) {
            this.agents = agents;
        }
        this.server = server;
    }

    public registerAgent(role: string, agent: IAgent) {
        this.agents.set(role, agent);
    }

    public setServer(server: IMCPServer) {
        this.server = server;
    }

    private async broadcast(type: string, payload: any) {
        if (this.server && this.server.broadcastRequest) {
            await this.server.broadcastRequest(type, payload);
        }
    }

    /**
     * Determine if a topic requires Meta-Architect Intervention
     */
    private isMetaTask(topic: string): boolean {
        const metaKeywords = ['create tool', 'modify system', 'refactor core', 'self-evolution', 'meta'];
        return metaKeywords.some(k => topic.toLowerCase().includes(k));
    }

    async runConsensusSession(topic: string): Promise<DebateResult> {
        console.log(`[Council] 🏛️ Convening Consensus Session on: "${topic}"`);
        const transcripts: { speaker: string, text: string }[] = [];
        let context = `Topic: "${topic}"\n\n`;

        await this.broadcast('COUNCIL_START', { topic });

        // Governance Check for Meta-Tasks
        if (this.isMetaTask(topic)) {
            console.log(`[Council] 🔒 Meta-Task Detected. Checking Autonomy Level...`);
            if (this.server && this.server.permissionManager.getAutonomyLevel() !== 'high') {
                const denial = "DENIED: Meta-Tasks require HIGH autonomy level.";
                await this.broadcast('COUNCIL_END', { result: { approved: false, summary: denial } });
                return { approved: false, transcripts: [], summary: denial };
            }
            context += `[SYSTEM NOTICE]: This is a META-TASK. The Meta-Architect will preside.\n\n`;
        }

        // 1. Product Manager (Strategy) - Default to 'product' agent or generic LLM
        // For now, we assume 'product' is handled by generic LLM if no specific agent
        // But let's check for registered agents first.

        // --- Round 1: Product Strategy ---
        // (If Meta-Task, Product Manager defines the User Need for the tool)
        const productResponse = await this.consult("Product Manager", "Strategy & Value", context, "Propose a high-level direction/task based on this topic. Be concrete.");
        context += `[Product Manager]: ${productResponse}\n\n`;
        transcripts.push({ speaker: "Product Manager", text: productResponse });
        await this.broadcast('COUNCIL_TRANSCRIPT', { speaker: "Product Manager", text: productResponse });

        // --- Round 2: Architect (Claude) OR Meta-Architect ---
        let architectRole = CouncilRole.ARCHITECT;
        let architectName = "The Architect";
        if (this.isMetaTask(topic)) {
            architectRole = CouncilRole.META_ARCHITECT;
            architectName = "The Meta-Architect";
        }

        const architectResponse = await this.consult(architectName, "System Design & Code Quality", context, "Critique the proposal for technical feasibility, structural elegance, and maintainability. Be strict.", architectRole);
        context += `[${architectName}]: ${architectResponse}\n\n`;
        transcripts.push({ speaker: architectName, text: architectResponse });
        await this.broadcast('COUNCIL_TRANSCRIPT', { speaker: architectName, text: architectResponse });

        // --- Round 3: Critic (Gemini) ---
        const criticResponse = await this.consult("The Critic", "Security & Risks", context, "Identify risks, edge cases, or security flaws in the architecture proposal. Be pessimistic.", CouncilRole.CRITIC);
        context += `[The Critic]: ${criticResponse}\n\n`;
        transcripts.push({ speaker: "The Critic", text: criticResponse });
        await this.broadcast('COUNCIL_TRANSCRIPT', { speaker: "The Critic", text: criticResponse });

        // --- Round 4: Product Synthesis ---
        const directiveResponse = await this.consult("Product Manager", "Strategy", context, "Synthesize the feedback into a single, actionable DIRECTIVE for the Agent. Start your response with 'DIRECTIVE: ...'", undefined, "DIRECTIVE: STANDBY");
        transcripts.push({ speaker: "Final Directive", text: directiveResponse });
        await this.broadcast('COUNCIL_TRANSCRIPT', { speaker: "Final Directive", text: directiveResponse });

        // Extract directive text
        let finalDirective = directiveResponse;
        const match = finalDirective.match(/DIRECTIVE:\s*(.*)/i);
        if (match) {
            finalDirective = match[1].trim();
        }

        console.log(`[Council] 🏁 Consensus Reached: ${finalDirective.substring(0, 100)}...`);

        const result = {
            approved: true,
            transcripts,
            summary: finalDirective
        };

        this.lastResult = result;

        await this.broadcast('COUNCIL_END', { result });
        return result;
    }

    private async consult(name: string, role: string, context: string, instruction: string, agentRole?: string, fallback: string = "I have no objections."): Promise<string> {
        const agent = agentRole ? this.agents.get(agentRole) : undefined;
        await this.broadcast('COUNCIL_THINKING', { speaker: name });

        const systemPrompt = `You are ${name}, a member of the AI Council.
Role: ${role}

Context of Debate:
${context}

Task: ${instruction}
Keep your response concise (under 4 sentences).`;

        try {
            if (agent && agent.isActive()) {
                console.log(`[Council] 👤 ${name} (Agent): Thinking...`);
                return await agent.send(systemPrompt + "\n\nPlease provide your input.");
            } else {
                // Fallback to Generic LLM
                console.log(`[Council] 👤 ${name} (LLM): Thinking... (No dedicated agent)`);

                if (this.server && (this.server as any).llmService) {
                    // 1. Select Model
                    const selection = await this.modelSelector.selectModel({
                        taskComplexity: 'high', // Council members should be smart
                        provider: 'anthropic' // Prefer Claude for reasoning, but selector will handle fallback
                    });

                    // 2. Generate
                    const response = await (this.server as any).llmService.generateText(
                        selection.provider,
                        selection.modelId,
                        systemPrompt,
                        instruction + "\nPlease provide your input."
                    );

                    return response.content;
                }

                return fallback;
            }
        } catch (e: any) {
            console.error(`[Council] ⚠️ Error consulting ${name}:`, e.message);
            return fallback;
        }
    }
}
