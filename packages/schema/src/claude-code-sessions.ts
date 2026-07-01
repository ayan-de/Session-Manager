export interface ClaudeProjectSessions {
  projectId: string;
  projectLabel: string;
  projectPathHint: string;
  sessionCount: number;
  lastUpdatedAt?: string;
  sessions: ClaudeSessionSummary[];
}

export interface ClaudeSessionSummary {
  sessionId: string;
  title: string;
  firstPrompt?: string;
  createdAt?: string;
  updatedAt?: string;
  messageCount?: number;
  gitBranch?: string;
  hasSubagents: boolean;
}
