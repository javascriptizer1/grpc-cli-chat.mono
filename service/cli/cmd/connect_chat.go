package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	"github.com/spf13/cobra"
)

func newConnectChatCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect-chat",
		Short: "Connect to a chat",
		Run: func(cmd *cobra.Command, args []string) {
			chatID, _ := cmd.Flags().GetString("chat_id")

			stream, err := sp.ChatClient(ctx).ConnectChat(context.Background(), chatID)

			if err != nil {
				log.Fatalf("Could not connect to chat: %v", err)
			}

			for {
				msg, err := stream.Recv()

				if err != nil {
					log.Fatalf("Error receiving message: %v", err)
				}

				fmt.Printf("Message from %s: %s\n", msg.Sender.Name, msg.Text)
			}
		},
	}

	cmd.Flags().String("chat_id", "", "ID of the chat")

	return cmd
}
