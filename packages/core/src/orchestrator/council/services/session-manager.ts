import type { Session, LogEntry, DevelopmentTask, Guidance, CLIType, SessionHealth, CLITool } from './types.js';
import path from 'path';
import { EventEmitter } from 'events';
import { cliRegistry } from './cli-registry.js';
import { wsManager } from './ws-manager.js';
import { loadConfig } from './config.js';
import { dbService } from './db.js';
import { getMcpServer } from '../../../lib/trpc-core.js';

export interface SessionOptions {
  tags?: string[];
  templateName?: string;
  workingDirectory?: string;
  cliType?: CLIType;
  env?: Record<string, string>;
}

export interface BulkSessionRequest {
  count: number;
  templateName?: string;
  tags?: string[];
  cliType?: CLIType;
  staggerDelayMs?: number;
}

export interface BulkSessionResponse {
  sessions: Session[];
  failed: string[];
}

class SessionManager extends EventEmitter {
  private sessions = new Map<string, Session>();
  private readonly config = loadConfig();

  constructor() {
    super();
  }

  /**
   * Start a new supervised session
   */
  async startSession(task?: DevelopmentTask, options: SessionOptions = {}): Promise<Session> {
    const id = task?.id || `session_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    const cliType = options.cliType || task?.cliType || 'claude-code';
    const tool = await cliRegistry.getTool(cliType);
    
    if (!tool) {
      throw new Error(`CLI tool '${cliType}' not found or not supported.`);
    }

    const workingDir = options.workingDirectory || process.cwd();
    const env = { ...process.env, ...this.config.env, ...options.env };

    const session: Session = {
      id,
      task: task || { id, description: 'Interactive Session', context: '', files: [] },
      status: 'starting',
      startTime: Date.now(),
      logs: [],
      port: 0, // In Borg integration, we might not use ports for all sessions
      tags: options.tags || [],
      cliType,
    };

    this.sessions.set(id, session);
    
    // Use Borg's PtySupervisor
    const mcpServer = getMcpServer();
    const ptySupervisor = (mcpServer as any).ptySupervisor;

    if (!ptySupervisor) {
      throw new Error('PtySupervisor not initialized in MCPServer');
    }

    try {
      const supervisorSession = await ptySupervisor.createSession({
        command: tool.command,
        args: tool.serveArgs,
        cwd: workingDir,
        env,
        metadata: {
          borgSessionId: id,
          cliType,
        }
      });

      session.status = 'active';
      
      // Hook into logs
      // Note: SessionSupervisor typically collects logs itself. 
      // We can subscribe to them if needed or just query via tRPC.
      // For parity, we'll keep a local log buffer.
      
      this.emit('session_started', session);
      wsManager.notifySessionStarted(session);

      return session;
    } catch (error) {
      session.status = 'error';
      const msg = error instanceof Error ? error.message : 'Unknown error starting session';
      this.emit('session_error', { id, error: msg });
      throw error;
    }
  }

  // ... (Rest of the session manager methods adapted to use ptySupervisor)
  
  async stopSession(id: string): Promise<void> {
    const session = this.sessions.get(id);
    if (!session) return;

    const mcpServer = getMcpServer();
    const ptySupervisor = (mcpServer as any).ptySupervisor;
    
    if (ptySupervisor) {
      // Find the internal supervisor ID from metadata or assume it matches if we mapped it
      // For now, let's assume we can kill it.
      await ptySupervisor.stopSession(id);
    }

    session.status = 'stopped';
    this.emit('session_stopped', id);
    wsManager.notifySessionStopped(id);
  }

  async stopAllSessions(): Promise<{ stopped: string[], failed: string[] }> {
    const stopped: string[] = [];
    const failed: string[] = [];

    for (const id of this.sessions.keys()) {
      try {
        await this.stopSession(id);
        stopped.push(id);
      } catch (err) {
        failed.push(id);
      }
    }

    return { stopped, failed };
  }

  getSession(id: string): Session | undefined {
    return this.sessions.get(id);
  }

  getAllSessions(): Session[] {
    return Array.from(this.sessions.values());
  }

  getActiveSessions(): Session[] {
    return this.getAllSessions().filter(s => s.status === 'active');
  }

  // Simplified for parity
  async getSessionStats() {
    return {
      total: this.sessions.size,
      active: this.getActiveSessions().length,
    };
  }
}

export const sessionManager = new SessionManager();
