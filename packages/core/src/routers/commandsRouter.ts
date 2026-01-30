import { z } from 'zod';
import { t, publicProcedure } from '../lib/trpc-core.js';

export const commandsRouter = t.router({
    execute: publicProcedure.input(z.object({
        input: z.string()
    })).mutation(async ({ input }) => {
        // @ts-ignore
        if (global.mcpServerInstance?.commandRegistry) {
            // @ts-ignore
            const result = await global.mcpServerInstance.commandRegistry.execute(input.input);
            if (result) {
                return {
                    handled: result.handled,
                    output: result.output,
                    error: result.error
                };
            }
            return { handled: false, output: 'Not a command', error: null };
        }
        return { handled: false, output: 'CommandRegistry not initialized', error: 'Service unavailable' };
    }),

    list: publicProcedure.query(() => {
        // @ts-ignore
        if (global.mcpServerInstance?.commandRegistry) {
            // @ts-ignore
            const commands = global.mcpServerInstance.commandRegistry.getCommands();
            return commands.map((cmd: any) => ({
                name: cmd.name,
                description: cmd.description
            }));
        }
        return [];
    }),
});
