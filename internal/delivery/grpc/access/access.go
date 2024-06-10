package accessgrpc

import (
	"context"
	"errors"

	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/access_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (impl *GrpcAccessImplementation) Check(ctx context.Context, request *access_v1.CheckRequest) (*emptypb.Empty, error) {

	ok := impl.authService.Check(ctx, request.GetEndpointAddress(), 1)

	if !ok {
		return &emptypb.Empty{}, errors.New("access denied")
	}

	return &emptypb.Empty{}, nil
}
