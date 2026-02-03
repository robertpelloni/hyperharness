import { CogneeClient } from './services/CogneeClient.js';

export interface HelperNode {
    id: string;
    label: string;
    properties: Record<string, any>;
}

export interface HelperEdge {
    source: string;
    target: string;
    relation: string;
    weight: number;
}

export class GraphMemory {
    private client: CogneeClient;

    constructor() {
        this.client = new CogneeClient();
    }

    /**
     * Adds text to the memory and triggers the 'cognitive layer' to extract graph nodes.
     * @param text The natural language text (fact, snippet, observation)
     * @param context Optional context identifier (e.g. 'auth_module')
     */
    async add(text: string, context: string = "borg_memory") {
        await this.client.add(text, context);
        // In background, or separate pipeline step:
        await this.client.cognify(context);
    }

    /**
     * Searches the graph for insights related to the query.
     * Use this when you want to "Connect the dots".
     * @param query 
     */
    async search(query: string): Promise<string[]> {
        return await this.client.search(query);
    }
}
