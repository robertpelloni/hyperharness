type OrchestratorConfigEnv = Record<string, string | undefined>;

export function resolveConfiguredOrchestratorBase(
  env: OrchestratorConfigEnv = process.env,
): string | null {
  const configured =
    env.NEXT_PUBLIC_BORG_ORCHESTRATOR_URL?.trim()
    || env.NEXT_PUBLIC_AUTOPILOT_URL?.trim()
    || '';

  if (!configured) {
    return null;
  }

  return configured.replace(/\/$/, '');
}

