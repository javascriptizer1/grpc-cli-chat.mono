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
			name, err := cmd.Flags().GetString("name")

			if err != nil {
				colog.Fatal("failed to get name: %s", err.Error())
			}

			userIDs, err := cmd.Flags().GetStringArray("user-ids")

			if err != nil {
				colog.Fatal("failed to get emails: %s", err.Error())
			}

			if len(userIDs) == 0 {
				colog.Fatal("id of users can not be empty: %v", err)
			}

			id, err := sp.HandlerService(ctx).CreateChat(ctx, name, userIDs)

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
	cmd.Flags().String("name", "", "Name of the chat")
	_ = cmd.Flags().StringArray("user-ids", []string{}, "ID of the users")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("user-ids")

}
