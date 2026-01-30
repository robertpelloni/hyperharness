console.log("aios Browser Extension Content Script Loaded");

// 1. Inject the Bridge Script into the Page Context
const script = document.createElement('script');
script.textContent = `
(function() {
    window.aios = {
        callTool: function(name, args) {
            return new Promise((resolve, reject) => {
                const id = Math.random().toString(36).substring(7);
                
                // Handler for Response
                const handler = (event) => {
                    if (event.source !== window || !event.data || event.data.type !== 'AIOS_MCP_RESPONSE') return;
                    if (event.data.payload.id !== id) return;
                    
                    window.removeEventListener('message', handler);
                    
                    if (event.data.payload.error) {
                        reject(event.data.payload.error);
                    } else {
                        resolve(event.data.payload.result);
                    }
                };
                
                window.addEventListener('message', handler);
                
                // Send Request
                window.postMessage({
                    type: 'AIOS_MCP_CALL',
                    payload: {
                        jsonrpc: '2.0',
                        method: 'tools/call',
                        params: { name, arguments: args },
                        id
                    }
                }, '*');
            });
        },
        
        listTools: function() {
             return new Promise((resolve, reject) => {
                const id = Math.random().toString(36).substring(7);
                const handler = (event) => {
                    if (event.source !== window || !event.data || event.data.type !== 'AIOS_MCP_RESPONSE') return;
                    if (event.data.payload.id !== id) return;
                    window.removeEventListener('message', handler);
                    if (event.data.payload.error) reject(event.data.payload.error);
                    else resolve(event.data.payload.result);
                };
                window.addEventListener('message', handler);
                window.postMessage({
                    type: 'AIOS_MCP_CALL',
                    payload: { jsonrpc: '2.0', method: 'tools/list', id }
                }, '*');
            });
        }
    };
    console.log("✅ window.aios injected");
})();
`;
(document.head || document.documentElement).appendChild(script);
script.remove();

// 2. Listen for Messages from Page
window.addEventListener('message', (event) => {
    if (event.source !== window || !event.data) return;

    if (event.data.type === 'AIOS_MCP_CALL') {
        const payload = event.data.payload;

        // Forward to Background
        chrome.runtime.sendMessage({ type: "MCP_REQUEST", payload }, (response) => {
            // Forward Response back to Page
            window.postMessage({
                type: 'AIOS_MCP_RESPONSE',
                payload: response
            }, '*');
        });
    }
});

// Listen for messages from background script (Reverse Control)
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    if (request.action === "ping") {
        sendResponse({ status: "pong" });
    }
});
