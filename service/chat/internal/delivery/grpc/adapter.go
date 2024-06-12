package grpc

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/domain"
)

type ChatService interface {
	Create(ctx context.Context, userIDs []string) (string, error)
	OneByID(ctx context.Context, id string) (*domain.Chat, error)
	CreateMessage(ctx context.Context, text string, chatID string, userInfo domain.UserInfo) (*domain.Message, error)
}

type UserClient interface {
	GetUserInfo(ctx context.Context) (*domain.UserInfo, error)
}
