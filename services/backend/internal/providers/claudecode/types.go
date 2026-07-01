package claudecode

import (
	"os"
	"path/filepath"
)

type ClaudeProjectSessions struct {
	ProjectID       string               `json:"projectId"`
	ProjectLabel    string               `json:"projectLabel"`
	ProjectPathHint string               `json:"projectPathHint"`
	SessionCount    int                  `json:"sessionCount"`
	LastUpdatedAt   string               `json:"lastUpdatedAt,omitempty"`
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

type sessionFileV1 struct {
	SessionID string `json:"sessionId"`
	Title     string `json:"title"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	GitBranch string `json:"gitBranch,omitempty"`
	GitRepo   string `json:"gitRepo,omitempty"`
	Error     any    `json:"error"`
}

type messageEntry struct {
	Type      string `json:"type"`
	MessageID string `json:"messageId"`
	Message   struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"message"`
}

type Provider struct {
	configRoot string
}

func NewProvider(root string) *Provider {
	if root == "" {
		home, _ := os.UserHomeDir()
		root = filepath.Join(home, ".claude", "projects")
	}
	return &Provider{configRoot: root}
}

func (p *Provider) DiscoverAllProjects() ([]ClaudeProjectSessions, error) {
	return DiscoverAllProjectsWithRoot(p.configRoot)
}
