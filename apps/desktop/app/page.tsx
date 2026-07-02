'use client';

import { useState } from 'react';
import ProviderTabBar from '@/src/components/ProviderTabBar';
import { ClaudeCodeTab } from '@/src/features/sessions/components/claude-code-tab';

export default function Home() {
  const [selectedProvider, setSelectedProvider] = useState('claude');

  return (
    <div className="flex flex-col h-screen">
      <ProviderTabBar selectedProvider={selectedProvider} onSelectProvider={setSelectedProvider} />
      <div className="flex-1 overflow-hidden">
        {selectedProvider === 'claude' && <ClaudeCodeTab />}
        {selectedProvider === 'opencode' && (
          <div className="flex flex-1 items-center justify-center text-text-muted">
            OpenCode sessions coming soon...
          </div>
        )}
      </div>
    </div>
  );
}
