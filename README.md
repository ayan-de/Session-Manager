# Session Manager

AI coding agent session tracker — import, browse, and transfer sessions across providers.

## Structure

```
session-manager/
├── apps/
│   ├── desktop/          # Tauri app (thin shell + React frontend)
│   └── mobile/           # React Native (future)
├── services/
│   └── backend/          # Go HTTP daemon — importer, exporter, SQLite storage
├── packages/
│   ├── schema/           # Canonical session schema (TypeScript)
│   ├── api-client/       # Typed HTTP client
│   └── ui/               # Shared React components
├── hooks/                # Provider hook scripts (Claude Code, OpenCode)
└── docs/
    ├── architecture.md     # This document
    └── session-tracker-design-doc.md
```

## Dev

```bash
# Backend
cd services/backend && go run cmd/backend/main.go

# Desktop (from root)
cd apps/desktop && npm run tauri dev
```
