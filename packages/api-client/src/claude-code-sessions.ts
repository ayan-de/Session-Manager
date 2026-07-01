import type { ClaudeProjectSessions } from "@session-manager/schema";

const BASE_URL = "http://localhost:8080";

export async function fetchClaudeCodeProjectSessions(): Promise<ClaudeProjectSessions[]> {
  const res = await fetch(`${BASE_URL}/api/claude-code/sessions`);
  if (!res.ok) {
    throw new Error(`Failed to fetch sessions: ${res.statusText}`);
  }
  return res.json();
}
