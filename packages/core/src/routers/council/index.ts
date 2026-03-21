import { t } from '../../lib/trpc-core.js';
import { sessionRouter } from './sessionRouter.js';
import { councilRouter as baseCouncilRouter } from './councilRouter.ts';
import { quotaRouter } from './quotaRouter.js';
import { historyRouter } from './historyRouter.js';
import { visualRouter } from './visualRouter.js';
import { smartPilotRouter } from './smartPilotRouter.js';
import { evolutionRouter } from './evolutionRouter.js';
import { hooksRouter } from './hooksRouter.js';
import { fineTuneRouter } from './fineTuneRouter.js';
import { ideRouter } from './ideRouter.js';

export const councilRouter = t.router({
  sessions: sessionRouter,
  base: baseCouncilRouter,
  quota: quotaRouter,
  history: historyRouter,
  visual: visualRouter,
  smartPilot: smartPilotRouter,
  evolution: evolutionRouter,
  hooks: hooksRouter,
  fineTuning: fineTuneRouter,
  ide: ideRouter,
});
