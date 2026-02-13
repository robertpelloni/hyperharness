
import { BaseAgent, AgentStatus } from './BaseAgent.js';
import { ResearchAgent, CodeAgent } from './SubAgents.js';
import type { MCPServer } from '../MCPServer.js';

export class SpawnerService {
    private static instance: SpawnerService;
    private agents: Map<string, BaseAgent> = new Map();
    private server: MCPServer | null = null;

    private constructor() { }

    public static getInstance(): SpawnerService {
        if (!SpawnerService.instance) {
            SpawnerService.instance = new SpawnerService();
        }
        return SpawnerService.instance;
    }

    public setServer(server: MCPServer) {
        this.server = server;
    }

    public spawn(type: 'research' | 'code' | 'custom', task: string): string {
        if (!this.server) {
            throw new Error("SpawnerService not initialized with MCPServer");
        }

        let agent: BaseAgent;

        switch (type) {
            case 'research':
                agent = new ResearchAgent(task, this.server);
                break;
            case 'code':
                agent = new CodeAgent(task, this.server);
                break;
            default:
                throw new Error(`Unknown agent type: ${type}`);
        }

        this.agents.set(agent.id, agent);

        // Start agent asynchronously
        agent.run().catch(err => {
            console.error(`[Spawner] Agent ${agent.id} crashed:`, err);
        });

        return agent.id;
    }

    public listAgents(): any[] {
        return Array.from(this.agents.values()).map(a => ({
            id: a.id,
            type: a.type,
            status: a.status,
            task: a.task,
            createdAt: a.createdAt,
            result: a.result
        }));
    }

    public getAgent(id: string): BaseAgent | undefined {
        return this.agents.get(id);
    }

    public killAgent(id: string): boolean {
        const agent = this.agents.get(id);
        if (agent && agent.status === AgentStatus.RUNNING) {
            // In a real implementation we'd need an abort controller
            // For now, we just mark as failed
            (agent as any).fail("Terminated by user");
            return true;
        }
        return false;
    }
}
