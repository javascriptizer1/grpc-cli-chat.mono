package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) OneByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.userRepo.OneByID(ctx, id)
}
