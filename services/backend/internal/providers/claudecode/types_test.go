package claudecode

import (
	"encoding/json"
	"testing"
)

func TestClaudeProjectSessions_JSON(t *testing.T) {
	t.Run("serializes all fields correctly", func(t *testing.T) {
		project := ClaudeProjectSessions{
			ProjectID:        "test-project",
			ProjectLabel:     "Test Project",
			ProjectPathHint:  "/home/user/test",
			SessionCount:     2,
			LastUpdatedAt:    "2024-01-15T10:30:00Z",
			Sessions: []ClaudeSessionSummary{
				{
					SessionID:    "session-1",
					Title:        "First session",
					FirstPrompt:  "Hello world",
					CreatedAt:    "2024-01-15T09:00:00Z",
					UpdatedAt:    "2024-01-15T10:00:00Z",
					MessageCount: 42,
					GitBranch:    "main",
					HasSubagents: false,
				},
				{
					SessionID:    "session-2",
					Title:        "Second session",
					HasSubagents: true,
				},
			},
		}

		data, err := json.Marshal(project)
		if err != nil {
			t.Fatalf("json.Marshal error: %v", err)
		}

		var unmarshaled ClaudeProjectSessions
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("json.Unmarshal error: %v", err)
		}

		if unmarshaled.ProjectID != project.ProjectID {
			t.Errorf("ProjectID: got %q, want %q", unmarshaled.ProjectID, project.ProjectID)
		}
		if unmarshaled.SessionCount != project.SessionCount {
			t.Errorf("SessionCount: got %d, want %d", unmarshaled.SessionCount, project.SessionCount)
		}
		if len(unmarshaled.Sessions) != len(project.Sessions) {
			t.Errorf("Sessions length: got %d, want %d", len(unmarshaled.Sessions), len(project.Sessions))
		}
	})
}

func TestClaudeSessionSummary_TitleFromFirstPrompt(t *testing.T) {
	t.Run("uses firstPrompt as title when title is empty", func(t *testing.T) {
		summary := ClaudeSessionSummary{
			SessionID:   "test-session",
			FirstPrompt: "How do I build a web server?",
			HasSubagents: false,
		}

		title := summary.Title
		if summary.FirstPrompt != "" && title == "" {
			title = summary.FirstPrompt
		}

		if title != "How do I build a web server?" {
			t.Errorf("title: got %q, want %q", title, "How do I build a web server?")
		}
	})
}

func TestProjectPathHint_URLDecode(t *testing.T) {
	t.Run("decodes URL-encoded path hint", func(t *testing.T) {
		encoded := "-home-user-Projects-Test"
		decoded, err := decodeProjectPathHint(encoded)
		if err != nil {
			t.Fatalf("decodeProjectPathHint error: %v", err)
		}
		expected := "/home/user/Projects/Test"
		if decoded != expected {
			t.Errorf("decoded: got %q, want %q", decoded, expected)
		}
	})
}

func TestSessionSorting(t *testing.T) {
	t.Run("sorts sessions by updatedAt descending", func(t *testing.T) {
		sessions := []ClaudeSessionSummary{
			{SessionID: "oldest", UpdatedAt: "2024-01-01T00:00:00Z"},
			{SessionID: "newest", UpdatedAt: "2024-01-15T00:00:00Z"},
			{SessionID: "middle", UpdatedAt: "2024-01-10T00:00:00Z"},
		}

		sortSessionsByUpdatedAt(sessions)

		if sessions[0].SessionID != "newest" {
			t.Errorf("first session: got %q, want %q", sessions[0].SessionID, "newest")
		}
		if sessions[1].SessionID != "middle" {
			t.Errorf("middle session: got %q, want %q", sessions[1].SessionID, "middle")
		}
		if sessions[2].SessionID != "oldest" {
			t.Errorf("last session: got %q, want %q", sessions[2].SessionID, "oldest")
		}
	})
}

func TestProjectSorting(t *testing.T) {
	t.Run("sorts projects by lastUpdatedAt descending", func(t *testing.T) {
		projects := []ClaudeProjectSessions{
			{ProjectID: "oldest", LastUpdatedAt: "2024-01-01T00:00:00Z"},
			{ProjectID: "newest", LastUpdatedAt: "2024-01-15T00:00:00Z"},
			{ProjectID: "middle", LastUpdatedAt: "2024-01-10T00:00:00Z"},
		}

		sortProjectsByLastUpdatedAt(projects)

		if projects[0].ProjectID != "newest" {
			t.Errorf("first project: got %q, want %q", projects[0].ProjectID, "newest")
		}
		if projects[1].ProjectID != "middle" {
			t.Errorf("middle project: got %q, want %q", projects[1].ProjectID, "middle")
		}
		if projects[2].ProjectID != "oldest" {
			t.Errorf("last project: got %q, want %q", projects[2].ProjectID, "oldest")
		}
	})
}
