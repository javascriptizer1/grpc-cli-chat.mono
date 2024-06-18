package handler

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/bubble/internal/domain"
)

func (h *Handler) GetChatList(ctx context.Context, p *pagination.Pagination) ([]*domain.ChatListInfo, uint32, error) {
	chats, total, err := h.chatClient.GetChatList(ctx, p)

	if err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}
