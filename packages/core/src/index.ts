export * from './orchestrator.js';
export { MCPServer } from './MCPServer.js';
export { appRouter } from './trpc.js';
export type { AppRouter } from './trpc.js';
export const name = "@borg/core";
export * from './services/EventBus.js';
export * from './services/ContextPruner.js';
export { SquadService } from './orchestrator/SquadService.js';
export { CouncilService, type CouncilSession } from './services/CouncilService.js';
