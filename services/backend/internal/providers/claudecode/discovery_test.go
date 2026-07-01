package claudecode

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverProjects(t *testing.T) {
	t.Run("discovers projects from config root", func(t *testing.T) {
		tempDir := t.TempDir()
		projectsDir := filepath.Join(tempDir, "projects")
		if err := os.MkdirAll(projectsDir, 0755); err != nil {
			t.Fatalf("mkdir failed: %v", err)
		}

		project1 := filepath.Join(projectsDir, "-home-user-Projects-App1")
		project2 := filepath.Join(projectsDir, "-home-user-Projects-App2")
		for _, p := range []string{project1, project2} {
			if err := os.MkdirAll(p, 0755); err != nil {
				t.Fatalf("mkdir failed: %v", err)
			}
		}

		projects, err := discoverProjects(projectsDir)
		if err != nil {
			t.Fatalf("discoverProjects error: %v", err)
		}

		if len(projects) != 2 {
			t.Errorf("project count: got %d, want 2", len(projects))
		}
	})

	t.Run("skips non-directory entries", func(t *testing.T) {
		tempDir := t.TempDir()
		projectsDir := filepath.Join(tempDir, "projects")
		if err := os.MkdirAll(projectsDir, 0755); err != nil {
			t.Fatalf("mkdir failed: %v", err)
		}

		projectDir := filepath.Join(projectsDir, "-home-user-Projects-App")
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			t.Fatalf("mkdir failed: %v", err)
		}

		file := filepath.Join(projectsDir, "not-a-project")
		if err := os.WriteFile(file, []byte("nope"), 0644); err != nil {
			t.Fatalf("writefile failed: %v", err)
		}

		projects, err := discoverProjects(projectsDir)
		if err != nil {
			t.Fatalf("discoverProjects error: %v", err)
		}

		if len(projects) != 1 {
			t.Errorf("project count: got %d, want 1", len(projects))
		}
	})

	t.Run("empty config root returns empty list", func(t *testing.T) {
		tempDir := t.TempDir()

		projects, err := discoverProjects(tempDir)
		if err != nil {
			t.Fatalf("discoverProjects error: %v", err)
		}

		if len(projects) != 0 {
			t.Errorf("project count: got %d, want 0", len(projects))
		}
	})
}

func TestDiscoverSessions(t *testing.T) {
	t.Run("finds main session jsonl files", func(t *testing.T) {
		tempDir := t.TempDir()
		projectDir := filepath.Join(tempDir, "-home-user-Projects-App")

		sessionsDir := filepath.Join(projectDir, "sessions")
		if err := os.MkdirAll(sessionsDir, 0755); err != nil {
			t.Fatalf("mkdir failed: %v", err)
		}

		mainSession := filepath.Join(sessionsDir, "main-abc123.jsonl")
		subagentSession := filepath.Join(sessionsDir, "subagent-def456.jsonl")

		for _, f := range []string{mainSession, subagentSession} {
			if err := os.WriteFile(f, []byte("{}"), 0644); err != nil {
				t.Fatalf("writefile failed: %v", err)
			}
		}

		sessions, err := discoverSessions(projectDir)
		if err != nil {
			t.Fatalf("discoverSessions error: %v", err)
		}

		if len(sessions) != 1 {
			t.Errorf("session count: got %d, want 1", len(sessions))
		}
	})

	t.Run("detects subagents directory", func(t *testing.T) {
		tempDir := t.TempDir()
		projectDir := filepath.Join(tempDir, "-home-user-Projects-App")

		sessionsDir := filepath.Join(projectDir, "sessions")
		if err := os.MkdirAll(sessionsDir, 0755); err != nil {
			t.Fatalf("mkdir failed: %v", err)
		}
		subagentsDir := filepath.Join(projectDir, "subagents")
		if err := os.MkdirAll(subagentsDir, 0755); err != nil {
			t.Fatalf("mkdir failed: %v", err)
		}

		mainSession := filepath.Join(sessionsDir, "main-abc123.jsonl")
		if err := os.WriteFile(mainSession, []byte("{}"), 0644); err != nil {
			t.Fatalf("writefile failed: %v", err)
		}

		hasSubagents := checkHasSubagents(projectDir)

		if !hasSubagents {
			t.Error("hasSubagents: got false, want true")
		}
	})

	t.Run("no subagents directory means no subagents", func(t *testing.T) {
		tempDir := t.TempDir()
		projectDir := filepath.Join(tempDir, "-home-user-Projects-App")

		sessionsDir := filepath.Join(projectDir, "sessions")
		if err := os.MkdirAll(sessionsDir, 0755); err != nil {
			t.Fatalf("mkdir failed: %v", err)
		}

		mainSession := filepath.Join(sessionsDir, "main-abc123.jsonl")
		if err := os.WriteFile(mainSession, []byte("{}"), 0644); err != nil {
			t.Fatalf("writefile failed: %v", err)
		}

		hasSubagents := checkHasSubagents(projectDir)

		if hasSubagents {
			t.Error("hasSubagents: got true, want false")
		}
	})
}

func TestResolveConfigRoot(t *testing.T) {
	t.Run("uses default path when root is empty", func(t *testing.T) {
		provider := NewProvider("")
		if provider.configRoot == "" {
			t.Error("configRoot should not be empty with default provider")
		}
	})

	t.Run("uses custom path when provided", func(t *testing.T) {
		provider := NewProvider("/custom/path")
		if provider.configRoot != "/custom/path" {
			t.Errorf("configRoot: got %q, want %q", provider.configRoot, "/custom/path")
		}
	})
}
