package grpc

import (
	"context"

	authv1 "github.com/javascriptizer1/grpc-cli-chat.mono/pkg/grpc/auth_v1"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/service/dto"
)

type AuthImplementation struct {
	authv1.UnimplementedAuthServiceServer
	authService AuthService
}

func NewGrpcAuthImplementation(authService AuthService) *AuthImplementation {
	return &AuthImplementation{
		authService: authService,
	}
}

func (impl *AuthImplementation) Register(ctx context.Context, request *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {

	id, err := impl.authService.Register(ctx, dto.RegisterInputDto{
		Name:            request.GetName(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		PasswordConfirm: request.GetPasswordConfirm(),
		Role:            domain.UserRole(request.GetRole()),
	})

	if err != nil {
		return nil, err
	}

	return &authv1.RegisterResponse{
		Id: id.String(),
	}, nil
}

func (impl *AuthImplementation) Login(ctx context.Context, request *authv1.LoginRequest) (*authv1.LoginResponse, error) {

	refreshToken, err := impl.authService.Login(ctx, request.GetLogin(), request.GetPassword())

	if err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}

func (impl *AuthImplementation) GetRefreshToken(ctx context.Context, request *authv1.GetRefreshTokenRequest) (*authv1.GetRefreshTokenResponse, error) {

	refreshToken, err := impl.authService.RefreshToken(ctx, request.GetOldRefreshToken())

	if err != nil {
		return nil, err
	}

	return &authv1.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}

func (impl *AuthImplementation) GetAccessToken(ctx context.Context, request *authv1.GetAccessTokenRequest) (*authv1.GetAccessTokenResponse, error) {

	accessToken, err := impl.authService.AccessToken(ctx, request.GetRefreshToken())

	if err != nil {
		return nil, err
	}

	return &authv1.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
