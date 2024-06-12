package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	"github.com/spf13/cobra"
)

func newConnectChatCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect-chat",
		Short: "Connect to a chat",
		Run: func(cmd *cobra.Command, _ []string) {
			chatID, err := cmd.Flags().GetString("chat_id")

			if err != nil {
				log.Print(color.RedString("failed to get chat id: %s\n", err.Error()))
			}

			stream, err := sp.ChatClient(ctx).ConnectChat(context.Background(), chatID)

			if err != nil {
				log.Print(color.RedString("Could not connect to chat: %v", err))
			}

			for {
				msg, err := stream.Recv()

				if err != nil {
					log.Print(color.RedString("Error receiving message: %v", err))
				} else {
					fmt.Print(color.WhiteString("%s: %s\n", msg.Sender.Name, msg.Text))
				}
			}
		},
	}

	cmd.Flags().String("chat_id", "", "ID of the chat")

	return cmd
}
