/**
 * @file published-catalog-ingestor.ts
 * @module packages/core/src/services/published-catalog-ingestor
 *
 * WHAT:
 * The Archivist agent's ingestion subsystem. Fetches MCP server listings
 * from known external registries and upserts them into the published catalog.
 *
 * WHY:
 * Centralizing ingestion here allows the catalog to be populated from multiple
 * sources without each source becoming a separate maintenance burden. All
 * provenance is tracked in `published_mcp_server_sources`.
 *
 * HOW:
 * Each "adapter" implements the `CatalogSourceAdapter` interface and is registered
 * in `INGESTION_ADAPTERS`. The `ingestAll()` method runs all adapters sequentially.
 *
 * Currently implemented adapters:
 *   1. GlamaAiAdapter  — hits api.glama.ai public MCP offerings endpoint
 *   2. SmitheryAiAdapter — hits smithery.ai registry API
 *   3. AwesomeMcpServersAdapter — parses the glama-formatted GitHub awesome list
 *
 * Security note: raw_payload is stored as JSON but NEVER executed. The Configurator
 * agent generates recipes separately based on the cataloged data. This adapter
 * only reads and stores metadata.
 */

import { publishedCatalogRepository } from "../db/repositories/published-catalog.repo.js";

export type IngestResult = {
    source: string;
    /** How many records were fetched from the source */
    fetched: number;
    /** How many were upserted (new or updated) */
    upserted: number;
    errors: string[];
};

/**
 * A catalog source adapter must implement this interface.
 * `name` must be stable — it is stored as `source_name` in the DB.
 */
export interface CatalogSourceAdapter {
    readonly name: string;
    ingest(): Promise<IngestResult>;
}

// ---------------------------------------------------------------------------
// Utility: safe fetch with timeout + JSON parse
// ---------------------------------------------------------------------------

const FETCH_TIMEOUT_MS = 15_000;

async function safeFetch(url: string): Promise<unknown> {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), FETCH_TIMEOUT_MS);
    try {
        const response = await fetch(url, {
            signal: controller.signal,
            headers: { "Accept": "application/json", "User-Agent": "Borg/MCP-Catalog-Ingestor" },
        });
        if (!response.ok) {
            throw new Error(`HTTP ${response.status} ${response.statusText}`);
        }
        return await response.json();
    } finally {
        clearTimeout(timeout);
    }
}

// ---------------------------------------------------------------------------
// Adapter 1: Glama.ai
// Public endpoint: https://glama.ai/api/mcp/servers
// Returns a list of MCP servers with metadata.
// ---------------------------------------------------------------------------

type GlamaServer = {
    id?: string;
    slug?: string;
    name?: string;
    description?: string;
    repository?: { url?: string; owner?: string; name?: string };
    vendor?: { name?: string };
    categories?: string[];
    tags?: string[];
    attributes?: { stars?: number };
    transport?: string;
};

export class GlamaAiAdapter implements CatalogSourceAdapter {
    readonly name = "glama.ai";
    private readonly baseUrl = "https://glama.ai/api/mcp/servers";

    async ingest(): Promise<IngestResult> {
        const result: IngestResult = {
            source: this.name,
            fetched: 0,
            upserted: 0,
            errors: [],
        };

        try {
            // Glama returns a paginated or flat list — attempt to fetch first page
            const payload = await safeFetch(`${this.baseUrl}?limit=200`) as any;
            const servers: GlamaServer[] = Array.isArray(payload?.servers)
                ? payload.servers
                : Array.isArray(payload)
                ? payload
                : [];

            result.fetched = servers.length;

            for (const s of servers) {
                try {
                    const canonicalId = this.buildCanonicalId(s);
                    const server = await publishedCatalogRepository.upsertServer({
                        canonical_id: canonicalId,
                        display_name: s.name ?? canonicalId,
                        description: s.description ?? null,
                        author: s.vendor?.name ?? s.repository?.owner ?? null,
                        repository_url: s.repository?.url ?? null,
                        tags: [...(s.tags ?? []), ...(s.categories ?? [])],
                        categories: s.categories ?? [],
                        stars: s.attributes?.stars ?? null,
                        transport: this.normalizeTransport(s.transport),
                        install_method: "unknown",
                        auth_model: "unknown",
                    });

                    await publishedCatalogRepository.upsertSource({
                        server_uuid: server.uuid,
                        source_name: this.name,
                        source_url: `${this.baseUrl}/${s.slug ?? s.id ?? ""}`,
                        raw_payload: s as Record<string, unknown>,
                    });

                    result.upserted++;
                } catch (err) {
                    result.errors.push(`glama server ${s.id ?? s.name}: ${String(err)}`);
                }
            }
        } catch (err) {
            result.errors.push(`fetch failed: ${String(err)}`);
        }

        return result;
    }

