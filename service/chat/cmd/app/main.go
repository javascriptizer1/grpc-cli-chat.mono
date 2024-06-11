package main

import (
	"context"
	"log"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/app"
)

func main() {

	ctx := context.Background()

	a, err := app.New(ctx)

	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	if err := a.Run(); err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
