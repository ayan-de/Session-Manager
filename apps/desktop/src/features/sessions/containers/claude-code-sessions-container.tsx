'use client';

import { useEffect, useState } from 'react';
import { getClaudeCodeProjectSessions } from '../queries/get-claude-code-project-sessions';
import type { ClaudeProjectSessions } from '@session-manager/schema';

interface UseClaudeCodeSessionsResult {
  data: ClaudeProjectSessions[] | undefined;
  isLoading: boolean;
  error: Error | null;
}

export function useClaudeCodeSessions(): UseClaudeCodeSessionsResult {
  const [data, setData] = useState<ClaudeProjectSessions[] | undefined>(undefined);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function fetchSessions() {
      setIsLoading(true);
      setError(null);

      try {
        const result = await getClaudeCodeProjectSessions();

        if (!cancelled) {
          setData(result);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err : new Error('Failed to fetch sessions'));
        }
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    }

    fetchSessions();

    return () => {
      cancelled = true;
    };
  }, []);

  return { data, isLoading, error };
}
