package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Token        string
	BotDir       string
	SessionsDir  string
	ExpiryHours  int
	ClaudeModel  string
	CmdTimeout   int // minutes
}

func Load() (*Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN not set")
	}

	botDir := os.Getenv("BOT_DIR")
	if botDir == "" {
		botDir, _ = os.Getwd()
	}

	cfg := &Config{
		Token:       token,
		BotDir:      botDir,
		SessionsDir: filepath.Join(botDir, "sessions"),
		ExpiryHours: 24,
		ClaudeModel: "sonnet",
		CmdTimeout:  5,
	}

	if v := os.Getenv("SESSION_EXPIRY_HOURS"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.ExpiryHours)
	}
	if v := os.Getenv("CLAUDE_MODEL"); v != "" {
		cfg.ClaudeModel = v
	}

	os.MkdirAll(cfg.SessionsDir, 0755)
	return cfg, nil
}
