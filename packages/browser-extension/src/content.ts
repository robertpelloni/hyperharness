import { Readability } from '@mozilla/readability';
import TurndownService from 'turndown';

console.log("Borg Browser Extension Content Script Loaded 🚀");

// 0. Console Hook Injection
// Intercept console.log/warn/error and forward to window
const consoleScript = document.createElement('script');
consoleScript.textContent = `
(function() {
    const originalLog = console.log;
    const originalWarn = console.warn;
    const originalError = console.error;

    function formatArgs(args) {
        return args.map(arg => {
            if (typeof arg === 'object') {
                try {
                    return JSON.stringify(arg);
                } catch(e) {
                    return String(arg);
                }
            }
            return String(arg);
        }).join(' ');
    }

    console.log = function(...args) {
        window.postMessage({ type: 'BORG_CONSOLE_LOG', level: 'info', content: formatArgs(args) }, '*');
        originalLog.apply(console, args);
    };

    console.warn = function(...args) {
        window.postMessage({ type: 'BORG_CONSOLE_LOG', level: 'warn', content: formatArgs(args) }, '*');
        originalWarn.apply(console, args);
    };

    console.error = function(...args) {
        window.postMessage({ type: 'BORG_CONSOLE_LOG', level: 'error', content: formatArgs(args) }, '*');
        originalError.apply(console, args);
    };
    
    // Also inject borg globals if needed
    window.borg = {
        callTool: function(name, args) {
             return new Promise((resolve, reject) => {
                const id = Math.random().toString(36).substring(7);
                const handler = (event) => {
                    if (event.source !== window || !event.data || event.data.type !== 'BORG_MCP_RESPONSE') return;
                    if (event.data.payload.id !== id) return;
                    window.removeEventListener('message', handler);
                    if (event.data.payload.error) reject(event.data.payload.error);
                    else resolve(event.data.payload.result);
                };
                window.addEventListener('message', handler);
                window.postMessage({
                    type: 'BORG_MCP_CALL',
                    payload: { jsonrpc: '2.0', method: 'tools/call', params: { name, arguments: args }, id }
                }, '*');
            });
        }
    };
    console.log("✅ window.borg injected");
    
})();
`;
(document.head || document.documentElement).appendChild(consoleScript);
consoleScript.remove();


// 1. Listen for Messages from Page (Console & MCP Calls)
window.addEventListener('message', (event) => {
    if (event.source !== window || !event.data) return;

    if (event.data.type === 'BORG_MCP_CALL') {
        const payload = event.data.payload;
        chrome.runtime.sendMessage({ type: "MCP_REQUEST", payload }, (response) => {
            window.postMessage({ type: 'BORG_MCP_RESPONSE', payload: response }, '*');
        });
    }

    if (event.data.type === 'BORG_CONSOLE_LOG') {
        // Forward to background (fire and forget)
        chrome.runtime.sendMessage({
            type: "CONSOLE_LOG",
            level: event.data.level,
            content: event.data.content,
            url: window.location.href,
            timestamp: new Date().toISOString()
        });
    }
});


// 2. Listen for messages from Background (Scraper & Reverse Control)
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    if (request.action === 'SCRAPE_PAGE') {
        try {
            const documentClone = document.cloneNode(true) as Document;
            const reader = new Readability(documentClone);
            const article = reader.parse();

            if (!article) {
                sendResponse({ success: false, error: 'Readability failed to parse article' });
                return;
            }

            const turndownService = new TurndownService({
                headingStyle: 'atx',
                codeBlockStyle: 'fenced',
                hr: '---',
                bulletListMarker: '-'
            });
            const markdown = turndownService.turndown(article.content);

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
    return true;
});
