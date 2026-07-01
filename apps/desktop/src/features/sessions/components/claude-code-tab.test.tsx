import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { ClaudeCodeTab } from '../components/claude-code-tab';
import type { ClaudeProjectSessions } from '@session-manager/schema';

vi.mock('../containers/claude-code-sessions-container', () => ({
  useClaudeCodeSessions: vi.fn(),
}));

import { useClaudeCodeSessions } from '../containers/claude-code-sessions-container';

const mockUseClaudeCodeSessions = useClaudeCodeSessions as ReturnType<typeof vi.fn>;

describe('ClaudeCodeTab', () => {
  const mockProjects: ClaudeProjectSessions[] = [
    {
      projectId: 'proj-1',
      projectLabel: 'My Project',
      projectPathHint: '/home/user/project',
      sessionCount: 2,
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
        {
          sessionId: 'sess-2',
          title: 'Fix bug Y',
          firstPrompt: 'Fix the bug',
          createdAt: '2024-01-14T09:00:00Z',
          updatedAt: '2024-01-14T15:00:00Z',
          messageCount: 10,
          gitBranch: 'bugfix-y',
          hasSubagents: false,
        },
      ],
    },
    {
      projectId: 'proj-2',
      projectLabel: 'Another Project',
      projectPathHint: '/home/user/another',
      sessionCount: 1,
      lastUpdatedAt: '2024-01-13T08:00:00Z',
      sessions: [
        {
          sessionId: 'sess-3',
          title: 'Initial work',
          firstPrompt: 'Start here',
          createdAt: '2024-01-13T08:00:00Z',
          updatedAt: '2024-01-13T08:00:00Z',
          messageCount: 5,
          hasSubagents: false,
        },
      ],
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders loading state', () => {
    mockUseClaudeCodeSessions.mockReturnValue({
      data: undefined,
      isLoading: true,
      error: null,
    });

    render(<ClaudeCodeTab />);

    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  it('renders error state', () => {
    mockUseClaudeCodeSessions.mockReturnValue({
      data: undefined,
      isLoading: false,
      error: new Error('Failed to load'),
    });

    render(<ClaudeCodeTab />);

    expect(screen.getByText(/error/i)).toBeInTheDocument();
    expect(screen.getByText(/failed to load/i)).toBeInTheDocument();
  });

  it('renders empty state', () => {
    mockUseClaudeCodeSessions.mockReturnValue({
      data: [],
      isLoading: false,
      error: null,
    });

    render(<ClaudeCodeTab />);

    expect(screen.getByText(/no sessions/i)).toBeInTheDocument();
  });

  it('renders projects and sessions', () => {
    mockUseClaudeCodeSessions.mockReturnValue({
      data: mockProjects,
      isLoading: false,
      error: null,
    });

    render(<ClaudeCodeTab />);

    expect(screen.getByText('My Project')).toBeInTheDocument();
    expect(screen.getByText('Another Project')).toBeInTheDocument();
    expect(screen.getByText('Work on feature X')).toBeInTheDocument();
    expect(screen.getByText('Fix bug Y')).toBeInTheDocument();
    expect(screen.getByText('Initial work')).toBeInTheDocument();
  });

  it('displays session metadata', () => {
    mockUseClaudeCodeSessions.mockReturnValue({
      data: mockProjects,
      isLoading: false,
      error: null,
    });

    render(<ClaudeCodeTab />);

    expect(screen.getByText('feature-x')).toBeInTheDocument();
    expect(screen.getByText('42 messages')).toBeInTheDocument();
  });

  it('shows subagents indicator when present', () => {
    mockUseClaudeCodeSessions.mockReturnValue({
      data: mockProjects,
      isLoading: false,
      error: null,
    });

    render(<ClaudeCodeTab />);

    const subagentBadges = screen.getAllByText(/subagent/i);
    expect(subagentBadges.length).toBe(1);
  });
});
