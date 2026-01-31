// Background Service Worker for Borg Director Link
console.log("Borg Browser Extension Background Service Worker Starting... 🚀");

let socket: WebSocket | null = null;
let requestIdCounter = 0;
const pendingRequests = new Map<number, (response: any) => void>();

function connect() {
  console.log("Connecting to Borg MCP Server (ws://localhost:3001)...");
  socket = new WebSocket("ws://localhost:3001");

  socket.onopen = () => {
    console.log("✅ Connected to Borg Hub");
  };

  socket.onmessage = async (event) => {
    try {
      const msg = JSON.parse(event.data);
      console.log("Received:", msg);

      // Handle Responses to our requests
      if (msg.id !== undefined && pendingRequests.has(msg.id)) {
        const resolver = pendingRequests.get(msg.id);
        resolver && resolver(msg);
        pendingRequests.delete(msg.id);
        return;
      }

      // Handle Server Requests (Reverse Control)
      if (msg.method) {
        await handleServerRequest(msg);
      }

    } catch (e) {
      console.error("Error parsing message:", e);
    }
  };

  socket.onclose = () => {
    console.log("❌ Disconnected. Reconnecting in 5s...");
    socket = null;
    setTimeout(connect, 5000);
  };

  socket.onerror = (err) => {
    console.error("WebSocket Error:", err);
  };
}

async function handleServerRequest(msg: any) {
  // JSON-RPC processing
  if (msg.method === "ping") {
    send({ jsonrpc: "2.0", result: "pong", id: msg.id });
    return;
  }

  try {
    let result;
    if (msg.method === "browser_navigate") {
      result = await navigate(msg.params.url);
    } else if (msg.method === "browser_evaluate") {
      result = await getActiveTabContent();
    } else if (msg.method === "browser_scrape") { // <--- NEW: Scrape Handler
      result = await scrapeActiveTab();
    } else {
      console.warn("Unknown method:", msg.method);
      // Don't error, just ignore for forward compatibility
      return;
    }

    if (msg.id !== undefined) {
      send({ jsonrpc: "2.0", result, id: msg.id });
    }
  } catch (e: any) {
    if (msg.id !== undefined) {
      send({ jsonrpc: "2.0", error: { code: -32000, message: e.message }, id: msg.id });
    }
  }
}

function send(msg: any) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(msg));
  } else {
    console.warn("Socket not open. Message dropped:", msg);
  }
}

// Scrape Logic
async function scrapeActiveTab(): Promise<any> {
  const [tab] = await chrome.tabs.query({ active: true, lastFocusedWindow: true });
  if (!tab?.id) throw new Error("No active tab found");

  // Send message to content script
  return chrome.tabs.sendMessage(tab.id, { action: 'SCRAPE_PAGE' });
}

// Browser Automation Helpers
async function getActiveTabContent() {
  const [tab] = await chrome.tabs.query({ active: true, lastFocusedWindow: true });
  if (!tab || !tab.id) throw new Error("No active tab found");

  const result = await chrome.scripting.executeScript({
    target: { tabId: tab.id },
    func: () => document.body.innerText,
  });

  return result[0].result;
}

async function navigate(url: string) {
  const [tab] = await chrome.tabs.query({ active: true, lastFocusedWindow: true });
  if (tab && tab.id) {
    await chrome.tabs.update(tab.id, { url });
    return `Navigated to ${url}`;
  } else {
    await chrome.tabs.create({ url });
    return `Opened ${url}`;
  }
}

// Start connection
connect();
