package app

import (
	"context"
	"log"

	mongoClient "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/client/mongo"
	accessv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/access_v1"
	authv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/auth_v1"
	userv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/user_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/helper/closer"
	grpcClient "github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/client/grpc"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/config"
	delivery "github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/delivery/grpc"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/repository"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serviceProvider struct {
	config *config.Config

	grpcClientConn *grpc.ClientConn

	dbClient *mongo.Database

	chatRepo    *repository.ChatRepository
	messageRepo *repository.MessageRepository

	chatService *service.ChatService

	authClient   *grpcClient.AuthClient
	accessClient *grpcClient.AccessClient
	userClient   *grpcClient.UserClient

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

func (s *serviceProvider) GRPCClientConn() *grpc.ClientConn {

	if s.grpcClientConn == nil {
		conn, err := grpc.NewClient(
			s.Config().GRPCClient.HostPort(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err != nil {
			log.Fatalf("failed to connect %s: %s", s.Config().GRPCClient.HostPort(), err.Error())
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
			log.Fatalf("failed to start up db: %v", err)
		}

		s.dbClient = db
	}

	return s.dbClient
}

func (s *serviceProvider) ChatRepository(ctx context.Context) *repository.ChatRepository {

	if s.chatRepo == nil {
		s.chatRepo = repository.NewChatRepository(s.MongoClient(ctx))
	}

	return s.chatRepo
}

func (s *serviceProvider) MessageRepository(ctx context.Context) *repository.MessageRepository {

	if s.messageRepo == nil {
		s.messageRepo = repository.NewMessageRepository(s.MongoClient(ctx))
	}

	return s.messageRepo
}

func (s *serviceProvider) ChatService(ctx context.Context) *service.ChatService {

	if s.chatService == nil {

		s.chatService = service.NewChatService(
			s.ChatRepository(ctx),
			s.MessageRepository(ctx),
			s.UserClient(ctx),
		)
	}

	return s.chatService
}

func (s *serviceProvider) AuthClient(_ context.Context) *grpcClient.AuthClient {

	if s.authClient == nil {
		s.authClient = grpcClient.NewAuthClient(authv1.NewAuthServiceClient(s.GRPCClientConn()))
	}

	return s.authClient
}

func (s *serviceProvider) AccessClient(_ context.Context) *grpcClient.AccessClient {

	if s.accessClient == nil {
		s.accessClient = grpcClient.NewAccessClient(accessv1.NewAccessServiceClient(s.GRPCClientConn()))
	}

	return s.accessClient
}

func (s *serviceProvider) UserClient(_ context.Context) *grpcClient.UserClient {

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
