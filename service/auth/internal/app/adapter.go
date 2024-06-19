package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/service/dto"
)

type UserRepository interface {
	Create(_ context.Context, input *domain.User) error
	OneByEmail(_ context.Context, email string) (u *domain.User, err error)
	OneByID(_ context.Context, id uuid.UUID) (u *domain.User, err error)
	List(ctx context.Context, filter *domain.UserListFilter) ([]*domain.User, error)
	Count(ctx context.Context, filter *domain.UserListFilter) (total uint32, err error)
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
	List(ctx context.Context, filter *domain.UserListFilter) ([]*domain.User, uint32, error)
}
