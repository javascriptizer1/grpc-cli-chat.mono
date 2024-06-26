package cmd

import (
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/cli/internal/app"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/cli/internal/tui"
	"github.com/spf13/cobra"
)

func newRootCommand(ctx context.Context, sp *app.ServiceProvider) *cobra.Command {
	return &cobra.Command{
		Use:     "gchat",
		Short:   "CLI client for chat service",
		Version: sp.Config().Version,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			err := sp.TokenManager(ctx).Load()

			if err != nil {
				log.Fatalf("load config error: %v", err)
			}

		},
		RunE: func(_ *cobra.Command, _ []string) error {
			var m tea.Model

			if sp.HandlerService(ctx).IsAccessValid(ctx) {
				m = tui.InitialChatListModel(ctx, sp, 0, 0)
			} else {
				m = tui.InitialAuthModel(ctx, sp)
			}

			p := tea.NewProgram(
				m,
				tea.WithAltScreen(),
				tea.WithContext(ctx),
			)

			_, err := p.Run()

			return err
		},
	}
}
