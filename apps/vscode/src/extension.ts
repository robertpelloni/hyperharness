import * as vscode from 'vscode';
import WebSocket from 'ws';

let socket: WebSocket | null = null;
let statusBarItem: vscode.StatusBarItem;
let outputChannel: vscode.OutputChannel;
let reconnectTimer: NodeJS.Timeout | null = null;
let lastActivityTime = Date.now();
let debounceTimer: NodeJS.Timeout | null = null;
let ignoreNextActivity = false; // Flag to prevent self-lockout

export function activate(context: vscode.ExtensionContext) {
    console.log('Borg VSCode Bridge is now active!');
    outputChannel = vscode.window.createOutputChannel('Borg Bridge');

    // Status Bar
    statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Right, 100);
    statusBarItem.command = 'borg.connect';
    updateStatusBar(false);
    statusBarItem.show();
    context.subscriptions.push(statusBarItem);

    // Commands
    context.subscriptions.push(vscode.commands.registerCommand('borg.connect', connectToCore));
    context.subscriptions.push(vscode.commands.registerCommand('borg.disconnect', disconnectFromCore));

    // Auto Connect
    const config = vscode.workspace.getConfiguration('borg');
    if (config.get('autoConnect', true)) {
        connectToCore();
    }

    // Track User Activity - Crucial for Anti-Hijack
    context.subscriptions.push(vscode.window.onDidChangeTextEditorSelection(() => {
        if (ignoreNextActivity) return;
        lastActivityTime = Date.now();
        if (debounceTimer) clearTimeout(debounceTimer);
        debounceTimer = setTimeout(sendActivity, 1000);
    }));

    context.subscriptions.push(vscode.workspace.onDidChangeTextDocument(() => {
        if (ignoreNextActivity) return;
        lastActivityTime = Date.now();
        if (debounceTimer) clearTimeout(debounceTimer);
        debounceTimer = setTimeout(sendActivity, 1000);
    }));
}

export function deactivate() {
    disconnectFromCore();
}

function sendActivity() {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({
            type: 'USER_ACTIVITY',
            lastActivityTime: Date.now()
        }));
    }
}

function updateStatusBar(connected: boolean) {
    if (connected) {
        statusBarItem.text = `$(plug) Borg: Connected`;
        statusBarItem.tooltip = `Connected to Borg Core`;
        statusBarItem.backgroundColor = undefined;
    } else {
        statusBarItem.text = `$(debug-disconnect) Borg: Disconnected`;
        statusBarItem.tooltip = 'Click to Connect';
        statusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.warningBackground');
    }
}

function log(message: string) {
    const timestamp = new Date().toISOString();
    outputChannel.appendLine(`[${timestamp}] ${message}`);
}

function connectToCore() {
    if (socket) return;

    const config = vscode.workspace.getConfiguration('borg');
    const url = config.get('coreUrl', 'ws://localhost:3001');

    log(`Connecting to ${url}...`);

    try {
        socket = new WebSocket(url);

        socket.on('open', () => {
            log('Connected to Borg Core');
            updateStatusBar(true);
            if (reconnectTimer) {
                clearTimeout(reconnectTimer);
                reconnectTimer = null;
            }
            vscode.window.showInformationMessage('Connected to Borg Core');
        });

        socket.on('message', async (data) => {
            try {
                const message = JSON.parse(data.toString());
                handleMessage(message);
            } catch (e) {
                log(`Failed to parse message: ${e}`);
            }
        });

        socket.on('close', () => {
            log('Disconnected from Borg Core');
            updateStatusBar(false);
            socket = null;
            // Auto reconnect
            reconnectTimer = setTimeout(connectToCore, 5000);
        });

        socket.on('error', (err) => {
            log(`Socket Error: ${err.message}`);
            socket?.close();
        });

    } catch (e: any) {
        log(`Connection setup error: ${e.message}`);
    }
}

function disconnectFromCore() {
    if (reconnectTimer) clearTimeout(reconnectTimer);
    if (socket) {
        socket.close();
        socket = null;
    }
    updateStatusBar(false);
}

