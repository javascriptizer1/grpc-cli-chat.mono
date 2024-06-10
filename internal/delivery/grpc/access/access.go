package accessgrpc

import (
	"context"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain/user"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/access_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (impl *GrpcAccessImplementation) Check(ctx context.Context, request *access_v1.CheckRequest) (*emptypb.Empty, error) {

	payload, ok := ctx.Value("payload").(jwt.Claims)

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

	access := impl.authService.Check(ctx, request.GetEndpointAddress(), user.Role(role))

	if !access {
		return &emptypb.Empty{}, status.Errorf(codes.PermissionDenied, "access denied")
	}

	return &emptypb.Empty{}, nil
}
