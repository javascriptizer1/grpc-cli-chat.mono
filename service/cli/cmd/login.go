// cmd/login.go
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	"github.com/spf13/cobra"
)

func newLoginCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login a user",
		Run: func(cmd *cobra.Command, args []string) {
			email, _ := cmd.Flags().GetString("email")
			password, _ := cmd.Flags().GetString("password")

			authClient := sp.AuthClient(ctx)

			refreshToken, err := authClient.Login(context.Background(), email, password)

			if err != nil {
				log.Fatalf("Could not login: %v", err)
			}

			fmt.Printf("Logged in, refresh token: %s\n", refreshToken)
		},
	}

	cmd.Flags().String("email", "", "Email of the user")
	cmd.Flags().String("password", "", "Password of the user")

	return cmd
}
