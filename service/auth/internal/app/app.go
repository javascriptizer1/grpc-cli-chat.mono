package app

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	accessv1 "github.com/javascriptizer1/grpc-cli-chat.mono/pkg/grpc/access_v1"
	authv1 "github.com/javascriptizer1/grpc-cli-chat.mono/pkg/grpc/auth_v1"
	userv1 "github.com/javascriptizer1/grpc-cli-chat.mono/pkg/grpc/user_v1"
	"github.com/javascriptizer1/grpc-cli-chat.mono/pkg/helper/closer"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/interceptor"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/logger"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			logger.Fatal("failed to run GRPC server: ", zap.Any("error", err))
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			logger.Fatal("failed to run HTTP server: ", zap.Any("error", err))
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initLogger,
		a.initGRPCServer,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)

		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.Init(a.serviceProvider.Config().Env)

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LogInterceptor,
			interceptor.NewAuthInterceptor(a.serviceProvider.Config().JWT.AccessSecretKey).Unary,
			interceptor.ValidateInterceptor,
		),
	)

	reflection.Register(a.grpcServer)

	authv1.RegisterAuthServiceServer(a.grpcServer, a.serviceProvider.GRPCAuthImpl(ctx))
	accessv1.RegisterAccessServiceServer(a.grpcServer, a.serviceProvider.GRPCAccessImpl(ctx))
	userv1.RegisterUserServiceServer(a.grpcServer, a.serviceProvider.GRPCUserImpl(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, a.serviceProvider.Config().GRPC.Host, opts)

	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.Config().HTTP.HostPort(),
		Handler:           corsMiddleware.Handler(mux),
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      a.serviceProvider.Config().HTTP.Timeout,
	}

	return nil
}

func (a *App) runGRPCServer() error {
	logger.Info("GRPC server is running on " + a.serviceProvider.Config().GRPC.HostPort())

	l, err := net.Listen("tcp", a.serviceProvider.Config().GRPC.HostPort())

	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(l)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	logger.Info("HTTP server is running on " + a.serviceProvider.Config().HTTP.HostPort())

	err := a.httpServer.ListenAndServe()

	if err != nil {
		return err
	}

	return nil
}
