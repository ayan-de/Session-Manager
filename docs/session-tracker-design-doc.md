# Session Tracker — Design Doc

## 1. What this is

target Platform => Linux (Arch + Derbian), Windows, Mac

A tool that indexes AI coding agent sessions (Claude Code, OpenCode, later Codex/Antigravity) into a **provider-agnostic canonical format**, so we can:

1. List/browse all sessions across tools in one UI (like a unified session log)
2. Convert a session from one provider's native format into another provider's native format, so it can be **resumed** in the target tool
3. Hide provider-specific tool names/schemas in the UI — show generic actions like "a file was read" instead of leaking `Read` vs whatever OpenCode calls it

**MVP scope: Claude Code ↔ OpenCode only.** Codex and Antigravity come later, once the core pipeline is proven.

## 2. Core architecture

```
[Claude Code JSONL] ──┐
                       ├─→ [Importer] → [Canonical Schema in DB] → [Exporter] → [Target native file]
[OpenCode SQLite]  ────┘                                                              │
                                                                                       ▼
                                                                          claude --resume <id>
                                                                          (or opencode equivalent)
```

**Why this shape, not N×N direct converters:** with 2 providers now and up to 4 later, we only ever need N importers + N exporters (import once → canonical → export many), not a converter per pair. This is the entire point of the canonical layer — it's the actual IP of the project, not the UI.

### Import is the easy, safe part
Both Claude Code and OpenCode already persist **complete, full-fidelity session data locally** — every message, tool call input/output, thinking block, timestamps, model used, etc. We are not trying to reconstruct anything lossy; we're parsing files that already have everything.

### Export/resume is the part that needs validation before we build UI around it
Resume, for Claude Code, works by replaying a well-formed message transcript back to the Anthropic Messages API (`claude --resume <session-id>` reads the JSONL and continues the conversation). This means export is really "produce a schema-valid transcript," which is bounded and testable — not a black box. But:

- If the original session was driven by a **different model family** (e.g. OpenCode running GPT/Grok), its tool_use blocks are shaped for that model's API. Converting to Claude's tool schema is real translation work, not just relabeling.
- If both sessions were run against the same underlying model family, cross-tool resume is much more tractable.

**Action item before deep UI work:** spike this. Hand-write/generate a synthetic Claude Code JSONL from a real session and confirm `claude --resume` actually accepts and continues it. Do the OpenCode-side equivalent too. This determines how much of "resume" is realistic for MVP vs a v2 goal.

## 3. Phase 1.5 — GitHub sync (same-provider, build this first)

Before tackling cross-provider conversion, build the much lower-risk version: **same-provider, cross-machine resume**, modeled on how LeetHub auto-pushes solved code to GitHub. If the session is Claude Code on PC A and you want to resume it as Claude Code on PC B, no canonical schema or taxonomy mapping is needed — it's the same file format on both ends. This validates the whole hook → capture → storage → restore pipeline end-to-end with zero conversion risk, and gives something usable almost immediately.

### Flow
1. **Push**: on session-end (the same hook trigger from Section 5), copy the raw session file as-is to a GitHub repo (private by default, like LeetHub) — e.g. `sessions/claude-code/<project>/<session-id>.jsonl`, `sessions/opencode/<project>/<session-id>.db-export`
2. **Pull**: on the target machine, `session-tracker pull <session-id>` fetches the file from GitHub and writes it back into the provider's expected local path (e.g. `~/.claude/projects/<encoded-path>/<session-id>.jsonl`), then prints the resume command (`claude --resume <session-id>`)

### Why build this before cross-provider export
- Zero conversion/taxonomy risk — it's a straight file copy, not a translation
- Proves the capture → transport → restore pipeline works before adding conversion complexity on top
- Delivers a usable feature (resume any session on any machine) even if cross-provider conversion turns out to be harder than expected

### Secrets/privacy risk — must handle before auto-pushing every session
Unlike LeetHub, which pushes clean solved code, a session transcript is raw context — it can contain API keys, tokens, `.env` contents, credentials typed into commands, internal file paths, etc. "Private repo" alone isn't sufficient (repos can be made public later, org access can change, tokens can leak). Before enabling auto-push by default:
- Add a redaction pass for common secret patterns (API key formats, `.env` file contents, cloud provider tokens) before commit
- Make push opt-in per session or per project rather than auto-push-everything
- Consider making push an explicit action rather than an automatic `SessionEnd` hook action, at least for v1

## 4. Capture mechanism — no MCP tool, use hooks + file parsing

**Important correction from earlier planning:** we are NOT building an MCP tool for capture. MCP tools are invoked by the model during a conversation; they can't fire on session-end lifecycle events. What we want instead:

