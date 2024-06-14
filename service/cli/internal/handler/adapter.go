package handler

import (
	"context"

	chatv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/chat_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/client/grpc/dto"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/domain"
)

type TokenManager interface {
	Save() error
	SetTokens(accessToken, refreshToken string) error
	AccessToken() string
	RefreshToken() string
}

type AuthClient interface {
	GetAccessToken(ctx context.Context, refreshToken string) (accessToken string, err error)
	GetRefreshToken(ctx context.Context, oldRefreshToken string) (refreshToken string, err error)
	Login(ctx context.Context, login string, password string) (refreshToken string, err error)
	Register(ctx context.Context, in dto.RegisterInputDto) (id string, err error)
}

type UserClient interface {
	GetUserList(ctx context.Context, options *domain.UserListOption) ([]*domain.UserInfo, uint32, error)
}

type ChatClient interface {
	ConnectChat(ctx context.Context, chatID string) (cha chatv1.ChatService_ConnectChatClient, err error)
	CreateChat(ctx context.Context, emails []string) (id string, err error)
	SendMessage(ctx context.Context, text string, chatID string) error
	GetChatList(ctx context.Context, p *pagination.Pagination) ([]*domain.ChatListInfo, uint32, error)
	GetChat(ctx context.Context, id string) (*domain.ChatInfo, error)
}
