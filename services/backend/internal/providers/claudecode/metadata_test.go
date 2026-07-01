package claudecode

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExtractSessionMetadata(t *testing.T) {
	t.Run("extracts basic metadata from session file", func(t *testing.T) {
		tempDir := t.TempDir()
		sessionFile := filepath.Join(tempDir, "main-abc123.jsonl")

		session := &sessionFileV1{
			SessionID: "abc123",
			CreatedAt: "2024-01-15T10:00:00Z",
			UpdatedAt: "2024-01-15T11:00:00Z",
			Title:     "Test Session",
		}
		data, _ := json.Marshal(session)
		if err := os.WriteFile(sessionFile, data, 0644); err != nil {
			t.Fatalf("writefile failed: %v", err)
		}

		metadata, err := extractSessionMetadata(sessionFile)
		if err != nil {
			t.Fatalf("extractSessionMetadata error: %v", err)
		}

		if metadata.SessionID != "abc123" {
			t.Errorf("SessionID: got %q, want %q", metadata.SessionID, "abc123")
		}
		if metadata.Title != "Test Session" {
			t.Errorf("Title: got %q, want %q", metadata.Title, "Test Session")
		}
	})

	t.Run("derives title from first user message when title is empty", func(t *testing.T) {
		tempDir := t.TempDir()
		sessionFile := filepath.Join(tempDir, "main-abc123.jsonl")

		session := &sessionFileV1{
			SessionID: "abc123",
			CreatedAt: "2024-01-15T10:00:00Z",
			UpdatedAt: "2024-01-15T11:00:00Z",
			Title:     "",
		}
		data, _ := json.Marshal(session)
		if err := os.WriteFile(sessionFile, data, 0644); err != nil {
			t.Fatalf("writefile failed: %v", err)
		}

		metadata, err := extractSessionMetadata(sessionFile)
		if err != nil {
			t.Fatalf("extractSessionMetadata error: %v", err)
		}

		if metadata.Title != "" {
			t.Errorf("Title (before deriving): got %q, want empty string", metadata.Title)
		}
	})

	t.Run("counts messages correctly", func(t *testing.T) {
		tempDir := t.TempDir()
		sessionFile := filepath.Join(tempDir, "main-abc123.jsonl")

		messages := []map[string]interface{}{
			{"type": "user", "message": map[string]interface{}{"content": []interface{}{map[string]interface{}{"type": "text", "text": "Hello"}}}},
			{"type": "assistant", "message": map[string]interface{}{"content": []interface{}{map[string]interface{}{"type": "text", "text": "Hi there!"}}}},
			{"type": "user", "message": map[string]interface{}{"content": []interface{}{map[string]interface{}{"type": "text", "text": "How are you?"}}}},
		}

		session := &sessionFileV1{
			SessionID: "abc123",
			CreatedAt: "2024-01-15T10:00:00Z",
			UpdatedAt: "2024-01-15T11:00:00Z",
		}

		f, _ := os.Create(sessionFile)
		data, _ := json.Marshal(session)
		f.Write(data)
		f.WriteString("\n")
		for _, msg := range messages {
			data, _ := json.Marshal(msg)
			f.Write(data)
			f.WriteString("\n")
		}
		f.Close()

		metadata, err := extractSessionMetadata(sessionFile)
		if err != nil {
			t.Fatalf("extractSessionMetadata error: %v", err)
		}

		if metadata.MessageCount != 3 {
			t.Errorf("MessageCount: got %d, want 3", metadata.MessageCount)
		}
	})

	t.Run("extracts first prompt from messages", func(t *testing.T) {
		tempDir := t.TempDir()
		sessionFile := filepath.Join(tempDir, "main-abc123.jsonl")

		messages := []map[string]interface{}{
			{"type": "user", "message": map[string]interface{}{"content": []interface{}{map[string]interface{}{"type": "text", "text": "Help me build a web server"}}}},
		}

		session := &sessionFileV1{
			SessionID: "abc123",
			CreatedAt: "2024-01-15T10:00:00Z",
			UpdatedAt: "2024-01-15T11:00:00Z",
		}

		f, _ := os.Create(sessionFile)
		data, _ := json.Marshal(session)
		f.Write(data)
		f.WriteString("\n")
		for _, msg := range messages {
			data, _ := json.Marshal(msg)
			f.Write(data)
			f.WriteString("\n")
		}
		f.Close()

		metadata, err := extractSessionMetadata(sessionFile)
		if err != nil {
			t.Fatalf("extractSessionMetadata error: %v", err)
		}

		if metadata.FirstPrompt != "Help me build a web server" {
			t.Errorf("FirstPrompt: got %q, want %q", metadata.FirstPrompt, "Help me build a web server")
		}
	})

	t.Run("handles session file with no messages", func(t *testing.T) {
		tempDir := t.TempDir()
		sessionFile := filepath.Join(tempDir, "main-abc123.jsonl")

		session := &sessionFileV1{
			SessionID: "abc123",
			CreatedAt: "2024-01-15T10:00:00Z",
			UpdatedAt: "2024-01-15T11:00:00Z",
		}
		data, _ := json.Marshal(session)
		if err := os.WriteFile(sessionFile, data, 0644); err != nil {
			t.Fatalf("writefile failed: %v", err)
		}

		metadata, err := extractSessionMetadata(sessionFile)
		if err != nil {
			t.Fatalf("extractSessionMetadata error: %v", err)
		}

		if metadata.FirstPrompt != "" {
			t.Errorf("FirstPrompt: got %q, want empty string", metadata.FirstPrompt)
		}
		if metadata.MessageCount != 0 {
			t.Errorf("MessageCount: got %d, want 0", metadata.MessageCount)
		}
	})
}

