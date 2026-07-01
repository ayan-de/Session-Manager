# Claude Code Session List Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the desktop MVP that shows a `Claude Code` tab with Claude Code sessions grouped under projects, loaded from the Go backend over HTTP.

**Architecture:** Keep Claude Code filesystem discovery and JSONL metadata parsing in the Go backend under `services/backend/internal/providers/claudecode/`, expose the grouped tree through `GET /api/claude-code/sessions`, and keep the Next.js desktop app as a thin container/presentation client. Share the response shape through hand-written TS types in `packages/schema`, and guard drift with a shared JSON contract fixture plus backend/frontend tests.

**Tech Stack:** Go 1.23 standard library HTTP server and tests, Next.js App Router, React 19, TypeScript 5, npm workspaces, Vitest, Testing Library.

---

## File Map

### Backend

- Create: `services/backend/internal/providers/claudecode/types.go`
- Create: `services/backend/internal/providers/claudecode/discovery.go`
- Create: `services/backend/internal/providers/claudecode/metadata.go`
- Create: `services/backend/internal/providers/claudecode/discovery_test.go`
- Create: `services/backend/internal/providers/claudecode/metadata_test.go`
- Create: `services/backend/internal/providers/claudecode/testdata/sample-session.jsonl`
- Create: `services/backend/internal/providers/claudecode/testdata/sample-session-with-subagents.jsonl`
- Create: `services/backend/internal/api/claude_code_sessions.go`
- Create: `services/backend/internal/api/claude_code_sessions_test.go`
- Modify: `services/backend/cmd/backend/main.go`

### Shared contract

- Create: `packages/schema/src/claude-code-sessions.ts`
- Modify: `packages/schema/src/index.ts`
- Create: `testdata/contracts/claude-code-sessions.json`

### Frontend client and test harness

- Create: `packages/api-client/src/claude-code-sessions.ts`
- Modify: `packages/api-client/src/index.ts`
- Modify: `apps/desktop/package.json`
- Create: `apps/desktop/vitest.config.ts`
- Create: `apps/desktop/test/setup.ts`

### Desktop feature

- Create: `apps/desktop/src/features/sessions/queries/get-claude-code-project-sessions.ts`
- Create: `apps/desktop/src/features/sessions/queries/get-claude-code-project-sessions.test.ts`
- Create: `apps/desktop/src/features/sessions/containers/claude-code-sessions-container.tsx`
- Create: `apps/desktop/src/features/sessions/components/claude-code-tab.tsx`
- Create: `apps/desktop/src/features/sessions/components/project-session-tree.tsx`
- Create: `apps/desktop/src/features/sessions/components/project-group-row.tsx`
- Create: `apps/desktop/src/features/sessions/components/session-row.tsx`
- Create: `apps/desktop/src/features/sessions/components/claude-code-tab.test.tsx`
- Modify: `apps/desktop/app/page.tsx`
- Modify: `apps/desktop/app/layout.tsx`

## Task 1: Define the shared contract and lock the backend endpoint shape first

**Files:**
- Create: `packages/schema/src/claude-code-sessions.ts`
- Modify: `packages/schema/src/index.ts`
- Create: `testdata/contracts/claude-code-sessions.json`
- Create: `services/backend/internal/api/claude_code_sessions.go`
- Create: `services/backend/internal/api/claude_code_sessions_test.go`
- Modify: `services/backend/cmd/backend/main.go`

- [ ] **Step 1: Write the failing backend handler test for `GET /api/claude-code/sessions`**

```go
package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	api "github.com/session-manager/backend/internal/api"
)

func TestClaudeCodeSessionsHandlerReturnsContractShape(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "..", "..", "testdata", "contracts", "claude-code-sessions.json")
	wantBytes, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	handler := api.NewClaudeCodeSessionsHandler(api.StaticClaudeCodeSessionTree())
	req := httptest.NewRequest(http.MethodGet, "/api/claude-code/sessions", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}

	var got any
	var want any
	if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if err := json.Unmarshal(wantBytes, &want); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	if diff := api.CompareJSON(got, want); diff != "" {
		t.Fatalf("response mismatch (-got +want):\n%s", diff)
	}
}
```

- [ ] **Step 2: Run the backend handler test to verify it fails**

Run: `go test ./internal/api -run TestClaudeCodeSessionsHandlerReturnsContractShape -v`
Expected: FAIL with missing package, missing handler, or missing helper symbols.

- [ ] **Step 3: Add the shared TypeScript contract and JSON fixture**

