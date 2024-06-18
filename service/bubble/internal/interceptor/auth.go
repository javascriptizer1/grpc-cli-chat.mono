package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	tokenManager TokenManager
}

func NewAuthInterceptor(tokenManager TokenManager) *AuthInterceptor {
	return &AuthInterceptor{tokenManager: tokenManager}
}

func (a *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		accessToken := a.tokenManager.AccessToken()

		if accessToken != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "Bearer "+accessToken)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (a *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		accessToken := a.tokenManager.AccessToken()

		if accessToken != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "Bearer "+accessToken)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}
