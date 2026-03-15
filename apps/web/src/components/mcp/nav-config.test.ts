import { describe, expect, it } from 'vitest';

import { SIDEBAR_SECTIONS } from './nav-config';

describe('SIDEBAR_SECTIONS', () => {
    it('uses unique href values within each section', () => {
        for (const section of SIDEBAR_SECTIONS) {
            const seen = new Set<string>();
            const duplicates: string[] = [];

            for (const item of section.items) {
                if (seen.has(item.href)) {
                    duplicates.push(item.href);
                    continue;
                }
                seen.add(item.href);
            }

            expect(
                duplicates,
                `Duplicate navigation hrefs found in section "${section.title}": ${duplicates.join(', ')}`,
            ).toEqual([]);
        }
    });
});
