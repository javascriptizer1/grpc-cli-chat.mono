package authsvc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain/user"
)

type AuthConfig struct {
	AccessTokenSecretKey  string
	AccessTokenDuration   time.Duration
	RefreshTokenSecretKey string
	RefreshTokenDuration  time.Duration
}

type AuthService struct {
	userRepo UserRepository
	config   AuthConfig
}

type UserRepository interface {
	Create(context.Context, *user.User) error
	OneById(context.Context, uuid.UUID) (*user.User, error)
	OneByEmail(context.Context, string) (*user.User, error)
}

func New(userRepo UserRepository, cfg AuthConfig) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}
