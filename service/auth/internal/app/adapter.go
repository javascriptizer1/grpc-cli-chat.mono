package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/service/dto"
)

type UserRepository interface {
	Create(_ context.Context, input *domain.User) error
	OneByEmail(_ context.Context, email string) (u *domain.User, err error)
	OneByID(_ context.Context, id uuid.UUID) (u *domain.User, err error)
}

type AuthService interface {
	AccessToken(ctx context.Context, refreshToken string) (string, error)
	Check(_ context.Context, endpoint string, role domain.UserRole) bool
	Login(ctx context.Context, login string, password string) (string, error)
	RefreshToken(ctx context.Context, oldRefreshToken string) (string, error)
	Register(ctx context.Context, input dto.RegisterInputDto) (id uuid.UUID, err error)
}

type UserService interface {
	OneByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}
