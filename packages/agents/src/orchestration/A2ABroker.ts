/**
 * @file A2ABroker.ts
 * @module packages/agents/src/orchestration/A2ABroker
 *
 * WHAT: Central message broker for the Agent-to-Agent (A2A) protocol.
 * Manages agent registration and routes messages between them.
 *
 * WHY: Decouples agents from each other, allowing them to communicate via
 * IDs and roles without knowing implementation details.
 */

import { A2AMessage, A2AMessageType, IA2AClient } from "@hypercode/adk";
import { EventEmitter } from "events";

export class A2ABroker extends EventEmitter {
    private agents: Map<string, IA2AClient> = new Map();
    private history: A2AMessage[] = [];
    private readonly MAX_HISTORY = 1000;

    constructor() {
        super();
    }

    /**
     * Register an agent with the broker
     */
    public registerAgent(id: string, client: IA2AClient) {
        this.agents.set(id, client);
        console.log(`[A2A Broker] Registered agent: ${id}`);
        
        client.onMessage((msg) => this.routeMessage(msg));
        
        this.emit('agent_registered', { id });
    }

    /**
     * Unregister an agent
     */
    public unregisterAgent(id: string) {
        this.agents.delete(id);
        console.log(`[A2A Broker] Unregistered agent: ${id}`);
        this.emit('agent_unregistered', { id });
    }

    /**
     * Routes a message to its intended recipient or broadcasts it
     */
    public async routeMessage(message: A2AMessage) {
        this.history.push(message);
        if (this.history.length > this.MAX_HISTORY) {
            this.history.shift();
        }

        this.emit('message_routed', message);

        if (message.recipient) {
            const agent = this.agents.get(message.recipient);
            if (agent) {
                await agent.sendMessage(message);
            } else {
                console.warn(`[A2A Broker] Recipient not found: ${message.recipient}. Bridging to Mesh...`);
                // Phase 102: Mesh Bridge - If recipient not local, emit for MeshService to handle
                this.emit('bridge_to_mesh', message);
            }
        } else {
            // Broadcast to all local agents except sender
            for (const [id, agent] of this.agents.entries()) {
                if (id !== message.sender) {
                    await agent.sendMessage(message);
                }
            }
            // Also broadcast to Mesh
            this.emit('bridge_to_mesh', message);
        }
    }

    public getHistory(): A2AMessage[] {
        return [...this.history];
    }

    public listAgents(): string[] {
        return Array.from(this.agents.keys());
    }
}

export const a2aBroker = new A2ABroker();
