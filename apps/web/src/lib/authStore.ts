import { randomUUID, scryptSync, timingSafeEqual } from 'crypto';
import { promises as fs } from 'fs';
import path from 'path';

type AuthUser = {
    id: string;
    name: string;
    email: string;
    passwordHash: string;
    createdAt: string;
    resetRequestedAt?: string;
};

type AuthDb = {
    users: AuthUser[];
};

const STORE_DIR = path.join(process.cwd(), '.borg-auth');
const STORE_FILE = path.join(STORE_DIR, 'users.json');
const DEFAULT_DB: AuthDb = { users: [] };

function normalizeEmail(email: string): string {
    return email.trim().toLowerCase();
}

function hashPassword(password: string): string {
    // Lightweight local hash for dev auth flow wiring.
    return scryptSync(password, 'borg-local-salt', 64).toString('hex');
}

function verifyPassword(password: string, passwordHash: string): boolean {
    const computed = Buffer.from(hashPassword(password), 'hex');
    const expected = Buffer.from(passwordHash, 'hex');
    if (computed.length !== expected.length) {
        return false;
    }
    return timingSafeEqual(computed, expected);
}

async function ensureStore(): Promise<void> {
    await fs.mkdir(STORE_DIR, { recursive: true });
    try {
        await fs.access(STORE_FILE);
    } catch {
        await fs.writeFile(STORE_FILE, JSON.stringify(DEFAULT_DB, null, 2), 'utf-8');
    }
}

async function readDb(): Promise<AuthDb> {
    await ensureStore();
    const raw = await fs.readFile(STORE_FILE, 'utf-8');
    try {
        const parsed = JSON.parse(raw) as AuthDb;
        return { users: parsed.users ?? [] };
    } catch {
        return DEFAULT_DB;
    }
}

async function writeDb(db: AuthDb): Promise<void> {
    await ensureStore();
    await fs.writeFile(STORE_FILE, JSON.stringify(db, null, 2), 'utf-8');
}

export async function createUser(input: { name: string; email: string; password: string }) {
    const db = await readDb();
    const normalized = normalizeEmail(input.email);
    const existing = db.users.find((u) => u.email === normalized);
    if (existing) {
        return { ok: false as const, reason: 'EXISTS' as const };
    }

    const user: AuthUser = {
        id: randomUUID(),
        name: input.name.trim(),
        email: normalized,
        passwordHash: hashPassword(input.password),
        createdAt: new Date().toISOString(),
    };

    db.users.push(user);
    await writeDb(db);

    return {
        ok: true as const,
        user: {
            id: user.id,
            name: user.name,
            email: user.email,
        },
    };
}

export async function authenticateUser(input: { email: string; password: string }) {
    const db = await readDb();
    const normalized = normalizeEmail(input.email);
    const user = db.users.find((u) => u.email === normalized);
    if (!user) {
        return { ok: false as const, reason: 'INVALID_CREDENTIALS' as const };
    }
    if (!verifyPassword(input.password, user.passwordHash)) {
        return { ok: false as const, reason: 'INVALID_CREDENTIALS' as const };
    }

    return {
        ok: true as const,
        user: {
            id: user.id,
            name: user.name,
            email: user.email,
        },
    };
}

export async function markResetRequested(email: string) {
    const db = await readDb();
    const normalized = normalizeEmail(email);
    const user = db.users.find((u) => u.email === normalized);
    if (user) {
        user.resetRequestedAt = new Date().toISOString();
        await writeDb(db);
    }
    // Always return success shape to avoid account enumeration.
    return { ok: true as const };
}
