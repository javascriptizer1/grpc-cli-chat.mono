package cmd

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	colog "github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/util"
	"github.com/spf13/cobra"
)

func newConnectChatCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect-chat",
		Short: "Connect to a chat",
		Run: func(cmd *cobra.Command, _ []string) {
			chatID, err := cmd.Flags().GetString("chat-id")

			if err != nil {
				colog.Fatal("failed to get chat id: %s", err.Error())
			}

			sp.HandlerService(ctx).ConnectChat(ctx, chatID)
		},
	}

	addConnectChatFlags(cmd)

	return cmd
}

func addConnectChatFlags(cmd *cobra.Command) {
	cmd.Flags().String("chat-id", "", "ID of the chat")
	cmd.MarkFlagRequired("chat-id")
}
