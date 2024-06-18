package client

import (
	"context"

	authv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/auth_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/bubble/internal/client/grpc/dto"
)

type AuthClient struct {
	client authv1.AuthServiceClient
}

func NewAuthClient(client authv1.AuthServiceClient) *AuthClient {
	return &AuthClient{client: client}
}

func (c *AuthClient) Register(ctx context.Context, in dto.RegisterInputDto) (id string, err error) {
	res, err := c.client.Register(ctx, &authv1.RegisterRequest{
		Name:            in.Name,
		Email:           in.Email,
		Password:        in.Password,
		PasswordConfirm: in.PasswordConfirm,
		Role:            authv1.Role(in.Role),
	})

	if err != nil {
		return id, err
	}

	return res.GetId(), nil
}

func (c *AuthClient) Login(ctx context.Context, login string, password string) (refreshToken string, err error) {
	res, err := c.client.Login(ctx, &authv1.LoginRequest{
		Login:    login,
		Password: password,
	})

	if err != nil {
		return refreshToken, err
	}

	return res.GetRefreshToken(), nil
}

func (c *AuthClient) GetRefreshToken(ctx context.Context, oldRefreshToken string) (refreshToken string, err error) {
	res, err := c.client.GetRefreshToken(ctx, &authv1.GetRefreshTokenRequest{OldRefreshToken: oldRefreshToken})

	if err != nil {
		return refreshToken, err
	}

	return res.GetRefreshToken(), nil
}

func (c *AuthClient) GetAccessToken(ctx context.Context, refreshToken string) (accessToken string, err error) {
	res, err := c.client.GetAccessToken(ctx, &authv1.GetAccessTokenRequest{RefreshToken: refreshToken})

	if err != nil {
		return accessToken, err
	}

	return res.GetAccessToken(), nil
}
