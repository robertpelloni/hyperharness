// Content Script
// Injects the 'borg' object into the page context
import { Readability } from '@mozilla/readability';
import TurndownService from 'turndown';

// Helper to inject a script tag for page-context access
function injectScript(file_path: string) {
    const script = document.createElement('script');
    script.setAttribute('type', 'text/javascript');
    script.setAttribute('src', file_path);
    (document.head || document.documentElement).appendChild(script);
}

// Listen for messages from Page (window.postMessage)
window.addEventListener("message", (event) => {
    if (event.source !== window) return;
    if (event.data.type && (event.data.type === "BORG_REQUEST")) {
        // Forward to Background
        chrome.runtime.sendMessage({
            type: 'EXECUTE_TOOL',
            tool: event.data.tool,
            args: event.data.args
        }, (response) => {
            // Send result back to Page
            window.postMessage({
                type: "BORG_RESPONSE",
                id: event.data.id,
                result: response
            }, "*");
        });
    }
});

console.log("[Borg Bridge] Content Script Loaded.");

// Inject the API
const injectAPI = () => {
    const script = document.createElement('script');
    script.textContent = `
    console.log("[Borg Bridge] API Injected");
    `;
    (document.head || document.documentElement).appendChild(script);
};
injectAPI();

// Console Log Interceptor - Inject into page context
const injectConsoleInterceptor = () => {
    const script = document.createElement('script');
    script.textContent = `
    (function() {
        const originalConsole = {
            log: console.log.bind(console),
            error: console.error.bind(console),
            warn: console.warn.bind(console),
            info: console.info.bind(console)
        };
        
        const forwardLog = (level, args) => {
            try {
                window.postMessage({
                    type: 'BORG_CONSOLE_LOG',
                    level: level,
                    message: Array.from(args).map(a => {
                        try { return typeof a === 'object' ? JSON.stringify(a) : String(a); }
                        catch { return '[Object]'; }
                    }).join(' ')
                }, '*');
            } catch (e) {}
        };
        
        console.log = function(...args) { forwardLog('log', args); return originalConsole.log(...args); };
        console.error = function(...args) { forwardLog('error', args); return originalConsole.error(...args); };
        console.warn = function(...args) { forwardLog('warn', args); return originalConsole.warn(...args); };
        console.info = function(...args) { forwardLog('info', args); return originalConsole.info(...args); };
    })();
    `;
    (document.head || document.documentElement).appendChild(script);
};
injectConsoleInterceptor();

let lastActivitySentAt = 0;
const ACTIVITY_THROTTLE_MS = 1500;

function emitUserActivity(trigger: string) {
    const now = Date.now();
    if ((now - lastActivitySentAt) < ACTIVITY_THROTTLE_MS) {
        return;
    }

    lastActivitySentAt = now;
    chrome.runtime.sendMessage({
        type: 'USER_ACTIVITY',
        lastActivityTime: now,
        trigger,
        activePage: {
            url: window.location.href,
            title: document.title,
            host: window.location.host,
        },
    });
}

window.addEventListener('focus', () => emitUserActivity('focus'));
window.addEventListener('click', () => emitUserActivity('click'), true);
window.addEventListener('keydown', () => emitUserActivity('keydown'), true);
window.addEventListener('scroll', () => emitUserActivity('scroll'), { passive: true });
document.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'visible') {
        emitUserActivity('visible');
    }
});
emitUserActivity('content_script_loaded');

// Listen for console logs from page context
window.addEventListener('message', (event) => {
    if (event.source !== window) return;
    if (event.data.type === 'BORG_CONSOLE_LOG') {
        chrome.runtime.sendMessage({
            type: 'CONSOLE_LOG',
            level: event.data.level,
            message: event.data.message
        });
    }
});

// Listen for messages from Background (WebSocket events or Popup)
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    if (message.type === 'PASTE_INTO_CHAT') {
        console.log("[Borg Bridge] Received Paste Request:", message.text);

        // Inject script to call window.injectDirectorMessage
        // We escape the text carefully to avoid syntax errors
        const safeText = JSON.stringify(message.text);
        const autoSubmit = message.submit ? 'true' : 'false';

        const script = document.createElement('script');
        script.textContent = `
            if (window.injectDirectorMessage) {
                window.injectDirectorMessage(${safeText}, ${autoSubmit});
            } else {
                console.warn("[Borg Bridge] window.injectDirectorMessage not found! Is DirectorChat mounted?");
            }
        `;
        (document.head || document.documentElement).appendChild(script);
        script.remove(); // Cleanup
        sendResponse({ success: true });
    }

    if (message.type === 'ABSORB_PAGE') {
        console.log("[Borg Bridge] Absorbing page...");
        try {
            const clone = document.cloneNode(true) as Document;
            const reader = new Readability(clone);
            const article = reader.parse();

            if (!article) {
                sendResponse({ success: false, error: 'Readability failed to parse page.' });
                return true; // async
            }

            const turndownService = new TurndownService({
                headingStyle: 'atx',
                codeBlockStyle: 'fenced'
            });

            const markdown = turndownService.turndown(article.content);
            const title = article.title || document.title;

            // Add metadata header
            const finalContent = `# ${title}\n\n> URL: ${window.location.href}\n> Excerpt: ${article.excerpt || 'None'}\n\n${markdown}`;

            sendResponse({ success: true, content: finalContent, title: title });
        } catch (error: any) {
            console.error('[Borg Bridge] Absorb Error:', error);
            sendResponse({ success: false, error: error.message });
        }
        return true; // Keep channel open for async response
    }

    if (message.type === 'SCRAPE_CHAT') {
        console.log("[Borg Bridge] Scraping chat history...");
        // This is highly dependent on the UI structure
        // Assuming there are elements with a specific class or we can just get all text
        const messages = Array.from(document.querySelectorAll('.chat-message, .message-content, [data-message-id]'))
            .map(el => el.textContent?.trim())
            .filter(t => !!t && t.length > 0);

        sendResponse({ success: true, history: messages });
        return true;
    }
});
