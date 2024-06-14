package grpc

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/domain"
)

type ChatService interface {
	Create(ctx context.Context, name string, userIDs []string) (string, error)
	OneByID(ctx context.Context, id string) (*domain.Chat, error)
	CreateMessage(ctx context.Context, text string, chatID string, userInfo domain.UserInfo) (*domain.Message, error)
	List(ctx context.Context, userID string, p pagination.Pagination) ([]*domain.Chat, uint32, error)
}

type UserClient interface {
	GetUserInfo(ctx context.Context) (*domain.UserInfo, error)
	GetUserList(ctx context.Context, in *domain.UserInfoListFilter) ([]*domain.UserInfo, uint32, error)
}
