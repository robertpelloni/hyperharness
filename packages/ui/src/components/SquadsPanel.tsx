'use client';

import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "./ui/card";
import { Badge } from "./ui/badge";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "./ui/dialog";
import { ScrollArea } from "./ui/scroll-area";
import { Users, Plus, Trash2, Terminal, GitBranch, Activity, RefreshCw } from 'lucide-react';
import { toast } from "sonner";
import { trpc } from '@/utils/trpc';
import { useRouter } from 'next/navigation';

export function SquadsPanel() {
    const router = useRouter();
    const [spawnOpen, setSpawnOpen] = useState(false);
    const [newBranch, setNewBranch] = useState('');
    const [newGoal, setNewGoal] = useState('');

    const utils = trpc.useUtils();
    const membersQuery = trpc.squad.list.useQuery(undefined, {
        refetchInterval: 5000 // Refresh every 5s
    });

    const spawnMutation = trpc.squad.spawn.useMutation({
        onSuccess: () => {
            toast.success("Squad member spawned successfully!");
            setSpawnOpen(false);
            setNewBranch('');
            setNewGoal('');
            utils.squad.list.invalidate();
        },
        onError: (err) => {
            toast.error(`Failed to spawn member: ${err.message}`);
        }
    });

    const killMutation = trpc.squad.kill.useMutation({
        onSuccess: () => {
            toast.success("Squad member terminated.");
            utils.squad.list.invalidate();
        },
        onError: (err) => {
            toast.error(`Failed to kill member: ${err.message}`);
        }
    });

    const handleSpawn = () => {
        if (!newBranch || !newGoal) {
            toast.error("Branch and Goal are required");
            return;
        }
        spawnMutation.mutate({ branch: newBranch, goal: newGoal });
    };

    return (
        <div className="container mx-auto p-6 max-w-6xl space-y-8">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
                        <Users className="w-8 h-8" />
                        Squads
                    </h1>
                    <p className="text-muted-foreground mt-1">
                        Manage autonomous agents running in parallel Git Worktrees.
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline" size="icon" onClick={() => membersQuery.refetch()}>
                        <RefreshCw className={`w-4 h-4 ${membersQuery.isFetching ? 'animate-spin' : ''}`} />
                    </Button>
                    <Dialog open={spawnOpen} onOpenChange={setSpawnOpen}>
                        <DialogTrigger asChild>
                            <Button>
                                <Plus className="w-4 h-4 mr-2" />
                                Spawn Member
                            </Button>
                        </DialogTrigger>
                        <DialogContent>
                            <DialogHeader>
                                <DialogTitle>Spawn New Squad Member</DialogTitle>
                                <DialogDescription>
                                    This will create a new Git Worktree and launch a Director agent to work on the specified task.
                                </DialogDescription>
                            </DialogHeader>
                            <div className="grid gap-4 py-4">
                                <div className="grid gap-2">
                                    <label htmlFor="branch" className="text-sm font-medium">Branch Name</label>
                                    <Input
                                        id="branch"
                                        placeholder="feat/my-feature"
                                        value={newBranch}
                                        onChange={(e) => setNewBranch(e.target.value)}
                                    />
                                </div>
                                <div className="grid gap-2">
                                    <label htmlFor="goal" className="text-sm font-medium">Goal / Task</label>
                                    <Input
                                        id="goal"
                                        placeholder="Refactor the login component..."
                                        value={newGoal}
                                        onChange={(e) => setNewGoal(e.target.value)}
                                    />
                                </div>
                            </div>
                            <DialogFooter>
                                <Button onClick={handleSpawn} disabled={spawnMutation.isPending}>
                                    {spawnMutation.isPending ? 'Spawning...' : 'Spawn'}
                                </Button>
                            </DialogFooter>
                        </DialogContent>
                    </Dialog>
                </div>
            </div>

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {membersQuery.isLoading ? (
                    <div className="col-span-full py-12 text-center text-muted-foreground">Loading members...</div>
                ) : membersQuery.data?.length === 0 ? (
                    <div className="col-span-full py-12 text-center border-2 border-dashed rounded-lg">
                        <Users className="w-12 h-12 mx-auto text-muted-foreground/50 mb-4" />
                        <h3 className="text-lg font-medium">No Active Squads</h3>
                        <p className="text-muted-foreground">Spawn a member to start parallel work.</p>
                    </div>
                ) : (
                    membersQuery.data?.map((member: any) => (
                        <Card key={member.id} className="relative overflow-hidden">
                            <CardHeader className="pb-3">
                                <div className="flex justify-between items-start">
                                    <CardTitle className="text-lg font-mono truncate mr-2" title={member.branch}>
                                        {member.branch}
                                    </CardTitle>
                                    <Badge variant={member.status === 'busy' ? "default" : "secondary"}>
                                        {member.status}
                                    </Badge>
                                </div>
                                <CardDescription className="font-mono text-xs text-muted-foreground truncate" title={member.id}>
                                    ID: {member.id}
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-4">
                                    <div className="flex items-center text-sm text-muted-foreground">
                                        <GitBranch className="w-4 h-4 mr-2" />
                                        Worktree Active
                                    </div>
                                    <div className="flex items-center text-sm text-muted-foreground">
                                        <Activity className="w-4 h-4 mr-2" />
                                        Director: {member.active ? 'Running' : 'Idle'}
                                    </div>

                                    <div className="flex justify-end pt-2">
                                        <Button
                                            variant="destructive"
                                            size="sm"
                                            onClick={() => {
                                                if (confirm(`Are you sure you want to kill ${member.branch}? This will delete the worktree.`)) {
                                                    killMutation.mutate({ branch: member.branch });
                                                }
                                            }}
                                            disabled={killMutation.isPending}
                                        >
                                            <Trash2 className="w-4 h-4 mr-2" />
                                            Kill
                                        </Button>
                                    </div>
                                </div>
                            </CardContent>
                            {/* Decorative background accent */}
                            <div className="absolute top-0 right-0 p-4 opacity-5 pointer-events-none">
                                <Terminal className="w-24 h-24" />
                            </div>
                        </Card>
                    ))
                )}
            </div>
        </div>
    );
}
