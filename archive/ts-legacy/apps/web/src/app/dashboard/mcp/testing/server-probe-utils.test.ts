import { describe, expect, it } from 'vitest';

import { buildServerProbeTargets, filterToolsForProbeTarget } from './server-probe-utils';

describe('server probe helpers', () => {
        expect(buildServerProbeTargets([
            { name: 'zeta', toolCount: 2, status: 'ready' },
            { name: 'alpha', toolCount: 4, status: 'connected' },
        ])).toEqual([
            {
                id: 'router',
                kind: 'router',
            },
            {
                id: 'server:alpha',
                kind: 'server',
                label: 'alpha',
                serverName: 'alpha',
                toolCount: 4,
                status: 'connected',
            },
            {
                id: 'server:zeta',
                kind: 'server',
                label: 'zeta',
                serverName: 'zeta',
                toolCount: 2,
                status: 'ready',
            },
        ]);
    });

    it('filters tool choices to the selected downstream server and keeps router targets broad', () => {
        const tools = [
            { name: 'github__search', server: 'github' },
            { name: 'memory__remember', server: 'memory' },
        ];

        expect(filterToolsForProbeTarget(tools, {
            id: 'router',
            kind: 'router',
            description: 'router',
        })).toEqual(tools);

        expect(filterToolsForProbeTarget(tools, {
            id: 'server:github',
            kind: 'server',
            label: 'github',
            description: 'downstream',
            serverName: 'github',
        })).toEqual([
            { name: 'github__search', server: 'github' },
        ]);
    });
});