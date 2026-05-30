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

export interface NormalizedDashboardEventsResult {
    data: DashboardEventRow[];
    invalid: boolean;
}

export interface NormalizedDashboardSystemStatusResult {
    data: DashboardSystemStatus;
    invalid: boolean;
}

            return acc;
        }

        const rawType = typeof item.type === 'string' ? item.type.trim() : '';
        if (rawType.length === 0) {
            invalid = true;
        }
        };
    }

    const status = payload.status === 'online' ? 'online' : 'offline';
    const uptime = typeof payload.uptime === 'number' && Number.isFinite(payload.uptime) ? payload.uptime : 0;
    const agents = Array.isArray(payload.agents)
        ? payload.agents.filter((agent): agent is string => typeof agent === 'string' && agent.trim().length > 0).map((agent) => agent.trim())
        : [];
    const invalid = (payload.status !== 'online' && payload.status !== 'offline')
        || !(typeof payload.uptime === 'number' && Number.isFinite(payload.uptime))
        || !Array.isArray(payload.agents);

    return {
        data: {
            status,
            uptime,
            agents,
        },
        invalid,
    };
}