async function handleMessage(msg: any) {
    if (msg.type === 'GET_USER_ACTIVITY') {
        if (socket) {
            socket.send(JSON.stringify({
                type: 'USER_ACTIVITY',
                requestId: msg.requestId,
                lastActivityTime,
                isIdle: (Date.now() - lastActivityTime) > 5000
            }));
        }
    }

    if (msg.type === 'VSCODE_COMMAND') {
        const { command, args } = msg;
        try {
            const result = await vscode.commands.executeCommand(command, ...(args || []));
            log(`Executed: ${command}`);
            sendResponse(msg.requestId, { success: true, result });
        } catch (e: any) {
            log(`Failed to execute ${command}: ${e.message}`);
            sendResponse(msg.requestId, { success: false, error: e.message });
        }
    }

    if (msg.type === 'PASTE_INTO_CHAT' || msg.type === 'SUBMIT_CHAT') {
        log(`[CHAT_ACTION] Received. type=${msg.type}`);
        try {
            // Anti-hijack: don't paste if user is typing
            if ((Date.now() - lastActivityTime) < 2000 && !ignoreNextActivity) {
                log(`[ABORT] User active within 2s, aborting auto-paste.`);
                return;
            }

            ignoreNextActivity = true;
            setTimeout(() => { ignoreNextActivity = false; }, 1500);

            if (msg.text) {
                await vscode.env.clipboard.writeText(msg.text);
                await vscode.commands.executeCommand('workbench.action.chat.open');
                await new Promise(r => setTimeout(r, 300));
                await vscode.commands.executeCommand('workbench.action.chat.focusInput');
                await new Promise(r => setTimeout(r, 200));
                await vscode.commands.executeCommand('editor.action.clipboardPasteAction');
            }

            if (msg.submit || msg.type === 'SUBMIT_CHAT') {
                await new Promise(r => setTimeout(r, 800));
                const commands = [
                    'workbench.action.chat.submit',
                    'workbench.action.chat.send',
                    'interactive.acceptChanges',
                    'workbench.action.terminal.chat.accept',
                    'inlineChat.accept'
                ];
                for (const cmd of commands) {
                    try { await vscode.commands.executeCommand(cmd); } catch (e) { }
                }
            }

        } catch (e: any) {
            log(`Failed chat action: ${e.message}`);
            ignoreNextActivity = false;
        }
    }

    if (msg.type === 'GET_STATUS') {
        const activeTerminal = vscode.window.activeTerminal;
        const status = {
            activeEditor: vscode.window.activeTextEditor?.document.fileName || null,
            activeTerminal: activeTerminal ? activeTerminal.name : null,
            workspace: vscode.workspace.workspaceFolders?.map((f: vscode.WorkspaceFolder) => f.uri.fsPath) || []
        };

        sendResponse(msg.requestId, { status });
    }

    if (msg.type === 'GET_SELECTION') {
        const editor = vscode.window.activeTextEditor;
        const text = editor?.document.getText(editor.selection) || "";
        sendResponse(msg.requestId, { content: text });
    }

    if (msg.type === 'GET_CHAT_HISTORY') {
        // [HEURISTIC] Since VSCode doesn't expose Chat API yet, we return a notice.
        // The Director will use this as a heartbeat.
        sendResponse(msg.requestId, { history: ["[System]: VS Code Extension Bridge Active. Monitoring session..."] });
    }

    if (msg.type === 'GET_TERMINAL') {
        const terminal = vscode.window.activeTerminal;
        if (terminal) {
            // No easy way to read terminal content directly from Extension API without creating a output file
            // But we can report the terminal name at least.
            sendResponse(msg.requestId, { content: `Terminal: ${terminal.name} (Active)` });
        } else {
            sendResponse(msg.requestId, { content: "No active terminal." });
        }
    }
}

function sendResponse(requestId: string, payload: any) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({
            type: 'RESPONSE',
            requestId,
            ...payload
        }));
    }
}
