package grpc

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	userv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/user_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/domain"
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

func (impl *UserImplementation) GetUserInfo(ctx context.Context, _ *emptypb.Empty) (*userv1.User, error) {
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

	return &userv1.User{
		Id:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		Role:      userv1.Role(u.Role),
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}, nil
}

func (impl *UserImplementation) GetUserList(ctx context.Context, request *userv1.GetUserListRequest) (*userv1.GetUserListResponse, error) {
	uuids := make([]uuid.UUID, len(request.GetUserIDs()))

	for i, v := range request.UserIDs {
		uuids[i] = uuid.MustParse(v)
	}

	users, total, err := impl.userService.List(ctx, &domain.UserListFilter{
		Pagination: *pagination.New(request.GetLimit(), request.GetPage()),
		UserIDs:    uuids,
	})

	if err != nil {
		return nil, err
	}

	protoUsers := make([]*userv1.User, len(users))

	for i, v := range users {
		protoUsers[i] = &userv1.User{
			Id:        v.ID.String(),
			Name:      v.Name,
			Email:     v.Email,
			Role:      userv1.Role(v.Role),
			CreatedAt: timestamppb.New(v.CreatedAt),
			UpdatedAt: timestamppb.New(v.UpdatedAt),
		}
	}

	return &userv1.GetUserListResponse{
		Users: protoUsers,
		Total: total,
	}, nil
}
