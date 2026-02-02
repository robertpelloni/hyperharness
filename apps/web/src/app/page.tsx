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
        <div className="flex items-center gap-4">
          <h1 className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-cyan-500 bg-clip-text text-transparent">Borg Mission Control</h1>
          <a href="/docs" className="text-xs px-2 py-1 bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 rounded hover:bg-blue-100 dark:hover:bg-blue-900 hover:text-blue-600 transition-colors" title="View feature documentation">
            📖 Docs
          </a>
          <a href="/guide" className="text-xs px-2 py-1 bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 rounded hover:bg-green-100 dark:hover:bg-green-900 hover:text-green-600 transition-colors" title="Complete user guide">
            📚 Guide
          </a>
        </div>
        <div className="flex items-center gap-4">
          <GlobalSearch />
          <div className="hidden lg:flex items-center gap-2 text-[10px] text-zinc-500">
            <span title="Drag widgets to rearrange">🔀 Draggable</span>
            <span>•</span>
            <span title="Layout saves automatically">💾 Auto-Save</span>
            <span>•</span>
            <span title="Hover widget titles for help">❓ Tooltips</span>
          </div>
        </div>
      </header>
      <main className="flex-1 bg-zinc-100 dark:bg-black/50 overflow-hidden">
        <DraggableDashboard />
      </main>
    </div>
  );
}
