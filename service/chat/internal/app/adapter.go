package app

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/client/grpc/dto"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/domain"
)

type ChatRepository interface {
	ContainUser(ctx context.Context, chatID string, userID string) bool
	Create(ctx context.Context, chat *domain.Chat) error
	List(ctx context.Context, userID string, p pagination.Pagination) ([]*domain.Chat, uint32, error)
	OneByID(ctx context.Context, id string) (*domain.Chat, error)
}

type MessageRepository interface {
	Create(ctx context.Context, message *domain.Message) error
	List(ctx context.Context, userID string) ([]*domain.Message, int, error)
}

type ChatService interface {
	Create(ctx context.Context, name string, userIDs []string) (string, error)
	CreateMessage(ctx context.Context, text string, chatID string, userInfo domain.UserInfo) (*domain.Message, error)
	List(ctx context.Context, userID string, p pagination.Pagination) ([]*domain.Chat, uint32, error)
	ListMessage(ctx context.Context, chatID string, userID string) ([]*domain.Message, int, error)
	OneByID(ctx context.Context, id string) (*domain.Chat, error)
}

type AuthClient interface {
	GetAccessToken(ctx context.Context, refreshToken string) (accessToken string, err error)
	GetRefreshToken(ctx context.Context, oldRefreshToken string) (refreshToken string, err error)
	Login(ctx context.Context, login string, password string) (refreshToken string, err error)
	Register(ctx context.Context, in dto.RegisterInputDto) (id string, err error)
}

type AccessClient interface {
	Check(ctx context.Context, endpoint string) (access bool, err error)
}

type UserClient interface {
	GetUserInfo(ctx context.Context) (*domain.UserInfo, error)
	GetUserList(ctx context.Context, in *domain.UserInfoListFilter) ([]*domain.UserInfo, uint32, error)
}
