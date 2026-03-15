import { describe, expect, it } from 'vitest';

import type { NavSection } from './nav-config';
import {
    buildNavItemsByHref,
    buildNavItemsByNormalizedHref,
    extractStringArray,
    hasNavValidationIssues,
    isNavHrefActive,
    matchesNavQuery,
    normalizeNavHref,
    normalizeNavHrefList,
    sanitizeCollapsedSections,
    sanitizeRecentSearches,
    validateSidebarSections,
} from './nav-validation';

describe('normalizeNavHref', () => {
    it('normalizes trailing slashes while preserving root', () => {
        expect(normalizeNavHref('/dashboard/library')).toBe('/dashboard/library');
        expect(normalizeNavHref('/dashboard/library/')).toBe('/dashboard/library');
        expect(normalizeNavHref('/')).toBe('/');
    });

    it('strips query/hash fragments and trims surrounding whitespace', () => {
        expect(normalizeNavHref(' /dashboard/library/?tab=overview#top ')).toBe('/dashboard/library');
        expect(normalizeNavHref('/dashboard/library?tab=overview')).toBe('/dashboard/library');
        expect(normalizeNavHref('/dashboard/library#summary')).toBe('/dashboard/library');
    });
});

describe('normalizeNavHrefList', () => {
    it('deduplicates semantic aliases while preserving first-seen order', () => {
        expect(normalizeNavHrefList([
            '/dashboard/library/?tab=overview',
            '/dashboard/library/',
            '/dashboard/library#top',
            '/dashboard/tools',
        ])).toEqual([
            '/dashboard/library',
            '/dashboard/tools',
        ]);
    });

    it('matches the canonical href produced by a pre-normalized toggle path', () => {
        const rawHref = '/dashboard/library/?tab=overview#top';

        expect(normalizeNavHrefList([rawHref])).toEqual([normalizeNavHref(rawHref)]);
    });
});

describe('extractStringArray', () => {
    it('returns only string entries from unknown input arrays', () => {
        expect(extractStringArray(['one', 2, null, 'two', false])).toEqual(['one', 'two']);
        expect(extractStringArray('not-an-array')).toEqual([]);
    });
});

describe('sanitizeRecentSearches', () => {
    it('trims, deduplicates, and limits recent search history', () => {
        expect(sanitizeRecentSearches([
            '  tools  ',
            'tools',
            '',
            ' billing ',
            'sessions',
            'logs',
        ], 3)).toEqual(['tools', 'billing', 'sessions']);
    });

    it('ignores non-string recent search values', () => {
        expect(sanitizeRecentSearches(['valid', 42, null, ' next '], 5)).toEqual(['valid', 'next']);
    });
});

describe('sanitizeCollapsedSections', () => {
    it('keeps only boolean collapse-state entries from plain objects', () => {
        expect(sanitizeCollapsedSections({
            Favorites: true,
            Recent: false,
            Broken: 'yes',
            Count: 1,
        })).toEqual({
            Favorites: true,
            Recent: false,
        });
    });

    it('returns an empty object for non-object values', () => {
        expect(sanitizeCollapsedSections(null)).toEqual({});
        expect(sanitizeCollapsedSections(['Favorites'])).toEqual({});
    });
});

describe('isNavHrefActive', () => {
    it('treats semantic aliases as active matches', () => {
        expect(isNavHrefActive('/dashboard/library?tab=overview', '/dashboard/library/')).toBe(true);
        expect(isNavHrefActive('/dashboard/library/#top', '/dashboard/library')).toBe(true);
    });

    it('supports nested-route activation but does not over-match root', () => {
        expect(isNavHrefActive('/dashboard/tools/detail', '/dashboard/tools')).toBe(true);
        expect(isNavHrefActive('/dashboard/tools', '/')).toBe(false);
    });
});

describe('matchesNavQuery', () => {
    const candidate = {
        title: 'Library',
        href: '/dashboard/library/?tab=overview',
        description: 'Browse tools and documents',
    };

    it('matches against title, description, section, and normalized href', () => {
        expect(matchesNavQuery('library', candidate, 'Knowledge')).toBe(true);
        expect(matchesNavQuery('documents', candidate, 'Knowledge')).toBe(true);
        expect(matchesNavQuery('knowledge', candidate, 'Knowledge')).toBe(true);
        expect(matchesNavQuery('/dashboard/library', candidate, 'Knowledge')).toBe(true);
    });

    it('treats blank queries as match-all and rejects unrelated terms', () => {
        expect(matchesNavQuery('', candidate, 'Knowledge')).toBe(true);
        expect(matchesNavQuery('totally-unrelated', candidate, 'Knowledge')).toBe(false);
    });
});

