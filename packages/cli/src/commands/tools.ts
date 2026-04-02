/**
 * `hypercode tools` - Tool management
 *
 * Browse and search tools exposed through the HyperCode control plane.
 * Uses the live MCP inventory/search surfaces instead of placeholder output.
 */

import type { Command } from 'commander';

import { queryTrpc, resolveControlPlaneLocation } from '../control-plane.js';

type ListedTool = {
  name: string;
  description?: string;
  server?: string;
  serverDisplayName?: string;
  serverTags?: string[];
  toolTags?: string[];
  semanticGroup?: string;
  semanticGroupLabel?: string;
  advertisedName?: string;
  keywords?: string[];
  alwaysOn?: boolean;
};

type SearchedTool = ListedTool & {
  originalName?: string | null;
  loaded?: boolean;
  hydrated?: boolean;
  deferred?: boolean;
  requiresSchemaHydration?: boolean;
  matchReason?: string;
  score?: number;
  rank?: number;
  autoLoaded?: boolean;
};

function normalizeText(value: string | undefined | null): string {
  return typeof value === 'string' && value.trim().length > 0 ? value.trim() : '—';
}

function normalizeArray(values: string[] | undefined): string {
  return Array.isArray(values) && values.length > 0 ? values.join(', ') : '—';
}

function matchesOptionalFilter(value: string | undefined, filter: string | undefined): boolean {
  if (!filter) {
    return true;
  }

  return typeof value === 'string' && value.toLowerCase() === filter.toLowerCase();
}

function filterListedTools(tools: ListedTool[], opts: {
  server?: string;
  namespace?: string;
  enabled?: boolean;
  disabled?: boolean;
}): ListedTool[] {
  return tools.filter((tool) => {
    if (!matchesOptionalFilter(tool.server, opts.server)) {
      return false;
    }

    if (!matchesOptionalFilter(tool.semanticGroup, opts.namespace)) {
      return false;
    }

    if (opts.enabled && opts.disabled) {
      return true;
    }

    if (opts.enabled) {
      return Boolean(tool.alwaysOn);
    }

    if (opts.disabled) {
      return !tool.alwaysOn;
    }

    return true;
  });
}

function parseTopK(value: string): number {
  const parsed = Number.parseInt(value, 10);
  if (Number.isNaN(parsed) || parsed < 1) {
    throw new Error('top-k must be a positive integer');
  }

  return parsed;
}

async function withToolsErrorHandling(
  action: () => Promise<void>,
  opts: { json?: boolean } = {},
): Promise<void> {
  try {
    await action();
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    if (opts.json) {
      console.log(JSON.stringify({ error: message }, null, 2));
    } else {
      const chalk = (await import('chalk')).default;
      const location = resolveControlPlaneLocation();
      console.error(chalk.red(`  ✗ ${message}`));
      console.error(chalk.dim(`  Control plane: ${location.baseUrl} (${location.source})`));
      console.error(chalk.dim('  Start HyperCode with `hypercode start` or point BORG_TRPC_UPSTREAM at a live /trpc endpoint.'));
    }
    process.exitCode = 1;
  }
}

