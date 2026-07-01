# Claude Code Session List MVP Design

## Goal

Build the first real desktop MVP around Claude Code session browsing.

The desktop app should show a `Claude Code` tab that groups sessions by project, then lists that project's sessions underneath. This MVP is read-only. It does not import into the canonical schema, export across providers, sync to GitHub, or resume sessions yet.

This slice is intentionally narrow. It provides immediate user value, validates the desktop architecture, and creates the right seams for later provider expansion.

## MVP Outcome

Users can open the desktop app, navigate to the `Claude Code` tab, and see:

1. Claude Code projects discovered from Claude's local session storage
2. Sessions grouped under each project
3. Lightweight session metadata for browsing

The UI should optimize for scanability and reliability, not full transcript rendering.

## Confirmed Claude Code Storage Contract

Based on the local Claude Code repository at `/home/ayan-de/Projects/claude-code`:

- Main session transcripts are stored under `~/.claude/projects/<project-dir>/<session-id>.jsonl`
- Subagent transcripts are stored under `~/.claude/projects/<project-dir>/<session-id>/subagents/...`
- Claude Code itself uses session storage helpers to resolve the projects root and transcript paths
- Claude Code's own resume flow progressively loads session metadata before loading full conversations

This MVP should follow the same spirit: browse metadata first, defer heavy transcript loading.

## Scope

### In scope

- Desktop `Claude Code` tab
- Grouping sessions under projects
- Reading Claude Code session files from local disk
- Parsing lightweight metadata for list rendering
- Loading, empty, and partial-error states
- Stable interfaces that support later detail view work

### Out of scope

- OpenCode support
- Canonical schema import
- Cross-provider conversion
- GitHub sync
- Full transcript rendering
- Session resume actions
- Session search across providers
- Background filesystem watching

## Product Shape

The primary information architecture is:

- Provider tab: `Claude Code`
- Inside the tab: project groups
- Inside each project group: session rows

This should behave like a project explorer, not a flat inbox.

### UX rules

- Projects are sorted by most recently updated session descending
- Sessions inside a project are sorted by last updated descending
- The most recently active project may be expanded by default
- Other projects may start collapsed when the list is large
- Corrupt or unreadable session files should not break the whole tab

## Architecture

The MVP should use a clear container/presentation pattern with provider-specific file logic in the Go backend service, not in the desktop frontend.

### Seams

#### 1. Go backend Claude Code provider module

Responsibility: know how Claude Code stores sessions on disk and return a list-ready model.

This module should live under `services/backend/internal/providers/claudecode/`.

This module should:

- resolve the Claude config root
- resolve the projects root
- enumerate project directories
- enumerate main session transcript files
- detect whether a session has subagents
- extract lightweight list metadata from JSONL
- keep list-view parsing separate from any future full transcript importer

This module should not:

- render UI
- know about desktop tabs or React state
- parse more transcript content than needed for list view

#### 2. Go backend HTTP handler

Responsibility: expose the Claude Code project/session tree to desktop clients.

This module should:

- call the Claude Code provider module
- map backend errors into stable HTTP responses
- return `ClaudeProjectSessions[]` as JSON

Recommended first endpoint:

- `GET /api/claude-code/sessions`

This endpoint is the use-case seam for the MVP.

#### 3. Desktop API client module

Responsibility: make a typed request to the backend and return the session tree.

This module should:

- call `GET /api/claude-code/sessions`
- expose typed response models to the desktop container
- avoid filesystem access, transcript parsing, and Claude-specific storage logic

#### 4. Desktop container module

Responsibility: orchestrate loading and map use-case results into UI state.

This module should:

- trigger loading on tab entry or page load
- hold loading, success, empty, and error states
- support refresh behavior later without changing provider logic
- keep selected project/session UI state if interaction grows later

#### 5. Presentational components

Responsibility: render grouped project/session data only.

These components should:

- receive already prepared view models
- display project rows and nested session rows
- display badges, timestamps, counts, and empty/error states

These components should not:

- access the filesystem
- parse paths
- parse JSONL
- know Claude Code storage rules

## Recommended Module Layout

The exact filenames can follow existing repo conventions, but the structure should resemble:

```text
apps/desktop/
  app/
    ...
  src/
    features/
      sessions/
        containers/
          claude-code-sessions-container.tsx
        components/
          claude-code-tab.tsx
          project-session-tree.tsx
          project-group-row.tsx
          session-row.tsx
        models/
          session-list-view-model.ts
        queries/
          get-claude-code-project-sessions.ts

packages/
  api-client/
    src/
      claude-code-sessions.ts
  schema/
    src/
      claude-code-sessions.ts

services/
  backend/
    internal/
      api/
        claude_code_sessions.go
      providers/
        claudecode/
          discovery.go
          metadata.go
          types.go
```

The frontend should not contain its own Claude Code filesystem reader for this MVP. If implementation speed requires temporary shortcuts, keep them inside the backend service, not in the desktop app.

## Data Contracts

The UI should not consume raw transcript files. It should consume a stable tree model.

```ts
export type ClaudeProjectSessions = {
  projectId: string
  projectLabel: string
  projectPathHint: string
  sessionCount: number
  lastUpdatedAt?: string
  sessions: ClaudeSessionSummary[]
}

export type ClaudeSessionSummary = {
  sessionId: string
  title: string
  firstPrompt?: string
  createdAt?: string
  updatedAt?: string
  messageCount?: number
  gitBranch?: string
  hasSubagents: boolean
}
```

