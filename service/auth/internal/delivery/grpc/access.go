package grpc

import (
	"context"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	accessv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/access_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/interceptor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AccessImplementation struct {
	accessv1.UnimplementedAccessServiceServer
	authService AuthService
}

func NewGrpcAccessImplementation(authService AuthService) *AccessImplementation {
	return &AccessImplementation{
		authService: authService,
	}
}

func (impl *AccessImplementation) Check(ctx context.Context, request *accessv1.CheckRequest) (*emptypb.Empty, error) {

	payload, ok := ctx.Value(interceptor.ContextKeyUserClaims).(jwt.Claims)

	if !ok {
		return nil, status.Errorf(codes.Internal, "missing required token")
	}

	audience, err := payload.GetAudience()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "extract role error")
	}

	role, err := strconv.Atoi(audience[0])

	if err != nil {
		return nil, status.Errorf(codes.Internal, "extract role error")
	}

	access := impl.authService.Check(ctx, request.GetEndpointAddress(), domain.UserRole(role))

	if !access {
		return &emptypb.Empty{}, status.Errorf(codes.PermissionDenied, "access denied")
	}

	return &emptypb.Empty{}, nil
}
