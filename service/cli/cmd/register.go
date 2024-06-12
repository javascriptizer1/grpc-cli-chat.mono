package cmd

import (
	"context"
	"log"

	"github.com/fatih/color"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/client/grpc/dto"
	"github.com/spf13/cobra"
)

func newRegisterCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new user",
		Run: func(cmd *cobra.Command, _ []string) {
			name, err := cmd.Flags().GetString("name")

			if err != nil {
				log.Print(color.RedString("failed to get name: %s\n", err.Error()))
			}

			email, err := cmd.Flags().GetString("email")

			if err != nil {
				log.Print(color.RedString("failed to get email: %s\n", err.Error()))
			}

			password, err := cmd.Flags().GetString("password")

			if err != nil {
				log.Print(color.RedString("failed to get password: %s\n", err.Error()))
			}

			passwordConfirm, err := cmd.Flags().GetString("password-confirm")

			if err != nil {
				log.Print(color.RedString("failed to get password confirmation: %s\n", err.Error()))
			}

			user := 1

			authClient := sp.AuthClient(ctx)

			req := dto.RegisterInputDto{
				Name:            name,
				Email:           email,
				Password:        password,
				PasswordConfirm: passwordConfirm,
				Role:            uint16(user),
			}

			id, err := authClient.Register(context.Background(), req)

			if err != nil {
				log.Print(color.RedString("Could not register: %v", err))
			} else {
				log.Print(color.GreenString("Registered user with ID: %s\n", id))
			}
		},
	}

	cmd.Flags().String("name", "", "Name of the user")
	cmd.Flags().String("email", "", "Email of the user")
	cmd.Flags().String("password", "", "Password of the user")
	cmd.Flags().String("password-confirm", "", "Password confirmation")

	return cmd
}
