package claudecode

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func extractSessionMetadata(sessionFilePath string) (ClaudeSessionSummary, error) {
	file, err := os.Open(sessionFilePath)
	if err != nil {
		return ClaudeSessionSummary{}, err
	}
	defer file.Close()

	var summary ClaudeSessionSummary
	var firstPrompt string
	var messageCount int
	var session sessionFileV1
	var skippedLines int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(string(scanner.Bytes()))
		if line == "" {
			continue
		}

		var base map[string]any
		if err := json.Unmarshal([]byte(line), &base); err != nil {
			skippedLines++
			continue
		}

		if sessionFile, ok := base["sessionId"]; ok && sessionFile != nil {
			if err := json.Unmarshal([]byte(line), &session); err == nil {
				summary.SessionID = session.SessionID
				summary.Title = session.Title
				summary.CreatedAt = session.CreatedAt
				summary.UpdatedAt = session.UpdatedAt
				summary.GitBranch = session.GitBranch
			}
			continue
		}

		if msgType, ok := base["type"].(string); ok && (msgType == "user" || msgType == "assistant" || msgType == "subagent") {
			messageCount++
			if firstPrompt == "" && msgType == "user" {
				if msg, ok := base["message"].(map[string]any); ok {
					if content, ok := msg["content"].([]any); ok && len(content) > 0 {
						if firstItem, ok := content[0].(map[string]any); ok {
							if text, ok := firstItem["text"].(string); ok {
								firstPrompt = text
							}
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return ClaudeSessionSummary{}, err
	}

	if skippedLines > 0 {
		log.Printf("debug: skipped %d malformed JSON lines in %s", skippedLines, sessionFilePath)
	}

	if summary.Title == "" && firstPrompt != "" {
		summary.Title = firstPrompt
	}
	summary.FirstPrompt = firstPrompt
	summary.MessageCount = messageCount

	return summary, nil
}

func getSessionMetadata(sessionFilePath string) (ClaudeSessionSummary, error) {
	return extractSessionMetadata(sessionFilePath)
}

func loadSessionsFromTestdata() ([]ClaudeSessionSummary, error) {
	var summaries []ClaudeSessionSummary

	samplePath := filepath.Join("testdata", "sample-session.jsonl")
	metadata, err := extractSessionMetadata(samplePath)
	if err == nil {
		summaries = append(summaries, metadata)
	}

	subagentPath := filepath.Join("testdata", "sample-session-with-subagents.jsonl")
	metadata, err = extractSessionMetadata(subagentPath)
	if err == nil {
		metadata.HasSubagents = true
		summaries = append(summaries, metadata)
	}

	return summaries, nil
}
