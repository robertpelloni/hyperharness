
function redirectProtocolUnsafeConsoleMethods(): void {
    const stderr = console.error.bind(console);

    console.log = stderr;
    console.info = stderr;
    console.debug = stderr;
    console.trace = stderr;
    console.dir = ((...args: unknown[]) => stderr(...args)) as typeof console.dir;
}

async function main() {
    // MCP stdio requires stdout to remain pristine JSON-RPC output only.
    redirectProtocolUnsafeConsoleMethods();

    process.on('unhandledRejection', (reason) => {
        console.error('[borg Core] Unhandled promise rejection:', reason);
    });

    process.on('uncaughtException', (error) => {
        console.error('[borg Core] Uncaught exception:', error);
        process.exit(1);
    });

    try {
        const { startStdioLoader } = await import('./stdioLoader.js');
        await startStdioLoader();
        console.error('[borg Core] MCP stdio loader started.');
    } catch (err) {
        console.error("Failed to start MCP server:", err);
        process.exit(1);
    }
}

main();
