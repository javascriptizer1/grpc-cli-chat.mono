package accessgrpc

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain/user"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/access_v1"
)

type GrpcAccessImplementation struct {
	access_v1.UnimplementedAccessServiceServer
	authService AuthService
}

type AuthService interface {
	Check(ctx context.Context, endpoint string, role user.Role) bool
}

func New(authService AuthService) *GrpcAccessImplementation {
	return &GrpcAccessImplementation{
		authService: authService,
	}
}
