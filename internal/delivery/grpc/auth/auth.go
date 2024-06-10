package authgrpc

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain/user"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/service/auth/dto"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/auth_v1"
)

func (impl *GrpcAuthImplementation) Register(ctx context.Context, request *auth_v1.RegisterRequest) (*auth_v1.RegisterResponse, error) {

	id, err := impl.authService.Register(ctx, dto.RegisterInputDto{
		Name:            request.GetName(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		PasswordConfirm: request.GetPasswordConfirm(),
		Role:            user.Role(request.GetRole()),
	})

	if err != nil {
		return nil, err
	}

	return &auth_v1.RegisterResponse{
		Id: id.String(),
	}, nil
}

func (impl *GrpcAuthImplementation) Login(ctx context.Context, request *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {

	refreshToken, err := impl.authService.Login(ctx, request.GetLogin(), request.GetPassword())

	if err != nil {
		return nil, err
	}

	return &auth_v1.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}

func (impl *GrpcAuthImplementation) GetRefreshToken(ctx context.Context, request *auth_v1.GetRefreshTokenRequest) (*auth_v1.GetRefreshTokenResponse, error) {

	refreshToken, err := impl.authService.RefreshToken(ctx, request.GetOldRefreshToken())

	if err != nil {
		return nil, err
	}

	return &auth_v1.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}

func (impl *GrpcAuthImplementation) GetAccessToken(ctx context.Context, request *auth_v1.GetAccessTokenRequest) (*auth_v1.GetAccessTokenResponse, error) {

	accessToken, err := impl.authService.AccessToken(ctx, request.GetRefreshToken())

	if err != nil {
		return nil, err
	}

	return &auth_v1.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