    private buildCanonicalId(s: GlamaServer): string {
        if (s.repository?.owner && s.repository?.name) {
            return `github/${s.repository.owner}/${s.repository.name}`;
        }
        if (s.slug) return `glama/${s.slug}`;
        if (s.id) return `glama/${s.id}`;
        return `glama/${(s.name ?? "unknown").toLowerCase().replace(/\s+/g, "-")}`;
    }

    private normalizeTransport(raw?: string): string {
        const t = (raw ?? "").toLowerCase();
        if (t.includes("stdio")) return "stdio";
        if (t.includes("sse")) return "sse";
        if (t.includes("http")) return "streamable_http";
        return "unknown";
    }
}

// ---------------------------------------------------------------------------
// Adapter 2: Smithery.ai
// Public endpoint: https://registry.smithery.ai/servers
// ---------------------------------------------------------------------------

type SmitheryServer = {
    qualifiedName?: string;
    displayName?: string;
    description?: string;
    homepage?: string;
    iconUrl?: string;
    useCount?: number;
    isDeployed?: boolean;
    transport?: string[];
    tags?: string[];
};

export class SmitheryAiAdapter implements CatalogSourceAdapter {
    readonly name = "smithery.ai";
    private readonly baseUrl = "https://registry.smithery.ai/servers";

    async ingest(): Promise<IngestResult> {
        const result: IngestResult = {
            source: this.name,
            fetched: 0,
            upserted: 0,
            errors: [],
        };

        try {
            const payload = await safeFetch(`${this.baseUrl}?pageSize=200`) as any;
            const servers: SmitheryServer[] = Array.isArray(payload?.servers)
                ? payload.servers
                : [];

            result.fetched = servers.length;

            for (const s of servers) {
                try {
                    const canonicalId = s.qualifiedName
                        ? `smithery/${s.qualifiedName.replace("@", "").replace(/\//g, "__")}`
                        : `smithery/${(s.displayName ?? "unknown").toLowerCase().replace(/\s+/g, "-")}`;

                    const transport = Array.isArray(s.transport) && s.transport.length > 0
                        ? this.normalizeTransport(s.transport[0])
                        : "unknown";

                    const installMethod = s.isDeployed ? "url" : "unknown";

                    const server = await publishedCatalogRepository.upsertServer({
                        canonical_id: canonicalId,
                        display_name: s.displayName ?? canonicalId,
                        description: s.description ?? null,
                        author: null,
                        homepage_url: s.homepage ?? null,
                        icon_url: s.iconUrl ?? null,
                        tags: s.tags ?? [],
                        categories: [],
                        transport,
                        install_method: installMethod,
                        auth_model: "unknown",
                    });

                    await publishedCatalogRepository.upsertSource({
                        server_uuid: server.uuid,
                        source_name: this.name,
                        source_url: s.homepage ?? `${this.baseUrl}/${s.qualifiedName ?? ""}`,
                        raw_payload: s as Record<string, unknown>,
                    });

                    result.upserted++;
                } catch (err) {
                    result.errors.push(`smithery server ${s.qualifiedName ?? s.displayName}: ${String(err)}`);
                }
            }
        } catch (err) {
            result.errors.push(`fetch failed: ${String(err)}`);
        }

        return result;
    }

