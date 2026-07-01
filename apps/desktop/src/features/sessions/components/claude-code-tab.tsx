'use client';

import { useClaudeCodeSessions } from '../containers/claude-code-sessions-container';
import { ProjectSessionTree } from './project-session-tree';

export function ClaudeCodeTab() {
  const { data, isLoading, error } = useClaudeCodeSessions();

  if (isLoading) {
    return (
      <div className="flex flex-1 items-center justify-center">
        <p>Loading sessions...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-1 items-center justify-center">
        <p>Error: {error.message}</p>
      </div>
    );
  }

  if (!data || data.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center">
        <p>No sessions found</p>
      </div>
    );
  }

  return (
    <div className="flex flex-1 flex-col overflow-auto">
      <ProjectSessionTree projects={data} />
    </div>
  );
}
