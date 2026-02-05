
import { KnowledgeGraph } from '@/components/visualizer/KnowledgeGraph';

export default function KnowledgePage() {
    return (
        <div className="p-6 space-y-6">
            <div className="flex flex-col gap-2">
                <h1 className="text-3xl font-bold tracking-tight">Knowledge Graph</h1>
                <p className="text-muted-foreground">
                    Visualize the structural relationships within the Borg's memory.
                </p>
            </div>

            <KnowledgeGraph />
        </div>
    );
}
