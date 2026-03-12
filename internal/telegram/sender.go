package telegram

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender struct {
	bot *tgbotapi.BotAPI
}

func NewSender(bot *tgbotapi.BotAPI) *Sender {
	return &Sender{bot: bot}
}

func (s *Sender) Send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	s.bot.Send(msg)
}

func (s *Sender) SendMarkdown(chatID int64, text string) {
	for _, chunk := range SplitMessage(text, 4000) {
		msg := tgbotapi.NewMessage(chatID, chunk)
		msg.ParseMode = "Markdown"
		if _, err := s.bot.Send(msg); err != nil {
			msg.ParseMode = ""
			s.bot.Send(msg)
		}
	}
}

func (s *Sender) SendTyping(chatID int64) {
	action := tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping)
	s.bot.Send(action)
}

func SplitMessage(text string, maxLen int) []string {
	if len(text) <= maxLen {
		return []string{text}
	}
	var chunks []string
	remaining := text
	for len(remaining) > 0 {
		if len(remaining) <= maxLen {
			chunks = append(chunks, remaining)
			break
		}
		idx := strings.LastIndex(remaining[:maxLen], "\n")
		if idx == -1 || idx < maxLen/2 {
			idx = maxLen
		}
		chunks = append(chunks, remaining[:idx])
		remaining = strings.TrimLeft(remaining[idx:], "\n ")
	}
	return chunks
}
