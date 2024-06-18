package client

import (
	"context"

	userv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/user_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/bubble/internal/domain"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserClient struct {
	client userv1.UserServiceClient
}

func NewUserClient(client userv1.UserServiceClient) *UserClient {
	return &UserClient{client: client}
}

func (c *UserClient) GetUserInfo(ctx context.Context) (*domain.UserInfo, error) {
	res, err := c.client.GetUserInfo(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	ui := domain.NewUserInfo(
		res.GetId(),
		res.GetName(),
		res.GetEmail(),
		uint16(res.GetRole()),
	)

	return ui, nil

}

func (c *UserClient) GetUserList(ctx context.Context, options *domain.UserListOption) ([]*domain.UserInfo, uint32, error) {
	res, err := c.client.GetUserList(ctx, &userv1.GetUserListRequest{
		Limit:   options.Limit(),
		Page:    options.Page(),
		UserIDs: options.UserIDs,
	})

	var users []*domain.UserInfo

	for _, v := range res.GetUsers() {
		ui := domain.NewUserInfo(
			v.Id,
			v.Name,
			v.Email,
			uint16(v.Role),
		)
		users = append(users, ui)
	}

	if err != nil {
		return []*domain.UserInfo{}, 0, err
	}

	return users, res.GetTotal(), nil

}
