
// Popup Logic
console.log("Popup loaded");

const statusSpan = document.getElementById('status')!;
const saveBtn = document.getElementById('saveBtn')!;
const msgDiv = document.getElementById('msg')!;

// 1. Check Connection
chrome.runtime.sendMessage({ type: 'CHECK_CONNECTION' }, (response) => {
    if (response && response.connected) {
        statusSpan.textContent = 'ONLINE';
        statusSpan.style.color = 'green';
        saveBtn.removeAttribute('disabled');
    } else {
        statusSpan.textContent = 'OFFLINE (Is borg start running?)';
        statusSpan.style.color = 'red';
        saveBtn.setAttribute('disabled', 'true');
    }
});

// 2. Save Context
saveBtn.addEventListener('click', async () => {
    msgDiv.textContent = 'Scraping Page...';

    // Get Active Tab
    const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
    if (!tab.id) return;

    // Send ABSORB_PAGE to content script
    try {
        chrome.tabs.sendMessage(tab.id, { type: 'ABSORB_PAGE' }, (response) => {
            if (chrome.runtime.lastError) {
                msgDiv.textContent = 'Error: Content script not ready. Try reloading page.';
                console.error(chrome.runtime.lastError);
                return;
            }

            if (!response || !response.success) {
                msgDiv.textContent = 'Failed: ' + (response?.error || 'Unknown error');
                return;
            }

            const { content, title } = response;
            msgDiv.textContent = 'Saving to Borg Core...';

            // Send to Background to save to Memory
            chrome.runtime.sendMessage({
                type: 'SAVE_CONTEXT',
                content: content,
                url: tab.url,
                title: title
            }, (res) => {
                if (res && res.success) {
                    msgDiv.textContent = '✅ Absorbed into Memory!';
                    setTimeout(() => msgDiv.textContent = '', 3000);
                } else {
                    msgDiv.textContent = 'Error Saving: ' + (res?.error || 'Core unreachable');
                }
            });
        });
    } catch (e: any) {
        msgDiv.textContent = 'Exception: ' + e.message;
    }
});
