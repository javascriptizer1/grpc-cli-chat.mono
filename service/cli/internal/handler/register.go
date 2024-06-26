package handler

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.mono/service/cli/internal/client/grpc/dto"
)

func (h *Handler) Register(ctx context.Context, in dto.RegisterInputDto) (string, error) {
	return h.authClient.Register(ctx, in)
}
