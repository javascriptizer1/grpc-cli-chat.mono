package cmd

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	colog "github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/util"
	"github.com/spf13/cobra"
)

func newGetChatCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-chat",
		Short: "Get chat info by id",
		Run: func(cmd *cobra.Command, _ []string) {
			chatID, err := cmd.Flags().GetString("chat-id")

			if err != nil {
				colog.Fatal("failed to get chat id: %s", err.Error())
			}

			chat, err := sp.HandlerService(ctx).GetChat(ctx, chatID)

			if err != nil {
				colog.Fatal("could not get chat: %v", err)
			}

			colog.Info("ID: %s", chat.ID)
			colog.Info("Name: %s", chat.Name)
			colog.Info("Participants")

			for i, v := range chat.Users {
				colog.Info("  %d. %s (%s)", i+1, v.Name, v.ID)
			}

		},
	}

	addGetChatFlags(cmd)

	return cmd
}

func addGetChatFlags(cmd *cobra.Command) {
	cmd.Flags().String("chat-id", "", "ID of the chat")
	_ = cmd.MarkFlagRequired("chat-id")
}
