package cmd

import (
	"context"
	"log"

	"github.com/fatih/color"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	"github.com/spf13/cobra"
)

func newCreateChatCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-chat",
		Short: "Create a new chat",
		Run: func(cmd *cobra.Command, _ []string) {
			emails, err := cmd.Flags().GetStringArray("emails")

			if err != nil {
				log.Print(color.RedString("failed to get emails: %s\n", err.Error()))
			}

			id, err := sp.ChatClient(ctx).CreateChat(context.Background(), emails)

			if err != nil {
				log.Print(color.RedString("Could not create chat: %v", err))
			} else {
				log.Print(color.GreenString("Chat created with ID: %s\n", id))
			}
		},
	}

	cmd.Flags().StringArray("emails", []string{}, "Emails of the users in the chat")

	return cmd
}
