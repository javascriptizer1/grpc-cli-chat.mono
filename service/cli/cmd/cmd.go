package cmd

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.mono/service/cli/internal/app"
)

func Execute() error {
	ctx := context.Background()
	sp := app.NewServiceProvider()

	rootCmd := newRootCommand(ctx, sp)

	return rootCmd.ExecuteContext(ctx)
}
