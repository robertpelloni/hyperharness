
import { LLMService } from '@borg/ai';
import type { MCPServer } from '../MCPServer.js';

export interface Mutation {
    id: string;
    originalPrompt: string;
    mutatedPrompt: string;
    reasoning: string;
    timestamp: number;
}

export interface Experiment {
    id: string;
    mutationId: string;
    task: string;
    status: 'PENDING' | 'RUNNING' | 'COMPLETED';
    resultA?: string;
    resultB?: string;
    winner?: 'A' | 'B' | 'TIE';
    judgeReasoning?: string;
}

export class DarwinService {
    private llm: LLMService;
    private server: MCPServer;
    private mutations: Mutation[] = [];
    private experiments: Experiment[] = [];

    constructor(llm: LLMService, server: MCPServer) {
        this.llm = llm;
        this.server = server;
    }

    public async proposeMutation(originalPrompt: string, goal: string): Promise<Mutation> {
        const prompt = `
        You are The Evolutionary Engine.
        Your goal is to mutate the following system prompt to better achieve: "${goal}".
        Make it more concise, more strict, or more creative as needed.
        
        Original Prompt:
        "${originalPrompt}"
        
        Return JSON:
        {
            "mutatedPrompt": "The new prompt text",
            "reasoning": "Why this might be better"
        }
        `;

        const response = await this.llm.generateText("openai", "gpt-4o", "You are an expert prompt engineer. JSON only.", prompt, {});

        let result;
        try {
            result = JSON.parse((response as any).text || (response as any).content || response);
        } catch (e) {
            throw new Error("Failed to parse mutation suggestion.");
        }

        const mutation: Mutation = {
            id: Math.random().toString(36).substring(7),
            originalPrompt,
            mutatedPrompt: result.mutatedPrompt,
            reasoning: result.reasoning,
            timestamp: Date.now()
        };

        this.mutations.push(mutation);
        return mutation;
    }

    public async startExperiment(mutationId: string, task: string): Promise<Experiment> {
        const mutation = this.mutations.find(m => m.id === mutationId);
        if (!mutation) throw new Error("Mutation not found");

        const experiment: Experiment = {
            id: Math.random().toString(36).substring(7),
            mutationId,
            task,
            status: 'PENDING'
        };
        this.experiments.push(experiment);

        // Asynchronously run the experiment
        this.runExperimentLogic(experiment, mutation).catch(console.error);

        return experiment;
    }

    private async runExperimentLogic(experiment: Experiment, mutation: Mutation) {
        experiment.status = 'RUNNING';

        // 1. Run Agent A (Control)
        console.log(`[Darwin] Running Control (A) for Exp ${experiment.id}...`);
        // Mock execution for V1 - in prod, instantiate ephemeral agents
        const resultA = await this.mockAgentExecution(mutation.originalPrompt, experiment.task);
        experiment.resultA = resultA;

        // 2. Run Agent B (Variant)
        console.log(`[Darwin] Running Variant (B) for Exp ${experiment.id}...`);
        const resultB = await this.mockAgentExecution(mutation.mutatedPrompt, experiment.task);
        experiment.resultB = resultB;

        // 3. Judge
        console.log(`[Darwin] Judging Exp ${experiment.id}...`);
        const judgePrompt = `
        Compare two AI outputs for the task: "${experiment.task}".
        
        Output A (Control):
        ${resultA}
        
        Output B (Variant):
        ${resultB}
        
        Which is better? Return JSON:
        {
            "winner": "A" or "B" or "TIE",
            "reasoning": "Explanation"
        }
        `;

        const judgeRes = await this.llm.generateText("openai", "gpt-4o", "", judgePrompt, {});
        try {
            const verdict = JSON.parse((judgeRes as any).text || (judgeRes as any).content || judgeRes);
            experiment.winner = verdict.winner;
            experiment.judgeReasoning = verdict.reasoning;
        } catch (e) {
            experiment.winner = 'TIE';
            experiment.judgeReasoning = "Failed to parse judge verdict.";
        }

        experiment.status = 'COMPLETED';
        console.log(`[Darwin] Experiment ${experiment.id} Complete. Winner: ${experiment.winner}`);
    }

    private async mockAgentExecution(sysPrompt: string, task: string): Promise<string> {
        // Simulate an agent running a task with a specific system prompt
        // We just ask the LLM directly contextually
        const res = await this.llm.generateText("openai", "gpt-4o", sysPrompt, task, {});
        return (res as any).text || (res as any).content || String(res);
    }

    public getStatus() {
        return {
            mutations: this.mutations,
            experiments: this.experiments
        };
    }
}
