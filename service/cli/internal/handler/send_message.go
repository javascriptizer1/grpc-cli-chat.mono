package handler

import "context"

func (h *Handler) SendMessage(ctx context.Context, text string, chatID string) error {
	return h.chatClient.SendMessage(ctx, text, chatID)
}
