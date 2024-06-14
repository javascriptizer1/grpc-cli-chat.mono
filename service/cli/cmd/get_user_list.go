package cmd

import (
	"context"
	"math"

	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/domain"
	colog "github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/util"
	"github.com/spf13/cobra"
)

func newGetUserListCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-user",
		Short: "Get list of users",
		Run: func(cmd *cobra.Command, _ []string) {
			limit, err := cmd.Flags().GetUint32("limit")

			if err != nil {
				colog.Fatal("failed to get limit: %s", err.Error())
			}

			page, err := cmd.Flags().GetUint32("page")

			if err != nil {
				colog.Fatal("failed to get page: %s", err.Error())
			}

			users, total, err := sp.HandlerService(ctx).GetUserList(
				ctx,
				&domain.UserListOption{
					Pagination: *pagination.New(limit, page),
					UserIDs:    []string{},
				})

			if err != nil {
				colog.Fatal("could not get users: %v", err)
			}

			colog.Info("")

			for i, v := range users {
				colog.Info("%d. %s (%s)", i+1, v.Name, v.ID)
			}

			colog.Info("")
			colog.Info("Total: %d", total)
			colog.Info("Page: %d/%d", page, int(math.Ceil(float64(total)/float64(limit))))

		},
	}

	addGetUserListFlags(cmd)

	return cmd
}

func addGetUserListFlags(cmd *cobra.Command) {
	cmd.Flags().Uint32("limit", 10, "Display limit")
	cmd.Flags().Uint32("page", 1, "Display page")
}