- **Hooks** (Claude Code: `SessionEnd`, `PreCompact`; OpenCode: whatever equivalent lifecycle hook exists) are shell commands defined in the tool's own config, and their only job is to **trigger our importer** — e.g. `session-tracker import --session-id <id> --provider claude-code`.
- The importer itself reads the raw file **directly off disk** (JSONL for Claude Code, SQLite for OpenCode) independent of the hook. A file-watcher/poller is a viable alternative to hooks entirely, since the data is already fully persisted by the tool itself.
- MCP may be added *later* as a read/query interface (letting an agent query the session archive on demand), but that's a separate feature, unrelated to capture.

### Known file locations (verify against current versions before building)
- Claude Code: `~/.claude/projects/<url-encoded-project-path>/<session-id>.jsonl` — one JSON object per line (user/assistant messages, tool_use, tool_result, thinking blocks). Subagent transcripts live in a sibling `subagents/` folder per session.
- OpenCode: local SQLite schema — already documented from the token tracker project; reuse that schema knowledge directly.

## 5. Canonical schema (draft — refine during implementation)

```jsonc
{
  "session_id": "uuid",
  "source_provider": "claude-code | opencode",
  "source_session_id": "original id in source tool",
  "model": "string",
  "working_directory": "string",
  "started_at": "iso8601",
  "ended_at": "iso8601",
  "git_state_start": { "branch": "...", "commit": "..." },
  "turns": [
    {
      "turn_id": "uuid",
      "role": "user | assistant | system",
      "timestamp": "iso8601",
      "content": [
        { "type": "text", "text": "..." },
        { "type": "thinking", "text": "...", "provider_specific": true },
        {
          "type": "tool_call",
          "tool_call_id": "uuid",
          "generic_type": "file_read | file_write | file_edit | shell_exec | search | web_fetch | subagent_spawn | other",
          "display_summary": "a file was read: src/index.ts",
          "raw_provider_payload": { "...": "verbatim original block, for lossless round-trip" }
        },
        {
          "type": "tool_result",
          "tool_call_id": "uuid",
          "status": "success | error",
          "content_summary": "...",
          "raw_provider_payload": { "...": "verbatim original" }
        }
      ]
    }
  ],
  "token_usage": { "input": 0, "output": 0, "cost_usd": 0 }
}
```

Key design decision: **store both** the generic taxonomy mapping (for UI display / cross-provider meaning) **and** the raw original payload (for lossless export back to the same provider, and as the source-of-truth input to translation logic when exporting to a *different* provider).

## 6. Generic tool taxonomy (starting set, expand as needed)

| Generic type | Claude Code | OpenCode |
|---|---|---|
| `file_read` | `Read` | (map from OpenCode's equivalent) |
| `file_write` | `Write` | ... |
| `file_edit` | `Edit` | ... |
| `shell_exec` | `Bash` | ... |
| `search` | `Grep`, `Glob` | ... |
| `web_fetch` | `WebFetch` | ... |
| `subagent_spawn` | `Task` | ... |

This table is the crux of the whole project — filling it in accurately for both providers is the main implementation task for the importer/exporter pair.

## 7. Components to build (MVP)

**Build order: Phase 1.5 (GitHub sync) first, then the canonical pipeline.**

Phase 1.5 — GitHub sync:
1. **Hook scripts** — thin wrappers registered in Claude Code's `settings.json` (`SessionEnd`/`PreCompact`) that trigger push. Equivalent for OpenCode.
2. **`session-tracker-push`** — copies the raw session file to a private GitHub repo, running the redaction pass first.
3. **`session-tracker-pull`** — fetches a session file from GitHub and writes it to the correct local path for the provider, prints the resume command.

Phase 2 — Canonical pipeline (cross-provider):
4. **`session-tracker-importer`** — CLI/daemon. Given a provider + session file path, parses it into the canonical schema and writes to local DB (SQLite to start).
5. **`session-tracker-exporter`** — given a canonical session + target provider, writes a native-format file. Start with Claude Code as export target since its format is best understood and resume is just an Anthropic Messages API replay.
6. **Resume validation spike** (do this *before* UI polish) — confirm exported files are actually accepted by `claude --resume`.
7. **UI** — session list (cross-provider), session detail view rendering generic tool-call summaries, a "transfer to..." action that runs the exporter.

## 8. Explicit non-goals for MVP

- Codex, Antigravity support (defer)
- MCP query interface (defer — separate feature)
- Perfect fidelity for cross-model-family resume (defer — flag as best-effort, may lose some tool-call semantics until taxonomy mapping matures)
- Real-time/streaming capture (session-end trigger is sufficient for v1)

## 9. Open risks to keep validating as we build

- Whether OpenCode's own resume mechanism will accept a synthetically-written session (same class of risk as Claude Code, needs its own spike)
- Tool schema drift — Claude Code / OpenCode tool definitions can change between versions; taxonomy mapping needs to be versioned, not hardcoded once
- Claude Code auto-deletes old session files over time — if we rely on re-reading the file post-hoc rather than capturing at hook-time, we need to grab it promptly, not assume it'll still be there later
- Secret/credential leakage via pushed session files (see Section 3) — redaction pass must be in place and tested before any default-on auto-push behavior ships
