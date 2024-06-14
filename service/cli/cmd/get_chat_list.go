package cmd

import (
	"context"
	"math"

	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	colog "github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/util"
	"github.com/spf13/cobra"
)

func newGetChatListCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-chat",
		Short: "Get list of user chats",
		Run: func(cmd *cobra.Command, _ []string) {
			limit, err := cmd.Flags().GetUint32("limit")

			if err != nil {
				colog.Fatal("failed to get limit: %s", err.Error())
			}

			page, err := cmd.Flags().GetUint32("page")

			if err != nil {
				colog.Fatal("failed to get page: %s", err.Error())
			}

			chats, total, err := sp.HandlerService(ctx).GetChatList(ctx, pagination.New(limit, page))

			if err != nil {
				colog.Fatal("could not get chats: %v", err)
			}

			colog.Info("")
			for i, v := range chats {
				colog.Info("%d. %s (%s)", i+1, v.Name, v.ID)
			}

			colog.Info("")
			colog.Info("Total: %d", total)
			colog.Info("Page: %d/%d", page, int(math.Ceil(float64(total)/float64(limit))))
		},
	}

	addGetChatListFlags(cmd)

	return cmd
}

func addGetChatListFlags(cmd *cobra.Command) {
	cmd.Flags().Uint32("limit", 10, "Display limit")
	cmd.Flags().Uint32("page", 1, "Display page")
}
