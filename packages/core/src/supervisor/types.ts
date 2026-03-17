import type { LocalExecutionEnvironment } from '../services/execution-environment.js';
import type { SessionExecutionPolicy, SessionExecutionProfile } from '../services/session-execution-policy.js';

export type SessionLogStream = 'stdout' | 'stderr' | 'system';
export type SupervisedSessionStatus = 'created' | 'starting' | 'running' | 'stopping' | 'stopped' | 'restarting' | 'error';
export type SupervisedSessionHealthStatus = 'healthy' | 'degraded' | 'unresponsive' | 'crashed';

export interface MinimalReadableLike {
    on(event: 'data', listener: (chunk: Buffer | string) => void): unknown;
}

export interface SupervisedProcessHandle {
    pid?: number;
    stdout?: MinimalReadableLike;
    stderr?: MinimalReadableLike;
    on(event: 'exit', listener: (code: number | null, signal: NodeJS.Signals | null) => void): unknown;
    kill(signal?: NodeJS.Signals | number): boolean;
}

export interface SpawnProcessOptions {
    cwd: string;
    env: NodeJS.ProcessEnv;
}

export type SpawnSupervisedProcess = (
    command: string,
    args: string[],
    options: SpawnProcessOptions,
) => SupervisedProcessHandle;

export interface SchedulerLike {
    setTimeout(callback: () => void, delayMs: number): unknown;
    clearTimeout(handle: unknown): void;
}

/**
 * Contract for an optional worktree isolation backend wired to the {@link SessionSupervisor}.
 *
 * Runtime availability:
 * - `createTaskEnvironment` — Real implementation via `GitWorktreeManager` (uses `git worktree
 *   add` internally). Requires the repository to be a Git working tree and Git to be present on
 *   PATH. The `MCPServer` passes a real `GitWorktreeManager` instance whenever the supervisor is
 *   initialised through the normal boot path, so sessions created with `isolateWorktree: true`
 *   WILL get a separate Git worktree directory at runtime.
 *
 * - `cleanupTaskEnvironment` — Called by the supervisor on session stop/error. Uses `git worktree
 *   remove --force` internally. This is real but requires Git to be accessible; errors are logged
 *   and swallowed so a cleanup failure does NOT block the session lifecycle.
 *
 * If the supervisor is constructed without a `worktreeManager` (e.g. in unit tests or when the
 * runtime boot path skips the manager), `isolateWorktree: true` sessions log a warning and fall
 * back to the session's requested working directory without true isolation.
 */
export interface WorktreeManagerLike {
    createTaskEnvironment(taskId: string): Promise<string>;
    cleanupTaskEnvironment(taskId: string): Promise<void>;
}

export interface SessionSupervisorLogEntry {
    timestamp: number;
    stream: SessionLogStream;
    message: string;
}

export interface SessionAttachInfo {
    id: string;
    pid?: number;
    command: string;
    args: string[];
    cwd: string;
    status: SupervisedSessionStatus;
    attachable: boolean;
}

export interface SupervisedSessionHealth {
    status: SupervisedSessionHealthStatus;
    lastCheck: number;
    consecutiveFailures: number;
    restartCount: number;
    lastRestartAt?: number;
    nextRestartAt?: number;
    lastExitCode?: number;
    lastExitSignal?: string;
    errorMessage?: string;
}

export interface CreateSupervisedSessionInput {
    name?: string;
    cliType: string;
    workingDirectory: string;
    command?: string;
    args?: string[];
    env?: Record<string, string>;
    executionProfile?: SessionExecutionProfile;
    autoRestart?: boolean;
    isolateWorktree?: boolean;
    metadata?: Record<string, unknown>;
    maxRestartAttempts?: number;
}

export interface SupervisedSessionSnapshot {
    id: string;
    name: string;
    cliType: string;
    command: string;
    args: string[];
    env: Record<string, string>;
    executionProfile: SessionExecutionProfile;
    executionPolicy: SessionExecutionPolicy | null;
    requestedWorkingDirectory: string;
    workingDirectory: string;
    worktreePath?: string;
    autoRestart: boolean;
    isolateWorktree: boolean;
    status: SupervisedSessionStatus;
    createdAt: number;
    startedAt?: number;
    stoppedAt?: number;
    scheduledRestartAt?: number;
    lastActivityAt: number;
    restartCount: number;
    maxRestartAttempts: number;
    lastExitCode?: number;
    lastExitSignal?: string;
    lastError?: string;
    metadata: Record<string, unknown>;
    logs: SessionSupervisorLogEntry[];
}

export interface PersistedSessionSupervisorState {
    sessions: SupervisedSessionSnapshot[];
    savedAt: number;
}

export interface SessionSupervisorOptions {
    rootDir?: string;
    persistencePath?: string;
    maxPersistedSessions?: number;
    maxLogEntries?: number;
    autoResumeOnStart?: boolean;
    restartDelayMs?: number;
    backoffMultiplier?: number;
    maxBackoffMs?: number;
    maxRestartAttempts?: number;
    now?: () => number;
    spawnProcess?: SpawnSupervisedProcess;
    scheduler?: SchedulerLike;
    worktreeManager?: WorktreeManagerLike;
    detectExecutionEnvironment?: () => Promise<LocalExecutionEnvironment>;
}
