package grpc

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	userv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/user_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/interceptor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserImplementation struct {
	userv1.UnimplementedUserServiceServer
	userService UserService
}

func NewGrpcUserImplementation(userService UserService) *UserImplementation {
	return &UserImplementation{
		userService: userService,
	}
}

func (impl *UserImplementation) GetUserInfo(ctx context.Context, _ *emptypb.Empty) (*userv1.GetUserResponse, error) {
	payload, ok := ctx.Value(interceptor.ContextKeyUserClaims).(jwt.Claims)

	if !ok {
		return nil, status.Errorf(codes.Internal, "missing required token")
	}

	subject, err := payload.GetSubject()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "extract subject error")
	}

	u, err := impl.userService.OneByID(ctx, uuid.MustParse(subject))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &userv1.GetUserResponse{
		Id:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		Role:      userv1.Role(u.Role),
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}, nil
}
