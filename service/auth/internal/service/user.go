package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/domain"
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

func (s *UserService) List(ctx context.Context, filter *domain.UserListFilter) ([]*domain.User, uint32, error) {
	users, err := s.userRepo.List(ctx, filter)

	if err != nil {
		return []*domain.User{}, 0, err
	}

	total, err := s.userRepo.Count(ctx, filter)

	if err != nil {
		return []*domain.User{}, 0, err
	}

	return users, total, nil

}
