/**
 * ragRouter.ts
 *
 * Exposes RAG endpoints for document ingestion via the Master Control Panel (dashboard).
 */

import { z } from 'zod';
import { t, publicProcedure, getMemoryManager } from '../lib/trpc-core.js';
import { DocumentIntakeService } from '../services/rag/DocumentIntakeService.js';
import { EmbeddingService } from '../services/rag/EmbeddingService.js';

export const ragRouter = t.router({
    ingestFile: publicProcedure.input(z.object({
        filePath: z.string(),
        userId: z.string().default('default'),
        chunkSize: z.number().optional().default(1000),
        chunkOverlap: z.number().optional().default(200),
        strategy: z.enum(['sliding_window', 'recursive', 'semantic']).optional().default('recursive')
    })).mutation(async ({ input }) => {
        const manager = getMemoryManager() as any;

        // Use local embedding by default
        const embedder = new EmbeddingService('local');
        const intakeService = new DocumentIntakeService(manager, embedder);

        try {
            const result = await intakeService.ingestFile(input.filePath, input.userId, {
                chunkSize: input.chunkSize,
                chunkOverlap: input.chunkOverlap,
                strategy: input.strategy
            });

            return { success: true, chunksIngested: result.chunks };
        } catch (e: any) {
            return { success: false, error: e.message };
        }
    }),

    ingestText: publicProcedure.input(z.object({
        text: z.string(),
        sourceName: z.string(),
        userId: z.string().default('default'),
        chunkSize: z.number().optional().default(1000),
        chunkOverlap: z.number().optional().default(200),
        strategy: z.enum(['sliding_window', 'recursive', 'semantic']).optional().default('recursive')
    })).mutation(async ({ input }) => {
        const manager = getMemoryManager() as any;

        const embedder = new EmbeddingService('local');
        const intakeService = new DocumentIntakeService(manager, embedder);

        try {
            const result = await intakeService.ingestText(input.text, input.sourceName, input.userId, {
                chunkSize: input.chunkSize,
                chunkOverlap: input.chunkOverlap,
                strategy: input.strategy
            });

            return { success: true, chunksIngested: result.chunks };
        } catch (e: any) {
            return { success: false, error: e.message };
        }
    })
});
