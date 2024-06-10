package usergrpc

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (impl *GrpcUserImplementation) GetUserInfo(ctx context.Context, request *emptypb.Empty) (*user_v1.GetUserResponse, error) {
	payload, ok := ctx.Value("payload").(jwt.Claims)

	if !ok {
		return nil, status.Errorf(codes.Internal, "missing required token")
	}

	subject, err := payload.GetSubject()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "extract subject error")
	}

	u, err := impl.userService.OneById(ctx, uuid.MustParse(subject))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
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
