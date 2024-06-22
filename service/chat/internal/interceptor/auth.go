package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type serverStreamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss *serverStreamWithContext) Context() context.Context {
	return ss.ctx
}

type AccessClient interface {
	Check(ctx context.Context, endpoint string) (bool, error)
}

type authInterceptor struct {
	accessClient AccessClient
}

func NewAuthInterceptor(accessClient AccessClient) *authInterceptor {
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

		ok, err = i.accessClient.Check(ctx, info.FullMethod)

		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, status.Errorf(codes.PermissionDenied, "access denied")
		}

		return handler(ctx, req)
	}
}

func (i *authInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		ok, err := i.accessClient.Check(ctx, info.FullMethod)

		if err != nil {
			return err
		}

		if !ok {
			return status.Errorf(codes.PermissionDenied, "access denied")
		}

		wrappedStream := &serverStreamWithContext{ss, ctx}

		return handler(srv, wrappedStream)
	}
}
