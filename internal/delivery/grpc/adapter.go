package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/service/dto"
)

type AuthService interface {
	Register(ctx context.Context, input dto.RegisterInputDto) (uuid.UUID, error)
	Login(ctx context.Context, login string, password string) (string, error)
	RefreshToken(ctx context.Context, oldRefreshToken string) (string, error)
	AccessToken(ctx context.Context, refreshToken string) (string, error)
	Check(ctx context.Context, endpoint string, role domain.UserRole) bool
}

type UserService interface {
	OneByID(context.Context, uuid.UUID) (*domain.User, error)
}
