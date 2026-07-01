package claudecode

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func discoverProjects(configRoot string) ([]ClaudeProjectSessions, error) {
	entries, err := os.ReadDir(configRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return []ClaudeProjectSessions{}, nil
		}
		return nil, err
	}

	var projects []ClaudeProjectSessions
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		projectPathHint := entry.Name()

		decodedPath, err := decodeProjectPathHint(projectPathHint)
		if err != nil {
			continue
		}

		label := filepath.Base(decodedPath)
		project := ClaudeProjectSessions{
			ProjectID:       projectPathHint,
			ProjectLabel:    label,
			ProjectPathHint: decodedPath,
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func decodeProjectPathHint(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	decoded := encoded
	if decoded[0] == '-' {
		decoded = "/" + decoded[1:]
	}
	decoded = strings.ReplaceAll(decoded, "--", "\x00")
	decoded = strings.ReplaceAll(decoded, "-", "/")
	decoded = strings.ReplaceAll(decoded, "\x00", "-")
	return decoded, nil
}

func discoverSessions(projectPath string) ([]string, error) {
	entries, err := os.ReadDir(projectPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var sessionFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".jsonl") {
			sessionFiles = append(sessionFiles, filepath.Join(projectPath, name))
		}
	}
	return sessionFiles, nil
}

func checkHasSubagents(projectPath string) bool {
	subagentsDir := filepath.Join(projectPath, "subagents")
	_, err := os.Stat(subagentsDir)
	return err == nil
}

func DiscoverAllProjects() ([]ClaudeProjectSessions, error) {
	home, _ := os.UserHomeDir()
	return DiscoverAllProjectsWithRoot(filepath.Join(home, ".claude", "projects"))
}

func DiscoverAllProjectsWithRoot(root string) ([]ClaudeProjectSessions, error) {
	configRootPath := root
	if configRootPath == "" {
		home, _ := os.UserHomeDir()
		configRootPath = filepath.Join(home, ".claude", "projects")
	}

	projects, err := discoverProjects(configRootPath)
	if err != nil {
		return nil, err
	}

	for i := range projects {
		projectDir := filepath.Join(configRootPath, projects[i].ProjectID)
		hasSubagents := checkHasSubagents(projectDir)

		sessions, err := discoverSessions(projectDir)
		if err != nil {
			continue
		}

		summaries := []ClaudeSessionSummary{}
		for _, sessionFile := range sessions {
			metadata, err := extractSessionMetadata(sessionFile)
			if err != nil {
				continue
			}
			metadata.HasSubagents = hasSubagents
			summaries = append(summaries, metadata)
		}

		sortSessionsByUpdatedAt(summaries)
		projects[i].Sessions = summaries
		projects[i].SessionCount = len(summaries)
		if len(summaries) > 0 {
			projects[i].LastUpdatedAt = summaries[0].UpdatedAt
		}
	}

	sortProjectsByLastUpdatedAt(projects)

	return projects, nil
}

func ExportProjectsToJSON() ([]byte, error) {
	projects, err := DiscoverAllProjects()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(projects, "", "  ")
}

func sortSessionsByUpdatedAt(sessions []ClaudeSessionSummary) {
	for i := 0; i < len(sessions)-1; i++ {
		for j := i + 1; j < len(sessions); j++ {
			if compareUpdatedAt(sessions[j].UpdatedAt, sessions[i].UpdatedAt) > 0 {
				sessions[i], sessions[j] = sessions[j], sessions[i]
			}
		}
	}
}

func sortProjectsByLastUpdatedAt(projects []ClaudeProjectSessions) {
	for i := 0; i < len(projects)-1; i++ {
		for j := i + 1; j < len(projects); j++ {
			if compareUpdatedAt(projects[j].LastUpdatedAt, projects[i].LastUpdatedAt) > 0 {
				projects[i], projects[j] = projects[j], projects[i]
			}
		}
	}
}

func compareUpdatedAt(a, b string) int {
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return -1
	}
	if b == "" {
		return 1
	}
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}

func getLastUpdatedAt(sessions []ClaudeSessionSummary) string {
	var latest string
	for _, s := range sessions {
		if s.UpdatedAt > latest {
			latest = s.UpdatedAt
		}
	}
	return latest
}
