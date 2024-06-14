package handler

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/domain"
)

func (h *Handler) GetChat(ctx context.Context, chatID string) (*domain.ChatInfo, error) {
	chat, err := h.chatClient.GetChat(ctx, chatID)

	if err != nil {
		return nil, err
	}

	return chat, nil
}
