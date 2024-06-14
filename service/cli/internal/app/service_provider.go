package app

import (
	"context"
	"log"

	authv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/auth_v1"
	chatv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/chat_v1"
	userv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/user_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/helper/closer"
	client "github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/client/grpc"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/config"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/handler"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/interceptor"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceProvider struct {
	config *config.Config

	authClient AuthClient
	chatClient ChatClient
	userClient UserClient

	tokenManager TokenManager

	handlerService Handler

	grpcAuthClientConn grpc.ClientConnInterface
	grpcChatClientConn grpc.ClientConnInterface
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) Config() *config.Config {

	if s.config == nil {
		cfg := config.MustLoad()
		s.config = cfg
	}

	return s.config
}

func (s *ServiceProvider) TokenManager(_ context.Context) TokenManager {

	if s.tokenManager == nil {
		tm := manager.NewFileTokenManager(s.Config().ClientConfigPath())

		s.tokenManager = tm
	}

	return s.tokenManager
}

func (s *ServiceProvider) GRPCAuthClientConn(ctx context.Context) grpc.ClientConnInterface {

	if s.grpcAuthClientConn == nil {

		authInterceptor := interceptor.NewAuthInterceptor(s.TokenManager(ctx))

		conn, err := grpc.NewClient(
			s.Config().GRPCAuth.HostPort(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(authInterceptor.Unary()),
		)

		if err != nil {
			log.Fatalf("failed to connect %s: %s", s.Config().GRPCAuth.HostPort(), err.Error())
		}

		s.grpcAuthClientConn = conn

		closer.Add(conn.Close)
	}

	return s.grpcAuthClientConn
}

func (s *ServiceProvider) GRPCChatClientConn(ctx context.Context) grpc.ClientConnInterface {

	if s.grpcChatClientConn == nil {

		authInterceptor := interceptor.NewAuthInterceptor(s.TokenManager(ctx))

		conn, err := grpc.NewClient(
			s.Config().GRPCChat.HostPort(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(authInterceptor.Unary()),
			grpc.WithChainStreamInterceptor(authInterceptor.Stream()),
		)

		if err != nil {
			log.Fatalf("failed to connect %s: %s", s.Config().GRPCChat.HostPort(), err.Error())
		}

		s.grpcChatClientConn = conn

		closer.Add(conn.Close)
	}

	return s.grpcChatClientConn
}

func (s *ServiceProvider) AuthClient(ctx context.Context) AuthClient {
	if s.authClient == nil {
		s.authClient = client.NewAuthClient(authv1.NewAuthServiceClient(s.GRPCAuthClientConn(ctx)))
	}

	return s.authClient
}

func (s *ServiceProvider) UserClient(ctx context.Context) UserClient {
	if s.userClient == nil {
		s.userClient = client.NewUserClient(userv1.NewUserServiceClient(s.GRPCAuthClientConn(ctx)))
	}

	return s.userClient
}

func (s *ServiceProvider) ChatClient(ctx context.Context) ChatClient {
	if s.chatClient == nil {
		s.chatClient = client.NewChatClient(chatv1.NewChatServiceClient(s.GRPCChatClientConn(ctx)))
	}

	return s.chatClient
}

func (s *ServiceProvider) HandlerService(ctx context.Context) Handler {
	if s.handlerService == nil {
		s.handlerService = handler.New(
			s.AuthClient(ctx),
			s.UserClient(ctx),
			s.ChatClient(ctx),
			s.TokenManager(ctx),
		)
	}

	return s.handlerService
}
