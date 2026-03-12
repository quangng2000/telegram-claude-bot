package session

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Manager struct {
	botDir      string
	sessionsDir string
	expiryHours int
	mutexes     sync.Map
}

func NewManager(botDir, sessionsDir string, expiryHours int) *Manager {
	return &Manager{
		botDir:      botDir,
		sessionsDir: sessionsDir,
		expiryHours: expiryHours,
	}
}

func (m *Manager) ExpiryHours() int {
	return m.expiryHours
}

func (m *Manager) GetChatDir(chatID int64) string {
	chatDir := filepath.Join(m.sessionsDir, fmt.Sprintf("%d", chatID))
	if _, err := os.Stat(chatDir); os.IsNotExist(err) {
		os.MkdirAll(chatDir, 0755)
		m.symlinkClaude(chatDir)
	}
	return chatDir
}

func (m *Manager) symlinkClaude(chatDir string) {
	src := filepath.Join(m.botDir, "CLAUDE.md")
	dst := filepath.Join(chatDir, "CLAUDE.md")
	if _, err := os.Stat(src); err == nil {
		os.Symlink(src, dst)
	}
}

func (m *Manager) Reset(chatID int64) {
	chatDir := filepath.Join(m.sessionsDir, fmt.Sprintf("%d", chatID))
	os.RemoveAll(chatDir)
}

func (m *Manager) GetMutex(chatID int64) *sync.Mutex {
	v, _ := m.mutexes.LoadOrStore(chatID, &sync.Mutex{})
	return v.(*sync.Mutex)
}

func (m *Manager) IsExpired(chatID int64) bool {
	logFile := filepath.Join(m.sessionsDir, fmt.Sprintf("%d", chatID), "history.log")
	info, err := os.Stat(logFile)
	if err != nil {
		return false
	}
	hours := time.Since(info.ModTime()).Hours()
	return hours >= float64(m.expiryHours)
}

func (m *Manager) GetSummaryFile(chatID int64) string {
	return filepath.Join(m.sessionsDir, fmt.Sprintf("%d", chatID), "previous-summary.md")
}

func (m *Manager) LoadSummary(chatID int64) string {
	data, err := os.ReadFile(m.GetSummaryFile(chatID))
	if err != nil {
		return ""
	}
	return string(data)
}
