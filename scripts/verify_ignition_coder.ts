
import { CoderAgent } from '../packages/core/src/agents/CoderAgent.js';
import { LLMService, ModelSelector } from '../packages/ai/src/index.js'; // Bypass package exports
import fs from 'fs/promises';
import path from 'path';

async function main() {
    console.log("🔥 IGNITION: Verifying CoderAgent...");

    // Mock/Real dependencies
    const selector = new ModelSelector();
    const llm = new LLMService(selector);
    const agent = new CoderAgent(llm);

    const testFile = 'ignition_test.txt';
    const task = `Create a file named ${testFile} with content 'Systems Online. Ignition Successful.'`;

    try {
        const result = await agent.handleTask({ task });
        console.log("Agent Result:", result);

        // Verify file existence
        const filePath = path.resolve(process.cwd(), testFile);
        const content = await fs.readFile(filePath, 'utf-8');

        if (content.includes('Ignition Successful')) {
            console.log("✅ SUCCESS: File created with correct content.");
            // Cleanup
            await fs.unlink(filePath);
        } else {
            console.error("❌ FAILURE: File content mismatch:", content);
            process.exit(1);
        }

    } catch (e) {
        console.error("❌ FAILURE: Agent execution failed", e);
        process.exit(1);
    }
}

main();
