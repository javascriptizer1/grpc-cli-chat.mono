package app

import (
	"context"

	chatv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/chat_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/client/grpc/dto"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/domain"
)

type AuthClient interface {
	GetAccessToken(ctx context.Context, refreshToken string) (accessToken string, err error)
	GetRefreshToken(ctx context.Context, oldRefreshToken string) (refreshToken string, err error)
	Login(ctx context.Context, login string, password string) (refreshToken string, err error)
	Register(ctx context.Context, in dto.RegisterInputDto) (id string, err error)
}

type ChatClient interface {
	ConnectChat(ctx context.Context, chatID string) (cha chatv1.ChatService_ConnectChatClient, err error)
	CreateChat(ctx context.Context, emails []string) (id string, err error)
	SendMessage(ctx context.Context, text string, chatID string) error
}

type UserClient interface {
	GetUserInfo(ctx context.Context) (*domain.UserInfo, error)
}
