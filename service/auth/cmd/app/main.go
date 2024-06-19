package main

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/app"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/logger"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	a, err := app.New(ctx)

	if err != nil {
		logger.Fatal("failed to init app ", zap.String("err", err.Error()))
	}

	if err := a.Run(); err != nil {
		logger.Fatal("failed to run app: ", zap.String("err", err.Error()))
	}

}
