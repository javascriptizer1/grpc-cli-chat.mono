package cmd

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	colog "github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/util"
	"github.com/spf13/cobra"
)

func newLoginCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login a user",
		Run: func(cmd *cobra.Command, _ []string) {
			login, err := cmd.Flags().GetString("login")

			if err != nil {
				colog.Fatal("failed to get login: %s", err.Error())
			}

			password, err := cmd.Flags().GetString("password")

			if err != nil {
				colog.Fatal("failed to get password: %s", err.Error())
			}

			refreshToken, err := sp.HandlerService(ctx).Login(ctx, login, password)

			if err != nil {
				colog.Error("could not login: %v", err)
			} else {
				colog.Success("logged in, refresh token: %s", refreshToken)
			}
		},
	}

	addLoginFlags(cmd)

	return cmd
}

func addLoginFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("login", "l", "", "Email of the user")
	cmd.Flags().StringP("password", "p", "", "Password of the user")
	cmd.MarkFlagRequired("login")
	cmd.MarkFlagRequired("password")
}
