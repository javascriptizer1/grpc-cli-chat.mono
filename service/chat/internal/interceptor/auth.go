package interceptor

import (
	"context"

	client "github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/client/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type authInterceptor struct {
	accessClient client.AccessClient
}

func NewAuthInterceptor(accessClient client.AccessClient) *authInterceptor {
	return &authInterceptor{
		accessClient: accessClient,
	}
}

func (i *authInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		ok = i.accessClient.Check(ctx, info.FullMethod)

		if !ok {
			return nil, err
		}

		return handler(ctx, req)
	}
}
