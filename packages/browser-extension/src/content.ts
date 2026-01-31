import { Readability } from '@mozilla/readability';
import TurndownService from 'turndown';

// Listen for messages from the background script
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    if (request.action === 'SCRAPE_PAGE') {
        try {
            // 1. Clone document to avoid modifying the live page
            const documentClone = document.cloneNode(true) as Document;

            // 2. Parse with Readability
            const reader = new Readability(documentClone);
            const article = reader.parse();

            if (!article) {
                sendResponse({ success: false, error: 'Readability failed to parse article' });
                return;
            }

            // 3. Convert to Markdown
            const turndownService = new TurndownService({
                headingStyle: 'atx',
                codeBlockStyle: 'fenced',
                hr: '---',
                bulletListMarker: '-'
            });
            const markdown = turndownService.turndown(article.content);

            // 4. Return result
            sendResponse({
                success: true,
                data: {
                    title: article.title,
                    byline: article.byline,
                    content: markdown,
                    url: window.location.href,
                    timestamp: new Date().toISOString()
                }
            });
        } catch (error) {
            sendResponse({ success: false, error: (error as Error).message });
        }
    }
    // Return true to indicate we will send a response asynchronously
    return true;
});

console.log('Borg Scraper Content Script Loaded 🚀');
