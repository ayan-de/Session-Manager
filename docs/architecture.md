# Session Manager — Architecture

## Project Structure

```
session-manager/
├── apps/
│   ├── desktop/                      # Tauri + Next.js (thin native shell + React UI)
│   │   ├── app/                      # Next.js App Router
│   │   ├── src-tauri/                # Rust: window, tray, sidecar spawn
│   │   │   ├── src/main.rs
│   │   │   └── src/lib.rs
│   │   ├── package.json
│   │   └── next.config.ts
│   └── mobile/                       # (future) React Native
│
├── services/
│   └── backend/                      # Go HTTP daemon
│       ├── cmd/backend/main.go       # Entry point
│       └── internal/
│           ├── api/                  # HTTP handlers
│           ├── importer/             # Session file parsers
│           ├── exporter/             # Native format writers
│           ├── taxonomy/             # Generic tool taxonomy mapping
│           ├── storage/              # SQLite (modernc.org/sqlite)
│           └── providers/
│               ├── claudecode/       # Claude Code JSONL parser
│               └── opencode/         # OpenCode SQLite parser
│
├── packages/                          # Shared across desktop + mobile
│   ├── schema/                       # Canonical session schema (TypeScript)
│   ├── api-client/                   # Typed HTTP client → backend
│   └── ui/                           # Shared React components
│
├── hooks/                             # Provider hook scripts
│   ├── claude-code/
│   └── opencode/
│
├── docs/
│   └── architecture.md
│
├── turbo.json
├── pnpm-workspace.yaml
└── package.json
```

## Key Decisions

### Go over Rust for backend
- Backend runs as standalone HTTP daemon, spawned as Tauri sidecar process
- Go chosen for familiarity (exec_runner.go, modernc.org/sqlite patterns from AgentBoard)
- Rust's borrow checker adds friction for AI-generated code; Go compiles cleanly with less debugging
- React Native mobile (future) requires network-callable backend anyway — sidecar architecture is correct regardless

### Tauri stays thin
- `src-tauri/` only handles: window management, system tray, spawning Go backend as sidecar
- All business logic (import, export, taxonomy, storage) lives in Go backend
- This keeps Rust complexity minimal and allows RN to hit the same HTTP API

### Canonical schema as shared package
- `packages/schema` defines the provider-agnostic session format
- Both Go backend (storage/import/export) and TypeScript frontend depend on it
- Future RN mobile imports the same schema — no duplication

### Backend is provider-agnostic
- Does not know about Tauri, desktop, or mobile
- Hooks trigger over HTTP directly — capture works even if no Tauri app is running
- Enables GitHub sync (Phase 1.5) and future remote access from RN

## Build Commands

```bash
# Backend
cd services/backend && go run cmd/backend/main.go   # starts :8080

# Desktop dev
cd apps/desktop && npm run tauri dev

# Frontend only
cd apps/desktop && npm run dev
```

## API (to implement)

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| POST | /import | Import a session file |
| GET | /sessions | List all sessions |
| GET | /sessions/:id | Get session detail |
| POST | /export | Export to provider format |
| POST | /hooks/capture | Called by provider hooks |
