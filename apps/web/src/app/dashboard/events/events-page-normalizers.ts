export interface DashboardEventRow {
    type: string;
    timestamp: number | null;
    source: string;
    dataPreview: string | null;
}

export interface DashboardSystemStatus {
    status: 'online' | 'offline';
    uptime: number;
    agents: string[];
}

function isObject(value: unknown): value is Record<string, unknown> {
    return typeof value === 'object' && value !== null;
}

function safePreview(value: unknown): string | null {
    if (typeof value === 'string') {
        return value;
    }

    if (value === undefined || value === null) {
        return null;
    }

    try {
        return JSON.stringify(value);
    } catch {
        return '[unserializable data]';
    }
}

function normalizeTimestamp(value: unknown): number | null {
    if (typeof value === 'number' && Number.isFinite(value)) {
        return value;
    }

    if (typeof value === 'string' && value.trim().length > 0) {
        const asNumber = Number(value);
        if (Number.isFinite(asNumber)) {
            return asNumber;
        }

        const asDate = Date.parse(value);
        return Number.isFinite(asDate) ? asDate : null;
    }

    return null;
}

export function normalizeDashboardEvents(payload: unknown): DashboardEventRow[] {
    if (!Array.isArray(payload)) {
        return [];
    }

    return payload.reduce<DashboardEventRow[]>((acc, item, index) => {
        if (!isObject(item)) {
            return acc;
        }

        const rawType = typeof item.type === 'string' ? item.type.trim() : '';
        const rawSource = typeof item.source === 'string' ? item.source.trim() : '';

        acc.push({
            type: rawType.length > 0 ? rawType : `unknown:${index}`,
            timestamp: normalizeTimestamp(item.timestamp),
            source: rawSource,
            dataPreview: safePreview(item.data),
        });

        return acc;
    }, []);
}

export function normalizeDashboardSystemStatus(payload: unknown): DashboardSystemStatus {
    if (!isObject(payload)) {
        return {
            status: 'offline',
            uptime: 0,
            agents: [],
        };
    }

    const status = payload.status === 'online' ? 'online' : 'offline';
    const uptime = typeof payload.uptime === 'number' && Number.isFinite(payload.uptime) ? payload.uptime : 0;
    const agents = Array.isArray(payload.agents)
        ? payload.agents.filter((agent): agent is string => typeof agent === 'string' && agent.trim().length > 0).map((agent) => agent.trim())
        : [];

    return {
        status,
        uptime,
        agents,
    };
}