
import fs from 'fs';
import path from 'path';
import { z } from 'zod';

export interface ProjectTask {
    id: string;
    description: string;
    status: 'TODO' | 'IN_PROGRESS' | 'DONE';
    sourceFile: string;
    lineNumber: number;
}

const TaskStateSchema = z.enum([' ', 'x', '/']);

const ParsedTaskLineSchema = z.object({
    indent: z.number().int().nonnegative(),
    state: TaskStateSchema,
    text: z.string().min(1),
    lineNumber: z.number().int().positive(),
});

type ParsedTaskLine = z.infer<typeof ParsedTaskLineSchema>;

export class ProjectTracker {
    private rootDir: string;

    constructor(rootDir: string) {
        this.rootDir = rootDir;
    }

    /**
     * Scans task.md and ROADMAP.md for the next actionable task.
     * Priority: 
     * 1. In Progress (but not done) in task.md
     * 2. Todo in task.md
     * 3. In Progress in ROADMAP.md
     * 4. Todo in ROADMAP.md
     */
    public getNextTask(): ProjectTask | null {
        // Default to local task.md in root
        const localTaskMd = path.join(this.rootDir, 'task.md');
        const docsTaskMd = path.join(this.rootDir, 'docs', 'task.md');

        let targetPath = fs.existsSync(localTaskMd) ? localTaskMd : (fs.existsSync(docsTaskMd) ? docsTaskMd : null);

        if (targetPath) {
            const task = this.findTaskInFile(targetPath);
            if (task) return task;
        }

        const roadmapPath = path.join(this.rootDir, 'ROADMAP.md');
        if (fs.existsSync(roadmapPath)) {
            const task = this.findTaskInFile(roadmapPath);
            if (task) return task;
        }

        return null;
    }

    public getStatus() {
        // Calculate progress based on checked items
        // For efficiency, maybe cache this or just scan on demand since files are small
        const taskMdPath = path.join(this.rootDir, 'task.md'); // Simplified
        let total = 0;
        let done = 0;

        if (fs.existsSync(taskMdPath)) {
            const content = fs.readFileSync(taskMdPath, 'utf-8');
            total += (content.match(/-\s*\[[ x/]\]/g) || []).length;
            done += (content.match(/-\s*\[x\]/g) || []).length;
        }

        const currentTask = this.getNextTask();

        return {
            taskId: currentTask ? currentTask.id : 'idle',
            status: currentTask ? 'busy' : 'idle',
            progress: total > 0 ? Math.round((done / total) * 100) : 0,
            currentTask: currentTask ? currentTask.description : 'Idle'
        };
    }

    private findTaskInFile(filePath: string): ProjectTask | null {
        const content = fs.readFileSync(filePath, 'utf-8');
        const lines = content.split('\n');

        // 1. Check for In Progress task: - [/]
        for (let i = 0; i < lines.length; i++) {
            const parsed = this.parseTaskLine(lines[i], i + 1);
            if (parsed?.state === '/') {
                const subTask = this.findFirstSubTask(lines, i + 1, parsed.indent);
                if (subTask) {
                    return {
                        id: `task-${path.basename(filePath)}-${subTask.line}`,
                        description: subTask.text,
                        status: 'TODO',
                        sourceFile: filePath,
                        lineNumber: subTask.line
                    };
                }

                return {
                    id: `task-${path.basename(filePath)}-${i + 1}`,
                    description: parsed.text,
                    status: 'IN_PROGRESS',
                    sourceFile: filePath,
                    lineNumber: i + 1
                };
            }
        }

        // 2. Check for Todo task: - [ ]
        for (let i = 0; i < lines.length; i++) {
            const parsed = this.parseTaskLine(lines[i], i + 1);
            if (parsed?.state === ' ') {
                return {
                    id: `task-${path.basename(filePath)}-${i + 1}`,
                    description: parsed.text,
                    status: 'TODO',
                    sourceFile: filePath,
                    lineNumber: i + 1
                };
            }
        }

        return null;
    }

    private parseTaskLine(line: string, lineNumber: number): ParsedTaskLine | null {
        const match = line.match(/^(\s*)-\s*\[([ x/])\]\s*(.+?)\s*$/);
        if (!match) {
            return null;
        }

        const parsed = ParsedTaskLineSchema.safeParse({
            indent: match[1].length,
            state: match[2],
            text: match[3].trim(),
            lineNumber,
        });

        return parsed.success ? parsed.data : null;
    }

    private findFirstSubTask(lines: string[], startIndex: number, parentIndent: number): { text: string, line: number } | null {

        for (let i = startIndex; i < lines.length; i++) {
            const line = lines[i];
            if (line.trim() === '') continue;

            const parsed = this.parseTaskLine(line, i + 1);
            if (!parsed) {
                if (line.match(/^\s*#/)) break;
                continue;
            }

            if (parsed.indent <= parentIndent) {
                break;
            }

            if (parsed.state === ' ' || parsed.state === '/') {
                return { text: parsed.text, line: parsed.lineNumber };
            }
        }

        return null;
    }

    public completeTask(task: ProjectTask): void {
        const content = fs.readFileSync(task.sourceFile, 'utf-8');
        const lines = content.split('\n');

        if (lines[task.lineNumber - 1]) {
            lines[task.lineNumber - 1] = lines[task.lineNumber - 1].replace('[ ]', '[x]').replace('[/]', '[x]');
            fs.writeFileSync(task.sourceFile, lines.join('\n'));
            console.log(`[ProjectTracker] Marked task as done: ${task.description}`);
        }
    }

    public startTask(task: ProjectTask): void {
        const content = fs.readFileSync(task.sourceFile, 'utf-8');
        const lines = content.split('\n');

        if (lines[task.lineNumber - 1]) {
            lines[task.lineNumber - 1] = lines[task.lineNumber - 1].replace('[ ]', '[/]');
            fs.writeFileSync(task.sourceFile, lines.join('\n'));
            console.log(`[ProjectTracker] Marked task as in-progress: ${task.description}`);
        }
    }
}
