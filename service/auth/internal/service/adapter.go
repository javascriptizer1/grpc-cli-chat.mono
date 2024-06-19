package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/domain"
)

type UserRepository interface {
	Create(context.Context, *domain.User) error
	OneByID(context.Context, uuid.UUID) (*domain.User, error)
	OneByEmail(context.Context, string) (*domain.User, error)
	List(ctx context.Context, filter *domain.UserListFilter) ([]*domain.User, error)
	Count(ctx context.Context, filter *domain.UserListFilter) (total uint32, err error)
}
