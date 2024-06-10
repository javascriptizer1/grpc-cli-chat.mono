package main

import (
	"fmt"
	"log"
	"net"

	"github.com/fatih/color"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/config"
	gs "github.com/javascriptizer1/grpc-cli-chat.backend/internal/delivery/grpc_server"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/interceptor"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/repository"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/service"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/client/postgres"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/authv1"
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

	userRepo := repository.New(db)
	userService := service.New(userRepo)
	grpcServer := gs.New(userService)

	authv1.RegisterUserServiceServer(s, grpcServer)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
