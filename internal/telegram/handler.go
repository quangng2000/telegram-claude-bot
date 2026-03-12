package telegram

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-claude-bot/internal/claude"
	"telegram-claude-bot/internal/session"
)

type Handler struct {
	sender   *Sender
	sessions *session.Manager
	runner   *claude.Runner
}

func NewHandler(sender *Sender, sessions *session.Manager, runner *claude.Runner) *Handler {
	return &Handler{
		sender:   sender,
		sessions: sessions,
		runner:   runner,
	}
}

func (h *Handler) Handle(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	text := msg.Text

	if h.handleCommand(chatID, text) {
		return
	}
	h.handleChat(chatID, text)
}

func (h *Handler) handleChat(chatID int64, text string) {
	if h.sessions.IsExpired(chatID) {
		h.sessions.Expire(chatID, h.runner.RunOneShot)
	}

	h.sessions.LogChat(chatID, "USER", text)
	stop := h.startTyping(chatID)
	defer close(stop)

	start := time.Now()
	log.Printf("Chat %d: \"%s\"", chatID, text)

	prompt := h.buildPrompt(chatID, text)

	mu := h.sessions.GetMutex(chatID)
	mu.Lock()
	chatDir := h.sessions.GetChatDir(chatID)
	response, err := h.runner.Run(chatDir, prompt)
	mu.Unlock()

	elapsed := time.Since(start).Seconds()

	if err != nil {
		log.Printf("Chat %d: error: %s", chatID, err)
		h.sessions.LogChat(chatID, "ERROR", err.Error())
		h.sender.Send(chatID, "Error: "+err.Error())
		return
	}

	log.Printf("Chat %d: responded in %.1fs", chatID, elapsed)
	h.sessions.LogChat(chatID, "BOT", response)
	h.sender.SendMarkdown(chatID, response)
}

func (h *Handler) buildPrompt(chatID int64, text string) string {
	summary := h.sessions.LoadSummary(chatID)
	if strings.TrimSpace(summary) == "" {
		return text
	}
	return fmt.Sprintf("[Previous conversation context:\n%s\n]\n\nUser message: %s",
		summary, text)
}

func (h *Handler) startTyping(chatID int64) chan struct{} {
	h.sender.SendTyping(chatID)
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(4 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				h.sender.SendTyping(chatID)
			}
		}
	}()
	return stop
}
