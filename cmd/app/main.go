package main

import (
	"fmt"
	"log"
	"net"

	"github.com/fatih/color"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/config"
	accessgrpc "github.com/javascriptizer1/grpc-cli-chat.backend/internal/delivery/grpc/access"
	authgrpc "github.com/javascriptizer1/grpc-cli-chat.backend/internal/delivery/grpc/auth"
	usergrpc "github.com/javascriptizer1/grpc-cli-chat.backend/internal/delivery/grpc/user"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/interceptor"
	userrepo "github.com/javascriptizer1/grpc-cli-chat.backend/internal/repository/user"
	authsvc "github.com/javascriptizer1/grpc-cli-chat.backend/internal/service/auth"
	usersvc "github.com/javascriptizer1/grpc-cli-chat.backend/internal/service/user"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/client/postgres"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/access_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/auth_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.MustLoad()

	db, err := postgres.New(postgres.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
	})

	if err != nil {
		log.Fatalf("failed to start up db: %v", err)
	}

	defer db.Close()

	fmt.Println(color.GreenString("Server is running on port: %d...", cfg.GRPC.Port))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.ValidateInterceptor))

	reflection.Register(s)

	userRepo := userrepo.New(db)
	userService := usersvc.New(userRepo)
	authService := authsvc.New(userRepo, authsvc.AuthConfig{
		AccessTokenSecretKey:  cfg.JWT.AccessSecretKey,
		AccessTokenDuration:   cfg.JWT.AccessDuration,
		RefreshTokenSecretKey: cfg.JWT.RefreshSecretKey,
		RefreshTokenDuration:  cfg.JWT.RefreshDuration,
	})

	authGRPCServer := authgrpc.New(authService)
	accessGRPCServer := accessgrpc.New(authService)
	userGRPCServer := usergrpc.New(userService)

	auth_v1.RegisterAuthServiceServer(s, authGRPCServer)
	access_v1.RegisterAccessServiceServer(s, accessGRPCServer)
	user_v1.RegisterUserServiceServer(s, userGRPCServer)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
