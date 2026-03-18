import fs from 'fs/promises';
import path from 'path';
import { LLMService } from '@borg/ai';
import AgentMemoryService from './AgentMemoryService.js';

export class SessionImportService {
    private llmService: LLMService;
    private memoryService: AgentMemoryService;
    private workspaceRoot: string;

    constructor(llmService: LLMService, memoryService: AgentMemoryService, workspaceRoot: string = process.cwd()) {
        this.llmService = llmService;
        this.memoryService = memoryService;
        this.workspaceRoot = workspaceRoot;
    }

    /**
     * Scan common locations for AI session logs and import them.
     */
    public async scanAndImport() {
        const locationsToScan = [
            { path: '.aider.chat.history.md', type: 'aider' },
            { path: '.docs/ai-logs', type: 'cursor-windsurf-opencode' },
            { path: '.claude', type: 'claude-code' }
            // Add more common global or local paths here
        ];

        for (const loc of locationsToScan) {
            const fullPath = path.join(this.workspaceRoot, loc.path);
            try {
                const stat = await fs.stat(fullPath);
                if (stat.isDirectory()) {
                    await this.scanDirectory(fullPath, loc.type);
                } else if (stat.isFile()) {
                    await this.importFile(fullPath, loc.type);
                }
            } catch (e: any) {
                // Ignore ENOENT (not found)
                if (e.code !== 'ENOENT') {
                    console.error(`[SessionImport] Error checking ${fullPath}:`, e.message);
                }
            }
        }
    }

    private async scanDirectory(dirPath: string, type: string) {
        try {
            const files = await fs.readdir(dirPath);
            for (const file of files) {
                const fullPath = path.join(dirPath, file);
                const stat = await fs.stat(fullPath);
                // Simple heuristic: only parse markdown, json, or text-based session logs
                if (stat.isFile() && (file.endsWith('.md') || file.endsWith('.json') || file.endsWith('.log') || file.endsWith('.txt'))) {
                    await this.importFile(fullPath, type);
                }
            }
        } catch (e: any) {
            console.error(`[SessionImport] Error scanning directory ${dirPath}:`, e.message);
        }
    }

    public async importFile(filePath: string, type: string) {
        try {
            const content = await fs.readFile(filePath, 'utf-8');
            if (!content || content.trim().length === 0) return;

            // Only process if we haven't already processed it recently (you could hash the file, but for demo we just process)
            console.log(`[SessionImport] Processing session log from ${filePath} (${type})...`);
            
            // Send to LLM to extract important memories
            await this.processSessionWithLLM(content, type, filePath);
            
            // Optionally, we could rename or mark the file as processed, but for now we just ingest.
        } catch (e: any) {
            console.error(`[SessionImport] Error importing file ${filePath}:`, e.message);
        }
    }

    /**
     * Use the LLM to extract actionable knowledge from a raw session transcript.
     */
    private async processSessionWithLLM(content: string, type: string, sourcePath: string) {
        // Truncate if too huge, grab the last chunks which usually contain the resolution
        const MAX_CHARS = 15000;
        const processableContent = content.length > MAX_CHARS ? content.slice(-MAX_CHARS) : content;

        const prompt = `
        You are a memory extraction unit for Borg.
        Analyze the following AI agent session transcript (${type}).
        Extract 1-3 key technical decisions, discovered facts, or persistent architectural rules that should be saved to long-term memory.
        If there is nothing of long-term value, return an empty JSON array [].
        
        Return JSON ONLY in this format:
        [
            {
                "fact": "The project uses XYZ pattern for...",
                "tags": ["architecture", "decision"]
            }
        ]
        
        Transcript:
        ${processableContent}
        `;

        try {
            const response = await this.llmService.generateText('openai', 'gpt-4o-mini', 'You extract memories.', prompt);
            
            // Extract JSON
            const textContent = Array.isArray(response.content) 
                ? (response.content[0] as any).text || '' 
                : typeof response.content === 'string' ? response.content : '';
                
            const start = textContent.indexOf('[');
            const end = textContent.lastIndexOf(']');
            
            if (start !== -1 && end !== -1) {
                const jsonStr = textContent.slice(start, end + 1);
                const facts = JSON.parse(jsonStr);
                
                if (Array.isArray(facts)) {
                    for (const f of facts) {
                        if (f.fact) {
                            console.log(`[SessionImport] ✨ Harvested memory: ${f.fact.substring(0, 50)}...`);
                            await this.memoryService.add(f.fact, 'long_term', 'project', { 
                                source: 'auto_import', 
                                type,
                                path: sourcePath,
                                tags: f.tags || []
                            });
                        }
                    }
                }
            }
        } catch (e: any) {
            console.error(`[SessionImport] LLM extraction failed:`, e.message);
        }
    }
}
