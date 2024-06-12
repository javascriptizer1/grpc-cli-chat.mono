// cmd/create_chat.go
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	"github.com/spf13/cobra"
)

func newCreateChatCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-chat",
		Short: "Create a new chat",
		Run: func(cmd *cobra.Command, args []string) {
			emails, _ := cmd.Flags().GetStringArray("emails")

			id, err := sp.ChatClient(ctx).CreateChat(context.Background(), emails)

			if err != nil {
				log.Fatalf("Could not create chat: %v", err)
			}

			fmt.Printf("Chat created with ID: %s\n", id)
		},
	}

	cmd.Flags().StringArray("emails", []string{}, "Emails of the users in the chat")

	return cmd
}
