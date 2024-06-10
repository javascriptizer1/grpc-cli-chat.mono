package authgrpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/service/auth/dto"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/auth_v1"
)

type GrpcAuthImplementation struct {
	auth_v1.UnimplementedAuthServiceServer
	authService AuthService
}

type AuthService interface {
	Register(ctx context.Context, input dto.RegisterInputDto) (uuid.UUID, error)
	Login(ctx context.Context, login string, password string) (string, error)
	RefreshToken(ctx context.Context, oldRefreshToken string) (string, error)
	AccessToken(ctx context.Context, refreshToken string) (string, error)
}

func New(authService AuthService) *GrpcAuthImplementation {
	return &GrpcAuthImplementation{
		authService: authService,
	}
}
