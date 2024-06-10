package grpc_server

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/repository/dto"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/authv1"
)

type AuthGrpcImplementation struct {
	authv1.UnimplementedUserServiceServer
	userService UserService
}

type UserService interface {
	Create(context.Context, dto.CreateUserDto) (uint64, error)
	One(context.Context, uint64) (*domain.User, error)
}

func New(userService UserService) *AuthGrpcImplementation {
	return &AuthGrpcImplementation{
		userService: userService,
	}
}