```ts
// packages/schema/src/claude-code-sessions.ts
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

// packages/schema/src/index.ts
export * from "./claude-code-sessions";

// testdata/contracts/claude-code-sessions.json
[
  {
    "projectId": "home-ayan-de-Projects-session-manager",
    "projectLabel": "Session-Manager",
    "projectPathHint": "/home/ayan-de/Projects/Session-Manager",
    "sessionCount": 2,
    "lastUpdatedAt": "2026-07-01T10:30:00Z",
    "sessions": [
      {
        "sessionId": "11111111-1111-1111-1111-111111111111",
        "title": "Review grouped session list MVP",
        "firstPrompt": "review the structure we will use",
        "createdAt": "2026-07-01T10:00:00Z",
        "updatedAt": "2026-07-01T10:30:00Z",
        "messageCount": 18,
        "gitBranch": "main",
        "hasSubagents": true
      },
      {
        "sessionId": "22222222-2222-2222-2222-222222222222",
        "title": "Sketch Claude Code storage reader",
        "firstPrompt": "check where they store it",
        "createdAt": "2026-07-01T09:10:00Z",
        "updatedAt": "2026-07-01T09:45:00Z",
        "messageCount": 9,
        "gitBranch": "main",
        "hasSubagents": false
      }
    ]
  }
]
```

- [ ] **Step 4: Add the minimal backend handler and wire it into `main.go`**

```go
// services/backend/internal/api/claude_code_sessions.go
package api

import (
	"encoding/json"
	"net/http"
)

type ClaudeSessionTreeProvider func() any

func NewClaudeCodeSessionsHandler(provider ClaudeSessionTreeProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(provider())
	})
}

func StaticClaudeCodeSessionTree() any {
	return []map[string]any{}
}

func CompareJSON(got any, want any) string {
	if got == nil && want == nil {
		return ""
	}
	if equalJSON(got, want) {
		return ""
	}
	return "json differs"
}

// services/backend/cmd/backend/main.go
http.Handle("/api/claude-code/sessions", api.NewClaudeCodeSessionsHandler(api.StaticClaudeCodeSessionTree))
```

- [ ] **Step 5: Update the handler implementation to return the fixture-shaped payload and make the test pass**

```go
func StaticClaudeCodeSessionTree() any {
	return []map[string]any{
		{
			"projectId":       "home-ayan-de-Projects-session-manager",
			"projectLabel":    "Session-Manager",
			"projectPathHint": "/home/ayan-de/Projects/Session-Manager",
			"sessionCount":    2,
			"lastUpdatedAt":   "2026-07-01T10:30:00Z",
			"sessions": []map[string]any{
				{
					"sessionId":    "11111111-1111-1111-1111-111111111111",
					"title":        "Review grouped session list MVP",
					"firstPrompt":  "review the structure we will use",
					"createdAt":    "2026-07-01T10:00:00Z",
					"updatedAt":    "2026-07-01T10:30:00Z",
					"messageCount": 18,
					"gitBranch":    "main",
					"hasSubagents": true,
				},
			},
		},
	}
}
```

- [ ] **Step 6: Run the backend handler test again**

Run: `go test ./internal/api -run TestClaudeCodeSessionsHandlerReturnsContractShape -v`
Expected: PASS

- [ ] **Step 7: Commit the contract-first baseline**

```bash
git add testdata/contracts/claude-code-sessions.json packages/schema/src/index.ts packages/schema/src/claude-code-sessions.ts services/backend/internal/api/claude_code_sessions.go services/backend/internal/api/claude_code_sessions_test.go services/backend/cmd/backend/main.go
git commit -m "feat: add claude code sessions contract"
```

## Task 2: Build Claude Code project discovery and JSONL metadata extraction in Go

**Files:**
- Create: `services/backend/internal/providers/claudecode/types.go`
- Create: `services/backend/internal/providers/claudecode/discovery.go`
- Create: `services/backend/internal/providers/claudecode/metadata.go`
- Create: `services/backend/internal/providers/claudecode/discovery_test.go`
- Create: `services/backend/internal/providers/claudecode/metadata_test.go`
- Create: `services/backend/internal/providers/claudecode/testdata/sample-session.jsonl`
- Create: `services/backend/internal/providers/claudecode/testdata/sample-session-with-subagents.jsonl`

- [ ] **Step 1: Write the failing discovery test for grouping sessions under projects**

