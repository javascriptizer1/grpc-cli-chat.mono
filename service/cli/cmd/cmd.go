package cmd

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	colog "github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gchat",
	Short:   "CLI client for chat service",
	Version: "0.0.1",
}

func Execute() {
	ctx := context.Background()

	sp := app.NewServiceProvider()

	rootCmd.AddCommand(
		newRegisterCommand(ctx, sp),
		newLoginCommand(ctx, sp),
		newCreateChatCommand(ctx, sp),
		newConnectChatCommand(ctx, sp),
		newSendMessageCommand(ctx, sp),
		// newListChatsCommand(ctx, sp),
		// newListUsersCommand(ctx, sp),
	)

	if err := rootCmd.Execute(); err != nil {
		colog.Fatal("failed to execute root command: %s", err.Error())
	}
}
