package telegram

import "fmt"

func (h *Handler) handleCommand(chatID int64, text string) bool {
	switch text {
	case "/start":
		h.sender.Send(chatID, h.startMessage())
		return true
	case "/reset":
		h.sessions.Reset(chatID)
		h.sender.Send(chatID, "Conversation fully reset.")
		return true
	case "/history":
		h.sendHistory(chatID)
		return true
	}
	return false
}

func (h *Handler) sendHistory(chatID int64) {
	recent := h.sessions.GetRecentHistory(chatID, 20)
	if recent == "" {
		h.sender.Send(chatID, "No history yet.")
		return
	}
	h.sender.SendMarkdown(chatID, "Recent history:\n\n"+recent)
}

func (h *Handler) startMessage() string {
	return "Hey! I'm your personal AI assistant. I can:\n" +
		"- Answer questions on any topic\n" +
		"- Help with coding, writing, and research\n" +
		"- Check and manage your K8S cluster\n" +
		"- Query and modify your PostgreSQL database\n" +
		"- Monitor services and deployments\n\n" +
		"Just ask me anything!\n\n" +
		"Commands:\n" +
		"/reset - Clear conversation and start fresh\n" +
		"/history - Show recent chat history\n" +
		fmt.Sprintf("Sessions expire after %dh of inactivity",
			h.sessions.ExpiryHours())
}