```go
package claudecode_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/session-manager/backend/internal/providers/claudecode"
)

func TestListProjectsWithSessionsGroupsByProjectAndSortsNewestFirst(t *testing.T) {
	root := t.TempDir()
	projectA := filepath.Join(root, "-home-ayan-de-Projects-Session-Manager")
	projectB := filepath.Join(root, "-home-ayan-de-Projects-claude-code")
	if err := os.MkdirAll(filepath.Join(projectA, "11111111-1111-1111-1111-111111111111", "subagents"), 0o755); err != nil {
		t.Fatalf("mkdir subagents: %v", err)
	}
	if err := os.MkdirAll(projectB, 0o755); err != nil {
		t.Fatalf("mkdir projectB: %v", err)
	}

	copyFixture(t, filepath.Join(projectA, "11111111-1111-1111-1111-111111111111.jsonl"), "testdata/sample-session-with-subagents.jsonl")
	copyFixture(t, filepath.Join(projectB, "22222222-2222-2222-2222-222222222222.jsonl"), "testdata/sample-session.jsonl")

	projects, err := claudecode.ListProjectsWithSessions(root)
	if err != nil {
		t.Fatalf("ListProjectsWithSessions: %v", err)
	}

	if len(projects) != 2 {
		t.Fatalf("len(projects) = %d, want 2", len(projects))
	}
	if projects[0].ProjectLabel != "Session-Manager" {
		t.Fatalf("first project = %q, want %q", projects[0].ProjectLabel, "Session-Manager")
	}
	if !projects[0].Sessions[0].HasSubagents {
		t.Fatalf("expected first session to report subagents")
	}
}
```

- [ ] **Step 2: Write the failing metadata extraction test for lightweight title and timestamps**

```go
func TestReadSessionMetadataExtractsFirstPromptAndTimes(t *testing.T) {
	path := filepath.Join("testdata", "sample-session.jsonl")
	meta, err := claudecode.ReadSessionMetadata(path, false)
	if err != nil {
		t.Fatalf("ReadSessionMetadata: %v", err)
	}

	if meta.SessionID != "22222222-2222-2222-2222-222222222222" {
		t.Fatalf("session id = %q", meta.SessionID)
	}
	if meta.FirstPrompt != "check where they store it" {
		t.Fatalf("first prompt = %q", meta.FirstPrompt)
	}
	if meta.Title == "" {
		t.Fatal("title should not be empty")
	}
	if meta.UpdatedAt == "" {
		t.Fatal("updatedAt should not be empty")
	}
}
```

- [ ] **Step 3: Run the provider tests to verify they fail**

Run: `go test ./internal/providers/claudecode -run 'TestListProjectsWithSessionsGroupsByProjectAndSortsNewestFirst|TestReadSessionMetadataExtractsFirstPromptAndTimes' -v`
Expected: FAIL with missing package, missing functions, or missing fixtures.

- [ ] **Step 4: Add the sample JSONL fixtures and provider types**

```go
// services/backend/internal/providers/claudecode/types.go
package claudecode

type ProjectSessions struct {
	ProjectID       string           `json:"projectId"`
	ProjectLabel    string           `json:"projectLabel"`
	ProjectPathHint string           `json:"projectPathHint"`
	SessionCount    int              `json:"sessionCount"`
	LastUpdatedAt   string           `json:"lastUpdatedAt,omitempty"`
	Sessions        []SessionSummary `json:"sessions"`
}

type SessionSummary struct {
	SessionID     string `json:"sessionId"`
	Title         string `json:"title"`
	FirstPrompt   string `json:"firstPrompt,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
	UpdatedAt     string `json:"updatedAt,omitempty"`
	MessageCount  int    `json:"messageCount,omitempty"`
	GitBranch     string `json:"gitBranch,omitempty"`
	HasSubagents  bool   `json:"hasSubagents"`
}