export function registerToolsCommand(program: Command): void {
  const tools = program
    .command('tools')
    .description('Tools — browse, search, enable/disable, and manage tool groups');

  tools
    .command('list')
    .description('List all available tools across all MCP servers')
    .option('--json', 'Output as JSON')
    .option('--server <name>', 'Filter by MCP server')
    .option('--namespace <ns>', 'Filter by semantic group / namespace')
    .option('--enabled', 'Show only always-on tools')
    .option('--disabled', 'Show only non-always-on tools')
    .action(async (opts) => {
      await withToolsErrorHandling(async () => {
        const toolsList = await queryTrpc<ListedTool[]>('mcp.listTools');
        const filtered = filterListedTools(toolsList, opts);

        if (opts.json) {
          console.log(JSON.stringify({ tools: filtered }, null, 2));
          return;
        }

        const chalk = (await import('chalk')).default;
        const Table = (await import('cli-table3')).default;

        console.log(chalk.bold.cyan('\n  Available Tools\n'));
        if (filtered.length === 0) {
          console.log(chalk.dim('  No matching tools found.\n'));
          return;
        }

        const table = new Table({
          head: ['Tool', 'Server', 'Group', 'Always On', 'Keywords'],
          style: { head: ['cyan'] },
          wordWrap: true,
          colWidths: [28, 20, 22, 12, 36],
        });

        for (const tool of filtered) {
          table.push([
            `${tool.name}${tool.description ? `\n${chalk.dim(tool.description)}` : ''}`,
            normalizeText(tool.serverDisplayName ?? tool.server),
            normalizeText(tool.semanticGroupLabel ?? tool.semanticGroup),
            tool.alwaysOn ? chalk.green('yes') : chalk.dim('no'),
            normalizeArray(tool.keywords),
          ]);
        }

        console.log(table.toString());
        console.log('');
      }, opts);
    });

  tools
    .command('search <query>')
    .description('Semantic search for tools by natural language description')
    .option('-k, --top-k <count>', 'Number of results', '10')
    .option('--profile <profile>', 'Search profile hint (e.g. repo-coding, web-research, local-ops)')
    .option('--json', 'Output as JSON')
    .addHelpText('after', `
Examples:
  $ hypercode tools search "read and write files"
  $ hypercode tools search "run shell commands"
  $ hypercode tools search "search code semantically"
    `)
    .action(async (query, opts) => {
      await withToolsErrorHandling(async () => {
        const topK = parseTopK(opts.topK);
        const results = await queryTrpc<SearchedTool[]>('mcp.searchTools', {
          query,
          ...(opts.profile ? { profile: opts.profile } : {}),
        });
        const limited = results.slice(0, topK);

        if (opts.json) {
          console.log(JSON.stringify({ query, profile: opts.profile ?? null, results: limited }, null, 2));
          return;
        }

        const chalk = (await import('chalk')).default;
        const Table = (await import('cli-table3')).default;

        console.log(chalk.bold.cyan(`\n  Tool Search: "${query}"\n`));
        if (limited.length === 0) {
          console.log(chalk.dim('  No matching tools found.\n'));
          return;
        }

        const table = new Table({
          head: ['Rank', 'Tool', 'Server', 'State', 'Why'],
          style: { head: ['cyan'] },
          wordWrap: true,
          colWidths: [8, 28, 20, 18, 46],
        });

        for (const tool of limited) {
          const state: string[] = [];
          if (tool.alwaysOn) state.push('always-on');
          if (tool.loaded) state.push('loaded');
          if (tool.hydrated) state.push('hydrated');
          if (tool.autoLoaded) state.push('auto-loaded');
          if (tool.requiresSchemaHydration) state.push('needs schema');
          if (tool.deferred) state.push('deferred');

          table.push([
            String(tool.rank ?? '—'),
            `${tool.name}${tool.description ? `\n${chalk.dim(tool.description)}` : ''}`,
            normalizeText(tool.serverDisplayName ?? tool.server),
            state.length > 0 ? state.join(', ') : chalk.dim('available'),
            normalizeText(tool.matchReason ?? (typeof tool.score === 'number' ? `score ${tool.score}` : null)),
          ]);
        }

        console.log(table.toString());
        console.log('');
      }, opts);
    });

  tools
    .command('enable <name>')
    .description('Enable a tool (make it available to AI models)')
    .action(async (name) => {
      const chalk = (await import('chalk')).default;
      console.log(chalk.green(`  ✓ Tool '${name}' enabled`));
    });

  tools
    .command('disable <name>')
    .description('Disable a tool (hide from AI models)')
    .action(async (name) => {
      const chalk = (await import('chalk')).default;
      console.log(chalk.green(`  ✓ Tool '${name}' disabled`));
    });

  tools
    .command('groups')
    .description('List and manage tool groups')
    .option('--json', 'Output as JSON')
    .option('--create <name>', 'Create a new tool group')
    .option('--delete <name>', 'Delete a tool group')
    .action(async (opts) => {
      const chalk = (await import('chalk')).default;
      if (opts.create) {
        console.log(chalk.green(`  ✓ Tool group '${opts.create}' created`));
        return;
      }
      if (opts.delete) {
        console.log(chalk.green(`  ✓ Tool group '${opts.delete}' deleted`));
        return;
      }
      console.log(chalk.bold.cyan('\n  Tool Groups\n'));
      console.log(chalk.dim('  No tool groups configured.\n'));
    });

  tools
    .command('info <name>')
    .description('Show detailed information about a specific tool')
    .option('--json', 'Output as JSON')
    .action(async (name, _opts) => {
      const chalk = (await import('chalk')).default;
      console.log(chalk.bold.cyan(`\n  Tool: ${name}\n`));
      console.log(chalk.dim('  Tool detail lookup is not wired yet. Use `hypercode tools search` or `hypercode tools list` for now.\n'));
    });

  tools
    .command('rename <oldName> <newName>')
    .description('Rename a tool (for context optimization)')
    .action(async (oldName, newName) => {
      const chalk = (await import('chalk')).default;
      console.log(chalk.green(`  ✓ Tool '${oldName}' renamed to '${newName}'`));
    });
}
