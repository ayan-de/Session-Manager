'use client';

import { memo, useCallback, useRef } from 'react';
import { ChevronLeft, ChevronRight } from 'lucide-react';

interface Provider {
  id: string;
  name: string;
  logo: string;
}

const PROVIDERS: Provider[] = [
  { id: 'claude', name: 'Claude', logo: '/logos/claude_code.svg' },
  { id: 'opencode', name: 'OpenCode', logo: '/logos/opencode.svg' },
];

interface ProviderTabBarProps {
  selectedProvider: string;
  onSelectProvider: (provider: string) => void;
}

const ProviderTab = memo(function ProviderTab({
  provider,
  isSelected,
  onSelect,
}: {
  provider: Provider;
  isSelected: boolean;
  onSelect: (id: string) => void;
}) {
  return (
    <button
      onClick={() => onSelect(provider.id)}
      className={`flex flex-col items-center justify-center gap-1 px-4 py-2 text-[10px] font-semibold whitespace-nowrap cursor-pointer border-0 outline-none focus:outline-none flex-none transition-colors ${
        isSelected ? 'text-white' : 'text-text-muted hover:text-text-main hover:bg-hover-subtle'
      }`}
      style={isSelected ? { backgroundColor: '#3b82f6' } : undefined}
    >
      <div className="relative flex items-center justify-center w-7 h-7">
        {isSelected && (
          <div
            className="absolute inset-0 rounded-full bg-white"
            style={{
              boxShadow:
                '0 0 10px 3px rgba(255, 255, 255, 0.4), 0 0 25px 8px rgba(255, 255, 255, 0.3), 0 0 50px 15px rgba(255, 255, 255, 0.2), 0 0 100px 30px rgba(255, 255, 255, 0.1)',
            }}
          />
        )}
        <img
          src={provider.logo}
          alt=""
          className={`w-5 h-5 object-contain relative z-10 ${isSelected ? 'opacity-100' : 'opacity-80'}`}
        />
      </div>
      <span>{provider.name}</span>
    </button>
  );
});

function ProviderTabBar({ selectedProvider, onSelectProvider }: ProviderTabBarProps) {
  const scrollRef = useRef<HTMLDivElement>(null);

  const scroll = useCallback((dir: 'left' | 'right') => {
    if (!scrollRef.current) return;
    scrollRef.current.scrollBy({ left: dir === 'right' ? 120 : -120, behavior: 'smooth' });
  }, []);

  return (
    <div className="flex items-center border-b border-border-subtle bg-secondary/10 will-change-transform">
      <button
        onClick={() => scroll('left')}
        className="flex-shrink-0 p-1.5 text-text-muted hover:text-text-main hover:bg-hover-subtle rounded transition-all cursor-pointer border-0"
        aria-label="Scroll left"
      >
        <ChevronLeft className="w-3.5 h-3.5" />
      </button>

      <div
        ref={scrollRef}
        className="flex-1 min-w-0 flex items-center overflow-x-auto px-1 py-2.5 scrollbar-none will-change-transform"
      >
        {PROVIDERS.map((p) => (
          <ProviderTab
            key={p.id}
            provider={p}
            isSelected={selectedProvider === p.id}
            onSelect={onSelectProvider}
          />
        ))}
      </div>

      <button
        onClick={() => scroll('right')}
        className="flex-shrink-0 p-1.5 text-text-muted hover:text-text-main hover:bg-hover-subtle rounded transition-all cursor-pointer border-0"
        aria-label="Scroll right"
      >
        <ChevronRight className="w-3.5 h-3.5" />
      </button>
    </div>
  );
}

export default memo(ProviderTabBar);
export { PROVIDERS };
