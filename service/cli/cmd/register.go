package cmd

import (
	"context"
	"log"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/client/grpc/dto"
	"github.com/spf13/cobra"
)

func newRegisterCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new user",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			email, _ := cmd.Flags().GetString("email")
			password, _ := cmd.Flags().GetString("password")
			passwordConfirm, _ := cmd.Flags().GetString("password-confirm")
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
				log.Fatalf("Could not register: %v", err)
			}

			log.Printf("Registered user with ID: %s\n", id)
		},
	}

	cmd.Flags().String("name", "", "Name of the user")
	cmd.Flags().String("email", "", "Email of the user")
	cmd.Flags().String("password", "", "Password of the user")
	cmd.Flags().String("password-confirm", "", "Password confirmation")

	return cmd
}
