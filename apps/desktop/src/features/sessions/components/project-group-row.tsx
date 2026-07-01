'use client';

import type { ClaudeProjectSessions } from '@session-manager/schema';
import { SessionRow } from './session-row';

interface ProjectGroupRowProps {
  project: ClaudeProjectSessions;
}

export function ProjectGroupRow({ project }: ProjectGroupRowProps) {
  const sortedSessions = [...(project.sessions ?? [])].sort((a, b) => {
    const aTime = a.updatedAt ? new Date(a.updatedAt).getTime() : 0;
    const bTime = b.updatedAt ? new Date(b.updatedAt).getTime() : 0;
    return bTime - aTime;
  });

  return (
    <div className="flex flex-col gap-1 rounded-lg border p-3">
      <div className="flex items-center justify-between">
        <h3 className="font-semibold">{project.projectLabel}</h3>
        <span className="text-sm text-zinc-500">{project.sessionCount} sessions</span>
      </div>
      <div className="flex flex-col gap-1">
        {sortedSessions.map((session) => (
          <SessionRow key={session.sessionId} session={session} />
        ))}
      </div>
    </div>
  );
}
