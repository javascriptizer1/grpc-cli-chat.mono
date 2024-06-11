package client

import (
	userv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/user_v1"
)

type UserClient struct {
	client userv1.UserServiceClient
}

func NewUserClient(client userv1.UserServiceClient) *UserClient {
	return &UserClient{client: client}
}

// func (c *AccessClient) ChGetUserInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetUserResponse, error) {

// }
