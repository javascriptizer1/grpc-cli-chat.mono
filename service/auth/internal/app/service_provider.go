package app

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/client/postgres"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/helper/closer"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/config"
	delivery "github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/delivery/grpc"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/logger"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/repository"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/service"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type serviceProvider struct {
	config *config.Config

	dbClient       *sqlx.DB
	userRepository *repository.UserRepository

	authService *service.AuthService
	userService *service.UserService

	authImpl   *delivery.AuthImplementation
	accessImpl *delivery.AccessImplementation
	userImpl   *delivery.UserImplementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetConfig() *config.Config {

	if s.config == nil {
		cfg := config.MustLoad()
		s.config = cfg
	}

	return s.config
}

func (s *serviceProvider) GetPostgresClient(_ context.Context) *sqlx.DB {
	if s.dbClient == nil {

		db, err := postgres.New(postgres.Config{
			Host:     s.GetConfig().DB.Host,
			Port:     s.GetConfig().DB.Port,
			User:     s.GetConfig().DB.User,
			Password: s.GetConfig().DB.Password,
			Name:     s.GetConfig().DB.Name,
		})

		if err != nil {
			logger.Fatal("failed to start up db: ", zap.String("err", err.Error()))
		}

		closer.Add(db.Close)

		s.dbClient = db
	}

	return s.dbClient
}

func (s *serviceProvider) GetUserRepository(ctx context.Context) *repository.UserRepository {

	if s.userRepository == nil {
		s.userRepository = repository.NewUserRepository(s.GetPostgresClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) GetAuthService(ctx context.Context) *service.AuthService {

	if s.authService == nil {
		s.authService = service.NewAuthService(
			s.GetUserRepository(ctx),
			service.AuthConfig{
				AccessTokenSecretKey:  s.GetConfig().JWT.AccessSecretKey,
				AccessTokenDuration:   s.GetConfig().JWT.AccessDuration,
				RefreshTokenSecretKey: s.GetConfig().JWT.RefreshSecretKey,
				RefreshTokenDuration:  s.GetConfig().JWT.RefreshDuration,
			})
	}

	return s.authService
}

func (s *serviceProvider) GetUserService(ctx context.Context) *service.UserService {

	if s.userService == nil {
		s.userService = service.NewUserService(
			s.GetUserRepository(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) GetGRPCAuthImpl(ctx context.Context) *delivery.AuthImplementation {

	if s.authImpl == nil {
		s.authImpl = delivery.NewGrpcAuthImplementation(s.GetAuthService(ctx))
	}

	return s.authImpl
}

func (s *serviceProvider) GetGRPCAccessImpl(ctx context.Context) *delivery.AccessImplementation {

	if s.accessImpl == nil {
		s.accessImpl = delivery.NewGrpcAccessImplementation(s.GetAuthService(ctx))
	}

	return s.accessImpl
}

func (s *serviceProvider) GetGRPCUserImpl(ctx context.Context) *delivery.UserImplementation {

	if s.userImpl == nil {
		s.userImpl = delivery.NewGrpcUserImplementation(s.GetUserService(ctx))
	}

	return s.userImpl
}
