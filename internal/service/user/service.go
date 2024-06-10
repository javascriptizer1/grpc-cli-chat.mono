package usersvc

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain/user"
)

type UserService struct {
	userRepo UserRepository
}

type UserRepository interface {
	Create(context.Context, *user.User) error
	OneById(context.Context, uuid.UUID) (*user.User, error)
}

func New(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}
