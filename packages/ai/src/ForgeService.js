import OpenAI from "openai";

export class ForgeService {
    client;
    baseUrl;

    constructor(config = { baseUrl: "http://localhost:8080/v1" }) {
        this.baseUrl = config.baseUrl;
        this.client = new OpenAI({
            baseURL: this.baseUrl,
            apiKey: config.apiKey || "forge-key",
            dangerouslyAllowBrowser: true
        });
    }

    setUrl(url) {
        this.baseUrl = url;
        this.client = new OpenAI({
            baseURL: this.baseUrl,
            apiKey: this.client.apiKey,
        });
    }

    async listModels() {
        try {
            const list = await this.client.models.list();
            return list.data.map((m) => ({
                id: m.id,
                provider: m.id.split('/')[0] || 'unknown',
                name: m.id
            }));
        }
        catch (error) {
            console.error("Failed to list models from Forge:", error);
            throw error;
        }
    }

    async chatCompletion(params) {
        return this.client.chat.completions.create(params);
    }
}