// services/backend/internal/providers/claudecode/testdata/sample-session.jsonl
{"type":"user","sessionId":"22222222-2222-2222-2222-222222222222","timestamp":"2026-07-01T09:10:00Z","message":{"role":"user","content":[{"type":"text","text":"check where they store it"}]},"cwd":"/home/ayan-de/Projects/claude-code","gitBranch":"main"}
{"type":"assistant","sessionId":"22222222-2222-2222-2222-222222222222","timestamp":"2026-07-01T09:45:00Z","message":{"role":"assistant","content":[{"type":"text","text":"They store transcripts under ~/.claude/projects/..."}]},"cwd":"/home/ayan-de/Projects/claude-code","gitBranch":"main"}
```

- [ ] **Step 5: Implement minimal discovery and metadata extraction to satisfy the tests**

```go
// services/backend/internal/providers/claudecode/discovery.go
package claudecode

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ListProjectsWithSessions(projectsRoot string) ([]ProjectSessions, error) {
	entries, err := os.ReadDir(projectsRoot)
	if err != nil {
		return nil, err
	}

	projects := make([]ProjectSessions, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		projectDir := filepath.Join(projectsRoot, entry.Name())
		sessions, err := readProjectSessions(projectDir)
		if err != nil || len(sessions) == 0 {
			continue
		}
		project := ProjectSessions{
			ProjectID:       sanitizeProjectID(entry.Name()),
			ProjectLabel:    humanizeProjectLabel(entry.Name()),
			ProjectPathHint: decodeProjectPath(entry.Name()),
			SessionCount:    len(sessions),
			LastUpdatedAt:   sessions[0].UpdatedAt,
			Sessions:        sessions,
		}
		projects = append(projects, project)
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].LastUpdatedAt > projects[j].LastUpdatedAt
	})
	return projects, nil
}

func humanizeProjectLabel(raw string) string {
	trimmed := strings.TrimPrefix(raw, "-")
	parts := strings.Split(trimmed, "-")
	if len(parts) == 0 {
		return raw
	}
	return parts[len(parts)-1]
}

