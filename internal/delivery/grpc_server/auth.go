package grpc_server

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/repository/dto"
	user "github.com/javascriptizer1/grpc-cli-chat.backend/internal/repository/model"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/authv1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (impl *AuthGrpcImplementation) CreateUser(ctx context.Context, request *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	id, err := impl.userService.Create(ctx, dto.CreateUserDto{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Role:     user.Role(request.Role),
	})

	if err != nil {
		return nil, err
	}

	return &authv1.CreateUserResponse{
		Id: id,
	}, nil
}

func (impl *AuthGrpcImplementation) GetUser(ctx context.Context, request *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
	user, err := impl.userService.One(ctx, request.GetId())

	if err != nil {
		return nil, err
	}

	return &authv1.GetUserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      authv1.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}
