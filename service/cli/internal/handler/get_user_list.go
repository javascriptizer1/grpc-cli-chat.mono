package handler

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/domain"
)

func (h *Handler) GetUserList(ctx context.Context, options *domain.UserListOption) ([]*domain.UserInfo, uint32, error) {
	users, total, err := h.userClient.GetUserList(ctx, options)

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
