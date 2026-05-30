export { t, publicProcedure, adminProcedure } from './lib/trpc-core.js';
import { t } from './lib/trpc-core.js';
import { suggestionsRouter } from './routers/suggestionsRouter.js';
import { squadRouter } from './routers/squadRouter.js';
import { councilRouter } from './routers/council/index.js';
import { graphRouter } from './routers/graphRouter.js';
import { workflowRouter } from './routers/workflowRouter.js';
import { testsRouter } from './routers/testsRouter.js';
import { contextRouter } from './routers/contextRouter.js';
import { commandsRouter } from './routers/commandsRouter.js';
import { symbolsRouter } from './routers/symbolsRouter.js';
import { autoDevRouter } from './routers/autoDevRouter.js';
import { shellRouter } from './routers/shellRouter.js';
import { memoryRouter } from './routers/memoryRouter.js';
import { skillsRouter } from './routers/skillsRouter.js';
import { researchRouter } from './routers/researchRouter.js';
import { pulseRouter } from './routers/pulseRouter.js';

import { knowledgeRouter } from './routers/knowledgeRouter.js';
import { agentMemoryRouter } from './routers/agentMemoryRouter.js';
import { planRouter as planServiceRouter } from './routers/planRouter.js';
import { metricsRouter as metricsServiceRouter } from './routers/metricsRouter.js';
import { supervisorRouter } from './routers/supervisorRouter.js';
import { lspRouter } from './routers/lspRouter.js';
import { settingsRouter } from './routers/settingsRouter.js';

import { sessionRouter } from './routers/sessionRouter.js';
import { billingRouter } from './routers/billingRouter.js';
import { mcpRouter } from './routers/mcpRouter.js';
import { healerRouter } from './routers/healerRouter.js';
import { darwinRouter } from './routers/darwinRouter.js';
import { autonomyRouter } from './routers/autonomyRouter.js';
import { directorRouter } from './routers/directorRouter.js';
import { directorConfigRouter } from './routers/directorConfigRouter.js';
import { gitRouter } from './routers/gitRouter.js';
import { auditRouter } from './routers/auditRouter.js';
import { submoduleRouter } from './routers/submoduleRouter.js';
import { expertRouter } from './routers/expertRouter.js';
import { systemProcedures } from './routers/systemProcedures.js';
import { mcpServersRouter } from './routers/mcpServersRouter.js';
import { apiKeysRouter } from './routers/apiKeysRouter.js';
import { toolsRouter } from './routers/toolsRouter.js';
import { toolSetsRouter } from './routers/toolSetsRouter.js';
import { logsRouter } from './routers/logsRouter.js';
import { configRouter } from './routers/configRouter.js';
import { serverHealthRouter } from './routers/serverHealthRouter.js';
import { policiesRouter } from './routers/policiesRouter.js';
import { savedScriptsRouter } from './routers/savedScriptsRouter.js';
import { oauthRouter } from './routers/oauthRouter.js';
import { agentRouter } from './routers/agentRouter.js';
import { browserRouter } from './routers/browserRouter.js';
import { meshRouter } from './routers/meshRouter.js';
import { marketplaceRouter } from './routers/marketplaceRouter.js';
import { deerFlowRouter } from './routers/deerFlowRouter.js';
import { cloudDevRouter } from './routers/cloudDevRouter.js';
import { swarmRouter } from './routers/swarmRouter.js';
import { openWebUIRouter } from './routers/openWebUIRouter.js';
import { infrastructureRouter } from './routers/infrastructureRouter.js';
import { ragRouter } from './routers/ragRouter.js';
import { catalogRouter } from './routers/catalogRouter.js';
import { projectRouter } from './routers/projectRouter.js';
import { linksBacklogRouter } from './routers/linksBacklogRouter.js';
import { unifiedDirectoryRouter } from './routers/unifiedDirectoryRouter.js';
import { codeModeRouter } from './routers/codeModeRouter.js';
import { secretsRouter } from './routers/secretsRouter.js';
import { workspaceRouter } from './routers/workspaceRouter.js';
import { browserExtensionRouter } from './routers/browserExtensionRouter.js';
import { sessionExportRouter } from './routers/sessionExportRouter.js';
import { toolChainingRouter } from './routers/toolChainingRouter.js';
import { browserControlsRouter } from './routers/browserControlsRouter.js';
=======
import { promptsRouter } from './routers/promptsRouter.js';
>>>>>>> origin/dependabot/cargo/packages/zed-extension/cargo-64b2a50fd2:packages/core/src/trpc.ts

// import { type AnyTRPCRouter } from '@trpc/server';

export const appRouter = t.router({
    project: projectRouter,
    graph: graphRouter,
    workflow: workflowRouter,
    tests: testsRouter,
});

export type AppRouter = typeof appRouter;
