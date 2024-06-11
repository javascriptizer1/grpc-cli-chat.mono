package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain"
)

type UserRepository interface {
	Create(context.Context, *domain.User) error
	OneByID(context.Context, uuid.UUID) (*domain.User, error)
	OneByEmail(context.Context, string) (*domain.User, error)
}
