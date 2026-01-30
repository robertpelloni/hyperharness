// Background Service Worker for Borg Director Link
console.log("aios Browser Extension Background Service Worker Starting...");

let socket: WebSocket | null = null;
let requestIdCounter = 0;
const pendingRequests = new Map<number, (response: any) => void>();

function connect() {
  console.log("Connecting to Borg MCP Server (ws://localhost:3001)...");
  socket = new WebSocket("ws://localhost:3001");

  socket.onopen = () => {
    console.log("✅ Connected to Borg Hub");
    // Initialize MCP Session? usually happens automatically or lazily
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

      // Handle Server Requests (Reverse Control - e.g. Server calling Browser Tools)
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
  // Basic Reverse Control (Server -> Browser)
  // Matches the json-rpc 2.0 format: { jsonrpc: "2.0", method: "...", params: {...}, id: ... }
  if (msg.method === "ping") {
    send({ jsonrpc: "2.0", result: "pong", id: msg.id });
    return;
  }

  // Commands implementation (similar to previous version)
  try {
    let result;
    if (msg.method === "browser_navigate") {
      result = await navigate(msg.params.url);
    } else if (msg.method === "browser_evaluate") {
      result = await getActiveTabContent();
    } else {
      // If unknown method, maybe it's a notification or we just ignore
      console.warn("Unknown method:", msg.method);
      return; // Or send error
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

// Forwarding from Content Script (Web Page -> Background -> Server)
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.type === "MCP_REQUEST") {
    // request.payload should be a JSON-RPC request object
    // We assume the Web Page constructs a valid JSON-RPC object
    // But we need to manage the ID mapping to reply to the correct tab/promise

    if (!socket || socket.readyState !== WebSocket.OPEN) {
      sendResponse({ error: "Borg Disconnected" });
      return true;
    }

    const internalId = ++requestIdCounter;
    const originalId = request.payload.id;

    // Wrap/Map ID?
    // Actually, we can just pass it through IF we trust the session is 1:1
    // But to be safe and handle the response back to `sendResponse`, we should track it.
    // However, `sendResponse` is for the one-off Chrome message.
    // We'll use pendingRequests to map back.

    const outboundMsg = { ...request.payload, id: internalId };

    pendingRequests.set(internalId, (responseMsg) => {
      // Send back to content script? 
      // Chrome message passing is async. 
      // We need to keep `sendResponse` valid by returning true.

      // Restore original ID for the client
      const clientResponse = { ...responseMsg, id: originalId };
      sendResponse(clientResponse);
    });

    send(outboundMsg);
    return true; // Keep channel open
  }

  if (request.action === "ping") return sendResponse("pong");
});

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
