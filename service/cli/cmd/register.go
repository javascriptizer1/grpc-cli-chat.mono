package cmd

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/client/grpc/dto"
	colog "github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/util"
	"github.com/spf13/cobra"
)

func newRegisterCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new user",
		Run: func(cmd *cobra.Command, _ []string) {
			name, err := cmd.Flags().GetString("name")

			if err != nil {
				colog.Fatal("failed to get name: %s", err.Error())
			}

			email, err := cmd.Flags().GetString("email")

			if err != nil {
				colog.Fatal("failed to get email: %s", err.Error())
			}

			password, err := cmd.Flags().GetString("password")

			if err != nil {
				colog.Fatal("failed to get password: %s", err.Error())
			}

			passwordConfirm, err := cmd.Flags().GetString("password-confirm")

			if err != nil {
				colog.Fatal("failed to get password confirmation: %s", err.Error())
			}

			user := 1

			req := dto.RegisterInputDto{
				Name:            name,
				Email:           email,
				Password:        password,
				PasswordConfirm: passwordConfirm,
				Role:            uint16(user),
			}

			id, err := sp.HandlerService(ctx).Register(ctx, req)

			if err != nil {
				colog.Error("could not register: %v", err)
			} else {
				colog.Success("registered user with ID: %s", id)
			}
		},
	}

	addRegisterFlags(cmd)

	return cmd
}

func addRegisterFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("name", "n", "", "Name of the user")
	cmd.Flags().StringP("email", "e", "", "Email of the user")
	cmd.Flags().StringP("password", "p", "", "Password of the user")
	cmd.Flags().StringP("password-confirm", "c", "", "Password confirmation")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("password")
	_ = cmd.MarkFlagRequired("password-confirm")
}
