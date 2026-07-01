import { fetchClaudeCodeProjectSessions } from '@session-manager/api-client';
import type { ClaudeProjectSessions } from '@session-manager/schema';

export async function getClaudeCodeProjectSessions(): Promise<ClaudeProjectSessions[]> {
  return fetchClaudeCodeProjectSessions();
}
