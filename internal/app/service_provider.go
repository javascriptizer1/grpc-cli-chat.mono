package app

import (
	"context"
	"log"

	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/config"
	accessgrpc "github.com/javascriptizer1/grpc-cli-chat.backend/internal/delivery/grpc/access"
	authgrpc "github.com/javascriptizer1/grpc-cli-chat.backend/internal/delivery/grpc/auth"
	usergrpc "github.com/javascriptizer1/grpc-cli-chat.backend/internal/delivery/grpc/user"
	userrepo "github.com/javascriptizer1/grpc-cli-chat.backend/internal/repository/user"
	authsvc "github.com/javascriptizer1/grpc-cli-chat.backend/internal/service/auth"
	usersvc "github.com/javascriptizer1/grpc-cli-chat.backend/internal/service/user"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/client/postgres"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/helper/closer"
	"github.com/jmoiron/sqlx"
)

type serviceProvider struct {
	config *config.Config

	dbClient       *sqlx.DB
	userRepository *userrepo.UserRepository

	authService *authsvc.AuthService
	userService *usersvc.UserService

	authImpl   *authgrpc.GrpcAuthImplementation
	accessImpl *accessgrpc.GrpcAccessImplementation
	userImpl   *usergrpc.GrpcUserImplementation
}

func new() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetConfig() *config.Config {

	if s.config == nil {
		cfg := config.MustLoad()
		s.config = cfg
	}

	return s.config
}

func (s *serviceProvider) GetPostgresClient(ctx context.Context) *sqlx.DB {
	if s.dbClient == nil {

		db, err := postgres.New(postgres.Config{
			Host:     s.GetConfig().DB.Host,
			Port:     s.GetConfig().DB.Port,
			User:     s.GetConfig().DB.User,
			Password: s.GetConfig().DB.Password,
			Name:     s.GetConfig().DB.Name,
		})

		if err != nil {
			log.Fatalf("failed to start up db: %v", err)
		}

		closer.Add(db.Close)

		s.dbClient = db
	}

	return s.dbClient
}

func (s *serviceProvider) GetUserRepository(ctx context.Context) *userrepo.UserRepository {

	if s.userRepository == nil {
		s.userRepository = userrepo.New(s.GetPostgresClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) GetAuthService(ctx context.Context) *authsvc.AuthService {

	if s.authService == nil {
		s.authService = authsvc.New(
			s.GetUserRepository(ctx),
			authsvc.AuthConfig{
				AccessTokenSecretKey:  s.GetConfig().JWT.AccessSecretKey,
				AccessTokenDuration:   s.GetConfig().JWT.AccessDuration,
				RefreshTokenSecretKey: s.GetConfig().JWT.RefreshSecretKey,
				RefreshTokenDuration:  s.GetConfig().JWT.RefreshDuration,
			})
	}

	return s.authService
}

func (s *serviceProvider) GetUserService(ctx context.Context) *usersvc.UserService {

	if s.userService == nil {
		s.userService = usersvc.New(
			s.GetUserRepository(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) GetGRPCAuthImpl(ctx context.Context) *authgrpc.GrpcAuthImplementation {

	if s.authImpl == nil {
		s.authImpl = authgrpc.New(s.GetAuthService(ctx))
	}

	return s.authImpl
}

func (s *serviceProvider) GetGRPCAccessImpl(ctx context.Context) *accessgrpc.GrpcAccessImplementation {

	if s.accessImpl == nil {
		s.accessImpl = accessgrpc.New(s.GetAuthService(ctx))
	}

	return s.accessImpl
}

func (s *serviceProvider) GetGRPCUserImpl(ctx context.Context) *usergrpc.GrpcUserImplementation {

	if s.userImpl == nil {
		s.userImpl = usergrpc.New(s.GetUserService(ctx))
	}

	return s.userImpl
}
