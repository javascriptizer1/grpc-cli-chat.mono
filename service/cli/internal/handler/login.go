package handler

import "context"

func (h *Handler) Login(ctx context.Context, login string, password string) (string, error) {
	return h.authClient.Login(ctx, login, password)
}
