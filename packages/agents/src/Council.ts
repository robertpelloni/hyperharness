import { ModelSelector } from "@borg/ai";
import { IAgent } from "@borg/ai";
import { IMCPServer } from "@borg/adk";

export enum CouncilRole {
    ARCHITECT = 'Architect',
    CRITIC = 'Critic',
    PRODUCT = 'Product',
    META_ARCHITECT = 'MetaArchitect'
}

export enum ConsensusMode {
    MAJORITY = 'majority',
    SUPERMAJORITY = 'supermajority',
    UNANIMOUS = 'unanimous',
    SINGLE_APPROVAL = 'single-approval'
}

interface TranscriptEntry {
    speaker: string;
    text: string;
    vote?: boolean;
    confidence?: number;
    round?: number;
}

interface DebateResult {
    approved: boolean;
    transcripts: TranscriptEntry[];
    summary: string;
    consensusMode: ConsensusMode;
    votes: { for: number, against: number };
    confidence: number;
}

interface DebateRound {
    roundNumber: number;
    phase: 'opening' | 'rebuttal' | 'closing';
    entries: TranscriptEntry[];
}

interface DebateConfig {
    maxRounds: number;
    allowRebuttals: boolean;
    requireSynthesis: boolean;
    earlyConsensusThreshold: number; // If confidence exceeds this, end early
}

export class Council {
    private agents: Map<string, IAgent> = new Map();
    private server: IMCPServer | undefined;
    public lastResult: DebateResult | null = null;
    private consensusMode: ConsensusMode = ConsensusMode.MAJORITY;

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

    async runConsensusSession(topic: string, config?: Partial<DebateConfig>): Promise<DebateResult> {
        const defaultConfig: DebateConfig = {
            maxRounds: 2,
            allowRebuttals: true,
            requireSynthesis: true,
            earlyConsensusThreshold: 0.95
        };
        const finalConfig = { ...defaultConfig, ...config };

        console.log(`[Council] 🏛️ Convening Debate Session on: "${topic}"`);
        const transcripts: TranscriptEntry[] = [];
        let context = `Topic: "${topic}"\n\n`;
        let votes = { for: 0, against: 0 };

        await this.broadcast('COUNCIL_START', { topic });

        // Governance Check
        if (this.isMetaTask(topic)) {
            if (this.server && this.server.permissionManager.getAutonomyLevel() !== 'high') {
                const denial = "DENIED: Meta-Tasks require HIGH autonomy level.";
                return { approved: false, transcripts: [], summary: denial, consensusMode: this.consensusMode, votes, confidence: 0 };
            }
            context += `[SYSTEM NOTICE]: This is a META-TASK. The Meta-Architect will preside.\n\n`;
        }

        const participants = [
            { name: "Architect", role: CouncilRole.ARCHITECT, instruction: "Propose or refine a technical implementation." },
            { name: "Security Expert", role: CouncilRole.CRITIC, instruction: "Review for security, risks, and edge cases." },
            { name: "QA Lead", role: CouncilRole.PRODUCT, instruction: "Review for quality, testability, and user impact." }
        ];

        // Multi-Round Debate
        for (let round = 1; round <= finalConfig.maxRounds; round++) {
            console.log(`[Council] --- Round ${round} ---`);
            const phase = round === 1 ? 'opening' : (round === finalConfig.maxRounds ? 'closing' : 'rebuttal');

            for (const p of participants) {
                const phaseInstruction = phase === 'opening'
                    ? p.instruction
                    : `Provide a ${phase} to the current debate context. Address points raised by others.`;

                const response = await this.consult(p.name, p.role, context, phaseInstruction, p.role);

                // Extract Vote & Confidence (only in later rounds or if they are very sure)
                const isApproved = !response.toLowerCase().includes("deny") && !response.toLowerCase().includes("reject");

                let confidence = 0.8;
                const confMatch = response.match(/confidence:?\s*(0?\.\d+|1)/i);
                if (confMatch) confidence = parseFloat(confMatch[1]);

                context += `[${p.name} | Round ${round}]: ${response}\n\n`;
                const entry = { speaker: p.name, text: response, vote: isApproved, confidence, round };
                transcripts.push(entry);
                await this.broadcast('COUNCIL_TRANSCRIPT', entry);

                // Early Exit if Unanimous High Confidence (after round 1)
                if (round > 1 && confidence >= finalConfig.earlyConsensusThreshold && isApproved) {
                    // This is a placeholder for more complex early-exit logic
                }
            }
        }

        // Final Synthesis if required
        if (finalConfig.requireSynthesis) {
            const synthesis = await this.synthesize(topic, context);
            context += `[Meta-Architect Synthesis]: ${synthesis}\n\n`;
            transcripts.push({ speaker: "Meta-Architect", text: synthesis, round: finalConfig.maxRounds + 1 });
            await this.broadcast('COUNCIL_SUMMARY', { summary: synthesis });
        }

        // --- Final Consensus Calculation ---
        // We use the last stances of each unique participant
        const lastStances = new Map<string, TranscriptEntry>();
        for (const t of transcripts) {
            if (t.speaker !== "Meta-Architect") {
                lastStances.set(t.speaker, t);
            }
        }

        votes = { for: 0, against: 0 };
        lastStances.forEach(t => {
            if (t.vote) votes.for++; else votes.against++;
        });

        let approved = false;
        const total = votes.for + votes.against;
        const avgConfidence = Array.from(lastStances.values()).reduce((acc, t) => acc + (t.confidence || 0), 0) / (lastStances.size || 1);

        switch (this.consensusMode) {
            case ConsensusMode.MAJORITY: approved = votes.for > total / 2; break;
            case ConsensusMode.SUPERMAJORITY: approved = votes.for >= total * 0.75; break;
            case ConsensusMode.UNANIMOUS: approved = votes.for === total; break;
            case ConsensusMode.SINGLE_APPROVAL: approved = votes.for >= 1; break;
        }

        const finalSummary = approved
            ? `APPROVED via ${this.consensusMode}. Majority votes in favor after ${finalConfig.maxRounds} rounds.`
            : `REJECTED via ${this.consensusMode}. Insufficient consensus reached after ${finalConfig.maxRounds} rounds.`;

        const result: DebateResult = {
            approved,
            transcripts,
            summary: finalSummary,
            consensusMode: this.consensusMode,
            votes,
            confidence: avgConfidence
        };

        this.lastResult = result;
        await this.broadcast('COUNCIL_END', { result });
        return result;
    }

    private async synthesize(topic: string, context: string): Promise<string> {
        console.log(`[Council] 🧠 Meta-Architect Synthesizing results...`);
        const instruction = `Synthesize the debate above into a final decision for topic: "${topic}". Resolve conflicts and provide a clear path forward.`;
        return await this.consult("Meta-Architect", CouncilRole.META_ARCHITECT, context, instruction, undefined);
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
