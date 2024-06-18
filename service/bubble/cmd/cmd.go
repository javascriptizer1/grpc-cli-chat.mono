package cmd

import (
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/bubble/internal/app"
)

func Launch() {
	ctx := context.Background()
	sp := app.NewServiceProvider()

	err := sp.TokenManager(ctx).Load()

	if err != nil {
		log.Fatalf("Failed to load token: %v", err)
	}

	p := tea.NewProgram(
		initialAuthModel(ctx, sp),
		tea.WithAltScreen(),
		tea.WithContext(ctx),
	)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}

}
