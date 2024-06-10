package usergrpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (impl *GrpcUserImplementation) GetUser(ctx context.Context, request *user_v1.GetUserRequest) (*user_v1.GetUserResponse, error) {

	u, err := impl.userService.OneById(ctx, uuid.MustParse(request.GetId()))

	if err != nil {
		return nil, err
	}

	return &user_v1.GetUserResponse{
		Id:        u.Id.String(),
		Name:      u.Name,
		Email:     u.Email,
		Role:      user_v1.Role(u.Role),
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}, nil
}