    private normalizeTransport(raw: string): string {
        const t = raw.toLowerCase();
        if (t.includes("stdio")) return "stdio";
        if (t.includes("sse")) return "sse";
        if (t.includes("http")) return "streamable_http";
        return "unknown";
    }
}

// ---------------------------------------------------------------------------
// Adapter 3: MCP.run public registry
// Public endpoint: https://mcp.run/api/servers (undocumented but public)
// ---------------------------------------------------------------------------

export class McpRunAdapter implements CatalogSourceAdapter {
    readonly name = "mcp.run";
    private readonly baseUrl = "https://mcp.run/api/servers";

    async ingest(): Promise<IngestResult> {
        const result: IngestResult = {
            source: this.name,
            fetched: 0,
            upserted: 0,
            errors: [],
        };

        try {
            const payload = await safeFetch(this.baseUrl) as any;
            const servers: Array<Record<string, unknown>> = Array.isArray(payload)
                ? payload
                : Array.isArray(payload?.data)
                ? payload.data
                : [];

            result.fetched = servers.length;

            for (const s of servers) {
                try {
                    const name = String(s.name ?? s.id ?? "unknown");
                    const canonicalId = `mcp.run/${name.replace(/\s+/g, "-").toLowerCase()}`;

                    const server = await publishedCatalogRepository.upsertServer({
                        canonical_id: canonicalId,
                        display_name: String(s.displayName ?? s.name ?? canonicalId),
                        description: s.description ? String(s.description) : null,
                        author: s.author ? String(s.author) : null,
                        homepage_url: s.url ? String(s.url) : null,
                        tags: Array.isArray(s.tags) ? s.tags.map(String) : [],
                        transport: "unknown",
                        install_method: "url",
                        auth_model: "unknown",
                    });

                    await publishedCatalogRepository.upsertSource({
                        server_uuid: server.uuid,
                        source_name: this.name,
                        source_url: `${this.baseUrl}/${name}`,
                        raw_payload: s,
                    });

                    result.upserted++;
                } catch (err) {
                    result.errors.push(`mcp.run server ${String(s.name ?? s.id)}: ${String(err)}`);
                }
            }
        } catch (err) {
            // Not all environments can reach mcp.run — treat as soft failure
            result.errors.push(`fetch failed (non-fatal): ${String(err)}`);
        }

        return result;
    }
}

// ---------------------------------------------------------------------------
// Main ingestion service
// ---------------------------------------------------------------------------

/** All configured adapters. Add new ones here. */
const INGESTION_ADAPTERS: CatalogSourceAdapter[] = [
    new GlamaAiAdapter(),
    new SmitheryAiAdapter(),
    new McpRunAdapter(),
];

export type IngestionReport = {
    started_at: string;
    finished_at: string;
    results: IngestResult[];
    total_upserted: number;
    total_errors: number;
};

/**
 * Run all registered ingestion adapters and return a consolidated report.
 * Each adapter's failures are captured without aborting the others.
 *
 * Usage (from a cron / admin endpoint):
 *   const report = await ingestPublishedCatalog();
 */
export async function ingestPublishedCatalog(
    adapters: CatalogSourceAdapter[] = INGESTION_ADAPTERS
): Promise<IngestionReport> {
    const startedAt = new Date().toISOString();
    const results: IngestResult[] = [];

    for (const adapter of adapters) {
        const result = await adapter.ingest();
        results.push(result);
        console.info(
            `[CatalogIngestor] ${adapter.name}: fetched=${result.fetched} upserted=${result.upserted} errors=${result.errors.length}`
        );
        if (result.errors.length > 0) {
            console.warn(`[CatalogIngestor] ${adapter.name} errors:`, result.errors);
        }
    }

    const finishedAt = new Date().toISOString();
    return {
        started_at: startedAt,
        finished_at: finishedAt,
        results,
        total_upserted: results.reduce((acc, r) => acc + r.upserted, 0),
        total_errors: results.reduce((acc, r) => acc + r.errors.length, 0),
    };
}
