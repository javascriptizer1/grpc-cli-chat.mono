package usergrpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain/user"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/user_v1"
)

type GrpcUserImplementation struct {
	user_v1.UnimplementedUserServiceServer
	userService UserService
}

type UserService interface {
	OneById(context.Context, uuid.UUID) (*user.User, error)
}

func New(userService UserService) *GrpcUserImplementation {
	return &GrpcUserImplementation{
		userService: userService,
	}
}
