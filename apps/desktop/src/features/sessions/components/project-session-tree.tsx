'use client';

import type { ClaudeProjectSessions } from '@session-manager/schema';
import { ProjectGroupRow } from './project-group-row';

interface ProjectSessionTreeProps {
  projects: ClaudeProjectSessions[];
}

export function ProjectSessionTree({ projects }: ProjectSessionTreeProps) {
  const sortedProjects = [...projects].sort((a, b) => {
    const aTime = a.lastUpdatedAt ? new Date(a.lastUpdatedAt).getTime() : 0;
    const bTime = b.lastUpdatedAt ? new Date(b.lastUpdatedAt).getTime() : 0;
    return bTime - aTime;
  });

  return (
    <div className="flex flex-col gap-2 p-4">
      {sortedProjects.map((project) => (
        <ProjectGroupRow key={project.projectId} project={project} />
      ))}
    </div>
  );
}
