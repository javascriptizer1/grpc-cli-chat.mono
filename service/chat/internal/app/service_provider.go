package app

import (
	"context"

	mongoClient "github.com/javascriptizer1/grpc-cli-chat.mono/pkg/client/mongo"
	accessv1 "github.com/javascriptizer1/grpc-cli-chat.mono/pkg/grpc/access_v1"
	authv1 "github.com/javascriptizer1/grpc-cli-chat.mono/pkg/grpc/auth_v1"
	userv1 "github.com/javascriptizer1/grpc-cli-chat.mono/pkg/grpc/user_v1"
	"github.com/javascriptizer1/grpc-cli-chat.mono/pkg/helper/closer"
	grpcClient "github.com/javascriptizer1/grpc-cli-chat.mono/service/chat/internal/client/grpc"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/chat/internal/config"
	delivery "github.com/javascriptizer1/grpc-cli-chat.mono/service/chat/internal/delivery/grpc"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/chat/internal/logger"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/chat/internal/repository"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/chat/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serviceProvider struct {
	config *config.Config

	grpcClientConn grpc.ClientConnInterface

	dbClient *mongo.Database

	chatRepo    ChatRepository
	messageRepo MessageRepository

	chatService ChatService

	authClient   AuthClient
	accessClient AccessClient
	userClient   UserClient

	chatImplementation *delivery.ChatImplementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) Config() *config.Config {

	if s.config == nil {
		cfg := config.MustLoad()
		s.config = cfg
	}

	return s.config
}

func (s *serviceProvider) GRPCClientConn() grpc.ClientConnInterface {

	if s.grpcClientConn == nil {
		conn, err := grpc.NewClient(
			s.Config().GRPCAuth.HostPort(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err != nil {
			logger.Fatal("failed to connect "+s.Config().GRPCAuth.HostPort(), zap.String("err", err.Error()))
		}

		s.grpcClientConn = conn

		closer.Add(conn.Close)
	}

	return s.grpcClientConn
}

func (s *serviceProvider) MongoClient(ctx context.Context) *mongo.Database {
	if s.dbClient == nil {

		db, err := mongoClient.New(ctx, mongoClient.Config{
			Host:     s.Config().DB.Host,
			Port:     s.Config().DB.Port,
			User:     s.Config().DB.User,
			Password: s.Config().DB.Password,
			Name:     s.Config().DB.Name,
		})

		if err != nil {
			logger.Fatal("failed to start up db", zap.Error(err))
		}

		s.dbClient = db
	}

	return s.dbClient
}

func (s *serviceProvider) ChatRepository(ctx context.Context) ChatRepository {

	if s.chatRepo == nil {
		s.chatRepo = repository.NewChatRepository(s.MongoClient(ctx))
	}

	return s.chatRepo
}

func (s *serviceProvider) MessageRepository(ctx context.Context) MessageRepository {

	if s.messageRepo == nil {
		s.messageRepo = repository.NewMessageRepository(s.MongoClient(ctx))
	}

	return s.messageRepo
}

func (s *serviceProvider) ChatService(ctx context.Context) ChatService {

	if s.chatService == nil {

		s.chatService = service.NewChatService(
			s.ChatRepository(ctx),
			s.MessageRepository(ctx),
			s.UserClient(ctx),
		)
	}

	return s.chatService
}

func (s *serviceProvider) AuthClient(_ context.Context) AuthClient {

	if s.authClient == nil {
		s.authClient = grpcClient.NewAuthClient(authv1.NewAuthServiceClient(s.GRPCClientConn()))
	}

	return s.authClient
}

func (s *serviceProvider) AccessClient(_ context.Context) AccessClient {

	if s.accessClient == nil {
		s.accessClient = grpcClient.NewAccessClient(accessv1.NewAccessServiceClient(s.GRPCClientConn()))
	}

	return s.accessClient
}

func (s *serviceProvider) UserClient(_ context.Context) UserClient {

	if s.userClient == nil {
		s.userClient = grpcClient.NewUserClient(userv1.NewUserServiceClient(s.GRPCClientConn()))
	}

	return s.userClient
}

func (s *serviceProvider) ChatImplementation(ctx context.Context) *delivery.ChatImplementation {
	if s.chatImplementation == nil {
		s.chatImplementation = delivery.NewGrpcChatImplementation(s.ChatService(ctx), s.UserClient(ctx))
	}

	return s.chatImplementation
}
