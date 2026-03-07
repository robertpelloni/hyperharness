
// Popup Logic
console.log("Popup loaded");

const statusSpan = document.getElementById('status')!;
const detailSpan = document.getElementById('detail')!;
const endpointSpan = document.getElementById('endpoints')!;
const urlInput = document.getElementById('urlInput') as HTMLInputElement;
const ingestUrlBtn = document.getElementById('ingestUrlBtn')!;
const saveBtn = document.getElementById('saveBtn')!;
const ragBtn = document.getElementById('ragBtn')!;
const settingsBtn = document.getElementById('settingsBtn')!;
const msgDiv = document.getElementById('msg')!;

function setActionState(enabled: boolean) {
    if (enabled) {
        ingestUrlBtn.removeAttribute('disabled');
        saveBtn.removeAttribute('disabled');
        ragBtn.removeAttribute('disabled');
    } else {
        ingestUrlBtn.setAttribute('disabled', 'true');
        saveBtn.setAttribute('disabled', 'true');
        ragBtn.setAttribute('disabled', 'true');
    }
}

function setStatus(response: { connected?: boolean; coreUrl?: string; wsUrl?: string; status?: string; error?: string } | undefined) {
    const connected = !!response?.connected;
    statusSpan.textContent = connected ? 'ONLINE' : 'OFFLINE';
    statusSpan.style.color = connected ? 'green' : 'red';
    detailSpan.textContent = connected
        ? `Core status: ${response?.status || 'online'}`
        : `Core unreachable${response?.error ? ` — ${response.error}` : '. Start Borg Core or update settings.'}`;
    endpointSpan.textContent = `HTTP ${response?.coreUrl || 'http://localhost:3001'} • WS ${response?.wsUrl || 'ws://localhost:3001'}`;
    setActionState(connected);
}

function refreshConnection() {
    chrome.runtime.sendMessage({ type: 'CHECK_CONNECTION' }, (response) => {
        setStatus(response);
    });
}

refreshConnection();

chrome.tabs.query({ active: true, currentWindow: true }, ([tab]) => {
    if (tab?.url && !urlInput.value) {
        urlInput.value = tab.url;
    }
});

settingsBtn.addEventListener('click', async () => {
    await chrome.runtime.openOptionsPage();
});

ingestUrlBtn.addEventListener('click', async () => {
    const url = urlInput.value.trim();
    if (!url) {
        msgDiv.textContent = 'Enter a URL to ingest first.';
        return;
    }

    msgDiv.textContent = 'Ingesting URL into Borg Knowledge...';
    chrome.runtime.sendMessage({ type: 'INGEST_URL', url }, (res) => {
        if (res && res.success) {
            msgDiv.textContent = '✅ URL ingested into Knowledge!';
            setTimeout(() => msgDiv.textContent = '', 3000);
        } else {
            msgDiv.textContent = 'URL Ingest Error: ' + (res?.error || 'Core unreachable');
        }
    });
});

async function scrapeActivePage() {
    const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
    if (!tab.id) {
        throw new Error('No active tab');
    }

    return new Promise<{ tab: chrome.tabs.Tab; content: string; title: string }>((resolve, reject) => {
        chrome.tabs.sendMessage(tab.id!, { type: 'ABSORB_PAGE' }, (response) => {
            if (chrome.runtime.lastError) {
                reject(new Error('Content script not ready. Try reloading page.'));
                return;
            }

            if (!response || !response.success) {
                reject(new Error(response?.error || 'Unknown scrape error'));
                return;
            }

            resolve({ tab, content: response.content, title: response.title });
        });
    });
}

// 2. Save Context
saveBtn.addEventListener('click', async () => {
    msgDiv.textContent = 'Scraping Page...';

    try {
        const { tab, content, title } = await scrapeActivePage();
        msgDiv.textContent = 'Saving to Borg Core...';

        chrome.runtime.sendMessage({
            type: 'SAVE_CONTEXT',
            content,
            url: tab.url,
            title
        }, (res) => {
            if (res && res.success) {
                msgDiv.textContent = '✅ Absorbed into Memory!';
                setTimeout(() => msgDiv.textContent = '', 3000);
            } else {
                msgDiv.textContent = 'Error Saving: ' + (res?.error || 'Core unreachable');
            }
        });
    } catch (e: any) {
        msgDiv.textContent = 'Exception: ' + e.message;
    }
});

// 3. Ingest into RAG
ragBtn.addEventListener('click', async () => {
    msgDiv.textContent = 'Scraping Page for RAG...';

    try {
        const { tab, content, title } = await scrapeActivePage();
        msgDiv.textContent = 'Ingesting into Borg RAG...';

        chrome.runtime.sendMessage({
            type: 'INGEST_RAG_TEXT',
            text: content,
            sourceName: title || tab.url || 'browser-extension-page'
        }, (res) => {
            if (res && res.success) {
                const chunks = res.data?.chunksIngested ?? 0;
                msgDiv.textContent = `✅ Ingested into RAG (${chunks} chunks)!`;
                setTimeout(() => msgDiv.textContent = '', 3000);
            } else {
                msgDiv.textContent = 'RAG Ingest Error: ' + (res?.error || 'Core unreachable');
            }
        });
    } catch (e: any) {
        msgDiv.textContent = 'Exception: ' + e.message;
    }
});

window.addEventListener('focus', refreshConnection);
