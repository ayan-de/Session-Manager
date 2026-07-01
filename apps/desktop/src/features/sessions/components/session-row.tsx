'use client';

import type { ClaudeSessionSummary } from '@session-manager/schema';

interface SessionRowProps {
  session: ClaudeSessionSummary;
}

export function SessionRow({ session }: SessionRowProps) {
  return (
    <div className="flex items-center gap-2 rounded px-2 py-1 hover:bg-zinc-100 dark:hover:bg-zinc-800">
      <div className="flex flex-1 flex-col">
        <span className="font-medium">{session.title}</span>
        <div className="flex items-center gap-2 text-sm text-zinc-500">
          {session.gitBranch && (
            <span className="rounded bg-zinc-200 px-1 text-xs dark:bg-zinc-700">
              {session.gitBranch}
            </span>
          )}
          {session.messageCount !== undefined && <span>{session.messageCount} messages</span>}
          {session.hasSubagents && (
            <span className="rounded bg-purple-100 px-1 text-xs text-purple-700 dark:bg-purple-900 dark:text-purple-300">
              Subagents
            </span>
          )}
        </div>
      </div>
    </div>
  );
}
