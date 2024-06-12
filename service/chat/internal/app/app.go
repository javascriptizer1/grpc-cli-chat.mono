package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	chatv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/chat_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/helper/closer"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, fmt.Errorf("init deps: %w", err)
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := a.runGRPCServer(); err != nil {
			log.Fatalf("run GRPC server: %s", err.Error())
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return fmt.Errorf("init: %w", err)
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.NewAuthInterceptor(a.serviceProvider.AccessClient(ctx)).Unary(),
			),
		),
	)

	reflection.Register(a.grpcServer)

	chatv1.RegisterChatServiceServer(a.grpcServer, a.serviceProvider.ChatImplementation(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s\n", a.serviceProvider.Config().GRPC.HostPort())

	l, err := net.Listen("tcp", a.serviceProvider.Config().GRPC.HostPort())

	if err != nil {
		return fmt.Errorf("failed to get listener: %s", err.Error())
	}

	if err = a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("failed to serve: %s", err.Error())
	}

	return nil
}