// services/backend/internal/providers/claudecode/metadata.go
func ReadSessionMetadata(path string, hasSubagents bool) (SessionSummary, error) {
	lines, err := os.ReadFile(path)
	if err != nil {
		return SessionSummary{}, err
	}
	return parseSessionSummary(lines, hasSubagents)
}
```

- [ ] **Step 6: Run the provider tests again**

Run: `go test ./internal/providers/claudecode -v`
Expected: PASS

- [ ] **Step 7: Commit the provider implementation**

```bash
git add services/backend/internal/providers/claudecode
git commit -m "feat: add claude code session discovery"
```

## Task 3: Replace the static handler with the real backend provider and cover error mapping

**Files:**
- Modify: `services/backend/internal/api/claude_code_sessions.go`
- Modify: `services/backend/internal/api/claude_code_sessions_test.go`
- Modify: `services/backend/cmd/backend/main.go`

- [ ] **Step 1: Write the failing handler test for provider errors and real provider wiring**

```go
func TestClaudeCodeSessionsHandlerReturnsServerErrorWhenProviderFails(t *testing.T) {
	handler := api.NewClaudeCodeSessionsHandler(func() (any, error) {
		return nil, errors.New("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/claude-code/sessions", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusInternalServerError)
	}
}
```

- [ ] **Step 2: Run the API tests to verify the new case fails**

Run: `go test ./internal/api -v`
Expected: FAIL because the handler signature does not yet support provider errors.

- [ ] **Step 3: Refactor the handler to depend on the real provider shape**

```go
type ClaudeSessionTreeProvider func() (any, error)

func NewClaudeCodeSessionsHandler(provider ClaudeSessionTreeProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, err := provider()
		if err != nil {
			http.Error(w, `{"error":"unable to load Claude Code sessions"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	})
}

func NewClaudeCodeSessionTreeProvider(projectsRoot string) ClaudeSessionTreeProvider {
	return func() (any, error) {
		return claudecode.ListProjectsWithSessions(projectsRoot)
	}
}
```

- [ ] **Step 4: Wire the real provider into `main.go`**

```go
projectsRoot := filepath.Join(os.Getenv("HOME"), ".claude", "projects")
http.Handle(
	"/api/claude-code/sessions",
	api.NewClaudeCodeSessionsHandler(api.NewClaudeCodeSessionTreeProvider(projectsRoot)),
)
```

- [ ] **Step 5: Run the backend API and provider tests together**

Run: `go test ./internal/api ./internal/providers/claudecode -v`
Expected: PASS

- [ ] **Step 6: Commit the real backend endpoint**

```bash
git add services/backend/internal/api/claude_code_sessions.go services/backend/internal/api/claude_code_sessions_test.go services/backend/cmd/backend/main.go
git commit -m "feat: expose claude code sessions endpoint"
```

## Task 4: Add the desktop test harness and typed API client

**Files:**
- Modify: `apps/desktop/package.json`
- Create: `apps/desktop/vitest.config.ts`
- Create: `apps/desktop/test/setup.ts`
- Create: `packages/api-client/src/claude-code-sessions.ts`
- Modify: `packages/api-client/src/index.ts`
- Create: `apps/desktop/src/features/sessions/queries/get-claude-code-project-sessions.ts`
- Create: `apps/desktop/src/features/sessions/queries/get-claude-code-project-sessions.test.ts`

- [ ] **Step 1: Write the failing frontend query test against the shared JSON contract**

```ts
import { describe, expect, it, vi } from 'vitest';
import contractFixture from '../../../../../testdata/contracts/claude-code-sessions.json';
import { getClaudeCodeProjectSessions } from './get-claude-code-project-sessions';

describe('getClaudeCodeProjectSessions', () => {
  it('returns the grouped project tree from the backend contract', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => contractFixture,
      }),
    );

    const result = await getClaudeCodeProjectSessions();

    expect(result[0]?.projectLabel).toBe('Session-Manager');
    expect(result[0]?.sessions[0]?.hasSubagents).toBe(true);
  });
});
```

- [ ] **Step 2: Install the test dependencies and add scripts**

```json
{
  "scripts": {
    "test": "vitest run",
    "test:watch": "vitest"
  },
  "devDependencies": {
    "@testing-library/jest-dom": "^6.6.3",
    "@testing-library/react": "^16.0.1",
    "jsdom": "^25.0.1",
    "vitest": "^2.0.5"
  }
}
```

Run: `npm install -w apps/desktop -D vitest @testing-library/react @testing-library/jest-dom jsdom`
Expected: packages installed successfully.

- [ ] **Step 3: Add the minimal Vitest configuration and setup file**

```ts
// apps/desktop/vitest.config.ts
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./test/setup.ts'],
  },
});

// apps/desktop/test/setup.ts
import '@testing-library/jest-dom/vitest';
```

- [ ] **Step 4: Implement the typed API client and desktop query wrapper**

```ts
// packages/api-client/src/claude-code-sessions.ts
import type { ClaudeProjectSessions } from '@session-manager/schema';

const BASE_URL = 'http://localhost:8080';

export async function fetchClaudeCodeProjectSessions(): Promise<ClaudeProjectSessions[]> {
  const res = await fetch(`${BASE_URL}/api/claude-code/sessions`);
  if (!res.ok) {
    throw new Error('Unable to load Claude Code sessions');
  }
  return (await res.json()) as ClaudeProjectSessions[];
}

// apps/desktop/src/features/sessions/queries/get-claude-code-project-sessions.ts
import { fetchClaudeCodeProjectSessions } from '@session-manager/api-client';

export function getClaudeCodeProjectSessions() {
  return fetchClaudeCodeProjectSessions();
}
```

- [ ] **Step 5: Run the frontend query test and make sure it passes**

Run: `npm run test -w apps/desktop -- get-claude-code-project-sessions.test.ts`
Expected: PASS

- [ ] **Step 6: Commit the desktop client foundation**

```bash
git add apps/desktop/package.json apps/desktop/vitest.config.ts apps/desktop/test/setup.ts packages/api-client/src/index.ts packages/api-client/src/claude-code-sessions.ts apps/desktop/src/features/sessions/queries/get-claude-code-project-sessions.ts apps/desktop/src/features/sessions/queries/get-claude-code-project-sessions.test.ts
git commit -m "feat: add claude code sessions api client"
```

## Task 5: Build the Claude Code tab container and grouped project/session UI

**Files:**
- Create: `apps/desktop/src/features/sessions/containers/claude-code-sessions-container.tsx`
- Create: `apps/desktop/src/features/sessions/components/claude-code-tab.tsx`
- Create: `apps/desktop/src/features/sessions/components/project-session-tree.tsx`
- Create: `apps/desktop/src/features/sessions/components/project-group-row.tsx`
- Create: `apps/desktop/src/features/sessions/components/session-row.tsx`
- Create: `apps/desktop/src/features/sessions/components/claude-code-tab.test.tsx`
- Modify: `apps/desktop/app/page.tsx`
- Modify: `apps/desktop/app/layout.tsx`

- [ ] **Step 1: Write the failing UI test for loading and grouped rendering**

```tsx
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import contractFixture from '../../../../../../testdata/contracts/claude-code-sessions.json';
import { ClaudeCodeSessionsContainer } from '../containers/claude-code-sessions-container';

vi.mock('../queries/get-claude-code-project-sessions', () => ({
  getClaudeCodeProjectSessions: vi.fn().mockResolvedValue(contractFixture),
}));

describe('ClaudeCodeSessionsContainer', () => {
  it('renders grouped projects and nested sessions', async () => {
    render(<ClaudeCodeSessionsContainer />);

    expect(screen.getByText(/loading claude code sessions/i)).toBeInTheDocument();

    expect(await screen.findByText('Session-Manager')).toBeInTheDocument();
    expect(await screen.findByText('Review grouped session list MVP')).toBeInTheDocument();
  });
});
```

- [ ] **Step 2: Run the UI test to verify it fails**

Run: `npm run test -w apps/desktop -- claude-code-tab.test.tsx`
Expected: FAIL with missing components or query modules.

- [ ] **Step 3: Implement the container and presentational components with clear responsibilities**

```tsx
// apps/desktop/src/features/sessions/containers/claude-code-sessions-container.tsx
'use client';

import { useEffect, useState } from 'react';
import type { ClaudeProjectSessions } from '@session-manager/schema';
import { getClaudeCodeProjectSessions } from '../queries/get-claude-code-project-sessions';
import { ClaudeCodeTab } from '../components/claude-code-tab';

export function ClaudeCodeSessionsContainer() {
  const [projects, setProjects] = useState<ClaudeProjectSessions[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getClaudeCodeProjectSessions()
      .then(setProjects)
      .catch(() => setError('Unable to load Claude Code sessions'))
      .finally(() => setLoading(false));
  }, []);

  return <ClaudeCodeTab projects={projects} loading={loading} error={error} />;
}

// apps/desktop/src/features/sessions/components/claude-code-tab.tsx
import type { ClaudeProjectSessions } from '@session-manager/schema';
import { ProjectSessionTree } from './project-session-tree';

export function ClaudeCodeTab({ projects, loading, error }: { projects: ClaudeProjectSessions[]; loading: boolean; error: string | null; }) {
  if (loading) return <p>Loading Claude Code sessions...</p>;
  if (error) return <p>{error}</p>;
  if (projects.length === 0) return <p>No Claude Code sessions found.</p>;
  return <ProjectSessionTree projects={projects} />;
}
```

- [ ] **Step 4: Replace the boilerplate landing page with the Claude Code tab shell**

```tsx
// apps/desktop/app/page.tsx
import { ClaudeCodeSessionsContainer } from '@/src/features/sessions/containers/claude-code-sessions-container';

export default function Home() {
  return (
    <main className="min-h-screen bg-stone-100 text-stone-950">
      <section className="mx-auto flex min-h-screen w-full max-w-6xl flex-col px-6 py-10">
        <header className="mb-8 border-b border-stone-300 pb-4">
          <p className="text-sm uppercase tracking-[0.2em] text-stone-500">Session Manager</p>
          <h1 className="font-mono text-3xl font-semibold">Claude Code</h1>
        </header>
        <ClaudeCodeSessionsContainer />
      </section>
    </main>
  );
}

// apps/desktop/app/layout.tsx
export const metadata: Metadata = {
  title: 'Session Manager',
  description: 'Browse Claude Code sessions by project',
};
```

- [ ] **Step 5: Run the frontend tests, lint, and build checks**

Run: `npm run test -w apps/desktop && npm run lint -w apps/desktop && npm run build -w apps/desktop`
Expected: all commands PASS

- [ ] **Step 6: Run the backend tests and start the server for manual verification**

Run: `go test ./... && go run ./cmd/backend/main.go`
Expected: tests PASS, then `Backend starting on :8080`

- [ ] **Step 7: Manually verify the desktop app against the live backend**

Run: `npm run dev -w apps/desktop`
Expected: the desktop page loads, shows the `Claude Code` heading, and renders project-grouped session rows from the backend.

- [ ] **Step 8: Commit the desktop MVP**

```bash
git add apps/desktop/app/layout.tsx apps/desktop/app/page.tsx apps/desktop/src/features/sessions packages/api-client/src packages/schema/src/claude-code-sessions.ts
git commit -m "feat: add claude code sessions desktop view"
```

## Self-Review Checklist

- Spec coverage: backend provider, backend endpoint, shared contract, drift guard, desktop container, grouped project UI, partial error handling, and tests are all mapped to tasks above
- Placeholder scan: no `TODO`, `TBD`, or implicit “write tests later” steps remain
- Type consistency: `ClaudeProjectSessions` and `ClaudeSessionSummary` naming is consistent across backend, shared TS contract, API client, query, and UI tasks
