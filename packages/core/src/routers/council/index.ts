import { t } from '../../lib/trpc-core.js';
import { sessionRouter } from './sessionRouter.js';
import { councilRouter as baseCouncilRouter } from './councilRouter.ts';
import { quotaRouter } from './quotaRouter.js';
import { historyRouter } from './historyRouter.js';
import { visualRouter } from './visualRouter.js';

export const councilRouter = t.router({
  sessions: sessionRouter,
  base: baseCouncilRouter,
  quota: quotaRouter,
  history: historyRouter,
  visual: visualRouter,
});
