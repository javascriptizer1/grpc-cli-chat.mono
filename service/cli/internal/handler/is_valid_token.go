package handler

import (
	"context"
)

func (h *Handler) IsAccessValid(ctx context.Context) bool {
	user, err := h.userClient.GetUserInfo(ctx)

	return user != nil && err == nil
}
