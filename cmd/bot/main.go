package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-claude-bot/internal/claude"
	"telegram-claude-bot/internal/config"
	"telegram-claude-bot/internal/session"
	"telegram-claude-bot/internal/telegram"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatal(err)
	}

	runner := claude.NewRunner(cfg.ClaudeModel, cfg.CmdTimeout)
	sessions := session.NewManager(cfg.BotDir, cfg.SessionsDir, cfg.ExpiryHours)
	sender := telegram.NewSender(bot)
	handler := telegram.NewHandler(sender, sessions, runner)

	log.Printf("Bot started: @%s (expiry: %dh, model: %s)",
		bot.Self.UserName, cfg.ExpiryHours, cfg.ClaudeModel)

	go runExpiryLoop(sessions, runner)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range bot.GetUpdatesChan(u) {
		if update.Message == nil || update.Message.Text == "" {
			continue
		}
		go handler.Handle(update.Message)
	}
}

func runExpiryLoop(sessions *session.Manager, runner *claude.Runner) {
	for {
		time.Sleep(30 * time.Minute)
		sessions.CheckAllExpired(runner.RunOneShot)
	}
}
