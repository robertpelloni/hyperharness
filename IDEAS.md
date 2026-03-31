# Borg IDEAS & Future Architectural Pivots

*This document captures ambitious, far-reaching feature ideas, refactoring opportunities, and conceptual pivots to ensure Borg remains the bleeding-edge pinnacle of autonomous IDE infrastructure.*

## 1. Decentralized / P2P Memory Swarm (The "Borg Collective")
We have standard local LanceDB and SQLite, but if someone is running Borg on their desktop and laptop, they should share a P2P context mesh. We could integrate libp2p or a lightweight WebRTC signaling layer so that memory subagents on one machine can securely broadcast their vector embeddings to authorized peers. 
*Impact:* True "omniscient memory" across all of a user's devices instantly.

## 2. Kernel-Level File System Watcher (eBPF/FSEvents)
Instead of relying strictly on IDE file-save hooks, Borg could deploy an eBPF (Linux) or native FSEvents (macOS/Windows) daemon. This daemon would monitor exactly which files are being compiled, opened, or crashed, and proactively feed those backtraces into the Roundtable Council *before* the developer even asks "why did it crash?"
*Impact:* Zero-friction debugging. The AI knows the error before you open the dashboard.

## 3. Sandboxed Wasm/Deno Execution for MCP Tools
Currently, we flag `stdio` MCP tools as "unsafe" during catalog ingestion. We could embed a Rust-based Wasmtime runtime or Deno sandbox inside of `@borg/core`. When an MCP server is fetched, we compile it to Wasm or run it via Deno with strict network/file capabilities.
*Impact:* We can safely install and test any MCP server from the public internet without risk of remote code execution escaping the vault.

## 4. Visual "Node-Graph" Workflow Editor
The dashboard currently relies on data tables and standard Next.js forms. We could integrate React Flow to create a visual node-based editor for chaining `SmartPilot` subagents, MCP tools, and memory filters. 
*Impact:* Empowers standard operators to build complex auto-orchestration chains without writing code.

## 5. Voice-Operated Local LLM (Whisper + Llama Integration)
Integrate an offline WebGPU-accelerated Whisper model in the web browser. The operator can literally just speak to the system ("Borg, fix the TS errors in the UI package") and the system transcribes it locally, passing it directly into the execution substrate.
*Impact:* Hands-free 10x coding operations.

## 6. Language Porting / Refactoring
*   **Rust rewrite of `core` event loop:** Port the PTY supervisor and session multiplexer from Node.js to Rust. Expose it back to Node via N-API (Neon/napi-rs). Node.js is excellent for TRPC and tying tools together, but heavy interactive PTY shell streams might suffer under the V8 event loop under heavy load.
*   **React Native / Expo Mobile App:** A dedicated iOS/Android control pane using React Native that connects to the local Borg server via generic WebSocket tunnels (e.g., ngrok or Cloudflare Tunnels), allowing the operator to monitor multi-agent debates while away from the keyboard.
