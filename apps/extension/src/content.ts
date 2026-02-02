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
});
