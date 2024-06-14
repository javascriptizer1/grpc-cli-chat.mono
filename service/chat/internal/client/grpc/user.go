package client

import (
	"context"

	userv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/user_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/domain"
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

func (c *UserClient) GetUserList(ctx context.Context, in *domain.UserInfoListFilter) ([]*domain.UserInfo, uint32, error) {
	res, err := c.client.GetUserList(ctx, &userv1.GetUserListRequest{
		Limit:   in.Limit,
		Page:    in.Page,
		UserIDs: in.UserIDs,
	})

	if err != nil {
		return []*domain.UserInfo{}, 0, err
	}

	users := make([]*domain.UserInfo, len(res.GetUsers()))

	for i, v := range res.GetUsers() {
		users[i] = domain.NewUserInfo(
			v.GetId(),
			v.GetName(),
			v.GetEmail(),
			uint16(v.GetRole()),
		)
	}

	return users, res.GetTotal(), nil
}
