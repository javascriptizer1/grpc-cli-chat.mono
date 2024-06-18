package handler

import (
	"context"

	chatv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/chat_v1"
)

func (h *Handler) ConnectChat(ctx context.Context, chatID string) (chatv1.ChatService_ConnectChatClient, error) {
	return h.chatClient.ConnectChat(ctx, chatID)
}
