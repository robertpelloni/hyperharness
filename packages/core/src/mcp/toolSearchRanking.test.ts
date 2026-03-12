import { describe, expect, it } from 'vitest';

import { pickAutoLoadCandidate, rankToolSearchCandidates } from './toolSearchRanking.js';

describe('toolSearchRanking auto-load decisions', () => {
    it('picks a high-confidence exact match for auto-load', () => {
        const rankedResults = rankToolSearchCandidates([
            {
                name: 'browser__open_tab',
                description: 'Open a browser tab',
            },
            {
                name: 'browser__close_tab',
                description: 'Close a browser tab',
            },
        ], 'browser__open_tab', 10);

        expect(pickAutoLoadCandidate(rankedResults, 'browser__open_tab')).toEqual({
            toolName: 'browser__open_tab',
            reason: 'auto-loaded after exact tool name match',
        });
    });

    it('does not auto-load when the top result is too ambiguous', () => {
        const rankedResults = rankToolSearchCandidates([
            {
                name: 'browser__open_tab',
                description: 'Open a browser tab',
            },
            {
                name: 'browser__open_window',
                description: 'Open a browser window',
            },
        ], 'open', 10);

        expect(pickAutoLoadCandidate(rankedResults, 'open')).toBeNull();
    });

    it('matches semantic tags and group labels for intent-style queries', () => {
        const rankedResults = rankToolSearchCandidates([
            {
                name: 'browser__open_tab',
                description: 'Open a browser tab',
                advertisedName: 'browser [browser, search] -> open_tab [browser, read]',
                serverTags: ['browser', 'search'],
                toolTags: ['browser', 'read'],
                semanticGroup: 'browser-automation',
                semanticGroupLabel: 'browser automation',
                keywords: ['browser', 'automation', 'tab', 'open'],
            },
            {
                name: 'filesystem__read_file',
                description: 'Read a file from disk',
                serverTags: ['filesystem'],
                toolTags: ['filesystem', 'read'],
                semanticGroup: 'filesystem-operations',
                semanticGroupLabel: 'filesystem operations',
                keywords: ['file', 'disk', 'read'],
            },
        ], 'browser automation', 10);

        expect(rankedResults[0]).toMatchObject({
            name: 'browser__open_tab',
            semanticGroup: 'browser-automation',
        });
        expect(rankedResults[0]?.matchReason).toMatch(/semantic group match|semantic tag match/);
    });
});