func TestGetLastUpdatedAt(t *testing.T) {
	t.Run("returns latest updatedAt from sessions", func(t *testing.T) {
		sessions := []ClaudeSessionSummary{
			{SessionID: "s1", UpdatedAt: "2024-01-10T00:00:00Z"},
			{SessionID: "s2", UpdatedAt: "2024-01-15T00:00:00Z"},
			{SessionID: "s3", UpdatedAt: "2024-01-12T00:00:00Z"},
		}

		latest := getLastUpdatedAt(sessions)

		expected := "2024-01-15T00:00:00Z"
		if latest != expected {
			t.Errorf("latest: got %q, want %q", latest, expected)
		}
	})

	t.Run("returns empty string for empty sessions", func(t *testing.T) {
		sessions := []ClaudeSessionSummary{}

		latest := getLastUpdatedAt(sessions)

		if latest != "" {
			t.Errorf("latest: got %q, want empty string", latest)
		}
	})

	t.Run("handles sessions with no updatedAt", func(t *testing.T) {
		sessions := []ClaudeSessionSummary{
			{SessionID: "s1"},
			{SessionID: "s2", UpdatedAt: "2024-01-15T00:00:00Z"},
		}

		latest := getLastUpdatedAt(sessions)

		expected := "2024-01-15T00:00:00Z"
		if latest != expected {
			t.Errorf("latest: got %q, want %q", latest, expected)
		}
	})
}

func TestLoadSampleSession(t *testing.T) {
	t.Run("loads and parses sample-session.jsonl", func(t *testing.T) {
		samplePath := filepath.Join("testdata", "sample-session.jsonl")

		metadata, err := extractSessionMetadata(samplePath)
		if err != nil {
			t.Fatalf("extractSessionMetadata error: %v", err)
		}

		if metadata.SessionID == "" {
			t.Error("SessionID should not be empty")
		}
		if metadata.Title == "" && metadata.FirstPrompt == "" {
			t.Error("At least one of Title or FirstPrompt should be set")
		}
	})

	t.Run("loads and parses sample-session-with-subagents.jsonl", func(t *testing.T) {
		samplePath := filepath.Join("testdata", "sample-session-with-subagents.jsonl")

		metadata, err := extractSessionMetadata(samplePath)
		if err != nil {
			t.Fatalf("extractSessionMetadata error: %v", err)
		}

		if metadata.SessionID == "" {
			t.Error("SessionID should not be empty")
		}
	})
}

var _ = time.Now // suppress unused warning
