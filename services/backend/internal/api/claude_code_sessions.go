package api

import (
	"encoding/json"
	"net/http"

	"github.com/session-manager/backend/internal/providers/claudecode"
)

type ClaudeProjectSessions struct {
	ProjectID       string                 `json:"projectId"`
	ProjectLabel    string                 `json:"projectLabel"`
	ProjectPathHint string                 `json:"projectPathHint"`
	SessionCount    int                    `json:"sessionCount"`
	LastUpdatedAt   string                 `json:"lastUpdatedAt,omitempty"`
	Sessions        []ClaudeSessionSummary `json:"sessions"`
}

type ClaudeSessionSummary struct {
	SessionID    string `json:"sessionId"`
	Title        string `json:"title"`
	FirstPrompt  string `json:"firstPrompt,omitempty"`
	CreatedAt    string `json:"createdAt,omitempty"`
	UpdatedAt    string `json:"updatedAt,omitempty"`
	MessageCount int    `json:"messageCount,omitempty"`
	GitBranch    string `json:"gitBranch,omitempty"`
	HasSubagents bool   `json:"hasSubagents"`
}

type ClaudeCodeSessionTreeProvider func() (any, error)

func NewClaudeCodeSessionsHandler(provider ClaudeCodeSessionTreeProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		payload, err := provider()
		if err != nil {
			http.Error(w, `{"error":"unable to load Claude Code sessions"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			http.Error(w, `{"error":"unable to load Claude Code sessions"}`, http.StatusInternalServerError)
		}
	})
}

func NewClaudeCodeSessionTreeProvider(projectsRoot string) ClaudeCodeSessionTreeProvider {
	return func() (any, error) {
		provider := claudecode.NewProvider(projectsRoot)
		return provider.DiscoverAllProjects()
	}
}
