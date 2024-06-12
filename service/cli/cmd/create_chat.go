package cmd

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	colog "github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/util"
	"github.com/spf13/cobra"
)

func newCreateChatCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-chat",
		Short: "Create a new chat",
		Run: func(cmd *cobra.Command, _ []string) {
			emails, err := cmd.Flags().GetStringArray("emails")

			if err != nil {
				colog.Fatal("failed to get emails: %s", err.Error())
			}

			if len(emails) == 0 {
				colog.Fatal("emails can not be empty: %v", err)
			}

			id, err := sp.HandlerService(ctx).CreateChat(ctx, emails)

			if err != nil {
				colog.Error("could not create chat: %v", err)
			} else {
				colog.Success("chat created with ID: %s", id)
			}
		},
	}

	addCreateChatFlags(cmd)

	return cmd
}

func addCreateChatFlags(cmd *cobra.Command) {
	cmd.Flags().StringArray("emails", []string{}, "Emails of the users in the chat")
	cmd.MarkFlagRequired("emails")
}
