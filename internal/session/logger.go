package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (m *Manager) LogChat(chatID int64, role, text string) {
	chatDir := m.GetChatDir(chatID)
	logFile := filepath.Join(chatDir, "history.log")

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	ts := time.Now().Format(time.RFC3339)
	fmt.Fprintf(f, "[%s] %s: %s\n", ts, role, text)
}

func (m *Manager) GetHistoryLog(chatID int64) string {
	logFile := filepath.Join(m.sessionsDir, fmt.Sprintf("%d", chatID), "history.log")
	data, err := os.ReadFile(logFile)
	if err != nil {
		return ""
	}
	return string(data)
}

func (m *Manager) GetRecentHistory(chatID int64, lines int) string {
	history := m.GetHistoryLog(chatID)
	if history == "" {
		return ""
	}
	all := strings.Split(strings.TrimSpace(history), "\n")
	start := 0
	if len(all) > lines {
		start = len(all) - lines
	}
	return strings.Join(all[start:], "\n")
}
