"use client";

/**
 * @file page.tsx
 * @route /dashboard/registry/[uuid]
 *
 * WHAT:
 * Server detail view for a single published MCP catalog entry.
 * Shows full metadata, validation history, active recipe, and provenance sources.
 *
 * WHY:
 * The registry table gives operators a bird's-eye view, but they need a
 * drill-down to make informed install/trust decisions. This page is the
 * "truth page" for a single published server — no guessing, no placeholders.
 *
 * SECTIONS:
 *   1. Header — name, status badge, transport, install method, confidence bar
 *   2. Metadata — description, author, repository, homepage, tags, categories
 *   3. Validation history — timestamp, outcome, tool_count, failure_class
 *   4. Active recipe — template template, required_secrets, required_env, confidence
 *   5. Provenance sources — which registries listed this server and when
 *   6. Actions — Validate, Install (if validated/certified)
 */

import { trpc } from "@/utils/trpc";
import { useParams, useRouter } from "next/navigation";
import {
    ArrowLeft,
    CheckCircle2,
    AlertCircle,
    Clock,
    Database,
    Package,
    Globe,
    Cpu,
    ExternalLink,
    PlayCircle,
    Download,
    Star,
    Tag,
    Link,
    ShieldCheck,
    FileCode,
    ListChecks,
    RefreshCw,
} from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import type { inferRouterOutputs } from "@trpc/server";
import type { AppRouter } from "@borg/core";

// ---- Types ----

type RouterOutput = inferRouterOutputs<AppRouter>;
type ServerDetail = RouterOutput["catalog"]["get"];
type ValidationRun = RouterOutput["catalog"]["listRuns"][number];

type ServerStatus =
    | "discovered"
    | "normalized"
    | "probeable"
    | "validated"
    | "certified"
    | "broken"
    | "archived";

// ---- Status config ----

const STATUS_CONFIG: Record<
    ServerStatus,
    { label: string; color: string; icon: React.ComponentType<{ className?: string }> }
> = {
    discovered: { label: "Discovered", color: "bg-zinc-800 text-zinc-400 border-zinc-700", icon: Clock },
    normalized: { label: "Normalized", color: "bg-blue-950 text-blue-400 border-blue-800", icon: Database },
    probeable: { label: "Probeable", color: "bg-amber-950 text-amber-400 border-amber-800", icon: Package },
    validated: { label: "Validated", color: "bg-emerald-950 text-emerald-400 border-emerald-800", icon: CheckCircle2 },
    certified: { label: "Certified", color: "bg-teal-950 text-teal-400 border-teal-800", icon: CheckCircle2 },
    broken: { label: "Broken", color: "bg-red-950 text-red-400 border-red-800", icon: AlertCircle },
    archived: { label: "Archived", color: "bg-zinc-900 text-zinc-600 border-zinc-800", icon: Clock },
};

function getStatusConfig(status: string) {
    return STATUS_CONFIG[status as ServerStatus] ?? STATUS_CONFIG.discovered;
}

// ---- Outcome config ----

const OUTCOME_COLOR: Record<string, string> = {
    passed: "text-emerald-400",
    failed: "text-red-400",
    error: "text-red-400",
    timeout: "text-amber-400",
    skipped: "text-zinc-500",
    pending: "text-blue-400",
};

function outcomeIcon(outcome: string) {
    if (outcome === "passed") return <CheckCircle2 className="w-3.5 h-3.5 text-emerald-400" />;
    if (outcome === "failed" || outcome === "error") return <AlertCircle className="w-3.5 h-3.5 text-red-400" />;
    if (outcome === "skipped") return <Clock className="w-3.5 h-3.5 text-zinc-500" />;
    return <Clock className="w-3.5 h-3.5 text-blue-400" />;
}

function formatDate(value: string | Date | null | undefined): string {
    if (!value) return "—";
    const d = typeof value === "string" ? new Date(value) : value;
    if (isNaN(d.getTime())) return "—";
    return d.toLocaleString();
}

