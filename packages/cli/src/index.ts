#!/usr/bin/env node
import './logger-patch.js';

console.log("[CLI] Bare Metal Startup..."); // Debug log

// Emulate 'start' command being default
const args = process.argv.slice(2);
const command = args[0] || 'start';

if (command === 'start') {
  (async () => {
    try {
      // CLI MARKER
      const fs = await import('fs');
      fs.writeFileSync('.cli_startup_marker', 'CLI Started at ' + new Date().toISOString());

      const React = (await import('react')).default;
      const { render } = await import('ink');
      // @ts-ignore
      const { App } = await import('./ui/App.js');

      console.log = function () { }; // Final silence before TUI take over, though App takes over pretty quick.
      render(React.createElement(App));
    } catch (e) {
      // original console error
      process.stderr.write(`[CLI] FATAL: ${e}\n`);
      process.exit(1);
    }
  })();
} else if (command === 'status') {
  // Lazy Load UI for status
  console.log("[CLI] Loading UI for status...");
  (async () => {
    const React = (await import('react')).default;
    const { render } = await import('ink');
    const { App } = await import('./ui/App.js');
    render(React.createElement(App, { view: 'status' }));
  })();
} else if (command === 'research') {
  (async () => {
    console.log("[CLI] Starting Research Pipeline...");
    const { MCPServer } = await import('@borg/core/MCPServer');
    // @ts-ignore
    const { runResearchPipeline } = await import('@borg/core/dist/scripts/researchPipeline.js');
    const server = new MCPServer();
    await runResearchPipeline(server);
  })();
} else {
  console.log("Unknown command. Usage: borg [start|status]");
  process.exit(1);
}
