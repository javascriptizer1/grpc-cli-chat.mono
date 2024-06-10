package service

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/repository/dto"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/helper/bcrypt"
)

type UserService struct {
	userRepo UserRepository
}

type UserRepository interface {
	Create(context.Context, dto.CreateUserDto) (uint64, error)
	One(context.Context, uint64) (*domain.User, error)
}

func New(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Create(ctx context.Context, input dto.CreateUserDto) (uint64, error) {
	passwordHash, err := bcrypt.Hash(input.Password)

	if err != nil {
		return 0, err
	}

	dto := dto.CreateUserDto{
		Name:     input.Name,
		Email:    input.Email,
		Password: passwordHash,
		Role:     input.Role,
	}

	return s.userRepo.Create(ctx, dto)
}

func (s *UserService) One(ctx context.Context, id uint64) (*domain.User, error) {
	return s.userRepo.One(ctx, id)
}