// ---- Main page ----

export default function ServerDetailPage() {
    const params = useParams();
    const router = useRouter();
    const uuid = typeof params?.uuid === "string" ? params.uuid : "";

    const [installName, setInstallName] = useState("");
    const [installEnv, setInstallEnv] = useState<Record<string, string>>({});
    const [showInstallModal, setShowInstallModal] = useState(false);
    const [showSecrets, setShowSecrets] = useState(false);

    const utils = trpc.useContext();

    // Server detail query
    const { data: detail, isLoading } = trpc.catalog.get.useQuery(
        { uuid },
        { enabled: Boolean(uuid), staleTime: 30_000 }
    );

    // Validation run history (up to 20)
    const { data: runs = [], isFetching: loadingRuns } = trpc.catalog.listRuns.useQuery(
        { server_uuid: uuid, limit: 20 },
        { enabled: Boolean(uuid), staleTime: 30_000 }
    );

    // Validate mutation
    const validateMutation = trpc.catalog.triggerValidation.useMutation({
        onSuccess: (result) => {
            const icon = result.outcome === "passed" ? "✅" : result.outcome === "skipped" ? "⏭️" : "❌";
            toast.success(
                `${icon} Validation ${result.outcome}${result.tool_count != null ? ` — ${result.tool_count} tools` : ""}`
            );
            utils.catalog.get.invalidate({ uuid });
            utils.catalog.listRuns.invalidate({ server_uuid: uuid });
        },
        onError: (err) => {
            toast.error("Validation error: " + err.message);
        },
    });

    // Install mutation
    const installMutation = trpc.catalog.installFromRecipe.useMutation({
        onSuccess: (result) => {
            toast.success(`Installed as "${result.name}"`);
            utils.catalog.get.invalidate({ uuid });
            utils.mcpServers.list.invalidate();
            setShowInstallModal(false);
            setInstallName("");
            setInstallEnv({});
            setShowSecrets(false);
        },
        onError: (err) => {
            toast.error("Install failed: " + err.message);
        },
    });

    if (!uuid) {
        return (
            <div className="p-8 text-center text-zinc-500">Invalid server UUID.</div>
        );
    }

    if (isLoading) {
        return (
            <div className="p-8 flex items-center justify-center gap-3 text-zinc-500">
                <RefreshCw className="w-5 h-5 animate-spin" /> Loading…
            </div>
        );
    }

    if (!detail?.server) {
        return (
            <div className="p-8 text-center space-y-4">
                <AlertCircle className="w-10 h-10 mx-auto text-red-500 opacity-60" />
                <p className="text-zinc-400">Server not found in catalog.</p>
                <button
                    onClick={() => router.back()}
                    className="text-sm text-indigo-400 hover:underline"
                >
                    ← Back to registry
                </button>
            </div>
        );
    }

    const { server, latestRun, activeRecipe, sources } = detail;
    const statusCfg = getStatusConfig(server.status);
    const StatusIcon = statusCfg.icon;
    const canInstall = server.status === "validated" || server.status === "certified";

    const requiredSecrets = activeRecipe?.required_secrets ?? [];
    const requiredEnv = activeRecipe?.required_env ?? {};
    const missingSecrets = requiredSecrets.filter((k) => !(installEnv[k] ?? "").trim());
    const canSubmit =
        !installMutation.isPending &&
        Boolean(activeRecipe) &&
        missingSecrets.length === 0;

    return (
        <div className="p-6 max-w-5xl mx-auto space-y-6">
            {/* Back + breadcrumb */}
            <div className="flex items-center gap-2 text-sm text-zinc-500">
                <button
                    onClick={() => router.push("/dashboard/registry")}
                    className="flex items-center gap-1.5 hover:text-zinc-300 transition-colors"
                >
                    <ArrowLeft className="w-4 h-4" />
                    MCP Registry
                </button>
                <span>/</span>
                <span className="text-zinc-300 truncate max-w-xs">{server.display_name}</span>
            </div>

            {/* Header card */}
            <div className="bg-zinc-900/40 border border-zinc-800 rounded-xl p-5">
                <div className="flex flex-wrap items-start justify-between gap-4">
                    <div className="flex-1 min-w-0">
                        <h1 className="text-xl font-bold text-zinc-100 truncate">{server.display_name}</h1>
                        {server.author && (
                            <p className="text-sm text-zinc-500 mt-0.5">by {server.author}</p>
                        )}
                        {server.description && (
                            <p className="text-zinc-400 text-sm mt-2 leading-relaxed">{server.description}</p>
                        )}
                    </div>

                    {/* Actions */}
                    <div className="flex items-center gap-2 flex-shrink-0">
                        {server.repository_url && (
                            <a
                                href={server.repository_url}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="flex items-center gap-1.5 px-3 py-1.5 text-xs rounded-lg border border-zinc-700 text-zinc-400 hover:text-zinc-200 hover:border-zinc-600 transition-colors"
                            >
                                <ExternalLink className="w-3.5 h-3.5" />
                                Repository
                            </a>
                        )}
                        <button
                            onClick={() => validateMutation.mutate({ server_uuid: uuid })}
                            disabled={validateMutation.isPending}
                            className="flex items-center gap-1.5 px-3 py-1.5 text-xs rounded-lg border border-zinc-700 text-zinc-400 hover:text-indigo-400 hover:border-indigo-700 disabled:opacity-40 transition-colors"
                        >
                            <PlayCircle className={`w-3.5 h-3.5 ${validateMutation.isPending ? "animate-pulse" : ""}`} />
                            {validateMutation.isPending ? "Validating…" : "Validate"}
                        </button>
                        {canInstall && (
                            <button
                                onClick={() => { setShowInstallModal(true); setInstallName(""); setInstallEnv({}); setShowSecrets(false); }}
                                className="flex items-center gap-1.5 px-3 py-1.5 text-xs rounded-lg bg-emerald-700 hover:bg-emerald-600 text-white transition-colors"
                            >
                                <Download className="w-3.5 h-3.5" />
                                Install
                            </button>
                        )}
                    </div>
                </div>

                {/* Badges row */}
                <div className="flex flex-wrap gap-2 mt-4">
                    <span className={`inline-flex items-center gap-1.5 px-2 py-1 rounded-full text-xs font-medium border ${statusCfg.color}`}>
                        <StatusIcon className="w-3 h-3" />
                        {statusCfg.label}
                    </span>
                    <span className="inline-flex items-center gap-1 px-2 py-1 rounded text-xs bg-zinc-800 text-zinc-400 border border-zinc-700">
                        <Globe className="w-3 h-3" />
                        {server.transport}
                    </span>
                    {server.install_method && server.install_method !== "unknown" && (
                        <span className="inline-flex items-center gap-1 px-2 py-1 rounded text-xs bg-zinc-800 text-zinc-400 border border-zinc-700">
                            <Cpu className="w-3 h-3" />
                            {server.install_method}
                        </span>
                    )}
                    {server.stars != null && (
                        <span className="inline-flex items-center gap-1 px-2 py-1 rounded text-xs bg-zinc-800 text-zinc-400 border border-zinc-700">
                            <Star className="w-3 h-3" />
                            {server.stars.toLocaleString()}
                        </span>
                    )}
                    {/* Confidence */}
                    <span className="inline-flex items-center gap-1.5 px-2 py-1 rounded text-xs bg-zinc-800 border border-zinc-700">
                        <ShieldCheck className={`w-3 h-3 ${server.confidence >= 70 ? "text-emerald-400" : server.confidence >= 40 ? "text-amber-400" : "text-red-400"}`} />
                        <span className="text-zinc-300">{server.confidence}% confidence</span>
                    </span>
                </div>

                {/* Tags */}
                {server.tags && server.tags.length > 0 && (
                    <div className="flex flex-wrap gap-1.5 mt-3">
                        {server.tags.map((tag: string) => (
                            <span key={tag} className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs bg-zinc-800/80 text-zinc-500 border border-zinc-800">
                                <Tag className="w-2.5 h-2.5" />
                                {tag}
                            </span>
                        ))}
                    </div>
                )}
            </div>

            {/* Two-column layout for recipe + history */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Active Recipe */}
                <Section icon={FileCode} title="Active Install Recipe">
                    {activeRecipe ? (
                        <div className="space-y-3 text-sm">
                            <div className="flex items-center justify-between text-xs">
                                <span className="text-zinc-500">Recipe version</span>
                                <span className="text-zinc-300">v{activeRecipe.recipe_version}</span>
                            </div>
                            <div className="flex items-center justify-between text-xs">
                                <span className="text-zinc-500">Confidence</span>
                                <span className={activeRecipe.confidence >= 70 ? "text-emerald-400" : "text-amber-400"}>
                                    {activeRecipe.confidence}%
                                </span>
                            </div>
                            {activeRecipe.explanation && (
                                <p className="text-xs text-zinc-500 italic leading-relaxed">{activeRecipe.explanation}</p>
                            )}
                            {requiredSecrets.length > 0 && (
                                <div>
                                    <p className="text-xs font-medium text-zinc-400 mb-1">Required secrets</p>
                                    <ul className="space-y-1">
                                        {requiredSecrets.map((k: string) => (
                                            <li key={k} className="text-xs font-mono text-amber-400 bg-amber-950/30 border border-amber-900/40 rounded px-2 py-1">
                                                {k}
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                            )}
                            {Object.keys(requiredEnv as Record<string, string>).length > 0 && (
                                <div>
                                    <p className="text-xs font-medium text-zinc-400 mb-1">Env defaults</p>
                                    <ul className="space-y-1">
                                        {Object.entries(requiredEnv as Record<string, string>).map(([k, v]) => (
                                            <li key={k} className="text-xs font-mono text-zinc-400 bg-zinc-800/50 border border-zinc-700/50 rounded px-2 py-1">
                                                <span className="text-zinc-300">{k}</span>
                                                {v && <span className="text-zinc-600"> = {v}</span>}
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                            )}
                            {/* Template preview (collapsed, code block) */}
                            <details className="group">
                                <summary className="text-xs text-zinc-500 cursor-pointer hover:text-zinc-300 select-none">
                                    View template JSON
                                </summary>
                                <pre className="mt-2 text-xs bg-zinc-900 border border-zinc-800 rounded p-3 overflow-x-auto text-zinc-400 leading-relaxed">
                                    {JSON.stringify(activeRecipe.template, null, 2)}
                                </pre>
                            </details>
                        </div>
                    ) : (
                        <p className="text-sm text-zinc-600 italic">No active recipe. Run validation to generate one.</p>
                    )}
                </Section>

                {/* Validation History */}
                <Section icon={ListChecks} title="Validation History">
                    {loadingRuns ? (
                        <div className="flex gap-2 items-center text-zinc-500 text-sm">
                            <RefreshCw className="w-4 h-4 animate-spin" /> Loading…
                        </div>
                    ) : runs.length === 0 ? (
                        <p className="text-sm text-zinc-600 italic">No validation runs yet.</p>
                    ) : (
                        <ul className="space-y-2">
                            {runs.map((run: ValidationRun) => (
                                <li key={run.uuid} className="text-xs border border-zinc-800 rounded-lg px-3 py-2 bg-zinc-900/30">
                                    <div className="flex items-center justify-between gap-2">
                                        <div className="flex items-center gap-1.5">
                                            {outcomeIcon(run.outcome)}
                                            <span className={`font-medium ${OUTCOME_COLOR[run.outcome] ?? "text-zinc-400"}`}>
                                                {run.outcome}
                                            </span>
                                            {run.tool_count != null && (
                                                <span className="text-zinc-600">· {run.tool_count} tools</span>
                                            )}
                                        </div>
                                        <span className="text-zinc-600 text-xs shrink-0">{formatDate(run.started_at)}</span>
                                    </div>
                                    {run.failure_class && (
                                        <p className="mt-1 text-red-500/80 font-mono">{run.failure_class}</p>
                                    )}
                                    <p className="mt-0.5 text-zinc-600">{run.run_mode} · by {run.performed_by}</p>
                                </li>
                            ))}
                        </ul>
                    )}
                </Section>
            </div>

            {/* Provenance Sources */}
            {sources && sources.length > 0 && (
                <Section icon={Link} title="Provenance Sources">
                    <ul className="space-y-2">
                        {sources.map((src) => (
                            <li key={src.uuid} className="flex items-center justify-between text-xs border border-zinc-800 rounded-lg px-3 py-2 bg-zinc-900/30">
                                <div>
                                    <span className="font-medium text-zinc-300">{src.source_name}</span>
                                    {src.source_url && (
                                        <a
                                            href={src.source_url}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className="ml-2 text-indigo-400 hover:underline"
                                        >
                                            {src.source_url}
                                        </a>
                                    )}
                                </div>
                                <span className="text-zinc-600 text-xs shrink-0">
                                    Last seen {formatDate(src.last_seen_at)}
                                </span>
                            </li>
                        ))}
                    </ul>
                </Section>
            )}

            {/* Metadata */}
            <Section icon={Database} title="Metadata">
                <dl className="grid grid-cols-2 gap-x-6 gap-y-2 text-xs">
                    <MetaRow label="Canonical ID" value={server.canonical_id} mono />
                    <MetaRow label="Auth model" value={server.auth_model} />
                    <MetaRow label="Created" value={formatDate(server.created_at)} />
                    <MetaRow label="Updated" value={formatDate(server.updated_at)} />
                    <MetaRow label="Last seen" value={formatDate(server.last_seen_at)} />
                    <MetaRow label="Last verified" value={formatDate(server.last_verified_at)} />
                    {server.homepage_url && (
                        <dt className="text-zinc-500">Homepage</dt>
                    )}
                    {server.homepage_url && (
                        <dd>
                            <a href={server.homepage_url} target="_blank" rel="noopener noreferrer"
                                className="text-indigo-400 hover:underline"
                            >
                                {server.homepage_url}
                            </a>
                        </dd>
                    )}
                    {server.categories && server.categories.length > 0 && (
                        <>
                            <dt className="text-zinc-500">Categories</dt>
                            <dd className="text-zinc-300">{server.categories.join(", ")}</dd>
                        </>
                    )}
                </dl>
            </Section>

            {/* Install modal */}
            {showInstallModal && (
                <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
                    <div className="w-full max-w-xl rounded-xl border border-zinc-800 bg-zinc-950 shadow-2xl">
                        <div className="px-5 py-4 border-b border-zinc-800">
                            <h2 className="text-base font-semibold text-zinc-100">Install from verified recipe</h2>
                            <p className="text-xs text-zinc-500 mt-1">
                                {server.display_name} · provide any required values before install.
                            </p>
                        </div>
                        <div className="p-5 space-y-4 max-h-[70vh] overflow-y-auto">
                            <div>
                                <label className="block text-xs text-zinc-400 mb-1">Optional name override</label>
                                <input
                                    className="w-full h-9 bg-zinc-900/50 border border-zinc-800 rounded-lg px-3 text-sm text-zinc-200 focus:outline-none focus:border-indigo-500/50"
                                    placeholder="auto-generated from catalog"
                                    value={installName}
                                    onChange={(e) => setInstallName(e.target.value)}
                                />
                            </div>
                            {Object.keys(requiredEnv as Record<string, string>).length > 0 && (
                                <div className="space-y-2">
                                    <p className="text-xs font-medium text-zinc-300">Recommended environment defaults</p>
                                    {Object.entries(requiredEnv as Record<string, string>).map(([key, value]) => (
                                        <div key={key}>
                                            <label className="block text-xs text-zinc-500 mb-1">{key}</label>
                                            <input
                                                className="w-full h-9 bg-zinc-900/50 border border-zinc-800 rounded-lg px-3 text-sm text-zinc-200 focus:outline-none focus:border-indigo-500/50"
                                                value={installEnv[key] ?? value}
                                                onChange={(e) => setInstallEnv((prev) => ({ ...prev, [key]: e.target.value }))}
                                            />
                                        </div>
                                    ))}
                                </div>
                            )}
                            {requiredSecrets.length > 0 && (
                                <div className="space-y-2">
                                    <div className="flex items-center justify-between">
                                        <p className="text-xs font-medium text-zinc-300">Required secrets</p>
                                        <button type="button" className="text-xs text-zinc-500 hover:text-zinc-300" onClick={() => setShowSecrets((v) => !v)}>
                                            {showSecrets ? "Hide" : "Show"}
                                        </button>
                                    </div>
                                    {requiredSecrets.map((key: string) => (
                                        <div key={key}>
                                            <label className="block text-xs text-zinc-500 mb-1">{key}</label>
                                            <input
                                                type={showSecrets ? "text" : "password"}
                                                className="w-full h-9 bg-zinc-900/50 border border-zinc-800 rounded-lg px-3 text-sm text-zinc-200 focus:outline-none focus:border-indigo-500/50"
                                                value={installEnv[key] ?? ""}
                                                onChange={(e) => setInstallEnv((prev) => ({ ...prev, [key]: e.target.value }))}
                                            />
                                        </div>
                                    ))}
                                    {missingSecrets.length > 0 && (
                                        <p className="text-xs text-amber-400">Missing: {missingSecrets.join(", ")}</p>
                                    )}
                                </div>
                            )}
                        </div>
                        <div className="px-5 py-4 border-t border-zinc-800 flex items-center justify-end gap-2">
                            <button
                                className="px-3 py-1.5 text-xs rounded border border-zinc-700 text-zinc-300 hover:bg-zinc-900"
                                onClick={() => { setShowInstallModal(false); setInstallName(""); setInstallEnv({}); setShowSecrets(false); }}
                            >
                                Cancel
                            </button>
                            <button
                                className="px-3 py-1.5 text-xs rounded bg-emerald-600 hover:bg-emerald-500 text-white disabled:bg-emerald-900 disabled:text-emerald-300"
                                disabled={!canSubmit}
                                onClick={() => {
                                    installMutation.mutate({
                                        server_uuid: uuid,
                                        name: installName.trim() || undefined,
                                        env: installEnv,
                                    });
                                }}
                            >
                                {installMutation.isPending ? "Installing…" : "Install"}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

// ---- Sub-components ----

function Section({
    icon: Icon,
    title,
    children,
}: {
    icon: React.ComponentType<{ className?: string }>;
    title: string;
    children: React.ReactNode;
}) {
    return (
        <div className="bg-zinc-900/30 border border-zinc-800 rounded-xl p-5">
            <div className="flex items-center gap-2 mb-4">
                <Icon className="w-4 h-4 text-zinc-500" />
                <h2 className="text-sm font-semibold text-zinc-300">{title}</h2>
            </div>
            {children}
        </div>
    );
}

function MetaRow({ label, value, mono }: { label: string; value?: string | null; mono?: boolean }) {
    if (!value) return null;
    return (
        <>
            <dt className="text-zinc-500">{label}</dt>
            <dd className={`text-zinc-300 truncate ${mono ? "font-mono text-indigo-400" : ""}`}>{value}</dd>
        </>
    );
}
