package session

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SummaryFunc func(workDir, prompt string) (string, error)

func (m *Manager) CheckAllExpired(summarize SummaryFunc) {
	entries, err := os.ReadDir(m.sessionsDir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		var chatID int64
		fmt.Sscanf(entry.Name(), "%d", &chatID)
		if chatID != 0 && m.IsExpired(chatID) {
			m.Expire(chatID, summarize)
		}
	}
}

func (m *Manager) Expire(chatID int64, summarize SummaryFunc) {
	chatDir := filepath.Join(m.sessionsDir, fmt.Sprintf("%d", chatID))
	if _, err := os.Stat(chatDir); os.IsNotExist(err) {
		return
	}

	log.Printf("Chat %d: session expired, generating summary...", chatID)

	history := m.GetHistoryLog(chatID)
	if len(strings.TrimSpace(history)) < 50 {
		m.Reset(chatID)
		return
	}

	summary := m.generateSummary(chatID, history, summarize)
	historyBackup := m.backupHistory(chatID)

	m.Reset(chatID)
	m.restoreWithSummary(chatID, historyBackup, summary)

	log.Printf("Chat %d: summary saved, session reset", chatID)
}

func (m *Manager) generateSummary(chatID int64, history string, summarize SummaryFunc) string {
	recentHistory := m.GetRecentHistory(chatID, 100)

	prompt := "Summarize the following conversation into a concise context brief (max 500 words). " +
		"Focus on: key topics discussed, important decisions made, system states observed, any ongoing issues. " +
		"This summary will be used to provide context for future conversations.\n\nConversation:\n" + recentHistory

	summary, err := summarize(m.botDir, prompt)
	if err != nil {
		log.Printf("Chat %d: summary failed: %s", chatID, err)
		return ""
	}
	return m.mergeSummary(chatID, summary)
}

func (m *Manager) mergeSummary(chatID int64, summary string) string {
	existing, _ := os.ReadFile(m.GetSummaryFile(chatID))
	merged := fmt.Sprintf("## Session ended: %s\n%s\n\n%s",
		time.Now().Format(time.RFC3339), summary, string(existing))

	const maxLen = 3000
	if len(merged) > maxLen {
		merged = merged[:maxLen] + "\n...(older history truncated)"
	}
	return merged
}

func (m *Manager) backupHistory(chatID int64) []byte {
	chatDir := filepath.Join(m.sessionsDir, fmt.Sprintf("%d", chatID))
	data, _ := os.ReadFile(filepath.Join(chatDir, "history.log"))
	return data
}

func (m *Manager) restoreWithSummary(chatID int64, history []byte, summary string) {
	chatDir := m.GetChatDir(chatID)
	os.WriteFile(filepath.Join(chatDir, "history.log"), history, 0644)
	if summary != "" {
		os.WriteFile(m.GetSummaryFile(chatID), []byte(summary), 0644)
	}
}
