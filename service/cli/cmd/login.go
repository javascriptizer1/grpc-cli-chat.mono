package cmd

import (
	"context"
	"log"

	"github.com/fatih/color"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	"github.com/spf13/cobra"
)

func newLoginCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login a user",
		Run: func(cmd *cobra.Command, _ []string) {
			login, err := cmd.Flags().GetString("login")

			if err != nil {
				log.Print(color.RedString("failed to get login: %s\n", err.Error()))
			}

			password, err := cmd.Flags().GetString("password")

			if err != nil {
				log.Print(color.RedString("failed to get password: %s\n", err.Error()))
			}

			authClient := sp.AuthClient(ctx)

			refreshToken, err := authClient.Login(context.Background(), login, password)

			if err != nil {
				log.Print(color.RedString("Could not login: %v", err))
			} else {
				log.Print(color.GreenString("Logged in, refresh token: %s\n", refreshToken))
			}
		},
	}

	cmd.Flags().String("email", "", "Email of the user")
	cmd.Flags().String("password", "", "Password of the user")

	return cmd
}
