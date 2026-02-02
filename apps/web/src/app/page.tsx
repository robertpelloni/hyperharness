"use client";
import React, { useState } from 'react';
import DraggableDashboard from "../components/DraggableDashboard";
import { LoginPage } from "../components/LoginPage";
import { GlobalSearch } from "../components/GlobalSearch";

export default function Home() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  if (!isAuthenticated) {
    return <LoginPage onLogin={() => setIsAuthenticated(true)} />;
  }

  return (
    <div className="flex min-h-screen flex-col bg-zinc-50 font-sans dark:bg-black">
      <header className="flex items-center justify-between px-6 py-4 border-b border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 z-10">
        <h1 className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-cyan-500 bg-clip-text text-transparent">Borg Mission Control</h1>
        <div className="flex items-center gap-6">
          <GlobalSearch />
          <div className="text-xs text-zinc-500 hidden md:block">
            Draggable Layout • Auto-Saving
          </div>
        </div>
      </header>
      <main className="flex-1 bg-zinc-100 dark:bg-black/50 overflow-hidden">
        <DraggableDashboard />
      </main>
    </div>
  );
}
