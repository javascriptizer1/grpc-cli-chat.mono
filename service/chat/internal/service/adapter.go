package service

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/domain"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *domain.Chat) error
	List(ctx context.Context, userID string) ([]*domain.Chat, error)
	ContainUser(ctx context.Context, chatID string, userID string) bool
	OneByID(ctx context.Context, id string) (*domain.Chat, error)
}

type MessageRepository interface {
	Create(ctx context.Context, message *domain.Message) error
	List(ctx context.Context, chatID string) ([]*domain.Message, int, error)
}

type UserClient interface {
	GetUserInfo(ctx context.Context) (*domain.UserInfo, error)
}
