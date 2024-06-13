package cmd

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	colog "github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/util"
	"github.com/spf13/cobra"
)

func newSendMessageCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-message",
		Short: "Send a message to a chat",
		Run: func(cmd *cobra.Command, _ []string) {
			chatID, err := cmd.Flags().GetString("chat-id")

			if err != nil {
				colog.Fatal("failed to get chat id: %s\n", err.Error())
			}

			text, err := cmd.Flags().GetString("text")

			if err != nil {
				colog.Fatal("failed to get message text: %s\n", err.Error())
			}

			err = sp.HandlerService(ctx).SendMessage(ctx, text, chatID)

			if err != nil {
				colog.Warn("could not send message: %v", err)
			}
		},
	}

	addSendMessageFlags(cmd)

	return cmd
}

func addSendMessageFlags(cmd *cobra.Command) {
	cmd.Flags().String("chat-id", "", "ID of the chat")
	cmd.Flags().String("text", "", "Text of the message")
	_ = cmd.MarkFlagRequired("chat-id")
	_ = cmd.MarkFlagRequired("text")
}
