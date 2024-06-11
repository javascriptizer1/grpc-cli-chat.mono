package interceptor

import (
	"context"
	"strings"

	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/helper/ctxkey"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/helper/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const ContextKeyUserClaims = ctxkey.ContextKey("userClaims")

type AuthInterceptor struct {
	accessTokenSecret string
}

func NewAuthInterceptor(accessTokenSecret string) *AuthInterceptor {
	return &AuthInterceptor{accessTokenSecret: accessTokenSecret}
}

func (i *AuthInterceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if !methodRequiresAuthentication(info.FullMethod) {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "metadata not found")
	}

	authTokens := md["authorization"]

	if len(authTokens) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is missing")
	}

	authHeader := authTokens[0]
	fields := strings.Fields(authHeader)

	if len(fields) < 2 {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth header format: %v", fields)
	}

	authType := strings.ToLower(fields[0])

	if authType != "bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization type: %v", authType)
	}

	t := fields[1]

	payload, err := jwt.VerifyToken(t, i.accessTokenSecret)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token %v", err)
	}

	ctx = context.WithValue(ctx, ContextKeyUserClaims, payload)

	return handler(ctx, req)
}

func methodRequiresAuthentication(fullMethod string) bool {
	m := extractMethodName(fullMethod)
	authRequiredMethods := []string{
		"user_v1.UserService/GetUserInfo",
		"access_v1.AccessService/Check",
	}

	for _, method := range authRequiredMethods {
		if m == method {
			return true
		}
	}

	return false
}

func extractMethodName(fullMethod string) string {
	method := strings.TrimLeft(fullMethod, "/")

	return method
}