### Contract notes

- `projectId` should be stable for rendering and expansion state
- `projectLabel` should be human-readable
- `projectPathHint` may be shortened for display safety
- `title` should prefer a user-facing label if one can be derived cheaply; otherwise fall back to first prompt or session id
- `messageCount` may be omitted if computing it cheaply is not reliable in the first version

These types should be shared across the backend response contract and the frontend client. `packages/schema` is the best current home if the repo wants one shared TypeScript definition for desktop and future React Native consumers.

### Type sync decision

For the MVP, Go response structs and TypeScript response types may be hand-written. This is accepted short-term duplication because the MVP has one endpoint and one small response shape.

To reduce silent drift without adding generation tooling yet:

- keep the response model intentionally small
- add a shared JSON fixture that represents the `GET /api/claude-code/sessions` response
- test the backend handler against that fixture shape
- test the desktop API client or schema consumer against the same fixture

Upgrade to generated types when the API grows beyond this initial surface. The trigger should be any of: a third backend endpoint, a second provider-specific response shape, or session detail payloads. At that point, prefer OpenAPI or Go-to-TypeScript generation over protobuf for this desktop HTTP API.

## Metadata Extraction Strategy

The list view should use lightweight parsing only.

### Read strategy

- enumerate only `*.jsonl` files representing main sessions
- ignore subagent transcript files in the main list
- read only the minimum content needed to derive list metadata
- prefer reading the first and last useful records instead of loading the whole transcript when possible

### Desired metadata

- session id
- created time
- updated time
- first prompt or best available title seed
- git branch if cheaply available
- approximate or exact message count if cheap
- whether subagent transcripts exist for that session

### Important rule

Do not couple the list MVP to the full transcript parser. A list metadata extractor and a future full importer or session detail loader should remain separate modules. This preserves depth and avoids a shallow "god parser".

## Error Handling

Failures should be isolated and typed.

### Project-level failures

- unreadable directory
- permission denied
- malformed project naming

These should produce a degraded project result or a provider warning, not a full app crash.

### Session-level failures

- unreadable JSONL file
- malformed JSONL lines
- missing expected fields

These should mark that session as unavailable or skipped while preserving the rest of the project list.

### Provider-level failures

- Claude config directory missing
- projects root missing
- filesystem inaccessible

These should map to stable backend errors, then to a user-friendly empty or warning state in the Claude Code tab.

## Testing Strategy

This MVP should be highly testable without needing the real Claude Code app running.

### Backend provider module tests

- project directory discovery
- session file discovery
- ignoring subagent transcripts as top-level sessions
- metadata extraction from representative JSONL fixtures
- corrupt or partial transcript handling

### Backend handler tests

- success response shape for `GET /api/claude-code/sessions`
- provider-level error mapping
- partial failure behavior where the endpoint still returns usable project/session data

### Desktop client and container tests

- grouping sessions under projects
- project sort order
- session sort order
- partial failure behavior
- empty-state behavior

### UI tests

- loading state
- empty state
- grouped project rendering
- nested session row rendering
- warning presentation for degraded results

Use fixtures instead of live reads for most tests. Real Claude Code directories may be used only for local manual verification.

## Design Principles Applied

### SOLID

- Single Responsibility: backend provider discovery, metadata extraction, HTTP transport, and UI rendering each have distinct reasons to change
- Open/Closed: later providers can be added through another backend adapter and endpoint shape without rewriting the Claude Code UI model
- Liskov: session-list contracts stay provider-agnostic enough to support future expansion
- Interface Segregation: list-view consumers receive only list-view data
- Dependency Inversion: the UI depends on a session-list use case, not filesystem details

### DRY

- Claude Code path resolution belongs in one backend provider seam
- JSONL metadata extraction belongs in one focused backend module
- transport and rendering stay thin and reusable

### YAGNI

- no canonical schema in this phase
- no background sync in this phase
- no cross-provider abstractions beyond what this MVP truly needs
- no transcript detail loading in the list view

## Future-Compatible Extensions

This design keeps straightforward seams for later work:

- add an `OpenCode` tab via another backend provider adapter
- add session detail loading when a session row is selected
- add resume actions after a reliable detail and restore flow exists
- add canonical import later without disturbing the list-view UI contract
- add background refresh or watcher support behind the same use case

This also avoids duplicating Claude Code parsing logic in both TypeScript and Go. The canonical pipeline can reuse or extend the same backend knowledge of Claude Code storage.

## Implementation Recommendation

Build this in the following order:

1. Go backend Claude Code provider discovery and metadata extraction
2. Go HTTP handler returning grouped project/session models
3. typed desktop API client
4. desktop container with loading and error states
5. presentational tree UI under the `Claude Code` tab
6. fixture-based tests and local manual verification against real Claude Code storage

## Success Criteria

The MVP is successful when:

- the desktop app displays a `Claude Code` tab
- projects are listed from Claude Code local storage
- sessions are grouped correctly under projects
- broken files do not crash the view
- the UI layer does not contain Claude-specific filesystem logic
- Claude Code list parsing logic lives only in the backend service
- the code structure is ready for future provider expansion without rework
