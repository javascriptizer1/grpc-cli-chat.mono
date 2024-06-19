package client

import (
	"context"

	accessv1 "github.com/javascriptizer1/grpc-cli-chat.mono/pkg/grpc/access_v1"
)

type AccessClient struct {
	client accessv1.AccessServiceClient
}

func NewAccessClient(client accessv1.AccessServiceClient) *AccessClient {
	return &AccessClient{client: client}
}

func (c *AccessClient) Check(ctx context.Context, endpoint string) (bool, error) {
	_, err := c.client.Check(ctx, &accessv1.CheckRequest{EndpointAddress: endpoint})

	return err == nil, err
}
