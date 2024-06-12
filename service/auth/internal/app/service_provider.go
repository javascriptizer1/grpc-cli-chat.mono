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
	userRepository UserRepository

	authService AuthService
	userService UserService

	authImpl   *delivery.AuthImplementation
	accessImpl *delivery.AccessImplementation
	userImpl   *delivery.UserImplementation
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

func (s *serviceProvider) PostgresClient(_ context.Context) *sqlx.DB {
	if s.dbClient == nil {

		db, err := postgres.New(postgres.Config{
			Host:     s.Config().DB.Host,
			Port:     s.Config().DB.Port,
			User:     s.Config().DB.User,
			Password: s.Config().DB.Password,
			Name:     s.Config().DB.Name,
		})

		if err != nil {
			logger.Fatal("failed to start up db: ", zap.String("err", err.Error()))
		}

		closer.Add(db.Close)

		s.dbClient = db
	}

	return s.dbClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) UserRepository {

	if s.userRepository == nil {
		s.userRepository = repository.NewUserRepository(s.PostgresClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) AuthService {

	if s.authService == nil {
		s.authService = service.NewAuthService(
			s.UserRepository(ctx),
			service.AuthConfig{
				AccessTokenSecretKey:  s.Config().JWT.AccessSecretKey,
				AccessTokenDuration:   s.Config().JWT.AccessDuration,
				RefreshTokenSecretKey: s.Config().JWT.RefreshSecretKey,
				RefreshTokenDuration:  s.Config().JWT.RefreshDuration,
			})
	}

	return s.authService
}

func (s *serviceProvider) UserService(ctx context.Context) UserService {

	if s.userService == nil {
		s.userService = service.NewUserService(
			s.UserRepository(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) GRPCAuthImpl(ctx context.Context) *delivery.AuthImplementation {

	if s.authImpl == nil {
		s.authImpl = delivery.NewGrpcAuthImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}

func (s *serviceProvider) GRPCAccessImpl(ctx context.Context) *delivery.AccessImplementation {

	if s.accessImpl == nil {
		s.accessImpl = delivery.NewGrpcAccessImplementation(s.AuthService(ctx))
	}

	return s.accessImpl
}

func (s *serviceProvider) GRPCUserImpl(ctx context.Context) *delivery.UserImplementation {

	if s.userImpl == nil {
		s.userImpl = delivery.NewGrpcUserImplementation(s.UserService(ctx))
	}

	return s.userImpl
}
