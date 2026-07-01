// Canonical session schema types
export * from "./claude-code-sessions";

export interface CanonicalSession {
  session_id: string;
  source_provider: "claude-code" | "opencode";
  source_session_id: string;
  model: string;
  working_directory: string;
  started_at: string;
  ended_at: string;
  git_state_start?: { branch: string; commit: string };
  turns: Turn[];
  token_usage?: { input: number; output: number; cost_usd: number };
}

export interface Turn {
  turn_id: string;
  role: "user" | "assistant" | "system";
  timestamp: string;
  content: ContentBlock[];
}

export type ContentBlock =
  | { type: "text"; text: string }
  | { type: "thinking"; text: string; provider_specific?: boolean }
  | ToolCallBlock
  | ToolResultBlock;

export interface ToolCallBlock {
  type: "tool_call";
  tool_call_id: string;
  generic_type: GenericToolType;
  display_summary: string;
  raw_provider_payload: Record<string, unknown>;
}

export interface ToolResultBlock {
  type: "tool_result";
  tool_call_id: string;
  status: "success" | "error";
  content_summary: string;
  raw_provider_payload: Record<string, unknown>;
}

export type GenericToolType =
  | "file_read"
  | "file_write"
  | "file_edit"
  | "shell_exec"
  | "search"
  | "web_fetch"
  | "subagent_spawn"
  | "other";
