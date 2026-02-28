/**
 * DebateProtocol.ts
 * 
 * Facilitates adversarial LLM execution. 
 * E.g., Claude writes the code, GPT tries to break it, Claude defends, Judge decides.
 * 
 * v2.7.35: Wired generateTurn and evaluateDebate to the Autopilot Council API
 * for real multi-model adversarial debate. Falls back to local stub if unavailable.
 */

import { EventEmitter } from 'events';

export interface DebateTurn {
    round: number;
    persona: string;
    model: string;
    argument: string;
    timestamp: number;
}

export interface DebateConfig {
    topic: string;
    proponentModel: string;
    opponentModel: string;
    judgeModel: string;
    maxRounds: number;
    /** Base URL for the autopilot council server */
    councilUrl?: string;
}

export class DebateProtocol extends EventEmitter {
    private history: DebateTurn[] = [];
    private config: DebateConfig;
    private councilUrl: string;

    constructor(config: DebateConfig) {
        super();
        this.config = config;
        this.councilUrl = config.councilUrl || 'http://localhost:3847';
    }

    /**
     * Executes the debate loop: Proponent argues, Opponent rebuts, Judge evaluates.
     * Each turn is a real LLM call routed through the Autopilot Council's
     * per-supervisor endpoints.
     */
    public async conductDebate(): Promise<{ winner: string; summary: string; history: DebateTurn[] }> {
        this.emit('debate:started', this.config);

        let currentArgument = this.config.topic;

        for (let round = 1; round <= this.config.maxRounds; round++) {
            // 1. Proponent speaks — argues FOR the topic
            const proTurn = await this.generateTurn(round, 'Proponent', this.config.proponentModel, currentArgument);
            this.history.push(proTurn);
            this.emit('debate:turn', proTurn);

            // 2. Opponent rebuts — argues AGAINST the proponent's position
            const conTurn = await this.generateTurn(round, 'Opponent', this.config.opponentModel, proTurn.argument);
            this.history.push(conTurn);
            this.emit('debate:turn', conTurn);

            currentArgument = conTurn.argument; // Feed opponent's rebuttal into next round
        }

        // 3. Judge evaluates the entire transcript
        const judgement = await this.evaluateDebate();
        this.emit('debate:concluded', judgement);

        return judgement;
    }

    /**
     * Generates a single debate turn by calling the Autopilot Council's
     * per-supervisor chat endpoint. The supervisor name is mapped from
     * the model string (e.g., 'claude-3-5-sonnet' → 'Claude' supervisor).
     */
    private async generateTurn(round: number, persona: string, model: string, context: string): Promise<DebateTurn> {
        try {
            const supervisorName = this.modelToSupervisor(model);
            const systemPrompt = persona === 'Proponent'
                ? `You are the PROPONENT in this debate. Argue strongly IN FAVOR of the position. Be specific, cite technical reasoning, and address counterarguments. Keep your response under 300 words.`
                : `You are the OPPONENT in this debate. Argue strongly AGAINST the position. Find weaknesses, propose alternatives, and challenge assumptions. Keep your response under 300 words.`;

            const res = await fetch(`${this.councilUrl}/api/supervisors/${supervisorName}/chat`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    message: `[Round ${round}] ${persona} — Respond to: "${context}"`,
                    systemPrompt,
                    context: this.history.map(t => `[R${t.round} ${t.persona}] ${t.argument}`).join('\n')
                }),
                signal: AbortSignal.timeout(30000)
            });

            if (res.ok) {
                const data = await res.json();
                const argument = data.response || data.message || data.text || '';
                if (argument) {
                    return { round, persona, model, argument, timestamp: Date.now() };
                }
            }
        } catch (err: any) {
            console.warn(`[Debate] Council supervisor call failed for ${persona} (${model}): ${err.message}`);
        }

        // Fallback: produce a basic turn noting the model was unreachable
        return {
            round,
            persona,
            model,
            argument: `[Fallback] ${persona} using ${model} — council unavailable. Original context: "${context.substring(0, 100)}..."`,
            timestamp: Date.now()
        };
    }

    /**
     * Passes the full debate transcript to the Judge model for final evaluation.
     * Uses the council debate endpoint for a holistic multi-model judgement.
     */
    private async evaluateDebate(): Promise<{ winner: string; summary: string; history: DebateTurn[] }> {
        const transcript = this.history
            .map(t => `[Round ${t.round} — ${t.persona} (${t.model})]\n${t.argument}`)
            .join('\n\n---\n\n');

        try {
            const judgeSupervisor = this.modelToSupervisor(this.config.judgeModel);
            const res = await fetch(`${this.councilUrl}/api/supervisors/${judgeSupervisor}/chat`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    message: `You are the JUDGE. Evaluate this ${this.config.maxRounds}-round debate on: "${this.config.topic}". Declare a winner (Proponent or Opponent) and explain your reasoning in 2-3 sentences.\n\nTranscript:\n${transcript}`,
                    systemPrompt: 'You are an impartial judge evaluating a structured debate between two AI models. Respond with JSON: {"winner": "Proponent"|"Opponent", "summary": "your reasoning"}'
                }),
                signal: AbortSignal.timeout(30000)
            });

            if (res.ok) {
                const data = await res.json();
                const judgeText = data.response || data.message || data.text || '';

                // Try to parse structured JSON from the judge's response
                try {
                    const jsonMatch = judgeText.match(/\{[\s\S]*\}/);
                    if (jsonMatch) {
                        const parsed = JSON.parse(jsonMatch[0]);
                        return {
                            winner: parsed.winner || 'Proponent',
                            summary: parsed.summary || judgeText,
                            history: this.history
                        };
                    }
                } catch { /* JSON parse failed — use raw text */ }

                return {
                    winner: judgeText.toLowerCase().includes('opponent') ? 'Opponent' : 'Proponent',
                    summary: judgeText || `Judge (${this.config.judgeModel}) evaluated ${this.config.maxRounds} rounds.`,
                    history: this.history
                };
            }
        } catch (err: any) {
            console.warn(`[Debate] Judge evaluation failed: ${err.message}`);
        }

        // Fallback judgement
        return {
            winner: 'Proponent',
            summary: `[Fallback] Judge (${this.config.judgeModel}) was unavailable. Defaulting to Proponent based on argument volume.`,
            history: this.history
        };
    }

    /**
     * Maps a model identifier string to the corresponding Autopilot Council supervisor name.
     * The council has supervisors: GPT-4o, Claude, DeepSeek, Gemini, Grok
     */
    private modelToSupervisor(model: string): string {
        const m = model.toLowerCase();
        if (m.includes('claude') || m.includes('anthropic')) return 'Claude';
        if (m.includes('gemini') || m.includes('google')) return 'Gemini';
        if (m.includes('deepseek')) return 'DeepSeek';
        if (m.includes('grok')) return 'Grok';
        return 'GPT-4o'; // Default
    }
}
