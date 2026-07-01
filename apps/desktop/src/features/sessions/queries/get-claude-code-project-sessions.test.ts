import { describe, it, expect, vi } from 'vitest';
import { getClaudeCodeProjectSessions } from './get-claude-code-project-sessions';

describe('getClaudeCodeProjectSessions', () => {
  it('returns parsed ClaudeProjectSessions array', async () => {
    const mockData = [
      {
        projectId: 'proj-1',
        projectLabel: 'My Project',
        projectPathHint: '/home/user/project',
        sessionCount: 3,
        lastUpdatedAt: '2024-01-15T10:30:00Z',
        sessions: [
          {
            sessionId: 'sess-1',
            title: 'Work on feature X',
            firstPrompt: 'Implement the new feature',
            createdAt: '2024-01-15T09:00:00Z',
            updatedAt: '2024-01-15T10:30:00Z',
            messageCount: 42,
            gitBranch: 'feature-x',
            hasSubagents: true,
          },
        ],
      },
    ];

    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockData),
      })
    ) as unknown as typeof fetch;

    const result = await getClaudeCodeProjectSessions();

    expect(result).toEqual(mockData);
    expect(result[0]).toMatchObject({
      projectId: 'proj-1',
      projectLabel: 'My Project',
      projectPathHint: '/home/user/project',
      sessionCount: 3,
    });
    expect(result[0].sessions[0]).toMatchObject({
      sessionId: 'sess-1',
      title: 'Work on feature X',
      hasSubagents: true,
    });
  });

  it('throws when fetch fails', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: false,
        statusText: 'Internal Server Error',
      })
    ) as unknown as typeof fetch;

    await expect(getClaudeCodeProjectSessions()).rejects.toThrow(
      'Failed to fetch sessions: Internal Server Error'
    );
  });
});
