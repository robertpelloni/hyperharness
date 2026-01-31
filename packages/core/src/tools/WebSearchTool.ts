
import { search, SafeSearchType, SearchTimeType } from 'duck-duck-scrape';

export const WebSearchTool = {
    name: "search_web",
    description: "Search the public web using DuckDuckGo. Returns top results with titles, links, and snippets.",
    inputSchema: {
        type: "object",
        properties: {
            query: { type: "string", description: "The search query" },
            limit: { type: "number", description: "Max results (default: 5)" }
        },
        required: ["query"]
    },
    handler: async (args: { query: string, limit?: number }) => {
        try {
            console.log(`[WebSearch] Searching for: ${args.query}`);
            const results = await search(args.query, {
                safeSearch: SafeSearchType.MODERATE,
                time: SearchTimeType.ALL,
            });

            const limit = args.limit || 5;
            const top = results.results.slice(0, limit).map(r => ({
                title: r.title,
                url: r.url,
                snippet: r.description
            }));

            return {
                content: [{
                    type: "text",
                    text: JSON.stringify(top, null, 2)
                }]
            };
        } catch (error: any) {
            console.error("[WebSearch] Error:", error);
            return {
                content: [{
                    type: "text",
                    text: `Error performing search: ${error.message}`
                }]
            };
        }
    }
};