describe('validateSidebarSections', () => {
    it('returns empty diagnostics when hrefs are unique', () => {
        const sections: NavSection[] = [
            {
                title: 'One',
                items: [
                    { title: 'A', href: '/a', icon: null, variant: 'ghost' },
                    { title: 'B', href: '/b', icon: null, variant: 'ghost' },
                ],
            },
            {
                title: 'Two',
                items: [
                    { title: 'C', href: '/c', icon: null, variant: 'ghost' },
                    { title: 'D', href: '/d', icon: null, variant: 'ghost' },
                ],
            },
        ];

        const diagnostics = validateSidebarSections(sections);

        expect(diagnostics.duplicateWithinSection).toEqual([]);
        expect(diagnostics.duplicateAcrossSections).toEqual([]);
        expect(diagnostics.normalizedHrefCollisions).toEqual([]);
        expect(hasNavValidationIssues(diagnostics)).toBe(false);
    });

    it('detects duplicates within and across sections', () => {
        const sections: NavSection[] = [
            {
                title: 'One',
                items: [
                    { title: 'A', href: '/repeat', icon: null, variant: 'ghost' },
                    { title: 'B', href: '/repeat', icon: null, variant: 'ghost' },
                    { title: 'C', href: '/x', icon: null, variant: 'ghost' },
                ],
            },
            {
                title: 'Two',
                items: [
                    { title: 'D', href: '/x', icon: null, variant: 'ghost' },
                    { title: 'E', href: '/y', icon: null, variant: 'ghost' },
                ],
            },
        ];

        const diagnostics = validateSidebarSections(sections);

        expect(diagnostics.duplicateWithinSection).toEqual([
            {
                section: 'One',
                duplicates: ['/repeat'],
            },
        ]);

        expect(diagnostics.duplicateAcrossSections).toEqual([
            {
                href: '/x',
                sections: ['One', 'Two'],
            },
        ]);

        expect(diagnostics.normalizedHrefCollisions).toEqual([]);
        expect(hasNavValidationIssues(diagnostics)).toBe(true);
    });

    it('detects normalized href collisions such as trailing slash variants', () => {
        const sections: NavSection[] = [
            {
                title: 'One',
                items: [
                    { title: 'A', href: '/dashboard/library', icon: null, variant: 'ghost' },
                ],
            },
            {
                title: 'Two',
                items: [
                    { title: 'B', href: '/dashboard/library/', icon: null, variant: 'ghost' },
                ],
            },
        ];

        const diagnostics = validateSidebarSections(sections);

        expect(diagnostics.normalizedHrefCollisions).toEqual([
            {
                normalizedHref: '/dashboard/library',
                hrefs: ['/dashboard/library', '/dashboard/library/'],
                sections: ['One', 'Two'],
            },
        ]);
        expect(hasNavValidationIssues(diagnostics)).toBe(true);
    });

    it('detects normalized href collisions for query and hash variants', () => {
        const sections: NavSection[] = [
            {
                title: 'One',
                items: [
                    { title: 'A', href: '/dashboard/library?tab=overview', icon: null, variant: 'ghost' },
                ],
            },
            {
                title: 'Two',
                items: [
                    { title: 'B', href: '/dashboard/library#top', icon: null, variant: 'ghost' },
                ],
            },
        ];

        const diagnostics = validateSidebarSections(sections);

        expect(diagnostics.normalizedHrefCollisions).toEqual([
            {
                normalizedHref: '/dashboard/library',
                hrefs: ['/dashboard/library#top', '/dashboard/library?tab=overview'],
                sections: ['One', 'Two'],
            },
        ]);
        expect(hasNavValidationIssues(diagnostics)).toBe(true);
    });

    it('builds href map with first-seen metadata when duplicates exist', () => {
        const sections: NavSection[] = [
            {
                title: 'One',
                items: [
                    { title: 'Original', href: '/same', icon: null, variant: 'ghost', description: 'first' },
                ],
            },
            {
                title: 'Two',
                items: [
                    { title: 'Shadowed', href: '/same', icon: null, variant: 'ghost', description: 'second' },
                ],
            },
        ];

        const map = buildNavItemsByHref(sections);

        expect(map.get('/same')?.title).toBe('Original');
        expect(map.get('/same')?.description).toBe('first');
    });

    it('builds normalized href map with first-seen metadata for semantic collisions', () => {
        const sections: NavSection[] = [
            {
                title: 'One',
                items: [
                    { title: 'Original', href: '/same/?tab=one#top', icon: null, variant: 'ghost', description: 'first' },
                ],
            },
            {
                title: 'Two',
                items: [
                    { title: 'Shadowed', href: '/same/', icon: null, variant: 'ghost', description: 'second' },
                ],
            },
        ];

        const map = buildNavItemsByNormalizedHref(sections);

        expect(map.get('/same')?.title).toBe('Original');
        expect(map.get('/same')?.description).toBe('first');
    });
});
