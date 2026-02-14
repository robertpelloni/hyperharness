
"use client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { httpBatchLink } from "@trpc/client";
import React, { useState } from "react";
import { trpc } from "./trpc";

function resolveTrpcUrl(): string {
    const fromEnv = process.env.NEXT_PUBLIC_TRPC_URL?.trim();
    if (fromEnv) {
        const normalized = fromEnv.replace(/\/$/, "");
        return normalized.endsWith("/trpc") ? normalized : `${normalized}/trpc`;
    }

    if (typeof window !== "undefined") {
        const { protocol, hostname } = window.location;
        const isLocal = hostname === "localhost" || hostname === "127.0.0.1";
        if (isLocal) {
            return `${protocol}//${hostname}:4000/trpc`;
        }
        return `${window.location.origin}/trpc`;
    }

    return "http://localhost:4000/trpc";
}

export function TRPCProvider({ children }: { children: React.ReactNode }) {
    const [queryClient] = useState(() => new QueryClient());
    const [trpcClient] = useState(() =>
        trpc.createClient({
            links: [
                httpBatchLink({
                    url: resolveTrpcUrl(),
                }),
            ],
        })
    );

    return (
        <trpc.Provider client={trpcClient} queryClient={queryClient}>
            <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
        </trpc.Provider>
    );
}
