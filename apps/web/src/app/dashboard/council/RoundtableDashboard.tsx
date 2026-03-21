'use client';

import { useState } from 'react';
import { 
  Gavel, 
  History, 
  Users, 
  Settings, 
  Zap, 
  BarChart3, 
  Layers, 
  Play, 
  Activity,
  Terminal,
  ShieldCheck,
  Eye,
  Plus,
  RefreshCw,
  Trash2
} from 'lucide-react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@borg/ui";
import { Badge } from "@borg/ui";
import { Button } from "@borg/ui";
import { Card, CardContent, CardHeader, CardTitle } from "@borg/ui";
import { trpc } from '@/utils/trpc';

import { SessionGrid } from './SessionGrid';
import { DebateHistoryTable } from './DebateHistoryTable';
import { CouncilMemberGrid } from './CouncilMemberGrid';
import { QuotaConfigPanel } from './QuotaConfigPanel';
import { SmartPilotPanel } from './SmartPilotPanel';
import { VisualArchitecture } from './VisualArchitecture';

export function RoundtableDashboard() {
  const [activeTab, setActiveTab] = useState('dashboard');

  // Queries
  const councilStatus = trpc.council.base.status.useQuery(undefined, {
    refetchInterval: 10000
  });
  
  const sessions = trpc.council.sessions.list.useQuery(undefined, {
    refetchInterval: 5000
  });

  const sessionStats = trpc.council.sessions.stats.useQuery(undefined, {
    refetchInterval: 5000
  });

  const history = trpc.council.history.list.useQuery({ limit: 20 }, {
    refetchInterval: 30000
  });

  const quotaStats = trpc.council.quota.allStats.useQuery(undefined, {
    refetchInterval: 15000
  });

  const smartPilotStatus = trpc.council.smartPilot.status.useQuery(undefined, {
    refetchInterval: 10000
  });

  const visualDiagram = trpc.council.visual.systemDiagram.useQuery(undefined, {
    refetchInterval: 20000
  });

  // Mutations
  const stopSession = trpc.council.sessions.stop.useMutation({
    onSuccess: () => sessions.refetch()
  });

  const resumeSession = trpc.council.sessions.resume.useMutation({
    onSuccess: () => sessions.refetch()
  });

  const deleteSession = trpc.council.sessions.delete.useMutation({
    onSuccess: () => sessions.refetch()
  });

  const toggleSmartPilot = trpc.council.smartPilot.updateConfig.useMutation({
    onSuccess: () => smartPilotStatus.refetch()
  });

  const stopAll = trpc.council.sessions.bulkStop.useMutation({
    onSuccess: () => sessions.refetch()
  });

  return (
    <div className="container mx-auto p-6 space-y-6 max-w-7xl">
      <header className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 border-b pb-6 mb-2">
        <div className="flex items-center gap-4">
          <div className="h-14 w-14 bg-purple-900/20 rounded-2xl flex items-center justify-center border border-purple-500/30 shadow-lg shadow-purple-500/10">
            <Gavel className="h-7 w-7 text-purple-400" />
          </div>
          <div>
            <div className="flex items-center gap-3">
              <h1 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-purple-400 via-blue-400 to-indigo-400 bg-clip-text text-transparent">
                Borg Roundtable
              </h1>
              <Badge variant="outline" className="border-green-500/30 text-green-400 bg-green-500/10 animate-pulse">
                {councilStatus.data?.enabled ? 'Online' : 'Standby'}
              </Badge>
            </div>
            <p className="text-muted-foreground text-sm mt-1 flex items-center gap-2">
              <ShieldCheck className="h-3 w-3 text-purple-500" />
              Multi-Model Orchestration Layer v2.1 (Active)
            </p>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <div className="flex -space-x-2 mr-4 hidden lg:flex">
            {councilStatus.data?.hierarchy.map((c: any) => (
              <div 
                key={c.id} 
                className="h-8 w-8 rounded-full border-2 border-background bg-zinc-800 flex items-center justify-center text-[10px] font-bold"
                title={`${c.name} (${c.supervisorCount} units)`}
              >
                {c.name.charAt(0)}
              </div>
            ))}
          </div>
          <Button variant="outline" size="sm" className="bg-background/50 backdrop-blur-sm border-border/50 gap-2 h-9">
            <Plus className="h-4 w-4" />
            New Session
          </Button>
          <Button className="bg-purple-600 hover:bg-purple-700 shadow-lg shadow-purple-600/20 gap-2 h-9">
            <Play className="h-4 w-4 fill-current" />
            Launch Swarm
          </Button>
        </div>
      </header>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="bg-card/20 border-border/40 backdrop-blur-sm">
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">Active Units</p>
                <h3 className="text-3xl font-black mt-1">{sessionStats.data?.active ?? 0}</h3>
              </div>
              <div className="h-12 w-12 rounded-xl bg-blue-500/10 flex items-center justify-center border border-blue-500/20">
                <Terminal className="h-6 w-6 text-blue-400" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card className="bg-card/20 border-border/40 backdrop-blur-sm">
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">Supervisor Pool</p>
                <h3 className="text-3xl font-black mt-1">{councilStatus.data?.supervisorCount ?? 0}</h3>
              </div>
              <div className="h-12 w-12 rounded-xl bg-purple-500/10 flex items-center justify-center border border-purple-500/20">
                <Users className="h-6 w-6 text-purple-400" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card className="bg-card/20 border-border/40 backdrop-blur-sm">
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">Historical Truth</p>
                <h3 className="text-3xl font-black mt-1">{history.data?.meta.totalRecords ?? 0}</h3>
              </div>
              <div className="h-12 w-12 rounded-xl bg-green-500/10 flex items-center justify-center border border-green-500/20">
                <Activity className="h-6 w-6 text-green-400" />
              </div>
            </div>
          </CardContent>
        </Card>
        <Card className="bg-card/20 border-border/40 backdrop-blur-sm">
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">Consensus Engine</p>
                <h3 className="text-lg font-black mt-1 truncate capitalize">
                  {councilStatus.data?.config.consensusMode?.replace('-', ' ') ?? 'Weighted'}
                </h3>
              </div>
              <div className="h-12 w-12 rounded-xl bg-orange-500/10 flex items-center justify-center border border-orange-500/20">
                <Layers className="h-6 w-6 text-orange-400" />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="bg-background/20 border-b border-border/50 w-full justify-start rounded-none h-auto p-0 mb-6 overflow-x-auto scrollbar-hide">
          <TabsTrigger value="dashboard" className="rounded-none border-b-2 border-transparent data-[state=active]:border-purple-500 data-[state=active]:bg-purple-500/5 px-6 py-3 gap-2">
            <Zap className="h-4 w-4" /> Dashboard
          </TabsTrigger>
          <TabsTrigger value="history" className="rounded-none border-b-2 border-transparent data-[state=active]:border-purple-500 data-[state=active]:bg-purple-500/5 px-6 py-3 gap-2">
            <History className="h-4 w-4" /> History
          </TabsTrigger>
          <TabsTrigger value="supervisors" className="rounded-none border-b-2 border-transparent data-[state=active]:border-purple-500 data-[state=active]:bg-purple-500/5 px-6 py-3 gap-2">
            <Users className="h-4 w-4" /> Supervisors
          </TabsTrigger>
          <TabsTrigger value="usage" className="rounded-none border-b-2 border-transparent data-[state=active]:border-purple-500 data-[state=active]:bg-purple-500/5 px-6 py-3 gap-2">
            <BarChart3 className="h-4 w-4" /> Usage
          </TabsTrigger>
          <TabsTrigger value="visual" className="rounded-none border-b-2 border-transparent data-[state=active]:border-purple-500 data-[state=active]:bg-purple-500/5 px-6 py-3 gap-2">
            <Eye className="h-4 w-4" /> Architecture
          </TabsTrigger>
          <TabsTrigger value="settings" className="rounded-none border-b-2 border-transparent data-[state=active]:border-purple-500 data-[state=active]:bg-purple-500/5 px-6 py-3 gap-2 ml-auto">
            <Settings className="h-4 w-4" />
          </TabsTrigger>
        </TabsList>

        <TabsContent value="dashboard" className="space-y-6 outline-none">
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
            <div className="lg:col-span-8 space-y-6">
              <div className="flex items-center justify-between px-1">
                <h3 className="text-sm font-bold uppercase tracking-wider flex items-center gap-2">
                  <Terminal className="h-4 w-4 text-blue-400" />
                  Terminal Fleet
                </h3>
                <div className="flex items-center gap-2">
                  <Button variant="ghost" size="xs" onClick={() => sessions.refetch()} className="gap-1.5 h-7">
                    <RefreshCw className="h-3 w-3" /> Sync
                  </Button>
                  <Button variant="ghost" size="xs" className="gap-1.5 h-7 text-red-400 hover:text-red-300 hover:bg-red-500/5" onClick={() => stopAll.mutate()}>
                    <Trash2 className="h-3 w-3" /> Purge
                  </Button>
                </div>
              </div>
              
              <SessionGrid 
                sessions={sessions.data ?? []} 
                onStop={(id) => stopSession.mutate({ id })}
                onResume={(id) => resumeSession.mutate({ id })}
                onDelete={(id) => deleteSession.mutate({ id })}
              />
              
              <div className="space-y-3">
                <h3 className="text-sm font-bold uppercase tracking-wider flex items-center gap-2 px-1">
                  <Activity className="h-4 w-4 text-green-400" />
                  Global Activity Log
                </h3>
                <div className="h-[200px] bg-black/40 border border-border/40 rounded-xl font-mono text-[10px] p-4 overflow-y-auto backdrop-blur-md">
                  <div className="text-blue-400/80 mb-1 flex gap-2">
                    <span className="opacity-50">[{new Date().toLocaleTimeString()}]</span>
                    <span>[INFO] Borg Roundtable initialized successfully</span>
                  </div>
                  <div className="text-purple-400/80 mb-1 flex gap-2">
                    <span className="opacity-50">[{new Date().toLocaleTimeString()}]</span>
                    <span>[COUNCIL] Verified {councilStatus.data?.supervisorCount} cognitive units across 3 specializations</span>
                  </div>
                  <div className="text-green-400/80 mb-1 flex gap-2">
                    <span className="opacity-50">[{new Date().toLocaleTimeString()}]</span>
                    <span>[SYSTEM] Local evidence lock synchronized with Master Index v2.1</span>
                  </div>
                </div>
              </div>
            </div>

            <div className="lg:col-span-4 space-y-6">
              <SmartPilotPanel 
                config={smartPilotStatus.data?.config ?? { enabled: false, autoApproveThreshold: 0.7, requireUnanimous: false, maxAutoApprovals: 5 }} 
                activePlans={smartPilotStatus.data?.activePlans ?? []}
                onToggle={(enabled) => toggleSmartPilot.mutate({ enabled })}
                onUpdateConfig={(updates) => toggleSmartPilot.mutate(updates)}
              />

              <Card className="bg-card/20 border-border/40 backdrop-blur-sm">
                <CardHeader className="pb-3 border-b border-border/20">
                  <CardTitle className="text-xs font-bold uppercase tracking-widest text-muted-foreground flex items-center gap-2">
                    <Layers className="h-3 w-3 text-orange-400" />
                    Specialized Hierarchies
                  </CardTitle>
                </CardHeader>
                <CardContent className="pt-4 space-y-3">
                  {councilStatus.data?.hierarchy.map((c: any) => (
                    <div key={c.id} className="group relative overflow-hidden flex items-center justify-between p-3 rounded-lg bg-accent/20 border border-border/30 hover:border-purple-500/30 transition-all">
                      <div className="flex flex-col gap-0.5">
                        <span className="text-xs font-bold">{c.name}</span>
                        <span className="text-[9px] text-muted-foreground uppercase">{c.specialties.slice(0, 2).join(' • ')}</span>
                      </div>
                      <Badge variant="secondary" className="text-[9px] h-5 bg-background/50 border-border/50 px-1.5 font-mono">
                        {c.supervisorCount}u
                      </Badge>
                    </div>
                  ))}
                </CardContent>
              </Card>
            </div>
          </div>
        </TabsContent>

        <TabsContent value="history" className="space-y-6 outline-none">
          <div className="flex items-center justify-between px-1">
            <h3 className="text-sm font-bold uppercase tracking-wider flex items-center gap-2">
              <History className="h-4 w-4 text-purple-400" />
              Verifiable Decision Log
            </h3>
            <Button variant="outline" size="xs" className="h-7 text-[10px] gap-1.5" onClick={() => history.refetch()}>
              <RefreshCw className="h-3 w-3" /> Refresh Audit
            </Button>
          </div>
          <DebateHistoryTable records={history.data?.records ?? []} />
        </TabsContent>

        <TabsContent value="supervisors" className="outline-none">
          <CouncilMemberGrid supervisors={[]} /> {/* To be populated from status.supervisors */}
        </TabsContent>

        <TabsContent value="usage" className="outline-none">
          <QuotaConfigPanel stats={[]} /> {/* To be populated from quotaStats */}
        </TabsContent>

        <TabsContent value="visual" className="h-[700px] outline-none">
          <VisualArchitecture mermaidCode={visualDiagram.data?.mermaid ?? 'graph TD\n  Start'} />
        </TabsContent>
      </Tabs>
    </div>
  );
}
